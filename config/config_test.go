package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigAndRealNodes(t *testing.T) {
	tests := []struct {
		configFile              string
		nodeListFiles           []string
		wantNodeLen             int
		wantRealNodeLen         int
		wantClusterLen          int
		wantGlobalComponentLen  int
		wantClusterComponentLen int
		wantNodeComponentLen    int
	}{
		{"1_ndiag.yml", []string{"1_nodes.yml"}, 3, 7, 1, 0, 0, 3},
		{"2_ndiag.yml", []string{"2_nodes.yml"}, 3, 7, 2, 1, 1, 3},
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
			t.Errorf("TestLoadConfigAndRealNodes(%d) got %v\nwant %v", i, got, tt.wantNodeLen)
		}
		if got := len(d.realNodes); got != tt.wantRealNodeLen {
			t.Errorf("TestLoadConfigAndRealNodes(%d) got %v\nwant %v", i, got, tt.wantRealNodeLen)
		}
		if got := len(d.Clusters()); got != tt.wantClusterLen {
			t.Errorf("TestLoadConfigAndRealNodes(%d) got %v\nwant %v", i, got, tt.wantClusterLen)
		}
		if got := len(d.GlobalComponents()); got != tt.wantGlobalComponentLen {
			t.Errorf("TestLoadConfigAndRealNodes(%d) got %v\nwant %v", i, got, tt.wantGlobalComponentLen)
		}
		if got := len(d.ClusterComponents()); got != tt.wantClusterComponentLen {
			t.Errorf("TestLoadConfigAndRealNodes(%d) got %v\nwant %v", i, got, tt.wantClusterComponentLen)
		}
		if got := len(d.NodeComponents()); got != tt.wantNodeComponentLen {
			t.Errorf("TestLoadConfigAndRealNodes(%d) got %v\nwant %v", i, got, tt.wantNodeComponentLen)
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
		{"2_ndiag.yml", []string{"2_nodes.yml"}, []string{"consul", "group"}, 1, 0, 4},
		{"2_ndiag.yml", []string{"2_nodes.yml"}, []string{"group"}, 1, 2, 4},
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

		cNodeLen := len(d.Nodes)
		cRealNodeLen := len(d.realNodes)
		cClusterLen := len(d.Clusters())
		cGlobalComponentLen := len(d.GlobalComponents())
		cClusterComponentLen := len(d.ClusterComponents())
		cNodeComponentLen := len(d.NodeComponents())
		cRelationLen := len(d.Relations)

		gotClusters, gotNodes, gotNEdges, err := d.BuildNestedClusters(tt.layers)
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

		if got := len(d.Nodes); got != cNodeLen {
			t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, cNodeLen)
		}
		if got := len(d.realNodes); got != cRealNodeLen {
			t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, cRealNodeLen)
		}
		if got := len(d.Clusters()); got != cClusterLen {
			t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, cClusterLen)
		}
		if got := len(d.GlobalComponents()); got != cGlobalComponentLen {
			t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, cGlobalComponentLen)
		}
		if got := len(d.ClusterComponents()); got != cClusterComponentLen {
			t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, cClusterComponentLen)
		}
		if got := len(d.ClusterComponents()); got != cClusterComponentLen {
			t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, cClusterComponentLen)
		}
		if got := len(d.NodeComponents()); got != cNodeComponentLen {
			t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, cNodeComponentLen)
		}
		if got := len(d.Relations); got != cRelationLen {
			t.Errorf("TestBuildNestedCluster(%d) got %v want %v", i, got, cRelationLen)
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
