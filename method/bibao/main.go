package main

import "fmt"

func main()  {

}

func f1(f func()) {
	fmt.Println("this is f1")
	f()
}

func f2(x, y int) {
	fmt.Println("this is f2")
	fmt.Println(x + y)
}

func f3(x, y int) func(f func(x, y int)) {
	return func(f func(x, y int)) {
		f()
	}
}


