package notifications

import "fmt"

// UserStore retrieves users from a data store.
type UserStore interface {
	GetUser(id string) (User, error)
}

// Emailer sends email messages.
type Emailer interface {
	Send(to, subject, body string) error
}

// User holds basic user info.
type User struct {
	ID    string
	Name  string
	Email string
}

// NotificationService sends notifications to users.
type NotificationService struct {
	store   UserStore
	emailer Emailer
}

// NewNotificationService creates a NotificationService with the given dependencies.
func NewNotificationService(store UserStore, emailer Emailer) *NotificationService {
	return &NotificationService{store: store, emailer: emailer}
}

// SendWelcome fetches the user by ID and sends them a welcome email.
// Returns an error if the user cannot be found or the email fails to send.
func (n *NotificationService) SendWelcome(userID string) error {
	user, err := n.store.GetUser(userID)
	if err != nil {
		return fmt.Errorf("could not get user: %w", err)
	}

	subject := fmt.Sprintf("Welcome, %s!", user.Name)
	body := fmt.Sprintf("Hi %s, thanks for joining!", user.Name)

	if err := n.emailer.Send(user.Email, subject, body); err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}

	return nil
}
