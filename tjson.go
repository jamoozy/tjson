// Package tjson implements the TJSON specification defined at https://github.com/tjson/tjson-spec.
//
// Note that this package attempts to copy the interface of "encoding/json".  If an element is
// missing from this API that exists in "encoding/json", just use the "encoding/json" version.
// Since TJSON is a subset of JSON, the API was omitted because the JSON version will work.
package tjson

import (
	"errors"
)

// Marshal converts an arbitrary interface into a byte array.
//
// See also encoding/json.Marshal
func Marshal(v interface{}) ([]byte, error) {
	return nil, errors.New("not impl")
}

// MarshalIndent marshals an arbitrary interface into an indented TJSON string returned as a byte
// slice.
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return nil, errors.New("not impl")
}

// Unmarshal unmarshals TJSON data into an interface.
func Unmarshal(data []byte, v interface{}) error {
	return errors.New("not impl")
}
