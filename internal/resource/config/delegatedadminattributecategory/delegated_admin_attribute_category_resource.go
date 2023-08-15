package delegatedadminattributecategory

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &delegatedAdminAttributeCategoryResource{}
	_ resource.ResourceWithConfigure   = &delegatedAdminAttributeCategoryResource{}
	_ resource.ResourceWithImportState = &delegatedAdminAttributeCategoryResource{}
	_ resource.Resource                = &defaultDelegatedAdminAttributeCategoryResource{}
	_ resource.ResourceWithConfigure   = &defaultDelegatedAdminAttributeCategoryResource{}
	_ resource.ResourceWithImportState = &defaultDelegatedAdminAttributeCategoryResource{}
)

// Create a Delegated Admin Attribute Category resource
func NewDelegatedAdminAttributeCategoryResource() resource.Resource {
	return &delegatedAdminAttributeCategoryResource{}
}

func NewDefaultDelegatedAdminAttributeCategoryResource() resource.Resource {
	return &defaultDelegatedAdminAttributeCategoryResource{}
}

// delegatedAdminAttributeCategoryResource is the resource implementation.
type delegatedAdminAttributeCategoryResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultDelegatedAdminAttributeCategoryResource is the resource implementation.
type defaultDelegatedAdminAttributeCategoryResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *delegatedAdminAttributeCategoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delegated_admin_attribute_category"
}

func (r *defaultDelegatedAdminAttributeCategoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_delegated_admin_attribute_category"
}

