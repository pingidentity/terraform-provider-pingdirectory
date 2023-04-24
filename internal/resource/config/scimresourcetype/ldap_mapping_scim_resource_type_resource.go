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
	_ resource.Resource                = &ldapMappingScimResourceTypeResource{}
	_ resource.ResourceWithConfigure   = &ldapMappingScimResourceTypeResource{}
	_ resource.ResourceWithImportState = &ldapMappingScimResourceTypeResource{}
	_ resource.Resource                = &defaultLdapMappingScimResourceTypeResource{}
	_ resource.ResourceWithConfigure   = &defaultLdapMappingScimResourceTypeResource{}
	_ resource.ResourceWithImportState = &defaultLdapMappingScimResourceTypeResource{}
)

// Create a Ldap Mapping Scim Resource Type resource
func NewLdapMappingScimResourceTypeResource() resource.Resource {
	return &ldapMappingScimResourceTypeResource{}
}

func NewDefaultLdapMappingScimResourceTypeResource() resource.Resource {
	return &defaultLdapMappingScimResourceTypeResource{}
}

// ldapMappingScimResourceTypeResource is the resource implementation.
type ldapMappingScimResourceTypeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultLdapMappingScimResourceTypeResource is the resource implementation.
type defaultLdapMappingScimResourceTypeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *ldapMappingScimResourceTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ldap_mapping_scim_resource_type"
}

func (r *defaultLdapMappingScimResourceTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_ldap_mapping_scim_resource_type"
}

