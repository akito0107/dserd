package dserd

import (
	"github.com/awalterschulze/gographviz"
	"golang.org/x/xerrors"
)

func MakeGraph(schemas map[string]*SchemaInfo, name string) (*gographviz.Graph, error) {
	g := gographviz.NewGraph()
	g.SetName(name)
	g.SetDir(true)

	for kind, _ := range schemas {
		attrs := make(map[string]string)
		// for _, p := range schema.Props {
		// 	attrs[p.Name] = p.Type
		// }
		if err := g.AddNode(name, kind, attrs); err != nil {
			return nil, xerrors.Errorf("addNode kind: %s failed: %w", kind, err)
		}
	}

	for kind, schema := range schemas {
		attrs := make(map[string]string)

		for _, ancestor := range schema.Ancestors {
			if err := g.AddEdge(kind, ancestor.Kind, true, attrs); err != nil {
				return nil, xerrors.Errorf("addEdge %s to %s failed: %w", kind, ancestor.Kind, err)
			}
		}
	}

	return g, nil
}
