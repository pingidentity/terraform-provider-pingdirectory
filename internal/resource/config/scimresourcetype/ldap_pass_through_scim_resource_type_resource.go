package scimresourcetype

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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
	_ resource.Resource                = &ldapPassThroughScimResourceTypeResource{}
	_ resource.ResourceWithConfigure   = &ldapPassThroughScimResourceTypeResource{}
	_ resource.ResourceWithImportState = &ldapPassThroughScimResourceTypeResource{}
	_ resource.Resource                = &defaultLdapPassThroughScimResourceTypeResource{}
	_ resource.ResourceWithConfigure   = &defaultLdapPassThroughScimResourceTypeResource{}
	_ resource.ResourceWithImportState = &defaultLdapPassThroughScimResourceTypeResource{}
)

// Create a Ldap Pass Through Scim Resource Type resource
func NewLdapPassThroughScimResourceTypeResource() resource.Resource {
	return &ldapPassThroughScimResourceTypeResource{}
}

func NewDefaultLdapPassThroughScimResourceTypeResource() resource.Resource {
	return &defaultLdapPassThroughScimResourceTypeResource{}
}

// ldapPassThroughScimResourceTypeResource is the resource implementation.
type ldapPassThroughScimResourceTypeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultLdapPassThroughScimResourceTypeResource is the resource implementation.
type defaultLdapPassThroughScimResourceTypeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *ldapPassThroughScimResourceTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ldap_pass_through_scim_resource_type"
}

func (r *defaultLdapPassThroughScimResourceTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_ldap_pass_through_scim_resource_type"
}

// Configure adds the provider configured client to the resource.
func (r *ldapPassThroughScimResourceTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultLdapPassThroughScimResourceTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type ldapPassThroughScimResourceTypeResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	LastUpdated                 types.String `tfsdk:"last_updated"`
	Notifications               types.Set    `tfsdk:"notifications"`
	RequiredActions             types.Set    `tfsdk:"required_actions"`
	Description                 types.String `tfsdk:"description"`
	Enabled                     types.Bool   `tfsdk:"enabled"`
	Endpoint                    types.String `tfsdk:"endpoint"`
	LookthroughLimit            types.Int64  `tfsdk:"lookthrough_limit"`
	SchemaCheckingOption        types.Set    `tfsdk:"schema_checking_option"`
	StructuralLDAPObjectclass   types.String `tfsdk:"structural_ldap_objectclass"`
	AuxiliaryLDAPObjectclass    types.Set    `tfsdk:"auxiliary_ldap_objectclass"`
	IncludeBaseDN               types.String `tfsdk:"include_base_dn"`
	IncludeFilter               types.Set    `tfsdk:"include_filter"`
	IncludeOperationalAttribute types.Set    `tfsdk:"include_operational_attribute"`
	CreateDNPattern             types.String `tfsdk:"create_dn_pattern"`
}

// GetSchema defines the schema for the resource.
func (r *ldapPassThroughScimResourceTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ldapPassThroughScimResourceTypeSchema(ctx, req, resp, false)
}

func (r *defaultLdapPassThroughScimResourceTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ldapPassThroughScimResourceTypeSchema(ctx, req, resp, true)
}

func ldapPassThroughScimResourceTypeSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Ldap Pass Through Scim Resource Type.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description: "A description for this SCIM Resource Type",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the SCIM Resource Type is enabled.",
				Required:    true,
			},
			"endpoint": schema.StringAttribute{
				Description: "The HTTP addressable endpoint of this SCIM Resource Type relative to the '/scim/v2' base URL. Do not include a leading '/'.",
				Required:    true,
			},
			"lookthrough_limit": schema.Int64Attribute{
				Description: "The maximum number of resources that the SCIM Resource Type should \"look through\" in the course of processing a search request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"schema_checking_option": schema.SetAttribute{
				Description: "Options to alter the way schema checking is performed during create or modify requests.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"structural_ldap_objectclass": schema.StringAttribute{
				Description: "Specifies the LDAP structural object class that should be exposed by this SCIM Resource Type.",
				Optional:    true,
			},
			"auxiliary_ldap_objectclass": schema.SetAttribute{
				Description: "Specifies an auxiliary LDAP object class that should be exposed by this SCIM Resource Type.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"include_base_dn": schema.StringAttribute{
				Description: "Specifies the base DN of the branch of the LDAP directory that can be accessed by this SCIM Resource Type.",
				Optional:    true,
			},
			"include_filter": schema.SetAttribute{
				Description: "The set of LDAP filters that define the LDAP entries that should be included in this SCIM Resource Type.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"include_operational_attribute": schema.SetAttribute{
				Description: "Specifies the set of operational LDAP attributes to be provided by this SCIM Resource Type.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"create_dn_pattern": schema.StringAttribute{
				Description: "Specifies the template to use for the DN when creating new entries.",
				Optional:    true,
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
func addOptionalLdapPassThroughScimResourceTypeFields(ctx context.Context, addRequest *client.AddLdapPassThroughScimResourceTypeRequest, plan ldapPassThroughScimResourceTypeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.LookthroughLimit) {
		addRequest.LookthroughLimit = plan.LookthroughLimit.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.SchemaCheckingOption) {
		var slice []string
		plan.SchemaCheckingOption.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumscimResourceTypeSchemaCheckingOptionProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumscimResourceTypeSchemaCheckingOptionPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.SchemaCheckingOption = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.StructuralLDAPObjectclass) {
		addRequest.StructuralLDAPObjectclass = plan.StructuralLDAPObjectclass.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AuxiliaryLDAPObjectclass) {
		var slice []string
		plan.AuxiliaryLDAPObjectclass.ElementsAs(ctx, &slice, false)
		addRequest.AuxiliaryLDAPObjectclass = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IncludeBaseDN) {
		addRequest.IncludeBaseDN = plan.IncludeBaseDN.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludeFilter) {
		var slice []string
		plan.IncludeFilter.ElementsAs(ctx, &slice, false)
		addRequest.IncludeFilter = slice
	}
	if internaltypes.IsDefined(plan.IncludeOperationalAttribute) {
		var slice []string
		plan.IncludeOperationalAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeOperationalAttribute = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CreateDNPattern) {
		addRequest.CreateDNPattern = plan.CreateDNPattern.ValueStringPointer()
	}
	return nil
}

// Read a LdapPassThroughScimResourceTypeResponse object into the model struct
func readLdapPassThroughScimResourceTypeResponse(ctx context.Context, r *client.LdapPassThroughScimResourceTypeResponse, state *ldapPassThroughScimResourceTypeResourceModel, expectedValues *ldapPassThroughScimResourceTypeResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Endpoint = types.StringValue(r.Endpoint)
	state.LookthroughLimit = internaltypes.Int64TypeOrNil(r.LookthroughLimit)
	state.SchemaCheckingOption = internaltypes.GetStringSet(
		client.StringSliceEnumscimResourceTypeSchemaCheckingOptionProp(r.SchemaCheckingOption))
	state.StructuralLDAPObjectclass = internaltypes.StringTypeOrNil(r.StructuralLDAPObjectclass, internaltypes.IsEmptyString(expectedValues.StructuralLDAPObjectclass))
	state.AuxiliaryLDAPObjectclass = internaltypes.GetStringSet(r.AuxiliaryLDAPObjectclass)
	state.IncludeBaseDN = internaltypes.StringTypeOrNil(r.IncludeBaseDN, internaltypes.IsEmptyString(expectedValues.IncludeBaseDN))
	state.IncludeFilter = internaltypes.GetStringSet(r.IncludeFilter)
	state.IncludeOperationalAttribute = internaltypes.GetStringSet(r.IncludeOperationalAttribute)
	state.CreateDNPattern = internaltypes.StringTypeOrNil(r.CreateDNPattern, internaltypes.IsEmptyString(expectedValues.CreateDNPattern))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createLdapPassThroughScimResourceTypeOperations(plan ldapPassThroughScimResourceTypeResourceModel, state ldapPassThroughScimResourceTypeResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.Endpoint, state.Endpoint, "endpoint")
	operations.AddInt64OperationIfNecessary(&ops, plan.LookthroughLimit, state.LookthroughLimit, "lookthrough-limit")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SchemaCheckingOption, state.SchemaCheckingOption, "schema-checking-option")
	operations.AddStringOperationIfNecessary(&ops, plan.StructuralLDAPObjectclass, state.StructuralLDAPObjectclass, "structural-ldap-objectclass")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AuxiliaryLDAPObjectclass, state.AuxiliaryLDAPObjectclass, "auxiliary-ldap-objectclass")
	operations.AddStringOperationIfNecessary(&ops, plan.IncludeBaseDN, state.IncludeBaseDN, "include-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeFilter, state.IncludeFilter, "include-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeOperationalAttribute, state.IncludeOperationalAttribute, "include-operational-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.CreateDNPattern, state.CreateDNPattern, "create-dn-pattern")
	return ops
}