// Configure adds the provider configured client to the resource.
func (r *ldapMappingScimResourceTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultLdapMappingScimResourceTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type ldapMappingScimResourceTypeResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	LastUpdated                 types.String `tfsdk:"last_updated"`
	Notifications               types.Set    `tfsdk:"notifications"`
	RequiredActions             types.Set    `tfsdk:"required_actions"`
	CoreSchema                  types.String `tfsdk:"core_schema"`
	RequiredSchemaExtension     types.Set    `tfsdk:"required_schema_extension"`
	OptionalSchemaExtension     types.Set    `tfsdk:"optional_schema_extension"`
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
func (r *ldapMappingScimResourceTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ldapMappingScimResourceTypeSchema(ctx, req, resp, false)
}

func (r *defaultLdapMappingScimResourceTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ldapMappingScimResourceTypeSchema(ctx, req, resp, true)
}

func ldapMappingScimResourceTypeSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Ldap Mapping Scim Resource Type.",
		Attributes: map[string]schema.Attribute{
			"core_schema": schema.StringAttribute{
				Description: "The core schema enforced on core attributes at the top level of a SCIM resource representation exposed by thisMapping SCIM Resource Type.",
				Required:    true,
			},
			"required_schema_extension": schema.SetAttribute{
				Description: "Required additive schemas that are enforced on extension attributes in a SCIM resource representation for this Mapping SCIM Resource Type.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"optional_schema_extension": schema.SetAttribute{
				Description: "Optional additive schemas that are enforced on extension attributes in a SCIM resource representation for this Mapping SCIM Resource Type.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
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
func addOptionalLdapMappingScimResourceTypeFields(ctx context.Context, addRequest *client.AddLdapMappingScimResourceTypeRequest, plan ldapMappingScimResourceTypeResourceModel) error {
	if internaltypes.IsDefined(plan.RequiredSchemaExtension) {
		var slice []string
		plan.RequiredSchemaExtension.ElementsAs(ctx, &slice, false)
		addRequest.RequiredSchemaExtension = slice
	}
	if internaltypes.IsDefined(plan.OptionalSchemaExtension) {
		var slice []string
		plan.OptionalSchemaExtension.ElementsAs(ctx, &slice, false)
		addRequest.OptionalSchemaExtension = slice
	}
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

// Read a LdapMappingScimResourceTypeResponse object into the model struct
func readLdapMappingScimResourceTypeResponse(ctx context.Context, r *client.LdapMappingScimResourceTypeResponse, state *ldapMappingScimResourceTypeResourceModel, expectedValues *ldapMappingScimResourceTypeResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.CoreSchema = types.StringValue(r.CoreSchema)
	state.RequiredSchemaExtension = internaltypes.GetStringSet(r.RequiredSchemaExtension)
	state.OptionalSchemaExtension = internaltypes.GetStringSet(r.OptionalSchemaExtension)
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
func createLdapMappingScimResourceTypeOperations(plan ldapMappingScimResourceTypeResourceModel, state ldapMappingScimResourceTypeResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.CoreSchema, state.CoreSchema, "core-schema")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RequiredSchemaExtension, state.RequiredSchemaExtension, "required-schema-extension")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.OptionalSchemaExtension, state.OptionalSchemaExtension, "optional-schema-extension")
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
func (r *ldapMappingScimResourceTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ldapMappingScimResourceTypeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddLdapMappingScimResourceTypeRequest(plan.Id.ValueString(),
		[]client.EnumldapMappingScimResourceTypeSchemaUrn{client.ENUMLDAPMAPPINGSCIMRESOURCETYPESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SCIM_RESOURCE_TYPELDAP_MAPPING},
		plan.CoreSchema.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Endpoint.ValueString())
	err := addOptionalLdapMappingScimResourceTypeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Ldap Mapping Scim Resource Type", err.Error())
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
		client.AddLdapMappingScimResourceTypeRequestAsAddScimResourceTypeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ScimResourceTypeApi.AddScimResourceTypeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Ldap Mapping Scim Resource Type", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state ldapMappingScimResourceTypeResourceModel
	readLdapMappingScimResourceTypeResponse(ctx, addResponse.LdapMappingScimResourceTypeResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultLdapMappingScimResourceTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ldapMappingScimResourceTypeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ScimResourceTypeApi.GetScimResourceType(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ldap Mapping Scim Resource Type", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state ldapMappingScimResourceTypeResourceModel
	readLdapMappingScimResourceTypeResponse(ctx, readResponse.LdapMappingScimResourceTypeResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ScimResourceTypeApi.UpdateScimResourceType(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createLdapMappingScimResourceTypeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ScimResourceTypeApi.UpdateScimResourceTypeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ldap Mapping Scim Resource Type", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLdapMappingScimResourceTypeResponse(ctx, updateResponse.LdapMappingScimResourceTypeResponse, &state, &plan, &resp.Diagnostics)
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
func (r *ldapMappingScimResourceTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLdapMappingScimResourceType(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLdapMappingScimResourceTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLdapMappingScimResourceType(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readLdapMappingScimResourceType(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state ldapMappingScimResourceTypeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ScimResourceTypeApi.GetScimResourceType(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ldap Mapping Scim Resource Type", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLdapMappingScimResourceTypeResponse(ctx, readResponse.LdapMappingScimResourceTypeResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *ldapMappingScimResourceTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLdapMappingScimResourceType(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLdapMappingScimResourceTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLdapMappingScimResourceType(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateLdapMappingScimResourceType(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan ldapMappingScimResourceTypeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state ldapMappingScimResourceTypeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ScimResourceTypeApi.UpdateScimResourceType(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createLdapMappingScimResourceTypeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ScimResourceTypeApi.UpdateScimResourceTypeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ldap Mapping Scim Resource Type", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLdapMappingScimResourceTypeResponse(ctx, updateResponse.LdapMappingScimResourceTypeResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultLdapMappingScimResourceTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *ldapMappingScimResourceTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ldapMappingScimResourceTypeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ScimResourceTypeApi.DeleteScimResourceTypeExecute(r.apiClient.ScimResourceTypeApi.DeleteScimResourceType(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Ldap Mapping Scim Resource Type", err, httpResp)
		return
	}
}

func (r *ldapMappingScimResourceTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLdapMappingScimResourceType(ctx, req, resp)
}

func (r *defaultLdapMappingScimResourceTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLdapMappingScimResourceType(ctx, req, resp)
}

func importLdapMappingScimResourceType(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
