package config

type Graph struct {
	Format string
	Attrs  Attrs
}

func (dest *Graph) Merge(src *Graph) (err error) {
	if src.Format != "" {
		dest.Format = src.Format
	}
	dest.Attrs = dest.Attrs.Merge(src.Attrs)
	return nil
}
