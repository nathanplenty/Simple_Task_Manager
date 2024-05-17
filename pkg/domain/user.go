package domain

// User represents a user in the system
type User struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// NewUser creates a new User instance
func NewUser(userID, userName, password string) *User {
	return &User{
		UserID:   userID,
		UserName: userName,
		Password: password,
	}
}
