package streamline

import "go.bobheadxi.dev/streamline/pipeline"

// safePipeline is safe to use when itself or its inner pipeline is nil.
type safePipeline struct{ pipeline pipeline.Pipeline }

func (p *safePipeline) isNil() bool    { return p == nil || p.pipeline == nil }
func (p *safePipeline) Inactive() bool { return p.isNil() || p.pipeline.Inactive() }
func (p *safePipeline) ProcessLine(line []byte) ([]byte, error) {
	if p.isNil() {
		return line, nil
	}
	return p.pipeline.ProcessLine(line)
}
