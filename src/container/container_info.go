package container

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gocker/src/record"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"
)

//
// RandContainerIDGenerator
// @Description: 随机生成容器id
// @param n
// @return string
//
func RandContainerIDGenerator(n int) string {
	if n < 0 || n > 32 {
		n = 32
	}
	hashBytes := sha256.Sum256([]byte(strconv.Itoa(int(time.Now().UnixNano()))))
	return fmt.Sprintf("%x", hashBytes[:n])
}

//
// InitContainerInfo
// @Description: 记录一个容器的信息
// @param cID
// @param cPID
// @param cmdArray
// @param containerName
// @param volume
// @param portMappings
// @return *record.ContainerInfo
// @return error
//
func InitContainerInfo(cID string, cPID int, cmdArray []string, containerName, volume string, portMappings []string) (*record.ContainerInfo, error) {
	createTime := time.Now().Format("2006-01-02 15:04:05")
	if containerName == "" {
		containerName = cID
	}
	containerInfo := &record.ContainerInfo{
		Pid:         strconv.Itoa(cPID),
		Id:          cID,
		Name:        containerName,
		Command:     strings.Join(cmdArray, ""),
		Volume:      volume,
		CreateTime:  createTime,
		Status:      record.RUNNING,
		PortMapping: portMappings,
	}
	//序列化info
	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("convert container info to json type error:%v", err)
		return nil, err
	}
	infoFilePath := path.Join(record.DefaultInfoLocation, cID)
	if err := os.MkdirAll(infoFilePath, 0622); err != nil {
		log.Errorf("create info file dir error:%v", err)
		return nil, err
	}
	//创建json文件
	jsonFileName := path.Join(infoFilePath, record.ConfigName)
	configFile, err := os.OpenFile(jsonFileName, syscall.O_CREAT|syscall.O_RDWR|syscall.O_APPEND, 0644)
	if err != nil {
		log.Errorf("create info file error:%v", err)
		return nil, err
	}
	defer configFile.Close()
	//写入文件
	if _, err := configFile.Write(jsonBytes); err != nil {
		log.Errorf("write container info to config file error:%v", err)
		return nil, err
	}
	return containerInfo, nil
}

//
// DeleteContainerInfo
// @Description: 删除cID的容器信息
// @param cID
//
func DeleteContainerInfo(cID string) {
	infoFilePath := path.Join(record.DefaultInfoLocation, cID)
	if err := os.RemoveAll(infoFilePath); err != nil {
		log.Errorf("delete [%s] container info error:%v", err)
		return
	}
}
