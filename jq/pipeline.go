package jq

import (
	"bytes"
	"context"
	"fmt"

	"github.com/itchyny/gojq"
	"go.bobheadxi.dev/streamline/pipeline"
)

// Pipeline builds a JQ query for a pipeline that runs the query against each line and
// maps the result to the output. If the query fails to build, Pipeline will return a
// pipeline that returns an error immediately on read - to handle query build errors, use
// BuildPipeline instead.
//
// Internally, Pipeline uses github.com/itchyny/gojq to build and run the query.
func Pipeline(query string) pipeline.Pipeline {
	return PipelineContext(context.Background(), query)
}

// PipelineContext is the same as Pipeline, but runs the generated JQ code in the given
// context.
func PipelineContext(ctx context.Context, query string) pipeline.Pipeline {
	p, err := BuildPipeline(ctx, query)
	if err != nil {
		return pipeline.MapErr(func(line []byte) ([]byte, error) { return nil, err })
	}
	return p
}

// Pipeline safely builds a pipeline that runs the JQ query against each line and maps the
// result to the output, returning an error if the query fails to build. The generated JQ
// code is run in the given context.
//
// Internally, BuildPipeline uses github.com/itchyny/gojq to build and run the query.
func BuildPipeline(ctx context.Context, query string) (pipeline.Pipeline, error) {
	jqCode, err := buildJQ(query)
	if err != nil {
		return nil, err
	}
	return &jqCodePipeline{
		ctx:          ctx,
		code:         jqCode,
		resultBuffer: &bytes.Buffer{},
	}, nil
}

type jqCodePipeline struct {
	ctx          context.Context
	code         *gojq.Code
	resultBuffer *bytes.Buffer
}

func (p *jqCodePipeline) ProcessLine(line []byte) ([]byte, error) {
	if len(line) == 0 {
		return line, nil
	}

	// Reset buffer - by the time a new line is processed, nobody should be
	// holding a reference to the previous results.
	p.resultBuffer.Reset()

	// Run code, populating resultBuffer with the output
	err := execJQ(p.ctx, p.code, line, p.resultBuffer)
	if err != nil {
		// Embed the consumed content for ease of debugging
		return nil, fmt.Errorf("%w: %s", err, string(line))
	}
	return p.resultBuffer.Bytes(), nil
}
