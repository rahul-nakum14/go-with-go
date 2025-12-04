package main

import (
	"fmt" 
)

type User struct {
	username  string
	number int
}


func test(data User){
	fmt.Println("Hello", data)
}
func main() {
	fmt.Println("Hello, World!")
	// result := concatStrings("Hello, ", "Go!")
	// fmt.Println(result)

	//  increment(5)
	// fmt.Println("Incremented Value:", incrementValue)

	fmt.Println("Is 4 even?", isEven(4))
	//  a:=[]int{1,2,3,4,5}
	//  var tmp = sliceAddition(a)
	//  fmt.Println(tmp)

	// x, y := getName();
	// fmt.Println(x, y)


	test(User{username: "sssssssssss", number: 1231231231})

}

func getName() (string, string){
	return  "Helo" , "Rahul"
}
func sliceAddition(a []int) int {
	sum:=0
	for _,tmp := range a{
		sum = sum+tmp
	}
	return  sum
}
func isEven(n int) bool {
	return n%2 == 0
}

func increment(x int)  {
	 x++
	fmt.Println("This line will never be executed",x)
}
func concatStrings(a, b string) string {
	return a + b
}