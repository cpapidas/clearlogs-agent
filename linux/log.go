package linux

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

// Log is responsible to describe the log actions we do
// in order to manage the messages from the  processes.
type Log struct {
	StdErr string
	StdOut string
}

// GetLogFromProcess is responsible to get logs from a specific process pid.
// The function gets the process pid as int parameter and two channels the
// message chan which contain the logs per line and the err channel which contains the
// errors per log.
//
// This is a linux implementation, the function will open the std named pipe and
// read all the content from there.
func (l Log) GetLogFromProcess(processPID int, message chan <-string, errChan chan <-error) {
	stdErr := fmt.Sprintf("/proc/%d/fd/1", processPID)
	stdOut := fmt.Sprintf("/proc/%d/fd/2", processPID)

	if l.StdErr != "" {
		stdErr = l.StdErr
	}
	if l.StdOut != "" {
		stdOut = l.StdOut
	}

	fileStdErr, err := os.OpenFile(stdErr, os.O_RDONLY, 0600)
	if err != nil {
		errChan <- fmt.Errorf("open named pipe file %s error: %v", stdErr, err)
		return
	}
	fileStdOut, err := os.OpenFile(stdOut, os.O_RDONLY, 0600)
	if err != nil {
		errChan <- fmt.Errorf("open named pipe file %s error: %v", stdOut, err)
		return
	}

	readerOut := bufio.NewReader(fileStdErr)
	readerErr := bufio.NewReader(fileStdOut)

	for {
		lineOut, err := readerOut.ReadBytes('\n')
		if err != nil {
			log.Printf("failed to readying bytes from %s: with err %v", stdOut, err)
			continue
		}
		message <- string(lineOut)
		lineErr, err := readerErr.ReadBytes('\n')
		if err != nil {
			log.Printf("failed to readying bytes from %s: with err %v", stdErr, err)
			continue
		}
		message <- string(lineErr)

		// We want to avoid CPU overheating that's why we have to add a time sleep here
		// read more at https://stackoverflow.com/questions/45443414/read-continuously-from-a-named-pipe
		time.Sleep(100 * time.Millisecond)
	}
}