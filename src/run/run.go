package run

import (
	log "github.com/sirupsen/logrus"
	"mini-docker/src/cgroups"
	"mini-docker/src/cgroups/subsystem"
	"mini-docker/src/container"
	"os"
	"strings"
)

//
// Run
// @Description:
// @param tty
// @param cmd
//
func Run(tty bool, cmdArray []string, config *subsystem.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if err := parent.Start(); err != nil {
		log.Error(err)
		return
	}
	cgroupManager := cgroups.NewCgroupManager("Gocker-cgroup")
	defer cgroupManager.Destroy()
	if err := cgroupManager.Apply(parent.Process.Pid); err != nil {
		log.Errorf("cgroup apply error:%v", err)
		return
	}
	if err := cgroupManager.Set(config); err != nil {
		log.Errorf("cgroup set error:%v", err)
		return
	}

	sendInitCommand(cmdArray, writePipe)

	log.Infof("parent process run")
	_ = parent.Wait()
	os.Exit(-1)
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
