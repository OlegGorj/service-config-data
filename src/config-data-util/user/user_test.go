package user

import (
	"memfilesystem"

	//"encoding/json"
	"fmt"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gotest.tools/assert"
	"io/ioutil"
	"os"
	"testing"
)

type TestingEnv struct {
	fs billy.Filesystem
	users Users
}

func generateUsers() *TestingEnv{

	path := "users.json"

	fs := memfs.New()
	f, _ := fs.Create(path)

	jsonFile, err := os.Open("../../gitutil/test_data/test_users1.json")

	if err != nil {
		fmt.Println("Issue with the file being read")
		print(err)
	}

	jsonBytes, _ := ioutil.ReadAll(jsonFile)
	listOfUsers := Users{}

	err = (&listOfUsers).CreateAllUsers(jsonBytes)

	if err != nil {
		fmt.Println(err)
	}
	_, err = f.Write(jsonBytes)

	if err != nil {
		fmt.Println(err)
	}

	testEnv := &TestingEnv{
		fs: fs,
		users:      listOfUsers,
	}
	return testEnv
}



func TestUser_Delete(t *testing.T) {

	testEnv := generateUsers()
	users := testEnv.users
	fs := testEnv.fs

	assert.Equal(t, len(users),3 ) // test initial size
	_ = users.Delete(fs,"abdullah@ibm.com")
	assert.Equal(t, len(users), 2)
	index, err := users.Read("abdullah@ibm.com")

	assert.Equal(t,index,-1) // test that correct user was deleted

	if err == nil {
		fmt.Println("error should be here, since user should not exist any more")
		t.Failed()
	}
}

func TestUser_Update(t *testing.T) {

	testEnv := generateUsers()
	users := testEnv.users
	fs := testEnv.fs

	newUserFile, _  := os.Open("../../gitutil/test_data/user.json")
	newUser, _ := ioutil.ReadAll(newUserFile)

	_ = users.Update(fs,"abdullah@ibm.com", newUser)

	index, _ := users.Read("abdullah@ibm.com")

	assert.Equal(t, index, -1)

}

func TestUser_Create(t *testing.T) {

	testEnv := generateUsers()
	users := testEnv.users
	fs := testEnv.fs

	newUserFile, _  := os.Open("../../gitutil/test_data/user.json")
	newUser, _ := ioutil.ReadAll(newUserFile)

	initialLength := len(users)

	originalBytes, err := memfilesystem.ReadFile(fs, "users.json")

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("ORIGINAL ", originalBytes)

	_ = users.Create(fs, newUser)


	afterBytes, err := memfilesystem.ReadFile(fs, "users.json")

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("AFTER ", afterBytes)



	//fmt.Println(string(afterBytes))


	assert.Equal(t, len(users), initialLength+1)
	assert.Equal(t, users[len(users)-1].Email, "newUser@ibm.com")

}

func TestUser_Read(t *testing.T) {

	testEnv := generateUsers()
	users := testEnv.users

	userToFind := "johnny@ibm.com"
	index, err := users.Read(userToFind)
	if err != nil {
		t.Fatal(err)
	}
	foundUser := testEnv.users[index]

	assert.Equal(t, foundUser.Email, userToFind)

}
