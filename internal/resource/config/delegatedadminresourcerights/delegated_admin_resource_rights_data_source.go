package delegatedadminresourcerights

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
	_ datasource.DataSource              = &delegatedAdminResourceRightsDataSource{}
	_ datasource.DataSourceWithConfigure = &delegatedAdminResourceRightsDataSource{}
)

// Create a Delegated Admin Resource Rights data source
func NewDelegatedAdminResourceRightsDataSource() datasource.DataSource {
	return &delegatedAdminResourceRightsDataSource{}
}

// delegatedAdminResourceRightsDataSource is the datasource implementation.
type delegatedAdminResourceRightsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *delegatedAdminResourceRightsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delegated_admin_resource_rights"
}

// Configure adds the provider configured client to the data source.
func (r *delegatedAdminResourceRightsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type delegatedAdminResourceRightsDataSourceModel struct {
	Id                       types.String `tfsdk:"id"`
	Type                     types.String `tfsdk:"type"`
	DelegatedAdminRightsName types.String `tfsdk:"delegated_admin_rights_name"`
	Description              types.String `tfsdk:"description"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	RestResourceType         types.String `tfsdk:"rest_resource_type"`
	AdminPermission          types.Set    `tfsdk:"admin_permission"`
	AdminScope               types.String `tfsdk:"admin_scope"`
	ResourceSubtree          types.Set    `tfsdk:"resource_subtree"`
	ResourcesInGroup         types.Set    `tfsdk:"resources_in_group"`
}

// GetSchema defines the schema for the datasource.
func (r *delegatedAdminResourceRightsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Delegated Admin Resource Rights.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Delegated Admin Resource Rights resource. Options are ['delegated-admin-resource-rights']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"delegated_admin_rights_name": schema.StringAttribute{
				Description: "Name of the parent Delegated Admin Rights",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Delegated Admin Resource Rights",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether these Delegated Admin Resource Rights are enabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"rest_resource_type": schema.StringAttribute{
				Description: "Specifies the resource type applicable to these Delegated Admin Resource Rights.",
				Required:    true,
			},
			"admin_permission": schema.SetAttribute{
				Description: "Specifies administrator(s) permissions.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"admin_scope": schema.StringAttribute{
				Description: "Specifies the scope of these Delegated Admin Resource Rights.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"resource_subtree": schema.SetAttribute{
				Description: "Specifies subtrees within the search base whose entries can be managed by the administrator(s). The admin-scope must be set to resources-in-specific-subtrees.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"resources_in_group": schema.SetAttribute{
				Description: "Specifies groups whose members can be managed by the administrator(s). The admin-scope must be set to resources-in-specific-groups.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a DelegatedAdminResourceRightsResponse object into the model struct
func readDelegatedAdminResourceRightsResponseDataSource(ctx context.Context, r *client.DelegatedAdminResourceRightsResponse, state *delegatedAdminResourceRightsDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("delegated-admin-resource-rights")
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RestResourceType = types.StringValue(r.RestResourceType)
	state.AdminPermission = internaltypes.GetStringSet(
		client.StringSliceEnumdelegatedAdminResourceRightsAdminPermissionProp(r.AdminPermission))
	state.AdminScope = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdelegatedAdminResourceRightsAdminScopeProp(r.AdminScope), false)
	state.ResourceSubtree = internaltypes.GetStringSet(r.ResourceSubtree)
	state.ResourcesInGroup = internaltypes.GetStringSet(r.ResourcesInGroup)
}

// Read resource information
func (r *delegatedAdminResourceRightsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state delegatedAdminResourceRightsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DelegatedAdminResourceRightsApi.GetDelegatedAdminResourceRights(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.RestResourceType.ValueString(), state.DelegatedAdminRightsName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Resource Rights", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDelegatedAdminResourceRightsResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
