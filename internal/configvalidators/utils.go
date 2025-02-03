// Copyright © 2025 Ping Identity Corporation

package configvalidators

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
)

func stringSliceToReadableString(slice []string) string {
	var output strings.Builder
	output.WriteRune('[')
	for i, str := range slice {
		if i > 0 {
			output.WriteString(", ")
		}
		output.WriteRune('"')
		output.WriteString(str)
		output.WriteRune('"')
	}
	output.WriteRune(']')
	return output.String()
}

func pathExpressionSliceToReadableString(slice []path.Expression) string {
	stringSlice := []string{}
	for _, pathExpr := range slice {
		stringSlice = append(stringSlice, pathExpr.String())
	}
	return stringSliceToReadableString(stringSlice)
}
