package version

import (
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// Supported PingDirectory versions
const (
	PingDirectory9100 = "9.1.0.0"
	PingDirectory9101 = "9.1.0.1"
	PingDirectory9102 = "9.1.0.2"
	PingDirectory9200 = "9.2.0.0"
	PingDirectory9201 = "9.2.0.1"
)

func IsValid(versionString string) bool {
	return getSortedVersionIndex(versionString) != -1
}

func getSortedVersionIndex(versionString string) int {
	for i, version := range getSortedVersions() {
		if version == versionString {
			return i
		}
	}
	return -1
}

func getSortedVersions() []string {
	return []string{
		PingDirectory9100,
		PingDirectory9101,
		PingDirectory9102,
		PingDirectory9200,
		PingDirectory9201,
	}
}

// Compare two PingDirectory versions. Returns a negative number if the first argument is less than the second,
// zero if they are equal, and a positive number if the first argument is greater than the second
func Compare(version1, version2 string) (int, error) {
	version1Index := getSortedVersionIndex(version1)
	if version1Index == -1 {
		return 0, errors.New("Invalid version: " + version1)
	}
	version2Index := getSortedVersionIndex(version2)
	if version2Index == -1 {
		return 0, errors.New("Invalid version: " + version2)
	}

	return version1Index - version2Index, nil
}

func Parse(versionString string) (string, error) {
	if len(versionString) == 0 {
		return versionString, errors.New("failed to parse PingDirectory version: empty version string")
	}

	var err error
	versionDigits := strings.Split(versionString, ".")
	// Expect a version like "x.x" or "x.x.x.x"
	// If only two digits are supplied, the last two will be assumed to be "0.0"
	if len(versionDigits) != 2 && len(versionDigits) != 4 {
		return versionString, errors.New("failed to parse PingDirectory version '" + versionString + "', Expected either two digits (e.g. '9.1') or four digits (e.g. '9.1.0.0')")
	}
	if len(versionDigits) == 2 {
		versionString += ".0.0"
	}
	if !IsValid(versionString) {
		err = errors.New("unsupported PingDirectory version: " + versionString)
	}
	return versionString, err
}

func CheckResourceSupported(diagnostics *diag.Diagnostics, minimumVersion, actualVersion, resourceName string) {
	// Check that the version is at least the minimum version
	compare, err := Compare(actualVersion, minimumVersion)
	if err != nil {
		diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare < 0 {
		diagnostics.AddError(resourceName+" is only supported for PingDirectory versions "+minimumVersion+" and later", "Found PD version "+actualVersion)
		return
	}
}
