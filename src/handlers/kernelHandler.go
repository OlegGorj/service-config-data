package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strings"

	log "github.com/oleggorj/service-common-lib/common/logging"
	conf "config-data-util"

)

type KernelHandler struct {
	Environments conf.MappingToEnv
}

// TODO: Implement the POST,DELETE, PUT request for Users
func (u *KernelHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	envValue := strings.ToLower(mux.Vars(req)["environment"])
	environment := u.Environments[envValue]

	if environment == nil {
		log.Error("ERROR: Environment does not exist")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	response := environment.Users
	err := json.NewEncoder(rw).Encode(response)
	if err != nil {
		log.Error("ERROR: The response variable in the user environment can't be encoded")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")

}
