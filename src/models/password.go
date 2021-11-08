package models

type Password struct {
	NewPassoword    string `json:"newPassword"`
	CurrentPassword string `json:"currentPassword"`
}
