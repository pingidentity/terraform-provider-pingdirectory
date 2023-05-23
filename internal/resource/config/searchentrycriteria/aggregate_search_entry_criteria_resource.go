package searchentrycriteria

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
	_ resource.Resource                = &aggregateSearchEntryCriteriaResource{}
	_ resource.ResourceWithConfigure   = &aggregateSearchEntryCriteriaResource{}
	_ resource.ResourceWithImportState = &aggregateSearchEntryCriteriaResource{}
	_ resource.Resource                = &defaultAggregateSearchEntryCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultAggregateSearchEntryCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultAggregateSearchEntryCriteriaResource{}
)

// Create a Aggregate Search Entry Criteria resource
func NewAggregateSearchEntryCriteriaResource() resource.Resource {
	return &aggregateSearchEntryCriteriaResource{}
}

func NewDefaultAggregateSearchEntryCriteriaResource() resource.Resource {
	return &defaultAggregateSearchEntryCriteriaResource{}
}

// aggregateSearchEntryCriteriaResource is the resource implementation.
type aggregateSearchEntryCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultAggregateSearchEntryCriteriaResource is the resource implementation.
type defaultAggregateSearchEntryCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *aggregateSearchEntryCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aggregate_search_entry_criteria"
}

func (r *defaultAggregateSearchEntryCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_aggregate_search_entry_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *aggregateSearchEntryCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultAggregateSearchEntryCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type aggregateSearchEntryCriteriaResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	LastUpdated                       types.String `tfsdk:"last_updated"`
	Notifications                     types.Set    `tfsdk:"notifications"`
	RequiredActions                   types.Set    `tfsdk:"required_actions"`
	AllIncludedSearchEntryCriteria    types.Set    `tfsdk:"all_included_search_entry_criteria"`
	AnyIncludedSearchEntryCriteria    types.Set    `tfsdk:"any_included_search_entry_criteria"`
	NotAllIncludedSearchEntryCriteria types.Set    `tfsdk:"not_all_included_search_entry_criteria"`
	NoneIncludedSearchEntryCriteria   types.Set    `tfsdk:"none_included_search_entry_criteria"`
	Description                       types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *aggregateSearchEntryCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	aggregateSearchEntryCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultAggregateSearchEntryCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	aggregateSearchEntryCriteriaSchema(ctx, req, resp, true)
}

func aggregateSearchEntryCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Aggregate Search Entry Criteria.",
		Attributes: map[string]schema.Attribute{
			"all_included_search_entry_criteria": schema.SetAttribute{
				Description: "Specifies a search entry criteria object that must match the associated search result entry in order to match the aggregate search entry criteria. If one or more all-included search entry criteria objects are provided, then a search result entry must match all of them in order to match the aggregate search entry criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"any_included_search_entry_criteria": schema.SetAttribute{
				Description: "Specifies a search entry criteria object that may match the associated search result entry in order to match the aggregate search entry criteria. If one or more any-included search entry criteria objects are provided, then a search result entry must match at least one of them in order to match the aggregate search entry criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"not_all_included_search_entry_criteria": schema.SetAttribute{
				Description: "Specifies a search entry criteria object that should not match the associated search result entry in order to match the aggregate search entry criteria. If one or more not-all-included search entry criteria objects are provided, then a search result entry must not match all of them (that is, it may match zero or more of them, but it must not match all of them) in order to match the aggregate search entry criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"none_included_search_entry_criteria": schema.SetAttribute{
				Description: "Specifies a search entry criteria object that must not match the associated search result entry in order to match the aggregate search entry criteria. If one or more none-included search entry criteria objects are provided, then a search result entry must not match any of them in order to match the aggregate search entry criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Search Entry Criteria",
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
func addOptionalAggregateSearchEntryCriteriaFields(ctx context.Context, addRequest *client.AddAggregateSearchEntryCriteriaRequest, plan aggregateSearchEntryCriteriaResourceModel) {
	if internaltypes.IsDefined(plan.AllIncludedSearchEntryCriteria) {
		var slice []string
		plan.AllIncludedSearchEntryCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedSearchEntryCriteria = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedSearchEntryCriteria) {
		var slice []string
		plan.AnyIncludedSearchEntryCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedSearchEntryCriteria = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedSearchEntryCriteria) {
		var slice []string
		plan.NotAllIncludedSearchEntryCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedSearchEntryCriteria = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedSearchEntryCriteria) {
		var slice []string
		plan.NoneIncludedSearchEntryCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedSearchEntryCriteria = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a AggregateSearchEntryCriteriaResponse object into the model struct
func readAggregateSearchEntryCriteriaResponse(ctx context.Context, r *client.AggregateSearchEntryCriteriaResponse, state *aggregateSearchEntryCriteriaResourceModel, expectedValues *aggregateSearchEntryCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.AllIncludedSearchEntryCriteria = internaltypes.GetStringSet(r.AllIncludedSearchEntryCriteria)
	state.AnyIncludedSearchEntryCriteria = internaltypes.GetStringSet(r.AnyIncludedSearchEntryCriteria)
	state.NotAllIncludedSearchEntryCriteria = internaltypes.GetStringSet(r.NotAllIncludedSearchEntryCriteria)
	state.NoneIncludedSearchEntryCriteria = internaltypes.GetStringSet(r.NoneIncludedSearchEntryCriteria)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createAggregateSearchEntryCriteriaOperations(plan aggregateSearchEntryCriteriaResourceModel, state aggregateSearchEntryCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedSearchEntryCriteria, state.AllIncludedSearchEntryCriteria, "all-included-search-entry-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedSearchEntryCriteria, state.AnyIncludedSearchEntryCriteria, "any-included-search-entry-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedSearchEntryCriteria, state.NotAllIncludedSearchEntryCriteria, "not-all-included-search-entry-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedSearchEntryCriteria, state.NoneIncludedSearchEntryCriteria, "none-included-search-entry-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *aggregateSearchEntryCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan aggregateSearchEntryCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddAggregateSearchEntryCriteriaRequest(plan.Id.ValueString(),
		[]client.EnumaggregateSearchEntryCriteriaSchemaUrn{client.ENUMAGGREGATESEARCHENTRYCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SEARCH_ENTRY_CRITERIAAGGREGATE})
	addOptionalAggregateSearchEntryCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.SearchEntryCriteriaApi.AddSearchEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddSearchEntryCriteriaRequest(
		client.AddAggregateSearchEntryCriteriaRequestAsAddSearchEntryCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.SearchEntryCriteriaApi.AddSearchEntryCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Aggregate Search Entry Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state aggregateSearchEntryCriteriaResourceModel
	readAggregateSearchEntryCriteriaResponse(ctx, addResponse.AggregateSearchEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultAggregateSearchEntryCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan aggregateSearchEntryCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SearchEntryCriteriaApi.GetSearchEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Aggregate Search Entry Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state aggregateSearchEntryCriteriaResourceModel
	readAggregateSearchEntryCriteriaResponse(ctx, readResponse.AggregateSearchEntryCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.SearchEntryCriteriaApi.UpdateSearchEntryCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createAggregateSearchEntryCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.SearchEntryCriteriaApi.UpdateSearchEntryCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Aggregate Search Entry Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAggregateSearchEntryCriteriaResponse(ctx, updateResponse.AggregateSearchEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *aggregateSearchEntryCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAggregateSearchEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAggregateSearchEntryCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAggregateSearchEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readAggregateSearchEntryCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state aggregateSearchEntryCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.SearchEntryCriteriaApi.GetSearchEntryCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Aggregate Search Entry Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAggregateSearchEntryCriteriaResponse(ctx, readResponse.AggregateSearchEntryCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *aggregateSearchEntryCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAggregateSearchEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAggregateSearchEntryCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAggregateSearchEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateAggregateSearchEntryCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan aggregateSearchEntryCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state aggregateSearchEntryCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.SearchEntryCriteriaApi.UpdateSearchEntryCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createAggregateSearchEntryCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.SearchEntryCriteriaApi.UpdateSearchEntryCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Aggregate Search Entry Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAggregateSearchEntryCriteriaResponse(ctx, updateResponse.AggregateSearchEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultAggregateSearchEntryCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *aggregateSearchEntryCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state aggregateSearchEntryCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.SearchEntryCriteriaApi.DeleteSearchEntryCriteriaExecute(r.apiClient.SearchEntryCriteriaApi.DeleteSearchEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Aggregate Search Entry Criteria", err, httpResp)
		return
	}
}

func (r *aggregateSearchEntryCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAggregateSearchEntryCriteria(ctx, req, resp)
}

func (r *defaultAggregateSearchEntryCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAggregateSearchEntryCriteria(ctx, req, resp)
}

func importAggregateSearchEntryCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
