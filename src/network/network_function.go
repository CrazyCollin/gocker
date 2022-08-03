package network

import (
	log "github.com/sirupsen/logrus"
	"gocker/src/record"
	"io/fs"
	"net"
	"os"
	"path"
	"path/filepath"
)

//
// Init
// @Description: 初始化网络配置
//
func Init() error {
	//加载bridge驱动
	var bridgeDriver = &BridgeNetworkDriver{}
	record.Drivers[bridgeDriver.Name()] = bridgeDriver
	//检测网络配置文件目录
	if _, err := os.Stat(record.DefaultNetworkPath); err != nil {
		//不存在网络配置文件目录就创建
		if os.IsNotExist(err) {
			if err := os.MkdirAll(record.DefaultNetworkPath, 0644); err != nil {
				log.Error(err)
				return err
			}
		} else {
			log.Error(err)
			return err
		}
	}
	if err := filepath.Walk(record.DefaultNetworkPath, func(networkPath string, info fs.FileInfo, err error) error {
		// 如果是目录则跳过
		if info.IsDir() {
			return nil
		}
		// 加载文件名作为网络名
		_, networkName := path.Split(networkPath)
		network := &Network{
			Name: networkName,
		}
		// 调用Network.load方法加载网络配置信息
		if err := network.load(networkPath); err != nil {
			log.Error(err)
			return err
		}
		// 将网络配置信息加入到networks字典中
		record.Networks[networkName] = network
		return nil
	}); err != nil {
		return err
	}
	return nil
}

//
// CreateNetwork
// @Description: 创建容器网络
//
func CreateNetwork(driver, subnet, name string) error {
	_, cidr, _ := net.ParseCIDR(subnet)
	//分配gateway的ip地址
	gatewayIP, err := allocator.Allocate(cidr)
	if err != nil {
		log.Error(err)
		return err
	}
	cidr.IP = gatewayIP

	//调用网络驱动创建网络
	network, err := record.Drivers[driver].Create(cidr.String(), name)
	if err != nil {
		log.Error(err)
		return err
	}
	return network.dump(record.DefaultNetworkPath)
}

//
// Connect
// @Description: 容器连接网络
//
func Connect(networkName string, containerInfo *record.ContainerInfo) error {
	return nil
}

//
// ListNetwork
// @Description: 遍历network
//
func ListNetwork() {

}

//
// DeleteNetwork
// @Description: 删除指定network
//
func DeleteNetwork(networkName string) error {
	return nil
}
