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
func (w *Worker) stop() {
	if _, ok := w.Containers[w.Container.Id]; !ok {
		log.Printf("container <%s> doesn't exist\n", w.Container.Id)
		msg := fmt.Sprintf("container <%s> doesn't exist\n", w.Container.Id)
		if err := ConnectionErrorResponse(w.Conn, msg); err != nil {
			log.Println(err)
			return
		}

		return
	}

	if w.Containers[w.Container.Id].Status != Running {
		log.Printf("this container <%s> is not running\n", w.Container.Id)
		msg := fmt.Sprintf("this container <%s> is not running\n", w.Container.Id)
		if err := ConnectionErrorResponse(w.Conn, msg); err != nil {
			log.Println(err)
			return
		}

		return
	}

	pid := w.Containers[w.Container.Id].Pid
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		log.Println(err)
		if err := ConnectionErrorResponse(w.Conn, err.Error()); err != nil {
			log.Println(err)
			return
		}
		return
	}
	time.Sleep(10 * time.Second)
	if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
		log.Println(err)
	}

	container := Container {
		Status: 	Stopped,
		Id: 		w.Containers[w.Container.Id].Id,
		FlistName: 	w.Containers[w.Container.Id].FlistName,
		Entrypoint: w.Containers[w.Container.Id].Entrypoint,
		Path: 		w.Containers[w.Container.Id].Path,
		Pid: 		w.Containers[w.Container.Id].Pid,
		fs: 		w.Containers[w.Container.Id].fs,
	}
	w.Containers[w.Container.Id] = container

	response := Response {
		Status: Success,
	}
	if err := ConnectionWrite(w.Conn, response); err != nil {
		log.Println(err)
		return
	}
}