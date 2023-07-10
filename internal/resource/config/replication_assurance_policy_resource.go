package config

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &replicationAssurancePolicyResource{}
	_ resource.ResourceWithConfigure   = &replicationAssurancePolicyResource{}
	_ resource.ResourceWithImportState = &replicationAssurancePolicyResource{}
	_ resource.Resource                = &defaultReplicationAssurancePolicyResource{}
	_ resource.ResourceWithConfigure   = &defaultReplicationAssurancePolicyResource{}
	_ resource.ResourceWithImportState = &defaultReplicationAssurancePolicyResource{}
)

// Create a Replication Assurance Policy resource
func NewReplicationAssurancePolicyResource() resource.Resource {
	return &replicationAssurancePolicyResource{}
}

func NewDefaultReplicationAssurancePolicyResource() resource.Resource {
	return &defaultReplicationAssurancePolicyResource{}
}

// replicationAssurancePolicyResource is the resource implementation.
type replicationAssurancePolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultReplicationAssurancePolicyResource is the resource implementation.
type defaultReplicationAssurancePolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *replicationAssurancePolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication_assurance_policy"
}

func (r *defaultReplicationAssurancePolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_replication_assurance_policy"
}

// Configure adds the provider configured client to the resource.
func (r *replicationAssurancePolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultReplicationAssurancePolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type replicationAssurancePolicyResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	LastUpdated          types.String `tfsdk:"last_updated"`
	Notifications        types.Set    `tfsdk:"notifications"`
	RequiredActions      types.Set    `tfsdk:"required_actions"`
	Description          types.String `tfsdk:"description"`
	Enabled              types.Bool   `tfsdk:"enabled"`
	EvaluationOrderIndex types.Int64  `tfsdk:"evaluation_order_index"`
	LocalLevel           types.String `tfsdk:"local_level"`
	RemoteLevel          types.String `tfsdk:"remote_level"`
	Timeout              types.String `tfsdk:"timeout"`
	ConnectionCriteria   types.String `tfsdk:"connection_criteria"`
	RequestCriteria      types.String `tfsdk:"request_criteria"`
}

// GetSchema defines the schema for the resource.
func (r *replicationAssurancePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	replicationAssurancePolicySchema(ctx, req, resp, false)
}

func (r *defaultReplicationAssurancePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	replicationAssurancePolicySchema(ctx, req, resp, true)
}

func replicationAssurancePolicySchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Replication Assurance Policy.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description: "Description of the Replication Assurance Policy.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Replication Assurance Policy is enabled for use in the server. If a Replication Assurance Policy is disabled, then no new operations will be associated with it.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"evaluation_order_index": schema.Int64Attribute{
				Description: "When multiple Replication Assurance Policies are defined, this property determines the evaluation order for finding a Replication Assurance Policy match against an operation. Policies are evaluated based on this index from least to greatest. Values of this property must be unique but not necessarily contiguous.",
				Required:    true,
			},
			"local_level": schema.StringAttribute{
				Description: "Specifies the assurance level used to replicate to local servers. A local server is defined as one with the same value for the location setting in the global configuration.  The local-level must be set to an assurance level at least as strict as the remote-level. In other words, if remote-level is set to \"received-any-remote-location\" or \"received-all-remote-locations\", then local-level must be either \"received-any-server\" or \"processed-all-servers\". If remote-level is \"processed-all-remote-servers\", then local-level must be \"processed-all-servers\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"remote_level": schema.StringAttribute{
				Description: "Specifies the assurance level used to replicate to remote servers. A remote server is defined as one with a different value for the location setting in the global configuration.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"timeout": schema.StringAttribute{
				Description: "Specifies the maximum length of time to wait for the replication assurance requirements to be met before timing out and replying to the client.",
				Required:    true,
			},
			"connection_criteria": schema.StringAttribute{
				Description: "Specifies a connection criteria used to indicate which operations from clients matching this criteria use this policy. If both a connection criteria and a request criteria are specified for a policy, then both must match an operation for the policy to be assigned.",
				Optional:    true,
			},
			"request_criteria": schema.StringAttribute{
				Description: "Specifies a request criteria used to indicate which operations from clients matching this criteria use this policy. If both a connection criteria and a request criteria are specified for a policy, then both must match an operation for the policy to be assigned.",
				Optional:    true,
			},
		},
	}
	if isDefault {
		// Add any default properties and set optional properties to computed where necessary
		SetAllAttributesToOptionalAndComputed(&schemaDef, []string{"id"})
	}
	AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add optional fields to create request for replication-assurance-policy replication-assurance-policy
