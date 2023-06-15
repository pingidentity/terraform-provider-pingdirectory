package connectioncriteria

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &connectionCriteriaResource{}
	_ resource.ResourceWithConfigure   = &connectionCriteriaResource{}
	_ resource.ResourceWithImportState = &connectionCriteriaResource{}
	_ resource.Resource                = &defaultConnectionCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultConnectionCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultConnectionCriteriaResource{}
)

// Create a Connection Criteria resource
func NewConnectionCriteriaResource() resource.Resource {
	return &connectionCriteriaResource{}
}

func NewDefaultConnectionCriteriaResource() resource.Resource {
	return &defaultConnectionCriteriaResource{}
}

// connectionCriteriaResource is the resource implementation.
type connectionCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultConnectionCriteriaResource is the resource implementation.
type defaultConnectionCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *connectionCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connection_criteria"
}

func (r *defaultConnectionCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_connection_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *connectionCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultConnectionCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type connectionCriteriaResourceModel struct {
	Id                               types.String `tfsdk:"id"`
	LastUpdated                      types.String `tfsdk:"last_updated"`
	Notifications                    types.Set    `tfsdk:"notifications"`
	RequiredActions                  types.Set    `tfsdk:"required_actions"`
	Type                             types.String `tfsdk:"type"`
	ExtensionClass                   types.String `tfsdk:"extension_class"`
	ExtensionArgument                types.Set    `tfsdk:"extension_argument"`
	AllIncludedConnectionCriteria    types.Set    `tfsdk:"all_included_connection_criteria"`
	AnyIncludedConnectionCriteria    types.Set    `tfsdk:"any_included_connection_criteria"`
	NotAllIncludedConnectionCriteria types.Set    `tfsdk:"not_all_included_connection_criteria"`
	NoneIncludedConnectionCriteria   types.Set    `tfsdk:"none_included_connection_criteria"`
	IncludedClientAddress            types.Set    `tfsdk:"included_client_address"`
	ExcludedClientAddress            types.Set    `tfsdk:"excluded_client_address"`
	IncludedConnectionHandler        types.Set    `tfsdk:"included_connection_handler"`
	ExcludedConnectionHandler        types.Set    `tfsdk:"excluded_connection_handler"`
	IncludedProtocol                 types.Set    `tfsdk:"included_protocol"`
	ExcludedProtocol                 types.Set    `tfsdk:"excluded_protocol"`
	CommunicationSecurityLevel       types.String `tfsdk:"communication_security_level"`
	UserAuthType                     types.Set    `tfsdk:"user_auth_type"`
	AuthenticationSecurityLevel      types.String `tfsdk:"authentication_security_level"`
	IncludedUserSASLMechanism        types.Set    `tfsdk:"included_user_sasl_mechanism"`
	ExcludedUserSASLMechanism        types.Set    `tfsdk:"excluded_user_sasl_mechanism"`
	IncludedUserBaseDN               types.Set    `tfsdk:"included_user_base_dn"`
	ExcludedUserBaseDN               types.Set    `tfsdk:"excluded_user_base_dn"`
	AllIncludedUserGroupDN           types.Set    `tfsdk:"all_included_user_group_dn"`
	AnyIncludedUserGroupDN           types.Set    `tfsdk:"any_included_user_group_dn"`
	NotAllIncludedUserGroupDN        types.Set    `tfsdk:"not_all_included_user_group_dn"`
	NoneIncludedUserGroupDN          types.Set    `tfsdk:"none_included_user_group_dn"`
	AllIncludedUserFilter            types.Set    `tfsdk:"all_included_user_filter"`
	AnyIncludedUserFilter            types.Set    `tfsdk:"any_included_user_filter"`
	NotAllIncludedUserFilter         types.Set    `tfsdk:"not_all_included_user_filter"`
	NoneIncludedUserFilter           types.Set    `tfsdk:"none_included_user_filter"`
	AllIncludedUserPrivilege         types.Set    `tfsdk:"all_included_user_privilege"`
	AnyIncludedUserPrivilege         types.Set    `tfsdk:"any_included_user_privilege"`
	NotAllIncludedUserPrivilege      types.Set    `tfsdk:"not_all_included_user_privilege"`
	NoneIncludedUserPrivilege        types.Set    `tfsdk:"none_included_user_privilege"`
	Description                      types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *connectionCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	connectionCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultConnectionCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	connectionCriteriaSchema(ctx, req, resp, true)
}

func connectionCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Connection Criteria.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Connection Criteria resource. Options are ['simple', 'aggregate', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"simple", "aggregate", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Connection Criteria.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Connection Criteria. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"all_included_connection_criteria": schema.SetAttribute{
				Description: "Specifies a connection criteria object that must match the associated client connection in order to match the aggregate connection criteria. If one or more all-included connection criteria objects are provided, then a client connection must match all of them in order to match the aggregate connection criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"any_included_connection_criteria": schema.SetAttribute{
				Description: "Specifies a connection criteria object that may match the associated client connection in order to match the aggregate connection criteria. If one or more any-included connection criteria objects are provided, then a client connection must match at least one of them in order to match the aggregate connection criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"not_all_included_connection_criteria": schema.SetAttribute{
				Description: "Specifies a connection criteria object that should not match the associated client connection in order to match the aggregate connection criteria. If one or more not-all-included connection criteria objects are provided, then a client connection must not match all of them (that is, it may match zero or more of them, but it must not match all of them) in order to match the aggregate connection criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"none_included_connection_criteria": schema.SetAttribute{
				Description: "Specifies a connection criteria object that must not match the associated client connection in order to match the aggregate connection criteria. If one or more none-included connection criteria objects are provided, then a client connection must not match any of them in order to match the aggregate connection criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"included_client_address": schema.SetAttribute{
				Description: "Specifies an address mask that may be used to specify a set of clients that should be included in this Simple Connection Criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"excluded_client_address": schema.SetAttribute{
				Description: "Specifies an address mask that may be used to specify a set of clients that should be excluded from this Simple Connection Criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"included_connection_handler": schema.SetAttribute{
				Description: "Specifies a connection handler for clients that should be included in this Simple Connection Criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"excluded_connection_handler": schema.SetAttribute{
				Description: "Specifies a connection handler for clients that should be excluded from this Simple Connection Criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"included_protocol": schema.SetAttribute{
				Description: "Specifies the name of a communication protocol that should be used by clients included in this Simple Connection Criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"excluded_protocol": schema.SetAttribute{
				Description: "Specifies the name of a communication protocol that should be used by clients excluded from this Simple Connection Criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"communication_security_level": schema.StringAttribute{
				Description: "Indicates whether this Simple Connection Criteria should require or allow clients using a secure communication channel.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_auth_type": schema.SetAttribute{
				Description: "Specifies the authentication types for client connections that may be included in this Simple Connection Criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"authentication_security_level": schema.StringAttribute{
				Description: "Indicates whether this Simple Connection Criteria should require or allow clients that authenticated using a secure manner. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"included_user_sasl_mechanism": schema.SetAttribute{
				Description: "Specifies the name of a SASL mechanism that should be used by clients included in this Simple Connection Criteria. This will only be taken into account for client connections that have authenticated to the server using a SASL mechanism and will be ignored for unauthenticated client connections and for client connections that authenticated using some other method (e.g., those performing simple or internal authentication).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"excluded_user_sasl_mechanism": schema.SetAttribute{
				Description: "Specifies the name of a SASL mechanism that should be used by clients excluded from this Simple Connection Criteria. This will only be taken into account for client connections that have authenticated to the server using a SASL mechanism and will be ignored for unauthenticated client connections and for client connections that authenticated using some other method (e.g., those performing simple or internal authentication).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"included_user_base_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which authenticated user entries may exist for clients included in this Simple Connection Criteria. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections. Refer to the authz version of this property in Simple Result Criteria if operations are being proxied (performed using proxied authorization), and you need to match the originating user of the operation rather than the proxy user (the user the proxy authenticated as).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"excluded_user_base_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which authenticated user entries may exist for clients excluded from this Simple Connection Criteria. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections. Refer to the authz version of this property in Simple Result Criteria if operations are being proxied (performed using proxied authorization), and you need to match the originating user of the operation rather than the proxy user (the user the proxy authenticated as).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"all_included_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authenticated users must exist for clients included in this Simple Connection Criteria. If any group DNs are provided, then the authenticated user must be a member of all of those groups. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections. Refer to the authz version of this property in Simple Result Criteria if operations are being proxied (performed using proxied authorization), and you need to match the originating user of the operation rather than the proxy user (the user the proxy authenticated as).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"any_included_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authenticated users may exist for clients included in this Simple Connection Criteria. If any group DNs are provided, then the authenticated user must be a member of at least one of those groups. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections. Refer to the authz version of this property in Simple Result Criteria if operations are being proxied (performed using proxied authorization), and you need to match the originating user of the operation rather than the proxy user (the user the proxy authenticated as).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"not_all_included_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authenticated users should not exist for clients included in this Simple Connection Criteria. If any group DNs are provided, then the authenticated user must not be a member of at least one of those groups (that is, the user may be a member of zero or more of those groups, but not of all of them). This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections. Refer to the authz version of this property in Simple Result Criteria if operations are being proxied (performed using proxied authorization), and you need to match the originating user of the operation rather than the proxy user (the user the proxy authenticated as).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"none_included_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authenticated users must not exist for clients included in this Simple Connection Criteria. If any group DNs are provided, then the authenticated user must not be a member any of those groups. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections. Refer to the authz version of this property in Simple Result Criteria if operations are being proxied (performed using proxied authorization), and you need to match the originating user of the operation rather than the proxy user (the user the proxy authenticated as).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"all_included_user_filter": schema.SetAttribute{
				Description: "Specifies a search filter that must match the entry of the authenticated user for clients included in this Simple Connection Criteria. If any filters are provided, then all of those filters must match the authenticated user entry. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"any_included_user_filter": schema.SetAttribute{
				Description: "Specifies a search filter that may match the entry of the authenticated user for clients included in this Simple Connection Criteria. If any filters are provided, then at least one of those filters must match the authenticated user entry. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"not_all_included_user_filter": schema.SetAttribute{
				Description: "Specifies a search filter that should not match the entry of the authenticated user for clients included in this Simple Connection Criteria. If any filters are provided, then at least one of those filters must not match the authenticated user entry (that is, the user entry may match zero or more of those filters, but not all of them). This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"none_included_user_filter": schema.SetAttribute{
				Description: "Specifies a search filter that must not match the entry of the authenticated user for clients included in this Simple Connection Criteria. If any filters are provided, then none of those filters may match the authenticated user entry. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"all_included_user_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that must be held by the authenticated user for clients included in this Simple Connection Criteria. If any privilege names are provided, then the authenticated user must have all of those privileges. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"any_included_user_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that may be held by the authenticated user for clients included in this Simple Connection Criteria. If any privilege names are provided, then the authenticated user must have at least one of those privileges. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"not_all_included_user_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that should not be held by the authenticated user for clients included in this Simple Connection Criteria. If any privilege names are provided, then the authenticated user must not have at least one of those privileges (that is, the user may hold zero or more of those privileges, but not all of them). This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"none_included_user_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that must not be held by the authenticated user for clients included in this Simple Connection Criteria. If any privilege names are provided, then the authenticated user must not have any of those privileges. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Connection Criteria",
				Optional:    true,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"simple", "aggregate", "third-party"}...),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAllAttributesToOptionalAndComputed(&schemaDef, []string{"id"})
	}
	config.AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *connectionCriteriaResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanConnectionCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultConnectionCriteriaResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanConnectionCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanConnectionCriteria(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var model connectionCriteriaResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.NoneIncludedConnectionCriteria) && model.Type.ValueString() != "aggregate" {
		resp.Diagnostics.AddError("Attribute 'none_included_connection_criteria' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'none_included_connection_criteria', the 'type' attribute must be one of ['aggregate']")
	}
	if internaltypes.IsDefined(model.ExcludedConnectionHandler) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'excluded_connection_handler' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'excluded_connection_handler', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.AllIncludedUserFilter) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'all_included_user_filter' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'all_included_user_filter', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.AllIncludedUserGroupDN) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'all_included_user_group_dn' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'all_included_user_group_dn', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.AllIncludedUserPrivilege) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'all_included_user_privilege' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'all_included_user_privilege', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.AuthenticationSecurityLevel) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'authentication_security_level' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'authentication_security_level', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.IncludedConnectionHandler) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'included_connection_handler' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'included_connection_handler', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.IncludedProtocol) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'included_protocol' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'included_protocol', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.NoneIncludedUserPrivilege) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'none_included_user_privilege' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'none_included_user_privilege', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.AnyIncludedUserFilter) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'any_included_user_filter' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'any_included_user_filter', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.ExtensionArgument) && model.Type.ValueString() != "third-party" {
		resp.Diagnostics.AddError("Attribute 'extension_argument' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'extension_argument', the 'type' attribute must be one of ['third-party']")
	}
	if internaltypes.IsDefined(model.IncludedUserBaseDN) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'included_user_base_dn' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'included_user_base_dn', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.AnyIncludedUserPrivilege) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'any_included_user_privilege' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'any_included_user_privilege', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.NotAllIncludedConnectionCriteria) && model.Type.ValueString() != "aggregate" {
		resp.Diagnostics.AddError("Attribute 'not_all_included_connection_criteria' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'not_all_included_connection_criteria', the 'type' attribute must be one of ['aggregate']")
	}
	if internaltypes.IsDefined(model.IncludedUserSASLMechanism) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'included_user_sasl_mechanism' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'included_user_sasl_mechanism', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.UserAuthType) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'user_auth_type' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'user_auth_type', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.IncludedClientAddress) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'included_client_address' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'included_client_address', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.NotAllIncludedUserGroupDN) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'not_all_included_user_group_dn' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'not_all_included_user_group_dn', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.AnyIncludedConnectionCriteria) && model.Type.ValueString() != "aggregate" {
		resp.Diagnostics.AddError("Attribute 'any_included_connection_criteria' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'any_included_connection_criteria', the 'type' attribute must be one of ['aggregate']")
	}
	if internaltypes.IsDefined(model.ExcludedProtocol) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'excluded_protocol' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'excluded_protocol', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.AllIncludedConnectionCriteria) && model.Type.ValueString() != "aggregate" {
		resp.Diagnostics.AddError("Attribute 'all_included_connection_criteria' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'all_included_connection_criteria', the 'type' attribute must be one of ['aggregate']")
	}
	if internaltypes.IsDefined(model.ExcludedUserSASLMechanism) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'excluded_user_sasl_mechanism' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'excluded_user_sasl_mechanism', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.ExcludedClientAddress) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'excluded_client_address' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'excluded_client_address', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.NotAllIncludedUserPrivilege) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'not_all_included_user_privilege' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'not_all_included_user_privilege', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.NotAllIncludedUserFilter) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'not_all_included_user_filter' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'not_all_included_user_filter', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.NoneIncludedUserGroupDN) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'none_included_user_group_dn' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'none_included_user_group_dn', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.CommunicationSecurityLevel) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'communication_security_level' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'communication_security_level', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.ExcludedUserBaseDN) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'excluded_user_base_dn' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'excluded_user_base_dn', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.ExtensionClass) && model.Type.ValueString() != "third-party" {
		resp.Diagnostics.AddError("Attribute 'extension_class' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'extension_class', the 'type' attribute must be one of ['third-party']")
	}
	if internaltypes.IsDefined(model.AnyIncludedUserGroupDN) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'any_included_user_group_dn' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'any_included_user_group_dn', the 'type' attribute must be one of ['simple']")
	}
	if internaltypes.IsDefined(model.NoneIncludedUserFilter) && model.Type.ValueString() != "simple" {
		resp.Diagnostics.AddError("Attribute 'none_included_user_filter' not supported by pingdirectory_connection_criteria resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'none_included_user_filter', the 'type' attribute must be one of ['simple']")
	}
}

