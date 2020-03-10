package main

import "fmt"

func main()  {
	var f1 = func(x, y int) {
		fmt.Println(x + y)
	}
	f1(10, 20)

	func(x, y int) {
		fmt.Println(x + y)
	}(1, 2)
}
