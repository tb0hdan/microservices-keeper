package utils

import (
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus" // nolint
)

func EnsureDir(name string) {
	_, err := os.Open(name)
	if os.IsNotExist(err) {
		log.Println("Making dir", name)
		err = os.MkdirAll(name, os.ModePerm)
		if err != nil {
			log.Fatalf("an error occured while running os.MkdirAll(): %+v\n", err)
		}
	}
}

func FindSSHKey() (sshPath string, err error) {
	cuser, err := user.Current()
	if err != nil {
		return
	}
	keysPath := path.Join(cuser.HomeDir, ".ssh")
	keys, err := filepath.Glob(path.Join(keysPath, "id_*"))
	for _, key := range keys {
		if strings.HasSuffix(key, ".pub") {
			continue
		}
		sshPath = key
		if strings.HasSuffix(key, "id_rsa") {
			break
		}
	}
	return
}
