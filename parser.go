package main

import (
	"bufio"
	"fmt"
	"os"
)

func ParseFile(path string) ([]Shift, error) {
	var shifts []Shift

	file, err := os.Open(path)
	if err != nil {
		return shifts, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return shifts, err
	}

	return shifts, nil
}
