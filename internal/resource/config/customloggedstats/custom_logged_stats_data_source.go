package customloggedstats

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
	_ datasource.DataSource              = &customLoggedStatsDataSource{}
	_ datasource.DataSourceWithConfigure = &customLoggedStatsDataSource{}
)

// Create a Custom Logged Stats data source
func NewCustomLoggedStatsDataSource() datasource.DataSource {
	return &customLoggedStatsDataSource{}
}

// customLoggedStatsDataSource is the datasource implementation.
type customLoggedStatsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *customLoggedStatsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_logged_stats"
}

// Configure adds the provider configured client to the data source.
func (r *customLoggedStatsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type customLoggedStatsDataSourceModel struct {
	Id                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Type                   types.String `tfsdk:"type"`
	PluginName             types.String `tfsdk:"plugin_name"`
	Description            types.String `tfsdk:"description"`
	Enabled                types.Bool   `tfsdk:"enabled"`
	MonitorObjectclass     types.String `tfsdk:"monitor_objectclass"`
	IncludeFilter          types.String `tfsdk:"include_filter"`
	AttributeToLog         types.Set    `tfsdk:"attribute_to_log"`
	ColumnName             types.Set    `tfsdk:"column_name"`
	StatisticType          types.Set    `tfsdk:"statistic_type"`
	HeaderPrefix           types.String `tfsdk:"header_prefix"`
	HeaderPrefixAttribute  types.String `tfsdk:"header_prefix_attribute"`
	RegexPattern           types.String `tfsdk:"regex_pattern"`
	RegexReplacement       types.String `tfsdk:"regex_replacement"`
	DivideValueBy          types.String `tfsdk:"divide_value_by"`
	DivideValueByAttribute types.String `tfsdk:"divide_value_by_attribute"`
	DecimalFormat          types.String `tfsdk:"decimal_format"`
	NonZeroImpliesNotIdle  types.Bool   `tfsdk:"non_zero_implies_not_idle"`
}

// GetSchema defines the schema for the datasource.
func (r *customLoggedStatsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Custom Logged Stats.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Custom Logged Stats resource. Options are ['custom-logged-stats']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_name": schema.StringAttribute{
				Description: "Name of the parent Plugin",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Custom Logged Stats",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Custom Logged Stats object is enabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"monitor_objectclass": schema.StringAttribute{
				Description: "The objectclass name of the monitor entries to examine for generating these statistics.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_filter": schema.StringAttribute{
				Description: "An optional LDAP filter that can be used restrict which monitor entries are used to produce the output.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"attribute_to_log": schema.SetAttribute{
				Description: "Specifies the attributes on the monitor entries that should be included in the output.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"column_name": schema.SetAttribute{
				Description: "Optionally, specifies an explicit name for each column header instead of having these names automatically generated from the monitored attribute name.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"statistic_type": schema.SetAttribute{
				Description: "Specifies the type of statistic to include in the output for each monitored attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"header_prefix": schema.StringAttribute{
				Description: "An optional prefix that is included in the header before the column name.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"header_prefix_attribute": schema.StringAttribute{
				Description: "An optional attribute from the monitor entry that is included as a prefix before the column name in the column header.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"regex_pattern": schema.StringAttribute{
				Description: "An optional regular expression pattern, that when used in conjunction with regex-replacement, can alter the value of the attribute being monitored.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"regex_replacement": schema.StringAttribute{
				Description: "An optional regular expression replacement value, that when used in conjunction with regex-pattern, can alter the value of the attribute being monitored.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"divide_value_by": schema.StringAttribute{
				Description: "An optional floating point value that can be used to scale the resulting value.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"divide_value_by_attribute": schema.StringAttribute{
				Description: "An optional property that can scale the resulting value by another attribute in the monitored entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"decimal_format": schema.StringAttribute{
				Description: "This provides a way to format the monitored attribute value in the output to control the precision for instance.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"non_zero_implies_not_idle": schema.BoolAttribute{
				Description: "If this property is set to true, then the value of any of the monitored attributes here can contribute to whether an interval is considered \"idle\" by the Periodic Stats Logger.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a CustomLoggedStatsResponse object into the model struct
func readCustomLoggedStatsResponseDataSource(ctx context.Context, r *client.CustomLoggedStatsResponse, state *customLoggedStatsDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("custom-logged-stats")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.MonitorObjectclass = types.StringValue(r.MonitorObjectclass)
	state.IncludeFilter = internaltypes.StringTypeOrNil(r.IncludeFilter, false)
	state.AttributeToLog = internaltypes.GetStringSet(r.AttributeToLog)
	state.ColumnName = internaltypes.GetStringSet(r.ColumnName)
	state.StatisticType = internaltypes.GetStringSet(
		client.StringSliceEnumcustomLoggedStatsStatisticTypeProp(r.StatisticType))
	state.HeaderPrefix = internaltypes.StringTypeOrNil(r.HeaderPrefix, false)
	state.HeaderPrefixAttribute = internaltypes.StringTypeOrNil(r.HeaderPrefixAttribute, false)
	state.RegexPattern = internaltypes.StringTypeOrNil(r.RegexPattern, false)
	state.RegexReplacement = internaltypes.StringTypeOrNil(r.RegexReplacement, false)
	state.DivideValueBy = internaltypes.StringTypeOrNil(r.DivideValueBy, false)
	state.DivideValueByAttribute = internaltypes.StringTypeOrNil(r.DivideValueByAttribute, false)
	state.DecimalFormat = internaltypes.StringTypeOrNil(r.DecimalFormat, false)
	state.NonZeroImpliesNotIdle = internaltypes.BoolTypeOrNil(r.NonZeroImpliesNotIdle)
}

// Read resource information
func (r *customLoggedStatsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state customLoggedStatsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CustomLoggedStatsAPI.GetCustomLoggedStats(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.PluginName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Custom Logged Stats", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readCustomLoggedStatsResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
