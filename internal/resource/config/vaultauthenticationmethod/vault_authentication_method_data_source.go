package vaultauthenticationmethod

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &vaultAuthenticationMethodDataSource{}
	_ datasource.DataSourceWithConfigure = &vaultAuthenticationMethodDataSource{}
)

// Create a Vault Authentication Method data source
func NewVaultAuthenticationMethodDataSource() datasource.DataSource {
	return &vaultAuthenticationMethodDataSource{}
}

// vaultAuthenticationMethodDataSource is the datasource implementation.
type vaultAuthenticationMethodDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *vaultAuthenticationMethodDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vault_authentication_method"
}

// Configure adds the provider configured client to the data source.
func (r *vaultAuthenticationMethodDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type vaultAuthenticationMethodDataSourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Type               types.String `tfsdk:"type"`
	Username           types.String `tfsdk:"username"`
	Password           types.String `tfsdk:"password"`
	VaultRoleID        types.String `tfsdk:"vault_role_id"`
	VaultSecretID      types.String `tfsdk:"vault_secret_id"`
	LoginMechanismName types.String `tfsdk:"login_mechanism_name"`
	VaultAccessToken   types.String `tfsdk:"vault_access_token"`
	Description        types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the datasource.
func (r *vaultAuthenticationMethodDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Vault Authentication Method.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Vault Authentication Method resource. Options are ['static-token', 'app-role', 'user-pass']",
				Required:    false,
				Optional:    false,
				Computed:    true,
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
			"vault_role_id": schema.StringAttribute{
				Description: "The role ID for the AppRole to authenticate.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"vault_secret_id": schema.StringAttribute{
				Description: "The secret ID for the AppRole to authenticate.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"login_mechanism_name": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `app-role`: The name used when enabling the desired AppRole authentication mechanism in the Vault server. When the `type` attribute is set to `user-pass`: The name used when enabling the desired UserPass authentication mechanism in the Vault server.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `app-role`: The name used when enabling the desired AppRole authentication mechanism in the Vault server.\n  - `user-pass`: The name used when enabling the desired UserPass authentication mechanism in the Vault server.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"vault_access_token": schema.StringAttribute{
				Description: "The static token used to authenticate to the Vault server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Vault Authentication Method",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a StaticTokenVaultAuthenticationMethodResponse object into the model struct
func readStaticTokenVaultAuthenticationMethodResponseDataSource(ctx context.Context, r *client.StaticTokenVaultAuthenticationMethodResponse, state *vaultAuthenticationMethodDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("static-token")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a AppRoleVaultAuthenticationMethodResponse object into the model struct
func readAppRoleVaultAuthenticationMethodResponseDataSource(ctx context.Context, r *client.AppRoleVaultAuthenticationMethodResponse, state *vaultAuthenticationMethodDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("app-role")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.VaultRoleID = types.StringValue(r.VaultRoleID)
	state.LoginMechanismName = internaltypes.StringTypeOrNil(r.LoginMechanismName, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a UserPassVaultAuthenticationMethodResponse object into the model struct
func readUserPassVaultAuthenticationMethodResponseDataSource(ctx context.Context, r *client.UserPassVaultAuthenticationMethodResponse, state *vaultAuthenticationMethodDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("user-pass")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Username = types.StringValue(r.Username)
	state.LoginMechanismName = internaltypes.StringTypeOrNil(r.LoginMechanismName, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read resource information
func (r *vaultAuthenticationMethodDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state vaultAuthenticationMethodDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.VaultAuthenticationMethodAPI.GetVaultAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Vault Authentication Method", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.StaticTokenVaultAuthenticationMethodResponse != nil {
		readStaticTokenVaultAuthenticationMethodResponseDataSource(ctx, readResponse.StaticTokenVaultAuthenticationMethodResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AppRoleVaultAuthenticationMethodResponse != nil {
		readAppRoleVaultAuthenticationMethodResponseDataSource(ctx, readResponse.AppRoleVaultAuthenticationMethodResponse, &state, &resp.Diagnostics)
	}
	if readResponse.UserPassVaultAuthenticationMethodResponse != nil {
		readUserPassVaultAuthenticationMethodResponseDataSource(ctx, readResponse.UserPassVaultAuthenticationMethodResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
