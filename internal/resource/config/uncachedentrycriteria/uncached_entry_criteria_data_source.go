package uncachedentrycriteria

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &uncachedEntryCriteriaDataSource{}
	_ datasource.DataSourceWithConfigure = &uncachedEntryCriteriaDataSource{}
)

// Create a Uncached Entry Criteria data source
func NewUncachedEntryCriteriaDataSource() datasource.DataSource {
	return &uncachedEntryCriteriaDataSource{}
}

// uncachedEntryCriteriaDataSource is the datasource implementation.
type uncachedEntryCriteriaDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *uncachedEntryCriteriaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_uncached_entry_criteria"
}

// Configure adds the provider configured client to the data source.
func (r *uncachedEntryCriteriaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type uncachedEntryCriteriaDataSourceModel struct {
	Id                              types.String `tfsdk:"id"`
	Name                            types.String `tfsdk:"name"`
	Type                            types.String `tfsdk:"type"`
	ExtensionClass                  types.String `tfsdk:"extension_class"`
	ExtensionArgument               types.Set    `tfsdk:"extension_argument"`
	ScriptClass                     types.String `tfsdk:"script_class"`
	ScriptArgument                  types.Set    `tfsdk:"script_argument"`
	Filter                          types.String `tfsdk:"filter"`
	FilterIdentifiesUncachedEntries types.Bool   `tfsdk:"filter_identifies_uncached_entries"`
	AccessTimeThreshold             types.String `tfsdk:"access_time_threshold"`
	Description                     types.String `tfsdk:"description"`
	Enabled                         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *uncachedEntryCriteriaDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Uncached Entry Criteria.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Uncached Entry Criteria resource. Options are ['default', 'last-access-time', 'filter-based', 'groovy-scripted', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Uncached Entry Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Uncached Entry Criteria. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Uncached Entry Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Uncached Entry Criteria. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"filter": schema.StringAttribute{
				Description: "Specifies the search filter that should be used to differentiate entries into cached and uncached sets.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"filter_identifies_uncached_entries": schema.BoolAttribute{
				Description: "Indicates whether the associated filter identifies those entries which should be stored in the uncached-id2entry database (if true) or entries which should be stored in the id2entry database (if false).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"access_time_threshold": schema.StringAttribute{
				Description: "Specifies the maximum length of time that has passed since an entry was last accessed that it should still be included in the id2entry database. Entries that have not been accessed in more than this length of time may be written into the uncached-id2entry database.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Uncached Entry Criteria",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Uncached Entry Criteria is enabled for use in the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a DefaultUncachedEntryCriteriaResponse object into the model struct
func readDefaultUncachedEntryCriteriaResponseDataSource(ctx context.Context, r *client.DefaultUncachedEntryCriteriaResponse, state *uncachedEntryCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("default")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a LastAccessTimeUncachedEntryCriteriaResponse object into the model struct
func readLastAccessTimeUncachedEntryCriteriaResponseDataSource(ctx context.Context, r *client.LastAccessTimeUncachedEntryCriteriaResponse, state *uncachedEntryCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("last-access-time")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AccessTimeThreshold = types.StringValue(r.AccessTimeThreshold)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a FilterBasedUncachedEntryCriteriaResponse object into the model struct
func readFilterBasedUncachedEntryCriteriaResponseDataSource(ctx context.Context, r *client.FilterBasedUncachedEntryCriteriaResponse, state *uncachedEntryCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("filter-based")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Filter = types.StringValue(r.Filter)
	state.FilterIdentifiesUncachedEntries = types.BoolValue(r.FilterIdentifiesUncachedEntries)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a GroovyScriptedUncachedEntryCriteriaResponse object into the model struct
func readGroovyScriptedUncachedEntryCriteriaResponseDataSource(ctx context.Context, r *client.GroovyScriptedUncachedEntryCriteriaResponse, state *uncachedEntryCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartyUncachedEntryCriteriaResponse object into the model struct
func readThirdPartyUncachedEntryCriteriaResponseDataSource(ctx context.Context, r *client.ThirdPartyUncachedEntryCriteriaResponse, state *uncachedEntryCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *uncachedEntryCriteriaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state uncachedEntryCriteriaDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.UncachedEntryCriteriaAPI.GetUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Uncached Entry Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.DefaultUncachedEntryCriteriaResponse != nil {
		readDefaultUncachedEntryCriteriaResponseDataSource(ctx, readResponse.DefaultUncachedEntryCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LastAccessTimeUncachedEntryCriteriaResponse != nil {
		readLastAccessTimeUncachedEntryCriteriaResponseDataSource(ctx, readResponse.LastAccessTimeUncachedEntryCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.FilterBasedUncachedEntryCriteriaResponse != nil {
		readFilterBasedUncachedEntryCriteriaResponseDataSource(ctx, readResponse.FilterBasedUncachedEntryCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedUncachedEntryCriteriaResponse != nil {
		readGroovyScriptedUncachedEntryCriteriaResponseDataSource(ctx, readResponse.GroovyScriptedUncachedEntryCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyUncachedEntryCriteriaResponse != nil {
		readThirdPartyUncachedEntryCriteriaResponseDataSource(ctx, readResponse.ThirdPartyUncachedEntryCriteriaResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
