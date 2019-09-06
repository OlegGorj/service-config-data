package key

import (
	_ "encoding/json"
	"fmt"
	_ "gopkg.in/src-d/go-billy.v4"
	_ "github.com/oleggorj/service-config-data/config-data-util/memfilesystem"
	_ "github.com/oleggorj/service-common-lib/common/logging"
)

type Keys []Key

type Key struct {
	Key string
	Val string
}

func (keys *Keys) Read(key string) (string, error){

	for i := range *keys{
		if (*keys)[i].Key == key {
			return (*keys)[i].Val, nil
		}
	}
	return "", fmt.Errorf("user not found")
}
