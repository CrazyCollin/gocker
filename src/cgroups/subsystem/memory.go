package subsystem

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

type MemorySubsystem struct {
}

func (m MemorySubsystem) Name() string {
	return "memory"
}

func (m MemorySubsystem) Set(cgroupName string, res *ResourceConfig) error {
	//获取自定义cgroup路径，如没有则创建
	cgroupPath, err := GetCgroupPath(m.Name(), cgroupName)
	if err != nil {
		return err
	}
	log.Infof("%s cgroup path:%s", m.Name(), cgroupPath)
	//写入限制资源配置
	limitFilePath := path.Join(cgroupPath, "memory.limit_in_bytes")
	if err := ioutil.WriteFile(limitFilePath, []byte(res.MemoryLimit), 0644); err != nil {
		return fmt.Errorf("set memory cgroup failed:%v", err)
	}
	return nil
}

func (m MemorySubsystem) Apply(cgroupName string, pid int) error {
	cgroupPath, err := GetCgroupPath(m.Name(), cgroupName)
	if err != nil {
		return err
	}
	log.Infof("%s cgroup path:%s", m.Name(), cgroupPath)
	//将pid加入cgroup
	limitFilePath := path.Join(cgroupPath, "tasks")
	if err := ioutil.WriteFile(limitFilePath, []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("add pid to cgroup failed:%v", err)
	}
	return nil
}

func (m MemorySubsystem) Remove(cgroupName string) error {
	cgroupPath, err := GetCgroupPath(m.Name(), cgroupName)
	if err != nil {
		return err
	}
	log.Infof("%s cgroup path:%s", m.Name(), cgroupPath)
	return os.RemoveAll(cgroupPath)
}

func FindCgroupMountPoint(subsystem string) (string, error) {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return "", fmt.Errorf("open /proc/self/mountinfo error: %v", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		log.Debugf("mount info txt fields:%s", fields)
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4], err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("file scanner error:%s", err)
	}
	return "", fmt.Errorf("FindCgroupMountPoint is empty")
}

func GetCgroupPath(subsystemName, cgroupName string) (string, error) {
	cgroupRoot, err := FindCgroupMountPoint(subsystemName)
	if err != nil {
		return "", err
	}
	cgroupPath := path.Join(cgroupRoot, cgroupName)
	_, err = os.Stat(cgroupPath)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("file status error:%v", err)
	}
	if os.IsNotExist(err) {
		if err := os.Mkdir(cgroupPath, os.ModePerm); err != nil {
			return "", fmt.Errorf("mkdir error:%v", err)
		}
	}
	return cgroupPath, nil
}
