// Copyright © 2025 Ping Identity Corporation

package sensitiveattribute

import (
	"context"

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
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &sensitiveAttributeResource{}
	_ resource.ResourceWithConfigure   = &sensitiveAttributeResource{}
	_ resource.ResourceWithImportState = &sensitiveAttributeResource{}
	_ resource.Resource                = &defaultSensitiveAttributeResource{}
	_ resource.ResourceWithConfigure   = &defaultSensitiveAttributeResource{}
	_ resource.ResourceWithImportState = &defaultSensitiveAttributeResource{}
)

// Create a Sensitive Attribute resource
func NewSensitiveAttributeResource() resource.Resource {
	return &sensitiveAttributeResource{}
}

func NewDefaultSensitiveAttributeResource() resource.Resource {
	return &defaultSensitiveAttributeResource{}
}

// sensitiveAttributeResource is the resource implementation.
type sensitiveAttributeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSensitiveAttributeResource is the resource implementation.
type defaultSensitiveAttributeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *sensitiveAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sensitive_attribute"
}

func (r *defaultSensitiveAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_sensitive_attribute"
}

// Configure adds the provider configured client to the resource.
func (r *sensitiveAttributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultSensitiveAttributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type sensitiveAttributeResourceModel struct {
	Id                                           types.String `tfsdk:"id"`
	Name                                         types.String `tfsdk:"name"`
	Notifications                                types.Set    `tfsdk:"notifications"`
	RequiredActions                              types.Set    `tfsdk:"required_actions"`
	Type                                         types.String `tfsdk:"type"`
	Description                                  types.String `tfsdk:"description"`
	AttributeType                                types.Set    `tfsdk:"attribute_type"`
	IncludeDefaultSensitiveOperationalAttributes types.Bool   `tfsdk:"include_default_sensitive_operational_attributes"`
	AllowInReturnedEntries                       types.String `tfsdk:"allow_in_returned_entries"`
	AllowInFilter                                types.String `tfsdk:"allow_in_filter"`
	AllowInAdd                                   types.String `tfsdk:"allow_in_add"`
	AllowInCompare                               types.String `tfsdk:"allow_in_compare"`
	AllowInModify                                types.String `tfsdk:"allow_in_modify"`
}

// GetSchema defines the schema for the resource.
func (r *sensitiveAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	sensitiveAttributeSchema(ctx, req, resp, false)
}

func (r *defaultSensitiveAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	sensitiveAttributeSchema(ctx, req, resp, true)
}

func sensitiveAttributeSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Sensitive Attribute.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Sensitive Attribute resource. Options are ['sensitive-attribute']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("sensitive-attribute"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"sensitive-attribute"}...),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Sensitive Attribute",
				Optional:    true,
			},
			"attribute_type": schema.SetAttribute{
				Description: "The name(s) or OID(s) of the attribute types for attributes whose values may be considered sensitive.",
				Required:    true,
				ElementType: types.StringType,
			},
			"include_default_sensitive_operational_attributes": schema.BoolAttribute{
				Description: "Indicates whether to automatically include any server-generated operational attributes that may contain sensitive data.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"allow_in_returned_entries": schema.StringAttribute{
				Description: "Indicates whether sensitive attributes should be included in entries returned to the client. This includes not only search result entries, but also other forms including in the values of controls like the pre-read, post-read, get authorization entry, and LDAP join response controls.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("secure-only"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"allow", "suppress", "secure-only"}...),
				},
			},
			"allow_in_filter": schema.StringAttribute{
				Description: "Indicates whether clients will be allowed to include sensitive attributes in search filters. This also includes filters that may be used in other forms, including assertion and LDAP join request controls.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("secure-only"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"allow", "reject", "secure-only"}...),
				},
			},
			"allow_in_add": schema.StringAttribute{
				Description: "Indicates whether clients will be allowed to include sensitive attributes in add requests.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("secure-only"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"allow", "reject", "secure-only"}...),
				},
			},
			"allow_in_compare": schema.StringAttribute{
				Description: "Indicates whether clients will be allowed to target sensitive attributes with compare requests.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("secure-only"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"allow", "reject", "secure-only"}...),
				},
			},
			"allow_in_modify": schema.StringAttribute{
				Description: "Indicates whether clients will be allowed to target sensitive attributes with modify requests.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("secure-only"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"allow", "reject", "secure-only"}...),
				},
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
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add optional fields to create request for sensitive-attribute sensitive-attribute
func addOptionalSensitiveAttributeFields(ctx context.Context, addRequest *client.AddSensitiveAttributeRequest, plan sensitiveAttributeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludeDefaultSensitiveOperationalAttributes) {
		addRequest.IncludeDefaultSensitiveOperationalAttributes = plan.IncludeDefaultSensitiveOperationalAttributes.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AllowInReturnedEntries) {
		allowInReturnedEntries, err := client.NewEnumsensitiveAttributeAllowInReturnedEntriesPropFromValue(plan.AllowInReturnedEntries.ValueString())
		if err != nil {
			return err
		}
		addRequest.AllowInReturnedEntries = allowInReturnedEntries
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AllowInFilter) {
		allowInFilter, err := client.NewEnumsensitiveAttributeAllowInFilterPropFromValue(plan.AllowInFilter.ValueString())
		if err != nil {
			return err
		}
		addRequest.AllowInFilter = allowInFilter
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AllowInAdd) {
		allowInAdd, err := client.NewEnumsensitiveAttributeAllowInAddPropFromValue(plan.AllowInAdd.ValueString())
		if err != nil {
			return err
		}
		addRequest.AllowInAdd = allowInAdd
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AllowInCompare) {
		allowInCompare, err := client.NewEnumsensitiveAttributeAllowInComparePropFromValue(plan.AllowInCompare.ValueString())
		if err != nil {
			return err
		}
		addRequest.AllowInCompare = allowInCompare
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AllowInModify) {
		allowInModify, err := client.NewEnumsensitiveAttributeAllowInModifyPropFromValue(plan.AllowInModify.ValueString())
		if err != nil {
			return err
		}
		addRequest.AllowInModify = allowInModify
	}
	return nil
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *sensitiveAttributeResourceModel) populateAllComputedStringAttributes() {
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.AllowInModify.IsUnknown() || model.AllowInModify.IsNull() {
		model.AllowInModify = types.StringValue("")
	}
	if model.AllowInReturnedEntries.IsUnknown() || model.AllowInReturnedEntries.IsNull() {
		model.AllowInReturnedEntries = types.StringValue("")
	}
	if model.AllowInFilter.IsUnknown() || model.AllowInFilter.IsNull() {
		model.AllowInFilter = types.StringValue("")
	}
	if model.AllowInAdd.IsUnknown() || model.AllowInAdd.IsNull() {
		model.AllowInAdd = types.StringValue("")
	}
	if model.AllowInCompare.IsUnknown() || model.AllowInCompare.IsNull() {
		model.AllowInCompare = types.StringValue("")
	}
}

