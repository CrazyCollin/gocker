package run

import (
	log "github.com/sirupsen/logrus"
	"mini-docker/src/container"
	"os"
)

//
// Run
// @Description:
// @param tty
// @param cmd
//
func Run(tty bool, cmd string) {
	parent := container.NewParentProcess(tty, cmd)
	if err := parent.Start(); err != nil {
		log.Error(err)
		return
	}
	log.Infof("parent process run")
	_ = parent.Wait()
	os.Exit(-1)
}
