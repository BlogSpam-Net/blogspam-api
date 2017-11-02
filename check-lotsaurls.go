//
// Check that there are not "too many" hyperlinks in a body.
//

package main

import (
	"regexp"
	"strings"
	"strconv"
	"index/suffixarray"
)

func init() {

	//
	// Add our plugin-method
	//
	x := Plugins{Name: "50-lotsaurls.js",
		Description: "Look for excessive numbers of HTTP links.",
		Author:      "Steve Kemp <steve@steve.org.uk>",
		Test:        checkHyperlinkCounts}
	plugins = append(plugins, x)

}

func checkHyperlinkCounts(x Submission) string {

	//
	// Map to store any options we find.
	//
	tmp := make(map[string]string)

	//
	// Default failure threshold.
	//
	tmp["max-links"] = "10"

	//
	// Do we have options?
	//
	if len(x.Options) > 0 {

		//
		// Split the string into an array, based on commas
		//
		options := strings.Split(x.Options, ",")

		//
		// Now look for key=val
		//
		for _, option := range options {
			re := regexp.MustCompile("^(.*)=([^=]+)$")
			match := re.FindStringSubmatch(option)

			if len(match) > 0 {
				tmp[match[1]] = match[2]
			}
		}
	}

	//
	// Now convert our (possibly updated) max value
	//
	max,err := strconv.Atoi(tmp["max-links"])
	if ( err != nil ) {
		return "Failed to parse max-links as a number"
	}
	if ( max <= 0 ) {
		return "Failed to parse max-links as a positive number"
	}


	//
	// Look for hyperlinks
	//
	r := regexp.MustCompile("https?://")

	//
	// Get the count
	//
	index := suffixarray.New([]byte(x.Comment))
	count := index.FindAllIndex(r, -1)

	if ( len(count) > max ) {
		return("Too many hyperlinks" )
	}
	//
	// All OK
	//
	return ""
}
