package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"  //nolint: depguard // import is necessary
	"github.com/stretchr/testify/require" //nolint: depguard // import is necessary
)

func TestRunCmd(t *testing.T) {
	type args struct {
		cmd []string
		env Environment
	}
	tests := []struct {
		name           string
		args           args
		wantReturnCode int
	}{
		{
			name: "ok case",
			args: args{
				cmd: []string{"echo", "hello world"},
				env: Environment{
					"key": {
						Value:      "value",
						NeedRemove: false,
					},
					"empty": {
						Value:      "",
						NeedRemove: true,
					},
				},
			},
			wantReturnCode: 0,
		},
		{
			name: "err case",
			args: args{
				cmd: []string{"false"},
				env: Environment{
					"key": {
						Value:      "value",
						NeedRemove: false,
					},
					"empty": {
						Value:      "",
						NeedRemove: true,
					},
				},
			},
			wantReturnCode: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.Setenv("key", "value-before")
			assert.NoError(t, err)

			err = os.Setenv("empty", "not-empty")
			assert.NoError(t, err)

			if gotReturnCode := RunCmd(tt.args.cmd, tt.args.env); gotReturnCode != tt.wantReturnCode {
				t.Errorf("RunCmd() = %v, want %v", gotReturnCode, tt.wantReturnCode)
			}

			envVar1 := os.Getenv("key")
			require.Equal(t, tt.args.env["key"].Value, envVar1)

			envVar2 := os.Getenv("empty")
			require.Equal(t, tt.args.env["empty"].Value, envVar2)
		})
	}
}
