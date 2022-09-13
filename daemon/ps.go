package main

import (
	"encoding/json"
	"log"
)

// ps lists containers in a tabular format. this is server side of ps command.
// it carries the work of collecting current containers data,
// and sends it to requesting client over network connection in response body.
func (s *Server) ps(conn Connection, request Request) {
	records, err := json.Marshal(s.Containers)
	if err != nil {
		log.Println(err)
		return
	}
	
	response := Response {
		Status: Success,
		Body: 	json.RawMessage(records),
	}
	if err := conn.SendResponse(response); err != nil {
		log.Println(err)
		return
	}
}