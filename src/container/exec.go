package container

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gocker/src/record"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"syscall"
)

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
