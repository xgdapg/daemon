package daemon

import (
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func Daemon() {
	envPpid := os.Getenv("__daemon_daemon_ppid__")
	os.Setenv("__daemon_daemon_ppid__", strconv.Itoa(os.Getpid()))
	ppid := os.Getppid()
	if ppid == 1 {
		return
	}
	if envPpid == strconv.Itoa(ppid) {
		return
	}
	cmd := getCmd()
	if err := cmd.Start(); err != nil {
		os.Exit(-1)
	}
	os.Exit(0)
}

func Monitor() {
	envPpid := os.Getenv("__daemon_monitor_ppid__")
	os.Setenv("__daemon_monitor_ppid__", strconv.Itoa(os.Getpid()))
	ppid := os.Getppid()
	if envPpid == strconv.Itoa(ppid) {
		return
	}
	var cmd *exec.Cmd
	sigChan := make(chan os.Signal, 4)
	signal.Notify(sigChan)
	go func() {
		for {
			sig := <-sigChan
			if sig == syscall.SIGKILL || sig == syscall.SIGINT || sig == syscall.SIGTERM {
				if cmd != nil && cmd.Process != nil {
					cmd.Process.Kill()
				}
				os.Exit(0)
			}
		}
	}()
	for {
		cmd = getCmd()
		if err := cmd.Start(); err != nil {
			os.Exit(-1)
		}
		cmd.Wait()
		time.Sleep(time.Second)
	}
}

func DaemonAndMonitor() {
	Daemon()
	Monitor()
}

func getCmd() *exec.Cmd {
	cmd := exec.Command(os.Args[0])
	if len(os.Args) > 1 {
		cmd.Args = append(cmd.Args, os.Args[1:]...)
	}
	cmd.Env = os.Environ()
	return cmd
}
