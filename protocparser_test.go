package main

import (
	"testing"
)

func TestStartComment(t *testing.T) {
	assertTrue(isStartComment("/**"), t)
	assertTrue(isStartComment("          /**"), t)
	assertFalse(isStartComment("/*"), t)
	assertFalse(isStartComment("//"), t)
}

func TestEndComment(t *testing.T) {
	assertTrue(isEndComment("*\\"), t)
	assertTrue(isEndComment("*\\    "), t)
	assertFalse(isEndComment("*\\    s"), t)
}

func TestIsRpc(t *testing.T) {
	assertTrue(isRpc("   RPC MyRPC"), t)
	assertFalse(isRpc("   RPC"), t)
}

func TestIsService(t *testing.T) {
	assertTrue(isService("   Service MyService"), t)
	assertFalse(isService("   Service"), t)
	assertFalse(isService("   ServiceWithTypo MyService"), t)
}

func assertTrue(expr bool, t *testing.T) {
	if !expr {
		t.Error("failure")
	}
}

func assertFalse(expr bool, t *testing.T) {
	if expr {
		t.Error("failure")
	}
}