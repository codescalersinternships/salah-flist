package main

import (
	"log"
	"syscall"
	"time"
)

// stop stops running container by sending SIGTERM signal,
// and after a grace period, SIGKILL to the process inside
// the container. this is server side of stop command, it
// carries the work of stopping entrypoint process inside
// the container. if command carried successfully, it sends SIGUSR1
// to requesting client, otherwise SIGUSR2 to represent failure.
func (w *Worker) stop() {
	if _, ok := w.Containers[w.Flist.ContainerName]; !ok {
		log.Printf("container name <%s> doesn't exist\n", w.Flist.ContainerName)
		w.reportFailureOperation()
		return
	}

	if w.Containers[w.Flist.ContainerName].Status != Running {
		log.Printf("this container <%s> is not running\n", w.Flist.ContainerName)
		w.reportFailureOperation()
		return
	}

	pid := w.Containers[w.Flist.ContainerName].Pid
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		log.Println(err)
		w.reportFailureOperation()
		return
	}
	time.Sleep(10 * time.Second)
	if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
		log.Println(err)
	}

	container := Container {
		Status: Stopped,
		Id: w.Containers[w.Flist.ContainerName].Id,
		FlistName: w.Containers[w.Flist.ContainerName].FlistName,
		Entrypoint: w.Containers[w.Flist.ContainerName].Entrypoint,
		Path: w.Containers[w.Flist.ContainerName].Path,
		Pid: w.Containers[w.Flist.ContainerName].Pid,
		fs: w.Containers[w.Flist.ContainerName].fs,
	}
	w.Containers[w.Flist.ContainerName] = container

	w.reportSuccessOperation()
}