package main

import (
	"fmt"
	"log"
)

// ps lists containers in a tabular format. this is server side of ps command.
// it carries the work of collecting current containers data,
// and sends it to requesting client over network connection. it then
// sends SIGUSR1 signal to client to represent success, otherwise it
// sends SIGUSR2 signal to represent failure.
func (w *Worker) ps() {
	for _, container := range w.Containers {
		record := fmt.Sprintf("%s\t%s\t%s\t%s,", container.Id, container.FlistName, container.Entrypoint, container.Status)
		
		n, err := w.Conn.Write([]byte(record))
		if err != nil || len([]byte(record)) != n {
			log.Println(err)
			w.reportFailureOperation()
			return
		}
	}
	w.Conn.Close()
	w.reportSuccessOperation()
}