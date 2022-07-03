package container

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gocker/src/record"
	"gocker/src/utils"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"
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
// ListContainerInfo
// @Description: 列出所有容器的基本信息
//
func ListContainerInfo() {
	filesPath := path.Join(record.DefaultInfoLocation)
	infoFiles, err := ioutil.ReadDir(filesPath)
	if err != nil {
		log.Errorf("read info file dir error:%v", err)
		return
	}
	var containersInfo []*record.ContainerInfo
	//遍历所有容器的元数据信息
	for _, file := range infoFiles {
		if file.Name() == "network" {
			continue
		}
		//获取单个容器的信息
		containerInfo, err := getContainerInfoFromFile(file)
		if err != nil {
			log.Errorf("get info from one container error:%v", err)
		}
		containersInfo = append(containersInfo, containerInfo)
	}
	//输出所有容器信息
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATE\n")
	for _, info := range containersInfo {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t\n", info.Id, info.Name, info.Pid, info.Status, info.Command, info.CreateTime)
	}
	if err := w.Flush(); err != nil {
		log.Errorf("display containers info error:%v", err)
		return
	}
}

//
// getContainerInfoFromFile
// @Description: 从info文件中获取单个容器的元数据
// @param file
// @return *record.ContainerInfo
// @return error
//
func getContainerInfoFromFile(file fs.FileInfo) (*record.ContainerInfo, error) {
	fileName := file.Name()
	filePath := path.Join(record.DefaultInfoLocation, fileName, record.ConfigName)
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("read container info file error:%v", err)
		return nil, err
	}
	containerInfo := &record.ContainerInfo{}
	err = json.Unmarshal(fileBytes, containerInfo)
	if err != nil {
		log.Errorf("unmarshal container info error:%v", err)
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

//
// exportContainerLogs
// @Description:创建容器日志文件，并将标准输出重定向至日志文件
// @param cID
// @param cmdOut
//
func exportContainerLogs(cID string, cmdOut *io.Writer) {
	logFilePath := path.Join(record.DefaultInfoLocation, cID)
	if has, err := utils.DirOrFileExist(logFilePath); err == nil && !has {
		if err := os.MkdirAll(logFilePath, 0622); err != nil {
			log.Errorf("create container logs dir error:%v", err)
			return
		}
	}
	logFilePath = path.Join(logFilePath, record.LogFileName)
	logFile, err := os.Create(logFilePath)
	if err != nil {
		log.Errorf("create container logs error:%v", err)
		return
	}
	//todo 不是很懂标准输入输出
	*cmdOut = logFile
}

//
// CheckLogsOfContainer
// @Description: 获取容器的日志信息
// @param cID
//
func CheckLogsOfContainer(cID string) {
	logFilePath := path.Join(record.DefaultInfoLocation, cID, record.LogFileName)
	logFile, err := os.Open(logFilePath)
	if err != nil {
		log.Errorf("open container logs file error:%v", err)
		return
	}
	defer logFile.Close()
	logs, err := ioutil.ReadAll(logFile)
	if err != nil {
		log.Errorf("read container logs error:%v", err)
		return
	}
	_, err = fmt.Fprint(os.Stdout, string(logs))
	if err != nil {
		log.Errorf("")
		return
	}
}
