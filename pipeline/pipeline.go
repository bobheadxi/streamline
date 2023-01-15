package pipeline

// Pipeline implementations are used to transform the data provided to a streamline.Stream.
// For example, they are useful for mapping and pruning data. To configure a Stream to use
// a Pipeline, use (*Stream).WithPipeline(...).
//
// Note that generally a Pipeline should not be used to implement handling of data - use
// (*Stream).Stream(...) and (*Stream).StreamBytes(...) instead.
type Pipeline interface {
	// Inactive indicates if this Pipeline does anything at all. It is used internally by
	// streamline.Stream to optimize for cases where accommodating a Pipeline adds overhead.
	Inactive() bool
	// ProcessLine returns a modified, unmodified, or omitted line. To omit a line, return
	// a nil []byte - an empty []byte will be retained.
	ProcessLine(line []byte) ([]byte, error)
}

// MultiPipeline is a Pipeline that applies all its Pipelines in serial.
type MultiPipeline []Pipeline

var _ Pipeline = (MultiPipeline)(nil)

func (mp MultiPipeline) Inactive() bool {
	var active bool
	for _, p := range mp {
		if !p.Inactive() {
			active = true
		}
	}
	return !active
}

func (mp MultiPipeline) ProcessLine(line []byte) ([]byte, error) {
	var err error
	for _, p := range mp {
		if p.Inactive() {
			continue
		}
		line, err = p.ProcessLine(line)
		if err != nil {
			break
		}
		if len(line) == 0 {
			break // nothing to pass on to next pipeline
		}
	}
	return line, err
}
