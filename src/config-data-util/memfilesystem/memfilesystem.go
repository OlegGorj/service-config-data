package memfilesystem

import (
	"fmt"
	"gopkg.in/src-d/go-billy.v4"
	"io/ioutil"
)

func OverWriteFile(fs billy.Filesystem, path string, dataToWrite []byte) error {

	//fsRef := *fs
	err := fs.Remove(path)
	if err != nil {

		return fmt.Errorf(fmt.Sprintf("File %s not found", path))
	}
	f, _ := fs.Create(path)
	_, _ = f.Write(dataToWrite)
	_ = f.Close()

	return nil

}

func ReadFile(fs billy.Filesystem, path string) ([]byte, error) {

	f, err := fs.Open(path)
	if err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("file does not exist")
	}

	bytes, err := ioutil.ReadAll(f)
	_ = f.Close()

	if err != nil {
		return nil, err
	}

	return bytes, err
}