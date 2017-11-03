//
// Test for our check-name plugin.
//

package main

import (
	"testing"
)

func TestHam(t *testing.T) {

	input := Submission{Name: "Steve Kemp"}

	result, detail := checkLinkName(input)
	if result != Undecided {
		t.Errorf("Unexpected response: %v", result)
	}
	if len(detail) != 0 {
		t.Errorf("Unexpected response: %v", detail)
	}
}

func TestSPAM(t *testing.T) {

	//
	// Test several link-based names.
	//
	inputs := []string{"http://example.com",
		"https://example.com",
		"HTTPS://string.com",
		"HtTPS://"}

	for _, input := range inputs {
		result, detail := checkLinkName(Submission{Name: input})
		if result != Spam {
			t.Errorf("Unexpected response: '%v'", result)
		}
		if len(detail) == 0 {
			t.Errorf("Unexpected response: '%v'", detail)
		}
	}
}
