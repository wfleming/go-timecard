go-timecard
===========

[![Circle CI](https://circleci.com/gh/wfleming/go-timecard.svg?style=svg)](https://circleci.com/gh/wfleming/go-timecard)


A simple timecard tool written in Go.

Installation
------------

`go get github.com/wfleming/go-timecard/punch`

Usage
-----

To start working on something:

`punch in my-fun-project`

To stop working on something:

`punch out`

When you want to see how much time you've been spending on things:

`punch summary`

When you want to see your times in a spreadsheet:

`punch summary --csv > time-report.csv`

FAQ
---

1. **Why?**

    Mostly because I wanted an excuse to dig into Go a bit more. But also
    because I wanted a simple time tracking tool I could use from my Terminal
    (since that's usually where I am anyway), and none of
    [the](http://taskwarrior.org)
    [existing](https://github.com/samg/timetrap)
    [tools](https://bitbucket.org/latestrevision/timebook/src)
    were quite what I wanted.

    That said, this is just a side project, mostly done as a learning exercise.
    Future support/development should not be expected. Caveat emptor and all
    that.

2. **I want syncing of my time across computers**

    `mkdir ~/Dropbox/.punch && ln -s ~/Dropbox/.punch ~/.punch`

    Boom.

3. **I want to filter projects I'm looking at in the summary view**

    `grep my-fun-project ~/.punch/entries.log | punch summary`

4. **I want to see only recent dates in my summary view**

    `tail -20 ~/.punch/entries.log | punch summary`

    NB: Right now `punch` expects all event streams to start with an `in` event.
    So, assuming you are currently punched out when requesting a summary, make
    sure to specify an even number of lines. Conversely, if you are currently
    punched in, make sure to specify an odd number.

5. **You didn't answer my question**

    Well, nobody actually asked me any questions, I was just making them up.
    File an issue or a pull request, if you feel the urge. Or message me
    [on Twitter](http://twitter.com/thorisalaptop).
