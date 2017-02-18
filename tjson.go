// Package tjson implements the TJSON specification defined at https://github.com/tjson/tjson-spec.
//
// Note that this package attempts to copy the interface of "encoding/json".  If an element is
// missing from this API that exists in "encoding/json", just use the "encoding/json" version.
// Since TJSON is a subset of JSON, the API was omitted because the JSON version will work.
package tjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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
	var sj interface{}
	if err := json.Unmarshal(data, &sj); err != nil {
		return err
	}

	switch v := v.(type) {
	case map[string]interface{}:
	case []map[string]interface{}:
	case *interface{}:
	default:
		return fmt.Errorf("Unrecognized type: %v.(%T)", v, v)
	}

	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.New("tjson: Non-pointer passed.")
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
