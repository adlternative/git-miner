package index

import (
	"encoding/binary"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

type File struct {
	Header
	indexBuf []byte
	offset   uint64
}

const HardCodeSignature = 0x44495243

type Header struct {
	Signature    uint32
	Version      uint32
	EntriesCount uint32
}

func (f *File) ParseHeader() error {
	f.Header.Signature = binary.BigEndian.Uint32(f.indexBuf[f.offset : f.offset+4])
	if f.Header.Signature != HardCodeSignature {
		return fmt.Errorf("invalid index header signature %0x", f.Header.Signature)
	}

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
	log.Printf("signature = %0x\n", f.Header.Signature)
	log.Printf("version = %d\n", f.Header.Version)
	log.Printf("entries count = %d\n", f.Header.EntriesCount)
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
