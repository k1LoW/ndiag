package config

import (
	"fmt"
	"strings"
)

type Network struct {
	NetworkId string
	Route     []*Component
	Desc      string
}

func (n *Network) FullName() string {
	return fmt.Sprintf(n.NetworkId)
}

func (n *Network) Id() string {
	return strings.ToLower(n.FullName())
}

type rawNetwork struct {
	Id    string
	Route []string
	Desc  string
}
