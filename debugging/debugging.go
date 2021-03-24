package debugging

import (
	"io"
	"os"
)

func DumpReaderToFile(r io.Reader, fn string) error {
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}

	return nil
}
