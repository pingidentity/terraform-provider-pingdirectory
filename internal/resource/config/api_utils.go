package config

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	internaltypes "terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdata-config-api-go-client"
)

// Get BasicAuth context with a username and password
func BasicAuthContext(ctx context.Context, providerConfig internaltypes.ProviderConfiguration) context.Context {
	return context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{
		UserName: providerConfig.Username,
		Password: providerConfig.Password,
	})
}

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

// Write out messages from the Config API response to tflog
func LogMessages(ctx context.Context, messages *client.MetaUrnPingidentitySchemasConfigurationMessages20) {
	if messages == nil {
		return
	}

	for _, message := range messages.Notifications {
		tflog.Warn(ctx, "Configuration API Notification: "+message)
	}

	for _, action := range messages.RequiredActions {
		actionJson, err := action.MarshalJSON()
		if err != nil {
			tflog.Warn(ctx, "Configuration API RequiredAction: "+string(actionJson))
		}
	}
}

// Read messages from the Configuration API response
func ReadMessages(ctx context.Context, messages *client.MetaUrnPingidentitySchemasConfigurationMessages20) (types.Set, types.Set) {
	// Report any notifications from the Config API
	var notifications types.Set
	var requiredActions types.Set
	if messages != nil {
		notifications = internaltypes.GetStringSet(messages.Notifications)
		requiredActions, _ = GetRequiredActionsSet(*messages)
		LogMessages(ctx, messages)
	} else {
		notifications, _ = types.SetValue(types.StringType, []attr.Value{})
		requiredActions, _ = types.SetValue(GetRequiredActionsObjectType(), []attr.Value{})
	}
	return notifications, requiredActions
}