// Add optional fields to create request for simple connection-criteria
func addOptionalSimpleConnectionCriteriaFields(ctx context.Context, addRequest *client.AddSimpleConnectionCriteriaRequest, plan connectionCriteriaResourceModel) error {
	if internaltypes.IsDefined(plan.IncludedClientAddress) {
		var slice []string
		plan.IncludedClientAddress.ElementsAs(ctx, &slice, false)
		addRequest.IncludedClientAddress = slice
	}
	if internaltypes.IsDefined(plan.ExcludedClientAddress) {
		var slice []string
		plan.ExcludedClientAddress.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedClientAddress = slice
	}
	if internaltypes.IsDefined(plan.IncludedConnectionHandler) {
		var slice []string
		plan.IncludedConnectionHandler.ElementsAs(ctx, &slice, false)
		addRequest.IncludedConnectionHandler = slice
	}
	if internaltypes.IsDefined(plan.ExcludedConnectionHandler) {
		var slice []string
		plan.ExcludedConnectionHandler.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedConnectionHandler = slice
	}
	if internaltypes.IsDefined(plan.IncludedProtocol) {
		var slice []string
		plan.IncludedProtocol.ElementsAs(ctx, &slice, false)
		addRequest.IncludedProtocol = slice
	}
	if internaltypes.IsDefined(plan.ExcludedProtocol) {
		var slice []string
		plan.ExcludedProtocol.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedProtocol = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CommunicationSecurityLevel) {
		communicationSecurityLevel, err := client.NewEnumconnectionCriteriaCommunicationSecurityLevelPropFromValue(plan.CommunicationSecurityLevel.ValueString())
		if err != nil {
			return err
		}
		addRequest.CommunicationSecurityLevel = communicationSecurityLevel
	}
	if internaltypes.IsDefined(plan.UserAuthType) {
		var slice []string
		plan.UserAuthType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumconnectionCriteriaUserAuthTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumconnectionCriteriaUserAuthTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.UserAuthType = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuthenticationSecurityLevel) {
		authenticationSecurityLevel, err := client.NewEnumconnectionCriteriaAuthenticationSecurityLevelPropFromValue(plan.AuthenticationSecurityLevel.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuthenticationSecurityLevel = authenticationSecurityLevel
	}
	if internaltypes.IsDefined(plan.IncludedUserSASLMechanism) {
		var slice []string
		plan.IncludedUserSASLMechanism.ElementsAs(ctx, &slice, false)
		addRequest.IncludedUserSASLMechanism = slice
	}
	if internaltypes.IsDefined(plan.ExcludedUserSASLMechanism) {
		var slice []string
		plan.ExcludedUserSASLMechanism.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedUserSASLMechanism = slice
	}
	if internaltypes.IsDefined(plan.IncludedUserBaseDN) {
		var slice []string
		plan.IncludedUserBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.IncludedUserBaseDN = slice
	}
	if internaltypes.IsDefined(plan.ExcludedUserBaseDN) {
		var slice []string
		plan.ExcludedUserBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedUserBaseDN = slice
	}
	if internaltypes.IsDefined(plan.AllIncludedUserGroupDN) {
		var slice []string
		plan.AllIncludedUserGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedUserGroupDN = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedUserGroupDN) {
		var slice []string
		plan.AnyIncludedUserGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedUserGroupDN = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedUserGroupDN) {
		var slice []string
		plan.NotAllIncludedUserGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedUserGroupDN = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedUserGroupDN) {
		var slice []string
		plan.NoneIncludedUserGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedUserGroupDN = slice
	}
	if internaltypes.IsDefined(plan.AllIncludedUserFilter) {
		var slice []string
		plan.AllIncludedUserFilter.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedUserFilter = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedUserFilter) {
		var slice []string
		plan.AnyIncludedUserFilter.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedUserFilter = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedUserFilter) {
		var slice []string
		plan.NotAllIncludedUserFilter.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedUserFilter = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedUserFilter) {
		var slice []string
		plan.NoneIncludedUserFilter.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedUserFilter = slice
	}
	if internaltypes.IsDefined(plan.AllIncludedUserPrivilege) {
		var slice []string
		plan.AllIncludedUserPrivilege.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumconnectionCriteriaAllIncludedUserPrivilegeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumconnectionCriteriaAllIncludedUserPrivilegePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AllIncludedUserPrivilege = enumSlice
	}
	if internaltypes.IsDefined(plan.AnyIncludedUserPrivilege) {
		var slice []string
		plan.AnyIncludedUserPrivilege.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumconnectionCriteriaAnyIncludedUserPrivilegeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumconnectionCriteriaAnyIncludedUserPrivilegePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AnyIncludedUserPrivilege = enumSlice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedUserPrivilege) {
		var slice []string
		plan.NotAllIncludedUserPrivilege.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumconnectionCriteriaNotAllIncludedUserPrivilegeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumconnectionCriteriaNotAllIncludedUserPrivilegePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.NotAllIncludedUserPrivilege = enumSlice
	}
	if internaltypes.IsDefined(plan.NoneIncludedUserPrivilege) {
		var slice []string
		plan.NoneIncludedUserPrivilege.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumconnectionCriteriaNoneIncludedUserPrivilegeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumconnectionCriteriaNoneIncludedUserPrivilegePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.NoneIncludedUserPrivilege = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for aggregate connection-criteria
func addOptionalAggregateConnectionCriteriaFields(ctx context.Context, addRequest *client.AddAggregateConnectionCriteriaRequest, plan connectionCriteriaResourceModel) error {
	if internaltypes.IsDefined(plan.AllIncludedConnectionCriteria) {
		var slice []string
		plan.AllIncludedConnectionCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedConnectionCriteria = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedConnectionCriteria) {
		var slice []string
		plan.AnyIncludedConnectionCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedConnectionCriteria = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedConnectionCriteria) {
		var slice []string
		plan.NotAllIncludedConnectionCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedConnectionCriteria = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedConnectionCriteria) {
		var slice []string
		plan.NoneIncludedConnectionCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedConnectionCriteria = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for third-party connection-criteria
func addOptionalThirdPartyConnectionCriteriaFields(ctx context.Context, addRequest *client.AddThirdPartyConnectionCriteriaRequest, plan connectionCriteriaResourceModel) error {
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
func populateConnectionCriteriaUnknownValues(ctx context.Context, model *connectionCriteriaResourceModel) {
	if model.ExcludedClientAddress.ElementType(ctx) == nil {
		model.ExcludedClientAddress = types.SetNull(types.StringType)
	}
	if model.IncludedClientAddress.ElementType(ctx) == nil {
		model.IncludedClientAddress = types.SetNull(types.StringType)
	}
	if model.AllIncludedUserPrivilege.ElementType(ctx) == nil {
		model.AllIncludedUserPrivilege = types.SetNull(types.StringType)
	}
	if model.NotAllIncludedUserFilter.ElementType(ctx) == nil {
		model.NotAllIncludedUserFilter = types.SetNull(types.StringType)
	}
	if model.AllIncludedUserGroupDN.ElementType(ctx) == nil {
		model.AllIncludedUserGroupDN = types.SetNull(types.StringType)
	}
	if model.ExcludedUserSASLMechanism.ElementType(ctx) == nil {
		model.ExcludedUserSASLMechanism = types.SetNull(types.StringType)
	}
	if model.AnyIncludedConnectionCriteria.ElementType(ctx) == nil {
		model.AnyIncludedConnectionCriteria = types.SetNull(types.StringType)
	}
	if model.ExcludedConnectionHandler.ElementType(ctx) == nil {
		model.ExcludedConnectionHandler = types.SetNull(types.StringType)
	}
	if model.NotAllIncludedConnectionCriteria.ElementType(ctx) == nil {
		model.NotAllIncludedConnectionCriteria = types.SetNull(types.StringType)
	}
	if model.AllIncludedUserFilter.ElementType(ctx) == nil {
		model.AllIncludedUserFilter = types.SetNull(types.StringType)
	}
	if model.NoneIncludedUserPrivilege.ElementType(ctx) == nil {
		model.NoneIncludedUserPrivilege = types.SetNull(types.StringType)
	}
	if model.NotAllIncludedUserPrivilege.ElementType(ctx) == nil {
		model.NotAllIncludedUserPrivilege = types.SetNull(types.StringType)
	}
	if model.ExcludedUserBaseDN.ElementType(ctx) == nil {
		model.ExcludedUserBaseDN = types.SetNull(types.StringType)
	}
	if model.NotAllIncludedUserGroupDN.ElementType(ctx) == nil {
		model.NotAllIncludedUserGroupDN = types.SetNull(types.StringType)
	}
	if model.IncludedProtocol.ElementType(ctx) == nil {
		model.IncludedProtocol = types.SetNull(types.StringType)
	}
	if model.AnyIncludedUserPrivilege.ElementType(ctx) == nil {
		model.AnyIncludedUserPrivilege = types.SetNull(types.StringType)
	}
	if model.AnyIncludedUserGroupDN.ElementType(ctx) == nil {
		model.AnyIncludedUserGroupDN = types.SetNull(types.StringType)
	}
	if model.UserAuthType.ElementType(ctx) == nil {
		model.UserAuthType = types.SetNull(types.StringType)
	}
	if model.IncludedUserSASLMechanism.ElementType(ctx) == nil {
		model.IncludedUserSASLMechanism = types.SetNull(types.StringType)
	}
	if model.ExtensionArgument.ElementType(ctx) == nil {
		model.ExtensionArgument = types.SetNull(types.StringType)
	}
	if model.IncludedConnectionHandler.ElementType(ctx) == nil {
		model.IncludedConnectionHandler = types.SetNull(types.StringType)
	}
	if model.IncludedUserBaseDN.ElementType(ctx) == nil {
		model.IncludedUserBaseDN = types.SetNull(types.StringType)
	}
	if model.NoneIncludedUserFilter.ElementType(ctx) == nil {
		model.NoneIncludedUserFilter = types.SetNull(types.StringType)
	}
	if model.AllIncludedConnectionCriteria.ElementType(ctx) == nil {
		model.AllIncludedConnectionCriteria = types.SetNull(types.StringType)
	}
	if model.ExcludedProtocol.ElementType(ctx) == nil {
		model.ExcludedProtocol = types.SetNull(types.StringType)
	}
	if model.AnyIncludedUserFilter.ElementType(ctx) == nil {
		model.AnyIncludedUserFilter = types.SetNull(types.StringType)
	}
	if model.NoneIncludedConnectionCriteria.ElementType(ctx) == nil {
		model.NoneIncludedConnectionCriteria = types.SetNull(types.StringType)
	}
	if model.NoneIncludedUserGroupDN.ElementType(ctx) == nil {
		model.NoneIncludedUserGroupDN = types.SetNull(types.StringType)
	}
}

// Read a SimpleConnectionCriteriaResponse object into the model struct
func readSimpleConnectionCriteriaResponse(ctx context.Context, r *client.SimpleConnectionCriteriaResponse, state *connectionCriteriaResourceModel, expectedValues *connectionCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("simple")
	state.Id = types.StringValue(r.Id)
	state.IncludedClientAddress = internaltypes.GetStringSet(r.IncludedClientAddress)
	state.ExcludedClientAddress = internaltypes.GetStringSet(r.ExcludedClientAddress)
	state.IncludedConnectionHandler = internaltypes.GetStringSet(r.IncludedConnectionHandler)
	state.ExcludedConnectionHandler = internaltypes.GetStringSet(r.ExcludedConnectionHandler)
	state.IncludedProtocol = internaltypes.GetStringSet(r.IncludedProtocol)
	state.ExcludedProtocol = internaltypes.GetStringSet(r.ExcludedProtocol)
	state.CommunicationSecurityLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumconnectionCriteriaCommunicationSecurityLevelProp(r.CommunicationSecurityLevel), internaltypes.IsEmptyString(expectedValues.CommunicationSecurityLevel))
	state.UserAuthType = internaltypes.GetStringSet(
		client.StringSliceEnumconnectionCriteriaUserAuthTypeProp(r.UserAuthType))
	state.AuthenticationSecurityLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumconnectionCriteriaAuthenticationSecurityLevelProp(r.AuthenticationSecurityLevel), internaltypes.IsEmptyString(expectedValues.AuthenticationSecurityLevel))
	state.IncludedUserSASLMechanism = internaltypes.GetStringSet(r.IncludedUserSASLMechanism)
	state.ExcludedUserSASLMechanism = internaltypes.GetStringSet(r.ExcludedUserSASLMechanism)
	state.IncludedUserBaseDN = internaltypes.GetStringSet(r.IncludedUserBaseDN)
	state.ExcludedUserBaseDN = internaltypes.GetStringSet(r.ExcludedUserBaseDN)
	state.AllIncludedUserGroupDN = internaltypes.GetStringSet(r.AllIncludedUserGroupDN)
	state.AnyIncludedUserGroupDN = internaltypes.GetStringSet(r.AnyIncludedUserGroupDN)
	state.NotAllIncludedUserGroupDN = internaltypes.GetStringSet(r.NotAllIncludedUserGroupDN)
	state.NoneIncludedUserGroupDN = internaltypes.GetStringSet(r.NoneIncludedUserGroupDN)
	state.AllIncludedUserFilter = internaltypes.GetStringSet(r.AllIncludedUserFilter)
	state.AnyIncludedUserFilter = internaltypes.GetStringSet(r.AnyIncludedUserFilter)
	state.NotAllIncludedUserFilter = internaltypes.GetStringSet(r.NotAllIncludedUserFilter)
	state.NoneIncludedUserFilter = internaltypes.GetStringSet(r.NoneIncludedUserFilter)
	state.AllIncludedUserPrivilege = internaltypes.GetStringSet(
		client.StringSliceEnumconnectionCriteriaAllIncludedUserPrivilegeProp(r.AllIncludedUserPrivilege))
	state.AnyIncludedUserPrivilege = internaltypes.GetStringSet(
		client.StringSliceEnumconnectionCriteriaAnyIncludedUserPrivilegeProp(r.AnyIncludedUserPrivilege))
	state.NotAllIncludedUserPrivilege = internaltypes.GetStringSet(
		client.StringSliceEnumconnectionCriteriaNotAllIncludedUserPrivilegeProp(r.NotAllIncludedUserPrivilege))
	state.NoneIncludedUserPrivilege = internaltypes.GetStringSet(
		client.StringSliceEnumconnectionCriteriaNoneIncludedUserPrivilegeProp(r.NoneIncludedUserPrivilege))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateConnectionCriteriaUnknownValues(ctx, state)
}

// Read a AggregateConnectionCriteriaResponse object into the model struct
func readAggregateConnectionCriteriaResponse(ctx context.Context, r *client.AggregateConnectionCriteriaResponse, state *connectionCriteriaResourceModel, expectedValues *connectionCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("aggregate")
	state.Id = types.StringValue(r.Id)
	state.AllIncludedConnectionCriteria = internaltypes.GetStringSet(r.AllIncludedConnectionCriteria)
	state.AnyIncludedConnectionCriteria = internaltypes.GetStringSet(r.AnyIncludedConnectionCriteria)
	state.NotAllIncludedConnectionCriteria = internaltypes.GetStringSet(r.NotAllIncludedConnectionCriteria)
	state.NoneIncludedConnectionCriteria = internaltypes.GetStringSet(r.NoneIncludedConnectionCriteria)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateConnectionCriteriaUnknownValues(ctx, state)
}

// Read a ThirdPartyConnectionCriteriaResponse object into the model struct
func readThirdPartyConnectionCriteriaResponse(ctx context.Context, r *client.ThirdPartyConnectionCriteriaResponse, state *connectionCriteriaResourceModel, expectedValues *connectionCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateConnectionCriteriaUnknownValues(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createConnectionCriteriaOperations(plan connectionCriteriaResourceModel, state connectionCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedConnectionCriteria, state.AllIncludedConnectionCriteria, "all-included-connection-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedConnectionCriteria, state.AnyIncludedConnectionCriteria, "any-included-connection-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedConnectionCriteria, state.NotAllIncludedConnectionCriteria, "not-all-included-connection-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedConnectionCriteria, state.NoneIncludedConnectionCriteria, "none-included-connection-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedClientAddress, state.IncludedClientAddress, "included-client-address")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedClientAddress, state.ExcludedClientAddress, "excluded-client-address")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedConnectionHandler, state.IncludedConnectionHandler, "included-connection-handler")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedConnectionHandler, state.ExcludedConnectionHandler, "excluded-connection-handler")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedProtocol, state.IncludedProtocol, "included-protocol")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedProtocol, state.ExcludedProtocol, "excluded-protocol")
	operations.AddStringOperationIfNecessary(&ops, plan.CommunicationSecurityLevel, state.CommunicationSecurityLevel, "communication-security-level")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.UserAuthType, state.UserAuthType, "user-auth-type")
	operations.AddStringOperationIfNecessary(&ops, plan.AuthenticationSecurityLevel, state.AuthenticationSecurityLevel, "authentication-security-level")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedUserSASLMechanism, state.IncludedUserSASLMechanism, "included-user-sasl-mechanism")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedUserSASLMechanism, state.ExcludedUserSASLMechanism, "excluded-user-sasl-mechanism")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedUserBaseDN, state.IncludedUserBaseDN, "included-user-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedUserBaseDN, state.ExcludedUserBaseDN, "excluded-user-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedUserGroupDN, state.AllIncludedUserGroupDN, "all-included-user-group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedUserGroupDN, state.AnyIncludedUserGroupDN, "any-included-user-group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedUserGroupDN, state.NotAllIncludedUserGroupDN, "not-all-included-user-group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedUserGroupDN, state.NoneIncludedUserGroupDN, "none-included-user-group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedUserFilter, state.AllIncludedUserFilter, "all-included-user-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedUserFilter, state.AnyIncludedUserFilter, "any-included-user-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedUserFilter, state.NotAllIncludedUserFilter, "not-all-included-user-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedUserFilter, state.NoneIncludedUserFilter, "none-included-user-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedUserPrivilege, state.AllIncludedUserPrivilege, "all-included-user-privilege")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedUserPrivilege, state.AnyIncludedUserPrivilege, "any-included-user-privilege")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedUserPrivilege, state.NotAllIncludedUserPrivilege, "not-all-included-user-privilege")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedUserPrivilege, state.NoneIncludedUserPrivilege, "none-included-user-privilege")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a simple connection-criteria
