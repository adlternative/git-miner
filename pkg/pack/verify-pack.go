package pack

import log "github.com/sirupsen/logrus"

func Verify(packPath string) error {
	packFile, err := NewPackFile(packPath)
	if err != nil {
		return err
	}
	err = packFile.ParseHeader()
	if err != nil {
		return err
	}
	err = packFile.ParseObjects()
	if err != nil {
		return err
	}
	for _, obj := range packFile.objects {
		log.Printf("index=%d offset=%d, type=%s, size=%d\n", obj.index, obj.offset, obj._type, obj.size)
	}

	return nil
}
