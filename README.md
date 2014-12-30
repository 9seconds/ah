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

ah supports Zsh and Bash.



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

Also, if you use [HomeBrew](http://brew.sh)/[LinuxBrew](http://brew.sh/linuxbrew), you may want to check the formula:

```bash
$ brew tap 9seconds/homebrew-ah
$ brew install ah
```

Update with Brews are trivial

```bash
$ brew update
$ brew reinstall ah
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

It also supports number argument. Let's say `ah s 10` will show latest 10
commands, `ah s 10 20` will show commands from 10 to 20. Also negative numbers
are supported (but with underscore prefix, not hyphen), they are mostly work
as Python slices. `ah s 10 _20` means literally "from 10 to the latest 20".
Basically `ah s 10` equal to `ah s _10 _1`



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


Automatic execution
-------------------

_(Only for zsh)_

Ah allows user to set it for automatic tracing outputs. To do that please do
the following:

1. You have to download script to [source](https://raw.githubusercontent.com/9seconds/ah/master/sourceit/zsh.sh)
   and place it wherever you want. Let's say, I put it to `~/.auto_ah.sh`
2. Add following line to your `~/.zshrc`: `source ~/.zshrc` and please be noticed
   that `ah` should be in your `PATH`. If `which ah` works, then you're done.
   Otherwise just add it to the path.

Basically, ah will track all your executions automatically. But since it is
dangerous to execute automatically everything around, there is a whitelist. ah
has 3 commands you should be interested in:

1. `al` shows the list of all whitelist command ah will automatically apply to.
2. `ar` removes command from the whitelist and ah won't execute it automatically.
3. `ad` add a command to the whitelist.

Let's check my current setup.

```bash
$ ah al
ag                   [interactive=false, pseudoTty=false]
aptg                 [interactive=true , pseudoTty=false]
awk                  [interactive=false, pseudoTty=false]
docker               [interactive=false, pseudoTty=true ]
docker_clean         [interactive=true , pseudoTty=false]
docker_stop          [interactive=true , pseudoTty=false]
docker_update        [interactive=true , pseudoTty=false]
find                 [interactive=false, pseudoTty=false]
grep                 [interactive=false, pseudoTty=false]
ipython              [interactive=true , pseudoTty=true ]
make                 [interactive=false, pseudoTty=false]
python               [interactive=false, pseudoTty=false]
sed                  [interactive=false, pseudoTty=false]
ssh                  [interactive=false, pseudoTty=false]
vagrant              [interactive=true , pseudoTty=true ]
```

As you can see, I have a mixed setup. I trace an output of `ag` or `find` command
and do it in non-interactive (interactive means `zsh -i -c`) way and do not
allocate pseudo TTY for them. There are several aliases (`docker_update` or `aptg`)
and to execute them I use interactive mode. And I use pseudo TTY for `ipython`.

Now let's add `go` for the list.

```bash
$ ah ad go
```

So simple. But how can I set interactiveness or pseudo TTYs? Pretty simple and
obvious (remember `t` command?)

```bash
$ ah ad -x go
```

for interactiveness. And for pseudo TTY

```bash
$ ah ad -y go
```

If you decide to use another set of options, just execute `ah ad` with another
set of options, it will override previous setting.

To remove command just use `ar`

```bash
$ ah ar go
```

No need to resource or do something more.


Configuration
-------------

ah supports configuration with YAML file. It should be placed in `~/.ah/config.yaml`.
Here is the full example (everything may be omit)

```
shell: zsh
histfile: /home/9seconds/.zsh_history
histtimeformat: "%d.%m.%y %H:%M:%S"

tmpdir: /tmp
```

That simple, yes. It is useful, if you bring a lot of commandline options in aliases
or if you want to execute ah automatically.

Here is the sequence of argument overriding:

1. Default options
2. Config options
3. Commandline options

So commandline options overrides config.
