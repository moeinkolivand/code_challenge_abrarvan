package dto

type SendMessageRequestDto struct {
	Message     string `json:"message" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
}
