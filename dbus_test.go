package dbus

import (
	"testing"
	"os"
)

func TestDialSystemBus(t *testing.T) {
	c, err := DialUnix("/var/run/dbus/system_bus_socket")
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Auth(); err != nil {
		t.Fatal(err)
	}
	c.Close()
}	

func TestDialSessionBus(t *testing.T) {
	// TODO(dfc) hard coded, need to decode DBUS_SESSION_BUS_ADDRESS
	c, err := DialUnix("@/tmp/dbus-dsfqbFM8Bp")
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Auth(); err != nil {
		t.Fatal(err)
	}
	c.Close()
}	

