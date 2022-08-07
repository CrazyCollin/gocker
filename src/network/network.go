package network

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"gocker/src/record"
	"io/ioutil"
	"net"
	"os"
	"path"
)

//
// Network
// @Description:容器网络集合，此网络上容器可以互相通信
//
type Network struct {
	Name    string     `json:"name"`
	IpRange *net.IPNet `json:"ip_range"`
	Driver  string
}

//
// Endpoint
// @Description:网络端点，用于连接容器和网络，保证容器内部和网络的通信
//
type Endpoint struct {
	//ID
	ID string `json:"id"`
	//veth设备
	Device netlink.Device `json:"dev"`
	//ip地址
	IpAddress net.IP `json:"ip"`
	//mac地址
	MacAddress net.HardwareAddr `json:"mac"`
	//端口映射
	PortMapping []string `json:"port_mapping"`
	//网络
	Network *Network
}

type NetworkDriver interface {
	//驱动名
	Name() string
	//创建网络
	Create(subnet, name string) (*Network, error)
	//删除网络
	Delete(network *Network) error
	//将指定容器网络端点连接至网络中
	Connect(network *Network, endpoint *Endpoint) error
	//将容器网络端点从网络中删除
	Disconnect(network *Network, endpoint *Endpoint) error
}

//
// configEndpointIpaddressAndRoute
// @Description: 设置容器网络设备
//
func configEndpointIpaddressAndRoute(endpoint *Endpoint, info *record.ContainerInfo) error {
	return nil
}

//
// enterContainerNetNamespace
// @Description: 进入容器设置veth
//
func enterContainerNetNamespace(link *netlink.Link, info *record.ContainerInfo) func() {
	return nil
}

//
// configPortMapping
// @Description: 配置端口映射
//
func configPortMapping(endpoint *Endpoint) error {
	return nil
}

//
// dump
// @Description: 保存网络配置文件
//
func (nw *Network) dump(dumpPath string) error {
	// 检查保存的目录是否存在
	if _, err := os.Stat(dumpPath); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dumpPath, 0644); err != nil {
				log.Error(err)
				return err
			}
		} else {
			log.Error(err)
			return err
		}
	}
	// 保存的文件名使用网络的名字
	nwPath := path.Join(dumpPath, nw.Name)
	// 打开保存的文件用于写入
	file, err := os.OpenFile(nwPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Error(err)
		return err
	}
	defer file.Close()

	// 通过json序列化
	nwBytes, err := json.Marshal(nw)
	if err != nil {
		log.Error(err)
		return err
	}
	// 写入
	if _, err := file.Write(nwBytes); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//
// load
// @Description: 加载网络配置文件
//
func (nw *Network) load(dumpPath string) error {
	// 打开配置文件
	nwConfigFile, err := os.Open(dumpPath)
	if err != nil {
		log.Error(err)
		return err
	}
	defer nwConfigFile.Close()
	// 从配置文件中读取网络的配置json
	jsonBytes, err := ioutil.ReadAll(nwConfigFile)
	if err != nil {
		log.Error(err)
		return err
	}
	// 反序列化
	if err := json.Unmarshal(jsonBytes, &nw); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//
// remove
// @Description: 删除网络配置文件
//
func (nw *Network) remove(dumpPath string) error {
	if _, err := os.Stat(path.Join(dumpPath, nw.Name)); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	} else {
		return os.Remove(path.Join(dumpPath, nw.Name))
	}
}
