package index

import (
	"fmt"
)

const IEOTSignature = 0x49454F54

type IEOTExtension struct {
	signature Signature
	size      uint32
	offset    uint32

	//endOffIndexEntriesOffset uint32
	//hash                     []byte
}

func (e *IEOTExtension) String() string {
	return fmt.Sprintf("[ieot] signature:%v, offset:%v size:%v", e.signature, e.offset, e.size)
}

func NewIEOTExtension(buf []byte, offset uint32, size uint32) (*IEOTExtension, error) {
	return &IEOTExtension{
		signature: IEOTSignature,
		size:      size,
		offset:    offset,
	}, nil
}
