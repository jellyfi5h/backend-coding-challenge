package server

import (
	"os"
	"os/exec"
	"syscall"
)

func forkProcess() error {
	cmd := exec.Command(os.Args[0])
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Start()
}

// release all parent process resources
func releaseProcess() error {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}
	return p.Release()
}

//run a child process cloned from it's parent
func runCloneProcess() (curr *os.Process, err error) {
	curr, _ = os.FindProcess(os.Getpid())
	if err = forkProcess(); err != nil {
		return
	}
	if err = releaseProcess(); err != nil {
		return
	}
	return
}

//CreateDaemon creates a child process and kill it's parent and set settings to turn a daemon
func CreateDaemon() (err error) {
	var parent *os.Process

	if os.Getppid() == 1 {
		setDaemonSettings()
		return
	}
	parent, err = runCloneProcess()
	if err != nil {
		return
	}
	if parent != nil {
		os.Exit(0) // kill the parent process
	}
	return
}

func setDaemonSettings() {
	var devNull *os.File

	devNull = os.NewFile(0, os.DevNull)
	defer devNull.Close()
	//os.Stdin, os.Stdout, os.Stderr = devNull, devNull, devNull

	syscall.Umask(0)   // set permissions
	syscall.Setsid()   //retrieve any controlling terminal of the daemon
	syscall.Chdir("/") //move to root directory
	//other settings ... (services)
}
