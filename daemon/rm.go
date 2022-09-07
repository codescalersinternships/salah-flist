package main

import (
	"log"
	"os"
)

func (w *Worker) rm() {
	if _, ok := w.Containers[w.Flist.ContainerName]; !ok {
		log.Printf("container name <%s> doesn't exist\n", w.Flist.ContainerName)
		w.reportFailureOperation()
		return
	}

	container := w.Containers[w.Flist.ContainerName]
	if container.Status == Stopped {
		if err := container.fs.Unmount(); err != nil {
			log.Println(err)
			w.reportFailureOperation()
		}
		if err := os.RemoveAll(container.Path); err != nil {
			log.Println(err)
			w.reportFailureOperation()
			return
		}

		w.reportSuccessOperation()
	}

	w.reportFailureOperation()
}