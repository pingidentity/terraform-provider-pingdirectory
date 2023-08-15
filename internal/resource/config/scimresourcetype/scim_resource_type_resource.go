package scimresourcetype

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &scimResourceTypeResource{}
	_ resource.ResourceWithConfigure   = &scimResourceTypeResource{}
	_ resource.ResourceWithImportState = &scimResourceTypeResource{}
	_ resource.Resource                = &defaultScimResourceTypeResource{}
	_ resource.ResourceWithConfigure   = &defaultScimResourceTypeResource{}
	_ resource.ResourceWithImportState = &defaultScimResourceTypeResource{}
)

// Create a Scim Resource Type resource
func NewScimResourceTypeResource() resource.Resource {
	return &scimResourceTypeResource{}
}

func NewDefaultScimResourceTypeResource() resource.Resource {
	return &defaultScimResourceTypeResource{}
}

// scimResourceTypeResource is the resource implementation.
type scimResourceTypeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultScimResourceTypeResource is the resource implementation.
type defaultScimResourceTypeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *scimResourceTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scim_resource_type"
}

func (r *defaultScimResourceTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_scim_resource_type"
}

// Configure adds the provider configured client to the resource.
func (r *scimResourceTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultScimResourceTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type scimResourceTypeResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	LastUpdated                 types.String `tfsdk:"last_updated"`
	Notifications               types.Set    `tfsdk:"notifications"`
	RequiredActions             types.Set    `tfsdk:"required_actions"`
	Type                        types.String `tfsdk:"type"`
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
func (r *scimResourceTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	scimResourceTypeSchema(ctx, req, resp, false)
}

func (r *defaultScimResourceTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	scimResourceTypeSchema(ctx, req, resp, true)
}

func scimResourceTypeSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Scim Resource Type.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of SCIM Resource Type resource. Options are ['ldap-pass-through', 'ldap-mapping']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"ldap-pass-through", "ldap-mapping"}...),
				},
			},
			"core_schema": schema.StringAttribute{
				Description: "The core schema enforced on core attributes at the top level of a SCIM resource representation exposed by thisMapping SCIM Resource Type.",
				Optional:    true,
			},
			"required_schema_extension": schema.SetAttribute{
				Description: "Required additive schemas that are enforced on extension attributes in a SCIM resource representation for this Mapping SCIM Resource Type.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"optional_schema_extension": schema.SetAttribute{
				Description: "Optional additive schemas that are enforced on extension attributes in a SCIM resource representation for this Mapping SCIM Resource Type.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
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
				Default:     int64default.StaticInt64(500),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"schema_checking_option": schema.SetAttribute{
				Description: "Options to alter the way schema checking is performed during create or modify requests.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"structural_ldap_objectclass": schema.StringAttribute{
				Description: "Specifies the LDAP structural object class that should be exposed by this SCIM Resource Type.",
				Optional:    true,
			},
			"auxiliary_ldap_objectclass": schema.SetAttribute{
				Description: "Specifies an auxiliary LDAP object class that should be exposed by this SCIM Resource Type.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"include_base_dn": schema.StringAttribute{
				Description: "Specifies the base DN of the branch of the LDAP directory that can be accessed by this SCIM Resource Type.",
				Optional:    true,
			},
			"include_filter": schema.SetAttribute{
				Description: "The set of LDAP filters that define the LDAP entries that should be included in this SCIM Resource Type.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"include_operational_attribute": schema.SetAttribute{
				Description: "Specifies the set of operational LDAP attributes to be provided by this SCIM Resource Type.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"create_dn_pattern": schema.StringAttribute{
				Description: "Specifies the template to use for the DN when creating new entries.",
				Optional:    true,
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

// Add config validators that apply to both default_ and non-default_
func configValidatorsScimResourceType() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("optional_schema_extension"),
			path.MatchRoot("type"),
			[]string{"ldap-mapping"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("core_schema"),
			path.MatchRoot("type"),
			[]string{"ldap-mapping"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("required_schema_extension"),
			path.MatchRoot("type"),
			[]string{"ldap-mapping"},
		),
	}
}

// Add config validators
func (r scimResourceTypeResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsScimResourceType()
}

// Add config validators
func (r defaultScimResourceTypeResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsScimResourceType()
}

// Add optional fields to create request for ldap-pass-through scim-resource-type
func addOptionalLdapPassThroughScimResourceTypeFields(ctx context.Context, addRequest *client.AddLdapPassThroughScimResourceTypeRequest, plan scimResourceTypeResourceModel) error {
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

// Add optional fields to create request for ldap-mapping scim-resource-type
func addOptionalLdapMappingScimResourceTypeFields(ctx context.Context, addRequest *client.AddLdapMappingScimResourceTypeRequest, plan scimResourceTypeResourceModel) error {
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

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateScimResourceTypeUnknownValues(ctx context.Context, model *scimResourceTypeResourceModel) {
	if model.RequiredSchemaExtension.ElementType(ctx) == nil {
		model.RequiredSchemaExtension = types.SetNull(types.StringType)
	}
	if model.OptionalSchemaExtension.ElementType(ctx) == nil {
		model.OptionalSchemaExtension = types.SetNull(types.StringType)
	}
}

// Read a LdapPassThroughScimResourceTypeResponse object into the model struct
func readLdapPassThroughScimResourceTypeResponse(ctx context.Context, r *client.LdapPassThroughScimResourceTypeResponse, state *scimResourceTypeResourceModel, expectedValues *scimResourceTypeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap-pass-through")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
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
	populateScimResourceTypeUnknownValues(ctx, state)
}

// Read a LdapMappingScimResourceTypeResponse object into the model struct
func readLdapMappingScimResourceTypeResponse(ctx context.Context, r *client.LdapMappingScimResourceTypeResponse, state *scimResourceTypeResourceModel, expectedValues *scimResourceTypeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap-mapping")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
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
	populateScimResourceTypeUnknownValues(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createScimResourceTypeOperations(plan scimResourceTypeResourceModel, state scimResourceTypeResourceModel) []client.Operation {
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

// Create a ldap-pass-through scim-resource-type
func (r *scimResourceTypeResource) CreateLdapPassThroughScimResourceType(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan scimResourceTypeResourceModel) (*scimResourceTypeResourceModel, error) {
	addRequest := client.NewAddLdapPassThroughScimResourceTypeRequest(plan.Name.ValueString(),
		[]client.EnumldapPassThroughScimResourceTypeSchemaUrn{client.ENUMLDAPPASSTHROUGHSCIMRESOURCETYPESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SCIM_RESOURCE_TYPELDAP_PASS_THROUGH},
		plan.Enabled.ValueBool(),
		plan.Endpoint.ValueString())
	err := addOptionalLdapPassThroughScimResourceTypeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Scim Resource Type", err.Error())
		return nil, err
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Scim Resource Type", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state scimResourceTypeResourceModel
	readLdapPassThroughScimResourceTypeResponse(ctx, addResponse.LdapPassThroughScimResourceTypeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a ldap-mapping scim-resource-type
func (r *scimResourceTypeResource) CreateLdapMappingScimResourceType(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan scimResourceTypeResourceModel) (*scimResourceTypeResourceModel, error) {
	addRequest := client.NewAddLdapMappingScimResourceTypeRequest(plan.Name.ValueString(),
		[]client.EnumldapMappingScimResourceTypeSchemaUrn{client.ENUMLDAPMAPPINGSCIMRESOURCETYPESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SCIM_RESOURCE_TYPELDAP_MAPPING},
		plan.CoreSchema.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Endpoint.ValueString())
	err := addOptionalLdapMappingScimResourceTypeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Scim Resource Type", err.Error())
		return nil, err
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Scim Resource Type", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state scimResourceTypeResourceModel
	readLdapMappingScimResourceTypeResponse(ctx, addResponse.LdapMappingScimResourceTypeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *scimResourceTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan scimResourceTypeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *scimResourceTypeResourceModel
	var err error
	if plan.Type.ValueString() == "ldap-pass-through" {
		state, err = r.CreateLdapPassThroughScimResourceType(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "ldap-mapping" {
		state, err = r.CreateLdapMappingScimResourceType(ctx, req, resp, plan)
		if err != nil {
			return
		}
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
func (r *defaultScimResourceTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan scimResourceTypeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ScimResourceTypeApi.GetScimResourceType(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Scim Resource Type", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state scimResourceTypeResourceModel
	if readResponse.LdapPassThroughScimResourceTypeResponse != nil {
		readLdapPassThroughScimResourceTypeResponse(ctx, readResponse.LdapPassThroughScimResourceTypeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LdapMappingScimResourceTypeResponse != nil {
		readLdapMappingScimResourceTypeResponse(ctx, readResponse.LdapMappingScimResourceTypeResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ScimResourceTypeApi.UpdateScimResourceType(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createScimResourceTypeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ScimResourceTypeApi.UpdateScimResourceTypeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Scim Resource Type", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.LdapPassThroughScimResourceTypeResponse != nil {
			readLdapPassThroughScimResourceTypeResponse(ctx, updateResponse.LdapPassThroughScimResourceTypeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LdapMappingScimResourceTypeResponse != nil {
			readLdapMappingScimResourceTypeResponse(ctx, updateResponse.LdapMappingScimResourceTypeResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *scimResourceTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readScimResourceType(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultScimResourceTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readScimResourceType(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readScimResourceType(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state scimResourceTypeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ScimResourceTypeApi.GetScimResourceType(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Scim Resource Type", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Scim Resource Type", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.LdapPassThroughScimResourceTypeResponse != nil {
		readLdapPassThroughScimResourceTypeResponse(ctx, readResponse.LdapPassThroughScimResourceTypeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LdapMappingScimResourceTypeResponse != nil {
		readLdapMappingScimResourceTypeResponse(ctx, readResponse.LdapMappingScimResourceTypeResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *scimResourceTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateScimResourceType(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultScimResourceTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateScimResourceType(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateScimResourceType(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan scimResourceTypeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state scimResourceTypeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ScimResourceTypeApi.UpdateScimResourceType(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createScimResourceTypeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ScimResourceTypeApi.UpdateScimResourceTypeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Scim Resource Type", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.LdapPassThroughScimResourceTypeResponse != nil {
			readLdapPassThroughScimResourceTypeResponse(ctx, updateResponse.LdapPassThroughScimResourceTypeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LdapMappingScimResourceTypeResponse != nil {
			readLdapMappingScimResourceTypeResponse(ctx, updateResponse.LdapMappingScimResourceTypeResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *defaultScimResourceTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *scimResourceTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state scimResourceTypeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ScimResourceTypeApi.DeleteScimResourceTypeExecute(r.apiClient.ScimResourceTypeApi.DeleteScimResourceType(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Scim Resource Type", err, httpResp)
		return
	}
}

func (r *scimResourceTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importScimResourceType(ctx, req, resp)
}

func (r *defaultScimResourceTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importScimResourceType(ctx, req, resp)
}

func importScimResourceType(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
