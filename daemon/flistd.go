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
	SockAddr = "/tmp/flist.sock"

	StorePath = "/var/lib/flist/store"
	ContainersPath = "/var/lib/flist/containers"
	flistsUnpackedPath = "/var/lib/flist/tmp"
	defaultStorageHubPath = "zdb://hub.grid.tf:9900"
	
	RunCmd = "run"
	StopCmd = "stop"
	RmCmd = "rm"
	PsCmd = "ps"

	Running = "RUNNING"
	Stopped = "STOPPED"
)

type Flist struct {
	Command string			`json:"command"`
	MetaURL string			`json:"metaURL"`
	Entrypoint string		`json:"entrypoint"`
	ContainerName string	`json:"containerName"`
	ProcessPid int			`json:"processPid"`
	Mountpoint string		`json:"mountpoint"`
}

type Container struct {
	Id string
	FlistName string
	Entrypoint string
	Path string
	Status string
	Pid int
	fs *g8ufs.G8ufs
}

type Worker struct {
	fs *g8ufs.G8ufs
	Conn net.Conn
	Flist Flist
	Containers map[string]Container
}

func new(containers map[string]Container) *Worker {
	return &Worker{ Containers: containers }
}

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

func (w *Worker) Signal(sigchnl chan os.Signal) {
	sigs := []os.Signal {syscall.SIGUSR1}
	
	signal.Notify(sigchnl, sigs...)

	go w.signalsHandler(sigchnl)
}

func (w *Worker) signalsHandler(sigchnl chan os.Signal) {
	for sig := range sigchnl {
		switch sig {
		case syscall.SIGUSR1:
			w.cleanUPContainer()
		}
	}
}
func (w *Worker) cleanUPContainer() {
	fs := w.Containers[w.Flist.ContainerName].fs
	if fs != nil {
		if err := fs.Unmount(); err != nil {
			log.Println(err)
		}
	}

	delete(w.Containers, w.Flist.ContainerName)
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