package config

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdata-config-api-go-client"
)

// Error returned from PingDirectory config API
type pingDirectoryError struct {
	Schemas []string `json:"schemas"`
	Status  string   `json:"status"`
	Detail  string   `json:"detail"`
}

// Report an HTTP error
func ReportHttpError(ctx context.Context, diagnostics *diag.Diagnostics, errorSummary string, err error, httpResp *http.Response) {
	httpErrorPrinted := false
	var internalError error
	if httpResp != nil {
		body, internalError := io.ReadAll(httpResp.Body)
		if internalError == nil {
			tflog.Debug(ctx, "Error HTTP response body: "+string(body))
			var pdError pingDirectoryError
			internalError = json.Unmarshal(body, &pdError)
			if internalError == nil {
				diagnostics.AddError(errorSummary, err.Error()+" - Detail: "+pdError.Detail)
				httpErrorPrinted = true
			}
		}
	}
	if !httpErrorPrinted {
		if internalError != nil {
			tflog.Warn(ctx, "Failed to unmarshal HTTP response body: "+internalError.Error())
		}
		diagnostics.AddError(errorSummary, err.Error())
	}
}

// Get BasicAuth context with a username and password
func BasicAuthContext(ctx context.Context, providerConfig types.ProviderConfiguration) context.Context {
	return context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{
		UserName: providerConfig.Username,
		Password: providerConfig.Password,
	})
}
