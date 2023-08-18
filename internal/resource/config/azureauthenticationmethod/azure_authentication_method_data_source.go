package azureauthenticationmethod

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &azureAuthenticationMethodDataSource{}
	_ datasource.DataSourceWithConfigure = &azureAuthenticationMethodDataSource{}
)

// Create a Azure Authentication Method data source
func NewAzureAuthenticationMethodDataSource() datasource.DataSource {
	return &azureAuthenticationMethodDataSource{}
}

// azureAuthenticationMethodDataSource is the datasource implementation.
type azureAuthenticationMethodDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *azureAuthenticationMethodDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_azure_authentication_method"
}

// Configure adds the provider configured client to the data source.
func (r *azureAuthenticationMethodDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type azureAuthenticationMethodDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	TenantID     types.String `tfsdk:"tenant_id"`
	ClientID     types.String `tfsdk:"client_id"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Description  types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the datasource.
func (r *azureAuthenticationMethodDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Azure Authentication Method.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Azure Authentication Method resource. Options are ['default', 'client-secret', 'username-password']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"tenant_id": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`client-secret`, `username-password`]: The tenant ID to use to authenticate. When the `type` attribute is set to `default`: The tenant ID to use to authenticate. If this is not provided, then it will be obtained from the AZURE_TENANT_ID environment variable.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`client-secret`, `username-password`]: The tenant ID to use to authenticate.\n  - `default`: The tenant ID to use to authenticate. If this is not provided, then it will be obtained from the AZURE_TENANT_ID environment variable.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"client_id": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`client-secret`, `username-password`]: The client ID to use to authenticate. When the `type` attribute is set to `default`: The client ID to use to authenticate. If this is not provided, then it will be obtained from the AZURE_CLIENT_ID",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`client-secret`, `username-password`]: The client ID to use to authenticate.\n  - `default`: The client ID to use to authenticate. If this is not provided, then it will be obtained from the AZURE_CLIENT_ID",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"username": schema.StringAttribute{
				Description: "The username for the user to authenticate.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password for the user to authenticate.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"client_secret": schema.StringAttribute{
				Description: "The client secret to use to authenticate.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Azure Authentication Method",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a DefaultAzureAuthenticationMethodResponse object into the model struct
func readDefaultAzureAuthenticationMethodResponseDataSource(ctx context.Context, r *client.DefaultAzureAuthenticationMethodResponse, state *azureAuthenticationMethodDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("default")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.TenantID = internaltypes.StringTypeOrNil(r.TenantID, false)
	state.ClientID = internaltypes.StringTypeOrNil(r.ClientID, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a ClientSecretAzureAuthenticationMethodResponse object into the model struct
func readClientSecretAzureAuthenticationMethodResponseDataSource(ctx context.Context, r *client.ClientSecretAzureAuthenticationMethodResponse, state *azureAuthenticationMethodDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("client-secret")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.TenantID = types.StringValue(r.TenantID)
	state.ClientID = types.StringValue(r.ClientID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a UsernamePasswordAzureAuthenticationMethodResponse object into the model struct
func readUsernamePasswordAzureAuthenticationMethodResponseDataSource(ctx context.Context, r *client.UsernamePasswordAzureAuthenticationMethodResponse, state *azureAuthenticationMethodDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("username-password")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.TenantID = types.StringValue(r.TenantID)
	state.ClientID = types.StringValue(r.ClientID)
	state.Username = types.StringValue(r.Username)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read resource information
func (r *azureAuthenticationMethodDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state azureAuthenticationMethodDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AzureAuthenticationMethodApi.GetAzureAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Azure Authentication Method", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.DefaultAzureAuthenticationMethodResponse != nil {
		readDefaultAzureAuthenticationMethodResponseDataSource(ctx, readResponse.DefaultAzureAuthenticationMethodResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ClientSecretAzureAuthenticationMethodResponse != nil {
		readClientSecretAzureAuthenticationMethodResponseDataSource(ctx, readResponse.ClientSecretAzureAuthenticationMethodResponse, &state, &resp.Diagnostics)
	}
	if readResponse.UsernamePasswordAzureAuthenticationMethodResponse != nil {
		readUsernamePasswordAzureAuthenticationMethodResponseDataSource(ctx, readResponse.UsernamePasswordAzureAuthenticationMethodResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
