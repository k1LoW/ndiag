package config

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMerge(t *testing.T) {
	rootDir := filepath.Join(testdataDir(t), "config_merge")
	dirs, err := os.ReadDir(rootDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}
		in, err := os.ReadDir(filepath.Join(rootDir, d.Name(), "in"))
		if err != nil {
			t.Fatal(err)
		}
		got := New()
		for _, i := range in {
			cfg := New()
			if err := cfg.LoadConfigFile(filepath.Join(rootDir, d.Name(), "in", i.Name())); err != nil {
				t.Fatal(err)
			}
			if err := got.Merge(cfg); err != nil {
				t.Fatal(err)
			}
		}
		want := New()
		if err := want.LoadConfigFile(filepath.Join(rootDir, d.Name(), "want.yml")); err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(got, want, cmp.AllowUnexported(Config{}, Node{}), cmpopts.IgnoreUnexported(regexp.Regexp{}), cmpopts.IgnoreFields(Config{}, "basePath", "Dict")); diff != "" {
			t.Errorf("%s", diff)
		}
	}
}
