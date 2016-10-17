package main

import (
	"testing"
)

func TestStartComment(t *testing.T) {
	assertTrue(startCommentRE.MatchString("/**"), t)
	assertTrue(startCommentRE.MatchString("          /**"), t)
	assertFalse(startCommentRE.MatchString("/*"), t)
	assertFalse(startCommentRE.MatchString("//"), t)
}

func TestEndComment(t *testing.T) {
	assertTrue(endCommentRE.MatchString("*\\"), t)
	assertTrue(endCommentRE.MatchString("*\\    "), t)
	assertFalse(endCommentRE.MatchString("*\\    s"), t)
}

func TestIsRpc(t *testing.T) {
	assertTrue(rpcRE.MatchString("   rpc MyRPC"), t)
	assertFalse(rpcRE.MatchString("   rpc"), t)
}

func TestIsService(t *testing.T) {
	assertTrue(serviceRE.MatchString("   service MyService"), t)
	assertFalse(serviceRE.MatchString("   service"), t)
	assertFalse(serviceRE.MatchString("   serviceWithTypo MyService"), t)
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