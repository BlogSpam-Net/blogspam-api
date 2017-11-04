//
// Check that mandatory fields are present.
//

package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func init() {
	//
	// Add our plugin-method.
	//
	var x = Plugins{Name: "30-mandatory.js",
		Description: "Look for any mandatory fields which might be missing.",
		Author:      "Steve Kemp <steve@steve.org.uk>",
		Test:        validateMandatory}
	registerPlugin(x)

}

//
// Test that mandatory fields are present.
//
func validateMandatory(x Submission) (PluginResult, string) {

	//
	// The mandatory fields we're going to insist upon by default
	//
	tmp := make(map[string]int)
	tmp["site"] = 1
	tmp["comment"] = 1
	tmp["ip"] = 1

	//
	// Do we have options?
	//
	if len(x.Options) > 0 {

		//
		// Split them into an array, based on commas
		//
		options := strings.Split(x.Options, ",")

		//
		// Now look for any additional mandatory fields
		//
		for _, option := range options {
			re := regexp.MustCompile("mandatory=([^=]+)$")
			match := re.FindStringSubmatch(option)

			if len(match) > 0 {
				tmp[match[1]] = 1
			}
		}
	}

	//
	// Now we can do the test for missing fields.
	//
	// There __must__ be a better way of doing this, by looking
	// at the subject field with reflection.
	//
	for field := range tmp {

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

			// Is this the field we're looking for?
			if strings.EqualFold(field, fieldName) {

				// Then raise an error if it is empty
				if len(fieldVal) < 1 {
					return Spam, fmt.Sprintf("Field %s is missing", fieldName)
				}
			}
		}
	}

	return Undecided, ""
}
