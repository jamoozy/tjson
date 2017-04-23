// Package tjson implements the TJSON specification defined at https://github.com/tjson/tjson-spec.
//
// Note that this package attempts to copy the interface of "encoding/json".  If an element is
// missing from this API that exists in "encoding/json", just use the "encoding/json" version.
// Since TJSON is a subset of JSON, the API was omitted because the JSON version will work.
package tjson

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const tFmt = "2006-01-02T15:04:05Z"

// Marshal converts an arbitrary interface into a byte array.
//
// See also encoding/json.Marshal
func Marshal(v interface{}) (_ []byte, err error) {
	m, ok := v.(map[string]interface{})
	if !ok {
		fmt.Errorf("Cannot handle non-maps :(")
	}

	m, err = tjsonify(m)
	if err != nil {
		return nil, err
	}

	fmt.Println(m)
	return json.Marshal(m)
}

func tjsonify(m map[string]interface{}) (_ map[string]interface{}, err error) {
	getSubTyp := func(v interface{}) string {
		switch v.(type) {
		case string:
			return "s"
		case int64:
			return "i"
		case uint64:
			return "u"
		case float64:
			return "f"
		case time.Time:
			return "t"
		case bool:
			return "v"
		case []interface{}:
			return "A"
		case map[string]interface{}:
			return "O"
		}
		panic("unrecognized type")
	}

	out := map[string]interface{}{}
	for k, v := range m {
		switch v := v.(type) {
		case string:
			out[k+":s"] = v

		// TODO b16,b32,b64?

		case int64:
			out[k+":i"] = strconv.FormatInt(v, 64)

		case uint64:
			out[k+":u"] = strconv.FormatUint(v, 64)

		case float64:
			out[k+":f"] = v

		case time.Time:
			out[k+":t"] = v.Format(tFmt)

		case bool:
			out[k+":b"] = v

		case []interface{}:
			var subTyp string
			a := make([]interface{}, 0, len(v))
			for _, val := range v {
				if subTyp == "" {
					subTyp = getSubTyp(val)
				} else if subTyp != getSubTyp(val) {
					return nil, fmt.Errorf("Heterogenous elemnets in slice.")
				}
				a = append(a, val)
			}
			out[k+":A<"+subTyp+">"] = a

		case map[string]struct{}:
			// FIXME support sets containing other than strings
			a := make([]string, 0, len(v))
			for key := range v {
				a = append(a, key)
			}
			out[k+":S<s>"] = a

		case map[string]interface{}:
			m[k+":O"], err = tjsonify(v)
			if err != nil {
				return
			}
		}
	}

	return out, nil
}

// MarshalIndent marshals an arbitrary interface into an indented TJSON string returned as a byte
// slice.
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// Unmarshal unmarshals TJSON data into an interface.
func Unmarshal(data []byte, v interface{}) error {
	// Surround the raw data s.t. the top level element is always an object.  This allows us to use
	// the JSON module to handle top-level arrays.
	data = []byte(fmt.Sprintf(`{"dat":%s}`, data))
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	fmt.Printf("From JSON: %+v\n", m)

	dat := m["dat"]
	switch dat := dat.(type) {
	case []interface{}:
		for _ = range dat {
			panic("oiy")
		}
		reflect.Indirect(reflect.ValueOf(v)).Set(reflect.ValueOf(make([]interface{}, 0)))

	case map[string]interface{}:
		if err := set("O", reflect.ValueOf(dat), reflect.ValueOf(v)); err != nil {
			return err
		}
	}

	return nil
}

