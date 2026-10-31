package main

import (
	"net"
	"testing"
)

func TestListenOnAvailablePortBindsLoopback(t *testing.T) {
	ln, err := listenOnAvailablePort(3000, 3999)
	if err != nil {
		t.Fatalf("listenOnAvailablePort: %v", err)
	}
	defer ln.Close()

	ip := ln.Addr().(*net.TCPAddr).IP
	if !ip.IsLoopback() {
		t.Errorf("bound to %s, want loopback address", ip)
	}
}
