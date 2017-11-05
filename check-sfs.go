//
//  Check for an IP that is in the stopforumspam.com blacklist.
//

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

//
// Register ourself as a blogspam-plugin.
//
func init() {
	registerPlugin(BlogspamPlugin{Name: "80-sfs.js",
		Description: "Look for blacklisted IPs via stopforumspam.com",
		Author:      "Steve Kemp <steve@steve.org.uk>",
		Test:        checkSFSBlacklist,
		RedisCache:  true})
}

//
// Lookup the IP address of the submitter in the stopforumspam.com blacklist.
//
func checkSFSBlacklist(x Submission) (PluginResult, string) {

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
	// Build a client with sane timeout
	//
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	//
	// The URL we'll fetch
	//
	url := fmt.Sprintf("http://www.stopforumspam.com/api?ip=%s", x.IP)

	//
	// Make the request
	//
	response, err := netClient.Get(url)

	//
	// Handle error
	//
	if err != nil {
		fmt.Printf("WARNING: HTTP-Error reading from %s - %s", url, err)
		return Error, err.Error()
	}

	//
	// Ensure we close the body
	//
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("WARNING: HTTP-Error reading body from %s - %s", url, err)
		return Error, err.Error()
	}

	//
	// Does it appear?
	//
	if strings.Contains(string(contents), "<appears>yes</appears>") {
		return Spam, "Listed in StopForumSpam.com"
	}

	//
	// Not listed
	//
	return Undecided, ""
}
