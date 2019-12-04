package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntcodeMachine_Execute(t *testing.T) {
	t.Parallel()
	separator := ","
	tests := []struct {
		name  string
		input io.Reader
		want  string
		err   error
	}{
		{
			name:  "happy path",
			input: strings.NewReader(`1,0,0,0,99`),
			want:  `2,0,0,0,99`,
			err:   EOX,
		},
		{
			name:  "happy path 2 - override halt",
			input: strings.NewReader(`1,1,1,4,99,5,6,0,99`),
			want:  `30,1,1,4,2,5,6,0,99`,
			err:   EOX,
		},
		//{
		//	name:  "happy path 3",
		//	input: strings.NewReader(`2,4,4,5,99,0`),
		//	want:  `2,4,4,5,99,9801`,
		//	err:   EOX,
		//},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewIntcodeMachineWithStdISet(tt.input, separator)
			assert.NoError(t, err)
			err = m.Execute()
			assert.Error(t, tt.err, err)
			assert.Equal(t, tt.want, m.String())
		})
	}
}
