package input_slack // nolint

import log "github.com/sirupsen/logrus" // nolint

func RunSlack(config *SlackConfiguration) {
	err := config.Application(*config)
	if err != nil {
		log.Fatalf("An error occured while running app: %+v", err)
	}
}
