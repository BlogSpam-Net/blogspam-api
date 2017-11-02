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
	plugins = append(plugins, x)

}

//
// Hacky test to drop comments with too-long name/subjects.
// test them.
//
func validateLength(x Submission) string {

	if len(x.Name) > 150 {
		return ("The submitted 'name' is too long.")
	}

	if len(x.Subject) > 150 {
		return ("The submitted 'subject' is too long.")
	}

	//
	// All OK
	//
	return ""
}
