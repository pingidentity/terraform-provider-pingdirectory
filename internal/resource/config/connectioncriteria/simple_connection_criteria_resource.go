package connectioncriteria

import (
	"context"
	"time"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &simpleConnectionCriteriaResource{}
	_ resource.ResourceWithConfigure   = &simpleConnectionCriteriaResource{}
	_ resource.ResourceWithImportState = &simpleConnectionCriteriaResource{}
	_ resource.Resource                = &defaultSimpleConnectionCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultSimpleConnectionCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultSimpleConnectionCriteriaResource{}
)

// Create a Simple Connection Criteria resource
func NewSimpleConnectionCriteriaResource() resource.Resource {
	return &simpleConnectionCriteriaResource{}
}

func NewDefaultSimpleConnectionCriteriaResource() resource.Resource {
	return &defaultSimpleConnectionCriteriaResource{}
}

// simpleConnectionCriteriaResource is the resource implementation.
type simpleConnectionCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSimpleConnectionCriteriaResource is the resource implementation.
type defaultSimpleConnectionCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *simpleConnectionCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_simple_connection_criteria"
}

func (r *defaultSimpleConnectionCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_simple_connection_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *simpleConnectionCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultSimpleConnectionCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type simpleConnectionCriteriaResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	LastUpdated                 types.String `tfsdk:"last_updated"`
	Notifications               types.Set    `tfsdk:"notifications"`
	RequiredActions             types.Set    `tfsdk:"required_actions"`
	IncludedClientAddress       types.Set    `tfsdk:"included_client_address"`
	ExcludedClientAddress       types.Set    `tfsdk:"excluded_client_address"`
	IncludedConnectionHandler   types.Set    `tfsdk:"included_connection_handler"`
	ExcludedConnectionHandler   types.Set    `tfsdk:"excluded_connection_handler"`
	IncludedProtocol            types.Set    `tfsdk:"included_protocol"`
	ExcludedProtocol            types.Set    `tfsdk:"excluded_protocol"`
	CommunicationSecurityLevel  types.String `tfsdk:"communication_security_level"`
	UserAuthType                types.Set    `tfsdk:"user_auth_type"`
	AuthenticationSecurityLevel types.String `tfsdk:"authentication_security_level"`
	IncludedUserSASLMechanism   types.Set    `tfsdk:"included_user_sasl_mechanism"`
	ExcludedUserSASLMechanism   types.Set    `tfsdk:"excluded_user_sasl_mechanism"`
	IncludedUserBaseDN          types.Set    `tfsdk:"included_user_base_dn"`
	ExcludedUserBaseDN          types.Set    `tfsdk:"excluded_user_base_dn"`
	AllIncludedUserGroupDN      types.Set    `tfsdk:"all_included_user_group_dn"`
	AnyIncludedUserGroupDN      types.Set    `tfsdk:"any_included_user_group_dn"`
	NotAllIncludedUserGroupDN   types.Set    `tfsdk:"not_all_included_user_group_dn"`
	NoneIncludedUserGroupDN     types.Set    `tfsdk:"none_included_user_group_dn"`
	AllIncludedUserFilter       types.Set    `tfsdk:"all_included_user_filter"`
	AnyIncludedUserFilter       types.Set    `tfsdk:"any_included_user_filter"`
	NotAllIncludedUserFilter    types.Set    `tfsdk:"not_all_included_user_filter"`
	NoneIncludedUserFilter      types.Set    `tfsdk:"none_included_user_filter"`
	AllIncludedUserPrivilege    types.Set    `tfsdk:"all_included_user_privilege"`
	AnyIncludedUserPrivilege    types.Set    `tfsdk:"any_included_user_privilege"`
	NotAllIncludedUserPrivilege types.Set    `tfsdk:"not_all_included_user_privilege"`
	NoneIncludedUserPrivilege   types.Set    `tfsdk:"none_included_user_privilege"`
	Description                 types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *simpleConnectionCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	simpleConnectionCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultSimpleConnectionCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	simpleConnectionCriteriaSchema(ctx, req, resp, true)
}

func simpleConnectionCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Simple Connection Criteria.",
		Attributes: map[string]schema.Attribute{
			"included_client_address": schema.SetAttribute{
				Description: "Specifies an address mask that may be used to specify a set of clients that should be included in this Simple Connection Criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"excluded_client_address": schema.SetAttribute{
				Description: "Specifies an address mask that may be used to specify a set of clients that should be excluded from this Simple Connection Criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"included_connection_handler": schema.SetAttribute{
				Description: "Specifies a connection handler for clients that should be included in this Simple Connection Criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"excluded_connection_handler": schema.SetAttribute{
				Description: "Specifies a connection handler for clients that should be excluded from this Simple Connection Criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"included_protocol": schema.SetAttribute{
				Description: "Specifies the name of a communication protocol that should be used by clients included in this Simple Connection Criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"excluded_protocol": schema.SetAttribute{
				Description: "Specifies the name of a communication protocol that should be used by clients excluded from this Simple Connection Criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
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
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
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
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"excluded_user_sasl_mechanism": schema.SetAttribute{
				Description: "Specifies the name of a SASL mechanism that should be used by clients excluded from this Simple Connection Criteria. This will only be taken into account for client connections that have authenticated to the server using a SASL mechanism and will be ignored for unauthenticated client connections and for client connections that authenticated using some other method (e.g., those performing simple or internal authentication).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"included_user_base_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which authenticated user entries may exist for clients included in this Simple Connection Criteria. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections. Refer to the authz version of this property in Simple Result Criteria if operations are being proxied (performed using proxied authorization), and you need to match the originating user of the operation rather than the proxy user (the user the proxy authenticated as).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"excluded_user_base_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which authenticated user entries may exist for clients excluded from this Simple Connection Criteria. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections. Refer to the authz version of this property in Simple Result Criteria if operations are being proxied (performed using proxied authorization), and you need to match the originating user of the operation rather than the proxy user (the user the proxy authenticated as).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"all_included_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authenticated users must exist for clients included in this Simple Connection Criteria. If any group DNs are provided, then the authenticated user must be a member of all of those groups. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections. Refer to the authz version of this property in Simple Result Criteria if operations are being proxied (performed using proxied authorization), and you need to match the originating user of the operation rather than the proxy user (the user the proxy authenticated as).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"any_included_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authenticated users may exist for clients included in this Simple Connection Criteria. If any group DNs are provided, then the authenticated user must be a member of at least one of those groups. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections. Refer to the authz version of this property in Simple Result Criteria if operations are being proxied (performed using proxied authorization), and you need to match the originating user of the operation rather than the proxy user (the user the proxy authenticated as).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"not_all_included_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authenticated users should not exist for clients included in this Simple Connection Criteria. If any group DNs are provided, then the authenticated user must not be a member of at least one of those groups (that is, the user may be a member of zero or more of those groups, but not of all of them). This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections. Refer to the authz version of this property in Simple Result Criteria if operations are being proxied (performed using proxied authorization), and you need to match the originating user of the operation rather than the proxy user (the user the proxy authenticated as).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"none_included_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authenticated users must not exist for clients included in this Simple Connection Criteria. If any group DNs are provided, then the authenticated user must not be a member any of those groups. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections. Refer to the authz version of this property in Simple Result Criteria if operations are being proxied (performed using proxied authorization), and you need to match the originating user of the operation rather than the proxy user (the user the proxy authenticated as).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"all_included_user_filter": schema.SetAttribute{
				Description: "Specifies a search filter that must match the entry of the authenticated user for clients included in this Simple Connection Criteria. If any filters are provided, then all of those filters must match the authenticated user entry. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"any_included_user_filter": schema.SetAttribute{
				Description: "Specifies a search filter that may match the entry of the authenticated user for clients included in this Simple Connection Criteria. If any filters are provided, then at least one of those filters must match the authenticated user entry. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"not_all_included_user_filter": schema.SetAttribute{
				Description: "Specifies a search filter that should not match the entry of the authenticated user for clients included in this Simple Connection Criteria. If any filters are provided, then at least one of those filters must not match the authenticated user entry (that is, the user entry may match zero or more of those filters, but not all of them). This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"none_included_user_filter": schema.SetAttribute{
				Description: "Specifies a search filter that must not match the entry of the authenticated user for clients included in this Simple Connection Criteria. If any filters are provided, then none of those filters may match the authenticated user entry. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"all_included_user_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that must be held by the authenticated user for clients included in this Simple Connection Criteria. If any privilege names are provided, then the authenticated user must have all of those privileges. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"any_included_user_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that may be held by the authenticated user for clients included in this Simple Connection Criteria. If any privilege names are provided, then the authenticated user must have at least one of those privileges. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"not_all_included_user_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that should not be held by the authenticated user for clients included in this Simple Connection Criteria. If any privilege names are provided, then the authenticated user must not have at least one of those privileges (that is, the user may hold zero or more of those privileges, but not all of them). This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"none_included_user_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that must not be held by the authenticated user for clients included in this Simple Connection Criteria. If any privilege names are provided, then the authenticated user must not have any of those privileges. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Connection Criteria",
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
func addOptionalSimpleConnectionCriteriaFields(ctx context.Context, addRequest *client.AddSimpleConnectionCriteriaRequest, plan simpleConnectionCriteriaResourceModel) error {
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
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	return nil
}

// Read a SimpleConnectionCriteriaResponse object into the model struct
func readSimpleConnectionCriteriaResponse(ctx context.Context, r *client.SimpleConnectionCriteriaResponse, state *simpleConnectionCriteriaResourceModel, expectedValues *simpleConnectionCriteriaResourceModel, diagnostics *diag.Diagnostics) {
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
}

// Create any update operations necessary to make the state match the plan
func createSimpleConnectionCriteriaOperations(plan simpleConnectionCriteriaResourceModel, state simpleConnectionCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
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

// Create a new resource
func (r *simpleConnectionCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan simpleConnectionCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddSimpleConnectionCriteriaRequest(plan.Id.ValueString(),
		[]client.EnumsimpleConnectionCriteriaSchemaUrn{client.ENUMSIMPLECONNECTIONCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CONNECTION_CRITERIASIMPLE})
	err := addOptionalSimpleConnectionCriteriaFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Simple Connection Criteria", err.Error())
		return
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Simple Connection Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state simpleConnectionCriteriaResourceModel
	readSimpleConnectionCriteriaResponse(ctx, addResponse.SimpleConnectionCriteriaResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultSimpleConnectionCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan simpleConnectionCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConnectionCriteriaApi.GetConnectionCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Simple Connection Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state simpleConnectionCriteriaResourceModel
	readSimpleConnectionCriteriaResponse(ctx, readResponse.SimpleConnectionCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ConnectionCriteriaApi.UpdateConnectionCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSimpleConnectionCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ConnectionCriteriaApi.UpdateConnectionCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Simple Connection Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSimpleConnectionCriteriaResponse(ctx, updateResponse.SimpleConnectionCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *simpleConnectionCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSimpleConnectionCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSimpleConnectionCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSimpleConnectionCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSimpleConnectionCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state simpleConnectionCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ConnectionCriteriaApi.GetConnectionCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Simple Connection Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSimpleConnectionCriteriaResponse(ctx, readResponse.SimpleConnectionCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *simpleConnectionCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSimpleConnectionCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSimpleConnectionCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSimpleConnectionCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSimpleConnectionCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan simpleConnectionCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state simpleConnectionCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ConnectionCriteriaApi.UpdateConnectionCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSimpleConnectionCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ConnectionCriteriaApi.UpdateConnectionCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Simple Connection Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSimpleConnectionCriteriaResponse(ctx, updateResponse.SimpleConnectionCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSimpleConnectionCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *simpleConnectionCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state simpleConnectionCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ConnectionCriteriaApi.DeleteConnectionCriteriaExecute(r.apiClient.ConnectionCriteriaApi.DeleteConnectionCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Simple Connection Criteria", err, httpResp)
		return
	}
}

func (r *simpleConnectionCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSimpleConnectionCriteria(ctx, req, resp)
}

func (r *defaultSimpleConnectionCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSimpleConnectionCriteria(ctx, req, resp)
}

func importSimpleConnectionCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
