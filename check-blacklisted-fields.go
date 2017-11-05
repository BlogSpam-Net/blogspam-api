//
//  Check for blacklisted fields.
//
//  We load a file for each possible field in our input, which means it
// is simple to blacklist regular expressions against particular fields.
//
//  For example we have a plugin which looks for hyperlinks in the name
// field.  We could replicate that here via:
//
//    echo ^https?:// >> ./blacklist.d/name
//
//  Simple.
//
//

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
)

//
// We store blacklisted field-data here
//
var blacklisted map[string][]string

//
// Store the data from the specified file into our blacklisted-map
//
func readData(path string, name string) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		blacklisted[name] = append(blacklisted[name], scanner.Text())
	}

	return scanner.Err()
}

//
// Process the given configuration-directory
//
func processDirectory(dir string) {

	files, err := ioutil.ReadDir(dir)
	if err == nil {
		for _, f := range files {
			readData(fmt.Sprintf("%s/%s", dir, f.Name()), f.Name())
		}
	}

}

//
// Register ourselves as a plugin, after setting up our config-files.
//
func init() {

	//
	// Create a map to hold our per-field lists
	//
	blacklisted = make(map[string][]string)

	//
	// Look for a set of field-based config-files.
	//
	processDirectory("./blacklist.d/")
	processDirectory("/etc/blogspam/blacklist.d/")

	registerPlugin(BlogspamPlugin{Name: "05-blacklisted-fields.js",
		Description: "Look for blacklisted patterns in fields",
		Author:      "Steve Kemp <steve@steve.org.uk>",
		Test:        checkBlacklistedFields})
}

//
// Test the incoming submission against our blacklist.
//
func checkBlacklistedFields(x Submission) (PluginResult, string) {

	//
	// We've got a list of fields, and a map of blacklists.
	//
	// Get all the fields of the structure, via reflection
	//
	s := reflect.ValueOf(&x).Elem()
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

		// Now we have an array of blacklisted items
		items := blacklisted[strings.ToLower(fieldName)]

		// We'll iterate over them.
		for _, val := range items {

			//
			// Each item is a regular expression.
			//
			// We make them case-insensitive with the "(?i)" prefix
			//
			re := regexp.MustCompile("(?i)" + val)
			match := re.FindStringSubmatch(fieldVal)

			if len(match) > 0 {
				return Spam, fmt.Sprintf("Blacklisted value in %s-field", fieldName)
			}

		}
	}

	return Undecided, ""
}
