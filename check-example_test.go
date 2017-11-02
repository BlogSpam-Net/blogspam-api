//
// Test for our example-detecting plugin.
//

package main

import (
	"testing"
)

func TestExampleMailHam(t *testing.T) {

	//
	// Test several emails.
	//
	inputs := []string{"steve@steve.fi",
		""}

	for _, input := range inputs {

		result := validateEmail(Submission{Email: input})
		if len(result) != 0 {
			t.Errorf("Unexpected response: '%v'", result)
		}
	}
}

func TestExampleSPAM(t *testing.T) {

	//
	// Test several emails.
	//
	inputs := []string{"sefdsfdsf@example.io",
		"bosfsdf@example.com"}

	for _, input := range inputs {

		result := validateEmail(Submission{Email: input})
		if len(result) == 0 {
			t.Errorf("Unexpected response: '%v'", result)
		}
	}
}
