//
// Test for our blacklisted-fields plugin.
//

package main

import (
	"testing"
)

//
// Valid comment.
//
func TestNotBlacklistedField(t *testing.T) {

	//
	// Test several Links.
	//
	inputs := []string{"https://steve.fi",
		"http://steve.fi"}

	for _, input := range inputs {

		result, detail := checkBlacklistedFields(Submission{Link: input})

		if result != Undecided {
			t.Errorf("Unexpected response: '%v'", result)
		}
		if len(detail) != 0 {
			t.Errorf("Unexpected response: '%v'", detail)
		}
	}
}

//
// These are each blacklisted
//
func TestBlacklistedLinkFields(t *testing.T) {

	inputs := []string{"https://zuschool.com/",
		"https://pcgle.com/moi"}

	for _, input := range inputs {

		result, detail := checkBlacklistedFields(Submission{Link: input})
		if result != Spam {
			t.Errorf("Unexpected response: '%v'", result)
		}
		if len(detail) == 0 {
			t.Errorf("Unexpected response: '%v'", detail)
		}
	}
}
