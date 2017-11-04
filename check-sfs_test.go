//
// Test for our stopforumspam.com-lookup plugin.
//

package main

import (
	"testing"
)

//
// Test IPs that should never be listed.
//
func TestNonSpammer(t *testing.T) {

	//
	// Test several IPs
	//
	inputs := []string{"127.0.0.1", "192.168.0.1", "x"}

	for _, input := range inputs {

		result, detail := checkSFSBlacklist(Submission{IP: input})

		if result != Undecided {
			t.Errorf("Unexpected response: '%v'", result)
		}
		if len(detail) != 0 {
			t.Errorf("Unexpected response: '%v'", detail)
		}
	}
}

//
// Test a known-bad IP
//
func TestSFSListed(t *testing.T) {

	result, detail := checkSFSBlacklist(Submission{IP: "37.115.125.139"})

	if result != Spam {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if len(detail) == 0 {
		t.Errorf("Unexpected response: '%v'", detail)
	}
}
