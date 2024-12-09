package types

import (
	"fmt"
)

type Link struct {
	NodeFrom      string
	NodeTo        string
	InterfaceFrom string
	InterfaceTo   string
}

func (l Link) MarshalYAML() (interface{}, error) {
	return map[string]interface{}{
		"endpoints": []string{
			fmt.Sprintf("%s:%s", l.NodeFrom, l.InterfaceFrom),
			fmt.Sprintf("%s:%s", l.NodeTo, l.InterfaceTo),
		},
	}, nil
}
