package index

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type Signature uint32

func (s Signature) String() string {
	buf := make([]byte, 5)
	binary.LittleEndian.PutUint32(buf, uint32(s))
	return string(buf)
}

var (
	ErrWrongSignature = errors.New("invalid signature")
)

// InvalidSignature is an error type that indicates a signature mismatch
type InvalidSignature struct {
	err    error
	expect uint32
	actual uint32
}

// Error implements the error interface for InvalidSignature
func (e *InvalidSignature) Error() string {
	return fmt.Sprintf("%e expect signature: %v, actual signature: %v", ErrWrongSignature, e.expect, e.actual)
}

func (e *InvalidSignature) Unwrap() error {
	return e.err
}

// NewInvalidSignature creates a new InvalidSignature error with a given message
func NewInvalidSignature(expect, actual uint32) *InvalidSignature {
	return &InvalidSignature{
		err:    ErrWrongSignature,
		expect: expect,
		actual: actual,
	}
}
