package gpk

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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

func PackageWalker(srcpath, startwith string, handler func(gpkpath string)) error {

	file, err := os.Open(srcpath)
	if err != nil {
		return err
	}
	subdir, err := file.Readdir(-1)

	if err != nil {
		return err
	}

	for _, fi := range subdir {
		if fi.IsDir() && strings.HasPrefix(fi.Name(), startwith) {
			PackageWalker(filepath.Join(srcpath, fi.Name()), "", handler)
		} else if fi.Name() == GpkFile {
			//then src path is the package/version directory
			handler(srcpath)
			break
		}
	}
	return nil
}

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

//Create a tar.gz using the best level compression
func MakeTarget() (err error) {
	return os.MkdirAll("target", os.ModeDir|os.ModePerm)
}
