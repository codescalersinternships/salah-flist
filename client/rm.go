package main

import (
	"log"
	"net"
	"os"
)

// rm removes stopped containers. this is client side of rm command,
// it firstly sends command data to daemon in request message,
// which does the work of removing container, then waits for a response
// from daemon to know whether the command was carried successfully or not.
func rm(conn net.Conn) {
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
		log.Printf("couldn't remove container <%s>, please check daemon logs\n", os.Args[2])
		return
	} else {
		log.Printf("container <%s> was removed successfully\n", os.Args[2])
	}
}