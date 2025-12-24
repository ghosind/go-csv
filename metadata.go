package csv

import (
	"reflect"
	"strings"
	"sync"
)

type fieldMeta struct {
	Index  int
	Name   string
	Type   reflect.Type
	Format string
}

var metadataCache sync.Map

func reflectMetadata(v reflect.Value) ([]*fieldMeta, error) {
	ty, err := getValueType(v)
	if err != nil {
		return nil, err
	}

	if meta, ok := metadataCache.Load(ty); ok {
		return meta.([]*fieldMeta), nil
	}

	metas := make([]*fieldMeta, 0, ty.NumField())
	for i := 0; i < ty.NumField(); i++ {
		f := ty.Field(i)
		if f.PkgPath != "" { // unexported
			continue
		}
		tag := f.Tag.Get("csv")
		if tag == "-" {
			continue
		}

		parts := strings.Split(tag, ",")
		name := strings.TrimSpace(parts[0])
		if name == "" {
			name = f.Name
		}
		format := ""
		if len(parts) > 1 {
			for _, part := range parts[1:] {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "format=") {
					format = strings.TrimPrefix(part, "format=")
				}
			}
		}

		fm := &fieldMeta{Index: i, Name: name, Type: f.Type, Format: format}
		metas = append(metas, fm)
	}

	metadataCache.Store(ty, metas)

	return metas, nil
}

func getValueType(v reflect.Value) (reflect.Type, error) {
	t := v.Type()
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		t = t.Elem()
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, ErrInvalidType
	}

	return t, nil
}
