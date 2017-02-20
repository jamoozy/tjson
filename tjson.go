// Package tjson implements the TJSON specification defined at https://github.com/tjson/tjson-spec.
//
// Note that this package attempts to copy the interface of "encoding/json".  If an element is
// missing from this API that exists in "encoding/json", just use the "encoding/json" version.
// Since TJSON is a subset of JSON, the API was omitted because the JSON version will work.
package tjson

import (
	"encoding/json"
	"reflect"
	"fmt"
)

// Marshal converts an arbitrary interface into a byte array.
//
// See also encoding/json.Marshal
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
	//return nil, errors.New("not impl")
}

// MarshalIndent marshals an arbitrary interface into an indented TJSON string returned as a byte
// slice.
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
	//return nil, errors.New("not impl")
}

// Unmarshal unmarshals TJSON data into an interface.
func Unmarshal(data []byte, v interface{}) error {
	vv := reflect.ValueOf(v)
	if vv.Kind() != reflect.Ptr || vv.IsNil() {
		return ErrInvalidArg
	}

	var u interface{}
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}

	return unpack(vv, reflect.ValueOf(&u))
}

// Unpacks the contents of `src' into `dst'.  `src' should be an interface{} value processed by the
// "json" module.  `dst' is the result after interpreting the contents of `src' in TJSON. `k' is the
// expected Kind (type) of the dst value.  If it isn't that type or
func unpack(dst, src reflect.Value) error {
	if dst.Kind() != reflect.Ptr || dst.IsNil() {
		return ErrInvalidArg
	}

	// "Drill down" to the actual source element we want to clone.
	for src.Kind() == reflect.Ptr || src.Kind() == reflect.Interface {
		src = src.Elem()
	}

	// For arrays, call this function recursively for each sub-element.
	switch src.Kind() {
	case reflect.Array, reflect.Slice:
		l := src.Len()
		dst.Elem().Set(reflect.ValueOf(make([]interface{}, l)))
		for i := 0; i < l; i++ {
			var subDstVal interface{}
			if err := unpack(reflect.ValueOf(&subDstVal), src.Index(i)); err != nil {
				return err
			}
			reflect.Append(dst.Elem(), reflect.ValueOf(subDstVal))
		}

	case reflect.Map:
		dst.Elem().Set(reflect.ValueOf(map[string]interface{}{}))
		for _, k := range src.MapKeys() {
			var subDstVal interface{}
			if err := unpack(reflect.ValueOf(&subDstVal), src.MapIndex(k)); err != nil {
				return err
			}
			dst.Elem().Elem().SetMapIndex(k, reflect.ValueOf(subDstVal))
		}

	case reflect.String:
		dst.Elem().Set(src)

	default:
		return fmt.Errorf("Unrecognized type: %+v", src.Kind())
	}

	return nil
}

//// simpleJSON is an abstraction used to
//type simpleJSON interface{}
//
//// Iterates through the keys and values in the
//func (sj *simpleJSON) convertTo(v interface{}) error {
//	if reflect.TypeOf(v).Kind() != reflect.Ptr {
//		return errors.New("tjson: Non-pointer passed.")
//	}
//	v = sj
//	return nil
//}
