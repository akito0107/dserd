package dserd

import (
	"fmt"

	"cloud.google.com/go/datastore"
)

type AnyEntity struct {
	Kind       string
	Properties map[string]interface{}
}

func NewAnyEntity(kind string) *AnyEntity {
	return &AnyEntity{
		Kind:       kind,
		Properties: make(map[string]interface{}),
	}
}

func (a *AnyEntity) Load(properties []datastore.Property) error {
	for _, p := range properties {
		a.Properties[p.Name] = p.Value
	}
	return nil
}

func (a *AnyEntity) Save() ([]datastore.Property, error) {
	// unused
	return nil, nil
}

func (a *AnyEntity) String() string {
	str := "Kind: " + a.Kind + "\n  props: \n"
	for name, p := range a.Properties {
		str += fmt.Sprintf("    name: %s, value: %+v \n", name, p)
	}
	return str
}