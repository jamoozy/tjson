package tjson

import (
	"errors"
)

type TJSONError error

var (
	ErrMissingTag   TJSONError = errors.New("Key missing a tag.")
	ErrDuplicateKey TJSONError = errors.New("Duplicate key.")

	ErrB16UpperCase TJSONError = errors.New("Base16 data contains upper case letters.")
	ErrB16Invalid   TJSONError = errors.New("Invalid base16-encoded data.")

	ErrB32UpperCase TJSONError = errors.New("Base32 data contains upper case letters.")
	ErrB32Invalid   TJSONError = errors.New("Invalid base32-encoded data.")
	ErrB32Padding   TJSONError = errors.New("Padding in base32-encoded data.")

	ErrB64UpperCase TJSONError = errors.New("Base64 data contains upper case letters.")
	ErrB64Invalid   TJSONError = errors.New("Invalid base64-encoded data.")
	ErrB64Padding   TJSONError = errors.New("Padding in base64-encoded data.")

	ErrIntOverflow  TJSONError = errors.New("Integer overflow.")
	ErrIntUnderflow TJSONError = errors.New("Integer underflow.")
	ErrIntInvalid   TJSONError = errors.New("Invalid integer.")

	ErrUintOverflow TJSONError = errors.New("Unsigned integer overflow.")
	ErrUintNegative TJSONError = errors.New("Unsigned integer negative.")
	ErrUintInvalid  TJSONError = errors.New("Invalid unsigned integer.")

	ErrTimeInvalid TJSONError = errors.New("Invalid time format.")
)
