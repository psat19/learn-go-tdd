# Critical Review of Test Implementation

Good effort overall — the structure is solid. Here's a critical breakdown:

---

## `MockUserStore` — significant issues

### 1. It's not a true mock — it's a fake.
A mock is configured per-test to return a specific value or error. You used a shared `map[string]User` pre-loaded with data. This means the mock encodes test knowledge (what users exist) rather than letting each test declare exactly what it needs. Compare:

```go
// Fake (yours)
store := &MockUserStore{store: map[string]User{"123": user}}

// Proper mock (controllable per-test)
store := &MockUserStore{userToReturn: user, errToReturn: nil}
```

### 2. Missing `calledWithID` tracking.
You never assert that `GetUser` was called with the right ID. This is one of the most valuable things a mock gives you — verifying the *call pathway*, not just the outcome.

---

## `MockEmailer` — significant issue

### 3. The "email send fails" test is incorrectly implemented.
You simulate failure by having user `"456"` with an empty `Email` field, then the mock checks `strings.TrimSpace(to) == ""`. This is fragile and wrong:
- The mock is now coupled to the data in the store.
- The mock shouldn't contain business logic. Its only job is to return `errToReturn` when told to.
- The right approach: set `errToReturn = errors.New("send failed")` on the mock directly.

### 4. Missing `to` field tracking.
`MockEmailer` never records the `to` address. You asserted `subject` and `body` but never verified the email went to `user.Email`.

---

## `TestSendWelcome` — structural issues

### 5. Shared service/mocks across subtests.
`mailer.called` from the "valid user" test is still `true` when "notification failed" runs. Subtests should each set up fresh mocks to be independent and not leak state.

### 6. The "valid user" test doesn't assert the return value.
`service.SendWelcome(id)` is called but the error is discarded. A successful path should also confirm `err == nil`.

### 7. Two assertions are commented out.
The error message assertions in "invalid user" and "notification failed" are commented out. These are meaningful checks — don't skip them.

---

## Summary

| Issue | Severity |
|---|---|
| MockUserStore is a fake, not a mock | Medium |
| No `calledWithID` assertion | Medium |
| Email failure simulated via data, not `errToReturn` | High |
| Missing `to` address assertion | Low |
| Shared state between subtests | High |
| Error ignored on happy path | Medium |
| Commented-out assertions | Low |
