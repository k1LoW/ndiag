package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/elliotchance/orderedmap"
	"github.com/goccy/go-yaml"
)

const Sep = ":"
const Esc = "\\"

var escRep = strings.NewReplacer(fmt.Sprintf("%s%s", Esc, Sep), "__NDIAG_REP__")
var unescRep = strings.NewReplacer("__NDIAG_REP__", fmt.Sprintf("%s%s", Esc, Sep))

const DefaultDocPath = "archdoc"

var DefaultConfigFilePaths = []string{"ndiag.yml"}
var DefaultDescPath = "ndiag.descriptions"

// DefaultDiagFormat is the default diagram format
const DefaultDiagFormat = "svg"

type NNode interface {
	Id() string
	FullName() string
}

type Attr struct {
	Key   string
	Value string
}

type NEdge struct {
	Src      *Component
	Dst      *Component
	Desc     string
	Relation *Relation
	Attrs    []*Attr
}

type Layer struct {
	Name string
	Desc string
}

type Config struct {
	Name              string      `yaml:"name"`
	Desc              string      `yaml:"desc,omitempty"`
	DocPath           string      `yaml:"docPath"`
	DescPath          string      `yaml:"descPath"`
	Diagrams          []*Diagram  `yaml:"diagrams"`
	Nodes             []*Node     `yaml:"nodes"`
	Relations         []*Relation `yaml:"relations"`
	rawRelations      []*rawRelation
	realNodes         []*RealNode
	layers            []*Layer
	clusters          Clusters
	globalComponents  []*Component
	clusterComponents []*Component
	nodeComponents    []*Component
	nEdges            []*NEdge
	tags              []*Tag
}

func New() *Config {
	return &Config{}
}

func (cfg *Config) DiagFormat() string {
	// TODO: jpg png
	return DefaultDiagFormat
}

func (cfg *Config) PrimaryDiagram() *Diagram {
	return cfg.Diagrams[0]
}

func (cfg *Config) Layers() []*Layer {
	return cfg.layers
}

func (cfg *Config) Clusters() Clusters {
	return cfg.clusters
}

func (cfg *Config) GlobalComponents() []*Component {
	return cfg.globalComponents
}

func (cfg *Config) ClusterComponents() []*Component {
	return cfg.clusterComponents
}

func (cfg *Config) NodeComponents() []*Component {
	return cfg.nodeComponents
}

func (cfg *Config) NEdges() []*NEdge {
	return cfg.nEdges
}

func (cfg *Config) Tags() []*Tag {
	return cfg.tags
}

func (cfg *Config) BuildNestedClusters(layers []string) (Clusters, []*Node, []*NEdge, error) {
	nEdges := []*NEdge{}
	if len(layers) == 0 {
		return Clusters{}, cfg.Nodes, cfg.nEdges, nil
	}
	clusters, globalNodes, err := buildNestedClusters(cfg.Clusters(), layers, cfg.Nodes)
	if err != nil {
		return clusters, globalNodes, nil, err
	}

	for _, e := range cfg.nEdges {
		hBelongTo := false
		tBelongTo := false
		for _, l := range layers {
			if e.Src.Cluster == nil || strings.EqualFold(e.Src.Cluster.Layer, l) {
				hBelongTo = true
			}
			if e.Dst.Cluster == nil || strings.EqualFold(e.Dst.Cluster.Layer, l) {
				tBelongTo = true
			}
		}
		if hBelongTo && tBelongTo {
			nEdges = append(nEdges, e)
		}
	}

	return clusters, globalNodes, nEdges, nil
}

func (cfg *Config) LoadConfig(in []byte) error {
	if err := yaml.Unmarshal(in, cfg); err != nil {
		return err
	}
	return nil
}

func (cfg *Config) LoadConfigFile(path string) error {
	buf, err := loadFile(path)
	if err != nil {
		return err
	}
	return cfg.LoadConfig(buf)
}

func (cfg *Config) LoadRealNodes(in []byte) error {
	if len(cfg.Nodes) == 0 {
		return errors.New("nodes not found")
	}
	if err := cfg.loadRealNodes(in); err != nil {
		return err
	}
	return nil
}

