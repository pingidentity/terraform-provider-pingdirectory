package searchreferencecriteria

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
	_ resource.Resource                = &aggregateSearchReferenceCriteriaResource{}
	_ resource.ResourceWithConfigure   = &aggregateSearchReferenceCriteriaResource{}
	_ resource.ResourceWithImportState = &aggregateSearchReferenceCriteriaResource{}
	_ resource.Resource                = &defaultAggregateSearchReferenceCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultAggregateSearchReferenceCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultAggregateSearchReferenceCriteriaResource{}
)

// Create a Aggregate Search Reference Criteria resource
func NewAggregateSearchReferenceCriteriaResource() resource.Resource {
	return &aggregateSearchReferenceCriteriaResource{}
}

func NewDefaultAggregateSearchReferenceCriteriaResource() resource.Resource {
	return &defaultAggregateSearchReferenceCriteriaResource{}
}

// aggregateSearchReferenceCriteriaResource is the resource implementation.
type aggregateSearchReferenceCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultAggregateSearchReferenceCriteriaResource is the resource implementation.
type defaultAggregateSearchReferenceCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *aggregateSearchReferenceCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aggregate_search_reference_criteria"
}

func (r *defaultAggregateSearchReferenceCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_aggregate_search_reference_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *aggregateSearchReferenceCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultAggregateSearchReferenceCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type aggregateSearchReferenceCriteriaResourceModel struct {
	Id                                    types.String `tfsdk:"id"`
	LastUpdated                           types.String `tfsdk:"last_updated"`
	Notifications                         types.Set    `tfsdk:"notifications"`
	RequiredActions                       types.Set    `tfsdk:"required_actions"`
	AllIncludedSearchReferenceCriteria    types.Set    `tfsdk:"all_included_search_reference_criteria"`
	AnyIncludedSearchReferenceCriteria    types.Set    `tfsdk:"any_included_search_reference_criteria"`
	NotAllIncludedSearchReferenceCriteria types.Set    `tfsdk:"not_all_included_search_reference_criteria"`
	NoneIncludedSearchReferenceCriteria   types.Set    `tfsdk:"none_included_search_reference_criteria"`
	Description                           types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *aggregateSearchReferenceCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	aggregateSearchReferenceCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultAggregateSearchReferenceCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	aggregateSearchReferenceCriteriaSchema(ctx, req, resp, true)
}

func aggregateSearchReferenceCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Aggregate Search Reference Criteria.",
		Attributes: map[string]schema.Attribute{
			"all_included_search_reference_criteria": schema.SetAttribute{
				Description: "Specifies a search reference criteria object that must match the associated search result reference in order to match the aggregate search reference criteria. If one or more all-included search reference criteria objects are provided, then a search result reference must match all of them in order to match the aggregate search reference criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"any_included_search_reference_criteria": schema.SetAttribute{
				Description: "Specifies a search reference criteria object that may match the associated search result reference in order to match the aggregate search reference criteria. If one or more any-included search reference criteria objects are provided, then a search result reference must match at least one of them in order to match the aggregate search reference criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"not_all_included_search_reference_criteria": schema.SetAttribute{
				Description: "Specifies a search reference criteria object that should not match the associated search result reference in order to match the aggregate search reference criteria. If one or more not-all-included search reference criteria objects are provided, then a search result reference must not match all of them (that is, it may match zero or more of them, but it must not match all of them) in order to match the aggregate search reference criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"none_included_search_reference_criteria": schema.SetAttribute{
				Description: "Specifies a search reference criteria object that must not match the associated search result reference in order to match the aggregate search reference criteria. If one or more none-included search reference criteria objects are provided, then a search result reference must not match any of them in order to match the aggregate search reference criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Search Reference Criteria",
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
func addOptionalAggregateSearchReferenceCriteriaFields(ctx context.Context, addRequest *client.AddAggregateSearchReferenceCriteriaRequest, plan aggregateSearchReferenceCriteriaResourceModel) {
	if internaltypes.IsDefined(plan.AllIncludedSearchReferenceCriteria) {
		var slice []string
		plan.AllIncludedSearchReferenceCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedSearchReferenceCriteria = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedSearchReferenceCriteria) {
		var slice []string
		plan.AnyIncludedSearchReferenceCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedSearchReferenceCriteria = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedSearchReferenceCriteria) {
		var slice []string
		plan.NotAllIncludedSearchReferenceCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedSearchReferenceCriteria = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedSearchReferenceCriteria) {
		var slice []string
		plan.NoneIncludedSearchReferenceCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedSearchReferenceCriteria = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a AggregateSearchReferenceCriteriaResponse object into the model struct
func readAggregateSearchReferenceCriteriaResponse(ctx context.Context, r *client.AggregateSearchReferenceCriteriaResponse, state *aggregateSearchReferenceCriteriaResourceModel, expectedValues *aggregateSearchReferenceCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.AllIncludedSearchReferenceCriteria = internaltypes.GetStringSet(r.AllIncludedSearchReferenceCriteria)
	state.AnyIncludedSearchReferenceCriteria = internaltypes.GetStringSet(r.AnyIncludedSearchReferenceCriteria)
	state.NotAllIncludedSearchReferenceCriteria = internaltypes.GetStringSet(r.NotAllIncludedSearchReferenceCriteria)
	state.NoneIncludedSearchReferenceCriteria = internaltypes.GetStringSet(r.NoneIncludedSearchReferenceCriteria)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createAggregateSearchReferenceCriteriaOperations(plan aggregateSearchReferenceCriteriaResourceModel, state aggregateSearchReferenceCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedSearchReferenceCriteria, state.AllIncludedSearchReferenceCriteria, "all-included-search-reference-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedSearchReferenceCriteria, state.AnyIncludedSearchReferenceCriteria, "any-included-search-reference-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedSearchReferenceCriteria, state.NotAllIncludedSearchReferenceCriteria, "not-all-included-search-reference-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedSearchReferenceCriteria, state.NoneIncludedSearchReferenceCriteria, "none-included-search-reference-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *aggregateSearchReferenceCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan aggregateSearchReferenceCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddAggregateSearchReferenceCriteriaRequest(plan.Id.ValueString(),
		[]client.EnumaggregateSearchReferenceCriteriaSchemaUrn{client.ENUMAGGREGATESEARCHREFERENCECRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SEARCH_REFERENCE_CRITERIAAGGREGATE})
	addOptionalAggregateSearchReferenceCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.SearchReferenceCriteriaApi.AddSearchReferenceCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddSearchReferenceCriteriaRequest(
		client.AddAggregateSearchReferenceCriteriaRequestAsAddSearchReferenceCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.SearchReferenceCriteriaApi.AddSearchReferenceCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Aggregate Search Reference Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state aggregateSearchReferenceCriteriaResourceModel
	readAggregateSearchReferenceCriteriaResponse(ctx, addResponse.AggregateSearchReferenceCriteriaResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultAggregateSearchReferenceCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan aggregateSearchReferenceCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SearchReferenceCriteriaApi.GetSearchReferenceCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Aggregate Search Reference Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state aggregateSearchReferenceCriteriaResourceModel
	readAggregateSearchReferenceCriteriaResponse(ctx, readResponse.AggregateSearchReferenceCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.SearchReferenceCriteriaApi.UpdateSearchReferenceCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createAggregateSearchReferenceCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.SearchReferenceCriteriaApi.UpdateSearchReferenceCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Aggregate Search Reference Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAggregateSearchReferenceCriteriaResponse(ctx, updateResponse.AggregateSearchReferenceCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *aggregateSearchReferenceCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAggregateSearchReferenceCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAggregateSearchReferenceCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAggregateSearchReferenceCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readAggregateSearchReferenceCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state aggregateSearchReferenceCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.SearchReferenceCriteriaApi.GetSearchReferenceCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Aggregate Search Reference Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAggregateSearchReferenceCriteriaResponse(ctx, readResponse.AggregateSearchReferenceCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *aggregateSearchReferenceCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAggregateSearchReferenceCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAggregateSearchReferenceCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAggregateSearchReferenceCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateAggregateSearchReferenceCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan aggregateSearchReferenceCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state aggregateSearchReferenceCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.SearchReferenceCriteriaApi.UpdateSearchReferenceCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createAggregateSearchReferenceCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.SearchReferenceCriteriaApi.UpdateSearchReferenceCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Aggregate Search Reference Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAggregateSearchReferenceCriteriaResponse(ctx, updateResponse.AggregateSearchReferenceCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultAggregateSearchReferenceCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *aggregateSearchReferenceCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state aggregateSearchReferenceCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.SearchReferenceCriteriaApi.DeleteSearchReferenceCriteriaExecute(r.apiClient.SearchReferenceCriteriaApi.DeleteSearchReferenceCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Aggregate Search Reference Criteria", err, httpResp)
		return
	}
}

func (r *aggregateSearchReferenceCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAggregateSearchReferenceCriteria(ctx, req, resp)
}

func (r *defaultAggregateSearchReferenceCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAggregateSearchReferenceCriteria(ctx, req, resp)
}

func importAggregateSearchReferenceCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
