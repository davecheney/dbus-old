package dbus

import (
	"testing"
)

func TestDial(t *testing.T) {
	c, err := DialUnix("/var/run/dbus/system_bus_socket")
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Auth(); err != nil {
		t.Fatal(err)
	}
	c.Close()
}	
