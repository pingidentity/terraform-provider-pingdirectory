package groupimplementation

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
	_ datasource.DataSource              = &groupImplementationsDataSource{}
	_ datasource.DataSourceWithConfigure = &groupImplementationsDataSource{}
)

// Create a Group Implementations data source
func NewGroupImplementationsDataSource() datasource.DataSource {
	return &groupImplementationsDataSource{}
}

// groupImplementationsDataSource is the datasource implementation.
type groupImplementationsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *groupImplementationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_implementations"
}

// Configure adds the provider configured client to the data source.
func (r *groupImplementationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type groupImplementationsDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *groupImplementationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Group Implementation objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Group Implementation objects found in the configuration",
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
func (r *groupImplementationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state groupImplementationsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.GroupImplementationAPI.ListGroupImplementations(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.GroupImplementationAPI.ListGroupImplementationsExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Group Implementation objects", err, httpResp)
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
		if response.VirtualStaticGroupImplementationResponse != nil {
			attributes["id"] = types.StringValue(response.VirtualStaticGroupImplementationResponse.Id)
			attributes["type"] = types.StringValue("virtual-static")
		}
		if response.DynamicGroupImplementationResponse != nil {
			attributes["id"] = types.StringValue(response.DynamicGroupImplementationResponse.Id)
			attributes["type"] = types.StringValue("dynamic")
		}
		if response.StaticGroupImplementationResponse != nil {
			attributes["id"] = types.StringValue(response.StaticGroupImplementationResponse.Id)
			attributes["type"] = types.StringValue("static")
		}
		if response.InvertedStaticGroupImplementationResponse != nil {
			attributes["id"] = types.StringValue(response.InvertedStaticGroupImplementationResponse.Id)
			attributes["type"] = types.StringValue("inverted-static")
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
