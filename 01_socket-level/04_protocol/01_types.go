package protocol

import (
	"errors"
	"fmt"
	"io"
)

const (
	BinaryType uint8 = iota + 1
	StringType
	MaxPayloadSize uint32 = 10 << 20 // 10 MB
)

var ErrMaxPayloadSize = errors.New("maximum payload size exceeded")

type Payload interface {
	fmt.Stringer   // Must have a String() string method
	io.ReaderFrom  // Must have a ReadFrom(r io.Reader) (n int64, err error)
	io.WriterTo    // Must have a WriteTo(w io.Writer) (n int64, err error)
	Bytes() []byte // Must have a Bytes() []byte method to retrieve raw bytes
}

type Binary []byte

type String string
