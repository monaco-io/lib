package ip

import (
	"net"
	"strings"
)

// TODO tmp PIP
// ExternalV4 ...
func ExternalV4() (ips []string) {
	inters, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, inter := range inters {
		if !strings.HasPrefix(inter.Name, "lo") {
			addrs, err := inter.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok {
					if ipnet.IP.IsLoopback() ||
						ipnet.IP.IsLinkLocalMulticast() ||
						ipnet.IP.IsLinkLocalUnicast() {
						continue
					}
					if ip4 := ipnet.IP.To4(); ip4 != nil {
						switch true {
						case ip4[0] == 10:
							continue
						case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
							continue
						case ip4[0] == 192 && ip4[1] == 168:
							continue
						default:
							ips = append(ips, ipnet.IP.String())
						}
					}
				}
			}
		}
	}
	return
}

// InternalV4 ...
func InternalV4() (ipv4 string) {
	inters, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, inter := range inters {
		if !strings.HasPrefix(inter.Name, "lo") {
			addrs, err := inter.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ipv4 = ipnet.IP.String()
						return
					}
				}
			}
		}
	}
	return
}

// AtoI conver ip addr to int.
func AtoI(s string) (ipInt int) {
	ip := net.ParseIP(s)
	if ip == nil {
		return
	}
	ip = ip.To4()
	if ip == nil {
		return
	}
	ipInt += int(ip[0]) << 24
	ipInt += int(ip[1]) << 16
	ipInt += int(ip[2]) << 8
	ipInt += int(ip[3])
	return
}

// ItoA conver int to ip addr.
func ItoA(ipInt int) string {
	ip := make(net.IP, net.IPv4len)
	ip[0] = byte((ipInt >> 24) & 0xFF)
	ip[1] = byte((ipInt >> 16) & 0xFF)
	ip[2] = byte((ipInt >> 8) & 0xFF)
	ip[3] = byte(ipInt & 0xFF)
	return ip.String()
}
