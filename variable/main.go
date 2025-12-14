package main

import "fmt"

func variableInitialization()  {
	// var username string
	// var hasPermission bool
	// var constData float64
	// var age int
	
	// fmt.Printf("%q %v %f %v\n", username, hasPermission, constData, age)


	/* 
	 Short Assignment Variable Declaration 
	*/

	shortAssignment := 123.234234
	fmt.Printf("Short Assignment ---- %v %T\n",shortAssignment, shortAssignment)

	/*
		Same Line Assignment operator
	*/
	x , y := 10,20
	fmt.Println(x,y)
	x , z := 50,60 // Note:  here x value will be update to 50
	fmt.Println(x,z)

	/*
		Type Conversion
	*/

	var constDetails  = 3.14
	fmt.Println(int(constDetails))

	/*
	Constants
	*/

	const demoPlan = "planName"
	fmt.Printf(demoPlan)

	/*
		PrintF and Sprintf  - Here Sprintf will return string instead print
	*/

	var name = "hello"
	var ageDetails  = 12

	msg := fmt.Sprintf("hi %s your age is %d",name,ageDetails)
	fmt.Println(msg)


	/*
	Conditional
	*/
	var data1 = 5000
	var data2 = 14
	var data3 = 70

    if data1 > data2 && data1 > data3 {

		fmt.Println(data1,"is a bigger")
	}	else if data2 >data1 && data2> data3{

		fmt.Println(data2,"is a bigger")
	}	else{

		fmt.Println(data3,"is a bigger")
	}

	/*
		If the arable is only used in the if then 
		we can do like if INITIAL_STATEMENT:  CONDITION
	*/
	if c:=999; c<data1 {
		fmt.Println("Hey How are you")
	}


	var a int;
	a=12;
	fmt.Println(a)

	b:=123
	fmt.Println(b)

	var c int =2342;
	fmt.Println(c)


	const (
		port = 5000
		host = "localhost"
	)

	fmt.Println(port)
}
func main() {
	variableInitialization()
}