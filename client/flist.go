package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	// Unix Domain Socket address
	SockAddr		= "/tmp/flist.sock"

	// "Success", "Error" are possible response statuses
	Success Status  = "success"
	Error 	Status  = "error"
)

var (	
	// Default buffers size
	BufSize 		= 1024
)

type Request struct {
	// Command can be "run", "stop", "rm", or "ps"
	Command string			`json:"command"`
	// Args are arguments for Command, if any
	Args []string			`json:"args"`
	// Body holds data of request, if any
	Body json.RawMessage	`json:"body"`
}

// Status of Response messages
type Status string
type Response struct {
	// Response Status can be "Success" or "Error"
	Status Status			`json:"status"`
	// ErrorMsg holds error message string, if any
	ErrorMsg string 		`json:"errorMsg"`
	// Body holds data of response, if any
	Body json.RawMessage	`json:"body"`
}

// Connection structure
type Connection struct {
	// Conn is network connection between client and daemon thread
	Conn net.Conn
}

// newRequest creates a new Request object
func newRequest(body ClientData, command string, args...string) (*Request, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request := &Request{
		Command: command,
		Args: args,
		Body: json.RawMessage(data),
	}

	return request, nil
}

// Signal init signal handling procedures
func Signal(sigchnl chan os.Signal, conn net.Conn) {
	sigs := []os.Signal {syscall.SIGINT}
	
	signal.Notify(sigchnl, sigs...)

	go signalsHandler(sigchnl, conn)
}

// signalsHandler helper function to run signals handlers
func signalsHandler(sigchnl chan os.Signal, conn net.Conn) {
	for sig := range sigchnl {
		switch sig {
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

func main() {
	if len(os.Args) < 2 {
        usage()
        os.Exit(1)
    }

	conn, err := net.Dial("unix", SockAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	sigchnl := make(chan os.Signal)
	Signal(sigchnl, conn)

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