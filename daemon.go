package main

import (
	"os"
	"os/exec"
	"strconv"
)

func init() {
	daemon()
}

func daemon() {
	ppid := os.Getppid()
	if ppid == 1 {
		return
	}
	if dppid := os.Getenv("__daemon_ppid"); dppid == strconv.Itoa(ppid) {
		return
	}

	cmd := exec.Command(os.Args[0])
	if len(os.Args) > 1 {
		cmd = exec.Command(os.Args[0], os.Args[1:]...)
	}
	cmd.Env = append(os.Environ(), "__daemon_ppid="+strconv.Itoa(os.Getpid()))
	if err := cmd.Start(); err != nil {
		os.Exit(-1)
	}
	os.Exit(0)
}