package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

// stop stops running container by sending SIGTERM signal,
// and after a grace period, SIGKILL to the process inside
// the container. this is client side of stop command. firstly,
// it sends command's data to daemon, which carry the work of
// applying the command, then waits for a signal from daemon to
// know whether the command was carried successfully or not.
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