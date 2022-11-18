package pingdirectory

import (
	"context"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	client "github.com/pingidentity/pingdata-config-api-go-client"
)

// Report an HTTP error
func ReportHttpError(diagnostics *diag.Diagnostics, errorPrefix string, err error, httpResp *http.Response) {
	diagnostics.AddError(errorPrefix, err.Error())
	if httpResp != nil {
		body, err := io.ReadAll(httpResp.Body)
		if err == nil {
			diagnostics.AddError("Response body: ", string(body))
		}
	}
}

// Get BasicAuth context with a username and password
//TODO maybe cache this somehow so it doesn't need to be done so often?
func BasicAuthContext(ctx context.Context, providerConfig pingdirectoryProviderModel) context.Context {
	return context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{
		UserName: providerConfig.Username.Value,
		Password: providerConfig.Password.Value,
	})
}
