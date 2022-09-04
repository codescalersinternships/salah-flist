package main

import (
	"bufio"
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
)

type Daemon struct {
	Conn net.Conn
}

func new() *Daemon {
	return &Daemon{ }
}

func (d *Daemon) serve() {
	scanner := bufio.NewScanner(d.Conn)
	if scanner.Scan() {
		command := scanner.Text()
		switch command {
		case RunCmd:
			d.run()
		case StopCmd:
			d.stop()
		case RmCmd:
			d.rm()
		case PsCmd:
			d.ps()
		}
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
	
	daemon := new()
	for {
        daemon.Conn, err = l.Accept()
        if err != nil {
            log.Fatal("accept error:", err)
        }

        go daemon.serve()
    }
}