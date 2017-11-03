//
//  Check for an IP that is in the dronebl.org blacklist.
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
	var x = Plugins{Name: "60-drone.js",
		Description: "Test IP of the comment-submitter against dronebl.org",
		Author:      "Steve Kemp <steve@steve.org.uk>",
		Test:        checkDroneBlacklist}
	plugins = append(plugins, x)

}

//
// Lookup the IP address of the submitter in the dronebl.org blacklist.
//
func checkDroneBlacklist(x Submission) (PluginResult, string) {

	//
	// See if we have an IPv4 address.
	//
	regex := regexp.MustCompile("^([0-9]+).([0-9]+).([0-9]+).([0-9]+)$")
	match := regex.FindStringSubmatch(x.IP)

	//
	// If that failed then we know we have an IPv6-address, or a missing
	// address, so we terminate
	//
	if len(match) <= 0 {
		return Undecided, ""
	}

	//
	// Convert the IPv4 address into an array of items
	//
	octets := strings.Split(x.IP, ".")

	//
	// Now calculate what we're going to lookup
	//
	lookup := fmt.Sprintf("%s.%s.%s.%s.dnsbl.dronebl.org",
		octets[3], octets[2], octets[1], octets[0])

	//
	// And look it up.
	//
	reply, _ := net.LookupHost(lookup)

	//
	// No reply?  Not spam
	//
	if len(reply) == 0 {
		return Undecided, ""
	}

	//
	// We got a listing, so we're SPAM.
	//
	return Spam, fmt.Sprintf("%s is listed in dronebl.org", x.IP)
}
