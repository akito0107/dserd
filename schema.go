package dserd

import (
	"context"

	"cloud.google.com/go/datastore"
	"golang.org/x/xerrors"
	"google.golang.org/api/iterator"
)

type Property struct {
	Name        string
	Type        string
	IsReference bool
}

type SchemaInfo struct {
	Kind      string
	Props     []*Property
	Ancestors []*SchemaInfo
}

func (s *SchemaInfo) HasAncestors() bool {
	for _, p := range s.Props {
		if p.IsReference {
			return true
		}
	}
	return false
}

func (s *SchemaInfo) ResolveAncestor(ctx context.Context, client *datastore.Client, schemas map[string]*SchemaInfo) error {
	var projections []string
	for _, p := range s.Props {
		if p.IsReference {
			projections = append(projections, p.Name)
		}
	}

	query := datastore.NewQuery(s.Kind).Project(projections...)
	for _, p := range projections {
		query = query.Filter(p+" > ", "")
	}
	query.Limit(1)

	it := client.Run(ctx, query)
	var ancestors []*SchemaInfo

	for {
		res := NewAnyEntity(s.Kind)
		_, err := it.Next(res)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return xerrors.Errorf("load failed: %w", err)
		}

		for _, p := range res.Properties {
			if k, ok := p.(*datastore.Key); ok {
				if ancestor, ok :=  schemas[k.Kind]; ok {
					ancestors = append(ancestors, ancestor)
				}
			}
		}
	}
	s.Ancestors = ancestors

	return nil
}
