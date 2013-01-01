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
    * list missing imports, and fix them
    * search for packages
* **Sharing** commands
    * download packages from remote locations
    * publish packages to remote locations
    * search packages or imports in remote locations

</big>


Getting Started
============

Version:   1.0.0.beta.1


<small>
<img alt="Under construction" src="http://upload.wikimedia.org/wikipedia/commons/thumb/5/54/Under_construction_icon-green.svg/200px-Under_construction_icon-green.svg.png" height="33" width="40"/>
This section is still under development, it is kind of sparse (sorry)
</small>


installing it 

Linux
--------

<pre> 
git clone git@github.com:eatienza/gopack.git
cd gopack/
export GOPATH=`pwd`
go install ./src/...
sudo cp ./bin/gpk /usr/bin/gpk
</pre>

Now you should have
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


Example
=========

Client/Server
-------------
On computer called 'server' lets start a server (it does not need to be on a specific directory)
<pre>
eric@server:$ gpk serve
    starting server :8080
</pre>
It will by default expose the local repository

On another computer called 'client' let's connect to this server
<pre>eric@client:$ gpk r+ server http://192.168.0.30:8080
    new remote: server http://192.168.0.30:8080
</pre>
Lets search for stuff in it
<pre>eric@client:$gpk search -r server ericaro.net</pre>
the result list is empty

Let's install the current project (gopack) as "0.0.0-master" in the local repository of client, and push it to server
<pre>eric@client:$ gpk install master</pre>
<pre>eric@client:$ gpk push server ericaro.net/gopack master
Success
</pre>
Note that "master" is a valid name for the [semantic version](http://semver.org) 0.0.0-master.

On the server side here is what has happened
<pre>
RECEIVING
       ericaro.net/gopack master GNU Lesser GPL INTO /home/eric/.gpkrepository/ericaro.net/gopack/master</pre>
<small>due to issue #4 the output is not exactly the one above</small>

Now, on the client side, if we search for package called ericaro.net we found one.
<pre>eric@client:$ gpk s -r server ericaro.net
    ericaro.net/gopack                       master</pre>
