package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
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
func getFlist(metaURL string) (string, error) {
	fileName, err := buildFileName(metaURL)
	if err != nil {
		return "", err
	}

	filePath := fmt.Sprintf("%s/%s", flistsPath, fileName)

	if _, err := os.Stat(filePath); err != nil {
		log.Printf("flist %s doesn't exist and needs to get downloaded from %s\n", fileName, metaURL)
		err := downloadFlist(metaURL, filePath)
		if err != nil {
			return "", err
		}
	}

	return filePath, nil
}

// Run subcommand to execute binary at given entrypoint after mounting given flist
func Run(metaURL string, entrypoint string) error {
	filePath, err := getFlist(metaURL)
	if err != nil {
		return err
	}


	return nil
}