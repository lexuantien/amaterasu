package utils

import "github.com/google/uuid"

// create new uuid
// 	string new uid string
func NewUuidString() string {
	return uuid.NewString()
}
