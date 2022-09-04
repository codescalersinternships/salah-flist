package main

import (
	"flag"
	"fmt"
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

	writeData(conn, fmt.Sprintf("%s\n", stopCmd.Name()))
	writeData(conn, fmt.Sprintf("%s\n", *stopContainerName))
}