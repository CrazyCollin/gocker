package network

import (
	"gocker/src/record"
	"net"
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

var ipAllocator = &IPAM{SubnetAllocatorPath: record.DefaultSubnetAllocatorPath}

//
// load
// @Description: 加载ipam的位图信息
// @receiver ipam
// @return error
//
func (ipam *IPAM) load() error {
	return nil
}

//
// dump
// @Description: 存储ipam的地址分配位图信息
// @receiver ipam
// @return error
//
func (ipam *IPAM) dump() error {
	return nil
}

//
// Allocate
// @Description: 使用bitmap分配一个ip地址（指定网段）
// @receiver ipam
// @param subnet
// @return ip
// @return err
//
func (ipam *IPAM) Allocate(subnet *net.IPNet) (ip net.IP, err error) {
	return nil, err
}

//
// Release
// @Description: 释放一个ip
// @receiver ipam
// @param subnet
// @param ipAddr
// @return error
//
func (ipam *IPAM) Release(subnet *net.IPNet, ipAddr net.IP) error {
	return nil
}
