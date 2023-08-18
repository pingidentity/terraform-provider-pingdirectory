package alarmmanager

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
	_ datasource.DataSource              = &alarmManagerDataSource{}
	_ datasource.DataSourceWithConfigure = &alarmManagerDataSource{}
)

// Create a Alarm Manager data source
func NewAlarmManagerDataSource() datasource.DataSource {
	return &alarmManagerDataSource{}
}

// alarmManagerDataSource is the datasource implementation.
type alarmManagerDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *alarmManagerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alarm_manager"
}

// Configure adds the provider configured client to the data source.
func (r *alarmManagerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type alarmManagerDataSourceModel struct {
	Id                     types.String `tfsdk:"id"`
	Type                   types.String `tfsdk:"type"`
	DefaultGaugeAlertLevel types.String `tfsdk:"default_gauge_alert_level"`
	GeneratedAlertTypes    types.Set    `tfsdk:"generated_alert_types"`
	SuppressedAlarm        types.Set    `tfsdk:"suppressed_alarm"`
}

// GetSchema defines the schema for the datasource.
func (r *alarmManagerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Alarm Manager.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Alarm Manager resource. Options are ['alarm-manager']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"default_gauge_alert_level": schema.StringAttribute{
				Description: "Specifies the level at which alerts are sent for alarms raised by the Alarm Manager.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"generated_alert_types": schema.SetAttribute{
				Description: "Indicates what kind of alert types should be generated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"suppressed_alarm": schema.SetAttribute{
				Description: "Specifies the names of the alarm alert types that should be suppressed. If the condition that triggers an alarm in this list occurs, then the alarm will not be raised and no alerts will be generated. Only a subset of alarms can be suppressed in this way. Alarms triggered by a gauge can be disabled by disabling the gauge.",
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

// Read a AlarmManagerResponse object into the model struct
func readAlarmManagerResponseDataSource(ctx context.Context, r *client.AlarmManagerResponse, state *alarmManagerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("alarm-manager")
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
	state.DefaultGaugeAlertLevel = types.StringValue(r.DefaultGaugeAlertLevel.String())
	state.GeneratedAlertTypes = internaltypes.GetStringSet(
		client.StringSliceEnumalarmManagerGeneratedAlertTypesProp(r.GeneratedAlertTypes))
	state.SuppressedAlarm = internaltypes.GetStringSet(
		client.StringSliceEnumalarmManagerSuppressedAlarmProp(r.SuppressedAlarm))
}

// Read resource information
func (r *alarmManagerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state alarmManagerDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AlarmManagerApi.GetAlarmManager(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Alarm Manager", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAlarmManagerResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
