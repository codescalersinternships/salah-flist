package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
)

var (
	stopCmd           = flag.NewFlagSet("stop", flag.ExitOnError)
	stopContainerName = stopCmd.String("container", "", "container name")
)

func stop(conn net.Conn) {
	if err := stopCmd.Parse(os.Args[2:]); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	flist := new()
	flist.Command = stopCmd.Name()
	flist.ContainerName = *stopContainerName

	data, err := json.Marshal(flist)
	if err != nil {
		log.Fatal(err)
	}

	writeData(conn, data)
}