package main


import (
	"encoding/json"
	"fmt"
	"github.ibm.com/AdvancedAnalyticsCanada/service-common-lib/common/config"
	log "github.ibm.com/AdvancedAnalyticsCanada/service-common-lib/common/logging"
	"github.ibm.com/AdvancedAnalyticsCanada/service-common-lib/common/util"
	"github.ibm.com/AdvancedAnalyticsCanada/service-common-lib/service"
	"github.ibm.com/AdvancedAnalyticsCanada/service-config-data/config-data-util/environment"
	"github.ibm.com/AdvancedAnalyticsCanada/service-config-data/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	conf "github.ibm.com/AdvancedAnalyticsCanada/service-config-data/config-data-util"

	"github.ibm.com/AdvancedAnalyticsCanada/service-config-data/gitutil"
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

	serviceAppName = "service-config-data"

	log.SetAppName(serviceAppName)

	args := os.Args

	if len(args) == 1 {

		githubRepoName = util.GetENV("REPO")
		if githubRepoName == "" {
			log.Fatal("ERROR: REPO name is required")
		}
		githubAccount = util.GetENV("GITACCOUNT")
		if githubAccount == "" {
			log.Warn("ERROR: GITACCOUNT is required")
		}
		githubApiToken = util.GetENV("APITOKEN")
		if githubApiToken == "" {
			log.Warn("ERROR: git APITOKEN is required")
		}
		serviceApiVersion = util.GetENV("APIVER")
		if serviceApiVersion == "" {
			log.Fatal("ERROR: service APIVER is required")
		}

		confEnvNames = []string{"dev", "sandbox"}

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

	servicePort = util.GetENV("PORT")
	if servicePort == "" {
		servicePort = "8000"
	}

	githubBranch = util.GetENV("GITBRANCH") ; if githubBranch == "" { githubBranch = "sandbox" }
	configFormat = util.GetENV("FORMAT")
	if configFormat == "" { configFormat = "json" }
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
	//service.RegisterHandler("/api/v2/configs", "GET", &handlers.ConfEnvHandler{Environments:ConfMappingOfEnvs})
	service.RegisterHandler("/api/v1/kernels/{environment}", "GET", &handlers.KernelHandler{Environments: ConfMappingOfEnvs})

	service.StartServer(servicePort)

	// TODO: Have to think how we want to log failures so that we can debug.
	//	 - If a pod dies in kubernetes the logs for it won't be available
	//	 - Need to write logs to some file and persist it probably
	//

	// TODO: Have better way to log and handle errors to avoid excessive prints
	// 	- https://hackernoon.com/golang-handling-errors-gracefully-8e27f1db729f

}

func ApiHandler(rw http.ResponseWriter, req *http.Request) {
	_, err := rw.Write([]byte(serviceApiVersion))
	if err != nil {
		log.Error("ERROR: Variable <g_api> is not defined properly")
	}
	rw.WriteHeader(http.StatusOK)
}


func initializeEnvironment()  {
	for _, envName := range confEnvNames {

		fs, repo, err := gitutil.GetRepoFromGit(githubAccount, githubApiToken, githubRepoName, envName)
		if err != nil {
			log.Info("Branch ", envName, " not intialized")
			continue
		}

		ConfMappingOfEnvs[envName] = &environment.Environment{
			FileSystem: fs,
			Repository: repo,
		}

		// Get the users
		arrayUserBytes, err := gitutil.GetFileFromRepo(fs, "users.json")

		if err == nil {
			_ = ConfMappingOfEnvs[envName].Users.CreateAllUsers(arrayUserBytes)
		}
	}


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
////@params
////
////@return
////
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

	if serviceDebugFlag == true {

		c := v.AllKeys()
		for i := 0; i < len(c) ; i++ {
			log.Debug( c[i] + " -> " + fmt.Sprintf("%s", v.Get(c[i])) )
		}
	}

	// look up the key and return value
	return fmt.Sprintf("%s", v.Get(key)), nil

}
