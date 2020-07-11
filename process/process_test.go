package process

import (
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"testing"
)

func TestFindPIDByGivenPortNumber(t *testing.T) {
	type args struct {
		port int32
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr bool
	}{
		{
			name:    "should return 0 with error for undefined pid",
			args:    args{port: 999999},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Process{}
			got, err := p.FindPIDByGivenPortNumber(tt.args.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindPIDByGivenPortNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FindPIDByGivenPortNumber() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindPIDByGivenPortNumberShouldFindTheCorrectProcess(t *testing.T) {
	pid, port := getProcessPidAndPort(t)
	p := Process{}
	result, err := p.FindPIDByGivenPortNumber(int32(port))
	if err != nil {
		t.Fatalf("expected nil; got err: %v", err)
	}
	if result != pid {
		t.Errorf("expected=%d; got=%d", result, pid)
	}
}

func TestFindProcessByNameShouldReturnAProcessForAValidName(t *testing.T) {
	processName, pid := getProcessPidFromName(t)
	p := Process{}
	results, err := p.FindProcessByName(processName)
	if err != nil {
		t.Fatalf("expected nil; got err: %v", err)
	}
	if results != pid {
		t.Errorf("expected=%d; got=%d", pid, results)
	}
}

func TestFindProcessByNameShouldReturnZeroForInvalidName(t *testing.T) {
	p := Process{}
	results, err := p.FindProcessByName("invalid_process_name")
	if err != nil {
		t.Fatalf("expected nil; got err: %v", err)
	}
	if results != 0 {
		t.Errorf("expected=%d; got=%d", 0, results)
	}
}

func TestKillProcessShouldReturnAnErrorForInvalidProcessPid(t *testing.T) {
	p := Process{}
	err := p.KillProcess(00000)
	if err == nil {
		t.Fatal("expected error to not be nil")
	}
}

// getProcessPidAndPort the function will return the pid as int32
// and the port as uint32. If for some reason cannot get the connections
// the function will make the test failed and if cannot find any opened
// connection then will mark the test skipped.
func getProcessPidAndPort(t *testing.T) (int32, uint32) {
	cc, err := net.Connections("")
	if err != nil {
		t.Fatalf("failed to get the connections with err: %v", err)
		return 0, 0
	}
	if len(cc) == 0 {
		t.Skip("no available connection")
		return 0, 0
	}
	c := cc[0]
	return c.Pid, c.Laddr.Port
}

// getProcessPidFromName the function return the process name as string
// if something occurred the function will make the test failed.
func getProcessPidFromName(t *testing.T) (string, int32) {
	pp, err := process.Processes()
	if err != nil {
		t.Fatalf("failed to get the processes with err: %v", err)
		return "", 0
	}
	p := pp[0]
	if len(pp) == 0 {
		t.Skip("no available processes")
		return "", 0
	}
	name, err := p.Name()
	if err != nil {
		t.Fatalf("failed to get process name with error: %v", err)
	}
	return name, p.Pid
}
