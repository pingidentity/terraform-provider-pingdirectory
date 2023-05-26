package uncachedattributecriteria

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
	_ resource.Resource                = &defaultUncachedAttributeCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultUncachedAttributeCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultUncachedAttributeCriteriaResource{}
	_ resource.Resource                = &defaultDefaultUncachedAttributeCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultDefaultUncachedAttributeCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultDefaultUncachedAttributeCriteriaResource{}
)

// Create a Default Uncached Attribute Criteria resource
func NewDefaultUncachedAttributeCriteriaResource() resource.Resource {
	return &defaultUncachedAttributeCriteriaResource{}
}

func NewDefaultDefaultUncachedAttributeCriteriaResource() resource.Resource {
	return &defaultDefaultUncachedAttributeCriteriaResource{}
}

// defaultUncachedAttributeCriteriaResource is the resource implementation.
type defaultUncachedAttributeCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultDefaultUncachedAttributeCriteriaResource is the resource implementation.
type defaultDefaultUncachedAttributeCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *defaultUncachedAttributeCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_uncached_attribute_criteria"
}

func (r *defaultDefaultUncachedAttributeCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_default_uncached_attribute_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *defaultUncachedAttributeCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultDefaultUncachedAttributeCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type defaultUncachedAttributeCriteriaResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	Description     types.String `tfsdk:"description"`
	Enabled         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *defaultUncachedAttributeCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	defaultUncachedAttributeCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultDefaultUncachedAttributeCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	defaultUncachedAttributeCriteriaSchema(ctx, req, resp, true)
}

func defaultUncachedAttributeCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Default Uncached Attribute Criteria.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description: "A description for this Uncached Attribute Criteria",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Uncached Attribute Criteria is enabled for use in the server.",
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
func addOptionalDefaultUncachedAttributeCriteriaFields(ctx context.Context, addRequest *client.AddDefaultUncachedAttributeCriteriaRequest, plan defaultUncachedAttributeCriteriaResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a DefaultUncachedAttributeCriteriaResponse object into the model struct
func readDefaultUncachedAttributeCriteriaResponse(ctx context.Context, r *client.DefaultUncachedAttributeCriteriaResponse, state *defaultUncachedAttributeCriteriaResourceModel, expectedValues *defaultUncachedAttributeCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createDefaultUncachedAttributeCriteriaOperations(plan defaultUncachedAttributeCriteriaResourceModel, state defaultUncachedAttributeCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *defaultUncachedAttributeCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan defaultUncachedAttributeCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddDefaultUncachedAttributeCriteriaRequest(plan.Id.ValueString(),
		[]client.EnumdefaultUncachedAttributeCriteriaSchemaUrn{client.ENUMDEFAULTUNCACHEDATTRIBUTECRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0UNCACHED_ATTRIBUTE_CRITERIADEFAULT},
		plan.Enabled.ValueBool())
	addOptionalDefaultUncachedAttributeCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.UncachedAttributeCriteriaApi.AddUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddUncachedAttributeCriteriaRequest(
		client.AddDefaultUncachedAttributeCriteriaRequestAsAddUncachedAttributeCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.AddUncachedAttributeCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Default Uncached Attribute Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state defaultUncachedAttributeCriteriaResourceModel
	readDefaultUncachedAttributeCriteriaResponse(ctx, addResponse.DefaultUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultDefaultUncachedAttributeCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan defaultUncachedAttributeCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.GetUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Default Uncached Attribute Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state defaultUncachedAttributeCriteriaResourceModel
	readDefaultUncachedAttributeCriteriaResponse(ctx, readResponse.DefaultUncachedAttributeCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.UncachedAttributeCriteriaApi.UpdateUncachedAttributeCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createDefaultUncachedAttributeCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.UpdateUncachedAttributeCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Default Uncached Attribute Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDefaultUncachedAttributeCriteriaResponse(ctx, updateResponse.DefaultUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultUncachedAttributeCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDefaultUncachedAttributeCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDefaultUncachedAttributeCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDefaultUncachedAttributeCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readDefaultUncachedAttributeCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state defaultUncachedAttributeCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.UncachedAttributeCriteriaApi.GetUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Default Uncached Attribute Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDefaultUncachedAttributeCriteriaResponse(ctx, readResponse.DefaultUncachedAttributeCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *defaultUncachedAttributeCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDefaultUncachedAttributeCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDefaultUncachedAttributeCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDefaultUncachedAttributeCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateDefaultUncachedAttributeCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan defaultUncachedAttributeCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state defaultUncachedAttributeCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.UncachedAttributeCriteriaApi.UpdateUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createDefaultUncachedAttributeCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.UncachedAttributeCriteriaApi.UpdateUncachedAttributeCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Default Uncached Attribute Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDefaultUncachedAttributeCriteriaResponse(ctx, updateResponse.DefaultUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultDefaultUncachedAttributeCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *defaultUncachedAttributeCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state defaultUncachedAttributeCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.DeleteUncachedAttributeCriteriaExecute(r.apiClient.UncachedAttributeCriteriaApi.DeleteUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Default Uncached Attribute Criteria", err, httpResp)
		return
	}
}

func (r *defaultUncachedAttributeCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDefaultUncachedAttributeCriteria(ctx, req, resp)
}

func (r *defaultDefaultUncachedAttributeCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDefaultUncachedAttributeCriteria(ctx, req, resp)
}

func importDefaultUncachedAttributeCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
