//
// Check that there are not too many different attempts at adding a URL
//

package main

import (
	"index/suffixarray"
	"regexp"
)

func init() {

	//
	// Add our plugin-method
	//
	x := Plugins{Name: "50-multilinks.js",
		Description: "Look for different linking strategies.",
		Author:      "Steve Kemp <steve@steve.org.uk>",
		Test:        checkLinkingTypes}
	plugins = append(plugins, x)

}

//
// Look for multiple linking strategies
//
func checkLinkingTypes(x Submission) (PluginResult, string) {

	//
	// The things we're looking for
	//
	patterns := []string{"<a href=\"https?:",
		"\\[?url=https?:",
		"\\[?link=https?:",
		"[ \t]https?:/"}

	//
	// Count of types we've found thus far
	//
	count := 0

	//
	// For each pattern.
	//
	for _, p := range patterns {

		// Look for matches
		r := regexp.MustCompile(p)

		index := suffixarray.New([]byte(x.Comment))
		matches := index.FindAllIndex(r, -1)

		if len(matches) > 0 {
			count += 1
		}
	}

	if count >= 3 {
		return Spam, "Multiple linking strategies"
	}

	//
	// All OK
	//
	return Undecided, ""
}
