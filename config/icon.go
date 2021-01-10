package config

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	issvg "github.com/h2non/go-is-svg"
	"github.com/k1LoW/glyph"
	"github.com/karrick/godirwalk"
)

const IconWidth = 80.0
const IconHeight = 80.0

type Icon struct {
	Path  string
	Glyph *glyph.Glyph
}

func (i *Icon) IsGlyph() bool {
	return (i.Glyph != nil)
}

type IconMap struct {
	tempIconDir string
	icons       map[string]*Icon
}

func NewIconMap(tempIconDir string) *IconMap {
	return &IconMap{
		tempIconDir: tempIconDir,
		icons:       map[string]*Icon{},
	}
}

func (m *IconMap) Keys() []string {
	keys := []string{}
	for k := range m.icons {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (m *IconMap) Get(k string) (*Icon, error) {
	i, ok := m.icons[k]
	if !ok {
		return nil, fmt.Errorf("invalid key: %s", k)
	}
	return i, nil
}

func (m *IconMap) Set(k string, i *Icon) {
	m.icons[k] = i
}

func (m *IconMap) GeneratePNGGlyphIcons() (e error) {
	if err := os.Mkdir(m.tempIconDir, 0750); err != nil {
		return err
	}
	for _, k := range m.Keys() {
		i, err := m.Get(k)
		if err != nil {
			return err
		}
		if !i.IsGlyph() {
			continue
		}
		f, err := os.OpenFile(i.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666) // #nosec
		if err != nil {
			return err
		}
		if err := i.Glyph.WriteImage(f); err != nil {
			e = f.Close()
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (m *IconMap) GenerateSVGGlyphIcons() (e error) {
	if err := os.Mkdir(m.tempIconDir, 0750); err != nil {
		return err
	}
	for _, k := range m.Keys() {
		i, err := m.Get(k)
		if err != nil {
			return err
		}
		if !i.IsGlyph() {
			continue
		}
		f, err := os.OpenFile(i.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666) // #nosec
		if err != nil {
			return err
		}
		if err := i.Glyph.Write(f); err != nil {
			e = f.Close()
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (m *IconMap) RemoveTempIconDir() error {
	return os.RemoveAll(m.tempIconDir)
}

func (cfg *Config) buildIconMap() error {
	tempIconDir := filepath.Join(os.TempDir(), fmt.Sprintf("ndiag.%06d", os.Getpid()))
	im := NewIconMap(tempIconDir)
	gm := glyph.NewMapWithIncluded(glyph.Width(IconWidth), glyph.Height(IconHeight))
	for _, i := range cfg.CustomIcons {
		g, k, err := i.ToGlyphAndKey()
		if err != nil {
			return err
		}
		gm.Set(k, g)
	}
	// set glyph.Glyph
	for _, k := range gm.Keys() {
		g, err := gm.Get(k)
		if err != nil {
			return err
		}
		i := &Icon{
			Path:  filepath.Join(tempIconDir, fmt.Sprintf("%s.%s", k, cfg.Format())),
			Glyph: g,
		}
		im.Set(k, i)
	}
	// set icon images
	depth := 5
	if _, err := os.Lstat(cfg.IconPath); err == nil {
		err := godirwalk.Walk(cfg.IconPath, &godirwalk.Options{
			Callback: func(path string, de *godirwalk.Dirent) error {
				if strings.Contains(path, ".git") {
					return godirwalk.SkipThis
				}
				d, err := de.IsDirOrSymlinkToDir()
				if err != nil {
					return err
				}
				rel, err := filepath.Rel(cfg.IconPath, path)
				if err != nil {
					return err
				}
				if d {
					if strings.Count(filepath.ToSlash(rel), "/") > depth {
						return godirwalk.SkipThis
					}
					return nil
				}
				if !isImg(path) {
					return nil
				}
				k := strings.ReplaceAll(filepath.ToSlash(strings.TrimSuffix(rel, filepath.Ext(rel))), "/", "-")
				i := &Icon{
					Path:  path,
					Glyph: nil,
				}
				im.Set(k, i)
				return nil
			},
			Unsorted: false,
		})
		if err != nil {
			return err
		}
	}
	cfg.iconMap = im
	return nil
}

func isImg(path string) bool {
	imgf, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return false
	}
	if issvg.Is(imgf) {
		return true
	}
	if _, _, err := image.DecodeConfig(bytes.NewReader(imgf)); err != nil {
		return false
	}
	return true
}
