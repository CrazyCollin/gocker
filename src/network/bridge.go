package network

import "net"

type BridgeNetworkDriver struct {
}

func (b *BridgeNetworkDriver) initBridge(network *Network) error {
	return nil
}

func (b *BridgeNetworkDriver) Name() string {
	return "bridge"
}

func (b *BridgeNetworkDriver) Create(subnet, name string) (*Network, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BridgeNetworkDriver) Delete(network *Network) error {
	//TODO implement me
	panic("implement me")
}

func (b *BridgeNetworkDriver) Connect(network *Network, endpoint *Endpoint) error {
	//TODO implement me
	panic("implement me")
}

func (b *BridgeNetworkDriver) Disconnect(network *Network, endpoint *Endpoint) error {
	//TODO implement me
	panic("implement me")
}

//
// createBridgeInterface
// @Description: 创建一个bridge driver/virtual dev
// @param bridge
// @return error
//
func createBridgeInterface(bridge string) error {
	return nil
}

//
// setInterfaceIP
// @Description: 设置bridge的地址和路由
// @param name
// @param rawIP
// @return error
//
func setInterfaceIP(name, rawIP string) error {
	return nil
}

//
// setInterfaceUP
// @Description: 设置bridge为启动状态
// @param interfaceName
// @return error
//
func setInterfaceUP(interfaceName string) error {
	return nil
}

func setupIPTables(bridgeName, subnet *net.IPNet) error {
	return nil
}
