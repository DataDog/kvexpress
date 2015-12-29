// +build linux darwin freebsd

package commands

import (
	"testing"
)

var exampleData = "This\nIs\nA\nMulti\nLine\nFile\nThat\nContains\nMultiple\nLines\nFor\nTesting.\n"
var trimmedExampleData = "A\nMulti\nLine\nFile\nThat\nContains\nMultiple\nLines\nFor\nTesting.\n"
var exampleDataSHA = "7820fd75cdaee5c6a53230e9a9ec213ca75fb69aa4bb74eb30194ab92decbb8d"

func TestLineCount(t *testing.T) {
	t.Log("Expecting 12 lines.")
	lines := LineCount(exampleData)
	if lines != 12 {
		t.Errorf("Got the wrong number of lines: %d", lines)
	}
}

func TestLengthCheck(t *testing.T) {
	t.Log("Expecting longEnough to be true.")
	longEnough := LengthCheck(exampleData, 10)
	if !longEnough {
		t.Errorf("Error: It certainly is long enough.")
	}
}

func TestComputeChecksum(t *testing.T) {
	t.Log("Expecting the checksum to match.")
	testSHA := ComputeChecksum(exampleData)
	if testSHA != exampleDataSHA {
		t.Errorf("The checksums didn't match.")
	}
}

func TestChecksumCompare(t *testing.T) {
	t.Log("Expecting the checksums to match.")
	checksumMatch := ChecksumCompare(exampleData, exampleDataSHA)
	if !checksumMatch {
		t.Errorf("The checksums should match.")
	}
}

func TestChecksumCompareNoMatch(t *testing.T) {
	t.Log("Expecting the checksums to NOT match.")
	checksumMatch := ChecksumCompare(exampleData, "this-is-a-fake-checksum")
	if checksumMatch {
		t.Errorf("The checksums should NOT match.")
	}
}

func TestRemoveLines(t *testing.T) {
	t.Log("Removing 2 lines.")
	leftoverLines := removeLines(exampleData, 2)
	if leftoverLines != trimmedExampleData {
		t.Errorf("We didn't trim the right lines.")
	}
}
