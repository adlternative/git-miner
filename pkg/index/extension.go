package index

import (
	"encoding/binary"
	"fmt"
)

type Extension interface {
	//Signature() string
	//Size() uint32
	String() string
}

func ParseExtensionHeader(buf []byte) (Signature, uint32, error) {
	if len(buf) != 8 {
		return 0, 0, fmt.Errorf("invalid extension header length")
	}

	return Signature(binary.BigEndian.Uint32(buf[:4])), binary.BigEndian.Uint32(buf[4:8]), nil
}
