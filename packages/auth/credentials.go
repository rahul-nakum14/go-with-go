package auth

import (
	"fmt"
)

func CredentialsInfo(username, password string) {
	fmt.Printf("Username: %s, Password: %s\n", username, password)
}