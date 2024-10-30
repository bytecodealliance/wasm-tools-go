package sockets

import (
	"net/netip"
	"testing"
	"unsafe"

	instancenetwork "tests/generated/wasi/sockets/v0.2.0/instance-network"
	ipnamelookup "tests/generated/wasi/sockets/v0.2.0/ip-name-lookup"
)

func TestResolveAddresses(t *testing.T) {
	const hostname = "localhost"

	network := instancenetwork.InstanceNetwork()
	defer network.ResourceDrop()

	result := ipnamelookup.ResolveAddresses(network, hostname)
	if result.IsErr() {
		t.Errorf("ResolveAddresses: %v", result.Err())
		return
	}

	addrs := result.OK()
	defer addrs.ResourceDrop()

	pollable := addrs.Subscribe()
	defer pollable.ResourceDrop()
	pollable.Block()

	for {
		result := addrs.ResolveNextAddress()
		if result.IsErr() {
			t.Errorf("ResolveNextAddress: %v", result.Err())
			return
		}
		ok := result.OK()
		if ok.None() {
			break
		}
		ip := ok.Some()
		var s string
		if ipv4 := ip.IPv4(); ipv4 != nil {
			s = netip.AddrFrom4(*ipv4).String()
		} else if ipv6 := ip.IPv4(); ipv6 != nil {
			t.Log("ipv6")
			s = netip.AddrFrom16(*(*[16]byte)(unsafe.Pointer(ipv6))).String()
		}
		t.Logf("%s: %s", hostname, s)
	}
}
