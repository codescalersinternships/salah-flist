package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"syscall"

	g8ufs "github.com/threefoldtech/0-fs"
	"github.com/threefoldtech/0-fs/meta"
	"github.com/threefoldtech/0-fs/storage"
)

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

// getFlist returns filepath of flist on disk, if file doesn't exist on disk
// downloads the flist from given metaURL and save to disk.
func getFlist(metaURL, fileName string) (string, error) {
	if err := os.MkdirAll(flistsStorePath, 0770); err != nil {
		return "", err
	}

	filePath := fmt.Sprintf("%s/%s", flistsStorePath, fileName)

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

	return fs, nil
}

// Run subcommand to execute binary at given entrypoint after mounting given flist
func Run(metaURL, entrypoint string) error {
	fileName, err := buildFileName(metaURL)
	if err != nil {
		return err
	}

	flistPath, err := getFlist(metaURL, fileName)
	if err != nil {
		return err
	}

	err = os.MkdirAll(flistsContainersPath, 0770)
	if err != nil {
		return err
	}

	containerDirPath, err := os.MkdirTemp(flistsContainersPath, fileName)
	if err != nil {
		return err
	}

	mountpoint := fmt.Sprintf("%s/%s", containerDirPath, "mnt")
	if err = os.MkdirAll(mountpoint, 0770); err != nil {
		return err
	}
	if err := os.MkdirAll(mountpoint + "/proc", 0777); err != nil {
		return err
	}
	if err := os.MkdirAll(mountpoint + "/tmp", 0777); err != nil {
		return err
	}

	fs, err := mountFlist(flistPath, fileName, containerDirPath, mountpoint)
	if err != nil {
		return err
	}

	// To-Do: binary execution
	go func () error {
		if err := syscall.Chroot(mountpoint); err != nil {
			log.Println(err)
			return err
		}
		if err := syscall.Chdir("/"); err != nil {
			log.Println(err)
			return err
		}

		if err := syscall.Mount("proc", "proc", "proc", 0, ""); err != nil {
			log.Println(err)
			return err
		}
		if err := syscall.Mount("tmp", "tmp", "tmpfs", 0, ""); err != nil {
			log.Println(err)
			return err
		}

		cmd := exec.Command(entrypoint)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	} ()

	
	fs.Wait()

	if err := syscall.Unmount("proc", 0); err != nil {
		return err
	}
	if err := syscall.Unmount("tmp", 0); err != nil {
		return err
	}
	if err := fs.Unmount(); err != nil {
		return err
	}

	return nil
}