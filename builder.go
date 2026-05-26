package phonenumbers

// ----------------------------------------------------------------------------
// Golang port of:
// https://github.com/googlei18n/libphonenumber/blob/master/tools/java/common/src/com/google/i18n/phonenumbers/BuildMetadataFromXml.java
// ----------------------------------------------------------------------------

func sp(value string) *string { _ = "STUB: not implemented"; return nil }

func bp(value bool) *bool { _ = "STUB: not implemented"; return nil }

func ip(value int32) *int32 { _ = "STUB: not implemented"; return nil }

func BuildPhoneMetadataCollection(inputXML []byte, liteBuild bool, specialBuild bool, isShortNumberMetadata bool) (*PhoneMetadataCollection, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

func buildPhoneMetadataFromElement(document *PhoneNumberMetadataE, liteBuild bool, specialBuild bool, isShortNumberMetadata bool, isAlternateFormatsMetadata bool) (*PhoneMetadataCollection, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

// Build a mapping from a country calling code to the region codes which denote the country/region
// represented by that country code. In the case of multiple countries sharing a calling code,
// such as the NANPA countries, the one indicated with "isMainCountryForCode" in the metadata
// should be first.
func BuildCountryCodeToRegionMap(metadataCollection *PhoneMetadataCollection) map[int][]string {
	_ = "STUB: not implemented"
	return nil
}

// For most countries, there will be only one region code for the country calling code.

// For alternate formats, there are no region codes at all.

func validateRE(re string, removeWhitespace bool) string {
	_ = "STUB: not implemented"
	// Removes all the whitespace and newline from the regexp. Not Ming pattern compile options to
	// make it work across programming languages.
	return ""
}

func loadTerritoryTagMetadata(regionCode string, territory *TerritoryE, nationalPrefix string) *PhoneMetadata {
	_ = "STUB: not implemented"
	return nil
}

func setLeadingDigitsPatterns(numberFormatElement *NumberFormatE, format *NumberFormat) {
	_ = "STUB: not implemented"
	return
}

/**
 * Extracts the pattern for international format. If there is no intlFormat, default to using the
 * national format. If the intlFormat is set to "NA" the intlFormat should be ignored.
 *
 * @throws  RuntimeException if multiple intlFormats have been encountered.
 * @return  whether an international number format is defined.
 */
func loadInternationalFormat(metadata *PhoneMetadata, numberFormatElement *NumberFormatE, nationalFormat *NumberFormat) bool {
	_ = "STUB: not implemented"
	return false
}

// Default to use the same as the national pattern if none is defined.

/**
 * Extracts the pattern for the national format.
 *
 * @throws  RuntimeException if multiple or no formats have been encountered.
 */
// @VisibleForTesting
func loadNationalFormat(metadata *PhoneMetadata, numberFormatElement *NumberFormatE, format *NumberFormat) {
	_ = "STUB: not implemented"
	return
}

func getDomesticCarrierCodeFormattingRule(carrierCodeFormattingRule string, nationalPrefix string) string {
	_ = "STUB: not implemented"
	// Replace $FG with the first group ($1) and $NP with the national prefix.
	return ""
}

func getNationalPrefixFormattingRule(nationalPrefixFormattingRule string, nationalPrefix string) string {
	_ = "STUB: not implemented"
	// Replace $NP with national prefix and $FG with the first group ($1).
	return ""
}

/**
 * Extracts the available formats from the provided DOM element. If it does not contain any
 * nationalPrefixFormattingRule, the one passed-in is retained; similarly for
 * nationalPrefixOptionalWhenFormatting. The nationalPrefix, nationalPrefixFormattingRule and
 * nationalPrefixOptionalWhenFormatting values are provided from the parent (territory) element.
 */
// @VisibleForTesting
func loadAvailableFormats(metadata *PhoneMetadata, element *TerritoryE, nationalPrefix string,
	nationalPrefixFormattingRule string, nationalPrefixOptionalWhenFormatting bool) {
	_ = "STUB: not implemented"
	return
}

// Only a small number of regions need to specify the intlFormats in the xml. For the majority
// of countries the intlNumberFormat metadata is an exact copy of the national NumberFormat
// metadata. To minimize the size of the metadata file, we only keep intlNumberFormats that
// actually differ in some way to the national formats.

/**
 * Checks if the possible lengths provided as a sorted set are equal to the possible lengths
 * stored already in the description pattern. Note that possibleLengths may be empty but must not
 * be null, and the PhoneNumberDesc passed in should also not be null.
 */
func arePossibleLengthsEqual(possibleLengths map[int32]bool, desc *PhoneNumberDesc) bool {
	_ = "STUB: not implemented"
	return false
}

// check whether the same elements exist

/**
 * Parses a possible length string into a set of the integers that are covered.
 *
 * @param possibleLengthString  a string specifying the possible lengths of phone numbers. Follows
 *     this syntax: ranges or elements are separated by commas, and ranges are specified in
 *     [min-max] notation, inclusive. For example, [3-5],7,9,[11-14] should be parsed to
 *     3,4,5,7,9,11,12,13,14.
 */
func parsePossibleLengthStringToSet(possibleLengthString string) map[int32]bool {
	_ = "STUB: not implemented"
	return nil
}

// Strip the leading and trailing [], and split on the -.

// We don't even accept [6-7] since we prefer the shorter 6,7 variant; for a range to be in
// use the hyphen needs to replace at least one digit.

/**
 * Reads the possible lengths present in the metadata and splits them into two sets: one for
 * full-length numbers, one for local numbers.
 *
 * @param data  one or more phone number descriptions, represented as XML nodes
 * @param lengths  a set to which to add possible lengths of full phone numbers
 * @param localOnlyLengths  a set to which to add possible lengths of phone numbers only diallable
 *     locally (e.g. within a province)
 */
func populatePossibleLengthSets(data []*PhoneNumberDescE, lengths map[int32]bool, localOnlyLengths map[int32]bool) {
	_ = "STUB: not implemented"
	return
}

// We don't add to the phone metadata yet, since we want to sort length elements found under
// different nodes first, make sure there are no duplicates between them and that the
// localOnly lengths don't overlap with the others.

// intersect our two maps

// We check again when we set these lengths on the metadata itself in setPossibleLengths
// that the elements in localOnly are not also in lengths. For e.g. the generalDesc, it
// might have a local-only length for one type that is a normal length for another type. We
// don't consider this an error, but we do want to remove the local-only lengths.

// It is okay if at this time we have duplicates, because the same length might be possible
// for e.g. fixed-line and for mobile numbers, and this method operates potentially on
// multiple phoneNumberDesc XML elements.

/**
 * Processes a phone number description element from the XML file and returns it as a
 * PhoneNumberDesc. If the description element is a fixed line or mobile number, the parent
 * description will be used to fill in the whole element if necessary, or any components that are
 * missing. For all other types, the parent description will only be used to fill in missing
 * components if the type has a partial definition. For example, if no "tollFree" element exists,
 * we assume there are no toll free numbers for that locale, and return a phone number description
 * with "NA" for both the national and possible number patterns. Note that the parent description
 * must therefore already be processed before this method is called on any child elements.
 *
 * @param parentDesc  a generic phone number description that will be used to fill in missing
 *     parts of the description, or null if this is the root node. This must be processed before
 *     this is run on any child elements.
 * @param countryElement  the XML element representing all the country information
 * @param numberType  the name of the number type, corresponding to the appropriate tag in the XML
 *     file with information about that type
 * @return  complete description of that phone number type
 */
// @VisibleForTesting
func processPhoneNumberDescElement(parentDesc *PhoneNumberDesc, element *PhoneNumberDescE) *PhoneNumberDesc {
	_ = "STUB: not implemented"
	return nil
}

// New way of handling possible number lengths. We don't do this for the general
// description, since these tags won't be present; instead we will calculate its values
// based on the values for all the other number type descriptions (see
// setPossibleLengthsGeneralDesc).

/**
 * Sets the possible length fields in the metadata from the sets of data passed in. Checks that
 * the length is covered by the "parent" phone number description element if one is present, and
 * if the lengths are exactly the same as this, they are not filled in for efficiency reasons.
 *
 * @param parentDesc  the "general description" element or null if desc is the generalDesc itself
 * @param desc  the PhoneNumberDesc object that we are going to set lengths for
 */
func setPossibleLengths(lengths map[int32]bool, localOnlyLengths map[int32]bool, parentDesc *PhoneNumberDesc, desc *PhoneNumberDesc) {
	_ = "STUB: not implemented"
	// We clear these fields since the metadata tends to inherit from the parent element for other
	// fields (via a mergeFrom).
	return
}

// Only add the lengths to this sub-type if they aren't exactly the same as the possible
// lengths in the general desc (for metadata size reasons).

// We shouldn't have possible lengths defined in a child element that are not covered by
// the general description. We check this here even though the general description is
// derived from child elements because it is only derived from a subset, and we need to
// ensure *all* child elements have a valid possible length.

// We check that the local-only length isn't also a normal possible length (only relevant for
// the general-desc, since within elements such as fixed-line we would throw an exception if we
// saw this) before adding it to the collection of possible local-only lengths.

// We check it is covered by either of the possible length sets of the parent
// PhoneNumberDesc, because for example 7 might be a valid localOnly length for mobile, but
// a valid national length for fixedLine, so the generalDesc would have the 7 removed from
// localOnly.

// Need to sort both lists, possible lengths need to be ordered

/**
 * Sets possible lengths in the general description, derived from certain child elements.
 */
func setPossibleLengthsGeneralDesc(generalDesc *PhoneNumberDesc, metadataId string, data *TerritoryE, isShortNumberMetadata bool) {
	_ = "STUB: not implemented"
	return
}

// The general description node should *always* be present if metadata for other types is
// present, aside from in some unit tests.
// (However, for e.g. formatting metadata in PhoneNumberAlternateFormats, no PhoneNumberDesc
// elements are present).

// We shouldn't have anything specified at the "general desc" level: we are going to
// calculate this ourselves from child elements.

// Make a copy here since we want to remove some nodes, but we don't want to do that on our actual data.
// We remove no-international dialing

func loadCountryMetadata(regionCode string, element *TerritoryE, isShortNumberMetadata bool, isAlternateFormatsMetadata bool) *PhoneMetadata {
	_ = "STUB: not implemented"
	return nil
}

// The alternate formats metadata does not need most of the patterns to be set.

func setRelevantDescPatterns(metadata *PhoneMetadata, element *TerritoryE, isShortNumberMetadata bool) {
	_ = "STUB: not implemented"
	return
}

// Calculate the possible lengths for the general description. This will be based on the
// possible lengths of the child elements.

// Set fields used by regular length phone numbers.

// Set fields used by short numbers.

// <!ELEMENT phoneNumberMetadata (territories)>
type PhoneNumberMetadataE struct {
	// <!ELEMENT territories (territory+)>
	Territories []TerritoryE `xml:"territories>territory"`
}

// <!ELEMENT territory (references?, availableFormats?, generalDesc, noInternationalDialling?,
// fixedLine?, mobile?, pager?, tollFree?, premiumRate?,
// sharedCost?, personalNumber?, voip?, uan?, voicemail?)>
type TerritoryE struct {
	// <!ATTLIST territory id CDATA #REQUIRED>
	ID string `xml:"id,attr"`

	// <!ATTLIST territory mainCountryForCode (true) #IMPLIED>
	MainCountryForCode bool `xml:"mainCountryForCode,attr"`

	// <!ATTLIST territory leadingDigits CDATA #IMPLIED>
	LeadingDigits string `xml:"leadingDigits,attr"`

	// <!ATTLIST territory countryCode CDATA #REQUIRED>
	CountryCode int32 `xml:"countryCode,attr"`

	// <!ATTLIST territory nationalPrefix CDATA #IMPLIED>
	NationalPrefix string `xml:"nationalPrefix,attr"`

	// <!ATTLIST territory internationalPrefix CDATA #IMPLIED>
	InternationalPrefix string `xml:"internationalPrefix,attr"`

	// <!ATTLIST territory preferredInternationalPrefix CDATA #IMPLIED>
	PreferredInternationalPrefix string `xml:"preferredInternationalPrefix,attr"`

	// <!ATTLIST territory nationalPrefixFormattingRule CDATA #IMPLIED>
	NationalPrefixFormattingRule string `xml:"nationalPrefixFormattingRule,attr"`

	// <!ATTLIST territory mobileNumberPortableRegion (true) #IMPLIED>
	MobileNumberPortableRegion bool `xml:"mobileNumberPortableRegion,attr"`

	// <!ATTLIST territory nationalPrefixForParsing CDATA #IMPLIED>
	NationalPrefixForParsing string `xml:"nationalPrefixForParsing,attr"`

	// <!ATTLIST territory nationalPrefixTransformRule CDATA #IMPLIED>
	NationalPrefixTransformRule string `xml:"nationalPrefixTransformRule,attr"`

	// <!ATTLIST territory preferredExtnPrefix CDATA #IMPLIED>
	PreferredExtnPrefix string `xml:"preferredExtnPrefix,attr"`

	// <!ATTLIST territory nationalPrefixOptionalWhenFormatting (true) #IMPLIED>
	NationalPrefixOptionalWhenFormatting bool `xml:"nationalPrefixOptionalWhenFormatting,attr"`

	// <!ATTLIST territory carrierCodeFormattingRule CDATA #IMPLIED>
	CarrierCodeFormattingRule string `xml:"carrierCodeFormattingRule,attr"`

	// <!ELEMENT references (sourceUrl+)>
	// <!ELEMENT sourceUrl (#PCDATA)>
	References []string `xml:"references>sourceUrl"`

	// <!ELEMENT availableFormats (numberFormat+)>
	AvailableFormats []NumberFormatE `xml:"availableFormats>numberFormat"`

	// <!ELEMENT generalDesc (nationalNumberPattern)>
	GeneralDesc *PhoneNumberDescE `xml:"generalDesc"`

	// <!ELEMENT noInternationalDialling (nationalNumberPattern, possibleLengths, exampleNumber)>
	NoInternationalDialing *PhoneNumberDescE `xml:"noInternationalDialing"`

	// <!ELEMENT fixedLine (nationalNumberPattern, possibleLengths, exampleNumber)>
	FixedLine *PhoneNumberDescE `xml:"fixedLine"`

	// <!ELEMENT mobile (nationalNumberPattern, possibleLengths, exampleNumber)>
	Mobile *PhoneNumberDescE `xml:"mobile"`

	// <!ELEMENT pager (nationalNumberPattern, possibleLengths, exampleNumber)>
	Pager *PhoneNumberDescE `xml:"pager"`

	// <!ELEMENT tollFree (nationalNumberPattern, possibleLengths, exampleNumber)>
	TollFree *PhoneNumberDescE `xml:"tollFree"`

	// <!ELEMENT premiumRate (nationalNumberPattern, possibleLengths, exampleNumber)>
	PremiumRate *PhoneNumberDescE `xml:"premiumRate"`

	// <!ELEMENT sharedCost (nationalNumberPattern, possibleLengths, exampleNumber)>
	SharedCost *PhoneNumberDescE `xml:"sharedCost"`

	// <!ELEMENT personalNumber (nationalNumberPattern, possibleLengths, exampleNumber)>
	PersonalNumber *PhoneNumberDescE `xml:"personalNumber"`

	// <!ELEMENT voip (nationalNumberPattern, possibleLengths, exampleNumber)>
	VOIP *PhoneNumberDescE `xml:"voip"`

	// <!ELEMENT uan (nationalNumberPattern, possibleLengths, exampleNumber)>
	UAN *PhoneNumberDescE `xml:"uan"`

	// <!ELEMENT voicemail (nationalNumberPattern, possibleLengths, exampleNumber)>
	VoiceMail *PhoneNumberDescE `xml:"voicemail"`

	// <!ELEMENT uan (nationalNumberPattern, possibleLengths, exampleNumber)>
	StandardRate *PhoneNumberDescE `xml:"standardRate"`

	// <!ELEMENT voicemail (nationalNumberPattern, possibleLengths, exampleNumber)>
	ShortCode *PhoneNumberDescE `xml:"shortCode"`

	// <!ELEMENT uan (nationalNumberPattern, possibleLengths, exampleNumber)>
	Emergency *PhoneNumberDescE `xml:"emergency"`

	// <!ELEMENT voicemail (nationalNumberPattern, possibleLengths, exampleNumber)>
	CarrierSpecific *PhoneNumberDescE `xml:"carrierSpecific"`
}

// <!ELEMENT numberFormat (leadingDigits*, format, intlFormat*)>
type NumberFormatE struct {
	// <!ELEMENT leadingDigits (#PCDATA)>
	LeadingDigits []string `xml:"leadingDigits"`

	// <!ELEMENT format (#PCDATA)>
	Format string `xml:"format"`

	// <!ELEMENT intlFormat (#PCDATA)>
	InternationalFormat []string `xml:"intlFormat"`

	// <!ATTLIST numberFormat nationalPrefixFormattingRule CDATA #IMPLIED>
	NationalPrefixFormattingRule string `xml:"nationalPrefixFormattingRule,attr"`

	// <!ATTLIST numberFormat nationalPrefixOptionalWhenFormatting (true) #IMPLIED>
	NationalPrefixOptionalWhenFormatting *bool `xml:"nationalPrefixOptionalWhenFormatting,attr"`

	// <!ATTLIST numberFormat carrierCodeFormattingRule CDATA #IMPLIED>
	CarrierCodeFormattingRule string `xml:"carrierCodeFormattingRule,attr"`

	// <!ATTLIST numberFormat pattern CDATA #REQUIRED>
	Pattern string `xml:"pattern,attr" validate:"required"`
}

type PossibleLengthE struct {
	// <!ATTLIST possibleLengths national CDATA #REQUIRED>
	National string `xml:"national,attr"`

	// <!ATTLIST possibleLengths localOnly CDATA #IMPLIED>
	LocalOnly string `xml:"localOnly,attr"`
}

type PhoneNumberDescE struct {
	// <!ELEMENT nationalNumberPattern (#PCDATA)>
	NationalNumberPattern string `xml:"nationalNumberPattern"`

	// <!ELEMENT possibleLengths EMPTY>
	PossibleLengths *PossibleLengthE `xml:"possibleLengths"`

	// <!ELEMENT exampleNumber (#PCDATA)>
	ExampleNumber string `xml:"exampleNumber"`
}
