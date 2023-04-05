package index

import "fmt"

const TreeSignature = 0x54524545

type TreeExtension struct {
	signature Signature
	size      uint32
	offset    uint32
}

func (e *TreeExtension) String() string {
	return fmt.Sprintf("[tree] signature:%v, offset:%v size:%v", e.signature, e.offset, e.size)
}

func NewTreeExtension(buf []byte, offset uint32, size uint32) (*TreeExtension, error) {
	return &TreeExtension{
		signature: TreeSignature,
		size:      size,
		offset:    offset,
	}, nil
}
