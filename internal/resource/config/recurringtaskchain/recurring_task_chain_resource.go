package recurringtaskchain

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &recurringTaskChainResource{}
	_ resource.ResourceWithConfigure   = &recurringTaskChainResource{}
	_ resource.ResourceWithImportState = &recurringTaskChainResource{}
	_ resource.Resource                = &defaultRecurringTaskChainResource{}
	_ resource.ResourceWithConfigure   = &defaultRecurringTaskChainResource{}
	_ resource.ResourceWithImportState = &defaultRecurringTaskChainResource{}
)

// Create a Recurring Task Chain resource
func NewRecurringTaskChainResource() resource.Resource {
	return &recurringTaskChainResource{}
}

func NewDefaultRecurringTaskChainResource() resource.Resource {
	return &defaultRecurringTaskChainResource{}
}

// recurringTaskChainResource is the resource implementation.
type recurringTaskChainResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultRecurringTaskChainResource is the resource implementation.
type defaultRecurringTaskChainResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *recurringTaskChainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_recurring_task_chain"
}

func (r *defaultRecurringTaskChainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_recurring_task_chain"
}

// Configure adds the provider configured client to the resource.
func (r *recurringTaskChainResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultRecurringTaskChainResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type recurringTaskChainResourceModel struct {
	Id                               types.String `tfsdk:"id"`
	Name                             types.String `tfsdk:"name"`
	LastUpdated                      types.String `tfsdk:"last_updated"`
	Notifications                    types.Set    `tfsdk:"notifications"`
	RequiredActions                  types.Set    `tfsdk:"required_actions"`
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

// GetSchema defines the schema for the resource.
func (r *recurringTaskChainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	recurringTaskChainSchema(ctx, req, resp, false)
}

func (r *defaultRecurringTaskChainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	recurringTaskChainSchema(ctx, req, resp, true)
}

func recurringTaskChainSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Recurring Task Chain.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description: "A description for this Recurring Task Chain",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Recurring Task Chain is enabled for use. Recurring Task Chains that are disabled will not have any new instances scheduled, but instances that are already scheduled will be preserved. Those instances may be manually canceled if desired.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"recurring_task": schema.SetAttribute{
				Description: "The set of recurring tasks that make up this chain. At least one value must be provided. If multiple values are given, then the task instances will be invoked in the order in which they are listed.",
				Required:    true,
				ElementType: types.StringType,
			},
			"scheduled_month": schema.SetAttribute{
				Description: "The months of the year in which instances of this Recurring Task Chain may be scheduled to start.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"scheduled_date_selection_type": schema.StringAttribute{
				Description: "The mechanism used to determine the dates on which instances of this Recurring Task Chain may be scheduled to start.",
				Required:    true,
			},
			"scheduled_day_of_the_week": schema.SetAttribute{
				Description: "The specific days of the week on which instances of this Recurring Task Chain may be scheduled to start. If the scheduled-day-selection-type property has a value of selected-days-of-the-week, then this property must have one or more values; otherwise, it must be left undefined.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"scheduled_day_of_the_month": schema.SetAttribute{
				Description: "The specific days of the month on which instances of this Recurring Task Chain may be scheduled to start. If the scheduled-day-selection-type property has a value of selected-days-of-the-month, then this property must have one or more values; otherwise, it must be left undefined.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"scheduled_time_of_day": schema.SetAttribute{
				Description: "The time of day at which instances of the Recurring Task Chain should be eligible to start running. Values should be in the format HH:MM (where HH is a two-digit representation of the hour of the day, between 00 and 23, inclusive), and MM is a two-digit representation of the minute of the hour (between 00 and 59, inclusive). Alternately, the value can be in the form *:MM, which indicates that the task should be eligible to start at the specified minute of every hour. At least one value must be provided, but multiple values may be given to indicate multiple start times within the same day.",
				Required:    true,
				ElementType: types.StringType,
			},
			"time_zone": schema.StringAttribute{
				Description: "The time zone that will be used to interpret the scheduled-time-of-day values. If no value is provided, then the JVM's default time zone will be used.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"interrupted_by_shutdown_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the server should exhibit if it is shut down or abnormally terminated while an instance of this Recurring Task Chain is running.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_offline_at_start_time_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the server should exhibit if it is offline when the start time arrives for the tasks in this Recurring Task Chain.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if isDefault {
		// Add any default properties and set optional properties to computed where necessary
		config.SetAllAttributesToOptionalAndComputed(&schemaDef)
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add optional fields to create request for recurring-task-chain recurring-task-chain
func addOptionalRecurringTaskChainFields(ctx context.Context, addRequest *client.AddRecurringTaskChainRequest, plan recurringTaskChainResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ScheduledMonth) {
		var slice []string
		plan.ScheduledMonth.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumrecurringTaskChainScheduledMonthProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumrecurringTaskChainScheduledMonthPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.ScheduledMonth = enumSlice
	}
	if internaltypes.IsDefined(plan.ScheduledDayOfTheWeek) {
		var slice []string
		plan.ScheduledDayOfTheWeek.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumrecurringTaskChainScheduledDayOfTheWeekProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumrecurringTaskChainScheduledDayOfTheWeekPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.ScheduledDayOfTheWeek = enumSlice
	}
	if internaltypes.IsDefined(plan.ScheduledDayOfTheMonth) {
		var slice []string
		plan.ScheduledDayOfTheMonth.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumrecurringTaskChainScheduledDayOfTheMonthProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumrecurringTaskChainScheduledDayOfTheMonthPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.ScheduledDayOfTheMonth = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeZone) {
		addRequest.TimeZone = plan.TimeZone.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.InterruptedByShutdownBehavior) {
		interruptedByShutdownBehavior, err := client.NewEnumrecurringTaskChainInterruptedByShutdownBehaviorPropFromValue(plan.InterruptedByShutdownBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.InterruptedByShutdownBehavior = interruptedByShutdownBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ServerOfflineAtStartTimeBehavior) {
		serverOfflineAtStartTimeBehavior, err := client.NewEnumrecurringTaskChainServerOfflineAtStartTimeBehaviorPropFromValue(plan.ServerOfflineAtStartTimeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.ServerOfflineAtStartTimeBehavior = serverOfflineAtStartTimeBehavior
	}
	return nil
}

// Read a RecurringTaskChainResponse object into the model struct
func readRecurringTaskChainResponse(ctx context.Context, r *client.RecurringTaskChainResponse, state *recurringTaskChainResourceModel, expectedValues *recurringTaskChainResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
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
	state.TimeZone = internaltypes.StringTypeOrNil(r.TimeZone, internaltypes.IsEmptyString(expectedValues.TimeZone))
	state.InterruptedByShutdownBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumrecurringTaskChainInterruptedByShutdownBehaviorProp(r.InterruptedByShutdownBehavior), internaltypes.IsEmptyString(expectedValues.InterruptedByShutdownBehavior))
	state.ServerOfflineAtStartTimeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumrecurringTaskChainServerOfflineAtStartTimeBehaviorProp(r.ServerOfflineAtStartTimeBehavior), internaltypes.IsEmptyString(expectedValues.ServerOfflineAtStartTimeBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createRecurringTaskChainOperations(plan recurringTaskChainResourceModel, state recurringTaskChainResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RecurringTask, state.RecurringTask, "recurring-task")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScheduledMonth, state.ScheduledMonth, "scheduled-month")
	operations.AddStringOperationIfNecessary(&ops, plan.ScheduledDateSelectionType, state.ScheduledDateSelectionType, "scheduled-date-selection-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScheduledDayOfTheWeek, state.ScheduledDayOfTheWeek, "scheduled-day-of-the-week")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScheduledDayOfTheMonth, state.ScheduledDayOfTheMonth, "scheduled-day-of-the-month")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScheduledTimeOfDay, state.ScheduledTimeOfDay, "scheduled-time-of-day")
	operations.AddStringOperationIfNecessary(&ops, plan.TimeZone, state.TimeZone, "time-zone")
	operations.AddStringOperationIfNecessary(&ops, plan.InterruptedByShutdownBehavior, state.InterruptedByShutdownBehavior, "interrupted-by-shutdown-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerOfflineAtStartTimeBehavior, state.ServerOfflineAtStartTimeBehavior, "server-offline-at-start-time-behavior")
	return ops
}

// Create a recurring-task-chain recurring-task-chain
func (r *recurringTaskChainResource) CreateRecurringTaskChain(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan recurringTaskChainResourceModel) (*recurringTaskChainResourceModel, error) {
	var RecurringTaskSlice []string
	plan.RecurringTask.ElementsAs(ctx, &RecurringTaskSlice, false)
	scheduledDateSelectionType, err := client.NewEnumrecurringTaskChainScheduledDateSelectionTypePropFromValue(plan.ScheduledDateSelectionType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for ScheduledDateSelectionType", err.Error())
		return nil, err
	}
	var ScheduledTimeOfDaySlice []string
	plan.ScheduledTimeOfDay.ElementsAs(ctx, &ScheduledTimeOfDaySlice, false)
	addRequest := client.NewAddRecurringTaskChainRequest(plan.Name.ValueString(),
		RecurringTaskSlice,
		*scheduledDateSelectionType,
		ScheduledTimeOfDaySlice)
	err = addOptionalRecurringTaskChainFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Recurring Task Chain", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RecurringTaskChainApi.AddRecurringTaskChain(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRecurringTaskChainRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.RecurringTaskChainApi.AddRecurringTaskChainExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Recurring Task Chain", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state recurringTaskChainResourceModel
	readRecurringTaskChainResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *recurringTaskChainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan recurringTaskChainResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateRecurringTaskChain(ctx, req, resp, plan)
	if err != nil {
		return
	}

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, *state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultRecurringTaskChainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan recurringTaskChainResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RecurringTaskChainApi.GetRecurringTaskChain(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Recurring Task Chain", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state recurringTaskChainResourceModel
	readRecurringTaskChainResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.RecurringTaskChainApi.UpdateRecurringTaskChain(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createRecurringTaskChainOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.RecurringTaskChainApi.UpdateRecurringTaskChainExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Recurring Task Chain", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readRecurringTaskChainResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *recurringTaskChainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readRecurringTaskChain(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultRecurringTaskChainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readRecurringTaskChain(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readRecurringTaskChain(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state recurringTaskChainResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.RecurringTaskChainApi.GetRecurringTaskChain(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
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
	readRecurringTaskChainResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *recurringTaskChainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateRecurringTaskChain(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultRecurringTaskChainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateRecurringTaskChain(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateRecurringTaskChain(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan recurringTaskChainResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state recurringTaskChainResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.RecurringTaskChainApi.UpdateRecurringTaskChain(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createRecurringTaskChainOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.RecurringTaskChainApi.UpdateRecurringTaskChainExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Recurring Task Chain", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readRecurringTaskChainResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// This config object is edit-only, so Terraform can't delete it.
// After running a delete, Terraform will just "forget" about this object and it can be managed elsewhere.
func (r *defaultRecurringTaskChainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *recurringTaskChainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state recurringTaskChainResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.RecurringTaskChainApi.DeleteRecurringTaskChainExecute(r.apiClient.RecurringTaskChainApi.DeleteRecurringTaskChain(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Recurring Task Chain", err, httpResp)
		return
	}
}

func (r *recurringTaskChainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importRecurringTaskChain(ctx, req, resp)
}

func (r *defaultRecurringTaskChainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importRecurringTaskChain(ctx, req, resp)
}

func importRecurringTaskChain(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
