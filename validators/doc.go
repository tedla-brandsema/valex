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
//	err := valex.ValidateStruct(&User{Email: "a@b.com", Age: 30})
//
// Parameters are supplied in the tag after the directive name
// (for example "rangeint,min=0,max=120"). Where a directive takes a list, the
// values are pipe-separated (for example "oneof,values=a|b|c").
//
// # Catalog
//
// The "val" tag directives, grouped by the Go type of the field they validate.
// The Registers column names the type to pass to valex.RegisterDirective.
//
//	Tag            Registers                     Params       Description
//	-------------- ----------------------------- ------------ ------------------------------------
//	-- int --
//	rangeint       IntRangeValidator             min, max     inclusive range
//	minint         MinIntValidator               min          value >= min
//	maxint         MaxIntValidator               max          value <= max
//	posint         NonNegativeIntValidator       -            non-negative
//	negint         NonPositiveIntValidator       -            non-positive
//	!zeroint       NonZeroIntValidator           -            not zero
//	oneofint       OneOfIntValidator             values       one of a pipe-separated list
//	-- float64 --
//	rangefloat     Float64RangeValidator         min, max     inclusive range
//	minfloat       MinFloat64Validator           min          value >= min
//	maxfloat       MaxFloat64Validator           max          value <= max
//	posfloat       NonNegativeFloat64Validator   -            non-negative
//	negfloat       NonPositiveFloat64Validator   -            non-positive
//	!zerofloat     NonZeroFloat64Validator       -            not zero
//	oneoffloat     OneOfFloat64Validator         values       one of a pipe-separated list
//	-- string --
//	url            UrlValidator                  -            valid absolute URL
//	email          EmailValidator                -            valid email address
//	!empty         NonEmptyStringValidator       -            non-empty
//	min            MinLengthValidator            size         length >= size
//	max            MaxLengthValidator            size         length <= size
//	len            LengthRangeValidator          min, max     length within [min, max]
//	regex          RegexValidator                pattern      matches regular expression
//	prefix         PrefixValidator               value        has prefix
//	suffix         SuffixValidator               value        has suffix
//	contains       ContainsValidator             value        contains substring
//	oneof          OneOfStringValidator          values       one of a pipe-separated list
//	alphanum       AlphaNumericValidator         -            alphanumeric
//	mac            MACAddressValidator           -            valid MAC address
//	ip             IpValidator                   -            valid IP address
//	ipv4           IPv4Validator                 -            valid IPv4 address
//	ipv6           IPv6Validator                 -            valid IPv6 address
//	hostname       HostnameValidator             -            valid hostname
//	cidr           IPCIDRValidator               -            valid CIDR notation
//	uuid           UUIDValidator                 version (4)  RFC 4122 UUID, optional version
//	base64         Base64Validator               -            valid base64 (standard or raw)
//	hex            HexValidator                  -            valid hex (optional 0x prefix)
//	xml            XMLValidator                  -            well-formed XML
//	json           JSONValidator                 -            valid JSON
//	time           TimeValidator                 format       valid time for the layout (default RFC3339)
//	-- time.Time --
//	!zerotime      NonZeroTimeValidator          -            not zero
//	beforetime     TimeBeforeValidator           before       before the given RFC3339 time
//	aftertime      TimeAfterValidator            after        after the given RFC3339 time
//	betweentime    TimeBetweenValidator          start, end   within [start, end] (RFC3339)
//	-- time.Duration --
//	posduration    PositiveDurationValidator     -            positive
//	!zeroduration  NonZeroDurationValidator      -            not zero
//	-- net.IP --
//	!zeroip        NonZeroIPValidator            -            not zero/unspecified
//	iprange        IPRangeValidator              start, end   within [start, end]
//	-- url.URL --
//	!zerourl       NonZeroURLValidator           -            not the zero value
//
// Alongside the tag directives, the package also offers generic programmatic
// validators that are not registered with the "val" tag: CmpRangeValidator and
// NonZeroValidator implement valex.Validator directly, and CompositeValidator
// chains several valex.Validator values into one.
package validators
