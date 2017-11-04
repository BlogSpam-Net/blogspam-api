[![Travis CI](https://img.shields.io/travis/skx/golang-blogspam/master.svg?style=flat-square)](https://travis-ci.org/skx/golang-blogspam)
[![Go Report Card](https://goreportcard.com/badge/github.com/skx/golang-blogspam)](https://goreportcard.com/report/github.com/skx/golang-blogspam)
[![license](https://img.shields.io/github/license/skx/golang-blogspam.svg)](https://github.com/skx/golang-blogspam/blob/master/LICENSE)


# Golang BlogSpam Server

I'm in the process of cutting down on the number of machines I host,
and one of the casualities is the [BlogSpam.net](https://blogspam.net/) service.

Rather than killing it outright I need to make it more efficient, such
that it can run upon an existing server I have, rather than requiring a dedicated machine of its own to run upon.  It has been suggested that I can reduce functionality, and port it to golang to make that happen.

> **NOTE** I did update the site to say the service will be retired, but I suspect nobody using the service will have noticed.  Oops.

I don't want to spend huge amounts of time on the process, as [the
existing code](https://github.com/skx/blogspam.js) is open-source and
available for users who wish to continue.  But I'm certainly willing to
spend a day or two trying to move it.

## Overview

The current codebase, version 2, is pretty simple:

* Read an incoming HTTP POST.
* The POST is a JSON object with a small number of fields:
    * Name
    * Email
    * Body
    * IP
    * etc
* Once the submission has been decoded from JSON invoke a series of plugins against it.
    * If any single plugin decides the submission was spam report that.
    * Otherwise return a "good" result.

I'm not going to change this basic setup, so really I need a golang HTTP-server that listens for POST-submissions, decodes the JSON, and invokes a series of "plguins" on it.  The plugins won't be _real_ plugins, but self-contained test-code.


## Proof of Concept

This repository implements the core of the BlogSPAM API, and can be deployed
easily.

Providing you have a working `golang` environment you can install and
launch like so:

    $ go get github.com/skx/golang-blogspam.git
    $ golang-blogspam -host 127.0.0.1 -port 9999 -redis localhost:6379

Once launched the original test-cases from the blogspam.js repository
should succeed when fired against it.

If a redis-server is specified it will be used, which means that previously
blacklist-IPs will be detected and rejected at low-cost.  In the current
implementation spam-submissions will not result in an IP being globally
blacklisted, that might come in the future though.

