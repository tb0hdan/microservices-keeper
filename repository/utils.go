package repository

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

func ensureDir(name string) {
	_, err := os.Open(name)
	if os.IsNotExist(err) {
		log.Println("Making dir", name)
		os.MkdirAll(name, os.ModePerm)
	}
}

func findSSHKey() (ssh_path string, err error) {
	cuser, err := user.Current()
	if err != nil {
		return
	}
	keys_path := path.Join(cuser.HomeDir, ".ssh")
	keys, err := filepath.Glob(path.Join(keys_path, "id_*"))
	for _, key := range keys {
		if strings.HasSuffix(key, ".pub") {
			continue
		}
		ssh_path = key
		if strings.HasSuffix(key, "id_rsa") {
			break
		}
	}
	return
}
