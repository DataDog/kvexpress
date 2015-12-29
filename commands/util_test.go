// +build linux darwin freebsd

package commands

import (
	"testing"
)

var testData = "This\nIs\nA\nMulti\nLine\nFile\nThat\nContains\nMultiple\nLines\nFor\nTesting.\n"
var compressedTestData = "H4sIAAAJbogA/wrJyCzm8izmcuTyLc0pyeTyycxL5XLLzEnlCslILOFyzs8rSczMK4bIFgCFQQqKudzyi7hCUotLMvPS9bgAAAAA//8BAAD//5xzJo1EAAAA"

func TestCompressData(t *testing.T) {
	compressed := CompressData(testData)
	if compressed != compressedTestData {
		t.Error("The compression is off.")
	}
}

func TestDecompressData(t *testing.T) {
	decompressed := DecompressData(compressedTestData)
	if decompressed != testData {
		t.Error("The decompression is off.")
	}
}
