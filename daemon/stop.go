package main

import (
	"log"
	"syscall"
	"time"
)

func (w *Worker) stop() {
	if _, ok := w.Containers[w.Flist.ContainerName]; !ok {
		log.Printf("container name <%s> doesn't exist", w.Flist.ContainerName)
		w.reportFailureOperation()
		return
	}

	if w.Containers[w.Flist.ContainerName].Status != Running {
		log.Printf("this container <%s> is not running", w.Flist.ContainerName)
		w.reportFailureOperation()
		return
	}

	pid := w.Containers[w.Flist.ContainerName].Pid
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		log.Println(err)
		return
	}
	time.Sleep(10 * time.Second)
	if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
		log.Println(err)
	}

	container := Container {
		Status: Stopped,
		Id: w.Containers[w.Flist.ContainerName].Id,
		Path: w.Containers[w.Flist.ContainerName].Path,
		Pid: w.Containers[w.Flist.ContainerName].Pid,
		fs: w.Containers[w.Flist.ContainerName].fs,
	}
	w.Containers[w.Flist.ContainerName] = container

	w.reportSuccessOperation()
}

func (w *Worker) reportSuccessOperation() {
	if err := syscall.Kill(w.Flist.ProcessPid, syscall.SIGUSR1); err != nil {
		log.Println(err)
		return
	}
}

func (w *Worker) reportFailureOperation() {
	if err := syscall.Kill(w.Flist.ProcessPid, syscall.SIGUSR2); err != nil {
		log.Println(err)
		return
	}
}