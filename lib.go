package lib

import (
	"github.com/monaco-io/lib/assert"
	"github.com/monaco-io/lib/ip"
	"github.com/monaco-io/lib/sys"
)

// network
var (
	// ExternalIPv4 get public ipv4
	ExternalIPv4 = ip.ExternalV4

	// InternalIPv4 get local ipv4
	InternalIPv4 = ip.InternalV4

	// IPAtoI ip string to int
	IPAtoI = ip.AtoI

	// IPItoA ip int to string
	IPItoA = ip.ItoA
)

// assert
var (
	// OK assert and fatal
	OK = assert.OK

	// Equal a == b
	Equal = assert.Equal

	// DeepEqual a deepequal b
	DeepEqual = assert.DeepEqual
)

// sys
var (
	// ExitGrace exec callback functions before exit (SIGINT/SIGQUIT/SIGTERM)
	ExitGrace = sys.ExitGrace
)
