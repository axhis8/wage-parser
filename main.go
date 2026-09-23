package main

import "fmt"

func main() {
	test, _ := ParseFile("./shifts.txt")
	fmt.Println(test)
}
