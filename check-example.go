//
//  Check for "@example" email-addresses.
//

package main

import (
	"regexp"
	"strings"
)

func init() {
	//
	// Add our plugin-method.
	//
	var x = Plugins{Name: "10-example.js",
		Description: "Look for example-domains in emails",
		Author:      "Steve Kemp <steve@steve.org.uk>",
		Test:        validateEmail}
	registerPlugin(x)
}

//
// Test that the email-field is non-empty and contains a non-example entry.
//
func validateEmail(x Submission) (PluginResult, string) {

	//
	// If we have no email-address we cannot do a test
	//
	if len(x.Email) <= 0 {
		return Undecided, ""
	}

	//
	// Get the domain from the email-address.
	//
	re := regexp.MustCompile("^.*@([^@]+)$")
	match := re.FindStringSubmatch(x.Email)

	//
	// If that worked.
	//
	if len(match) > 0 {

		//
		// Does it start with example?
		//
		if strings.HasPrefix(strings.ToLower(match[1]), "example.") {
			return Spam, "Example-based email-address"
		}
	}
	return Undecided, ""
}