func (cfg *Config) Build() error {
	for _, n := range cfg.Nodes {
		if len(n.RealNodes) == 0 {
			return fmt.Errorf("'%s' does not have any real nodes", n.FullName())
		}
	}
	if err := cfg.buildClusters(); err != nil {
		return err
	}
	if err := cfg.buildComponents(); err != nil {
		return err
	}
	if err := cfg.buildRelations(); err != nil {
		return err
	}
	if err := cfg.checkUnique(); err != nil {
		return err
	}
	if len(cfg.Diagrams) == 0 {
		cfg.Diagrams = append(cfg.Diagrams, &Diagram{
			Name:   "Nodes",
			Layers: []string{},
		})
	}
	// set default
	if cfg.DocPath == "" {
		cfg.DocPath = DefaultDocPath
	}
	if cfg.DescPath == "" {
		cfg.DescPath = DefaultDescPath
	}
	if err := cfg.buildDescriptions(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) loadRealNodes(in []byte) error {
	rNodes := []string{}
	if err := yaml.Unmarshal(in, &rNodes); err == nil {
		for _, rn := range rNodes {
			belongTo := false
			newRn := &RealNode{
				Node: Node{
					Name: rn,
				},
			}
		N:
			for _, n := range cfg.Nodes {
				if n.nameRe.MatchString(rn) {
					belongTo = true
					newRn.BelongTo = n
					n.RealNodes = append(n.RealNodes, newRn)
					break N
				}
			}
			if !belongTo {
				return fmt.Errorf("there is a real node '%s' that does not belong to a node", newRn.Name)
			}
			cfg.realNodes = append(cfg.realNodes, newRn)
		}
	} else {
		// config format
		rConfig := New()
		if err := yaml.Unmarshal(in, rConfig); err != nil {
			return err
		}
		for _, rn := range rConfig.Nodes {
			belongTo := false
			newRn := &RealNode{
				Node: *rn,
			}
		NN:
			for _, n := range cfg.Nodes {
				if n.nameRe.MatchString(rn.Name) {
					belongTo = true
					newRn.BelongTo = n
					n.RealNodes = append(n.RealNodes, newRn)
					break NN
				}
			}
			if !belongTo {
				return fmt.Errorf("there is a real node '%s' that does not belong to a node", newRn.Name)
			}
			cfg.realNodes = append(cfg.realNodes, newRn)
		}
	}
	return nil
}

func (cfg *Config) LoadRealNodesFile(path string) error {
	buf, err := loadFile(path)
	if err != nil {
		return err
	}
	return cfg.LoadRealNodes(buf)
}

func (cfg *Config) FindNode(name string) (*Node, error) {
	for _, n := range cfg.Nodes {
		if strings.EqualFold(n.FullName(), name) {
			return n, nil
		}
	}
	return nil, fmt.Errorf("node not found: %s", name)
}

func (cfg *Config) FindComponent(name string) (*Component, error) {
	var components []*Component
	switch sepCount(name) {
	case 2: // cluster components
		components = cfg.clusterComponents
	case 1: // node components
		components = cfg.nodeComponents
	case 0: // global components
		components = cfg.globalComponents
	}
	for _, c := range components {
		if strings.EqualFold(c.FullName(), name) {
			return c, nil
		}
	}
	return nil, fmt.Errorf("component not found: %s", name)
}

func (cfg *Config) buildComponents() error {
	gc := orderedmap.NewOrderedMap()
	nc := orderedmap.NewOrderedMap()
	cc := orderedmap.NewOrderedMap()
	for _, rel := range cfg.rawRelations {
		for _, r := range rel.Components {
			switch sepCount(r) {
			case 2: // cluster components
				cc.Set(r, struct{}{})
			case 1: // node components
				nc.Set(r, struct{}{})
			case 0: // global components
				gc.Set(r, struct{}{})
			}
		}
	}

	// global components
	for _, c := range gc.Keys() {
		// create global component from relations
		cfg.globalComponents = append(cfg.globalComponents, &Component{
			Name: c.(string),
		})
	}

	// node components
	for _, n := range cfg.Nodes {
		cfg.nodeComponents = append(cfg.nodeComponents, n.Components...)
	}

	for _, c := range nc.Keys() {
		belongTo := false
		splitted := sepSplit(c.(string))
		nodeName := splitted[0]
		comName := splitted[1]
		n, err := cfg.FindNode(nodeName)
		if err != nil {
			return fmt.Errorf("node '%s' not found: %s", nodeName, c)
		}
		for _, com := range n.Components {
			if strings.EqualFold(com.FullName(), c.(string)) {
				belongTo = true
				break
			}
		}
		if !belongTo {
			// create node component from relations
			component := &Component{
				Name: comName,
				Node: n,
			}
			n.Components = append(n.Components, component)
			cfg.nodeComponents = append(cfg.nodeComponents, component)
		}
	}

	// cluster components
	for _, c := range cc.Keys() {
		splitted := sepSplit(c.(string))
		clName := fmt.Sprintf("%s:%s", splitted[0], splitted[1])
		comName := splitted[2]
		belongTo := false
		for _, cl := range cfg.Clusters() {
			if strings.EqualFold(cl.FullName(), clName) {
				// create cluster component from relations
				com := &Component{
					Cluster: cl,
					Name:    comName,
				}
				cl.Components = append(cl.Components, com)
				cfg.clusterComponents = append(cfg.clusterComponents, com)
				belongTo = true
				break
			}
		}
		if !belongTo {
			return fmt.Errorf("cluster '%s' not found: %s", clName, c)
		}
	}
	return nil
}

func (cfg *Config) buildClusters() error {
	for _, n := range cfg.Nodes {
		for _, c := range n.rawClusters {
			cluster, err := cfg.parseClusterLabel(c)
			if err != nil {
				return err
			}
			cluster.Nodes = append(cluster.Nodes, n)
			n.Clusters = append(n.Clusters, cluster)
		}
	}
	return nil
}

func (cfg *Config) parseClusterLabel(label string) (*Cluster, error) {
	if !strings.Contains(label, Sep) {
		return nil, fmt.Errorf("invalid cluster id: %s", label)
	}
	splitted := sepSplit(label)
	if len(splitted) != 2 {
		return nil, fmt.Errorf("invalid cluster id: %s", label)
	}
	layer := splitted[0]
	name := splitted[1]
	current := cfg.clusters.Find(layer, name)
	if current != nil {
		return current, nil
	}
	newC := &Cluster{
		Layer: layer,
		Name:  name,
	}
	cfg.clusters = append(cfg.clusters, newC)
	if !layerContains(cfg.layers, layer) {
		cfg.layers = append(cfg.layers, &Layer{Name: layer})
	}
	return newC, nil
}

func (cfg *Config) buildRelations() error {
	relTags := orderedmap.NewOrderedMap()
	for _, rel := range cfg.rawRelations {
		nrel := &Relation{
			RelationId: rel.Id,
			Type:       rel.Type,
			Tags:       rel.Tags,
			Attrs:      rel.Attrs,
		}
		for _, r := range rel.Components {
			c, err := cfg.FindComponent(r)
			if err != nil {
				return err
			}
			nrel.Components = append(nrel.Components, c)
		}
		cfg.Relations = append(cfg.Relations, nrel)

		// tags
		for _, t := range rel.Tags {
			var nt *Tag
			nti, ok := relTags.Get(t)
			if ok {
				nt = nti.(*Tag)
			} else {
				nt = &Tag{
					Name: t,
				}
				relTags.Set(t, nt)
			}
			nt.Relations = append(nt.Relations, nrel)
		}
	}
	cfg.nEdges = SplitRelations(cfg.Relations)

	for _, k := range relTags.Keys() {
		nt, _ := relTags.Get(k)
		cfg.tags = append(cfg.tags, nt.(*Tag))
	}

	return nil
}

func (cfg *Config) buildDescriptions() error {
	if cfg.DescPath == "" {
		return nil
	}
	err := os.MkdirAll(cfg.DescPath, 0755) // #nosec
	if err != nil {
		return err
	}

	// diagrams
	for _, d := range cfg.Diagrams {
		if d.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MdPath("_diagram", d.Layers))
		if err != nil {
			return err
		}
		d.Desc = desc
	}

	// clusters
	for _, c := range cfg.Clusters() {
		if c.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MdPath("_cluster", []string{c.Id()}))
		if err != nil {
			return err
		}
		c.Desc = desc
	}

	// layers
	for _, l := range cfg.Layers() {
		if l.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MdPath("_layer", []string{l.Name}))
		if err != nil {
			return err
		}
		l.Desc = desc
	}

	// nodes
	for _, n := range cfg.Nodes {
		if n.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MdPath("_node", []string{n.Id()}))
		if err != nil {
			return err
		}
		n.Desc = desc
	}

	// components
	for _, c := range cfg.GlobalComponents() {
		if c.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MdPath("_component", []string{c.Id()}))
		if err != nil {
			return err
		}
		c.Desc = desc
	}
	for _, c := range cfg.ClusterComponents() {
		if c.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MdPath("_component", []string{c.Id()}))
		if err != nil {
			return err
		}
		c.Desc = desc
	}
	for _, c := range cfg.NodeComponents() {
		if c.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MdPath("_component", []string{c.Id()}))
		if err != nil {
			return err
		}
		c.Desc = desc
	}

	// tags
	for _, t := range cfg.tags {
		if t.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MdPath("_tag", []string{t.Id()}))
		if err != nil {
			return err
		}
		t.Desc = desc
	}

	return nil
}

