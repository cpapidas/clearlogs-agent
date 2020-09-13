package mock

// Log mock implementation of a process.
type Log struct {
	GetLogFromProcessFnc func(input string, message chan<- string, errChan chan<- error)
}

// GetLogFromProcess mock implementation.
func (l Log) GetLogFromProcess(input string, message chan<- string, errChan chan<- error) {
	if l.GetLogFromProcessFnc != nil {
		l.GetLogFromProcessFnc(input, message, errChan)
	}
	message <- ""
}
