package main

import "fmt"


func main (){
	c:= make(map[string]string)

	c["helo"] = "123"
	c["price"] ="sdfs"
	fmt.Println(c["helo"])
	fmt.Println(c["price"])

	fmt.Println(len(c))

	fmt.Println(c)
	delete(c, "price")
	clear(c)
	fmt.Println(c)

	k:= map[string]int{"asd":19}

	err,ok:= k["assd"]

	if ok {
		fmt.Println(" ok")
	}else {
		fmt.Println("Not  ok,err",err)
	}
}