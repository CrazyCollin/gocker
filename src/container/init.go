package container

import (
	log "github.com/sirupsen/logrus"
	"os"
	"syscall"
)

//
// RunContainerInitProcess
// @Description:
// @param cmd
// @param args
// @return error
//
func RunContainerInitProcess(cmd string, args []string) error {
	log.Infof("command %s,args %s", cmd, args)
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		return err
	}
	argv := []string{cmd}
	if err := syscall.Exec(cmd, argv, os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}
