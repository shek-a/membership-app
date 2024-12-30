package models

type Member struct {
	ID          int    `json:"id"`
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Email       string `json:"email" binding:"required"`
	DateOfBirth string `json:"dateOfBirth" binding:"required"`
}

type UpdateMember struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	DateOfBirth string `json:"dateOfBirth"`
}
