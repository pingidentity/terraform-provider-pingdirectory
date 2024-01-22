package ldapcorrelationattributepair

import (
	"context"
	"strings"

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
	_ resource.Resource                = &ldapCorrelationAttributePairResource{}
	_ resource.ResourceWithConfigure   = &ldapCorrelationAttributePairResource{}
	_ resource.ResourceWithImportState = &ldapCorrelationAttributePairResource{}
	_ resource.Resource                = &defaultLdapCorrelationAttributePairResource{}
	_ resource.ResourceWithConfigure   = &defaultLdapCorrelationAttributePairResource{}
	_ resource.ResourceWithImportState = &defaultLdapCorrelationAttributePairResource{}
)

// Create a Ldap Correlation Attribute Pair resource
func NewLdapCorrelationAttributePairResource() resource.Resource {
	return &ldapCorrelationAttributePairResource{}
}

func NewDefaultLdapCorrelationAttributePairResource() resource.Resource {
	return &defaultLdapCorrelationAttributePairResource{}
}

// ldapCorrelationAttributePairResource is the resource implementation.
type ldapCorrelationAttributePairResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultLdapCorrelationAttributePairResource is the resource implementation.
type defaultLdapCorrelationAttributePairResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *ldapCorrelationAttributePairResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ldap_correlation_attribute_pair"
}

func (r *defaultLdapCorrelationAttributePairResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_ldap_correlation_attribute_pair"
}

