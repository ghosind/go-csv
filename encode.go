package csv

import (
	"bytes"
	"encoding"
	"encoding/csv"
	"io"
	"reflect"
	"strconv"
	"sync"
)

type Marshaler interface {
	MarshalCSV() ([]byte, error)
}

func Marshal(v any) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	e := newEncodeState(buf)
	defer encodeStatePool.Put(e)

	if err := e.marshal(v); err != nil {
		return nil, err
	}
	e.writer.Flush()

	return buf.Bytes(), nil
}

type encodeState struct {
	writer *csv.Writer
}

var encodeStatePool sync.Pool = sync.Pool{
	New: func() any {
		return &encodeState{}
	},
}

func newEncodeState(writer io.Writer) *encodeState {
	csvWriter := csv.NewWriter(writer)
	if v := encodeStatePool.Get(); v != nil {
		e := v.(*encodeState)
		e.writer = csvWriter
		return e
	}
	return &encodeState{writer: csvWriter}
}

func (e *encodeState) marshal(v any) (err error) {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return nil
	}

	meta, err := reflectMetadata(rv)
	if err != nil {
		return err
	}

	// write header
	header := make([]string, len(meta))
	for i, m := range meta {
		header[i] = m.Name
	}
	if err := e.writer.Write(header); err != nil {
		return err
	}

	// write rows
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return e.writeRow(rv, meta)
	}

	for i := 0; i < rv.Len(); i++ {
		if err := e.writeRow(rv.Index(i), meta); err != nil {
			return err
		}
	}

	return nil
}

func (e *encodeState) writeRow(v reflect.Value, meta []*fieldMeta) error {
	row := make([]string, len(meta))

	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return e.writer.Write(row)
		}
		v = v.Elem()
	}

	for i, m := range meta {
		fv := v.Field(m.Index)
		str, err := valueEncoder(m)(fv)
		if err != nil {
			return err
		}
		row[i] = str
	}

	return e.writer.Write(row)
}

type encoderFunc func(reflect.Value) (string, error)

var encoderCache sync.Map

func valueEncoder(meta *fieldMeta) encoderFunc {
	return typeEncoder(meta.Type)
}

type ptrEncoder struct {
	elemEnc encoderFunc
}

func newPtrEncoder(t reflect.Type) encoderFunc {
	enc := ptrEncoder{typeEncoder(t.Elem())}
	return enc.encode
}

func (pe ptrEncoder) encode(v reflect.Value) (string, error) {
	if v.IsNil() {
		return "", nil
	}

	return pe.elemEnc(v.Elem())
}

func typeEncoder(t reflect.Type) encoderFunc {
	if fi, ok := encoderCache.Load(t); ok {
		return fi.(encoderFunc)
	}

	f := newTypeEncoder(t)
	encoderCache.Store(t, f)
	return f
}

var (
	marshalerType     = reflect.TypeFor[Marshaler]()
	textMarshalerType = reflect.TypeFor[encoding.TextMarshaler]()
)

func newTypeEncoder(t reflect.Type) encoderFunc {
	if t.Implements(marshalerType) {
		return marshalerEncoder
	}
	if t.Implements(textMarshalerType) {
		return textMarshalerEncoder
	}

	switch t.Kind() {
	case reflect.Bool:
		return boolEncoder
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uintEncoder
	case reflect.Float32, reflect.Float64:
		return floatEncoder
	case reflect.String:
		return stringEncoder
	case reflect.Ptr:
		return newPtrEncoder(t)
	default:
		return unsupportedTypeEncoder
	}
}

func boolEncoder(v reflect.Value) (string, error) {
	s := strconv.FormatBool(v.Bool())
	return s, nil
}

func intEncoder(v reflect.Value) (string, error) {
	s := strconv.FormatInt(v.Int(), 10)
	return s, nil
}

func uintEncoder(v reflect.Value) (string, error) {
	s := strconv.FormatUint(v.Uint(), 10)
	return s, nil
}

func floatEncoder(v reflect.Value) (string, error) {
	s := strconv.FormatFloat(v.Float(), 'f', -1, 64)
	return s, nil
}

func stringEncoder(v reflect.Value) (string, error) {
	return v.String(), nil
}

func marshalerEncoder(v reflect.Value) (string, error) {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return "", nil
	}
	m, ok := v.Interface().(Marshaler)
	if !ok {
		return "", ErrUnsupportedType
	}
	b, err := m.MarshalCSV()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func textMarshalerEncoder(v reflect.Value) (string, error) {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return "", nil
	}
	m, ok := v.Interface().(encoding.TextMarshaler)
	if !ok {
		return "", ErrUnsupportedType
	}
	b, err := m.MarshalText()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func unsupportedTypeEncoder(v reflect.Value) (string, error) {
	return "", ErrUnsupportedType
}
