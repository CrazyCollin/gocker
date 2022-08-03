package network

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"net"
	"os/exec"
	"strings"
)

//
// BridgeNetworkDriver
// @Description: 网络驱动
//
type BridgeNetworkDriver struct {
}

//
// initBridge
// @Description: 初始化网桥
//
func (b *BridgeNetworkDriver) initBridge(network *Network) error {
	bridgeName := network.Name
	//1.创建虚拟bridge设备
	if err := createBridgeInterface(bridgeName); err != nil {
		return fmt.Errorf("create bridge error:%v", err)
	}
	//2.配置bridge的地址和路由
	gatewayIP := *network.IpRange
	gatewayIP.IP = network.IpRange.IP
	if err := setBridgeIP(bridgeName, gatewayIP.String()); err != nil {
		return fmt.Errorf("config bridge error:%v", err)
	}
	//3.启动bridge
	if err := setBridgeUP(bridgeName); err != nil {
		return fmt.Errorf("start bridge error:%v", err)
	}
	//4.设置iptables的SNAT规则
	if err := setupIPTables(bridgeName, network.IpRange); err != nil {
		return fmt.Errorf("setup bridge iptables error:%v", err)
	}
	return nil
}

func (b *BridgeNetworkDriver) Name() string {
	return "bridge"
}

func (b *BridgeNetworkDriver) Create(subnet, name string) (*Network, error) {
	//TODO implement me
	panic("implement me")
}

//
// Delete
// @Description: 删除对应的network，即删除指定bridge
//
func (b *BridgeNetworkDriver) Delete(network *Network) error {
	bridgeName := network.Name
	link, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return fmt.Errorf("get link error:%v", err)
	}
	return netlink.LinkDel(link)
}

func (b *BridgeNetworkDriver) Connect(network *Network, endpoint *Endpoint) error {

	return nil
}

func (b *BridgeNetworkDriver) Disconnect(network *Network, endpoint *Endpoint) error {
	//TODO implement me
	panic("implement me")
}

//
// createBridgeInterface
// @Description: 创建一个bridge driver/virtual dev
//
func createBridgeInterface(bridgeName string) error {
	iface, err := net.InterfaceByName(bridgeName)
	//存在同名bridge
	if iface != nil || err == nil {
		return fmt.Errorf("exist same bridge:%s", iface.Name)
	}
	//初始化一个netlink的link基础对象
	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = bridgeName
	//创建bridge
	bridge := &netlink.Bridge{
		LinkAttrs: linkAttrs,
	}
	//创建bridge虚拟网络设备
	if err := netlink.LinkAdd(bridge); err != nil {
		return fmt.Errorf("create virtual bridge error:%v", err)
	}
	return nil
}

//
// setInterfaceIP
// @Description: 设置bridge的地址和路由
//
func setBridgeIP(bridgeName, rawIP string) error {
	link, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return fmt.Errorf("get link error:%v", err)
	}
	ipNet, err := netlink.ParseIPNet(rawIP)
	if err != nil {
		return err
	}
	//ip addr add xxx
	//为bridge配置ip地址
	addr := &netlink.Addr{
		IPNet: ipNet,
		Label: "",
		Flags: 0,
		Scope: 0,
	}
	return netlink.AddrAdd(link, addr)
}

//
// setInterfaceUP
// @Description: 设置bridge为启动状态
//
func setBridgeUP(bridgeName string) error {
	link, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return fmt.Errorf("get link error:%v", link)
	}
	//启动bridge
	if err := netlink.LinkSetUp(link); err != nil {
		return fmt.Errorf("enable bridge-%s- up error:%v", bridgeName, err)
	}
	return nil
}

//
// setupIPTables
// @Description: 设置特定bridge的MASQUERADE规则
//
func setupIPTables(bridgeName string, subnet *net.IPNet) error {
	iptablesCMD := fmt.Sprintf("-t nat -A POSTROUTING -s %s ! -o %s -j MASQUERADE", subnet.String(), bridgeName)
	cmd := exec.Command("iptables", strings.Split(iptablesCMD, " ")...)
	output, err := cmd.Output()
	if err != nil {
		log.Errorf("iptables output:%v", output)
	}
	return nil
}
