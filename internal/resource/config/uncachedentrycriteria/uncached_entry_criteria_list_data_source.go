package uncachedentrycriteria

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
	_ datasource.DataSource              = &uncachedEntryCriteriaListDataSource{}
	_ datasource.DataSourceWithConfigure = &uncachedEntryCriteriaListDataSource{}
)

// Create a Uncached Entry Criteria List data source
func NewUncachedEntryCriteriaListDataSource() datasource.DataSource {
	return &uncachedEntryCriteriaListDataSource{}
}

// uncachedEntryCriteriaListDataSource is the datasource implementation.
type uncachedEntryCriteriaListDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *uncachedEntryCriteriaListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_uncached_entry_criteria_list"
}

// Configure adds the provider configured client to the data source.
func (r *uncachedEntryCriteriaListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type uncachedEntryCriteriaListDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *uncachedEntryCriteriaListDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Uncached Entry Criteria objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Uncached Entry Criteria objects found in the configuration",
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
func (r *uncachedEntryCriteriaListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state uncachedEntryCriteriaListDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.UncachedEntryCriteriaAPI.ListUncachedEntryCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.UncachedEntryCriteriaAPI.ListUncachedEntryCriteriaExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Uncached Entry Criteria objects", err, httpResp)
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
		if response.DefaultUncachedEntryCriteriaResponse != nil {
			attributes["id"] = types.StringValue(response.DefaultUncachedEntryCriteriaResponse.Id)
			attributes["type"] = types.StringValue("default")
		}
		if response.LastAccessTimeUncachedEntryCriteriaResponse != nil {
			attributes["id"] = types.StringValue(response.LastAccessTimeUncachedEntryCriteriaResponse.Id)
			attributes["type"] = types.StringValue("last-access-time")
		}
		if response.FilterBasedUncachedEntryCriteriaResponse != nil {
			attributes["id"] = types.StringValue(response.FilterBasedUncachedEntryCriteriaResponse.Id)
			attributes["type"] = types.StringValue("filter-based")
		}
		if response.GroovyScriptedUncachedEntryCriteriaResponse != nil {
			attributes["id"] = types.StringValue(response.GroovyScriptedUncachedEntryCriteriaResponse.Id)
			attributes["type"] = types.StringValue("groovy-scripted")
		}
		if response.ThirdPartyUncachedEntryCriteriaResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyUncachedEntryCriteriaResponse.Id)
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
