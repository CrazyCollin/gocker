package container

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
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
	log.Infof("RunContainerInitProcess command %s,args %s", cmd, args)
	err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		log.Errorf("private mount error:%v", err)
		return err
	}
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		log.Errorf("proc mount error: %v", err)
		return err
	}

	path, err := exec.LookPath(cmd)
	if err != nil {
		log.Errorf("can't find exec path:%s %v", cmd, err)
		return err
	}
	log.Infof("find path:%s", path)
	if err := syscall.Exec(path, args, os.Environ()); err != nil {
		log.Errorf("syscall exec error:%v", err.Error())
	}
	return nil
}