// Read a SensitiveAttributeResponse object into the model struct
func readSensitiveAttributeResponse(ctx context.Context, r *client.SensitiveAttributeResponse, state *sensitiveAttributeResourceModel, expectedValues *sensitiveAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("sensitive-attribute")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.AttributeType = internaltypes.GetStringSet(r.AttributeType)
	state.IncludeDefaultSensitiveOperationalAttributes = internaltypes.BoolTypeOrNil(r.IncludeDefaultSensitiveOperationalAttributes)
	state.AllowInReturnedEntries = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsensitiveAttributeAllowInReturnedEntriesProp(r.AllowInReturnedEntries), true)
	state.AllowInFilter = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsensitiveAttributeAllowInFilterProp(r.AllowInFilter), true)
	state.AllowInAdd = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsensitiveAttributeAllowInAddProp(r.AllowInAdd), true)
	state.AllowInCompare = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsensitiveAttributeAllowInCompareProp(r.AllowInCompare), true)
	state.AllowInModify = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsensitiveAttributeAllowInModifyProp(r.AllowInModify), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSensitiveAttributeOperations(plan sensitiveAttributeResourceModel, state sensitiveAttributeResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AttributeType, state.AttributeType, "attribute-type")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeDefaultSensitiveOperationalAttributes, state.IncludeDefaultSensitiveOperationalAttributes, "include-default-sensitive-operational-attributes")
	operations.AddStringOperationIfNecessary(&ops, plan.AllowInReturnedEntries, state.AllowInReturnedEntries, "allow-in-returned-entries")
	operations.AddStringOperationIfNecessary(&ops, plan.AllowInFilter, state.AllowInFilter, "allow-in-filter")
	operations.AddStringOperationIfNecessary(&ops, plan.AllowInAdd, state.AllowInAdd, "allow-in-add")
	operations.AddStringOperationIfNecessary(&ops, plan.AllowInCompare, state.AllowInCompare, "allow-in-compare")
	operations.AddStringOperationIfNecessary(&ops, plan.AllowInModify, state.AllowInModify, "allow-in-modify")
	return ops
}

