package utils

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func DirOrFileExist(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, err
	} else if os.IsNotExist(err) {
		return false, err
	} else {
		log.Error(err)
		return false, err
	}
}
