package index

func Verify(fileName string) error {
	file, err := NewFile(fileName)
	if err != nil {
		return err
	}
	err = file.ParseHeader()
	if err != nil {
		return err
	}
	file.ShowHeader()

	return nil
}
