package subsystem

//
// ResourceConfig
// @Description: 资源配置，对内存cpu等的具体限制
//
type ResourceConfig struct {
	MemoryLimit string
	CpuSet      string
	CpuShare    string
}

//
// Subsystem
// @Description: 资源限制接口
//
type Subsystem interface {
	Name() string
	Set(cgroupName string, res *ResourceConfig) error
	Apply(cgroupName string, pid int) error
	Remove(cgroupName string) error
}

var SubsystemIns = []Subsystem{
	&CpuSetSubsystem{},
	&CpuShareSubsystem{},
	&MemorySubsystem{},
}
