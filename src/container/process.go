package container

import (
	log "github.com/sirupsen/logrus"
	"gocker/src/record"
	"os"
	"os/exec"
	"path"
	"syscall"
)

//
// NewParentProcess
// @Description: fork进程启动时，创建出管道，返回写管道用于写数据；读管道传入新进程
// @param tty
// @return *exec.Cmd
// @return *os.File
//
func NewParentProcess(tty bool, volume, ImageTarPath, cID string, envSlice []string) (*exec.Cmd, *os.File) {
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
	} else {
		exportContainerLogs(cID, &cmd.Stdout)
	}
	//创建容器工作空间
	mntURL := path.Join(record.RootURL, "mnt", cID)
	NewWorkSpace(record.RootURL, mntURL, ImageTarPath, volume, cID)
	//设置容器进程启动路径
	cmd.Dir = mntURL
	//传入管道文件句炳
	//指定要由新进程的打开文件
	cmd.ExtraFiles = []*os.File{readPipe}
	cmd.Env = append(os.Environ(), envSlice...)
	return cmd, writePipe
}
