package dbus

// package dbus implements a dbus client described in
// http://dbus.freedesktop.org/doc/dbus-specification.html

import (	
	"os"
	"net"
	"syscall"
	"fmt"
)

type Client struct {
	*net.UnixConn
}	

func DialUnix(path string) (*Client, error) {
	addr, err := net.ResolveUnixAddr("unix", path)
	conn, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		return nil, err
	}
	return NewClient(conn)
}

func NewClient(conn *net.UnixConn) (*Client, error) {
	return &Client{ conn }, nil
} 

func (c *Client) Auth() error {
	cred := syscall.UnixCredentials(&syscall.Ucred{
		Pid: int32(os.Getpid()),
		Uid: uint32(os.Getuid()),
		Gid: uint32(os.Getgid()),
	})	
	_, _, err := c.WriteMsgUnix([]byte{ 0 }, cred, nil)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(c, "AUTH EXTERNAL %x\r\n", fmt.Sprintf("%d", os.Getuid())); err != nil {
		return err
	}

	var code, cookie string
	if _, err := fmt.Fscanln(c, &code, &cookie); err != nil {
		return err
	}

	fmt.Println(code, cookie)
	
	if _, err := fmt.Fprintf(c, "BEGIN\r\n"); err != nil {
		return err
	}
	return nil
}
