package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/elliotchance/orderedmap"
	"github.com/goccy/go-yaml"
	"github.com/k1LoW/glyph"
	"github.com/k1LoW/tbls/dict"
	"github.com/pasztorpisti/qs"
)

const Sep = ":"
const Esc = "\\"
const Q = "?"

var escRep = strings.NewReplacer(fmt.Sprintf("%s%s", Esc, Sep), "__NDIAG_REP__")
var unescRep = strings.NewReplacer("__NDIAG_REP__", fmt.Sprintf("%s%s", Esc, Sep))
var qRep = strings.NewReplacer(fmt.Sprintf("%s%s", Esc, Q), "__NDIAG_REP__")
var unqRep = strings.NewReplacer("__NDIAG_REP__", fmt.Sprintf("%s%s", Esc, Q))

const DefaultDocPath = "archdoc"

var DefaultConfigFilePaths = []string{"ndiag.yml"}
var DefaultDescPath = "ndiag.descriptions"
var DefaultIconPath = "ndiag.icons"

// DefaultFormat is the default view format
const DefaultFormat = "svg"

// Element is graph element
type Element interface {
	Id() string
	FullName() string
}

// Attr is attribute of graph element/edge
type Attr struct {
	Key   string
	Value string
}

// Edge is graph edge
type Edge struct {
	Src      *Component
	Dst      *Component
	Desc     string
	Relation *Relation
	Attrs    []*Attr
}

type Config struct {
	Name              string             `yaml:"name"`
	Desc              string             `yaml:"desc,omitempty"`
	DocPath           string             `yaml:"docPath"`
	DescPath          string             `yaml:"descPath,omitempty"`
	IconPath          string             `yaml:"iconPath,omitempty"`
	Graph             *Graph             `yaml:"graph,omitempty"`
	HideViews         bool               `yaml:"hideViews,omitempty"`
	HideLayers        bool               `yaml:"hideLayers,omitempty"`
	HideRealNodes     bool               `yaml:"hideRealNodes,omitempty"`
	HideLabelGroups   bool               `yaml:"hideLabelGroups,omitempty"`
	Views             []*View            `yaml:"views"`
	Nodes             []*Node            `yaml:"nodes"`
	Relations         []*Relation        `yaml:"relations,omitempty"`
	Dict              *dict.Dict         `yaml:"dict,omitempty"`
	BaseColor         string             `yaml:"baseColor,omitempty"`
	TextColor         string             `yaml:"textColor,omitempty"`
	CustomIcons       []*glyph.Blueprint `yaml:"customIcons,omitempty"`
	basePath          string
	rawRelations      []*rawRelation
	realNodes         []*RealNode
	layers            []*Layer
	clusters          Clusters
	globalComponents  []*Component
	clusterComponents []*Component
	nodeComponents    []*Component
	edges             []*Edge
	labels            Labels
	colorSets         ColorSets
	iconMap           *IconMap
}

type Graph struct {
	Format        string        `yaml:"format,omitempty"`
	MapSliceAttrs yaml.MapSlice `yaml:"attrs,omitempty"`
}

func (g *Graph) Attrs() []*Attr {
	attrs := []*Attr{}
	for _, kv := range g.MapSliceAttrs {
		attrs = append(attrs, &Attr{
			Key:   kv.Key.(string),
			Value: kv.Value.(string),
		})
	}
	return attrs
}

func New() *Config {
	return &Config{
		Graph: &Graph{},
		Dict:  &dict.Dict{},
	}
}

func (cfg *Config) Format() string {
	if cfg.Graph.Format != "" {
		return cfg.Graph.Format
	}
	return DefaultFormat
}

func (cfg *Config) IconMap() *IconMap {
	return cfg.iconMap
}

