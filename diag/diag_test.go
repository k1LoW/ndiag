package diag

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
	for _, tt := range tests {
		d := New()
		if err := d.LoadConfigFile(filepath.Join(testdataDir(t), tt.configFile)); err != nil {
			t.Fatal(err)
		}
		for _, n := range tt.nodeListFiles {
			if err := d.LoadRealNodesFile(filepath.Join(testdataDir(t), n)); err != nil {
				t.Fatal(err)
			}
		}
		if got := len(d.Nodes); got != tt.wantNodeLen {
			t.Errorf("got %v\nwant %v", got, tt.wantNodeLen)
		}
		if got := len(d.realNodes); got != tt.wantRealNodeLen {
			t.Errorf("got %v\nwant %v", got, tt.wantRealNodeLen)
		}
		if got := len(d.Clusters()); got != tt.wantClusterLen {
			t.Errorf("got %v\nwant %v", got, tt.wantClusterLen)
		}
		if got := len(d.GlobalComponents()); got != tt.wantGlobalComponentLen {
			t.Errorf("got %v\nwant %v", got, tt.wantGlobalComponentLen)
		}
		if got := len(d.ClusterComponents()); got != tt.wantClusterComponentLen {
			t.Errorf("got %v\nwant %v", got, tt.wantClusterComponentLen)
		}
		if got := len(d.NodeComponents()); got != tt.wantNodeComponentLen {
			t.Errorf("got %v\nwant %v", got, tt.wantNodeComponentLen)
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
