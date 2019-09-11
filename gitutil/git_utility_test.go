package gitutil

import (
	"encoding/json"
	"fmt"
	conf "github.com/oleggorj/service-config-data/config-data-util"
	userutil "github.com/oleggorj/service-config-data/config-data-util/user"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	// "io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type config_git struct {
	name_value      string
	api_token_value string
}

type Testwhitelist struct {
	Whitelist []string `json:"whitelist"`
	Admin     []string `json:"admin"`
}

/*
 * Grab the jhub-whitelist.json from git repo config-data
 * compare that with the sample sandbox_whitelist.json which is the same file as jhub-whitelist.json
 * Expect the two to be the same
 */

 // TODO: Rewrrite this test because it fails now due to repo change
 //	  - Should use the users.json possibly or some file for testing
func TestGetGitRepoConfigFile(t *testing.T) {
	// read file that is the exact same as one coming from repository whitelist.json
	absPath, _ := filepath.Abs("./test_data/sandbox_whitelist.json")
	jsonFile, err := os.Open(absPath)
	if err != nil {
		fmt.Println("Issue with the file being read")
		print(err)
		t.Fail()
	}
	jsonByteValue, _ := ioutil.ReadAll(jsonFile)

	var GitCredentials GitCredentials
	credentialsFile, _ := ioutil.ReadFile("../credentials.json")
	json.Unmarshal(credentialsFile, &GitCredentials)
	fmt.Println(GitCredentials)
	var whiteList Testwhitelist
	json.Unmarshal(jsonByteValue, &whiteList)
	fmt.Println(whiteList)
	fs, _, _ := GetRepoFromGit(GitCredentials.Account, GitCredentials.ApiToken, GitCredentials.RepoName, "sandbox")

	//// TODO: change this to work from the test branch, since the following code doesn't exist any more
	contents, _ := GetFileFromRepo(fs, "jhub-whitelist.json")
	var gitWhiteList Testwhitelist
	json.Unmarshal(contents, &gitWhiteList)
	fmt.Println(gitWhiteList)
	if gitWhiteList.Whitelist[0] != whiteList.Whitelist[0] {
		t.Errorf("The white lists do not equal to one another")
	}
	if len(gitWhiteList.Whitelist) != len(whiteList.Whitelist) {
		t.Errorf("The white lists do not equal to one another")
	}
	if len(gitWhiteList.Admin) != len(whiteList.Admin) {
		t.Errorf("The white lists do not equal to one another")
	}
	if gitWhiteList.Admin[2] != whiteList.Admin[2] {
		t.Errorf("The white lists do not equal to one another")
	}
}

func CreateTestingEnvironment(repo *git.Repository, fs billy.Filesystem, filePath string) map[string]*conf.Environment {
	absPath, _ := filepath.Abs(filePath)
	jsonFile, err := os.Open(absPath)
	if err != nil {
		fmt.Println("Issue with the file being read")
		print(err)
	}
	json_byte_value, _ := ioutil.ReadAll(jsonFile)
	users := userutil.CreateAllUsers(json_byte_value, nil)
	testEnvironment := make(map[string]*conf.Environment)
	testEnvironment["test"] = &conf.Environment{
		Name:       "test",
		FileSystem: fs,
		Repository: repo,
		Users:      users,
	}

	return testEnvironment
}

func TestUpdateFileOnGitRepo(t *testing.T) {
	var GitCredentials GitCredentials
	credentialsFile, _ := ioutil.ReadFile("../credentials.json")
	err := json.Unmarshal(credentialsFile, &GitCredentials)
	if err != nil {
		_ = fmt.Errorf("some issue with unmarshalling credentials")
	}
	fs, repo, _ := GetRepoFromGit(GitCredentials.Account, GitCredentials.ApiToken, GitCredentials.RepoName, "test")
	testEnvironment := CreateTestingEnvironment(repo, fs, "./test_data/test_users2.json")
	fmt.Println(len(testEnvironment["test"].Users))
	newUser := conf.User{
		Email: "newuser@ibm.com",
		UserMetadata: conf.UserMetaData{
			IsAdmin:                 true,
			ListOfTeamBucketMapping: nil,
			UserBucket:              "",
		},
	}
	testEnvironment["test"].Users = append(testEnvironment["test"].Users, newUser)
	fmt.Println(len(testEnvironment["test"].Users))

	err = UpdateFileOnGitRepo(testEnvironment["test"], "users.json")
	if err != nil {
		_ = fmt.Errorf("some issue with my udate")
	}
}
