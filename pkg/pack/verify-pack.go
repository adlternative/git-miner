package pack

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
	return nil
}
