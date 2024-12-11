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

type FileData struct {
	FileID       string
	FileType     string
	FileName     string
	FileExtension string
}

// Response структура для ответа от API
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data"` // Тут может быть ваша конкретная структура данных
}