//
//  Check for a repetitive link.
//

package main

import (
	"strings"
)

//
// Register ourself as a blogspam-plugin.
//
func init() {
	registerPlugin(BlogspamPlugin{Name: "33-link-body.js",
		Description: "Look for the link repeated in the body.",
		Author:      "Steve Kemp <steve@steve.org.uk>",
		Test:        checkRepetitiveLinks})
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
