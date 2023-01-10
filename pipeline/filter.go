package pipeline

// Filter is a Pipeline that allows omission of individual lines from streamline.Stream
// by returning false on lines that should be skipped.
type Filter func(line []byte) bool

var _ Pipeline = (Filter)(nil)

func (f Filter) Inactive() bool { return f == nil }

func (f Filter) ProcessLine(line []byte) ([]byte, error) {
	if f.Inactive() {
		return line, nil
	}
	if f(line) {
		return line, nil
	}
	return nil, nil
}
