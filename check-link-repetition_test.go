//
// Test for our link-repetition plugin.
//

package main

import (
	"testing"
)

func TestRepetitiveOK(t *testing.T) {

	result, detail := checkRepetitiveLinks(Submission{Link: "example.com",
		Comment: "This is Fine"})

	if result != Undecided {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if len(detail) != 0 {
		t.Errorf("Unexpected response: '%v'", detail)
	}
}

func TestLinkRepeated(t *testing.T) {

	result, detail := checkRepetitiveLinks(Submission{Link: "http://example.com/", Comment: "This is a comment with our same link in it http://example.com/"})
	if result != Spam {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if len(detail) == 0 {
		t.Errorf("Unexpected response: '%v'", detail)
	}
}
