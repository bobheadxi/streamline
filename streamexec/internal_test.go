package streamexec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModeSet(t *testing.T) {
	for _, tc := range []struct {
		name            string
		modes           modeSet
		wantMode        StreamMode
		wantModeMatches []StreamMode
		wantModeNot     []StreamMode
	}{
		{
			name:            "default",
			modes:           modeSet{},
			wantMode:        Combined,
			wantModeMatches: []StreamMode{Stdout, Stderr},
			wantModeNot:     []StreamMode{ErrWithStderr},
		},
		{
			name:            "only stdout",
			modes:           modeSet{Stdout},
			wantMode:        Stdout,
			wantModeMatches: []StreamMode{Stdout},
			wantModeNot:     []StreamMode{Stderr, ErrWithStderr},
		},
		{
			name:            "combine multiple",
			modes:           modeSet{Stdout, ErrWithStderr},
			wantMode:        Stdout | ErrWithStderr,
			wantModeMatches: []StreamMode{Stdout, ErrWithStderr},
			wantModeNot:     []StreamMode{Stderr},
		},
		{
			name:            "pre-combined",
			modes:           modeSet{Stdout | ErrWithStderr},
			wantMode:        Stdout | ErrWithStderr,
			wantModeMatches: []StreamMode{Stdout, ErrWithStderr},
			wantModeNot:     []StreamMode{Stderr},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			m := tc.modes.getMode()
			assert.Equal(t, tc.wantMode, m)
			for _, want := range tc.wantModeMatches {
				assert.True(t, m&want != 0)
			}
			for _, not := range tc.wantModeNot {
				assert.True(t, m&not == 0)
			}
		})
	}

}
