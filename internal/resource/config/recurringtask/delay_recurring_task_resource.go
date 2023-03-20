package recurringtask

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
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &delayRecurringTaskResource{}
	_ resource.ResourceWithConfigure   = &delayRecurringTaskResource{}
	_ resource.ResourceWithImportState = &delayRecurringTaskResource{}
	_ resource.Resource                = &defaultDelayRecurringTaskResource{}
	_ resource.ResourceWithConfigure   = &defaultDelayRecurringTaskResource{}
	_ resource.ResourceWithImportState = &defaultDelayRecurringTaskResource{}
)

// Create a Delay Recurring Task resource
func NewDelayRecurringTaskResource() resource.Resource {
	return &delayRecurringTaskResource{}
}

func NewDefaultDelayRecurringTaskResource() resource.Resource {
	return &defaultDelayRecurringTaskResource{}
}

// delayRecurringTaskResource is the resource implementation.
type delayRecurringTaskResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultDelayRecurringTaskResource is the resource implementation.
type defaultDelayRecurringTaskResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *delayRecurringTaskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delay_recurring_task"
}

func (r *defaultDelayRecurringTaskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_delay_recurring_task"
}

// Configure adds the provider configured client to the resource.
func (r *delayRecurringTaskResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultDelayRecurringTaskResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type delayRecurringTaskResourceModel struct {
	Id                                      types.String `tfsdk:"id"`
	LastUpdated                             types.String `tfsdk:"last_updated"`
	Notifications                           types.Set    `tfsdk:"notifications"`
	RequiredActions                         types.Set    `tfsdk:"required_actions"`
	SleepDuration                           types.String `tfsdk:"sleep_duration"`
	DurationToWaitForWorkQueueIdle          types.String `tfsdk:"duration_to_wait_for_work_queue_idle"`
	LdapURLForSearchExpectedToReturnEntries types.Set    `tfsdk:"ldap_url_for_search_expected_to_return_entries"`
	SearchInterval                          types.String `tfsdk:"search_interval"`
	SearchTimeLimit                         types.String `tfsdk:"search_time_limit"`
	DurationToWaitForSearchToReturnEntries  types.String `tfsdk:"duration_to_wait_for_search_to_return_entries"`
	TaskReturnStateIfTimeoutIsEncountered   types.String `tfsdk:"task_return_state_if_timeout_is_encountered"`
	Description                             types.String `tfsdk:"description"`
	CancelOnTaskDependencyFailure           types.Bool   `tfsdk:"cancel_on_task_dependency_failure"`
	EmailOnStart                            types.Set    `tfsdk:"email_on_start"`
	EmailOnSuccess                          types.Set    `tfsdk:"email_on_success"`
	EmailOnFailure                          types.Set    `tfsdk:"email_on_failure"`
	AlertOnStart                            types.Bool   `tfsdk:"alert_on_start"`
	AlertOnSuccess                          types.Bool   `tfsdk:"alert_on_success"`
	AlertOnFailure                          types.Bool   `tfsdk:"alert_on_failure"`
}

// GetSchema defines the schema for the resource.
func (r *delayRecurringTaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	delayRecurringTaskSchema(ctx, req, resp, false)
}

func (r *defaultDelayRecurringTaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	delayRecurringTaskSchema(ctx, req, resp, true)
}

func delayRecurringTaskSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Delay Recurring Task.",
		Attributes: map[string]schema.Attribute{
			"sleep_duration": schema.StringAttribute{
				Description: "The length of time to sleep before the task completes.",
				Optional:    true,
			},
			"duration_to_wait_for_work_queue_idle": schema.StringAttribute{
				Description: "Indicates that task should wait for up to the specified length of time for the work queue to report that all worker threads are idle and there are no pending operations. Note that this primarily monitors operations that use worker threads, which does not include internal operations (for example, those invoked by extensions), and may not include requests from non-LDAP clients (for example, HTTP-based clients).",
				Optional:    true,
			},
			"ldap_url_for_search_expected_to_return_entries": schema.SetAttribute{
				Description: "An LDAP URL that provides the criteria for a search request that is expected to return at least one entry. The search will be performed internally, and only the base DN, scope, and filter from the URL will be used; any host, port, or requested attributes included in the URL will be ignored.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"search_interval": schema.StringAttribute{
				Description: "The length of time the server should sleep between searches performed using the criteria from the ldap-url-for-search-expected-to-return-entries property.",
				Optional:    true,
			},
			"search_time_limit": schema.StringAttribute{
				Description: "The length of time that the server will wait for a response to each internal search performed using the criteria from the ldap-url-for-search-expected-to-return-entries property.",
				Optional:    true,
			},
			"duration_to_wait_for_search_to_return_entries": schema.StringAttribute{
				Description: "The maximum length of time that the server will continue to perform internal searches using the criteria from the ldap-url-for-search-expected-to-return-entries property.",
				Optional:    true,
			},
			"task_return_state_if_timeout_is_encountered": schema.StringAttribute{
				Description: "The return state to use if a timeout is encountered while waiting for the server work queue to become idle (if the duration-to-wait-for-work-queue-idle property has a value), or if the time specified by the duration-to-wait-for-search-to-return-entries elapses without the associated search returning any entries.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Recurring Task",
				Optional:    true,
			},
			"cancel_on_task_dependency_failure": schema.BoolAttribute{
				Description: "Indicates whether an instance of this Recurring Task should be canceled if the task immediately before it in the recurring task chain fails to complete successfully (including if it is canceled by an administrator before it starts or while it is running).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"email_on_start": schema.SetAttribute{
				Description: "The email addresses to which a message should be sent whenever an instance of this Recurring Task starts running. If this option is used, then at least one smtp-server must be configured in the global configuration.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"email_on_success": schema.SetAttribute{
				Description: "The email addresses to which a message should be sent whenever an instance of this Recurring Task completes successfully. If this option is used, then at least one smtp-server must be configured in the global configuration.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"email_on_failure": schema.SetAttribute{
				Description: "The email addresses to which a message should be sent if an instance of this Recurring Task fails to complete successfully. If this option is used, then at least one smtp-server must be configured in the global configuration.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"alert_on_start": schema.BoolAttribute{
				Description: "Indicates whether the server should generate an administrative alert whenever an instance of this Recurring Task starts running.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"alert_on_success": schema.BoolAttribute{
				Description: "Indicates whether the server should generate an administrative alert whenever an instance of this Recurring Task completes successfully.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"alert_on_failure": schema.BoolAttribute{
				Description: "Indicates whether the server should generate an administrative alert whenever an instance of this Recurring Task fails to complete successfully.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalDelayRecurringTaskFields(ctx context.Context, addRequest *client.AddDelayRecurringTaskRequest, plan delayRecurringTaskResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SleepDuration) {
		stringVal := plan.SleepDuration.ValueString()
		addRequest.SleepDuration = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DurationToWaitForWorkQueueIdle) {
		stringVal := plan.DurationToWaitForWorkQueueIdle.ValueString()
		addRequest.DurationToWaitForWorkQueueIdle = &stringVal
	}
	if internaltypes.IsDefined(plan.LdapURLForSearchExpectedToReturnEntries) {
		var slice []string
		plan.LdapURLForSearchExpectedToReturnEntries.ElementsAs(ctx, &slice, false)
		addRequest.LdapURLForSearchExpectedToReturnEntries = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchInterval) {
		stringVal := plan.SearchInterval.ValueString()
		addRequest.SearchInterval = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchTimeLimit) {
		stringVal := plan.SearchTimeLimit.ValueString()
		addRequest.SearchTimeLimit = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DurationToWaitForSearchToReturnEntries) {
		stringVal := plan.DurationToWaitForSearchToReturnEntries.ValueString()
		addRequest.DurationToWaitForSearchToReturnEntries = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TaskReturnStateIfTimeoutIsEncountered) {
		taskReturnStateIfTimeoutIsEncountered, err := client.NewEnumrecurringTaskTaskReturnStateIfTimeoutIsEncounteredPropFromValue(plan.TaskReturnStateIfTimeoutIsEncountered.ValueString())
		if err != nil {
			return err
		}
		addRequest.TaskReturnStateIfTimeoutIsEncountered = taskReturnStateIfTimeoutIsEncountered
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	if internaltypes.IsDefined(plan.CancelOnTaskDependencyFailure) {
		boolVal := plan.CancelOnTaskDependencyFailure.ValueBool()
		addRequest.CancelOnTaskDependencyFailure = &boolVal
	}
	if internaltypes.IsDefined(plan.EmailOnStart) {
		var slice []string
		plan.EmailOnStart.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnStart = slice
	}
	if internaltypes.IsDefined(plan.EmailOnSuccess) {
		var slice []string
		plan.EmailOnSuccess.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnSuccess = slice
	}
	if internaltypes.IsDefined(plan.EmailOnFailure) {
		var slice []string
		plan.EmailOnFailure.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnFailure = slice
	}
	if internaltypes.IsDefined(plan.AlertOnStart) {
		boolVal := plan.AlertOnStart.ValueBool()
		addRequest.AlertOnStart = &boolVal
	}
	if internaltypes.IsDefined(plan.AlertOnSuccess) {
		boolVal := plan.AlertOnSuccess.ValueBool()
		addRequest.AlertOnSuccess = &boolVal
	}
	if internaltypes.IsDefined(plan.AlertOnFailure) {
		boolVal := plan.AlertOnFailure.ValueBool()
		addRequest.AlertOnFailure = &boolVal
	}
	return nil
}

// Read a DelayRecurringTaskResponse object into the model struct
func readDelayRecurringTaskResponse(ctx context.Context, r *client.DelayRecurringTaskResponse, state *delayRecurringTaskResourceModel, expectedValues *delayRecurringTaskResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.SleepDuration = internaltypes.StringTypeOrNil(r.SleepDuration, internaltypes.IsEmptyString(expectedValues.SleepDuration))
	config.CheckMismatchedPDFormattedAttributes("sleep_duration",
		expectedValues.SleepDuration, state.SleepDuration, diagnostics)
	state.DurationToWaitForWorkQueueIdle = internaltypes.StringTypeOrNil(r.DurationToWaitForWorkQueueIdle, internaltypes.IsEmptyString(expectedValues.DurationToWaitForWorkQueueIdle))
	config.CheckMismatchedPDFormattedAttributes("duration_to_wait_for_work_queue_idle",
		expectedValues.DurationToWaitForWorkQueueIdle, state.DurationToWaitForWorkQueueIdle, diagnostics)
	state.LdapURLForSearchExpectedToReturnEntries = internaltypes.GetStringSet(r.LdapURLForSearchExpectedToReturnEntries)
	state.SearchInterval = internaltypes.StringTypeOrNil(r.SearchInterval, internaltypes.IsEmptyString(expectedValues.SearchInterval))
	config.CheckMismatchedPDFormattedAttributes("search_interval",
		expectedValues.SearchInterval, state.SearchInterval, diagnostics)
	state.SearchTimeLimit = internaltypes.StringTypeOrNil(r.SearchTimeLimit, internaltypes.IsEmptyString(expectedValues.SearchTimeLimit))
	config.CheckMismatchedPDFormattedAttributes("search_time_limit",
		expectedValues.SearchTimeLimit, state.SearchTimeLimit, diagnostics)
	state.DurationToWaitForSearchToReturnEntries = internaltypes.StringTypeOrNil(r.DurationToWaitForSearchToReturnEntries, internaltypes.IsEmptyString(expectedValues.DurationToWaitForSearchToReturnEntries))
	config.CheckMismatchedPDFormattedAttributes("duration_to_wait_for_search_to_return_entries",
		expectedValues.DurationToWaitForSearchToReturnEntries, state.DurationToWaitForSearchToReturnEntries, diagnostics)
	state.TaskReturnStateIfTimeoutIsEncountered = internaltypes.StringTypeOrNil(
		client.StringPointerEnumrecurringTaskTaskReturnStateIfTimeoutIsEncounteredProp(r.TaskReturnStateIfTimeoutIsEncountered), internaltypes.IsEmptyString(expectedValues.TaskReturnStateIfTimeoutIsEncountered))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createDelayRecurringTaskOperations(plan delayRecurringTaskResourceModel, state delayRecurringTaskResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.SleepDuration, state.SleepDuration, "sleep-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.DurationToWaitForWorkQueueIdle, state.DurationToWaitForWorkQueueIdle, "duration-to-wait-for-work-queue-idle")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.LdapURLForSearchExpectedToReturnEntries, state.LdapURLForSearchExpectedToReturnEntries, "ldap-url-for-search-expected-to-return-entries")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchInterval, state.SearchInterval, "search-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchTimeLimit, state.SearchTimeLimit, "search-time-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.DurationToWaitForSearchToReturnEntries, state.DurationToWaitForSearchToReturnEntries, "duration-to-wait-for-search-to-return-entries")
	operations.AddStringOperationIfNecessary(&ops, plan.TaskReturnStateIfTimeoutIsEncountered, state.TaskReturnStateIfTimeoutIsEncountered, "task-return-state-if-timeout-is-encountered")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.CancelOnTaskDependencyFailure, state.CancelOnTaskDependencyFailure, "cancel-on-task-dependency-failure")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EmailOnStart, state.EmailOnStart, "email-on-start")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EmailOnSuccess, state.EmailOnSuccess, "email-on-success")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EmailOnFailure, state.EmailOnFailure, "email-on-failure")
	operations.AddBoolOperationIfNecessary(&ops, plan.AlertOnStart, state.AlertOnStart, "alert-on-start")
	operations.AddBoolOperationIfNecessary(&ops, plan.AlertOnSuccess, state.AlertOnSuccess, "alert-on-success")
	operations.AddBoolOperationIfNecessary(&ops, plan.AlertOnFailure, state.AlertOnFailure, "alert-on-failure")
	return ops
}

