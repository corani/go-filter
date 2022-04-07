// Package filter defines a filter expression according to the CTO Circle REST API Harmonization
// section 3.3.17.2.
package filter

import (
	"fmt"
	"reflect"
)

type Filter interface {
	Apply(interface{}) (interface{}, error)
}

type filter struct {
	raw  string
	root *node
}

func Parse(raw string) (Filter, error) {
	root, err := parseExpression(raw)
	if err != nil {
		return nil, err
	}

	f := &filter{
		raw:  raw,
		root: root,
	}

	return f, nil
}

func MustParse(raw string) Filter {
	f, err := Parse(raw)
	if err != nil {
		panic(err)
	}

	return f
}

func Parsef(format string, args ...interface{}) (Filter, error) {
	return Parse(fmt.Sprintf(format, args...))
}

func MustParsef(format string, args ...interface{}) Filter {
	return MustParse(fmt.Sprintf(format, args...))
}

func (f *filter) String() string {
	if f == nil {
		return "<nil>"
	}

	return f.raw
}

func (f *filter) Apply(list interface{}) (interface{}, error) {
	if reflect.TypeOf(list).Kind() != reflect.Slice {
		return nil, fmt.Errorf("expected slice, got %v", reflect.TypeOf(list).Kind())
	}

	count := reflect.ValueOf(list).Len()
	if count == 0 {
		return list, nil
	}

	// NOTE(daniel): assume all elements are of the same type. @generics
	program, err := compile(f.root, reflect.ValueOf(list).Index(0).Type())
	if err != nil {
		return nil, err
	}

	var result []interface{}

	for i := 0; i < count; i++ {
		data := reflect.ValueOf(list).Index(i)

		ok, err := program.run(data)
		if err != nil {
			return nil, err
		}

		if ok {
			result = append(result, data.Interface())
		}
	}

	return result, nil
}
