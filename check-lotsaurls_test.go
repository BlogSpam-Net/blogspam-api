//
// Test for our hyperlink-counting plugin.
//

package main

import (
	"strings"
	"testing"
)

func TestHyperLinkHam(t *testing.T) {

	//
	// Test a simple comment.
	//
	result, detail := checkHyperlinkCounts(Submission{Comment: "I http://steve.fi/"})
	if result != Undecided {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if len(detail) != 0 {
		t.Errorf("Unexpected response: '%v'", detail)
	}
}

//
// Test broken options
//
func TestHyperLinkBadCount(t *testing.T) {

	//
	// Test several options
	//
	inputs := []string{"max-links=b",
		"max-links=-1"}

	for _, input := range inputs {

		result, detail := checkHyperlinkCounts(Submission{Comment: "Foo, bar", Options: input})
		if result != Error {
			t.Errorf("Unexpected response: '%v'", result)
		}
		if len(detail) == 0 {
			t.Errorf("Unexpected response: '%v'", detail)
		}
		if !strings.Contains(detail, "Failed to parse") {
			t.Errorf("Unexpected response: '%v'", detail)
		}
	}
}

//
// Test with > 10 links
//
func TestHyperLinkDefaults(t *testing.T) {

	//
	// Too many links.
	//
	input := "http://steve.com/ http://steve.com/ http://steve.com/ http://steve.com/ http://steve.com/ http://steve.com/ http://steve.com/ http://steve.com/ http://steve.com/ http://steve.com/ http://steve.com/ http://steve.com/ http://steve.com/"

	result, detail := checkHyperlinkCounts(Submission{Comment: input})
	if result != Spam {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if len(detail) == 0 {
		t.Errorf("Unexpected response: '%v'", detail)
	}
	if !strings.Contains(detail, "Too many") {
		t.Errorf("Unexpected response: '%v'", detail)
	}
}
