package main

import (
	"errors"
	"testing"
)

func TestSearch(t *testing.T) {
	dictionary := Dictionary{"test": "this is just a test"}

	t.Run("Has the key", func(t *testing.T) {

		got, _ := dictionary.Search("test")
		want := "this is just a test"

		assertStrings(t, got, want)
	})

	t.Run("Does not have the key", func(t *testing.T) {

		_, err := dictionary.Search("test123")
		want := "Key not found"

		if err == nil {
			t.Errorf("Should have got error because the key is not in the dictionary")
		}

		assertError(t, err, errors.New(want))
		assertStrings(t, err.Error(), want)
	})
}

func TestAdd(t *testing.T) {
	dictionary := Dictionary{}

	t.Run("new add", func(t *testing.T) {
		err := dictionary.Add("test", "this is just a test")

		if err != nil {
			t.Errorf("Should not have got error when adding new key")
		}

		assertError(t, err, nil)

		got, _ := dictionary.Search("test")
		want := "this is just a test"

		assertStrings(t, got, want)
	})

	t.Run("existing key", func(t *testing.T) {
		err := dictionary.Add("test", "this is just a test")

		if err == nil {
			t.Errorf("Should have got error when adding existing key")
		}

		assertError(t, err, errors.New("Key already exists"))
	})
}

func TestUpdate(t *testing.T) {
	t.Run("existing word", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dictionary := Dictionary{word: definition}
		newDefinition := "new definition"

		err := dictionary.Update(word, newDefinition)

		assertError(t, err, nil)
		assertDefinition(t, dictionary, word, newDefinition)
	})

	t.Run("new word", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dictionary := Dictionary{}

		err := dictionary.Update(word, definition)

		assertError(t, err, errors.New("The key does not exist"))
	})
}

func TestDelete(t *testing.T) {
	t.Run("existing word", func(t *testing.T) {
		word := "test"
		dictionary := Dictionary{word: "test definition"}

		err := dictionary.Delete(word)

		assertError(t, err, nil)

		_, err = dictionary.Search(word)

		assertError(t, err, errors.New("Key not found"))
	})

	t.Run("non-existing word", func(t *testing.T) {
		word := "test"
		dictionary := Dictionary{}

		err := dictionary.Delete(word)

		assertError(t, err, errors.New("Key not found"))
	})
}

func assertStrings(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertError(t testing.TB, got error, want error) {
	t.Helper()

	if got == nil && want == nil {
		return
	}

	if got == nil || want == nil {
		t.Errorf("got %v want %v", got, want)
		return
	}

	if got.Error() != want.Error() {
		t.Errorf("got %q want %q", got.Error(), want.Error())
	}
}

func assertDefinition(t testing.TB, dictionary Dictionary, word, definition string) {
	t.Helper()
	got, err := dictionary.Search(word)
	if err != nil {
		t.Fatal("should find added word:", err)
	}

	if got != definition {
		t.Errorf("got %q want %q", got, definition)
	}
}
