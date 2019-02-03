package main

import (
	"os"
	"goask/core/adapter/fakeadapter"
	"goask/graphqlhelper"
	"goask/resolver"
	"log"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

func main() {
	schemas, err := graphqlhelper.ReadSchemas(
		"./resolver/schema/schema.graphql",
		"./resolver/schema/query.graphql",
		"./resolver/schema/mutation.graphql",
	)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("./data.json", os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	data, err := fakeadapter.NewData(file)
	if err != nil {
		panic(err)
	}

	schema, err := graphql.ParseSchema(schemas, &resolver.Root{
		Query: resolver.Query{
			Data: &data,
		},
		Mutation: resolver.Mutation{
			Data: &data,
		},
	})
	if err != nil {
		panic(err)
	}

	http.Handle("/query", &relay.Handler{Schema: schema})
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