// Configure adds the provider configured client to the resource.
func (r *ldapCorrelationAttributePairResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultLdapCorrelationAttributePairResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type ldapCorrelationAttributePairResourceModel struct {
	Id                            types.String `tfsdk:"id"`
	Name                          types.String `tfsdk:"name"`
	Notifications                 types.Set    `tfsdk:"notifications"`
	RequiredActions               types.Set    `tfsdk:"required_actions"`
	Type                          types.String `tfsdk:"type"`
	CorrelatedLdapDataViewName    types.String `tfsdk:"correlated_ldap_data_view_name"`
	ScimResourceTypeName          types.String `tfsdk:"scim_resource_type_name"`
	PrimaryCorrelationAttribute   types.String `tfsdk:"primary_correlation_attribute"`
	SecondaryCorrelationAttribute types.String `tfsdk:"secondary_correlation_attribute"`
}

// GetSchema defines the schema for the resource.
func (r *ldapCorrelationAttributePairResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ldapCorrelationAttributePairSchema(ctx, req, resp, false)
}

func (r *defaultLdapCorrelationAttributePairResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ldapCorrelationAttributePairSchema(ctx, req, resp, true)
}

func ldapCorrelationAttributePairSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Ldap Correlation Attribute Pair.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of LDAP Correlation Attribute Pair resource. Options are ['ldap-correlation-attribute-pair']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("ldap-correlation-attribute-pair"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"ldap-correlation-attribute-pair"}...),
				},
			},
			"correlated_ldap_data_view_name": schema.StringAttribute{
				Description: "Name of the parent Correlated LDAP Data View",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"scim_resource_type_name": schema.StringAttribute{
				Description: "Name of the parent SCIM Resource Type",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"primary_correlation_attribute": schema.StringAttribute{
				Description: "The LDAP attribute from the base SCIM Resource Type whose value will be used to match objects in the Correlated LDAP Data View.",
				Required:    true,
			},
			"secondary_correlation_attribute": schema.StringAttribute{
				Description: "The LDAP attribute from the Correlated LDAP Data View whose value will be matched.",
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
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type", "correlated_ldap_data_view_name", "scim_resource_type_name"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add optional fields to create request for ldap-correlation-attribute-pair ldap-correlation-attribute-pair
func addOptionalLdapCorrelationAttributePairFields(ctx context.Context, addRequest *client.AddLdapCorrelationAttributePairRequest, plan ldapCorrelationAttributePairResourceModel) {
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *ldapCorrelationAttributePairResourceModel) populateAllComputedStringAttributes() {
	if model.PrimaryCorrelationAttribute.IsUnknown() || model.PrimaryCorrelationAttribute.IsNull() {
		model.PrimaryCorrelationAttribute = types.StringValue("")
	}
	if model.SecondaryCorrelationAttribute.IsUnknown() || model.SecondaryCorrelationAttribute.IsNull() {
		model.SecondaryCorrelationAttribute = types.StringValue("")
	}
}

// Read a LdapCorrelationAttributePairResponse object into the model struct
func readLdapCorrelationAttributePairResponse(ctx context.Context, r *client.LdapCorrelationAttributePairResponse, state *ldapCorrelationAttributePairResourceModel, expectedValues *ldapCorrelationAttributePairResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap-correlation-attribute-pair")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PrimaryCorrelationAttribute = types.StringValue(r.PrimaryCorrelationAttribute)
	state.SecondaryCorrelationAttribute = types.StringValue(r.SecondaryCorrelationAttribute)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *ldapCorrelationAttributePairResourceModel) setStateValuesNotReturnedByAPI(expectedValues *ldapCorrelationAttributePairResourceModel) {
	if !expectedValues.CorrelatedLdapDataViewName.IsUnknown() {
		state.CorrelatedLdapDataViewName = expectedValues.CorrelatedLdapDataViewName
	}
	if !expectedValues.ScimResourceTypeName.IsUnknown() {
		state.ScimResourceTypeName = expectedValues.ScimResourceTypeName
	}
}

// Create any update operations necessary to make the state match the plan
func createLdapCorrelationAttributePairOperations(plan ldapCorrelationAttributePairResourceModel, state ldapCorrelationAttributePairResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.PrimaryCorrelationAttribute, state.PrimaryCorrelationAttribute, "primary-correlation-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.SecondaryCorrelationAttribute, state.SecondaryCorrelationAttribute, "secondary-correlation-attribute")
	return ops
}

// Create a ldap-correlation-attribute-pair ldap-correlation-attribute-pair
func (r *ldapCorrelationAttributePairResource) CreateLdapCorrelationAttributePair(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan ldapCorrelationAttributePairResourceModel) (*ldapCorrelationAttributePairResourceModel, error) {
	addRequest := client.NewAddLdapCorrelationAttributePairRequest(plan.PrimaryCorrelationAttribute.ValueString(),
		plan.SecondaryCorrelationAttribute.ValueString(),
		plan.Name.ValueString())
	addOptionalLdapCorrelationAttributePairFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LdapCorrelationAttributePairAPI.AddLdapCorrelationAttributePair(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.CorrelatedLdapDataViewName.ValueString(), plan.ScimResourceTypeName.ValueString())
	apiAddRequest = apiAddRequest.AddLdapCorrelationAttributePairRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.LdapCorrelationAttributePairAPI.AddLdapCorrelationAttributePairExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Ldap Correlation Attribute Pair", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state ldapCorrelationAttributePairResourceModel
	readLdapCorrelationAttributePairResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *ldapCorrelationAttributePairResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ldapCorrelationAttributePairResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateLdapCorrelationAttributePair(ctx, req, resp, plan)
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
func (r *defaultLdapCorrelationAttributePairResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ldapCorrelationAttributePairResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LdapCorrelationAttributePairAPI.GetLdapCorrelationAttributePair(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.CorrelatedLdapDataViewName.ValueString(), plan.ScimResourceTypeName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ldap Correlation Attribute Pair", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state ldapCorrelationAttributePairResourceModel
	readLdapCorrelationAttributePairResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LdapCorrelationAttributePairAPI.UpdateLdapCorrelationAttributePair(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.CorrelatedLdapDataViewName.ValueString(), plan.ScimResourceTypeName.ValueString())
	ops := createLdapCorrelationAttributePairOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LdapCorrelationAttributePairAPI.UpdateLdapCorrelationAttributePairExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ldap Correlation Attribute Pair", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLdapCorrelationAttributePairResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *ldapCorrelationAttributePairResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLdapCorrelationAttributePair(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultLdapCorrelationAttributePairResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLdapCorrelationAttributePair(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readLdapCorrelationAttributePair(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state ldapCorrelationAttributePairResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LdapCorrelationAttributePairAPI.GetLdapCorrelationAttributePair(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString(), state.CorrelatedLdapDataViewName.ValueString(), state.ScimResourceTypeName.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Ldap Correlation Attribute Pair", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ldap Correlation Attribute Pair", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLdapCorrelationAttributePairResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *ldapCorrelationAttributePairResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLdapCorrelationAttributePair(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLdapCorrelationAttributePairResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLdapCorrelationAttributePair(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateLdapCorrelationAttributePair(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan ldapCorrelationAttributePairResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state ldapCorrelationAttributePairResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LdapCorrelationAttributePairAPI.UpdateLdapCorrelationAttributePair(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString(), plan.CorrelatedLdapDataViewName.ValueString(), plan.ScimResourceTypeName.ValueString())

	// Determine what update operations are necessary
	ops := createLdapCorrelationAttributePairOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LdapCorrelationAttributePairAPI.UpdateLdapCorrelationAttributePairExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ldap Correlation Attribute Pair", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLdapCorrelationAttributePairResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultLdapCorrelationAttributePairResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *ldapCorrelationAttributePairResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ldapCorrelationAttributePairResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LdapCorrelationAttributePairAPI.DeleteLdapCorrelationAttributePairExecute(r.apiClient.LdapCorrelationAttributePairAPI.DeleteLdapCorrelationAttributePair(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.CorrelatedLdapDataViewName.ValueString(), state.ScimResourceTypeName.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Ldap Correlation Attribute Pair", err, httpResp)
		return
	}
}

func (r *ldapCorrelationAttributePairResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLdapCorrelationAttributePair(ctx, req, resp)
}

func (r *defaultLdapCorrelationAttributePairResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLdapCorrelationAttributePair(ctx, req, resp)
}

func importLdapCorrelationAttributePair(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 3 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [scim-resource-type-name]/[correlated-ldap-data-view-name]/[ldap-correlation-attribute-pair-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scim_resource_type_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("correlated_ldap_data_view_name"), split[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), split[2])...)
}
