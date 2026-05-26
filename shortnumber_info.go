package phonenumbers

var (
	shortNumberRegionToMetadataMap = make(map[string]*PhoneMetadata)
)

func readFromShortNumberRegionToMetadataMap(key string) (*PhoneMetadata, bool) {
	_ = "STUB: not implemented"
	return nil, false
}

func writeToShortNumberRegionToMetadataMap(key string, val *PhoneMetadata) {
	_ = "STUB: not implemented"
	return
}

func init() {
	err := loadShortNumberMetadataFromFile()
	if err != nil {
		panic(err)
	}
}

var (
	currShortNumberMetadataColl *PhoneMetadataCollection
	shortNumberReloadMetadata   = true
)

func ShortNumberMetadataCollection() (*PhoneMetadataCollection, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

func loadShortNumberMetadataFromFile() error { _ = "STUB: not implemented"; return nil }

// it's a non geographical entity, unused

// Check whether a short number is a possible number. If a country calling code is shared by
// multiple regions, this returns true if it's possible in any of them. This provides a more
// lenient check than #isValidShortNumber.
// See IsPossibleShortNumberForRegion(PhoneNumber, string) for details.
func IsPossibleShortNumber(number *PhoneNumber) bool { _ = "STUB: not implemented"; return false }

// Check whether a short number is a possible number when dialed from the given region. This
// provides a more lenient check than IsValidShortNumberForRegion.
func IsPossibleShortNumberForRegion(number *PhoneNumber, regionDialingFrom string) bool {
	_ = "STUB: not implemented"
	return false
}

// Tests whether a short number matches a valid pattern. If a country calling code is shared by
// multiple regions, this returns true if it's valid in any of them. Note that this doesn't verify
// the number is actually in use, which is impossible to tell by just looking at the number
// itself. See IsValidShortNumberForRegion(PhoneNumber, String) for details.
func IsValidShortNumber(number *PhoneNumber) bool { _ = "STUB: not implemented"; return false }

// If a matching region had been found for the phone number from among two or more regions,
// then we have already implicitly verified its validity for that region.

// Tests whether a short number matches a valid pattern in a region. Note that this doesn't verify
// the number is actually in use, which is impossible to tell by just looking at the number itself.
func IsValidShortNumberForRegion(number *PhoneNumber, regionDialingFrom string) bool {
	_ = "STUB: not implemented"
	return false
}

func getShortNumberMetadataForRegion(regionCode string) *PhoneMetadata {
	_ = "STUB: not implemented"
	return nil
}

func getRegionCodeForShortNumberFromRegionList(number *PhoneNumber, regionCodes []string) string {
	_ = "STUB: not implemented"
	return ""
}

// The number is valid for this region.

// Helper method to check that the country calling code of the number matches the region it's
// being dialed from.
func regionDialingFromMatchesNumber(number *PhoneNumber, regionDialingFrom string) bool {
	_ = "STUB: not implemented"
	return false
}

// TODO: Once we have benchmarked ShortNumberInfo, consider if it is worth keeping
// this performance optimization.
func matchesPossibleNumberAndNationalNumber(number string, numberDesc *PhoneNumberDesc) bool {
	_ = "STUB: not implemented"
	return false
}

// In these countries, if extra digits are added to an emergency number, it no longer connects
// to the emergency service.
var REGIONS_WHERE_EMERGENCY_NUMBERS_MUST_BE_EXACT = []string{"BR", "CL", "NI"}

func matchesEmergencyNumber(number string, regionCode string, allowPrefixMatch bool) bool {
	_ = "STUB: not implemented"
	return false
}

// Returns false if the number starts with a plus sign. We don't believe dialing the country
// code before emergency numbers (e.g. +1911) works, but later, if that proves to work, we can
// add additional logic here to handle it.

// Returns true if the given number exactly matches an emergency service number in the given
// region.
//
// This method takes into account cases where the number might contain formatting, but doesn't
// allow additional digits to be appended. Note that isEmergencyNumber(number, region)
// implies connectsToEmergencyNumber(number, region).
//
// number: the phone number to test
// regionCode: the region where the phone number is being dialed
// return: whether the number exactly matches an emergency services number in the given region
func IsEmergencyNumber(number string, regionCode string) bool {
	_ = "STUB: not implemented"
	return false
}

// Returns true if the given number, exactly as dialed, might be used to connect to an emergency
// service in the given region.
//
// This method accepts a string, rather than a PhoneNumber, because it needs to distinguish
// cases such as "+1 911" and "911", where the former may not connect to an emergency service in
// all cases but the latter would. This method takes into account cases where the number might
// contain formatting, or might have additional digits appended (when it is okay to do that in
// the specified region).
//
// number: the phone number to test
// regionCode: the region where the phone number is being dialed
// return: whether the number might be used to connect to an emergency service in the given region
func ConnectsToEmergencyNumber(number string, regionCode string) bool {
	_ = "STUB: not implemented"
	return false
}
