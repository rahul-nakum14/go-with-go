package main

import (
	"fmt"
	"errors"
)

type CustomError struct {
	Message string
	Code    int
}

func (e CustomError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}
func main() {
	result, err := Divide(4, 0)

	if err != nil {
		fmt.Println(err)
		// Do something with the error
		return
	}

	fmt.Println(result)
}
var ErrDivideByZero = errors.New("cannot divide by zero")


func Divide(a, b int) (int, error) {
	if b == 0 {
		// return 0, errors.New("cannot divide by zero") -- way 1
		// return 0, fmt.Errorf("cannot divide %d by zero", a) -- way 2
		// return 0, ErrDivideByZero // way 3
		return 0, CustomError{Message: "cannot divide by zero", Code: 400} // way 4
	}

	return a/b, nil
}


/**
	old code 
**/
// type User struct{
// 	username string
// 	password string
// }

// type userError struct{
// 	User
// }

// func (e userError)Error() string  {
// 	return fmt.Sprintf("something went wrong with user: %v", e.User)
// }

// func (u User) getUserDetails() (string, error) {
// 	if u.username == "" {
// 		return "", userError{User: u}
// 	}

// 	return "hello", nil
// }

// func calcuSumOfN(n int) int {
// 	sum :=0
// 	// for i:=0;i<n;i++{
// 	// 	sum=sum+i
// 	// }
// 	i:=0
// 	for i<n{
// 		sum=sum+i;
// 		i++
// 	}
// 	return sum
// }
// func main() {
// 	details := User{username: "", password: ""}
// 	user, err := details.getUserDetails()
// 	if (err != nil){
// 		fmt.Println("error occured",err)
// 	}
// 	fmt.Println(user)

// 	// fmt.Println(calcuSumOfN(5))
// }