func (cfg *Config) readDescFile(f string) (string, error) {
	descPath := filepath.Join(cfg.DescPath, f)
	file, err := os.OpenFile(descPath, os.O_RDONLY|os.O_CREATE, 0644) // #nosec
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	if err := file.Close(); err != nil {
		return "", err
	}
	return string(b), err
}

func buildNestedClusters(clusters Clusters, layers []string, nodes []*Node) (Clusters, []*Node, error) {
	if len(layers) == 0 {
		return clusters, nodes, nil
	}
	leaf := layers[len(layers)-1]
	layers = layers[:len(layers)-1]

	remain := []*Node{}
	belongTo := []*Node{}
	for _, n := range nodes {
		c := n.Clusters.FindByLayer(leaf)
		if len(c) == 0 {
			remain = append(remain, n)
			continue
		}
		if len(c) > 1 {
			return nil, nil, fmt.Errorf("duplicate layer %s", leaf)
		}
		belongTo = append(belongTo, n)
		if len(layers) == 0 {
			continue
		}

		// build cluster tree using Node.Clusters
		parent := ""
		var pc Clusters
		for i := 1; i <= len(layers); i++ {
			parent = layers[len(layers)-i]
			pc = n.Clusters.FindByLayer(parent)
			if len(pc) > 1 {
				return nil, nil, fmt.Errorf("duplicate layer %s", parent)
			}
			if len(pc) == 0 {
				continue
			}
			// _, _ = fmt.Fprintf(os.Stderr, "build cluster tree %v->%v\n", pc[0].FullName(), c[0].FullName())
			c[0].Parent = pc[0]
			pc[0].Children = append(pc[0].Children, c[0])
			c = pc
		}
	}

	// build a direct member node of a cluster
	for _, c := range clusters {
		if c.Layer == leaf {
			continue
		}
		nodes := []*Node{}
	N:
		for _, n := range c.Nodes {
			for _, b := range belongTo {
				if n.FullName() == b.FullName() {
					continue N
				}
			}
			nodes = append(nodes, n)
		}
		c.Nodes = nodes
	}

	// build root clusters
	if len(layers) == 0 {
		root := Clusters{}
	NN:
		for _, c := range clusters {
			if c.Parent == nil && (c.Layer == leaf || len(c.Nodes) > 0) {
				for _, n := range c.Nodes {
					for _, rn := range remain {
						if n == rn {
							continue NN
						}
					}
				}
				root = append(root, c)
			}
		}
		clusters = root
	}

	return buildNestedClusters(clusters, layers, remain)
}

