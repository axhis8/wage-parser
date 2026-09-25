package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func ParseFile(path string) (TotalShift, error) {
	var totalShift TotalShift

	file, err := os.Open(path)
	if err != nil {
		return totalShift, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}

		shift, err := parseLine(scanner.Text(), len(totalShift.Shifts)+1)
		if err != nil {
			return totalShift, err
		}

		totalShift.Shifts = append(totalShift.Shifts, shift)
	}

	if err := scanner.Err(); err != nil {
		return totalShift, err
	}

	totalShift = newTotalShift(totalShift.Shifts)
	return totalShift, nil
}

func parseLine(line string, lineNum int) (Shift, error) {
	var (
		shift    Shift
		breakMin int
	)

	date, day, remainder, err := parseHeader(line)
	if err != nil {
		return shift, err
	}

	start := strings.Index(remainder, "(")
	end := strings.Index(remainder[start+1:], ")")
	if start != -1 && end != -1 {
		breakStr := remainder[start : start+2+end]
		breakMin, err = parseBreakLine(breakStr)
		if err != nil {
			return shift, err
		}

		remainder = strings.ReplaceAll(remainder, breakStr, "")
	}

	startTime, endTime, err := splitTimeRange(remainder)
	if err != nil {
		return shift, err
	}

	shift, err = newShift(lineNum, day, date, startTime, endTime, breakMin)
	if err != nil {
		return shift, err
	}

	return shift, nil
}

func parseBreakLine(breakStr string) (int, error) {
	normalizedBreakStr := breakStr[1 : len(breakStr)-1]
	normalizedBreakStr = strings.ReplaceAll(normalizedBreakStr, " ", "")
	normalizedBreakStr = strings.TrimFunc(normalizedBreakStr, unicode.IsLetter)

	breakMin, err := strconv.Atoi(normalizedBreakStr)
	if err != nil {
		return breakMin, err
	}

	return breakMin, nil
}

func parseHeader(line string) (date, day, remainder string, err error) {
	beforeColon, afterColon, found := strings.Cut(line, ":")
	if !found {
		return date, day, remainder, fmt.Errorf("invalid header %q", line)
	}

	beforeLineSplit := strings.Split(beforeColon, " ")
	if len(beforeLineSplit) != 2 {
		return date, day, remainder, fmt.Errorf("invalid header %q", line)
	}

	remainder = afterColon
	date, day = beforeLineSplit[0], beforeLineSplit[1]

	return date, day, remainder, nil
}

func splitTimeRange(rangeStr string) (start, end string, err error) {
	splitStr := strings.Split(rangeStr, "-")
	if len(splitStr) != 2 {
		return start, end, fmt.Errorf("invalid range %q", rangeStr)
	}

	for i := range splitStr {
		splitStr[i], err = normalizeTimeToken(splitStr[i])
		if err != nil {
			return start, end, err
		}
	}

	start = splitStr[0]
	end = splitStr[1]

	return start, end, nil
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
	case 5:
		return normalizedTok, nil
	}

	return normalizedTok, fmt.Errorf("invalid time token %q", tok)
}
