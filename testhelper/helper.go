//nolint:godoclint,wrapcheck // These are for tests
package testhelper

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func extractFile(zf *zip.File, fn string) (e error) {
	var fDst *os.File
	var fSrc io.ReadCloser

	// Open compressed file
	if fSrc, e = zf.Open(); e != nil {
		return e
	}
	defer func() {
		if e2 := fSrc.Close(); e == nil {
			e = e2
		}
	}()

	// Open destination file
	if fDst, e = os.Create(filepath.Clean(fn)); e != nil {
		return e
	}
	defer func() {
		if e2 := fDst.Close(); e == nil {
			e = e2
		}
	}()

	// Write decompressed contents to file
	for {
		_, e = io.CopyN(fDst, fSrc, 1024*1024) //nolint:mnd // 1MB

		switch e {
		case nil:
		case io.EOF:
			return nil
		default:
			return e
		}
	}
}

func extractZip(fn string, dir string) (e error) {
	var fi fs.FileInfo
	var z *zip.ReadCloser

	// Open zip file
	if z, e = zip.OpenReader(fn); e != nil {
		return e
	}
	defer func() {
		if e2 := z.Close(); e == nil {
			e = e2
		}
	}()

	// Loop thru the contained files
	for _, zf := range z.File {
		//nolint:gosec // G305 - I made the zip file and I trust me
		fn = filepath.Clean(filepath.Join(dir, zf.Name))

		if fi = zf.FileInfo(); fi.IsDir() {
			if e = os.MkdirAll(fn, fi.Mode()); e != nil {
				return e
			}
		} else {
			if e = extractFile(zf, fn); e != nil {
				return e
			}
		}
	}

	return nil
}

func Setup(t *testing.T, prefix string) error {
	t.Helper()

	var repo string = filepath.Join("testdata", "dummy")
	var z string = filepath.Join(prefix, "testhelper", "dummy.zip")

	if e := os.RemoveAll(repo); e != nil {
		return e
	}

	if e := extractZip(z, "testdata"); e != nil {
		return e
	}

	return nil
}
