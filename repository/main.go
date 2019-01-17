package repository

import (
	"flag"
	"fmt"
	git "gopkg.in/src-d/go-git.v4"
	"os"
)

var directory = flag.String("dir", ".", "Directory to clone to")


func Run() {
	flag.Parse()
	fmt.Println(*directory)
	url := "https://github.com/src-d/go-git"
	r, err := git.PlainClone(*directory, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(127)
	}
	fmt.Println(r)
}