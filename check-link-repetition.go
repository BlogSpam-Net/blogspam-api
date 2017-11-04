//
//  Check for a repetitive link.
//

package main

import (
	"strings"
)

func init() {
	//
	// Add our plugin-method.
	//
	var x = Plugins{Name: "33-link-body.js",
		Description: "Look for the link repeated in the body.",
		Author:      "Steve Kemp <steve@steve.org.uk>",
		Test:        checkRepetitiveLinks}
	registerPlugin(x)

}

//
// Block spam which has the "Link" repeated in the "Comment".
//
func checkRepetitiveLinks(x Submission) (PluginResult, string) {

	//
	// If we have no Link we cannot do a test
	//
	if len(x.Link) <= 0 {
		return Undecided, ""
	}

	//
	// Does the same link show up in the body?
	//
	if strings.Contains(x.Comment, x.Link) {
		return Spam, "Repetition of links"
	}

	//
	// Continue.
	//
	return Undecided, ""
}
