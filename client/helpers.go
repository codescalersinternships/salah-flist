package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func writeData(conn net.Conn, data string) {
	if n, err := conn.Write([]byte(fmt.Sprintf("%s\n", data)));
	err != nil || n != len(data) + 1 {
		log.Println(err)
		os.Exit(1)
	}
}