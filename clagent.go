package clagent

import (
	"fmt"
	"log"
	"strconv"
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
	// from a running process. It takes the input string as first argument
	// the message chan which is the logs per line and an error channel in case
	// something happen.
	GetLogFromProcess(input string, message chan<- string, errChan chan<- error)
}

// HTTPClient describe the actions over HTTP connection.
type HTTPClient interface {
	SendLogs(token, log string) error
}

// Config describe all the properties for the program.
//
// Token describe the server token in order to send the data
// Pid is an optional parameter if Port is set, if skipped
// 	then the program will use the Port to find the process pid
// Pod is an optional parameter if Pid is set
type Config struct {
	Token   string
	Pid     int32
	Port    int32
	BaseUrl string
	UseUDP  bool
	UDPAddress string
}

// ListenToPid function is responsible to find the pid and listen
// for logs. The function will send the log to the service by
// the given implementation (net sockets).
//
// If for some reason the pid in the config file is null the
// program will try to find the pid by the given port.
func ListenToPid(
	conf Config,
	proc Process,
	lg Log,
	stop <-chan bool,
	httpClient HTTPClient,
	token string,
) error {
	var (
		pid int32
		err error
	)
	if !conf.UseUDP {
		pid, err = getPid(conf, proc)
		if err != nil {
			return err
		}
	}

	log.Printf("find the process with pid: %d", pid)

	message := make(chan string)
	errChan := make(chan error)

	spid := strconv.Itoa(int(pid))
	go lg.GetLogFromProcess(spid, message, errChan)

	// we need an infinity loop in order to stay connected and get all
	// the messages from the process.
	for {
		//time.Sleep(1 * time.Second)
		log.Println("Checking for a message")
		select {
		case mess := <-message:
			log.Printf("%s\n\n", mess)
			//err := httpClient.SendLogs(token, mess)
			//if err != nil {
			//	log.Printf("error on send data: %v", err)
			//}
		case errc := <-errChan:
			log.Printf("error on reading messages, error: %v", errc)
		case <-stop:
			log.Println("stop the process")
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
