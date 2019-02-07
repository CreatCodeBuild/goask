package main

import (
	"bytes"
	"io/ioutil"
	"goask/core/adapter/fakeadapter"
	"goask/graphqlhelper"
	"goask/resolver"
	"log"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/gorilla/mux"
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

	data, err := fakeadapter.NewData(fakeadapter.NewFileSerializer("./data.json"))
	if err != nil {
		panic(err)
	}

	schema, err := graphql.ParseSchema(schemas, &resolver.Root{
		Query: resolver.Query{
			Data: data,
		},
		Mutation: resolver.Mutation{
			Data: data,
		},
	})
	if err != nil {
		panic(err)
	}

	handler := &relay.Handler{Schema: schema}
	
	r := mux.NewRouter()
	r.Handle("/query", handler)

	r.Use(func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, _ := ioutil.ReadAll(r.Body)
			r.Body = ioutil.NopCloser(bytes.NewReader(b))
			log.Println(r.RequestURI, string(b))
			next.ServeHTTP(&ResponseWriterLogger{w}, r)
		})
	})

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

type ResponseWriterLogger struct {
	rw http.ResponseWriter
}

func (r *ResponseWriterLogger) Header() http.Header {
	return r.rw.Header()
}

func (r *ResponseWriterLogger) Write(b []byte) (int, error) {
	return r.rw.Write(b)
}

func (r  *ResponseWriterLogger) WriteHeader(statusCode int) {
	r.rw.WriteHeader(statusCode)
}