// Create a sensitive-attribute sensitive-attribute
func (r *sensitiveAttributeResource) CreateSensitiveAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan sensitiveAttributeResourceModel) (*sensitiveAttributeResourceModel, error) {
	var AttributeTypeSlice []string
	plan.AttributeType.ElementsAs(ctx, &AttributeTypeSlice, false)
	addRequest := client.NewAddSensitiveAttributeRequest(AttributeTypeSlice,
		plan.Name.ValueString())
	err := addOptionalSensitiveAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Sensitive Attribute", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.SensitiveAttributeAPI.AddSensitiveAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddSensitiveAttributeRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.SensitiveAttributeAPI.AddSensitiveAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Sensitive Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state sensitiveAttributeResourceModel
	readSensitiveAttributeResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *sensitiveAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan sensitiveAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateSensitiveAttribute(ctx, req, resp, plan)
	if err != nil {
		return
	}

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
func (r *defaultSensitiveAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan sensitiveAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SensitiveAttributeAPI.GetSensitiveAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Sensitive Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state sensitiveAttributeResourceModel
	readSensitiveAttributeResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.SensitiveAttributeAPI.UpdateSensitiveAttribute(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createSensitiveAttributeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.SensitiveAttributeAPI.UpdateSensitiveAttributeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Sensitive Attribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSensitiveAttributeResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
	}

	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *sensitiveAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSensitiveAttribute(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultSensitiveAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSensitiveAttribute(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readSensitiveAttribute(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state sensitiveAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.SensitiveAttributeAPI.GetSensitiveAttribute(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Sensitive Attribute", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Sensitive Attribute", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSensitiveAttributeResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *sensitiveAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSensitiveAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSensitiveAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSensitiveAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSensitiveAttribute(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan sensitiveAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state sensitiveAttributeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.SensitiveAttributeAPI.UpdateSensitiveAttribute(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createSensitiveAttributeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.SensitiveAttributeAPI.UpdateSensitiveAttributeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Sensitive Attribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSensitiveAttributeResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSensitiveAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *sensitiveAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state sensitiveAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.SensitiveAttributeAPI.DeleteSensitiveAttributeExecute(r.apiClient.SensitiveAttributeAPI.DeleteSensitiveAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Sensitive Attribute", err, httpResp)
		return
	}
}

func (r *sensitiveAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSensitiveAttribute(ctx, req, resp)
}

func (r *defaultSensitiveAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSensitiveAttribute(ctx, req, resp)
}

func importSensitiveAttribute(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
