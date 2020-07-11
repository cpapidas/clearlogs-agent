package process

import (
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

// AgentName defines the program's process name.
const AgentName string = "clagent"

// Process defines the process implementation
type Process struct {}

// FindProcessByName is responsible to find the the process
// by name and return the pid. The function will return
// an error if something occurred.
func (Process) FindProcessByName(name string) ([]int32, error) {
	var results []int32

	// Get all processes
	pp, err := process.Processes()
	if err != nil {
		return results, fmt.Errorf("failed to get all processes with err: %v", err)
	}
	for _, p := range pp {
		pName, err := p.Name()
		if err != nil {
			return results, fmt.Errorf("failed to get process name with error: %v", err)
		}
		if pName == name {
			results := append(results, p.Pid)
			return results, nil
		}

		// Get all children processes.
		cc, err := p.Children()
		if err != nil {
			// could not find the children processes for this pid
			continue
		}
		for _, c := range cc {
			cName, err := c.Name()
			if err != nil {
				return results, fmt.Errorf("failed to get name from child with error: %v", err)
			}
			if cName == name {
				results := append(results, c.Pid)
				return results, nil
			}
		}
	}

	return results, nil
}

// KillProcess is responsible to kill a process, the
// function will return an error if something occurred.
func (Process) KillProcess(pid int32) error {
	ff := process.Process{Pid: pid}
	err := ff.Kill()
	if err != nil {
		return fmt.Errorf("failed to kill a process with error: %v", err)
	}
	return nil
}

// FindPIDByGivenPortNumber will parse all the available
// processes in order to find the pid which is listening to
// given port.
//
// The function will return the pid as int32 or an error
// if something occur.
func (pro Process)  FindPIDByGivenPortNumber(port int32) (int32, error) {
	// Get all main process
	pp, err := process.Processes()
	if err != nil {
		return 0, fmt.Errorf("failed to get all processes with err: %v", err)
	}
	for _, p := range pp {
		pid, err := pro.checkProcessForPort(p, port)
		if err != nil {
			return 0, fmt.Errorf("failed to get the process id by port with error: %v", err)
		}
		if pid != 0 {
			return pid, nil
		}

		// Check also the children processes
		dd, err := p.Children()
		if err != nil {
			// could not find the children processes for this pid
			continue
		}
		for _, d := range dd {
			pid, err := pro.checkProcessForPort(d, port)
			if err != nil {
				return 0, fmt.Errorf("failed to get the process id by port with error: %v", err)
			}
			if pid != 0 {
				return pid, nil
			}
		}
	}
	return 0, errors.New(fmt.Sprintf("failed to find any process for port: %d", port))
}

// checkProcessForPort is responsible to check the given process
// for the port number and return the process id or an error object
// if something occurred.
func (Process) checkProcessForPort(p *process.Process, port int32) (int32, error) {
	cc, err := net.ConnectionsPid("", p.Pid)
	if err != nil {
		return 0, fmt.Errorf("failed to get the net data for pid: %d with err: %v", p.Pid, err)
	}
	for _, c := range cc {
		if c.Laddr.Port == uint32(port) {
			return c.Pid, nil
		}
	}

	return 0, nil
}