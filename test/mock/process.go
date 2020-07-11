package mock

// Process mock implementation of a process.
type Process struct {
	FindPIDByGivenPortNumberFnc func(port int32) (int32, error)
	FindProcessByNameFnc func(name string) ([]int32, error)
	KillProcessFnc func(pid int32) error
}

// FindPIDByGivenPortNumber mock implementation
func (p Process) FindPIDByGivenPortNumber(port int32) (int32, error) {
	if p.FindPIDByGivenPortNumberFnc != nil {
		return p.FindPIDByGivenPortNumberFnc(port)
	}
	return 0, nil
}

// FindProcessByName mock implementation
func (p Process) FindProcessByName(name string) ([]int32, error) {
	if p.FindProcessByNameFnc != nil {
		return p.FindProcessByNameFnc(name)
	}
	return []int32{}, nil
}

// KillProcess mock implementation
func (p Process) KillProcess(pid int32) error {
	if p.KillProcessFnc != nil {
		return p.KillProcessFnc(pid)
	}
	return nil
}
