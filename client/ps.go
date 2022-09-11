package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"text/tabwriter"
)

type Records struct {
	Records string
}

// ps lists containers in a tabular format. this is the client side
// of ps command, at first client sends a request message, then waits daemon to send back response.
// client manipulates data in response body then present it on STDOUT in tabular format.
func ps(conn net.Conn) {
	request := newRequest(os.Args[1], os.Args[2:]...)

	if err := ConnectionWrite(conn, request); err != nil {
		log.Println(err)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 3, 3, 3, ' ', 0)
	
	if _, err := fmt.Fprintln(w, "CONTAINER ID\tFLIST\tCOMMAND\tSTATUS"); err != nil {
		log.Println(err)
		return
	}

	var response Response
	if err := ConnectionRead(conn, &response); err != nil {
		log.Println(err)
		return
	}

	if response.Status == Success {
		var recordsData Records
		if err := json.Unmarshal(response.Body, &recordsData); err != nil {
			log.Println(err)
			return
		}
		records := strings.Split(recordsData.Records, ",")
		
		for _, record := range records {
			_, err := fmt.Fprintln(w, record)
			if err != nil {
				log.Println(err)
				return
			}
		}
	} else {
		log.Println(response.ErrorMsg)
		return
	}
	
	w.Flush()
}