package subsystem

type CpuShareSubsystem struct {
}

func (c CpuShareSubsystem) Name() string {
	//TODO implement me
	panic("implement me")
}

func (c CpuShareSubsystem) Set(cgroupName string, res *ResourceConfig) error {
	//TODO implement me
	panic("implement me")
}

func (c CpuShareSubsystem) Apply(cgroupName string, pid int) error {
	//TODO implement me
	panic("implement me")
}

func (c CpuShareSubsystem) Remove(cgroupName string) error {
	//TODO implement me
	panic("implement me")
}
