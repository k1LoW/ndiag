package aws

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

const archiveURL = "https://d1.awsstatic.com/webteam/architecture-icons/Q32020/AWS-Architecture-Assets-For-Light-and-Dark-BG_20200911.478ff05b80f909792f7853b1a28de8e28eac67f4.zip"
const logoURL = "https://upload.wikimedia.org/wikipedia/commons/9/93/Amazon_Web_Services_Logo.svg"

type AWSIcon struct{}

var rep = strings.NewReplacer("_Light", "", "_48", "", "loT", "iot", "IoT", "iot", "FSx", "fsx", "AMIs", "amis", "_", "-", "&", "and", "VMware", "vmware")
var rep2 = strings.NewReplacer("res-amazon", "res", "res-aws", "res", "arch-aws-", "", "arch-amazon-", "")

func (f *AWSIcon) Fetch(iconPath, prefix string) error {
	dir, err := ioutil.TempDir("", "ndiag-icon-aws")
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
	if err := os.MkdirAll(filepath.Join(iconPath, prefix), 0750); err != nil {
		return err
	}
	counter := map[string]struct{}{}
	for _, f := range r.File {
		if strings.Contains(f.Name, "_Dark") {
			continue
		}
		if strings.Contains(f.Name, "_64") || strings.Contains(f.Name, "_32") || strings.Contains(f.Name, "_16") {
			continue
		}
		if !strings.Contains(f.Name, ".svg") {
			continue
		}
		if f.FileInfo().IsDir() {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}
		fn := rep2.Replace(strcase.KebabCase(rep.Replace(filepath.Base(f.Name))))

		path := filepath.Join(iconPath, prefix, fn)

		b := make([]byte, f.UncompressedSize)
		_, err = io.ReadFull(rc, b)
		if err != nil {
			_ = rc.Close()
			return err
		}

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
	if err := ioutil.WriteFile(path, b, 0600); err != nil {
		return err
	}
	counter[path] = struct{}{}

	_, _ = fmt.Fprintf(os.Stderr, "%d icons fetched\n", len(counter))
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", "Done.")
	return nil
}
