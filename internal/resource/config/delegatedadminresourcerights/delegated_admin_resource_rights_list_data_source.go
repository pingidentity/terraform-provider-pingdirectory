package delegatedadminresourcerights

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &delegatedAdminResourceRightsListDataSource{}
	_ datasource.DataSourceWithConfigure = &delegatedAdminResourceRightsListDataSource{}
)

// Create a Delegated Admin Resource Rights List data source
func NewDelegatedAdminResourceRightsListDataSource() datasource.DataSource {
	return &delegatedAdminResourceRightsListDataSource{}
}

// delegatedAdminResourceRightsListDataSource is the datasource implementation.
type delegatedAdminResourceRightsListDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *delegatedAdminResourceRightsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delegated_admin_resource_rights_list"
}

// Configure adds the provider configured client to the data source.
func (r *delegatedAdminResourceRightsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type delegatedAdminResourceRightsListDataSourceModel struct {
	Id                       types.String `tfsdk:"id"`
	Filter                   types.String `tfsdk:"filter"`
	Ids                      types.Set    `tfsdk:"ids"`
	DelegatedAdminRightsName types.String `tfsdk:"delegated_admin_rights_name"`
}

// GetSchema defines the schema for the datasource.
func (r *delegatedAdminResourceRightsListDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Delegated Admin Resource Rights objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"delegated_admin_rights_name": schema.StringAttribute{
				Description: "Name of the parent Delegated Admin Rights",
				Required:    true,
			},
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"ids": schema.SetAttribute{
				Description: "Delegated Admin Resource Rights IDs found in the configuration",
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

// Read resource information
func (r *delegatedAdminResourceRightsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state delegatedAdminResourceRightsListDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.DelegatedAdminResourceRightsAPI.ListDelegatedAdminResourceRights(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.DelegatedAdminRightsName.ValueString())
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.DelegatedAdminResourceRightsAPI.ListDelegatedAdminResourceRightsExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Delegated Admin Resource Rights objects", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	ids := []attr.Value{}
	for _, response := range readResponse.Resources {
		ids = append(ids, types.StringValue(response.Id))
	}

	state.Ids, diags = types.SetValue(types.StringType, ids)
	resp.Diagnostics.Append(diags...)
	state.Id = types.StringValue("id")

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