// Configure adds the provider configured client to the resource.
func (r *delegatedAdminAttributeCategoryResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultDelegatedAdminAttributeCategoryResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type delegatedAdminAttributeCategoryResourceModel struct {
	Id                types.String `tfsdk:"id"`
	LastUpdated       types.String `tfsdk:"last_updated"`
	Notifications     types.Set    `tfsdk:"notifications"`
	RequiredActions   types.Set    `tfsdk:"required_actions"`
	Type              types.String `tfsdk:"type"`
	Description       types.String `tfsdk:"description"`
	DisplayName       types.String `tfsdk:"display_name"`
	DisplayOrderIndex types.Int64  `tfsdk:"display_order_index"`
}

// GetSchema defines the schema for the resource.
func (r *delegatedAdminAttributeCategoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	delegatedAdminAttributeCategorySchema(ctx, req, resp, false)
}

func (r *defaultDelegatedAdminAttributeCategoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	delegatedAdminAttributeCategorySchema(ctx, req, resp, true)
}

func delegatedAdminAttributeCategorySchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Delegated Admin Attribute Category.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Delegated Admin Attribute Category resource. Options are ['delegated-admin-attribute-category']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("delegated-admin-attribute-category"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"delegated-admin-attribute-category"}...),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Delegated Admin Attribute Category",
				Optional:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "A human readable display name for this Delegated Admin Attribute Category.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_order_index": schema.Int64Attribute{
				Description: "Delegated Admin Attribute Categories are ordered for display based on this index from least to greatest.",
				Required:    true,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type", "display_name"})
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Add optional fields to create request for delegated-admin-attribute-category delegated-admin-attribute-category
func addOptionalDelegatedAdminAttributeCategoryFields(ctx context.Context, addRequest *client.AddDelegatedAdminAttributeCategoryRequest, plan delegatedAdminAttributeCategoryResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a DelegatedAdminAttributeCategoryResponse object into the model struct
func readDelegatedAdminAttributeCategoryResponse(ctx context.Context, r *client.DelegatedAdminAttributeCategoryResponse, state *delegatedAdminAttributeCategoryResourceModel, expectedValues *delegatedAdminAttributeCategoryResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("delegated-admin-attribute-category")
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.DisplayName = types.StringValue(r.DisplayName)
	state.DisplayOrderIndex = types.Int64Value(r.DisplayOrderIndex)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createDelegatedAdminAttributeCategoryOperations(plan delegatedAdminAttributeCategoryResourceModel, state delegatedAdminAttributeCategoryResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.DisplayName, state.DisplayName, "display-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.DisplayOrderIndex, state.DisplayOrderIndex, "display-order-index")
	return ops
}

// Create a delegated-admin-attribute-category delegated-admin-attribute-category
func (r *delegatedAdminAttributeCategoryResource) CreateDelegatedAdminAttributeCategory(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan delegatedAdminAttributeCategoryResourceModel) (*delegatedAdminAttributeCategoryResourceModel, error) {
	addRequest := client.NewAddDelegatedAdminAttributeCategoryRequest(plan.DisplayName.ValueString(),
		plan.DisplayOrderIndex.ValueInt64())
	addOptionalDelegatedAdminAttributeCategoryFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DelegatedAdminAttributeCategoryApi.AddDelegatedAdminAttributeCategory(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddDelegatedAdminAttributeCategoryRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.DelegatedAdminAttributeCategoryApi.AddDelegatedAdminAttributeCategoryExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Delegated Admin Attribute Category", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state delegatedAdminAttributeCategoryResourceModel
	readDelegatedAdminAttributeCategoryResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *delegatedAdminAttributeCategoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan delegatedAdminAttributeCategoryResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateDelegatedAdminAttributeCategory(ctx, req, resp, plan)
	if err != nil {
		return
	}

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, *state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultDelegatedAdminAttributeCategoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan delegatedAdminAttributeCategoryResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DelegatedAdminAttributeCategoryApi.GetDelegatedAdminAttributeCategory(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.DisplayName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Attribute Category", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state delegatedAdminAttributeCategoryResourceModel
	readDelegatedAdminAttributeCategoryResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.DelegatedAdminAttributeCategoryApi.UpdateDelegatedAdminAttributeCategory(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.DisplayName.ValueString())
	ops := createDelegatedAdminAttributeCategoryOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.DelegatedAdminAttributeCategoryApi.UpdateDelegatedAdminAttributeCategoryExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Delegated Admin Attribute Category", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDelegatedAdminAttributeCategoryResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *delegatedAdminAttributeCategoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDelegatedAdminAttributeCategory(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultDelegatedAdminAttributeCategoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDelegatedAdminAttributeCategory(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readDelegatedAdminAttributeCategory(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state delegatedAdminAttributeCategoryResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.DelegatedAdminAttributeCategoryApi.GetDelegatedAdminAttributeCategory(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.DisplayName.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Attribute Category", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Attribute Category", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDelegatedAdminAttributeCategoryResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *delegatedAdminAttributeCategoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDelegatedAdminAttributeCategory(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDelegatedAdminAttributeCategoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDelegatedAdminAttributeCategory(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateDelegatedAdminAttributeCategory(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan delegatedAdminAttributeCategoryResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state delegatedAdminAttributeCategoryResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.DelegatedAdminAttributeCategoryApi.UpdateDelegatedAdminAttributeCategory(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.DisplayName.ValueString())

	// Determine what update operations are necessary
	ops := createDelegatedAdminAttributeCategoryOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.DelegatedAdminAttributeCategoryApi.UpdateDelegatedAdminAttributeCategoryExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Delegated Admin Attribute Category", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDelegatedAdminAttributeCategoryResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultDelegatedAdminAttributeCategoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *delegatedAdminAttributeCategoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state delegatedAdminAttributeCategoryResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DelegatedAdminAttributeCategoryApi.DeleteDelegatedAdminAttributeCategoryExecute(r.apiClient.DelegatedAdminAttributeCategoryApi.DeleteDelegatedAdminAttributeCategory(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.DisplayName.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Delegated Admin Attribute Category", err, httpResp)
		return
	}
}

func (r *delegatedAdminAttributeCategoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDelegatedAdminAttributeCategory(ctx, req, resp)
}

func (r *defaultDelegatedAdminAttributeCategoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDelegatedAdminAttributeCategory(ctx, req, resp)
}

func importDelegatedAdminAttributeCategory(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to display_name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("display_name"), req, resp)
}
