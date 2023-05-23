package resultcriteria

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	_ resource.Resource                = &replicationAssuranceResultCriteriaResource{}
	_ resource.ResourceWithConfigure   = &replicationAssuranceResultCriteriaResource{}
	_ resource.ResourceWithImportState = &replicationAssuranceResultCriteriaResource{}
	_ resource.Resource                = &defaultReplicationAssuranceResultCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultReplicationAssuranceResultCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultReplicationAssuranceResultCriteriaResource{}
)

// Create a Replication Assurance Result Criteria resource
func NewReplicationAssuranceResultCriteriaResource() resource.Resource {
	return &replicationAssuranceResultCriteriaResource{}
}

func NewDefaultReplicationAssuranceResultCriteriaResource() resource.Resource {
	return &defaultReplicationAssuranceResultCriteriaResource{}
}

// replicationAssuranceResultCriteriaResource is the resource implementation.
type replicationAssuranceResultCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultReplicationAssuranceResultCriteriaResource is the resource implementation.
type defaultReplicationAssuranceResultCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *replicationAssuranceResultCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication_assurance_result_criteria"
}

func (r *defaultReplicationAssuranceResultCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_replication_assurance_result_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *replicationAssuranceResultCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultReplicationAssuranceResultCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type replicationAssuranceResultCriteriaResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	LastUpdated                       types.String `tfsdk:"last_updated"`
	Notifications                     types.Set    `tfsdk:"notifications"`
	RequiredActions                   types.Set    `tfsdk:"required_actions"`
	LocalAssuranceLevel               types.Set    `tfsdk:"local_assurance_level"`
	RemoteAssuranceLevel              types.Set    `tfsdk:"remote_assurance_level"`
	AssuranceTimeoutCriteria          types.String `tfsdk:"assurance_timeout_criteria"`
	AssuranceTimeoutValue             types.String `tfsdk:"assurance_timeout_value"`
	ResponseDelayedByAssurance        types.String `tfsdk:"response_delayed_by_assurance"`
	AssuranceBehaviorAlteredByControl types.String `tfsdk:"assurance_behavior_altered_by_control"`
	AssuranceSatisfied                types.String `tfsdk:"assurance_satisfied"`
	Description                       types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *replicationAssuranceResultCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	replicationAssuranceResultCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultReplicationAssuranceResultCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	replicationAssuranceResultCriteriaSchema(ctx, req, resp, true)
}

func replicationAssuranceResultCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Replication Assurance Result Criteria.",
		Attributes: map[string]schema.Attribute{
			"local_assurance_level": schema.SetAttribute{
				Description: "The local assurance level values that will be allowed to match this Replication Assurance Result Criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"remote_assurance_level": schema.SetAttribute{
				Description: "The local assurance level values that will be allowed to match this Replication Assurance Result Criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"assurance_timeout_criteria": schema.StringAttribute{
				Description: "The criteria to use when performing matching based on the assurance timeout.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"assurance_timeout_value": schema.StringAttribute{
				Description: "The value to use for performing matching based on the assurance timeout. This will be ignored if the assurance-timeout-criteria is \"any\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"response_delayed_by_assurance": schema.StringAttribute{
				Description: "Indicates whether this Replication Assurance Result Criteria should match operations based on whether the response to the client was delayed by assurance processing.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"assurance_behavior_altered_by_control": schema.StringAttribute{
				Description: "Indicates whether this Replication Assurance Result Criteria should match operations based on whether the assurance requirements were altered by a control included in the request from the client.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"assurance_satisfied": schema.StringAttribute{
				Description: "Indicates whether this Replication Assurance Result Criteria should match operations based on whether the assurance requirements have been satisfied.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Result Criteria",
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
func addOptionalReplicationAssuranceResultCriteriaFields(ctx context.Context, addRequest *client.AddReplicationAssuranceResultCriteriaRequest, plan replicationAssuranceResultCriteriaResourceModel) error {
	if internaltypes.IsDefined(plan.LocalAssuranceLevel) {
		var slice []string
		plan.LocalAssuranceLevel.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumresultCriteriaLocalAssuranceLevelProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumresultCriteriaLocalAssuranceLevelPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.LocalAssuranceLevel = enumSlice
	}
	if internaltypes.IsDefined(plan.RemoteAssuranceLevel) {
		var slice []string
		plan.RemoteAssuranceLevel.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumresultCriteriaRemoteAssuranceLevelProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumresultCriteriaRemoteAssuranceLevelPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.RemoteAssuranceLevel = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AssuranceTimeoutCriteria) {
		assuranceTimeoutCriteria, err := client.NewEnumresultCriteriaAssuranceTimeoutCriteriaPropFromValue(plan.AssuranceTimeoutCriteria.ValueString())
		if err != nil {
			return err
		}
		addRequest.AssuranceTimeoutCriteria = assuranceTimeoutCriteria
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AssuranceTimeoutValue) {
		addRequest.AssuranceTimeoutValue = plan.AssuranceTimeoutValue.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ResponseDelayedByAssurance) {
		responseDelayedByAssurance, err := client.NewEnumresultCriteriaResponseDelayedByAssurancePropFromValue(plan.ResponseDelayedByAssurance.ValueString())
		if err != nil {
			return err
		}
		addRequest.ResponseDelayedByAssurance = responseDelayedByAssurance
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AssuranceBehaviorAlteredByControl) {
		assuranceBehaviorAlteredByControl, err := client.NewEnumresultCriteriaAssuranceBehaviorAlteredByControlPropFromValue(plan.AssuranceBehaviorAlteredByControl.ValueString())
		if err != nil {
			return err
		}
		addRequest.AssuranceBehaviorAlteredByControl = assuranceBehaviorAlteredByControl
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AssuranceSatisfied) {
		assuranceSatisfied, err := client.NewEnumresultCriteriaAssuranceSatisfiedPropFromValue(plan.AssuranceSatisfied.ValueString())
		if err != nil {
			return err
		}
		addRequest.AssuranceSatisfied = assuranceSatisfied
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Read a ReplicationAssuranceResultCriteriaResponse object into the model struct
func readReplicationAssuranceResultCriteriaResponse(ctx context.Context, r *client.ReplicationAssuranceResultCriteriaResponse, state *replicationAssuranceResultCriteriaResourceModel, expectedValues *replicationAssuranceResultCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.LocalAssuranceLevel = internaltypes.GetStringSet(
		client.StringSliceEnumresultCriteriaLocalAssuranceLevelProp(r.LocalAssuranceLevel))
	state.RemoteAssuranceLevel = internaltypes.GetStringSet(
		client.StringSliceEnumresultCriteriaRemoteAssuranceLevelProp(r.RemoteAssuranceLevel))
	state.AssuranceTimeoutCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaAssuranceTimeoutCriteriaProp(r.AssuranceTimeoutCriteria), internaltypes.IsEmptyString(expectedValues.AssuranceTimeoutCriteria))
	state.AssuranceTimeoutValue = internaltypes.StringTypeOrNil(r.AssuranceTimeoutValue, internaltypes.IsEmptyString(expectedValues.AssuranceTimeoutValue))
	config.CheckMismatchedPDFormattedAttributes("assurance_timeout_value",
		expectedValues.AssuranceTimeoutValue, state.AssuranceTimeoutValue, diagnostics)
	state.ResponseDelayedByAssurance = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaResponseDelayedByAssuranceProp(r.ResponseDelayedByAssurance), internaltypes.IsEmptyString(expectedValues.ResponseDelayedByAssurance))
	state.AssuranceBehaviorAlteredByControl = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaAssuranceBehaviorAlteredByControlProp(r.AssuranceBehaviorAlteredByControl), internaltypes.IsEmptyString(expectedValues.AssuranceBehaviorAlteredByControl))
	state.AssuranceSatisfied = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaAssuranceSatisfiedProp(r.AssuranceSatisfied), internaltypes.IsEmptyString(expectedValues.AssuranceSatisfied))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createReplicationAssuranceResultCriteriaOperations(plan replicationAssuranceResultCriteriaResourceModel, state replicationAssuranceResultCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.LocalAssuranceLevel, state.LocalAssuranceLevel, "local-assurance-level")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RemoteAssuranceLevel, state.RemoteAssuranceLevel, "remote-assurance-level")
	operations.AddStringOperationIfNecessary(&ops, plan.AssuranceTimeoutCriteria, state.AssuranceTimeoutCriteria, "assurance-timeout-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.AssuranceTimeoutValue, state.AssuranceTimeoutValue, "assurance-timeout-value")
	operations.AddStringOperationIfNecessary(&ops, plan.ResponseDelayedByAssurance, state.ResponseDelayedByAssurance, "response-delayed-by-assurance")
	operations.AddStringOperationIfNecessary(&ops, plan.AssuranceBehaviorAlteredByControl, state.AssuranceBehaviorAlteredByControl, "assurance-behavior-altered-by-control")
	operations.AddStringOperationIfNecessary(&ops, plan.AssuranceSatisfied, state.AssuranceSatisfied, "assurance-satisfied")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *replicationAssuranceResultCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan replicationAssuranceResultCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddReplicationAssuranceResultCriteriaRequest(plan.Id.ValueString(),
		[]client.EnumreplicationAssuranceResultCriteriaSchemaUrn{client.ENUMREPLICATIONASSURANCERESULTCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RESULT_CRITERIAREPLICATION_ASSURANCE})
	err := addOptionalReplicationAssuranceResultCriteriaFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Replication Assurance Result Criteria", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ResultCriteriaApi.AddResultCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddResultCriteriaRequest(
		client.AddReplicationAssuranceResultCriteriaRequestAsAddResultCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ResultCriteriaApi.AddResultCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Replication Assurance Result Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state replicationAssuranceResultCriteriaResourceModel
	readReplicationAssuranceResultCriteriaResponse(ctx, addResponse.ReplicationAssuranceResultCriteriaResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultReplicationAssuranceResultCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan replicationAssuranceResultCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ResultCriteriaApi.GetResultCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Replication Assurance Result Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state replicationAssuranceResultCriteriaResourceModel
	readReplicationAssuranceResultCriteriaResponse(ctx, readResponse.ReplicationAssuranceResultCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ResultCriteriaApi.UpdateResultCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createReplicationAssuranceResultCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ResultCriteriaApi.UpdateResultCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Replication Assurance Result Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readReplicationAssuranceResultCriteriaResponse(ctx, updateResponse.ReplicationAssuranceResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *replicationAssuranceResultCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readReplicationAssuranceResultCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultReplicationAssuranceResultCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readReplicationAssuranceResultCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readReplicationAssuranceResultCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state replicationAssuranceResultCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ResultCriteriaApi.GetResultCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Replication Assurance Result Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readReplicationAssuranceResultCriteriaResponse(ctx, readResponse.ReplicationAssuranceResultCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *replicationAssuranceResultCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateReplicationAssuranceResultCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultReplicationAssuranceResultCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateReplicationAssuranceResultCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateReplicationAssuranceResultCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan replicationAssuranceResultCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state replicationAssuranceResultCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ResultCriteriaApi.UpdateResultCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createReplicationAssuranceResultCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ResultCriteriaApi.UpdateResultCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Replication Assurance Result Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readReplicationAssuranceResultCriteriaResponse(ctx, updateResponse.ReplicationAssuranceResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultReplicationAssuranceResultCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *replicationAssuranceResultCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state replicationAssuranceResultCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ResultCriteriaApi.DeleteResultCriteriaExecute(r.apiClient.ResultCriteriaApi.DeleteResultCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Replication Assurance Result Criteria", err, httpResp)
		return
	}
}

func (r *replicationAssuranceResultCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importReplicationAssuranceResultCriteria(ctx, req, resp)
}

func (r *defaultReplicationAssuranceResultCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importReplicationAssuranceResultCriteria(ctx, req, resp)
}

func importReplicationAssuranceResultCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
