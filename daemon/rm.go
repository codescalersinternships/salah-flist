package main

import (
	"fmt"
	"log"
	"os"
)

// rm removes stopped containers. this is server side of rm command.
// it carries the work of removing container and its data.
func (w *Worker) rm() {
	if _, ok := w.Containers[w.Container.Id]; !ok {
		log.Printf("container name <%s> doesn't exist\n", w.Container.Id)
		msg := fmt.Sprintf("container name <%s> doesn't exist\n", w.Container.Id)
		ConnectionErrorResponse(w.Conn, msg)

		return
	}

	container := w.Containers[w.Container.Id]
	if container.Status == Stopped {
		delete(w.Containers, w.Container.Id)

		if err := os.RemoveAll(container.Path); err != nil {
			log.Println(err)
			ConnectionErrorResponse(w.Conn, err.Error())
			return
		}

		if err := ConnectionWrite(w.Conn, Response{Status: Success}); err != nil {
			log.Println(err)
			return
		}
	}

	msg := fmt.Sprintf("couldn't remove container <%s>", w.Container.Id)
	ConnectionErrorResponse(w.Conn, msg)
}