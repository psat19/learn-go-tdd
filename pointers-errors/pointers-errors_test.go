package main

import (
	"testing"
)

func TestWallet(t *testing.T) {
	t.Run("Deposit", func(t *testing.T) {
		thisWallet := Wallet{}

		thisWallet.Deposit(Bitcoin(10))

		got := thisWallet.Balance()
		want := Bitcoin(10)

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("Withdraw", func(t *testing.T) {
		thisWallet := Wallet{balance: Bitcoin(20)}

		err := thisWallet.Withdraw(Bitcoin(10))
		if err != nil {
			t.Fatal("got an error but didn't want one")
		}

		got := thisWallet.Balance()

		want := Bitcoin(10)

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("withdraw insufficient funds", func(t *testing.T) {
		startingBalance := Bitcoin(20)
		wallet := Wallet{startingBalance}
		err := wallet.Withdraw(Bitcoin(100))

		got := wallet.Balance()
		want := startingBalance

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}

		if err == nil {
			t.Fatal("wanted an error but didn't get one")
		}

		if err.Error() != "cannot withdraw, insufficient funds" {
			t.Errorf("got %q, want %q", err, "cannot withdraw, insufficient funds")
		}
	})

}
