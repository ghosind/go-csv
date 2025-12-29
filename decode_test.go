package csv_test

import (
	"testing"

	"github.com/ghosind/go-assert"
	"github.com/ghosind/go-csv"
)

func TestDecodeStruct(t *testing.T) {
	a := assert.New(t)
	data := "id,name,age,salary,is_manager\n1,John Doe,30,5500,true\n"
	var sample SampleStruct
	err := csv.Unmarshal([]byte(data), &sample)
	a.NilNow(err)
	expected := SampleStruct{
		ID:        1,
		Name:      "John Doe",
		Age:       30,
		Salary:    5500,
		IsManager: true,
	}
	a.EqualNow(expected, sample)
}

func TestDecodeStructPointer(t *testing.T) {
	a := assert.New(t)
	data := "id,name,age,salary,is_manager\n1,John Doe,30,5500,true\n"
	sample := new(SampleStruct)
	err := csv.Unmarshal([]byte(data), sample)
	a.NilNow(err)
	expected := SampleStruct{
		ID:        1,
		Name:      "John Doe",
		Age:       30,
		Salary:    5500,
		IsManager: true,
	}
	a.EqualNow(expected, *sample)
}

func TestEmptyCSVToStruct(t *testing.T) {
	a := assert.New(t)
	data := "id,name,age,salary,is_manager\n"
	var sample SampleStruct
	err := csv.Unmarshal([]byte(data), &sample)
	a.NilNow(err)
	expected := SampleStruct{}
	a.EqualNow(expected, sample)
}

func TestDecodeStructSlice(t *testing.T) {
	a := assert.New(t)
	data := "id,name,age,salary,is_manager\n1,John Doe,30,5500,true\n2,Jane Smith,25,3000,false\n"
	var samples []SampleStruct
	err := csv.Unmarshal([]byte(data), &samples)
	a.NilNow(err)
	expected := []SampleStruct{
		{ID: 1, Name: "John Doe", Age: 30, Salary: 5500, IsManager: true},
		{ID: 2, Name: "Jane Smith", Age: 25, Salary: 3000, IsManager: false},
	}
	a.EqualNow(expected, samples)
}

func TestDecodeStructSliceLargerThanCap(t *testing.T) {
	a := assert.New(t)
	data := "id,name,age,salary,is_manager\n1,John Doe,30,5500,true\n2,Jane Smith,25,3000,false\n"
	samples := make([]SampleStruct, 0, 1)
	err := csv.Unmarshal([]byte(data), &samples)
	a.NilNow(err)
	expected := []SampleStruct{
		{ID: 1, Name: "John Doe", Age: 30, Salary: 5500, IsManager: true},
		{ID: 2, Name: "Jane Smith", Age: 25, Salary: 3000, IsManager: false},
	}
	a.EqualNow(expected, samples)
}

func TestDecodeStructSliceWithoutHeader(t *testing.T) {
	a := assert.New(t)
	data := "1,John Doe,30,5500,true\n2,Jane Smith,25,3000,false\n"
	samples := make([]SampleStruct, 0, 1)
	err := csv.Unmarshal([]byte(data), &samples)
	a.NilNow(err)
	expected := []SampleStruct{
		{ID: 1, Name: "John Doe", Age: 30, Salary: 5500, IsManager: true},
		{ID: 2, Name: "Jane Smith", Age: 25, Salary: 3000, IsManager: false},
	}
	a.EqualNow(expected, samples)
}

func TestDecodeStructArray(t *testing.T) {
	a := assert.New(t)
	data := "id,name,age,salary,is_manager\n1,John Doe,30,5500,true\n2,Jane Smith,25,3000,false\n"
	var samples [2]SampleStruct
	err := csv.Unmarshal([]byte(data), &samples)
	a.NilNow(err)
	expected := [2]SampleStruct{
		{ID: 1, Name: "John Doe", Age: 30, Salary: 5500, IsManager: true},
		{ID: 2, Name: "Jane Smith", Age: 25, Salary: 3000, IsManager: false},
	}
	a.EqualNow(expected, samples)
}

func TestDecodeStructArrayLargerThanSize(t *testing.T) {
	a := assert.New(t)
	data := "id,name,age,salary,is_manager\n1,John Doe,30,5500,true\n2,Jane Smith,25,3000,false\n"
	var samples [1]SampleStruct
	err := csv.Unmarshal([]byte(data), &samples)
	a.NilNow(err)
	expected := [1]SampleStruct{
		{ID: 1, Name: "John Doe", Age: 30, Salary: 5500, IsManager: true},
	}
	a.EqualNow(expected, samples)
}
