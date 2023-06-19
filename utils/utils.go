package utils

import "github.com/google/uuid"

// function to generate id that returns uuid
func GenerateId() string {
	return uuid.New().String()
}
