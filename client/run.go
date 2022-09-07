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

var (
	DaemonPid = 0
	
	BufSize = 1024
)

func run(conn net.Conn) {
	flist := new(os.Args[1], os.Args[2], os.Args[3], "", os.Args[4:]...)

	flistData, err := json.Marshal(flist)
	if err != nil {
		log.Fatal(err)
	}
	writeData(conn, flistData)
	
	DaemonPid = getDaemonPid(conn)

	runtime.LockOSThread()

	cmd := exec.Command(flist.Entrypoint, flist.Arg...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Chroot: flist.Mountpoint,
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWNET | syscall.CLONE_NEWNS |
					syscall.CLONE_NEWUSER | syscall.CLONE_NEWPID,
		Unshareflags: syscall.CLONE_NEWNS,
		UidMappings: []syscall.SysProcIDMap{{
			ContainerID: 0,
			HostID: syscall.Getuid(),
			Size: 1,
		}},
		GidMappings: []syscall.SysProcIDMap{{
			ContainerID: 0,
			HostID: syscall.Getgid(),
			Size: 1,
		}},
		Credential: &syscall.Credential{
			Uid: uint32(syscall.Getuid()),
			Gid: uint32(syscall.Getuid()),
		},
		Pdeathsig: syscall.SIGKILL,
	}
	cmd.Stdin  = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	<-done
	
	if err := cmd.Run(); err != nil {
		log.Printf("Command %s returned error %v", flist.Entrypoint, err)
	}

	runtime.UnlockOSThread()

	tellDaemonToCleanupContainer()
}