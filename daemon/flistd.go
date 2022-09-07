package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	g8ufs "github.com/threefoldtech/0-fs"
)

const (
	// Unix Domain Socket address
	SockAddr = "/tmp/flist.sock"

	// default paths for downloaded flists, mounted containers, and tmp data
	StorePath = "/var/lib/flist/store"
	ContainersPath = "/var/lib/flist/containers"
	flistsUnpackedPath = "/var/lib/flist/tmp"
	defaultStorageHubPath = "zdb://hub.grid.tf:9900"
	
	// RunCmd, StopCmd, RmCmd, PsCmd are flist sub-commands
	RunCmd = "run"
	StopCmd = "stop"
	RmCmd = "rm"
	PsCmd = "ps"

	// Running, Stopped are available container states
	Running = "RUNNING"
	Stopped = "STOPPED"
)

// Flist contains container's data sent in responses
// from daemon to clients
type Flistd struct {
	// Command can be "run", "stop", "rm", or "ps"
	Command string			`json:"command"`
	// MetaURL is URL for meta flist file on the internet
	MetaURL string			`json:"metaURL"`
	// Entrypoint is file path for binary to execute as entrypoint in container
	Entrypoint string		`json:"entrypoint"`
	// name of the container, it's set by "stop" and "rm" commands
	// to specify which container to stop or remove.
	ContainerName string	`json:"containerName"`
	// pid of client process, daemon needs it to send signals to clients
	ClientPid int			`json:"clientPid"`
	// path of container's mountpoint
	Mountpoint string		`json:"mountpoint"`
}

// Container contains data about created container
type Container struct {
	// Id is container's unique ID, it's same as ContainerName
	Id string
	// name of downloaded flist meta file used to mount container
	FlistName string
	// Entrypoint is file path for binary to execute as entrypoint in container
	Entrypoint string
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
	// filesystem object of mounted container
	fs *g8ufs.G8ufs
	// Conn is network connection between client and daemon thread
	Conn net.Conn
	Flist Flistd
	// all mounted "RUNNING", "STOPPED" containers
	Containers map[string]Container
}

// new creates a new Worker object
func new(containers map[string]Container) *Worker {
	return &Worker{ Containers: containers }
}

// serve serves the command requested by clients by running one of
// the command "run", "stop", "rm", "ps"
func (w *Worker) serve() {
	if err := json.NewDecoder(w.Conn).Decode(&w.Flist); err != nil {
		log.Println(err)
		return
	}

	switch w.Flist.Command {
	case RunCmd:
		w.run()
	case StopCmd:
		w.stop()
	case RmCmd:
		w.rm()
	case PsCmd:
		w.ps()
	}
}

// Signal init signal handling procedures
func (w *Worker) Signal(sigchnl chan os.Signal) {
	sigs := []os.Signal {syscall.SIGUSR1}
	
	signal.Notify(sigchnl, sigs...)

	go w.signalsHandler(sigchnl)
}

// signalsHandler helper function to run signals handlers
func (w *Worker) signalsHandler(sigchnl chan os.Signal) {
	for sig := range sigchnl {
		switch sig {
		// clients send SIGUSR1 signals to tell daemon cleanup stopped container
		case syscall.SIGUSR1:
			w.cleanUPContainer()
		}
	}
}

// cleanUPContainer helper function used as SIGUSR1 handler
// to clean up container data on catching SIGUSR1 signals from clients
func (w *Worker) cleanUPContainer() {
	fs := w.Containers[w.Flist.ContainerName].fs
	if fs != nil {
		if err := fs.Unmount(); err != nil {
			log.Println(err)
		}
	}

	delete(w.Containers, w.Flist.ContainerName)
}

// reportSuccessOperation sends SIGUSR1 signal to client to tell client that
// the requested command was carried successfully
func (w *Worker) reportSuccessOperation() {
	if err := syscall.Kill(w.Flist.ClientPid, syscall.SIGUSR1); err != nil {
		log.Println(err)
		return
	}
}

// reportSuccessOperation sends SIGUSR2 signal to client to tell client that
// the requested command failed to complete
func (w *Worker) reportFailureOperation() {
	if err := syscall.Kill(w.Flist.ClientPid, syscall.SIGUSR2); err != nil {
		log.Println(err)
		return
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
		worker := new(containers)
        worker.Conn, err = l.Accept()
        if err != nil {
            log.Fatal("accept error:", err)
        }

        go worker.serve()
    }
}