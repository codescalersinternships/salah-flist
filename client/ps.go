package main

import (
	"fmt"
	"net"
)

func ps(conn net.Conn) {
	writeData(conn, fmt.Sprintf("%s\n", "ps"))
}