// helpers package contains helper functions that are commonly used across this application
package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ReadFromFileToBytes(path string) []byte {
	absPath, _ := filepath.Abs(path)
	jsonFile, err := os.Open(absPath)
	if err != nil {
		fmt.Println("Issue with the file being read")
		print(err)
	}
	json_byte_value, _ := ioutil.ReadAll(jsonFile)
	return json_byte_value
}