func (r *connectionCriteriaResource) CreateSimpleConnectionCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan connectionCriteriaResourceModel) (*connectionCriteriaResourceModel, error) {
	addRequest := client.NewAddSimpleConnectionCriteriaRequest(plan.Id.ValueString(),
		[]client.EnumsimpleConnectionCriteriaSchemaUrn{client.ENUMSIMPLECONNECTIONCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CONNECTION_CRITERIASIMPLE})
	err := addOptionalSimpleConnectionCriteriaFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Connection Criteria", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ConnectionCriteriaApi.AddConnectionCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddConnectionCriteriaRequest(
		client.AddSimpleConnectionCriteriaRequestAsAddConnectionCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ConnectionCriteriaApi.AddConnectionCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Connection Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state connectionCriteriaResourceModel
	readSimpleConnectionCriteriaResponse(ctx, addResponse.SimpleConnectionCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a aggregate connection-criteria
func (r *connectionCriteriaResource) CreateAggregateConnectionCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan connectionCriteriaResourceModel) (*connectionCriteriaResourceModel, error) {
	addRequest := client.NewAddAggregateConnectionCriteriaRequest(plan.Id.ValueString(),
		[]client.EnumaggregateConnectionCriteriaSchemaUrn{client.ENUMAGGREGATECONNECTIONCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CONNECTION_CRITERIAAGGREGATE})
	err := addOptionalAggregateConnectionCriteriaFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Connection Criteria", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ConnectionCriteriaApi.AddConnectionCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddConnectionCriteriaRequest(
		client.AddAggregateConnectionCriteriaRequestAsAddConnectionCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ConnectionCriteriaApi.AddConnectionCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Connection Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state connectionCriteriaResourceModel
	readAggregateConnectionCriteriaResponse(ctx, addResponse.AggregateConnectionCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party connection-criteria
func (r *connectionCriteriaResource) CreateThirdPartyConnectionCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan connectionCriteriaResourceModel) (*connectionCriteriaResourceModel, error) {
	addRequest := client.NewAddThirdPartyConnectionCriteriaRequest(plan.Id.ValueString(),
		[]client.EnumthirdPartyConnectionCriteriaSchemaUrn{client.ENUMTHIRDPARTYCONNECTIONCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CONNECTION_CRITERIATHIRD_PARTY},
		plan.ExtensionClass.ValueString())
	err := addOptionalThirdPartyConnectionCriteriaFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Connection Criteria", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ConnectionCriteriaApi.AddConnectionCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddConnectionCriteriaRequest(
		client.AddThirdPartyConnectionCriteriaRequestAsAddConnectionCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ConnectionCriteriaApi.AddConnectionCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Connection Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state connectionCriteriaResourceModel
	readThirdPartyConnectionCriteriaResponse(ctx, addResponse.ThirdPartyConnectionCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *connectionCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan connectionCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *connectionCriteriaResourceModel
	var err error
	if plan.Type.ValueString() == "simple" {
		state, err = r.CreateSimpleConnectionCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "aggregate" {
		state, err = r.CreateAggregateConnectionCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyConnectionCriteria(ctx, req, resp, plan)
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
func (r *defaultConnectionCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan connectionCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConnectionCriteriaApi.GetConnectionCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Connection Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state connectionCriteriaResourceModel
	if plan.Type.ValueString() == "simple" {
		readSimpleConnectionCriteriaResponse(ctx, readResponse.SimpleConnectionCriteriaResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "aggregate" {
		readAggregateConnectionCriteriaResponse(ctx, readResponse.AggregateConnectionCriteriaResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "third-party" {
		readThirdPartyConnectionCriteriaResponse(ctx, readResponse.ThirdPartyConnectionCriteriaResponse, &state, &plan, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ConnectionCriteriaApi.UpdateConnectionCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createConnectionCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ConnectionCriteriaApi.UpdateConnectionCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Connection Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "simple" {
			readSimpleConnectionCriteriaResponse(ctx, updateResponse.SimpleConnectionCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "aggregate" {
			readAggregateConnectionCriteriaResponse(ctx, updateResponse.AggregateConnectionCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party" {
			readThirdPartyConnectionCriteriaResponse(ctx, updateResponse.ThirdPartyConnectionCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *connectionCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConnectionCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultConnectionCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConnectionCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readConnectionCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state connectionCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ConnectionCriteriaApi.GetConnectionCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Connection Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.SimpleConnectionCriteriaResponse != nil {
		readSimpleConnectionCriteriaResponse(ctx, readResponse.SimpleConnectionCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AggregateConnectionCriteriaResponse != nil {
		readAggregateConnectionCriteriaResponse(ctx, readResponse.AggregateConnectionCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyConnectionCriteriaResponse != nil {
		readThirdPartyConnectionCriteriaResponse(ctx, readResponse.ThirdPartyConnectionCriteriaResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *connectionCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConnectionCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultConnectionCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConnectionCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateConnectionCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan connectionCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state connectionCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ConnectionCriteriaApi.UpdateConnectionCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createConnectionCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ConnectionCriteriaApi.UpdateConnectionCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Connection Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "simple" {
			readSimpleConnectionCriteriaResponse(ctx, updateResponse.SimpleConnectionCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "aggregate" {
			readAggregateConnectionCriteriaResponse(ctx, updateResponse.AggregateConnectionCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party" {
			readThirdPartyConnectionCriteriaResponse(ctx, updateResponse.ThirdPartyConnectionCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultConnectionCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *connectionCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state connectionCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ConnectionCriteriaApi.DeleteConnectionCriteriaExecute(r.apiClient.ConnectionCriteriaApi.DeleteConnectionCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Connection Criteria", err, httpResp)
		return
	}
}

func (r *connectionCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConnectionCriteria(ctx, req, resp)
}

func (r *defaultConnectionCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConnectionCriteria(ctx, req, resp)
}

func importConnectionCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
