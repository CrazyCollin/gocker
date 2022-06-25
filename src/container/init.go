package container

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

//
// RunContainerInitProcess
// @Description:
// @param cmd
// @param args
// @return error
//
func RunContainerInitProcess() error {
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

	cmdArray := readUserCommand()
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("can't find exec path:%s %v", cmdArray[0], err)
		return err
	}
	log.Infof("find path:%s", path)
	if err := syscall.Exec(path, cmdArray, os.Environ()); err != nil {
		log.Errorf("syscall exec error:%v", err.Error())
	}
	return nil
}

func readUserCommand() []string {
	readPipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(readPipe)
	if err != nil {
		log.Errorf("read init argv pipe error:%v", err)
		return nil
	}
	return strings.Split(string(msg), " ")
}
