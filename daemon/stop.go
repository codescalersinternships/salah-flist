package main

import (
	"fmt"
	"log"
	"syscall"
	"time"
)

// stop stops running container by sending SIGTERM signal,
// and after a grace period, SIGKILL to the process inside
// the container. this is server side of stop command, it
// carries the work of stopping entrypoint process inside
// the container.
func (s *Server) stop(conn Connection, request Request) {
	containerID := request.Args[0]

	containerToStop, ok := s.Containers[containerID]
	if !ok {
		log.Printf("container <%s> doesn't exist\n", containerID)
		msg := fmt.Sprintf("container <%s> doesn't exist\n", containerID)
		if err := conn.SendErrorResponse(msg); err != nil {
			log.Println(err)
			return
		}

		return
	}

	if containerToStop.Status != Running {
		log.Printf("this container <%s> is not running\n", containerID)
		msg := fmt.Sprintf("this container <%s> is not running\n", containerID)
		if err := conn.SendErrorResponse(msg); err != nil {
			log.Println(err)
			return
		}

		return
	}

	pid := containerToStop.Pid
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		log.Println(err)
		if err := conn.SendErrorResponse(err.Error()); err != nil {
			log.Println(err)
			return
		}
		return
	}
	time.Sleep(10 * time.Second)
	if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
		log.Println(err)
	}

	containerToStop.Status = Stopped
	s.Containers[containerID] = containerToStop

	response := Response {
		Status: Success,
	}
	if err := conn.SendResponse(response); err != nil {
		log.Println(err)
		return
	}
}