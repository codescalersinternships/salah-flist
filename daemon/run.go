package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"

	g8ufs "github.com/threefoldtech/0-fs"
	"github.com/threefoldtech/0-fs/meta"
	"github.com/threefoldtech/0-fs/storage"
)

// buildFileName builds flist's file name from metaURL
func buildFileName(metaURL string) (string, error) {
	u, err := url.Parse(metaURL)
	if err != nil {
		return "", err
	}

	urlPath := u.Path
	pathSlice := strings.Split(urlPath, "/")

	return pathSlice[len(pathSlice)-1], nil
}

// downloadFlist downloads flist form metaURL into filepath
func downloadFlist(metaURL, filePath string) error {
	resp, err := http.Get(metaURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dest, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dest.Close()

	size, err := io.Copy(dest, resp.Body)
	if err != nil {
		if err := os.Remove(filePath); err != nil {
			return err
		}
		return err

	} else if resp.ContentLength != -1 && size != resp.ContentLength {
		if err := os.Remove(filePath); err != nil {
			return err
		}
		return fmt.Errorf("error happened while copying data into disk")
	}

	return nil
}

// getFlist returns filepath of flist on disk, if file doesn't exist on disk
// downloads the flist from given metaURL and save to disk.
func getFlist(metaURL, fileName string) (string, error) {
	if err := os.MkdirAll(StorePath, 0770); err != nil {
		return "", err
	}

	filePath := fmt.Sprintf("%s/%s", StorePath, fileName)

	if _, err := os.Stat(filePath); err != nil {
		log.Printf("flist %s doesn't exist and needs to get downloaded from %s\n", fileName, metaURL)
		err := downloadFlist(metaURL, filePath)
		if err != nil {
			return "", err
		}
	}

	return filePath, nil
}

// unpackFlistArchive unpacks a tgz flist (archive) from flistPath
// to a tmp location at "/var/lib/flist/tmp/"
func unpackFlistArchive(flistPath, fileName string) (string, error) {
	f, err := os.Open(flistPath)
	if err != nil {
		return "", err
	}

	tmpFlistDir := fmt.Sprintf("%s/%s", flistsUnpackedPath, fileName)

	err = os.MkdirAll(tmpFlistDir, 0770)
	if err != nil {
		return "", err
	}

	if err := meta.Unpack(f, tmpFlistDir); err != nil {
		return "", err
	}
	return tmpFlistDir, nil
}

// mountFlist mounts flist stored on disk at flistPath, then
// runs the executable at entrypoint
func mountFlist(flistPath, fileName, containerDirPath, mountpoint string) (*g8ufs.G8ufs, error) {
	tmpFlistDir, err := unpackFlistArchive(flistPath, fileName)
	if err != nil {
		return nil, err
	}

	metaStore, err := meta.NewStore(tmpFlistDir)
	if err != nil {
		return nil, err
	}

	storageHub, err := storage.NewSimpleStorage(defaultStorageHubPath)
	if err != nil {
		return nil, err
	}

	opt := g8ufs.Options {
		Backend: fmt.Sprintf("%s/%s", containerDirPath, "backend"),
		Target: mountpoint,
		Store: metaStore,
		Storage: storageHub,
		Reset: true,
	}

	fs, err := g8ufs.Mount(&opt)
	if err != nil {
		return nil, err
	}

	procPath := fmt.Sprintf("%s/proc", mountpoint)
	if err := os.MkdirAll(procPath, 0777); err != nil {
		return nil, err
	}
	tmpPath := fmt.Sprintf("%s/tmp", mountpoint)
	if err := os.MkdirAll(tmpPath, 0777); err != nil {
		return nil, err
	}

	if err := syscall.Mount(procPath, procPath, "proc", 0, ""); err != nil {
		return nil, err
	}
	if err := syscall.Mount(tmpPath, tmpPath, "tmpfs", 0, ""); err != nil {
		return nil, err
	}

	return fs, nil
}

// run mounts container and runs entrypoint process inside it. this is
// server side of run command, it carries the work of mounting flist,
// after mount success it sends Response message with Success Status, otherwise,
// it sends Response message with Error status to represent failure.
func (s *Server) run(conn Connection, request Request) {
	var container Container
	if err := json.Unmarshal(request.Body, &container); err != nil {
		log.Println(err)
		if err := conn.SendErrorResponse(err.Error()); err != nil {
			log.Println(err)
			return
		}
		return
	}
	container.MetaURL 		= request.Args[0]
	container.Entrypoint 	= request.Args[1]
	container.Args 			= request.Args[2:]
	
	fileName, err := buildFileName(container.MetaURL)
	if err != nil {
		log.Println(err)
		if err := conn.SendErrorResponse(err.Error()); err != nil {
			log.Println(err)
			return
		}
		return
	}

	flistPath, err := getFlist(container.MetaURL, fileName)
	if err != nil {
		log.Println(err)
		if err := conn.SendErrorResponse(err.Error()); err != nil {
			log.Println(err)
			return
		}
		return
	}

	if err := os.MkdirAll(ContainersPath, 0770); err != nil {
		log.Println(err)
		if err := conn.SendErrorResponse(err.Error()); err != nil {
			log.Println(err)
			return
		}
		return
	}

	containerDirPath, err := os.MkdirTemp(ContainersPath, fileName)
	if err != nil {
		log.Println(err)
		if err := conn.SendErrorResponse(err.Error()); err != nil {
			log.Println(err)
			return
		}
		return
	}
	pathSlice := strings.Split(containerDirPath, "/")
	containerId := pathSlice[len(pathSlice)-1]

	mountpoint := fmt.Sprintf("%s/%s", containerDirPath, "mnt")
	if err = os.MkdirAll(mountpoint, 0770); err != nil {
		log.Println(err)
		if err := conn.SendErrorResponse(err.Error()); err != nil {
			log.Println(err)
			return
		}
		return
	}

	fs, err := mountFlist(flistPath, fileName, containerDirPath, mountpoint)
	if err != nil {
		log.Println(err)
		if err := conn.SendErrorResponse(err.Error()); err != nil {
			log.Println(err)
			return
		}		
		return
	}

	container.Status 	 = Running
	container.Id	 	 = containerId
	container.FlistName  = fileName
	container.Path		 = containerDirPath
	container.Mountpoint = mountpoint
	container.Fs		 = fs

	s.Containers[container.Id] = container

	defer log.Printf("Container at %v unmounted successfully", containerDirPath)
	defer container.Fs.Unmount()

	responseBody, err := json.Marshal(container)
	if err != nil {
		log.Println(err)
		return
	}
	response := Response {
		Status: Success,
		Body: 	json.RawMessage(responseBody),
	}
	if err := conn.SendResponse(response); err != nil {
		log.Println(err)
		return
	}

	container.Status = Stopped
	if err := conn.ReadResponse(&response); err != nil {
		log.Println(err)
		s.Containers[container.Id] = container
		return
	}

	s.Containers[container.Id] = container
}