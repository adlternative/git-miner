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

func NewEOIEExtension(offset uint32, buf []byte) (*EOIEExtension, error) {
	signature := binary.BigEndian.Uint32(buf[:4])
	if signature != EOIESignature {
		return nil, NewInvalidSignature(EOIESignature, signature)
	}
	buf = buf[4:]
	size := binary.BigEndian.Uint32(buf[:4])
	if size != EOIESize {
		return nil, fmt.Errorf("invalid eoie size: %d", size)
	}
	buf = buf[4:]
	if size != uint32(len(buf)) {
		return nil, fmt.Errorf("invalid eoie data size: %d", len(buf))
	}

	//offset 4B
	endOffIndexEntriesOffset := binary.BigEndian.Uint32(buf[:4])
	buf = buf[4:]
	//hash 20B

	return &EOIEExtension{
		signature:                Signature(signature),
		size:                     size,
		offset:                   offset,
		endOffIndexEntriesOffset: endOffIndexEntriesOffset,
		hash:                     buf,
	}, nil
}
