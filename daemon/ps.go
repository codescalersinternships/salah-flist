package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ps lists containers in a tabular format. this is server side of ps command.
// it carries the work of collecting current containers data,
// and sends it to requesting client over network connection in response body.
func (w *Worker) ps() {
	var records string
	for _, container := range w.Containers {
		record := fmt.Sprintf("%s\t%s\t%s\t%s,", container.Id, container.FlistName, container.Entrypoint, container.Status)
		
		records = fmt.Sprintf("%s%s", records, record)
	}
	
	response := Response {
		Status: Success,
		Body: 	json.RawMessage([]byte(fmt.Sprintf("{\"records\": %q}", records))),
	}
	if err := ConnectionWrite(w.Conn, response); err != nil {
		log.Println(err)
		return
	}
}