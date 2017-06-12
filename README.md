# tiny-go-server

A minimalistic, single-module, dependency-less two-directory Web server in Go.

# PURPOSE

I wanted a no-frills Web server to act as a repository librarian for a small set of 
[Lua](https://www.lua.org/) scripts written for an [OpenComputers](http://ocdoc.cil.li/) 
environment. Virtual "robots" in MC, running a very bare-bones 
[wget](https://www.gnu.org/software/wget/) knockoff will be interacting with this server.

# DESCRIPTION

Specific requirements:

* A landing page (to provide the other URLs);
* Serve directory listings and files from two directories `lib` and `upl`;
* Upload files to `upl` (but not `lib`);
* Accept single-line "log" messages and display them in the server's log.

# IMPLEMENTATION

Go's [standard packages](https://golang.org/pkg/) provide enough scaffolding for a simple Web server
from a single source module. This code does exactly what I want (modulo the [ToDo](ToDo.md)s), 
nothing more and nothing less, and I can control its feature set.

The [LICENSE](LICENSE) is MIT, which means that anyone can do whatever they want with the code.

**THIS SERVER MAY HAVE SECURITY ISSUES ** (though certainly not intentionally), 
so I recommend you don't use it in a critical environment.

# DOWNLOADS

Two files you may be interested in:

* The source, [tiny-go-server.go](blob/master/src/github.com/CarlSmotricz/tiny-go-server/tiny-go-server.go])
* The binary, [tiny-go-server](blob/master/bin/tiny-go-server), for Linux X86-64 only.
