package repository

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"path"

	"github.com/tb0hdan/microservices-keeper/repository/input/slack" // nolint
	"github.com/tb0hdan/microservices-keeper/repository/output/ccd"
	"github.com/tb0hdan/microservices-keeper/repository/output/git"
	"github.com/tb0hdan/microservices-keeper/repository/runner"
	"github.com/tb0hdan/microservices-keeper/repository/structs"

	log "github.com/sirupsen/logrus" // nolint

	cfg "github.com/tb0hdan/microservices-keeper/repository/configuration"
	"github.com/tb0hdan/microservices-keeper/repository/utils"
)

var (
	directory = flag.String("dir", ".", "Directory to clone to")
	sshKey    = flag.String("ssh", "", "Path to SSH key")
	url       = flag.String("url", "", "CCD repository URL")
	gitUser   = flag.String("gituser", "", "Git user. Defaults to currently logged in one.")
	message   = flag.String("message", "", "Message to add. Ignored when pipe is used.")
	version   = flag.Bool("version", false, "Print version")
	// slack
	token    = flag.String("slack-token", "", "Slack Token")
	modes    = flag.Int("slack-modes", 0, "Slack modes: 01 - Events, 10 - WebSockets, 11 - Both")
	endpoint = flag.String("slack-endpoint", "/events-endpoint", "HTTP endpoint for Slack Events API")
)

func Run(bversion, buildID string) { // nolint
	var (
		err error
		usr *user.User
	)
	// Flags are parsed inside NewConfiguration()
	configuration := cfg.NewConfiguration()
	configuration.Init(flag.CommandLine)

	if *version {
		sname := path.Base(os.Args[0])
		if sname == "main" {
			sname = "microservices-keeper"
		}
		log.Printf("%s version %s-%s\n", sname, bversion, buildID)
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
		*sshKey, err = utils.FindSSHKey()
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

	if *message == "" && *token == "" {
		log.Fatal("At least one of [message|token] is required")
	}

	if *token == "" {
		msgfunc := func() (string, error) {
			return *message, nil
		}

		entity := &structs.RunnerEntity{
			Git:             git.NewGit(*url, *directory, *sshKey),
			CCD:             ccd.NewCCD(),
			Configuration:   configuration,
			MessageFunction: msgfunc,
			Directory:       *directory,
		}

		runner.RunWithAbstractGit(entity)
	} else {
		msgHandler := func(msg string) (r string, err error) {

			msgfunc := func() (string, error) {
				return msg, nil
			}

			entity := &structs.RunnerEntity{
				Git:             git.NewGit(*url, *directory, *sshKey),
				CCD:             ccd.NewCCD(),
				Configuration:   configuration,
				MessageFunction: msgfunc,
				Directory:       *directory,
			}

			runner.RunWithAbstractGit(entity)

			return
		}
		slackCfg := &input_slack.SlackConfiguration{
			APIToken:       *token,
			Endpoint:       *endpoint,
			MessageHandler: msgHandler,
			Application:    nil,
		}
		input_slack.RunSlackLoop(slackCfg, *modes)
	}

}
