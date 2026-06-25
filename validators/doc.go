// Package validators provides a catalog of ready-made validation directives for
// the valex engine.
//
// Directives are opt-in: importing this package registers nothing on its own.
// Register the ones you need against valex's "val" tag with
// valex.RegisterDirective, then validate with valex.ValidateStruct (or the
// valex/forms helpers):
//
//	valex.RegisterDirective(&validators.EmailValidator{})
//	valex.RegisterDirective(&validators.IntRangeValidator{})
//
//	type User struct {
//		Email string `val:"email"`
//		Age   int    `val:"rangeint,min=0,max=120"`
//	}
//
//	ok, err := valex.ValidateStruct(&User{Email: "a@b.com", Age: 30})
//
// Parameters are supplied in the tag after the directive name
// (for example "rangeint,min=0,max=120"). Where a directive takes a list, the
// values are pipe-separated (for example "oneof,values=a|b|c"). Directives whose
// names begin with "!" also have a plain-word alias, noted in the catalog below.
//
// # Catalog
//
// The "val" tag directives, grouped by the Go type of the field they validate:
//
//	Tag            Type            Params       Description
//	-------------- --------------- ------------ -------------------------------------
//	-- int --
//	rangeint       int             min, max     inclusive range
//	minint         int             min          value >= min
//	maxint         int             max          value <= max
//	posint         int             -            non-negative
//	negint         int             -            non-positive
//	!zeroint       int             -            not zero (alias: nonzeroint)
//	oneofint       int             values       one of a pipe-separated list
//	-- float64 --
//	rangefloat     float64         min, max     inclusive range
//	minfloat       float64         min          value >= min
//	maxfloat       float64         max          value <= max
//	posfloat       float64         -            non-negative
//	negfloat       float64         -            non-positive
//	!zerofloat     float64         -            not zero (alias: nonzerofloat)
//	oneoffloat     float64         values       one of a pipe-separated list
//	-- string --
//	url            string          -            valid absolute URL
//	email          string          -            valid email address
//	!empty         string          -            non-empty (alias: nonempty)
//	min            string          size         length >= size
//	max            string          size         length <= size
//	len            string          min, max     length within [min, max]
//	regex          string          pattern      matches regular expression
//	prefix         string          value        has prefix
//	suffix         string          value        has suffix
//	contains       string          value        contains substring
//	oneof          string          values       one of a pipe-separated list
//	alphanum       string          -            alphanumeric
//	mac            string          -            valid MAC address
//	ip             string          -            valid IP address
//	ipv4           string          -            valid IPv4 address
//	ipv6           string          -            valid IPv6 address
//	hostname       string          -            valid hostname
//	cidr           string          -            valid CIDR notation
//	uuid           string          version (4)  RFC 4122 UUID, optional version
//	base64         string          -            valid base64 (standard or raw)
//	hex            string          -            valid hex (optional 0x prefix)
//	xml            string          -            well-formed XML
//	json           string          -            valid JSON
//	time           string          format       valid time for the layout (default RFC3339)
//	-- time.Time --
//	!zerotime      time.Time       -            not zero (alias: nonzerotime)
//	beforetime     time.Time       before       before the given RFC3339 time
//	aftertime      time.Time       after        after the given RFC3339 time
//	betweentime    time.Time       start, end   within [start, end] (RFC3339)
//	-- time.Duration --
//	posduration    time.Duration   -            positive
//	!zeroduration  time.Duration   -            not zero (alias: nonzeroduration)
//	-- net.IP --
//	!zeroip        net.IP          -            not zero/unspecified (alias: nonzeroip)
//	iprange        net.IP          start, end   within [start, end]
//	-- url.URL --
//	!zerourl       url.URL         -            not the zero value (alias: nonzerourl)
//
// Alongside the tag directives, the package also offers generic programmatic
// validators that are not registered with the "val" tag: CmpRangeValidator and
// NonZeroValidator implement valex.Validator directly, and CompositeValidator
// chains several valex.Validator values into one.
package validators
