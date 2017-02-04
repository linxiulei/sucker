package utils

import (
	"io"
	"os"
)

func CopyFile(src, dst string, modeSync bool) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}

	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	if modeSync {
		fileInfo, err := in.Stat()
		if err != nil {
			return err
		}

		if err := out.Chmod(fileInfo.Mode()); err != nil {
			return err
		}
	}
	err = out.Sync()
	return
}
