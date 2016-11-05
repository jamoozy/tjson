# TJSON

This library implements the TJSON standard as defined at https://github.com/tjson/tjson-spec.

Note that this package attempts to copy the interface of "encoding/json".
If an element is missing from this API that exists in "encoding/json", just use the "encoding/json" version.
Since TJSON is a stricter (non-strict) subset of JSON, the API was omitted because the JSON function will work.
