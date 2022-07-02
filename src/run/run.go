package run

import (
	log "github.com/sirupsen/logrus"
	"mini-docker/src/cgroups"
	"mini-docker/src/cgroups/subsystem"
	"mini-docker/src/container"
	"mini-docker/src/record"
	"os"
	"path"
	"strings"
)

//
// Run
// @Description:
// @param tty
// @param cmd
//
func Run(tty bool, cmdArray []string, config *subsystem.ResourceConfig, cgroupName string,
	volume, containerName, ImageTarPath, cID string, envSlice, portMappings []string, networkName string) {
	parent, writePipe := container.NewParentProcess(tty)
	//执行命令但不等待其结束
	//fork子进程，在/proc/self/exe中调用自己
	if err := parent.Start(); err != nil {
		log.Error(err)
		return
	}

	//记录容器信息

	//设置容器网络
	if networkName != "" {

	}

	cgroupManager := cgroups.NewCgroupManager(cgroupName + "-" + cID)
	if err := cgroupManager.Apply(parent.Process.Pid); err != nil {
		log.Errorf("cgroup apply error:%v", err)
		return
	}
	if err := cgroupManager.Set(config); err != nil {
		log.Errorf("cgroup set error:%v", err)
		return
	}

	sendInitCommand(cmdArray, writePipe)

	//等待结束
	if tty {
		//tty模式下父进程等待子进程结束
		if err := parent.Wait(); err != nil {
			log.Error(err)
		}
		cgroupManager.Destroy()
		mntURL := path.Join(record.RootURL, "mnt", cID)
		container.DeleteWorkspace(record.RootURL, mntURL, volume, cID)

		os.Exit(1)
	} else {
		//交由system pid=1的进程接管
		//fmt.Printf()
		//todo return container id
	}
}

func sendInitCommand(array []string, writePipe *os.File) {
	command := strings.Join(array, " ")
	log.Infof("all command is:%s", command)
	if _, err := writePipe.WriteString(command); err != nil {
		log.Errorf("writepipe write string error:%v", err)
		return
	}
	if err := writePipe.Close(); err != nil {
		log.Errorf("writepipe close error:%v", err)
	}
}
