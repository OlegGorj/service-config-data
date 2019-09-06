package handlers

import (
	//"bytes"
	"encoding/json"
	"fmt"
	conf "github.com/oleggorj/service-config-data/config-data-util"
	"github.com/oleggorj/service-config-data/config-data-util/environment"
	"github.com/oleggorj/service-config-data/config-data-util/user"
	"gotest.tools/assert"

	//"github.com/oleggorj/service-config-data/helpers"
	"github.com/gorilla/mux"
	//"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	is "gotest.tools/assert/cmp"
)



func generateUsers() *environment.Environment{

	path := "users.json"

	fs := memfs.New()
	f, _ := fs.Create(path)

	jsonFile, err := os.Open("../gitutil/test_data/test_users1.json")

	if err != nil {
		fmt.Println("Issue with the file being read")
		fmt.Println(err)
		return nil
	}

	jsonBytes, _ := ioutil.ReadAll(jsonFile)
	listOfUsers := user.Users{}

	err = (&listOfUsers).CreateAllUsers(jsonBytes)

	if err != nil {
		fmt.Println(err)
	}
	_, err = f.Write(jsonBytes)

	if err != nil {
		fmt.Println(err)
	}

	testEnv := &environment.Environment{
		FileSystem: fs,
		Users:      listOfUsers,
	}
	return testEnv
}



func TestUsersHandler_ServeHTTP(t *testing.T) {

	sampleEnvironmentMapping := make(conf.MappingToEnv)
	sampleEnvironmentMapping["test"] = generateUsers()

	req, err := http.NewRequest("GET", "/api/v1/configs/test/users", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.Handler(&UsersHandler{Environments: sampleEnvironmentMapping})

	router := mux.NewRouter()
	router.Handle("/api/v1/configs/{environment}/users", handler).Methods("GET")
	router.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)


	var response UserRequestResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	userBytes, _ := json.Marshal(response.Data)
	var usersArray user.Users
	_ = json.Unmarshal(userBytes, &usersArray)

	assert.Assert(t, is.Len(usersArray, 3))
}

//func TestUserHandlerGet_ServeHTTP(t *testing.T) {
//
//	sampleEnvironmentMapping := usersEnvGenHelper()
//
//	req, err := http.NewRequest("GET", "/api/v1/configs/test/users/abdullah@ibm.com", nil)
//
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	rr := httptest.NewRecorder()
//	handler := http.Handler(&UserHandler{Environments: sampleEnvironmentMapping})
//
//	router := mux.NewRouter()
//	router.Handle("/api/v1/configs/{environment}/users/{email}", handler).Methods("GET")
//	router.ServeHTTP(rr, req)
//
//	var usersArray conf.User
//	err = json.Unmarshal(rr.Body.Bytes(), &usersArray)
//
//
//	assert.Equal(t, usersArray.Email, sampleEnvironmentMapping["test"].Users[0].Email)
//
//
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	fmt.Println(usersArray)
//	fmt.Println(rr.Body)
//
//}

//func TestUserHandlerPost_ServeHTTP(t *testing.T) {
//
//	reqType := "POST"
//
//	userToUpdate := []byte(`{"email":"abdullah@ibm.com","metadata":{"is_admin":false,"list_of_team_bucket_mapping":[{"cos_instance":"client2_instance","team_name":"fun_team_2","bucket_name":"sample_bucket_2"}],"user_bucket":"abdullas_bucket"}}`)
//
//
//	req, err := http.NewRequest(reqType, "/api/v1/configs/test/users/abdullah@ibm.com", bytes.NewBuffer(userToUpdate))
//	if err != nil {
//		t.Fatal(err)
//	}
//	req.Header.Set("Content-Type", "application/json")
//
//	sampleEnvironmentMapping := usersEnvGenHelper()
//
//	rr := httptest.NewRecorder()
//	handler := http.Handler(&UserHandler{Environments: sampleEnvironmentMapping})
//
//	router := mux.NewRouter()
//	router.Handle("/api/v1/configs/{environment}/users/{email}", handler).Methods(reqType)
//	router.ServeHTTP(rr, req)
//
//	fmt.Println("TEST ", rr.Body)
//	var usersArray conf.User
//	err = json.Unmarshal(rr.Body.Bytes(), &usersArray)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	fmt.Println(usersArray)
//	fmt.Println(rr.Body)
	//
	//testEnv := generateUsers()

	//newUser = conf.User{
	//	Email: "bobby@ibm.com",
	//	UserMetadata: conf.UserMetaData{
	//		IsAdmin: false,
	//		UserBucket: "",
	//	},
	//}
	//userBytes, _ := json.Marshal(user)

	// Have to add user to array as well
	//testEnv.Users = append(testEnv.Users, newUser)
	//
	//usersBytes, _ := json.Marshal(testEnv.Users)
	//_ = testEnv.FileSystem.Remove("users.json")
	//f, _ := testEnv.FileSystem.Create("users.json")
	//_, _ = f.Write(usersBytes)
	//
	//f, _ = testEnv.FileSystem.Open("users.json")
	//contents, _ := ioutil.ReadAll(f)
	//var usersStruct []conf.User
	//json.Unmarshal(contents, &usersStruct)
	//fmt.Println(usersStruct)

	// Must also has them to the filesystem so that a git commit can occur


	//req, err := http.NewRequest("GET", "/api/v1/configs/test/users/abdullah@ibm.com", nil)
	//
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//rr := httptest.NewRecorder()
	//handler := http.Handler(&UserHandler{Environments: sampleEnvironmentMapping})
	//
	//router := mux.NewRouter()
	//router.Handle("/api/v1/configs/{environment}/users/{email}", handler).Methods("GET")
	//router.ServeHTTP(rr, req)
	//var usersArray conf.User
	//err = json.Unmarshal(rr.Body.Bytes(), &usersArray)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//fmt.Println(usersArray)
	//fmt.Println(rr.Body)

//}

//
//func usersEnvGenHelper() map[string]*conf.Environment {
//	userFile := helpers.ReadFromFileToBytes("../gitutil/test_data/test_users1.json")
//	sampleEnvironmentMapping := make(map[string]*conf.Environment)
//	sampleEnvironmentMapping["test"] = &conf.Environment{
//		Name:  "test",
//		Users: userutil.CreateAllUsers(userFile, nil),
//	}
//	return sampleEnvironmentMapping
//}
