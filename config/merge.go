package config

import (
	"errors"
)

func (dest *Config) Merge(src *Config) (err error) {
	if len(dest.Relations) > 0 || len(src.Relations) > 0 {
		return errors.New("it should be before the Config.Build phase")
	}
	if len(dest.realNodes) > 0 || len(src.realNodes) > 0 {
		return errors.New("it should be before the Config.Build phase")
	}
	if len(dest.layers) > 0 || len(src.layers) > 0 {
		return errors.New("it should be before the Config.Build phase")
	}
	if len(dest.clusters) > 0 || len(src.clusters) > 0 {
		return errors.New("it should be before the Config.Build phase")
	}
	if len(dest.globalComponents) > 0 || len(src.globalComponents) > 0 {
		return errors.New("it should be before the Config.Build phase")
	}
	if len(dest.clusterComponents) > 0 || len(src.clusterComponents) > 0 {
		return errors.New("it should be before the Config.Build phase")
	}
	if len(dest.nodeComponents) > 0 || len(src.nodeComponents) > 0 {
		return errors.New("it should be before the Config.Build phase")
	}
	if len(dest.edges) > 0 || len(src.edges) > 0 {
		return errors.New("it should be before the Config.Build phase")
	}
	if len(dest.labels) > 0 || len(src.labels) > 0 {
		return errors.New("it should be before the Config.Build phase")
	}
	if len(dest.colorSets) > 0 || len(src.colorSets) > 0 {
		return errors.New("it should be before the Config.Build phase")
	}
	if dest.iconMap != nil || src.iconMap != nil {
		return errors.New("it should be before the Config.Build phase")
	}

	if src.Name != "" {
		dest.Name = src.Name
	}
	if src.Desc != "" {
		dest.Desc = src.Desc
	}
	if src.DocPath != "" {
		dest.DocPath = src.DocPath
	}
	if src.DescPath != "" {
		dest.DescPath = src.DescPath
	}
	if src.IconPath != "" {
		dest.IconPath = src.IconPath
	}
	if err := dest.Graph.Merge(src.Graph); err != nil {
		return err
	}
	if src.HideViews {
		dest.HideViews = src.HideViews
	}
	if src.HideLayers {
		dest.HideLayers = src.HideLayers
	}
	if src.HideRealNodes {
		dest.HideRealNodes = src.HideRealNodes
	}
	if src.HideLabels {
		dest.HideLabels = src.HideLabels
	}
	dest.Views = dest.Views.Merge(src.Views)
	dest.Nodes, err = dest.Nodes.Merge(src.Nodes)
	if err != nil {
		return err
	}
	dest.Dict.Merge(src.Dict.Dump())
	if src.BaseColor != "" {
		dest.BaseColor = src.BaseColor
	}
	if src.TextColor != "" {
		dest.TextColor = src.TextColor
	}
	dest.CustomIcons = dest.CustomIcons.Merge(src.CustomIcons)
	if src.basePath != "" {
		dest.basePath = src.basePath
	}
	dest.rawRelations = dest.rawRelations.Merge(src.rawRelations)

	return nil
}
