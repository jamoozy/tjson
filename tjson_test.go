package tjson

import (
	"testing"

	"reflect"
	"time"
)

// TJSONTestCase represents a test case pulled from:
//   https://github.com/tjson/tjson-spec/blob/master/draft-tjson-examples.txt
type TJSONTestCase struct {
	// Name and description of this test case.
	name, description string

	// Input string to test.
	input string

	// Either an error or the expected result.
	err error
	res interface{}
}

func (tc *TJSONTestCase) Run(t *testing.T) {
	t.Logf("Running test case: %q", tc.name)
	t.Logf("%s", tc.description)

	// Unmarshal the TJSON and verify that either its error or its result is the expected one.
	var res interface{}
	if err := Unmarshal([]byte(tc.input), &res); err != tc.err {
		t.Fatalf("Errors differ: %v != %v", tc.err, err)
	} else if !reflect.DeepEqual(res, tc.res) {
		t.Fatalf("Expected %#v, got %#v", tc.res, res)
	}

	// Re-marshal the result to verify that the the reverse works.
	//
	// Ideally we'd compare the resultant TJSON with our tc.input var.  TJSON give no guarantee of
	// order though, so we can't :-(
	if dat, err := Marshal(tc.res); err != nil {
		t.Fatalf("Could not re-marshal: %s", err.Error())
	} else if len(dat) <= 0 {
		t.Fatalf("Returned empty data.")
	} else if input := string(dat); len(input) != len(tc.input) {
		t.Fatalf("Returned different strings:\n  %s\n  %s", tc.input, input)
	}

	t.Logf("Success!")
}