// Create a new resource
func (r *delayRecurringTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan delayRecurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddDelayRecurringTaskRequest(plan.Id.ValueString(),
		[]client.EnumdelayRecurringTaskSchemaUrn{client.ENUMDELAYRECURRINGTASKSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RECURRING_TASKDELAY})
	err := addOptionalDelayRecurringTaskFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Delay Recurring Task", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RecurringTaskApi.AddRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRecurringTaskRequest(
		client.AddDelayRecurringTaskRequestAsAddRecurringTaskRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RecurringTaskApi.AddRecurringTaskExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Delay Recurring Task", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state delayRecurringTaskResourceModel
	readDelayRecurringTaskResponse(ctx, addResponse.DelayRecurringTaskResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultDelayRecurringTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan delayRecurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RecurringTaskApi.GetRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delay Recurring Task", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state delayRecurringTaskResourceModel
	readDelayRecurringTaskResponse(ctx, readResponse.DelayRecurringTaskResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.RecurringTaskApi.UpdateRecurringTask(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createDelayRecurringTaskOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.RecurringTaskApi.UpdateRecurringTaskExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Delay Recurring Task", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDelayRecurringTaskResponse(ctx, updateResponse.DelayRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
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
func (r *delayRecurringTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDelayRecurringTask(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDelayRecurringTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDelayRecurringTask(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readDelayRecurringTask(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state delayRecurringTaskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.RecurringTaskApi.GetRecurringTask(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delay Recurring Task", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDelayRecurringTaskResponse(ctx, readResponse.DelayRecurringTaskResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *delayRecurringTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDelayRecurringTask(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDelayRecurringTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDelayRecurringTask(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateDelayRecurringTask(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan delayRecurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state delayRecurringTaskResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.RecurringTaskApi.UpdateRecurringTask(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createDelayRecurringTaskOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.RecurringTaskApi.UpdateRecurringTaskExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Delay Recurring Task", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDelayRecurringTaskResponse(ctx, updateResponse.DelayRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultDelayRecurringTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *delayRecurringTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state delayRecurringTaskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.RecurringTaskApi.DeleteRecurringTaskExecute(r.apiClient.RecurringTaskApi.DeleteRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Delay Recurring Task", err, httpResp)
		return
	}
}

func (r *delayRecurringTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDelayRecurringTask(ctx, req, resp)
}

func (r *defaultDelayRecurringTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDelayRecurringTask(ctx, req, resp)
}

func importDelayRecurringTask(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
