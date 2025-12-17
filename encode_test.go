package csv_test

import (
	"testing"

	"github.com/foodieats/go-csv"
	"github.com/ghosind/go-assert"
)

type SampleStruct struct {
	ID    int    `csv:"id"`
	Name  string `csv:"name"`
	Email string `csv:"email"`
}

func TestEncodeStruct(t *testing.T) {
	a := assert.New(t)
	sample := SampleStruct{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	data, err := csv.Marshal(sample)
	a.NilNow(err)
	expected := "id,name,email\n1,John Doe,john@example.com\n"
	a.EqualNow(expected, string(data))
}

func TestEncodeStructSlice(t *testing.T) {
	a := assert.New(t)
	samples := []SampleStruct{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
	}

	data, err := csv.Marshal(samples)
	a.NilNow(err)
	expected := "id,name,email\n1,John Doe,john@example.com\n2,Jane Smith,jane@example.com\n"
	a.EqualNow(expected, string(data))
}
