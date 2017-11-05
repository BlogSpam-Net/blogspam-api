//
// Check that the "Name" field of our incoming submission doesn't contain
// a hyperlink
//

package main

import "strings"

//
// Register ourself as a blogspam-plugin.
//
func init() {
	registerPlugin(BlogspamPlugin{Name: "35-name.js",
		Description: "Check there is no hyper-link in the name-field",
		Author:      "Steve Kemp <steve@steve.org.uk>",
		Test:        checkLinkName})

}

func checkLinkName(x Submission) (PluginResult, string) {
	if strings.HasPrefix(strings.ToLower(x.Name), "http") {
		return Spam, "Hyperlink detected in name-field"
	}

	return Undecided, ""
}
