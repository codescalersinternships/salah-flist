package main

import (
	"encoding/json"
	"log"
	"net"
)

func ps(conn net.Conn) {
	flist := new("ps", "", "", "")
	
	data, err := json.Marshal(flist)
	if err != nil {
		log.Fatal(err)
	}

	writeData(conn, data)
}