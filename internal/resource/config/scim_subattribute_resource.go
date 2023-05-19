package config

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &scimSubattributeResource{}
	_ resource.ResourceWithConfigure   = &scimSubattributeResource{}
	_ resource.ResourceWithImportState = &scimSubattributeResource{}
	_ resource.Resource                = &defaultScimSubattributeResource{}
	_ resource.ResourceWithConfigure   = &defaultScimSubattributeResource{}
	_ resource.ResourceWithImportState = &defaultScimSubattributeResource{}
)

// Create a Scim Subattribute resource
func NewScimSubattributeResource() resource.Resource {
	return &scimSubattributeResource{}
}

func NewDefaultScimSubattributeResource() resource.Resource {
	return &defaultScimSubattributeResource{}
}

// scimSubattributeResource is the resource implementation.
type scimSubattributeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultScimSubattributeResource is the resource implementation.
type defaultScimSubattributeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *scimSubattributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scim_subattribute"
}

func (r *defaultScimSubattributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_scim_subattribute"
}

// Configure adds the provider configured client to the resource.
func (r *scimSubattributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultScimSubattributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type scimSubattributeResourceModel struct {
	Id                types.String `tfsdk:"id"`
	LastUpdated       types.String `tfsdk:"last_updated"`
	Notifications     types.Set    `tfsdk:"notifications"`
	RequiredActions   types.Set    `tfsdk:"required_actions"`
	ScimAttributeName types.String `tfsdk:"scim_attribute_name"`
	ScimSchemaName    types.String `tfsdk:"scim_schema_name"`
	Description       types.String `tfsdk:"description"`
	Type              types.String `tfsdk:"type"`
	Required          types.Bool   `tfsdk:"required"`
	CaseExact         types.Bool   `tfsdk:"case_exact"`
	MultiValued       types.Bool   `tfsdk:"multi_valued"`
	CanonicalValue    types.Set    `tfsdk:"canonical_value"`
	Mutability        types.String `tfsdk:"mutability"`
	Returned          types.String `tfsdk:"returned"`
	ReferenceType     types.Set    `tfsdk:"reference_type"`
}

// GetSchema defines the schema for the resource.
func (r *scimSubattributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	scimSubattributeSchema(ctx, req, resp, false)
}

func (r *defaultScimSubattributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	scimSubattributeSchema(ctx, req, resp, true)
}

func scimSubattributeSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Scim Subattribute.",
		Attributes: map[string]schema.Attribute{
			"scim_attribute_name": schema.StringAttribute{
				Description: "Name of the parent SCIM Attribute",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"scim_schema_name": schema.StringAttribute{
				Description: "Name of the parent SCIM Schema",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this SCIM Subattribute",
				Optional:    true,
			},
			"type": schema.StringAttribute{
				Description: "Specifies the data type for this sub-attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"required": schema.BoolAttribute{
				Description: "Specifies whether this sub-attribute is required.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"case_exact": schema.BoolAttribute{
				Description: "Specifies whether the sub-attribute values are case sensitive.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"multi_valued": schema.BoolAttribute{
				Description: "Specifies whether this attribute may have multiple values.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"canonical_value": schema.SetAttribute{
				Description: "Specifies the suggested canonical type values for the sub-attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"mutability": schema.StringAttribute{
				Description: "Specifies the circumstances under which the values of the sub-attribute can be written.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"returned": schema.StringAttribute{
				Description: "Specifies the circumstances under which the values of the sub-attribute are returned in response to a request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"reference_type": schema.SetAttribute{
				Description: "Specifies the SCIM resource types that may be referenced. This property is only applicable for sub-attributes that are of type 'reference'. Valid values are: A SCIM resource type (e.g., 'User' or 'Group'), 'external' - indicating the resource is an external resource (e.g., such as a photo), or 'uri' - indicating that the reference is to a service endpoint or an identifier (such as a schema urn).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
		},
	}
	if setOptionalToComputed {
		SetAllAttributesToOptionalAndComputed(&schema, []string{"id", "scim_attribute_name", "scim_schema_name"})
	}
	AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalScimSubattributeFields(ctx context.Context, addRequest *client.AddScimSubattributeRequest, plan scimSubattributeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Type) {
		typeValue, err := client.NewEnumscimSubattributeTypePropFromValue(plan.Type.ValueString())
		if err != nil {
			return err
		}
		addRequest.Type = typeValue
	}
	if internaltypes.IsDefined(plan.Required) {
		addRequest.Required = plan.Required.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.CaseExact) {
		addRequest.CaseExact = plan.CaseExact.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MultiValued) {
		addRequest.MultiValued = plan.MultiValued.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.CanonicalValue) {
		var slice []string
		plan.CanonicalValue.ElementsAs(ctx, &slice, false)
		addRequest.CanonicalValue = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Mutability) {
		mutability, err := client.NewEnumscimSubattributeMutabilityPropFromValue(plan.Mutability.ValueString())
		if err != nil {
			return err
		}
		addRequest.Mutability = mutability
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Returned) {
		returned, err := client.NewEnumscimSubattributeReturnedPropFromValue(plan.Returned.ValueString())
		if err != nil {
			return err
		}
		addRequest.Returned = returned
	}
	if internaltypes.IsDefined(plan.ReferenceType) {
		var slice []string
		plan.ReferenceType.ElementsAs(ctx, &slice, false)
		addRequest.ReferenceType = slice
	}
	return nil
}

// Read a ScimSubattributeResponse object into the model struct
func readScimSubattributeResponse(ctx context.Context, r *client.ScimSubattributeResponse, state *scimSubattributeResourceModel, expectedValues *scimSubattributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ScimAttributeName = expectedValues.ScimAttributeName
	state.ScimSchemaName = expectedValues.ScimSchemaName
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Type = types.StringValue(r.Type.String())
	state.Required = types.BoolValue(r.Required)
	state.CaseExact = types.BoolValue(r.CaseExact)
	state.MultiValued = types.BoolValue(r.MultiValued)
	state.CanonicalValue = internaltypes.GetStringSet(r.CanonicalValue)
	state.Mutability = types.StringValue(r.Mutability.String())
	state.Returned = types.StringValue(r.Returned.String())
	state.ReferenceType = internaltypes.GetStringSet(r.ReferenceType)
	state.Notifications, state.RequiredActions = ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createScimSubattributeOperations(plan scimSubattributeResourceModel, state scimSubattributeResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.Type, state.Type, "type")
	operations.AddBoolOperationIfNecessary(&ops, plan.Required, state.Required, "required")
	operations.AddBoolOperationIfNecessary(&ops, plan.CaseExact, state.CaseExact, "case-exact")
	operations.AddBoolOperationIfNecessary(&ops, plan.MultiValued, state.MultiValued, "multi-valued")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.CanonicalValue, state.CanonicalValue, "canonical-value")
	operations.AddStringOperationIfNecessary(&ops, plan.Mutability, state.Mutability, "mutability")
	operations.AddStringOperationIfNecessary(&ops, plan.Returned, state.Returned, "returned")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ReferenceType, state.ReferenceType, "reference-type")
	return ops
}

