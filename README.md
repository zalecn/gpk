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


get it
<pre>go get github.com/eatienza/gopack</pre>

and install it. That's it, gpk comes as a standalone, executable. Type <pre>gpk help</pre> to use the built-in help.
