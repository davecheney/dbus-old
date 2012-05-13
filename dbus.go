package dbus

// package dbus implements a dbus client described in
// http://dbus.freedesktop.org/doc/dbus-specification.html

import (
	"encoding/binary"
)

var endian = struct {
	binary.ByteOrder
	char byte
}{
	binary.BigEndian,
	'B',
}

type headerField struct {
	Name  byte
	Value interface{}
}

type header struct {
	Endianness byte
	Type       byte
	Flags      byte
	Version    byte
	BodyLen    uint32
	Serial     uint32
	Fields     []headerField
}
