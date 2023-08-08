package resultcodemap

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
	_ datasource.DataSource              = &resultCodeMapDataSource{}
	_ datasource.DataSourceWithConfigure = &resultCodeMapDataSource{}
)

// Create a Result Code Map data source
func NewResultCodeMapDataSource() datasource.DataSource {
	return &resultCodeMapDataSource{}
}

// resultCodeMapDataSource is the datasource implementation.
type resultCodeMapDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *resultCodeMapDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_result_code_map"
}

// Configure adds the provider configured client to the data source.
func (r *resultCodeMapDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type resultCodeMapDataSourceModel struct {
	Id                            types.String `tfsdk:"id"`
	Name                          types.String `tfsdk:"name"`
	Type                          types.String `tfsdk:"type"`
	Description                   types.String `tfsdk:"description"`
	BindAccountLockedResultCode   types.Int64  `tfsdk:"bind_account_locked_result_code"`
	BindMissingUserResultCode     types.Int64  `tfsdk:"bind_missing_user_result_code"`
	BindMissingPasswordResultCode types.Int64  `tfsdk:"bind_missing_password_result_code"`
	ServerErrorResultCode         types.Int64  `tfsdk:"server_error_result_code"`
}

// GetSchema defines the schema for the datasource.
func (r *resultCodeMapDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Result Code Map.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Result Code Map resource. Options are ['result-code-map']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Result Code Map",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"bind_account_locked_result_code": schema.Int64Attribute{
				Description: "Specifies the result code that should be returned if a bind attempt fails because the user's account is locked as a result of too many failed authentication attempts.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"bind_missing_user_result_code": schema.Int64Attribute{
				Description: "Specifies the result code that should be returned if a bind attempt fails because the target user entry does not exist in the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"bind_missing_password_result_code": schema.Int64Attribute{
				Description: "Specifies the result code that should be returned if a password-based bind attempt fails because the target user entry does not have a password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_error_result_code": schema.Int64Attribute{
				Description: "Specifies the result code that should be returned if a generic error occurs within the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a ResultCodeMapResponse object into the model struct
func readResultCodeMapResponseDataSource(ctx context.Context, r *client.ResultCodeMapResponse, state *resultCodeMapDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("result-code-map")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.BindAccountLockedResultCode = internaltypes.Int64TypeOrNil(r.BindAccountLockedResultCode)
	state.BindMissingUserResultCode = internaltypes.Int64TypeOrNil(r.BindMissingUserResultCode)
	state.BindMissingPasswordResultCode = internaltypes.Int64TypeOrNil(r.BindMissingPasswordResultCode)
	state.ServerErrorResultCode = internaltypes.Int64TypeOrNil(r.ServerErrorResultCode)
}

// Read resource information
func (r *resultCodeMapDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state resultCodeMapDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ResultCodeMapApi.GetResultCodeMap(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Result Code Map", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readResultCodeMapResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
