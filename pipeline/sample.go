package pipeline

// Sample is a Pipeline that only includes every Nth line from streamline.Stream.
type Sample struct {
	// N indicates that this Sample pipeline should only retain every Nth line from the
	// output.
	//
	// If set to 0 or 1, this Pipeline is marked as inactive.
	N int

	// current is an internal counter.
	current int
}

var _ Pipeline = (Filter)(nil)

func (s *Sample) Inactive() bool { return s == nil || s.N == 0 || s.N == 1 }

func (s *Sample) ProcessLine(line []byte) ([]byte, error) {
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
