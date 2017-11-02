package main

import (
	"testing"
	"strings"
)

func TestIPBlacklistOK(t *testing.T) {

	//
	// Just submit something that is not blacklisted.
	//
	result := checkBlacklist(Submission{Email: "moi@exampl.fi",
	Subject: "Hello", IP: "127.0.0.1" })
	if len(result) != 0 {
		t.Errorf("Unexpected response: '%v'", result)
	}
}


func TestIPBlacklistOne(t *testing.T) {

	//
	// This should pass - the IP is outside the CIDR range
	//
	result := checkBlacklist(Submission{Email: "moi@exampl.fi",
	Subject: "Hello", IP: "10.20.30.48", Options: "blacklist=10.20.30.40/29" })
	if len(result) != 0 {
		t.Errorf("Unexpected response: '%v'", result)
	}
}


func TestIPBlacklistTwo(t *testing.T) {

	//
	// This should fail the source IP is inside the blacklist.
	//
	result := checkBlacklist(Submission{Email: "moi@exampl.fi",
	Subject: "Hello", IP: "10.20.30.47", Options: "blacklist=10.20.30.40/29" })
	if len(result) == 0 {
		t.Errorf("Unexpected response: '%v'", result)
	}
}

func TestIPBlacklistBogusCIDR(t *testing.T) {

	//
	// This should fail, as the CIDR is bogus.
	//
	result := checkBlacklist(Submission{Email: "moi@exampl.fi",
	Subject: "Hello", IP: "10.20.30.47", Options: "blacklist=10.20.30.40/329" })
	if len(result) == 0 {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if !strings.Contains(result, "parse CIDR") {
		t.Errorf("Unexpected response: '%v'", result)
	}

}
