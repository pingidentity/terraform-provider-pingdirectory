package uncachedattributecriteria

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
	_ datasource.DataSource              = &uncachedAttributeCriteriaDataSource{}
	_ datasource.DataSourceWithConfigure = &uncachedAttributeCriteriaDataSource{}
)

// Create a Uncached Attribute Criteria data source
func NewUncachedAttributeCriteriaDataSource() datasource.DataSource {
	return &uncachedAttributeCriteriaDataSource{}
}

// uncachedAttributeCriteriaDataSource is the datasource implementation.
type uncachedAttributeCriteriaDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *uncachedAttributeCriteriaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_uncached_attribute_criteria"
}

// Configure adds the provider configured client to the data source.
func (r *uncachedAttributeCriteriaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type uncachedAttributeCriteriaDataSourceModel struct {
	Id                types.String `tfsdk:"id"`
	Type              types.String `tfsdk:"type"`
	ExtensionClass    types.String `tfsdk:"extension_class"`
	ExtensionArgument types.Set    `tfsdk:"extension_argument"`
	AttributeType     types.Set    `tfsdk:"attribute_type"`
	MinValueCount     types.Int64  `tfsdk:"min_value_count"`
	MinTotalValueSize types.String `tfsdk:"min_total_value_size"`
	ScriptClass       types.String `tfsdk:"script_class"`
	ScriptArgument    types.Set    `tfsdk:"script_argument"`
	Description       types.String `tfsdk:"description"`
	Enabled           types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *uncachedAttributeCriteriaDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Describes a Uncached Attribute Criteria.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Name of this object.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of Uncached Attribute Criteria resource. Options are ['default', 'groovy-scripted', 'simple', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Uncached Attribute Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Uncached Attribute Criteria. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"attribute_type": schema.SetAttribute{
				Description: "Specifies the attribute types for attributes that may be written to the uncached-id2entry database.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"min_value_count": schema.Int64Attribute{
				Description: "Specifies the minimum number of values that an attribute must have before it will be written into the uncached-id2entry database.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"min_total_value_size": schema.StringAttribute{
				Description: "Specifies the minimum total value size (i.e., the sum of the sizes of all values) that an attribute must have before it will be written into the uncached-id2entry database.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Uncached Attribute Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Uncached Attribute Criteria. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Uncached Attribute Criteria",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Uncached Attribute Criteria is enabled for use in the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
}

// Read a DefaultUncachedAttributeCriteriaResponse object into the model struct
func readDefaultUncachedAttributeCriteriaResponseDataSource(ctx context.Context, r *client.DefaultUncachedAttributeCriteriaResponse, state *uncachedAttributeCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("default")
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a GroovyScriptedUncachedAttributeCriteriaResponse object into the model struct
func readGroovyScriptedUncachedAttributeCriteriaResponseDataSource(ctx context.Context, r *client.GroovyScriptedUncachedAttributeCriteriaResponse, state *uncachedAttributeCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a SimpleUncachedAttributeCriteriaResponse object into the model struct
func readSimpleUncachedAttributeCriteriaResponseDataSource(ctx context.Context, r *client.SimpleUncachedAttributeCriteriaResponse, state *uncachedAttributeCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("simple")
	state.Id = types.StringValue(r.Id)
	state.AttributeType = internaltypes.GetStringSet(r.AttributeType)
	state.MinValueCount = internaltypes.Int64TypeOrNil(r.MinValueCount)
	state.MinTotalValueSize = internaltypes.StringTypeOrNil(r.MinTotalValueSize, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartyUncachedAttributeCriteriaResponse object into the model struct
func readThirdPartyUncachedAttributeCriteriaResponseDataSource(ctx context.Context, r *client.ThirdPartyUncachedAttributeCriteriaResponse, state *uncachedAttributeCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *uncachedAttributeCriteriaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state uncachedAttributeCriteriaDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.GetUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Uncached Attribute Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.DefaultUncachedAttributeCriteriaResponse != nil {
		readDefaultUncachedAttributeCriteriaResponseDataSource(ctx, readResponse.DefaultUncachedAttributeCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedUncachedAttributeCriteriaResponse != nil {
		readGroovyScriptedUncachedAttributeCriteriaResponseDataSource(ctx, readResponse.GroovyScriptedUncachedAttributeCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SimpleUncachedAttributeCriteriaResponse != nil {
		readSimpleUncachedAttributeCriteriaResponseDataSource(ctx, readResponse.SimpleUncachedAttributeCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyUncachedAttributeCriteriaResponse != nil {
		readThirdPartyUncachedAttributeCriteriaResponseDataSource(ctx, readResponse.ThirdPartyUncachedAttributeCriteriaResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