func addOptionalReplicationAssurancePolicyFields(ctx context.Context, addRequest *client.AddReplicationAssurancePolicyRequest, plan replicationAssurancePolicyResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LocalLevel) {
		localLevel, err := client.NewEnumreplicationAssurancePolicyLocalLevelPropFromValue(plan.LocalLevel.ValueString())
		if err != nil {
			return err
		}
		addRequest.LocalLevel = localLevel
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RemoteLevel) {
		remoteLevel, err := client.NewEnumreplicationAssurancePolicyRemoteLevelPropFromValue(plan.RemoteLevel.ValueString())
		if err != nil {
			return err
		}
		addRequest.RemoteLevel = remoteLevel
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		addRequest.ConnectionCriteria = plan.ConnectionCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	return nil
}

// Read a ReplicationAssurancePolicyResponse object into the model struct
func readReplicationAssurancePolicyResponse(ctx context.Context, r *client.ReplicationAssurancePolicyResponse, state *replicationAssurancePolicyResourceModel, expectedValues *replicationAssurancePolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.EvaluationOrderIndex = types.Int64Value(r.EvaluationOrderIndex)
	state.LocalLevel = types.StringValue(r.LocalLevel.String())
	state.RemoteLevel = types.StringValue(r.RemoteLevel.String())
	state.Timeout = types.StringValue(r.Timeout)
	CheckMismatchedPDFormattedAttributes("timeout",
		expectedValues.Timeout, state.Timeout, diagnostics)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.Notifications, state.RequiredActions = ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createReplicationAssurancePolicyOperations(plan replicationAssurancePolicyResourceModel, state replicationAssurancePolicyResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddInt64OperationIfNecessary(&ops, plan.EvaluationOrderIndex, state.EvaluationOrderIndex, "evaluation-order-index")
	operations.AddStringOperationIfNecessary(&ops, plan.LocalLevel, state.LocalLevel, "local-level")
	operations.AddStringOperationIfNecessary(&ops, plan.RemoteLevel, state.RemoteLevel, "remote-level")
	operations.AddStringOperationIfNecessary(&ops, plan.Timeout, state.Timeout, "timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionCriteria, state.ConnectionCriteria, "connection-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.RequestCriteria, state.RequestCriteria, "request-criteria")
	return ops
}

// Create a replication-assurance-policy replication-assurance-policy
func (r *replicationAssurancePolicyResource) CreateReplicationAssurancePolicy(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan replicationAssurancePolicyResourceModel) (*replicationAssurancePolicyResourceModel, error) {
	addRequest := client.NewAddReplicationAssurancePolicyRequest(plan.Id.ValueString(),
		plan.EvaluationOrderIndex.ValueInt64(),
		plan.Timeout.ValueString())
	err := addOptionalReplicationAssurancePolicyFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Replication Assurance Policy", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ReplicationAssurancePolicyApi.AddReplicationAssurancePolicy(
		ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddReplicationAssurancePolicyRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.ReplicationAssurancePolicyApi.AddReplicationAssurancePolicyExecute(apiAddRequest)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Replication Assurance Policy", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state replicationAssurancePolicyResourceModel
	readReplicationAssurancePolicyResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *replicationAssurancePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan replicationAssurancePolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateReplicationAssurancePolicy(ctx, req, resp, plan)
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
func (r *defaultReplicationAssurancePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan replicationAssurancePolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ReplicationAssurancePolicyApi.GetReplicationAssurancePolicy(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Replication Assurance Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state replicationAssurancePolicyResourceModel
	readReplicationAssurancePolicyResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ReplicationAssurancePolicyApi.UpdateReplicationAssurancePolicy(ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createReplicationAssurancePolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ReplicationAssurancePolicyApi.UpdateReplicationAssurancePolicyExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Replication Assurance Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readReplicationAssurancePolicyResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *replicationAssurancePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readReplicationAssurancePolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultReplicationAssurancePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readReplicationAssurancePolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readReplicationAssurancePolicy(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state replicationAssurancePolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ReplicationAssurancePolicyApi.GetReplicationAssurancePolicy(
		ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Replication Assurance Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readReplicationAssurancePolicyResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *replicationAssurancePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateReplicationAssurancePolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultReplicationAssurancePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateReplicationAssurancePolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateReplicationAssurancePolicy(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan replicationAssurancePolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state replicationAssurancePolicyResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ReplicationAssurancePolicyApi.UpdateReplicationAssurancePolicy(
		ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createReplicationAssurancePolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ReplicationAssurancePolicyApi.UpdateReplicationAssurancePolicyExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Replication Assurance Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readReplicationAssurancePolicyResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultReplicationAssurancePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *replicationAssurancePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state replicationAssurancePolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ReplicationAssurancePolicyApi.DeleteReplicationAssurancePolicyExecute(r.apiClient.ReplicationAssurancePolicyApi.DeleteReplicationAssurancePolicy(
		ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Replication Assurance Policy", err, httpResp)
		return
	}
}

func (r *replicationAssurancePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importReplicationAssurancePolicy(ctx, req, resp)
}

func (r *defaultReplicationAssurancePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importReplicationAssurancePolicy(ctx, req, resp)
}

func importReplicationAssurancePolicy(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
