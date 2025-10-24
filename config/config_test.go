package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigAndRealNodes(t *testing.T) {
	tests := []struct {
		desc                    string
		configFile              string
		nodeListFiles           []string
		wantNodeLen             int
		wantRealNodeLen         int
		wantClusterLen          int
		wantGlobalComponentLen  int
		wantClusterComponentLen int
		wantNodeComponentLen    int
		wantEdgeLen             int
		wantLabelLen            int
		wantElementLen          int
	}{
		{
			configFile:              "1_ndiag.yml",
			nodeListFiles:           []string{"1_nodes.yml"},
			wantNodeLen:             3,
			wantRealNodeLen:         7,
			wantClusterLen:          1,
			wantGlobalComponentLen:  0,
			wantClusterComponentLen: 0,
			wantNodeComponentLen:    3,
			wantEdgeLen:             0,
			wantLabelLen:            0,
			wantElementLen:          9,
		},
		{
			configFile:              "2_ndiag.yml",
			nodeListFiles:           []string{"2_nodes.yml"},
			wantNodeLen:             3,
			wantRealNodeLen:         7,
			wantClusterLen:          2,
			wantGlobalComponentLen:  1,
			wantClusterComponentLen: 1,
			wantNodeComponentLen:    4,
			wantEdgeLen:             5,
			wantLabelLen:            1,
			wantElementLen:          20,
		},
		{
			configFile:              "3_ndiag.yml",
			nodeListFiles:           []string{"3_nodes.yml"},
			wantNodeLen:             3,
			wantRealNodeLen:         7,
			wantClusterLen:          2,
			wantGlobalComponentLen:  1,
			wantClusterComponentLen: 1,
			wantNodeComponentLen:    4,
			wantEdgeLen:             5,
			wantLabelLen:            1,
			wantElementLen:          20,
		},
		{
			desc:                    "Determine that it is the same component with or without query parameters.",
			configFile:              "4_ndiag.yml",
			nodeListFiles:           []string{},
			wantNodeLen:             3,
			wantRealNodeLen:         0,
			wantClusterLen:          2,
			wantGlobalComponentLen:  1,
			wantClusterComponentLen: 1,
			wantNodeComponentLen:    4,
			wantEdgeLen:             6,
			wantLabelLen:            1,
			wantElementLen:          21,
		},
		{
			desc:                    "No nodes",
			configFile:              "6_ndiag.yml",
			nodeListFiles:           []string{},
			wantNodeLen:             0,
			wantRealNodeLen:         0,
			wantClusterLen:          1,
			wantGlobalComponentLen:  1,
			wantClusterComponentLen: 1,
			wantNodeComponentLen:    0,
			wantEdgeLen:             1,
			wantLabelLen:            0,
			wantElementLen:          6,
		},
		{
			desc:                    "Labels",
			configFile:              "8_ndiag.yml",
			nodeListFiles:           []string{},
			wantNodeLen:             1,
			wantRealNodeLen:         0,
			wantClusterLen:          1,
			wantGlobalComponentLen:  1,
			wantClusterComponentLen: 1,
			wantNodeComponentLen:    1,
			wantEdgeLen:             2,
			wantLabelLen:            5,
			wantElementLen:          13,
		},
	}
	for i, tt := range tests {
		func() {
			tempDir, err := os.MkdirTemp("", "ndiag")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tempDir)
			d := New()
			if err := d.LoadConfigFile(filepath.Join(testdataDir(t), tt.configFile)); err != nil {
				t.Fatal(err)
			}
			for _, n := range tt.nodeListFiles {
				if err := d.LoadRealNodesFile(filepath.Join(testdataDir(t), n)); err != nil {
					t.Fatal(err)
				}
			}
			d.DescPath = tempDir
			if err := d.Build(); err != nil {
				t.Fatal(err)
			}
			if got := len(d.Nodes); got != tt.wantNodeLen {
				t.Errorf("TestLoadConfigAndRealNodes(%d) node len got %v\nwant %v", i, got, tt.wantNodeLen)
			}
			if got := len(d.realNodes); got != tt.wantRealNodeLen {
				t.Errorf("TestLoadConfigAndRealNodes(%d) real node len got %v\nwant %v", i, got, tt.wantRealNodeLen)
			}
			if got := len(d.Clusters()); got != tt.wantClusterLen {
				t.Errorf("TestLoadConfigAndRealNodes(%d) cluster len got %v\nwant %v", i, got, tt.wantClusterLen)
			}
			if got := len(d.GlobalComponents()); got != tt.wantGlobalComponentLen {
				t.Errorf("TestLoadConfigAndRealNodes(%d) global component len got %v\nwant %v", i, got, tt.wantGlobalComponentLen)
			}
			if got := len(d.ClusterComponents()); got != tt.wantClusterComponentLen {
				t.Errorf("TestLoadConfigAndRealNodes(%d) cluster component len got %v\nwant %v", i, got, tt.wantClusterComponentLen)
			}
			if got := len(d.NodeComponents()); got != tt.wantNodeComponentLen {
				t.Errorf("TestLoadConfigAndRealNodes(%d) node component len got %v\nwant %v", i, got, tt.wantNodeComponentLen)
			}
			if got := len(d.Edges()); got != tt.wantEdgeLen {
				t.Errorf("TestLoadConfigAndRealNodes(%d) nedge len got %v\nwant %v", i, got, tt.wantEdgeLen)
			}
			if got := len(d.Labels()); got != tt.wantLabelLen {
				t.Errorf("TestLoadConfigAndRealNodes(%d) label len got %v\nwant %v", i, got, tt.wantLabelLen)
			}
			if got := len(d.Elements()); got != tt.wantElementLen {
				t.Errorf("TestLoadConfigAndRealNodes(%d) ndiag element len got %v\nwant %v", i, got, tt.wantElementLen)
			}
		}()
	}
}

