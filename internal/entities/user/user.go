package user

import "time"

// User represents a user (employee) entity with JSON tags
type User struct {
	ID        string    `json:"id"`         // UUID of the user
	Username  string    `json:"username"`   // Unique username
	FirstName string    `json:"first_name"` // First name of the user
	LastName  string    `json:"last_name"`  // Last name of the user
	CreatedAt time.Time `json:"created_at"` // Timestamp when the user was created
	UpdatedAt time.Time `json:"updated_at"` // Timestamp when the user was last updated
}
