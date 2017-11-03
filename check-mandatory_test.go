//
// Test for our mandatory-fields plugin.
//

package main

import (
	"strings"
	"testing"
)

//
// Test that all is OK when the minimum required fields are present.
//
func TestMandatoryOK(t *testing.T) {

	result, _ := validateMandatory(Submission{Site: "example",
		IP: "1.2.3.4", Comment: "This is a test"})
	if result != Undecided {
		t.Errorf("Unexpected response: '%v'", result)
	}
}

//
// Test that all is OK when the minimum required fields are present, as well
// as a single extra mandatory field
//
func TestMandatoryOKExtra(t *testing.T) {

	result, _ := validateMandatory(Submission{Site: "example",
		IP: "1.2.3.4", Comment: "This is a test",
		Options: "mandatory=agent", Agent: "foo"})
	if result != Undecided {
		t.Errorf("Unexpected response: '%v'", result)
	}
}

//
// Test that we receive an alert if we're missing a `site` parameter.
//
func TestMandatoryMissingSite(t *testing.T) {

	result, detail := validateMandatory(Submission{Site: "",
		IP: "1.2.3.4", Comment: "This is a test"})

	if result != Spam {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if !strings.Contains(detail, "is missing") {
		t.Errorf("Unexpected response: '%v'", detail)
	}
}

//
// Test that we receive an alert if we're missing a `ip` parameter.
//
func TestMandatoryMissingIP(t *testing.T) {

	result, detail := validateMandatory(Submission{Site: "steve.fi",
		IP: "", Comment: "This is a test"})

	if result != Spam {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if !strings.Contains(detail, "is missing") {
		t.Errorf("Unexpected response: '%v'", detail)
	}
}

//
// Test that we receive an alert if we're missing a `comment` parameter.
//
func TestMandatoryMissingComment(t *testing.T) {

	result, detail := validateMandatory(Submission{Site: "fsdf",
		IP: "1.2.3.4"})

	if result != Spam {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if !strings.Contains(detail, "is missing") {
		t.Errorf("Unexpected response: '%v'", detail)
	}
}

//
// Test that we receive an alert if we're missing an extra `agent` parameter.
//
func TestMandatoryMissingAgent(t *testing.T) {

	result, detail := validateMandatory(Submission{Site: "fsdf",
		IP:      "1.2.3.4",
		Options: "mandatory=agent"})

	if result != Spam {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if !strings.Contains(detail, "is missing") {
		t.Errorf("Unexpected response: '%v'", detail)
	}
}
