package logrotationpolicy

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &fixedTimeLogRotationPolicyResource{}
	_ resource.ResourceWithConfigure   = &fixedTimeLogRotationPolicyResource{}
	_ resource.ResourceWithImportState = &fixedTimeLogRotationPolicyResource{}
	_ resource.Resource                = &defaultFixedTimeLogRotationPolicyResource{}
	_ resource.ResourceWithConfigure   = &defaultFixedTimeLogRotationPolicyResource{}
	_ resource.ResourceWithImportState = &defaultFixedTimeLogRotationPolicyResource{}
)

// Create a Fixed Time Log Rotation Policy resource
func NewFixedTimeLogRotationPolicyResource() resource.Resource {
	return &fixedTimeLogRotationPolicyResource{}
}

func NewDefaultFixedTimeLogRotationPolicyResource() resource.Resource {
	return &defaultFixedTimeLogRotationPolicyResource{}
}

// fixedTimeLogRotationPolicyResource is the resource implementation.
type fixedTimeLogRotationPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultFixedTimeLogRotationPolicyResource is the resource implementation.
type defaultFixedTimeLogRotationPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *fixedTimeLogRotationPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fixed_time_log_rotation_policy"
}

func (r *defaultFixedTimeLogRotationPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_fixed_time_log_rotation_policy"
}

// Configure adds the provider configured client to the resource.
func (r *fixedTimeLogRotationPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultFixedTimeLogRotationPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type fixedTimeLogRotationPolicyResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	TimeOfDay       types.Set    `tfsdk:"time_of_day"`
	Description     types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *fixedTimeLogRotationPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	fixedTimeLogRotationPolicySchema(ctx, req, resp, false)
}

func (r *defaultFixedTimeLogRotationPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	fixedTimeLogRotationPolicySchema(ctx, req, resp, true)
}

func fixedTimeLogRotationPolicySchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Fixed Time Log Rotation Policy.",
		Attributes: map[string]schema.Attribute{
			"time_of_day": schema.SetAttribute{
				Description: "Specifies the time of day at which log rotation should occur.",
				Required:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Rotation Policy",
				Optional:    true,
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
func addOptionalFixedTimeLogRotationPolicyFields(ctx context.Context, addRequest *client.AddFixedTimeLogRotationPolicyRequest, plan fixedTimeLogRotationPolicyResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a FixedTimeLogRotationPolicyResponse object into the model struct
func readFixedTimeLogRotationPolicyResponse(ctx context.Context, r *client.FixedTimeLogRotationPolicyResponse, state *fixedTimeLogRotationPolicyResourceModel, expectedValues *fixedTimeLogRotationPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.TimeOfDay = internaltypes.GetStringSet(r.TimeOfDay)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createFixedTimeLogRotationPolicyOperations(plan fixedTimeLogRotationPolicyResourceModel, state fixedTimeLogRotationPolicyResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.TimeOfDay, state.TimeOfDay, "time-of-day")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *fixedTimeLogRotationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan fixedTimeLogRotationPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var TimeOfDaySlice []string
	plan.TimeOfDay.ElementsAs(ctx, &TimeOfDaySlice, false)
	addRequest := client.NewAddFixedTimeLogRotationPolicyRequest(plan.Id.ValueString(),
		[]client.EnumfixedTimeLogRotationPolicySchemaUrn{client.ENUMFIXEDTIMELOGROTATIONPOLICYSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_ROTATION_POLICYFIXED_TIME},
		TimeOfDaySlice)
	addOptionalFixedTimeLogRotationPolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogRotationPolicyApi.AddLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogRotationPolicyRequest(
		client.AddFixedTimeLogRotationPolicyRequestAsAddLogRotationPolicyRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogRotationPolicyApi.AddLogRotationPolicyExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Fixed Time Log Rotation Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state fixedTimeLogRotationPolicyResourceModel
	readFixedTimeLogRotationPolicyResponse(ctx, addResponse.FixedTimeLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultFixedTimeLogRotationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan fixedTimeLogRotationPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogRotationPolicyApi.GetLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Fixed Time Log Rotation Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state fixedTimeLogRotationPolicyResourceModel
	readFixedTimeLogRotationPolicyResponse(ctx, readResponse.FixedTimeLogRotationPolicyResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogRotationPolicyApi.UpdateLogRotationPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createFixedTimeLogRotationPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogRotationPolicyApi.UpdateLogRotationPolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Fixed Time Log Rotation Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFixedTimeLogRotationPolicyResponse(ctx, updateResponse.FixedTimeLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
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
func (r *fixedTimeLogRotationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFixedTimeLogRotationPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFixedTimeLogRotationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFixedTimeLogRotationPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readFixedTimeLogRotationPolicy(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state fixedTimeLogRotationPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogRotationPolicyApi.GetLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Fixed Time Log Rotation Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readFixedTimeLogRotationPolicyResponse(ctx, readResponse.FixedTimeLogRotationPolicyResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *fixedTimeLogRotationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFixedTimeLogRotationPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFixedTimeLogRotationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFixedTimeLogRotationPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateFixedTimeLogRotationPolicy(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan fixedTimeLogRotationPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state fixedTimeLogRotationPolicyResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogRotationPolicyApi.UpdateLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createFixedTimeLogRotationPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogRotationPolicyApi.UpdateLogRotationPolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Fixed Time Log Rotation Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFixedTimeLogRotationPolicyResponse(ctx, updateResponse.FixedTimeLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultFixedTimeLogRotationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *fixedTimeLogRotationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state fixedTimeLogRotationPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogRotationPolicyApi.DeleteLogRotationPolicyExecute(r.apiClient.LogRotationPolicyApi.DeleteLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Fixed Time Log Rotation Policy", err, httpResp)
		return
	}
}

func (r *fixedTimeLogRotationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFixedTimeLogRotationPolicy(ctx, req, resp)
}

func (r *defaultFixedTimeLogRotationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFixedTimeLogRotationPolicy(ctx, req, resp)
}

func importFixedTimeLogRotationPolicy(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
