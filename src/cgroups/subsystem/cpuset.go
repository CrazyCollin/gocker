package subsystem

type CpuSetSubsystem struct {
}

func (c CpuSetSubsystem) Name() string {
	//TODO implement me
	panic("implement me")
}

func (c CpuSetSubsystem) Set(cgroupName string, res *ResourceConfig) error {
	//TODO implement me
	panic("implement me")
}

func (c CpuSetSubsystem) Apply(cgroupName string, pid int) error {
	//TODO implement me
	panic("implement me")
}

func (c CpuSetSubsystem) Remove(cgroupName string) error {
	//TODO implement me
	panic("implement me")
}
