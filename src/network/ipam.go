package network

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gocker/src/record"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
)

//
// IPAM
// @Description: 网络功能中的一个组件，负责ip地址的分配和释放
//
type IPAM struct {
	//分配文件存放的位置
	SubnetAllocatorPath string
	//存储网段和位图
	Subnets map[string]string
}

var allocator = &IPAM{SubnetAllocatorPath: record.DefaultSubnetAllocatorPath}

//
// load
// @Description: 加载ipam的位图信息
//
func (ipam *IPAM) load() error {
	if _, err := os.Stat(record.DefaultSubnetAllocatorPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}
	configFile, err := os.Open(record.DefaultSubnetAllocatorPath)
	defer configFile.Close()
	if err != nil {
		return err
	}
	dataBytes, err := ioutil.ReadAll(configFile)
	err = json.Unmarshal(dataBytes, &ipam.Subnets)
	if err != nil {
		log.Errorf("unmarshal ipam config file error:%v", err)
		return err
	}
	return nil
}

//
// dump
// @Description: 存储ipam的地址分配位图信息
//
func (ipam *IPAM) dump() error {
	ipamConfigDir, _ := path.Split(record.DefaultSubnetAllocatorPath)
	if _, err := os.Stat(ipamConfigDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(ipamConfigDir, 0644); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	configFile, err := os.OpenFile(record.DefaultSubnetAllocatorPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	defer configFile.Close()
	if err != nil {
		return err
	}
	configJson, err := json.Marshal(ipam.Subnets)
	if err != nil {
		return err
	}
	_, err = configFile.Write(configJson)
	if err != nil {
		return err
	}
	return nil
}

//
// Allocate
// @Description: 使用bitmap分配一个ip地址（指定网段）
//
func (ipam *IPAM) Allocate(subnet *net.IPNet) (ip net.IP, err error) {
	ipam.Subnets = make(map[string]string)
	if err := ipam.load(); err != nil {
		log.Error(err)
	}
	//重新生成子网段实例
	_, subnet, _ = net.ParseCIDR(subnet.String())
	one, size := subnet.Mask.Size()

	ipAddr := subnet.String()
	//当前网段未分配，初始化网段的分段配置
	if _, has := ipam.Subnets[ipAddr]; !has {
		ipam.Subnets[ipAddr] = strings.Repeat("0", 1<<uint8(size-one))
	}

	var AllocateIP net.IP

	for c := range ipam.Subnets[ipAddr] {
		//找到未分配ip
		if ipam.Subnets[ipAddr][c] == '0' {
			ipAllocate := []byte(ipam.Subnets[ipAddr])
			ipAllocate[c] = '1'
			ipam.Subnets[ipAddr] = string(ipAllocate)
			resIP := subnet.IP
			for t := uint(4); t > 0; t -= 1 {
				[]byte(resIP)[4-t] += uint8(c >> ((t - 1) * 8))
			}
			resIP[3] += 1
			AllocateIP = resIP
			break
		}
	}
	if err := ipam.dump(); err != nil {
		return nil, err
	}
	return AllocateIP, err
}

//
// Release
// @Description: 释放一个ip
//
func (ipam *IPAM) Release(subnet *net.IPNet, ipAddr net.IP) error {
	ipam.Subnets = make(map[string]string)
	_, subnet, _ = net.ParseCIDR(subnet.String())
	//加载网段分配信息
	if err := ipam.load(); err != nil {
		log.Error(err)
		return err
	}
	c := 0
	releaseIP := ipAddr.To4()
	releaseIP[3] -= 1
	for t := uint(4); t > 0; t -= 1 {
		c += int(releaseIP[t-1]-subnet.IP[t-1]) << ((4 - t) * 8)
	}
	ipAlloc := []byte(ipam.Subnets[subnet.String()])
	ipAlloc[c] = '0'
	ipam.Subnets[subnet.String()] = string(ipAlloc)
	if err := ipam.dump(); err != nil {
		return err
	}
	return nil
}
