package ext

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// IP2Int ip to int
func IP2Int(ip string) int64 {
	array := strings.Split(ip, ".")
	if len(array) != 4 {
		return 0
	}

	A, err := strconv.Atoi(array[0])
	if err != nil {
		return 0
	}

	B, err := strconv.Atoi(array[1])
	if err != nil {
		return 0
	}

	C, err := strconv.Atoi(array[2])
	if err != nil {
		return 0
	}

	D, err := strconv.Atoi(array[3])
	if err != nil {
		return 0
	}

	return int64(((A*256+B)*256+C)*256 + D)
}

// Int2IP int to ip
func Int2IP(ip int64) string {
	ulMask := [4]int64{0x000000FF, 0x0000FF00, 0x00FF0000, 0xFF000000}
	var result [4]string
	for i := 0; i < 4; i++ {
		result[3-i] = strconv.FormatInt((ip&ulMask[i])>>(uint(i)*8), 10)
	}

	return strings.Join(result[:], ".")
}

// ServerIP external ip
func ServerIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}

	return "", errors.New("ARE YOU CONNECTED TO THE NETWORK?")
}

var (
	privateBlocks []*net.IPNet
)

func init() {
	for _, b := range []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"} {
		if _, block, err := net.ParseCIDR(b); err == nil {
			privateBlocks = append(privateBlocks, block)
		}
	}
}

// ExtractIP server ip
func ExtractIP(addr string) (string, error) {
	if len(addr) > 0 && (addr != "0.0.0.0" && addr != "[::]") {
		return addr, nil
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("Failed to get interface addresses! Err: %v", err)
	}

	var ipAddr []byte

	for _, rawAddr := range addrs {
		var ip net.IP
		switch addr := rawAddr.(type) {
		case *net.IPAddr:
			ip = addr.IP
		case *net.IPNet:
			ip = addr.IP
		default:
			continue
		}

		if ip.To4() == nil {
			continue
		}

		if !isPrivateIP(ip.String()) {
			continue
		}

		ipAddr = ip
		break
	}

	if ipAddr == nil {
		return "", fmt.Errorf("No private IP address found, and explicit IP not provided")
	}

	return net.IP(ipAddr).String(), nil
}

func isPrivateIP(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	for _, priv := range privateBlocks {
		if priv.Contains(ip) {
			return true
		}
	}
	return false
}
