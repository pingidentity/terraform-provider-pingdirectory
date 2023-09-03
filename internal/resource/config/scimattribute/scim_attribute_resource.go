package scimattribute

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
	_ resource.Resource                = &scimAttributeResource{}
	_ resource.ResourceWithConfigure   = &scimAttributeResource{}
	_ resource.ResourceWithImportState = &scimAttributeResource{}
	_ resource.Resource                = &defaultScimAttributeResource{}
	_ resource.ResourceWithConfigure   = &defaultScimAttributeResource{}
	_ resource.ResourceWithImportState = &defaultScimAttributeResource{}
)

// Create a Scim Attribute resource
func NewScimAttributeResource() resource.Resource {
	return &scimAttributeResource{}
}

func NewDefaultScimAttributeResource() resource.Resource {
	return &defaultScimAttributeResource{}
}

// scimAttributeResource is the resource implementation.
type scimAttributeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultScimAttributeResource is the resource implementation.
type defaultScimAttributeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *scimAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scim_attribute"
}

func (r *defaultScimAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_scim_attribute"
}

// Configure adds the provider configured client to the resource.
func (r *scimAttributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultScimAttributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type scimAttributeResourceModel struct {
	Id              types.String `tfsdk:"id"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	ResourceType    types.String `tfsdk:"resource_type"`
	ScimSchemaName  types.String `tfsdk:"scim_schema_name"`
	Description     types.String `tfsdk:"description"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	Required        types.Bool   `tfsdk:"required"`
	CaseExact       types.Bool   `tfsdk:"case_exact"`
	MultiValued     types.Bool   `tfsdk:"multi_valued"`
	CanonicalValue  types.Set    `tfsdk:"canonical_value"`
	Mutability      types.String `tfsdk:"mutability"`
	Returned        types.String `tfsdk:"returned"`
	ReferenceType   types.Set    `tfsdk:"reference_type"`
}

// GetSchema defines the schema for the resource.
func (r *scimAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	scimAttributeSchema(ctx, req, resp, false)
}

func (r *defaultScimAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	scimAttributeSchema(ctx, req, resp, true)
}

func scimAttributeSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Scim Attribute.",
		Attributes: map[string]schema.Attribute{
			"resource_type": schema.StringAttribute{
				Description: "The type of SCIM Attribute resource. Options are ['scim-attribute']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("scim-attribute"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"scim-attribute"}...),
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
				Description: "A description for this SCIM Attribute",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the attribute.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "Specifies the data type for this attribute.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("string"),
			},
			"required": schema.BoolAttribute{
				Description: "Specifies whether this attribute is required.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"case_exact": schema.BoolAttribute{
				Description: "Specifies whether the attribute values are case sensitive.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"multi_valued": schema.BoolAttribute{
				Description: "Specifies whether this attribute may have multiple values.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"canonical_value": schema.SetAttribute{
				Description: "Specifies the suggested canonical type values for the attribute.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"mutability": schema.StringAttribute{
				Description: "Specifies the circumstances under which the values of the attribute can be written.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("read-write"),
			},
			"returned": schema.StringAttribute{
				Description: "Specifies the circumstances under which the values of the attribute are returned in response to a request.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("by-default"),
			},
			"reference_type": schema.SetAttribute{
				Description: "Specifies the SCIM resource types that may be referenced. This property is only applicable for attributes that are of type 'reference'. Valid values are: A SCIM resource type (e.g., 'User' or 'Group'), 'external' - indicating the resource is an external resource (e.g., such as a photo), or 'uri' - indicating that the reference is to a service endpoint or an identifier (such as a schema urn).",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["resource_type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		}
		schemaDef.Attributes["resource_type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"resource_type", "name", "scim_schema_name"})
	} else {
		// Add RequiresReplace modifier for read-only attributes
		nameAttr := schemaDef.Attributes["name"].(schema.StringAttribute)
		nameAttr.PlanModifiers = append(nameAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["name"] = nameAttr
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Add optional fields to create request for scim-attribute scim-attribute
func addOptionalScimAttributeFields(ctx context.Context, addRequest *client.AddScimAttributeRequest, plan scimAttributeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Type) {
		typeValue, err := client.NewEnumscimAttributeTypePropFromValue(plan.Type.ValueString())
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
		mutability, err := client.NewEnumscimAttributeMutabilityPropFromValue(plan.Mutability.ValueString())
		if err != nil {
			return err
		}
		addRequest.Mutability = mutability
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Returned) {
		returned, err := client.NewEnumscimAttributeReturnedPropFromValue(plan.Returned.ValueString())
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

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *scimAttributeResourceModel) populateAllComputedStringAttributes() {
	if model.Returned.IsUnknown() || model.Returned.IsNull() {
		model.Returned = types.StringValue("")
	}
	if model.Type.IsUnknown() || model.Type.IsNull() {
		model.Type = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.Mutability.IsUnknown() || model.Mutability.IsNull() {
		model.Mutability = types.StringValue("")
	}
	if model.Name.IsUnknown() || model.Name.IsNull() {
		model.Name = types.StringValue("")
	}
}

// Read a ScimAttributeResponse object into the model struct
func readScimAttributeResponse(ctx context.Context, r *client.ScimAttributeResponse, state *scimAttributeResourceModel, expectedValues *scimAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("scim-attribute")
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Name = types.StringValue(r.Name)
	state.Type = types.StringValue(r.Type.String())
	state.Required = types.BoolValue(r.Required)
	state.CaseExact = types.BoolValue(r.CaseExact)
	state.MultiValued = types.BoolValue(r.MultiValued)
	state.CanonicalValue = internaltypes.GetStringSet(r.CanonicalValue)
	state.Mutability = types.StringValue(r.Mutability.String())
	state.Returned = types.StringValue(r.Returned.String())
	state.ReferenceType = internaltypes.GetStringSet(r.ReferenceType)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *scimAttributeResourceModel) setStateValuesNotReturnedByAPI(expectedValues *scimAttributeResourceModel) {
	if !expectedValues.ScimSchemaName.IsUnknown() {
		state.ScimSchemaName = expectedValues.ScimSchemaName
	}
}

// Create any update operations necessary to make the state match the plan
func createScimAttributeOperations(plan scimAttributeResourceModel, state scimAttributeResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.Name, state.Name, "name")
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

// Create a scim-attribute scim-attribute
func (r *scimAttributeResource) CreateScimAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan scimAttributeResourceModel) (*scimAttributeResourceModel, error) {
	addRequest := client.NewAddScimAttributeRequest(plan.Name.ValueString(),
		plan.Name.ValueString())
	err := addOptionalScimAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Scim Attribute", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ScimAttributeApi.AddScimAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.ScimSchemaName.ValueString())
	apiAddRequest = apiAddRequest.AddScimAttributeRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.ScimAttributeApi.AddScimAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Scim Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state scimAttributeResourceModel
	readScimAttributeResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *scimAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan scimAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateScimAttribute(ctx, req, resp, plan)
	if err != nil {
		return
	}

	// Populate Computed attribute values
	state.setStateValuesNotReturnedByAPI(&plan)
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
func (r *defaultScimAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan scimAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ScimAttributeApi.GetScimAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.ScimSchemaName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Scim Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state scimAttributeResourceModel
	readScimAttributeResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ScimAttributeApi.UpdateScimAttribute(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.ScimSchemaName.ValueString())
	ops := createScimAttributeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ScimAttributeApi.UpdateScimAttributeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Scim Attribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readScimAttributeResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *scimAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readScimAttribute(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultScimAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readScimAttribute(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readScimAttribute(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state scimAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ScimAttributeApi.GetScimAttribute(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString(), state.ScimSchemaName.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Scim Attribute", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Scim Attribute", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readScimAttributeResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *scimAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateScimAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultScimAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateScimAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateScimAttribute(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan scimAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state scimAttributeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ScimAttributeApi.UpdateScimAttribute(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString(), plan.ScimSchemaName.ValueString())

	// Determine what update operations are necessary
	ops := createScimAttributeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ScimAttributeApi.UpdateScimAttributeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Scim Attribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readScimAttributeResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// This config object is edit-only, so Terraform can't delete it.
// After running a delete, Terraform will just "forget" about this object and it can be managed elsewhere.
func (r *defaultScimAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *scimAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state scimAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ScimAttributeApi.DeleteScimAttributeExecute(r.apiClient.ScimAttributeApi.DeleteScimAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.ScimSchemaName.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Scim Attribute", err, httpResp)
		return
	}
}

func (r *scimAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importScimAttribute(ctx, req, resp)
}

func (r *defaultScimAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importScimAttribute(ctx, req, resp)
}

func importScimAttribute(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [scim-schema-name]/[scim-attribute-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scim_schema_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), split[1])...)
}
