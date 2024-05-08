package domain

// User struct represents the user model
type User struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// NewUser creates and returns a new User instance
func NewUser(userid, username, password string) *User {
	return &User{
		UserID:   userid,
		UserName: username,
		Password: password,
	}
}
