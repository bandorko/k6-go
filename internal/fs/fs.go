package fs

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
)

//CreateTempDir creates a temporary directory
func CreateTempDir(baseDir string) (string, error) {
	return ioutil.TempDir(baseDir, "k6gotmp")
}

//ReadFile reads a file into a []byte.
func ReadFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// Read the file into a byte slice
	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return bs, nil

}

func CopyDirectory(src string, dst string) error {
	return copy.Copy(src, dst)
}

func MoveFilesByExtension(srcDir string, dstDir string, ext string) error {
	directory, err := os.Open(srcDir)
	if err != nil {
		return err
	}
	defer directory.Close()

	files, err := directory.Readdir(-1)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ext {
				err = os.Rename(filepath.Join(srcDir, file.Name()), filepath.Join(dstDir, file.Name()))
				if err != nil {
					return err
				}

			}
		}
	}
	return nil
}
