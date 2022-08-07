package network

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gocker/src/record"
	"io/fs"
	"net"
	"os"
	"path"
	"path/filepath"
	"text/tabwriter"
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
	network, ok := record.Networks[networkName]
	if !ok {
		err := fmt.Errorf("no such network %s", networkName)
		log.Error(err)
		return err
	}
	//分配给容器一个ip
	containerIP, err := allocator.Allocate(network.IpRange)
	if err != nil {
		log.Error(err)
		return err
	}
	endpoint := &Endpoint{
		ID:          fmt.Sprintf("%s-%s", containerInfo.Id, networkName),
		IpAddress:   containerIP,
		PortMapping: containerInfo.PortMapping,
		Network:     network,
	}
	//调用网络驱动连接endpoint与network
	if err := record.Drivers[network.Driver].Connect(network, endpoint); err != nil {
		log.Error(err)
		return err
	}
	//进入容器net namespace配置容器网络设备的ip地址和路由
	if err := configEndpointIpaddressAndRoute(endpoint, containerInfo); err != nil {
		log.Error(err)
		return err
	}
	return configPortMapping(endpoint)
}

//
// ListNetwork
// @Description: 遍历network
//
func ListNetwork() {
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "NAME\tIpRange\tDriver\n")
	for _, v := range record.Networks {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			v.Name,
			v.IpRange.String(),
			v.Driver,
		)
	}
	if err := w.Flush(); err != nil {
		log.Error(err)
		return
	}
}

//
// DeleteNetwork
// @Description: 删除指定network
//
func DeleteNetwork(networkName string) error {
	nw, ok := record.Networks[networkName]
	if !ok {
		err := fmt.Errorf(" no such network: %s", networkName)
		log.Error(err)
		return err
	}
	// 调用IPAM的实例释放网络网关的IP
	if err := allocator.Release(nw.IpRange, &nw.IpRange.IP); err != nil {
		return fmt.Errorf(" error remove network gateway ip: %s", err)
	}
	// 调用网络驱动删除网络创建的设备与配置
	if err := record.Drivers[nw.Driver].Delete(nw); err != nil {
		return fmt.Errorf(" Error Remove Network DriverError: %s", err)
	}
	// 从网络的配置目录中删除该网络对应的配置文件
	return nw.remove(record.DefaultNetworkPath)
}
