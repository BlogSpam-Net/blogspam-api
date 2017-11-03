//
//  Check for min/max size of comment.
//

package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	//
	// Add our plugin-method.
	//
	var x = Plugins{Name: "40-size.js",
		Description: "Look at the size of the body",
		Author:      "Steve Kemp <steve@steve.org.uk>",
		Test:        validateSize}
	plugins = append(plugins, x)

}

//
// If there are options which specify the min/max-size of the body, then
// test them.
//
func validateSize(x Submission) (PluginResult, string) {

	//
	// Map to store any options we find.
	//
	tmp := make(map[string]string)

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
	// Do we have a min-size?
	//
	if len(tmp["min-size"]) > 0 {
		i, err := strconv.Atoi(tmp["min-size"])
		if err != nil {
			return Error, "Failed to parse max-size as a number"

		}
		if i <= 0 {
			return Error, "Failed to parse max-size as a positive number"
		}

		if len(x.Comment) < i {
			return Spam, fmt.Sprintf("Comment size is %d which is less than the minimum size %s", len(x.Comment), tmp["min-size"])
		}
	}

	//
	// Do we have a max-size?
	//
	if len(tmp["max-size"]) > 0 {
		i, err := strconv.Atoi(tmp["max-size"])
		if err != nil {
			return Error, "Failed to parse max-size as a number"
		}
		if i <= 0 {
			return Error, "Failed to parse max-size as a positive number"
		}

		if len(x.Comment) > i {
			return Spam, fmt.Sprintf("Comment size is %d which is more than the maximum size %s", len(x.Comment), tmp["min-size"])
		}
	}

	//
	// All OK
	//
	return Undecided, ""
}
