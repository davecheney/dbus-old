package dbus

import (
	"os"
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
	addrs := ParseAddress(os.Getenv("DBUS_SESSION_BUS_ADDRESS"))
	if len(addrs) < 1 {
		t.Fatalf("Could not parse DBUS_SESSION_BUS_ADDRESS")
	}
	c, err := addrs[0].Dial()
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Auth(); err != nil {
		t.Fatal(err)
	}
	c.Close()
}
