package main

import (
	"testing"
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
			if gotReturnCode := RunCmd(tt.args.cmd, tt.args.env); gotReturnCode != tt.wantReturnCode {
				t.Errorf("RunCmd() = %v, want %v", gotReturnCode, tt.wantReturnCode)
			}
		})
	}
}
