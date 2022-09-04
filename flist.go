package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	flistsStorePath = "/var/lib/flist/store"
	flistsContainersPath = "/var/lib/flist/containers"
	flistsUnpackedPath = "/var/lib/flist/tmp"
	defaultStorageHubPath = "zdb://hub.grid.tf:9900"
)

func main() {
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	runMetaURL := runCmd.String("meta", "", "URL for flist meta file")
	runEntryPoint := runCmd.String("entrypoint", "", "set executable to run when container is initiated")

	if len(os.Args) < 2 {
        usage()
        os.Exit(1)
    }

	switch os.Args[1] {
	case "run":
		if err := runCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err := Run(*runMetaURL, *runEntryPoint)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {

}