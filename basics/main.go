package main

import "fmt"

func main() {
	i := 3
	switch i {
	case 1, 2:
		fmt.Println("Hleoll")
	case 3:
		fmt.Println("asdasd")
	}

	whoami := func(i interface{}) {
		switch t := i.(type) {
		case bool:
			fmt.Println("this si the int")
		default:
			fmt.Println("its others,,,",t)
		}
	}
	whoami(123123)


}















// package main

// import "fmt"

// type message struct {
// 	msg1 string
// 	msg2 string
// 	msg3 string
// }

// func acceptArrayArg(m []message) []message {
// 	return m
// }

// func arrayDemo() [3]string {
// 	return [3]string{"qww", "sf", "SDFsd"}
// }
// func main() {
// 	// data :=[6]int{1,2,3,4,5,6}
// 	// for _,data := range data{
// 	// 	if data == 5 {
// 	// 		break
// 	// 	}
// 	// 	// fmt.Println(data)
// 	// }

// 	/**Passing Array as Arguments **/
// // 	var data = []message{
// // 		{msg1: "hello world", msg2: "hello world 2", msg3: "hello world 3"},
// // 	}
// // fmt.Println(data[0].msg1, data[0].msg2)
// 	// fmt.Println(arrayDemo())
// 	// fmt.Println(acceptArrayArg(data))

// 	var demo=324
// 	fmt.Println(demo)
// }
