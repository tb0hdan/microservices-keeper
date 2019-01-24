package repository

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"path"

	log "github.com/sirupsen/logrus" // nolint
)

var (
	directory = flag.String("dir", ".", "Directory to clone to")
	sshKey    = flag.String("ssh", "", "Path to SSH key")
	url       = flag.String("url", "", "CCD repository URL")
	gitUser   = flag.String("gituser", "", "Git user. Defaults to currently logged in one.")
	message   = flag.String("message", "", "Message to add. Ignored when pipe is used.")
	version   = flag.Bool("version", false, "Print version")
)

func Run(bversion, buildID string) { // nolint
	var (
		err error
		usr *user.User
	)
	// Flags are parsed inside NewConfiguration()
	configuration := NewConfiguration()
	configuration.Init(flag.CommandLine)

	if *version {
		sname := path.Base(os.Args[0])
		if sname == "main" {
			sname = "microservices-keeper"
		}
		fmt.Printf("%s version %s-%s\n", sname, bversion, buildID)
		os.Exit(1)
	}

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

	if *message == "" {
		log.Fatal("Cannot run with empty message")
	}

	msgfunc := func() (string, error) {
		return *message, nil
	}

	RunWithAbstractGit(NewGit(*url, *directory, *sshKey), NewCCD(), configuration, msgfunc)
}
