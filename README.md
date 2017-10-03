# boiling-api
*spicy hot stuff so spicy oh wow*

[![Build Status](https://api.travis-ci.org/boilingrip/boiling-api.svg?branch=master)](https://travis-ci.org/boilingrip/boiling-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/boilingrip/boiling-api)](https://goreportcard.com/report/github.com/boilingrip/boiling-api)
[![GoDoc](https://godoc.org/github.com/boilingrip/boiling-api?status.svg)](https://godoc.org/github.com/boilingrip/boiling-api)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

This is the home of the *boiling* project.
Boiling aims to be an open-source software for private trackers.

## Goals of the project

- Develop a music-themed private tracker API as __open-source__ software
- Have strict abstraction between the API (this project) and any frontend (not this project)
- Use the [chihaya](https://github.com/chihaya/chihaya) tracker backend
- Produce well-tested, reliable (and hopefully fast) code
- (Maybe) Have some fancy plugin architecture so people can extend it (honestly lowest priority for now)

## Motivation

Although I always wanted to run a private tracker, I figured it's a lot of risk for a small reward.
Even if you developed the coolest tracker ever, you couldn't put it on your CV or maybe even tell your friends about it.
Running the infrastructure, managing the community, and doing all that while staying out of sight seems like a time-consuming task.

Nonetheless, I wanted to write a tracker software.
After experimenting with a bunch of frameworks, languages, frontend-technologies and whatnot I got fed up with frontend development.
So I had the idea of just writing a "backend" for the frontend - an API.

I noticed that most development of tracker software is done behind closed doors, in secret.
This might make sense for some features, but is detrimental to the quality of the software in general.

So I came up with the idea of writing the API as an open-source project.
I am not running a private tracker with this, I am not writing a frontend for this.
I hope that staff of sites who might eventually run this software understand how open-source development works and communicate their ideas with me.
This can only work properly if we can keep the development as centralized as possible.
If every tracker grabs some alpha release and starts months of development on top of that, the community will not benefit from that work at all.

## Technology

- Language: Go
- Database: Postgres
- Transport: JSON
- Web framework: Iris, for now

## State of the project

- Pre-everything
- No tracker integration yet
- Nothing is being cached yet
- All errors are just spilled to the user - need distinct error types at some point
- Database layer basics somewhat fleshed out, proof-of-concept shows that unit-testing the DB layer works, is not too much work and is useful
- API layer draft. End-to-end testing of API methods is a little cumbersome but works
- API has no concept of rights, everyone can do everything for now
- Need a class/rights/permissions/... system

## Ideas

- Maybe allow bot accounts?
- All content is written in markdown and compiled... where?
- Have own markdown flavor, to make linking easier? (something like `:torrent:12345:` ?)

## FAQ

- Why is this written in Go??
    > Because I cba to write anything well-tested, reliable, and fast in PHP.
        Go has strong typing, is easy to deploy and runs pretty fast.
        We get out-of-the-box HTTP/2, a stdlib that does most of what  we need and arguably the best tooling for any language available.
- Why is this just an API??
    > Because that way it's much easier to write safe and clean code.
        We can guarantee that a middleware to check permissions runs for every request.
        We can ensure that every interaction with the data has to pass through the API.
        Things like rate limiting or permissions are trivial to implement.
- Are you running a private tracker??
    > No, I just write this software. Don't send me DMCA notices.
- Can I test it?
    > A public instance will be running at `api.boiling.rip:8234`, resets every 24 hours or so.
        You'll need an API key. 
        Or just run your own.
- Is this project affiliated with chihaya?
    > No. I think they're doing a great job, which is why I chose their tracker and their contributing guidelines.
- OMG why is this a music-only tracker??
    > We're using a relational database with actual types and structure.
        No, we're not gonna switch to NoSQL or Mongo or whatever.

## API Documentation

See `API.md`.

## Development

(This applies more and more the closer we get to a v1 release)  
We're using pull requests.
Adhere to the [chihaya contributing guidelines](https://github.com/chihaya/chihaya/blob/master/CONTRIBUTING.md).
PRs must pass travis and be reviewed.
Versioning will be dealt with when we have a v1 release.
No breaking changes then.

## License
MIT