gopack
======

Version:   1.0.0.beta.1

<big>Gopack is a software dependency management tool for Go.</big>

Gopack keeps packages *dependency* information to build the GOPATH variable, it also keeps *remote* locations where to publish and get packages.

It can then deliver:

* Building commands
    * compile
    * test
* Managing commands
    * list all dependencies (recursively)
    * list missing imports, and fix them
    * search for packages
* Sharing commands
    * download packages from remote locations
    * publish packages to remote locations
    * search packages or imports in remote locations
    
Gopack provides support for *Building*, *Managing*, *Sharing* libraries in [Go](http://golang.org).



Definitions
-----------

**Project**:  A [Golang](http://golang.org) package project, is the directory containing the source code for a given package. Usually it contains a src, a bin and a pkg directory.

**Local repository**: A gopack local repository, is a directory containing package, downloaded or installed. They are organize by package name then version.

**Remote repository**: A remote repository is defined by its URL. Gopack supports `file://` and `http://`. Gopack can connect to several remotes, or start a server that turn any local repository into a remote repository for others.

**Package**: A package is a snapshot of a project' sources. It is fully defined by its name, and a version.

**Package name**: Gopack follows [Golang](http://golang.org) rules. Any string is suitable as a package name, but we recommend the host/path parts of a URL you own, and that you can guarantee to be unique.

**Package Version**: Gopack follows the [semantic version 2.0-rc.1](http://semver/org).

**Project dependencies**: Dependencies are references to other packages. They are formed of the package name and a version.



Getting Started
---------------

<small>
<img alt="Under construction" src="http://upload.wikimedia.org/wikipedia/commons/thumb/5/54/Under_construction_icon-green.svg/200px-Under_construction_icon-green.svg.png" height="33" width="40"/>
This section is still under development, it is kind of sparse (sorry)
</small>


get it
<pre>go get github.com/eatienza/gopack</pre>

and install it. That's it, gpk comes as a standalone, executable. Type <pre>gpk help</pre> to use the built-in help.

