package record

import "gocker/src/network"

type ContainerInfo struct {
	Pid         string   `json:"pid"`
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Command     string   `json:"command"`
	Volume      string   `json:"volume"`
	CreateTime  string   `json:"create_time"`
	Status      string   `json:"status"`
	PortMapping []string `json:"port_mapping"`
}

var (
	RUNNING = "running"
	STOP    = "stopped"
	EXIT    = "exited"
)

const (
	DefaultInfoLocation        = "/var/run/gocker/"
	RootURL                    = "/var/lib/gocker/aufs/"
	DefaultNetworkPath         = "/var/run/gocker/network/network/"
	DefaultSubnetAllocatorPath = "/var/run/gocker/network/ipam/subnet.json"
)

const (
	EnvExecPid = "gocker_pid"
	EnvExecCmd = "gocker_cmd"
)

var (
	ConfigName  = "containerInfo.json"
	LogFileName = "container.log"
)

var (
	// Drivers 网络驱动映射
	Drivers = map[string]network.NetworkDriver{}
	// Networks 所有网络映射
	Networks = map[string]*network.Network{}
)
