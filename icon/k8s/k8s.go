package k8s

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
)

const prefix = "k8s"
const archiveURL = "https://github.com/kubernetes/community/archive/master.zip"

var pathRe = regexp.MustCompile(`\A.+/([^/]+)/([^/]+)/([^/]+)\.svg\z`)

type K8sIcon struct{}

func (f *K8sIcon) Fetch(iconPath string) error {
	_, _ = fmt.Fprintf(os.Stderr, "Fetching from %s ...\n", archiveURL)
	dir, err := ioutil.TempDir("", "ndiag-icon-k8s")
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
	for _, f := range r.File {
		if !strings.Contains(f.Name, "icons/svg") {
			continue
		}
		if f.FileInfo().IsDir() {
			continue
		}
		matched := pathRe.FindStringSubmatch(f.Name)

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
		var path string
		if matched[2] == "labeled" {
			path = filepath.Join(iconPath, prefix, matched[1], fmt.Sprintf("%s.%s", matched[3], "svg"))
		} else {
			path = filepath.Join(iconPath, prefix, matched[1], matched[3], fmt.Sprintf("%s.%s", matched[2], "svg"))
		}
		if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
			return err
		}
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
