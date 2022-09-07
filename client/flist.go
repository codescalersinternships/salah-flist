package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Flist struct {
	Command string			`json:"command"`
	MetaURL string			`json:"metaURL"`
	Entrypoint string		`json:"entrypoint"`
	Arg []string			`json:"arg"`
	ContainerName string	`json:"containerName"`
	ProcessPid int			`json:"processPid"`
	Mountpoint string		`json:"mountpoint"`
}

func new(command, meta, entrypoint, containerName string, arg ...string) *Flist {
	return &Flist{
		Command: command,
		MetaURL: meta,
		Entrypoint: entrypoint,
		Arg: arg,
		ContainerName: containerName,
		ProcessPid: os.Getpid(),
		Mountpoint: "",
	}
}

func Signal(sigchnl chan os.Signal) {
	sigs := []os.Signal {syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2}
	
	signal.Notify(sigchnl, sigs...)

	go signalsHandler(sigchnl)
}

func signalsHandler(sigchnl chan os.Signal) {
	for sig := range sigchnl {
		switch sig {
		case syscall.SIGUSR1:
			handleSuccessOperation()
		case syscall.SIGUSR2:
			handleFailureOperation()
		case syscall.SIGINT:
			tellDaemonToCleanupContainer()
		}
	}
}

func tellDaemonToCleanupContainer() {
	syscall.Kill(DaemonPid, syscall.SIGUSR1)

	os.Exit(1)
}

func handleSuccessOperation() {
	done <- true
}

func handleFailureOperation() {
	done <- false
}

const SockAddr = "/tmp/flist.sock"

func main() {
	if len(os.Args) < 2 {
        usage()
        os.Exit(1)
    }

	sigchnl := make(chan os.Signal)
	Signal(sigchnl)

	conn, err := net.Dial("unix", SockAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	switch os.Args[1] {
	case "run":
		run(conn)
	case "stop":
		stop(conn)
	case "ps":
		ps(conn)
	case "rm":
		rm(conn)
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {

}