func TestBuildNestedCluster(t *testing.T) {
	tests := []struct {
		configFile        string
		nodeListFiles     []string
		layers            []string
		wantClusterLen    int
		wantGlobalNodeLen int
		wantEdgeLen       int
	}{
		{"1_ndiag.yml", []string{"1_nodes.yml"}, []string{}, 0, 3, 0},
		{"1_ndiag.yml", []string{"1_nodes.yml"}, []string{"consul"}, 1, 0, 0},
		{"2_ndiag.yml", []string{"2_nodes.yml"}, []string{"consul"}, 1, 0, 2},
		{"2_ndiag.yml", []string{"2_nodes.yml"}, []string{"consul", "group"}, 1, 0, 5},
		{"2_ndiag.yml", []string{"2_nodes.yml"}, []string{"group"}, 1, 2, 5},
	}
	for i, tt := range tests {
		func() {
			tempDir, err := os.MkdirTemp("", "ndiag")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tempDir)
			cfg := New()
			if err := cfg.LoadConfigFile(filepath.Join(testdataDir(t), tt.configFile)); err != nil {
				t.Fatal(err)
			}
			for _, n := range tt.nodeListFiles {
				if err := cfg.LoadRealNodesFile(filepath.Join(testdataDir(t), n)); err != nil {
					t.Fatal(err)
				}
			}
			cfg.DescPath = tempDir
			if err := cfg.Build(); err != nil {
				t.Fatal(err)
			}

			cNodeLen := len(cfg.Nodes)
			cRealNodeLen := len(cfg.realNodes)
			cClusterLen := len(cfg.Clusters())
			cGlobalComponentLen := len(cfg.GlobalComponents())
			cClusterComponentLen := len(cfg.ClusterComponents())
			cNodeComponentLen := len(cfg.NodeComponents())
			cRelationLen := len(cfg.Relations)

			gotClusters, gotNodes, gotEdges, err := cfg.BuildNestedClusters(tt.layers)
			if err != nil {
				t.Fatal(err)
			}
			if got := len(gotClusters); got != tt.wantClusterLen {
				t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, tt.wantClusterLen)
			}
			if got := len(gotNodes); got != tt.wantGlobalNodeLen {
				t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, tt.wantGlobalNodeLen)
			}
			if got := len(gotEdges); got != tt.wantEdgeLen {
				t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, tt.wantEdgeLen)
			}

			if got := len(cfg.Nodes); got != cNodeLen {
				t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, cNodeLen)
			}
			if got := len(cfg.realNodes); got != cRealNodeLen {
				t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, cRealNodeLen)
			}
			if got := len(cfg.Clusters()); got != cClusterLen {
				t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, cClusterLen)
			}
			if got := len(cfg.GlobalComponents()); got != cGlobalComponentLen {
				t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, cGlobalComponentLen)
			}
			if got := len(cfg.ClusterComponents()); got != cClusterComponentLen {
				t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, cClusterComponentLen)
			}
			if got := len(cfg.ClusterComponents()); got != cClusterComponentLen {
				t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, cClusterComponentLen)
			}
			if got := len(cfg.NodeComponents()); got != cNodeComponentLen {
				t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, cNodeComponentLen)
			}
			if got := len(cfg.Relations); got != cRelationLen {
				t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, cRelationLen)
			}
		}()
	}
}

func TestHideDetails(t *testing.T) {
	func() {
		tempDir, err := os.MkdirTemp("", "ndiag")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tempDir)
		cfg := New()
		if err := cfg.LoadConfigFile(filepath.Join(testdataDir(t), "..", "example", "3-tier", "input", "ndiag.yml")); err != nil {
			t.Fatal(err)
		}
		cfg.DocPath = tempDir
		if err := cfg.Build(); err != nil {
			t.Fatal(err)
		}
		want := len(cfg.Components())
		if err := cfg.HideDetails(); err != nil {
			t.Fatal(err)
		}
		if got := len(cfg.Components()); got != want {
			t.Errorf("got %v\nwant %v", got, want)
		}
		for _, c := range cfg.Components() {
			if !strings.Contains(c.Name, "component") {
				t.Errorf("got %v\nwant %v", c.Name, "component*")
			}
		}
	}()
}

func testdataDir(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir, err := filepath.Abs(filepath.Join(filepath.Dir(wd), "testdata"))
	if err != nil {
		t.Fatal(err)
	}
	return dir
}
