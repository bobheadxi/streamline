package pipeline

// Map is a Pipeline that allows modifications of individual lines from streamline.Stream.
// Implementations can return a nil []byte to indicate a line is to be skipped.
type Map func(line []byte) []byte

var _ Pipeline = (Map)(nil)

func (m Map) Inactive() bool { return m == nil }

func (m Map) ProcessLine(line []byte) ([]byte, error) {
	if m.Inactive() {
		return line, nil
	}
	return m(line), nil
}

// MapErr is a Pipeline that allows modifications of individual lines from
// streamline.Stream with error handling. Implementations can return a nil []byte to
// indicate a line is to be skipped.
//
// Errors interrupt line processing and are propagated to streamline.Stream.
type MapErr func(line []byte) ([]byte, error)

var _ Pipeline = (MapErr)(nil)

func (m MapErr) Inactive() bool { return m == nil }

func (m MapErr) ProcessLine(line []byte) ([]byte, error) {
	if m.Inactive() {
		return line, nil
	}
	return m(line)
}
