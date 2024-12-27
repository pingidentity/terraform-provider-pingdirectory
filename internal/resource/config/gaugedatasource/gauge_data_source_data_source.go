package gaugedatasource

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
	_ datasource.DataSource              = &gaugeDataSourceDataSource{}
	_ datasource.DataSourceWithConfigure = &gaugeDataSourceDataSource{}
)

// Create a Gauge Data Source data source
func NewGaugeDataSourceDataSource() datasource.DataSource {
	return &gaugeDataSourceDataSource{}
}

// gaugeDataSourceDataSource is the datasource implementation.
type gaugeDataSourceDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *gaugeDataSourceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gauge_data_source"
}

// Configure adds the provider configured client to the data source.
func (r *gaugeDataSourceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type gaugeDataSourceDataSourceModel struct {
	Id                            types.String  `tfsdk:"id"`
	Name                          types.String  `tfsdk:"name"`
	Type                          types.String  `tfsdk:"type"`
	DataOrientation               types.String  `tfsdk:"data_orientation"`
	StatisticType                 types.String  `tfsdk:"statistic_type"`
	DivideValueBy                 types.Float64 `tfsdk:"divide_value_by"`
	DivideValueByAttribute        types.String  `tfsdk:"divide_value_by_attribute"`
	DivideValueByCounterAttribute types.String  `tfsdk:"divide_value_by_counter_attribute"`
	Description                   types.String  `tfsdk:"description"`
	AdditionalText                types.String  `tfsdk:"additional_text"`
	MonitorObjectclass            types.String  `tfsdk:"monitor_objectclass"`
	MonitorAttribute              types.String  `tfsdk:"monitor_attribute"`
	IncludeFilter                 types.String  `tfsdk:"include_filter"`
	ResourceAttribute             types.String  `tfsdk:"resource_attribute"`
	ResourceType                  types.String  `tfsdk:"resource_type"`
	MinimumUpdateInterval         types.String  `tfsdk:"minimum_update_interval"`
}

// GetSchema defines the schema for the datasource.
func (r *gaugeDataSourceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Gauge Data Source.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Gauge Data Source resource. Options are ['indicator', 'numeric']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"data_orientation": schema.StringAttribute{
				Description: "Indicates whether a higher or lower value is a more severe condition.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"statistic_type": schema.StringAttribute{
				Description: "Specifies the type of statistic to include in the output for the monitored attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"divide_value_by": schema.Float64Attribute{
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
			"divide_value_by_counter_attribute": schema.StringAttribute{
				Description: "An optional property that can scale the resulting value by another attribute whose value represents a counter in the monitored entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Gauge Data Source",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"additional_text": schema.StringAttribute{
				Description: "Additional information about the source of this data that is added to alerts sent as a result of gauges that use this Gauge Data Source.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"monitor_objectclass": schema.StringAttribute{
				Description: "The object class name of the monitor entries to examine for generating gauge data.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"monitor_attribute": schema.StringAttribute{
				Description: "Specifies the attribute on the monitor entries from which to derive the current gauge value.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_filter": schema.StringAttribute{
				Description: "An optional LDAP filter that can be used restrict which monitor entries are used to compute output.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"resource_attribute": schema.StringAttribute{
				Description: "Specifies the attribute whose value is used to identify the specific resource being monitored (e.g. device name).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"resource_type": schema.StringAttribute{
				Description: "A string indicating the type of resource being monitored.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"minimum_update_interval": schema.StringAttribute{
				Description: "The minimum frequency with which gauges using this Gauge Data Source can be configured for update. In order to prevent undesirable side effects, some Gauge Data Sources may use this property to impose a higher bound on the update frequency of gauges.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a IndicatorGaugeDataSourceResponse object into the model struct
func readIndicatorGaugeDataSourceResponseDataSource(ctx context.Context, r *client.IndicatorGaugeDataSourceResponse, state *gaugeDataSourceDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("indicator")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.AdditionalText = internaltypes.StringTypeOrNil(r.AdditionalText, false)
	state.MonitorObjectclass = types.StringValue(r.MonitorObjectclass)
	state.MonitorAttribute = types.StringValue(r.MonitorAttribute)
	state.IncludeFilter = internaltypes.StringTypeOrNil(r.IncludeFilter, false)
	state.ResourceAttribute = internaltypes.StringTypeOrNil(r.ResourceAttribute, false)
	state.ResourceType = internaltypes.StringTypeOrNil(r.ResourceType, false)
	state.MinimumUpdateInterval = internaltypes.StringTypeOrNil(r.MinimumUpdateInterval, false)
}

// Read a NumericGaugeDataSourceResponse object into the model struct
func readNumericGaugeDataSourceResponseDataSource(ctx context.Context, r *client.NumericGaugeDataSourceResponse, state *gaugeDataSourceDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("numeric")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.DataOrientation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeDataSourceDataOrientationProp(r.DataOrientation), false)
	state.StatisticType = types.StringValue(r.StatisticType.String())
	state.DivideValueBy = internaltypes.Float64TypeOrNil(r.DivideValueBy)
	state.DivideValueByAttribute = internaltypes.StringTypeOrNil(r.DivideValueByAttribute, false)
	state.DivideValueByCounterAttribute = internaltypes.StringTypeOrNil(r.DivideValueByCounterAttribute, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.AdditionalText = internaltypes.StringTypeOrNil(r.AdditionalText, false)
	state.MonitorObjectclass = types.StringValue(r.MonitorObjectclass)
	state.MonitorAttribute = types.StringValue(r.MonitorAttribute)
	state.IncludeFilter = internaltypes.StringTypeOrNil(r.IncludeFilter, false)
	state.ResourceAttribute = internaltypes.StringTypeOrNil(r.ResourceAttribute, false)
	state.ResourceType = internaltypes.StringTypeOrNil(r.ResourceType, false)
	state.MinimumUpdateInterval = internaltypes.StringTypeOrNil(r.MinimumUpdateInterval, false)
}

// Read resource information
func (r *gaugeDataSourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state gaugeDataSourceDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.GaugeDataSourceAPI.GetGaugeDataSource(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Gauge Data Source", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.IndicatorGaugeDataSourceResponse != nil {
		readIndicatorGaugeDataSourceResponseDataSource(ctx, readResponse.IndicatorGaugeDataSourceResponse, &state, &resp.Diagnostics)
	}
	if readResponse.NumericGaugeDataSourceResponse != nil {
		readNumericGaugeDataSourceResponseDataSource(ctx, readResponse.NumericGaugeDataSourceResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
