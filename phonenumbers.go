package phonenumbers

import (
	"embed"
	"errors"
	"regexp"
	"strconv"
	"sync"
)

func init() {
	initMetadata()
}

func initMetadata() {
	_ = "STUB: not implemented"
	// load our regions
	return
}

// then our metadata

// We can assume that if the county calling code maps to the
// non-geo entity region code then that's the only region code
// it maps to.

// This is the subset of all country codes that map to the
// non-geo entity region code.

// The supported regions set does not include the "001"
// non-geo entity region code.

// If the non-geo entity still got added to the set of supported
// regions it must be because there are entries that list the non-geo
// entity alongside normal regions (which is wrong). If we discover
// this, remove the non-geo entity from the set of supported regions
// and log (or not log).

const (
	// MIN_LENGTH_FOR_NSN is the minimum and maximum length of the national significant number.
	MIN_LENGTH_FOR_NSN = 2
	// MAX_LENGTH_FOR_NSN: The ITU says the maximum length should be 15, but we have
	// found longer numbers in Germany.
	MAX_LENGTH_FOR_NSN = 17
	// MAX_LENGTH_COUNTRY_CODE is the maximum length of the country calling code.
	MAX_LENGTH_COUNTRY_CODE = 3
	// MAX_INPUT_STRING_LENGTH caps input strings for parsing at 250 chars.
	// This prevents malicious input from overflowing the regular-expression
	// engine.
	MAX_INPUT_STRING_LENGTH = 250

	// UNKNOWN_REGION is the region-code for the unknown region.
	UNKNOWN_REGION = "ZZ"

	NANPA_COUNTRY_CODE = 1

	// The PLUS_SIGN signifies the international prefix.
	PLUS_SIGN = '+'

	STAR_SIGN = '*'

	RFC3966_EXTN_PREFIX     = ";ext="
	RFC3966_PREFIX          = "tel:"
	RFC3966_PHONE_CONTEXT   = ";phone-context="
	RFC3966_ISDN_SUBADDRESS = ";isub="

	// Regular expression of acceptable punctuation found in phone
	// numbers. This excludes punctuation found as a leading character
	// only. This consists of dash characters, white space characters,
	// full stops, slashes, square brackets, parentheses and tildes. It
	// also includes the letter 'x' as that is found as a placeholder
	// for carrier information in some phone numbers. Full-width variants
	// are also present.
	VALID_PUNCTUATION = "-x\u2010-\u2015\u2212\u30FC\uFF0D-\uFF0F " +
		"\u00A0\u00AD\u200B\u2060\u3000()\uFF08\uFF09\uFF3B\uFF3D." +
		"\\[\\]/~\u2053\u223C\uFF5E"

	DIGITS = "\\p{Nd}"

	// We accept alpha characters in phone numbers, ASCII only, upper
	// and lower case.
	VALID_ALPHA = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	PLUS_CHARS  = "+\uFF0B"

	// This is defined by ICU as the unknown time zone.
	UNKNOWN_TIMEZONE = "Etc/Unknown"
)

