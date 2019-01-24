package repository

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"

	log "github.com/sirupsen/logrus" // nolint
)

var (
	directory = flag.String("dir", ".", "Directory to clone to")
	sshKey    = flag.String("ssh", "", "Path to SSH key")
	url       = flag.String("url", "", "CCD repository URL")
	gitUser   = flag.String("gituser", "", "Git user. Defaults to currently logged in one.")
	message   = flag.String("message", "", "Message to add. Ignored when pipe is used.")
)

func Run() { // nolint
	var (
		err error
		usr *user.User
	)
	configuration := NewConfiguration()
	configuration.Init(flag.CommandLine)

	*directory, _ = configuration.Get("directory")

	// Populate directory
	if *directory == "." || *directory == "" {
		*directory = "microservices-keeper-ccds"
	}

	*sshKey, _ = configuration.Get("ssh_key")

	// Populate SSH key
	if *sshKey == "" {
		*sshKey, err = findSSHKey()
		if err != nil {
			log.Fatal(err)
		}
	}

	if *gitUser == "" {
		usr, err = user.Current()
		if err != nil {
			log.Fatalf("could not get current user: %+v\n", err)
		}
		*gitUser = usr.Name
	}

	*url, _ = configuration.Get("url")

	if *url == "" {
		*url = fmt.Sprintf("git@github.com:%s/%s.git", *gitUser, *directory)
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Size() > 0 {
		// there's data from pipe
		reader := bufio.NewReader(os.Stdin)

		var output []rune

		for {
			input, _, err := reader.ReadRune()
			if err != nil && err == io.EOF {
				break
			}
			output = append(output, input)
		}
		*message = string(output)
	}

	msgfunc := func() (string, error) {
		return *message, nil
	}

	RunWithAbstractGit(NewGit(*url, *directory, *sshKey), NewCCD(), configuration, msgfunc)
}
