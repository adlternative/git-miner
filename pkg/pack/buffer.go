package pack

import (
	"fmt"
	"io"
)

const DefaultBufferSize = 1 << 12

type buffer struct {
	length uint32
	offset uint32
	buf    []byte
	reader io.Reader
}

func newBuffer(reader io.Reader) *buffer {
	return &buffer{
		buf:    make([]byte, DefaultBufferSize),
		reader: reader,
	}
}

// flush move the data to buffer head
func (b *buffer) flush() {
	copy(b.buf, b.buf[b.offset:b.offset+b.length])
	b.offset = 0
}

func (b *buffer) Buffer() []byte {
	return b.buf[b.offset : b.offset+b.length]
}

func (b *buffer) Fill(min uint32) ([]byte, error) {
	if min <= b.length {
		return b.buf[b.offset : b.offset+b.length], nil
	}
	if min > DefaultBufferSize {
		return nil, fmt.Errorf("cannot fill %d bytes", min)
	}

	b.flush()

	for b.length < min {
		ret, err := b.reader.Read(b.buf[b.offset+b.length:])
		if err != nil {
			if err == io.EOF {
				return nil, err
			}
		}
		b.length += uint32(ret)
	}

	return b.buf[:min], nil
}

func (b *buffer) Use(length uint32) {
	b.length -= length
	b.offset += length
}
