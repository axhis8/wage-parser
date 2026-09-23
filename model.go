package main

import (
	"fmt"
	"strconv"
	"strings"
)

const hourlyWage float64 = 12.00

type Shift struct {
	Line           int     `json:"line"`
	Day            string  `json:"day"`
	Date           string  `json:"date"`
	StartTime      string  `json:"startTime"`
	EndTime        string  `json:"endTime"`
	HoursInTime    string  `json:"hoursInTime"`
	HoursInDecimal float64 `json:"hoursInDecimal"`
}

type TotalShift struct {
	Shifts              []Shift `json:"shifts"`
	TotalHoursInDecimal float64 `json:"totalHours"`
	TotalHoursInTime    string  `json:"totalHoursInTime"`
	EstimatedSalary     float64 `json:"estimatedSalary"`
}

func NewShift(line int, day, date, startTime, endTime string) (Shift, error) {
	startTimeDecimal, err := getDecimalFromTime(startTime)
	if err != nil {
		return Shift{}, err
	}

	endTimeDecimal, err := getDecimalFromTime(endTime)
	if err != nil {
		return Shift{}, err
	}

	totalHours := endTimeDecimal - startTimeDecimal
	totalHoursTime := getTimeFromDecimal(totalHours)

	return Shift{
		Line:           line,
		Day:            day,
		Date:           date,
		StartTime:      startTime,
		EndTime:        endTime,
		HoursInTime:    totalHoursTime,
		HoursInDecimal: totalHours,
	}, nil
}

func NewTotalShift(shifts []Shift) TotalShift {
	totalHoursInDecimal := calcTotalHours(shifts)
	totalHoursInTime := getTimeFromDecimal(totalHoursInDecimal)
	estimatedSalary := totalHoursInDecimal * hourlyWage

	return TotalShift{
		Shifts:              shifts,
		TotalHoursInDecimal: totalHoursInDecimal,
		TotalHoursInTime:    totalHoursInTime,
		EstimatedSalary:     estimatedSalary,
	}
}

func calcTotalHours(shifts []Shift) float64 {
	sum := 0.0
	for i := range shifts {
		sum += shifts[i].HoursInDecimal
	}

	return sum
}

func getTimeFromDecimal(hours float64) string {
	hour := int(hours)
	minutes := int((hours - float64(hour)) * 60)

	return fmt.Sprintf("%02d:%02d", hour, minutes)
}

func getDecimalFromTime(timeStr string) (float64, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0.0, fmt.Errorf("invalid Time Format: %s", timeStr)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0.0, err
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0.0, err
	}

	return float64(hours) + (float64(minutes) / 60.0), nil
}
