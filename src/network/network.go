package network

import (
	"github.com/vishvananda/netlink"
	"gocker/src/record"
	"net"
)

//
// Network
// @Description:容器网络集合，此网络上容器可以h互相通信
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
	ID          string           `json:"id"`
	Device      netlink.Device   `json:"dev"`
	IpAddress   net.IP           `json:"ip"`
	MacAddress  net.HardwareAddr `json:"mac"`
	PortMapping []string         `json:"port_mapping"`
	Network     *Network
}

type NetworkDriver interface {
	Name() string
	Create(subnet, name string) (*Network, error)
	Delete(network *Network) error
	Connect(network *Network, endpoint *Endpoint) error
	Disconnect(network *Network, endpoint *Endpoint) error
}

//
// Init
// @Description: 从配置中加载网络配置信息
// @return error
//
func Init() error {
	//加载bridge驱动
	var bridgeDriver = &BridgeNetworkDriver{}
	record.Drivers[bridgeDriver.Name()] = bridgeDriver

	return nil
}

//
// CreateNetwork
// @Description: 创建网络
// @param driver
// @param subnet
// @param name
// @return error
//
func CreateNetwork(driver, subnet, name string) error {
	return nil
}

//
// Connect
// @Description: 容器连接网络
// @param networkName
// @param containerInfo
// @return error
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
// @Description: 删除一个network
// @param networkName
// @return error
//
func DeleteNetwork(networkName string) error {
	return nil
}

func configEndpointIpaddressAndRoute(link *netlink.Link, info *record.ContainerInfo) error {
	return nil
}

//
// enterContainerNetNamespace
// @Description: 进入容器设置veth
// @param link
// @param info
// @return func()
//
func enterContainerNetNamespace(link *netlink.Link, info *record.ContainerInfo) func() {
	return nil
}

//
// configPortMapping
// @Description: 配置端口映射
// @param endpoint
// @return error
//
func configPortMapping(endpoint *Endpoint) error {
	return nil
}

//
// dump
// @Description: 保存网络配置文件
// @receiver nw
// @param dumpPath
// @return error
//
func (nw *Network) dump(dumpPath string) error {
	return nil
}

//
// load
// @Description: 加载网络配置文件
// @receiver nw
// @param dumpPath
// @return error
//
func (nw *Network) load(dumpPath string) error {
	return nil
}

//
// remove
// @Description: 删除网络配置文件
// @receiver nw
// @param dumpPath
// @return error
//
func (nw *Network) remove(dumpPath string) error {
	return nil
}
