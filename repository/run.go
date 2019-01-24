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

func RunWithAbstractGit(mygit GitInterface, myccd CCDInterface, configuration Configuration, msgfunc func() (string, error)) {
	// Configure git instance
	mygit.SetConfiguration(configuration)

	// Clone here
	err := mygit.Clone()
	if err != nil && err.Error() != "repository already exists"{
		log.Fatalf("an error occured while running mygit.Clone(): %+v\n", err)
	}

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

	text, err := msgfunc()
	if err != nil {
		log.Fatalf("an error occured while getting message: %+v\n", err)
	}

	fname, err := myccd.CreateCCD(text, ccdPath)
	if err != nil {
		log.Fatal(err)
	}
	err = myccd.WriteToCCD(text, path.Join(ccdPath, fname))
	if err != nil {
		log.Fatalf("an error occured while writing CCD: %+v\n", err)
	}

	err = mygit.PullAll()
	if err != nil && err.Error() != "already up-to-date"{
		log.Fatalf("an error occured while running mygit.PullAll(): %+v\n", err)
	}
	err = mygit.AddAllFiles()
	if err != nil {
		log.Fatalf("an error occured while running mygit.AddAllFiles(): %+v\n", err)
	}
	err = mygit.CommitAll()
	if err != nil {
		log.Fatalf("an error occured while running mygit.CommitAll(): %+v\n", err)
	}
	mygit.Push()
}