var (
	// Map of country calling codes that use a mobile token before the
	// area code. One example of when this is relevant is when determining
	// the length of the national destination code, which should be the
	// length of the area code plus the length of the mobile token.
	MOBILE_TOKEN_MAPPINGS = map[int]string{
		54: "9",
	}

	// A map that contains characters that are essential when dialling.
	// That means any of the characters in this map must not be removed
	// from a number when dialling, otherwise the call will not reach
	// the intended destination.
	DIALLABLE_CHAR_MAPPINGS = map[rune]rune{
		'1':       '1',
		'2':       '2',
		'3':       '3',
		'4':       '4',
		'5':       '5',
		'6':       '6',
		'7':       '7',
		'8':       '8',
		'9':       '9',
		'0':       '0',
		PLUS_SIGN: PLUS_SIGN,
		'*':       '*',
		'#':       '#',
	}

	// Only upper-case variants of alpha characters are stored.
	ALPHA_MAPPINGS = map[rune]rune{
		'A': '2',
		'B': '2',
		'C': '2',
		'D': '3',
		'E': '3',
		'F': '3',
		'G': '4',
		'H': '4',
		'I': '4',
		'J': '5',
		'K': '5',
		'L': '5',
		'M': '6',
		'N': '6',
		'O': '6',
		'P': '7',
		'Q': '7',
		'R': '7',
		'S': '7',
		'T': '8',
		'U': '8',
		'V': '8',
		'W': '9',
		'X': '9',
		'Y': '9',
		'Z': '9',
	}

	// For performance reasons, amalgamate both into one map.
	ALPHA_PHONE_MAPPINGS = map[rune]rune{
		'0': '0',
		'1': '1',
		'2': '2',
		'3': '3',
		'4': '4',
		'5': '5',
		'6': '6',
		'7': '7',
		'8': '8',
		'9': '9',
		'A': '2',
		'B': '2',
		'C': '2',
		'D': '3',
		'E': '3',
		'F': '3',
		'G': '4',
		'H': '4',
		'I': '4',
		'J': '5',
		'K': '5',
		'L': '5',
		'M': '6',
		'N': '6',
		'O': '6',
		'P': '7',
		'Q': '7',
		'R': '7',
		'S': '7',
		'T': '8',
		'U': '8',
		'V': '8',
		'W': '9',
		'X': '9',
		'Y': '9',
		'Z': '9',
	}

	// Separate map of all symbols that we wish to retain when formatting
	// alpha numbers. This includes digits, ASCII letters and number
	// grouping symbols such as "-" and " ".
	ALL_PLUS_NUMBER_GROUPING_SYMBOLS = map[rune]rune{
		'1':       '1',
		'2':       '2',
		'3':       '3',
		'4':       '4',
		'5':       '5',
		'6':       '6',
		'7':       '7',
		'8':       '8',
		'9':       '9',
		'0':       '0',
		PLUS_SIGN: PLUS_SIGN,
		'*':       '*',
		'A':       'A',
		'B':       'B',
		'C':       'C',
		'D':       'D',
		'E':       'E',
		'F':       'F',
		'G':       'G',
		'H':       'H',
		'I':       'I',
		'J':       'J',
		'K':       'K',
		'L':       'L',
		'M':       'M',
		'N':       'N',
		'O':       'O',
		'P':       'P',
		'Q':       'Q',
		'R':       'R',
		'S':       'S',
		'T':       'T',
		'U':       'U',
		'V':       'V',
		'W':       'W',
		'X':       'X',
		'Y':       'Y',
		'Z':       'Z',
		'a':       'A',
		'b':       'B',
		'c':       'C',
		'd':       'D',
		'e':       'E',
		'f':       'F',
		'g':       'G',
		'h':       'H',
		'i':       'I',
		'j':       'J',
		'k':       'K',
		'l':       'L',
		'm':       'M',
		'n':       'N',
		'o':       'O',
		'p':       'P',
		'q':       'Q',
		'r':       'R',
		's':       'S',
		't':       'T',
		'u':       'U',
		'v':       'V',
		'w':       'W',
		'x':       'X',
		'y':       'Y',
		'z':       'Z',
		'-':       '-',
		'\uFF0D':  '-',
		'\u2010':  '-',
		'\u2011':  '-',
		'\u2012':  '-',
		'\u2013':  '-',
		'\u2014':  '-',
		'\u2015':  '-',
		'\u2212':  '-',
		'/':       '/',
		'\uFF0F':  '/',
		' ':       ' ',
		'\u3000':  ' ',
		'\u2060':  ' ',
		'.':       '.',
		'\uFF0E':  '.',
	}

	// Pattern that makes it easy to distinguish whether a region has a
	// unique international dialing prefix or not. If a region has a
	// unique international prefix (e.g. 011 in USA), it will be
	// represented as a string that contains a sequence of ASCII digits.
	// If there are multiple available international prefixes in a
	// region, they will be represented as a regex string that always
	// contains character(s) other than ASCII digits.
	// Note this regex also includes tilde, which signals waiting for the tone.
	UNIQUE_INTERNATIONAL_PREFIX = regexp.MustCompile("[\\d]+(?:[~\u2053\u223C\uFF5E][\\d]+)?")

	PLUS_CHARS_PATTERN      = regexp.MustCompile("[" + PLUS_CHARS + "]+")
	SEPARATOR_PATTERN       = regexp.MustCompile("[" + VALID_PUNCTUATION + "]+")
	NOT_SEPARATOR_PATTERN   = regexp.MustCompile("[^" + VALID_PUNCTUATION + "]+")
	CAPTURING_DIGIT_PATTERN = regexp.MustCompile("(" + DIGITS + ")")

	// Regular expression of acceptable characters that may start a
	// phone number for the purposes of parsing. This allows us to
	// strip away meaningless prefixes to phone numbers that may be
	// mistakenly given to us. This consists of digits, the plus symbol
	// and arabic-indic digits. This does not contain alpha characters,
	// although they may be used later in the number. It also does not
	// include other punctuation, as this will be stripped later during
	// parsing and is of no information value when parsing a number.
	VALID_START_CHAR         = "[" + PLUS_CHARS + DIGITS + "]"
	VALID_START_CHAR_PATTERN = regexp.MustCompile(VALID_START_CHAR)

	// Regular expression of characters typically used to start a second
	// phone number for the purposes of parsing. This allows us to strip
	// off parts of the number that are actually the start of another
	// number, such as for: (530) 583-6985 x302/x2303 -> the second
	// extension here makes this actually two phone numbers,
	// (530) 583-6985 x302 and (530) 583-6985 x2303. We remove the second
	// extension so that the first number is parsed correctly.
	SECOND_NUMBER_START         = "[\\\\/] *x"
	SECOND_NUMBER_START_PATTERN = regexp.MustCompile(SECOND_NUMBER_START)

	// Regular expression of trailing characters that we want to remove.
	// We remove all characters that are not alpha or numerical characters.
	// The hash character is retained here, as it may signify the previous
	// block was an extension.
	UNWANTED_END_CHARS        = "[[\\P{N}&&\\P{L}]&&[^#]]+$"
	UNWANTED_END_CHAR_PATTERN = regexp.MustCompile(UNWANTED_END_CHARS)

	// We use this pattern to check if the phone number has at least three
	// letters in it - if so, then we treat it as a number where some
	// phone-number digits are represented by letters.
	VALID_ALPHA_PHONE_PATTERN = regexp.MustCompile("^(?:.*?[A-Za-z]){3}.*$")

	// Regular expression of viable phone numbers. This is location
	// independent. Checks we have at least three leading digits, and
	// only valid punctuation, alpha characters and digits in the phone
	// number. Does not include extension data. The symbol 'x' is allowed
	// here as valid punctuation since it is often used as a placeholder
	// for carrier codes, for example in Brazilian phone numbers. We also
	// allow multiple "+" characters at the start.
	// Corresponds to the following:
	// [digits]{minLengthNsn}|
	// plus_sign*(
	//    ([punctuation]|[star])*[digits]
	// ){3,}([punctuation]|[star]|[digits]|[alpha])*
	//
	// The first reg-ex is to allow short numbers (two digits long) to be
	// parsed if they are entered as "15" etc, but only if there is no
	// punctuation in them. The second expression restricts the number of
	// digits to three or more, but then allows them to be in
	// international form, and to have alpha-characters and punctuation.
	//
	// Note VALID_PUNCTUATION starts with a -, so must be the first in the range.
	VALID_PHONE_NUMBER = DIGITS + "{" + strconv.Itoa(MIN_LENGTH_FOR_NSN) + "}" + "|" +
		"[" + PLUS_CHARS + "]*(?:[" + VALID_PUNCTUATION + string(STAR_SIGN) +
		"]*" + DIGITS + "){3,}[" +
		VALID_PUNCTUATION + string(STAR_SIGN) + VALID_ALPHA + DIGITS + "]*"

	// Default extension prefix to use when formatting. This will be put
	// in front of any extension component of the number, after the main
	// national number is formatted. For example, if you wish the default
	// extension formatting to be " extn: 3456", then you should specify
	// " extn: " here as the default extension prefix. This can be
	// overridden by region-specific preferences.
	DEFAULT_EXTN_PREFIX = " ext. "

	// We cap the maximum length of an extension based on the ambiguity of
	// the way the extension is prefixed. As per ITU, the officially allowed
	// length for extensions is actually 40, but we don't support this since
	// we haven't seen real examples and this introduces many false
	// interpretations as the extension labels are not standardized.

	possibleSeparatorsBetweenNumberAndExtLabel = "[ \u00A0\\t,]*"
	// Optional full stop (.) or colon, followed by zero or more spaces/tabs/commas.
	possibleCharsAfterExtLabel = "[:\\.\uFF0E]?[ \u00A0\\t,-]*"
	optionalExtnSuffix         = "#?"

	// Here the extension is called out in a more explicit way, i.e.
	// mentioning it with obvious patterns like "ext.".
	// Canonical-equivalence doesn't seem to be an option with Android
	// java, so we allow two options for representing the accented o -
	// the character itself, and one in the unicode decomposed form
	// with the combining acute accent.
	explicitExtLabels = "(?:e?xt(?:ensi(?:o\u0301?|\u00F3))?n?|" +
		"\uFF45?\uFF58\uFF54\uFF4E?|\u0434\u043E\u0431|anexo)"
	// One-character symbols that can be used to indicate an extension,
	// and less commonly used or more ambiguous extension labels.
	ambiguousExtLabels = "(?:[x\uFF58#\uFF03~\uFF5E]|int|\uFF49\uFF4E\uFF54)"
	ambiguousSeparator = "[- ]+"

	// Regexp of all possible ways to write extensions, for use when
	// parsing. This will be run as a case-insensitive regexp match.
	// maybeStripExtension iterates over all submatches of this pattern
	// and assumes that every capturing group corresponds to the
	// extension digits. When modifying this regexp, ensure that any
	// non-extension grouping uses non-capturing groups (?:...) so as
	// not to introduce additional capturing groups that would break
	// extension parsing.
	//
	// The first regular expression covers RFC 3966 format, where the
	// extension is added using ";ext=". The second is a more generic
	// expression where the extension is mentioned with explicit labels like "ext:". In both
	// cases we allow more digits. The third captures when single
	// character extension labels or less commonly used labels are used,
	// with fewer extension digits to reduce false positives. The fourth
	// covers American numbers where the extension is written with a
	// hash at the end, such as "- 503#".
	EXTN_PATTERNS_FOR_MATCHING = RFC3966_EXTN_PREFIX + "(" + DIGITS + "{1,20})" + "|" +
		possibleSeparatorsBetweenNumberAndExtLabel + explicitExtLabels +
		possibleCharsAfterExtLabel + "(" + DIGITS + "{1,20})" + optionalExtnSuffix + "|" +
		possibleSeparatorsBetweenNumberAndExtLabel + ambiguousExtLabels +
		possibleCharsAfterExtLabel + "(" + DIGITS + "{1,9})" + optionalExtnSuffix + "|" +
		ambiguousSeparator + "(" + DIGITS + "{1,6})#"

	// Additional patterns supported when parsing extensions, not when matching.
	// ",," is commonly used for auto dialling the extension when connected.
	// Semi-colon works on iPhone and Android to pop up a button with the
	// extension number following.
	possibleSeparatorsNumberExtLabelNoComma = "[ \u00A0\\t]*"
	autoDiallingAndExtLabelsFound           = "(?:,{2}|;)"

	EXTN_PATTERNS_FOR_PARSING = EXTN_PATTERNS_FOR_MATCHING + "|" +
		possibleSeparatorsNumberExtLabelNoComma + autoDiallingAndExtLabelsFound +
		possibleCharsAfterExtLabel + "(" + DIGITS + "{1,15})" + optionalExtnSuffix + "|" +
		possibleSeparatorsNumberExtLabelNoComma + "(?:,)+" +
		possibleCharsAfterExtLabel + "(" + DIGITS + "{1,9})" + optionalExtnSuffix

	// Regexp of all known extension prefixes used by different regions
	// followed by 1 or more valid digits, for use when parsing.
	EXTN_PATTERN = regexp.MustCompile("(?i)(?:" + EXTN_PATTERNS_FOR_PARSING + ")$")

	// We append optionally the extension pattern to the end here, as a
	// valid phone number may have an extension prefix appended,
	// followed by 1 or more digits.
	VALID_PHONE_NUMBER_PATTERN = regexp.MustCompile(
		"(?i)^(" + VALID_PHONE_NUMBER + "(?:" + EXTN_PATTERNS_FOR_PARSING + ")?)$")

	NON_DIGITS_PATTERN = regexp.MustCompile(`(\D+)`)
	DIGITS_PATTERN     = regexp.MustCompile(`(\d+)`)

	// The FIRST_GROUP_PATTERN was originally set to $1 but there are some
	// countries for which the first group is not used in the national
	// pattern (e.g. Argentina) so the $1 group does not match correctly.
	// Therefore, we use \d, so that the first group actually used in the
	// pattern will be matched.
	FIRST_GROUP_PATTERN = regexp.MustCompile(`(\$\d)`)
	NP_PATTERN          = regexp.MustCompile(`\$NP`)
	FG_PATTERN          = regexp.MustCompile(`\$FG`)
	CC_PATTERN          = regexp.MustCompile(`\$CC`)

	// A pattern that is used to determine if the national prefix
	// formatting rule has the first group only, i.e., does not start
	// with the national prefix. Note that the pattern explicitly allows
	// for unbalanced parentheses.
	FIRST_GROUP_ONLY_PREFIX_PATTERN = regexp.MustCompile(`^\(?\$1\)?$`)

	REGION_CODE_FOR_NON_GEO_ENTITY = "001"

	// Regular expression of valid global-number-digits for the phone-context parameter, following the
	// syntax defined in RFC3966.
	RFC3966_VISUAL_SEPARATOR             = "[\\-\\.\\(\\)]?"
	RFC3966_PHONE_DIGIT                  = "(" + DIGITS + "|" + RFC3966_VISUAL_SEPARATOR + ")"
	RFC3966_GLOBAL_NUMBER_DIGITS         = "^\\" + string(PLUS_SIGN) + RFC3966_PHONE_DIGIT + "*" + DIGITS + RFC3966_PHONE_DIGIT + "*$"
	RFC3966_GLOBAL_NUMBER_DIGITS_PATTERN = regexp.MustCompile(RFC3966_GLOBAL_NUMBER_DIGITS)

	// Regular expression of valid domainname for the phone-context parameter, following the syntax
	// defined in RFC3966.
	ALPHANUM                   = VALID_ALPHA + DIGITS
	RFC3966_DOMAINLABEL        = "[" + ALPHANUM + "]+((\\-)*[" + ALPHANUM + "])*"
	RFC3966_TOPLABEL           = "[" + VALID_ALPHA + "]+((\\-)*[" + ALPHANUM + "])*"
	RFC3966_DOMAINNAME         = "^(" + RFC3966_DOMAINLABEL + "\\.)*" + RFC3966_TOPLABEL + "\\.?$"
	RFC3966_DOMAINNAME_PATTERN = regexp.MustCompile(RFC3966_DOMAINNAME)

	// Set of country codes that have geographically assigned mobile numbers (see GEO_MOBILE_COUNTRIES
	// below) which are not based on *area codes*. For example, in China mobile numbers start with a
	// carrier indicator, and beyond that are geographically assigned: this carrier indicator is not
	// considered to be an area code.
	GEO_MOBILE_COUNTRIES_WITHOUT_MOBILE_AREA_CODES = map[int32]bool{
		86: true, // China
	}

	// Set of country codes that doesn't have national prefix, but it has area codes.
	COUNTRIES_WITHOUT_NATIONAL_PREFIX_WITH_AREA_CODES = map[int32]bool{
		52: true, // Mexico
	}

	// Set of country calling codes that have geographically assigned mobile numbers. This may not be
	// complete; we add calling codes case by case, as we find geographical mobile numbers or hear
	// from user reports. Note that countries like the US, where we can't distinguish between
	// fixed-line or mobile numbers, are not listed here, since we consider FIXED_LINE_OR_MOBILE to be
	// a possibly geographically-related type anyway (like FIXED_LINE).
	GEO_MOBILE_COUNTRIES = map[int32]bool{
		52: true, // Mexico
		54: true, // Argentina
		55: true, // Brazil
		62: true, // Indonesia: some prefixes only (fixed CMDA wireless)
		86: true, // China
	}

	dataLoadMutex = sync.Mutex{}
)

