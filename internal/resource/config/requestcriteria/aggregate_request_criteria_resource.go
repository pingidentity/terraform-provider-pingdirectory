package requestcriteria

import (
	"context"
	"time"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &aggregateRequestCriteriaResource{}
	_ resource.ResourceWithConfigure   = &aggregateRequestCriteriaResource{}
	_ resource.ResourceWithImportState = &aggregateRequestCriteriaResource{}
)

// Create a Aggregate Request Criteria resource
func NewAggregateRequestCriteriaResource() resource.Resource {
	return &aggregateRequestCriteriaResource{}
}

// aggregateRequestCriteriaResource is the resource implementation.
type aggregateRequestCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *aggregateRequestCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aggregate_request_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *aggregateRequestCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type aggregateRequestCriteriaResourceModel struct {
	Id                            types.String `tfsdk:"id"`
	LastUpdated                   types.String `tfsdk:"last_updated"`
	Notifications                 types.Set    `tfsdk:"notifications"`
	RequiredActions               types.Set    `tfsdk:"required_actions"`
	AllIncludedRequestCriteria    types.Set    `tfsdk:"all_included_request_criteria"`
	AnyIncludedRequestCriteria    types.Set    `tfsdk:"any_included_request_criteria"`
	NotAllIncludedRequestCriteria types.Set    `tfsdk:"not_all_included_request_criteria"`
	NoneIncludedRequestCriteria   types.Set    `tfsdk:"none_included_request_criteria"`
	Description                   types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *aggregateRequestCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Aggregate Request Criteria.",
		Attributes: map[string]schema.Attribute{
			"all_included_request_criteria": schema.SetAttribute{
				Description: "Specifies a request criteria object that must match the associated operation request in order to match the aggregate request criteria. If one or more all-included request criteria objects are provided, then an operation request must match all of them in order to match the aggregate request criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_request_criteria": schema.SetAttribute{
				Description: "Specifies a request criteria object that may match the associated operation request in order to the this aggregate request criteria. If one or more any-included request criteria objects are provided, then an operation request must match at least one of them in order to match the aggregate request criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_request_criteria": schema.SetAttribute{
				Description: "Specifies a request criteria object that should not match the associated operation request in order to match the aggregate request criteria. If one or more not-all-included request criteria objects are provided, then an operation request must not match all of them (that is, it may match zero or more of them, but it must not match all of them) in order to match the aggregate request criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_request_criteria": schema.SetAttribute{
				Description: "Specifies a request criteria object that must not match the associated operation request in order to match the aggregate request criteria. If one or more none-included request criteria objects are provided, then an operation request must not match any of them in order to match the aggregate request criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Request Criteria",
				Optional:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalAggregateRequestCriteriaFields(ctx context.Context, addRequest *client.AddAggregateRequestCriteriaRequest, plan aggregateRequestCriteriaResourceModel) {
	if internaltypes.IsDefined(plan.AllIncludedRequestCriteria) {
		var slice []string
		plan.AllIncludedRequestCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedRequestCriteria = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedRequestCriteria) {
		var slice []string
		plan.AnyIncludedRequestCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedRequestCriteria = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedRequestCriteria) {
		var slice []string
		plan.NotAllIncludedRequestCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedRequestCriteria = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedRequestCriteria) {
		var slice []string
		plan.NoneIncludedRequestCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedRequestCriteria = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
}

// Read a AggregateRequestCriteriaResponse object into the model struct
func readAggregateRequestCriteriaResponse(ctx context.Context, r *client.AggregateRequestCriteriaResponse, state *aggregateRequestCriteriaResourceModel, expectedValues *aggregateRequestCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.AllIncludedRequestCriteria = internaltypes.GetStringSet(r.AllIncludedRequestCriteria)
	state.AnyIncludedRequestCriteria = internaltypes.GetStringSet(r.AnyIncludedRequestCriteria)
	state.NotAllIncludedRequestCriteria = internaltypes.GetStringSet(r.NotAllIncludedRequestCriteria)
	state.NoneIncludedRequestCriteria = internaltypes.GetStringSet(r.NoneIncludedRequestCriteria)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createAggregateRequestCriteriaOperations(plan aggregateRequestCriteriaResourceModel, state aggregateRequestCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedRequestCriteria, state.AllIncludedRequestCriteria, "all-included-request-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedRequestCriteria, state.AnyIncludedRequestCriteria, "any-included-request-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedRequestCriteria, state.NotAllIncludedRequestCriteria, "not-all-included-request-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedRequestCriteria, state.NoneIncludedRequestCriteria, "none-included-request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *aggregateRequestCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan aggregateRequestCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddAggregateRequestCriteriaRequest(plan.Id.ValueString(),
		[]client.EnumaggregateRequestCriteriaSchemaUrn{client.ENUMAGGREGATEREQUESTCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0REQUEST_CRITERIAAGGREGATE})
	addOptionalAggregateRequestCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RequestCriteriaApi.AddRequestCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRequestCriteriaRequest(
		client.AddAggregateRequestCriteriaRequestAsAddRequestCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RequestCriteriaApi.AddRequestCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Aggregate Request Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state aggregateRequestCriteriaResourceModel
	readAggregateRequestCriteriaResponse(ctx, addResponse.AggregateRequestCriteriaResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *aggregateRequestCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state aggregateRequestCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RequestCriteriaApi.GetRequestCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Aggregate Request Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAggregateRequestCriteriaResponse(ctx, readResponse.AggregateRequestCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *aggregateRequestCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan aggregateRequestCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state aggregateRequestCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.RequestCriteriaApi.UpdateRequestCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createAggregateRequestCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.RequestCriteriaApi.UpdateRequestCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Aggregate Request Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAggregateRequestCriteriaResponse(ctx, updateResponse.AggregateRequestCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *aggregateRequestCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state aggregateRequestCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.RequestCriteriaApi.DeleteRequestCriteriaExecute(r.apiClient.RequestCriteriaApi.DeleteRequestCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Aggregate Request Criteria", err, httpResp)
		return
	}
}

func (r *aggregateRequestCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
