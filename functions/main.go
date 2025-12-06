package main
import "fmt"

type rectangle struct{
	width int
	height int
}

func (area rectangle)getArea() int {
	return  area.width * area.height
}

func main (){
    area := rectangle{}   // create a value of type rectangle

	area.height =20;
	area.width =20;
	fmt.Println(area.getArea())
}