package icon

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/xmlquery"
)

type Fetcher interface {
	Fetch(iconPath, prefix string) error
}

func Download(src, dest string) (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	resp, err := client.Get(src)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	p := filepath.Join(dest, filepath.Base(src))
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(f, resp.Body); err != nil {
		_ = f.Close()
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}
	return p, nil
}

func OptimizeSVG(buf []byte, width, height float64) ([]byte, error) {
	imgdoc, err := xmlquery.Parse(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	s := xmlquery.FindOne(imgdoc, "//svg")
	attrs := []xml.Attr{}
	hasSize := false
	for _, a := range s.Attr {
		switch {
		case a.Name.Local == "width":
			hasSize = true
			a.Value = fmt.Sprintf("%spx", strconv.FormatFloat(width, 'f', 2, 64))
		case a.Name.Local == "height":
			hasSize = true
			a.Value = fmt.Sprintf("%spx", strconv.FormatFloat(height, 'f', 2, 64))
		}
		attrs = append(attrs, a)
	}
	if hasSize {
		s.Attr = attrs
	} else {
		s.Attr = append([]xml.Attr{
			xml.Attr{
				Name:  xml.Name{Local: "width"},
				Value: fmt.Sprintf("%spx", strconv.FormatFloat(width, 'f', 2, 64)),
			},
			xml.Attr{
				Name:  xml.Name{Local: "height"},
				Value: fmt.Sprintf("%spx", strconv.FormatFloat(height, 'f', 2, 64)),
			},
		}, attrs...)
	}

	// If there are no line breaks, Graphviz will not recognize it as SVG.
	docstr := strings.Replace(strings.Replace(imgdoc.OutputXML(false), "?>", "?>\n", 1), "-->", "-->\n", 1)
	if strings.Contains(docstr, "<?xml?>") {
		docstr = strings.Replace(docstr, "<?xml?>", `<?xml version="1.0"?>`, 1)
	}

	return []byte(docstr), nil
}
