package main

import (
	"encoding/json"
	"net"
)

// ConnectionWrite writes data to connection
func ConnectionWrite(conn net.Conn, data any) error {
	enc := json.NewEncoder(conn)
	if err := enc.Encode(data); err != nil {
		return err
	}

	return nil
}

// ConnectionRead reads data from connection into buf buffer
func ConnectionRead(conn net.Conn, buf any) error {
	dec := json.NewDecoder(conn)

	if err := dec.Decode(buf); err != nil {
		return err
	}

	return nil
}

// ConnectionErrorResponse writes Response with Status Error
// and sets ErrorMsg field to msg argument
func ConnectionErrorResponse(conn net.Conn, msg string) error {
	err := ConnectionWrite(conn, Response {
		Status: 	Error,
		ErrorMsg: 	msg,
	})
	if err != nil {
		return err
	}

	return nil
}