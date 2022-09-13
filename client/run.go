package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

type ClientData struct {
	Mountpoint string	`json:"mountpoint"`
	Pid int				`json:"pid"`
}

// run mounts container and runs entrypoint process inside it. this is
// client side of run command.  it firstly sends command data to daemon
// in a request, which does the work of mounting container, then waits for a response
// from daemon to know whether the container was mounted successfully or not.
// if mounted successfully, client execute the entrypoint process isolated
// inside container.
func run(conn net.Conn) {
	clientData := ClientData {Pid: os.Getpid()}
	request, err := newRequest(clientData, os.Args[1], os.Args[2:]...)
	if err != nil {
		log.Println(err)
		return
	}
	
	if err := ConnectionWrite(conn, request); err != nil {
		log.Println(err)
		return
	}

	var response Response
	if err := ConnectionRead(conn, &response); err != nil {
		log.Println(err)
		if err := ConnectionErrorResponse(conn, err.Error()); err != nil {
			log.Println(err)
			return
		}
		return
	}

	if response.Status == Error {
		log.Println(response.ErrorMsg)
		return
	}
	
	if err := json.Unmarshal(response.Body, &clientData); err != nil {
		log.Println(err)
		if err := ConnectionErrorResponse(conn, err.Error()); err != nil {
			log.Println(err)
			return
		}
		return
	}

	runtime.LockOSThread()

	cmd := exec.Command(os.Args[3], os.Args[4:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Chroot: 		clientData.Mountpoint,

		Cloneflags: 	syscall.CLONE_NEWUTS | syscall.CLONE_NEWNET | syscall.CLONE_NEWNS |
						syscall.CLONE_NEWUSER | syscall.CLONE_NEWPID,

		Unshareflags: 	syscall.CLONE_NEWNS,

		UidMappings: 	[]syscall.SysProcIDMap{{
			ContainerID: 	0,
			HostID: 		syscall.Getuid(),
			Size: 			1,
		}},

		GidMappings: 	[]syscall.SysProcIDMap{{
			ContainerID: 	0,
			HostID: 		syscall.Getgid(),
			Size: 			1,
		}},

		Credential: 	&syscall.Credential{
			Uid: 			uint32(syscall.Getuid()),
			Gid: 			uint32(syscall.Getuid()),
		},

		Pdeathsig: 		syscall.SIGKILL,
	}
	cmd.Stdin  = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		log.Printf("Command %s returned error %v", os.Args[3], err)
		if err := ConnectionErrorResponse(conn, err.Error()); err != nil {
			log.Println(err)
			return
		}
		return
	}

	if err := ConnectionWrite(conn, Response {Status: Success}); err != nil {
		log.Println(err)
		return
	}

	runtime.UnlockOSThread()
}