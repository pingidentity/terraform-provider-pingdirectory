package logretentionpolicy

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
	_ resource.Resource                = &timeLimitLogRetentionPolicyResource{}
	_ resource.ResourceWithConfigure   = &timeLimitLogRetentionPolicyResource{}
	_ resource.ResourceWithImportState = &timeLimitLogRetentionPolicyResource{}
	_ resource.Resource                = &defaultTimeLimitLogRetentionPolicyResource{}
	_ resource.ResourceWithConfigure   = &defaultTimeLimitLogRetentionPolicyResource{}
	_ resource.ResourceWithImportState = &defaultTimeLimitLogRetentionPolicyResource{}
)

// Create a Time Limit Log Retention Policy resource
func NewTimeLimitLogRetentionPolicyResource() resource.Resource {
	return &timeLimitLogRetentionPolicyResource{}
}

func NewDefaultTimeLimitLogRetentionPolicyResource() resource.Resource {
	return &defaultTimeLimitLogRetentionPolicyResource{}
}

// timeLimitLogRetentionPolicyResource is the resource implementation.
type timeLimitLogRetentionPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultTimeLimitLogRetentionPolicyResource is the resource implementation.
type defaultTimeLimitLogRetentionPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *timeLimitLogRetentionPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_time_limit_log_retention_policy"
}

func (r *defaultTimeLimitLogRetentionPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_time_limit_log_retention_policy"
}

// Configure adds the provider configured client to the resource.
func (r *timeLimitLogRetentionPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultTimeLimitLogRetentionPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type timeLimitLogRetentionPolicyResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	RetainDuration  types.String `tfsdk:"retain_duration"`
	Description     types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *timeLimitLogRetentionPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	timeLimitLogRetentionPolicySchema(ctx, req, resp, false)
}

func (r *defaultTimeLimitLogRetentionPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	timeLimitLogRetentionPolicySchema(ctx, req, resp, true)
}

func timeLimitLogRetentionPolicySchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Time Limit Log Retention Policy.",
		Attributes: map[string]schema.Attribute{
			"retain_duration": schema.StringAttribute{
				Description: "Specifies the desired minimum length of time that each log file should be retained.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Retention Policy",
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
func addOptionalTimeLimitLogRetentionPolicyFields(ctx context.Context, addRequest *client.AddTimeLimitLogRetentionPolicyRequest, plan timeLimitLogRetentionPolicyResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a TimeLimitLogRetentionPolicyResponse object into the model struct
func readTimeLimitLogRetentionPolicyResponse(ctx context.Context, r *client.TimeLimitLogRetentionPolicyResponse, state *timeLimitLogRetentionPolicyResourceModel, expectedValues *timeLimitLogRetentionPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.RetainDuration = types.StringValue(r.RetainDuration)
	config.CheckMismatchedPDFormattedAttributes("retain_duration",
		expectedValues.RetainDuration, state.RetainDuration, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createTimeLimitLogRetentionPolicyOperations(plan timeLimitLogRetentionPolicyResourceModel, state timeLimitLogRetentionPolicyResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.RetainDuration, state.RetainDuration, "retain-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *timeLimitLogRetentionPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan timeLimitLogRetentionPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddTimeLimitLogRetentionPolicyRequest(plan.Id.ValueString(),
		[]client.EnumtimeLimitLogRetentionPolicySchemaUrn{client.ENUMTIMELIMITLOGRETENTIONPOLICYSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_RETENTION_POLICYTIME_LIMIT},
		plan.RetainDuration.ValueString())
	addOptionalTimeLimitLogRetentionPolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogRetentionPolicyApi.AddLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogRetentionPolicyRequest(
		client.AddTimeLimitLogRetentionPolicyRequestAsAddLogRetentionPolicyRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogRetentionPolicyApi.AddLogRetentionPolicyExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Time Limit Log Retention Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state timeLimitLogRetentionPolicyResourceModel
	readTimeLimitLogRetentionPolicyResponse(ctx, addResponse.TimeLimitLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultTimeLimitLogRetentionPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan timeLimitLogRetentionPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogRetentionPolicyApi.GetLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Time Limit Log Retention Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state timeLimitLogRetentionPolicyResourceModel
	readTimeLimitLogRetentionPolicyResponse(ctx, readResponse.TimeLimitLogRetentionPolicyResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogRetentionPolicyApi.UpdateLogRetentionPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createTimeLimitLogRetentionPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogRetentionPolicyApi.UpdateLogRetentionPolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Time Limit Log Retention Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readTimeLimitLogRetentionPolicyResponse(ctx, updateResponse.TimeLimitLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
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
func (r *timeLimitLogRetentionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readTimeLimitLogRetentionPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultTimeLimitLogRetentionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readTimeLimitLogRetentionPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readTimeLimitLogRetentionPolicy(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state timeLimitLogRetentionPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogRetentionPolicyApi.GetLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Time Limit Log Retention Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readTimeLimitLogRetentionPolicyResponse(ctx, readResponse.TimeLimitLogRetentionPolicyResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *timeLimitLogRetentionPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateTimeLimitLogRetentionPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultTimeLimitLogRetentionPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateTimeLimitLogRetentionPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateTimeLimitLogRetentionPolicy(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan timeLimitLogRetentionPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state timeLimitLogRetentionPolicyResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogRetentionPolicyApi.UpdateLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createTimeLimitLogRetentionPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogRetentionPolicyApi.UpdateLogRetentionPolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Time Limit Log Retention Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readTimeLimitLogRetentionPolicyResponse(ctx, updateResponse.TimeLimitLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultTimeLimitLogRetentionPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *timeLimitLogRetentionPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state timeLimitLogRetentionPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogRetentionPolicyApi.DeleteLogRetentionPolicyExecute(r.apiClient.LogRetentionPolicyApi.DeleteLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Time Limit Log Retention Policy", err, httpResp)
		return
	}
}

func (r *timeLimitLogRetentionPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importTimeLimitLogRetentionPolicy(ctx, req, resp)
}

func (r *defaultTimeLimitLogRetentionPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importTimeLimitLogRetentionPolicy(ctx, req, resp)
}

func importTimeLimitLogRetentionPolicy(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
