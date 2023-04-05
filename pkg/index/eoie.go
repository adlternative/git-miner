package index

import (
	"encoding/binary"
	"fmt"
)

const EOIESignature = 0x454F4945

type EOIEExtension struct {
	signature Signature
	size      uint32
	offset    uint32

	endOffIndexEntriesOffset uint32
	hash                     []byte
}

func (eoie *EOIEExtension) String() string {
	return fmt.Sprintf("[eoie] signature:%v, offset:%v size:%v, eoieOffset:%v, hash:%x", eoie.signature, eoie.offset, eoie.size, eoie.endOffIndexEntriesOffset, eoie.hash)
}

func (eoie *EOIEExtension) Signature() string {
	return eoie.signature.String()
}

func (eoie *EOIEExtension) Size() uint32 {
	return eoie.size
}

func NewEOIEExtension(buf []byte, offset uint32) (*EOIEExtension, error) {
	curOffset := offset

	if offset+8 > uint32(len(buf)) {
		return nil, fmt.Errorf("too short IEOT extention header")
	}

	signature, size, err := ParseExtensionHeader(buf[offset : offset+8])
	if err != nil {
		return nil, err
	}
	if signature != EOIESignature {
		return nil, NewInvalidSignature(EOIESignature, signature)
	}

	if size != EOIESize {
		return nil, fmt.Errorf("invalid eoie size: %d", size)
	}

	offset += 8

	//offset 4B
	endOffIndexEntriesOffset := binary.BigEndian.Uint32(buf[offset : offset+4])
	offset += 4
	//hash 20B
	hash := buf[offset : offset+20]

	return &EOIEExtension{
		signature:                signature,
		size:                     size,
		offset:                   curOffset,
		endOffIndexEntriesOffset: endOffIndexEntriesOffset,
		hash:                     hash,
	}, nil
}
