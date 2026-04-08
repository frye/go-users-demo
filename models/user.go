package models

// UserProfile represents user profile data
type UserProfile struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	FullName string `json:"fullName"`
	Emoji    string `json:"emoji"`
}

// CreateUserRequest is the input for creating a user
type CreateUserRequest struct {
	FullName string `json:"fullName" binding:"required"`
	Emoji    string `json:"emoji" binding:"required"`
}

// UpdateUserRequest is the input for updating a user
type UpdateUserRequest struct {
	FullName string `json:"fullName" binding:"required"`
	Emoji    string `json:"emoji" binding:"required"`
}
