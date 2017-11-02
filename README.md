# BlogSpam v3

I'm in the process of cutting down on the number of machines I host,
and one of the casualities is the BlogSpam.net service.

Rather than killing it outright I need to make it more efficient, such
that it can run upon an existing server I have.  It has been suggested
that I can reduce functionality, and port to golang to make that work.

I don't want to spend huge amounts of time on the process, as [the
existing code](https://github.com/skx/blogspam.js) is open-source and
available for users who wish to continue.  But I'm certainly willing to
spend two days trying to move it.

## Overview

The current codebase, version 2, is pretty simple:

* Read an incoming HTTP POST.
* The POST is a JSON object with a small number of fields
    * Name
    * Email
    * Body
    * IP
    * etc
* Once the submission has been decoded from JSON invoke a series of plugins against it.
    * If any single plugin decides the submission was spam report that
    * Otherwise return a "good" result.

I'm not going to change this basic setup, so really I need a golang HTTP-server that listens for POST-submissions, decodes the JSON, and invokes a series of "plguins" on it.  The plugins won't be _real_ plugins, but self-contained test-code.

As a proof of concept I "just" need to handle the POST-receiving, the decoding, and reimplement a couple of the plugins.   I know from my logs which ones are the most useful, so that's an easy thing to do.


## TODO

* Port more plugins
