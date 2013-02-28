package gopack

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"log"
)

// Some anti-pattern ioutils: so kept private to this package

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func JsonReadFile(path string, v interface{}) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}

func JsonWriteFile(path string, v interface{}) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(v)
}

//walkDir recursively walk into src directory  and fire callbacks to dirHandler, and fileHandler
// it changes the dst to reflect the src relative path, meaning that if there is a src/foo dir, the handlers will be called with a dst/foo, and src/foo , this way
// it is easy to either copy dir, or tar.gz dir (not changing the structure though)
func walkDir(dst, src string, dirHandler, fileHandler func(dst, src string) error) error {
	if dirHandler != nil {
		dirHandler(dst, src)
	}

	file, err := os.Open(src)
	if err != nil {
		return err
	}
	subdir, err := file.Readdir(-1)
	if err != nil {
		return err
	}

	for _, fi := range subdir {
		ndst, nsrc := filepath.Join(dst, fi.Name()), filepath.Join(src, fi.Name())
		if fi.IsDir() {
			err = walkDir(ndst, nsrc, dirHandler, fileHandler)
			if err != nil {
				return err
			}
		} else {
			err := fileHandler(ndst, nsrc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ScanDir(src string) (dir []string, err error) {
	file, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	subdir, err := file.Readdir(-1)
	if err != nil {
		return dir, err
	}
	dir = make([]string, 0)

	for _, fi := range subdir {
		if fi.IsDir() {
			nsrc := filepath.Join(src, fi.Name())
			dir = append(dir, nsrc)
			ndir, err := ScanDir(nsrc)
			if err != nil {
				return dir, err
			}
			dir = append(dir, ndir...)
		}
	}
	return
}

//ScanPackages scan for subdirectories containing .go files
func ScanPackages(src string) (dir []string, err error) {
	log.Printf("scanning for Package in %s", src)
	file, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	appended := false // current src will be appended only if it contains at least one go file
	subdir, err := file.Readdir(-1)
	if err != nil {
		return dir, err
	}
	dir = make([]string, 0)

	for _, fi := range subdir {
		if fi.IsDir() {
			nsrc := filepath.Join(src, fi.Name())
			ndir, err := ScanPackages(nsrc)
			if err != nil {
				return dir, err
			}
			dir = append(dir, ndir...)
		} else {
			if !appended && strings.HasSuffix(fi.Name(), ".go") {
				// there is a go file // append the parent's (only once)
				log.Printf("Current package %s has a .go file (%s) ", src, fi.Name() )
				dir = append(dir, src)
				appended = true
			}
		}
	}
	return
}

//PackageWalker recursively scan a directory for packages ( identified as directory containing a .gpk file
// calls the handler with those directory until the handler returns false, or the directory tree has been exhausted.
// when he has found a package, it no longer look into it.
func PackageWalker(srcpath, startwith string, handler func(gpkpath string) bool) (c bool, err error) {
	c = true
	file, err := os.Open(srcpath)

	if err != nil {
		return true, err
	}
	subdir, err := file.Readdir(-1)

	if err != nil {
		return true, err
	}

	for _, fi := range subdir {
		if fi.IsDir() && strings.HasPrefix(fi.Name(), startwith) {
			c, err = PackageWalker(filepath.Join(srcpath, fi.Name()), "", handler)
			if !c {
				break
			}
		} else {
			if fi.Name() == GpkFile {
				//then src path is the package/version directory
				c = handler(srcpath)
				break
			}
		}
	}
	return
}

//Untar reads the .gopackage file within the tar in memory. It does not set the Root
func Unpack(dst string, in io.Reader) (err error) {
	gz, err := gzip.NewReader(in)
	if err != nil {
		return
	}

	defer gz.Close()
	tr := tar.NewReader(gz)
	//fmt.Printf("unpacking to %s\n", dst)
	os.MkdirAll(dst, os.ModeDir|os.ModePerm) // mkdir -p
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		// make the target file
		ndst := filepath.Join(dst, hdr.Name)
		os.MkdirAll(filepath.Dir(ndst), os.ModeDir|os.ModePerm) // mkdir -p
		//fmt.Printf("%s\n", ndst)
		df, err := os.Create(ndst)
		if err != nil {
			break
		}
		io.Copy(df, tr)
		df.Close()
	}
	return
}

//TarFile tar src file into a dst file in the tar writer 
func TarFile(dst, src string, tw *tar.Writer) (err error) {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()
	fi, _ := sf.Stat()
	hdr := new(tar.Header)
	hdr.Size = fi.Size()
	hdr.Name = dst
	hdr.Mode = int64(fi.Mode())
	hdr.ModTime = fi.ModTime()

	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	_, err = io.Copy(tw, sf)
	//log.Printf("%v %d\n", hdr.Name, hdr.Size)
	return
}

//TarBuff copy a buffer content into the dst path in the tar writer
func TarBuff(dst string, src *bytes.Buffer, tw *tar.Writer) (err error) {
	hdr := new(tar.Header)
	hdr.Size = int64(src.Len())
	hdr.Name = dst
	hdr.Mode = int64(os.ModePerm)
	hdr.ModTime = time.Now()
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	_, err = io.Copy(tw, src)
	return
}

//CopyFile copies a single file at once
func CopyFile(dst, src string) (int64, error) {
	sf, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer sf.Close()
	df, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer df.Close()
	return io.Copy(df, sf)
}
