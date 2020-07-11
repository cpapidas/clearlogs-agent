package mac

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
)

// Log is responsible to describe the log actions we do
// in order to manage the messages from the  processes.
type Log struct {}

// GetLogFromProcess is responsible to get logs from a specific process pid.
// The function gets the process pid as int parameter and two channels the
// message chan which contain the logs per line and the err channel which contains the
// errors per log.
//
// This is a mac implementation so, we are using dtrace command in order to
// have access to process logs.
func (Log) GetLogFromProcess(processPID int, message chan <-string, errChan chan <-error) {
	cmd := exec.Command( "dtrace",
		"-p", strconv.Itoa(processPID),
		"-qn", "syscall::write*:entry /pid == $target && arg0 == 1/ { printf(\"%s\", copyinstr(arg1, arg2));}",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errChan <- fmt.Errorf("error on stdout pipe, error: %v", err)
		return
	}

	err = cmd.Start()
	if err != nil {
		errChan <- fmt.Errorf("error on stdout pipe, error: %v", err)
		return
	}

	buf := bufio.NewReader(stdout)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			errChan <- err
		}
		message <- string(line)
	}
}
