package main

import (
	"fmt"
)

func main() {

	// var num1 []int;
	// num2:=[]int{}
	// var num3 = make([]int,3,5)
	// var num4 = make([]int,0,5)


	// fmt.Println("num 1",num1)
	// fmt.Println("num 2",num2)
	// fmt.Println("num 3",num3)
	// fmt.Println("num 4",num4)
	// fmt.Println(nums)

	// var nums1 = make([]int,3,5)
	// nums1[0] =123
	// nums1 = append(nums1,2)
	// 	nums1 = append(nums1,2)
	// 		nums1 = append(nums1,2)
	// 			nums1 = append(nums1,2)
	// 				nums1 = append(nums1,2)
	// 					nums1 = append(nums1,2)
	// 						nums1 = append(nums1,2)
	// 							nums1 = append(nums1,2)
	// 								nums1 = append(nums1,2)

	// nums1 := make([]int, 0, 5)
	// nums1 = append(nums1,2)
	// nums2 := make([]int, len(nums1))
	// fmt.Println("Length of the nums1 ---->", len(nums1))
	// fmt.Println("Capacity of the nums1 ---->", cap(nums1))
	// fmt.Println("THis is the num1----> ", nums1)


	// copy(nums2,nums1)
	// fmt.Println("Length of the nums2 ----->", len(nums2))
	// fmt.Println("Capacity of the nums2 ---->", cap(nums1))
	// fmt.Println("THis is the num2 ------>", nums2)

	// var nums = []int{1,2,3}
	// fmt.Println(nums[:])

	//slice package
	// 	var num1 = []int{1,2,3}
	// var num2 = []int{1,2,3}
	// fmt.Println(slices.Equal(num1,num2))

	//2d slice
			var num1 = [][]int{{1,2,3},{4,5,6}}
	fmt.Println((num1))
}
