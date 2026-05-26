package main

import (
	"fmt"
	"os"

	"github.com/nyaruka/phonenumbers"
)

const distPath = "data"

func main() {
	if err := buildMetadata(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func buildMetadata() error { _ = "STUB: not implemented"; return nil }

func cloneUpstreamRepo(url string) error { _ = "STUB: not implemented"; return nil }

func buildNumberMetadata(srcFile, varName, dstFile string, short bool) (*phonenumbers.PhoneMetadataCollection, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

func buildRegionMetadata(metadata *phonenumbers.PhoneMetadataCollection, varName, dstFile string) error {
	_ = "STUB: not implemented"
	return nil
}

// generate our map data

func buildTimezoneMetadata(srcFile, varName, dstFile string) error {
	_ = "STUB: not implemented"
	return nil
}

// build our map of prefix to timezones

// parse our prefix

// generate our map data

func buildPrefixMetadata(srcDir, varName, dstDir string) error {
	_ = "STUB: not implemented"
	// get our top level language directories
	return nil
}

// for each directory

// only look at directories

// build a map for that directory

// save it for our language

// iterate through our map, creating our full set of values and prefixes

// make sure we won't overrun uint16s

// need sorted prefixes for our diff writing to work

// sorted values compress better

// build our reverse mapping from value to offset

// write our map

// first write our values, as length of string and raw bytes

// then then number of prefix / value pairs

// we write our prefix / value pairs as a varint of the difference of the previous prefix
// and a uint16 of the value index

func renderMap(prefixMap map[int][]string) ([]byte, error) {
	_ = "STUB: not implemented"
	// build lists of our keys and values
	return nil, nil
}

// first write our values, as length of string and raw bytes

// then the number of keys

// we write our key / value pairs as a varint of the difference of the previous prefix
// and a uint16 of the value index

// first write our prefix

// then our values

// write our number of values

// then each value as the interned index

// generates the file contents for a data file
func generateBinFile(varName string, data []byte) []byte { _ = "STUB: not implemented"; return nil }

func readMappingsForDir(dir string) (map[int]string, error) {
	_ = "STUB: not implemented"
	return nil, nil
}
