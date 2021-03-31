package lib

import (
	"github.com/monaco-io/lib/ip"
)

// network
var (
	// ExternalIP get public ipv4
	ExternalIP = ip.ExternalV4

	// InternalIPv4 get local ipv4
	InternalIPv4 = ip.InternalV4

	// IPAtoI ip string to int
	IPAtoI = ip.AtoI

	// IPItoA ip int to string
	IPItoA = ip.ItoA
)
