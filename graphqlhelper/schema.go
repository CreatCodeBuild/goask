package graphqlhelper

import (
	"errors"
	"io/ioutil"
	"strings"
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
