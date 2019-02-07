package input_trello // nolint

import (
	"github.com/adlio/trello"
	"github.com/tb0hdan/microservices-keeper/repository/logs"
)

func RunTrello(cfg *TrelloConfiguration) {
	client := trello.NewClient(cfg.APIKey, cfg.Token)
	client.Logger = logs.Logger
	board, err := client.GetBoard("bOaRdID", trello.Defaults())
	if err != nil {
		logs.Logger.Fatalf("An error occured while getting Trello board: %+v", err)
	}
	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		logs.Logger.Fatalf("An error occured while getting Trello cards: %+v", err)
	}
	logs.Logger.Println(cards)
}
