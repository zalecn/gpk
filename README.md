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

