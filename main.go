package main

import (
	"encoding/json"
	"fmt"
)

type Shift struct {
	Line           int     `json:"line"`
	Day            string  `json:"day"`
	Date           string  `json:"date"`
	StartTime      string  `json:"startTime"`
	EndTime        string  `json:"endTime"`
	HoursInTime    string  `json:"hoursInTime"`
	HoursInDecimal float64 `json:"hours"`
}

type TotalShift struct {
	Shifts              []Shift `json:"shifts"`
	TotalHoursInDecimal float64 `json:"totalHours"`
	TotalHoursInTime    float64 `json:"totalHoursInTime"`
	EstimatedSalary     float64 `json:"estimatedSalary"`
}

func main() {
	s := Shift{
		Line:           1,
		Day:            "Montag",
		Date:           "27.07",
		StartTime:      "09:50",
		EndTime:        "13:50",
		HoursInTime:    "4:00",
		HoursInDecimal: 4,
	}

	jsonData, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return
	}

	fmt.Println(string(jsonData))
}
