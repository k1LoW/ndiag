package config

type Graph struct {
	Format string
	Attrs  Attrs
}

func (dest *Graph) Merge(src *Graph) error {
	if src.Format != "" {
		dest.Format = src.Format
	}
	dest.Attrs = append(dest.Attrs, src.Attrs...)
	return nil
}
