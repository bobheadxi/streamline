package pipeline

import "fmt"

// Sampler is a Pipeline that only includes every Nth line from streamline.Stream.
type Sampler struct {
	// N indicates that this Sample pipeline should only retain every Nth line from the
	// output. If N is 0, all lines are skipped; if N is 1, all lines are retained.
	N int

	// current is an internal counter.
	current int
}

var _ Pipeline = (*Sampler)(nil)

// Sample creates a Sampler pipeline that only includes every nth line from
// streamline.Stream. If N is 0, all lines are skipped; if N is 1, all lines are retained.
// Negative values will result in an error.
func Sample(n int) *Sampler {
	return &Sampler{N: n}
}

func (s *Sampler) ProcessLine(line []byte) ([]byte, error) {
	switch s.N {
	case -1:
		return nil, fmt.Errorf("invalid n %d", s.N)
	case 0:
		return nil, nil // never sample
	case 1:
		return line, nil // always sample
	default:
		s.current += 1
		if s.current%s.N == 0 {
			s.current = 0 // reset counter
			return line, nil
		}
		return nil, nil
	}
}
