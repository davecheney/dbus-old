package dbus

// bus address parsing support

import (
	"strings"
)

// Address represents a remote D-Bus server.
type Address interface {

	// Dial initiates a connection to the specified D-Bus address.
	Dial() (*Client, error)

	// parse parses the remainder of the address specification 
	// according to the addresses rules.
	parse(string) bool
}

// ParseAddress parses a D-BUS connection string and returns a slice
// of potential server addresses.
func ParseAddress(addr string) (addrs []Address) {
	for _, a := range strings.Split(addr, ";") {
		if tr := strings.Split(a, ":"); len(tr) == 2 {
			var addr Address
			switch tr[0] {
			case "unix":
				addr = new(unixAddress)
			default:
				continue
			}
			if addr.parse(tr[1]) {
				addrs = append(addrs, addr)
			}
		}
	}
	return
}

type unixAddress struct {
	path string
	guid string // TODO(dfc) should this be an integer ?
}

func (a *unixAddress) parse(rest string) bool {
	for _, kv := range strings.Split(rest, ",") {
		if k := strings.Split(kv, "="); len(k) == 2 {
			switch k[0] {
			case "path":
				if a.path != "" {
					return false
				}
				a.path = k[1]
			case "abstract":
				if a.path != "" {
					return false
				}
				a.path = "@" + k[1]
			case "guid":
				a.guid = k[1]
			default:
				// unknown key
				return false
			}
		} else {
			return false
		}
	}
	return a.path != ""
}

func (a *unixAddress) Dial() (*Client, error) {
	transport, err := dialUnix(a.path)
	return &Client{transport}, err
}
