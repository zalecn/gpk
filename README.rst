===
gpk
===

:Author: eric@ericaro.net
:Date:   2012-12-22
:Version: 1.0.0.beta.1
:Manual section: 1
:Manual group: compiler


.. _Golang: http://golang.org
.. _Gopack: http://gpk.ericaro.net
.. _Semver: http://semver.org/



Gopack is a software project management tool for Golang_.

Gopack can:

  * Resolve dependencies and install a golang project
  * Share a project
    * between sub projects
    * within a team
    * worldwide as a reusable library

Gopack main features are:

  * Fixes the GOPATH issue.
  * Simple go project start
  * Simple and safe dependency import
  * Multiple project integration

SYNOPSIS
========

  gpk [options] command

DEFINITIONS
===========



Project
    A Golang_ package project, is the directory containing the source code for a given package. 
    Usually is contains a src, a bin and a pkg directory.
Local repository
    A gopack local repository, is a directory containing only packages. 
    They are organize by package name then version.

Remote repository
    A remote repository is defined by its URL. Gopack_ supports `file://` and `http://`. 
    Gopack can start a server that turn any local repository into a remote repository for others.

Package
    A package is a snapshot of a project' sources. It is fully defined by its name, and a version

Package name
    Gopack_ follows Golang_ rules. Any string is suitable as a package name, 
    but we recommend the host/path parts of a URL you own, and that you can guarantee to be unique.

Package Version
    Gopack_ follow the `semantic version 2.0-rc.1`__. 
    Gopack makes a special case of version `0.0.0-something`. They are considered to be *snapshot* version. 
    Every version parser involved in Gopack_ parse `something` into `0.0.0-something`. 
    Gopack also displays `0.0.0-something`  as `something` to keep it short.

__ Semver_

Project dependencies
    Dependencies are references to other packages. 
    They are formed of the package name + a version.


COMMANDS
========


gpk ! -n NAME -l LICENSE
    Initialize or Edit the current project

gpk i VERSION
    Install into the local repository

gpk s QUERY
    Search Packages .

gpk ? 
    Print status

gpk d+ NAME VERSION
    Add dependency

gpk d- NAME
    Remove dependency

gpk c 
    Compile project

gpk t 
    Run go test

gpk h [COMMAND]
    Display help information about commands

gpk ld 
    List declared Dependencies.

gpk lr 
    List Remotes.

gpk lp 
    List all packages dependencies (recursive)

gpk lm 
    Analyse the current directory and report or fix missing dependencies

gpk serve ADDR
    Serve local repository as an http server

gpk push REMOTE PACKAGE VERSION
    Push a project in a remote repository

gpk r+ NAME URL [TOKEN]
    Add a remote server.

gpk r- NAME
    Remove a Remote




