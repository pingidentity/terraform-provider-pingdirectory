package rootdnuser

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
	_ datasource.DataSource              = &rootDnUsersDataSource{}
	_ datasource.DataSourceWithConfigure = &rootDnUsersDataSource{}
)

// Create a Root Dn Users data source
func NewRootDnUsersDataSource() datasource.DataSource {
	return &rootDnUsersDataSource{}
}

// rootDnUsersDataSource is the datasource implementation.
type rootDnUsersDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *rootDnUsersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_root_dn_users"
}

// Configure adds the provider configured client to the data source.
func (r *rootDnUsersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type rootDnUsersDataSourceModel struct {
	Id     types.String `tfsdk:"id"`
	Filter types.String `tfsdk:"filter"`
	Ids    types.Set    `tfsdk:"ids"`
}

// GetSchema defines the schema for the datasource.
func (r *rootDnUsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists Root Dn User objects in the server configuration.",
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
			"ids": schema.SetAttribute{
				Description: "Root Dn User IDs found in the configuration",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Read resource information
func (r *rootDnUsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state rootDnUsersDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.RootDnUserApi.ListRootDnUsers(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.RootDnUserApi.ListRootDnUsersExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Root Dn User objects", err, httpResp)
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