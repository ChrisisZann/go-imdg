package comms

import (
	"fmt"
	"net"
	"strconv"
)

// NodeAddr implements net.Addr interface

// NodeAddr implements net.Addr interface
type NodeAddr struct {
	network string
	addr    string
}

func (m NodeAddr) Network() string {
	return m.network
}

func (m NodeAddr) String() string {
	return m.addr
}

// NodeAddr implements net.Addr interface

func NewNodeAddr(n, a string) (NodeAddr, error) {
	if err := validateNetwork(n); err != nil {
		return NodeAddr{}, err
	}
	if err := validateAddress(a); err != nil {
		return NodeAddr{}, err
	}
	return NodeAddr{
		network: n,
		addr:    a,
	}, nil
}

// validateNetwork checks if the network string is valid.
func validateNetwork(network string) error {
	if network != "tcp" {
		return fmt.Errorf("invalid network: only 'tcp' is supported")
	}
	return nil
}

// validateAddress checks if the address string is valid.
func validateAddress(address string) error {
	if address == "" {
		return fmt.Errorf("address cannot be empty")
	}

	// Check if the address is in the format IP:PORT
	var ip string
	var port int
	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("invalid address format: %v", err)
	}
	ip = host
	port, err = strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid port: %v", err)
	}

	if port <= 0 || port > 65535 {
		return fmt.Errorf("invalid port: must be between 1 and 65535, got: %d", port)
	}
	// If IP is "localhost", it's valid, but still check the port
	if ip == "localhost" {
		return nil
	}
	// Additional validation for non-"localhost" IPs can be added here if needed
	if ip != "localhost" {
		// Check if the IP is a valid IPv4 or IPv6 address
		if net.ParseIP(ip) == nil {
			return fmt.Errorf("invalid IP address: %s", ip)
		}
	}
	return nil
}
