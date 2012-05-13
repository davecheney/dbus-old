package dbus

// package dbus implements a dbus client described in
// http://dbus.freedesktop.org/doc/dbus-specification.html

import (	
	"io"
	"os"
	"net"
	"syscall"
	"fmt"
)

type transport interface {
	// Auth performs the AUTH handshake.
	Auth() error 

	io.Closer
}

type unixTransport struct {
	*net.UnixConn
}

func dialUnix(path string) (transport, error) {
	addr, err := net.ResolveUnixAddr("unix", path)
	conn, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		return nil, err
	}
	return &unixTransport{ conn }, nil
}

func (t *unixTransport) Auth() error {
	cred := syscall.UnixCredentials(&syscall.Ucred{
		Pid: int32(os.Getpid()),
		Uid: uint32(os.Getuid()),
		Gid: uint32(os.Getgid()),
	})	
	_, _, err := t.WriteMsgUnix([]byte{ 0 }, cred, nil)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(t, "AUTH EXTERNAL %x\r\n", fmt.Sprintf("%d", os.Getuid())); err != nil {
		return err
	}

	var code, cookie string
	if _, err := fmt.Fscanln(t, &code, &cookie); err != nil {
		return err
	}

	fmt.Println(code, cookie)
	
	if _, err := fmt.Fprintf(t, "BEGIN\r\n"); err != nil {
		return err
	}
	return nil
}
