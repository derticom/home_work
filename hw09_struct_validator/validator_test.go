package hw09structvalidator

import (
	"encoding/json"
	"testing"

	//nolint:depguard // Применение 'require' необходимо для тестирования.
	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		in      interface{}
		wantErr error
	}{
		{
			name: "Ok case - User",
			in: User{
				ID:     "163746583764hdt38495483u2hdtyrb3748f",
				Name:   "Mike",
				Age:    35,
				Email:  "mile@gmail.com",
				Role:   "admin",
				Phones: []string{"12345678901", "77345678901"},
			},
			wantErr: ValidationErrors{},
		},
		{
			name: "Age older then max - User",
			in: User{
				ID:     "163746583764hdt38495483u2hdtyrb3748f",
				Name:   "Mike",
				Age:    55,
				Email:  "mile@gmail.com",
				Role:   "admin",
				Phones: []string{"12345678901", "77345678901"},
			},
			wantErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   errBiggerThanMax,
				},
			},
		},
		{
			name: "All incorrect - User",
			in: User{
				ID:     "123",
				Name:   "Mike",
				Age:    88,
				Email:  "gmail.com",
				Role:   "user",
				Phones: []string{"123456701"},
			},
			wantErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   errLengthNotEqual,
				},
				ValidationError{
					Field: "Age",
					Err:   errBiggerThanMax,
				},
				ValidationError{
					Field: "Email",
					Err:   errRegexpMismatch,
				},
				ValidationError{
					Field: "Role",
					Err:   errValueNotInSet,
				},
				ValidationError{
					Field: "Phones",
					Err:   errLengthNotEqual,
				},
			},
		},
		{
			name: "Ok case - App",
			in: App{
				Version: "12345",
			},
			wantErr: ValidationErrors{},
		},
		{
			name: "Small len - App",
			in: App{
				Version: "1",
			},
			wantErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   errLengthNotEqual,
				},
			},
		},
		{
			name: "Ok case - Response",
			in: Response{
				Code: 404,
				Body: "Smth",
			},
			wantErr: ValidationErrors{},
		},
		{
			name: "Not in - Response",
			in: Response{
				Code: 999,
			},
			wantErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   errValueNotInSet,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.in)
			require.Equal(t, tt.wantErr, err)
		})
	}
}
