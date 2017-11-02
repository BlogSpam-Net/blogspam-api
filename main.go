//
// This is a simple proof-of-concept port of theBlogSpam.net API to golang.
//
// * We receive a JSON POST which we'll convert into a simple structure.
//
// * Then we run a bunch of "plugins" over the submission.
//
// * If any single plugin decides the comment is spam, we drop it.
//
// * Otherwise we're all OK.
//
// Steve
// --
//

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"regexp"
	"strings"
)

//
// The incoming JSON struct
//
type Submission struct {
	//
	// The user-agent that submitted the comment - optional
	//
	Agent string

	//
	// The actual comment - mandatory
	//
	Comment string

	//
	// The email of the comment submitter - optional
	//
	Email string

	//
	// The IP that submitted the comment - mandatory
	//
	IP string

	//
	// The link the comment-submitter supplied - optional
	//
	Link string

	//
	// The author-name  of the comment - optional
	//
	Name string

	//
	// Any options - optional
	//
	Options string

	//
	// The site this comment was for - mandatory
	//
	Site string

	//
	// The subject the author supplied - optional
	//
	Subject string

	//
	// The version of your plugin, if any - optional
	//
	Version string
}

//
// The spam-testing method
//
// This function is given a Submission structure and should return
// nil if the submission looks OK, otherwise an error-string
//
type PluginTest func(Submission) string

//
// A structure to describe each known-plugin.
//
type Plugins struct {
	//
	// The author of the plugin.
	//
	Author string

	//
	// The name of the plugin.
	//
	Name string

	//
	// A description of the plugin.
	//
	Description string

	//
	// The function to invoke to use the plugin.
	//
	Test PluginTest
}

//
// The global list of plugins we have loaded.
//
// Since we're using golang everything is static, but we could have
// chosen to use the new plugin API.  For the moment we'll avoid that
// to simplify the compilation, and also because nobody ever contributed
// a plugin of any kind.
//
//
var plugins []Plugins

//
// HTTP-Handler: Re-train input.  [NOP]
//
func ClassifyHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "OK")
}

//
// HTTP-Handler: Dump statistics.  [NOP]
//
func StatsHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "{\"spam\":\"0\",\"ok\":\"0\",\"nop\":\"true\"}")
}

//
// Our spam-test handler
//
// Parse the incoming JSON-structure, and if there are no errors
// in doing so then test the comment with all known plugins.
//
func SpamTestHandler(res http.ResponseWriter, req *http.Request) {
	var (
		status int
		err    error
	)
	defer func() {
		if nil != err {
			http.Error(res, err.Error(), status)
			// Don't spam stdout when running test-cases.
			if flag.Lookup("test.v") == nil {
				fmt.Printf("Error: %s\n", err.Error())
			}
		}
	}()

	//
	// Ensure this was a POST-request
	//
	if req.Method != "POST" {
		err = errors.New("Must be called via HTTP-POST")
		status = http.StatusInternalServerError
		return
	}

	//
	// Decode the submitted JSON body
	//
	decoder := json.NewDecoder(req.Body)

	//
	// This is what we'll decode
	//
	var input Submission
	err = decoder.Decode(&input)

	//
	// If decoding the JSON failed then we'll abort
	//
	if err != nil {
		status = http.StatusInternalServerError
		return
	}

	//
	// If we decoded then pretty-print it - unless running `make test`.
	//
	if flag.Lookup("test.v") == nil {
		fmt.Printf("\t%+v\n", input)
	}

	//
	// We might have options which will disable upcoming plugins.
	//
	// If so we'll keep track of the plugins that are excluded here.
	//
	var exclude []string

	//
	// Do we have options?
	//
	if len(input.Options) > 0 {

		//
		// Split the string into an array, based on commas
		//
		options := strings.Split(input.Options, ",")

		//
		// Now look for exclusions
		//
		for _, option := range options {
			re := regexp.MustCompile("^exclude=(.*)$")
			match := re.FindStringSubmatch(option)

			if len(match) > 0 {
				exclude = append(exclude, match[1])
			}
		}
	}

	//
	// Now we invoke each known-plugin
	//
	for _, obj := range plugins {

		//
		// The name of this plugin, and whether we should skip it
		//
		name := obj.Name
		var skip = false

		//
		// Look for exclusion(s)
		//
		for _, ex := range exclude {

			//
			// TODO: Regexp-Check
			//
			if strings.Contains(name, ex) || name == ex {
				if flag.Lookup("test.v") == nil {
					fmt.Printf("\tPlugin skipped: %s\n", name)
				}
				skip = true
			}
		}

		if skip {
			continue
		}

		//
		// Call the test
		//
		result := obj.Test(input)

		//
		// If the plugin-method decided this submission was
		// SPAM then we immediately reutrn that result to the
		// caller of our service.
		//
		if len(result) > 0 {

			//
			// This plugin-test resulted in a spam result
			//
			ret := make(map[string]string)
			ret["result"] = "SPAM"
			ret["blocker"] = obj.Name
			ret["reason"] = result
			ret["version"] = "2.0"

			//
			// Covert this temporary hash to a JSON
			// hash we can send to the caller.
			//
			jsonString, err := json.Marshal(ret)
			if err != nil {
				status = http.StatusInternalServerError
				return
			} else {
				res.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(res, "%s", jsonString)
			}

			if flag.Lookup("test.v") == nil {
				fmt.Printf("\nXXXX SPAM - %s: %s\n", obj.Name, result)
			}

			return
		}
	}

	//
	// Here we've invoked each of our plugins, and we didn't get
	// a SPAM result, so we assume we're good.
	//
	fmt.Fprintf(res, "{\"result\":\"OK\", \"version\":\"3.0\"}")
}

