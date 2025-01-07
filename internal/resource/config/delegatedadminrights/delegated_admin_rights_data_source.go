package delegatedadminrights

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &delegatedAdminRightsDataSource{}
	_ datasource.DataSourceWithConfigure = &delegatedAdminRightsDataSource{}
)

// Create a Delegated Admin Rights data source
func NewDelegatedAdminRightsDataSource() datasource.DataSource {
	return &delegatedAdminRightsDataSource{}
}

// delegatedAdminRightsDataSource is the datasource implementation.
type delegatedAdminRightsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *delegatedAdminRightsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delegated_admin_rights"
}

// Configure adds the provider configured client to the data source.
func (r *delegatedAdminRightsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type delegatedAdminRightsDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	Description  types.String `tfsdk:"description"`
	Enabled      types.Bool   `tfsdk:"enabled"`
	AdminUserDN  types.String `tfsdk:"admin_user_dn"`
	AdminGroupDN types.String `tfsdk:"admin_group_dn"`
}

// GetSchema defines the schema for the datasource.
func (r *delegatedAdminRightsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Delegated Admin Rights.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Delegated Admin Rights resource. Options are ['delegated-admin-rights']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Delegated Admin Rights",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Delegated Admin Rights is enabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"admin_user_dn": schema.StringAttribute{
				Description: "Specifies the DN of an administrative user who has authority to manage resources. Either admin-user-dn or admin-group-dn must be specified, but not both.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"admin_group_dn": schema.StringAttribute{
				Description: "Specifies the DN of a group of administrative users who have authority to manage resources. Either admin-user-dn or admin-group-dn must be specified, but not both.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a DelegatedAdminRightsResponse object into the model struct
func readDelegatedAdminRightsResponseDataSource(ctx context.Context, r *client.DelegatedAdminRightsResponse, state *delegatedAdminRightsDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("delegated-admin-rights")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AdminUserDN = internaltypes.StringTypeOrNil(r.AdminUserDN, false)
	state.AdminGroupDN = internaltypes.StringTypeOrNil(r.AdminGroupDN, false)
}

// Read resource information
func (r *delegatedAdminRightsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state delegatedAdminRightsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DelegatedAdminRightsAPI.GetDelegatedAdminRights(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Rights", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDelegatedAdminRightsResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
