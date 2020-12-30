package config

import (
	"os"
	"path/filepath"
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
		wantNEdgeLen            int
		wantLabelLen            int
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
			wantNEdgeLen:            0,
			wantLabelLen:            0,
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
			wantNEdgeLen:            5,
			wantLabelLen:            2,
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
			wantNEdgeLen:            5,
			wantLabelLen:            2,
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
			wantNEdgeLen:            6,
			wantLabelLen:            1,
		},
	}
	for i, tt := range tests {
		d := New()
		if err := d.LoadConfigFile(filepath.Join(testdataDir(t), tt.configFile)); err != nil {
			t.Fatal(err)
		}
		for _, n := range tt.nodeListFiles {
			if err := d.LoadRealNodesFile(filepath.Join(testdataDir(t), n)); err != nil {
				t.Fatal(err)
			}
		}
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
		if got := len(d.NEdges()); got != tt.wantNEdgeLen {
			t.Errorf("TestLoadConfigAndRealNodes(%d) nedge len got %v\nwant %v", i, got, tt.wantNEdgeLen)
		}
		if got := len(d.Labels()); got != tt.wantLabelLen {
			t.Errorf("TestLoadConfigAndRealNodes(%d) label len got %v\nwant %v", i, got, tt.wantLabelLen)
		}
	}
}

func TestBuildNestedCluster(t *testing.T) {
	tests := []struct {
		configFile        string
		nodeListFiles     []string
		layers            []string
		wantClusterLen    int
		wantGlobalNodeLen int
		wantNEdgeLen      int
	}{
		{"1_ndiag.yml", []string{"1_nodes.yml"}, []string{}, 0, 3, 0},
		{"1_ndiag.yml", []string{"1_nodes.yml"}, []string{"consul"}, 1, 0, 0},
		{"2_ndiag.yml", []string{"2_nodes.yml"}, []string{"consul"}, 1, 0, 2},
		{"2_ndiag.yml", []string{"2_nodes.yml"}, []string{"consul", "group"}, 1, 0, 5},
		{"2_ndiag.yml", []string{"2_nodes.yml"}, []string{"group"}, 1, 2, 5},
	}
	for i, tt := range tests {
		cfg := New()
		if err := cfg.LoadConfigFile(filepath.Join(testdataDir(t), tt.configFile)); err != nil {
			t.Fatal(err)
		}
		for _, n := range tt.nodeListFiles {
			if err := cfg.LoadRealNodesFile(filepath.Join(testdataDir(t), n)); err != nil {
				t.Fatal(err)
			}
		}
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

		gotClusters, gotNodes, gotNEdges, err := cfg.BuildNestedClusters(tt.layers)
		if err != nil {
			t.Fatal(err)
		}
		if got := len(gotClusters); got != tt.wantClusterLen {
			t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, tt.wantClusterLen)
		}
		if got := len(gotNodes); got != tt.wantGlobalNodeLen {
			t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, tt.wantGlobalNodeLen)
		}
		if got := len(gotNEdges); got != tt.wantNEdgeLen {
			t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, tt.wantNEdgeLen)
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
	}
}

func TestComponentIcon(t *testing.T) {
	cfg := New()
	if err := cfg.LoadConfigFile(filepath.Join(testdataDir(t), "4_ndiag.yml")); err != nil {
		t.Fatal(err)
	}
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
