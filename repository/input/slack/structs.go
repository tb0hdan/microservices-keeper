package input_slack // nolint

type SlackConfiguration struct {
	APIToken       string
	Endpoint       string
	MessageHandler func(string) (string, error)
	Application    func(configuration SlackConfiguration) error
}
