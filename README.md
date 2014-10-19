ah
==

*This is totally incomplete description. Read it just to get an idea what is going on here*

ah is a simple externder of your terminal session history. Basically it has
to be used to capture outputs of your terminal command executions but it
also provides you with a couple of niceties.

ah has several commands. Let's do some short crash course where I will show
you all of possibilities.

Crash course
------------

So you want ~~to be a hero~~ to capture an output of some of your commands

```bash
$ find . -name "*.go"
./ah.go
./app/history_entries/parsers.go
./app/history_entries/get_commands.go
./app/history_entries/history_entry.go
./app/environments/environments.go
./app/utils/utils.go
./app/commands/bookmark.go
./app/commands/list_trace.go
./app/commands/tee.go
./app/commands/execute.go
./app/commands/show.go
./app/slices/slices.go
```

Nice but what are you going to do if you want to store an output? Use `tee`!

```bash
$ find . -name "*.go" 2>&1 | tee output.log
```

Okay. But what are you going to do if you delete `output.log` for some reason?
Or if you close your terminal? Use `ah`!

```bash
$ ah t -- find . -name "*.go"
./ah.go
./app/history_entries/parsers.go
./app/history_entries/get_commands.go
./app/history_entries/history_entry.go
./app/environments/environments.go
./app/utils/utils.go
./app/commands/bookmark.go
./app/commands/list_trace.go
./app/commands/tee.go
./app/commands/execute.go
./app/commands/show.go
./app/slices/slices.go
```

Okay, the same output. The same exit code. What is the profit? Let's check the
next command, `ah s`

```bash
$ ah s
...
!10045 *	./ah t -- find . -name "*.go"
!10046  	ah s
```

Okay, it looks like a normal `history` output. What is the profit? ah has
internal grep compatibility with fuzzy matching.

```bash
$ ah -g docker
...
!9482  	sudo docker rmi $(sudo docker images -q)
!9483  	sudo docker rm $(sudo docker ps -a -q)
!9484  	sudo docker rmi $(sudo docker images -q)
...
```

By the way, if you use `$HISTTIMEFORMAT` then ah will interpret it as you expect

Let's checkout fuzzy matching

```bash
$ ah s -z -g find
...
!10042  	find . -name "*.go"
!10045 *	./ah t -- find . -name "*.go"
!10049  	./ah -z -g find
...
```

Btw, have you noticed this little star mark? It means that log has been kept.

```bash
$ ah l 10045
./ah.go
./app/history_entries/parsers.go
./app/history_entries/get_commands.go
./app/history_entries/history_entry.go
./app/environments/environments.go
./app/utils/utils.go
./app/commands/bookmark.go
./app/commands/list_trace.go
./app/commands/tee.go
./app/commands/execute.go
./app/commands/show.go
./app/slices/slices.go
```

Cool right?

ah also has some cool features I will describe later.