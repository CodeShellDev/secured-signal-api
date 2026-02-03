package netutils

import (
	"errors"
	"net"
	"strings"
)

func ParseIPorNet(str string) (*net.IPNet, error) {
    if !strings.Contains(str, "/") {
        ip := net.ParseIP(str)

        if ip == nil {
            return nil, errors.New("invalid ip: " + str)
        }

        var mask net.IPMask

        if ip.To4() != nil {
            mask = net.CIDRMask(32, 32) // IPv4 /32
        } else {
            mask = net.CIDRMask(128, 128) // IPv6 /128
        }
		
        return &net.IPNet{IP: ip, Mask: mask}, nil
    }

    _, network, err := net.ParseCIDR(str)

    if err != nil {
        return nil, err
    }

    return network, nil
}

func IsIPIn(ip net.IP, list []*net.IPNet) bool {
	for _, net := range list {
		if net.Contains(ip) {
			return true
		}
	}
	return false
}