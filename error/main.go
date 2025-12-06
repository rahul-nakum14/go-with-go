package main

import (
	"fmt"
)

type User struct{
	username string
	password string
}

type userError struct{
	User
}

func (e userError)Error() string  {
	return fmt.Sprintf("something went wrong with user: %v", e.User)
}

func (u User) getUserDetails() (string, error) {
	if u.username == "" {
		return "", userError{User: u}
	}

	return "hello", nil
}

func calcuSumOfN(n int) int {
	sum :=0
	// for i:=0;i<n;i++{
	// 	sum=sum+i
	// }
	i:=0
	for i<n{
		sum=sum+i;
		i++
	}
	return sum
}
func main() {
	// details := User{username: "", password: ""}
	// user, err := details.getUserDetails()
	// if (err != nil){
	// 	fmt.Println("error occured",err)
	// }
	// fmt.Println(user)

	fmt.Println(calcuSumOfN(5))
}