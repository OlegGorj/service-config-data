package gitutil


import (
	"fmt"
	"github.com/oleggorj/service-config-data/config-data-util/memfilesystem"
	log "github.com/oleggorj/service-common-lib/common/logging"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"time"
	//"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"io/ioutil"
	"os"
	//"time"
)

// TODO: Write a stub for both functions
type GitCredentials struct {
	RepoName string `json:"repo_name"`
	Account  string `json:"account"`
	ApiToken string `json:"api_token"`
}

func GetRepoFromGit(gitAccount, apiToken, repoName, branch string) (billy.Filesystem, *git.Repository, error) {
	url := ""
	if gitAccount != "" {
		url = fmt.Sprintf("https://%s:%s@%s", gitAccount, apiToken, repoName)
	}else{
		url = fmt.Sprintf("https://%s", repoName)
	}
	log.Info("Cloning " + url)

	fs := memfs.New()
	storer := memory.NewStorage()
	repo, err := git.Clone(storer, fs, &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		SingleBranch:  true,
		Progress:      os.Stdout,
	})
	if err != nil {
		return nil, nil, err
	}

	return fs, repo, nil
}

func GetFileFromRepo(fs billy.Filesystem, file_name string) ([]byte, error) {

	f, err := fs.Open(file_name)
	if err != nil {
		log.Info("File ", file_name, " doesn't exist")
		log.Fatal(err)
		return nil, err
	}

	contents, err := ioutil.ReadAll(f)
	if err != nil {
		log.Info("Problem with reading file contents")
		log.Info(err)
	}

	return contents, nil
}

// TODO: Finish this
//func UpdateGitRepo(environment *conf.Environment, envName string) error {
//
//	log.Info("Updating it repo")
//
//
//	w, _ := (*environment).Repository.Worktree()
//
//	err := w.Pull(&git.PullOptions{
//		RemoteName: "test",
//	})
//	log.Info("Updating it repo")
//
//	return err
//}

// TODO: Add extract kernels utility from memfs
// 	- Might be better to but this in a different package
func UpdateFileOnGitRepo(repo git.Repository, fs billy.Filesystem, filePath string) error {
	w, _ := (repo).Worktree()
	bytes, _ := memfilesystem.ReadFile(fs, filePath)
	fmt.Println("AFTER ", string(bytes))
	_ , _ = w.Add(filePath)

	_, err := w.Commit("test commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "serviceaccount-config-data",
			Email: "serviceaccount-config-data@gmail.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("Error comitting data to repo")
	}

	err = (repo).Push(&git.PushOptions{})
	if err != nil {
		log.Info(err)
		return fmt.Errorf("Error pushing data to repo")
	}
	return nil
}
