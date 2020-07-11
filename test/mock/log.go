package mock

// Log mock implementation of a process.
type Log struct {
	GetLogFromProcessFnc func(processPID int, message chan<- string, errChan chan<- error)
}

// GetLogFromProcess mock implementation.
func (l Log) GetLogFromProcess(processPID int, message chan<- string, errChan chan<- error) {
	if l.GetLogFromProcessFnc != nil {
		l.GetLogFromProcessFnc(processPID, message, errChan)
	}
	message <- ""
}