// Create a new resource
func (r *ldapPassThroughScimResourceTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ldapPassThroughScimResourceTypeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddLdapPassThroughScimResourceTypeRequest(plan.Id.ValueString(),
		[]client.EnumldapPassThroughScimResourceTypeSchemaUrn{client.ENUMLDAPPASSTHROUGHSCIMRESOURCETYPESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SCIM_RESOURCE_TYPELDAP_PASS_THROUGH},
		plan.Enabled.ValueBool(),
		plan.Endpoint.ValueString())
	err := addOptionalLdapPassThroughScimResourceTypeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Ldap Pass Through Scim Resource Type", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ScimResourceTypeApi.AddScimResourceType(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddScimResourceTypeRequest(
		client.AddLdapPassThroughScimResourceTypeRequestAsAddScimResourceTypeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ScimResourceTypeApi.AddScimResourceTypeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Ldap Pass Through Scim Resource Type", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state ldapPassThroughScimResourceTypeResourceModel
	readLdapPassThroughScimResourceTypeResponse(ctx, addResponse.LdapPassThroughScimResourceTypeResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultLdapPassThroughScimResourceTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ldapPassThroughScimResourceTypeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ScimResourceTypeApi.GetScimResourceType(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ldap Pass Through Scim Resource Type", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state ldapPassThroughScimResourceTypeResourceModel
	readLdapPassThroughScimResourceTypeResponse(ctx, readResponse.LdapPassThroughScimResourceTypeResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ScimResourceTypeApi.UpdateScimResourceType(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createLdapPassThroughScimResourceTypeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ScimResourceTypeApi.UpdateScimResourceTypeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ldap Pass Through Scim Resource Type", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLdapPassThroughScimResourceTypeResponse(ctx, updateResponse.LdapPassThroughScimResourceTypeResponse, &state, &plan, &resp.Diagnostics)
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
func (r *ldapPassThroughScimResourceTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLdapPassThroughScimResourceType(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLdapPassThroughScimResourceTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLdapPassThroughScimResourceType(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readLdapPassThroughScimResourceType(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state ldapPassThroughScimResourceTypeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ScimResourceTypeApi.GetScimResourceType(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ldap Pass Through Scim Resource Type", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLdapPassThroughScimResourceTypeResponse(ctx, readResponse.LdapPassThroughScimResourceTypeResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *ldapPassThroughScimResourceTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLdapPassThroughScimResourceType(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLdapPassThroughScimResourceTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLdapPassThroughScimResourceType(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateLdapPassThroughScimResourceType(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan ldapPassThroughScimResourceTypeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state ldapPassThroughScimResourceTypeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ScimResourceTypeApi.UpdateScimResourceType(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createLdapPassThroughScimResourceTypeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ScimResourceTypeApi.UpdateScimResourceTypeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ldap Pass Through Scim Resource Type", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLdapPassThroughScimResourceTypeResponse(ctx, updateResponse.LdapPassThroughScimResourceTypeResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultLdapPassThroughScimResourceTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *ldapPassThroughScimResourceTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ldapPassThroughScimResourceTypeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ScimResourceTypeApi.DeleteScimResourceTypeExecute(r.apiClient.ScimResourceTypeApi.DeleteScimResourceType(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Ldap Pass Through Scim Resource Type", err, httpResp)
		return
	}
}

func (r *ldapPassThroughScimResourceTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLdapPassThroughScimResourceType(ctx, req, resp)
}

func (r *defaultLdapPassThroughScimResourceTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLdapPassThroughScimResourceType(ctx, req, resp)
}

func importLdapPassThroughScimResourceType(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
