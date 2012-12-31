package gopack

import (
	"ericaro.net/gopack/gocmd"
	"ericaro.net/gopack/semver"

	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// provide missing analysis elements

//MissingImports computes the import missing in this project.
func (r *LocalRepository) MissingImports(p *Project, offline bool) (missing []string) {
	missing = make([]string, 0)
	known := make(map[string]bool)

	// run the go build command for local src, and with the appropriate gopath
	dependencies, err := r.ResolveDependencies(p, offline, false)
	if err != nil {
		return
	}
	gopath, err := r.GoPath(dependencies)
	if err != nil {
		return
	}

	build.Default.GOPATH = gocmd.Join(gopath, p.WorkingDir())
	srcDirs := build.Default.SrcDirs()
	// srcDirs is where to lookup at packages 

	srcDir := filepath.Join(p.WorkingDir(), "src")
	// srcDir is where to scan for sources imports 
	packages, err := ScanDir(srcDir)
	if err != nil {
		ErrorStyle.Printf("error %s\n", err)
		return
	}

	fset := token.NewFileSet() // positions are relative to fset
	for _, dir := range packages {
		pkgs, _ := parser.ParseDir(fset, dir, nil, parser.ImportsOnly)
		//			if err != nil {
		//				ErrorStyle.Printf("scan error %s\n", err)
		//			}
		for _, pk := range pkgs {
			for _, f := range pk.Files {
				for _, imp := range f.Imports {

					i, _ := strconv.Unquote(imp.Path.Value)
					if _, ok := known[i]; !ok { // its a new import to process
						if !isContained(srcDirs, i) {
							// it's a hit ! this is a missing import
							missing = append(missing, i)
						}
						known[i] = true // always mark the import to skip later queries
					}
				}
			}
		}

	}
	return
}

func isContained(srcDirs []string, imp string) bool {

	for _, src := range srcDirs {
		dst := filepath.Join(src, imp)
		if s, err := os.Stat(dst); err == nil {
			if s.IsDir() {
				return true
			}
		}

	}
	return false
}

//MissingPackages merge all missing import into what should be packages
func (r *LocalRepository) MissingPackages(missingImports []string) (missing []string) {
	missing = make([]string, 0)
	for _, m := range missingImports {
		dominated := false
		for _, p := range missingImports {
			if strings.HasPrefix(m, p) && m != p { // it mean dominated
				dominated = true
				break // we found what should be the package
			}
		}
		if !dominated {
			missing = append(missing, m)
		}
	}
	return
}

func (r *LocalRepository) ResolvePackages(missingPackages []string, offline bool) (missing []*ProjectID) {
	missing = make([]*ProjectID, len(missingPackages))

	//idea is to lookup in the local repo for the packages

	return
}

func (r *LocalRepository) ImportSearch(imp string) (pkg []ProjectID) {

	pkg = make([]ProjectID, 0)

	handler := func(srcpath string) bool {
		dst := filepath.Join(srcpath, "src", imp)
		_, err := os.Stat(dst)
		if err == nil {
			pk, _ := filepath.Rel(r.Root(), filepath.Dir(srcpath))
			version := filepath.Base(srcpath)
			v, _ := semver.ParseVersion(version)
			if err == nil {
				pkg = append(pkg, *NewProjectID(pk, v))
			}
			return true // stop searching
		}
		return true
	}

	PackageWalker(r.root, "", handler) // call s handler on every 'package'
	sort.Sort(reverse{ProjectIDs(pkg)})
	return

}

type reverse struct {
	sort.Interface
}

func (r reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}
