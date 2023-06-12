package workqueue

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	_ resource.Resource                = &workQueueResource{}
	_ resource.ResourceWithConfigure   = &workQueueResource{}
	_ resource.ResourceWithImportState = &workQueueResource{}
)

// Create a Work Queue resource
func NewWorkQueueResource() resource.Resource {
	return &workQueueResource{}
}

// workQueueResource is the resource implementation.
type workQueueResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *workQueueResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_work_queue"
}

// Configure adds the provider configured client to the resource.
func (r *workQueueResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type workQueueResourceModel struct {
	// Id field required for acceptance testing framework
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
}

type defaultWorkQueueResourceModel struct {
	// Id field required for acceptance testing framework
	Id                                       types.String `tfsdk:"id"`
	LastUpdated                              types.String `tfsdk:"last_updated"`
	Notifications                            types.Set    `tfsdk:"notifications"`
	RequiredActions                          types.Set    `tfsdk:"required_actions"`
	NumWorkerThreads                         types.Int64  `tfsdk:"num_worker_threads"`
	NumWriteWorkerThreads                    types.Int64  `tfsdk:"num_write_worker_threads"`
	NumAdministrativeSessionWorkerThreads    types.Int64  `tfsdk:"num_administrative_session_worker_threads"`
	NumQueues                                types.Int64  `tfsdk:"num_queues"`
	NumWriteQueues                           types.Int64  `tfsdk:"num_write_queues"`
	MaxWorkQueueCapacity                     types.Int64  `tfsdk:"max_work_queue_capacity"`
	MaxOfferTime                             types.String `tfsdk:"max_offer_time"`
	MonitorQueueTime                         types.Bool   `tfsdk:"monitor_queue_time"`
	MaxQueueTime                             types.String `tfsdk:"max_queue_time"`
	ExpensiveOperationCheckInterval          types.String `tfsdk:"expensive_operation_check_interval"`
	ExpensiveOperationMinimumConcurrentCount types.Int64  `tfsdk:"expensive_operation_minimum_concurrent_count"`
	ExpensiveOperationMinimumDumpInterval    types.String `tfsdk:"expensive_operation_minimum_dump_interval"`
}

// GetSchema defines the schema for the resource.
func (r *workQueueResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a Work Queue.",
		Attributes:  map[string]schema.Attribute{},
	}
	config.AddCommonSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a HighThroughputWorkQueueResponse object into the model struct
func readHighThroughputWorkQueueResponseDefault(ctx context.Context, r *client.HighThroughputWorkQueueResponse, state *defaultWorkQueueResourceModel, expectedValues *defaultWorkQueueResourceModel, diagnostics *diag.Diagnostics) {
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
	state.NumWorkerThreads = internaltypes.Int64TypeOrNil(r.NumWorkerThreads)
	state.NumWriteWorkerThreads = internaltypes.Int64TypeOrNil(r.NumWriteWorkerThreads)
	state.NumAdministrativeSessionWorkerThreads = internaltypes.Int64TypeOrNil(r.NumAdministrativeSessionWorkerThreads)
	state.NumQueues = internaltypes.Int64TypeOrNil(r.NumQueues)
	state.NumWriteQueues = internaltypes.Int64TypeOrNil(r.NumWriteQueues)
	state.MaxWorkQueueCapacity = internaltypes.Int64TypeOrNil(r.MaxWorkQueueCapacity)
	state.MaxOfferTime = internaltypes.StringTypeOrNil(r.MaxOfferTime, true)
	config.CheckMismatchedPDFormattedAttributes("max_offer_time",
		expectedValues.MaxOfferTime, state.MaxOfferTime, diagnostics)
	state.MonitorQueueTime = internaltypes.BoolTypeOrNil(r.MonitorQueueTime)
	state.MaxQueueTime = internaltypes.StringTypeOrNil(r.MaxQueueTime, true)
	config.CheckMismatchedPDFormattedAttributes("max_queue_time",
		expectedValues.MaxQueueTime, state.MaxQueueTime, diagnostics)
	state.ExpensiveOperationCheckInterval = internaltypes.StringTypeOrNil(r.ExpensiveOperationCheckInterval, true)
	config.CheckMismatchedPDFormattedAttributes("expensive_operation_check_interval",
		expectedValues.ExpensiveOperationCheckInterval, state.ExpensiveOperationCheckInterval, diagnostics)
	state.ExpensiveOperationMinimumConcurrentCount = internaltypes.Int64TypeOrNil(r.ExpensiveOperationMinimumConcurrentCount)
	state.ExpensiveOperationMinimumDumpInterval = internaltypes.StringTypeOrNil(r.ExpensiveOperationMinimumDumpInterval, true)
	config.CheckMismatchedPDFormattedAttributes("expensive_operation_minimum_dump_interval",
		expectedValues.ExpensiveOperationMinimumDumpInterval, state.ExpensiveOperationMinimumDumpInterval, diagnostics)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createWorkQueueOperations(plan workQueueResourceModel, state workQueueResourceModel) []client.Operation {
	var ops []client.Operation
	return ops
}

// Create any update operations necessary to make the state match the plan
func createWorkQueueOperationsDefault(plan defaultWorkQueueResourceModel, state defaultWorkQueueResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddInt64OperationIfNecessary(&ops, plan.NumWorkerThreads, state.NumWorkerThreads, "num-worker-threads")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumWriteWorkerThreads, state.NumWriteWorkerThreads, "num-write-worker-threads")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumAdministrativeSessionWorkerThreads, state.NumAdministrativeSessionWorkerThreads, "num-administrative-session-worker-threads")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumQueues, state.NumQueues, "num-queues")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumWriteQueues, state.NumWriteQueues, "num-write-queues")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxWorkQueueCapacity, state.MaxWorkQueueCapacity, "max-work-queue-capacity")
	operations.AddStringOperationIfNecessary(&ops, plan.MaxOfferTime, state.MaxOfferTime, "max-offer-time")
	operations.AddBoolOperationIfNecessary(&ops, plan.MonitorQueueTime, state.MonitorQueueTime, "monitor-queue-time")
	operations.AddStringOperationIfNecessary(&ops, plan.MaxQueueTime, state.MaxQueueTime, "max-queue-time")
	operations.AddStringOperationIfNecessary(&ops, plan.ExpensiveOperationCheckInterval, state.ExpensiveOperationCheckInterval, "expensive-operation-check-interval")
	operations.AddInt64OperationIfNecessary(&ops, plan.ExpensiveOperationMinimumConcurrentCount, state.ExpensiveOperationMinimumConcurrentCount, "expensive-operation-minimum-concurrent-count")
	operations.AddStringOperationIfNecessary(&ops, plan.ExpensiveOperationMinimumDumpInterval, state.ExpensiveOperationMinimumDumpInterval, "expensive-operation-minimum-dump-interval")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *workQueueResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan defaultWorkQueueResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.WorkQueueApi.GetWorkQueue(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Work Queue", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state defaultWorkQueueResourceModel
	readHighThroughputWorkQueueResponseDefault(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.WorkQueueApi.UpdateWorkQueue(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	ops := createWorkQueueOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.WorkQueueApi.UpdateWorkQueueExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Work Queue", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readHighThroughputWorkQueueResponseDefault(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *workQueueResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state workQueueResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.WorkQueueApi.GetWorkQueue(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Work Queue", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readHighThroughputWorkQueueResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *workQueueResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan workQueueResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state workQueueResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.WorkQueueApi.UpdateWorkQueue(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))

	// Determine what update operations are necessary
	ops := createWorkQueueOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.WorkQueueApi.UpdateWorkQueueExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Work Queue", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readHighThroughputWorkQueueResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *workQueueResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *workQueueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Set a placeholder id value to appease terraform.
	// The real attributes will be imported when terraform performs a read after the import.
	// If no value is set here, Terraform will error out when importing.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), "id")...)
}
