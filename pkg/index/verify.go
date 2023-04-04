package index

import log "github.com/sirupsen/logrus"

func Verify(fileName string) error {
	file, err := NewFile(fileName)
	if err != nil {
		return err
	}
	file.ShowFileInfo()
	err = file.ParseHeader()
	if err != nil {
		return err
	}
	file.ShowHeader()
	find, err := file.ParseEndOfIndexEntriesExtension()
	if err != nil {
		return err
	}
	if find {
		log.Print(file.eoie)
	} else {
		log.Printf("[eoie] not found")
	}

	return nil
}
