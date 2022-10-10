package main

import (
	"log"
	"net"
	"os"
)

// stop stops running container by sending SIGTERM signal,
// and after a grace period, SIGKILL to the process inside
// the container. this is client side of stop command. firstly,
// it sends command's data to daemon in a request, which carry the work of
// applying the command, then waits for a response from daemon to
// know whether the command was carried successfully or not.
func stop(conn net.Conn) {
	request, err := newRequest(ClientData{}, os.Args[1], os.Args[2:]...)
	if err != nil {
		log.Println(err)
		return
	}

	if err := ConnectionWrite(conn, request); err != nil {
		log.Println(err)
		return
	}

	var response Response
	if err := ConnectionRead(conn, &response); err != nil {
		log.Println(err)
		if err := ConnectionErrorResponse(conn, err.Error()); err != nil {
			log.Println(err)
			return
		}
		return
	}

	if response.Status == Error {
		log.Printf("couldn't stop container <%s>: %s\n", os.Args[2], response.ErrorMsg)
		return
	} else {
		log.Printf("container <%s> was stopped successfully\n", os.Args[2])
	}
}