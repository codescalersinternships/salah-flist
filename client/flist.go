package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// Unix Domain Socket address
const SockAddr = "/tmp/flist.sock"

var (
	// pid of daemon process, clients needs it to send signals to daemon
	DaemonPid = 0
	
	// Default buffers size
	BufSize = 1024
)

var (
	// done is a channel used to synchronize procedures in client and daemon
	done = make(chan bool, 1)
)


// Flist contains container's data sent in requests
// from client to daemon
type Flist struct {
	// Command can be "run", "stop", "rm", or "ps"
	Command string			`json:"command"`
	// MetaURL is URL for meta flist file on the internet
	MetaURL string			`json:"metaURL"`
	// Entrypoint is file path for binary to execute as entrypoint in container
	Entrypoint string		`json:"entrypoint"`
	Arg []string			`json:"arg"`
	// name of the container, it's set by "stop" and "rm" commands
	// to specify which container to stop or remove.
	ContainerName string	`json:"containerName"`
	// pid of client process, daemon needs it to send signals to clients
	ClientPid int			`json:"clientPid"`
	// path of container's mountpoint
	Mountpoint string		`json:"mountpoint"`
}

// new creates a new flist object
func new(command, meta, entrypoint, containerName string, arg ...string) *Flist {
	return &Flist{
		Command: command,
		MetaURL: meta,
		Entrypoint: entrypoint,
		Arg: arg,
		ContainerName: containerName,
		ClientPid: os.Getpid(),
	}
}

// Signal init signal handling procedures
func Signal(sigchnl chan os.Signal, conn net.Conn) {
	sigs := []os.Signal {syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2}
	
	signal.Notify(sigchnl, sigs...)

	go signalsHandler(sigchnl, conn)
}

// signalsHandler helper function to run signals handlers
func signalsHandler(sigchnl chan os.Signal, conn net.Conn) {
	for sig := range sigchnl {
		switch sig {
		case syscall.SIGUSR1:
			handleSuccessOperation()
		case syscall.SIGUSR2:
			handleFailureOperation()
		case syscall.SIGINT:
			tellDaemonToCleanupContainer(conn)
		}
	}
}

// tellDaemonToCleanupContainer closes connection with daemon.
// when connection close, daemon know that client exited and clean up container.
func tellDaemonToCleanupContainer(conn net.Conn) {
	if err := conn.Close(); err != nil {
		log.Println(err)
		return
	}

	os.Exit(1)
}

// handleSuccessOperation is a SIGUSR1 handler. it sends "true" to done
// channel. daemon tells clients that command carried successfully by sending
// SIGUSR1 signal
func handleSuccessOperation() {
	done <- true
}

// handleFailureOperation is a SIGUSR2 handler. it sends "false" to done
// channel. daemon tells clients that command failed to complete by sending
// SIGUSR2 signal
func handleFailureOperation() {
	done <- false
}

func main() {
	if len(os.Args) < 2 {
        usage()
        os.Exit(1)
    }

	conn, err := net.Dial("unix", SockAddr)

	sigchnl := make(chan os.Signal)
	Signal(sigchnl, conn)

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
	fmt.Printf("usage: %s COMMAND ARGS...\n", os.Args[0])
}