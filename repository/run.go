package repository

import (
	"fmt"
	"log"
	"path"
)

type GitInterface interface {
	Clone() error
	GetBranches() ([]string, error)
	PullAll() error
	AddAllFiles() error
	CommitAll() error
	Push()
	SetConfiguration(configuration Configuration)
}

type CCDInterface interface {
	ReadCCDs(directory string) (ccds []string)
	CreateCCD(text, dpath string) (name string, err error)
	WriteToCCD(text, filename string) error
}

func RunWithAbstractGit(mygit GitInterface, myccd CCDInterface, configuration Configuration) {
	// Configure git instance
	mygit.SetConfiguration(configuration)

	// Clone here
	mygit.Clone()


	branches, err := mygit.GetBranches()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(branches)

	docPath := path.Join(*directory, "doc")
	ccdPath := path.Join(docPath, "ccds")
	reportPath := path.Join(docPath, "reports")

	ensureDir(docPath)
	ensureDir(ccdPath)
	ensureDir(reportPath)

	text := "We\nHave\nDecided\nTo have some serious fun tomorrow.\nWe'll build a rocket!\n"

	fname, err := myccd.CreateCCD(text, ccdPath)
	if err != nil {
		log.Fatal(err)
	}
	err = myccd.WriteToCCD(text, path.Join(ccdPath, fname))

	mygit.PullAll()
	mygit.AddAllFiles()
	mygit.CommitAll()
	mygit.Push()
}
