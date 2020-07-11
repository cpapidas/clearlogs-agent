// +build !windows,!android

package main

import (
	"errors"
	"flag"
	"github.com/cpapidas/clagent"
	"github.com/cpapidas/clagent/mac"
	"github.com/cpapidas/clagent/process"
	"log"
	"runtime"
)

func main() {
	// Set up the program flags/parameters
	token, port, pid := setUpFlags()

	// Create the config file according to program flags
	conf := setUpConfig(token, port, pid)

	// Validate the config object and see if we have
	// all the data in order to continue
	err := validateConfig(conf)
	if err != nil {
		log.Fatalf("the application does not have all the nesecery configurations, error: %v", err)
	}

	// Create a process object, used for the dependency injection
	pro := process.Process{}

	// Check if we have another clagent process that is already running,
	// if yes then we need to have only one process up in order to avoid
	// duplicated data in server. Kill the old process and continue.
	cpid, err := pro.FindProcessByName(process.AgentName)
	if err != nil {
		log.Fatalf("failed to search for a process by name with error: %v", err)
	}
	if cpid != 0 {
		log.Println("a similar process already is running, killing the old process")
		err := pro.KillProcess(cpid)
		if err != nil {
			log.Fatalf("failed to kill the old process with error: %v", err)
		}
	}

	// According to client operating system return the proper
	// provider to handle the logs.
	lg, err := findTheLogProvider()
	if err != nil {
		log.Fatalf("failed to create log provider with error: %v", err)
	}

	stop := make(chan bool, 1)

	// Start listen to a specific pid and send the data to the server.
	err = clagent.ListenToPid(conf, pro, lg, stop)
	if err != nil {
		log.Fatalf("application error: %v", err)
	}
}

// setUpFlags is responsible to setup the flags
func setUpFlags() (*string, *int, *int) {
	token := flag.String("token", "", "The client token")
	port := flag.Int("port", 0, "The process port")
	pid := flag.Int("pid", 0, "The process pid")

	flag.Parse()

	return token, port, pid
}

// setUpConfig is responsible to setup and return a domain Config object,
// with the given program parameters
func setUpConfig(token *string, port, pid *int) clagent.Config {
	return clagent.Config{
		Token: *token,
		Pid:   int32(*pid),
		Port:  int32(*port),
	}
}

// validateConfig is responsible to describe if the application has all
// the necessary information to start or not. The function returns an
// error in case config is wrong or something missing.
func validateConfig(conf clagent.Config) error {
	if conf.Token == "" {
		return errors.New("missing the Token property from config file")
	}
	if conf.Port == 0 && conf.Pid == 0 {
		return errors.New("port or pid should present in config file, both are missing")
	}
	return nil
}

// findTheLogProvider according to the user's system this function
// is responsible to find the proper implementation in order to manage
// the logs.
func findTheLogProvider() (clagent.Log, error) {
	switch runtime.GOOS {
	case "darwin":
		return mac.Log{}, nil
	default:
		return nil, errors.New("cannot find supported operation system")
	}
}