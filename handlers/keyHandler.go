package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	log "github.com/oleggorj/service-common-lib/common/logging"
	conf "github.com/oleggorj/service-config-data/config-data-util"

	"github.com/oleggorj/service-config-data/config-data-util/memfilesystem"
	"github.com/oleggorj/service-config-data/config-data-util/key"
)

type KeyHandler struct {
	Environments conf.MappingToEnv
}

func (u *KeyHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	//serviceDebugFlag := false
	rw.Header().Set("Content-Type", "application/json")

	appValue := strings.ToLower(mux.Vars(req)["app"])
	envValue := strings.ToLower(mux.Vars(req)["env"])
	keyValue := mux.Vars(req)["key"]
	if appValue == "" || envValue == "" || keyValue == "" {
		log.Error("ERROR: <app>, <env> or <key> can not be empty.\n")
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	environment := u.Environments[envValue]
	if environment == nil {
		log.Error("ERROR: Environment " + envValue + " does not exist")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Method == http.MethodGet {
		// read a key
		var k key.Keys
		fs := &environment.FileSystem
		bytes, _ := memfilesystem.ReadFile(*fs, appValue + ".json")
		environment.Keys , _ = k.Init( bytes )
		val, _ := environment.Keys.Read(keyValue)
		//fmt.Println("Keys Read:  " + keyValue + ", val:" + val )
		rw.Write([]byte(val))
		rw.WriteHeader(http.StatusOK)
	}

}
