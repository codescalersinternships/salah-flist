package main

import (
	"encoding/json"
	"flag"
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

	flist := new()
	flist.Command = runCmd.Name()
	flist.MetaURL = *runMetaURL
	flist.Entrypoint = *runEntryPoint

	data, err := json.Marshal(flist)
	if err != nil {
		log.Fatal(err)
	}

	writeData(conn, data)
}