func TestAll(t *testing.T) {
	testcases := []TJSONTestCase{
		{
			name:        "Empty Array",
			description: "Arrays are allowed as a toplevel value and can be empty",

			input: `[]`,
			res:   []interface{}{},
		},

		{
			name:        "Empty Object",
			description: "Objects are allowed as a toplevel value and can be empty",

			input: `{}`,
			res:   map[string]interface{}{},
		},

		{
			name:        "Object with UTF-8 String Key",
			description: "Strings are allowed as names of object members",

			input: `{"example:s":"foobar"}`,
			res: map[string]interface{}{
				"example": "foobar",
			},
		},

		{
			name:        "Invalid Object with Untagged Name",
			description: "All strings in TJSON must be tagged",

			input: `{"example":"foobar"}`,
			err:   ErrMissingTag,
		},

		{
			name:        "Invalid Object with Empty Tag",
			description: "All strings in TJSON must be tagged",

			input: `{"example:":"foobar"}`,
			err:   ErrMissingTag,
		},

		{
			name:        "Invalid Object with Repeated Member Names",
			description: "Names of the members of objects must be distinct",

			input: `{"example:i":"1","example:i":"2"}`,
			err:   ErrDuplicateKey,
		},

		{
			name:        "Array of integers",
			description: "Arrays are parameterized by the types of their contents",

			input: `{"example:A<i>": ["1", "2", "3"]}`,
			res: map[string]interface{}{
				"example": []int64{1, 2, 3},
			},
		},

		{
			name:        "Array of objects",
			description: "Objects are the only allowed terminal non-scalar",

			input: `{"example:A<O>": [{"a:i": "1"}, {"b:i": "2"}]}`,
			res: map[string]interface{}{
				"example": []map[string]interface{}{
					{"a": int64(1)},
					{"b": int64(2)},
				},
			},
		},

		{
			name:        "Multidimensional array of integers",
			description: "Arrays can contain other arrays",

			input: `{"example:A<A<i>>": [["1", "2"], ["3", "4"], ["5", "6"]]}`,
			res: map[string]interface{}{
				"example": [][]int64{
					{int64(1), int64(2)},
					{int64(3), int64(4)},
					{int64(5), int64(6)},
				},
			},
		},

		{
			name:        "Base16 Binary Data",
			description: "Base16 data begins with the 'b16:' prefix",

			input: `{"example:b16":"48656c6c6f2c20776f726c6421"}`,
			res: map[string]interface{}{
				"example": []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x21},
			},
		},

		{
			name:        "Invalid Base16 Binary Data with bad case",
			description: "Base16 data MUST be lower case",

			input: `{"example:b16":"48656C6C6F2C20776F726C6421"}`,
			err:   ErrB16UpperCase,
		},

		{
			name:        "Invalid Base16 Binary Data",
			description: "Base16 data must be valid hexadecimal",

			input: `{"example:b16":"This is not a valid hexadecimal string"}`,
			err:   ErrB16Invalid,
		},

		{
			name:        "Base32 Binary Data",
			description: "Base32 data begins with the 'b32:' prefix",

			input: `{"example:b32":"jbswy3dpfqqho33snrscc"}`,
			res: map[string]interface{}{
				"example": []byte("wat" /*TODO*/),
			},
		},

		{
			name:        "Invalid Base32 Binary Data with bad case",
			description: "Base32 data MUST be lower case",

			input: `{"example:b32":"JBSWY3DPFQQHO33SNRSCC"}`,
			err:   ErrB32UpperCase,
		},

		{
			name:        "Invalid Base32 Binary Data with padding",
			description: "Base64 data MUST NOT include padding",

			input: `{"example:b32":"jbswy3dpfqqho33snrscc==="}`,
			err:   ErrB64Padding,
		},

		{
			name:        "Invalid Base32 Binary Data",
			description: "Base32 data must be valid",

			input: `{"example:b32":"This is not a valid base32 string"}`,
			err:   ErrB32Invalid,
		},

		{
			name:        "Base64url Binary Data",
			description: "Base64 data begins with the 'b64:' prefix",

			input: `{"example:b64":"SGVsbG8sIHdvcmxkIQ"}`,
			res: map[string]interface{}{
				"example": []byte("wat" /*TODO*/),
			},
		},

		{
			name:        "Invalid Base64url Binary Data with padding",
			description: "Base64 data MUST NOT include padding",

			input: `{"example:b64":"SGVsbG8sIHdvcmxkIQ=="}`,
			err:   ErrB64Padding,
		},

		{
			name:        "Invalid Base64url Binary Data with non-URL safe characters",
			description: "Traditional Base64 is expressly disallowed",

			input: `{"example:b64":"+/+/"}`,
			err:   ErrB64Invalid,
		},

		{
			name:        "Invalid Base64url Binary Data",
			description: "Garbage data MUST be rejected",

			input: `{"example:b64":"This is not a valid base64url string"}`,
			err:   ErrB64Invalid,
		},

		{
			name:        "Signed Integer",
			description: "Signed integers are represented as tagged strings",

			input: `{"example:i":"42"}`,
			res: map[string]interface{}{
				"example": int64(42),
			},
		},

		{
			name:        "Signed Integer Range Test",
			description: "It should be possible to represent the full range of a signed 64-bit integer",

			input: `{"min:i":"-9223372036854775808", "max:i:":"9223372036854775807"}`,
			res: map[string]interface{}{
				"min": -9223372036854775808,
				"max": 9223372036854775807,
			},
		},

		{
			name:        "Oversized Signed Integer Test",
			description: "Values larger than can be represented by a signed 64-bit integer MUST be rejected",

			input: `{"oversize:i":"9223372036854775808"}`,
			err:   ErrIntOverflow,
		},

		{
			name:        "Undersized Signed Integer Test",
			description: "Values smaller than can be represented by a signed 64-bit integer MUST be rejected",

			input: `{"undersize:i":"-9223372036854775809"}`,
			err:   ErrIntUnderflow,
		},

		{
			name:        "Invalid Signed Integer",
			description: "Garbage data after the integer tag should be rejected",

			input: `{"invalid:i":"This is not a valid integer"}`,
			err:   ErrIntInvalid,
		},

		{
			name:        "Unsigned Integer",
			description: "Unsigned integers are represented as tagged strings",

			input: `{"example:u":"42"}`,
			res: map[string]interface{}{
				"example": uint64(42),
			},
		},

		{
			name:        "Unsigned Integer Range Test",
			description: "It should be possible to represent the full range of a signed 64-bit integer",

			input: `{"maxint:u":"18446744073709551615"}`,
			res: map[string]interface{}{
				"maxint": uint64(18446744073709551615),
			},
		},

		{
			name:        "Oversized Unsigned Integer Test",
			description: "Values larger than can be represented by an unsigned 64-bit integer MUST be rejected",

			input: `{"oversized:u":"18446744073709551616"}`,
			err:   ErrUintOverflow,
		},

		{
			name:        "Negative Unsigned Integer Test",
			description: "Unsigned integers cannot be negative",

			input: `{"negative:u":"-1"}`,
			err:   ErrUintNegative,
		},

		{
			name:        "Invalid Unsigned Integer",
			description: "Garbage data after the integer tag should be rejected",

			input: `{"invalid:u":"This is not a valid integer"}`,
			err:   ErrUintInvalid,
		},

		{
			name:        "Timestamp",
			description: "A valid RFC3339 timestamp example",

			input: `{"example:t":"2016-10-02T07:31:51Z"}`,
			res: map[string]interface{}{
				"example": time.Date(2016, time.October, 2, 7, 31, 51, 0, time.UTC),
			},
		},

		{
			name:        "Timestamp With Invalid Time Zone",
			description: "All timestamps must be in the UTC ('Z') time zone",

			input: `{"invalid:t":"2016-10-02T07:31:51-08:00"}`,
			err:   ErrTimeInvalid,
		},

		{
			name:        "Invalid Timestamp",
			description: "Garbage data after the timestamp tag should be rejected",

			input: `{"invalid:t":"This is not a valid timestamp"}`,
			err:   ErrTimeInvalid,
		},

		// --- custom tests ---
		// The following tests are Go-specific.  They're here to maintain the common paradigms that were
		// established in the "encoding/json" module, e.g., `json:"name"` tag handling.
		{
			name:        "Pointer to JSON-tagged struct",
			description: "Should be able to handle `json` tags.",

			input: `{"example:t","2016-10-02T07:31:51Z"}`,
			res: struct {
				Example time.Time `json:"example"`
			}{
				Example: time.Date(2016, time.October, 2, 7, 31, 51, 0, time.UTC),
			},
		},

		{
			name:        "Pointer to TJSON-tagged struct",
			description: "Should be able to handle `tjson` tags.",

			input: `{"example:t","2016-10-02T07:31:51Z"}`,
			res: struct {
				Example time.Time `tjson:"example"`
			}{
				Example: time.Date(2016, time.October, 2, 7, 31, 51, 0, time.UTC),
			},
		},
	}

	for i := range testcases {
		testcases[i].Run(t)
	}
}