func (cfg *Config) PrimaryView() *View {
	return cfg.Views[0]
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

func (cfg *Config) Components() []*Component {
	return append(append(cfg.globalComponents, cfg.nodeComponents...), cfg.clusterComponents...)
}

func (cfg *Config) Edges() []*Edge {
	return cfg.edges
}

func (cfg *Config) Labels() Labels {
	return cfg.labels
}

func (cfg *Config) ColorSets() ColorSets {
	return cfg.colorSets
}

func (cfg *Config) BuildNestedClusters(layers []string) (Clusters, []*Node, []*Edge, error) {
	edges := []*Edge{}
	if len(layers) == 0 {
		return Clusters{}, cfg.Nodes, cfg.edges, nil
	}
	clusters, globalNodes, err := buildNestedClusters(cfg.Clusters(), layers, cfg.Nodes)
	if err != nil {
		return clusters, globalNodes, nil, err
	}

	for _, e := range cfg.edges {
		hBelongTo := false
		tBelongTo := false
		for _, l := range layers {
			if e.Src.Cluster == nil || strings.EqualFold(e.Src.Cluster.Layer.Id(), l) {
				hBelongTo = true
			}
			if e.Dst.Cluster == nil || strings.EqualFold(e.Dst.Cluster.Layer.Id(), l) {
				tBelongTo = true
			}
		}
		if hBelongTo && tBelongTo {
			edges = append(edges, e)
		}
	}

	return clusters, globalNodes, edges, nil
}

func (cfg *Config) PruneClustersByLabels(clusters Clusters, globalNodes []*Node, globalComponents []*Component, edges []*Edge, labels []string) (Clusters, []*Node, []*Component, []*Edge, error) {
	filteredEdges := []*Edge{}
	nIds := orderedmap.NewOrderedMap()
	cIds := orderedmap.NewOrderedMap()
	comIds := orderedmap.NewOrderedMap()

	allowLabels := Labels{}
	for _, s := range labels {
		l, err := cfg.FindLabel(s)
		if err != nil {
			return clusters, globalNodes, globalComponents, edges, err
		}
		allowLabels = append(allowLabels, l)
	}

	// collect filtered nodes
	for _, n := range cfg.Nodes {
		if len(allowLabels) == 0 || len(n.Labels.Subtract(allowLabels)) > 0 {
			nIds.Set(n.Id(), n)
			for _, c := range n.Components {
				comIds.Set(c.Id(), c)
			}
			continue
		}
		for _, c := range n.Components {
			if len(allowLabels) == 0 || len(c.Labels.Subtract(allowLabels)) > 0 {
				comIds.Set(c.Id(), c)
				nIds.Set(n.Id(), n)
			}
		}
	}

	// collect filtered components
	for _, c := range cfg.Components() {
		if len(allowLabels) == 0 || len(c.Labels.Subtract(allowLabels)) > 0 {
			comIds.Set(c.Id(), c)
			if c.Cluster != nil {
				cIds.Set(c.Cluster.Id(), c.Cluster)
			}
		}
	}

	// collect filtered cluster/nodes/components/edges using edges
	for _, e := range edges {
		if len(allowLabels) != 0 && len(e.Relation.Labels.Subtract(allowLabels)) == 0 {
			continue
		}

		filteredEdges = append(filteredEdges, e)

		// src
		comIds.Set(e.Src.Id(), e.Src)
		switch {
		case e.Src.Node != nil:
			// node component
			nIds.Set(e.Src.Node.Id(), e.Src.Node)
		case e.Src.Cluster != nil:
			// cluster component
			cIds.Set(e.Src.Cluster.Id(), e.Src.Cluster)
		}

		// dst
		comIds.Set(e.Dst.Id(), e.Dst)
		switch {
		case e.Dst.Node != nil:
			// node component
			nIds.Set(e.Dst.Node.Id(), e.Dst.Node)
		case e.Dst.Cluster != nil:
			// cluster component
			cIds.Set(e.Dst.Cluster.Id(), e.Dst.Cluster)
		}
	}

	for _, k := range cIds.Keys() {
		v, _ := cIds.Get(k)
		c := v.(*Cluster)
		if !clusters.Contains(c) {
			clusters = append(clusters, c)
		}
	}

	// prune cluster nodes
	pruneClusters(clusters, nIds, comIds)

	// global nodes
	filteredNodes := []*Node{}
	for _, n := range globalNodes {
		_, ok := nIds.Get(n.Id())
		if ok {
			filteredComponents := []*Component{}
			for _, c := range n.Components {
				_, ok := comIds.Get(c.Id())
				if ok {
					filteredComponents = append(filteredComponents, c)
				}
			}
			n.Components = filteredComponents
			filteredNodes = append(filteredNodes, n)
		}
	}
	globalNodes = filteredNodes

	// global components
	filteredComponents := []*Component{}
	for _, c := range globalComponents {
		_, ok := comIds.Get(c.Id())
		if ok {
			filteredComponents = append(filteredComponents, c)
		}
	}
	globalComponents = filteredComponents

	return clusters, globalNodes, globalComponents, filteredEdges, nil
}

func (cfg *Config) PruneNodesByLabels(nodes []*Node, labels []string) ([]*Node, error) {
	if len(labels) == 0 {
		return nodes, nil
	}
	nIds := orderedmap.NewOrderedMap()
	comIds := orderedmap.NewOrderedMap()

	for _, name := range labels {
		l, err := cfg.FindLabel(name)
		if err != nil {
			return nodes, nil
		}
		edges := SplitRelations(l.Relations)

		for _, e := range edges {
			switch {
			case e.Src.Node != nil:
				nIds.Set(e.Src.Node.Id(), e.Src.Node)
			}
			comIds.Set(e.Src.Id(), e.Src)

			switch {
			case e.Dst.Node != nil:
				nIds.Set(e.Dst.Node.Id(), e.Dst.Node)
			}
			comIds.Set(e.Dst.Id(), e.Dst)
		}
	}

	filteredNodes := []*Node{}
	for _, n := range nodes {
		_, ok := nIds.Get(n.Id())
		if ok {
			filteredComponents := []*Component{}
			for _, c := range n.Components {
				_, ok := comIds.Get(c.Id())
				if ok {
					filteredComponents = append(filteredComponents, c)
				}
			}
			n.Components = filteredComponents
			filteredNodes = append(filteredNodes, n)
		}
	}
	nodes = filteredNodes

	return nodes, nil
}

func (cfg *Config) loadConfig(in []byte) error {
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
	cfg.basePath, err = filepath.Abs(filepath.Dir(path))
	if err != nil {
		return err
	}
	return cfg.loadConfig(buf)
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
	if err := cfg.checkFormat(); err != nil {
		return err
	}
	if err := cfg.buildDefault(); err != nil {
		return err
	}
	if err := cfg.buildIconMap(); err != nil {
		return err
	}
	if err := cfg.buildClusters(); err != nil {
		return err
	}
	if err := cfg.buildNodes(); err != nil {
		return err
	}
	if err := cfg.buildComponents(); err != nil {
		return err
	}
	if err := cfg.buildRelations(); err != nil {
		return err
	}
	if err := cfg.buildColors(); err != nil {
		return err
	}
	if err := cfg.checkUnique(); err != nil {
		return err
	}
	if len(cfg.Views) == 0 {
		cfg.Views = append(cfg.Views, &View{
			Name:   "Nodes",
			Layers: []string{},
		})
	}
	if err := cfg.buildDescriptions(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) BuildForIcons() error {
	if err := cfg.checkFormat(); err != nil {
		return err
	}
	if err := cfg.buildDefault(); err != nil {
		return err
	}
	if err := cfg.buildIconMap(); err != nil {
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
					n.rawClusters = merge(n.rawClusters, rn.rawClusters)
					n.rawComponents = merge(n.rawComponents, rn.rawComponents)
					break NN
				}
			}
			if !belongTo {
				return fmt.Errorf("there is a real node '%s' that does not belong to a node", newRn.Name)
			}
			cfg.realNodes = append(cfg.realNodes, newRn)
		}
		// replace component id ( real node name -> node name )
		for _, rel := range rConfig.rawRelations {
			replaced := []string{}
			for _, c := range rel.Components {
				splitted := sepSplit(c)
				if len(splitted) == 2 {
				RL:
					for _, n := range cfg.Nodes {
						if n.nameRe.MatchString(splitted[0]) {
							splitted[0] = n.Name
							break RL
						}
					}
				}
				replaced = append(replaced, sepJoin(splitted))
			}
			rel.Components = replaced
		}
		cfg.rawRelations = append(cfg.rawRelations, uniqueRawRelations(rConfig.rawRelations)...)
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

func (cfg *Config) FindComponent(s string) (*Component, error) {
	name := queryTrim(s)
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
		if strings.EqualFold(c.Id(), name) {
			return c, nil
		}
	}
	return nil, fmt.Errorf("component not found: %s", name)
}

