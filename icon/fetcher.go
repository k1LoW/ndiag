package icon

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"image/png"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/xmlquery"
	"github.com/nfnt/resize"
)

var sizeRe = regexp.MustCompile(`\A([0-9.]+)`)

type Fetcher interface {
	Fetch(iconPath, prefix string) error
}

func Download(src, dest string) (string, error) {
	client := &http.Client{
		Timeout: 60 * time.Second,
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
	f, err := os.OpenFile(filepath.Clean(p), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
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

func OptimizeSVG(b []byte, width, height float64) ([]byte, error) {
	size := height // At the moment we only use height out of width and height.

	imgdoc, err := xmlquery.Parse(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	s := xmlquery.FindOne(imgdoc, "//svg")
	attrs := []xml.Attr{}
	hasViewBox := false
	cw := 0.0
	ch := 0.0
	nw := size
	nh := size
	for _, a := range s.Attr {
		switch {
		case a.Name.Local == "width":
			if !hasViewBox {
				matched := sizeRe.FindStringSubmatch(a.Value)
				if len(matched) > 0 {
					cw, _ = strconv.ParseFloat(matched[1], 64)
				}
			}
		case a.Name.Local == "height":
			if !hasViewBox {
				matched := sizeRe.FindStringSubmatch(a.Value)
				if len(matched) > 0 {
					ch, _ = strconv.ParseFloat(matched[1], 64)
				}
			}
		case a.Name.Local == "viewBox":
			splitted := strings.Split(a.Value, " ")
			if len(splitted) == 4 {
				hasViewBox = true
				cw, _ = strconv.ParseFloat(splitted[2], 64)
				ch, _ = strconv.ParseFloat(splitted[3], 64)
			}
			attrs = append(attrs, a)
		default:
			attrs = append(attrs, a)
		}
	}
	if cw > 0 && ch > 0 {
		if cw > ch {
			// Extend the size horizontally only if width > height
			nw = size * (cw / ch)
		}
	}

	s.Attr = append([]xml.Attr{
		xml.Attr{
			Name:  xml.Name{Local: "width"},
			Value: fmt.Sprintf("%spx", strconv.FormatFloat(nw, 'f', 2, 64)),
		},
		xml.Attr{
			Name:  xml.Name{Local: "height"},
			Value: fmt.Sprintf("%spx", strconv.FormatFloat(nh, 'f', 2, 64)),
		},
	}, attrs...)

	if !hasViewBox {
		s.Attr = append(s.Attr, xml.Attr{
			Name:  xml.Name{Local: "viewBox"},
			Value: fmt.Sprintf("0 0 %g %g", cw, ch),
		})
	}

	// If there are no line breaks, Graphviz will not recognize it as SVG.
	docstr := strings.Replace(strings.Replace(imgdoc.OutputXML(false), "?>", "?>\n", 1), "-->", "-->\n", 1)
	if strings.Contains(docstr, "<?xml?>") {
		docstr = strings.Replace(docstr, "<?xml?>", `<?xml version="1.0"?>`, 1)
	}

	return []byte(docstr), nil
}

func ResizePNG(b []byte, width, height float64) ([]byte, error) {
	size := height // At the moment we only use height out of width and height.

	i, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	rct := i.Bounds()
	w := uint(size)
	h := uint(size)
	if rct.Dx() > rct.Dy() {
		// Extend the size horizontally only if width > height
		w = uint(size * (float64(rct.Dx()) / float64(rct.Dy())))
	}
	r := resize.Resize(w, h, i, resize.Bilinear)
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, r); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
