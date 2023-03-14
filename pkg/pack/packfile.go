package pack

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"unsafe"

	gitzlib "github.com/adlternative/git-zlib-cgo"
)

const headerSize = 12
const Signature = 0x5041434b
const GitSha1Rawsz = 20

type PackFile struct {
	file       *os.File
	version    uint32
	objectNums uint32
	curOffset  uint64
	objects    []*Object

	inputBuf *buffer
}

func (pf *PackFile) fill(min uint64) ([]byte, error) {
	return pf.inputBuf.Fill(min)
}

func (pf *PackFile) buffer() []byte {
	return pf.inputBuf.Buffer()
}

func (pf *PackFile) use(length uint64) {
	pf.inputBuf.Use(length)
	pf.curOffset += length
}

func NewPackFile(packPath string) (*PackFile, error) {
	file, err := os.Open(packPath)
	if err != nil {
		return nil, err
	}
	return &PackFile{
		file:     file,
		inputBuf: newBuffer(file),
	}, nil
}

func (pf *PackFile) ParseHeader() error {
	header, err := pf.fill(headerSize)
	if err != nil {
		return err
	}
	defer pf.use(headerSize)

	if binary.BigEndian.Uint32(header[0:4]) != Signature {
		return fmt.Errorf("bad signature %v", header[0:4])
	}

	version := binary.BigEndian.Uint32(header[4:8])
	if version != 2 && version != 3 {
		return fmt.Errorf("bad version %d", version)
	}
	pf.version = version
	log.Printf("version = %d\n", version)
	objectNums := binary.BigEndian.Uint32(header[8:12])

	pf.objectNums = objectNums
	log.Printf("object nums = %d\n", objectNums)

	return nil
}

func MSB64(value uint64) uint8 {
	size := unsafe.Sizeof(value) * 8
	return uint8(value >> (size - 8))
}

func (pf *PackFile) ParseObject(index uint32) error {
	curOffset := pf.curOffset

	b, err := pf.readByte()
	if err != nil {
		return err
	}

	_type := ObjectType((b >> 4) & 7)
	size := uint64(b & 15)
	shift := 4

	for b&0x80 != 0 {
		b, err = pf.readByte()
		if err != nil {
			return err
		}

		size += (uint64(b) & 0x7f) << shift
		shift += 7
	}

	switch _type {
	case ObjRefDelta:
		_, err = pf.fill(GitSha1Rawsz)
		if err != nil {
			return err
		}

		// handle ref delta

		pf.use(GitSha1Rawsz)
	case ObjOfsDelta:
		b, err = pf.readByte()
		if err != nil {
			return err
		}

		baseOffset := uint64(b & 127)
		for b&128 != 0 {
			baseOffset++
			if baseOffset == 0 || (MSB64(baseOffset) != 0) {
				return fmt.Errorf("bad delta base object offset value")
			}

			if b, err = pf.readByte(); err != nil {
				return err
			}

			baseOffset = (baseOffset << 7) + uint64(b&127)
		}
		ofsOffset := curOffset - baseOffset
		if ofsOffset <= 0 || ofsOffset >= curOffset {
			return fmt.Errorf("delta base offset is out of bound: curOffset=%d, baseOffet=%d, b=%d", curOffset, baseOffset, b)
		}
	case ObjCommit, ObjTree, ObjBlob, ObjTag:
	default:
		return fmt.Errorf("bad type %v", _type)
	}

	obj := &Object{
		offset: curOffset,
		_type:  _type,
		size:   size,
	}
	pf.objects = append(pf.objects, obj)

	log.Printf("index=%d offset=%d, type=%s, size=%d\n", index, obj.offset, obj._type, obj.size)
	_, err = pf.unpackEntryData(int(obj.size), obj._type)
	if err != nil {
		return err
	}
	return nil
}

func (pf *PackFile) ParseObjects() error {
	for i := uint32(0); i < pf.objectNums; i++ {
		if err := pf.ParseObject(i); err != nil {
			return err
		}
	}

	return nil
}

func (pf *PackFile) readByte() (byte, error) {
	buf, err := pf.fill(1)
	if err != nil {
		return 0, err
	}
	c := buf[0]
	pf.use(1)
	return c, nil
}

func (pf *PackFile) Close() error {
	return pf.file.Close()
}

func (pf *PackFile) unpackEntryData(size int, _type ObjectType) ([]byte, error) {
	var err error
	outBuf := make([]byte, size)
	zstream := &gitzlib.GitZStream{}
	status := gitzlib.Z_OK

	err = zstream.InflateInit()
	if err != nil {
		return nil, err
	}
	zstream.SetOutBuf(outBuf, size)

	for status == gitzlib.Z_OK {
		_, err = pf.fill(1)
		if err != nil {
			return nil, err
		}

		allInputBuf := pf.buffer()
		inputLength := len(allInputBuf)
		//log.Printf("curoff=%d, inputlen=%d curdata=%d", pf.curOffset, inputLength, allInputBuf[0])
		zstream.SetInBuf(allInputBuf, inputLength)

		status, err = zstream.Inflate(0)
		if err != nil {
			return nil, err
		}

		pf.use(uint64(inputLength - zstream.AvailIn()))
	}
	if status != gitzlib.Z_STREAM_END || zstream.TotalOut() != size {
		return nil, fmt.Errorf("inflate returned %d", status)
	}

	err = zstream.InflateEnd()
	if err != nil {
		return nil, err
	}

	return outBuf, nil
}
