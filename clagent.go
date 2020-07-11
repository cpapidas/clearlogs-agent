package clagent

import (
	"fmt"
	"log"
)

// Process describes the process handler
type Process interface {
	// FindPIDByGivenPortNumberType describe the type of function
	// which is responsible to return the pid by given port number.
	FindPIDByGivenPortNumber(port int32) (int32, error)

	// FindProcessByName is responsible to find the the process
	// by name and return the pid. The function will return
	// an error if something occurred.
	FindProcessByName(name string) ([]int32, error)

	// KillProcess is responsible to kill a process, the
	// function will return an error if something occurred.
	KillProcess(pid int32) error
}

// Log is responsible to describe the action we do in the
// log level. For example when we want to get the logs
// from a process.
type Log interface {
	// GetLogFromProcessType describes the function which will get the logs
	// from a running process. It takes the process pid as first argument
	// the message chan which is the logs per line and an error chanel in case
	// something happen.
	GetLogFromProcess(processPID int, message chan<- string, errChan chan<- error)
}

// Config describe all the properties for the program.
//
// Token describe the server token in order to send the data
// Pid is an optional parameter if Port is set, if skipped
// 	then the program will use the Port to find the process pid
// Pod is an optional parameter if Pid is set
type Config struct {
	Token string
	Pid   int32
	Port  int32
}

// ListenToPid function is responsible to find the pid and listen
// for logs. The function will send the log to the service by
// the given implementation (tcp sockets).
//
// If for some reason the pid in the config file is null the
// program will try to find the pid by the given port.
func ListenToPid(
	conf Config,
	proc Process,
	lg Log,
	stop <-chan bool,
) error {
	pid, err := getPid(conf, proc)
	if err != nil {
		return err
	}
	message := make(chan string)
	errChan := make(chan error)

	go lg.GetLogFromProcess(int(pid), message, errChan)

	// we need an infinity loop in order to stay connected and get all
	// the messages from the process.
	for {
		select {
		case mess := <-message:
			_ = mess
			fmt.Println(mess)
		case errc := <-errChan:
			log.Printf("error on reading messages, error: %v", errc)
		case <-stop:
			return nil
		}
	}
}

// getPid function is responsible to return the pid from config file,
// if this is not exist then will try to find the pid from the given
// port.
func getPid(conf Config, proc Process) (int32, error) {
	if conf.Pid != 0 {
		return conf.Pid, nil
	}

	pid, err := proc.FindPIDByGivenPortNumber(conf.Port)

	if err != nil {
		return 0, fmt.Errorf("failed to get pid by port number with err: %v", err)
	}

	return pid, nil
}
