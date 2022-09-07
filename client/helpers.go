package main

import (
	"log"
	"net"
	"os"
	"strconv"
)

// writeData writes data on network connection
func writeData(conn net.Conn, data []byte) {
	if n, err := conn.Write(data);
	err != nil || n != len(data) {
		log.Println(err)
		os.Exit(1)
	}
}

// getDaemonPid gets pid of daemon by communication
// it on network connection
func getDaemonPid(conn net.Conn) int {
	buf := make([]byte, BufSize)

	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	DaemonPid, err := strconv.Atoi(string(buf[:n]))
	if err != nil {
		log.Fatal(err)
	}

	return DaemonPid
}