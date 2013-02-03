gopack
======

<big>**Gopack** is a software dependency management tool for **Go**.

**Gopack** keeps packages **dependency** information to build the GOPATH variable, it also keeps **remote** locations where to publish and get packages.

It can then deliver:

* **Building** commands
    * compile
    * test
* **Managing** commands
    * list all dependencies (recursively)
    * find and list missing imports, and fix them
    * search for packages by name ( find the right version)
* **Sharing** commands
    * download packages from remote locations
    * publish packages to remote locations
    * search packages or imports in remote locations

</big>


Getting Started
============

Version:   1.0.0.beta.3


<small>
<img alt="Under construction" src="http://upload.wikimedia.org/wikipedia/commons/thumb/5/54/Under_construction_icon-green.svg/200px-Under_construction_icon-green.svg.png" height="33" width="40"/>
This section is still under development, it is kind of sparse (sorry)
</small>


installing it 

Linux
--------

<pre> 
git clone https://github.com/gopack/gpk.git
cd gpk/
export GOPATH=`pwd`
go install ./src/...
sudo cp ./bin/gpk /usr/bin/gpk
</pre>

Now you should get something like
<pre>$>gpk help</pre>

<pre>

NAME
       gpk - Gopack is a software dependency management tool for Golang.
             It help Managing, Building, and Sharing libraries in Go.

SYNOPSIS
       gpk [general options] <command> [options]  

OPTIONS
       option   default              usage
       -local   .gpkrepository       path to the local repository to be used by default.
       -v       false                Print the version number.


COMMANDS

       h        help                 Display help information about commands

       !        init                 Initialize or Edit the current project
       ?        status               Print status

       c        compile              Compile project
       t        test                 Run go test

       d+       dadd                 Add dependency
       d-       dremove              Remove dependency
       ld       list-dependencies    List declared Dependencies.
       lm       list-missing         Analyse the current directory and report or fix missing dependencies
       lp       list-package         List all packages dependencies (recursive)

       lr       list-remotes         List Remotes.
       r+       radd                 Add a remote server.
       r-       rremove              Remove a Remote

       i        install              Install into the local repository
       push     push                 Push a project in a remote repository
       s        search               Search Packages .
       serve    serve                Serve local repository as an http server


       Type 'gpk help [COMMAND]' for more details about a command.
</pre>
 but prettier if you are using a vterm console

Examples
=========

Hello World Project
---------------

The purpose is to show the smallest possible Helloworld project.
Create the Workspace layout.
<pre>
$> mkdir test
$> cd test
$> gpk init -c -n mypath/test -l ASF
    new name:mypath/test
    new license:"Apache License 2.0"
</pre>
Creates the workspace, and the directory layout. Its time to populate it with the helloworld.go file
<pre>
$> vi src/mypath/test/helloworld.go
</pre>
Edit the helloworld.go file and make it:
<pre>
package main
import "fmt"
func main() {
    fmt.Println("Hello, World")
}
</pre>

Then Compile, and run
<pre>
$> gpk compile
$> ./bin/test
</pre>


Client/Server
-------------------

Purpose: show how to start a basic gopack package server.

On computer called 'server' lets start a server (it does not need to be on a specific directory)
<pre>
eric@server:$ gpk serve
    starting server :8080
</pre>
It will by default expose the local repository

On another computer called 'client' let's connect to this server
<pre>eric@client:$ gpk r+ server http://192.168.0.30:8080
    new remote: quere http://192.168.0.30:8080
</pre>
Lets search for stuff in it
<pre>eric@client:$gpk search -r server ericaro.net</pre>
the result list is empty

Let's install the current project (gopack) as "0.0.0-master" in the local repository of "client", and push it to "server"
I assume that I'm on a project called ericaro.net/gopack
<pre>eric@client:$ gpk install master</pre>
<pre>eric@server:$ gpk push server ericaro.net/gopack master
Success
</pre>
Note that "master" is a valid name for the [semantic version](http://semver.org) 0.0.0-master.

On the server side here is what has happened
<pre>
RECEIVING
       ericaro.net/gopack master GNU Lesser GPL INTO /home/eric/.gpkrepository/ericaro.net/gopack/master</pre>
</pre>

Now, on the client side, if we search for package called ericaro.net we found one.

<pre>eric@client:$ gpk s -r server ericaro.net</pre>


<h2>Dependencies</h2>
 

Lets work on another project, and we added some imports in the code:
<pre>import (
    "ericaro.net/gopack/"
    "ericaro.net/gopack/protocol"
    "ericaro.net/gopack/semver"

)</pre>

is it one two or three packages or only one ?

The project will not compile, right ?
<pre>
$ gpk c
src/myproject/gae/services.go:6:2: import "ericaro.net/gopack": cannot find package
src/myproject/gae/entities.go:8:2: import "ericaro.net/gopack/protocol": cannot find package
src/myproject/gae/entities.go:9:2: import "ericaro.net/gopack/semver": cannot find package
exit status 1
</pre>
<small>'gpk c' is short for 'gpk compile'</small>

Let's fix the project

<pre>$gpk lm -f
Missing imports (3), missing packages (1)
Missing packages ericaro.net/gopack                       -> ☑ ericaro.net/gopack 1.0.0-beta.1 
                                                          -> ☐ ericaro.net/gopack master
Project Updated
</pre>
<small>

* 'gpk lm' is short for 'gpk list-missing'
* -f stands for 'fix'

</small>

Now your project compiles, and the dependency is under control:
<pre>$ gpk ld

LIST OF DECLARED DEPENDENCIES:
        google.com/appengine                     1.7.3
        ericaro.net/gopack                       1.0.0-beta.1
</pre>
<small>'gpk ld' is short for 'gpk list-dependencies'
</small>


<pre>$ gpk lp -l

LIST OF PACKAGES:
        google.com/appengine                     1.7.3
        code.google.com/p/goprotobuf/proto       default
        ericaro.net/gopack                       1.0.0-beta.1

</pre>
<small>

* 'gpk lp' is short for 'gpk list-package'
* -l stands for 'list' by default 'gpk lp' prints the dependencies in a GOPATH way
</small>

**Tip** 

Typing
           <pre>alias GP='export GOPATH=`gpk lp`'</pre>
In your shell is an easy way to to get an automatic GOPATH setter.


