package record

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
	DefaultInfoLocation = "/var/run/gocker/"
	RootURL             = "/var/lib/gocker/aufs/"
)

const (
	EnvExecPid = "gocker_pid"
	EnvExecCmd = "gocker_cmd"
)

var (
	ConfigName  = "containerInfo.json"
	LogFileName = "container.log"
)
