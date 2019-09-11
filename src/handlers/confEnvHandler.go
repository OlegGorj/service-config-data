package handlers
//
//import (
//
//	"encoding/json"
//	log "github.com/oleggorj/service-common-lib/common/logging"

//	"net/http"
//)
//
//type ConfEnvHandler struct {
//	Environments conf.MappingToEnv
//}
//
//func (c *ConfEnvHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
//	response := make(map[string]conf.Environment)
//	for k, v := range c.Environments {
//		response[k] = conf.Environment{
//			Name: v.Name,
//			Users: v.Users,
//		}
//	}
//	err := json.NewEncoder(rw).Encode(response)
//	if err != nil {
//		log.Error("ERROR: The response variable in the user environment can't be encoded")
//		rw.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	rw.WriteHeader(http.StatusOK)
//	rw.Header().Set("Content-Type", "application/json")
//
//}
