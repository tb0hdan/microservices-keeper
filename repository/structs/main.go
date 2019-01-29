package structs

import "github.com/tb0hdan/microservices-keeper/repository/configuration"

type GitInterface interface {
	Clone() error
	GetBranches() ([]string, error)
	PullAll() error
	AddAllFiles() error
	CommitAll() error
	Push()
	SetConfiguration(cfg configuration.Configuration)
}

type CCDInterface interface {
	ReadCCDs(directory string) (ccds []string)
	CreateCCD(text, dpath string) (name string, err error)
	WriteToCCD(text, filename string) error
	RebuildIndex(directory string) (err error)
}

type RunnerEntity struct {
	Git             GitInterface
	CCD             CCDInterface
	Directory       string
	Configuration   configuration.Configuration
	MessageFunction func() (string, error)
}
