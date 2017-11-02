//
// Test for our check-mx plugin.
//

package main

import (
	"testing"
)

func TestMXHam(t *testing.T) {

	//
	// Test several emails.
	//
	inputs := []string{"steve@steve.fi",
		""}

	for _, input := range inputs {

		result := validateMX(Submission{Email: input})
		if len(result) != 0 {
			t.Errorf("Unexpected response: '%v'", result)
		}
	}
}

func TestMXSPAM(t *testing.T) {

	//
	// Test several emails.
	//
	inputs := []string{"sefdsfdsf@fdsdfsdfs.dfsdcom",
		"bosfsdf@exmpla.edfsfsdf.com"}

	for _, input := range inputs {

		result := validateMX(Submission{Email: input})
		if len(result) == 0 {
			t.Errorf("Unexpected response: '%v'", result)
		}
	}
}
