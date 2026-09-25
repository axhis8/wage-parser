package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	totalShifts, err := ParseFile("./shifts.txt")
	if err != nil {
		panic(err)
	}

	data, err := json.MarshalIndent(totalShifts, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
