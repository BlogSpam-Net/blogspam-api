//
// Test for our length-detecting plugin.
//

package main

import (
	"strings"
	"testing"
)

func TestLengthHam(t *testing.T) {

	//
	// Test a simple comment.
	//
	result, detail := validateLength(Submission{Name: "Steve Kemp",
		Subject: "Hello, world!"})
	if result != Undecided {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if len(detail) != 0 {
		t.Errorf("Unexpected response: '%v'", detail)
	}
}

//
// Test too long name
//
func TestLengthName(t *testing.T) {

	name := "Steve"
	for len(name) < 200 {
		name = name + " "
	}
	result, detail := validateLength(Submission{Name: name})
	if result != Spam {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if len(detail) == 0 {
		t.Errorf("Unexpected response: '%v'", detail)
	}
	if !strings.Contains(detail, "'name'") {
		t.Errorf("Unexpected response: '%v'", detail)
	}
}

//
// Test too long "subject"
//
func TestLengthSubject(t *testing.T) {

	subject := "Steve"
	for len(subject) <= 200 {
		subject = subject + " "
	}
	result, detail := validateLength(Submission{Subject: subject})
	if result != Spam {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if len(detail) == 0 {
		t.Errorf("Unexpected response: '%v'", detail)
	}
	if !strings.Contains(detail, "'subject'") {
		t.Errorf("Unexpected response: '%v'", detail)
	}
}
