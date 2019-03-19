package id

type Generator struct {
	counter uint64
}

func (g *Generator) ID() (uint64, error) {
	g.counter += 1
	return g.counter, nil
}

func NewGenerator() *Generator {
	return &Generator{counter: 0}
}

type PersistentGenerator struct {
	Generator
	filename string
}

func (g PersistentGenerator) ID() (uint64, error) {
	defer func() {
		// todo: sync the file
	}()
	return g.Generator.ID()
}

func NewGeneratorFromFile(filename string) (PersistentGenerator, error) {
	// todo: read the file
	var cur uint64 = 0
	return PersistentGenerator{filename: filename, Generator: Generator{counter: cur}}, nil
}
