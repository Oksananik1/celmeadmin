package qst

import (
	"fmt"
	"net/url"
	"reflect"
)

// Query конвертирует объект состояния в URL query
func Query(state interface{}) (url.Values, error) {
	type encoder interface {
		Encode(string, url.Values)
	}
	switch state := state.(type) {
	case encoder:
		query := url.Values{}
		state.Encode("", query)
		return query, nil
	default:
		return reflectEncode(state)
	}
}

// State заполняет состояние из URL query
func State(state interface{}, query url.Values) error {
	type decoder interface {
		Decode(string, url.Values)
	}
	switch state := state.(type) {
	case decoder:
		state.Decode("", query)
	default:
		panic("not implemented yet")
	}
	return nil
}

// Clean удаляет параметры по-умолчанию из URL query
func Clean(query, defaults url.Values) {
	skip := make([]string, 0, len(query))
	for k, xs := range query {
		if ys, ok := defaults[k]; ok {
			if aryEq(xs, ys) {
				skip = append(skip, k)
			}
		}
	}
	for _, k := range skip {
		query.Del(k)
	}
}

func reflectEncode(state interface{}) (url.Values, error) {
	// NOTE: Функция не работает
	root := reflect.ValueOf(state)
	result := url.Values{}
	if root.Kind() == reflect.Ptr {
		root = root.Elem()
	}
	if root.Kind() != reflect.Struct {
		return nil, fmt.Errorf("unsupported kind: %v", root.Kind())
	}
	rt := root.Type()
	for i := 0; i < root.NumField(); i++ {
		field := rt.Field(i)
		qs := field.Tag.Get("qs")
		if qs == "" {
			qs = field.Name
		}
		result[qs] = nil
	}
	return result, nil
}

func aryEq(left, right []string) bool {
	for i := 0; i < len(left) && i < len(right); i++ {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}
