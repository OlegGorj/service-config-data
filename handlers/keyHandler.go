package handlers

import (
	_ "fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"encoding/json"
	"encoding/xml"
	"github.com/tidwall/gjson"

	log "github.com/oleggorj/service-common-lib/common/logging"
	conf "github.com/oleggorj/service-config-data/config-data-util"
	"github.com/oleggorj/service-config-data/config-data-util/memfilesystem"
)

type KeyHandler struct {
	Environments conf.MappingToEnv
}

// move this to utilis
func IsJSON(str string) bool {
    var js json.RawMessage
    return json.Unmarshal([]byte(str), &js) == nil
}

func (u *KeyHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	outformat := "plain"
	v := req.URL.Query()
	outformat =  v.Get("out")
	//serviceDebugFlag := false

	appValue := strings.ToLower(mux.Vars(req)["app"])
	envValue := strings.ToLower(mux.Vars(req)["env"])
	keyValue := strings.Replace( mux.Vars(req)["key"] , "@","#",-1) // mux.Vars(req)["key"]
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
		// read config json
		fs := &environment.FileSystem
		bytes, err := memfilesystem.ReadFile(*fs, appValue + ".json")
		if err != nil {
		    log.Error(err)
		}

		// cleanup new lines
		environment.JsonData = strings.Replace( string(bytes), "\n","",-1 )
		// get the value for the key (keyValue)
		val := gjson.Get( environment.JsonData , keyValue )

		var byteData []byte = []byte( val.String() )
		if  outformat == "json" {
			rw.Header().Set("Content-Type", "application/json")
			// don't marshal output if its already json
			if IsJSON( val.String() ) == false { byteData, err = json.Marshal( val.String()  ) }
		}else if outformat == "xml" {
			rw.Header().Set("Content-Type", "application/xml")
			byteData, err = xml.Marshal( val.String()  )
		}else{
			rw.Header().Set("Content-Type", "application/text")
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write( byteData )

	}

}
