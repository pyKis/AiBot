package models

type Role string

const (
	Admin  Role = "Admin"
	Client Role = "Client"
	Team   Role = "Team"
	None   Role = "None"
)

type User struct {
	UserID      int64
	Username    string
	FirstName   string
	LastName    string
	PhoneNumber string
	Role        Role
}
