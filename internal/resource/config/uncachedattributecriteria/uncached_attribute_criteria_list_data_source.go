package uncachedattributecriteria

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
	_ datasource.DataSource              = &uncachedAttributeCriteriaListDataSource{}
	_ datasource.DataSourceWithConfigure = &uncachedAttributeCriteriaListDataSource{}
)

// Create a Uncached Attribute Criteria List data source
func NewUncachedAttributeCriteriaListDataSource() datasource.DataSource {
	return &uncachedAttributeCriteriaListDataSource{}
}

// uncachedAttributeCriteriaListDataSource is the datasource implementation.
type uncachedAttributeCriteriaListDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *uncachedAttributeCriteriaListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_uncached_attribute_criteria_list"
}

// Configure adds the provider configured client to the data source.
func (r *uncachedAttributeCriteriaListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type uncachedAttributeCriteriaListDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *uncachedAttributeCriteriaListDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists Uncached Attribute Criteria objects in the server configuration.",
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
			"objects": schema.SetAttribute{
				Description: "Uncached Attribute Criteria objects found in the configuration",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: internaltypes.ObjectsObjectType(),
			},
		},
	}
}

// Read resource information
func (r *uncachedAttributeCriteriaListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state uncachedAttributeCriteriaListDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.UncachedAttributeCriteriaApi.ListUncachedAttributeCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.ListUncachedAttributeCriteriaExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Uncached Attribute Criteria objects", err, httpResp)
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
		if response.DefaultUncachedAttributeCriteriaResponse != nil {
			attributes["id"] = types.StringValue(response.DefaultUncachedAttributeCriteriaResponse.Id)
			attributes["type"] = types.StringValue("default")
		}
		if response.SimpleUncachedAttributeCriteriaResponse != nil {
			attributes["id"] = types.StringValue(response.SimpleUncachedAttributeCriteriaResponse.Id)
			attributes["type"] = types.StringValue("simple")
		}
		if response.GroovyScriptedUncachedAttributeCriteriaResponse != nil {
			attributes["id"] = types.StringValue(response.GroovyScriptedUncachedAttributeCriteriaResponse.Id)
			attributes["type"] = types.StringValue("groovy-scripted")
		}
		if response.ThirdPartyUncachedAttributeCriteriaResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyUncachedAttributeCriteriaResponse.Id)
			attributes["type"] = types.StringValue("third-party")
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