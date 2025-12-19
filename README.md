# csv

![test](https://github.com/ghosind/go-csv/workflows/test/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/ghosind/go-csv)](https://goreportcard.com/report/github.com/ghosind/go-csv)
[![codecov](https://codecov.io/gh/ghosind/go-csv/branch/main/graph/badge.svg)](https://codecov.io/gh/ghosind/go-csv)
![Version Badge](https://img.shields.io/github/v/release/ghosind/go-csv)
![License Badge](https://img.shields.io/github/license/ghosind/go-csv)
[![Go Reference](https://pkg.go.dev/badge/github.com/ghosind/go-csv.svg)](https://pkg.go.dev/github.com/ghosind/go-csv)

The Go library for serializing/deserializing CSV to/from Go structs.

## Features

- Map CSV headers to struct fields using `csv` tags (or field name when tag omitted).
- Support basic types: string, ints, uints, floats, bool.
- Easy to use API for marshaling and unmarshaling.

## Installation

Run the following command to install the package:

```bash
go get github.com/ghosind/go-csv
```

## Getting Started

Here's a quick example to get you started:

1. Define a struct with `csv` tags:

```go
type Person struct {
  Name string `csv:"name"`
  Age  int    `csv:"age"`
}
```

2. Prepare CSV data:

```csv
name,age
Alice,30
Bob,25
```

3. Unmarshal CSV into a slice:

```go
var people []Person
err := csv.Unmarshal([]byte(data), &people)
log.Print(people)
// Output:
// [{Alice 30} {Bob 25}]
```

4. Marshal a slice back to CSV:

```go
data, err := csv.Marshal(people)
log.Print(string(data))
// Output:
// name,age
// Alice,30
// Bob,25
```

## Tests

Run tests using the following command:

```bash
go test ./...
```

## License

This project is licensed under the MIT License, see the [LICENSE](./LICENSE) file for details.
