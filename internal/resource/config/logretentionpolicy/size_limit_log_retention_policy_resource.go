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
	_ resource.Resource                = &sizeLimitLogRetentionPolicyResource{}
	_ resource.ResourceWithConfigure   = &sizeLimitLogRetentionPolicyResource{}
	_ resource.ResourceWithImportState = &sizeLimitLogRetentionPolicyResource{}
	_ resource.Resource                = &defaultSizeLimitLogRetentionPolicyResource{}
	_ resource.ResourceWithConfigure   = &defaultSizeLimitLogRetentionPolicyResource{}
	_ resource.ResourceWithImportState = &defaultSizeLimitLogRetentionPolicyResource{}
)

// Create a Size Limit Log Retention Policy resource
func NewSizeLimitLogRetentionPolicyResource() resource.Resource {
	return &sizeLimitLogRetentionPolicyResource{}
}

func NewDefaultSizeLimitLogRetentionPolicyResource() resource.Resource {
	return &defaultSizeLimitLogRetentionPolicyResource{}
}

// sizeLimitLogRetentionPolicyResource is the resource implementation.
type sizeLimitLogRetentionPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSizeLimitLogRetentionPolicyResource is the resource implementation.
type defaultSizeLimitLogRetentionPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *sizeLimitLogRetentionPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_size_limit_log_retention_policy"
}

func (r *defaultSizeLimitLogRetentionPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_size_limit_log_retention_policy"
}

// Configure adds the provider configured client to the resource.
func (r *sizeLimitLogRetentionPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultSizeLimitLogRetentionPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type sizeLimitLogRetentionPolicyResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	DiskSpaceUsed   types.String `tfsdk:"disk_space_used"`
	Description     types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *sizeLimitLogRetentionPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	sizeLimitLogRetentionPolicySchema(ctx, req, resp, false)
}

func (r *defaultSizeLimitLogRetentionPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	sizeLimitLogRetentionPolicySchema(ctx, req, resp, true)
}

func sizeLimitLogRetentionPolicySchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Size Limit Log Retention Policy.",
		Attributes: map[string]schema.Attribute{
			"disk_space_used": schema.StringAttribute{
				Description: "Specifies the maximum total disk space used by the log files.",
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
func addOptionalSizeLimitLogRetentionPolicyFields(ctx context.Context, addRequest *client.AddSizeLimitLogRetentionPolicyRequest, plan sizeLimitLogRetentionPolicyResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a SizeLimitLogRetentionPolicyResponse object into the model struct
func readSizeLimitLogRetentionPolicyResponse(ctx context.Context, r *client.SizeLimitLogRetentionPolicyResponse, state *sizeLimitLogRetentionPolicyResourceModel, expectedValues *sizeLimitLogRetentionPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.DiskSpaceUsed = types.StringValue(r.DiskSpaceUsed)
	config.CheckMismatchedPDFormattedAttributes("disk_space_used",
		expectedValues.DiskSpaceUsed, state.DiskSpaceUsed, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSizeLimitLogRetentionPolicyOperations(plan sizeLimitLogRetentionPolicyResourceModel, state sizeLimitLogRetentionPolicyResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.DiskSpaceUsed, state.DiskSpaceUsed, "disk-space-used")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *sizeLimitLogRetentionPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan sizeLimitLogRetentionPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddSizeLimitLogRetentionPolicyRequest(plan.Id.ValueString(),
		[]client.EnumsizeLimitLogRetentionPolicySchemaUrn{client.ENUMSIZELIMITLOGRETENTIONPOLICYSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_RETENTION_POLICYSIZE_LIMIT},
		plan.DiskSpaceUsed.ValueString())
	addOptionalSizeLimitLogRetentionPolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogRetentionPolicyApi.AddLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogRetentionPolicyRequest(
		client.AddSizeLimitLogRetentionPolicyRequestAsAddLogRetentionPolicyRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogRetentionPolicyApi.AddLogRetentionPolicyExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Size Limit Log Retention Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state sizeLimitLogRetentionPolicyResourceModel
	readSizeLimitLogRetentionPolicyResponse(ctx, addResponse.SizeLimitLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultSizeLimitLogRetentionPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan sizeLimitLogRetentionPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogRetentionPolicyApi.GetLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Size Limit Log Retention Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state sizeLimitLogRetentionPolicyResourceModel
	readSizeLimitLogRetentionPolicyResponse(ctx, readResponse.SizeLimitLogRetentionPolicyResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogRetentionPolicyApi.UpdateLogRetentionPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSizeLimitLogRetentionPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogRetentionPolicyApi.UpdateLogRetentionPolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Size Limit Log Retention Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSizeLimitLogRetentionPolicyResponse(ctx, updateResponse.SizeLimitLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
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
func (r *sizeLimitLogRetentionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSizeLimitLogRetentionPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSizeLimitLogRetentionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSizeLimitLogRetentionPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSizeLimitLogRetentionPolicy(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state sizeLimitLogRetentionPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogRetentionPolicyApi.GetLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Size Limit Log Retention Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSizeLimitLogRetentionPolicyResponse(ctx, readResponse.SizeLimitLogRetentionPolicyResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *sizeLimitLogRetentionPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSizeLimitLogRetentionPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSizeLimitLogRetentionPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSizeLimitLogRetentionPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSizeLimitLogRetentionPolicy(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan sizeLimitLogRetentionPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state sizeLimitLogRetentionPolicyResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogRetentionPolicyApi.UpdateLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSizeLimitLogRetentionPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogRetentionPolicyApi.UpdateLogRetentionPolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Size Limit Log Retention Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSizeLimitLogRetentionPolicyResponse(ctx, updateResponse.SizeLimitLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSizeLimitLogRetentionPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *sizeLimitLogRetentionPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state sizeLimitLogRetentionPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogRetentionPolicyApi.DeleteLogRetentionPolicyExecute(r.apiClient.LogRetentionPolicyApi.DeleteLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Size Limit Log Retention Policy", err, httpResp)
		return
	}
}

func (r *sizeLimitLogRetentionPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSizeLimitLogRetentionPolicy(ctx, req, resp)
}

func (r *defaultSizeLimitLogRetentionPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSizeLimitLogRetentionPolicy(ctx, req, resp)
}

func importSizeLimitLogRetentionPolicy(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
