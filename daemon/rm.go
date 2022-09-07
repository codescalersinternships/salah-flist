package main

import (
	"log"
	"os"
)

// rm removes stopped containers. this is server side of rm command.
// it carries the work of removing container and its data and sends it
// to requesting client over network connection. it then sends SIGUSR1
// signal to client to represent success, otherwise it sends SIGUSR2
// signal to represent failure.
func (w *Worker) rm() {
	if _, ok := w.Containers[w.Flist.ContainerName]; !ok {
		log.Printf("container name <%s> doesn't exist\n", w.Flist.ContainerName)
		w.reportFailureOperation()
		return
	}

	container := w.Containers[w.Flist.ContainerName]
	if container.Status == Stopped {
		container.fs.Unmount();

		delete(w.Containers, w.Flist.ContainerName)

		if err := os.RemoveAll(container.Path); err != nil {
			log.Println(err)
			w.reportFailureOperation()
			return
		}

		w.reportSuccessOperation()
	}

	w.reportFailureOperation()
}