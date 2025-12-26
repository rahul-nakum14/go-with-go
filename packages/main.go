package main

import (
	"fmt"
	"github.com/rahu-nakum14/learn-gp/auth"
	"github.com/rahu-nakum14/learn-gp/user"
	"github.com/fatih/color"
)

func main(){
	auth.CredentialsInfo("rahul","test123")
	fmt.Println(auth.GetSessionInfo())

	user := user.User{
		Username: "rahulnakum14",
		Email:    "rahul@example.com",
	}
	fmt.Printf("User: %+v\n", user)
	color.Cyan("Prints text in cyan.")
}