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
	"time"
)

type Unmarshaler interface {
	UnmarshalCSV([]byte) error
}

func Unmarshal(data []byte, v any) error {
	e := NewDecoder(bytes.NewReader(data))
	defer decoderPool.Put(e)

	return e.unmarshal(v)
}

type Decoder struct {
	reader     *csv.Reader
	lastRecord []string
	useLast    bool
	noHeader   bool
}

var decoderPool sync.Pool = sync.Pool{
	New: func() any {
		return &Decoder{}
	},
}

func NewDecoder(reader io.Reader, opts ...CSVOption) *Decoder {
	builder := newCSVBuilder(opts...)

	csvReader := csv.NewReader(reader)
	v := decoderPool.Get()
	if v == nil {
		v = &Decoder{}
	}
	d := v.(*Decoder)
	csvReader.Comma = builder.comma
	d.reader = csvReader
	d.noHeader = builder.noHeader
	return d
}

func (d *Decoder) Decode(v any) error {
	return d.unmarshal(v)
}

func (d *Decoder) unmarshal(v any) error {
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

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
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
	case reflect.Chan:
		for {
			elem := reflect.New(rv.Type().Elem()).Elem()
			ok, err := d.readRecord(meta, elem)
			if err != nil {
				return err
			}
			if !ok {
				break
			}
			rv.Send(elem)
		}
	default:
		_, err := d.readRecord(meta, rv)
		return err
	}

	return nil
}

func (d *Decoder) getMetaFields(rv reflect.Value) ([]*fieldMeta, error) {
	meta, err := reflectMetadata(rv)
	if err != nil {
		return nil, err
	}

	if d.noHeader {
		return meta, nil
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

func (d *Decoder) readLine() ([]string, error) {
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

func (d *Decoder) readRecord(meta []*fieldMeta, v reflect.Value) (bool, error) {
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
		if i >= len(meta) {
			break
		}

		m := meta[i]
		if m == nil {
			continue
		}

		fv := v.Field(m.Index)
		if err := d.marshalValue(col, fv, m); err != nil {
			return false, err
		}
	}

	return true, nil
}

var (
	unmarshalerType     = reflect.TypeFor[Unmarshaler]()
	textUnmarshalerType = reflect.TypeFor[encoding.TextUnmarshaler]()
)

func (d *Decoder) newPointerDecoder(s string, v reflect.Value, m *fieldMeta) error {
	// string pointer can be empty
	if s == "" && v.Type().Elem().Kind() != reflect.String {
		return nil
	}

	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}
	v = v.Elem()

	return d.marshalValue(s, v, m)
}

func (d *Decoder) marshalValue(col string, v reflect.Value, meta *fieldMeta) error {
	if v.CanAddr() && v.Addr().Type().Implements(unmarshalerType) {
		return unmarshalerDecoder(col, v.Addr(), meta)
	}

	t := v.Type()

	switch t.Kind() {
	case reflect.Bool:
		return boolDecoder(col, v, meta)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intDecoder(col, v, meta)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uintDecoder(col, v, meta)
	case reflect.Float32, reflect.Float64:
		return floatDecoder(col, v, meta)
	case reflect.String:
		return stringDecoder(col, v, meta)
	case reflect.Pointer:
		return d.newPointerDecoder(col, v, meta)
	case reflect.Struct:
		if t.ConvertibleTo(timeType) {
			return timeDecoder(col, v, meta)
		}
	}

	if v.CanAddr() && v.Addr().Type().Implements(textUnmarshalerType) {
		return textUnmarshalerDecoder(col, v.Addr(), meta)
	}

	return unsupportedDecoder(col, v, meta)
}

func boolDecoder(s string, v reflect.Value, _ *fieldMeta) error {
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

func intDecoder(s string, v reflect.Value, _ *fieldMeta) error {
	intVal, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	v.SetInt(intVal)
	return nil
}

func uintDecoder(s string, v reflect.Value, _ *fieldMeta) error {
	uintVal, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	v.SetUint(uintVal)
	return nil
}

func floatDecoder(s string, v reflect.Value, _ *fieldMeta) error {
	floatVal, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	v.SetFloat(floatVal)
	return nil
}

func stringDecoder(s string, v reflect.Value, _ *fieldMeta) error {
	v.SetString(s)
	return nil
}

func timeDecoder(s string, v reflect.Value, m *fieldMeta) error {
	if s == "" {
		return nil
	}

	layout := time.RFC3339Nano
	if m.Format != "" {
		layout = m.Format
	}

	tm, err := time.Parse(layout, s)
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(tm))
	return nil
}

func unmarshalerDecoder(s string, v reflect.Value, _ *fieldMeta) error {
	um, ok := v.Interface().(Unmarshaler)
	if !ok {
		return ErrUnsupportedType
	}
	return um.UnmarshalCSV([]byte(s))
}

func textUnmarshalerDecoder(s string, v reflect.Value, _ *fieldMeta) error {
	tum, ok := v.Interface().(encoding.TextUnmarshaler)
	if !ok {
		return ErrUnsupportedType
	}
	return tum.UnmarshalText([]byte(s))
}

func unsupportedDecoder(_ string, _ reflect.Value, _ *fieldMeta) error {
	return ErrUnsupportedType
}
