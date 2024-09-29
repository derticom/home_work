package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"  //nolint: depguard // import is necessary
	"github.com/stretchr/testify/require" //nolint: depguard // import is necessary
)

const dirPermission = 0o755

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
		wantFile string
		wantErr  error
	}{
		{
			name: "offset 0, limit 0",
			args: args{
				fromPath: fromPath,
				toPath:   toPath,
				offset:   0,
				limit:    0,
			},
			wantFile: "testdata/out_offset0_limit0.txt",
			wantErr:  nil,
		},
		{
			name: "offset 0, limit 10",
			args: args{
				fromPath: fromPath,
				toPath:   toPath,
				offset:   0,
				limit:    10,
			},
			wantFile: "testdata/out_offset0_limit10.txt",
			wantErr:  nil,
		},
		{
			name: "offset 0, limit 1000",
			args: args{
				fromPath: fromPath,
				toPath:   toPath,
				offset:   0,
				limit:    1000,
			},
			wantFile: "testdata/out_offset0_limit1000.txt",
			wantErr:  nil,
		},
		{
			name: "offset 0, limit 10000",
			args: args{
				fromPath: fromPath,
				toPath:   toPath,
				offset:   0,
				limit:    10000,
			},
			wantFile: "testdata/out_offset0_limit10000.txt",
			wantErr:  nil,
		},
		{
			name: "offset 100, limit 1000",
			args: args{
				fromPath: fromPath,
				toPath:   toPath,
				offset:   100,
				limit:    1000,
			},
			wantFile: "testdata/out_offset100_limit1000.txt",
			wantErr:  nil,
		},
		{
			name: "offset 6000, limit 1000",
			args: args{
				fromPath: fromPath,
				toPath:   toPath,
				offset:   6000,
				limit:    1000,
			},
			wantFile: "testdata/out_offset6000_limit1000.txt",
			wantErr:  nil,
		},
		{
			name: "error - the same path",
			args: args{
				fromPath: fromPath,
				toPath:   fromPath,
				offset:   0,
				limit:    0,
			},
			wantErr: ErrTheSamePath,
		},
		{
			name: "error - invalid limit",
			args: args{
				fromPath: fromPath,
				toPath:   toPath,
				offset:   0,
				limit:    -1000,
			},
			wantErr: ErrInvalidInput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.Mkdir("./tmp", dirPermission)
			require.NoError(t, err)

			defer func() {
				err = os.RemoveAll("tmp")
				require.NoError(t, err)
			}()

			err = Copy(tt.args.fromPath, tt.args.toPath, tt.args.offset, tt.args.limit)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)

			gotFile, err := os.ReadFile(tt.args.toPath)
			require.NoError(t, err)

			wantFile, err := os.ReadFile(tt.wantFile)
			require.NoError(t, err)

			assert.Equal(t, wantFile, gotFile)
		})
	}
}
