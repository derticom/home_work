package main

import (
	"reflect"
	"testing"
)

func TestReadDir2(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    Environment
		wantErr bool
	}{
		{
			name: "ok case",
			args: args{
				dir: "./testdata/env",
			},
			want: Environment{
				"BAR": {
					Value:      "bar",
					NeedRemove: false,
				},
				"EMPTY": {
					Value:      "",
					NeedRemove: false,
				},
				"FOO": {
					Value:      "   foo\nwith new line",
					NeedRemove: false,
				},
				"HELLO": {
					Value:      "\"hello\"",
					NeedRemove: false,
				},
				"UNSET": {
					Value:      "",
					NeedRemove: true,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadDir(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}
