package wageparser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// ParseFile reads shift data from the text file at path, parses each line,
// and returns the aggregated totals as a TotalShift, using hourlyWage to
// calculate the estimated salary.
func ParseFile(path string, hourlyWage float64) (totalShift TotalShift, err error) {
	var shifts []Shift

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

		shift, err := parseLine(scanner.Text(), len(shifts)+1)
		if err != nil {
			return totalShift, err
		}

		shifts = append(shifts, shift)
	}

	if err := scanner.Err(); err != nil {
		return totalShift, err
	}

	totalShift = newTotalShift(shifts, hourlyWage)
	return totalShift, nil
}

func parseLine(line string, lineNum int) (shift Shift, err error) {
	date, day, remainder, err := parseHeader(line)
	if err != nil {
		return shift, err
	}

	remainder, breakStr, err := splitRemainder(remainder)
	if err != nil {
		return shift, err
	}

	breakMin, err := parseBreakLine(breakStr)
	if err != nil {
		return shift, err
	}

	startTime, endTime, err := splitTimeRange(remainder)
	if err != nil {
		return shift, err
	}

	shift, err = newShift(lineNum, day, date, startTime, endTime, breakMin)
	return shift, err
}

func parseBreakLine(breakStr string) (int, error) {
	if breakStr == "" {
		return 0, nil
	}

	normalizedBreakStr := breakStr[1 : len(breakStr)-1]
	normalizedBreakStr = strings.ReplaceAll(normalizedBreakStr, " ", "")
	normalizedBreakStr = strings.TrimFunc(normalizedBreakStr, unicode.IsLetter)

	breakMin, err := strconv.Atoi(normalizedBreakStr)
	return breakMin, err
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

// Splits the remainder and break as two strings
func splitRemainder(remainder string) (cleanRemainder, breakStr string, err error) {
	start := strings.Index(remainder, "(")
	end := strings.Index(remainder, ")")

	if start == -1 && end == -1 {
		return remainder, breakStr, nil
	} else if start == -1 || end == -1 {
		return remainder, breakStr, fmt.Errorf("invalid break %q", remainder)
	}

	breakStr = remainder[start : end+1]
	cleanRemainder = strings.ReplaceAll(remainder, breakStr, "")

	return cleanRemainder, breakStr, nil
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
