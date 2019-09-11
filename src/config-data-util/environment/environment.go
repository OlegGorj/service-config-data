package environment

import (
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-git.v4"

	"config-data-util/user"
	"config-data-util/kernel"
	"config-data-util/key"
)


type Environment struct {
	Name       	string
	Repository 	*git.Repository
	FileSystem 	billy.Filesystem
	Users     	user.Users
	Kernels			[]kernel.Kernel
	Keys 				key.Keys
	JsonData		string

}
