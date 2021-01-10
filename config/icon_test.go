package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestComponentIcon(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "ndiag")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	cfg := New()
	if err := cfg.LoadConfigFile(filepath.Join(testdataDir(t), "4_ndiag.yml")); err != nil {
		t.Fatal(err)
	}
	cfg.DescPath = tempDir
	if err := cfg.Build(); err != nil {
		t.Fatal(err)
	}
	for _, c := range cfg.GlobalComponents() {
		if c.Metadata.Icon == "" {
			t.Errorf("icon does not set: %s", c.Id())
		}
	}
	for _, c := range cfg.ClusterComponents() {
		if c.Metadata.Icon == "" {
			t.Errorf("icon does not set: %s", c.Id())
		}
	}
	for _, c := range cfg.NodeComponents() {
		if c.Metadata.Icon == "" {
			t.Errorf("icon does not set: %s", c.Id())
		}
	}
}

func TestClusterIcon(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "ndiag")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	cfg := New()
	if err := cfg.LoadConfigFile(filepath.Join(testdataDir(t), "7_ndiag.yml")); err != nil {
		t.Fatal(err)
	}
	cfg.DescPath = tempDir
	if err := cfg.Build(); err != nil {
		t.Fatal(err)
	}
	for _, c := range cfg.Clusters() {
		if c.Metadata.Icon == "" {
			t.Errorf("icon does not set: %s", c.Id())
		}
	}
}

func TestIconImage(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "ndiag")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	cfg := New()
	if err := cfg.LoadConfigFile(filepath.Join(testdataDir(t), "5_ndiag.yml")); err != nil {
		t.Fatal(err)
	}
	cfg.DescPath = tempDir
	if err := cfg.Build(); err != nil {
		t.Fatal(err)
	}
	{
		i, err := cfg.IconMap().Get("extra")
		if err != nil {
			t.Fatal(err)
		}
		if i.IsGlyph() {
			t.Fatal(fmt.Errorf("%s should not be glyph.Glyph", i.Path))
		}
	}
	{
		i, err := cfg.IconMap().Get("path-to-extra")
		if err != nil {
			t.Fatal(err)
		}
		if i.IsGlyph() {
			t.Fatal(fmt.Errorf("%s is not glyph.Glyph", i.Path))
		}
	}
	{
		i, err := cfg.IconMap().Get("invalid")
		if err == nil {
			t.Fatal(fmt.Errorf("%s is not image file", i.Path))
		}
	}
}
