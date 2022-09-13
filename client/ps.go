package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"text/tabwriter"
)

type Record struct {
	// Id is container's unique ID, it's same as ContainerName
	Id string				`json:"id"`
	// name of downloaded flist meta file used to mount container
	FlistName string		`json:"flistName"`
	// Entrypoint is file path for binary to execute as entrypoint in container
	Entrypoint string		`json:"entrypoint"`
	// Status is current status of container, it can be "RUNNING" or "STOPPED"
	Status string			`json:"status"`
}

// ps lists containers in a tabular format. this is the client side
// of ps command, at first client sends a request message, then waits daemon to send back response.
// client manipulates data in response body then present it on STDOUT in tabular format.
func ps(conn net.Conn) {
	request, err := newRequest(ClientData{}, os.Args[1], os.Args[2:]...)
	if err != nil {
		log.Println(err)
		return
	}

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
		var records map[string]Record
		if err := json.Unmarshal(response.Body, &records); err != nil {
			log.Println(err)
			return
		}		
		for _, record := range records {
			_, err := fmt.Fprintf(w, 
				"%s\t%s\t%s\t%s\t\n",
				record.Id, record.FlistName, record.Entrypoint, record.Status)
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