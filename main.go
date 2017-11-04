//
// The BlogSpam.net API, written in golang.
//
// * We receive a JSON POST which we'll convert into a simple structure.
// * Then we run a bunch of "plugins" over the submission.
// * If any single plugin decides the comment is spam, we drop it.
// * Otherwise we're all OK.
//
// This is a bit ropy.
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
	"github.com/go-redis/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"
)

//
// The Submission structure is what we parse incoming JSON into.
//
// Each plugin which is implemented will operate solely on an instance
// of this structure.
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
// PluginResult is the return-code of each plugin-method.
//
// Each plugin will return a result which is "spam", "ham", "undecided",
// or error.  These are defined next.
//
type PluginResult int

//
// There are several possible plugin-results:
//
//   Spam:
//    Stop processing and inform the caller.
//   Ham:
//    Stop processing and inform the caller.
//   Undecided:
//    Continue running further plugins.
//   Error:
//    Internal error running a plugin.
//
const (
	Spam PluginResult = iota
	Ham
	Undecided
	Error
)

//
// PluginTest is the function which each plugin implements to check
// an incoming Submission instance for SPAM.
//
// This function is given a Submission structure and should return
// one of the enum-results noted above, as well as an optional detail
// field in the case of a SPAM-result.
//
type PluginTest func(Submission) (PluginResult, string)

//
// A Plugins object is present for each plugin which is implemented,
// and bundled with this repository.
//
// There is no provision for external plugins.
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

	//
	// Should SPAM-results be recorded in Redis?
	//
	// This is a dangerous setting, which is designed to cache
	// the results of expensive plugins.
	//
	RedisCache bool
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
// The global Redis client, if redis is enabled.
//
var redisHandle *redis.Client

//
// Should we be verbose?
//
var verbose bool

//
// Register a plugin - we use this method to ensure that the plugins
// are sorted by name, which means the lighter-weight plugins run
// first.
//
func registerPlugin(addition Plugins) {

	plugins = append(plugins, addition)

	sort.Slice(plugins[:], func(i, j int) bool {
		return plugins[i].Name < plugins[j].Name
	})
}

//
// ClassifyHandler is a HTTP-Handler which should re-train the given input.
//
// However it is not implemented.
//
func ClassifyHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "OK")
}

