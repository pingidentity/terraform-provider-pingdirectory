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
	_ resource.Resource                = &timeLimitLogRotationPolicyResource{}
	_ resource.ResourceWithConfigure   = &timeLimitLogRotationPolicyResource{}
	_ resource.ResourceWithImportState = &timeLimitLogRotationPolicyResource{}
	_ resource.Resource                = &defaultTimeLimitLogRotationPolicyResource{}
	_ resource.ResourceWithConfigure   = &defaultTimeLimitLogRotationPolicyResource{}
	_ resource.ResourceWithImportState = &defaultTimeLimitLogRotationPolicyResource{}
)

// Create a Time Limit Log Rotation Policy resource
func NewTimeLimitLogRotationPolicyResource() resource.Resource {
	return &timeLimitLogRotationPolicyResource{}
}

func NewDefaultTimeLimitLogRotationPolicyResource() resource.Resource {
	return &defaultTimeLimitLogRotationPolicyResource{}
}

// timeLimitLogRotationPolicyResource is the resource implementation.
type timeLimitLogRotationPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultTimeLimitLogRotationPolicyResource is the resource implementation.
type defaultTimeLimitLogRotationPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *timeLimitLogRotationPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_time_limit_log_rotation_policy"
}

func (r *defaultTimeLimitLogRotationPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_time_limit_log_rotation_policy"
}

// Configure adds the provider configured client to the resource.
func (r *timeLimitLogRotationPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultTimeLimitLogRotationPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type timeLimitLogRotationPolicyResourceModel struct {
	Id               types.String `tfsdk:"id"`
	LastUpdated      types.String `tfsdk:"last_updated"`
	Notifications    types.Set    `tfsdk:"notifications"`
	RequiredActions  types.Set    `tfsdk:"required_actions"`
	RotationInterval types.String `tfsdk:"rotation_interval"`
	Description      types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *timeLimitLogRotationPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	timeLimitLogRotationPolicySchema(ctx, req, resp, false)
}

func (r *defaultTimeLimitLogRotationPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	timeLimitLogRotationPolicySchema(ctx, req, resp, true)
}

func timeLimitLogRotationPolicySchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Time Limit Log Rotation Policy.",
		Attributes: map[string]schema.Attribute{
			"rotation_interval": schema.StringAttribute{
				Description: "Specifies the time interval between rotations.",
				Required:    true,
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
func addOptionalTimeLimitLogRotationPolicyFields(ctx context.Context, addRequest *client.AddTimeLimitLogRotationPolicyRequest, plan timeLimitLogRotationPolicyResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a TimeLimitLogRotationPolicyResponse object into the model struct
func readTimeLimitLogRotationPolicyResponse(ctx context.Context, r *client.TimeLimitLogRotationPolicyResponse, state *timeLimitLogRotationPolicyResourceModel, expectedValues *timeLimitLogRotationPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.RotationInterval = types.StringValue(r.RotationInterval)
	config.CheckMismatchedPDFormattedAttributes("rotation_interval",
		expectedValues.RotationInterval, state.RotationInterval, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createTimeLimitLogRotationPolicyOperations(plan timeLimitLogRotationPolicyResourceModel, state timeLimitLogRotationPolicyResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.RotationInterval, state.RotationInterval, "rotation-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *timeLimitLogRotationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan timeLimitLogRotationPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddTimeLimitLogRotationPolicyRequest(plan.Id.ValueString(),
		[]client.EnumtimeLimitLogRotationPolicySchemaUrn{client.ENUMTIMELIMITLOGROTATIONPOLICYSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_ROTATION_POLICYTIME_LIMIT},
		plan.RotationInterval.ValueString())
	addOptionalTimeLimitLogRotationPolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogRotationPolicyApi.AddLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogRotationPolicyRequest(
		client.AddTimeLimitLogRotationPolicyRequestAsAddLogRotationPolicyRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogRotationPolicyApi.AddLogRotationPolicyExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Time Limit Log Rotation Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state timeLimitLogRotationPolicyResourceModel
	readTimeLimitLogRotationPolicyResponse(ctx, addResponse.TimeLimitLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultTimeLimitLogRotationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan timeLimitLogRotationPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogRotationPolicyApi.GetLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Time Limit Log Rotation Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state timeLimitLogRotationPolicyResourceModel
	readTimeLimitLogRotationPolicyResponse(ctx, readResponse.TimeLimitLogRotationPolicyResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogRotationPolicyApi.UpdateLogRotationPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createTimeLimitLogRotationPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogRotationPolicyApi.UpdateLogRotationPolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Time Limit Log Rotation Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readTimeLimitLogRotationPolicyResponse(ctx, updateResponse.TimeLimitLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
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
func (r *timeLimitLogRotationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readTimeLimitLogRotationPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultTimeLimitLogRotationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readTimeLimitLogRotationPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readTimeLimitLogRotationPolicy(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state timeLimitLogRotationPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogRotationPolicyApi.GetLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Time Limit Log Rotation Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readTimeLimitLogRotationPolicyResponse(ctx, readResponse.TimeLimitLogRotationPolicyResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *timeLimitLogRotationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateTimeLimitLogRotationPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultTimeLimitLogRotationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateTimeLimitLogRotationPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateTimeLimitLogRotationPolicy(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan timeLimitLogRotationPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state timeLimitLogRotationPolicyResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogRotationPolicyApi.UpdateLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createTimeLimitLogRotationPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogRotationPolicyApi.UpdateLogRotationPolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Time Limit Log Rotation Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readTimeLimitLogRotationPolicyResponse(ctx, updateResponse.TimeLimitLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultTimeLimitLogRotationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *timeLimitLogRotationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state timeLimitLogRotationPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogRotationPolicyApi.DeleteLogRotationPolicyExecute(r.apiClient.LogRotationPolicyApi.DeleteLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Time Limit Log Rotation Policy", err, httpResp)
		return
	}
}

func (r *timeLimitLogRotationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importTimeLimitLogRotationPolicy(ctx, req, resp)
}

func (r *defaultTimeLimitLogRotationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importTimeLimitLogRotationPolicy(ctx, req, resp)
}

func importTimeLimitLogRotationPolicy(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
