//
// Test for our dronebl.org-lookup plugin.
//

package main

import (
	"testing"
)

//
// Test IPs that should never be listed.
//
func TestNonDrone(t *testing.T) {

	//
	// Test several IPs
	//
	inputs := []string{"127.0.0.1",
		"192.168.0.1",
		"10.11.12.13"}

	for _, input := range inputs {

		result, detail := checkDroneBlacklist(Submission{IP: input})

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
func TestDroneListed(t *testing.T) {

	result, detail := checkDroneBlacklist(Submission{IP: "116.255.241.111"})

	if result != Spam {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if len(detail) == 0 {
		t.Errorf("Unexpected response: '%v'", detail)
	}
}
