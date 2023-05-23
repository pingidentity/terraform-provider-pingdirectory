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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &aggregateResultCriteriaResource{}
	_ resource.ResourceWithConfigure   = &aggregateResultCriteriaResource{}
	_ resource.ResourceWithImportState = &aggregateResultCriteriaResource{}
	_ resource.Resource                = &defaultAggregateResultCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultAggregateResultCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultAggregateResultCriteriaResource{}
)

// Create a Aggregate Result Criteria resource
func NewAggregateResultCriteriaResource() resource.Resource {
	return &aggregateResultCriteriaResource{}
}

func NewDefaultAggregateResultCriteriaResource() resource.Resource {
	return &defaultAggregateResultCriteriaResource{}
}

// aggregateResultCriteriaResource is the resource implementation.
type aggregateResultCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultAggregateResultCriteriaResource is the resource implementation.
type defaultAggregateResultCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *aggregateResultCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aggregate_result_criteria"
}

func (r *defaultAggregateResultCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_aggregate_result_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *aggregateResultCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultAggregateResultCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type aggregateResultCriteriaResourceModel struct {
	Id                           types.String `tfsdk:"id"`
	LastUpdated                  types.String `tfsdk:"last_updated"`
	Notifications                types.Set    `tfsdk:"notifications"`
	RequiredActions              types.Set    `tfsdk:"required_actions"`
	AllIncludedResultCriteria    types.Set    `tfsdk:"all_included_result_criteria"`
	AnyIncludedResultCriteria    types.Set    `tfsdk:"any_included_result_criteria"`
	NotAllIncludedResultCriteria types.Set    `tfsdk:"not_all_included_result_criteria"`
	NoneIncludedResultCriteria   types.Set    `tfsdk:"none_included_result_criteria"`
	Description                  types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *aggregateResultCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	aggregateResultCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultAggregateResultCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	aggregateResultCriteriaSchema(ctx, req, resp, true)
}

func aggregateResultCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Aggregate Result Criteria.",
		Attributes: map[string]schema.Attribute{
			"all_included_result_criteria": schema.SetAttribute{
				Description: "Specifies a result criteria object that must match the associated operation result in order to match the aggregate result criteria. If one or more all-included result criteria objects are provided, then an operation result must match all of them in order to match the aggregate result criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"any_included_result_criteria": schema.SetAttribute{
				Description: "Specifies a result criteria object that may match the associated operation result in order to match the aggregate result criteria. If one or more any-included result criteria objects are provided, then an operation result must match at least one of them in order to match the aggregate result criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"not_all_included_result_criteria": schema.SetAttribute{
				Description: "Specifies a result criteria object that should not match the associated operation result in order to match the aggregate result criteria. If one or more not-all-included result criteria objects are provided, then an operation result must not match all of them (that is, it may match zero or more of them, but it must not match all of them) in order to match the aggregate result criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"none_included_result_criteria": schema.SetAttribute{
				Description: "Specifies a result criteria object that must not match the associated operation result in order to match the aggregate result criteria. If one or more none-included result criteria objects are provided, then an operation result must not match any of them in order to match the aggregate result criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
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
func addOptionalAggregateResultCriteriaFields(ctx context.Context, addRequest *client.AddAggregateResultCriteriaRequest, plan aggregateResultCriteriaResourceModel) {
	if internaltypes.IsDefined(plan.AllIncludedResultCriteria) {
		var slice []string
		plan.AllIncludedResultCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedResultCriteria = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedResultCriteria) {
		var slice []string
		plan.AnyIncludedResultCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedResultCriteria = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedResultCriteria) {
		var slice []string
		plan.NotAllIncludedResultCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedResultCriteria = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedResultCriteria) {
		var slice []string
		plan.NoneIncludedResultCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedResultCriteria = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a AggregateResultCriteriaResponse object into the model struct
func readAggregateResultCriteriaResponse(ctx context.Context, r *client.AggregateResultCriteriaResponse, state *aggregateResultCriteriaResourceModel, expectedValues *aggregateResultCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.AllIncludedResultCriteria = internaltypes.GetStringSet(r.AllIncludedResultCriteria)
	state.AnyIncludedResultCriteria = internaltypes.GetStringSet(r.AnyIncludedResultCriteria)
	state.NotAllIncludedResultCriteria = internaltypes.GetStringSet(r.NotAllIncludedResultCriteria)
	state.NoneIncludedResultCriteria = internaltypes.GetStringSet(r.NoneIncludedResultCriteria)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createAggregateResultCriteriaOperations(plan aggregateResultCriteriaResourceModel, state aggregateResultCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedResultCriteria, state.AllIncludedResultCriteria, "all-included-result-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedResultCriteria, state.AnyIncludedResultCriteria, "any-included-result-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedResultCriteria, state.NotAllIncludedResultCriteria, "not-all-included-result-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedResultCriteria, state.NoneIncludedResultCriteria, "none-included-result-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *aggregateResultCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan aggregateResultCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddAggregateResultCriteriaRequest(plan.Id.ValueString(),
		[]client.EnumaggregateResultCriteriaSchemaUrn{client.ENUMAGGREGATERESULTCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RESULT_CRITERIAAGGREGATE})
	addOptionalAggregateResultCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ResultCriteriaApi.AddResultCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddResultCriteriaRequest(
		client.AddAggregateResultCriteriaRequestAsAddResultCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ResultCriteriaApi.AddResultCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Aggregate Result Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state aggregateResultCriteriaResourceModel
	readAggregateResultCriteriaResponse(ctx, addResponse.AggregateResultCriteriaResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultAggregateResultCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan aggregateResultCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ResultCriteriaApi.GetResultCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Aggregate Result Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state aggregateResultCriteriaResourceModel
	readAggregateResultCriteriaResponse(ctx, readResponse.AggregateResultCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ResultCriteriaApi.UpdateResultCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createAggregateResultCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ResultCriteriaApi.UpdateResultCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Aggregate Result Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAggregateResultCriteriaResponse(ctx, updateResponse.AggregateResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *aggregateResultCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAggregateResultCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAggregateResultCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAggregateResultCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readAggregateResultCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state aggregateResultCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ResultCriteriaApi.GetResultCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Aggregate Result Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAggregateResultCriteriaResponse(ctx, readResponse.AggregateResultCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *aggregateResultCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAggregateResultCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAggregateResultCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAggregateResultCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateAggregateResultCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan aggregateResultCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state aggregateResultCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ResultCriteriaApi.UpdateResultCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createAggregateResultCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ResultCriteriaApi.UpdateResultCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Aggregate Result Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAggregateResultCriteriaResponse(ctx, updateResponse.AggregateResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultAggregateResultCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *aggregateResultCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state aggregateResultCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ResultCriteriaApi.DeleteResultCriteriaExecute(r.apiClient.ResultCriteriaApi.DeleteResultCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Aggregate Result Criteria", err, httpResp)
		return
	}
}

func (r *aggregateResultCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAggregateResultCriteria(ctx, req, resp)
}

func (r *defaultAggregateResultCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAggregateResultCriteria(ctx, req, resp)
}

func importAggregateResultCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
