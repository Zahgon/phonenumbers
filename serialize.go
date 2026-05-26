package phonenumbers

// intStringMap is our data structure for maps from prefixes to a single string
// this is used for our carrier and geocoding maps
type intStringMap struct {
	Map       map[int]string
	MaxLength int
}

func loadPrefixMap(data []byte) (*intStringMap, error) { _ = "STUB: not implemented"; return nil, nil }

// ok, first read in our number of values

// then our values

// read our # of mappings

// first read our diff

// then our map

// return our values

func digitCount(n int) int { _ = "STUB: not implemented"; return 0 }

// Special case for 0

// Handle negative numbers

// intStringArrayMap is our map from an int to an array of strings
// this is used for our timezone and region maps
type intStringArrayMap struct {
	Map       map[int][]string
	MaxLength int
}

func loadIntArrayMap(data []byte) (*intStringArrayMap, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

// ok, first read in our number of values

// then our values

// read our # of mappings

// first read our diff

// then our values

// return our values

func decodeUnzip(data []byte) ([]byte, error) { _ = "STUB: not implemented"; return nil, nil }
