package hashicorp

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/k1LoW/ndiag/config"
	"github.com/k1LoW/ndiag/icon"
	"github.com/stoewer/go-strcase"
)

type HashicorpIcon struct{}

var archiveURLs = []string{
	"https://drive.google.com/uc?export=download&id=1EtZa2tnvRJoSk5kqJtuToS8mLnQudoiM", // Vagant
	"https://drive.google.com/uc?export=download&id=1MgDnPnLnTAUQPGwDskbQns1R7KwkjDot", // Packer
	"https://drive.google.com/uc?export=download&id=18sAMuWBoedXWDDZwfuTzaUUmDrpezr7s", // Terraform
	"https://drive.google.com/uc?export=download&id=1MWU8ODoxNFebk_MXlNlUWs98lqxqEqtN", // Vault
	"https://drive.google.com/uc?export=download&id=1cC7MFfkq3sCj_NK4sBCC8hgUPWbLxval", // Nomad
	"https://drive.google.com/uc?export=download&id=17tlD38R-KQmAN1NnQ0RL6ROeXOEwiWRx", // Consul
	"https://drive.google.com/uc?export=download&id=1GA8kNh_8QYNHRP-kdMvbWH_AqVY2hDV7", // Waypoint
}

var rep = strings.NewReplacer("_FullColor_RGB", "", "_VerticalLogo", "", "_PrimaryLogo", "-h")

func (f *HashicorpIcon) Fetch(iconPath, prefix string) error {
	dir, err := ioutil.TempDir("", "ndiag-icon-hashicorp")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	counter := map[string]struct{}{}
	for _, u := range archiveURLs {
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
			if f.FileInfo().IsDir() {
				continue
			}
			if !strings.Contains(f.Name, ".svg") {
				continue
			}
			if !strings.Contains(f.Name, "FullColor") || strings.Contains(f.Name, "Enterprise") {
				continue
			}

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

			path := filepath.Join(iconPath, prefix, strcase.KebabCase(rep.Replace(filepath.Base(f.Name))))

			b, err = icon.OptimizeSVG(b, config.IconWidth, config.IconHeight)
			if err != nil {
				return err
			}

			if err := ioutil.WriteFile(path, b, f.Mode()); err != nil {
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
