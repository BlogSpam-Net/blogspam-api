//
//  Check for local IP blacklist.
//

package main

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

func init() {
	//
	// Add our plugin-method.
	//
	var x = Plugins{Name: "20-ip.js",
		Description: "Look for blacklisted IP addresses",
		Author:      "Steve Kemp <steve@steve.org.uk>",
		Test:        checkBlacklist}
	plugins = append(plugins, x)

}

//
// Test that the submitter isn't blacklisted, by IP.
//
func checkBlacklist(x Submission) (PluginResult, string) {

	//
	// Map to store any IPs we're to blacklist
	//
	tmp := make(map[string]int)

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
			re := regexp.MustCompile("blacklist=([^=]+)$")
			match := re.FindStringSubmatch(option)

			if len(match) > 0 {
				tmp[match[1]] = 1
			}
		}
	}

	//
	// The source IP we're going to test against the blacklisted entries.
	//
	source := net.ParseIP(x.IP)

	//
	// If we have some blacklisted IPs..
	//
	for ip := range tmp {

		// Only parse the CIDR if it looks like one.
		if strings.Contains( ip, "/" ) {
			// Parse the network
			_, subnet, err := net.ParseCIDR(ip)
			if err != nil {
				return Error, fmt.Sprintf("Failed to parse CIDR %s", ip)
			}

			// Is it in there?
			if subnet.Contains(source) {
				return Spam, "IP blacklisted"
			}
		} else {

			// Is it a literal match?
			if x.IP == ip  {
				return Spam, "IP blacklisted"
			}
		}
	}

	//
	// If Redis is not available we're done
	//
	if redisHandle == nil {
		return Undecided, ""
	}

	//
	// Since we have redis-enabled we'll now look for the remote IP too.
	//
	// The key is named `blacklist-$IP`
	//
	key := fmt.Sprintf("blacklist-%s", x.IP)

	//
	// Run the lookup
	//
	result, err := redisHandle.Get(key).Result()
	if err != nil {
		return Error, err.Error()
	}

	//
	// If there was a result then it is spam
	//
	if len(result) > 0 {
		return Spam, result
	}

	//
	// Not blocked by options, or previous attempts
	//
	return Undecided, ""
}
