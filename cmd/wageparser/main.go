package main

import (
	"encoding/json"
	"fmt"
	"github.com/axhis8/wage-parser/wageparser"
)

func main() {
	totalShifts, err := wageparser.ParseFile("./shifts.txt", 12)
	if err != nil {
		panic(err)
	}

	data, err := json.MarshalIndent(totalShifts, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
