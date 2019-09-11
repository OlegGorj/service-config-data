// Package handlers holds all of the handlers used by the service-config-data service
package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/oleggorj/service-common-lib/common/logging"
	conf "config-data-util"
	"config-data-util/user"
	"config-data-util/memfilesystem"
	"gitutil"
	"net/http"
	"strings"
)

type UsersHandler struct {
	Environments conf.MappingToEnv
}

type UserHandler struct {
	Environments conf.MappingToEnv
}

type UserRequestResponse struct {
	Data interface{} `json:"data"`
	Err string	`json:"err"`
}

func (u *UsersHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	rw.Header().Set("Content-Type", "application/json")

	envValue := strings.ToLower(mux.Vars(req)["environment"])
	environment := u.Environments[envValue]

	resp := UserRequestResponse{}

	if environment == nil {
		resp.Err = fmt.Sprintf("Environment '%s' not found", envValue)
		_ = json.NewEncoder(rw).Encode(resp)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	resp.Data = environment.Users
	_ = json.NewEncoder(rw).Encode(resp)

	rw.WriteHeader(http.StatusOK)

}

// TODO: Notice some redundancies in how error responses are being thrown, need to fix that

func (u *UserHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// TODO: Some repeated code here might want to refactor this later
	//	 - Env checking is going to be common between all handlers
	rw.Header().Set("Content-Type", "application/json")

	envValue := strings.ToLower(mux.Vars(req)["environment"])
	environment := u.Environments[envValue]

	resp := UserRequestResponse{}

	if environment == nil {
		resp.Err = fmt.Sprintf("Environment '%s' not found", envValue)
		_ = json.NewEncoder(rw).Encode(resp)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	userEmail := strings.ToLower(mux.Vars(req)["email"])
	log.Info("User: ", userEmail, " was Requested")

	users := &environment.Users
	fs := &environment.FileSystem
	repo := environment.Repository

	bytes, _ := memfilesystem.ReadFile(*fs, "users.json")
	fmt.Println("AFTER ", string(bytes))

	indexOfUser, err := users.Read(userEmail)

	// Check if user requested exists if not a POST request
	if req.Method != http.MethodPost && err != nil {
		resp.Err = fmt.Sprintf("User '%s' not found", userEmail)
		_ = json.NewEncoder(rw).Encode(resp)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Method == http.MethodGet {
		// read a user

		resp.Data = environment.Users[indexOfUser]
		_ = json.NewEncoder(rw).Encode(resp)
		rw.WriteHeader(http.StatusOK)

	} else {

		switch req.Method {

		// TODO: Some redundancies in code for put and post need to write method to fix that
		case http.MethodPost:
			// create a user
			// Questions: Do we want to create a user with just the new email. Or should we accept the whole body?

			// Check if already created

			newUser := user.User{}
			_ = json.NewDecoder(req.Body).Decode(&newUser)

			if newUser.Email == "" || newUser.Email != userEmail{
				resp.Err = fmt.Sprintf("User '%s' in route does not match '%s' in body", newUser.Email, userEmail)
				_ = json.NewEncoder(rw).Encode(resp)
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			newUserBytes, _ := json.Marshal(newUser)
			err := users.Create(fs, newUserBytes)

			if err != nil {
				resp.Err = fmt.Sprintf("Error: %s", err)
				_ = json.NewEncoder(rw).Encode(resp)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			resp.Data = fmt.Sprintf("Successfully created '%s'", userEmail)

		case http.MethodPut:

			updatedUser := user.User{}
			_ = json.NewDecoder(req.Body).Decode(&updatedUser)

			// Check that the email in route matches body of request
			if updatedUser.Email != userEmail {
				resp.Err = fmt.Sprintf("User '%s' in route does not match '%s' in body", updatedUser.Email, userEmail)
				_ = json.NewEncoder(rw).Encode(resp)
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			updatedUserBytes, _ := json.Marshal(updatedUser)
			err := users.Update(fs, userEmail, updatedUserBytes)

			if err != nil {
				resp.Err = fmt.Sprintf("Something went wrong updating user")
				_ = json.NewEncoder(rw).Encode(resp)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			resp.Data = fmt.Sprintf("Successfully updated '%s'", userEmail)

		case http.MethodDelete:

			err := users.Delete(fs, userEmail)

			if err != nil {
				resp.Err = fmt.Sprintf("There was an internal error")
				_ = json.NewEncoder(rw).Encode(resp)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			resp.Data = fmt.Sprintf("Successfully deleted '%s'", userEmail)
		}

		// TODO: Fix git push issue
		err = gitutil.UpdateFileOnGitRepo(*repo, *fs, "users.json")
		if err != nil {
			resp.Data = fmt.Sprint(resp.Data) + fmt.Sprintf("but %s", err)
		}
		_ = json.NewEncoder(rw).Encode(resp)
		rw.WriteHeader(http.StatusOK)

	}
}
