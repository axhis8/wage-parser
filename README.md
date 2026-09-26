# wage-parser

Small Go tool that parses a specific handwritten work shift log into calculated worked hours, breaks, and an estimated wage as JSON.
This was my first real Go project. I used it to learn the language hands-on.

## What it does

`wage-parser` reads a plain text file of shifts and turns it into structured, computed output:

- Parses shift lines (date, day, start/end time, optional break) from a simple text format
- Computes worked hours per shift and in total, net of breaks
- Does all duration math in exact integer minutes internally to avoid floating-point drift
- Calculates an estimated salary from an hourly wage
- Outputs everything as JSON, shaped to double as a future backend API response

## Input format

```
25.07 Saturday: 11:07 - 14:05
03.08 Monday: 10:00 - 16:00 (41min Break)
```

- `DD.MM Weekday: HH:MM - HH:MM` per shift, one per line
- Optional `(XXmin Break)` for a break, subtracted from the worked time
- Blank lines are ignored

See [FORMAT.md](FORMAT.md) for the exact parsing rules and edge cases, and [shifts.example.txt](shifts.example.txt) for a ready-to-use sample file.

## Installation / Usage

```bash
git clone https://github.com/axhis8/wage-parser.git
cd wage-parser
go run ./cmd/wageparser
```

Build a binary instead:

```bash
go build -o bin/wageparser ./cmd/wageparser
./bin/wageparser
```

As a library:

```go
import "github.com/axhis8/wage-parser/wageparser"

result, err := wageparser.ParseFile("shifts.txt", 12.00)
```

## Example output

```json
{
  "shifts": [
    {
      "line": 1,
      "day": "Saturday",
      "date": "25.07",
      "startTime": "11:07",
      "endTime": "14:05",
      "hoursInTime": "02:58",
      "hoursInDecimal": 2.97,
      "breakInMinutes": 0
    }
  ],
  "totalHours": 2.97,
  "totalHoursInTime": "02:58",
  "totalBreakInMinutes": 0,
  "estimatedSalary": 35.64
}
```

## Project structure

```
wage-parser/
├── cmd/
│   └── wageparser/
│       └── main.go      # CLI entry point
├── wageparser/
│   ├── doc.go            # package doc comment
│   ├── model.go          # Shift/TotalShift types + calculations
│   └── parser.go         # text file parsing
└── go.mod
```

## Tech stack

- Go 1.27
- Standard library only

## What I learned

First Go project, so a good chunk of the learning happened here:

- Go's slice semantics
- Avoiding floating-point accumulation errors by doing all duration math in integer minutes, converting to decimal hours only at the final display step
- Structuring a repo as an importable library (`wageparser/`)
