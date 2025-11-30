package k8s

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/k1LoW/ndiag/config"
	"github.com/k1LoW/ndiag/icon"
)

const archiveURL = "https://github.com/kubernetes/community/archive/master.zip"
const logoURL = "https://raw.githubusercontent.com/kubernetes/kubernetes/master/logo/logo_with_border.svg"

var pathRe = regexp.MustCompile(`\A.+/([^/]+)/([^/]+)/([^/]+)\.svg\z`)
var rep = strings.NewReplacer("control_plane_components", "control-plane", "infrastructure_components", "infra", "_", "-")

type Icon struct{}

func (f *Icon) Fetch(iconPath, prefix string) error {
	dir, err := os.MkdirTemp("", "ndiag-icon-k8s")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	_, _ = fmt.Fprintf(os.Stderr, "Fetching icons from %s ...\n", archiveURL)
	ap, err := icon.Download(archiveURL, dir)
	if err != nil {
		return err
	}
	r, err := zip.OpenReader(ap)
	if err != nil {
		return err
	}
	counter := map[string]struct{}{}
	for _, f := range r.File {
		if !strings.Contains(f.Name, "icons/svg") {
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
		var path string
		if matched[2] == "labeled" {
			path = rep.Replace(filepath.Join(prefix, matched[1], fmt.Sprintf("%s.%s", matched[3], "svg")))
		} else {
			path = rep.Replace(filepath.Join(prefix, matched[1], fmt.Sprintf("%s-%s.%s", matched[3], matched[2], "svg")))
		}
		path = filepath.Join(iconPath, path)
		if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
			return err
		}

		buf, err = icon.OptimizeSVG(buf, config.IconWidth, config.IconHeight)
		if err != nil {
			return err
		}

		if err := os.WriteFile(path, buf, 0600); err != nil {
			_ = rc.Close()
			return err
		}
		counter[path] = struct{}{}
		if err := rc.Close(); err != nil {
			return err
		}
	}

	// logo
	_, _ = fmt.Fprintf(os.Stderr, "Fetching icon from %s ...\n", logoURL)
	lp, err := icon.Download(logoURL, dir)
	if err != nil {
		return err
	}
	b, err := os.ReadFile(filepath.Clean(lp))
	if err != nil {
		return err
	}
	b, err = icon.OptimizeSVG(b, config.IconWidth, config.IconHeight)
	if err != nil {
		return err
	}
	path := filepath.Join(iconPath, prefix, "logo.svg")
	if err := os.WriteFile(path, b, 0600); err != nil {
		return err
	}
	counter[path] = struct{}{}

	_, _ = fmt.Fprintf(os.Stderr, "%d icons fetched\n", len(counter))
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", "Done.")
	return nil
}
