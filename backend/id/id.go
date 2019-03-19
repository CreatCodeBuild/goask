package id

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Generator struct {
	Counter uint64
}

func (g *Generator) ID() (uint64, error) {
	g.Counter += 1
	return g.Counter, nil
}

func NewGenerator() *Generator {
	return &Generator{Counter: 0}
}

type PersistentGenerator struct {
	Generator
	filename string
}

func (g PersistentGenerator) ID() (id uint64, err error) {
	defer func() {
		b, err2 := json.Marshal(g.Generator)
		if err2 != nil {
			err = err2
			return
		}
		err = ioutil.WriteFile(g.filename, b, os.ModePerm)
		if err != nil {
			return
		}
	}()
	return g.Generator.ID()
}

func NewGeneratorFromFile(filename string) (*PersistentGenerator, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		if err.Error() == `open id.json: no such file or directory` {
			b = []byte(`{"Counter":0}`)
		} else {
			return nil, err
		}
	}

	gen := Generator{}
	err = json.Unmarshal(b, &gen)
	if err != nil {
		return nil, err
	}

	return &PersistentGenerator{filename: filename, Generator: gen}, nil
}
