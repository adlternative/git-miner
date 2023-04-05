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
	//ieot *IEOTExtension

	exts []Extension
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
		return fmt.Errorf("index parse header failed: %w", NewInvalidSignature(IndexSignature, Signature(signature)))
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

func (f *File) parseEndOfIndexEntriesExtension() error {
	fileSize := len(f.indexBuf)

	extOffset := fileSize - EOIESizeWithHeader - SHA1Size
	if extOffset < 0 {
		return nil
	}

	eoie, err := NewEOIEExtension(f.indexBuf, uint32(extOffset))
	if err != nil {
		if errors.Is(err, ErrWrongSignature) {
			return nil
		}
		return err
	}
	f.eoie = eoie
	f.exts = append(f.exts, eoie)
	return nil
}

func (f *File) Parse() error {
	err := f.ParseHeader()
	if err != nil {
		return err
	}
	// eoie
	err = f.parseEndOfIndexEntriesExtension()
	if err != nil {
		return err
	}
	if f.eoie != nil {
		err := f.parseExtensions()
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *File) Show() error {
	f.ShowFileInfo()
	f.ShowHeader()

	for _, ext := range f.exts {
		log.Println(ext)
	}

	return nil
}

func (f *File) parseExtensions() error {
	offset := f.eoie.endOffIndexEntriesOffset

	for offset < f.eoie.offset {
		var ext Extension
		if offset+8 > uint32(len(f.indexBuf)) {
			return fmt.Errorf("short extension header length")
		}
		signature, size, err := ParseExtensionHeader(f.indexBuf[offset : offset+8])
		if err != nil {
			return err
		}

		switch signature {
		case IEOTSignature:
			ext, err = NewIEOTExtension(f.indexBuf, offset, size)
			if err != nil {
				return err
			}
		case TreeSignature:
			ext, err = NewTreeExtension(f.indexBuf, offset, size)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown Signature: %v", signature)
		}
		f.exts = append(f.exts, ext)
		offset += size + 8
	}

	return nil
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
