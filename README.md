gopack
======

Version:   1.0.0.beta.1

Gopack is a software dependency management tool for Go.

Gopack provides support for *Managing*, *Building*, *Sharing* libraries in [Go](http://golang.org).


Gopack provides a *distributed* system to share libraries.

Once connected to a node you can link your project with the rest of the world, and use those libraries.

You can be part of this world too, by serving your packages too.

DEFINITIONS
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

get it
<pre>go get github.com/eatienza/gopack</pre>
and install it. That's it, gpk comes as a standalone, executable. Type <pre>gpk help</pre> to use the built-in help.

