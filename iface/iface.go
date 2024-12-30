package iface

type Iface interface {
	GetName() string
	SetConfig(map[string]string)

	StartIface()
	StopIface()
}