//
// StatsHandler is a HTTP-handler which should return the per-site
// statistics to the caller for the given site.
//
func StatsHandler(res http.ResponseWriter, req *http.Request) {
	var (
		status int
		err    error
	)
	defer func() {
		if nil != err {
			http.Error(res, err.Error(), status)
			// Don't spam stdout when running test-cases.
			if flag.Lookup("test.v") == nil {
				fmt.Printf("WARNING - Error returned from /stats handler - %s\n", err.Error())
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
	// Create a map for returning our results to the caller.
	//
	// We default to having zero for both counts.  This ensures
	// we populate the return-value(s) in the event of an error,
	// or if redis is disabled
	//
	ret := make(map[string]string)
	ret["spam"] = "0"
	ret["ok"] = "0"

	//
	// If we have a site then we're good
	//
	site := input.Site

	//
	// Get the spam-count, and assuming no error then we
	// update our map.
	//
	if redisHandle != nil {
		spamCount, err := redisHandle.Get(fmt.Sprintf("site-%s-spam", site)).Result()
		if err != nil {
			ret["error"] = err.Error()
		} else {
			ret["spam"] = spamCount
		}
	}

	//
	// Get the ham-count, and assuming no error then we
	// update our map.
	//
	if redisHandle != nil {
		hamCount, err := redisHandle.Get(fmt.Sprintf("site-%s-ok", site)).Result()
		if err != nil {
			ret["error"] = err.Error()
		} else {
			ret["ok"] = hamCount
		}
	}

	//
	// Convert this temporary hash to a JSON object we can return
	//
	jsonString, err := json.Marshal(ret)
	if err != nil {
		status = http.StatusInternalServerError
		return
	}

	//
	// Send it.
	//
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(res, "%s", jsonString)
}

//
// SendSpamResult informs the caller of a SPAM result.
//
// Bump our global and per-site count, if redis is available.
//
func SendSpamResult(res http.ResponseWriter, input Submission, plugin Plugins, detail string) {

	if redisHandle != nil {
		//
		// Bump the global count of SPAM.
		//
		redisHandle.Incr("global-spam")

		//
		// Bump the per-site count of SPAM.
		//
		redisHandle.Incr(fmt.Sprintf("site-%s-spam", input.Site))
	}

	//
	// This plugin-test resulted in a spam result, and we'll
	// return that to the caller as JSON.
	//
	// Create a map to hold the details for now.
	//
	ret := make(map[string]string)
	ret["result"] = "SPAM"
	ret["blocker"] = plugin.Name
	ret["reason"] = detail
	ret["version"] = "2.0"

	//
	// Covert the temporary hash to a JSON-object.
	//
	jsonString, err := json.Marshal(ret)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	//
	// Send to the caller.
	//
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(res, "%s", jsonString)

}

//
// SendOKResult tells the caller their submission was not Spam.
//
// Bump our global and per-site count, if redis is available.
//
func SendOKResult(res http.ResponseWriter, input Submission) {

	if redisHandle != nil {
		//
		// Bump the global Ham-count
		//
		redisHandle.Incr("global-ok")

		//
		// Bump the per-site Ham-count
		//
		redisHandle.Incr(fmt.Sprintf("site-%s-ok", input.Site))
	}

	//
	// Send the result to the caller.
	//
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(res, "{\"result\":\"OK\", \"version\":\"3.0\"}")
}

//
// SpamTestHandler is the meant of our server, it reads incoming
// JSON submissions and invokes plugins to determine if the submission
// represented a SPAM comment.
//
// Parse the incoming JSON-structure, and if there are no errors
// in doing so then test the comment with all known plugins.
//
// Once complete send the appropriate result to the caller.
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
				fmt.Printf("WARNING - Error returned from / handler - %s\n", err.Error())
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
	// Dump the incoming request to STDOUT if running verbosely.
	//
	if verbose {

		//
		// Get all the fields of the structure, via reflection
		//
		s := reflect.ValueOf(&input).Elem()
		typeOfT := s.Type()

		//
		// Iterate over the fields.
		//
		for i := 0; i < s.NumField(); i++ {

			// The specific field
			f := s.Field(i)

			// The name/value of the field
			fieldName := typeOfT.Field(i).Name
			fieldVal := fmt.Sprintf("%s", f.Interface())

			// Print non-empty fields
			if len(fieldVal) > 0 {
				fmt.Printf("\t%s : %s\n", fieldName, fieldVal)
			}
		}
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
	// Now we invoke each known-plugin, unless we're to exclude
	// any specific one.
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
			// Look for this plugin being excluded.
			//
			if strings.Contains(name, ex) || name == ex {
				skip = true
			}
		}

		if skip {
			continue
		}

		//
		// Call the plugin method to run the test.
		//
		result, detail := obj.Test(input)

		if verbose {
			fmt.Printf("Plugin %s returned: %d,%s\n",
				obj.Name, result, detail)
		}

		if result == Spam {
			//
			// If the plugin-method decided this submission was
			// SPAM then we immediately return that result to the
			// caller of our service.
			//
			SendSpamResult(res, input, obj, detail)

			//
			// If we should cache in redis, and redis
			// is enabled, do so
			//
			if (obj.RedisCache == true) && (redisHandle != nil) {
				key := fmt.Sprintf("blacklist-%s", input.IP)
				period := time.Hour * 48
				err := redisHandle.Set(key, detail, period).Err()
				if err != nil {
					fmt.Printf("WARNING redis-error blacklisting IP %s - %s\n", input.IP, err.Error())
				}
			}

			return
		}
		if result == Ham {

			//
			// The result is definitely OK - tell the caller.
			//
			SendOKResult(res, input)
			return

		}
		if result == Undecided {

			// Nop
		}
		if result == Error {

			// Nop
			fmt.Printf("Error running plugin: %s\n\t%s\n",
				obj.Name, detail)
		}
	}

	//
	// If we reached this point no plugin decided this was SPAM,
	// so we default to saying it was Ham.
	//
	SendOKResult(res, input)
}

//
// PluginListHandler is a HTTP-handler to export our list of known-plugins.
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
func serve(host string, port int) {

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
	bind := fmt.Sprintf("%s:%d", host, port)
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

	//
	// The command-line flags we support
	//

	//
	// Host/Port for binding upon
	//
	host := flag.String("host", "127.0.0.1", "The IP to bind upon")
	port := flag.Int("port", 9999, "The port number to listen upon")
	verb := flag.Bool("verbose", false, "Should we be verbose")

	//
	// Optional redis-server address
	//
	rserver := flag.String("redis", "",
		"The host:port of the optional redis-server to use.")

	//
	// Parse the flags
	//
	flag.Parse()

	//
	// Set the global verbose flag.
	//
	verbose = (*verb == true)

	//
	// If redis host/port was specified then open the connection now.
	//
	if len(*rserver) > 0 {
		fmt.Printf("Using redis-server %s\n", *rserver)
		redisHandle = redis.NewClient(&redis.Options{
			Addr:     *rserver,
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	} else {
		redisHandle = nil
	}

	//
	// And finally start our server
	//
	serve(*host, *port)
}
