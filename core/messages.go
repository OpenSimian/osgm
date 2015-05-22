package core

import "github.com/M-O-S-E-S/mgm/mgm"

// NetworkMessage is a wrapper for sending multiple message types accross a single wire
type NetworkMessage struct {
	MessageType string
	HStats      mgm.HostStat `json:",omitempty"`
	Request     string       `json:",omitempty"`
}
