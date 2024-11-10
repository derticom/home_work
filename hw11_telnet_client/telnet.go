package main

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/pkg/errors" //nolint: depguard // import is necessary
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TelClient struct {
	Address string
	Timeout time.Duration
	In      io.ReadCloser
	Out     io.Writer
	Conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelClient{
		Address: address,
		Timeout: timeout,
		In:      in,
		Out:     out,
	}
}

func (tc *TelClient) Connect() error {
	var err error
	tc.Conn, err = net.DialTimeout("tcp", tc.Address, tc.Timeout)
	if err != nil {
		return errors.Wrap(err, "failed to net.DialTimeout")
	}

	log.Printf("connected to %s\n", tc.Address)

	return nil
}

func (tc *TelClient) Send() error {
	_, err := io.Copy(tc.Conn, tc.In)
	if err != nil {
		return errors.Wrap(err, "failed to io.Copy")
	}

	return nil
}

func (tc *TelClient) Receive() error {
	_, err := io.Copy(tc.Out, tc.Conn)
	if err != nil {
		return errors.Wrap(err, "failed to io.Copy")
	}

	return nil
}

func (tc *TelClient) Close() error {
	if tc.Conn != nil {
		err := tc.Conn.Close()
		if err != nil {
			return errors.Wrap(err, "failed to Conn.Close")
		}
	}

	return nil
}
