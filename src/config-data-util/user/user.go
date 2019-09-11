package user

import (
	"encoding/json"
	"fmt"
	"gopkg.in/src-d/go-billy.v4"
	"config-data-util/memfilesystem"
	log "github.com/oleggorj/service-common-lib/common/logging"
)

type Users []User

type User struct {
	Email        string       `json:"email"`
	UserMetadata MetaData `json:"metadata"`
}

type MetaData struct {
	IsAdmin                 bool                `json:"is_admin"`
	ListOfTeamBucketMapping []TeamBucketMapping `json:"list_of_team_bucket_mapping"`
	UserBucket              string              `json:"user_bucket"`
}

type TeamBucketMapping struct {
	CosInstance string `json:"cos_instance"`
	TeamName    string `json:"team_name"`
	BucketName  string `json:"bucket_name"`
}

func (users *Users) CreateAllUsers(usersBytes []byte) error {

	var tempUsers Users
	err :=  json.Unmarshal(usersBytes, &tempUsers)
	if err != nil {
		log.Info("Problem unmarshalling new user data")
		log.Info(err)
		return err
	}
	*users = tempUsers
	return nil

}

func (users *Users) Create(fs *billy.Filesystem, data []byte) error {

	var newUser User
	err := json.Unmarshal(data,&newUser)
	if err != nil {
		log.Info("Problem unmarshalling new user data")
		log.Info(err)
		log.Fatal(err)
		return err
	}

	index, _ := users.Read(newUser.Email)
	if index != -1 {
		return fmt.Errorf("the user already exists")
	}

	*users = append(*users, newUser)
	userBytes, _ := json.Marshal(users)

	err = memfilesystem.OverWriteFile(*fs, "users.json", userBytes)
	if err != nil {
		*users = (*users)[:len(*users) - 1]
		return  err
	}

	return nil
}

func (users *Users) Read(email string) (int, error){

	for i := range *users{
		if (*users)[i].Email == email {
			return i,nil
		}
	}
	return -1, fmt.Errorf("user not found")
}

func (users *Users) Update(fs *billy.Filesystem, email string, data []byte) error {

	index, err  := users.Read(email)
	// TODO: Handle error better
	if err != nil {
		return err
	}

	var updatedUser User
	err = json.Unmarshal(data,&updatedUser)

	if err != nil {
		return err
	}

	(*users)[index] = updatedUser
	userBytes, _ := json.Marshal(users)
	err = memfilesystem.OverWriteFile(*fs, "users.json", userBytes)
	if err != nil {
		return  err
	}

	return nil
}

func (users *Users) Delete(fs *billy.Filesystem, email string) error {

	i, err := users.Read(email)
	if err != nil {
		return err
	}
	// Order doesn't matter so we just swap the last one in for the oen we're deleting and shorten the slice
	(*users)[i] = (*users)[len(*users) - 1]
	*users = (*users)[:len(*users) - 1]
	userBytes, _ := json.Marshal(users)
	err = memfilesystem.OverWriteFile(*fs, "users.json", userBytes)
	if err != nil {
		return  err
	}

	return nil
}
