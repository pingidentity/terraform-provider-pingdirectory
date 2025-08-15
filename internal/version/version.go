// Copyright © 2025 Ping Identity Corporation

package version

import (
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// Supported PingDirectory versions
const (
	PingDirectory9300  = "9.3.0.0"
	PingDirectory9301  = "9.3.0.1"
	PingDirectory9302  = "9.3.0.2"
	PingDirectory9303  = "9.3.0.3"
	PingDirectory9304  = "9.3.0.4"
	PingDirectory9305  = "9.3.0.5"
	PingDirectory9306  = "9.3.0.6"
	PingDirectory9307  = "9.3.0.7"
	PingDirectory9308  = "9.3.0.8"
	PingDirectory10000 = "10.0.0.0"
	PingDirectory10001 = "10.0.0.1"
	PingDirectory10002 = "10.0.0.2"
	PingDirectory10003 = "10.0.0.3"
	PingDirectory10004 = "10.0.0.4"
	PingDirectory10100 = "10.1.0.0"
	PingDirectory10102 = "10.1.0.2"
	PingDirectory10103 = "10.1.0.3"
	PingDirectory10200 = "10.2.0.0"
	PingDirectory10201 = "10.2.0.1"
	PingDirectory10300 = "10.3.0.0"
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
		PingDirectory9300,
		PingDirectory9301,
		PingDirectory9302,
		PingDirectory9303,
		PingDirectory9304,
		PingDirectory9305,
		PingDirectory9306,
		PingDirectory9307,
		PingDirectory9308,
		PingDirectory10000,
		PingDirectory10001,
		PingDirectory10002,
		PingDirectory10003,
		PingDirectory10004,
		PingDirectory10100,
		PingDirectory10102,
		PingDirectory10103,
		PingDirectory10200,
		PingDirectory10201,
		PingDirectory10300,
	}
}

func getSortedVersionsMessage() string {
	message := "Supported versions are: "
	for i, version := range getSortedVersions() {
		message += version
		if i < len(getSortedVersions())-1 {
			message += ", "
		}
	}
	return message
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

func Parse(versionString string) (string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(versionString) == 0 {
		diags.AddError("failed to parse PingDirectory version", "empty version string")
		return "", diags
	}

	versionDigits := strings.Split(versionString, ".")
	// Expect a version like "x.x" or "x.x.x.x"
	// If only two digits are supplied, the last two will be assumed to be "0.0"
	if len(versionDigits) != 2 && len(versionDigits) != 4 {
		diags.AddError("failed to parse PingDirectory version '"+versionString+"'", "Expected either two digits (e.g. '9.3') or four digits (e.g. '9.3.0.0')")
		return "", diags
	}
	if len(versionDigits) == 2 {
		versionString += ".0.0"
	}
	if !IsValid(versionString) {
		// Check if the major-minor version is valid
		majorMinorVersionString := versionDigits[0] + "." + versionDigits[1] + ".0.0"
		if !IsValid(majorMinorVersionString) {
			diags.AddError("PingDirectory version '"+versionString+"' is not supported in this version of the PingDirectory terraform provider", getSortedVersionsMessage())
			return "", diags
		}
		// The major-minor version is valid, only the patch is invalid. Warn but do not fail, assume the lastest patch version
		sortedVersions := getSortedVersions()
		versionIndex := -1
		switch majorMinorVersionString {
		case "9.3.0.0":
			// Use the first version prior to 10.0.0.0
			versionIndex = getSortedVersionIndex(PingDirectory10000) - 1
		case "10.0.0.0":
			// Use the first version prior to 10.1.0.0
			versionIndex = getSortedVersionIndex(PingDirectory10100) - 1
		case "10.1.0.0":
			// Use the first version prior to 10.2.0.0
			versionIndex = getSortedVersionIndex(PingDirectory10200) - 1
		case "10.2.0.0":
			// Use the first version prior to 10.3.0.0
			versionIndex = getSortedVersionIndex(PingDirectory10300) - 1
		case "10.3.0.0":
			// This is the latest major-minor version, so just use the latest patch version available
			versionIndex = len(sortedVersions) - 1
		}
		if versionIndex < 0 || versionIndex >= len(sortedVersions) {
			// This should never happen
			diags.AddError("Unexpected failure determining major-minor PingDirectory version", "")
			return "", diags
		}
		assumedVersion := string(sortedVersions[versionIndex])
		diags.AddWarning("PingDirectory patch version '"+versionString+"' is not recognized by this version of the PingDirectory terraform provider",
			"Assuming the latest patch version supported by the provider: '"+assumedVersion+"'")
		versionString = assumedVersion
	}
	return versionString, diags
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
