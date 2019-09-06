package key

import (
	"encoding/json"
	"fmt"
	_ "gopkg.in/src-d/go-billy.v4"
	_ "github.com/oleggorj/service-config-data/config-data-util/memfilesystem"
	log "github.com/oleggorj/service-common-lib/common/logging"
)

type Keys []Key

type Key struct {
	Key string
	Val string
}

func (keys *Keys) Init(jsonBuffer []byte) (Keys, error){
	//fmt.Println("AFTER ", string(jsonBuffer))
	keys_arr := make(map[string]interface{})
	err := json.Unmarshal(jsonBuffer, &keys_arr)
	if err != nil {
	    log.Error(err)
	}
	var allkeys Keys
	for j_key, j_value := range keys_arr {
			var k Key
			k.Key = j_key
			k.Val = fmt.Sprintf("%s", j_value )
			allkeys = append(allkeys, k)
			//fmt.Println("key : ", k.Key, " value : ", k.Val)
  }

	return allkeys, nil
}

func (keys *Keys) Read(key string) (string, error){

	for i := range *keys{
		if (*keys)[i].Key == key {
			return (*keys)[i].Val, nil
		}
	}
	return "", fmt.Errorf("key not found")
}
