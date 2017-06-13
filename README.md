# tiny-go-server

A minimalistic, single-module, dependency-less two-directory Web server in Go.

## Purpose

I wanted a no-frills Web server to act as a repository librarian for a small set of 
[Lua](https://www.lua.org/) scripts written for an [OpenComputers](http://ocdoc.cil.li/) 
environment. Virtual "robots" in MC, running a very bare-bones 
[wget](https://www.gnu.org/software/wget/) knockoff, will be interacting with this server.

## Description

Specific requirements:

* A landing page (to provide the other URLs);
* Serve directory listings and files from two directories `lib` and `upl`;
* Upload files to `upl` (but not `lib`);
* Accept single-line "log" messages and display them in the server's log.

## Implementation

Go's [standard packages](https://golang.org/pkg/) provide enough scaffolding for a simple Web server
from a single source module. This code does exactly what I want (modulo the [ToDo](ToDo.md)s), 
nothing more and nothing less, and I can control its feature set.

## The fine print

The [LICENSE](LICENSE) is MIT, which means that anyone can do whatever they want with the code.

**SECURITY WAS NOT CONSIDERED IN THIS SOFTWARE** 
so it might have (unintended!) security issues.

As a result, I recommend you don't use it in a security-critical environment.

## Downloads

Two files you may be interested in:

* The source, [tiny-go-server.go](src/github.com/CarlSmotricz/tiny-go-server/tiny_go_server.go)
* The binary, [tiny-go-server](bin/tiny-go-server), for Linux X86-64 only.

## [ToDo-s](ToDo.md)
