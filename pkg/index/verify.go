package index

func Verify(fileName string) error {
	file, err := NewFile(fileName)
	if err != nil {
		return err
	}

	err = file.Parse()
	if err != nil {
		return err
	}

	err = file.Show()
	if err != nil {
		return err
	}

	return nil
}
