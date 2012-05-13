package dbus

import (
	"testing"
)

func TestDialUnixSystemBus(t *testing.T) {
	c, err := dialUnix("/var/run/dbus/system_bus_socket")
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Auth(); err != nil {
		t.Fatal(err)
	}
	c.Close()
}	

func TestDialUnixSessionBus(t *testing.T) {
	// TODO(dfc) hard coded, need to decode DBUS_SESSION_BUS_ADDRESS
	c, err := dialUnix("@/tmp/dbus-dsfqbFM8Bp")
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Auth(); err != nil {
		t.Fatal(err)
	}
	c.Close()
}	

