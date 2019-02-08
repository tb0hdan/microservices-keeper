package input_stdin // nolint

import "github.com/tb0hdan/microservices-keeper/repository/logs"

func RunSTDIn(cfg *STDInConfiguration) {
	_, err := cfg.MessageHandler(cfg.Message)
	if err != nil {
		logs.Logger.Fatalf("An error occured while running STDIn method: %+v", err)
	}
}
