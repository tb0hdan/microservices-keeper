package utils

import (
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/tb0hdan/microservices-keeper/repository/logs"
)

func EnsureDir(name string) {
	_, err := os.Open(name)
	if os.IsNotExist(err) {
		logs.Logger.Println("Making dir", name)
		err = os.MkdirAll(name, os.ModePerm)
		if err != nil {
			logs.Logger.Fatalf("an error occured while running os.MkdirAll(): %+v\n", err)
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
