package monitorprovider

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
	_ resource.Resource                = &diskSpaceUsageMonitorProviderResource{}
	_ resource.ResourceWithConfigure   = &diskSpaceUsageMonitorProviderResource{}
	_ resource.ResourceWithImportState = &diskSpaceUsageMonitorProviderResource{}
)

// Create a Disk Space Usage Monitor Provider resource
func NewDiskSpaceUsageMonitorProviderResource() resource.Resource {
	return &diskSpaceUsageMonitorProviderResource{}
}

// diskSpaceUsageMonitorProviderResource is the resource implementation.
type diskSpaceUsageMonitorProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *diskSpaceUsageMonitorProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_disk_space_usage_monitor_provider"
}

// Configure adds the provider configured client to the resource.
func (r *diskSpaceUsageMonitorProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type diskSpaceUsageMonitorProviderResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	LastUpdated                     types.String `tfsdk:"last_updated"`
	Notifications                   types.Set    `tfsdk:"notifications"`
	RequiredActions                 types.Set    `tfsdk:"required_actions"`
	LowSpaceWarningSizeThreshold    types.String `tfsdk:"low_space_warning_size_threshold"`
	LowSpaceWarningPercentThreshold types.Int64  `tfsdk:"low_space_warning_percent_threshold"`
	LowSpaceErrorSizeThreshold      types.String `tfsdk:"low_space_error_size_threshold"`
	LowSpaceErrorPercentThreshold   types.Int64  `tfsdk:"low_space_error_percent_threshold"`
	OutOfSpaceErrorSizeThreshold    types.String `tfsdk:"out_of_space_error_size_threshold"`
	OutOfSpaceErrorPercentThreshold types.Int64  `tfsdk:"out_of_space_error_percent_threshold"`
	AlertFrequency                  types.String `tfsdk:"alert_frequency"`
	Description                     types.String `tfsdk:"description"`
	Enabled                         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *diskSpaceUsageMonitorProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Disk Space Usage Monitor Provider.",
		Attributes: map[string]schema.Attribute{
			"low_space_warning_size_threshold": schema.StringAttribute{
				Description: "Specifies the low space warning threshold value as an absolute amount of space. If the amount of usable disk space drops below this amount, then the Directory Server will begin generating warning alert notifications.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"low_space_warning_percent_threshold": schema.Int64Attribute{
				Description: "Specifies the low space warning threshold value as a percentage of total space. If the amount of usable disk space drops below this amount, then the Directory Server will begin generating warning alert notifications.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"low_space_error_size_threshold": schema.StringAttribute{
				Description: "Specifies the low space error threshold value as an absolute amount of space. If the amount of usable disk space drops below this amount, then the Directory Server will start rejecting operations requested by non-root users.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"low_space_error_percent_threshold": schema.Int64Attribute{
				Description: "Specifies the low space error threshold value as a percentage of total space. If the amount of usable disk space drops below this amount, then the Directory Server will start rejecting operations requested by non-root users.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"out_of_space_error_size_threshold": schema.StringAttribute{
				Description: "Specifies the out of space error threshold value as an absolute amount of space. If the amount of usable disk space drops below this amount, then the Directory Server will shut itself down to avoid problems that may occur from complete exhaustion of usable space.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"out_of_space_error_percent_threshold": schema.Int64Attribute{
				Description: "Specifies the out of space error threshold value as a percentage of total space. If the amount of usable disk space drops below this amount, then the Directory Server will shut itself down to avoid problems that may occur from complete exhaustion of usable space.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"alert_frequency": schema.StringAttribute{
				Description: "Specifies the length of time between administrative alerts generated in response to lack of usable disk space. Administrative alerts will be generated whenever the amount of usable space drops below any threshold, and they will also be generated at regular intervals as long as the amount of usable space remains below the threshold value. A value of zero indicates that alerts should only be generated when the amount of usable space drops below a configured threshold.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Monitor Provider",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Monitor Provider is enabled for use.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Read a DiskSpaceUsageMonitorProviderResponse object into the model struct
func readDiskSpaceUsageMonitorProviderResponse(ctx context.Context, r *client.DiskSpaceUsageMonitorProviderResponse, state *diskSpaceUsageMonitorProviderResourceModel, expectedValues *diskSpaceUsageMonitorProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.LowSpaceWarningSizeThreshold = internaltypes.StringTypeOrNil(r.LowSpaceWarningSizeThreshold, true)
	config.CheckMismatchedPDFormattedAttributes("low_space_warning_size_threshold",
		expectedValues.LowSpaceWarningSizeThreshold, state.LowSpaceWarningSizeThreshold, diagnostics)
	state.LowSpaceWarningPercentThreshold = internaltypes.Int64TypeOrNil(r.LowSpaceWarningPercentThreshold)
	state.LowSpaceErrorSizeThreshold = internaltypes.StringTypeOrNil(r.LowSpaceErrorSizeThreshold, true)
	config.CheckMismatchedPDFormattedAttributes("low_space_error_size_threshold",
		expectedValues.LowSpaceErrorSizeThreshold, state.LowSpaceErrorSizeThreshold, diagnostics)
	state.LowSpaceErrorPercentThreshold = internaltypes.Int64TypeOrNil(r.LowSpaceErrorPercentThreshold)
	state.OutOfSpaceErrorSizeThreshold = internaltypes.StringTypeOrNil(r.OutOfSpaceErrorSizeThreshold, true)
	config.CheckMismatchedPDFormattedAttributes("out_of_space_error_size_threshold",
		expectedValues.OutOfSpaceErrorSizeThreshold, state.OutOfSpaceErrorSizeThreshold, diagnostics)
	state.OutOfSpaceErrorPercentThreshold = internaltypes.Int64TypeOrNil(r.OutOfSpaceErrorPercentThreshold)
	state.AlertFrequency = types.StringValue(r.AlertFrequency)
	config.CheckMismatchedPDFormattedAttributes("alert_frequency",
		expectedValues.AlertFrequency, state.AlertFrequency, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createDiskSpaceUsageMonitorProviderOperations(plan diskSpaceUsageMonitorProviderResourceModel, state diskSpaceUsageMonitorProviderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.LowSpaceWarningSizeThreshold, state.LowSpaceWarningSizeThreshold, "low-space-warning-size-threshold")
	operations.AddInt64OperationIfNecessary(&ops, plan.LowSpaceWarningPercentThreshold, state.LowSpaceWarningPercentThreshold, "low-space-warning-percent-threshold")
	operations.AddStringOperationIfNecessary(&ops, plan.LowSpaceErrorSizeThreshold, state.LowSpaceErrorSizeThreshold, "low-space-error-size-threshold")
	operations.AddInt64OperationIfNecessary(&ops, plan.LowSpaceErrorPercentThreshold, state.LowSpaceErrorPercentThreshold, "low-space-error-percent-threshold")
	operations.AddStringOperationIfNecessary(&ops, plan.OutOfSpaceErrorSizeThreshold, state.OutOfSpaceErrorSizeThreshold, "out-of-space-error-size-threshold")
	operations.AddInt64OperationIfNecessary(&ops, plan.OutOfSpaceErrorPercentThreshold, state.OutOfSpaceErrorPercentThreshold, "out-of-space-error-percent-threshold")
	operations.AddStringOperationIfNecessary(&ops, plan.AlertFrequency, state.AlertFrequency, "alert-frequency")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *diskSpaceUsageMonitorProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan diskSpaceUsageMonitorProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.MonitorProviderApi.GetMonitorProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Disk Space Usage Monitor Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state diskSpaceUsageMonitorProviderResourceModel
	readDiskSpaceUsageMonitorProviderResponse(ctx, readResponse.DiskSpaceUsageMonitorProviderResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.MonitorProviderApi.UpdateMonitorProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createDiskSpaceUsageMonitorProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.MonitorProviderApi.UpdateMonitorProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Disk Space Usage Monitor Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDiskSpaceUsageMonitorProviderResponse(ctx, updateResponse.DiskSpaceUsageMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *diskSpaceUsageMonitorProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state diskSpaceUsageMonitorProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.MonitorProviderApi.GetMonitorProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Disk Space Usage Monitor Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDiskSpaceUsageMonitorProviderResponse(ctx, readResponse.DiskSpaceUsageMonitorProviderResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *diskSpaceUsageMonitorProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan diskSpaceUsageMonitorProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state diskSpaceUsageMonitorProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.MonitorProviderApi.UpdateMonitorProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createDiskSpaceUsageMonitorProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.MonitorProviderApi.UpdateMonitorProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Disk Space Usage Monitor Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDiskSpaceUsageMonitorProviderResponse(ctx, updateResponse.DiskSpaceUsageMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *diskSpaceUsageMonitorProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *diskSpaceUsageMonitorProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
