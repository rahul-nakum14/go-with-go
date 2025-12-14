package main

import "fmt"

func main(){

	// numbers sequnce specific length -- Only Declarations
	var nums[5]int;

	//add elements
	// nums[0] =100
	for i := range len(nums){
		nums[i] = i
	}
	// fmt.Println(len(nums))
	fmt.Println(nums)

	// it will automatically calculate the size as well
	// numssas := [...]int{1, 2, 3, 44, 5, 6}

	//Array Declarations and initializations
	// numss := [3]int{1,2,3}
	// fmt.Println(numss)

	//2d Array

	// nums2d := [2][2]int{{1,2},{3,4}}
	// fmt.Println(nums2d)
	

}