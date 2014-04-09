# Finch

[![GoDoc](https://godoc.org/github.com/BrianHicks/finch?status.png)](https://godoc.org/github.com/BrianHicks/finch)
[![Build Status](https://travis-ci.org/BrianHicks/finch.png?branch=master)](https://travis-ci.org/BrianHicks/finch)
[![Download](https://api.bintray.com/packages/brianhicks/finch/finch/images/download.png)][download]

A task manager, written in Go. Finch implements the [Final Version][fv]
algorithm by [Mark Forster][mf].

## Installation

    go get github.com/BrianHicks/finch/finch

or [Download from BinTray][download]

## Usage

To sum up how this works: add tasks to a list, mark some as "selected" (always
starting with the oldest), then work through your selected list
newest-to-oldest. In Finch, the sequence would look like this:

```
$ finch add email Jane about the meeting
Added: be: email Jane about the meeting

$ finch add fix the holodeck
Added: bi: fix the holodeck

$ finch add eat lunch
Added: bo: eat lunch

$ finch available
be: email Jane about the meeting
bi: fix the holodeck
bo: eat lunch

$ finch select 0 2
Saved: be, bo

$ finch next
bo: eat lunch (*)

# after lunch...

$ finch done
Marked done: bo: eat lunch (done)

$ finch next
be: email Jane about the meeting (*)

# but you decide that "email Jane about the meeting" can't be done now, or you
# do a little bit of work on it and are interrupted. FV says you should
# re-enter the task at the end of the list. So we have "delay" as well.

$ finch delay
Delaying: be: email Jane about the meeting (*)
Delayed until 2014-04-09 09:53:10.119418346 -0500 CDT

$ finch select
bi: fix the holodeck
be: email Jane about the meeting
```

To get descriptions of the commands, run `finch help` or `finch help [command]`.

## Contributing/Development/Hacking

Finch is licensed under the [Apache License](LICENSE.txt). By contributing,
you agree to licensing your contributions under the Apache License and bound to
provide only code you have the right to license in that manner.

If you just want something to do, check out the [low-hanging fruit][lhf] tag on
the [issues][issues] for this repository.

With that out of the way, `go get github.com/BrianHicks/finch` (or `git clone`,
of course) should do the trick. You can run `make deps` to get everything.
Before you commit, please run `make test`, `make lint`, and `gofmt`. It'll make
everyone's lives easier.

Work on unreleased versions is done on "develop" and feature branches. If
anyone (including me) commits to master I'll be very unhappy. So to start
working on this using [hub][hub]:

```
go get github.com/BrianHicks/finch
cd $GOPATH/src/github.com/BrianHicks/finch
hub fork
git checkout develop
git checkout -b feature/yourfeature
```

If you have any thoughts about build tools please tell me! My email address is
in my Github profile if you don't want to just open an issue. This is my first
"real" go project and I'm open to changing things to fit with what the
community usually expects.

[mf]: http://markforster.squarespace.com/ "Get Everything Done - Mark Forster"
[fv]: http://archive.constantcontact.com/fs004/1100358239599/archive/1109980854493.html "Final Version"
[hub]: https://github.com/github/hub "Hub - hub helps you win at git"
[lhf]: https://github.com/brianhicks/finch/issues?labels=low-hanging+fruit&page=1&state=open "Issues • BrianHicks/finch"
[issues]: https://github.com/brianhicks/finch/issues?state=open "Issues • BrianHicks/finch"
[download]: https://bintray.com/brianhicks/finch/finch/_latestVersion
