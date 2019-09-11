package main

import (
	"encoding/json"
	"fmt"
	"github.com/oleggorj/service-common-lib/common/config"
	log "github.com/oleggorj/service-common-lib/common/logging"
	"github.com/oleggorj/service-common-lib/common/util"
	"github.com/oleggorj/service-common-lib/service"
	"github.com/oleggorj/service-config-data/config-data-util/environment"
	"github.com/oleggorj/service-config-data/handlers"
	conf "github.com/oleggorj/service-config-data/config-data-util"
	"github.com/oleggorj/service-config-data/gitutil"

	"github.com/spf13/viper"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	githubRepoName    string
	githubAccount     string
	githubApiToken    string
	servicePort       string
	serviceApiVersion string
	serviceAppName    string
	// legacy
	serviceDebugFlag  bool
	githubConfigFile  string
	githubBranch 	  string
	configFormat	  string

	ConfMappingOfEnvs conf.MappingToEnv
	confEnvNames       []string
	)


func init() {

	viper.SetConfigName("config")
	viper.AddConfigPath(".") ; viper.AddConfigPath("/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	serviceAppName = viper.GetString("service.name")
	log.SetAppName(serviceAppName)
	args := os.Args
	if len(args) == 1 {

		githubRepoName = viper.GetString("service.backend.repo")
		if githubRepoName == "" {
			log.Fatal("ERROR: REPO name is required")
		}
		githubAccount = viper.GetString("service.backend.account")
		if githubAccount == "" {
			log.Warn("warning: GITACCOUNT is required")
		}
		githubApiToken = viper.GetString("service.backend.token")
		if githubApiToken == "" {
			log.Warn("warning: git APITOKEN is required")
		}
		serviceApiVersion = viper.GetString("service.apiver")
		if serviceApiVersion == "" {
			log.Warn("warning: service APIVER is required")
		}

		// init list of branches
		confEnvNames = []string{}
		err := viper.UnmarshalKey("service.backend.branches", &confEnvNames)
		log.Debug("Branches from config: %v, %#v", err, confEnvNames)

	} else if args[1] == "dev" {

		serviceApiVersion = "test"

		var GitCredentials gitutil.GitCredentials
		jsonFile, _ := ioutil.ReadFile("credentials.json")
		err := json.Unmarshal(jsonFile, &GitCredentials)

		if err != nil {
			log.Info("No credentials file found")
			log.Fatal(err)
		}

		githubRepoName = GitCredentials.RepoName
		githubAccount = GitCredentials.Account
		githubApiToken = GitCredentials.ApiToken

		confEnvNames = []string{"test"}

	} else {
		log.Fatal("ERROR: Invalid arguments passed")
	}

	servicePort = viper.GetString("service.port")
	if servicePort == "" {
		servicePort = "8000"
	}

	githubBranch = util.GetENV("GITBRANCH") ; if githubBranch == "" { githubBranch = "sandbox" }
	configFormat = util.GetENV("FORMAT") ; if configFormat == "" { configFormat = "json" }
	configFile := util.GetENV("CONFIGFILE") ; if configFile == "" { configFile = "services" }

	ConfMappingOfEnvs = make(conf.MappingToEnv)
	initializeEnvironment()
}


func main() {

	log.Info("Starting service '" + serviceAppName + "'...")
	log.Info("INFO: Starting config data service for '", githubRepoName, "' environment.. ")
	// legacy handlers
	service.RegisterHandlerFunction("/api", "GET", ApiHandler)
	service.RegisterHandlerFunction("/api/v1/{app}/{env}/{key}", "GET", KeyHandler)
	service.RegisterHandlerFunction("/api/v1/{app}/{env}/{key}/{debug}", "GET", KeyHandler)

	// v2 handlers
	service.RegisterHandler("/api/v2/configs/{environment}/users", "GET", &handlers.UsersHandler{Environments: ConfMappingOfEnvs})
	service.RegisterHandler("/api/v2/configs/{environment}/users/{email}", "GET", &handlers.UserHandler{Environments: ConfMappingOfEnvs})
	service.RegisterHandler("/api/v2/configs/{environment}/users/{email}", "DELETE", &handlers.UserHandler{Environments: ConfMappingOfEnvs})
	service.RegisterHandler("/api/v2/configs/{environment}/users/{email}", "POST", &handlers.UserHandler{Environments: ConfMappingOfEnvs})
	service.RegisterHandler("/api/v2/configs/{environment}/users/{email}", "PUT", &handlers.UserHandler{Environments: ConfMappingOfEnvs})
	service.RegisterHandler("/api/v1/kernels/{environment}", "GET", &handlers.KernelHandler{Environments: ConfMappingOfEnvs})
	service.RegisterHandler("/api/v2/{app}/{env}/{key}", "GET", &handlers.KeyHandler{Environments: ConfMappingOfEnvs})
	service.RegisterHandler("/api/v2/{app}/{env}/{key}/debug", "GET", &handlers.KeyHandler{Environments: ConfMappingOfEnvs})

	service.StartServer(servicePort)
	// TODO: Have to think how we want to log failures so that we can debug.
	//	 - If a pod dies in kubernetes the logs for it won't be available
	//	 - Need to write logs to some file and persist it probably
	//
	// TODO: Have better way to log and handle errors to avoid excessive prints
	// 	- https://hackernoon.com/golang-handling-errors-gracefully-8e27f1db729f
}

func initializeEnvironment()  {
	for _, envName := range confEnvNames {

		fs, repo, err := gitutil.GetRepoFromGit(githubAccount, githubApiToken, githubRepoName, envName)
		if err != nil {
			log.Error("-- Branch ", envName, " not intialized. Does it exist?")
			log.Error(err)
			continue
		} else {
			log.Info("-- Branch ", envName, " is intialized.")
		}

		ConfMappingOfEnvs[envName] = &environment.Environment{
			FileSystem: fs,
			Repository: repo,
		}
	}
}



// Depricated
// Legacy code - needs to be cleaned up
//

func ApiHandler(rw http.ResponseWriter, req *http.Request) {
	_, err := rw.Write([]byte(serviceApiVersion))
	if err != nil {
		log.Error("ERROR: Variable <g_api> is not defined properly")
	}
	rw.WriteHeader(http.StatusOK)
}


//
func KeyHandler(rw http.ResponseWriter, req *http.Request) {
	serviceDebugFlag = false
	rw.Header().Set("Content-Type", "application/json")

	the_app := strings.ToLower(mux.Vars(req)["app"])
	the_env := strings.ToLower(mux.Vars(req)["env"])
	the_key := mux.Vars(req)["key"]

	log.Info("the_app " + the_app)

	if the_app == "" || the_env == "" || the_key == "" {
		log.Error("ERROR: <app>, <env> or <key> can not be empty.\n")
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}
	if debug := strings.ToLower(mux.Vars(req)["debug"]); debug == "debug" { serviceDebugFlag = true }

	if serviceDebugFlag == true {
		_, _ = rw.Write( []byte( fmt.Sprintf("param `app` = %s\n", the_app) ) )
		_, _ = rw.Write( []byte( fmt.Sprintf("param `env` = %s\n", the_env) ) )
		_, _ = rw.Write( []byte( fmt.Sprintf("param `key` = %s\n", the_key) ) )
		log.Debug(fmt.Sprintf("param `env` = %s, param `key` = %s", the_env, the_key))
	}

	// set config file
	githubConfigFile = the_app + ".json"
	val, err := getValue(the_key)
	if err != nil { rw.WriteHeader(http.StatusInternalServerError) ; return }
	_, _ = rw.Write( []byte( val ))

	rw.WriteHeader(http.StatusOK)
}
//@params
//
//@return
//
func getValue(key string) (string, error){
	if serviceDebugFlag == true {
		log.Debug( githubAccount," ", githubApiToken, " ",githubRepoName, " ",githubBranch, " ",githubConfigFile)
	}
	configFile, err := config.GetGitRepoConfigFile(githubAccount, githubApiToken, githubRepoName, githubBranch, githubConfigFile)
	if err != nil { return "", fmt.Errorf("ERROR: error retriving configuration: %v", err) }
	if configFile == "" { return "", fmt.Errorf("Can not resolve temp file name.") }

	// reading config file into Viper interface
	v, err := config.ReadConfig(configFile)
	if err != nil { return "", fmt.Errorf("Error when reading config: %v\n", err) }
	// debug part of endpoint
	if serviceDebugFlag == true {
		c := v.AllKeys()
		for i := 0; i < len(c) ; i++ {
			log.Debug( c[i] + " -> " + fmt.Sprintf("%s", v.Get(c[i])) )
		}
	}

	// look up the key and return value
	return fmt.Sprintf("%s", v.Get(key)), nil

}
