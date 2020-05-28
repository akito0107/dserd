package dserd

import (
	"github.com/awalterschulze/gographviz"
	"golang.org/x/xerrors"
)

func MakeGraph(schemas map[string]*SchemaInfo) (*gographviz.Graph, error) {
	g := gographviz.NewGraph()

	for kind, schema := range schemas {
		attrs := make(map[string]string)
		for _, p := range schema.Props {
			attrs[p.Name] = p.Type
		}
		if err := g.AddNode("", kind, attrs); err != nil {
			return nil, xerrors.Errorf("addNode kind: %s failed: %w", kind, err)
		}
	}

	return g, nil
}

