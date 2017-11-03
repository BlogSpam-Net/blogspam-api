[![Travis CI](https://img.shields.io/travis/skx/golang-blogspam-server/master.svg?style=flat-square)](https://travis-ci.org/skx/golang-blogspam-server)
[![Go Report Card](https://goreportcard.com/badge/github.com/skx/golang-blogspam-server)](https://goreportcard.com/report/github.com/skx/puppet-summary)
[![license](https://img.shields.io/github/license/skx/golang-blogspam-server.svg)](https://github.com/skx/golang-blogspam-server/blob/master/LICENSE)


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
    * If any single plugin decides the submission was spam report that
    * Otherwise return a "good" result.

I'm not going to change this basic setup, so really I need a golang HTTP-server that listens for POST-submissions, decodes the JSON, and invokes a series of "plguins" on it.  The plugins won't be _real_ plugins, but self-contained test-code.


## Proof of Concept

This repository contains a proof of concept - it launches a service, reads the JSON-bodies, and invokes plugins against those submissions.

I've ported several existing plugins and confirmed they pass the tests in the original repositories test-suit, and added 100% golang test-coverage.

I've also hardwired a redis-connection to localhost, which will store the state
of spam/ham counts for each site - as well as globally.

## Missing Features

The existing server uses redis to maintain state:

* IPs that send "bad" comments will often be blacklisted for a period of hours.

For the moment I've ignored both of those features.  The stats are useful to the site-owners, and myself, but I think we can live without them.

## TODO

* Port more plugins.
* Breakout the code in `main.go` so we can add tests
* Make redis optional
* Lookup blacklisted IPs in redis.
* The return values for "plugins" should be finer-grained:
   * OK - This is valid comment and cease now
   * Spam - this is spam, stop now
   * Continue - Continue testing.


## Benchmarks

Fake benchmarks:

    $ cat ~/site.json
    {"site": "http://example.com" }


Localhost + golang:

     $ ab -p ~/site.json -T application/json -c 10 -n 2000 http://localhost:9999/stats
     Time taken for tests:   0.210 seconds
     Complete requests:      2000
     Failed requests:        0
     Total transferred:      262000 bytes
     Total body sent:        344000
     HTML transferred:       46000 bytes
     Requests per second:    9505.16 [#/sec] (mean)
     Time per request:       1.052 [ms] (mean)
     Time per request:       0.105 [ms] (mean, across all concurrent requests)

Remote server + node.js - **NOTE** this is doing a quarter the the number of tests:

     $ ab -p ~/site.json -T application/json -c 10 -n 500 http://test.blogspam.net:9999/stats
     Time taken for tests:   85.771 seconds
     Complete requests:      500
     Failed requests:        0
     Total transferred:      65000 bytes
     Total body sent:        90000
     HTML transferred:       11500 bytes
     Requests per second:    5.83 [#/sec] (mean)
     Time per request:       1715.421 [ms] (mean)
     Time per request:       171.542 [ms] (mean, across all concurrent requests)

I think that says it all...
