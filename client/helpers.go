package main

import (
	"log"
	"net"
	"os"
)

func writeData(conn net.Conn, data []byte) {
	if n, err := conn.Write(data);
	err != nil || n != len(data) {
		log.Println(err)
		os.Exit(1)
	}
}