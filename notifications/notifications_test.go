package notifications

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// TODO: Implement the following mocks and tests.
//
// --- What to Mock ---
//
// 1. MockUserStore (implements UserStore)
//    Fields to track:
//      - calledWithID string        — the ID that was passed to GetUser
//      - userToReturn User          — the user to hand back
//      - errToReturn  error         — set this to simulate a store failure
//    Method: GetUser(id string) (User, error)
//
// 2. MockEmailer (implements Emailer)
//    Fields to track:
//      - to, subject, body string   — arguments passed to Send
//      - called bool                — whether Send was invoked at all
//      - errToReturn  error         — set this to simulate a send failure
//    Method: Send(to, subject, body string) error
//
// --- Tests to Write ---
//
// TestSendWelcome
//   "sends correct email to user"
//     - set up MockUserStore to return a known User
//     - set up MockEmailer (no error)
//     - call service.SendWelcome(userID)
//     - assert no error was returned
//     - assert MockEmailer received the correct to/subject/body
//     - assert MockUserStore was called with the right ID
//
//   "returns error when user not found"
//     - set up MockUserStore to return an error
//     - call service.SendWelcome(userID)
//     - assert an error was returned
//     - assert MockEmailer.Send was NOT called
//
//   "returns error when email send fails"
//     - set up MockUserStore to return a valid user
//     - set up MockEmailer to return an error
//     - call service.SendWelcome(userID)
//     - assert an error was returned

type MockUserStore struct {
	store map[string]User
}

func (s *MockUserStore) GetUser(id string) (User, error) {
	user, ok := s.store[id]

	if !ok {
		return User{}, fmt.Errorf("could not get user: %s", id)
	}

	return user, nil
}

type MockEmailer struct {
	subject string
	body    string
	called  bool
}

func (e *MockEmailer) Send(to, subject, body string) error {
	e.called = true

	if strings.TrimSpace(to) == "" {
		return errors.New("Cannot send email")
	}

	e.body = body
	e.subject = subject

	return nil
}

func TestSendWelcome(t *testing.T) {
	user := &MockUserStore{}
	user.store = map[string]User{
		"123": {
			ID:    "123",
			Name:  "Prathap",
			Email: "test@123",
		},
		"456": {
			ID:   "456",
			Name: "Ram",
		},
	}
	mailer := &MockEmailer{}
	service := NewNotificationService(user, mailer)

	t.Run("valid user", func(t *testing.T) {
		id := "123"
		service.SendWelcome(id)

		userName := user.store[id].Name

		if !mailer.called {
			t.Fatalf("notification wasn't called at all")
		}

		AssertString(t, mailer.subject, fmt.Sprintf("Welcome, %s!", userName))
		AssertString(t, mailer.body, fmt.Sprintf("Hi %s, thanks for joining!", userName))
	})
	t.Run("invalid user", func(t *testing.T) {
		id := "1234"
		err := service.SendWelcome(id)

		if err == nil {
			t.Fatalf("Should have got error")
		}

		// AssertString(t, err.Error(), fmt.Sprintf("could not get user: %s", id))
	})
	t.Run("notification failed", func(t *testing.T) {
		id := "456"
		err := service.SendWelcome(id)

		if err == nil {
			t.Fatalf("Should have got error")
		}

		// AssertString(t, err.Error(), "Cannot send email")
	})
}

func AssertString(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("Got: %s \n Want: %s \n ----", got, want)
	}
}
