package auth

type User struct {
	ID       int64
	Username string
	PasswordHash string
	Role     string  // e.g. "admin", "engineer", "viewer"
}
