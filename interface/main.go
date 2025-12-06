package main

import "fmt"

type shape interface {
	area() int
}

type Greeter interface {
	Greet()
}

type Person struct {
	name string
}
type rectangle struct {
	width  int
	height int
}
type circle struct {
	radius int
}

func (r rectangle) area() int {
	return r.height * r.width
}

func (c circle) area() int {
	return c.radius * 2
}

func (p Person) Greet() {
	fmt.Println("Hello, my name is", p.name)
}

func main() {
	// p := Person{name: "Alice"}

	// var g Greeter
	// g = p

	// g.Greet()
    // Create values
    r := rectangle{width: 10, height: 5}
    c := circle{radius: 20}

    // Use interface variable
    var s shape

    s = r
    fmt.Println("Rectangle area:", s.area())

    s = c
    fmt.Println("Circle area:", s.area())

}
