package daemon

import (
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	Daemon uint = 1 << iota
	Monitor
)

var (
	once  = new(sync.Once)
	sPid  = strconv.Itoa(os.Getpid())
	sPpid = strconv.Itoa(os.Getppid())
	env   = os.Environ()
)

func Exec(mode uint) {
	once.Do(func() {
		if mode&Daemon == Daemon {
			daemon()
		}
		if mode&Monitor == Monitor {
			monitor()
		}
	})
}

func isDaemoned() bool {
	envPpid := getEnv("__daemon_daemon_ppid__")
	setEnv("__daemon_daemon_ppid__", sPid)
	if sPpid == "1" {
		return true
	}
	if envPpid == sPpid {
		return true
	}
	return false
}

func daemon() {
	if isDaemoned() {
		return
	}
	cmd := getCmd()
	if err := cmd.Start(); err != nil {
		os.Exit(-1)
	}
	os.Exit(0)
}

func isMonitored() bool {
	envPpid := getEnv("__daemon_monitor_ppid__")
	setEnv("__daemon_monitor_ppid__", sPid)
	if envPpid == sPpid {
		return true
	}
	return false
}

func monitor() {
	if isMonitored() {
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

func getCmd() *exec.Cmd {
	cmd := exec.Command(os.Args[0])
	if len(os.Args) > 1 {
		cmd.Args = append(cmd.Args, os.Args[1:]...)
	}
	cmd.Env = env
	return cmd
}

func setEnv(k, v string) {
	k = k + "="
	for i, e := range env {
		if strings.HasPrefix(e, k) {
			env[i] = k + v
			return
		}
	}
	env = append(env, k+v)
}

func getEnv(k string) string {
	k = k + "="
	for _, e := range env {
		if strings.HasPrefix(e, k) {
			return e[len(k):]
		}
	}
	return ""
}
