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
	Id string				`json:"id"`
	// MetaURL is URL for meta flist file on the internet
	MetaURL string			`json:"metaUrl"`
	// name of downloaded flist meta file used to mount container
	FlistName string		`json:"flistName"`
	// Entrypoint is file path for binary to execute as entrypoint in container
	Entrypoint string		`json:"entrypoint"`
	// Args are arguments of Entrypoint command
	Args []string			`json:"args"`
	// Path is location path of container directory
	Path string				`json:"path"`
	// Status is current status of container, it can be "RUNNING" or "STOPPED"
	Status string			`json:"status"`
	// pid of entrypoint process running in container
	Pid int					`json:"pid"`
	// Mountpoint is path of mounted container's filesystem
	Mountpoint string		`json:"mountpoint"`
	// filesystem object of mounted container
	Fs *g8ufs.G8ufs			`json:"fs"`
}

type Connection struct {
	// Conn is network connection between client and daemon thread
	Conn net.Conn
}

type Server struct {
	Listener net.Listener
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

// newServer Creates a new Server object
func newServer() (*Server, error) {
	l, err := net.Listen("unix", SockAddr)
    if err != nil {
		return nil, err
    }

	server := Server {
		Listener: l,
		Containers: make(map[string]Container, 0),
	}

	return &server, nil
}

// Accept is a Server method waits for and returns the next connection
func (server *Server) Accept() (*Connection, error) {
	conn, err := server.Listener.Accept()
	if err != nil {
		return nil, err
	}

	return &Connection{Conn: conn}, nil
}

// serve serves the command requested by clients by running one of
// the command "run", "stop", "rm", "ps"
func (s *Server) serve(conn Connection) {
	var request Request
	if err := conn.ReadRequest(&request); err != nil {
		log.Println(err)
		return
	}

	switch request.Command {
	case RunCmd:
		s.run(conn, request)
	case StopCmd:
		s.stop(conn, request)
	case RmCmd:
		s.rm(conn, request)
	case PsCmd:
		s.ps(conn, request)
	}
}

func main() {
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
    }

	server, err := newServer()
    if err != nil {
		log.Fatal("listen error:", err)
    }
    defer server.Listener.Close()

	for {
		conn, err := server.Accept()
        if err != nil {
			log.Fatal("accept error:", err)
        }

        go server.serve(*conn)
    }
}