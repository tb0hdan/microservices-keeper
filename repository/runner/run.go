package runner

import (
	"log"
	"path"

	"github.com/tb0hdan/microservices-keeper/repository/structs"
	"github.com/tb0hdan/microservices-keeper/repository/utils"
)

func RunWithAbstractGit(entity *structs.RunnerEntity) { // nolint

	// Configure git instance
	entity.Git.SetConfiguration(entity.Configuration)

	// Clone here
	err := entity.Git.Clone()
	if err != nil && err.Error() != "repository already exists" {
		log.Fatalf("an error occured while running mygit.Clone(): %+v\n", err)
	}

	branches, err := entity.Git.GetBranches()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(branches)

	docPath := path.Join(entity.Directory, "doc")
	ccdPath := path.Join(docPath, "ccds")
	reportPath := path.Join(docPath, "reports")

	utils.EnsureDir(docPath)
	utils.EnsureDir(ccdPath)
	utils.EnsureDir(reportPath)

	text, err := entity.MessageFunction()
	if err != nil {
		log.Fatalf("an error occured while getting message: %+v\n", err)
	}

	fname, err := entity.CCD.CreateCCD(text, ccdPath)
	if err != nil {
		log.Fatal(err)
	}
	err = entity.CCD.WriteToCCD(text, path.Join(ccdPath, fname))
	if err != nil {
		log.Fatalf("an error occured while writing CCD: %+v\n", err)
	}

	err = entity.CCD.RebuildIndex(ccdPath)
	if err != nil {
		log.Fatalf("an error occured while rebuilding CCD index: %+v\n", err)
	}

	err = entity.Git.PullAll()
	if err != nil && err.Error() != "already up-to-date" {
		log.Fatalf("an error occured while running mygit.PullAll(): %+v\n", err)
	}
	err = entity.Git.AddAllFiles()
	if err != nil {
		log.Fatalf("an error occured while running mygit.AddAllFiles(): %+v\n", err)
	}
	err = entity.Git.CommitAll()
	if err != nil {
		log.Fatalf("an error occured while running mygit.CommitAll(): %+v\n", err)
	}
	entity.Git.Push()
}
