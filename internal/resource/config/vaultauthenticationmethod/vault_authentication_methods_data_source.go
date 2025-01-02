package vaultauthenticationmethod

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &vaultAuthenticationMethodsDataSource{}
	_ datasource.DataSourceWithConfigure = &vaultAuthenticationMethodsDataSource{}
)

// Create a Vault Authentication Methods data source
func NewVaultAuthenticationMethodsDataSource() datasource.DataSource {
	return &vaultAuthenticationMethodsDataSource{}
}

// vaultAuthenticationMethodsDataSource is the datasource implementation.
type vaultAuthenticationMethodsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *vaultAuthenticationMethodsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vault_authentication_methods"
}

// Configure adds the provider configured client to the data source.
func (r *vaultAuthenticationMethodsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type vaultAuthenticationMethodsDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *vaultAuthenticationMethodsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Vault Authentication Method objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Vault Authentication Method objects found in the configuration",
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
func (r *vaultAuthenticationMethodsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state vaultAuthenticationMethodsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.VaultAuthenticationMethodAPI.ListVaultAuthenticationMethods(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.VaultAuthenticationMethodAPI.ListVaultAuthenticationMethodsExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Vault Authentication Method objects", err, httpResp)
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
		if response.AppRoleVaultAuthenticationMethodResponse != nil {
			attributes["id"] = types.StringValue(response.AppRoleVaultAuthenticationMethodResponse.Id)
			attributes["type"] = types.StringValue("app-role")
		}
		if response.UserPassVaultAuthenticationMethodResponse != nil {
			attributes["id"] = types.StringValue(response.UserPassVaultAuthenticationMethodResponse.Id)
			attributes["type"] = types.StringValue("user-pass")
		}
		if response.StaticTokenVaultAuthenticationMethodResponse != nil {
			attributes["id"] = types.StringValue(response.StaticTokenVaultAuthenticationMethodResponse.Id)
			attributes["type"] = types.StringValue("static-token")
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
