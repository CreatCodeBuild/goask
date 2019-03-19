package main

import (
	"encoding/json"
	"goask/core/adapter"
	"goask/core/adapter/fakeadapter"
	"goask/graphqlhelper"
	"goask/id"
	"goask/resolver"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

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
		appLogger.ErrorExit(errors.WithStack(err))
	}

	// Initialzie Server
	server := prepareServer(&graphqlhelper.LoggableSchema{Schema: schema, Logger: appLogger}, appLogger)

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
	questionDAO := fakeadapter.NewQuestionDAO(data, userDAO, id.NewGenerator()) // todo: use persistent implementation
	searcher := fakeadapter.NewSearcher(data)
	tagDAO := fakeadapter.NewTagDAO(data)
	return userDAO, answerDAO, questionDAO, searcher, tagDAO, err
}

func prepareServer(schema *graphqlhelper.LoggableSchema, logger *logger.Logger) *http.Server {
	// Initialzie GraphQL Relay Server Handler
	handler := &Handler{Schema: schema} // probably need to reimplement a relay handler to accept a interface of Schame so that I can wrap the graphql.Schema with my logging functionality

	// Initialize mux router
	r := mux.NewRouter()
	r.Handle("/query", handler)

	return &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

// Handler is a customized relay handler which accepts a Schema object wrapper which logs.
type Handler struct {
	Schema *graphqlhelper.LoggableSchema
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := h.Schema.Exec(r.Context(), params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
