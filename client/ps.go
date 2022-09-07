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


func ps(conn net.Conn) {
	flist := new("ps", "", "", "")
	
	data, err := json.Marshal(flist)
	if err != nil {
		log.Fatal(err)
	}

	writeData(conn, data)

	w := tabwriter.NewWriter(os.Stdout, 3, 3, 3, ' ', 0)
	
	if _, err := fmt.Fprintln(w, "CONTAINER ID\tFLIST\tCOMMAND\tSTATUS"); err != nil {
		log.Println(err)
		return
	}

	buf := make([]byte, BufSize)
	var d string
	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			break
		}
		
		d = fmt.Sprintf("%s%s", d, buf[:n])
	}
	records := strings.Split(d, ",")

	success := <-done
	if success {
		for _, record := range records {
			_, err := fmt.Fprintln(w, record)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
	w.Flush()
}