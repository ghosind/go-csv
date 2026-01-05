package csv

import (
	"bytes"
	"encoding"
	"encoding/csv"
	"io"
	"reflect"
	"strconv"
	"sync"
	"time"
)

// Marshaler is the interface implemented by types that can marshal a CSV
// record representation of themselves.
type Marshaler interface {
	MarshalCSV() ([]byte, error)
}

type Encoder struct {
	writer   *csv.Writer
	noHeader bool
}

var encoderPool sync.Pool = sync.Pool{
	New: func() any {
		return &Encoder{}
	},
}

// Marshal returns the CSV encoding of v.
func Marshal(v any) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	e := NewEncoder(buf)
	defer encoderPool.Put(e)

	if err := e.marshal(v); err != nil {
		return nil, err
	}
	e.writer.Flush()

	return buf.Bytes(), nil
}

func NewEncoder(writer io.Writer, opts ...CSVOption) *Encoder {
	builder := newCSVBuilder(opts...)

	csvWriter := csv.NewWriter(writer)
	csvWriter.Comma = builder.comma
	csvWriter.UseCRLF = builder.useCRLF

	v := encoderPool.Get()
	if v == nil {
		v = &Encoder{}
	}
	e := v.(*Encoder)
	e.writer = csvWriter
	e.noHeader = builder.noHeader
	return e
}

func (e *Encoder) Encode(v any) error {
	defer e.writer.Flush()

	return e.marshal(v)
}

func (e *Encoder) marshal(v any) error {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return nil
	}

	meta, err := reflectMetadata(rv)
	if err != nil {
		return err
	}

	// write header
	if !e.noHeader {
		header := make([]string, len(meta))
		for i, m := range meta {
			header[i] = m.Name
		}
		if err := e.writer.Write(header); err != nil {
			return err
		}
	}

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		if rv.Len() == 0 {
			return nil
		}

		for i := 0; i < rv.Len(); i++ {
			if err := e.writeRow(rv.Index(i), meta); err != nil {
				return err
			}
		}
	case reflect.Chan:
		for {
			elem, ok := rv.Recv()
			if !ok {
				break
			}
			if err := e.writeRow(elem, meta); err != nil {
				return err
			}
		}
	default:
		return e.writeRow(rv, meta)
	}

	return nil
}

func (e *Encoder) writeRow(v reflect.Value, meta []*fieldMeta) error {
	row := make([]string, len(meta))

	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return e.writer.Write(row)
		}
		v = v.Elem()
	}

	for i, m := range meta {
		fv := v.Field(m.Index)
		str, err := valueEncoder(m)(fv, m)
		if err != nil {
			return err
		}
		row[i] = str
	}

	return e.writer.Write(row)
}

type encoderFunc func(reflect.Value, *fieldMeta) (string, error)

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

func (pe ptrEncoder) encode(v reflect.Value, m *fieldMeta) (string, error) {
	if v.IsNil() {
		return "", nil
	}

	return pe.elemEnc(v.Elem(), m)
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
	timeType          = reflect.TypeFor[time.Time]()
)

func newTypeEncoder(t reflect.Type) encoderFunc {
	if t.Kind() != reflect.Ptr && reflect.PointerTo(t).Implements(marshalerType) {
		return marshalerEncoder
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
	case reflect.Struct:
		if t.ConvertibleTo(timeType) {
			return timeEncoder
		}
	}

	if reflect.PointerTo(t).Implements(textMarshalerType) {
		return textMarshalerEncoder
	}

	return unsupportedTypeEncoder
}

func boolEncoder(v reflect.Value, _ *fieldMeta) (string, error) {
	s := strconv.FormatBool(v.Bool())
	return s, nil
}

func intEncoder(v reflect.Value, _ *fieldMeta) (string, error) {
	s := strconv.FormatInt(v.Int(), 10)
	return s, nil
}

func uintEncoder(v reflect.Value, _ *fieldMeta) (string, error) {
	s := strconv.FormatUint(v.Uint(), 10)
	return s, nil
}

func floatEncoder(v reflect.Value, _ *fieldMeta) (string, error) {
	s := strconv.FormatFloat(v.Float(), 'f', -1, 64)
	return s, nil
}

func stringEncoder(v reflect.Value, _ *fieldMeta) (string, error) {
	return v.String(), nil
}

func timeEncoder(v reflect.Value, m *fieldMeta) (string, error) {
	tm := v.Interface().(time.Time)

	if m.Format != "" {
		return tm.Format(m.Format), nil
	}

	// fallback to TextMarshalerEncoder
	return textMarshalerEncoder(v, m)
}

func marshalerEncoder(v reflect.Value, _ *fieldMeta) (string, error) {
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

func textMarshalerEncoder(v reflect.Value, _ *fieldMeta) (string, error) {
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

func unsupportedTypeEncoder(_ reflect.Value, _ *fieldMeta) (string, error) {
	return "", ErrUnsupportedType
}
