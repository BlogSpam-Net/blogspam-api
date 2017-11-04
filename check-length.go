//
//  Check for min/max size of name/subject
//

package main

func init() {
	//
	// Add our plugin-method.
	//
	var x = Plugins{Name: "42-length.js",
		Description: "Look at the size of subject and name",
		Author:      "Steve Kemp <steve@steve.org.uk>",
		Test:        validateLength}
	registerPlugin(x)

}

//
// Hacky test to drop comments with too-long name/subjects.
// test them.
//
func validateLength(x Submission) (PluginResult, string) {

	if len(x.Name) > 140 {
		return Spam, "The submitted 'name' is too long."
	}

	if len(x.Subject) > 140 {
		return Spam, "The submitted 'subject' is too long."
	}

	//
	// All OK
	//
	return Undecided, ""
}
