package repository

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"time"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"gopkg.in/src-d/go-git.v4"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

const (
	MAX_CCD_WORDS = 8
)

var (
	directory = flag.String("dir", ".", "Directory to clone to")
	ssh_key   = flag.String("ssh", "", "Path to SSH key")
	err       error
)

func Run() {
	flag.Parse()

	// Populate directory
	if *directory == "." || *directory == "" {
		*directory = "microservices-keeper-ccds"
	}

	// Populate SSH key
	if *ssh_key == "" {
		*ssh_key, err = findSSHKey()
		if err != nil {
			log.Fatal(err)
		}
	}

	url := "git@github.com:tb0hdan/microservices-keeper-ccds.git"

	pk, err := gitssh.NewPublicKeysFromFile("git", *ssh_key, "")

	_, err = git.PlainClone(*directory, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth:              pk,
	})

	if err == git.ErrRepositoryAlreadyExists {
		// git pull
	}

	if err != nil && err != git.ErrRepositoryAlreadyExists {
		log.Println(err)
		os.Exit(127)
	}

	r, err := git.PlainOpen(*directory)
	if err != nil {
		log.Println(err)
		os.Exit(129)
	}
	br, err := r.Branches()
	log.Println(br.Next())

	docPath := path.Join(*directory, "doc")
	ccdPath := path.Join(docPath, "ccds")
	reportPath := path.Join(docPath, "reports")

	ensureDir(docPath)
	ensureDir(ccdPath)
	ensureDir(reportPath)

	text := "We\nHave\nDecided\nTo have some serious fun tomorrow.\nWe'll build a rocket!\n"

	ccd := &CCD{}

	fname, err := ccd.createCCD(text, ccdPath)
	if err != nil {
		log.Fatal(err)
	}
	err = ccd.writeToCCD(text, path.Join(ccdPath, fname))

	w, err := r.Worktree()

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})

	status, err := w.Status()

	for k := range status {
		// TODO: Verify that file is actually a CCD
		log.Println(k)
		_, err = w.Add(k)
		if err != nil {
			log.Printf("Could not add file: %s\n", k)
			break
		}
	}

	commit, err := w.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})

	obj, err := r.CommitObject(commit)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(obj)

	po := &git.PushOptions{
		RemoteName: "origin",
		Auth:       pk,
	}

	err = r.Push(po)
	log.Println(err)
}