// Create a new resource
func (r *scimSubattributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan scimSubattributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddScimSubattributeRequest(plan.Id.ValueString())
	err := addOptionalScimSubattributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Scim Subattribute", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ScimSubattributeApi.AddScimSubattribute(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.ScimAttributeName.ValueString(), plan.ScimSchemaName.ValueString())
	apiAddRequest = apiAddRequest.AddScimSubattributeRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.ScimSubattributeApi.AddScimSubattributeExecute(apiAddRequest)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Scim Subattribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state scimSubattributeResourceModel
	readScimSubattributeResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultScimSubattributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan scimSubattributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ScimSubattributeApi.GetScimSubattribute(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.ScimAttributeName.ValueString(), plan.ScimSchemaName.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Scim Subattribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state scimSubattributeResourceModel
	readScimSubattributeResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ScimSubattributeApi.UpdateScimSubattribute(ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.ScimAttributeName.ValueString(), plan.ScimSchemaName.ValueString())
	ops := createScimSubattributeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ScimSubattributeApi.UpdateScimSubattributeExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Scim Subattribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readScimSubattributeResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *scimSubattributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readScimSubattribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultScimSubattributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readScimSubattribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readScimSubattribute(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state scimSubattributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ScimSubattributeApi.GetScimSubattribute(
		ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString(), state.ScimAttributeName.ValueString(), state.ScimSchemaName.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Scim Subattribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readScimSubattributeResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *scimSubattributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateScimSubattribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultScimSubattributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateScimSubattribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateScimSubattribute(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan scimSubattributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state scimSubattributeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ScimSubattributeApi.UpdateScimSubattribute(
		ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString(), plan.ScimAttributeName.ValueString(), plan.ScimSchemaName.ValueString())

	// Determine what update operations are necessary
	ops := createScimSubattributeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ScimSubattributeApi.UpdateScimSubattributeExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Scim Subattribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readScimSubattributeResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultScimSubattributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *scimSubattributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state scimSubattributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ScimSubattributeApi.DeleteScimSubattributeExecute(r.apiClient.ScimSubattributeApi.DeleteScimSubattribute(
		ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString(), state.ScimAttributeName.ValueString(), state.ScimSchemaName.ValueString()))
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Scim Subattribute", err, httpResp)
		return
	}
}

func (r *scimSubattributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importScimSubattribute(ctx, req, resp)
}

func (r *defaultScimSubattributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importScimSubattribute(ctx, req, resp)
}

func importScimSubattribute(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 3 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [scim-schema-name]/[scim-attribute-name]/[scim-subattribute-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scim_schema_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scim_attribute_name"), split[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), split[2])...)
}
