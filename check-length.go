//
//  Check for min/max size of name/subject
//

package main

//
// Register ourself as a blogspam-plugin.
//
func init() {
	registerPlugin(BlogspamPlugin{Name: "42-length.js",
		Description: "Look at the size of subject and name",
		Author:      "Steve Kemp <steve@steve.org.uk>",
		Test:        validateLength})

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
