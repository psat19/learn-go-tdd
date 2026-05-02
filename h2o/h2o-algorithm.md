# H2O Synchronization: Barrier + Ordering

## Problem

Hydrogen and oxygen goroutines call `Hydrogen()` and `Oxygen()` concurrently.
Two guarantees are required:

1. **Barrier** — all 3 atoms (2H + 1O) must be present before any of them releases.
2. **Ordering** — within each molecule, release order must be H → H → O.

---

## Why simpler approaches fail

### Ordering only (counter-based)

```go
for b.counter >= 2 { b.cond.Wait() }  // H waits
for b.counter < 2  { b.cond.Wait() }  // O waits
```

Enforces sequencing but breaks the barrier: H1 and H2 release before oxygen has
arrived at all. The atoms never "bond together".

### Generic cyclic barrier (n=3)

```go
barrier.Await()  // trips when any 3 goroutines arrive
```

Does not distinguish types. Three hydrogen goroutines trip the barrier with zero oxygen.

### Generation-based typed barrier

```go
// trips when hCount≥2 AND oCount≥1, then increments generation
// goroutines exit when their captured gen != current generation
```

**The excess-goroutine bug.** With input `HHHHOO`, all four H's arrive and increment
`hCount` before any O. When O1 arrives, `tryTrip` fires once — but all four H's captured
the same old `gen` value, so all four exit the barrier loop. H3 and H4 pass through
without their O2 being present, breaking the barrier guarantee.

---

## The Solution: Slot-based Barrier + Step Ordering

Two independent phases share a single `sync.Cond`.

### Phase 1 — Slot-based barrier

`tryTrip` mints exactly **2 H-tokens** and **1 O-token** per molecule:

```go
func (b *H2O) tryTrip() {
    if b.hCount >= 2 && b.oCount >= 1 {
        b.hCount -= 2
        b.oCount -= 1
        b.hSlots += 2   // exactly 2 passes
        b.oSlots += 1   // exactly 1 pass
        b.cond.Broadcast()
    }
}
```

A goroutine only exits the barrier by claiming a slot:

```go
// Hydrogen — barrier phase
for b.hSlots == 0 { b.cond.Wait() }
b.hSlots--

// Oxygen — barrier phase
for b.oSlots == 0 { b.cond.Wait() }
b.oSlots--
```

Because `tryTrip` decrements `hCount` and `oCount` before granting slots, it can only
fire again when a *new* pair of 2H + 1O has accumulated. H3 and H4 cannot claim slots
until O2 arrives and fires a second `tryTrip`.

### Phase 2 — Step ordering

`step` (0 → 1 → 2 → reset) gates releases within each molecule:

```go
// Hydrogen — ordering phase
for b.step >= 2 { b.cond.Wait() }   // wait while it's O's turn
b.step++
b.cond.Broadcast()

// Oxygen — ordering phase
for b.step < 2 { b.cond.Wait() }    // wait until both H's have gone
b.step = 0                           // reset for next molecule
b.cond.Broadcast()
```

`step` resets only when O releases, so O is structurally last every time.

---

## Data structures

| Variable | Role |
|---|---|
| `hCount` | H atoms waiting at the barrier (not yet slotted) |
| `oCount` | O atoms waiting at the barrier (not yet slotted) |
| `hSlots` | Available H passes (minted by `tryTrip`) |
| `oSlots` | Available O passes (minted by `tryTrip`) |
| `step`   | Position within current molecule: 0 = first H, 1 = second H, 2 = O |

---

## Full algorithm

```go
func (b *H2O) AwaitHydrogen() {
    b.mu.Lock()
    b.hCount++
    b.tryTrip()
    for b.hSlots == 0 {      // Phase 1: barrier
        b.cond.Wait()
    }
    b.hSlots--
    for b.step >= 2 {        // Phase 2: ordering
        b.cond.Wait()
    }
    b.step++
    b.cond.Broadcast()
    b.mu.Unlock()
}

func (b *H2O) AwaitOxygen() {
    b.mu.Lock()
    b.oCount++
    b.tryTrip()
    for b.oSlots == 0 {      // Phase 1: barrier
        b.cond.Wait()
    }
    b.oSlots--
    for b.step < 2 {         // Phase 2: ordering
        b.cond.Wait()
    }
    b.step = 0
    b.cond.Broadcast()
    b.mu.Unlock()
}
```

---

## Worked example: `HHHHOO`

| Event   | hCount | oCount | hSlots | oSlots | step | Notes                        |
|---------|--------|--------|--------|--------|------|------------------------------|
| H1      | 1      | 0      | 0      | 0      | 0    | no trip                      |
| H2      | 2      | 0      | 0      | 0      | 0    | no trip (oCount=0)           |
| H3      | 3      | 0      | 0      | 0      | 0    | no trip                      |
| H4      | 4      | 0      | 0      | 0      | 0    | no trip                      |
| O1      | 2      | 0      | 2      | 1      | 0    | **trip!** H1,H2 get slots    |
| H1 ord  | 2      | 0      | 1      | 1      | 1    | step 0 → 1                   |
| H2 ord  | 2      | 0      | 0      | 1      | 2    | step 1 → 2                   |
| O1 ord  | 2      | 0      | 0      | 0      | 0    | step 2 → 0 → **molecule 1 💧** |
| O2      | 0      | 0      | 2      | 1      | 0    | **trip!** H3,H4 get slots    |
| H3 ord  | 0      | 0      | 1      | 1      | 1    | step 0 → 1                   |
| H4 ord  | 0      | 0      | 0      | 1      | 2    | step 1 → 2                   |
| O2 ord  | 0      | 0      | 0      | 0      | 0    | step 2 → 0 → **molecule 2 💧** |

H3 and H4 are blocked in `hSlots == 0` until O2 arrives — the barrier guarantee holds.

---

## Key invariants

- `hSlots` is always a multiple of 2; `oSlots` always equals `hSlots / 2` (at the moment of minting).
- At most 2 H goroutines and 1 O goroutine are ever in the ordering phase for the same molecule simultaneously.
- `step` only resets to 0 when O releases, making O structurally last in every molecule.
- `Broadcast` (not `Signal`) is used because H and O goroutines share a single condition variable and you cannot predict which type will be the relevant next waiter.
