package main

import (
	"fmt"
	"log"
)

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