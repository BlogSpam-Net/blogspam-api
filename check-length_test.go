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
	result := validateLength(Submission{Name: "Steve Kemp",
		Subject: "Hello, world!"})
	if len(result) != 0 {
		t.Errorf("Unexpected response: '%v'", result)
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
	result := validateLength(Submission{Name: name})
	if len(result) == 0 {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if !strings.Contains(result, "'name'") {
		t.Errorf("Unexpected response: '%v'", result)
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
	result := validateLength(Submission{Subject: subject})
	if len(result) == 0 {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if !strings.Contains(result, "'subject'") {
		t.Errorf("Unexpected response: '%v'", result)
	}
}
