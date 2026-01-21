package valex

import (
	"strings"
	"testing"
)

func TestValidateStruct_int(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		wantValid bool
		errSubstr string
	}{
		{
			name: "Valid int range and min length string",
			data: &struct {
				Age  int    `val:"range,min=0,max=120"`
				Name string `val:"min,size=3"`
			}{Age: 30, Name: "John"},
			wantValid: true,
		},
		{
			name: "Invalid int range",
			data: &struct {
				Age int `val:"range,min=0,max=120"`
			}{Age: -1},
			wantValid: false,
			errSubstr: "out of range",
		},
		{
			name: "Unknown directive id",
			data: &struct {
				Field int `val:"foobar"`
			}{Field: 10},
			wantValid: false,
			errSubstr: "unknown directive \"foobar\"",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			valid, err := ValidateStruct(tc.data)
			if valid != tc.wantValid {
				t.Errorf("expected valid=%v, got %v (error: %v)", tc.wantValid, valid, err)
			}
			if !tc.wantValid && err != nil && tc.errSubstr != "" {
				if !strings.Contains(err.Error(), tc.errSubstr) {
					t.Errorf("expected error to contain %q, got %q", tc.errSubstr, err.Error())
				}
			}
		})
	}
}

func TestValidateStruct_string(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		wantValid bool
		errSubstr string
	}{
		{
			name: "Invalid string min length",
			data: &struct {
				Name string `val:"min,size=3"`
			}{Name: "Al"},
			wantValid: false,
			errSubstr: "error processing field \"Name\"",
		},
		{
			name: "Malformed parameter for string validator",
			data: &struct {
				Name string `val:"min,"`
			}{Name: "Alice"},
			wantValid: false,
			errSubstr: "error processing field",
		},
		{
			name: "Valid length range for string",
			data: &struct {
				Code string `val:"len,min=3,max=5"`
			}{Code: "abcd"},
			wantValid: true,
		},
		{
			name: "Invalid length range (too short)",
			data: &struct {
				Code string `val:"length,min=3,max=5"`
			}{Code: "ab"},
			wantValid: false,
			errSubstr: "error processing field \"Code\"",
		},
		{
			name: "Valid regex match",
			data: &struct {
				Code string `val:"regex,pattern=^\\d+$"`
			}{Code: "12345"},
			wantValid: true,
		},
		{
			name: "Invalid regex match",
			data: &struct {
				Code string `val:"regex,pattern=^\\d+$"`
			}{Code: "abc"},
			wantValid: false,
			errSubstr: "does not match pattern",
		},
		{
			name: "Invalid regex pattern",
			data: &struct {
				Code string `val:"regex,pattern=["`
			}{Code: "123"},
			wantValid: false,
			errSubstr: "invalid regex pattern",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			valid, err := ValidateStruct(tc.data)
			if valid != tc.wantValid {
				t.Errorf("expected valid=%v, got %v (error: %v)", tc.wantValid, valid, err)
			}
			if !tc.wantValid && err != nil && tc.errSubstr != "" {
				if !strings.Contains(err.Error(), tc.errSubstr) {
					t.Errorf("expected error to contain %q, got %q", tc.errSubstr, err.Error())
				}
			}
		})
	}
}
