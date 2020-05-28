package dserd

import (
	"context"

	"cloud.google.com/go/datastore"
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

func LoadKind(ctx context.Context, client *datastore.Client, kind string) (*SchemaInfo, error) {
	kindKey := datastore.NameKey("__kind__", kind, nil)
	q := datastore.NewQuery("__property__").Ancestor(kindKey)

	type Prop struct {
		Repr []string `datastore:"property_representation"`
	}
	var props []Prop
	keys, err := client.GetAll(ctx, q, &props)
	if err != nil {
		return nil, err
	}

	var properties []*Property
	for i, v := range props {
		properties = append(properties, &Property{
			Name:        keys[i].Name,
			Type:        v.Repr[0],
			IsReference: containsReference(v.Repr),
		})
	}

	return &SchemaInfo{
		Kind:  kind,
		Props: properties,
	}, nil
}

func containsReference(refs []string) bool {
	for _, r := range refs {
		if r == "REFERENCE" {
			return true
		}
	}
	return false
}
