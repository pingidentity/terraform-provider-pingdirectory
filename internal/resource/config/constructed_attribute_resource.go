package config

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &constructedAttributeResource{}
	_ resource.ResourceWithConfigure   = &constructedAttributeResource{}
	_ resource.ResourceWithImportState = &constructedAttributeResource{}
	_ resource.Resource                = &defaultConstructedAttributeResource{}
	_ resource.ResourceWithConfigure   = &defaultConstructedAttributeResource{}
	_ resource.ResourceWithImportState = &defaultConstructedAttributeResource{}
)

// Create a Constructed Attribute resource
func NewConstructedAttributeResource() resource.Resource {
	return &constructedAttributeResource{}
}

func NewDefaultConstructedAttributeResource() resource.Resource {
	return &defaultConstructedAttributeResource{}
}

// constructedAttributeResource is the resource implementation.
type constructedAttributeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultConstructedAttributeResource is the resource implementation.
type defaultConstructedAttributeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *constructedAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_constructed_attribute"
}

func (r *defaultConstructedAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_constructed_attribute"
}

// Configure adds the provider configured client to the resource.
func (r *constructedAttributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultConstructedAttributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type constructedAttributeResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	Description     types.String `tfsdk:"description"`
	AttributeType   types.String `tfsdk:"attribute_type"`
	ValuePattern    types.Set    `tfsdk:"value_pattern"`
}

// GetSchema defines the schema for the resource.
func (r *constructedAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	constructedAttributeSchema(ctx, req, resp, false)
}

func (r *defaultConstructedAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	constructedAttributeSchema(ctx, req, resp, true)
}

func constructedAttributeSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Constructed Attribute.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description: "A description for this Constructed Attribute",
				Optional:    true,
			},
			"attribute_type": schema.StringAttribute{
				Description: "Specifies the attribute type for the attribute whose values are to be constructed.",
				Required:    true,
			},
			"value_pattern": schema.SetAttribute{
				Description: "Specifies a pattern for constructing the attribute value using fixed text and attribute values from the entry.",
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
	if isDefault {
		// Add any default properties and set optional properties to computed where necessary
		SetAllAttributesToOptionalAndComputed(&schemaDef, []string{"id"})
	}
	AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add optional fields to create request for constructed-attribute constructed-attribute
func addOptionalConstructedAttributeFields(ctx context.Context, addRequest *client.AddConstructedAttributeRequest, plan constructedAttributeResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a ConstructedAttributeResponse object into the model struct
func readConstructedAttributeResponse(ctx context.Context, r *client.ConstructedAttributeResponse, state *constructedAttributeResourceModel, expectedValues *constructedAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.AttributeType = types.StringValue(r.AttributeType)
	state.ValuePattern = internaltypes.GetStringSet(r.ValuePattern)
	state.Notifications, state.RequiredActions = ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createConstructedAttributeOperations(plan constructedAttributeResourceModel, state constructedAttributeResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.AttributeType, state.AttributeType, "attribute-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ValuePattern, state.ValuePattern, "value-pattern")
	return ops
}

// Create a constructed-attribute constructed-attribute
func (r *constructedAttributeResource) CreateConstructedAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan constructedAttributeResourceModel) (*constructedAttributeResourceModel, error) {
	var ValuePatternSlice []string
	plan.ValuePattern.ElementsAs(ctx, &ValuePatternSlice, false)
	addRequest := client.NewAddConstructedAttributeRequest(plan.Id.ValueString(),
		plan.AttributeType.ValueString(),
		ValuePatternSlice)
	addOptionalConstructedAttributeFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ConstructedAttributeApi.AddConstructedAttribute(
		ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddConstructedAttributeRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.ConstructedAttributeApi.AddConstructedAttributeExecute(apiAddRequest)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Constructed Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state constructedAttributeResourceModel
	readConstructedAttributeResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *constructedAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan constructedAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateConstructedAttribute(ctx, req, resp, plan)
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
func (r *defaultConstructedAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan constructedAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConstructedAttributeApi.GetConstructedAttribute(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Constructed Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state constructedAttributeResourceModel
	readConstructedAttributeResponse(ctx, readResponse, &state, &plan, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ConstructedAttributeApi.UpdateConstructedAttribute(ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createConstructedAttributeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ConstructedAttributeApi.UpdateConstructedAttributeExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Constructed Attribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConstructedAttributeResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *constructedAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConstructedAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultConstructedAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConstructedAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readConstructedAttribute(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state constructedAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ConstructedAttributeApi.GetConstructedAttribute(
		ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Constructed Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readConstructedAttributeResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *constructedAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConstructedAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultConstructedAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConstructedAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateConstructedAttribute(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan constructedAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state constructedAttributeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ConstructedAttributeApi.UpdateConstructedAttribute(
		ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createConstructedAttributeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ConstructedAttributeApi.UpdateConstructedAttributeExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Constructed Attribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConstructedAttributeResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultConstructedAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *constructedAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state constructedAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ConstructedAttributeApi.DeleteConstructedAttributeExecute(r.apiClient.ConstructedAttributeApi.DeleteConstructedAttribute(
		ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Constructed Attribute", err, httpResp)
		return
	}
}

func (r *constructedAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConstructedAttribute(ctx, req, resp)
}

func (r *defaultConstructedAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConstructedAttribute(ctx, req, resp)
}

func importConstructedAttribute(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
