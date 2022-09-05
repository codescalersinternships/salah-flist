package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

const (
	SockAddr = "/tmp/flist.sock"

	RunCmd = "run"
	StopCmd = "stop"
	RmCmd = "rm"
	PsCmd = "ps"

	StorePath = "/var/lib/flist/store" // redis.flist
	ContainersPath = "/var/lib/flist/containers"
	flistsUnpackedPath = "/var/lib/flist/tmp"
	defaultStorageHubPath = "zdb://hub.grid.tf:9900"
)

type Flist struct {
	Command string			`json:"command"`
	MetaURL string			`json:"metaURL"`
	Entrypoint string		`json:"entrypoint"`
	ContainerName string	`json:"containerName"`
}

type Container struct {
	Id string
	Path string
	Status string
}

type Worker struct {
	Conn net.Conn
	Flist Flist
	Containers []Container
}

func new(containers []Container) *Worker {
	return &Worker{ Containers: containers }
}

func (w *Worker) serve() {
	if err := json.NewDecoder(w.Conn).Decode(&w.Flist); err != nil {
		log.Fatal(err)
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

func main() {
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
    }
	
	l, err := net.Listen("unix", SockAddr)
    if err != nil {
		log.Fatal("listen error:", err)
    }
    defer l.Close()
	
	containers := make([]Container, 0)

	for {
		worker := new(containers)
        worker.Conn, err = l.Accept()
        if err != nil {
            log.Fatal("accept error:", err)
        }

        go worker.serve()
    }
}