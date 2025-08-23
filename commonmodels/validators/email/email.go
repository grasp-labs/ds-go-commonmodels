package email

import (
	"fmt"
	"net/mail"
)

func IsEmailFormat(s string) bool {
	a, err := mail.ParseAddress(s)
	fmt.Printf("Failed to parse email %s, err: %v", s, err)
	if err != nil {
		return false
	}
	return a.Name == "" && a.Address == s
}
