# Shift File Format

The exact input format `wageparser.ParseFile` expects, one shift per line.

## Line format

```
DD.MM Weekday: HH:MM - HH:MM
```

Optionally followed by a break annotation:

```
DD.MM Weekday: HH:MM - HH:MM (XXmin Break)
```

- Blank lines are ignored (e.g. to separate weeks)
- `DD.MM` and `Weekday` aren't validated against a real calendar - they're just the two space-separated words before the `:`. Write whatever you want there, as long as neither one contains a space.

## Time tokens

Each side of the `-` is normalized independently:

| Input    | Becomes  | Notes                        |
|----------|----------|-------------------------------|
| `9`      | `09:00`  | single-digit hour             |
| `09`     | `09:00`  | two-digit hour                |
| `9:50`   | `09:50`  | single-digit hour + minutes   |
| `09:50`  | `09:50`  | used as-is                    |

Anything else (non-numeric, or a different length) is a parse error.

## Break annotation

Optional `(...)` at the end of the line, e.g. `(41min Break)` or `(30 Min Break)`.

- Spaces are removed, then letters are trimmed off both ends of what's left — whatever remains must be a plain number, interpreted as minutes
- Wording and spacing are flexible (`41min`, `41 min`, `Break 41min` all work), as long as exactly one number is in there
- An opening `(` without a matching closing `)` (or vice versa) is a parse error, not silently ignored
- An empty `()` is also a parse error — there's no number to parse

## Error handling

`ParseFile` stops at the first invalid line and returns an error. There's no partial result for a file with one broken line, so fix the line and re-run.

## Try it yourself

```bash
cp shifts.example.txt shifts.txt
go run ./cmd/wageparser
```

`cmd/wageparser/main.go` currently points `ParseFile` at a hardcoded `./shifts.txt`, so the example file needs to be copied there first.
