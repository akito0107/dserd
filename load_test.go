package dserd

import (
	"context"
	"testing"

	"cloud.google.com/go/datastore"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

const projectId = "local-test"

type User struct {
	Name string
}

type Post struct {
	Title   string
	Content string
	User    *datastore.Key
}

func putTestEntity(t *testing.T, ctx context.Context, client *datastore.Client) {
	t.Helper()

	k := datastore.NameKey("User", "aaa", nil)
	u := &User{Name: "aaa"}
	if _, err := client.Put(ctx, k, u); err != nil {
		t.Logf("%+v", err)
	}

	pk := datastore.NameKey("Post", uuid.New().String(), k)
	p := &Post{
		Title:   "aaa",
		Content: "bbb",
		User:    k,
	}
	if _, err := client.Put(ctx, pk, p); err != nil {
		t.Logf("%+v", err)
	}
}

func TestLoadKind(t *testing.T) {
	ctx := context.Background()

	opts := []option.ClientOption{
		option.WithEndpoint("localhost:8081"),
		option.WithGRPCDialOption(grpc.WithInsecure()),
		option.WithoutAuthentication(),
	}
	client, err := datastore.NewClient(context.Background(), projectId, opts...)
	if err != nil {
		t.Fatal(err)
	}
	putTestEntity(t, ctx, client)

	schemaInfos := make(map[string]*SchemaInfo)
	u, err := LoadKind(ctx, client, "User")
	if err != nil {
		t.Fatal(err)
	}
	schemaInfos["User"] = u

	si, err := LoadKind(ctx, client, "Post")
	if err != nil {
		t.Fatal(err)
	}
	schemaInfos["Post"] = si

	expected := &SchemaInfo{
		Kind:      "Post",
		Props:     []*Property{
			{
				Name: "Content",
				Type: "STRING",
			},
			{
				Name: "Title",
				Type: "STRING",
			},
			{
				Name: "User",
				Type: "REFERENCE",
				IsReference: true,
			},
		},
		Ancestors: nil,
	}

	if diff := cmp.Diff(si, expected); diff != "" {
		t.Errorf("diff: %s", diff)
	}

	if err := si.ResolveAncestor(ctx, client, schemaInfos); err != nil {
		t.Fatalf("%+v", err)
	}
}