func set(typ string, src, dst reflect.Value) error {
	fmt.Printf("set(%q, %v.(%v), %v.(%v)\n", typ, src, src.Kind(), dst, dst.Kind())
	if dst.Kind() != reflect.Ptr {
		return fmt.Errorf("Expected pointer")
	}
	ind := reflect.Indirect(dst)

	// Get at underlying value.
	if !src.CanInterface() {
		return fmt.Errorf("Can't interface.")
	}
	src = reflect.ValueOf(src.Interface())

	switch typ {
	case "s":
		if src.Kind() != reflect.String {
			fmt.Errorf("Expected string")
		}
		ind.Set(src)
		return nil

	case "b16":
		if src.Kind() != reflect.String {
			return fmt.Errorf("Expected string")
		}
		b, err := hex.DecodeString(src.String())
		if err != nil {
			return err
		}
		ind.Set(reflect.ValueOf(b))
		return nil

	case "b32":
		if src.Kind() != reflect.String {
			return fmt.Errorf("Expected string")
		}
		b, err := base32.StdEncoding.DecodeString(src.String())
		if err != nil {
			return err
		}
		ind.Set(reflect.ValueOf(b))
		return nil

	case "b64":
		if src.Kind() != reflect.String {
			return fmt.Errorf("Expected string")
		}
		b, err := base64.RawStdEncoding.DecodeString(src.String())
		if err != nil {
			return err
		}
		ind.Set(reflect.ValueOf(b))
		return nil

	case "i":
		if src.Kind() != reflect.String {
			return fmt.Errorf("Expected string")
		}
		i, err := strconv.ParseInt(src.String(), 10, 64)
		if err != nil {
			return err
		}
		ind.Set(reflect.ValueOf(i))
		return nil

	case "u":
		if src.Kind() != reflect.String {
			return fmt.Errorf("Expected string")
		}
		i, err := strconv.ParseUint(src.String(), 10, 64)
		if err != nil {
			return err
		}
		ind.Set(reflect.ValueOf(i))
		return nil

	case "f":
		if src.Kind() != reflect.Float64 {
			return fmt.Errorf("Expected float64")
		}
		ind.Set(src)
		return nil

	case "t":
		if src.Kind() != reflect.String {
			return fmt.Errorf("Expected string")
		}
		tim, err := time.Parse(tFmt, src.String())
		if err != nil {
			return err
		}
		ind.Set(reflect.ValueOf(tim))
		return nil

	case "v":
		if src.Kind() != reflect.Bool {
			return fmt.Errorf("Expected bool")
		}
		ind.Set(src)
		return nil

	case "O":
		if src.Kind() != reflect.Map {
			return fmt.Errorf("Expected map")
		}

		out := map[string]interface{}{}
		for _, key := range src.MapKeys() {
			parts := strings.Split(key.String(), ":")
			if len(parts) != 2 || parts[1] == "" {
				return ErrMissingTag
			}

			val := src.MapIndex(key)
			var subOut interface{}
			if err := set(parts[1], val, reflect.ValueOf(&subOut)); err != nil {
				return err
			}

			out[parts[0]] = subOut
		}

		ind.Set(reflect.ValueOf(out))
		return nil
	}

	if strings.HasPrefix(typ, "A<") {
		if !strings.HasSuffix(typ, ">") {
			return fmt.Errorf("Malformatted key: %q", typ)
		}

		iface := src.Interface()
		v, ok := iface.([]interface{})
		if !ok {
			return fmt.Errorf("Bad type: %v.(%v)\n", iface, reflect.TypeOf(iface))
		}

		src = reflect.ValueOf(v)
		if k := src.Kind(); k != reflect.Slice && k != reflect.Array {
			fmt.Printf("Got: %v.(%v)\n", src.String(), src.Kind())
			return fmt.Errorf("Expected []interface{}")
		}

		subTyp := typ[2 : len(typ)-1]

		var a []interface{}
		for i := 0; i < src.Len(); i++ {
			var next interface{}
			if err := set(subTyp, src.Index(i), reflect.ValueOf(&next)); err != nil {
				return err
			}
			a = append(a, next)
		}

		ind.Set(reflect.ValueOf(a))
		return nil
	}

	// FIXME Go only supports string keys in map[]s.  Can we extend this to arbitrary objects?
	if strings.HasPrefix(typ, "S<") {
		if !strings.HasSuffix(typ, ">") {
			return fmt.Errorf("Malformatted key: %q", typ)
		}

		if k := src.Kind(); k != reflect.Slice && k != reflect.Array {
			return fmt.Errorf("Expected []interface{}")
		}

		subTyp := typ[2 : len(typ)-2]

		var m map[string]struct{}
		for i := 0; i < src.Len(); i++ {
			var v interface{}
			if err := set(subTyp, src.Index(i), reflect.ValueOf(&v)); err != nil {
				return err
			}
			m[v.(string)] = struct{}{}
		}

		ind.Set(reflect.ValueOf(m))
		return nil
	}

	return fmt.Errorf("Unrecognized type: %q", typ)
}
