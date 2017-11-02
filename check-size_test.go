//
// Test for our size-detecting plugin.
//

package main

import (
	"strings"
	"testing"
)

func TestSizeHam(t *testing.T) {

	//
	// Test a simple comment.
	//
	result := validateEmail(Submission{Comment: "I like to eat cakes",
		Options: "min-size=1,max-size=100"})
	if len(result) != 0 {
		t.Errorf("Unexpected response: '%v'", result)
	}
}

//
// Test broken size-options
//
func TestSizeBrokenOptions(t *testing.T) {

	//
	// Test several options
	//
	inputs := []string{"min-size=b",
		"min-size=-1",
		"max-size=pi",
		"max-size=-5"}

	for _, input := range inputs {

		result := validateSize(Submission{Comment: "Foo, bar", Options: input})
		if len(result) == 0 {
			t.Errorf("Unexpected response: '%v'", result)
		}
		if !strings.Contains(result, "Failed to parse") {
			t.Errorf("Unexpected response: '%v'", result)
		}
	}
}

//
// Test too small.
//
func TestSizeTooSmall(t *testing.T) {

	//
	// Test several options
	//
	inputs := []string{"moi",
		"hello"}

	for _, input := range inputs {

		result := validateSize(Submission{Comment: input,
			Options: "min-size=100"})
		if len(result) == 0 {
			t.Errorf("Unexpected response: '%v'", result)
		}
		if !strings.Contains(result, "minimum size") {
			t.Errorf("Unexpected response: '%v'", result)
		}
	}
}

//
// Test too large
//
func TestSizeTooLarge(t *testing.T) {

	//
	// Test several comments
	//
	inputs := []string{"This is a huge comment, honest",
		"I like big books and I cannot lie"}

	for _, input := range inputs {

		result := validateSize(Submission{Comment: input,
			Options: "max-size=10"})
		if len(result) == 0 {
			t.Errorf("Unexpected response: '%v'", result)
		}
		if !strings.Contains(result, "maximum size") {
			t.Errorf("Unexpected response: '%v'", result)
		}
	}
}
