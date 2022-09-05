package main

import (
	"log"
	"net"
	"os"
)

type Flist struct {
	Command string			`json:"command"`
	MetaURL string			`json:"metaURL"`
	Entrypoint string		`json:"entrypoint"`
	ContainerName string	`json:"containerName"`
}

func new() *Flist {
	return &Flist{ }
}

const SockAddr = "/tmp/flist.sock"

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