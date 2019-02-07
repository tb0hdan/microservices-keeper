package input_trello // nolint

type TrelloConfiguration struct {
	APIKey         string
	Token          string
	MessageHandler func(string) (string, error)
}
