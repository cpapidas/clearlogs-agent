// +build !windows,!android, !zos

package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/cpapidas/clagent"
	"github.com/cpapidas/clagent/linux"
	"github.com/cpapidas/clagent/mac"
	"github.com/cpapidas/clagent/net"
	"github.com/cpapidas/clagent/process"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

func main() {
	log.Println("Setting up the program flags")
	// Set up the program flags/parameters
	token, port, pid := setUpFlags()

	log.Println("Creating the config object")
	// Create the config file according to program flags
	conf := setUpConfig(token, port, pid)

	log.Println("Validating the config object")
	// Validate the config object and see if we have
	// all the data in order to continue
	err := validateConfig(conf)
	if err != nil {
		log.Fatalf("the application does not have all the nesecery configurations, error: %v", err)
	}

	// Create a process object, used for the dependency injection
	pro := process.Process{}

	log.Println("Checking for similar processes")
	err = killSimilarProcess(pro)
	if err != nil {
		log.Fatalf("failed to kill simalar process with error: %v", err)
	}

	// According to client operating system return the proper
	// provider to handle the logs.
	log.Println("Find the log provider for the current OS")
	lg, err := findTheLogProvider()
	if err != nil {
		log.Fatalf("failed to create log provider with error: %v", err)
	}

	stop := make(chan bool, 1)

	httpclient := net.NewHTTPClient("", &http.Client{Timeout: 5 * time.Second})

	// Start listen to a specific pid and send the data to the server.
	log.Println("Listening for process logs")
	err = clagent.ListenToPid(conf, pro, lg, stop, httpclient, conf.Token)
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
		return mac.Log{ShouldUseSudo: true}, nil
	case "plan9", "linux", "dragonfly", "freebsd", "hurd", "illumos", "solaris", "nacl", "netbsd", "openbsd":
		return linux.Log{}, nil
	case "zos":
		return nil, errors.New("cannot find supported operation system: zos")
	case "windows":
		return nil, errors.New("cannot find supported operation system: window")
	case "android":
		return nil, errors.New("cannot find supported operation system: android")
	default:
		return nil, errors.New("cannot find supported operation system: unknow")
	}
}

// killSimilarProcess checks if we have another clagent process that is already running,
// if yes then we need to have only one process up in order to avoid
// duplicated data in server. Kill the old process and continue.
func killSimilarProcess(pro process.Process) error {
	cpids, err := pro.FindProcessByName(process.AgentName)
	if err != nil {
		return fmt.Errorf("failed to search for a process by name with error: %v", err)
	}
	for _, cpid := range cpids {
		if cpid != 0 && cpid != int32(os.Getpid()) {
			log.Println("a similar process already is running, killing the old process")
			err := pro.KillProcess(cpid)
			if err != nil {
				return fmt.Errorf("failed to kill the old process with error: %v", err)
			}
		}
	}
	return nil
}
