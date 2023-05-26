package uncachedattributecriteria

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	_ resource.Resource                = &simpleUncachedAttributeCriteriaResource{}
	_ resource.ResourceWithConfigure   = &simpleUncachedAttributeCriteriaResource{}
	_ resource.ResourceWithImportState = &simpleUncachedAttributeCriteriaResource{}
	_ resource.Resource                = &defaultSimpleUncachedAttributeCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultSimpleUncachedAttributeCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultSimpleUncachedAttributeCriteriaResource{}
)

// Create a Simple Uncached Attribute Criteria resource
func NewSimpleUncachedAttributeCriteriaResource() resource.Resource {
	return &simpleUncachedAttributeCriteriaResource{}
}

func NewDefaultSimpleUncachedAttributeCriteriaResource() resource.Resource {
	return &defaultSimpleUncachedAttributeCriteriaResource{}
}

// simpleUncachedAttributeCriteriaResource is the resource implementation.
type simpleUncachedAttributeCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSimpleUncachedAttributeCriteriaResource is the resource implementation.
type defaultSimpleUncachedAttributeCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *simpleUncachedAttributeCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_simple_uncached_attribute_criteria"
}

func (r *defaultSimpleUncachedAttributeCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_simple_uncached_attribute_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *simpleUncachedAttributeCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultSimpleUncachedAttributeCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type simpleUncachedAttributeCriteriaResourceModel struct {
	Id                types.String `tfsdk:"id"`
	LastUpdated       types.String `tfsdk:"last_updated"`
	Notifications     types.Set    `tfsdk:"notifications"`
	RequiredActions   types.Set    `tfsdk:"required_actions"`
	AttributeType     types.Set    `tfsdk:"attribute_type"`
	MinValueCount     types.Int64  `tfsdk:"min_value_count"`
	MinTotalValueSize types.String `tfsdk:"min_total_value_size"`
	Description       types.String `tfsdk:"description"`
	Enabled           types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *simpleUncachedAttributeCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	simpleUncachedAttributeCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultSimpleUncachedAttributeCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	simpleUncachedAttributeCriteriaSchema(ctx, req, resp, true)
}

func simpleUncachedAttributeCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Simple Uncached Attribute Criteria.",
		Attributes: map[string]schema.Attribute{
			"attribute_type": schema.SetAttribute{
				Description: "Specifies the attribute types for attributes that may be written to the uncached-id2entry database.",
				Required:    true,
				ElementType: types.StringType,
			},
			"min_value_count": schema.Int64Attribute{
				Description: "Specifies the minimum number of values that an attribute must have before it will be written into the uncached-id2entry database.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"min_total_value_size": schema.StringAttribute{
				Description: "Specifies the minimum total value size (i.e., the sum of the sizes of all values) that an attribute must have before it will be written into the uncached-id2entry database.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
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
func addOptionalSimpleUncachedAttributeCriteriaFields(ctx context.Context, addRequest *client.AddSimpleUncachedAttributeCriteriaRequest, plan simpleUncachedAttributeCriteriaResourceModel) {
	if internaltypes.IsDefined(plan.MinValueCount) {
		addRequest.MinValueCount = plan.MinValueCount.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MinTotalValueSize) {
		addRequest.MinTotalValueSize = plan.MinTotalValueSize.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a SimpleUncachedAttributeCriteriaResponse object into the model struct
func readSimpleUncachedAttributeCriteriaResponse(ctx context.Context, r *client.SimpleUncachedAttributeCriteriaResponse, state *simpleUncachedAttributeCriteriaResourceModel, expectedValues *simpleUncachedAttributeCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.AttributeType = internaltypes.GetStringSet(r.AttributeType)
	state.MinValueCount = internaltypes.Int64TypeOrNil(r.MinValueCount)
	state.MinTotalValueSize = internaltypes.StringTypeOrNil(r.MinTotalValueSize, internaltypes.IsEmptyString(expectedValues.MinTotalValueSize))
	config.CheckMismatchedPDFormattedAttributes("min_total_value_size",
		expectedValues.MinTotalValueSize, state.MinTotalValueSize, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSimpleUncachedAttributeCriteriaOperations(plan simpleUncachedAttributeCriteriaResourceModel, state simpleUncachedAttributeCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AttributeType, state.AttributeType, "attribute-type")
	operations.AddInt64OperationIfNecessary(&ops, plan.MinValueCount, state.MinValueCount, "min-value-count")
	operations.AddStringOperationIfNecessary(&ops, plan.MinTotalValueSize, state.MinTotalValueSize, "min-total-value-size")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *simpleUncachedAttributeCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan simpleUncachedAttributeCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var AttributeTypeSlice []string
	plan.AttributeType.ElementsAs(ctx, &AttributeTypeSlice, false)
	addRequest := client.NewAddSimpleUncachedAttributeCriteriaRequest(plan.Id.ValueString(),
		[]client.EnumsimpleUncachedAttributeCriteriaSchemaUrn{client.ENUMSIMPLEUNCACHEDATTRIBUTECRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0UNCACHED_ATTRIBUTE_CRITERIASIMPLE},
		AttributeTypeSlice,
		plan.Enabled.ValueBool())
	addOptionalSimpleUncachedAttributeCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.UncachedAttributeCriteriaApi.AddUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddUncachedAttributeCriteriaRequest(
		client.AddSimpleUncachedAttributeCriteriaRequestAsAddUncachedAttributeCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.AddUncachedAttributeCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Simple Uncached Attribute Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state simpleUncachedAttributeCriteriaResourceModel
	readSimpleUncachedAttributeCriteriaResponse(ctx, addResponse.SimpleUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultSimpleUncachedAttributeCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan simpleUncachedAttributeCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.GetUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Simple Uncached Attribute Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state simpleUncachedAttributeCriteriaResourceModel
	readSimpleUncachedAttributeCriteriaResponse(ctx, readResponse.SimpleUncachedAttributeCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.UncachedAttributeCriteriaApi.UpdateUncachedAttributeCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSimpleUncachedAttributeCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.UpdateUncachedAttributeCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Simple Uncached Attribute Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSimpleUncachedAttributeCriteriaResponse(ctx, updateResponse.SimpleUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *simpleUncachedAttributeCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSimpleUncachedAttributeCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSimpleUncachedAttributeCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSimpleUncachedAttributeCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSimpleUncachedAttributeCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state simpleUncachedAttributeCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.UncachedAttributeCriteriaApi.GetUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Simple Uncached Attribute Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSimpleUncachedAttributeCriteriaResponse(ctx, readResponse.SimpleUncachedAttributeCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *simpleUncachedAttributeCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSimpleUncachedAttributeCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSimpleUncachedAttributeCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSimpleUncachedAttributeCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSimpleUncachedAttributeCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan simpleUncachedAttributeCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state simpleUncachedAttributeCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.UncachedAttributeCriteriaApi.UpdateUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSimpleUncachedAttributeCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.UncachedAttributeCriteriaApi.UpdateUncachedAttributeCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Simple Uncached Attribute Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSimpleUncachedAttributeCriteriaResponse(ctx, updateResponse.SimpleUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSimpleUncachedAttributeCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *simpleUncachedAttributeCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state simpleUncachedAttributeCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.DeleteUncachedAttributeCriteriaExecute(r.apiClient.UncachedAttributeCriteriaApi.DeleteUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Simple Uncached Attribute Criteria", err, httpResp)
		return
	}
}

func (r *simpleUncachedAttributeCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSimpleUncachedAttributeCriteria(ctx, req, resp)
}

func (r *defaultSimpleUncachedAttributeCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSimpleUncachedAttributeCriteria(ctx, req, resp)
}

func importSimpleUncachedAttributeCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
