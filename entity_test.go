package dserd

import (
	"context"
	"testing"

	"cloud.google.com/go/datastore"
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
	}
	if _, err := client.Put(ctx, pk, p); err != nil {
		t.Logf("%+v", err)
	}
}

func TestEntityDumper(t *testing.T) {
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

	t.Logf("\n")

	e, err := Load(ctx, client, "User")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(e.String())
	e, err = Load(ctx, client, "Post")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(e.String())
}
