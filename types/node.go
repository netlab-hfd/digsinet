package types

import "github.com/Lachstec/digsinet-ng/iface"

type Node struct {
	Name   string
	Kind   string
	Ifaces []iface.Iface
}
