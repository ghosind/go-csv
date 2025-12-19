package csv

import (
	"reflect"
	"sync"
)

type fieldMeta struct {
	Index int
	Name  string
	Type  reflect.Type
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

		name := tag
		if name == "" {
			name = f.Name
		}
		fm := &fieldMeta{Index: i, Name: name, Type: f.Type}
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
