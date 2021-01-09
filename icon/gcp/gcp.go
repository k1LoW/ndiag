package gcp

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/k1LoW/ndiag/icon"
	"github.com/stoewer/go-strcase"
)

const archiveURL = "https://cloud.google.com/icons/files/google-cloud-icons.zip"

var pathRe = regexp.MustCompile(`\A.+/([^/]+)\.svg\z`)
var rep = strings.NewReplacer("-512-color", "", "-521-color", "", "-color", "", " (1)", "", "_", "-")

type GCPIcon struct{}

func (f *GCPIcon) Fetch(iconPath, prefix string) error {
	_, _ = fmt.Fprintf(os.Stderr, "Fetching from %s ...\n", archiveURL)
	dir, err := ioutil.TempDir("", "ndiag-icon-gcp")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	ap, err := icon.Download(archiveURL, dir)
	if err != nil {
		return err
	}
	r, err := zip.OpenReader(ap)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(iconPath, prefix), 0750); err != nil {
		return err
	}

	for _, f := range r.File {
		if strings.Contains(f.Name, "__MACOSX") {
			continue
		}
		if strings.Contains(f.Name, "Expanded Product Card Icons") {
			continue
		}
		if f.FileInfo().IsDir() {
			continue
		}
		matched := pathRe.FindStringSubmatch(f.Name)
		if len(matched) == 0 {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		buf := make([]byte, f.UncompressedSize)
		_, err = io.ReadFull(rc, buf)
		if err != nil {
			_ = rc.Close()
			return err
		}
		path := filepath.Join(iconPath, prefix, fmt.Sprintf("%s.%s", strcase.KebabCase(rep.Replace(matched[1])), "svg"))
		if err := ioutil.WriteFile(path, buf, f.Mode()); err != nil {
			_ = rc.Close()
			return err
		}
		if err := rc.Close(); err != nil {
			return err
		}
	}
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", "Done.")
	return nil
}
