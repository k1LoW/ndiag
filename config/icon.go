package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/k1LoW/glyph"
)

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
	gm := glyph.NewMapWithIncluded(glyph.Width(80.0), glyph.Height(80.0))
	for _, i := range cfg.CustomIcons {
		g, k, err := i.ToGlyphAndKey()
		if err != nil {
			return err
		}
		gm.Set(k, g)
	}
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
	cfg.iconMap = im
	return nil
}
