package run

import (
	log "github.com/sirupsen/logrus"
	"mini-docker/src/cgroups"
	"mini-docker/src/cgroups/subsystem"
	"mini-docker/src/container"
	"os"
)

//
// Run
// @Description:
// @param tty
// @param cmd
//
func Run(tty bool, cmdArray []string, config *subsystem.ResourceConfig) {
	parent := container.NewParentProcess(tty, cmdArray)
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
	log.Infof("parent process run")
	_ = parent.Wait()
	os.Exit(-1)
}
