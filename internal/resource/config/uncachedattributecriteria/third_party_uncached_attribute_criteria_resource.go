package uncachedattributecriteria

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
	_ resource.Resource                = &thirdPartyUncachedAttributeCriteriaResource{}
	_ resource.ResourceWithConfigure   = &thirdPartyUncachedAttributeCriteriaResource{}
	_ resource.ResourceWithImportState = &thirdPartyUncachedAttributeCriteriaResource{}
	_ resource.Resource                = &defaultThirdPartyUncachedAttributeCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultThirdPartyUncachedAttributeCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultThirdPartyUncachedAttributeCriteriaResource{}
)

// Create a Third Party Uncached Attribute Criteria resource
func NewThirdPartyUncachedAttributeCriteriaResource() resource.Resource {
	return &thirdPartyUncachedAttributeCriteriaResource{}
}

func NewDefaultThirdPartyUncachedAttributeCriteriaResource() resource.Resource {
	return &defaultThirdPartyUncachedAttributeCriteriaResource{}
}

// thirdPartyUncachedAttributeCriteriaResource is the resource implementation.
type thirdPartyUncachedAttributeCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultThirdPartyUncachedAttributeCriteriaResource is the resource implementation.
type defaultThirdPartyUncachedAttributeCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *thirdPartyUncachedAttributeCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_third_party_uncached_attribute_criteria"
}

func (r *defaultThirdPartyUncachedAttributeCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_third_party_uncached_attribute_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *thirdPartyUncachedAttributeCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultThirdPartyUncachedAttributeCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type thirdPartyUncachedAttributeCriteriaResourceModel struct {
	Id                types.String `tfsdk:"id"`
	LastUpdated       types.String `tfsdk:"last_updated"`
	Notifications     types.Set    `tfsdk:"notifications"`
	RequiredActions   types.Set    `tfsdk:"required_actions"`
	ExtensionClass    types.String `tfsdk:"extension_class"`
	ExtensionArgument types.Set    `tfsdk:"extension_argument"`
	Description       types.String `tfsdk:"description"`
	Enabled           types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *thirdPartyUncachedAttributeCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	thirdPartyUncachedAttributeCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultThirdPartyUncachedAttributeCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	thirdPartyUncachedAttributeCriteriaSchema(ctx, req, resp, true)
}

func thirdPartyUncachedAttributeCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Third Party Uncached Attribute Criteria.",
		Attributes: map[string]schema.Attribute{
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Uncached Attribute Criteria.",
				Required:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Uncached Attribute Criteria. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
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
func addOptionalThirdPartyUncachedAttributeCriteriaFields(ctx context.Context, addRequest *client.AddThirdPartyUncachedAttributeCriteriaRequest, plan thirdPartyUncachedAttributeCriteriaResourceModel) {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a ThirdPartyUncachedAttributeCriteriaResponse object into the model struct
func readThirdPartyUncachedAttributeCriteriaResponse(ctx context.Context, r *client.ThirdPartyUncachedAttributeCriteriaResponse, state *thirdPartyUncachedAttributeCriteriaResourceModel, expectedValues *thirdPartyUncachedAttributeCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createThirdPartyUncachedAttributeCriteriaOperations(plan thirdPartyUncachedAttributeCriteriaResourceModel, state thirdPartyUncachedAttributeCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *thirdPartyUncachedAttributeCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan thirdPartyUncachedAttributeCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddThirdPartyUncachedAttributeCriteriaRequest(plan.Id.ValueString(),
		[]client.EnumthirdPartyUncachedAttributeCriteriaSchemaUrn{client.ENUMTHIRDPARTYUNCACHEDATTRIBUTECRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0UNCACHED_ATTRIBUTE_CRITERIATHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalThirdPartyUncachedAttributeCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.UncachedAttributeCriteriaApi.AddUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddUncachedAttributeCriteriaRequest(
		client.AddThirdPartyUncachedAttributeCriteriaRequestAsAddUncachedAttributeCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.AddUncachedAttributeCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Third Party Uncached Attribute Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state thirdPartyUncachedAttributeCriteriaResourceModel
	readThirdPartyUncachedAttributeCriteriaResponse(ctx, addResponse.ThirdPartyUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultThirdPartyUncachedAttributeCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan thirdPartyUncachedAttributeCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.GetUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Third Party Uncached Attribute Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state thirdPartyUncachedAttributeCriteriaResourceModel
	readThirdPartyUncachedAttributeCriteriaResponse(ctx, readResponse.ThirdPartyUncachedAttributeCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.UncachedAttributeCriteriaApi.UpdateUncachedAttributeCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createThirdPartyUncachedAttributeCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.UpdateUncachedAttributeCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Third Party Uncached Attribute Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readThirdPartyUncachedAttributeCriteriaResponse(ctx, updateResponse.ThirdPartyUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *thirdPartyUncachedAttributeCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readThirdPartyUncachedAttributeCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultThirdPartyUncachedAttributeCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readThirdPartyUncachedAttributeCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readThirdPartyUncachedAttributeCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state thirdPartyUncachedAttributeCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.UncachedAttributeCriteriaApi.GetUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Third Party Uncached Attribute Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readThirdPartyUncachedAttributeCriteriaResponse(ctx, readResponse.ThirdPartyUncachedAttributeCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *thirdPartyUncachedAttributeCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateThirdPartyUncachedAttributeCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultThirdPartyUncachedAttributeCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateThirdPartyUncachedAttributeCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateThirdPartyUncachedAttributeCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan thirdPartyUncachedAttributeCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state thirdPartyUncachedAttributeCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.UncachedAttributeCriteriaApi.UpdateUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createThirdPartyUncachedAttributeCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.UncachedAttributeCriteriaApi.UpdateUncachedAttributeCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Third Party Uncached Attribute Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readThirdPartyUncachedAttributeCriteriaResponse(ctx, updateResponse.ThirdPartyUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultThirdPartyUncachedAttributeCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *thirdPartyUncachedAttributeCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state thirdPartyUncachedAttributeCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.DeleteUncachedAttributeCriteriaExecute(r.apiClient.UncachedAttributeCriteriaApi.DeleteUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Third Party Uncached Attribute Criteria", err, httpResp)
		return
	}
}

func (r *thirdPartyUncachedAttributeCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importThirdPartyUncachedAttributeCriteria(ctx, req, resp)
}

func (r *defaultThirdPartyUncachedAttributeCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importThirdPartyUncachedAttributeCriteria(ctx, req, resp)
}

func importThirdPartyUncachedAttributeCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
