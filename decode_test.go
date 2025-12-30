package csv_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ghosind/go-assert"
	"github.com/ghosind/go-csv"
)

func TestDecodeEmpty(t *testing.T) {
	a := assert.New(t)
	data := ""
	var sample *SampleStruct
	err := csv.Unmarshal([]byte(data), sample)
	a.NilNow(err)
	a.NilNow(sample)
}

func TestDecodeStructToNonPointer(t *testing.T) {
	a := assert.New(t)
	data := "id,name,age,salary,is_manager\n1,John Doe,30,5500,true\n"
	var sample SampleStruct
	err := csv.Unmarshal([]byte(data), sample)
	a.NotNilNow(err)
}

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

func TestDecodeStructWithErrorValue(t *testing.T) {
	a := assert.New(t)
	data := "id,name,age,salary,is_manager\n1,John Doe,thirty,5500,true\n"
	var sample SampleStruct
	err := csv.Unmarshal([]byte(data), &sample)
	a.NotNilNow(err)
}

func TestDecodeStructWithInvalidBoolValue(t *testing.T) {
	a := assert.New(t)
	data := "id,name,age,salary,is_manager\n1,John Doe,30,5500,unknown\n"
	var sample SampleStruct
	err := csv.Unmarshal([]byte(data), &sample)
	a.NilNow(err)
	expected := SampleStruct{
		ID:        1,
		Name:      "John Doe",
		Age:       30,
		Salary:    5500,
		IsManager: false,
	}
	a.EqualNow(expected, sample)
}

