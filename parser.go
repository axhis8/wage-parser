package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	}

	if err := scanner.Err(); err != nil {
		return shifts, err
	}

	return shifts, nil
}

func normalizeTimeToken(tok string) (string, error) {
	normalizedTok := strings.TrimSpace(tok)

	// tok != 12:a4
	replacedTok := strings.Replace(normalizedTok, ":", "", 1)
	if _, err := strconv.ParseInt(replacedTok, 10, 64); err != nil {
		return normalizedTok, fmt.Errorf("%q is not a number", tok)
	}

	switch len(normalizedTok) {
	case 1:
		// tok == 9
		return "0" + normalizedTok + ":00", nil
	case 2:
		// tok == 09
		return normalizedTok + ":00", nil
	case 4:
		// tok != 9:5 and tok != 09:5 and tok == 9:50
		if s := strings.Split(normalizedTok, ":"); len(s[0]) == 1 && len(s[1]) == 2 {
			return "0" + normalizedTok, nil
		}
		return normalizedTok, fmt.Errorf("invalid time token %q", tok)
	case 5:
		return normalizedTok, nil
	default:
		return normalizedTok, fmt.Errorf("invalid time token %q", tok)
	}
}