func (cfg *Config) checkUnique() error {

	ids := map[string]string{}

	// nodes
	for _, n := range cfg.Nodes {
		if t, exist := ids[n.Id()]; exist {
			return fmt.Errorf("duplicate id: %s[%s] <-> %s[%s]", t, n.Id(), "node", n.Id())
		}
		ids[n.Id()] = "node"
	}

	// components
	for _, c := range cfg.GlobalComponents() {
		if t, exist := ids[c.Id()]; exist {
			return fmt.Errorf("duplicate id: %s[%s] <-> %s[%s]", t, c.Id(), "component", c.Id())
		}
		ids[c.Id()] = "component"
	}
	for _, c := range cfg.ClusterComponents() {
		if t, exist := ids[c.Id()]; exist {
			return fmt.Errorf("duplicate id: %s[%s] <-> %s[%s]", t, c.Id(), "component", c.Id())
		}
		ids[c.Id()] = "component"
	}
	for _, c := range cfg.NodeComponents() {
		if t, exist := ids[c.Id()]; exist {
			return fmt.Errorf("duplicate id: %s[%s] <-> %s[%s]", t, c.Id(), "component", c.Id())
		}
		ids[c.Id()] = "component"
	}

	// clusters
	for _, c := range cfg.Clusters() {
		if t, exist := ids[c.Id()]; exist {
			return fmt.Errorf("duplicate id: %s[%s] <-> %s[%s]", t, c.Id(), "cluster", c.Id())
		}
		ids[c.Id()] = "cluster"
	}

	// read nodes
	m := map[string]struct{}{}
	for _, rn := range cfg.realNodes {
		if _, exist := m[rn.Name]; exist {
			return fmt.Errorf("duplicate real node name: %s", rn.Name)
		}
		m[rn.Name] = struct{}{}
	}

	return nil
}

func loadFile(path string) ([]byte, error) {
	if path == "" {
		return nil, nil
	}
	fullPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadFile(filepath.Clean(fullPath))
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func sepCount(s string) int {
	return strings.Count(escRep.Replace(s), Sep)
}

func sepSplit(s string) []string {
	splitted := strings.Split(escRep.Replace(s), Sep)
	unescaped := []string{}
	for _, ss := range splitted {
		unescaped = append(unescaped, unescRep.Replace(ss))
	}
	return unescaped
}

func sepContains(s string) bool {
	return strings.Contains(escRep.Replace(s), Sep)
}

func layerContains(s []*Layer, e string) bool {
	for _, v := range s {
		if e == v.Name {
			return true
		}
	}
	return false
}
