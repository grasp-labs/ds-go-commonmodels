package email

import (
	"net/mail"
)

func IsEmailFormat(s string) bool {
	a, err := mail.ParseAddress(s)
	if err != nil {
		return false
	}
	return a.Name == "" && a.Address == s
}
