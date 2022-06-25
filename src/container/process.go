package container

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

//
// NewParentProcess
// @Description: fork进程启动时，创建出管道，返回写管道用于写数据；读管道传入新进程
// @param tty
// @return *exec.Cmd
// @return *os.File
//
func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	//生成管道
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		log.Errorf("create pipe error:%v", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	return cmd, writePipe
}
