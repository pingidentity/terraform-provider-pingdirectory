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
	_ resource.Resource                = &simpleSearchReferenceCriteriaResource{}
	_ resource.ResourceWithConfigure   = &simpleSearchReferenceCriteriaResource{}
	_ resource.ResourceWithImportState = &simpleSearchReferenceCriteriaResource{}
	_ resource.Resource                = &defaultSimpleSearchReferenceCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultSimpleSearchReferenceCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultSimpleSearchReferenceCriteriaResource{}
)

// Create a Simple Search Reference Criteria resource
func NewSimpleSearchReferenceCriteriaResource() resource.Resource {
	return &simpleSearchReferenceCriteriaResource{}
}

func NewDefaultSimpleSearchReferenceCriteriaResource() resource.Resource {
	return &defaultSimpleSearchReferenceCriteriaResource{}
}

// simpleSearchReferenceCriteriaResource is the resource implementation.
type simpleSearchReferenceCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSimpleSearchReferenceCriteriaResource is the resource implementation.
type defaultSimpleSearchReferenceCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *simpleSearchReferenceCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_simple_search_reference_criteria"
}

func (r *defaultSimpleSearchReferenceCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_simple_search_reference_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *simpleSearchReferenceCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultSimpleSearchReferenceCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type simpleSearchReferenceCriteriaResourceModel struct {
	Id                             types.String `tfsdk:"id"`
	LastUpdated                    types.String `tfsdk:"last_updated"`
	Notifications                  types.Set    `tfsdk:"notifications"`
	RequiredActions                types.Set    `tfsdk:"required_actions"`
	RequestCriteria                types.String `tfsdk:"request_criteria"`
	AllIncludedReferenceControl    types.Set    `tfsdk:"all_included_reference_control"`
	AnyIncludedReferenceControl    types.Set    `tfsdk:"any_included_reference_control"`
	NotAllIncludedReferenceControl types.Set    `tfsdk:"not_all_included_reference_control"`
	NoneIncludedReferenceControl   types.Set    `tfsdk:"none_included_reference_control"`
	Description                    types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *simpleSearchReferenceCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	simpleSearchReferenceCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultSimpleSearchReferenceCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	simpleSearchReferenceCriteriaSchema(ctx, req, resp, true)
}

func simpleSearchReferenceCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Simple Search Reference Criteria.",
		Attributes: map[string]schema.Attribute{
			"request_criteria": schema.StringAttribute{
				Description: "Specifies a request criteria object that must match the associated request for references included in this Simple Search Reference Criteria.",
				Optional:    true,
			},
			"all_included_reference_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must be present in search result references included in this Simple Search Reference Criteria. If any control OIDs are provided, then the reference must contain all of those controls.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"any_included_reference_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that may be present in search result references included in this Simple Search Reference Criteria. If any control OIDs are provided, then the reference must contain at least one of those controls.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"not_all_included_reference_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that should not be present in search result references included in this Simple Search Reference Criteria. If any control OIDs are provided, then the reference must not contain at least one of those controls (that is, it may contain zero or more of those controls, but not all of them).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"none_included_reference_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must not be present in search result references included in this Simple Search Reference Criteria. If any control OIDs are provided, then the reference must not contain any of those controls.",
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
func addOptionalSimpleSearchReferenceCriteriaFields(ctx context.Context, addRequest *client.AddSimpleSearchReferenceCriteriaRequest, plan simpleSearchReferenceCriteriaResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AllIncludedReferenceControl) {
		var slice []string
		plan.AllIncludedReferenceControl.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedReferenceControl = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedReferenceControl) {
		var slice []string
		plan.AnyIncludedReferenceControl.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedReferenceControl = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedReferenceControl) {
		var slice []string
		plan.NotAllIncludedReferenceControl.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedReferenceControl = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedReferenceControl) {
		var slice []string
		plan.NoneIncludedReferenceControl.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedReferenceControl = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a SimpleSearchReferenceCriteriaResponse object into the model struct
func readSimpleSearchReferenceCriteriaResponse(ctx context.Context, r *client.SimpleSearchReferenceCriteriaResponse, state *simpleSearchReferenceCriteriaResourceModel, expectedValues *simpleSearchReferenceCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.AllIncludedReferenceControl = internaltypes.GetStringSet(r.AllIncludedReferenceControl)
	state.AnyIncludedReferenceControl = internaltypes.GetStringSet(r.AnyIncludedReferenceControl)
	state.NotAllIncludedReferenceControl = internaltypes.GetStringSet(r.NotAllIncludedReferenceControl)
	state.NoneIncludedReferenceControl = internaltypes.GetStringSet(r.NoneIncludedReferenceControl)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSimpleSearchReferenceCriteriaOperations(plan simpleSearchReferenceCriteriaResourceModel, state simpleSearchReferenceCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.RequestCriteria, state.RequestCriteria, "request-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedReferenceControl, state.AllIncludedReferenceControl, "all-included-reference-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedReferenceControl, state.AnyIncludedReferenceControl, "any-included-reference-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedReferenceControl, state.NotAllIncludedReferenceControl, "not-all-included-reference-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedReferenceControl, state.NoneIncludedReferenceControl, "none-included-reference-control")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *simpleSearchReferenceCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan simpleSearchReferenceCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddSimpleSearchReferenceCriteriaRequest(plan.Id.ValueString(),
		[]client.EnumsimpleSearchReferenceCriteriaSchemaUrn{client.ENUMSIMPLESEARCHREFERENCECRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SEARCH_REFERENCE_CRITERIASIMPLE})
	addOptionalSimpleSearchReferenceCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.SearchReferenceCriteriaApi.AddSearchReferenceCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddSearchReferenceCriteriaRequest(
		client.AddSimpleSearchReferenceCriteriaRequestAsAddSearchReferenceCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.SearchReferenceCriteriaApi.AddSearchReferenceCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Simple Search Reference Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state simpleSearchReferenceCriteriaResourceModel
	readSimpleSearchReferenceCriteriaResponse(ctx, addResponse.SimpleSearchReferenceCriteriaResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultSimpleSearchReferenceCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan simpleSearchReferenceCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SearchReferenceCriteriaApi.GetSearchReferenceCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Simple Search Reference Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state simpleSearchReferenceCriteriaResourceModel
	readSimpleSearchReferenceCriteriaResponse(ctx, readResponse.SimpleSearchReferenceCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.SearchReferenceCriteriaApi.UpdateSearchReferenceCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSimpleSearchReferenceCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.SearchReferenceCriteriaApi.UpdateSearchReferenceCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Simple Search Reference Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSimpleSearchReferenceCriteriaResponse(ctx, updateResponse.SimpleSearchReferenceCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *simpleSearchReferenceCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSimpleSearchReferenceCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSimpleSearchReferenceCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSimpleSearchReferenceCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSimpleSearchReferenceCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state simpleSearchReferenceCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.SearchReferenceCriteriaApi.GetSearchReferenceCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Simple Search Reference Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSimpleSearchReferenceCriteriaResponse(ctx, readResponse.SimpleSearchReferenceCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *simpleSearchReferenceCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSimpleSearchReferenceCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSimpleSearchReferenceCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSimpleSearchReferenceCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSimpleSearchReferenceCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan simpleSearchReferenceCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state simpleSearchReferenceCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.SearchReferenceCriteriaApi.UpdateSearchReferenceCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSimpleSearchReferenceCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.SearchReferenceCriteriaApi.UpdateSearchReferenceCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Simple Search Reference Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSimpleSearchReferenceCriteriaResponse(ctx, updateResponse.SimpleSearchReferenceCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSimpleSearchReferenceCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *simpleSearchReferenceCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state simpleSearchReferenceCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.SearchReferenceCriteriaApi.DeleteSearchReferenceCriteriaExecute(r.apiClient.SearchReferenceCriteriaApi.DeleteSearchReferenceCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Simple Search Reference Criteria", err, httpResp)
		return
	}
}

func (r *simpleSearchReferenceCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSimpleSearchReferenceCriteria(ctx, req, resp)
}

func (r *defaultSimpleSearchReferenceCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSimpleSearchReferenceCriteria(ctx, req, resp)
}

func importSimpleSearchReferenceCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
