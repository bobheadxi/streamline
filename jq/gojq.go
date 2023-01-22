package jq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/itchyny/gojq"
)

// buildJQ parses and compiles a jq query.
func buildJQ(query string) (*gojq.Code, error) {
	jq, err := gojq.Parse(query)
	if err != nil {
		return nil, fmt.Errorf("jq.Parse: %w", err)
	}
	jqCode, err := gojq.Compile(jq)
	if err != nil {
		return nil, fmt.Errorf("jq.Compile: %w", err)
	}
	return jqCode, nil
}

// execJQ executes the compiled jq query against content from reader.
func execJQ(ctx context.Context, jqCode *gojq.Code, data []byte) ([]byte, error) {
	var input interface{}
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, fmt.Errorf("json: %w", err)
	}

	var result bytes.Buffer
	iter := jqCode.RunWithContext(ctx, input)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		if err, ok := v.(error); ok {
			return nil, fmt.Errorf("jq: %w", err)
		}

		encoded, err := gojq.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("jq: %w", err)
		}
		result.Write(encoded)
	}
	return result.Bytes(), nil
}
