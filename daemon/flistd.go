package main

import (
	"encoding/json"
	"log"
	"net"
	"os"

	g8ufs "github.com/threefoldtech/0-fs"
)

const (
	// Unix Domain Socket address
	SockAddr 				= "/tmp/flist.sock"

	// Default buffers size
	BufSize 				= 1024

	// default paths for downloaded flists, mounted containers, and tmp data
	StorePath 				= "/var/lib/flist/store"
	ContainersPath			= "/var/lib/flist/containers"
	flistsUnpackedPath 		= "/var/lib/flist/tmp"
	defaultStorageHubPath   = "zdb://hub.grid.tf:9900"
	
	// RunCmd, StopCmd, RmCmd, PsCmd are flist sub-commands
	RunCmd 					= "run"
	StopCmd 				= "stop"
	RmCmd 					= "rm"
	PsCmd 					= "ps"

	// Running, Stopped are available container states
	Running 				= "RUNNING"
	Stopped 				= "STOPPED"

	// "Success", "Error" are possible response statuses
	Success Status 			= "success"
	Error 	Status 			= "error"
)

// Container contains data about created container
type Container struct {
	// Id is container's unique ID, it's same as ContainerName
	Id string
	// MetaURL is URL for meta flist file on the internet
	MetaURL string
	// name of downloaded flist meta file used to mount container
	FlistName string
	// Entrypoint is file path for binary to execute as entrypoint in container
	Entrypoint string
	// Args are arguments of Entrypoint command
	Args []string
	// Path is location path of container directory
	Path string
	// Status is current status of container, it can be "RUNNING" or "STOPPED"
	Status string

	// pid of entrypoint process running in container
	Pid int

	// filesystem object of mounted container
	fs *g8ufs.G8ufs
}

// Worker represents a thread that handles commands from clients
type Worker struct {
	// Conn is network connection between client and daemon thread
	Conn net.Conn
	// Container is the container that worker is responsible for
	Container Container
	// all mounted "RUNNING", "STOPPED" containers
	Containers map[string]Container
}

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

// new creates a new Worker object
func newWorker(conn net.Conn, containers map[string]Container) *Worker {
	return &Worker {
		Conn: conn,
		Containers: containers,
	}
}

// serve serves the command requested by clients by running one of
// the command "run", "stop", "rm", "ps"
func (w *Worker) serve() {
	var request Request
	if err := ConnectionRead(w.Conn, &request); err != nil {
		log.Println(err)
		return
	}

	switch request.Command {
	case RunCmd:
		var container Container
		if err := json.Unmarshal(request.Body, &container); err != nil {
			log.Println(err)
			if err := ConnectionErrorResponse(w.Conn, err.Error()); err != nil {
				log.Println(err)
				return
			}
			return
		}
		container.MetaURL 		= request.Args[0]
		container.Entrypoint 	= request.Args[1]
		container.Args 			= request.Args[2:]
		w.Container 			= container

		w.run()
	case StopCmd:
		w.Container = Container{Id: request.Args[0]}
		w.stop()
	case RmCmd:
		w.Container = Container{Id: request.Args[0]}
		w.rm()
	case PsCmd:
		w.ps()
	}
}


func main() {
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
    }
	
	l, err := net.Listen("unix", SockAddr)
    if err != nil {
		log.Fatal("listen error:", err)
    }
    defer l.Close()
	
	containers := make(map[string]Container, 0)

	for {
		conn, err := l.Accept()
        if err != nil {
			log.Fatal("accept error:", err)
        }

		worker := newWorker(conn, containers)

        go worker.serve()
    }
}