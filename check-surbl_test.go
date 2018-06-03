//
// Test for our surbl.org-lookup plugin.
//

package main

import (
	"testing"
)

//
// Test IPs that should never be listed.
//
func TestNonSurbl(t *testing.T) {

	result, detail := checkSurblBlacklist(Submission{Comment: "Moi kissa, no URLs here steve.fi/anal.rape https://gibberish.steve.fi/fuck.you"})

	if result != Undecided {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if len(detail) != 0 {
		t.Errorf("Unexpected response: '%v'", detail)
	}
}

//
// Test an example that is currently listed in Surbl
//
func TestSurblListed(t *testing.T) {
	//	return
	result, detail := checkSurblBlacklist(Submission{Comment: "Listed link: http://pornapps.xblog.in"})

	if result != Spam {
		t.Errorf("Unexpected response: '%v'", result)
	}
	if len(detail) == 0 {
		t.Errorf("Unexpected response: '%v'", detail)
	}
}
