package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

var (
	done = make(chan bool, 1)
)

func stop(conn net.Conn) {
	flist := new(os.Args[1], "", "", os.Args[2], os.Args[3:]...)
	data, err := json.Marshal(flist)
	if err != nil {
		log.Fatal(err)
	}

	writeData(conn, data)

	success := <-done
	if success {
		log.Printf("container <%s> was stopped successfully\n", flist.ContainerName)
	} else {
		log.Printf("couldn't stop container <%s>, please check daemon logs\n", flist.ContainerName)
	}
}