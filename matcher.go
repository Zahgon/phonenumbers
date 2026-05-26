package phonenumbers

import (
	"regexp"
)

type PhoneNumberMatcher struct {
}

func NewPhoneNumberMatcher(seq string) *PhoneNumberMatcher {
	_ = "STUB: not implemented"
	// TODO(ttacon): to be implemented
	return nil
}

func ContainsOnlyValidXChars(number *PhoneNumber, candidate string) bool {
	_ = "STUB: not implemented"
	// The characters 'x' and 'X' can be (1) a carrier code, in which
	// case they always precede the national significant number or (2)
	// an extension sign, in which case they always precede the extension
	// number. We assume a carrier code is more than 1 digit, so the first
	// case has to have more than 1 consecutive 'x' or 'X', whereas the
	// second case can only have exactly 1 'x' or 'X'. We ignore the
	// character if it appears as the last character of the string.
	return false
}

// This is the carrier code case, in which the 'X's
// always precede the national significant number.

// This is the extension sign case, in which the 'x'
// or 'X' should always precede the extension number.

func IsNationalPrefixPresentIfRequired(number *PhoneNumber) bool {
	_ = "STUB: not implemented"
	// First, check how we deduced the country code. If it was written
	// in international format, then the national prefix is not required.
	return false
}

// Check if a national prefix should be present when formatting this number.

// To do this, we check that a national prefix formatting rule was
// present and that it wasn't just the first-group symbol ($1) with
// punctuation.

// The national-prefix is optional in these cases, so we
// don't need to check if it was present.

// National Prefix not needed for this number.

// Normalize the remainder.

// Check if we found a national prefix and/or carrier code at
// the start of the raw input, and return the result.

func ContainsMoreThanOneSlashInNationalNumber(
	number *PhoneNumber,
	candidate string) bool {
	_ = "STUB: not implemented"
	return false
}

// No slashes, this is okay.

// Now look for a second one.

// Only one slash, this is okay.

// If the first slash is after the country calling code, this is permitted.

// Any more slashes and this is illegal.

func CheckNumberGroupingIsValid(
	number *PhoneNumber,
	candidate string,
	fn func(*PhoneNumber, string, []string) bool) bool {
	_ = "STUB: not implemented"
	// TODO(ttacon): to be implemented
	return false
}

func AllNumberGroupsRemainGrouped(
	number *PhoneNumber,
	normalizedCandidate string,
	formattedNumberGroups []string) bool {
	_ = "STUB: not implemented"
	return false
}

// First skip the country code if the normalized candidate contained it.

// Check each group of consecutive digits are not broken into
// separate groupings in the normalizedCandidate string.

// Fails if the substring of normalizedCandidate starting
// from fromIndex doesn't contain the consecutive digits
// in formattedNumberGroups[i].

// Moves fromIndex forward.

// We are at the position right after the NDC. We get
// the region used for formatting information based on
// the country code in the phone number, rather than the
// number itself, as we do not need to distinguish between
// different countries with the same country calling code
// and this is faster.

// This means there is no formatting symbol after the
// NDC. In this case, we only accept the number if there
// is no formatting symbol at all in the number, except
// for extensions. This is only important for countries
// with national prefixes.

// The check here makes sure that we haven't mistakenly already
// used the extension to match the last group of the subscriber
// number. Note the extension cannot have formatting in-between digits.

func AllNumberGroupsAreExactlyPresent(
	number *PhoneNumber,
	normalizedCandidate string,
	formattedNumberGroups []string) bool {
	_ = "STUB: not implemented"
	return false
}

// Set this to the last group, skipping it if the number has an extension.

// First we check if the national significant number is formatted
// as a block. We use contains and not equals, since the national
// significant number may be present with a prefix such as a national
// number prefix, or the country code itself.

// Starting from the end, go through in reverse, excluding the first
// group, and check the candidate and number groups are the same.

// Now check the first group. There may be a national prefix at
// the start, so we only check that the candidate group ends with
// the formatted number group.

// Returns whether the given national number (a string containing only decimal digits) matches
// the national number pattern defined in the given PhoneNumberDesc message.
func MatchNationalNumber(number string, numberDesc *PhoneNumberDesc, allowPrefixMatch bool) bool {
	_ = "STUB: not implemented"
	return false
}

// We don't want to consider it a prefix match when matching non-empty input against an empty pattern.

func match(number string, pattern *regexp.Regexp, allowPrefixMatch bool) bool {
	_ = "STUB: not implemented"
	return false
}

// Strictly match
