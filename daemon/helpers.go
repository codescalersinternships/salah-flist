package main

import (
	"encoding/json"
)

// SendRequest writes request data to connection
func (c *Connection) SendRequest(req Request) error {
	enc := json.NewEncoder(c.Conn)
	if err := enc.Encode(req); err != nil {
		return err
	}

	return nil
}

// ReadRequest reads request data from connection into buf buffer
func (c *Connection) ReadRequest(buf *Request) error {
	dec := json.NewDecoder(c.Conn)

	if err := dec.Decode(buf); err != nil {
		return err
	}

	return nil
}

// SendResponse writes response data to connection
func (c *Connection) SendResponse(res Response) error {
	enc := json.NewEncoder(c.Conn)
	if err := enc.Encode(res); err != nil {
		return err
	}

	return nil
}

// ReadResponse reads response data from connection into buf buffer
func (c *Connection) ReadResponse(buf *Response) error {
	dec := json.NewDecoder(c.Conn)

	if err := dec.Decode(buf); err != nil {
		return err
	}

	return nil
}

// SendErrorResponse writes Response with Status Error
// and sets ErrorMsg field to msg argument
func (c *Connection) SendErrorResponse(msg string) error {
	err := c.SendResponse(Response {
		Status: 	Error,
		ErrorMsg: 	msg,
	})
	if err != nil {
		return err
	}

	return nil
}