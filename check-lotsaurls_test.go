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
	result := checkHyperlinkCounts(Submission{Comment: "I http://steve.fi/"})
	if len(result) != 0 {
		t.Errorf("Unexpected response: '%v'", result)
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

		result := checkHyperlinkCounts(Submission{Comment: "Foo, bar", Options: input})
		if len(result) == 0 {
			t.Errorf("Unexpected response: '%v'", result)
		}
		if !strings.Contains(result, "Failed to parse") {
			t.Errorf("Unexpected response: '%v'", result)
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

	result := checkHyperlinkCounts(Submission{Comment: input })
	if len(result) == 0 {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if !strings.Contains(result, "Too many") {
		t.Errorf("Unexpected response: '%v'", result)
	}
}
