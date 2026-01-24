package valex

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/tedla-brandsema/tagex"
)

func TestIntRangeValidator(t *testing.T) {
	v := &IntRangeValidator{Min: 10, Max: 20}
	tests := []struct {
		input int
		ok    bool
	}{
		{15, true},
		{10, true},
		{20, true},
		{9, false},
		{21, false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%d): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestFloat64RangeValidator(t *testing.T) {
	v := &Float64RangeValidator{Min: 1.5, Max: 3.5}
	tests := []struct {
		input float64
		ok    bool
	}{
		{1.5, true},
		{2.1, true},
		{3.5, true},
		{1.4, false},
		{3.6, false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%g): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestNonZeroValidator(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		v := &NonZeroValidator[int]{}
		tests := []struct {
			input int
			ok    bool
		}{
			{0, false},
			{1, true},
			{-1, true},
		}
		for _, tc := range tests {
			ok, err := v.Validate(tc.input)
			if ok != tc.ok {
				t.Errorf("%T(%d): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
			}
		}
	})

	t.Run("string", func(t *testing.T) {
		v := &NonZeroValidator[string]{}
		tests := []struct {
			input string
			ok    bool
		}{
			{"", false},
			{"value", true},
		}
		for _, tc := range tests {
			ok, err := v.Validate(tc.input)
			if ok != tc.ok {
				t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
			}
		}
	})

	t.Run("struct", func(t *testing.T) {
		type sample struct {
			ID   int
			Name string
		}
		v := &NonZeroValidator[sample]{}
		tests := []struct {
			input sample
			ok    bool
		}{
			{sample{}, false},
			{sample{ID: 1}, true},
			{sample{Name: "name"}, true},
		}
		for _, tc := range tests {
			ok, err := v.Validate(tc.input)
			if ok != tc.ok {
				t.Errorf("%T(%v): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
			}
		}
	})

	t.Run("pointer", func(t *testing.T) {
		v := &NonZeroValidator[*int]{}
		var zero *int
		value := 5
		tests := []struct {
			input *int
			ok    bool
		}{
			{zero, false},
			{&value, true},
		}
		for _, tc := range tests {
			ok, err := v.Validate(tc.input)
			if ok != tc.ok {
				t.Errorf("%T(%v): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
			}
		}
	})

	t.Run("nil interface", func(t *testing.T) {
		v := &NonZeroValidator[any]{}
		tests := []struct {
			input any
			ok    bool
		}{
			{nil, false},
			{0, false},
			{"", false},
			{"ok", true},
		}
		for _, tc := range tests {
			ok, err := v.Validate(tc.input)
			if ok != tc.ok {
				t.Errorf("%T(%v): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
			}
		}
	})
}

func TestNonNegativeIntValidator(t *testing.T) {
	v := &NonNegativeIntValidator{}
	tests := []struct {
		input int
		ok    bool
	}{
		{-1, false},
		{0, true},
		{1, true},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%d): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestNonNegativeFloat64Validator(t *testing.T) {
	v := &NonNegativeFloat64Validator{}
	tests := []struct {
		input float64
		ok    bool
	}{
		{-1.5, false},
		{0, true},
		{1.5, true},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%g): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestNonPositiveIntValidator(t *testing.T) {
	v := &NonPositiveIntValidator{}
	tests := []struct {
		input int
		ok    bool
	}{
		{-1, true},
		{0, true},
		{1, false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%d): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestNonPositiveFloat64Validator(t *testing.T) {
	v := &NonPositiveFloat64Validator{}
	tests := []struct {
		input float64
		ok    bool
	}{
		{-1.5, true},
		{0, true},
		{1.5, false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%g): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestUrlValidator(t *testing.T) {
	v := &UrlValidator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{"https://www.example.com", true},
		{"ftp://example.com", true}, // Acceptable as a valid URL scheme.
		{"invalid-url", false},
		{"", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestEmailValidator(t *testing.T) {
	v := &EmailValidator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{"user@example.com", true},
		{"invalid-email", false},
		{"", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestNonEmptyStringValidator(t *testing.T) {
	v := &NonEmptyStringValidator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{"hello", true},
		{"", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestMinLengthValidator(t *testing.T) {
	v := &MinLengthValidator{Size: 3}
	tests := []struct {
		input string
		ok    bool
	}{
		{"abc", true},
		{"abcd", true},
		{"ab", false},
		{"", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestMaxLengthValidator(t *testing.T) {
	v := &MaxLengthValidator{Size: 3}
	tests := []struct {
		input string
		ok    bool
	}{
		{"abc", true},
		{"ab", true},
		{"abcd", false},
		{"", true},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestLengthRangeValidator(t *testing.T) {
	v := &LengthRangeValidator{Min: 3, Max: 6}
	tests := []struct {
		input string
		ok    bool
	}{
		{"abcd", true},
		{"ab", false},
		{"abcdefg", false},
		{"abc", true},    // exactly at minimum
		{"abcdef", true}, // exactly at maximum
	}
	for _, tt := range tests {
		ok, err := v.Validate(tt.input)
		if ok != tt.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v, err=%v", *v, tt.input, tt.ok, ok, err)
		}
	}
}

func TestLengthRangeValidatorBounds(t *testing.T) {
	tests := []struct {
		name string
		v    *LengthRangeValidator
	}{
		{name: "min greater than max", v: &LengthRangeValidator{Min: 5, Max: 2}},
		{name: "min negative", v: &LengthRangeValidator{Min: -1, Max: 2}},
		{name: "max negative", v: &LengthRangeValidator{Min: 1, Max: -2}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := tc.v.Validate("abc")
			if ok || err == nil {
				t.Fatalf("expected bounds error, got ok=%v err=%v", ok, err)
			}
		})
	}
}

func TestMinLengthValidatorNegative(t *testing.T) {
	v := &MinLengthValidator{Size: -1}
	ok, err := v.Validate("abc")
	if ok || err == nil {
		t.Fatalf("expected negative size error, got ok=%v err=%v", ok, err)
	}
}

func TestMaxLengthValidatorNegative(t *testing.T) {
	v := &MaxLengthValidator{Size: -1}
	ok, err := v.Validate("abc")
	if ok || err == nil {
		t.Fatalf("expected negative size error, got ok=%v err=%v", ok, err)
	}
}

func TestRegexValidator(t *testing.T) {
	pattern := regexp.MustCompile(`^\d+$`)
	v := &RegexValidator{Pattern: pattern}
	tests := []struct {
		input string
		ok    bool
	}{
		{"123", true},
		{"abc", false},
		{"123abc", false},
		{"", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestMinIntValidator(t *testing.T) {
	v := &MinIntValidator{Min: 10}
	tests := []struct {
		input int
		ok    bool
	}{
		{10, true},
		{11, true},
		{9, false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%d): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestMinFloat64Validator(t *testing.T) {
	v := &MinFloat64Validator{Min: 1.5}
	tests := []struct {
		input float64
		ok    bool
	}{
		{1.5, true},
		{1.6, true},
		{1.4, false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%g): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestMaxIntValidator(t *testing.T) {
	v := &MaxIntValidator{Max: 10}
	tests := []struct {
		input int
		ok    bool
	}{
		{10, true},
		{9, true},
		{11, false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%d): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestMaxFloat64Validator(t *testing.T) {
	v := &MaxFloat64Validator{Max: 1.5}
	tests := []struct {
		input float64
		ok    bool
	}{
		{1.5, true},
		{1.4, true},
		{1.6, false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%g): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestNonZeroIntValidator(t *testing.T) {
	v := &NonZeroIntValidator{}
	tests := []struct {
		input int
		ok    bool
	}{
		{0, false},
		{1, true},
		{-1, true},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%d): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestNonZeroFloat64Validator(t *testing.T) {
	v := &NonZeroFloat64Validator{}
	tests := []struct {
		input float64
		ok    bool
	}{
		{0, false},
		{1.5, true},
		{-1.5, true},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%g): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestNonZeroTimeValidator(t *testing.T) {
	v := &NonZeroTimeValidator{}
	tests := []struct {
		input time.Time
		ok    bool
	}{
		{time.Time{}, false},
		{time.Now(), true},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%v): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestOneOfStringValidator(t *testing.T) {
	v := &OneOfStringValidator{Values: []string{"red", "green", "blue"}}
	tests := []struct {
		input string
		ok    bool
	}{
		{"red", true},
		{"green", true},
		{"yellow", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestOneOfIntValidator(t *testing.T) {
	v := &OneOfIntValidator{Values: []int{1, 3, 5}}
	tests := []struct {
		input int
		ok    bool
	}{
		{1, true},
		{2, false},
		{5, true},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%d): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestOneOfFloat64Validator(t *testing.T) {
	v := &OneOfFloat64Validator{Values: []float64{1.5, 3.5, 5.5}}
	tests := []struct {
		input float64
		ok    bool
	}{
		{1.5, true},
		{2.5, false},
		{5.5, true},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%g): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestPrefixValidator(t *testing.T) {
	v := &PrefixValidator{Value: "pre"}
	tests := []struct {
		input string
		ok    bool
	}{
		{"prefix", true},
		{"nopre", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestSuffixValidator(t *testing.T) {
	v := &SuffixValidator{Value: "fix"}
	tests := []struct {
		input string
		ok    bool
	}{
		{"suffix", true},
		{"fixes", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestContainsValidator(t *testing.T) {
	v := &ContainsValidator{Value: "mid"}
	tests := []struct {
		input string
		ok    bool
	}{
		{"amidb", true},
		{"nomatch", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestUUIDValidator(t *testing.T) {
	v := &UUIDValidator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", true},  // v4
		{"6ba7b810-9dad-11d1-80b4-00c04fd430c8", false}, // v1 default expects v4
		{"550e8400e29b41d4a716446655440000", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestUUIDValidatorVersion(t *testing.T) {
	v := &UUIDValidator{Version: 1}
	tests := []struct {
		input string
		ok    bool
	}{
		{"6ba7b810-9dad-11d1-80b4-00c04fd430c8", true},  // v1
		{"550e8400-e29b-41d4-a716-446655440000", false}, // v4
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestHostnameValidator(t *testing.T) {
	v := &HostnameValidator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{"example.com", true},
		{"localhost", true},
		{"http://example.com", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestIPCIDRValidator(t *testing.T) {
	v := &IPCIDRValidator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{"192.168.0.0/24", true},
		{"2001:db8::/32", true},
		{"invalid", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestBase64Validator(t *testing.T) {
	v := &Base64Validator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{"aGVsbG8=", true},
		{"aGVsbG8", true},
		{"not-base64", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestHexValidator(t *testing.T) {
	v := &HexValidator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{"deadbeef", true},
		{"0xdeadbeef", true},
		{"xyz", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestTimeValidator(t *testing.T) {
	v := &TimeValidator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{"2020-01-02T03:04:05Z", true},
		{"2020-01-02", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestTimeValidatorFormat(t *testing.T) {
	v := &TimeValidator{Format: "2006-01-02"}
	tests := []struct {
		input string
		ok    bool
	}{
		{"2020-01-02", true},
		{"2020-01-02T03:04:05Z", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

type evenDirectiveTest struct{}

func (d *evenDirectiveTest) Name() string {
	return "even"
}

func (d *evenDirectiveTest) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

func (d *evenDirectiveTest) Handle(val int) (int, error) {
	if val%2 != 0 {
		return val, fmt.Errorf("value %d is not even", val)
	}
	return val, nil
}

func TestRegisterDirective(t *testing.T) {
	RegisterDirective(&evenDirectiveTest{})

	tests := []struct {
		name      string
		data      interface{}
		wantValid bool
		errSubstr string
	}{
		{
			name: "Valid even value",
			data: &struct {
				Count int `val:"even"`
			}{Count: 4},
			wantValid: true,
		},
		{
			name: "Invalid odd value",
			data: &struct {
				Count int `val:"even"`
			}{Count: 3},
			wantValid: false,
			errSubstr: "not even",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := ValidateStruct(tt.data)
			if ok != tt.wantValid {
				t.Fatalf("expected ok=%v, got ok=%v (err: %v)", tt.wantValid, ok, err)
			}
			if !tt.wantValid && err != nil && tt.errSubstr != "" {
				if !strings.Contains(err.Error(), tt.errSubstr) {
					t.Fatalf("expected error to contain %q, got %q", tt.errSubstr, err.Error())
				}
			}
		})
	}
}

func TestValidatedValue(t *testing.T) {
	t.Run("set and get success", func(t *testing.T) {
		v := &NonNegativeIntValidator{}
		val := &ValidatedValue[int]{Validator: v}
		if err := val.Set(10); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got := val.Get(); got != 10 {
			t.Fatalf("expected value=10, got %d", got)
		}
		if got := val.String(); got != "10" {
			t.Fatalf("expected String()=10, got %q", got)
		}
	})

	t.Run("set with nil validator", func(t *testing.T) {
		val := &ValidatedValue[int]{}
		if err := val.Set(1); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("set with invalid value", func(t *testing.T) {
		v := &NonPositiveIntValidator{}
		val := &ValidatedValue[int]{Validator: v}
		err := val.Set(1)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "positive integer") {
			t.Fatalf("expected error to mention positive integer, got %q", err.Error())
		}
	})
}

func TestMustValidate(t *testing.T) {
	t.Run("returns value on success", func(t *testing.T) {
		v := &NonNegativeIntValidator{}
		got := MustValidate(5, v)
		if got != 5 {
			t.Fatalf("expected 5, got %d", got)
		}
	})

	t.Run("panics on failure", func(t *testing.T) {
		v := &NonNegativeIntValidator{}
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic, got nil")
			}
		}()
		_ = MustValidate(-1, v)
	})
}

func TestAlphaNumericValidator(t *testing.T) {
	v := &AlphaNumericValidator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{"abc123", true},
		{"ABC", true},
		{"abc 123", false},
		{"abc-123", false},
		{"", false},
	}
	for _, tt := range tests {
		ok, err := v.Validate(tt.input)
		if ok != tt.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v, err=%v", *v, tt.input, tt.ok, ok, err)
		}
	}
}

func TestMACAddressValidator(t *testing.T) {
	v := &MACAddressValidator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{"00:00:5e:00:53:01", true},
		{"02:00:5e:10:00:00:00:01", true},
		{"00:00:00:00:fe:80:00:00:00:00:00:00:02:00:5e:10:00:00:00:01", true},
		{"00-00-5e-00-53-01", true},
		{"02-00-5e-10-00-00-00-01", true},
		{"00-00-00-00-fe-80-00-00-00-00-00-00-02-00-5e-10-00-00-00-01", true},
		{"0000.5e00.5301", true},
		{"0200.5e10.0000.0001", true},
		{"0000.0000.fe80.0000.0000.0000.0200.5e10.0000.0001", true}, {"01:23:45:67:89:ab", true},
		{"01-23-45-67-89-ab", true},
		{"0123456789ab", false},
		{"invalid-mac", false},
	}
	for _, tt := range tests {
		ok, err := v.Validate(tt.input)
		if ok != tt.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v, err=%v", *v, tt.input, tt.ok, ok, err)
		}
	}
}

func TestIpValidator(t *testing.T) {
	v := &IpValidator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{"192.168.1.1", true},    // valid IPv4
		{"2001:0db8::1", true},   // valid IPv6
		{"invalid-ip", false},    // invalid IP
		{"123.456.789.0", false}, // invalid IPv4
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v (err: %v)", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestIPv4Validator(t *testing.T) {
	v := &IPv4Validator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{"192.168.0.1", true},
		{"2001:db8::1", false},
		{"invalid", false},
	}
	for _, tt := range tests {
		ok, err := v.Validate(tt.input)
		if ok != tt.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v, err=%v", *v, tt.input, tt.ok, ok, err)
		}
	}
}

func TestIPv6Validator(t *testing.T) {
	v := &IPv6Validator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{"2001:db8::1", true},
		{"192.168.0.1", false},
		{"invalid", false},
	}
	for _, tt := range tests {
		ok, err := v.Validate(tt.input)
		if ok != tt.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v, err=%v", *v, tt.input, tt.ok, ok, err)
		}
	}
}

func TestXMLValidator(t *testing.T) {
	v := &XMLValidator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{
			input: `<note>
					<to>Tove</to>
					<from>Jani</from>
					<heading>Reminder</heading>
					<body>Don't forget me this weekend!</body>
				</note>`,
			ok: true,
		},
		{
			input: `<root><child>value</child></root>`,
			ok:    true,
		},
		{
			input: `<root><child>value</child>`,
			ok:    false,
		},
		{
			input: "Just plain text",
			ok:    false,
		},
		{
			input: "",
			ok:    false,
		},
	}

	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v, error: %v", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestJSONValidator(t *testing.T) {
	v := &JSONValidator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{`{"key": "value"}`, true},
		{`[1, 2, 3]`, true},
		{`"a simple string"`, true},
		{`123`, true},
		{`true`, true},
		{`false`, true},
		{`null`, true},
		{"  { \"key\": 123 }  ", true}, // Whitespace is allowed.
		{``, false},
		{`invalid json`, false},
		{`{"key": "value" extra}`, false},
		{`[1, 2,]`, false},
	}

	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T(%q): expected ok=%v, got ok=%v, error: %v", *v, tc.input, tc.ok, ok, err)
		}
	}
}

func TestCompositeValidator_String(t *testing.T) {
	nonEmpty := &NonEmptyStringValidator{}
	minLength := &MinLengthValidator{Size: 3}
	composite := &CompositeValidator[string]{Validators: []Validator[string]{nonEmpty, minLength}}

	tests := []struct {
		input string
		ok    bool
	}{
		{"abc", true},
		{"ab", false}, // Fails min length
		{"", false},   // Fails non-empty
	}

	for _, tc := range tests {
		ok, err := composite.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T (string) for input %q: expected ok=%v, got ok=%v (err: %v)", *composite, tc.input, tc.ok, ok, err)
		}
	}
}

func TestCompositeValidator_Int(t *testing.T) {
	nonNegative := &NonNegativeIntValidator{}
	rangeValidator := &IntRangeValidator{Min: 0, Max: 100}
	composite := &CompositeValidator[int]{Validators: []Validator[int]{nonNegative, rangeValidator}}

	tests := []struct {
		input int
		ok    bool
	}{
		{-5, false}, // Fails non-negative check
		{50, true},
		{150, false}, // Fails range check
	}

	for _, tc := range tests {
		ok, err := composite.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("%T (int) for input %d: expected ok=%v, got ok=%v (err: %v)", *composite, tc.input, tc.ok, ok, err)
		}
	}
}
