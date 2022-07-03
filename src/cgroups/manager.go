package cgroups

import "gocker/src/cgroups/subsystem"

//
// CgroupManager
// @Description:
//
type CgroupManager struct {
	CgroupName string
	Resource   *subsystem.ResourceConfig
}

func NewCgroupManager(cgroupName string) *CgroupManager {
	return &CgroupManager{CgroupName: cgroupName}
}

//
// Apply
// @Description: 将pid注册进cgroup
// @receiver cm
// @param pid
// @return error
//
func (cm *CgroupManager) Apply(pid int) error {
	for _, ins := range subsystem.SubsystemIns {
		err := ins.Apply(cm.CgroupName, pid)
		if err != nil {
			return err
		}
	}
	return nil
}

//
// Set
// @Description: 设置限制
// @receiver cm
// @param res
// @return error
//
func (cm *CgroupManager) Set(res *subsystem.ResourceConfig) error {
	for _, ins := range subsystem.SubsystemIns {
		err := ins.Set(cm.CgroupName, res)
		if err != nil {
			return err
		}
	}
	return nil
}

//
// Destroy
// @Description: 释放group
// @receiver cm
// @return error
//
func (cm *CgroupManager) Destroy() error {
	for _, ins := range subsystem.SubsystemIns {
		err := ins.Remove(cm.CgroupName)
		if err != nil {
			return err
		}
	}
	return nil
}
