// Copyright © 2025 Ping Identity Corporation

package workqueue

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
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
	r.apiClient = providerCfg.ApiClient
}

type workQueueResourceModel struct {
	Id                                       types.String `tfsdk:"id"`
	Notifications                            types.Set    `tfsdk:"notifications"`
	RequiredActions                          types.Set    `tfsdk:"required_actions"`
	Type                                     types.String `tfsdk:"type"`
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
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Work Queue resource. Options are ['high-throughput']",
				Optional:    false,
				Required:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"high-throughput"}...),
				},
			},
			"num_worker_threads": schema.Int64Attribute{
				Description: "Specifies the total number of worker threads that should be used within the server in order to process requested operations. The worker threads will be split evenly across all of the configured queues.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"num_write_worker_threads": schema.Int64Attribute{
				Description: "Specifies the number of worker threads that should be used within the server to process write (add, delete, modify, and modify DN) operations. If this is specified, then separate sets of worker threads will be used for processing read and write operations, and the value of the num-worker-threads property will reflect the number of threads to use to process read operations.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"num_administrative_session_worker_threads": schema.Int64Attribute{
				Description: "Specifies the number of worker threads that should be used to process operations as part of an administrative session. These threads may be reserved only for special use by management applications like dsconfig, the administration console, and other administrative tools, so that these applications may be used to diagnose problems and take any necessary corrective action even if all \"normal\" worker threads are busy processing other requests.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"num_queues": schema.Int64Attribute{
				Description: "Specifies the number of blocking queues that should be maintained. A value of zero indicates that the server should attempt to automatically select an optimal value (one queue for every two worker threads).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"num_write_queues": schema.Int64Attribute{
				Description: "Specifies the number of blocking queues that should be maintained for write operations. This will only be used if a value is specified for the num-write-worker-threads property, in which case the num-queues property will specify the number of queues for read operations. Otherwise, all operations will be processed by a common set of worker threads and the value of the num-queues property will specify the number of queues for all types of operations.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"max_work_queue_capacity": schema.Int64Attribute{
				Description: "Specifies the maximum number of pending operations that may be held in any of the queues at any given time. The total number of pending requests may be as large as this value times the total number of queues.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"max_offer_time": schema.StringAttribute{
				Description: "Specifies the maximum length of time that the connection handler should be allowed to wait to enqueue a request if the work queue is full. If the attempt to enqueue an operation does not succeed within this period of time, then the operation will be rejected and an error response will be returned to the client. A value of zero indicates that operations should be rejected immediately if the work queue is already at its maximum capacity.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"monitor_queue_time": schema.BoolAttribute{
				Description: "Indicates whether the work queue should monitor the length of time that operations are held in the queue. When enabled the queue time will be included with access log messages as \"qtime\" in milliseconds.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"max_queue_time": schema.StringAttribute{
				Description: "Specifies the maximum length of time that an operation should be allowed to wait on the work queue. If an operation has been waiting on the queue longer than this period of time, then it will receive an immediate failure result rather than being processed once it has been handed off to a worker thread. A value of zero seconds indicates that there should not be any maximum queue time imposed. This setting will only be used if the monitor-queue-time property has a value of true.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"expensive_operation_check_interval": schema.StringAttribute{
				Description: "The interval that the work queue should use when checking for potentially expensive operations. If at least expensive-operation-minimum-concurrent-count worker threads are found to be processing the same operation on two consecutive polls separated by this time interval (i.e., the worker thread has been processing that operation for at least this length of time, and potentially up to twice this length of time), then a stack trace of all running threads will be written to a file for analysis to provide potentially useful information that may help better understand the reason it is taking so long. It may be that the operation is simply an expensive one to process, but there may be other external factors (e.g., a database checkpoint, a log rotation, lock contention, etc.) that could be to blame. This option is primarily intended for debugging purposes and should generally be used under the direction of Ping Identity support.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"expensive_operation_minimum_concurrent_count": schema.Int64Attribute{
				Description: "The minimum number of concurrent expensive operations that should be detected to trigger dumping stack traces for all threads. If at least this number of worker threads are seen processing the same operations in two consecutive intervals, then the server will dump a stack trace of all threads to a file. This option is primarily intended for debugging purposes and should generally be used under the direction of Ping Identity support.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"expensive_operation_minimum_dump_interval": schema.StringAttribute{
				Description: "The minimum length of time that should be required to pass after dumping stack trace information for all threads before the server should be allowed to create a second dump. This will help prevent the server from dumping stack traces too frequently and eventually consuming all available disk space with stack trace log output. This option is primarily intended for debugging purposes and should generally be used under the direction of Ping Identity support.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a HighThroughputWorkQueueResponse object into the model struct
func readHighThroughputWorkQueueResponse(ctx context.Context, r *client.HighThroughputWorkQueueResponse, state *workQueueResourceModel, expectedValues *workQueueResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("high-throughput")
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
	var plan workQueueResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.WorkQueueAPI.GetWorkQueue(
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
	var state workQueueResourceModel
	readHighThroughputWorkQueueResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.WorkQueueAPI.UpdateWorkQueue(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	ops := createWorkQueueOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.WorkQueueAPI.UpdateWorkQueueExecute(updateRequest)
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

	readResponse, httpResp, err := r.apiClient.WorkQueueAPI.GetWorkQueue(
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
	updateRequest := r.apiClient.WorkQueueAPI.UpdateWorkQueue(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))

	// Determine what update operations are necessary
	ops := createWorkQueueOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.WorkQueueAPI.UpdateWorkQueueExecute(updateRequest)
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
