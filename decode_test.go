package csv_test

import (
	"testing"

	"github.com/foodieats/go-csv"
	"github.com/ghosind/go-assert"
)

func TestDecodeStruct(t *testing.T) {
	a := assert.New(t)
	data := "id,name,email\n1,John Doe,john@example.com\n"
	var sample SampleStruct
	err := csv.Unmarshal([]byte(data), &sample)
	a.NilNow(err)
	expected := SampleStruct{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}
	a.EqualNow(expected, sample)
}

func TestDecodeStructPointer(t *testing.T) {
	a := assert.New(t)
	data := "id,name,email\n1,John Doe,john@example.com\n"
	sample := new(SampleStruct)
	err := csv.Unmarshal([]byte(data), sample)
	a.NilNow(err)
	expected := SampleStruct{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}
	a.EqualNow(expected, *sample)
}

func TestDecodeStructSlice(t *testing.T) {
	a := assert.New(t)
	data := "id,name,email\n1,John Doe,john@example.com\n2,Jane Smith,jane@example.com\n"
	var samples []SampleStruct
	err := csv.Unmarshal([]byte(data), &samples)
	a.NilNow(err)
	expected := []SampleStruct{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
	}
	a.EqualNow(expected, samples)
}

func TestDecodeStructSliceLargerThanCap(t *testing.T) {
	a := assert.New(t)
	data := "id,name,email\n1,John Doe,john@example.com\n2,Jane Smith,jane@example.com\n"
	samples := make([]SampleStruct, 0, 1)
	err := csv.Unmarshal([]byte(data), &samples)
	a.NilNow(err)
	expected := []SampleStruct{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
	}
	a.EqualNow(expected, samples)
}

func TestDecodeStructArray(t *testing.T) {
	a := assert.New(t)
	data := "id,name,email\n1,John Doe,john@example.com\n2,Jane Smith,jane@example.com\n"
	var samples [2]SampleStruct
	err := csv.Unmarshal([]byte(data), &samples)
	a.NilNow(err)
	expected := [2]SampleStruct{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
	}
	a.EqualNow(expected, samples)
}

func TestDecodeStructArrayLargerThanSize(t *testing.T) {
	a := assert.New(t)
	data := "id,name,email\n1,John Doe,john@example.com\n2,Jane Smith,jane@example.com\n"
	var samples [1]SampleStruct
	err := csv.Unmarshal([]byte(data), &samples)
	a.NilNow(err)
	expected := [1]SampleStruct{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
	}
	a.EqualNow(expected, samples)
}
