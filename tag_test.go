package valex

import (
	"net"
	neturl "net/url"
	"strings"
	"testing"
	"time"
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
				Age  int    `val:"rangeint,min=0,max=120"`
				Name string `val:"min,size=3"`
			}{Age: 30, Name: "John"},
			wantValid: true,
		},
		{
			name: "Invalid int range",
			data: &struct {
				Age int `val:"rangeint,min=0,max=120"`
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
			errSubstr: "directive processing field \"Name\"",
		},
		{
			name: "Malformed parameter for string validator",
			data: &struct {
				Name string `val:"min,"`
			}{Name: "Alice"},
			wantValid: false,
			errSubstr: "parameter processing field \"Name\"",
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
			errSubstr: "unknown directive \"length\"",
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

func TestValidateStruct_float64(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		wantValid bool
		errSubstr string
	}{
		{
			name: "Valid float64 range",
			data: &struct {
				Score float64 `val:"rangefloat,min=0.5,max=1.5"`
			}{Score: 1.0},
			wantValid: true,
		},
		{
			name: "Invalid float64 min",
			data: &struct {
				Score float64 `val:"minfloat,min=0.5"`
			}{Score: 0.25},
			wantValid: false,
			errSubstr: "less than minimum",
		},
		{
			name: "Valid float64 oneof",
			data: &struct {
				Score float64 `val:"oneoffloat,values=1.5|2.5|3.5"`
			}{Score: 2.5},
			wantValid: true,
		},
		{
			name: "Invalid float64 oneof value",
			data: &struct {
				Score float64 `val:"oneoffloat,values=1.5|bad"`
			}{Score: 1.5},
			wantValid: false,
			errSubstr: "invalid float",
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

func TestValidateStruct_timeDurationIPURL(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		wantValid bool
		errSubstr string
	}{
		{
			name: "Valid time before",
			data: &struct {
				When time.Time `val:"beforetime,before=2024-01-02T03:04:05Z"`
			}{When: time.Date(2024, 1, 2, 3, 4, 4, 0, time.UTC)},
			wantValid: true,
		},
		{
			name: "Invalid time after",
			data: &struct {
				When time.Time `val:"aftertime,after=2024-01-02T03:04:05Z"`
			}{When: time.Date(2024, 1, 2, 3, 4, 4, 0, time.UTC)},
			wantValid: false,
			errSubstr: "not after",
		},
		{
			name: "Valid time between",
			data: &struct {
				When time.Time `val:"betweentime,start=2024-01-02T00:00:00Z,end=2024-01-03T00:00:00Z"`
			}{When: time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC)},
			wantValid: true,
		},
		{
			name: "Invalid time between",
			data: &struct {
				When time.Time `val:"betweentime,start=2024-01-02T00:00:00Z,end=2024-01-03T00:00:00Z"`
			}{When: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC)},
			wantValid: false,
			errSubstr: "not in range",
		},
		{
			name: "Valid positive duration",
			data: &struct {
				Delay time.Duration `val:"posduration"`
			}{Delay: time.Second},
			wantValid: true,
		},
		{
			name: "Invalid positive duration",
			data: &struct {
				Delay time.Duration `val:"posduration"`
			}{Delay: 0},
			wantValid: false,
			errSubstr: "not positive",
		},
		{
			name: "Invalid non-zero duration",
			data: &struct {
				Delay time.Duration `val:"!zeroduration"`
			}{Delay: 0},
			wantValid: false,
			errSubstr: "duration is zero",
		},
		{
			name: "Valid non-zero IP",
			data: &struct {
				Addr net.IP `val:"!zeroip"`
			}{Addr: net.ParseIP("192.168.1.1")},
			wantValid: true,
		},
		{
			name: "Invalid non-zero IP",
			data: &struct {
				Addr net.IP `val:"!zeroip"`
			}{Addr: net.ParseIP("0.0.0.0")},
			wantValid: false,
			errSubstr: "ip is zero",
		},
		{
			name: "Valid IP range",
			data: &struct {
				Addr net.IP `val:"iprange,start=192.168.1.10,end=192.168.1.20"`
			}{Addr: net.ParseIP("192.168.1.15")},
			wantValid: true,
		},
		{
			name: "Invalid IP range",
			data: &struct {
				Addr net.IP `val:"iprange,start=192.168.1.10,end=192.168.1.20"`
			}{Addr: net.ParseIP("192.168.1.30")},
			wantValid: false,
			errSubstr: "not in range",
		},
		{
			name: "Valid non-zero URL",
			data: &struct {
				Addr neturl.URL `val:"!zerourl"`
			}{Addr: neturl.URL{Host: "example.com"}},
			wantValid: true,
		},
		{
			name: "Invalid non-zero URL",
			data: &struct {
				Addr neturl.URL `val:"!zerourl"`
			}{Addr: neturl.URL{}},
			wantValid: false,
			errSubstr: "url is zero",
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
