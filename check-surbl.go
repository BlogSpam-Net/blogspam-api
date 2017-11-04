//
//  Check for an IP that is in the surbl.org blacklist.
//

package main

import (
	"net"
	"net/url"
	"strings"
	"mvdan.cc/xurls"
)

func init() {
	//
	// Add our plugin-method.
	//
	var x = Plugins{Name: "60-surbl.js",
		Description: "Test links in messages against surbl.org",
		Author:      "Steve Kemp <steve@steve.org.uk>",
		Test:        checkSurblBlacklist}
	plugins = append(plugins, x)

}

//
// Lookup the hyperlinks in the Surbl.org blacklist.
//
func checkSurblBlacklist(x Submission) (PluginResult, string) {

	//
	// We'll store lookups to perform here.
	//
	lookups := make(map[string]int)

	//
	// Find the links in the body of our comment.
	//
	urlsRe := xurls.Relaxed()
	links := urlsRe.FindAllString(x.Comment, -1)

	//
	// If we got any we must process them.
	//
	// This means removing any `https?://` prefix
	// and any URL part.
	//
	for _,link := range( links ) {

		//
		// If we don't have a protocol-prefix, add it.
		//
		if ( ! strings.HasPrefix( link, "http://" ) &&
			! strings.HasPrefix( link, "https://" ) ) {
			link = "http://" + link
		}

		//
		// Now parse out the hostname of the link.
		//
		u, err := url.Parse(link)
		if err == nil {
			hostname := u.Hostname()

			//
			// The thing we lookup - stored in our map
			// to ensure we don't have duplicates
			//
			lookups[ hostname + ".multi.surbl.org" ] = 1
		}
	}

	//
	// Now we have a list of things to lookup.
	//
	// Let us do that, if any result in a result we know we've got spam.
	//
	for host, _ := range( lookups ) {

		reply, _ := net.LookupHost(host)
		if len(reply) != 0 {
			return Spam, "Posted link(s) listed in surbl.org"
		}
	}


	//
	// We got no listing, so we're OK.
	//
	return Undecided, ""
}
