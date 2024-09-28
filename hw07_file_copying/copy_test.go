package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	fromPath = "./testdata/input.txt"
	toPath   = "./tmp/input.txt"
)

func TestCopy(t *testing.T) {
	type args struct {
		fromPath string
		toPath   string
		offset   int64
		limit    int64
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantFile string
	}{
		{
			name: "offset 0, limit 0",
			args: args{
				fromPath: fromPath,
				toPath:   toPath,
				offset:   0,
				limit:    0,
			},
			wantErr:  false,
			wantFile: "testdata/out_offset0_limit0.txt",
		},
		{
			name: "offset 0, limit 10",
			args: args{
				fromPath: fromPath,
				toPath:   toPath,
				offset:   0,
				limit:    10,
			},
			wantErr:  false,
			wantFile: "testdata/out_offset0_limit10.txt",
		},
		{
			name: "offset 0, limit 1000",
			args: args{
				fromPath: fromPath,
				toPath:   toPath,
				offset:   0,
				limit:    1000,
			},
			wantErr:  false,
			wantFile: "testdata/out_offset0_limit1000.txt",
		},
		{
			name: "offset 0, limit 10000",
			args: args{
				fromPath: fromPath,
				toPath:   toPath,
				offset:   0,
				limit:    10000,
			},
			wantErr:  false,
			wantFile: "testdata/out_offset0_limit10000.txt",
		},
		{
			name: "offset 100, limit 1000",
			args: args{
				fromPath: fromPath,
				toPath:   toPath,
				offset:   100,
				limit:    1000,
			},
			wantErr:  false,
			wantFile: "testdata/out_offset100_limit1000.txt",
		},
		{
			name: "offset 6000, limit 1000",
			args: args{
				fromPath: fromPath,
				toPath:   toPath,
				offset:   6000,
				limit:    1000,
			},
			wantErr:  false,
			wantFile: "testdata/out_offset6000_limit1000.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.Mkdir("./tmp", 0755)
			require.NoError(t, err)

			defer func() {
				err = os.RemoveAll("tmp")
				require.NoError(t, err)
			}()

			err = Copy(tt.args.fromPath, tt.args.toPath, tt.args.offset, tt.args.limit)
			assert.NoError(t, err)

			gotFile, err := os.ReadFile(tt.args.toPath)
			require.NoError(t, err)

			wantFile, err := os.ReadFile(tt.wantFile)
			require.NoError(t, err)

			assert.Equal(t, wantFile, gotFile)
		})
	}
}
