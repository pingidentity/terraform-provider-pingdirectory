package conjurauthenticationmethod

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
	_ datasource.DataSource              = &conjurAuthenticationMethodDataSource{}
	_ datasource.DataSourceWithConfigure = &conjurAuthenticationMethodDataSource{}
)

// Create a Conjur Authentication Method data source
func NewConjurAuthenticationMethodDataSource() datasource.DataSource {
	return &conjurAuthenticationMethodDataSource{}
}

// conjurAuthenticationMethodDataSource is the datasource implementation.
type conjurAuthenticationMethodDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *conjurAuthenticationMethodDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_conjur_authentication_method"
}

// Configure adds the provider configured client to the data source.
func (r *conjurAuthenticationMethodDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type conjurAuthenticationMethodDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	ApiKey      types.String `tfsdk:"api_key"`
	Description types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the datasource.
func (r *conjurAuthenticationMethodDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Conjur Authentication Method.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Conjur Authentication Method resource. Options are ['api-key']",
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
				Description: "The password for the user to authenticate. This will be used to obtain an API key for the target user.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"api_key": schema.StringAttribute{
				Description: "The API key for the user to authenticate.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Conjur Authentication Method",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a ApiKeyConjurAuthenticationMethodResponse object into the model struct
func readApiKeyConjurAuthenticationMethodResponseDataSource(ctx context.Context, r *client.ApiKeyConjurAuthenticationMethodResponse, state *conjurAuthenticationMethodDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("api-key")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Username = types.StringValue(r.Username)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read resource information
func (r *conjurAuthenticationMethodDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state conjurAuthenticationMethodDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConjurAuthenticationMethodApi.GetConjurAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Conjur Authentication Method", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readApiKeyConjurAuthenticationMethodResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
