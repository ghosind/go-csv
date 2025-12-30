package csv

import (
	"bytes"
	"encoding"
	"encoding/csv"
	"errors"
	"io"
	"reflect"
	"strconv"
	"sync"
)

type Unmarshaler interface {
	UnmarshalCSV([]byte) error
}

func Unmarshal(data []byte, v any) error {
	e := newDecodeState(bytes.NewReader(data))
	defer decodeStatePool.Put(e)

	return e.unmarshal(v)
}

type decodeState struct {
	reader     *csv.Reader
	lastRecord []string
	useLast    bool
}

var decodeStatePool sync.Pool = sync.Pool{
	New: func() any {
		return &decodeState{}
	},
}

func newDecodeState(reader io.Reader) *decodeState {
	csvReader := csv.NewReader(reader)
	if v := decodeStatePool.Get(); v != nil {
		d := v.(*decodeState)
		d.reader = csvReader
		return d
	}
	return &decodeState{reader: csvReader}
}

func (d *decodeState) unmarshal(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return newInvalidUnmarshalError(rv)
	}

	meta, err := d.getMetaFields(rv)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		_, err := d.readRecord(meta, rv)
		return err
	}

	for i := 0; ; i++ {
		if rv.Kind() == reflect.Array && i >= rv.Len() {
			break
		}

		var elem reflect.Value
		if i < rv.Len() {
			elem = rv.Index(i)
		} else {
			elem = reflect.New(rv.Type().Elem())
		}

		ok, err := d.readRecord(meta, elem)
		if err != nil {
			return err
		}
		if !ok {
			break
		}

		if i < rv.Len() {
			rv.Index(i).Set(elem)
		} else {
			rv.Set(reflect.Append(rv, elem.Elem()))
		}
	}

	return nil
}

func (d *decodeState) getMetaFields(rv reflect.Value) ([]*fieldMeta, error) {
	meta, err := reflectMetadata(rv)
	if err != nil {
		return nil, err
	}

	header, err := d.readLine()
	if err != nil {
		return nil, err
	}

	// reorder meta according to header
	orderedMeta := make([]*fieldMeta, len(header))
	matched := false
	for i, colName := range header {
		for _, m := range meta {
			if m.Name == colName {
				orderedMeta[i] = m
				matched = true
				break
			}
		}
	}

	// skip header parse if no field matched
	if !matched {
		// mark the last record to reuse
		d.useLast = true
		return meta, nil
	}

	return orderedMeta, nil
}

func (d *decodeState) readLine() ([]string, error) {
	if d.useLast {
		d.useLast = false
		return d.lastRecord, nil
	}

	records, err := d.reader.Read()
	if err != nil {
		return nil, err
	}

	d.lastRecord = records
	return records, nil
}

func (d *decodeState) readRecord(meta []*fieldMeta, v reflect.Value) (bool, error) {
	record, err := d.readLine()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return false, nil
		}
		return false, err
	}

	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			if !v.CanSet() {
				return false, ErrCannotSet
			}
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}

	for i, col := range record {
		m := meta[i]
		if m == nil {
			continue
		}

		fv := v.Field(m.Index)
		if err := valueDecoder(m)(col, fv); err != nil {
			return false, err
		}
	}

	return true, nil
}

type decoderFunc func(string, reflect.Value) error

var decoderCache sync.Map

func valueDecoder(meta *fieldMeta) decoderFunc {
	return typeDecoder(meta.Type)
}

func typeDecoder(t reflect.Type) decoderFunc {
	if fn, ok := decoderCache.Load(t); ok {
		return fn.(decoderFunc)
	}

	f := newTypeDecoder(t)
	decoderCache.Store(t, f)
	return f
}

var (
	unmarshalerType     = reflect.TypeFor[Unmarshaler]()
	textUnmarshalerType = reflect.TypeFor[encoding.TextUnmarshaler]()
)

func newTypeDecoder(t reflect.Type) decoderFunc {
	if t.Implements(unmarshalerType) {
		return unmarshalerDecoder
	}
	if t.Implements(textUnmarshalerType) {
		return textUnmarshalerDecoder
	}

	switch t.Kind() {
	case reflect.Bool:
		return boolDecoder
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intDecoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uintDecoder
	case reflect.Float32, reflect.Float64:
		return floatDecoder
	case reflect.String:
		return stringDecoder
	default:
		return unsupportedDecoder
	}
}

func boolDecoder(s string, v reflect.Value) error {
	switch s {
	case "true", "1":
		v.SetBool(true)
	case "false", "0":
		v.SetBool(false)
	default:
		v.SetBool(false)
	}
	return nil
}

func intDecoder(s string, v reflect.Value) error {
	intVal, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	v.SetInt(intVal)
	return nil
}

func uintDecoder(s string, v reflect.Value) error {
	uintVal, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	v.SetUint(uintVal)
	return nil
}

func floatDecoder(s string, v reflect.Value) error {
	floatVal, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	v.SetFloat(floatVal)
	return nil
}

func stringDecoder(s string, v reflect.Value) error {
	v.SetString(s)
	return nil
}

func unmarshalerDecoder(s string, v reflect.Value) error {
	um := v.Interface().(Unmarshaler)
	return um.UnmarshalCSV([]byte(s))
}

func textUnmarshalerDecoder(s string, v reflect.Value) error {
	tum := v.Interface().(encoding.TextUnmarshaler)
	return tum.UnmarshalText([]byte(s))
}

func unsupportedDecoder(s string, v reflect.Value) error {
	return nil
}
