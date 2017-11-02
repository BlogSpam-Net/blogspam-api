//
// Test for our check-name plugin.
//

package main

import (
	"testing"
)

func TestHam(t *testing.T) {

	input := Submission{Name: "Steve Kemp"}

	result := checkLinkName(input)
	if len(result) != 0 {
		t.Errorf("Unexpected response: %v", result)
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
		result := checkLinkName(Submission{Name: input})
		if len(result) == 0 {
			t.Errorf("Unexpected response: '%v'", result)
		}
	}
}
