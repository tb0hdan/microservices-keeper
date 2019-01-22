package repository

import (
	"flag"
	log "github.com/sirupsen/logrus"
)

var (
	directory = flag.String("dir", ".", "Directory to clone to")
	ssh_key   = flag.String("ssh", "", "Path to SSH key")
	url       = flag.String("url", "", "CCD repository URL")
	err       error
)

func Run() {

	configuration := NewConfiguration()
	configuration.Init(flag.CommandLine)


	*directory, _ = configuration.Get("directory")

	// Populate directory
	if *directory == "." || *directory == "" {
		*directory = "microservices-keeper-ccds"
	}

	*ssh_key, _ = configuration.Get("ssh_key")

	// Populate SSH key
	if *ssh_key == "" {
		*ssh_key, err = findSSHKey()
		if err != nil {
			log.Fatal(err)
		}
	}

	*url, _ = configuration.Get("url")

	if *url == "" {
		*url = "git@github.com:tb0hdan/microservices-keeper-ccds.git"
	}

	RunWithAbstractGit(NewGit(*url, *directory, *ssh_key), NewCCD(), configuration)
}