//
// Our plugin-list handler
//
func PluginListHandler(res http.ResponseWriter, req *http.Request) {
	var (
		status int
		err    error
	)
	defer func() {
		if nil != err {
			http.Error(res, err.Error(), status)
			fmt.Printf("Error: %s\n", err.Error())
		}
	}()

	//
	// Make a map.
	//
	m := make(map[string](map[string](string)))

	//
	// Populate it, from our known-plugins.
	//
	for _, obj := range plugins {
		m[obj.Name] = make(map[string]string)

		m[obj.Name]["author"] = obj.Author
		m[obj.Name]["description"] = obj.Description
	}

	//
	// Convert to JSON.
	//
	jsonString, err := json.Marshal(m)

	//
	// If that failed abort.
	//
	if err != nil {
		status = http.StatusInternalServerError
		return
	}

	//
	// Otherwise send back to the caller.
	//
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(res, "%s", jsonString)
}

//
// Launch our HTTP server
//
func serve(port int) {

	//
	// Create a new router and our route-mappings.
	//
	router := mux.NewRouter()

	//
	// API end-points.
	//
	//  1. Spam-Test
	//
	router.HandleFunc("/", SpamTestHandler).Methods("POST")
	//
	//  2. Plugin-List
	//
	router.HandleFunc("/plugins", PluginListHandler).Methods("GET")
	router.HandleFunc("/plugins/", PluginListHandler).Methods("GET")
	//
	//  3.  Stats
	//
	router.HandleFunc("/stats", StatsHandler).Methods("POST")
	router.HandleFunc("/stats/", StatsHandler).Methods("POST")
	//
	//  4.  Classify/Train comments
	//
	router.HandleFunc("/classify", ClassifyHandler).Methods("POST")
	router.HandleFunc("/classify/", ClassifyHandler).Methods("POST")
	//

	//
	// Bind the router.
	//
	http.Handle("/", router)

	//
	// Show where we'll bind
	//
	bind := fmt.Sprintf("%s:%d", "0.0.0.0", port)
	fmt.Printf("Launching the server on http://%s\n", bind)

	//
	// Wire up logging.
	//
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	//
	// Launch the server.
	//
	err := http.ListenAndServe(bind, loggedRouter)
	if err != nil {
		fmt.Printf("\nError: %s\n", err.Error())
	}
}

func main() {

	// //
	// // Create a sample submission
	// //
	// v := Submission{Name: "http://example.com",
	// 	Email: "non@non.org"}

	// //
	// // Call each known-plugin
	// //
	// for _, obj := range plugins {
	// 	res := obj.Test(v)
	// 	fmt.Printf("Plugin: %s -> Result: %s\n", obj.Name, res)
	// 	if len(res) > 0 {
	// 		fmt.Printf("\tSPAM!\n")
	// 	} else {
	// 		fmt.Printf("\tOK\n")
	// 	}
	// }

	serve(9999)
}
