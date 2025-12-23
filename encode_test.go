package csv_test

import (
	"strconv"
	"testing"

	"github.com/foodieats/go-csv"
	"github.com/ghosind/go-assert"
)

type SampleStruct struct {
	ID        int     `csv:"id"`
	Name      string  `csv:"name"`
	Age       uint    `csv:"age"`
	Salary    float64 `csv:"salary"`
	IsManager bool    `csv:"is_manager"`
}

func TestEncodeStruct(t *testing.T) {
	a := assert.New(t)
	sample := SampleStruct{
		ID:        1,
		Name:      "John Doe",
		Age:       30,
		Salary:    5500,
		IsManager: true,
	}

	data, err := csv.Marshal(sample)
	a.NilNow(err)
	expected := "id,name,age,salary,is_manager\n1,John Doe,30,5500,true\n"
	a.EqualNow(expected, string(data))
}

func TestEncodeStructZeroValues(t *testing.T) {
	a := assert.New(t)
	sample := SampleStruct{}

	data, err := csv.Marshal(sample)
	a.NilNow(err)
	expected := "id,name,age,salary,is_manager\n0,,0,0,false\n"
	a.EqualNow(expected, string(data))
}

func TestEncodeStructPointer(t *testing.T) {
	a := assert.New(t)
	sample := &SampleStruct{
		ID:        1,
		Name:      "John Doe",
		Age:       30,
		Salary:    5500,
		IsManager: true,
	}

	data, err := csv.Marshal(sample)
	a.NilNow(err)
	expected := "id,name,age,salary,is_manager\n1,John Doe,30,5500,true\n"
	a.EqualNow(expected, string(data))
}

func TestEncodeStructPointerNil(t *testing.T) {
	a := assert.New(t)
	var sample *SampleStruct = nil

	data, err := csv.Marshal(sample)
	a.NilNow(err)
	expected := "id,name,age,salary,is_manager\n,,,,\n"
	a.EqualNow(expected, string(data))
}

func TestEncodeStructSlice(t *testing.T) {
	a := assert.New(t)
	samples := []SampleStruct{
		{ID: 1, Name: "John Doe", Age: 30, Salary: 5500, IsManager: true},
		{ID: 2, Name: "Jane Smith", Age: 25, Salary: 3000, IsManager: false},
	}

	data, err := csv.Marshal(samples)
	a.NilNow(err)
	expected := "id,name,age,salary,is_manager\n1,John Doe,30,5500,true\n2,Jane Smith,25,3000,false\n"
	a.EqualNow(expected, string(data))
}

func TestEncodeEmptySlice(t *testing.T) {
	a := assert.New(t)
	expected := "id,name,age,salary,is_manager\n"

	data, err := csv.Marshal([]SampleStruct{})
	a.NilNow(err)
	a.EqualNow(expected, string(data))
}

func TestEncodeStructSlicePointer(t *testing.T) {
	a := assert.New(t)
	samples := []*SampleStruct{
		{ID: 1, Name: "John Doe", Age: 30, Salary: 5500, IsManager: true},
		{ID: 2, Name: "Jane Smith", Age: 25, Salary: 3000, IsManager: false},
	}

	data, err := csv.Marshal(samples)
	a.NilNow(err)
	expected := "id,name,age,salary,is_manager\n1,John Doe,30,5500,true\n2,Jane Smith,25,3000,false\n"
	a.EqualNow(expected, string(data))
}

func TestEncodeStructArray(t *testing.T) {
	a := assert.New(t)
	samples := [2]SampleStruct{
		{ID: 1, Name: "John Doe", Age: 30, Salary: 5500, IsManager: true},
		{ID: 2, Name: "Jane Smith", Age: 25, Salary: 3000, IsManager: false},
	}

	data, err := csv.Marshal(samples)
	a.NilNow(err)
	expected := "id,name,age,salary,is_manager\n1,John Doe,30,5500,true\n2,Jane Smith,25,3000,false\n"
	a.EqualNow(expected, string(data))
}

func TestEncodeStructArrayPointer(t *testing.T) {
	a := assert.New(t)
	samples := [2]*SampleStruct{
		{ID: 1, Name: "John Doe", Age: 30, Salary: 5500, IsManager: true},
		{ID: 2, Name: "Jane Smith", Age: 25, Salary: 3000, IsManager: false},
	}

	data, err := csv.Marshal(samples)
	a.NilNow(err)
	expected := "id,name,age,salary,is_manager\n1,John Doe,30,5500,true\n2,Jane Smith,25,3000,false\n"
	a.EqualNow(expected, string(data))
}

func TestEncodeNil(t *testing.T) {
	a := assert.New(t)

	data, err := csv.Marshal(nil)
	a.NilNow(err)
	a.EqualNow("", string(data))
}

type SimplePointerStruct struct {
	ID        *int     `csv:"id"`
	Name      *string  `csv:"name"`
	Age       *uint    `csv:"age"`
	Salary    *float64 `csv:"salary"`
	IsManager *bool    `csv:"is_manager"`
}

func TestEncodeStructWithPointerFields(t *testing.T) {
	a := assert.New(t)
	id, name, age, salary, isManager := 1, "John Doe", uint(30), 5500.0, true
	sample := SimplePointerStruct{
		ID:        &id,
		Name:      &name,
		Age:       &age,
		Salary:    &salary,
		IsManager: &isManager,
	}

	data, err := csv.Marshal(sample)
	a.NilNow(err)
	expected := "id,name,age,salary,is_manager\n1,John Doe,30,5500,true\n"
	a.EqualNow(expected, string(data))
}

