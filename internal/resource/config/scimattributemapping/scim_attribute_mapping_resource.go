package scimattributemapping

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
	_ resource.Resource                = &scimAttributeMappingResource{}
	_ resource.ResourceWithConfigure   = &scimAttributeMappingResource{}
	_ resource.ResourceWithImportState = &scimAttributeMappingResource{}
	_ resource.Resource                = &defaultScimAttributeMappingResource{}
	_ resource.ResourceWithConfigure   = &defaultScimAttributeMappingResource{}
	_ resource.ResourceWithImportState = &defaultScimAttributeMappingResource{}
)

// Create a Scim Attribute Mapping resource
func NewScimAttributeMappingResource() resource.Resource {
	return &scimAttributeMappingResource{}
}

func NewDefaultScimAttributeMappingResource() resource.Resource {
	return &defaultScimAttributeMappingResource{}
}

// scimAttributeMappingResource is the resource implementation.
type scimAttributeMappingResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultScimAttributeMappingResource is the resource implementation.
type defaultScimAttributeMappingResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *scimAttributeMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scim_attribute_mapping"
}

func (r *defaultScimAttributeMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_scim_attribute_mapping"
}

// Configure adds the provider configured client to the resource.
func (r *scimAttributeMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultScimAttributeMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type scimAttributeMappingResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Notifications             types.Set    `tfsdk:"notifications"`
	RequiredActions           types.Set    `tfsdk:"required_actions"`
	Type                      types.String `tfsdk:"type"`
	ScimResourceTypeName      types.String `tfsdk:"scim_resource_type_name"`
	CorrelatedLDAPDataView    types.String `tfsdk:"correlated_ldap_data_view"`
	ScimResourceTypeAttribute types.String `tfsdk:"scim_resource_type_attribute"`
	LdapAttribute             types.String `tfsdk:"ldap_attribute"`
	Readable                  types.Bool   `tfsdk:"readable"`
	Writable                  types.Bool   `tfsdk:"writable"`
	Searchable                types.Bool   `tfsdk:"searchable"`
	Authoritative             types.Bool   `tfsdk:"authoritative"`
}

// GetSchema defines the schema for the resource.
func (r *scimAttributeMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	scimAttributeMappingSchema(ctx, req, resp, false)
}

func (r *defaultScimAttributeMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	scimAttributeMappingSchema(ctx, req, resp, true)
}

func scimAttributeMappingSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Scim Attribute Mapping.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of SCIM Attribute Mapping resource. Options are ['scim-attribute-mapping']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("scim-attribute-mapping"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"scim-attribute-mapping"}...),
				},
			},
			"scim_resource_type_name": schema.StringAttribute{
				Description: "Name of the parent SCIM Resource Type",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"correlated_ldap_data_view": schema.StringAttribute{
				Description: "The Correlated LDAP Data View that persists the mapped SCIM Resource Type attribute(s).",
				Optional:    true,
			},
			"scim_resource_type_attribute": schema.StringAttribute{
				Description: "The attribute path of SCIM Resource Type attributes to be mapped.",
				Required:    true,
			},
			"ldap_attribute": schema.StringAttribute{
				Description: "The LDAP attribute to be mapped, or the path to a specific field of an LDAP attribute with the JSON object attribute syntax.",
				Required:    true,
			},
			"readable": schema.BoolAttribute{
				Description: "Specifies whether the mapping is used to map from LDAP attribute to SCIM Resource Type attribute in a read operation.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"writable": schema.BoolAttribute{
				Description: "Specifies that the mapping is used to map from SCIM Resource Type attribute to LDAP attribute in a write operation.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"searchable": schema.BoolAttribute{
				Description: "Specifies that the mapping is used to map from SCIM Resource Type attribute to LDAP attribute in a search filter.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"authoritative": schema.BoolAttribute{
				Description: "Specifies that the mapping is authoritative over other mappings for the same SCIM Resource Type attribute (for read operations).",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
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
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type", "scim_resource_type_name"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add optional fields to create request for scim-attribute-mapping scim-attribute-mapping
func addOptionalScimAttributeMappingFields(ctx context.Context, addRequest *client.AddScimAttributeMappingRequest, plan scimAttributeMappingResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CorrelatedLDAPDataView) {
		addRequest.CorrelatedLDAPDataView = plan.CorrelatedLDAPDataView.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Readable) {
		addRequest.Readable = plan.Readable.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Writable) {
		addRequest.Writable = plan.Writable.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Searchable) {
		addRequest.Searchable = plan.Searchable.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Authoritative) {
		addRequest.Authoritative = plan.Authoritative.ValueBoolPointer()
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *scimAttributeMappingResourceModel) populateAllComputedStringAttributes() {
	if model.CorrelatedLDAPDataView.IsUnknown() || model.CorrelatedLDAPDataView.IsNull() {
		model.CorrelatedLDAPDataView = types.StringValue("")
	}
	if model.ScimResourceTypeAttribute.IsUnknown() || model.ScimResourceTypeAttribute.IsNull() {
		model.ScimResourceTypeAttribute = types.StringValue("")
	}
	if model.LdapAttribute.IsUnknown() || model.LdapAttribute.IsNull() {
		model.LdapAttribute = types.StringValue("")
	}
}

// Read a ScimAttributeMappingResponse object into the model struct
func readScimAttributeMappingResponse(ctx context.Context, r *client.ScimAttributeMappingResponse, state *scimAttributeMappingResourceModel, expectedValues *scimAttributeMappingResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("scim-attribute-mapping")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.CorrelatedLDAPDataView = internaltypes.StringTypeOrNil(r.CorrelatedLDAPDataView, internaltypes.IsEmptyString(expectedValues.CorrelatedLDAPDataView))
	state.ScimResourceTypeAttribute = types.StringValue(r.ScimResourceTypeAttribute)
	state.LdapAttribute = types.StringValue(r.LdapAttribute)
	state.Readable = internaltypes.BoolTypeOrNil(r.Readable)
	state.Writable = internaltypes.BoolTypeOrNil(r.Writable)
	state.Searchable = internaltypes.BoolTypeOrNil(r.Searchable)
	state.Authoritative = internaltypes.BoolTypeOrNil(r.Authoritative)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *scimAttributeMappingResourceModel) setStateValuesNotReturnedByAPI(expectedValues *scimAttributeMappingResourceModel) {
	if !expectedValues.ScimResourceTypeName.IsUnknown() {
		state.ScimResourceTypeName = expectedValues.ScimResourceTypeName
	}
}

// Create any update operations necessary to make the state match the plan
func createScimAttributeMappingOperations(plan scimAttributeMappingResourceModel, state scimAttributeMappingResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.CorrelatedLDAPDataView, state.CorrelatedLDAPDataView, "correlated-ldap-data-view")
	operations.AddStringOperationIfNecessary(&ops, plan.ScimResourceTypeAttribute, state.ScimResourceTypeAttribute, "scim-resource-type-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.LdapAttribute, state.LdapAttribute, "ldap-attribute")
	operations.AddBoolOperationIfNecessary(&ops, plan.Readable, state.Readable, "readable")
	operations.AddBoolOperationIfNecessary(&ops, plan.Writable, state.Writable, "writable")
	operations.AddBoolOperationIfNecessary(&ops, plan.Searchable, state.Searchable, "searchable")
	operations.AddBoolOperationIfNecessary(&ops, plan.Authoritative, state.Authoritative, "authoritative")
	return ops
}

// Create a scim-attribute-mapping scim-attribute-mapping
func (r *scimAttributeMappingResource) CreateScimAttributeMapping(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan scimAttributeMappingResourceModel) (*scimAttributeMappingResourceModel, error) {
	addRequest := client.NewAddScimAttributeMappingRequest(plan.ScimResourceTypeAttribute.ValueString(),
		plan.LdapAttribute.ValueString(),
		plan.Name.ValueString())
	addOptionalScimAttributeMappingFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ScimAttributeMappingAPI.AddScimAttributeMapping(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.ScimResourceTypeName.ValueString())
	apiAddRequest = apiAddRequest.AddScimAttributeMappingRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.ScimAttributeMappingAPI.AddScimAttributeMappingExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Scim Attribute Mapping", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state scimAttributeMappingResourceModel
	readScimAttributeMappingResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *scimAttributeMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan scimAttributeMappingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateScimAttributeMapping(ctx, req, resp, plan)
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
func (r *defaultScimAttributeMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan scimAttributeMappingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ScimAttributeMappingAPI.GetScimAttributeMapping(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.ScimResourceTypeName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Scim Attribute Mapping", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state scimAttributeMappingResourceModel
	readScimAttributeMappingResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ScimAttributeMappingAPI.UpdateScimAttributeMapping(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.ScimResourceTypeName.ValueString())
	ops := createScimAttributeMappingOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ScimAttributeMappingAPI.UpdateScimAttributeMappingExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Scim Attribute Mapping", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readScimAttributeMappingResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *scimAttributeMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readScimAttributeMapping(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultScimAttributeMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readScimAttributeMapping(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readScimAttributeMapping(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state scimAttributeMappingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ScimAttributeMappingAPI.GetScimAttributeMapping(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString(), state.ScimResourceTypeName.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Scim Attribute Mapping", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Scim Attribute Mapping", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readScimAttributeMappingResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *scimAttributeMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateScimAttributeMapping(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultScimAttributeMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateScimAttributeMapping(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateScimAttributeMapping(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan scimAttributeMappingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state scimAttributeMappingResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ScimAttributeMappingAPI.UpdateScimAttributeMapping(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString(), plan.ScimResourceTypeName.ValueString())

	// Determine what update operations are necessary
	ops := createScimAttributeMappingOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ScimAttributeMappingAPI.UpdateScimAttributeMappingExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Scim Attribute Mapping", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readScimAttributeMappingResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultScimAttributeMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *scimAttributeMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state scimAttributeMappingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ScimAttributeMappingAPI.DeleteScimAttributeMappingExecute(r.apiClient.ScimAttributeMappingAPI.DeleteScimAttributeMapping(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.ScimResourceTypeName.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Scim Attribute Mapping", err, httpResp)
		return
	}
}

func (r *scimAttributeMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importScimAttributeMapping(ctx, req, resp)
}

func (r *defaultScimAttributeMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importScimAttributeMapping(ctx, req, resp)
}

func importScimAttributeMapping(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [scim-resource-type-name]/[scim-attribute-mapping-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scim_resource_type_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), split[1])...)
}
