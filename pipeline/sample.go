package pipeline

// Sampler is a Pipeline that only includes every Nth line from streamline.Stream.
type Sampler struct {
	// N indicates that this Sample pipeline should only retain every Nth line from the
	// output.
	//
	// If set to 0, 1, or a negative value, this Pipeline is marked as inactive.
	N int

	// current is an internal counter.
	current int
}

var _ Pipeline = (*Sampler)(nil)

// Sample creates a Sampler pipeline that only includes every nth line from
// streamline.Stream.
func Sample(n int) *Sampler {
	return &Sampler{N: n}
}

func (s *Sampler) Inactive() bool { return s == nil || s.N <= 1 }

func (s *Sampler) ProcessLine(line []byte) ([]byte, error) {
	if s.Inactive() {
		return line, nil
	}

	s.current += 1
	if s.current%s.N == 0 {
		s.current = 0
		return line, nil
	}

	return nil, nil
}
