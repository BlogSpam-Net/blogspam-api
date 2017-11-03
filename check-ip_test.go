package main

import (
	"strings"
	"testing"
)

func TestIPBlacklistOK(t *testing.T) {

	//
	// Just submit something that is not blacklisted.
	//
	result, detail := checkBlacklist(Submission{Email: "moi@exampl.fi",
		Subject: "Hello", IP: "127.0.0.1"})
	if result != Undecided {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if len(detail) != 0 {
		t.Errorf("Unexpected response: '%v'", detail)
	}
}

func TestIPBlacklistOne(t *testing.T) {

	//
	// This should pass - the IP is outside the CIDR range
	//
	result, detail := checkBlacklist(Submission{Email: "moi@exampl.fi",
		Subject: "Hello", IP: "10.20.30.48", Options: "blacklist=10.20.30.40/29"})
	if result != Undecided {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if len(detail) != 0 {
		t.Errorf("Unexpected response: '%v'", detail)
	}
}

func TestIPBlacklistTwo(t *testing.T) {

	//
	// This should fail the source IP is inside the blacklist.
	//
	result, detail := checkBlacklist(Submission{Email: "moi@exampl.fi",
		Subject: "Hello", IP: "10.20.30.47", Options: "blacklist=10.20.30.40/29"})
	if result != Spam {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if len(detail) == 0 {
		t.Errorf("Unexpected response: '%v'", detail)
	}
}

func TestIPBlacklistBogusCIDR(t *testing.T) {

	//
	// This should fail, as the CIDR is bogus.
	//
	result, detail := checkBlacklist(Submission{Email: "moi@exampl.fi",
		Subject: "Hello", IP: "10.20.30.47", Options: "blacklist=10.20.30.40/329"})
	if result != Error {
		t.Errorf("Unexpected result: '%v'", result)
	}
	if len(detail) == 0 {
		t.Errorf("Unexpected response: '%v'", detail)
	}
	if !strings.Contains(detail, "parse CIDR") {
		t.Errorf("Unexpected response: '%v'", detail)
	}

}
