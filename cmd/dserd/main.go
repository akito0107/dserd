package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/grpc"

	"github.com/akito0107/dserd"
)

const projectId = "local-test"

func main() {
	ctx := context.Background()

	opts := []option.ClientOption{
		option.WithEndpoint("localhost:8081"),
		option.WithGRPCDialOption(grpc.WithInsecure()),
		option.WithoutAuthentication(),
	}
	client, err := datastore.NewClient(context.Background(), projectId, opts...)
	if err != nil {
		log.Fatal(err)
	}
	putTestEntity(ctx, client)

	schemaInfos := make(map[string]*dserd.SchemaInfo)
	loadKinds := []string{"User", "Post"}

	for _, k := range loadKinds {
		si, err := dserd.LoadKind(ctx, client, k)
		if err != nil {
			log.Fatal(err)
		}
		schemaInfos[k] = si
	}

	for _, s := range schemaInfos {
		if err := s.ResolveAncestor(ctx, client, schemaInfos); err != nil {
			log.Fatal(err)
		}
	}

	graph, err := dserd.MakeGraph(schemaInfos)
	if err != nil {
		log.Fatal(err)
	}

	str := graph.String()
	fmt.Fprint(os.Stdout, str)
}


type User struct {
	Name string
}

type Post struct {
	Title   string
	Content string
	User    *datastore.Key
}

func putTestEntity(ctx context.Context, client *datastore.Client) {

	k := datastore.NameKey("User", "aaa", nil)
	u := &User{Name: "aaa"}
	if _, err := client.Put(ctx, k, u); err != nil {
		log.Printf("%+v", err)
	}

	pk := datastore.NameKey("Post", uuid.New().String(), k)
	p := &Post{
		Title:   "aaa",
		Content: "bbb",
		User:    k,
	}
	if _, err := client.Put(ctx, pk, p); err != nil {
		log.Printf("%+v", err)
	}
}