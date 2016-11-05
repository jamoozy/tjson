package tjson

import (
	"errors"
)

// Error is a TJSON-specific error that Marshal can return.
type Error error

// These are all the errors that malformatted TJSON can result in.  Note that JSON errors can also
// be returned by the functions in this library.
var (
	ErrMissingTag   Error = errors.New("Key missing a tag.")
	ErrDuplicateKey Error = errors.New("Duplicate key.")

	ErrB16UpperCase Error = errors.New("Base16 data contains upper case letters.")
	ErrB16Invalid   Error = errors.New("Invalid base16-encoded data.")

	ErrB32UpperCase Error = errors.New("Base32 data contains upper case letters.")
	ErrB32Invalid   Error = errors.New("Invalid base32-encoded data.")
	ErrB32Padding   Error = errors.New("Padding in base32-encoded data.")

	ErrB64UpperCase Error = errors.New("Base64 data contains upper case letters.")
	ErrB64Invalid   Error = errors.New("Invalid base64-encoded data.")
	ErrB64Padding   Error = errors.New("Padding in base64-encoded data.")

	ErrIntOverflow  Error = errors.New("Integer overflow.")
	ErrIntUnderflow Error = errors.New("Integer underflow.")
	ErrIntInvalid   Error = errors.New("Invalid integer.")

	ErrUintOverflow Error = errors.New("Unsigned integer overflow.")
	ErrUintNegative Error = errors.New("Unsigned integer negative.")
	ErrUintInvalid  Error = errors.New("Invalid unsigned integer.")

	ErrTimeInvalid Error = errors.New("Invalid time format.")
)
