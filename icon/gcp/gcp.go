package gcp

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
	"github.com/stoewer/go-strcase"
)

const archiveURL = "https://cloud.google.com/icons/files/google-cloud-icons.zip"
const iconArchiveURL = "https://cloud.google.com/files/logos/logos-cloud.zip"

var pathRe = regexp.MustCompile(`\A.+/([^/]+)\.(svg)\z`)
var iconRe = regexp.MustCompile(`\Alogos_cloud/([^/]+)\.(png)\z`)
var rep = strings.NewReplacer(
	"-512-color", "", "-521-color", "", "-color", "", " (1)", "", "_", "-",
	"icon_cloud_192pt_clr", "logo",
	"lockup_cloud_main", "logo-lockup-v",
	"lockup_cloud_sm_v", "logo-lockup-v-sm",
	"logo_lockup_cloud_rgb", "logo-lockup-h",
)

type Icon struct{}

func (f *Icon) Fetch(iconPath, prefix string) error {
	dir, err := os.MkdirTemp("", "ndiag-icon-gcp")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	counter := map[string]struct{}{}
	for _, u := range []string{archiveURL, iconArchiveURL} {
		_, _ = fmt.Fprintf(os.Stderr, "Fetching icons from %s ...\n", u)
		ap, err := icon.Download(u, dir)
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
				matched = iconRe.FindStringSubmatch(f.Name)
				if len(matched) == 0 {
					continue
				}
			}
			format := matched[2]
			rc, err := f.Open()
			if err != nil {
				return err
			}
			b := make([]byte, f.UncompressedSize)
			_, err = io.ReadFull(rc, b)
			if err != nil {
				_ = rc.Close()
				return err
			}
			path := filepath.Join(iconPath, prefix, fmt.Sprintf("%s.%s", strcase.KebabCase(rep.Replace(matched[1])), format))

			switch format {
			case "svg":
				b, err = icon.OptimizeSVG(b, config.IconWidth, config.IconHeight)
				if err != nil {
					return err
				}
			case "png":
				b, err = icon.ResizePNG(b, config.IconWidth, config.IconHeight)
				if err != nil {
					return err
				}
			}

			if err := os.WriteFile(path, b, 0644); err != nil {
				_ = rc.Close()
				return err
			}
			counter[path] = struct{}{}
			if err := rc.Close(); err != nil {
				return err
			}
		}
	}
	_, _ = fmt.Fprintf(os.Stderr, "%d icons fetched\n", len(counter))
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", "Done.")
	return nil
}
