package main

import (
	"fmt"
	"math"
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
	HoursInMinutes int     `json:"-"`
	BreakInMinutes int     `json:"breakInMinutes"`
}

type TotalShift struct {
	Shifts              []Shift `json:"shifts"`
	TotalHoursInDecimal float64 `json:"totalHours"`
	TotalHoursInTime    string  `json:"totalHoursInTime"`
	TotalBreakInMinutes int     `json:"totalBreakInMinutes"`
	EstimatedSalary     float64 `json:"estimatedSalary"`
}

func newShift(num int, day, date, startTime, endTime string, breakMin int) (Shift, error) {
	startTimeDecimal, err := getMinutesFromTime(startTime)
	if err != nil {
		return Shift{}, err
	}

	endTimeDecimal, err := getMinutesFromTime(endTime)
	if err != nil {
		return Shift{}, err
	}

	totalMinutes := (endTimeDecimal - startTimeDecimal) - breakMin
	totalHoursDecimal := float64(totalMinutes) / 60
	totalHoursTime := getTimeFromMinutes(totalMinutes)

	return Shift{
		Line:           num,
		Day:            day,
		Date:           date,
		StartTime:      startTime,
		EndTime:        endTime,
		HoursInTime:    totalHoursTime,
		HoursInDecimal: roundToTwoDecimals(totalHoursDecimal),
		HoursInMinutes: totalMinutes,
		BreakInMinutes: breakMin,
	}, nil
}

func newTotalShift(shifts []Shift) TotalShift {
	totalHoursInMinutes := calcFieldInt(shifts, func(s Shift) int { return s.HoursInMinutes })
	totalBreakInMinutes := calcFieldInt(shifts, func(s Shift) int { return s.BreakInMinutes })
	totalHoursInTime := getTimeFromMinutes(totalHoursInMinutes)

	totalHoursInDecimal := float64(totalHoursInMinutes) / 60
	estimatedSalary := totalHoursInDecimal * hourlyWage

	return TotalShift{
		Shifts:              shifts,
		TotalHoursInDecimal: roundToTwoDecimals(totalHoursInDecimal),
		TotalHoursInTime:    totalHoursInTime,
		TotalBreakInMinutes: totalBreakInMinutes,
		EstimatedSalary:     roundToTwoDecimals(estimatedSalary),
	}
}

func calcFieldInt(shifts []Shift, field func(Shift) int) int {
	sum := 0
	for i := range shifts {
		sum += field(shifts[i])
	}

	return sum
}

func getTimeFromMinutes(minutes int) string {
	hours := minutes / 60

	return fmt.Sprintf("%02d:%02d", hours, minutes%60)
}

func getMinutesFromTime(time string) (int, error) {
	parts := strings.Split(time, ":")
	if len(parts) != 2 {
		return 0.0, fmt.Errorf("invalid Time Format: %s", time)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0.0, err
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0.0, err
	}

	return hours*60 + minutes, nil
}

func roundToTwoDecimals(x float64) float64 {
	return math.Round(x*100) / 100
}
