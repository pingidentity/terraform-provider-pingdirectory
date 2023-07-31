package passphraseprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &passphraseProvidersDataSource{}
	_ datasource.DataSourceWithConfigure = &passphraseProvidersDataSource{}
)

// Create a Passphrase Providers data source
func NewPassphraseProvidersDataSource() datasource.DataSource {
	return &passphraseProvidersDataSource{}
}

// passphraseProvidersDataSource is the datasource implementation.
type passphraseProvidersDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *passphraseProvidersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_passphrase_providers"
}

// Configure adds the provider configured client to the data source.
func (r *passphraseProvidersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type passphraseProvidersDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *passphraseProvidersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists Passphrase Provider objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder name of this object required by Terraform.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Passphrase Provider objects found in the configuration",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: internaltypes.ObjectsObjectType(),
			},
		},
	}
}

// Read resource information
func (r *passphraseProvidersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state passphraseProvidersDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.PassphraseProviderApi.ListPassphraseProviders(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.PassphraseProviderApi.ListPassphraseProvidersExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Passphrase Provider objects", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	objects := []attr.Value{}
	for _, response := range readResponse.Resources {
		attributes := map[string]attr.Value{}
		if response.EnvironmentVariablePassphraseProviderResponse != nil {
			attributes["id"] = types.StringValue(response.EnvironmentVariablePassphraseProviderResponse.Id)
			attributes["type"] = types.StringValue("environment-variable")
		}
		if response.AmazonSecretsManagerPassphraseProviderResponse != nil {
			attributes["id"] = types.StringValue(response.AmazonSecretsManagerPassphraseProviderResponse.Id)
			attributes["type"] = types.StringValue("amazon-secrets-manager")
		}
		if response.ObscuredValuePassphraseProviderResponse != nil {
			attributes["id"] = types.StringValue(response.ObscuredValuePassphraseProviderResponse.Id)
			attributes["type"] = types.StringValue("obscured-value")
		}
		if response.AzureKeyVaultPassphraseProviderResponse != nil {
			attributes["id"] = types.StringValue(response.AzureKeyVaultPassphraseProviderResponse.Id)
			attributes["type"] = types.StringValue("azure-key-vault")
		}
		if response.FileBasedPassphraseProviderResponse != nil {
			attributes["id"] = types.StringValue(response.FileBasedPassphraseProviderResponse.Id)
			attributes["type"] = types.StringValue("file-based")
		}
		if response.ConjurPassphraseProviderResponse != nil {
			attributes["id"] = types.StringValue(response.ConjurPassphraseProviderResponse.Id)
			attributes["type"] = types.StringValue("conjur")
		}
		if response.VaultPassphraseProviderResponse != nil {
			attributes["id"] = types.StringValue(response.VaultPassphraseProviderResponse.Id)
			attributes["type"] = types.StringValue("vault")
		}
		if response.ThirdPartyPassphraseProviderResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyPassphraseProviderResponse.Id)
			attributes["type"] = types.StringValue("third-party")
		}
		obj, diags := types.ObjectValue(internaltypes.ObjectsAttrTypes(), attributes)
		resp.Diagnostics.Append(diags...)
		objects = append(objects, obj)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	state.Objects, diags = types.SetValue(internaltypes.ObjectsObjectType(), objects)
	resp.Diagnostics.Append(diags...)
	state.Id = types.StringValue("id")

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
