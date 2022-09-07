package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

// rm removes stopped containers. this is client side of rm command,
// it firstly sends command data to daemon, which does the work of removing container,
// then waits for a signal from daemon to know whether the command was
// carried successfully or not.
func rm(conn net.Conn) {
	flist := new(os.Args[1], "", "", os.Args[2], os.Args[3:]...)

	data, err := json.Marshal(flist)
	if err != nil {
		log.Fatal(err)
	}

	writeData(conn, data)

	success := <-done
	if success {
		log.Printf("container <%s> was removed successfully\n", flist.ContainerName)
	} else {
		log.Printf("couldn't remove container <%s>, please check daemon logs\n", flist.ContainerName)
	}
}