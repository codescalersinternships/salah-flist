package main

import (
	"encoding/json"
	"flag"
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

	flist := new(rmCmd.Name(), "", "", *rmContainerName)

	data, err := json.Marshal(flist)
	if err != nil {
		log.Fatal(err)
	}

	writeData(conn, data)
}