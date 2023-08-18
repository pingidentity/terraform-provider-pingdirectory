package cipherstreamprovider

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
	_ datasource.DataSource              = &cipherStreamProvidersDataSource{}
	_ datasource.DataSourceWithConfigure = &cipherStreamProvidersDataSource{}
)

// Create a Cipher Stream Providers data source
func NewCipherStreamProvidersDataSource() datasource.DataSource {
	return &cipherStreamProvidersDataSource{}
}

// cipherStreamProvidersDataSource is the datasource implementation.
type cipherStreamProvidersDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *cipherStreamProvidersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cipher_stream_providers"
}

// Configure adds the provider configured client to the data source.
func (r *cipherStreamProvidersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type cipherStreamProvidersDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *cipherStreamProvidersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Cipher Stream Provider objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Cipher Stream Provider objects found in the configuration",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: internaltypes.ObjectsObjectType(),
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read resource information
func (r *cipherStreamProvidersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state cipherStreamProvidersDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.CipherStreamProviderApi.ListCipherStreamProviders(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.ListCipherStreamProvidersExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Cipher Stream Provider objects", err, httpResp)
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
		if response.AmazonKeyManagementServiceCipherStreamProviderResponse != nil {
			attributes["id"] = types.StringValue(response.AmazonKeyManagementServiceCipherStreamProviderResponse.Id)
			attributes["type"] = types.StringValue("amazon-key-management-service")
		}
		if response.DefaultCipherStreamProviderResponse != nil {
			attributes["id"] = types.StringValue(response.DefaultCipherStreamProviderResponse.Id)
			attributes["type"] = types.StringValue("default")
		}
		if response.AmazonSecretsManagerCipherStreamProviderResponse != nil {
			attributes["id"] = types.StringValue(response.AmazonSecretsManagerCipherStreamProviderResponse.Id)
			attributes["type"] = types.StringValue("amazon-secrets-manager")
		}
		if response.AzureKeyVaultCipherStreamProviderResponse != nil {
			attributes["id"] = types.StringValue(response.AzureKeyVaultCipherStreamProviderResponse.Id)
			attributes["type"] = types.StringValue("azure-key-vault")
		}
		if response.FileBasedCipherStreamProviderResponse != nil {
			attributes["id"] = types.StringValue(response.FileBasedCipherStreamProviderResponse.Id)
			attributes["type"] = types.StringValue("file-based")
		}
		if response.WaitForPassphraseCipherStreamProviderResponse != nil {
			attributes["id"] = types.StringValue(response.WaitForPassphraseCipherStreamProviderResponse.Id)
			attributes["type"] = types.StringValue("wait-for-passphrase")
		}
		if response.ConjurCipherStreamProviderResponse != nil {
			attributes["id"] = types.StringValue(response.ConjurCipherStreamProviderResponse.Id)
			attributes["type"] = types.StringValue("conjur")
		}
		if response.Pkcs11CipherStreamProviderResponse != nil {
			attributes["id"] = types.StringValue(response.Pkcs11CipherStreamProviderResponse.Id)
			attributes["type"] = types.StringValue("pkcs11")
		}
		if response.VaultCipherStreamProviderResponse != nil {
			attributes["id"] = types.StringValue(response.VaultCipherStreamProviderResponse.Id)
			attributes["type"] = types.StringValue("vault")
		}
		if response.ThirdPartyCipherStreamProviderResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyCipherStreamProviderResponse.Id)
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
