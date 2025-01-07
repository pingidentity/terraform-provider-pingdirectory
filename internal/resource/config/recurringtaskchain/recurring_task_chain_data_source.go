package recurringtaskchain

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
	_ datasource.DataSource              = &recurringTaskChainDataSource{}
	_ datasource.DataSourceWithConfigure = &recurringTaskChainDataSource{}
)

// Create a Recurring Task Chain data source
func NewRecurringTaskChainDataSource() datasource.DataSource {
	return &recurringTaskChainDataSource{}
}

// recurringTaskChainDataSource is the datasource implementation.
type recurringTaskChainDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *recurringTaskChainDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_recurring_task_chain"
}

// Configure adds the provider configured client to the data source.
func (r *recurringTaskChainDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type recurringTaskChainDataSourceModel struct {
	Id                               types.String `tfsdk:"id"`
	Name                             types.String `tfsdk:"name"`
	Type                             types.String `tfsdk:"type"`
	Description                      types.String `tfsdk:"description"`
	Enabled                          types.Bool   `tfsdk:"enabled"`
	RecurringTask                    types.Set    `tfsdk:"recurring_task"`
	ScheduledMonth                   types.Set    `tfsdk:"scheduled_month"`
	ScheduledDateSelectionType       types.String `tfsdk:"scheduled_date_selection_type"`
	ScheduledDayOfTheWeek            types.Set    `tfsdk:"scheduled_day_of_the_week"`
	ScheduledDayOfTheMonth           types.Set    `tfsdk:"scheduled_day_of_the_month"`
	ScheduledTimeOfDay               types.Set    `tfsdk:"scheduled_time_of_day"`
	TimeZone                         types.String `tfsdk:"time_zone"`
	InterruptedByShutdownBehavior    types.String `tfsdk:"interrupted_by_shutdown_behavior"`
	ServerOfflineAtStartTimeBehavior types.String `tfsdk:"server_offline_at_start_time_behavior"`
}

// GetSchema defines the schema for the datasource.
func (r *recurringTaskChainDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Recurring Task Chain.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Recurring Task Chain resource. Options are ['recurring-task-chain']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Recurring Task Chain",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Recurring Task Chain is enabled for use. Recurring Task Chains that are disabled will not have any new instances scheduled, but instances that are already scheduled will be preserved. Those instances may be manually canceled if desired.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"recurring_task": schema.SetAttribute{
				Description: "The set of recurring tasks that make up this chain. At least one value must be provided. If multiple values are given, then the task instances will be invoked in the order in which they are listed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"scheduled_month": schema.SetAttribute{
				Description: "The months of the year in which instances of this Recurring Task Chain may be scheduled to start.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"scheduled_date_selection_type": schema.StringAttribute{
				Description: "The mechanism used to determine the dates on which instances of this Recurring Task Chain may be scheduled to start.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"scheduled_day_of_the_week": schema.SetAttribute{
				Description: "The specific days of the week on which instances of this Recurring Task Chain may be scheduled to start. If the scheduled-day-selection-type property has a value of selected-days-of-the-week, then this property must have one or more values; otherwise, it must be left undefined.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"scheduled_day_of_the_month": schema.SetAttribute{
				Description: "The specific days of the month on which instances of this Recurring Task Chain may be scheduled to start. If the scheduled-day-selection-type property has a value of selected-days-of-the-month, then this property must have one or more values; otherwise, it must be left undefined.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"scheduled_time_of_day": schema.SetAttribute{
				Description: "The time of day at which instances of the Recurring Task Chain should be eligible to start running. Values should be in the format HH:MM (where HH is a two-digit representation of the hour of the day, between 00 and 23, inclusive), and MM is a two-digit representation of the minute of the hour (between 00 and 59, inclusive). Alternately, the value can be in the form *:MM, which indicates that the task should be eligible to start at the specified minute of every hour. At least one value must be provided, but multiple values may be given to indicate multiple start times within the same day.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"time_zone": schema.StringAttribute{
				Description: "The time zone that will be used to interpret the scheduled-time-of-day values. If no value is provided, then the JVM's default time zone will be used.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"interrupted_by_shutdown_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the server should exhibit if it is shut down or abnormally terminated while an instance of this Recurring Task Chain is running.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_offline_at_start_time_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the server should exhibit if it is offline when the start time arrives for the tasks in this Recurring Task Chain.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a RecurringTaskChainResponse object into the model struct
func readRecurringTaskChainResponseDataSource(ctx context.Context, r *client.RecurringTaskChainResponse, state *recurringTaskChainDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("recurring-task-chain")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RecurringTask = internaltypes.GetStringSet(r.RecurringTask)
	state.ScheduledMonth = internaltypes.GetStringSet(
		client.StringSliceEnumrecurringTaskChainScheduledMonthProp(r.ScheduledMonth))
	state.ScheduledDateSelectionType = types.StringValue(r.ScheduledDateSelectionType.String())
	state.ScheduledDayOfTheWeek = internaltypes.GetStringSet(
		client.StringSliceEnumrecurringTaskChainScheduledDayOfTheWeekProp(r.ScheduledDayOfTheWeek))
	state.ScheduledDayOfTheMonth = internaltypes.GetStringSet(
		client.StringSliceEnumrecurringTaskChainScheduledDayOfTheMonthProp(r.ScheduledDayOfTheMonth))
	state.ScheduledTimeOfDay = internaltypes.GetStringSet(r.ScheduledTimeOfDay)
	state.TimeZone = internaltypes.StringTypeOrNil(r.TimeZone, false)
	state.InterruptedByShutdownBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumrecurringTaskChainInterruptedByShutdownBehaviorProp(r.InterruptedByShutdownBehavior), false)
	state.ServerOfflineAtStartTimeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumrecurringTaskChainServerOfflineAtStartTimeBehaviorProp(r.ServerOfflineAtStartTimeBehavior), false)
}

// Read resource information
func (r *recurringTaskChainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state recurringTaskChainDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RecurringTaskChainAPI.GetRecurringTaskChain(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Recurring Task Chain", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readRecurringTaskChainResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
