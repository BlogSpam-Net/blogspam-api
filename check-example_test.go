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

		result, detail := validateEmail(Submission{Email: input})

		if result != Undecided {
			t.Errorf("Unexpected response: '%v'", result)
		}
		if len(detail) != 0 {
			t.Errorf("Unexpected response: '%v'", detail)
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

		result, detail := validateEmail(Submission{Email: input})
		if result != Spam {
			t.Errorf("Unexpected response: '%v'", result)
		}
		if len(detail) == 0 {
			t.Errorf("Unexpected response: '%v'", detail)
		}
	}
}