func TestEncodeStructWithNilPointerFields(t *testing.T) {
	a := assert.New(t)
	id, name := 1, "John Doe"
	sample := SimplePointerStruct{
		ID:   &id,
		Name: &name,
	}

	data, err := csv.Marshal(sample)
	a.NilNow(err)
	expected := "id,name,age,salary,is_manager\n1,John Doe,,,\n"
	a.EqualNow(expected, string(data))
}

type MarshalableStruct struct {
	Country string
	ZipCode int
}

func (m MarshalableStruct) MarshalCSV() ([]byte, error) {
	return []byte(m.Country + " (" + strconv.Itoa(m.ZipCode) + ")"), nil
}

type WrapMarshalableStruct struct {
	City     string            `csv:"city"`
	Location MarshalableStruct `csv:"location"`
}

func TestEncodeStructWithMarshalableField(t *testing.T) {
	a := assert.New(t)
	sample := WrapMarshalableStruct{
		City:     "New York",
		Location: MarshalableStruct{Country: "USA", ZipCode: 10001},
	}

	data, err := csv.Marshal(sample)
	a.NilNow(err)
	expected := "city,location\nNew York,USA (10001)\n"
	a.EqualNow(expected, string(data))
}

type WrapMarshalablePointerStruct struct {
	City     string             `csv:"city"`
	Location *MarshalableStruct `csv:"location"`
}

func TestEncodeStructWithMarshalablePointerField(t *testing.T) {
	a := assert.New(t)
	sample := WrapMarshalablePointerStruct{
		City:     "New York",
		Location: &MarshalableStruct{Country: "USA", ZipCode: 10001},
	}

	data, err := csv.Marshal(sample)
	a.NilNow(err)
	expected := "city,location\nNew York,USA (10001)\n"
	a.EqualNow(expected, string(data))

	sample = WrapMarshalablePointerStruct{
		City:     "New York",
		Location: nil,
	}

	data, err = csv.Marshal(sample)
	a.NilNow(err)
	expected = "city,location\nNew York,\n"
	a.EqualNow(expected, string(data))
}

type TextMarshalStruct struct {
	Country string
	ZipCode int
}

func (m TextMarshalStruct) MarshalText() ([]byte, error) {
	return []byte(m.Country + " [" + strconv.Itoa(m.ZipCode) + "]"), nil
}

type WrapTextMarshalStruct struct {
	City     string            `csv:"city"`
	Location TextMarshalStruct `csv:"location"`
}

func TestEncodeStructWithTextMarshalField(t *testing.T) {
	a := assert.New(t)
	sample := WrapTextMarshalStruct{
		City:     "New York",
		Location: TextMarshalStruct{Country: "USA", ZipCode: 10001},
	}

	data, err := csv.Marshal(sample)
	a.NilNow(err)
	expected := "city,location\nNew York,USA [10001]\n"
	a.EqualNow(expected, string(data))
}

type WrapTextMarshalPointerStruct struct {
	City     string             `csv:"city"`
	Location *TextMarshalStruct `csv:"location"`
}

func TestEncodeStructWithTextMarshalPointerField(t *testing.T) {
	a := assert.New(t)
	sample := WrapTextMarshalPointerStruct{
		City:     "New York",
		Location: &TextMarshalStruct{Country: "USA", ZipCode: 10001},
	}

	data, err := csv.Marshal(sample)
	a.NilNow(err)
	expected := "city,location\nNew York,USA [10001]\n"
	a.EqualNow(expected, string(data))

	sample = WrapTextMarshalPointerStruct{
		City:     "New York",
		Location: nil,
	}

	data, err = csv.Marshal(sample)
	a.NilNow(err)
	expected = "city,location\nNew York,\n"
	a.EqualNow(expected, string(data))
}

type UnexportedFieldStruct struct {
	ID    int    `csv:"id"`
	name  string `csv:"name"`
	Email string `csv:"email"`
}

func TestEncodeStructWithUnexportedField(t *testing.T) {
	a := assert.New(t)
	sample := UnexportedFieldStruct{
		ID:    1,
		name:  "John Doe",
		Email: "john@example.com",
	}

	data, err := csv.Marshal(sample)
	a.NilNow(err)
	expected := "id,email\n1,john@example.com\n"
	a.EqualNow(expected, string(data))
}

type UnsupportedStruct struct {
	Data map[string]int
}

func TestEncodeStructWithUnsupportedField(t *testing.T) {
	a := assert.New(t)
	sample := UnsupportedStruct{
		Data: map[string]int{"a": 1, "b": 2},
	}

	_, err := csv.Marshal(sample)
	a.IsErrorNow(err, csv.ErrUnsupportedType)
}

type NoTagStruct struct {
	ID   int
	Name string
}

func TestEncodeStructWithoutTags(t *testing.T) {
	a := assert.New(t)
	sample := NoTagStruct{
		ID:   1,
		Name: "John Doe",
	}

	data, err := csv.Marshal(sample)
	a.NilNow(err)
	expected := "ID,Name\n1,John Doe\n"
	a.EqualNow(expected, string(data))
}

type emptyStruct struct{}

func TestEncodeEmptyStruct(t *testing.T) {
	a := assert.New(t)
	sample := emptyStruct{}

	data, err := csv.Marshal(sample)
	a.NilNow(err)
	expected := "\n\n"
	a.EqualNow(expected, string(data))
}

type IgnoreFieldStruct struct {
	ID    int    `csv:"id"`
	Name  string `csv:"-"`
	Email string `csv:"email"`
}

func TestEncodeStructWithIgnoredField(t *testing.T) {
	a := assert.New(t)
	sample := IgnoreFieldStruct{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	data, err := csv.Marshal(sample)
	a.NilNow(err)
	expected := "id,email\n1,john@example.com\n"
	a.EqualNow(string(data), expected)
}
