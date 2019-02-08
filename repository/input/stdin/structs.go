package input_stdin // nolint

type STDInConfiguration struct {
	MessageHandler func(string) (string, error)
	Message        string
}
