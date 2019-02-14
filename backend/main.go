package main

import (
	"bytes"
	"goask/core/adapter"
	"goask/core/adapter/fakeadapter"
	"goask/graphqlhelper"
	"goask/resolver"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"

	logger "goask/log"
)

func main() {
	// Create application loggger
	appLogger := &logger.Logger{}

	// Read multiple schema files and combine them
	schemas, err := prepareSchmea()
	if err != nil {
		appLogger.ErrorExit(err)
	}

	// Initialize adapters
	userDAO, answerDAO, questionDAO, searcher, tagDAO, err := prepareDataAccessObjects()
	if err != nil {
		appLogger.ErrorExit(err)
	}

	// Initialize standard resolver with correct dependencies
	standardResolver, err := resolver.NewStdResolver(questionDAO, answerDAO, userDAO, searcher, tagDAO, appLogger)
	if err != nil {
		appLogger.ErrorExit(err)
	}

	// Initialize schema
	schema, err := graphql.ParseSchema(schemas, &resolver.Root{
		Query:    resolver.NewQuery(standardResolver),
		Mutation: resolver.NewMutation(standardResolver),
	})
	if err != nil {
		appLogger.ErrorExit(err)
	}

	// Initialzie Server
	server := prepareServer(schema)

	// Start the server
	if err := server.ListenAndServe(); err != nil {
		appLogger.ErrorExit(err)
	}
}

func prepareSchmea() (string, error) {
	return graphqlhelper.ReadSchemas(
		"./resolver/schema/schema.graphql",
		"./resolver/schema/query.graphql",
		"./resolver/schema/mutation.graphql",
		"./resolver/schema/types.graphql",
	)
}

func prepareDataAccessObjects() (adapter.UserDAO, adapter.AnswerDAO, adapter.QuestionDAO, adapter.Searcher, adapter.TagDAO, error) {
	// Initialize adapters
	data, err := fakeadapter.NewData(fakeadapter.NewFileSerializer("./data.json"))
	userDAO := fakeadapter.NewUserDAO(data)
	answerDAO := fakeadapter.NewAnswerDAO(data)
	questionDAO := fakeadapter.NewQuestionDAO(data, userDAO)
	searcher := fakeadapter.NewSearcher(data)
	tagDAO := fakeadapter.NewTagDAO(data)
	return userDAO, answerDAO, questionDAO, searcher, tagDAO, err
}

func prepareServer(schema *graphql.Schema) *http.Server {
	// Initialzie GraphQL Relay Server Handler
	handler := &relay.Handler{Schema: schema}

	// Initialize mux router
	r := mux.NewRouter()
	r.Handle("/query", handler)

	// Register a logging middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, _ := ioutil.ReadAll(r.Body)
			r.Body = ioutil.NopCloser(bytes.NewReader(b))
			log.Println(r.RequestURI, string(b))
			next.ServeHTTP(&ResponseWriterLogger{w}, r)
		})
	})

	return &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
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

func (r *ResponseWriterLogger) WriteHeader(statusCode int) {
	r.rw.WriteHeader(statusCode)
}
