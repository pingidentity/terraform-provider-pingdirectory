package uncachedentrycriteria

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
	_ resource.Resource                = &lastAccessTimeUncachedEntryCriteriaResource{}
	_ resource.ResourceWithConfigure   = &lastAccessTimeUncachedEntryCriteriaResource{}
	_ resource.ResourceWithImportState = &lastAccessTimeUncachedEntryCriteriaResource{}
	_ resource.Resource                = &defaultLastAccessTimeUncachedEntryCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultLastAccessTimeUncachedEntryCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultLastAccessTimeUncachedEntryCriteriaResource{}
)

// Create a Last Access Time Uncached Entry Criteria resource
func NewLastAccessTimeUncachedEntryCriteriaResource() resource.Resource {
	return &lastAccessTimeUncachedEntryCriteriaResource{}
}

func NewDefaultLastAccessTimeUncachedEntryCriteriaResource() resource.Resource {
	return &defaultLastAccessTimeUncachedEntryCriteriaResource{}
}

// lastAccessTimeUncachedEntryCriteriaResource is the resource implementation.
type lastAccessTimeUncachedEntryCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultLastAccessTimeUncachedEntryCriteriaResource is the resource implementation.
type defaultLastAccessTimeUncachedEntryCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *lastAccessTimeUncachedEntryCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_last_access_time_uncached_entry_criteria"
}

func (r *defaultLastAccessTimeUncachedEntryCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_last_access_time_uncached_entry_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *lastAccessTimeUncachedEntryCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultLastAccessTimeUncachedEntryCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type lastAccessTimeUncachedEntryCriteriaResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	LastUpdated         types.String `tfsdk:"last_updated"`
	Notifications       types.Set    `tfsdk:"notifications"`
	RequiredActions     types.Set    `tfsdk:"required_actions"`
	AccessTimeThreshold types.String `tfsdk:"access_time_threshold"`
	Description         types.String `tfsdk:"description"`
	Enabled             types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *lastAccessTimeUncachedEntryCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	lastAccessTimeUncachedEntryCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultLastAccessTimeUncachedEntryCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	lastAccessTimeUncachedEntryCriteriaSchema(ctx, req, resp, true)
}

func lastAccessTimeUncachedEntryCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Last Access Time Uncached Entry Criteria.",
		Attributes: map[string]schema.Attribute{
			"access_time_threshold": schema.StringAttribute{
				Description: "Specifies the maximum length of time that has passed since an entry was last accessed that it should still be included in the id2entry database. Entries that have not been accessed in more than this length of time may be written into the uncached-id2entry database.",
				Required:    true,
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
func addOptionalLastAccessTimeUncachedEntryCriteriaFields(ctx context.Context, addRequest *client.AddLastAccessTimeUncachedEntryCriteriaRequest, plan lastAccessTimeUncachedEntryCriteriaResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a LastAccessTimeUncachedEntryCriteriaResponse object into the model struct
func readLastAccessTimeUncachedEntryCriteriaResponse(ctx context.Context, r *client.LastAccessTimeUncachedEntryCriteriaResponse, state *lastAccessTimeUncachedEntryCriteriaResourceModel, expectedValues *lastAccessTimeUncachedEntryCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.AccessTimeThreshold = types.StringValue(r.AccessTimeThreshold)
	config.CheckMismatchedPDFormattedAttributes("access_time_threshold",
		expectedValues.AccessTimeThreshold, state.AccessTimeThreshold, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createLastAccessTimeUncachedEntryCriteriaOperations(plan lastAccessTimeUncachedEntryCriteriaResourceModel, state lastAccessTimeUncachedEntryCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.AccessTimeThreshold, state.AccessTimeThreshold, "access-time-threshold")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *lastAccessTimeUncachedEntryCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan lastAccessTimeUncachedEntryCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddLastAccessTimeUncachedEntryCriteriaRequest(plan.Id.ValueString(),
		[]client.EnumlastAccessTimeUncachedEntryCriteriaSchemaUrn{client.ENUMLASTACCESSTIMEUNCACHEDENTRYCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0UNCACHED_ENTRY_CRITERIALAST_ACCESS_TIME},
		plan.AccessTimeThreshold.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalLastAccessTimeUncachedEntryCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.UncachedEntryCriteriaApi.AddUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddUncachedEntryCriteriaRequest(
		client.AddLastAccessTimeUncachedEntryCriteriaRequestAsAddUncachedEntryCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.UncachedEntryCriteriaApi.AddUncachedEntryCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Last Access Time Uncached Entry Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state lastAccessTimeUncachedEntryCriteriaResourceModel
	readLastAccessTimeUncachedEntryCriteriaResponse(ctx, addResponse.LastAccessTimeUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultLastAccessTimeUncachedEntryCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan lastAccessTimeUncachedEntryCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.UncachedEntryCriteriaApi.GetUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Last Access Time Uncached Entry Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state lastAccessTimeUncachedEntryCriteriaResourceModel
	readLastAccessTimeUncachedEntryCriteriaResponse(ctx, readResponse.LastAccessTimeUncachedEntryCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.UncachedEntryCriteriaApi.UpdateUncachedEntryCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createLastAccessTimeUncachedEntryCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.UncachedEntryCriteriaApi.UpdateUncachedEntryCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Last Access Time Uncached Entry Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLastAccessTimeUncachedEntryCriteriaResponse(ctx, updateResponse.LastAccessTimeUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *lastAccessTimeUncachedEntryCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLastAccessTimeUncachedEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLastAccessTimeUncachedEntryCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLastAccessTimeUncachedEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readLastAccessTimeUncachedEntryCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state lastAccessTimeUncachedEntryCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.UncachedEntryCriteriaApi.GetUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Last Access Time Uncached Entry Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLastAccessTimeUncachedEntryCriteriaResponse(ctx, readResponse.LastAccessTimeUncachedEntryCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *lastAccessTimeUncachedEntryCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLastAccessTimeUncachedEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLastAccessTimeUncachedEntryCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLastAccessTimeUncachedEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateLastAccessTimeUncachedEntryCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan lastAccessTimeUncachedEntryCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state lastAccessTimeUncachedEntryCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.UncachedEntryCriteriaApi.UpdateUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createLastAccessTimeUncachedEntryCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.UncachedEntryCriteriaApi.UpdateUncachedEntryCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Last Access Time Uncached Entry Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLastAccessTimeUncachedEntryCriteriaResponse(ctx, updateResponse.LastAccessTimeUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultLastAccessTimeUncachedEntryCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *lastAccessTimeUncachedEntryCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state lastAccessTimeUncachedEntryCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.UncachedEntryCriteriaApi.DeleteUncachedEntryCriteriaExecute(r.apiClient.UncachedEntryCriteriaApi.DeleteUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Last Access Time Uncached Entry Criteria", err, httpResp)
		return
	}
}

func (r *lastAccessTimeUncachedEntryCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLastAccessTimeUncachedEntryCriteria(ctx, req, resp)
}

func (r *defaultLastAccessTimeUncachedEntryCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLastAccessTimeUncachedEntryCriteria(ctx, req, resp)
}

func importLastAccessTimeUncachedEntryCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
