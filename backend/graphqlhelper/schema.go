package graphqlhelper

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/graph-gophers/graphql-go"
)

var ErrNoSchemaSupplied = errors.New("no schema files are supplied")

// ReadSchemas reads multiple files and concatenate their content into one string.
func ReadSchemas(files ...string) (string, error) {
	if len(files) == 0 {
		return "", ErrNoSchemaSupplied
	}
	builder := strings.Builder{}
	for _, file := range files {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return "", err
		}
		_, err = builder.Write(b)
		if err != nil {
			return "", err
		}
	}
	return builder.String(), nil
}

type logger interface {
	Printf(format string, v ...interface{})
}

type LoggableSchema struct {
	Logger logger
	Schema *graphql.Schema
}

func (l *LoggableSchema) Exec(ctx context.Context, queryString string, operationName string, variables map[string]interface{}) *graphql.Response {
	if l.Logger == nil {
		l.Logger = log.New(os.Stdout, "", log.LstdFlags)
	}
	resp := l.Schema.Exec(ctx, queryString, operationName, variables)
	l.Logger.Printf("query: %s", queryString)
	return resp
}
