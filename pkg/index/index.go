package index

import (
	"encoding/binary"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

type File struct {
	Header
	indexBuf []byte
	offset   uint64

	eoie *EOIEExtension
}

const IndexSignature = 0x44495243

type Header struct {
	signature    Signature
	Version      uint32
	EntriesCount uint32
}

func (h *Header) String() string {
	return fmt.Sprintf("[header] signature:%v, version:%v, entriesCount:%v", h.signature, h.Version, h.EntriesCount)
}

func (f *File) ParseHeader() error {
	signature := binary.BigEndian.Uint32(f.indexBuf[f.offset : f.offset+4])
	if signature != IndexSignature {
		return fmt.Errorf("index parse header failed: %w", NewInvalidSignature(IndexSignature, signature))
	}

	f.Header.signature = Signature(signature)
	f.offset += 4
	f.Header.Version = binary.BigEndian.Uint32(f.indexBuf[f.offset : f.offset+4])
	if f.Header.Version > 4 || f.Header.Version < 2 {
		return fmt.Errorf("invalid index header version %d", f.Header.Version)
	}
	f.offset += 4
	f.Header.EntriesCount = binary.BigEndian.Uint32(f.indexBuf[f.offset : f.offset+4])
	f.offset += 4
	return nil
}

func (f *File) ShowHeader() {
	log.Println(&f.Header)
}

func (f *File) ShowFileInfo() {
	log.Println("[file] size:", len(f.indexBuf))
}

const SHA1Size = 20
const EOIESize = 4 + SHA1Size
const EOIESizeWithHeader = EOIESize + 4 + 4

func (f *File) ParseEndOfIndexEntriesExtension() (bool, error) {
	fileSize := len(f.indexBuf)

	extOffset := fileSize - EOIESizeWithHeader - SHA1Size
	if extOffset < 0 {
		return false, nil
	}

	eoie, err := NewEOIEExtension(uint32(extOffset), f.indexBuf[extOffset:extOffset+EOIESizeWithHeader])
	if err != nil {
		if errors.Is(err, ErrWrongSignature) {
			return false, nil
		}
		return false, err
	}
	f.eoie = eoie

	return true, nil
}

func NewFile(fileName string) (*File, error) {
	indexBuf, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return &File{
		indexBuf: indexBuf,
	}, nil
}
