//
// Check that the "Name" field of our incoming submission doesn't contain
// a hyperlink
//

package main

import "strings"

func init() {

	//
	// Add our plugin-method
	//
	x := Plugins{Name: "35-name.js",
		Description: "Check there is no hyper-link in the name-field",
		Author:      "Steve Kemp <steve@steve.org.uk>",
		Test:        checkLinkName}
	plugins = append(plugins, x)

}

func checkLinkName(x Submission) (PluginResult, string) {
	if strings.HasPrefix(strings.ToLower(x.Name), "http") {
		return Spam, "Hyperlink detected in name-field"
	}

	return Undecided, ""
}
