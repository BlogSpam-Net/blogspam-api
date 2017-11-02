//
//  Check that the incoming submission has an MX-record for the specified
// email-address
//

package main

import (
	"fmt"
	"net"
	"regexp"
)

func init() {
	//
	// Add our plugin-method.
	//
	var x = Plugins{Name: "requiremx",
		Description: "Validates that an incoming submission has an MX record",
		Author:      "Steve Kemp <steve@steve.org.uk>",
		Test:        validateMX}
	plugins = append(plugins, x)

}

//
// Test that the email-field is non-empty and contains an MX-record
//
func validateMX(x Submission) string {

	//
	// If we have no email-address we cannot do an MX-lookup.
	//
	if len(x.Email) <= 0 {
		return ""
	}

	//
	// Get the email-address
	//
	re := regexp.MustCompile("^.*@([^@]+)$")
	match := re.FindStringSubmatch(x.Email)

	//
	// If that worked.
	//
	if len(match) > 0 {

		//
		// Lookup the MX-record of the domain.
		//
		// We're only looking for an error-here.
		//
		_, err := net.LookupMX(match[1])

		if err != nil {
			return (fmt.Sprintf("Failed to lookup MX-record of %s", match[1]))
		}
	}
	return ""
}
