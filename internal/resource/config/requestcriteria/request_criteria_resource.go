package requestcriteria

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &requestCriteriaResource{}
	_ resource.ResourceWithConfigure   = &requestCriteriaResource{}
	_ resource.ResourceWithImportState = &requestCriteriaResource{}
	_ resource.Resource                = &defaultRequestCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultRequestCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultRequestCriteriaResource{}
)

// Create a Request Criteria resource
func NewRequestCriteriaResource() resource.Resource {
	return &requestCriteriaResource{}
}

func NewDefaultRequestCriteriaResource() resource.Resource {
	return &defaultRequestCriteriaResource{}
}

// requestCriteriaResource is the resource implementation.
type requestCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultRequestCriteriaResource is the resource implementation.
type defaultRequestCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *requestCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_criteria"
}

func (r *defaultRequestCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_request_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *requestCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultRequestCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type requestCriteriaResourceModel struct {
	Id                                     types.String `tfsdk:"id"`
	Name                                   types.String `tfsdk:"name"`
	Notifications                          types.Set    `tfsdk:"notifications"`
	RequiredActions                        types.Set    `tfsdk:"required_actions"`
	Type                                   types.String `tfsdk:"type"`
	ExtensionClass                         types.String `tfsdk:"extension_class"`
	ExtensionArgument                      types.Set    `tfsdk:"extension_argument"`
	AllIncludedRequestCriteria             types.Set    `tfsdk:"all_included_request_criteria"`
	AnyIncludedRequestCriteria             types.Set    `tfsdk:"any_included_request_criteria"`
	NotAllIncludedRequestCriteria          types.Set    `tfsdk:"not_all_included_request_criteria"`
	NoneIncludedRequestCriteria            types.Set    `tfsdk:"none_included_request_criteria"`
	OperationType                          types.Set    `tfsdk:"operation_type"`
	OperationOrigin                        types.Set    `tfsdk:"operation_origin"`
	ConnectionCriteria                     types.String `tfsdk:"connection_criteria"`
	AllIncludedRequestControl              types.Set    `tfsdk:"all_included_request_control"`
	AnyIncludedRequestControl              types.Set    `tfsdk:"any_included_request_control"`
	NotAllIncludedRequestControl           types.Set    `tfsdk:"not_all_included_request_control"`
	NoneIncludedRequestControl             types.Set    `tfsdk:"none_included_request_control"`
	IncludedTargetEntryDN                  types.Set    `tfsdk:"included_target_entry_dn"`
	ExcludedTargetEntryDN                  types.Set    `tfsdk:"excluded_target_entry_dn"`
	AllIncludedTargetEntryFilter           types.Set    `tfsdk:"all_included_target_entry_filter"`
	AnyIncludedTargetEntryFilter           types.Set    `tfsdk:"any_included_target_entry_filter"`
	NotAllIncludedTargetEntryFilter        types.Set    `tfsdk:"not_all_included_target_entry_filter"`
	NoneIncludedTargetEntryFilter          types.Set    `tfsdk:"none_included_target_entry_filter"`
	AllIncludedTargetEntryGroupDN          types.Set    `tfsdk:"all_included_target_entry_group_dn"`
	AnyIncludedTargetEntryGroupDN          types.Set    `tfsdk:"any_included_target_entry_group_dn"`
	NotAllIncludedTargetEntryGroupDN       types.Set    `tfsdk:"not_all_included_target_entry_group_dn"`
	NoneIncludedTargetEntryGroupDN         types.Set    `tfsdk:"none_included_target_entry_group_dn"`
	TargetBindType                         types.Set    `tfsdk:"target_bind_type"`
	IncludedTargetSASLMechanism            types.Set    `tfsdk:"included_target_sasl_mechanism"`
	ExcludedTargetSASLMechanism            types.Set    `tfsdk:"excluded_target_sasl_mechanism"`
	IncludedTargetAttribute                types.Set    `tfsdk:"included_target_attribute"`
	ExcludedTargetAttribute                types.Set    `tfsdk:"excluded_target_attribute"`
	IncludedExtendedOperationOID           types.Set    `tfsdk:"included_extended_operation_oid"`
	ExcludedExtendedOperationOID           types.Set    `tfsdk:"excluded_extended_operation_oid"`
	IncludedSearchScope                    types.Set    `tfsdk:"included_search_scope"`
	UsingAdministrativeSessionWorkerThread types.String `tfsdk:"using_administrative_session_worker_thread"`
	IncludedApplicationName                types.Set    `tfsdk:"included_application_name"`
	ExcludedApplicationName                types.Set    `tfsdk:"excluded_application_name"`
	Description                            types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *requestCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	requestCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultRequestCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	requestCriteriaSchema(ctx, req, resp, true)
}

func requestCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Request Criteria.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Request Criteria resource. Options are ['root-dse', 'simple', 'aggregate', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"root-dse", "simple", "aggregate", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Request Criteria.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Request Criteria. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"all_included_request_criteria": schema.SetAttribute{
				Description: "Specifies a request criteria object that must match the associated operation request in order to match the aggregate request criteria. If one or more all-included request criteria objects are provided, then an operation request must match all of them in order to match the aggregate request criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"any_included_request_criteria": schema.SetAttribute{
				Description: "Specifies a request criteria object that may match the associated operation request in order to the this aggregate request criteria. If one or more any-included request criteria objects are provided, then an operation request must match at least one of them in order to match the aggregate request criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"not_all_included_request_criteria": schema.SetAttribute{
				Description: "Specifies a request criteria object that should not match the associated operation request in order to match the aggregate request criteria. If one or more not-all-included request criteria objects are provided, then an operation request must not match all of them (that is, it may match zero or more of them, but it must not match all of them) in order to match the aggregate request criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"none_included_request_criteria": schema.SetAttribute{
				Description: "Specifies a request criteria object that must not match the associated operation request in order to match the aggregate request criteria. If one or more none-included request criteria objects are provided, then an operation request must not match any of them in order to match the aggregate request criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"operation_type": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `root-dse`: The types of operations that may be matched by this Root DSE Request Criteria. When the `type` attribute is set to `simple`: Specifies the operation type(s) for operations that should be included in this Simple Request Criteria.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `root-dse`: The types of operations that may be matched by this Root DSE Request Criteria.\n  - `simple`: Specifies the operation type(s) for operations that should be included in this Simple Request Criteria.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"operation_origin": schema.SetAttribute{
				Description: "Specifies the origin for operations to be included in this Simple Request Criteria. If no values are provided, then the operation origin will not be taken into consideration when determining whether an operation matches this Simple Request Criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"connection_criteria": schema.StringAttribute{
				Description: "Specifies a connection criteria object that must match the associated client connection for operations included in this Simple Request Criteria.",
				Optional:    true,
			},
			"all_included_request_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must be present in the request from the client for operations included in this Simple Request Criteria. If any control OIDs are provided, then the request must contain all of those controls.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"any_included_request_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that may be present in the request from the client for operations included in this Simple Request Criteria. If any control OIDs are provided, then the request must contain at least one of those controls.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"not_all_included_request_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that should not be present in the request from the client for operations included in this Simple Request Criteria. If any control OIDs are provided, then the request must not contain at least one of those controls (that is, the request may contain zero or more of those controls, but not all of them).",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"none_included_request_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must not be present in the request from the client for operations included in this Simple Request Criteria. If any control OIDs are provided, then the request must not contain any of those controls.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"included_target_entry_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which targeted entries may exist for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"excluded_target_entry_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which targeted entries may not exist for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"all_included_target_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that must match the target entry for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any filters are provided, then the target entry must match all of those filters.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"any_included_target_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that may match the target entry for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any filters are provided, then the target entry must match at least one of those filters.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"not_all_included_target_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that should not match the target entry for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any filters are provided, then the target entry must not match at least one of those filters (that is, the request may match zero or more of those filters, but not of all of them).",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"none_included_target_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that must not match the target entry for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any filters are provided, then the target entry must not match any of those filters.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"all_included_target_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the target entry must be a member for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any group DNs are provided, then the target entry must be a member of all of those groups.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"any_included_target_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the target entry may be a member for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any group DNs are provided, then the target entry must be a member of at least one of those groups.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"not_all_included_target_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the target entry should not be a member for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any group DNs are provided, then the target entry must not be a member of at least one of those groups (that is, the target entry may be a member of zero or more of those groups, but not all of them).",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"none_included_target_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the target entry must not be a member for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any group DNs are provided, then the target entry must not be a member of any of those groups.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"target_bind_type": schema.SetAttribute{
				Description: "Specifies the authentication type for bind requests included in this Simple Request Criteria. This will only be taken into account for bind operations and will be ignored for any other type of operation. If no values are provided, then the authentication type will not be considered when determining whether the request should be included in this Simple Request Criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"included_target_sasl_mechanism": schema.SetAttribute{
				Description: "Specifies the name of a SASL mechanism for bind requests included in this Simple Request Criteria. This will only be taken into account for SASL bind operations and will be ignored for other types of operations and for bind operations that do not use SASL authentication.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"excluded_target_sasl_mechanism": schema.SetAttribute{
				Description: "Specifies the name of a SASL mechanism for bind requests excluded from this Simple Request Criteria. This will only be taken into account for SASL bind operations and will be ignored for other types of operations and for bind operations that do not use SASL authentication.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"included_target_attribute": schema.SetAttribute{
				Description: "Specifies the name or OID of an attribute type which must be targeted by requests included in this Simple Request Criteria. This will only be taken into account for add, compare, modify, modify DN, and search operations. It will be ignored for abandon, bind, delete, extended, and unbind operations.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"excluded_target_attribute": schema.SetAttribute{
				Description: "Specifies the name or OID of an attribute type which must not be targeted by requests included in this Simple Request Criteria. This will only be taken into account for add, compare, modify, modify DN, and search operations. It will be ignored for abandon, bind, delete, extended, and unbind operations.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"included_extended_operation_oid": schema.SetAttribute{
				Description: "Specifies the request OID for extended requests included in this Simple Request Criteria. This will only be taken into account for extended requests and will be ignored for all other types of requests.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"excluded_extended_operation_oid": schema.SetAttribute{
				Description: "Specifies the request OID for extended requests excluded from this Simple Request Criteria. This will only be taken into account for extended requests and will be ignored for all other types of requests.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"included_search_scope": schema.SetAttribute{
				Description: "Specifies the search scope values included in this Simple Request Criteria. This will only be taken into account for search requests and will be ignored for all other types of requests.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"using_administrative_session_worker_thread": schema.StringAttribute{
				Description: "Indicates whether operations being processed using a dedicated administrative session worker thread should be included in this Simple Request Criteria.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"true", "false", "any"}...),
				},
			},
			"included_application_name": schema.SetAttribute{
				Description: "Specifies an application name for requests included in this Simple Request Criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"excluded_application_name": schema.SetAttribute{
				Description: "Specifies an application name for requests excluded from this Simple Request Criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Request Criteria",
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
	} else {
		// Add RequiresReplace modifier for read-only attributes
		extensionClassAttr := schemaDef.Attributes["extension_class"].(schema.StringAttribute)
		extensionClassAttr.PlanModifiers = append(extensionClassAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["extension_class"] = extensionClassAttr
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *requestCriteriaResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var planModel, configModel requestCriteriaResourceModel
	req.Config.Get(ctx, &configModel)
	req.Plan.Get(ctx, &planModel)
	resourceType := planModel.Type.ValueString()
	anyDefaultsSet := false
	// Set defaults for root-dse type
	if resourceType == "root-dse" {
		if !internaltypes.IsDefined(configModel.OperationType) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("compare")})
			if !planModel.OperationType.Equal(defaultVal) {
				planModel.OperationType = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for simple type
	if resourceType == "simple" {
		if !internaltypes.IsDefined(configModel.OperationType) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("abandon"), types.StringValue("add"), types.StringValue("bind"), types.StringValue("compare"), types.StringValue("delete"), types.StringValue("extended"), types.StringValue("modify"), types.StringValue("modify-dn"), types.StringValue("search"), types.StringValue("unbind")})
			if !planModel.OperationType.Equal(defaultVal) {
				planModel.OperationType = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludedSearchScope) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("base-object"), types.StringValue("single-level"), types.StringValue("whole-subtree"), types.StringValue("subordinate-subtree")})
			if !planModel.IncludedSearchScope.Equal(defaultVal) {
				planModel.IncludedSearchScope = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.UsingAdministrativeSessionWorkerThread) {
			defaultVal := types.StringValue("any")
			if !planModel.UsingAdministrativeSessionWorkerThread.Equal(defaultVal) {
				planModel.UsingAdministrativeSessionWorkerThread = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	if anyDefaultsSet {
		planModel.Notifications = types.SetUnknown(types.StringType)
		planModel.RequiredActions = types.SetUnknown(config.GetRequiredActionsObjectType())
	}
	planModel.setNotApplicableAttrsNull()
	resp.Plan.Set(ctx, &planModel)
}

func (model *requestCriteriaResourceModel) setNotApplicableAttrsNull() {
	resourceType := model.Type.ValueString()
	// Set any not applicable computed attributes to null for each type
	if resourceType == "root-dse" {
		model.UsingAdministrativeSessionWorkerThread = types.StringNull()
		model.IncludedSearchScope, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "aggregate" {
		model.UsingAdministrativeSessionWorkerThread = types.StringNull()
		model.IncludedSearchScope, _ = types.SetValue(types.StringType, []attr.Value{})
		model.OperationType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "third-party" {
		model.UsingAdministrativeSessionWorkerThread = types.StringNull()
		model.IncludedSearchScope, _ = types.SetValue(types.StringType, []attr.Value{})
		model.OperationType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsRequestCriteria() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"simple"},
			resourcevalidator.Conflicting(
				path.MatchRoot("included_application_name"),
				path.MatchRoot("excluded_application_name"),
			),
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("operation_type"),
			path.MatchRoot("type"),
			[]string{"root-dse", "simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("operation_origin"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("connection_criteria"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("all_included_request_control"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("any_included_request_control"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("not_all_included_request_control"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("none_included_request_control"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("included_target_entry_dn"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("excluded_target_entry_dn"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("all_included_target_entry_filter"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("any_included_target_entry_filter"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("not_all_included_target_entry_filter"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("none_included_target_entry_filter"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("all_included_target_entry_group_dn"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("any_included_target_entry_group_dn"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("not_all_included_target_entry_group_dn"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("none_included_target_entry_group_dn"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("target_bind_type"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("included_target_sasl_mechanism"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("excluded_target_sasl_mechanism"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("included_target_attribute"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("excluded_target_attribute"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("included_extended_operation_oid"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("excluded_extended_operation_oid"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("included_search_scope"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("using_administrative_session_worker_thread"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("included_application_name"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("excluded_application_name"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("all_included_request_criteria"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("any_included_request_criteria"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("not_all_included_request_criteria"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("none_included_request_criteria"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_class"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_argument"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party",
			[]path.Expression{path.MatchRoot("extension_class")},
		),
	}
}

// Add config validators
func (r requestCriteriaResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsRequestCriteria()
}

// Add config validators
func (r defaultRequestCriteriaResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsRequestCriteria()
}

// Add optional fields to create request for root-dse request-criteria
func addOptionalRootDseRequestCriteriaFields(ctx context.Context, addRequest *client.AddRootDseRequestCriteriaRequest, plan requestCriteriaResourceModel) error {
	if internaltypes.IsDefined(plan.OperationType) {
		var slice []string
		plan.OperationType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumrequestCriteriaRootDseOperationTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumrequestCriteriaRootDseOperationTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.OperationType = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for simple request-criteria
func addOptionalSimpleRequestCriteriaFields(ctx context.Context, addRequest *client.AddSimpleRequestCriteriaRequest, plan requestCriteriaResourceModel) error {
	if internaltypes.IsDefined(plan.OperationType) {
		var slice []string
		plan.OperationType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumrequestCriteriaSimpleOperationTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumrequestCriteriaSimpleOperationTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.OperationType = enumSlice
	}
	if internaltypes.IsDefined(plan.OperationOrigin) {
		var slice []string
		plan.OperationOrigin.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumrequestCriteriaOperationOriginProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumrequestCriteriaOperationOriginPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.OperationOrigin = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		addRequest.ConnectionCriteria = plan.ConnectionCriteria.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AllIncludedRequestControl) {
		var slice []string
		plan.AllIncludedRequestControl.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedRequestControl = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedRequestControl) {
		var slice []string
		plan.AnyIncludedRequestControl.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedRequestControl = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedRequestControl) {
		var slice []string
		plan.NotAllIncludedRequestControl.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedRequestControl = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedRequestControl) {
		var slice []string
		plan.NoneIncludedRequestControl.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedRequestControl = slice
	}
	if internaltypes.IsDefined(plan.IncludedTargetEntryDN) {
		var slice []string
		plan.IncludedTargetEntryDN.ElementsAs(ctx, &slice, false)
		addRequest.IncludedTargetEntryDN = slice
	}
	if internaltypes.IsDefined(plan.ExcludedTargetEntryDN) {
		var slice []string
		plan.ExcludedTargetEntryDN.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedTargetEntryDN = slice
	}
	if internaltypes.IsDefined(plan.AllIncludedTargetEntryFilter) {
		var slice []string
		plan.AllIncludedTargetEntryFilter.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedTargetEntryFilter = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedTargetEntryFilter) {
		var slice []string
		plan.AnyIncludedTargetEntryFilter.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedTargetEntryFilter = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedTargetEntryFilter) {
		var slice []string
		plan.NotAllIncludedTargetEntryFilter.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedTargetEntryFilter = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedTargetEntryFilter) {
		var slice []string
		plan.NoneIncludedTargetEntryFilter.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedTargetEntryFilter = slice
	}
	if internaltypes.IsDefined(plan.AllIncludedTargetEntryGroupDN) {
		var slice []string
		plan.AllIncludedTargetEntryGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedTargetEntryGroupDN = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedTargetEntryGroupDN) {
		var slice []string
		plan.AnyIncludedTargetEntryGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedTargetEntryGroupDN = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedTargetEntryGroupDN) {
		var slice []string
		plan.NotAllIncludedTargetEntryGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedTargetEntryGroupDN = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedTargetEntryGroupDN) {
		var slice []string
		plan.NoneIncludedTargetEntryGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedTargetEntryGroupDN = slice
	}
	if internaltypes.IsDefined(plan.TargetBindType) {
		var slice []string
		plan.TargetBindType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumrequestCriteriaTargetBindTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumrequestCriteriaTargetBindTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.TargetBindType = enumSlice
	}
	if internaltypes.IsDefined(plan.IncludedTargetSASLMechanism) {
		var slice []string
		plan.IncludedTargetSASLMechanism.ElementsAs(ctx, &slice, false)
		addRequest.IncludedTargetSASLMechanism = slice
	}
	if internaltypes.IsDefined(plan.ExcludedTargetSASLMechanism) {
		var slice []string
		plan.ExcludedTargetSASLMechanism.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedTargetSASLMechanism = slice
	}
	if internaltypes.IsDefined(plan.IncludedTargetAttribute) {
		var slice []string
		plan.IncludedTargetAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludedTargetAttribute = slice
	}
	if internaltypes.IsDefined(plan.ExcludedTargetAttribute) {
		var slice []string
		plan.ExcludedTargetAttribute.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedTargetAttribute = slice
	}
	if internaltypes.IsDefined(plan.IncludedExtendedOperationOID) {
		var slice []string
		plan.IncludedExtendedOperationOID.ElementsAs(ctx, &slice, false)
		addRequest.IncludedExtendedOperationOID = slice
	}
	if internaltypes.IsDefined(plan.ExcludedExtendedOperationOID) {
		var slice []string
		plan.ExcludedExtendedOperationOID.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedExtendedOperationOID = slice
	}
	if internaltypes.IsDefined(plan.IncludedSearchScope) {
		var slice []string
		plan.IncludedSearchScope.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumrequestCriteriaIncludedSearchScopeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumrequestCriteriaIncludedSearchScopePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.IncludedSearchScope = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UsingAdministrativeSessionWorkerThread) {
		usingAdministrativeSessionWorkerThread, err := client.NewEnumrequestCriteriaUsingAdministrativeSessionWorkerThreadPropFromValue(plan.UsingAdministrativeSessionWorkerThread.ValueString())
		if err != nil {
			return err
		}
		addRequest.UsingAdministrativeSessionWorkerThread = usingAdministrativeSessionWorkerThread
	}
	if internaltypes.IsDefined(plan.IncludedApplicationName) {
		var slice []string
		plan.IncludedApplicationName.ElementsAs(ctx, &slice, false)
		addRequest.IncludedApplicationName = slice
	}
	if internaltypes.IsDefined(plan.ExcludedApplicationName) {
		var slice []string
		plan.ExcludedApplicationName.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedApplicationName = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for aggregate request-criteria
func addOptionalAggregateRequestCriteriaFields(ctx context.Context, addRequest *client.AddAggregateRequestCriteriaRequest, plan requestCriteriaResourceModel) error {
	if internaltypes.IsDefined(plan.AllIncludedRequestCriteria) {
		var slice []string
		plan.AllIncludedRequestCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedRequestCriteria = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedRequestCriteria) {
		var slice []string
		plan.AnyIncludedRequestCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedRequestCriteria = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedRequestCriteria) {
		var slice []string
		plan.NotAllIncludedRequestCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedRequestCriteria = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedRequestCriteria) {
		var slice []string
		plan.NoneIncludedRequestCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedRequestCriteria = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for third-party request-criteria
func addOptionalThirdPartyRequestCriteriaFields(ctx context.Context, addRequest *client.AddThirdPartyRequestCriteriaRequest, plan requestCriteriaResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateRequestCriteriaUnknownValues(model *requestCriteriaResourceModel) {
	if model.AnyIncludedRequestControl.IsUnknown() || model.AnyIncludedRequestControl.IsNull() {
		model.AnyIncludedRequestControl, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludedSearchScope.IsUnknown() || model.IncludedSearchScope.IsNull() {
		model.IncludedSearchScope, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.TargetBindType.IsUnknown() || model.TargetBindType.IsNull() {
		model.TargetBindType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.NoneIncludedTargetEntryGroupDN.IsUnknown() || model.NoneIncludedTargetEntryGroupDN.IsNull() {
		model.NoneIncludedTargetEntryGroupDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AnyIncludedRequestCriteria.IsUnknown() || model.AnyIncludedRequestCriteria.IsNull() {
		model.AnyIncludedRequestCriteria, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExcludedTargetEntryDN.IsUnknown() || model.ExcludedTargetEntryDN.IsNull() {
		model.ExcludedTargetEntryDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludedTargetEntryDN.IsUnknown() || model.IncludedTargetEntryDN.IsNull() {
		model.IncludedTargetEntryDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.NoneIncludedRequestCriteria.IsUnknown() || model.NoneIncludedRequestCriteria.IsNull() {
		model.NoneIncludedRequestCriteria, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.NotAllIncludedRequestControl.IsUnknown() || model.NotAllIncludedRequestControl.IsNull() {
		model.NotAllIncludedRequestControl, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AllIncludedRequestControl.IsUnknown() || model.AllIncludedRequestControl.IsNull() {
		model.AllIncludedRequestControl, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AnyIncludedTargetEntryGroupDN.IsUnknown() || model.AnyIncludedTargetEntryGroupDN.IsNull() {
		model.AnyIncludedTargetEntryGroupDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludedApplicationName.IsUnknown() || model.IncludedApplicationName.IsNull() {
		model.IncludedApplicationName, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.NotAllIncludedTargetEntryGroupDN.IsUnknown() || model.NotAllIncludedTargetEntryGroupDN.IsNull() {
		model.NotAllIncludedTargetEntryGroupDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludedTargetSASLMechanism.IsUnknown() || model.IncludedTargetSASLMechanism.IsNull() {
		model.IncludedTargetSASLMechanism, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AllIncludedRequestCriteria.IsUnknown() || model.AllIncludedRequestCriteria.IsNull() {
		model.AllIncludedRequestCriteria, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.NotAllIncludedRequestCriteria.IsUnknown() || model.NotAllIncludedRequestCriteria.IsNull() {
		model.NotAllIncludedRequestCriteria, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExcludedTargetAttribute.IsUnknown() || model.ExcludedTargetAttribute.IsNull() {
		model.ExcludedTargetAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExcludedTargetSASLMechanism.IsUnknown() || model.ExcludedTargetSASLMechanism.IsNull() {
		model.ExcludedTargetSASLMechanism, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AllIncludedTargetEntryFilter.IsUnknown() || model.AllIncludedTargetEntryFilter.IsNull() {
		model.AllIncludedTargetEntryFilter, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludedExtendedOperationOID.IsUnknown() || model.IncludedExtendedOperationOID.IsNull() {
		model.IncludedExtendedOperationOID, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.OperationOrigin.IsUnknown() || model.OperationOrigin.IsNull() {
		model.OperationOrigin, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AnyIncludedTargetEntryFilter.IsUnknown() || model.AnyIncludedTargetEntryFilter.IsNull() {
		model.AnyIncludedTargetEntryFilter, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.NoneIncludedTargetEntryFilter.IsUnknown() || model.NoneIncludedTargetEntryFilter.IsNull() {
		model.NoneIncludedTargetEntryFilter, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AllIncludedTargetEntryGroupDN.IsUnknown() || model.AllIncludedTargetEntryGroupDN.IsNull() {
		model.AllIncludedTargetEntryGroupDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludedTargetAttribute.IsUnknown() || model.IncludedTargetAttribute.IsNull() {
		model.IncludedTargetAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.OperationType.IsUnknown() || model.OperationType.IsNull() {
		model.OperationType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExcludedExtendedOperationOID.IsUnknown() || model.ExcludedExtendedOperationOID.IsNull() {
		model.ExcludedExtendedOperationOID, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExcludedApplicationName.IsUnknown() || model.ExcludedApplicationName.IsNull() {
		model.ExcludedApplicationName, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.NotAllIncludedTargetEntryFilter.IsUnknown() || model.NotAllIncludedTargetEntryFilter.IsNull() {
		model.NotAllIncludedTargetEntryFilter, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.NoneIncludedRequestControl.IsUnknown() || model.NoneIncludedRequestControl.IsNull() {
		model.NoneIncludedRequestControl, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *requestCriteriaResourceModel) populateAllComputedStringAttributes() {
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.ExtensionClass.IsUnknown() || model.ExtensionClass.IsNull() {
		model.ExtensionClass = types.StringValue("")
	}
	if model.UsingAdministrativeSessionWorkerThread.IsUnknown() || model.UsingAdministrativeSessionWorkerThread.IsNull() {
		model.UsingAdministrativeSessionWorkerThread = types.StringValue("")
	}
	if model.ConnectionCriteria.IsUnknown() || model.ConnectionCriteria.IsNull() {
		model.ConnectionCriteria = types.StringValue("")
	}
}

// Read a RootDseRequestCriteriaResponse object into the model struct
func readRootDseRequestCriteriaResponse(ctx context.Context, r *client.RootDseRequestCriteriaResponse, state *requestCriteriaResourceModel, expectedValues *requestCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("root-dse")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.OperationType = internaltypes.GetStringSet(
		client.StringSliceEnumrequestCriteriaRootDseOperationTypeProp(r.OperationType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateRequestCriteriaUnknownValues(state)
}

// Read a SimpleRequestCriteriaResponse object into the model struct
func readSimpleRequestCriteriaResponse(ctx context.Context, r *client.SimpleRequestCriteriaResponse, state *requestCriteriaResourceModel, expectedValues *requestCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("simple")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.OperationType = internaltypes.GetStringSet(
		client.StringSliceEnumrequestCriteriaSimpleOperationTypeProp(r.OperationType))
	state.OperationOrigin = internaltypes.GetStringSet(
		client.StringSliceEnumrequestCriteriaOperationOriginProp(r.OperationOrigin))
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.AllIncludedRequestControl = internaltypes.GetStringSet(r.AllIncludedRequestControl)
	state.AnyIncludedRequestControl = internaltypes.GetStringSet(r.AnyIncludedRequestControl)
	state.NotAllIncludedRequestControl = internaltypes.GetStringSet(r.NotAllIncludedRequestControl)
	state.NoneIncludedRequestControl = internaltypes.GetStringSet(r.NoneIncludedRequestControl)
	state.IncludedTargetEntryDN = internaltypes.GetStringSet(r.IncludedTargetEntryDN)
	state.ExcludedTargetEntryDN = internaltypes.GetStringSet(r.ExcludedTargetEntryDN)
	state.AllIncludedTargetEntryFilter = internaltypes.GetStringSet(r.AllIncludedTargetEntryFilter)
	state.AnyIncludedTargetEntryFilter = internaltypes.GetStringSet(r.AnyIncludedTargetEntryFilter)
	state.NotAllIncludedTargetEntryFilter = internaltypes.GetStringSet(r.NotAllIncludedTargetEntryFilter)
	state.NoneIncludedTargetEntryFilter = internaltypes.GetStringSet(r.NoneIncludedTargetEntryFilter)
	state.AllIncludedTargetEntryGroupDN = internaltypes.GetStringSet(r.AllIncludedTargetEntryGroupDN)
	state.AnyIncludedTargetEntryGroupDN = internaltypes.GetStringSet(r.AnyIncludedTargetEntryGroupDN)
	state.NotAllIncludedTargetEntryGroupDN = internaltypes.GetStringSet(r.NotAllIncludedTargetEntryGroupDN)
	state.NoneIncludedTargetEntryGroupDN = internaltypes.GetStringSet(r.NoneIncludedTargetEntryGroupDN)
	state.TargetBindType = internaltypes.GetStringSet(
		client.StringSliceEnumrequestCriteriaTargetBindTypeProp(r.TargetBindType))
	state.IncludedTargetSASLMechanism = internaltypes.GetStringSet(r.IncludedTargetSASLMechanism)
	state.ExcludedTargetSASLMechanism = internaltypes.GetStringSet(r.ExcludedTargetSASLMechanism)
	state.IncludedTargetAttribute = internaltypes.GetStringSet(r.IncludedTargetAttribute)
	state.ExcludedTargetAttribute = internaltypes.GetStringSet(r.ExcludedTargetAttribute)
	state.IncludedExtendedOperationOID = internaltypes.GetStringSet(r.IncludedExtendedOperationOID)
	state.ExcludedExtendedOperationOID = internaltypes.GetStringSet(r.ExcludedExtendedOperationOID)
	state.IncludedSearchScope = internaltypes.GetStringSet(
		client.StringSliceEnumrequestCriteriaIncludedSearchScopeProp(r.IncludedSearchScope))
	state.UsingAdministrativeSessionWorkerThread = internaltypes.StringTypeOrNil(
		client.StringPointerEnumrequestCriteriaUsingAdministrativeSessionWorkerThreadProp(r.UsingAdministrativeSessionWorkerThread), true)
	state.IncludedApplicationName = internaltypes.GetStringSet(r.IncludedApplicationName)
	state.ExcludedApplicationName = internaltypes.GetStringSet(r.ExcludedApplicationName)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateRequestCriteriaUnknownValues(state)
}

// Read a AggregateRequestCriteriaResponse object into the model struct
func readAggregateRequestCriteriaResponse(ctx context.Context, r *client.AggregateRequestCriteriaResponse, state *requestCriteriaResourceModel, expectedValues *requestCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("aggregate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllIncludedRequestCriteria = internaltypes.GetStringSet(r.AllIncludedRequestCriteria)
	state.AnyIncludedRequestCriteria = internaltypes.GetStringSet(r.AnyIncludedRequestCriteria)
	state.NotAllIncludedRequestCriteria = internaltypes.GetStringSet(r.NotAllIncludedRequestCriteria)
	state.NoneIncludedRequestCriteria = internaltypes.GetStringSet(r.NoneIncludedRequestCriteria)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateRequestCriteriaUnknownValues(state)
}

// Read a ThirdPartyRequestCriteriaResponse object into the model struct
func readThirdPartyRequestCriteriaResponse(ctx context.Context, r *client.ThirdPartyRequestCriteriaResponse, state *requestCriteriaResourceModel, expectedValues *requestCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateRequestCriteriaUnknownValues(state)
}

// Create any update operations necessary to make the state match the plan
func createRequestCriteriaOperations(plan requestCriteriaResourceModel, state requestCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedRequestCriteria, state.AllIncludedRequestCriteria, "all-included-request-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedRequestCriteria, state.AnyIncludedRequestCriteria, "any-included-request-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedRequestCriteria, state.NotAllIncludedRequestCriteria, "not-all-included-request-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedRequestCriteria, state.NoneIncludedRequestCriteria, "none-included-request-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.OperationType, state.OperationType, "operation-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.OperationOrigin, state.OperationOrigin, "operation-origin")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionCriteria, state.ConnectionCriteria, "connection-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedRequestControl, state.AllIncludedRequestControl, "all-included-request-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedRequestControl, state.AnyIncludedRequestControl, "any-included-request-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedRequestControl, state.NotAllIncludedRequestControl, "not-all-included-request-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedRequestControl, state.NoneIncludedRequestControl, "none-included-request-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedTargetEntryDN, state.IncludedTargetEntryDN, "included-target-entry-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedTargetEntryDN, state.ExcludedTargetEntryDN, "excluded-target-entry-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedTargetEntryFilter, state.AllIncludedTargetEntryFilter, "all-included-target-entry-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedTargetEntryFilter, state.AnyIncludedTargetEntryFilter, "any-included-target-entry-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedTargetEntryFilter, state.NotAllIncludedTargetEntryFilter, "not-all-included-target-entry-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedTargetEntryFilter, state.NoneIncludedTargetEntryFilter, "none-included-target-entry-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedTargetEntryGroupDN, state.AllIncludedTargetEntryGroupDN, "all-included-target-entry-group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedTargetEntryGroupDN, state.AnyIncludedTargetEntryGroupDN, "any-included-target-entry-group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedTargetEntryGroupDN, state.NotAllIncludedTargetEntryGroupDN, "not-all-included-target-entry-group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedTargetEntryGroupDN, state.NoneIncludedTargetEntryGroupDN, "none-included-target-entry-group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.TargetBindType, state.TargetBindType, "target-bind-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedTargetSASLMechanism, state.IncludedTargetSASLMechanism, "included-target-sasl-mechanism")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedTargetSASLMechanism, state.ExcludedTargetSASLMechanism, "excluded-target-sasl-mechanism")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedTargetAttribute, state.IncludedTargetAttribute, "included-target-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedTargetAttribute, state.ExcludedTargetAttribute, "excluded-target-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedExtendedOperationOID, state.IncludedExtendedOperationOID, "included-extended-operation-oid")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedExtendedOperationOID, state.ExcludedExtendedOperationOID, "excluded-extended-operation-oid")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedSearchScope, state.IncludedSearchScope, "included-search-scope")
	operations.AddStringOperationIfNecessary(&ops, plan.UsingAdministrativeSessionWorkerThread, state.UsingAdministrativeSessionWorkerThread, "using-administrative-session-worker-thread")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedApplicationName, state.IncludedApplicationName, "included-application-name")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedApplicationName, state.ExcludedApplicationName, "excluded-application-name")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a root-dse request-criteria
func (r *requestCriteriaResource) CreateRootDseRequestCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan requestCriteriaResourceModel) (*requestCriteriaResourceModel, error) {
	addRequest := client.NewAddRootDseRequestCriteriaRequest([]client.EnumrootDseRequestCriteriaSchemaUrn{client.ENUMROOTDSEREQUESTCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0REQUEST_CRITERIAROOT_DSE},
		plan.Name.ValueString())
	err := addOptionalRootDseRequestCriteriaFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Request Criteria", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RequestCriteriaAPI.AddRequestCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRequestCriteriaRequest(
		client.AddRootDseRequestCriteriaRequestAsAddRequestCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RequestCriteriaAPI.AddRequestCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Request Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state requestCriteriaResourceModel
	readRootDseRequestCriteriaResponse(ctx, addResponse.RootDseRequestCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a simple request-criteria
func (r *requestCriteriaResource) CreateSimpleRequestCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan requestCriteriaResourceModel) (*requestCriteriaResourceModel, error) {
	addRequest := client.NewAddSimpleRequestCriteriaRequest([]client.EnumsimpleRequestCriteriaSchemaUrn{client.ENUMSIMPLEREQUESTCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0REQUEST_CRITERIASIMPLE},
		plan.Name.ValueString())
	err := addOptionalSimpleRequestCriteriaFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Request Criteria", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RequestCriteriaAPI.AddRequestCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRequestCriteriaRequest(
		client.AddSimpleRequestCriteriaRequestAsAddRequestCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RequestCriteriaAPI.AddRequestCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Request Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state requestCriteriaResourceModel
	readSimpleRequestCriteriaResponse(ctx, addResponse.SimpleRequestCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a aggregate request-criteria
func (r *requestCriteriaResource) CreateAggregateRequestCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan requestCriteriaResourceModel) (*requestCriteriaResourceModel, error) {
	addRequest := client.NewAddAggregateRequestCriteriaRequest([]client.EnumaggregateRequestCriteriaSchemaUrn{client.ENUMAGGREGATEREQUESTCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0REQUEST_CRITERIAAGGREGATE},
		plan.Name.ValueString())
	err := addOptionalAggregateRequestCriteriaFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Request Criteria", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RequestCriteriaAPI.AddRequestCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRequestCriteriaRequest(
		client.AddAggregateRequestCriteriaRequestAsAddRequestCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RequestCriteriaAPI.AddRequestCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Request Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state requestCriteriaResourceModel
	readAggregateRequestCriteriaResponse(ctx, addResponse.AggregateRequestCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party request-criteria
func (r *requestCriteriaResource) CreateThirdPartyRequestCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan requestCriteriaResourceModel) (*requestCriteriaResourceModel, error) {
	addRequest := client.NewAddThirdPartyRequestCriteriaRequest([]client.EnumthirdPartyRequestCriteriaSchemaUrn{client.ENUMTHIRDPARTYREQUESTCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0REQUEST_CRITERIATHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Name.ValueString())
	err := addOptionalThirdPartyRequestCriteriaFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Request Criteria", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RequestCriteriaAPI.AddRequestCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRequestCriteriaRequest(
		client.AddThirdPartyRequestCriteriaRequestAsAddRequestCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RequestCriteriaAPI.AddRequestCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Request Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state requestCriteriaResourceModel
	readThirdPartyRequestCriteriaResponse(ctx, addResponse.ThirdPartyRequestCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *requestCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan requestCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *requestCriteriaResourceModel
	var err error
	if plan.Type.ValueString() == "root-dse" {
		state, err = r.CreateRootDseRequestCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "simple" {
		state, err = r.CreateSimpleRequestCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "aggregate" {
		state, err = r.CreateAggregateRequestCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyRequestCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
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
func (r *defaultRequestCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan requestCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RequestCriteriaAPI.GetRequestCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Request Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state requestCriteriaResourceModel
	if readResponse.RootDseRequestCriteriaResponse != nil {
		readRootDseRequestCriteriaResponse(ctx, readResponse.RootDseRequestCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SimpleRequestCriteriaResponse != nil {
		readSimpleRequestCriteriaResponse(ctx, readResponse.SimpleRequestCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AggregateRequestCriteriaResponse != nil {
		readAggregateRequestCriteriaResponse(ctx, readResponse.AggregateRequestCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyRequestCriteriaResponse != nil {
		readThirdPartyRequestCriteriaResponse(ctx, readResponse.ThirdPartyRequestCriteriaResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.RequestCriteriaAPI.UpdateRequestCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createRequestCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.RequestCriteriaAPI.UpdateRequestCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Request Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.RootDseRequestCriteriaResponse != nil {
			readRootDseRequestCriteriaResponse(ctx, updateResponse.RootDseRequestCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SimpleRequestCriteriaResponse != nil {
			readSimpleRequestCriteriaResponse(ctx, updateResponse.SimpleRequestCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AggregateRequestCriteriaResponse != nil {
			readAggregateRequestCriteriaResponse(ctx, updateResponse.AggregateRequestCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyRequestCriteriaResponse != nil {
			readThirdPartyRequestCriteriaResponse(ctx, updateResponse.ThirdPartyRequestCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
	}

	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *requestCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readRequestCriteria(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultRequestCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readRequestCriteria(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readRequestCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state requestCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.RequestCriteriaAPI.GetRequestCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Request Criteria", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Request Criteria", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.RootDseRequestCriteriaResponse != nil {
		readRootDseRequestCriteriaResponse(ctx, readResponse.RootDseRequestCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SimpleRequestCriteriaResponse != nil {
		readSimpleRequestCriteriaResponse(ctx, readResponse.SimpleRequestCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AggregateRequestCriteriaResponse != nil {
		readAggregateRequestCriteriaResponse(ctx, readResponse.AggregateRequestCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyRequestCriteriaResponse != nil {
		readThirdPartyRequestCriteriaResponse(ctx, readResponse.ThirdPartyRequestCriteriaResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *requestCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateRequestCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultRequestCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateRequestCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateRequestCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan requestCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state requestCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.RequestCriteriaAPI.UpdateRequestCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createRequestCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.RequestCriteriaAPI.UpdateRequestCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Request Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.RootDseRequestCriteriaResponse != nil {
			readRootDseRequestCriteriaResponse(ctx, updateResponse.RootDseRequestCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SimpleRequestCriteriaResponse != nil {
			readSimpleRequestCriteriaResponse(ctx, updateResponse.SimpleRequestCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AggregateRequestCriteriaResponse != nil {
			readAggregateRequestCriteriaResponse(ctx, updateResponse.AggregateRequestCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyRequestCriteriaResponse != nil {
			readThirdPartyRequestCriteriaResponse(ctx, updateResponse.ThirdPartyRequestCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *defaultRequestCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *requestCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state requestCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.RequestCriteriaAPI.DeleteRequestCriteriaExecute(r.apiClient.RequestCriteriaAPI.DeleteRequestCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Request Criteria", err, httpResp)
		return
	}
}

func (r *requestCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importRequestCriteria(ctx, req, resp)
}

func (r *defaultRequestCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importRequestCriteria(ctx, req, resp)
}

func importRequestCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
