package dserd

import (
	"context"

	"cloud.google.com/go/datastore"
	"golang.org/x/xerrors"
	"google.golang.org/api/iterator"
)

func Load(ctx context.Context, client *datastore.Client, kind string) (*AnyEntity, error) {
	q := datastore.NewQuery(kind).Limit(1)
	it := client.Run(ctx, q)
	e := NewAnyEntity(kind)
	_, err := it.Next(e)
	if err == iterator.Done {
		return nil, xerrors.Errorf("empty entity sets")
	}
	if err != nil {
		return nil, xerrors.Errorf("must entity load failed: %w", err)
	}
	return e, nil
}