// INTERNATIONAL and NATIONAL formats are consistent with the definition
// in ITU-T Recommendation E123. For example, the number of the Google
// Switzerland office will be written as "+41 44 668 1800" in
// INTERNATIONAL format, and as "044 668 1800" in NATIONAL format. E164
// format is as per INTERNATIONAL format but with no formatting applied,
// e.g. "+41446681800". RFC3966 is as per INTERNATIONAL format, but with
// all spaces and other separating symbols replaced with a hyphen, and
// with any phone number extension appended with ";ext=". It also will
// have a prefix of "tel:" added, e.g. "tel:+41-44-668-1800".
//
// Note: If you are considering storing the number in a neutral format,
// you are highly advised to use the PhoneNumber class.

type PhoneNumberFormat int

const (
	E164 PhoneNumberFormat = iota
	INTERNATIONAL
	NATIONAL
	RFC3966
)

type PhoneNumberType int

const (
	// NOTES:
	//
	// FIXED_LINE_OR_MOBILE:
	//     In some regions (e.g. the USA), it is impossible to distinguish
	//     between fixed-line and mobile numbers by looking at the phone
	//     number itself.
	// SHARED_COST:
	//     The cost of this call is shared between the caller and the
	//     recipient, and is hence typically less than PREMIUM_RATE calls.
	//     See // http://en.wikipedia.org/wiki/Shared_Cost_Service for
	//     more information.
	// VOIP:
	//     Voice over IP numbers. This includes TSoIP (Telephony Service over IP).
	// PERSONAL_NUMBER:
	//     A personal number is associated with a particular person, and may
	//     be routed to either a MOBILE or FIXED_LINE number. Some more
	//     information can be found here:
	//     http://en.wikipedia.org/wiki/Personal_Numbers
	// UAN:
	//     Used for "Universal Access Numbers" or "Company Numbers". They
	//     may be further routed to specific offices, but allow one number
	//     to be used for a company.
	// VOICEMAIL:
	//     Used for "Voice Mail Access Numbers".
	// UNKNOWN:
	//     A phone number is of type UNKNOWN when it does not fit any of
	// the known patterns for a specific region.
	FIXED_LINE PhoneNumberType = iota
	MOBILE
	FIXED_LINE_OR_MOBILE
	TOLL_FREE
	PREMIUM_RATE
	SHARED_COST
	VOIP
	PERSONAL_NUMBER
	PAGER
	UAN
	VOICEMAIL
	UNKNOWN
)

type MatchType int

const (
	NOT_A_NUMBER MatchType = iota
	NO_MATCH
	SHORT_NSN_MATCH
	NSN_MATCH
	EXACT_MATCH
)

type ValidationResult int

const (
	IS_POSSIBLE ValidationResult = iota
	INVALID_COUNTRY_CODE
	TOO_SHORT
	TOO_LONG
	IS_POSSIBLE_LOCAL_ONLY
	INVALID_LENGTH
)

// TODO(ttacon): leniency comments?
type Leniency int

const (
	POSSIBLE Leniency = iota
	VALID
	STRICT_GROUPING
	EXACT_GROUPING
)

func (l Leniency) Verify(number *PhoneNumber, candidate string) bool {
	_ = "STUB: not implemented"
	return false
}

var (
	// golang map is not go routine safe. Sometimes process exiting
	// because of panic. So adding mutex to synchronize the operation.

	// The set of regions that share country calling code 1.
	// There are roughly 26 regions.
	nanpaRegions = make(map[string]struct{})

	// A mapping from a region code to the PhoneMetadata for that region.
	// Note: Synchronization, though only needed for the Android version
	// of the library, is used in all versions for consistency.
	regionToMetadataMap = make(map[string]*PhoneMetadata)

	// A mapping from a country calling code for a non-geographical
	// entity to the PhoneMetadata for that country calling code.
	// Examples of the country calling codes include 800 (International
	// Toll Free Service) and 808 (International Shared Cost Service).
	// Note: Synchronization, though only needed for the Android version
	// of the library, is used in all versions for consistency.
	countryCodeToNonGeographicalMetadataMap = make(map[int]*PhoneMetadata)

	// A cache for frequently used region-specific regular expressions.
	// The initial capacity is set to 100 as this seems to be an optimal
	// value for Android, based on performance measurements.
	regexCache    = make(map[string]*regexp.Regexp)
	regCacheMutex sync.RWMutex

	// The set of regions the library supports.
	// There are roughly 240 of them and we set the initial capacity of
	// the HashSet to 320 to offer a load factor of roughly 0.75.
	supportedRegions = make(map[string]bool, 320)

	// The set of calling codes that map to the non-geo entity
	// region ("001"). This set currently contains < 12 elements so the
	// default capacity of 16 (load factor=0.75) is fine.
	countryCodesForNonGeographicalRegion = make(map[int]bool, 16)

	// These are maps for our prefix to carrier maps
	carrierPrefixMap = make(map[string]*intStringMap)

	// These are maps for our prefix to geocoding maps
	geocodingPrefixMap = make(map[string]*intStringMap)

	// All the calling codes we support
	supportedCallingCodes = make(map[int]bool, 320)

	// Our once and map for prefix to timezone lookups
	timezoneOnce sync.Once
	timezoneMap  *intStringArrayMap

	// Our map from country code (as integer) to two letter region codes
	countryCodeToRegion map[int][]string
)

var ErrEmptyMetadata = errors.New("empty metadata")

func readFromRegexCache(key string) (*regexp.Regexp, bool) {
	_ = "STUB: not implemented"
	return nil, false
}

func writeToRegexCache(key string, value *regexp.Regexp) { _ = "STUB: not implemented"; return }

func regexFor(pattern string) *regexp.Regexp { _ = "STUB: not implemented"; return nil }

func readFromNanpaRegions(key string) (struct{}, bool) {
	_ = "STUB: not implemented"
	return nil, false
}

func writeToNanpaRegions(key string, val struct{}) { _ = "STUB: not implemented"; return }

func readFromRegionToMetadataMap(key string) (*PhoneMetadata, bool) {
	_ = "STUB: not implemented"
	return nil, false
}

func writeToRegionToMetadataMap(key string, val *PhoneMetadata) { _ = "STUB: not implemented"; return }

func readFromCountryCodeToNonGeographicalMetadataMap(key int) (*PhoneMetadata, bool) {
	_ = "STUB: not implemented"
	return nil, false
}

func writeToCountryCodeToNonGeographicalMetadataMap(key int, v *PhoneMetadata) {
	_ = "STUB: not implemented"
	return
}

func loadMetadataFromFile() error { _ = "STUB: not implemented"; return nil }

// it's a non geographical entity

var (
	currMetadataColl *PhoneMetadataCollection
	reloadMetadata   = true
)

