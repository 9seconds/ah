ah
==

[![Build Status](https://travis-ci.org/9seconds/ah.svg?branch=master)](https://travis-ci.org/9seconds/ah)

ah is a complementary software for a builtin shell `history` command
you've used to use for years and I hope you've dreamt about it as I did.

It is not a replacement for `history` but anyway it perfectly matches
a common `history | grep` pattern of usage but it allows you to do a bit more.
It allows you to trace an output of a command, to fetch it from the archive,
to bookmark some commands and to execute them.

How often do you kick yourself for loosing important output of your SSH session?
Sometimes you use screen or tmux for that purposes but it is pretty awkward to
search through it. Or how often do you find yourself typing CTRL+R or juggling
with Up and Down buttons to find something like
`make && mv -f coolapp $DIR/bin && coolapp`. Stop it for a great good, ah will
likely help you here.

Currently it supports following features:
* Tracing an output
* Fetching the output trace
* Bookmarking command
* Executing of a command by number or bookmark
* Showing a history with greping on regular expression or fuzzy match

ah does not maintains its own history file, it uses your regular `~/.bash_history`
or `.zsh_history`. So no worries here: bash or zsh maintains a history and
ah gives you several features on the top.



Installation
------------

You may build ah from sources or just download proper binary from
[releases](https://github.com/9seconds/ah/releases).

To install it from sources, just do following:

```bash
$ git clone https://github.com/9seconds/ah.git $GOPATH/src/github.com/9seconds/ah
$ cd $GOPATH/src/github.com/9seconds/ah
$ make install
```

It will copy the binary into `$GOBIN/ah`. Or you may do that:

```bash
$ git clone https://github.com/9seconds/ah.git $GOPATH/src/github.com/9seconds/ah
$ cd $GOPATH/src/github.com/9seconds/ah
$ make
$ mv ah /wherever/you/want
```


Tracing an output
-----------------

So you want ~~to be a hero~~ to capture an output of some of your commands.
Usually if you know that output is rather important, you use `tee` command.

```bash
$ find . -name "*.go" -type f | tee files.log
```

The main problem here is that only stdout will go to the `files.log`. You will
lose stderr, right? The common way of solving that is redirecting a streams

```bash
$ find . -name "*.go" -type f 2>&1 | tee files.log
```

or

```bash
$ find . -name "*.go" -type f |& tee files.log
```

in a recent Bash-compatible shells. The problem here, you are streams are
dangerously mixed into one and there is no way to pipe them into different processes.
Let's say, you want to have output stored persistently but in the same time you
want to have only filtered log messages on the screen. You can't do something
like this.

```bash
$ find . -name "*.go" -type f |& tee files.log > /dev/null 2> grep -i localhost
```

Okay, this a rare case but please notice that ah knows how to handle that. Let's
talk on how to store this output, where to keep it. You may store it somehow
but it is a way better to have a tool which allows you not to think on how to
name this output and how to remember which was the command.

Let's talk about ah now. Ah has its main `t` command which works rather simple.

```bash
$ ah t -- find . -name "*.go" -type f
```

Thats all. You will see output on the screen and you may pipe both streams wherever
you want! Ah will store it persistently. And it will finish execution with
precisely the same exit code as the original command does. Neat, right?

If you want to run a program which requires a pseudo TTY, just use `-y` option.
And if you want to have your aliases to work, just run it with `-x` option!

Ah supports SSH and you may even run curses apps there, they will work, no worries.



Show the history
----------------

Now let's talk about viewing the history. ah does that with `s` command.

```bash
$ ah s
...
!10024  (01.11.14 16:01:09) *  ah t -- find . -name "*.go" -type f
```

What do we have here: we have banged command number (gues why it has ! here),
we have a date (yes, `HISTTIMEFORMAT` supported!), we have a rather strange
star mark and a command. What does that star mark mean? Basically it just shows
that ah keeps a mixed output of that command and you may fetch it on demand.

ah has `-g` options which allows you to grep this list. Argument - a regular
expression. It also has a convenient flag `-z` which activates fuzzy match. It
works like this

```bash
$ ah s -z -g doigreREPOsoru
```

And I will see matched `docker images | egrep -v 'REPOSITORY|<none>' | cut -d' ' -f1 | sort -u`
I bold important letters here. Basically I do it thinking like this:
*"I want __do__ cker __i__ mages, it was __gre__ p __REPO__ SITORY
and __sor__ ted with -__u__"* typing just a few letters.



Show an output
--------------

Output could be checked with `l` command. Just type `ah l 10024` and you are
good.



Bookmarks
---------

You may pin any command number with bookmark using `b` command. After that
you may execute it with `e` command. To fetch a list of bookmarks use `lb` commands,
to remove several, use `rb` command.

So simple.


Garbage collecting
------------------

If you do not need a lot of traces or bookmarks, you may get rid of them using
`gt` (garbage collect traces) and `gb` (garbage collect bookmarks) commands.
