package input_slack // nolint

import (
	"fmt"

	log "github.com/sirupsen/logrus" // nolint
)

const (
	SlackEvents     = 01
	SlackWebsockets = 10
)

func RunSlackLoop(config *SlackConfiguration, modes int) {
	c1 := make(chan error)
	c2 := make(chan error)
	switch modes {
	case SlackEvents:
		fmt.Println("Running Slack Events...")
		if err := RunEvents(config); err != nil {
			log.Fatalf("%+v", err)
		}
	case SlackWebsockets:
		fmt.Println("Running slack websockets...")
		if err := RunWebsockets(config); err != nil {
			log.Fatalf("%+v", err)
		}
	case SlackEvents | SlackWebsockets:
		fmt.Println("Running events and websockets...")
		go func() {
			c1 <- RunEvents(config)
		}()
		go func() {
			c2 <- RunWebsockets(config)
		}()
		select {
		case err1 := <-c1:
			log.Fatalf("RunEvents failed with: %+v", err1)
		case err2 := <-c2:
			log.Fatalf("RunWebsockets failed with: %+v", err2)
		}

	default:
		log.Fatalf("Unkwown mode: %d", modes)
	}

}
