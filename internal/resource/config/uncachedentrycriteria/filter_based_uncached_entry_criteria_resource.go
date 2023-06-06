package uncachedentrycriteria

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &filterBasedUncachedEntryCriteriaResource{}
	_ resource.ResourceWithConfigure   = &filterBasedUncachedEntryCriteriaResource{}
	_ resource.ResourceWithImportState = &filterBasedUncachedEntryCriteriaResource{}
	_ resource.Resource                = &defaultFilterBasedUncachedEntryCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultFilterBasedUncachedEntryCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultFilterBasedUncachedEntryCriteriaResource{}
)

// Create a Filter Based Uncached Entry Criteria resource
func NewFilterBasedUncachedEntryCriteriaResource() resource.Resource {
	return &filterBasedUncachedEntryCriteriaResource{}
}

func NewDefaultFilterBasedUncachedEntryCriteriaResource() resource.Resource {
	return &defaultFilterBasedUncachedEntryCriteriaResource{}
}

// filterBasedUncachedEntryCriteriaResource is the resource implementation.
type filterBasedUncachedEntryCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultFilterBasedUncachedEntryCriteriaResource is the resource implementation.
type defaultFilterBasedUncachedEntryCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *filterBasedUncachedEntryCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filter_based_uncached_entry_criteria"
}

func (r *defaultFilterBasedUncachedEntryCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_filter_based_uncached_entry_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *filterBasedUncachedEntryCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultFilterBasedUncachedEntryCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type filterBasedUncachedEntryCriteriaResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	LastUpdated                     types.String `tfsdk:"last_updated"`
	Notifications                   types.Set    `tfsdk:"notifications"`
	RequiredActions                 types.Set    `tfsdk:"required_actions"`
	Filter                          types.String `tfsdk:"filter"`
	FilterIdentifiesUncachedEntries types.Bool   `tfsdk:"filter_identifies_uncached_entries"`
	Description                     types.String `tfsdk:"description"`
	Enabled                         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *filterBasedUncachedEntryCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	filterBasedUncachedEntryCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultFilterBasedUncachedEntryCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	filterBasedUncachedEntryCriteriaSchema(ctx, req, resp, true)
}

func filterBasedUncachedEntryCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Filter Based Uncached Entry Criteria.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "Specifies the search filter that should be used to differentiate entries into cached and uncached sets.",
				Required:    true,
			},
			"filter_identifies_uncached_entries": schema.BoolAttribute{
				Description: "Indicates whether the associated filter identifies those entries which should be stored in the uncached-id2entry database (if true) or entries which should be stored in the id2entry database (if false).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Uncached Entry Criteria",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Uncached Entry Criteria is enabled for use in the server.",
				Required:    true,
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
func addOptionalFilterBasedUncachedEntryCriteriaFields(ctx context.Context, addRequest *client.AddFilterBasedUncachedEntryCriteriaRequest, plan filterBasedUncachedEntryCriteriaResourceModel) {
	if internaltypes.IsDefined(plan.FilterIdentifiesUncachedEntries) {
		addRequest.FilterIdentifiesUncachedEntries = plan.FilterIdentifiesUncachedEntries.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a FilterBasedUncachedEntryCriteriaResponse object into the model struct
func readFilterBasedUncachedEntryCriteriaResponse(ctx context.Context, r *client.FilterBasedUncachedEntryCriteriaResponse, state *filterBasedUncachedEntryCriteriaResourceModel, expectedValues *filterBasedUncachedEntryCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Filter = types.StringValue(r.Filter)
	state.FilterIdentifiesUncachedEntries = types.BoolValue(r.FilterIdentifiesUncachedEntries)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createFilterBasedUncachedEntryCriteriaOperations(plan filterBasedUncachedEntryCriteriaResourceModel, state filterBasedUncachedEntryCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Filter, state.Filter, "filter")
	operations.AddBoolOperationIfNecessary(&ops, plan.FilterIdentifiesUncachedEntries, state.FilterIdentifiesUncachedEntries, "filter-identifies-uncached-entries")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *filterBasedUncachedEntryCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan filterBasedUncachedEntryCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddFilterBasedUncachedEntryCriteriaRequest(plan.Id.ValueString(),
		[]client.EnumfilterBasedUncachedEntryCriteriaSchemaUrn{client.ENUMFILTERBASEDUNCACHEDENTRYCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0UNCACHED_ENTRY_CRITERIAFILTER_BASED},
		plan.Filter.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalFilterBasedUncachedEntryCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.UncachedEntryCriteriaApi.AddUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddUncachedEntryCriteriaRequest(
		client.AddFilterBasedUncachedEntryCriteriaRequestAsAddUncachedEntryCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.UncachedEntryCriteriaApi.AddUncachedEntryCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Filter Based Uncached Entry Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state filterBasedUncachedEntryCriteriaResourceModel
	readFilterBasedUncachedEntryCriteriaResponse(ctx, addResponse.FilterBasedUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultFilterBasedUncachedEntryCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan filterBasedUncachedEntryCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.UncachedEntryCriteriaApi.GetUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Filter Based Uncached Entry Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state filterBasedUncachedEntryCriteriaResourceModel
	readFilterBasedUncachedEntryCriteriaResponse(ctx, readResponse.FilterBasedUncachedEntryCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.UncachedEntryCriteriaApi.UpdateUncachedEntryCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createFilterBasedUncachedEntryCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.UncachedEntryCriteriaApi.UpdateUncachedEntryCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Filter Based Uncached Entry Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFilterBasedUncachedEntryCriteriaResponse(ctx, updateResponse.FilterBasedUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *filterBasedUncachedEntryCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFilterBasedUncachedEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFilterBasedUncachedEntryCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFilterBasedUncachedEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readFilterBasedUncachedEntryCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state filterBasedUncachedEntryCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.UncachedEntryCriteriaApi.GetUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Filter Based Uncached Entry Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readFilterBasedUncachedEntryCriteriaResponse(ctx, readResponse.FilterBasedUncachedEntryCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *filterBasedUncachedEntryCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFilterBasedUncachedEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFilterBasedUncachedEntryCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFilterBasedUncachedEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateFilterBasedUncachedEntryCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan filterBasedUncachedEntryCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state filterBasedUncachedEntryCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.UncachedEntryCriteriaApi.UpdateUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createFilterBasedUncachedEntryCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.UncachedEntryCriteriaApi.UpdateUncachedEntryCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Filter Based Uncached Entry Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFilterBasedUncachedEntryCriteriaResponse(ctx, updateResponse.FilterBasedUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultFilterBasedUncachedEntryCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *filterBasedUncachedEntryCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state filterBasedUncachedEntryCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.UncachedEntryCriteriaApi.DeleteUncachedEntryCriteriaExecute(r.apiClient.UncachedEntryCriteriaApi.DeleteUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Filter Based Uncached Entry Criteria", err, httpResp)
		return
	}
}

func (r *filterBasedUncachedEntryCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFilterBasedUncachedEntryCriteria(ctx, req, resp)
}

func (r *defaultFilterBasedUncachedEntryCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFilterBasedUncachedEntryCriteria(ctx, req, resp)
}

func importFilterBasedUncachedEntryCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
