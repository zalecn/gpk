

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

    
    
.. include:: links.rst
    