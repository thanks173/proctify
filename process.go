package proctify

const (
	StateRunning rune = 'R'
	StateDead    rune = 'X'
	StateSleep   rune = 'S'
	StateZombie  rune = 'Z'
)

type Process interface {
	// The process ID
	Pid() int

	// The parent process ID
	PPid() int

	// Executable name
	Executable() string

	// The command that was used to launch the process
	Cmd() string

	// The path to the executable
	Path() string

	// The process state
	State() rune

	// The user that started the process
	Uid() int
}

// Processes returns the list of all running processes.
func Processes() ([]Process, error) {
	return runningProcesses()
}

// LookupPid looks up a process by given pid.
func LookupPid(pid int) (Process, error) {
	return lookupPid(pid)
}
