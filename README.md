# gopack


<big>**Gopack** is a software dependency management tool for **Go**.

**Gopack** keeps package's **dependency** information to build the GOPATH variable. It also keeps **remote** locations where to publish and get packages.

It can then deliver:

* **Building** commands
    * compile. Single platform or cross compile
    * test. Run them on the go or cross compile them for later tests.
* **Managing** commands
    * list all dependencies (recursively)
    * find and list missing imports, and fix them
    * search for packages by name ( find the right version)
* **Sharing** commands
    * download packages from remote locations
    * publish packages and binaries to remote locations
    * search packages or imports in remote locations

</big>

# Documentation

[Read the Wiki](wiki)

# Getting Started


<small>
<img alt="Under construction" src="http://upload.wikimedia.org/wikipedia/commons/thumb/5/54/Under_construction_icon-green.svg/200px-Under_construction_icon-green.svg.png" height="33" width="40"/>
This section is still under development, it is kind of sparse (sorry)
</small>


installing it
-----------------

Very soon, gopack will be available for direct download in every supported platform. Meanwhile, you still need to build it.

### Linux


<pre> 
git clone https://github.com/gopack/gpk.git
cd gpk/
GOPATH=`pwd` go install ./src/...
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


Hello World
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


# News

2013-02: 

1.0.0-beta.5
   Add support for cross platform binaries, and test binaries compilation too
   Add support for binaries upload/download from [gopack servers](wiki/Setting-a-gopack-LAN-server).

2012-12: Project inception

