package main
import (
	"fmt";
	"errors"
)

type User struct{
	username string
	password string
}

func (data User)getUserDetails() (string, error) {
	if data.username != ""{
		return  fmt.Sprint(data.username, ""),nil
	}
	return "Broken", errors.New("It is not working")
}

func main() { 
	data := User{username: "", password: "qqwe"}
	// fmt.Println(data.password)
	fmt.Println(data.getUserDetails())
}
