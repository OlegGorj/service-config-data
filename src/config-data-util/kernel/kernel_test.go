package kernel

import (
	"fmt"
	"os"
	"testing"
)

func TestGetKernels(t *testing.T) {
	folder, err  := os.Open("../../gitutil/test_data/kernels")
	if err != nil {
		fmt.Println(err)
	}

	info, _ := folder.Readdir(0)

	for index, fileName := range info {
		fmt.Println(fileName.Name())
		fmt.Println(index)
	}
}