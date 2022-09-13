package main

import (
	"fmt"
	"log"
	"os"

)

// rm removes stopped containers. this is server side of rm command.
// it carries the work of removing container and its data.
func (s *Server) rm(conn Connection, request Request) {
	containerID := request.Args[0]
	
	container, ok := s.Containers[containerID]
	if !ok {
		log.Printf("container name <%s> doesn't exist\n", containerID)
		msg := fmt.Sprintf("container name <%s> doesn't exist\n", containerID)
		if err := conn.SendErrorResponse(msg); err != nil {
			log.Println(err)
			return
		}

		return
	}

	if container.Status == Stopped {
		delete(s.Containers, containerID)

		if err := os.RemoveAll(container.Path); err != nil {
			log.Println(err)
			if err := conn.SendErrorResponse(err.Error()); err != nil {
				log.Println(err)
				return
			}
			return
		}

		if err := conn.SendResponse(Response{Status: Success}); err != nil {
			log.Println(err)
			return
		}
	}

	msg := fmt.Sprintf("couldn't remove container <%s>", containerID)
	if err := conn.SendErrorResponse(msg); err != nil {
		log.Println(err)
		return
	}
}