func MetadataCollection() (*PhoneMetadataCollection, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

// Attempts to extract a possible number from the string passed in.
// This currently strips all leading characters that cannot be used to
// start a phone number. Characters that can be used to start a phone
// number are defined in the VALID_START_CHAR_PATTERN. If none of these
// characters are found in the number passed in, an empty string is
// returned. This function also attempts to strip off any alternative
// extensions or endings if two or more are present, such as in the case
// of: (530) 583-6985 x302/x2303. The second extension here makes this
// actually two phone numbers, (530) 583-6985 x302 and (530) 583-6985 x2303.
// We remove the second extension so that the first number is parsed correctly.
func extractPossibleNumber(number string) string { _ = "STUB: not implemented"; return "" }

// Remove trailing non-alpha non-numerical characters.

// Check for extra numbers at the end.

// Checks to see if the string of characters could possibly be a phone
// number at all. At the moment, checks to see that the string begins
// with at least 2 digits, ignoring any punctuation commonly found in
// phone numbers. This method does not require the number to be
// normalized in advance - but does assume that leading non-number symbols
// have been removed, such as by the method extractPossibleNumber.
// @VisibleForTesting
func isViablePhoneNumber(number string) bool { _ = "STUB: not implemented"; return false }

// Normalizes a string of characters representing a phone number. This
// performs the following conversions:
//
//   - Punctuation is stripped.
//
//   - For ALPHA/VANITY numbers:
//
//   - Letters are converted to their numeric representation on a telephone
//     keypad. The keypad used here is the one defined in ITU Recommendation
//     E.161. This is only done if there are 3 or more letters in the
//     number, to lessen the risk that such letters are typos.
//
//   - For other numbers:
//
//   - Wide-ascii digits are converted to normal ASCII (European) digits.
//
//   - Arabic-Indic numerals are converted to European numerals.
//
//   - Spurious alpha characters are stripped.
func normalize(number string) string { _ = "STUB: not implemented"; return "" }

// Normalizes a string of characters representing a phone number. This
// converts wide-ascii and arabic-indic numerals to European numerals,
// and strips punctuation and alpha characters.
func NormalizeDigitsOnly(number string) string { _ = "STUB: not implemented"; return "" }

/* strip non-digits */

// ugly hack still, but fills out the functionality (sort of)
// TODO(ttacon): more completely/elegantly solve this
var arabicIndicNumberals = map[rune]rune{
	'٠':      '0',
	'۰':      '0',
	'١':      '1',
	'۱':      '1',
	'٢':      '2',
	'۲':      '2',
	'٣':      '3',
	'۳':      '3',
	'٤':      '4',
	'۴':      '4',
	'٥':      '5',
	'۵':      '5',
	'٦':      '6',
	'۶':      '6',
	'٧':      '7',
	'۷':      '7',
	'٨':      '8',
	'۸':      '8',
	'٩':      '9',
	'۹':      '9',
	'\uFF10': '0',
	'\uFF11': '1',
	'\uFF12': '2',
	'\uFF13': '3',
	'\uFF14': '4',
	'\uFF15': '5',
	'\uFF16': '6',
	'\uFF17': '7',
	'\uFF18': '8',
	'\uFF19': '9',
}

func normalizeDigits(number string, keepNonDigits bool) string {
	_ = "STUB: not implemented"
	return ""
}

// Normalizes a string of characters representing a phone number. This
// strips all characters which are not diallable on a mobile phone
// keypad (including all non-ASCII digits).
func normalizeDiallableCharsOnly(number string) string { _ = "STUB: not implemented"; return "" }

/* remove non matches */

// Converts all alpha characters in a number to their respective digits
// on a keypad, but retains existing formatting.
func ConvertAlphaCharactersInNumber(number string) string { _ = "STUB: not implemented"; return "" }

// Gets the length of the geographical area code from the PhoneNumber
// object passed in, so that clients could use it to split a national
// significant number into geographical area code and subscriber number. It
// works in such a way that the resultant subscriber number should be
// diallable, at least on some devices. An example of how this could be used:
//
//	number, err := Parse("16502530000", "US");
//	// ... deal with err appropriately ...
//	nationalSignificantNumber := GetNationalSignificantNumber(number);
//	var areaCode, subscriberNumber;
//
//	int areaCodeLength = GetLengthOfGeographicalAreaCode(number);
//	if (areaCodeLength > 0) {
//	  areaCode = nationalSignificantNumber[0:areaCodeLength];
//	  subscriberNumber = nationalSignificantNumber[areaCodeLength:];
//	} else {
//	  areaCode = "";
//	  subscriberNumber = nationalSignificantNumber;
//	}
//
// N.B.: area code is a very ambiguous concept, so the I18N team generally
// recommends against using it for most purposes, but recommends using the
// more general national_number instead. Read the following carefully before
// deciding to use this method:
//
//   - geographical area codes change over time, and this method honors those changes;
//     therefore, it doesn't guarantee the stability of the result it produces.
//   - subscriber numbers may not be diallable from all devices (notably mobile
//     devices, which typically requires the full national_number to be dialled
//     in most regions).
//   - most non-geographical numbers have no area codes, including numbers from
//     non-geographical entities
//   - some geographical numbers have no area codes.
func GetLengthOfGeographicalAreaCode(number *PhoneNumber) int { _ = "STUB: not implemented"; return 0 }

// If a country doesn't use a national prefix, and this number doesn't have an Italian leading
// zero, we assume it is a closed dialling plan with no area codes.
// Note:this is our general assumption, but there are exceptions which are tracked in
// COUNTRIES_WITHOUT_NATIONAL_PREFIX_WITH_AREA_CODES.

// Note this is a rough heuristic; it doesn't cover Indonesia well, for example, where area
// codes are present for some mobile phones but not for others. We have no better way of
// representing this in the metadata at this point.

// Gets the length of the national destination code (NDC) from the
// PhoneNumber object passed in, so that clients could use it to split a
// national significant number into NDC and subscriber number. The NDC of
// a phone number is normally the first group of digit(s) right after the
// country calling code when the number is formatted in the international
// format, if there is a subscriber number part that follows. An example
// of how this could be used:
//
//	PhoneNumberUtil phoneUtil = PhoneNumberUtil.getInstance();
//	PhoneNumber number = phoneUtil.parse("18002530000", "US");
//	String nationalSignificantNumber = phoneUtil.GetNationalSignificantNumber(number);
//	String nationalDestinationCode;
//	String subscriberNumber;
//
//	int nationalDestinationCodeLength =
//	    phoneUtil.GetLengthOfNationalDestinationCode(number);
//	if nationalDestinationCodeLength > 0 {
//	    nationalDestinationCode = nationalSignificantNumber.substring(0,
//	        nationalDestinationCodeLength);
//	    subscriberNumber = nationalSignificantNumber.substring(
//	        nationalDestinationCodeLength);
//	} else {
//	    nationalDestinationCode = "";
//	    subscriberNumber = nationalSignificantNumber;
//	}
//
// Refer to the unittests to see the difference between this function and
// GetLengthOfGeographicalAreaCode().
func GetLengthOfNationalDestinationCode(number *PhoneNumber) int {
	_ = "STUB: not implemented"
	return 0
}

// We don't want to alter the proto given to us, but we don't
// want to include the extension when we format it, so we copy
// it and clear the extension here.

// The pattern will start with "+COUNTRY_CODE " so the first group
// will always be the empty string (before the + symbol) and the
// second group will be the country calling code. The third group
// will be area code if it is not the last group.

// For example Argentinian mobile numbers, when formatted in
// the international format, are in the form of +54 9 NDC XXXX....
// As a result, we take the length of the third group (NDC) and
// add the length of the second group (which is the mobile token),
// which also forms part of the national significant number. This
// assumes that the mobile token is always formatted separately
// from the rest of the phone number.

// Returns the mobile token for the provided country calling code if it
// has one, otherwise returns an empty string. A mobile token is a number
// inserted before the area code when dialing a mobile number from that
// country from abroad.
func GetCountryMobileToken(countryCallingCode int) string { _ = "STUB: not implemented"; return "" }

// Normalizes a string of characters representing a phone number by replacing
// all characters found in the accompanying map with the values therein,
// and stripping all other characters if removeNonMatches is true.
func normalizeHelper(number string,
	normalizationReplacements map[rune]rune,
	removeNonMatches bool) string {
	_ = "STUB: not implemented"
	return ""
}

// If neither of the above are true, we remove this character.

// GetSupportedRegions returns all regions the library has metadata for.
func GetSupportedRegions() map[string]bool { _ = "STUB: not implemented"; return nil }

// GetSupportedCallingCodes returns all country calling codes the library has metadata for, covering both non-geographical
// entities (global network calling codes) and those used for geographical entities. This could be
// used to populate a drop-down box of country calling codes for a phone-number widget, for
// instance.
func GetSupportedCallingCodes() map[int]bool { _ = "STUB: not implemented"; return nil }

// GetSupportedGlobalNetworkCallingCodes returns all global network calling codes the library has metadata for.
func GetSupportedGlobalNetworkCallingCodes() map[int]bool { _ = "STUB: not implemented"; return nil }

// Helper function to check if the national prefix formatting rule has the
// first group only, i.e., does not start with the national prefix.
func formattingRuleHasFirstGroupOnly(nationalPrefixFormattingRule string) bool {
	_ = "STUB: not implemented"
	return false
}

// Tests whether a phone number has a geographical association. It checks
// if the number is associated to a certain region in the country where it
// belongs to. Note that this doesn't verify if the number is actually in use.
//
// A similar method is implemented as PhoneNumberOfflineGeocoder.canBeGeocoded,
// which performs a looser check, since it only prevents cases where prefixes
// overlap for geocodable and non-geocodable numbers. Also, if new phone
// number types were added, we should check if this other method should be
// updated too.
func isNumberGeographical(phoneNumber *PhoneNumber) bool { _ = "STUB: not implemented"; return false }

// Helper function to check region code is not unknown or null.
func isValidRegionCode(regionCode string) bool { _ = "STUB: not implemented"; return false }

// Helper function to check the country calling code is valid.
func hasValidCountryCallingCode(countryCallingCode int) bool {
	_ = "STUB: not implemented"
	return false
}

// Formats a phone number in the specified format using default rules. Note
// that this does not promise to produce a phone number that the user can
// dial from where they are - although we do format in either 'national' or
// 'international' format depending on what the client asks for, we do not
// currently support a more abbreviated format, such as for users in the
// same "area" who could potentially dial the number without area code.
// Note that if the phone number has a country calling code of 0 or an
// otherwise invalid country calling code, we cannot work out which
// formatting rules to apply so we return the national significant number
// with no formatting applied.
func Format(number *PhoneNumber, numberFormat PhoneNumberFormat) string {
	_ = "STUB: not implemented"
	return ""
}

// Unparseable numbers that kept their raw input just use that.
// This is the only case where a number can be formatted as E164
// without a leading '+' symbol (but the original number wasn't
// parseable anyway).
// TODO: Consider removing the 'if' above so that unparseable
// strings without raw input format to the empty string instead of "+00"

// Same as Format(PhoneNumber, PhoneNumberFormat), but accepts a mutable
// StringBuilder as a parameter to decrease object creation when invoked
// many times.
func FormatWithBuf(number *PhoneNumber, numberFormat PhoneNumberFormat, formattedNumber *Builder) {
	_ = "STUB: not implemented"
	// Clear the StringBuilder first.
	return
}

// Early exit for E164 case (even if the country calling code
// is invalid) since no formatting of the national number needs
// to be applied. Extensions are not formatted.

// Note GetRegionCodeForCountryCode() is used because formatting
// information for regions which share a country calling code is
// contained by only one region for performance reasons. For
// example, for NANPA regions it will be contained in the metadata for US.

// Metadata cannot be null because the country calling code is
// valid (which means that the region code cannot be ZZ and must
// be one of our supported region codes).

// FormatByPattern formats a phone number in the specified format using client-defined
// formatting rules. Note that if the phone number has a country calling
// code of zero or an otherwise invalid country calling code, we cannot
// work out things like whether there should be a national prefix applied,
// or how to format extensions, so we return the national significant
// number with no formatting applied.
func FormatByPattern(number *PhoneNumber,
	numberFormat PhoneNumberFormat,
	userDefinedFormats []*NumberFormat) string {
	_ = "STUB: not implemented"
	return ""
}

// Note GetRegionCodeForCountryCode() is used because formatting
// information for regions which share a country calling code is
// contained by only one region for performance reasons. For example,
// for NANPA regions it will be contained in the metadata for US.

// Metadata cannot be null because the country calling code is valid

// If no pattern above is matched, we format the number as a whole.

// Before we do a replacement of the national prefix pattern
// $NP with the national prefix, we need to copy the rule so
// that subsequent replacements for different numbers have the
// appropriate national prefix.

// Replace $NP with national prefix and $FG with the
// first group ($1).

// We don't want to have a rule for how to format the
// national prefix if there isn't one.

// Formats a phone number in national format for dialing using the carrier
// as specified in the carrierCode. The carrierCode will always be used
// regardless of whether the phone number already has a preferred domestic
// carrier code stored. If carrierCode contains an empty string, returns
// the number in national format without any carrier code.
func FormatNationalNumberWithCarrierCode(number *PhoneNumber, carrierCode string) string {
	_ = "STUB: not implemented"
	return ""
}

// Note GetRegionCodeForCountryCode() is used because formatting
// information for regions which share a country calling code is
// contained by only one region for performance reasons. For
// example, for NANPA regions it will be contained in the metadata for US.

// Metadata cannot be null because the country calling code is valid.

func getMetadataForRegionOrCallingCode(countryCallingCode int, regionCode string) *PhoneMetadata {
	_ = "STUB: not implemented"
	return nil
}

// Formats a phone number in national format for dialing using the carrier
// as specified in the preferredDomesticCarrierCode field of the PhoneNumber
// object passed in. If that is missing, use the fallbackCarrierCode passed
// in instead. If there is no preferredDomesticCarrierCode, and the
// fallbackCarrierCode contains an empty string, return the number in
// national format without any carrier code.
//
// Use formatNationalNumberWithCarrierCode instead if the carrier code
// passed in should take precedence over the number's
// preferredDomesticCarrierCode when formatting.
func FormatNationalNumberWithPreferredCarrierCode(
	number *PhoneNumber,
	fallbackCarrierCode string) string {
	_ = "STUB: not implemented"
	return ""
}

// Returns a number formatted in such a way that it can be dialed from a
// mobile phone in a specific region. If the number cannot be reached from
// the region (e.g. some countries block toll-free numbers from being
// called outside of the country), the method returns an empty string.
func FormatNumberForMobileDialing(
	number *PhoneNumber,
	regionCallingFrom string,
	withFormatting bool) string {
	_ = "STUB: not implemented"
	return ""
}

// go impl defaults to ""

// Clear the extension, as that part cannot normally be dialed
// together with the main number.

// can we assume this is safe? (no nil-pointer?)

// Carrier codes may be needed in some countries. We handle this here.

// Historically, we set this to an empty string when parsing with
// raw input if none was found in the input string. However, this
// doesn't result in a number we can dial. For this reason, we
// treat the empty string the same as if it isn't set at all.

// Brazilian fixed line and mobile numbers need to be dialed
// with a carrier code when called within Brazil. Without
// that, most of the carriers won't connect the call.
// Because of that, we return an empty string here.

// For NANPA countries, we output international format for
// numbers that can be dialed internationally, since that
// always works, except for numbers which might potentially be
// short numbers, which are always dialled in national format.

// For non-geographical countries, and Mexican and Chilean fixed
// line and mobile numbers, we output international format for
// numbers that can be dialed internationally as that always
// works.

// MX fixed line and mobile numbers should always be formatted
// in international format, even when dialed within MX. For
// national format to work, a carrier code needs to be used,
// and the correct carrier code depends on if the caller and
// callee are from the same local area. It is trickier to get
// that to work correctly than using international format, which
// is tested to work fine on all carriers. CL fixed line
// numbers need the national prefix when dialing in the national
// format, but don't have it when used for display. The reverse
// is true for mobile numbers. As a result, we output them in
// the international format to make it work.

// We assume that short numbers are not diallable from outside
// their region, so if a number is not a valid regular length
// phone number, we treat it as if it cannot be internationally
// dialled.

// Formats a phone number for out-of-country dialing purposes. If no
// regionCallingFrom is supplied, we format the number in its
// INTERNATIONAL format. If the country calling code is the same as that
// of the region where the number is from, then NATIONAL formatting will
// be applied.
//
// If the number itself has a country calling code of zero or an otherwise
// invalid country calling code, then we return the number with no
// formatting applied.
//
// Note this function takes care of the case for calling inside of NANPA and
// between Russia and Kazakhstan (who share the same country calling code).
// In those cases, no international prefix is used. For regions which have
// multiple international prefixes, the number in its INTERNATIONAL format
// will be returned instead.
func FormatOutOfCountryCallingNumber(
	number *PhoneNumber,
	regionCallingFrom string) string {
	_ = "STUB: not implemented"
	return ""
}

// For NANPA regions, return the national format for these
// regions but prefix it with the country calling code.

// If regions share a country calling code, the country calling
// code need not be dialled. This also applies when dialling
// within a region, so this if clause covers both these cases.
// Technically this is the case for dialling from La Reunion to
// other overseas departments of France (French Guiana, Martinique,
// Guadeloupe), but not vice versa - so we don't cover this edge
// case for now and for those cases return the version including
// country calling code.
// Details here: http://www.petitfute.com/voyage/225-info-pratiques-reunion

// Metadata cannot be null because we checked 'isValidRegionCode()' above.

// For regions that have multiple international prefixes, the
// international format of the number is returned, unless there is
// a preferred international prefix.

// Metadata cannot be null because the country calling code is valid.

// Formats a phone number using the original phone number format that the
// number is parsed from. The original format is embedded in the
// country_code_source field of the PhoneNumber object passed in. If such
// information is missing, the number will be formatted into the NATIONAL
// format by default. When the number contains a leading zero and this is
// unexpected for this country, or we don't have a formatting pattern for
// the number, the method returns the raw input when it is available.
//
// Note this method guarantees no digit will be inserted, removed or
// modified as a result of formatting.
func FormatInOriginalFormat(number *PhoneNumber, regionCallingFrom string) string {
	_ = "STUB: not implemented"
	return ""
}

// We check if we have the formatting pattern because without that, we might format the number
// as a group without national prefix.

// Fall-through to default case.

// We strip non-digits from the NDD here, and from the raw
// input later, so that we can compare them easily.

/* strip non-digits */

// If the region doesn't have a national prefix at all,
// we can safely return the national format without worrying
// about a national prefix being added.

// Otherwise, we check if the original number was entered with
// a national prefix.

// If so, we can safely return the national format.

// Metadata cannot be null here because GetNddPrefixForRegion()
// (above) returns null if there is no metadata for the region.

// The format rule could still be null here if the national
// number was 0 and there was no raw input (this should not
// be possible for numbers generated by the phonenumber library
// as they would also not have a country calling code and we
// would have exited earlier).

// When the format we apply to this number doesn't contain
// national prefix, we can just return the national format.
// TODO: Refactor the code below with the code in
// isNationalPrefixPresentIfRequired.

// We assume that the first-group symbol will never be _before_
// the national prefix.

// National prefix not used when formatting this number.

// Otherwise, we need to remove the national prefix from our output.

// If no digit is inserted/removed/modified as a result of our
// formatting, we return the formatted phone number; otherwise we
// return the raw input the user entered.

// Check if rawInput, which is assumed to be in the national format, has
// a national prefix. The national prefix is assumed to be in digits-only
// form.
func rawInputContainsNationalPrefix(rawInput, nationalPrefix, regionCode string) bool {
	_ = "STUB: not implemented"
	return false
}

// Some Japanese numbers (e.g. 00777123) might be mistaken to
// contain the national prefix when written without it
// (e.g. 0777123) if we just do prefix matching. To tackle that,
// we check the validity of the number if the assumed national
// prefix is removed (777123 won't be valid in Japan).

func hasFormattingPatternForNumber(number *PhoneNumber) bool {
	_ = "STUB: not implemented"
	return false
}

// Formats a phone number for out-of-country dialing purposes.
//
// Note that in this version, if the number was entered originally using
// alpha characters and this version of the number is stored in raw_input,
// this representation of the number will be used rather than the digit
// representation. Grouping information, as specified by characters
// such as "-" and " ", will be retained.
//
// Caveats:
//
//   - This will not produce good results if the country calling code is
//     both present in the raw input _and_ is the start of the national
//     number. This is not a problem in the regions which typically use
//     alpha numbers.
//   - This will also not produce good results if the raw input has any
//     grouping information within the first three digits of the national
//     number, and if the function needs to strip preceding digits/words
//     in the raw input before these digits. Normally people group the
//     first three digits together so this is not a huge problem - and will
//     be fixed if it proves to be so.
func FormatOutOfCountryKeepingAlphaChars(
	number *PhoneNumber,
	regionCallingFrom string) string {
	_ = "STUB: not implemented"
	return ""
}

// If there is no raw input, then we can't keep alpha characters
// because there aren't any. In this case, we return
// formatOutOfCountryCallingNumber.

// Strip any prefix such as country calling code, IDD, that was
// present. We do this by comparing the number in raw_input with
// the parsed number. To do this, first we normalize punctuation.
// We retain number grouping symbols such as " " only.

// Now we trim everything before the first three digits in the
// parsed number. We choose three because all valid alpha numbers
// have 3 digits at the start - if it does not, then we don't trim
// anything at all. Similarly, if the national number was less than
// three digits, we don't trim anything at all.

// If no pattern above is matched, we format the original input.

// The first group is the first group of digits that the user
// wrote together.

// Here we just concatenate them back together after the national
// prefix has been fixed.

// Now we format using this pattern instead of the default pattern,
// but with the national prefix prefixed if necessary. This will not
// work in the cases where the pattern (and not the leading digits)
// decide whether a national prefix needs to be used, since we
// have overridden the pattern to match anything, but that is not
// the case in the metadata to date.

// If an unsupported region-calling-from is entered, or a country
// with multiple international prefixes, the international format
// of the number is returned, unless there is a preferred international
// prefix.

// Metadata cannot be null because the country calling code is valid.

// Strip any extension from the raw input before appending the formatted extension.

// Invalid region entered as country-calling-from (so no metadata
// was found for it) or the region chosen has multiple international
// dialling prefixes.

// Gets the national significant number of the a phone number. Note a
// national significant number doesn't contain a national prefix or
// any formatting.
func GetNationalSignificantNumber(number *PhoneNumber) string {
	_ = "STUB: not implemented"
	// If leading zero(s) have been set, we prefix this now. Note this
	// is not a national prefix.
	return ""
}

// A helper function that is used by format and formatByPattern.
func prefixNumberWithCountryCallingCode(
	countryCallingCode int,
	numberFormat PhoneNumberFormat,
	formattedNumber *Builder) {
	_ = "STUB: not implemented"

	// TODO(ttacon): add some sort of BulkWrite builder to Builder
	// also that name isn't too awesome...:)
	return
}

// Simple wrapper of formatNsn for the common case of no carrier code.
func formatNsn(
	number string, metadata *PhoneMetadata, numberFormat PhoneNumberFormat) string {
	_ = "STUB: not implemented"
	return ""
}

// Note in some regions, the national number can be written in two
// completely different ways depending on whether it forms part of the
// NATIONAL format or INTERNATIONAL format. The numberFormat parameter
// here is used to specify which format to use for those cases. If a
// carrierCode is specified, this will be inserted into the formatted
// string to replace $CC.
func formatNsnWithCarrier(number string, metadata *PhoneMetadata, numberFormat PhoneNumberFormat, carrierCode string) string {
	_ = "STUB: not implemented"
	return ""
}

// When the intlNumberFormats exists, we use that to format national
// number for the INTERNATIONAL format instead of using the
// numberDesc.numberFormats.

func chooseFormattingPatternForNumber(
	availableFormats []*NumberFormat,
	nationalNumber string) *NumberFormat {
	_ = "STUB: not implemented"
	return nil
}

// Strictly match

// We always use the last leading_digits_pattern, as it is the
// most detailed.

// inds[0] == 0 ensures strict match of leading digits

// Simple wrapper of formatNsnUsingPattern for the common case of no carrier code.
func formatNsnUsingPattern(
	nationalNumber string,
	formattingPattern *NumberFormat,
	numberFormat PhoneNumberFormat) string {
	_ = "STUB: not implemented"
	return ""
}

// Note that carrierCode is optional - if null or an empty string, no
// carrier code replacement will take place.
func formatNsnUsingPatternWithCarrier(
	nationalNumber string,
	formattingPattern *NumberFormat,
	numberFormat PhoneNumberFormat,
	carrierCode string) string {
	_ = "STUB: not implemented"
	return ""
}

// Replace the $CC in the formatting rule with the desired carrier code.

// Now replace the $FG in the formatting rule with the first group
// and the carrier code combined in the appropriate way.

// Use the national prefix formatting rule instead.

// Strip any leading punctuation.

// Gets a valid number for the specified region.
func GetExampleNumber(regionCode string) *PhoneNumber { _ = "STUB: not implemented"; return nil }

// Gets a valid number for the specified region and number type.
func GetExampleNumberForType(regionCode string, typ PhoneNumberType) *PhoneNumber {
	_ = "STUB: not implemented"
	// Check the region code is valid.
	return nil
}

// PhoneNumberDesc (pointer?)

// Gets a valid number for the specified country calling code for a non-geographical entity.
func GetExampleNumberForNonGeoEntity(countryCallingCode int) *PhoneNumber {
	_ = "STUB: not implemented"
	return nil
}

// For geographical entities, fixed-line data is always present. However, for non-geographical
// entities, this is not the case, so we have to go through different types to find the
// example number.

// Appends the formatted extension of a phone number to formattedNumber,
// if the phone number had an extension specified.
func maybeAppendFormattedExtension(
	number *PhoneNumber,
	metadata *PhoneMetadata,
	numberFormat PhoneNumberFormat,
	formattedNumber *Builder) {
	_ = "STUB: not implemented"
	return
}

func getNumberDescByType(
	metadata *PhoneMetadata,
	typ PhoneNumberType) *PhoneNumberDesc {
	_ = "STUB: not implemented"
	return nil
}

// Gets the type of a phone number.
func GetNumberType(number *PhoneNumber) PhoneNumberType {
	_ = "STUB: not implemented"
	return *new(PhoneNumberType)
}

func getNumberTypeHelper(nationalNumber string, metadata *PhoneMetadata) PhoneNumberType {
	_ = "STUB: not implemented"
	return *new(PhoneNumberType)
}

// Otherwise, test to see if the number is mobile. Only do this if
// certain that the patterns for mobile and fixed line aren't the same.

// Returns the metadata for the given region code or nil if the region
// code is invalid or unknown.
func getMetadataForRegion(regionCode string) *PhoneMetadata { _ = "STUB: not implemented"; return nil }

func getMetadataForNonGeographicalRegion(countryCallingCode int) *PhoneMetadata {
	_ = "STUB: not implemented"
	return nil
}

func isNumberPossibleForDesc(nationalNumber string, numberDesc *PhoneNumberDesc) bool {
	_ = "STUB: not implemented"
	// Check if any possible number lengths are present; if so, we use them to avoid checking the
	// validation pattern if they don't match. If they are absent, this means they match the general
	// description, which we have already checked before checking a specific number type.
	return false
}

// Strictly match

func isNumberMatchingDesc(nationalNumber string, numberDesc *PhoneNumberDesc) bool {
	_ = "STUB: not implemented"
	return false
}

// Strictly match

// Tests whether a phone number matches a valid pattern. Note this doesn't
// verify the number is actually in use, which is impossible to tell by
// just looking at a number itself.
func IsValidNumber(number *PhoneNumber) bool { _ = "STUB: not implemented"; return false }

// Tests whether a phone number is valid for a certain region. Note this
// doesn't verify the number is actually in use, which is impossible to
// tell by just looking at a number itself. If the country calling code is
// not the same as the country calling code for the region, this immediately
// exits with false. After this, the specific number pattern rules for the
// region are examined. This is useful for determining for example whether
// a particular number is valid for Canada, rather than just a valid NANPA
// number.
// Warning: In most cases, you want to use IsValidNumber() instead. For
// example, this method will mark numbers from British Crown dependencies
// such as the Isle of Man as invalid for the region "GB" (United Kingdom),
// since it has its own region code, "IM", which may be undesirable.
func IsValidNumberForRegion(number *PhoneNumber, regionCode string) bool {
	_ = "STUB: not implemented"
	return false
}

// Either the region code was invalid, or the country calling
// code for this number does not match that of the region code.

// Returns the region where a phone number is from. This could be used for
// geocoding at the region level.
func GetRegionCodeForNumber(number *PhoneNumber) string { _ = "STUB: not implemented"; return "" }

func getRegionCodeForNumberFromRegionList(
	number *PhoneNumber,
	regionCodes []string) string {
	_ = "STUB: not implemented"
	return ""
}

// If leadingDigits is present, use this. Otherwise, do
// full validation. Metadata cannot be null because the
// region codes come from the country calling code map.

// Non capturing grouping to support OR'ed alternatives (e.g. 555|1[78]|2)

// Returns the region code that matches the specific country calling code.
// In the case of no region code being found, ZZ will be returned. In the
// case of multiple regions, the one designated in the metadata as the
// "main" region for this calling code will be returned. If the
// countryCallingCode entered is valid but doesn't match a specific region
// (such as in the case of non-geographical calling codes like 800) the
// value "001" will be returned (corresponding to the value for World in
// the UN M.49 schema).
func GetRegionCodeForCountryCode(countryCallingCode int) string {
	_ = "STUB: not implemented"
	return ""
}

// Returns a list with the region codes that match the specific country
// calling code. For non-geographical country calling codes, the region
// code 001 is returned. Also, in the case of no region code being found,
// an empty list is returned.
func GetRegionCodesForCountryCode(countryCallingCode int) []string {
	_ = "STUB: not implemented"
	return nil
}

// Returns the country calling code for a specific region. For example, this
// would be 1 for the United States, and 64 for New Zealand.
func GetCountryCodeForRegion(regionCode string) int { _ = "STUB: not implemented"; return 0 }

// Returns the country calling code for a specific region. For example,
// this would be 1 for the United States, and 64 for New Zealand. Assumes
// the region is already valid.
func getCountryCodeForValidRegion(regionCode string) int { _ = "STUB: not implemented"; return 0 }

// Returns the national dialling prefix for a specific region. For example,
// this would be 1 for the United States, and 0 for New Zealand. Set
// stripNonDigits to true to strip symbols like "~" (which indicates a
// wait for a dialling tone) from the prefix returned. If no national prefix
// is present, we return null.
//
// Warning: Do not use this method for do-your-own formatting - for some
// regions, the national dialling prefix is used only for certain types
// of numbers. Use the library's formatting functions to prefix the
// national prefix when required.
func GetNddPrefixForRegion(regionCode string, stripNonDigits bool) string {
	_ = "STUB: not implemented"
	return ""
}

// If no national prefix was found, we return an empty string.

// Note: if any other non-numeric symbols are ever used in
// national prefixes, these would have to be removed here as well.

// Checks if this is a region under the North American Numbering Plan
// Administration (NANPA).
func IsNANPACountry(regionCode string) bool { _ = "STUB: not implemented"; return false }

// Checks if the number is a valid vanity (alpha) number such as 800
// MICROSOFT. A valid vanity number will start with at least 3 digits and
// will have three or more alpha characters. This does not do
// region-specific checks - to work out if this number is actually valid
// for a region, it should be parsed and methods such as
// IsPossibleNumberWithReason() and IsValidNumber() should be used.
func IsAlphaNumber(number string) bool { _ = "STUB: not implemented"; return false }

// Number is too short, or doesn't match the basic phone
// number pattern.

// Convenience wrapper around IsPossibleNumberWithReason(). Instead of
// returning the reason for failure, this method returns a boolean value.
func IsPossibleNumber(number *PhoneNumber) bool { _ = "STUB: not implemented"; return false }

func descHasPossibleNumberData(desc *PhoneNumberDesc) bool { _ = "STUB: not implemented"; return false }

func mergeLengths(l1 []int32, l2 []int32) []int32 { _ = "STUB: not implemented"; return nil }

// Helper method to check a number against possible lengths for this number type, and determine
// whether it matches, or is too short or too long.
func testNumberLength(number string, metadata *PhoneMetadata, numberType PhoneNumberType) ValidationResult {
	_ = "STUB: not implemented"
	return *new(ValidationResult)
}

// There should always be "possibleLengths" set for every element. This is declared in the XML
// schema which is verified by PhoneNumberMetadataSchemaTest.
// For size efficiency, where a sub-description (e.g. fixed-line) has the same possibleLengths
// as the parent, this is missing, so we fall back to the general desc (where no numbers of the
// type exist at all, there is one possible length (-1) which is guaranteed not to match the
// length of any real phone number).

// The rare case has been encountered where no fixedLine data is available (true for some
// non-geographical entities), so we just check mobile.

// Note that when adding the possible lengths from mobile, we have to again check they
// aren't empty since if they are this indicates they are the same as the general desc and
// should be obtained from there.

// If the type is not supported at all (indicated by the possible lengths containing -1 at this
// point) we return invalid length.

// This is safe because there is never an overlap between the possible lengths and the local-only
// lengths; this is checked at build time.

// We skip the first element; we've already checked it.

// Check whether a phone number is a possible number. It provides a more
// lenient check than IsValidNumber() in the following sense:
//
//   - It only checks the length of phone numbers. In particular, it
//     doesn't check starting digits of the number.
//   - It doesn't attempt to figure out the type of the number, but uses
//     general rules which applies to all types of phone numbers in a
//     region. Therefore, it is much faster than isValidNumber.
//   - For fixed line numbers, many regions have the concept of area code,
//     which together with subscriber number constitute the national
//     significant number. It is sometimes okay to dial the subscriber number
//     only when dialing in the same area. This function will return true
//     if the subscriber-number-only version is passed in. On the other hand,
//     because isValidNumber validates using information on both starting
//     digits (for fixed line numbers, that would most likely be area codes)
//     and length (obviously includes the length of area codes for fixed
//     line numbers), it will return false for the subscriber-number-only
//     version.
func IsPossibleNumberWithReason(number *PhoneNumber) ValidationResult {
	_ = "STUB: not implemented"
	return *new(ValidationResult)
}

// Note: For Russian Fed and NANPA numbers, we just use the rules
// from the default region (US or Russia) since the
// getRegionCodeForNumber will not work if the number is possible
// but not valid. This would need to be revisited if the possible
// number pattern ever differed between various regions within
// those plans.

// Metadata cannot be null because the country calling code is valid.

// Handling case of numbers with no metadata.

// Attempts to extract a valid number from a phone number that is too long
// to be valid, and resets the PhoneNumber object passed in to that valid
// version. If no valid number could be extracted, the PhoneNumber object
// passed in will not be modified.
func TruncateTooLongNumber(number *PhoneNumber) bool { _ = "STUB: not implemented"; return false }

// Gets an AsYouTypeFormatter for the specific region.
// TODO(ttacon): uncomment once we do asyoutypeformatter.go
// public AsYouTypeFormatter getAsYouTypeFormatter(String regionCode) {
//    return new AsYouTypeFormatter(regionCode);
// }

// Extracts country calling code from fullNumber, returns it and places
// the remaining number in nationalNumber. It assumes that the leading plus
// sign or IDD has already been removed. Returns 0 if fullNumber doesn't
// start with a valid country calling code, and leaves nationalNumber
// unmodified.
func extractCountryCode(fullNumber, nationalNumber *Builder) int {
	_ = "STUB: not implemented"
	return 0
}

// Country codes do not begin with a '0'.

var ErrTooShortAfterIDD = errors.New("phone number had an IDD, but " +
	"after this was not long enough to be a viable phone number")

// Tries to extract a country calling code from a number. This method will
// return zero if no country calling code is considered to be present.
// Country calling codes are extracted in the following ways:
//
//   - by stripping the international dialing prefix of the region the
//     person is dialing from, if this is present in the number, and looking
//     at the next digits
//   - by stripping the '+' sign if present and then looking at the next digits
//   - by comparing the start of the number and the country calling code of
//     the default region. If the number is not considered possible for the
//     numbering plan of the default region initially, but starts with the
//     country calling code of this region, validation will be reattempted
//     after stripping this country calling code. If this number is considered a
//     possible number, then the first digits will be considered the country
//     calling code and removed as such.
//
// It will throw a NumberParseException if the number starts with a '+' but
// the country calling code supplied after this does not match that of any
// known region.
func maybeExtractCountryCode(
	number string,
	defaultRegionMetadata *PhoneMetadata,
	nationalNumber *Builder,
	keepRawInput bool,
	phoneNumber *PhoneNumber) (int, error) {
	_ = "STUB: not implemented"
	return 0, nil
}

// Set the default prefix to be something that will never match.

// If this fails, they must be using a strange country calling code
// that we don't recognize, or that doesn't exist.

// Check to see if the number starts with the country calling code
// for the default region. If so, we remove the country calling
// code, and do some checks on the validity of the number before
// and after.

// Strictly match

/* Don't need the carrier code */

// If the number was not valid before but is valid now, or
// if it was too long before, we consider the number with
// the country calling code stripped to be a better result and
// keep that instead.

// No country calling code present.

// Strips the IDD from the start of the number if present. Helper function
// used by maybeStripInternationalPrefixAndNormalize.
func parsePrefixAsIdd(iddPattern *regexp.Regexp, number *Builder) bool {
	_ = "STUB: not implemented"
	return false
}

// ind is a two element slice
// Only strip this if the first digit after the match is not
// a 0, since country calling codes cannot begin with 0.

// Strips any international prefix (such as +, 00, 011) present in the
// number provided, normalizes the resulting number, and indicates if
// an international prefix was present.
func maybeStripInternationalPrefixAndNormalize(
	number *Builder,
	possibleIddPrefix string) PhoneNumber_CountryCodeSource {
	_ = "STUB: not implemented"
	return *new(PhoneNumber_CountryCodeSource)
}

// Check to see if the number begins with one or more plus signs.
// Return is an int pair [start,end]
// Strictly match from string start

// Can now normalize the rest of the number since we've consumed
// the "+" sign at the start.

// Attempt to parse the first digits as an international prefix.

// Strips any national prefix (such as 0, 1) present in the number provided.
// @VisibleForTesting
func maybeStripNationalPrefixAndCarrierCode(
	number *Builder,
	metadata *PhoneMetadata,
	carrierCode *Builder) bool {
	_ = "STUB: not implemented"
	return false
}

// Early return for numbers of zero length.

// Strictly match from string start
// Attempt to parse the first digits as a national prefix.

// Strictly match

// Check if the original number is viable.

// prefixMatcher.group(numOfGroups) == null implies nothing was
// captured by the capturing groups in possibleNationalPrefix;
// therefore, no transformation is necessary, and we just
// remove the national prefix.

// groups is a list of index pairs, idx0,idx1 defines the whole match, idx2+ submatches.
// Subtract one to ignore group(0) in count

// Negative idx means subgroup did not match
// If the original number was viable, and the resultant number
// is not, we return.

// groups[1] == last match idx

// Ensure group(1) matched before slicing
// always extract group(1) as carrier code

// Check that the resultant number is still viable. If not,
// return. Check this by copying the string buffer and
// making the transformation on the copy first.

// Check group(1) got a submatch
// group(1) idxs

// MaybeSeparateExtensionFromPhone will extract any extension (as in, the part
// of the number dialled after the call is connected, usually indicated with
// extn, ext, x or similar) from the end of the number and returns it along
// with the proceeding phone number. The phone number will maintain its
// original formatting including alpha characters.
func MaybeSeparateExtensionFromPhone(rawPhone string) (phoneNumber string, extension string) {
	_ = "STUB: not implemented"
	return "", ""
}

func splitAtExtensionSeparator(rawPhone string) (phoneNumber string, extWithSeparator string) {
	_ = "STUB: not implemented"
	return "", ""
}

func removeLeadingExtensionSeparator(extWithSeparator string) string {
	_ = "STUB: not implemented"
	return ""
}

// Strips any extension (as in, the part of the number dialled after the
// call is connected, usually indicated with extn, ext, x or similar) from
// the end of the number, and returns it.
// @VisibleForTesting
func maybeStripExtension(number *Builder) string {
	_ = "STUB: not implemented"
	// If we find a potential extension, and the number preceding this is
	// a viable number, we assume it is an extension.
	return ""
}

// The numbers are captured into groups in the regular expression.

// We go through the capturing groups until we find one
// that captured some digits. If none did, then we will
// return the empty string.

// Checks to see that the region code used is valid, or if it is not valid,
// that the number to parse starts with a + symbol so that we can attempt
// to infer the region from the number. Returns false if it cannot use the
// region provided and the region cannot be inferred.
func checkRegionForParsing(numberToParse, defaultRegion string) bool {
	_ = "STUB: not implemented"
	return false
}

// If the number is null or empty, we can't infer the region.

// Parses a string and returns it in proto buffer format. This method will
// throw a NumberParseException if the number is not considered to be a
// possible number. Note that validation of whether the number is actually
// a valid number for a particular region is not performed. This can be
// done separately with IsValidNumber().
func Parse(numberToParse, defaultRegion string) (*PhoneNumber, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

// Same as Parse(string, string), but accepts mutable PhoneNumber as a
// parameter to decrease object creation when invoked many times.
func ParseToNumber(numberToParse, defaultRegion string, phoneNumber *PhoneNumber) error {
	_ = "STUB: not implemented"
	return nil
}

// Parses a string and returns it in proto buffer format. This method
// differs from Parse() in that it always populates the raw_input field of
// the protocol buffer with numberToParse as well as the country_code_source
// field.
func ParseAndKeepRawInput(
	numberToParse, defaultRegion string) (*PhoneNumber, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

// Same as ParseAndKeepRawInput(String, String), but accepts a mutable
// PhoneNumber as a parameter to decrease object creation when invoked many
// times.
func ParseAndKeepRawInputToNumber(
	numberToParse, defaultRegion string,
	phoneNumber *PhoneNumber) error {
	_ = "STUB: not implemented"
	return nil
}

// Returns an iterable over all PhoneNumberMatch PhoneNumberMatches in text.
// This is a shortcut for findNumbers(CharSequence, String, Leniency, long)
// getMatcher(text, defaultRegion, Leniency.VALID, Long.MAX_VALUE)}.
// public Iterable<PhoneNumberMatch> findNumbers(CharSequence text, String defaultRegion) {
//    return findNumbers(text, defaultRegion, Leniency.VALID, Long.MAX_VALUE);
// }

// Returns an iterable over all PhoneNumberMatch PhoneNumberMatches in text.
// public Iterable<PhoneNumberMatch> findNumbers(
//	final CharSequence text, final String defaultRegion, final Leniency leniency,
//	final long maxTries) {
//
//		return new Iterable<PhoneNumberMatch>() {
//			public Iterator<PhoneNumberMatch> iterator() {
//				return new PhoneNumberMatcher(
//					PhoneNumberUtil.this, text, defaultRegion, leniency, maxTries);
//			}
//		};
//	}

// A helper function to set the values related to leading zeros in a
// PhoneNumber.
func setItalianLeadingZerosForPhoneNumber(
	nationalNum string, phoneNumber *PhoneNumber) {
	_ = "STUB: not implemented"
	return
}

// Note that if the national number is all "0"s, the last "0"
// is not counted as a leading zero.

var (
	ErrInvalidCountryCode = errors.New("invalid country code")
	ErrNotANumber         = errors.New("the phone number supplied is not a number")
	ErrTooShortNSN        = errors.New("the string supplied is too short to be a phone number")
)

// Parses a string and fills up the phoneNumber. This method is the same
// as the public Parse() method, with the exception that it allows the
// default region to be null, for use by IsNumberMatch(). checkRegion should
// be set to false if it is permitted for the default region to be null or
// unknown ("ZZ").
func parseHelper(
	numberToParse, defaultRegion string,
	keepRawInput, checkRegion bool,
	phoneNumber *PhoneNumber) error {
	_ = "STUB: not implemented"
	return nil
}

// Check the region supplied is valid, or that the extracted number
// starts with some sort of + sign so the number's region can be determined.

// Attempt to parse extension first, since it doesn't require
// region-specific data and we want to have the non-normalised
// number here.

// Check to see if the number is given in international format so we
// know whether this number is from the default region or not.

// TODO: This method should really just take in the string buffer that
// has already been created, and just remove the prefix, rather than
// taking in a string and then outputting a string buffer.

// There might be a plus at the beginning

// Strip the plus-char, and try again.

// Metadata cannot be null because the country calling
// code is valid.

// If no extracted country calling code, use the region supplied
// instead. The national number is just the normalized version of
// the number we were given to parse.

// We require that the NSN remaining after stripping the national
// prefix and carrier code be of a possible length for the region.
// Otherwise, we don't do the stripping, since the original number
// could be a valid short number.

var ErrNumTooLong = errors.New("the string supplied is too long to be a phone number")

// Extracts the value of the phone-context parameter of numberToExtractFrom where the index of
// ";phone-context=" is the parameter indexOfPhoneContext, following the syntax defined in
// RFC3966.
func extractPhoneContext(numberToExtractFrom string, indexOfPhoneContext int) string {
	_ = "STUB: not implemented"
	// If no phone-context parameter is present
	return ""
}

// If phone-context parameter is empty

// find end of this phone-context (go doesn't have a indexOf(s, after))

// If phone-context is not the last parameter

// Returns whether the value of phoneContext follows the syntax defined in RFC3966.
func isPhoneContextValid(phoneContext string) bool { _ = "STUB: not implemented"; return false }

// Does phone-context value match pattern of global-number-digits or domainname

// Converts numberToParse to a form that we can parse and write it to
// nationalNumber if it is written in RFC3966; otherwise extract a possible
// number out of it and write to nationalNumber.
func buildNationalNumberForParsing(
	numberToParse string,
	nationalNumber *Builder) error {
	_ = "STUB: not implemented"
	return nil
}

// If the phone context contains a phone number prefix, we need to capture it, whereas domains
// will be ignored.

// Additional parameters might follow the phone context. If so, we will remove them here
// because the parameters after phone context are not important for parsing the phone
// number.

// Now append everything between the "tel:" prefix and the phone-context. This should include
// the national number, an optional extension or isdn-subaddress component. Note we also
// handle the case when "tel:" is missing, as we have seen in some of the phone number inputs.
// In that case, we append everything from the beginning.

// Extract a possible number from the string passed in (this
// strips leading characters that could not be the start of a
// phone number.)

// Delete the isdn-subaddress and everything after it if it is present.
// Note extension won't appear at the same time with isdn-subaddress
// according to paragraph 5.3 of the RFC3966 spec,

// If both phone context and isdn-subaddress are absent but other
// parameters are present, the parameters are left in nationalNumber.
// This is because we are concerned about deleting content from a
// potential number string when there is no strong evidence that the
// number is actually written in RFC3966.

// Takes two phone numbers and compares them for equality.
//
// Returns EXACT_MATCH if the country_code, NSN, presence of a leading zero
// for Italian numbers and any extension present are the same.
// Returns NSN_MATCH if either or both has no region specified, and the NSNs
// and extensions are the same.
// Returns SHORT_NSN_MATCH if either or both has no region specified, or the
// region specified is the same, and one NSN could be a shorter version of
// the other number. This includes the case where one has an extension
// specified, and the other does not.
// Returns NO_MATCH otherwise.
// For example, the numbers +1 345 657 1234 and 657 1234 are a SHORT_NSN_MATCH.
// The numbers +1 345 657 1234 and 345 657 are a NO_MATCH.
func IsNumberMatchWithNumbers(firstNumberIn, secondNumberIn *PhoneNumber) MatchType {
	_ = "STUB: not implemented"
	// Make copies of the phone number so that the numbers passed in are not edited.
	return *new(MatchType)
}

// First clear raw_input, country_code_source and
// preferred_domestic_carrier_code fields and any empty-string
// extensions so that we can use the proto-buffer equality method.

// NOTE(ttacon): don't think we need this in go land...

// Early exit if both had extensions and these are different.

// Both had country_code specified.

// TODO(ttacon): remove when make gen-equals

// A SHORT_NSN_MATCH occurs if there is a difference because of
// the presence or absence of an 'Italian leading zero', the
// presence or absence of an extension, or one NSN being a
// shorter variant of the other.

// This is not a match.

// Checks cases where one or both country_code fields were not
// specified. To make equality checks easier, we first set the
// country_code fields to be equal.

// If all else was the same, then this is an NSN_MATCH.
// TODO(ttacon): remove when make gen-equals

// Returns true when one national number is the suffix of the other or both
// are the same.
func isNationalNumberSuffixOfTheOther(firstNumber, secondNumber *PhoneNumber) bool {
	_ = "STUB: not implemented"
	return false
}

// Note that endsWith returns true if the numbers are equal.

// Takes two phone numbers as strings and compares them for equality. This is
// a convenience wrapper for IsNumberMatch(PhoneNumber, PhoneNumber). No
// default region is known.
func IsNumberMatch(firstNumber, secondNumber string) MatchType {
	_ = "STUB: not implemented"
	return *new(MatchType)
}

// Takes two phone numbers and compares them for equality. This is a
// convenience wrapper for IsNumberMatch(PhoneNumber, PhoneNumber). No
// default region is known.
func IsNumberMatchWithOneNumber(
	firstNumber *PhoneNumber, secondNumber string) MatchType {
	_ = "STUB: not implemented"
	// First see if the second number has an implicit country calling
	// code, by attempting to parse it.
	return *new(MatchType)
}

// The second number has no country calling code. EXACT_MATCH is no
// longer possible. We parse it as if the region was the same as that
// for the first number, and if EXACT_MATCH is returned, we replace
// this with NSN_MATCH.

// If the first number didn't have a valid country calling
// code, then we parse the second number without one as well.

// Returns true if the number can be dialled from outside the region, or
// unknown. If the number can only be dialled from within the region,
// returns false. Does not check the number is a valid number. Note that,
// at the moment, this method does not handle short numbers.
// TODO: Make this method public when we have enough metadata to make it worthwhile.
func canBeInternationallyDialled(number *PhoneNumber) bool { _ = "STUB: not implemented"; return false }

// Note numbers belonging to non-geographical entities
// (e.g. +800 numbers) are always internationally diallable,
// and will be caught here.

// Returns true if the supplied region supports mobile number portability.
// Returns false for invalid, unknown or regions that don't support mobile
// number portability.
func IsMobileNumberPortableRegion(regionCode string) bool { _ = "STUB: not implemented"; return false }

// GetTimezonesForPrefix returns a slice of Timezones corresponding to the number passed
// or error when it is impossible to convert the string to int
// The algorithm tries to match the timezones starting from the maximum
// number of phone number digits and decreasing until it finds one or reaches 0
func GetTimezonesForPrefix(number string) ([]string, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

// strip any leading +

// GetTimezonesForNumber returns the names of timezones which we believe maps to the
// passed in number.
func GetTimezonesForNumber(number *PhoneNumber) ([]string, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

func lazyLoadPrefixes(langMap map[string]*intStringMap, dataFS embed.FS, dir, language string) (*intStringMap, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

// if we already have prefixes (or nil if they don't exist) return that

// try to load the data file for this language

// language doesn't have data

func getValueForNumber(langMap map[string]*intStringMap, dataFS embed.FS, dir, language string, maxLength int, number *PhoneNumber) (string, int, error) {
	_ = "STUB: not implemented"
	return "", 0, nil
}

// GetCarrierForNumber returns the carrier we believe the number belongs to. Note due
// to number porting this is only a guess, there is no guarantee to its accuracy.
func GetCarrierForNumber(number *PhoneNumber, lang string) (string, error) {
	_ = "STUB: not implemented"
	return "", nil
}

// GetSafeCarrierDisplayNameForNumber Gets the name of the carrier for the given phone number
// only when it is 'safe' to display to users.
// A carrier name is considered safe if the number is valid and
// for a region that doesn't support mobile number portability .
func GetSafeCarrierDisplayNameForNumber(phoneNumber *PhoneNumber, lang string) (string, error) {
	_ = "STUB: not implemented"
	return "", nil
}

// GetCarrierWithPrefixForNumber returns the carrier we believe the number belongs to, as well as
// its prefix. Note due to number porting this is only a guess, there is no guarantee to its accuracy.
func GetCarrierWithPrefixForNumber(number *PhoneNumber, lang string) (string, int, error) {
	_ = "STUB: not implemented"
	return "", 0, nil
}

// fallback to english

// GetGeocodingForNumber returns the location we think the number was first acquired in. This is
// just our best guess, there is no guarantee to its accuracy.
func GetGeocodingForNumber(number *PhoneNumber, lang string) (string, error) {
	_ = "STUB: not implemented"
	return "", nil
}

// fallback to english

// fallback to locale

// fallback to english
