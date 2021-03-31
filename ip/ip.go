package ip

import (
	"fmt"
	"net"
)

// TODO tmp PIP
// ExternalV4 ...
func ExternalV4() (ips []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipnet.IP.IsLoopback() ||
				ipnet.IP.IsLinkLocalMulticast() ||
				ipnet.IP.IsLinkLocalUnicast() {
				continue
			}
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				if (ip4[0] == 10) ||
					(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) ||
					(ip4[0] == 192 && ip4[1] == 168) {
					continue
				}
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return
}

// InternalV4 ...
func InternalV4() (ipv4 string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipv4 = ipnet.IP.String()
				return
			}
		}
	}
	return
}

// AtoI conver ip addr to int.
func AtoI(ipStr string) (ipInt int) {
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		return
	}
	ipInt = int(ip[3]) | int(ip[2])<<8 | int(ip[1])<<16 | int(ip[0])<<24
	return
}

// ItoA conver int to ip addr.
func ItoA(ipInt int) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		uint8(ipInt>>24), uint8(ipInt>>16),
		uint8(ipInt>>8), uint8(ipInt),
	)
}