func TestDecodeStructWithoutHeader(t *testing.T) {
	a := assert.New(t)
	data := "1,John Doe,30,5500,true\n"
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

func TestDecodeStructWithUnmatchedFields(t *testing.T) {
	a := assert.New(t)
	data := "id,name,age,salary,email\n1,John Doe,30,5500,john@example.com\n"
	var sample SampleStruct
	err := csv.Unmarshal([]byte(data), &sample)
	a.NilNow(err)
	expected := SampleStruct{
		ID:        1,
		Name:      "John Doe",
		Age:       30,
		Salary:    5500,
		IsManager: false,
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

func TestEmptyCSVToStructPointer(t *testing.T) {
	a := assert.New(t)
	data := "id,name,age,salary,is_manager\n"
	var sample *SampleStruct
	err := csv.Unmarshal([]byte(data), &sample)
	a.NilNow(err)
	a.NilNow(sample)
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

func TestDecodeStructPointerSlice(t *testing.T) {
	a := assert.New(t)
	data := "id,name,age,salary,is_manager\n1,John Doe,30,5500,true\n2,Jane Smith,25,3000,false\n"
	var samples []*SampleStruct
	err := csv.Unmarshal([]byte(data), &samples)
	a.NilNow(err)
	expected := []*SampleStruct{
		{ID: 1, Name: "John Doe", Age: 30, Salary: 5500, IsManager: true},
		{ID: 2, Name: "Jane Smith", Age: 25, Salary: 3000, IsManager: false},
	}
	a.DeepEqualNow(expected, samples)
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

func TestDecodePointerFieldsStruct(t *testing.T) {
	a := assert.New(t)
	data := "id,name,age,salary,is_manager\n1,John Doe,30,5500,true\n"
	var sample SimplePointerStruct
	err := csv.Unmarshal([]byte(data), &sample)
	a.NilNow(err)
	id, name, age, salary, isManager := 1, "John Doe", uint(30), 5500.0, true
	expected := SimplePointerStruct{
		ID:        &id,
		Name:      &name,
		Age:       &age,
		Salary:    &salary,
		IsManager: &isManager,
	}
	a.DeepEqualNow(expected, sample)
}

func TestDecodePointerFieldsStructWithEmptyField(t *testing.T) {
	a := assert.New(t)
	data := "id,name,age,salary,is_manager\n,,,,\n"
	var sample SimplePointerStruct
	err := csv.Unmarshal([]byte(data), &sample)
	a.NilNow(err)
	name := ""
	expected := SimplePointerStruct{
		ID:        nil,
		Name:      &name,
		Age:       nil,
		Salary:    nil,
		IsManager: nil,
	}
	a.DeepEqualNow(expected, sample)
}

func (m *MarshalableStruct) UnmarshalCSV(b []byte) error {
	_, err := fmt.Sscanf(string(b), "%s (%d)", &m.Country, &m.ZipCode)
	return err
}

func TestDecodeStructWithUnmarshalableField(t *testing.T) {
	a := assert.New(t)
	var sample WrapMarshalableStruct
	data := "city,location\nNew York,USA (10001)\n"
	err := csv.Unmarshal([]byte(data), &sample)
	a.NilNow(err)
	expected := WrapMarshalableStruct{
		City: "New York",
		Location: MarshalableStruct{
			Country: "USA",
			ZipCode: 10001,
		},
	}
	a.EqualNow(sample, expected)
}

func TestDecodeStructWithUnmarshalablePointerField(t *testing.T) {
	a := assert.New(t)
	var sample WrapMarshalablePointerStruct
	data := "city,location\nNew York,USA (10001)\n"
	err := csv.Unmarshal([]byte(data), &sample)
	a.NilNow(err)
	expected := WrapMarshalablePointerStruct{
		City: "New York",
		Location: &MarshalableStruct{
			Country: "USA",
			ZipCode: 10001,
		},
	}
	a.DeepEqualNow(sample, expected)
}

func (m *TextMarshalStruct) UnmarshalText(b []byte) error {
	_, err := fmt.Sscanf(string(b), "%s [%d]", &m.Country, &m.ZipCode)
	return err
}

func TestDecodeStructWithTextUnmarshalField(t *testing.T) {
	a := assert.New(t)
	data := "city,location\nNew York,USA [10001]\n"
	var sample WrapTextMarshalStruct

	err := csv.Unmarshal([]byte(data), &sample)
	a.NilNow(err)
	expected := WrapTextMarshalStruct{
		City:     "New York",
		Location: TextMarshalStruct{Country: "USA", ZipCode: 10001},
	}
	a.EqualNow(sample, expected)
}

func TestDecodeStructWithTextUnmarshalPointerField(t *testing.T) {
	a := assert.New(t)
	data := "city,location\nNew York,USA [10001]\n"
	var sample WrapTextMarshalPointerStruct

	err := csv.Unmarshal([]byte(data), &sample)
	a.NilNow(err)
	expected := WrapTextMarshalPointerStruct{
		City:     "New York",
		Location: &TextMarshalStruct{Country: "USA", ZipCode: 10001},
	}
	a.DeepEqualNow(sample, expected)
}

func TestDecodeStructWithUnexportedField(t *testing.T) {
	a := assert.New(t)
	data := "id,email\n1,john@example.com\n"
	var sample UnexportedFieldStruct

	err := csv.Unmarshal([]byte(data), &sample)
	a.NilNow(err)
	expected := UnexportedFieldStruct{
		ID:    1,
		Email: "john@example.com",
	}

	a.EqualNow(sample, expected)
}

func TestDecodeStructWithUnsupportedField(t *testing.T) {
	a := assert.New(t)
	data := "Data,b\n1,b\n"
	var sample UnsupportedStruct

	err := csv.Unmarshal([]byte(data), &sample)
	a.NotNilNow(err)
	a.IsErrorNow(err, csv.ErrUnsupportedType)
}

func TestDecodeStructWithoutTags(t *testing.T) {
	a := assert.New(t)
	data := "ID,Name\n1,John Doe\n"
	var sample NoTagStruct

	err := csv.Unmarshal([]byte(data), &sample)
	a.NilNow(err)
	expected := NoTagStruct{
		ID:   1,
		Name: "John Doe",
	}
	a.EqualNow(sample, expected)
}

func TestDecodeStructWithIgnoredField(t *testing.T) {
	a := assert.New(t)
	data := "id,name,email\n1,John Doe,john@example.com\n"
	var sample IgnoreFieldStruct

	err := csv.Unmarshal([]byte(data), &sample)
	a.NilNow(err)
	expected := IgnoreFieldStruct{
		ID:    1,
		Email: "john@example.com",
	}
	a.EqualNow(sample, expected)
}

func TestDecodeTimeStruct(t *testing.T) {
	a := assert.New(t)
	tm := time.Date(2025, 10, 1, 11, 30, 00, 0, time.UTC)
	data := "no_fmt_time,fmt_time,no_fmt_time_ptr,fmt_time_ptr\n2025-10-01T11:30:00Z,2025-10-01T11:30:00,2025-10-01T11:30:00Z,2025-10-01T11:30:00\n"
	var sample TimeStruct

	err := csv.Unmarshal([]byte(data), &sample)
	a.NilNow(err)
	expected := TimeStruct{
		NoFmtTime:    tm,
		FmtTime:      tm,
		NoFmtTimePtr: &tm,
		FmtTimePtr:   &tm,
	}
	a.DeepEqualNow(sample, expected)
}