func (cfg *Config) FindLabel(name string) (*Label, error) {
	for _, l := range cfg.Labels() {
		if l.Name == name {
			return l, nil
		}
	}
	return nil, fmt.Errorf("label not found: %s", name)
}

func (cfg *Config) FindOrCreateLabel(name string) *Label {
	for _, l := range cfg.Labels() {
		if l.Name == name {
			return l
		}
	}
	l := &Label{
		Name: name,
	}
	cfg.labels = append(cfg.labels, l)
	return l
}

func (cfg *Config) FindLayer(s string) (*Layer, error) {
	for _, l := range cfg.Layers() {
		if strings.EqualFold(s, l.Id()) {
			return l, nil
		}
	}
	return nil, fmt.Errorf("layer not found: %s", s)
}

func buildNestedClusters(clusters Clusters, layers []string, nodes []*Node) (Clusters, []*Node, error) {
	if len(layers) == 0 {
		return clusters, nodes, nil
	}
	leaf := layers[len(layers)-1]
	layers = layers[:len(layers)-1]

	globalNodes := []*Node{}
	belongTo := []*Node{}
	for _, n := range nodes {
		c := n.Clusters.FindByLayer(leaf)
		if len(c) == 0 {
			globalNodes = append(globalNodes, n)
			continue
		}
		if len(c) > 1 {
			return nil, nil, fmt.Errorf("duplicate layer: %s", leaf)
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
				return nil, nil, fmt.Errorf("duplicate layer: %s", parent)
			}
			if len(pc) == 0 {
				continue
			}
			if c[0].Parent != nil && c[0].Parent.Id() != pc[0].Id() {
				return nil, nil, fmt.Errorf("belong to two or more clusters: '%s' belongs to '%s' and '%s'", c[0].FullName(), c[0].Parent.FullName(), pc[0].FullName())
			}
			c[0].Parent = pc[0]
			pc[0].Children = append(pc[0].Children, c[0])
			c = pc
		}
	}

	// build a direct member node of a cluster
	for _, c := range clusters {
		if strings.EqualFold(c.Layer.Id(), leaf) {
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
			if c.Parent == nil && (strings.EqualFold(c.Layer.Id(), leaf) || len(c.Nodes) > 0) {
				for _, n := range c.Nodes {
					for _, rn := range globalNodes {
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

	return buildNestedClusters(clusters, layers, globalNodes)
}

func (cfg *Config) checkFormat() error {
	if cfg.HideViews && len(cfg.Views) > 1 {
		return errors.New("can't make hideViews true if you have more than one views defined")
	}
	if len(cfg.realNodes) > 0 {
		for _, n := range cfg.Nodes {
			if len(n.RealNodes) == 0 {
				return fmt.Errorf("'%s' does not have any real nodes", n.FullName())
			}
		}
	}
	if cfg.Format() != "svg" && cfg.Format() != "png" {
		return fmt.Errorf("invalid format: %s", cfg.Format())
	}
	return nil
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

func pruneClusters(clusters []*Cluster, nIds, comIds *orderedmap.OrderedMap) {
	for _, c := range clusters {
		filteredNodes := []*Node{}
		for _, n := range c.Nodes {
			_, ok := nIds.Get(n.Id())
			if ok {
				// node component
				filteredComponents := []*Component{}
				for _, com := range n.Components {
					_, ok := comIds.Get(com.Id())
					if ok {
						filteredComponents = append(filteredComponents, com)
					}
				}
				n.Components = filteredComponents

				filteredNodes = append(filteredNodes, n)
			}
		}
		c.Nodes = filteredNodes
		filteredComponents := []*Component{}
		for _, com := range c.Components {
			_, ok := comIds.Get(com.Id())
			if ok {
				filteredComponents = append(filteredComponents, com)
			}
		}
		c.Components = filteredComponents

		pruneClusters(c.Children, nIds, comIds)
	}
}

func (cfg *Config) parseComponent(comName string) (*Component, error) {
	c := &Component{}
	if queryContains(comName) {
		var m ComponentMetadata
		splited := querySplit(comName)
		c.Name = splited[0]
		if err := qs.Unmarshal(&m, splited[1]); err != nil {
			return nil, err
		}
		if m.Icon != "" {
			_, err := cfg.IconMap().Get(m.Icon)
			if err != nil {
				return nil, fmt.Errorf("not found icon: %s", m.Icon)
			}
		}
		c.Metadata = m
	} else {
		c.Name = comName
	}
	return c, nil
}

func querySplit(s string) []string {
	splitted := strings.Split(qRep.Replace(s), Q)
	unescaped := []string{}
	for _, ss := range splitted {
		unescaped = append(unescaped, unqRep.Replace(ss))
	}
	return unescaped
}

func queryContains(s string) bool {
	return strings.Contains(qRep.Replace(s), Q)
}

func queryTrim(s string) string {
	splitted := sepSplit(s)
	trimed := []string{}
	for _, ss := range splitted {
		splitted := querySplit(ss)
		trimed = append(trimed, splitted[0])
	}
	return sepJoin(trimed)
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

func sepJoin(ss []string) string {
	return strings.Join(ss, Sep)
}

func sepContains(s string) bool {
	return strings.Contains(escRep.Replace(s), Sep)
}

func merge(a, b []string) []string {
	m := orderedmap.NewOrderedMap()
	for _, s := range a {
		m.Set(s, s)
	}
	for _, s := range b {
		m.Set(s, s)
	}
	o := []string{}
	for _, k := range m.Keys() {
		s, _ := m.Get(k)
		o = append(o, s.(string))
	}
	return o
}
