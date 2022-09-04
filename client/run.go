package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	runCmd        = flag.NewFlagSet("run", flag.ExitOnError)
	runMetaURL    = runCmd.String("meta", "", "URL for flist meta file")
	runEntryPoint = runCmd.String("entrypoint", "", "set executable to run when container is initiated")
)

func run(conn net.Conn) {
	if err := runCmd.Parse(os.Args[2:]); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	writeData(conn, fmt.Sprintf("%s\n", runCmd.Name()))
	writeData(conn, fmt.Sprintf("%s\n", *runMetaURL))
	writeData(conn, fmt.Sprintf("%s\n", *runEntryPoint))
}