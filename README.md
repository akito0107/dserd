WIP: Datastore Kind Relationship Diagram generator
---

```shell
$ gcloud beta emulators datastore start
$ go run cmd/dserd/main.go > graph.dot
$ dot -Tsvg graph.dot -o dot.svg
```