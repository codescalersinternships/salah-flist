package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	rmCmd           = flag.NewFlagSet("stop", flag.ExitOnError)
	rmContainerName = rmCmd.String("container", "", "container name")
)

func rm(conn net.Conn) {
	if err := rmCmd.Parse(os.Args[2:]); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	writeData(conn, fmt.Sprintf("%s\n", rmCmd.Name()))
	writeData(conn, fmt.Sprintf("%s\n", *rmContainerName))
}