package network

import (
	"github.com/vishvananda/netlink"
	"gocker/src/record"
	"net"
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

func configEndpointIpaddressAndRoute(link *netlink.Link, info *record.ContainerInfo) error {
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
	return nil
}

//
// load
// @Description: 加载网络配置文件
//
func (nw *Network) load(dumpPath string) error {
	return nil
}

//
// remove
// @Description: 删除网络配置文件
//
func (nw *Network) remove(dumpPath string) error {
	return nil
}
