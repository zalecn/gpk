Gopack Documentation
====================

Copyright
:   public domain

Version
:   0.1

introduction
------------

[Gopack](http://gpk.ericaro.net) is a command line tool for dependencies
management in [Golang](http://golang.org)

To fulfill his duty, [Gopack](http://gpk.ericaro.net) defines a few
concepts we are going to introduce here.

Project
:   A [Golang](http://golang.org) package project, is the directory
    containing the source code for a given package. Usually is contains
    a src, a bin and a pkg directory.

Local repository
:   A gopack local repository, is a directory containing only packages.
    They are organize by package name then version. coucou papa

Remote repository
:   A remote repository is defined by its URL.
    [Gopack](http://gpk.ericaro.net) supports \`[file://](file://)\` and
    \`[http://](http://)\`. Gopack can start a server that turn any
    local repository into a remote repository for others.

Package
:   A package is a snapshot of a project' sources. It is fully defined
    by its name, and a version

Package name
:   [Gopack](http://gpk.ericaro.net) follows [Golang](http://golang.org)
    rules. Any string is suitable as a package name, but we recommend
    the host/path parts of a URL you own, and that you can guarantee to
    be unique.

Package Version
:   [Gopack](http://gpk.ericaro.net) follow the [semantic version
    2.0-rc.1](Semver_). Gopack makes a special case of version
    \`0.0.0-something\`. They are considered to be *snapshot* version.
    Every version parser involved in [Gopack](http://gpk.ericaro.net)
    parse \`something\` into \`0.0.0-something\`. Gopack also displays
    \`0.0.0-something\` as \`something\` to keep it short.

Project dependencies
:   Dependencies are references to other packages. They are formed of
    the package name + a version.

Gopack stores Project dependencies into a file (.gpk) at the root of the
project's directory. Every time a command is executed,
[Gopack](http://gpk.ericaro.net) opens this small file and retrieve the
information. [Gopack](http://gpk.ericaro.net) provides several commands
to edit this file.

Gopack resolves recursively every project's dependencies. Resolving a
dependency is finding a directory in the local repository containing the
appropriate version of the package. This resolution can be used directly
to manage your GOPATH variable. A good trick is to define

This way you when entering a project you can set your local GOPATH
variable by typing

If the dependencies is missing from the local repository, Gopack can try
to find it in remote repositories. There is a central repository, but
you can add other remote repository. For instance your corp central
repository.

Gopack uses this resolve capacity to run several go commands, like
\`install\`, or \`test\`

Usually the next step after compiling and test a project having other's
dependencies, is to make your project available to others, in your
organization, or even publicly.

Gopack provides a set of command to install your projects in your local
repository. This is the first step to work with several projects all
together.

Gopack offers also to ability to "push" a package from your local
repository to any arbitrary remote repository. You can build the
"dependency management" you want. If you don't want to push your
project, you can serve your local repository, and share your IP address
so that others can pull packages from your repository.

Every workflow is possible, push, pull, central.

Getting started
---------------

Follow those simple steps to get started.

### Install

> -   download binary (TODO provide link to the binaries)
> -   put in your path

### Start using

> -   create a small project (find a dependency, ensure that this
>     dependency is on the central server).

\>\>\> gpk init

\>\>\> gpk + ericaro.net/xxx

\>\>\> gpk x (note that this will download the package from central)

\>\>\> tester (it works)

### Share

\>\>\> gpk install (make it available to your other projects)

\>\>\> gpk push central your.package/test 1.0.0 (cave at this is for
every)

\>\>\> gpk r+ corp
[http://central.mycorp.com](http://central.mycorp.com)

\>\>\> gpk push corp mypackage 1.0.0

Commands
--------
