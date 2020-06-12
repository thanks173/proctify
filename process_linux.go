// +build linux

package proctify

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

// linuxProcess is an implementation of Process from Linux
type linuxProcess struct {
	pid   int
	ppid  int
	state rune

	executable string
	cmd        string
	path       string

	uid int
}

func (p *linuxProcess) Pid() int {
	return p.pid
}

func (p *linuxProcess) PPid() int {
	return p.ppid
}

func (p *linuxProcess) Executable() string {
	return p.executable
}

func (p *linuxProcess) Cmd() string {
	return p.cmd
}

func (p *linuxProcess) Path() string {
	return p.path
}

func (p *linuxProcess) State() rune {
	switch p.state {
	case StateRunning:
		return StateRunning
	case StateDead:
		return StateDead
	case StateZombie:
		return StateZombie
	}
	return StateSleep
}

func (p *linuxProcess) Uid() int {
	return p.uid
}

func runningProcesses() ([]Process, error) {
	pDir, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer pDir.Close()

	res := make([]Process, 0, 50)
	for {
		names, err := pDir.Readdirnames(10)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, n := range names {
			// Filter only numberic dir name
			if n[0] < '0' || n[0] > '9' {
				continue
			}

			pid, err := strconv.ParseInt(n, 10, 0)
			if err != nil {
				continue
			}

			p, err := lookupPid(int(pid))
			if err != nil {
				continue
			}

			res = append(res, p)
		}
	}

	return res, nil
}

func lookupPid(pid int) (*linuxProcess, error) {
	p := &linuxProcess{pid: pid}
	err := p.reload()
	return p, err
}

// reload refresh all the process information from the pid
func (p *linuxProcess) reload() error {
	info, err := os.Stat(fmt.Sprintf("/proc/%d", p.pid))
	if err != nil {
		return err
	}

	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		p.uid = int(stat.Uid)
	}

	cmdPath := fmt.Sprintf("/proc/%d/cmdline", p.pid)
	dataBytes, err := ioutil.ReadFile(cmdPath)
	if err != nil {
		return err
	}
	p.cmd = string(dataBytes)

	exeLink := fmt.Sprintf("/proc/%d/exe", p.pid)
	exePath, err := os.Readlink(exeLink)
	if err == nil {
		p.path = exePath
		idx := strings.LastIndex(exePath, "/") + 1
		p.executable = exePath[idx:]
	}

	statPath := fmt.Sprintf("/proc/%d/stat", p.pid)
	dataBytes, err = ioutil.ReadFile(statPath)
	if err != nil {
		return err
	}

	data := string(dataBytes)
	r := regexp.MustCompile(`\(+(.+?)\)+\s*(\w+)\s*(\d+)`)
	match := r.FindStringSubmatch(data)

	if len(match) != 4 {
		return errors.New("can't parse stat info")
	}

	cmd := match[1]
	p.state = rune(match[2][0])
	p.ppid, err = strconv.Atoi(match[3])
	if err != nil {
		return err
	}

	if p.cmd == "" {
		p.cmd = cmd
	}

	return err
}
