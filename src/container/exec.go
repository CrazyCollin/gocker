package container

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gocker/src/record"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
)

//
// ExecContainer
// @Description: exec容器
// @param cID
// @param cmdArray
//
func ExecContainer(cID string, cmdArray []string) {
	containerInfo, err := getContainerInfo(cID)
	if err != nil {
		log.Errorf("exec container error:%v", err)
		return
	}
	cmdStr := strings.Join(cmdArray, " ")
	log.Infof("env container pid %s", containerInfo.Pid)
	log.Infof("env contianer ccommand:%s", cmdStr)

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//设置pid和command
	if err := os.Setenv(record.EnvExecPid, containerInfo.Pid); err != nil {
		log.Errorf("exec container of set pid env error:%v", err)
		return
	}
	if err := os.Setenv(record.EnvExecCmd, cmdStr); err != nil {
		log.Errorf("exec container of set command env error:%v", err)
		return
	}
	envs := getEnvsById(containerInfo.Pid)
	cmd.Env = append(os.Environ(), envs...)
	if err := cmd.Run(); err != nil {
		log.Errorf("run container command error:%v", err)
		return
	}
}

//
// getContainerInfo
// @Description: 根据cID获取容器信息
// @param cID
// @return *record.ContainerInfo
// @return error
//
func getContainerInfo(cID string) (*record.ContainerInfo, error) {
	infoFilePath := path.Join(record.DefaultInfoLocation, cID, record.ConfigName)
	info, err := ioutil.ReadFile(infoFilePath)
	if err != nil {
		log.Errorf("read contianer info error:%v", err)
		return nil, err
	}
	containerInfo := &record.ContainerInfo{}
	err = json.Unmarshal(info, containerInfo)
	if err != nil {
		log.Errorf("decode contianer info file error:%v", err)
		return nil, err
	}
	return containerInfo, nil
}

//
// StopContainer
// @Description: 关闭容器
// @param cID
//
func StopContainer(cID string) {
	containerInfo, err := getContainerInfo(cID)
	if err != nil {
		log.Errorf("stop container error:%v", err)
		return
	}
	//kill容器进程
	pid, _ := strconv.Atoi(containerInfo.Pid)
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		log.Errorf("kill container process error:%v", err)
		return
	}
	//写入容器info
	containerInfo.Status = record.STOP
	containerInfo.Pid = " "
	info, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("stop container error:%v", err)
		return
	}
	infoFilePath := path.Join(record.DefaultInfoLocation, cID, record.ConfigName)
	if err := ioutil.WriteFile(infoFilePath, info, 0622); err != nil {
		log.Errorf("kill container process error:%v", err)
		return
	}
}

//
// RemoveContainer
// @Description: 删除指定容器
// @param cID
//
func RemoveContainer(cID string) {
	containerInfo, err := getContainerInfo(cID)
	if err != nil {
		log.Errorf("remove contianer error:%v", err)
		return
	}
	//删除容器相关信息
	if containerInfo.Status == record.STOP {
		infoFilePath := path.Join(record.DefaultInfoLocation, cID)
		if err := os.RemoveAll(infoFilePath); err != nil {
			log.Errorf("remove container error:%v", err)
			return
		}
		mntURL := path.Join(record.RootURL, "mnt", cID)
		DeleteWorkspace(record.RootURL, mntURL, containerInfo.Volume, cID)
	} else {
		log.Warnf("please set contianer stopped first")
	}
}

//
// getEnvsById
// @Description: 根据容器pid得到其环境变量
// @param cPID
// @return []string
//
func getEnvsById(cPID string) []string {
	//进程的环境变量都存在/proc/pid/environ中
	envsPath := fmt.Sprintf("/proc/%s/environ", cPID)
	envsInfo, err := ioutil.ReadFile(envsPath)
	if err != nil {
		log.Errorf("get container environment error:%v", err)
		return nil
	}
	envs := strings.Split(string(envsInfo), "\u0000")
	return envs
}
