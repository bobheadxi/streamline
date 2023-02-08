package pipeline

// Indexed is a pipeline that provides a line handler the index of each line (i.e. how
// many lines the pipeline has processed). The first line to be processed has an index
// of 1.
func Indexed(handler func(i int, line []byte) ([]byte, error)) Pipeline {
	return &indexer{Handler: handler}
}

type indexer struct {
	Handler func(i int, line []byte) ([]byte, error)
	index   int
}

func (i *indexer) ProcessLine(line []byte) ([]byte, error) {
	i.index += 1
	return i.Handler(i.index, line)
}
