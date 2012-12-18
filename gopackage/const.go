//Package gopackage handles programmatic API for go package dependencies (gpk)  operations.
// It defines a file .gpk that contains the current project definition, and the current project dependencies.
// It provides tools to edit/read this file.
//
//It provides tools to package, download upload project to a local or remote server. 
package gopackage

import (
"net/url"
)

const (
	Cmd               = "gpk"
	GopackageFile     = ".gpk"
	GopackageVersion  = "0.0.0.1"
	DefaultRepository = ".gpkrepository"
	Release           = "Release"
	Snapshot          = "Snapshot"
)

// TODO parse in a yaml or xml format


// I would like to revamp the whole API:
// there is :
// version := root + 4 digits
// Project ref = name + version 
// Project = project ref (self) + (projec ref)*
// project pack := project (self) + snapshot/release + ( extra stuff to come, like sha1, signature etc).

// there is a tree of repo
// there are all connected for "get"
// I can push to any of them as far as I've got the key !

// some repo support rewrite, some don't, the root does not, 
// a package is release if it belongs to a read only repo.
// when pushing a package to a repo, it checks if it already existing in the tree
// and it checks for the permission to write.
// => that' the way a "release" is made:
// => a stagging repo is a repo that accept only one write



// I would also like to remove the local repo but it ain't possible, I need a layout where to build "gopath" from.
// I need a way to build "proxies" so that I can freely push to it.
