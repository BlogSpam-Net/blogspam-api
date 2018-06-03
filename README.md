[![Travis CI](https://img.shields.io/travis/BlogSpam-Net/blogspam-api/master.svg?style=flat-square)](https://travis-ci.org/BlogSpam-Net/blogspam-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/BlogSpam-Net/blogspam-api)](https://goreportcard.com/report/github.com/BlogSpam-net/blogspam-api)
[![license](https://img.shields.io/github/license/BlogSpam-Net/blogspam-api.svg)](https://github.com/BlogSpam-net/blogspam-api/blob/master/LICENSE)


# Golang BlogSpam Server

The [BlogSpam.net](https://blogspam.net/) service presents an API which allows
incoming blog/forum comments to be tested for SPAM in real-time.

This repository contains an implementation of the API in golang, which allows you to run your own instance of the service, this superceeds the previous [implementation in node.js](https://github.com/BlogSpam-net/blogspam.js).


## Overview

The service presents a simple API over HTTP.  There are only a small number
of end-points:

* `POST /`
    * Test the incoming submission for SPAM.
* `POST /stats`
    * Retrieve the per-site SPAM/HAM statistics
* `GET /global-stats`
    * Retrieve global SPAM/HAM stats.
* `GET /plugins`
    * Retrieve the list of plugins.
* `POST /classify`
    * Retrain a comment.

These endpoints, and the parameters they require, are documented upon the website:

* [https://blogspam.net/api/2.0/](https://blogspam.net/api/2.0/)


## Plugin Implementation

Although we refer to them as "plugins" the individual tests which are applied to incoming submissions are all in-process and hardwired - there is nothing dynamic about them.

Each plugin has a name, and an order, and each is invoked in turn upon the incoming submission.  If any single plugin determines an incoming comment is SPAM then it is rejected, similarly any single plugin may decided a comment is definitely-HAM.  Otherwise processing continues until all plugins have been invoked.


## Installation

Providing you have a working `golang` environment you can install and
launch like so:

    $ go get github.com/BlogSpam-Net/blogspam-api
    $ blogspam-api -host 127.0.0.1 -port 9999 -redis localhost:6379

As hinted in the command-line arguments you'll want to install [redis](https://redis.io/) upon the local-host, but otherwise there is no configuration or setup required.


Steve
--
