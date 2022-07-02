package container

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

//
// RunContainerInitProcess
// @Description:容器执行第一个进程
// @param cmd
// @param args
// @return error
//
func RunContainerInitProcess() error {

	if err := setUpMount(); err != nil {
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

//
// readUserCommand
// @Description: 读取程序传入参数
// @return []string
//
func readUserCommand() []string {
	readPipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(readPipe)
	if err != nil {
		log.Errorf("read init argv pipe error:%v", err)
		return nil
	}
	return strings.Split(string(msg), " ")
}

//
// setUpMount
// @Description: 初始化挂载点
// @return error
//
func setUpMount() error {
	if err := syscall.Mount("/", "/", "", syscall.MS_REC|syscall.MS_PRIVATE, ""); err != nil {
		return fmt.Errorf("set up mount proc error:%v", err)
	}
	//获取当前路径
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get current location error:%v", err)
	}
	log.Infof("current location:%s", pwd)
	err = pivotRoot(pwd)
	if err != nil {
		return err
	}
	//挂载/proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		log.Errorf("mount proc error:%v", err)
		return err
	}
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=775")
	return nil
}

func pivotRoot(root string) error {
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount rootfs to itself error:%v", err)
	}
	//当前路径下创建.pivot_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if _, err := os.Stat(pivotDir); err == nil {
		if err := os.Remove(pivotDir); err == nil {
			return err
		}
	}
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return fmt.Errorf("mkdir of pivot_root error:%v", err)
	}

	//pivot_root到新的rootfs
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root error:%v", err)
	}

	//更改当前目录到根目录
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir root error:%v", err)
	}

	//取消.pivot_root的挂载并删除
	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir error:%v", err)
	}
	return os.Remove(pivotDir)
}
