package repository

import (
	"log"
	"time"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"gopkg.in/src-d/go-git.v4"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type GoGit struct {
	Directory  string
	URL        string
	pk         *gitssh.PublicKeys
	r          *git.Repository
	w          *git.Worktree
	SSHKeyPath string
	Configuration Configuration
}


func (gg *GoGit) SetConfiguration(configuration Configuration) {
	gg.Configuration = configuration
}

//Auth
func (gg *GoGit) Auth() (err error) {
	gg.pk, err = gitssh.NewPublicKeysFromFile("git", gg.SSHKeyPath, "")
	if err != nil {
		log.Println("Auth", err)
	}
	return
}

//Clone
func (gg *GoGit) Clone() (err error) {

	_, err = git.PlainClone(gg.Directory, false, &git.CloneOptions{
		URL:               gg.URL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth:              gg.pk,
	})

	return err
}

//Open
func (gg *GoGit) Open() (err error) {
	gg.r, err = git.PlainOpen(gg.Directory)
	if err != nil {
		log.Println("Open", err)
	}
	return
}

//Push
func (gg *GoGit) Push() {
	po := &git.PushOptions{
		RemoteName: "origin",
		Auth:       gg.pk,
	}

	err = gg.r.Push(po)
	log.Println(err)
}

//GetRepository
func (gg *GoGit) GetRepository() *git.Repository {
	return gg.r
}

//PrepareWorkTree
func (gg *GoGit) PrepareWorkTree() error {
	gg.w, err = gg.r.Worktree()
	if err != nil {
		log.Println("PrepareWorkTree", err)
	}
	return err
}

//GetWorkTree
func (gg *GoGit) GetWorkTree() *git.Worktree {
	return gg.w
}

//CommitAll
func (gg *GoGit) CommitAll() error {
	name, err := gg.Configuration.Get("name")
	email, err := gg.Configuration.Get("email")
	commit, err := gg.w.Commit("CCD update", &git.CommitOptions{
		Author: &object.Signature{
			Name:  name,
			Email: email,
			When:  time.Now(),
		},
	})

	obj, err := gg.r.CommitObject(commit)
	if err != nil {
		log.Fatal("CommitAll", err)
	}

	log.Println(obj)
	return err
}

//AddAllFiles
func (gg *GoGit) AddAllFiles() error {
	status, err := gg.w.Status()

	for k := range status {
		// TODO: Verify that file is actually a CCD
		log.Println(k)
		_, err = gg.w.Add(k)
		if err != nil {
			log.Printf("Could not add file: %s\n", k)
			break
		}
	}
	return err
}

//PullAll
func (gg *GoGit) PullAll() error {
	err = gg.w.Pull(&git.PullOptions{RemoteName: "origin", Auth: gg.pk})
	if err != nil {
		log.Println("PullAll", err)
	}
	return err
}

//GetBranches
func (gg *GoGit) GetBranches() (branches []string, err error) {
	br, err := gg.r.Branches()
	for {
		branch, err := br.Next()
		if err != nil {
			break
		}
		branches = append(branches, branch.String())
	}
	return
}

//NewGit
func NewGit(url, directory, ssh_key string) (ggn *GoGit) {
	ggn = &GoGit{URL: url, Directory: directory, SSHKeyPath: ssh_key}

	err = ggn.Auth()
	if err != nil {
		log.Fatalf("Auth: %+v\n", err)
	}


	pull := false
	err = ggn.Clone()
	if err == git.ErrRepositoryAlreadyExists {
		// git pull
		pull = true
	} else if err != nil {
		log.Fatalf("Clone: %+v\n", err)
	}

	err = ggn.Open()

	if err != nil {
		log.Fatalf("Open: %+v\n", err)
	}


	err = ggn.PrepareWorkTree()
	if err != nil {
		log.Fatalf("PrepareWorkTree: %+v\n", err)
	}

	if pull {
		err = ggn.PullAll()
		if err != nil {
			log.Println("PullAll: %+v\n", err)
		}
	}
	return
}
