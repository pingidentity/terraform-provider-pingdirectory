package passthroughauthenticationhandler

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &passThroughAuthenticationHandlerResource{}
	_ resource.ResourceWithConfigure   = &passThroughAuthenticationHandlerResource{}
	_ resource.ResourceWithImportState = &passThroughAuthenticationHandlerResource{}
	_ resource.Resource                = &defaultPassThroughAuthenticationHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultPassThroughAuthenticationHandlerResource{}
	_ resource.ResourceWithImportState = &defaultPassThroughAuthenticationHandlerResource{}
)

// Create a Pass Through Authentication Handler resource
func NewPassThroughAuthenticationHandlerResource() resource.Resource {
	return &passThroughAuthenticationHandlerResource{}
}

func NewDefaultPassThroughAuthenticationHandlerResource() resource.Resource {
	return &defaultPassThroughAuthenticationHandlerResource{}
}

// passThroughAuthenticationHandlerResource is the resource implementation.
type passThroughAuthenticationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultPassThroughAuthenticationHandlerResource is the resource implementation.
type defaultPassThroughAuthenticationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *passThroughAuthenticationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pass_through_authentication_handler"
}

func (r *defaultPassThroughAuthenticationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_pass_through_authentication_handler"
}

// Configure adds the provider configured client to the resource.
func (r *passThroughAuthenticationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultPassThroughAuthenticationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type passThroughAuthenticationHandlerResourceModel struct {
	Id                                          types.String `tfsdk:"id"`
	Name                                        types.String `tfsdk:"name"`
	LastUpdated                                 types.String `tfsdk:"last_updated"`
	Notifications                               types.Set    `tfsdk:"notifications"`
	RequiredActions                             types.Set    `tfsdk:"required_actions"`
	Type                                        types.String `tfsdk:"type"`
	ExtensionClass                              types.String `tfsdk:"extension_class"`
	ExtensionArgument                           types.Set    `tfsdk:"extension_argument"`
	SubordinatePassThroughAuthenticationHandler types.Set    `tfsdk:"subordinate_pass_through_authentication_handler"`
	ContinueOnFailureType                       types.Set    `tfsdk:"continue_on_failure_type"`
	Server                                      types.Set    `tfsdk:"server"`
	ServerAccessMode                            types.String `tfsdk:"server_access_mode"`
	DnMap                                       types.Set    `tfsdk:"dn_map"`
	BindDNPattern                               types.String `tfsdk:"bind_dn_pattern"`
	SearchBaseDN                                types.String `tfsdk:"search_base_dn"`
	SearchFilterPattern                         types.String `tfsdk:"search_filter_pattern"`
	InitialConnections                          types.Int64  `tfsdk:"initial_connections"`
	MaxConnections                              types.Int64  `tfsdk:"max_connections"`
	UseLocation                                 types.Bool   `tfsdk:"use_location"`
	MaximumAllowedLocalResponseTime             types.String `tfsdk:"maximum_allowed_local_response_time"`
	MaximumAllowedNonlocalResponseTime          types.String `tfsdk:"maximum_allowed_nonlocal_response_time"`
	UsePasswordPolicyControl                    types.Bool   `tfsdk:"use_password_policy_control"`
	ApiURL                                      types.String `tfsdk:"api_url"`
	AuthURL                                     types.String `tfsdk:"auth_url"`
	OAuthClientID                               types.String `tfsdk:"oauth_client_id"`
	OAuthClientSecret                           types.String `tfsdk:"oauth_client_secret"`
	OAuthClientSecretPassphraseProvider         types.String `tfsdk:"oauth_client_secret_passphrase_provider"`
	EnvironmentID                               types.String `tfsdk:"environment_id"`
	HttpProxyExternalServer                     types.String `tfsdk:"http_proxy_external_server"`
	UserMappingLocalAttribute                   types.Set    `tfsdk:"user_mapping_local_attribute"`
	UserMappingRemoteJSONField                  types.Set    `tfsdk:"user_mapping_remote_json_field"`
	AdditionalUserMappingSCIMFilter             types.String `tfsdk:"additional_user_mapping_scim_filter"`
	Description                                 types.String `tfsdk:"description"`
	IncludedLocalEntryBaseDN                    types.Set    `tfsdk:"included_local_entry_base_dn"`
	ConnectionCriteria                          types.String `tfsdk:"connection_criteria"`
	RequestCriteria                             types.String `tfsdk:"request_criteria"`
}

// GetSchema defines the schema for the resource.
func (r *passThroughAuthenticationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	passThroughAuthenticationHandlerSchema(ctx, req, resp, false)
}

func (r *defaultPassThroughAuthenticationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	passThroughAuthenticationHandlerSchema(ctx, req, resp, true)
}

func passThroughAuthenticationHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Pass Through Authentication Handler.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Pass Through Authentication Handler resource. Options are ['ping-one', 'ldap', 'aggregate', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"ping-one", "ldap", "aggregate", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Pass Through Authentication Handler.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Pass Through Authentication Handler. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"subordinate_pass_through_authentication_handler": schema.SetAttribute{
				Description: "The set of subordinate pass-through authentication handlers that may be used to perform the authentication processing. Handlers will be invoked in order until one is found for which the bind operation matches the associated criteria and either succeeds or fails in a manner that should not be ignored.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"continue_on_failure_type": schema.SetAttribute{
				Description: "The set of pass-through authentication failure types that should not result in an immediate failure, but should instead allow the aggregate handler to proceed with the next configured subordinate handler.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"server": schema.SetAttribute{
				Description: "Specifies the LDAP external server(s) to which authentication attempts should be forwarded.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"server_access_mode": schema.StringAttribute{
				Description: "Specifies the manner in which external servers should be used for pass-through authentication attempts if multiple servers are defined.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dn_map": schema.SetAttribute{
				Description: "Specifies one or more DN mappings that may be used to transform bind DNs before attempting to bind to the external servers.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"bind_dn_pattern": schema.StringAttribute{
				Description: "A pattern to use to construct the bind DN for the simple bind request to send to the remote server. This may consist of a combination of static text and attribute values and other directives enclosed in curly braces.  For example, the value \"cn={cn},ou=People,dc=example,dc=com\" indicates that the remote bind DN should be constructed from the text \"cn=\" followed by the value of the local entry's cn attribute followed by the text \"ou=People,dc=example,dc=com\". If an attribute contains the value to use as the bind DN for pass-through authentication, then the pattern may simply be the name of that attribute in curly braces (e.g., if the seeAlso attribute contains the bind DN for the target user, then a bind DN pattern of \"{seeAlso}\" would be appropriate).  Note that a bind DN pattern can be used to construct a bind DN that is not actually a valid LDAP distinguished name. For example, if authentication is being passed through to a Microsoft Active Directory server, then a bind DN pattern could be used to construct a user principal name (UPN) as an alternative to a distinguished name.",
				Optional:    true,
			},
			"search_base_dn": schema.StringAttribute{
				Description: "The base DN to use when searching for the user entry using a filter constructed from the pattern defined in the search-filter-pattern property. If no base DN is specified, the null DN will be used as the search base DN.",
				Optional:    true,
			},
			"search_filter_pattern": schema.StringAttribute{
				Description: "A pattern to use to construct a filter to use when searching an external server for the entry of the user as whom to bind. For example, \"(mail={uid:ldapFilterEscape}@example.com)\" would construct a search filter to search for a user whose entry in the local server contains a uid attribute whose value appears before \"@example.com\" in the mail attribute in the external server. Note that the \"ldapFilterEscape\" modifier should almost always be used with attributes specified in the pattern.",
				Optional:    true,
			},
			"initial_connections": schema.Int64Attribute{
				Description: "Specifies the initial number of connections to establish to each external server against which authentication may be attempted.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"max_connections": schema.Int64Attribute{
				Description: "Specifies the maximum number of connections to maintain to each external server against which authentication may be attempted. This value must be greater than or equal to the value for the initial-connections property.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"use_location": schema.BoolAttribute{
				Description: "Indicates whether to take server locations into account when prioritizing the servers to use for pass-through authentication attempts.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"maximum_allowed_local_response_time": schema.StringAttribute{
				Description: "The maximum length of time to wait for a response from an external server in the same location as this Directory Server before considering it unavailable.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"maximum_allowed_nonlocal_response_time": schema.StringAttribute{
				Description: "The maximum length of time to wait for a response from an external server in a different location from this Directory Server before considering it unavailable.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"use_password_policy_control": schema.BoolAttribute{
				Description: "Indicates whether to include the password policy request control (as defined in draft-behera-ldap-password-policy-10) in bind requests sent to the external server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"api_url": schema.StringAttribute{
				Description: "Specifies the API endpoint for the PingOne web service.",
				Optional:    true,
			},
			"auth_url": schema.StringAttribute{
				Description: "Specifies the API endpoint for the PingOne authentication service.",
				Optional:    true,
			},
			"oauth_client_id": schema.StringAttribute{
				Description: "Specifies the OAuth Client ID used to authenticate connections to the PingOne API.",
				Optional:    true,
			},
			"oauth_client_secret": schema.StringAttribute{
				Description: "Specifies the OAuth Client Secret used to authenticate connections to the PingOne API.",
				Optional:    true,
				Sensitive:   true,
			},
			"oauth_client_secret_passphrase_provider": schema.StringAttribute{
				Description: "Specifies a passphrase provider that can be used to obtain the OAuth Client Secret used to authenticate connections to the PingOne API.",
				Optional:    true,
			},
			"environment_id": schema.StringAttribute{
				Description: "Specifies the PingOne Environment that will be associated with this PingOne Pass Through Authentication Handler.",
				Optional:    true,
			},
			"http_proxy_external_server": schema.StringAttribute{
				Description: "A reference to an HTTP proxy server that should be used for requests sent to the PingOne service.",
				Optional:    true,
			},
			"user_mapping_local_attribute": schema.SetAttribute{
				Description: "The names of the attributes in the local user entry whose values must match the values of the corresponding fields in the PingOne service.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"user_mapping_remote_json_field": schema.SetAttribute{
				Description: "The names of the fields in the PingOne service whose values must match the values of the corresponding attributes in the local user entry, as specified in the user-mapping-local-attribute property.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"additional_user_mapping_scim_filter": schema.StringAttribute{
				Description: "An optional SCIM filter that will be ANDed with the filter created to identify the account in the PingOne service that corresponds to the local entry. Only the \"eq\", \"sw\", \"and\", and \"or\" filter types may be used.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Pass Through Authentication Handler",
				Optional:    true,
			},
			"included_local_entry_base_dn": schema.SetAttribute{
				Description: "The base DNs for the local users whose authentication attempts may be passed through to the external authentication service.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_criteria": schema.StringAttribute{
				Description: "A reference to connection criteria that will be used to indicate which bind requests should be passed through to the external authentication service.",
				Optional:    true,
			},
			"request_criteria": schema.StringAttribute{
				Description: "A reference to request criteria that will be used to indicate which bind requests should be passed through to the external authentication service.",
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

// Validate that any restrictions are met in the plan
func (r *passThroughAuthenticationHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanPassThroughAuthenticationHandler(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_pass_through_authentication_handler")
}

func (r *defaultPassThroughAuthenticationHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanPassThroughAuthenticationHandler(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_default_pass_through_authentication_handler")
}

func modifyPlanPassThroughAuthenticationHandler(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, resourceName string) {
	compare, err := version.Compare(providerConfig.ProductVersion, version.PingDirectory9300)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model passThroughAuthenticationHandlerResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.Type) && model.Type.ValueString() == "ping-one" {
		version.CheckResourceSupported(&resp.Diagnostics, version.PingDirectory9300,
			providerConfig.ProductVersion, resourceName+" with type \"ping_one\"")
	}
	if internaltypes.IsDefined(model.Type) && model.Type.ValueString() == "aggregate" {
		version.CheckResourceSupported(&resp.Diagnostics, version.PingDirectory9300,
			providerConfig.ProductVersion, resourceName+" with type \"aggregate\"")
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsPassThroughAuthenticationHandler() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"ping-one"},
			resourcevalidator.ExactlyOneOf(
				path.MatchRoot("oauth_client_secret"),
				path.MatchRoot("oauth_client_secret_passphrase_provider"),
			),
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("additional_user_mapping_scim_filter"),
			path.MatchRoot("type"),
			[]string{"ping-one"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("server"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("api_url"),
			path.MatchRoot("type"),
			[]string{"ping-one"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("use_location"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("auth_url"),
			path.MatchRoot("type"),
			[]string{"ping-one"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("oauth_client_id"),
			path.MatchRoot("type"),
			[]string{"ping-one"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("server_access_mode"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("continue_on_failure_type"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("max_connections"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("search_base_dn"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("dn_map"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("oauth_client_secret_passphrase_provider"),
			path.MatchRoot("type"),
			[]string{"ping-one"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("user_mapping_remote_json_field"),
			path.MatchRoot("type"),
			[]string{"ping-one"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("user_mapping_local_attribute"),
			path.MatchRoot("type"),
			[]string{"ping-one"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("subordinate_pass_through_authentication_handler"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("environment_id"),
			path.MatchRoot("type"),
			[]string{"ping-one"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_argument"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("initial_connections"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("search_filter_pattern"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("http_proxy_external_server"),
			path.MatchRoot("type"),
			[]string{"ping-one"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("maximum_allowed_local_response_time"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("oauth_client_secret"),
			path.MatchRoot("type"),
			[]string{"ping-one"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("maximum_allowed_nonlocal_response_time"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_class"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("use_password_policy_control"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("bind_dn_pattern"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
	}
}

// Add config validators
func (r passThroughAuthenticationHandlerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsPassThroughAuthenticationHandler()
}

// Add config validators
func (r defaultPassThroughAuthenticationHandlerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsPassThroughAuthenticationHandler()
}

// Add optional fields to create request for ping-one pass-through-authentication-handler
func addOptionalPingOnePassThroughAuthenticationHandlerFields(ctx context.Context, addRequest *client.AddPingOnePassThroughAuthenticationHandlerRequest, plan passThroughAuthenticationHandlerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OAuthClientSecret) {
		addRequest.OAuthClientSecret = plan.OAuthClientSecret.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OAuthClientSecretPassphraseProvider) {
		addRequest.OAuthClientSecretPassphraseProvider = plan.OAuthClientSecretPassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HttpProxyExternalServer) {
		addRequest.HttpProxyExternalServer = plan.HttpProxyExternalServer.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AdditionalUserMappingSCIMFilter) {
		addRequest.AdditionalUserMappingSCIMFilter = plan.AdditionalUserMappingSCIMFilter.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludedLocalEntryBaseDN) {
		var slice []string
		plan.IncludedLocalEntryBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.IncludedLocalEntryBaseDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		addRequest.ConnectionCriteria = plan.ConnectionCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for ldap pass-through-authentication-handler
func addOptionalLdapPassThroughAuthenticationHandlerFields(ctx context.Context, addRequest *client.AddLdapPassThroughAuthenticationHandlerRequest, plan passThroughAuthenticationHandlerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ServerAccessMode) {
		serverAccessMode, err := client.NewEnumpassThroughAuthenticationHandlerServerAccessModePropFromValue(plan.ServerAccessMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.ServerAccessMode = serverAccessMode
	}
	if internaltypes.IsDefined(plan.DnMap) {
		var slice []string
		plan.DnMap.ElementsAs(ctx, &slice, false)
		addRequest.DnMap = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BindDNPattern) {
		addRequest.BindDNPattern = plan.BindDNPattern.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchBaseDN) {
		addRequest.SearchBaseDN = plan.SearchBaseDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchFilterPattern) {
		addRequest.SearchFilterPattern = plan.SearchFilterPattern.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InitialConnections) {
		addRequest.InitialConnections = plan.InitialConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxConnections) {
		addRequest.MaxConnections = plan.MaxConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.UseLocation) {
		addRequest.UseLocation = plan.UseLocation.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaximumAllowedLocalResponseTime) {
		addRequest.MaximumAllowedLocalResponseTime = plan.MaximumAllowedLocalResponseTime.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaximumAllowedNonlocalResponseTime) {
		addRequest.MaximumAllowedNonlocalResponseTime = plan.MaximumAllowedNonlocalResponseTime.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.UsePasswordPolicyControl) {
		addRequest.UsePasswordPolicyControl = plan.UsePasswordPolicyControl.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludedLocalEntryBaseDN) {
		var slice []string
		plan.IncludedLocalEntryBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.IncludedLocalEntryBaseDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		addRequest.ConnectionCriteria = plan.ConnectionCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for aggregate pass-through-authentication-handler
func addOptionalAggregatePassThroughAuthenticationHandlerFields(ctx context.Context, addRequest *client.AddAggregatePassThroughAuthenticationHandlerRequest, plan passThroughAuthenticationHandlerResourceModel) error {
	if internaltypes.IsDefined(plan.ContinueOnFailureType) {
		var slice []string
		plan.ContinueOnFailureType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpassThroughAuthenticationHandlerContinueOnFailureTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpassThroughAuthenticationHandlerContinueOnFailureTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.ContinueOnFailureType = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludedLocalEntryBaseDN) {
		var slice []string
		plan.IncludedLocalEntryBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.IncludedLocalEntryBaseDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		addRequest.ConnectionCriteria = plan.ConnectionCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for third-party pass-through-authentication-handler
func addOptionalThirdPartyPassThroughAuthenticationHandlerFields(ctx context.Context, addRequest *client.AddThirdPartyPassThroughAuthenticationHandlerRequest, plan passThroughAuthenticationHandlerResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludedLocalEntryBaseDN) {
		var slice []string
		plan.IncludedLocalEntryBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.IncludedLocalEntryBaseDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		addRequest.ConnectionCriteria = plan.ConnectionCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populatePassThroughAuthenticationHandlerUnknownValues(model *passThroughAuthenticationHandlerResourceModel) {
	if model.UserMappingRemoteJSONField.IsUnknown() || model.UserMappingRemoteJSONField.IsNull() {
		model.UserMappingRemoteJSONField, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.SubordinatePassThroughAuthenticationHandler.IsUnknown() || model.SubordinatePassThroughAuthenticationHandler.IsNull() {
		model.SubordinatePassThroughAuthenticationHandler, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.UserMappingLocalAttribute.IsUnknown() || model.UserMappingLocalAttribute.IsNull() {
		model.UserMappingLocalAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.DnMap.IsUnknown() || model.DnMap.IsNull() {
		model.DnMap, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.Server.IsUnknown() || model.Server.IsNull() {
		model.Server, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ContinueOnFailureType.IsUnknown() || model.ContinueOnFailureType.IsNull() {
		model.ContinueOnFailureType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.MaximumAllowedNonlocalResponseTime.IsUnknown() || model.MaximumAllowedNonlocalResponseTime.IsNull() {
		model.MaximumAllowedNonlocalResponseTime = types.StringValue("")
	}
	if model.MaximumAllowedLocalResponseTime.IsUnknown() || model.MaximumAllowedLocalResponseTime.IsNull() {
		model.MaximumAllowedLocalResponseTime = types.StringValue("")
	}
	if model.ServerAccessMode.IsUnknown() || model.ServerAccessMode.IsNull() {
		model.ServerAccessMode = types.StringValue("")
	}
	if model.OAuthClientSecret.IsUnknown() {
		model.OAuthClientSecret = types.StringNull()
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *passThroughAuthenticationHandlerResourceModel) populateAllComputedStringAttributes() {
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.ExtensionClass.IsUnknown() || model.ExtensionClass.IsNull() {
		model.ExtensionClass = types.StringValue("")
	}
	if model.AuthURL.IsUnknown() || model.AuthURL.IsNull() {
		model.AuthURL = types.StringValue("")
	}
	if model.HttpProxyExternalServer.IsUnknown() || model.HttpProxyExternalServer.IsNull() {
		model.HttpProxyExternalServer = types.StringValue("")
	}
	if model.SearchBaseDN.IsUnknown() || model.SearchBaseDN.IsNull() {
		model.SearchBaseDN = types.StringValue("")
	}
	if model.ConnectionCriteria.IsUnknown() || model.ConnectionCriteria.IsNull() {
		model.ConnectionCriteria = types.StringValue("")
	}
	if model.SearchFilterPattern.IsUnknown() || model.SearchFilterPattern.IsNull() {
		model.SearchFilterPattern = types.StringValue("")
	}
	if model.OAuthClientSecretPassphraseProvider.IsUnknown() || model.OAuthClientSecretPassphraseProvider.IsNull() {
		model.OAuthClientSecretPassphraseProvider = types.StringValue("")
	}
	if model.EnvironmentID.IsUnknown() || model.EnvironmentID.IsNull() {
		model.EnvironmentID = types.StringValue("")
	}
	if model.ApiURL.IsUnknown() || model.ApiURL.IsNull() {
		model.ApiURL = types.StringValue("")
	}
	if model.RequestCriteria.IsUnknown() || model.RequestCriteria.IsNull() {
		model.RequestCriteria = types.StringValue("")
	}
	if model.AdditionalUserMappingSCIMFilter.IsUnknown() || model.AdditionalUserMappingSCIMFilter.IsNull() {
		model.AdditionalUserMappingSCIMFilter = types.StringValue("")
	}
	if model.BindDNPattern.IsUnknown() || model.BindDNPattern.IsNull() {
		model.BindDNPattern = types.StringValue("")
	}
	if model.OAuthClientID.IsUnknown() || model.OAuthClientID.IsNull() {
		model.OAuthClientID = types.StringValue("")
	}
}

// Read a PingOnePassThroughAuthenticationHandlerResponse object into the model struct
func readPingOnePassThroughAuthenticationHandlerResponse(ctx context.Context, r *client.PingOnePassThroughAuthenticationHandlerResponse, state *passThroughAuthenticationHandlerResourceModel, expectedValues *passThroughAuthenticationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ping-one")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ApiURL = types.StringValue(r.ApiURL)
	state.AuthURL = types.StringValue(r.AuthURL)
	state.OAuthClientID = types.StringValue(r.OAuthClientID)
	state.OAuthClientSecretPassphraseProvider = internaltypes.StringTypeOrNil(r.OAuthClientSecretPassphraseProvider, internaltypes.IsEmptyString(expectedValues.OAuthClientSecretPassphraseProvider))
	state.EnvironmentID = types.StringValue(r.EnvironmentID)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, internaltypes.IsEmptyString(expectedValues.HttpProxyExternalServer))
	state.UserMappingLocalAttribute = internaltypes.GetStringSet(r.UserMappingLocalAttribute)
	state.UserMappingRemoteJSONField = internaltypes.GetStringSet(r.UserMappingRemoteJSONField)
	state.AdditionalUserMappingSCIMFilter = internaltypes.StringTypeOrNil(r.AdditionalUserMappingSCIMFilter, internaltypes.IsEmptyString(expectedValues.AdditionalUserMappingSCIMFilter))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.IncludedLocalEntryBaseDN = internaltypes.GetStringSet(r.IncludedLocalEntryBaseDN)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePassThroughAuthenticationHandlerUnknownValues(state)
}

// Read a LdapPassThroughAuthenticationHandlerResponse object into the model struct
func readLdapPassThroughAuthenticationHandlerResponse(ctx context.Context, r *client.LdapPassThroughAuthenticationHandlerResponse, state *passThroughAuthenticationHandlerResourceModel, expectedValues *passThroughAuthenticationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Server = internaltypes.GetStringSet(r.Server)
	state.ServerAccessMode = types.StringValue(r.ServerAccessMode.String())
	state.DnMap = internaltypes.GetStringSet(r.DnMap)
	state.BindDNPattern = internaltypes.StringTypeOrNil(r.BindDNPattern, internaltypes.IsEmptyString(expectedValues.BindDNPattern))
	state.SearchBaseDN = internaltypes.StringTypeOrNil(r.SearchBaseDN, internaltypes.IsEmptyString(expectedValues.SearchBaseDN))
	state.SearchFilterPattern = internaltypes.StringTypeOrNil(r.SearchFilterPattern, internaltypes.IsEmptyString(expectedValues.SearchFilterPattern))
	state.InitialConnections = types.Int64Value(r.InitialConnections)
	state.MaxConnections = types.Int64Value(r.MaxConnections)
	state.UseLocation = internaltypes.BoolTypeOrNil(r.UseLocation)
	state.MaximumAllowedLocalResponseTime = internaltypes.StringTypeOrNil(r.MaximumAllowedLocalResponseTime, true)
	config.CheckMismatchedPDFormattedAttributes("maximum_allowed_local_response_time",
		expectedValues.MaximumAllowedLocalResponseTime, state.MaximumAllowedLocalResponseTime, diagnostics)
	state.MaximumAllowedNonlocalResponseTime = internaltypes.StringTypeOrNil(r.MaximumAllowedNonlocalResponseTime, true)
	config.CheckMismatchedPDFormattedAttributes("maximum_allowed_nonlocal_response_time",
		expectedValues.MaximumAllowedNonlocalResponseTime, state.MaximumAllowedNonlocalResponseTime, diagnostics)
	state.UsePasswordPolicyControl = internaltypes.BoolTypeOrNil(r.UsePasswordPolicyControl)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.IncludedLocalEntryBaseDN = internaltypes.GetStringSet(r.IncludedLocalEntryBaseDN)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePassThroughAuthenticationHandlerUnknownValues(state)
}

// Read a AggregatePassThroughAuthenticationHandlerResponse object into the model struct
func readAggregatePassThroughAuthenticationHandlerResponse(ctx context.Context, r *client.AggregatePassThroughAuthenticationHandlerResponse, state *passThroughAuthenticationHandlerResourceModel, expectedValues *passThroughAuthenticationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("aggregate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SubordinatePassThroughAuthenticationHandler = internaltypes.GetStringSet(r.SubordinatePassThroughAuthenticationHandler)
	state.ContinueOnFailureType = internaltypes.GetStringSet(
		client.StringSliceEnumpassThroughAuthenticationHandlerContinueOnFailureTypeProp(r.ContinueOnFailureType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.IncludedLocalEntryBaseDN = internaltypes.GetStringSet(r.IncludedLocalEntryBaseDN)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePassThroughAuthenticationHandlerUnknownValues(state)
}

// Read a ThirdPartyPassThroughAuthenticationHandlerResponse object into the model struct
func readThirdPartyPassThroughAuthenticationHandlerResponse(ctx context.Context, r *client.ThirdPartyPassThroughAuthenticationHandlerResponse, state *passThroughAuthenticationHandlerResourceModel, expectedValues *passThroughAuthenticationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.IncludedLocalEntryBaseDN = internaltypes.GetStringSet(r.IncludedLocalEntryBaseDN)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePassThroughAuthenticationHandlerUnknownValues(state)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *passThroughAuthenticationHandlerResourceModel) setStateValuesNotReturnedByAPI(expectedValues *passThroughAuthenticationHandlerResourceModel) {
	if !expectedValues.OAuthClientSecret.IsUnknown() {
		state.OAuthClientSecret = expectedValues.OAuthClientSecret
	}
}

// Create any update operations necessary to make the state match the plan
func createPassThroughAuthenticationHandlerOperations(plan passThroughAuthenticationHandlerResourceModel, state passThroughAuthenticationHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SubordinatePassThroughAuthenticationHandler, state.SubordinatePassThroughAuthenticationHandler, "subordinate-pass-through-authentication-handler")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ContinueOnFailureType, state.ContinueOnFailureType, "continue-on-failure-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Server, state.Server, "server")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerAccessMode, state.ServerAccessMode, "server-access-mode")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DnMap, state.DnMap, "dn-map")
	operations.AddStringOperationIfNecessary(&ops, plan.BindDNPattern, state.BindDNPattern, "bind-dn-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchBaseDN, state.SearchBaseDN, "search-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchFilterPattern, state.SearchFilterPattern, "search-filter-pattern")
	operations.AddInt64OperationIfNecessary(&ops, plan.InitialConnections, state.InitialConnections, "initial-connections")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxConnections, state.MaxConnections, "max-connections")
	operations.AddBoolOperationIfNecessary(&ops, plan.UseLocation, state.UseLocation, "use-location")
	operations.AddStringOperationIfNecessary(&ops, plan.MaximumAllowedLocalResponseTime, state.MaximumAllowedLocalResponseTime, "maximum-allowed-local-response-time")
	operations.AddStringOperationIfNecessary(&ops, plan.MaximumAllowedNonlocalResponseTime, state.MaximumAllowedNonlocalResponseTime, "maximum-allowed-nonlocal-response-time")
	operations.AddBoolOperationIfNecessary(&ops, plan.UsePasswordPolicyControl, state.UsePasswordPolicyControl, "use-password-policy-control")
	operations.AddStringOperationIfNecessary(&ops, plan.ApiURL, state.ApiURL, "api-url")
	operations.AddStringOperationIfNecessary(&ops, plan.AuthURL, state.AuthURL, "auth-url")
	operations.AddStringOperationIfNecessary(&ops, plan.OAuthClientID, state.OAuthClientID, "oauth-client-id")
	operations.AddStringOperationIfNecessary(&ops, plan.OAuthClientSecret, state.OAuthClientSecret, "oauth-client-secret")
	operations.AddStringOperationIfNecessary(&ops, plan.OAuthClientSecretPassphraseProvider, state.OAuthClientSecretPassphraseProvider, "oauth-client-secret-passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.EnvironmentID, state.EnvironmentID, "environment-id")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpProxyExternalServer, state.HttpProxyExternalServer, "http-proxy-external-server")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.UserMappingLocalAttribute, state.UserMappingLocalAttribute, "user-mapping-local-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.UserMappingRemoteJSONField, state.UserMappingRemoteJSONField, "user-mapping-remote-json-field")
	operations.AddStringOperationIfNecessary(&ops, plan.AdditionalUserMappingSCIMFilter, state.AdditionalUserMappingSCIMFilter, "additional-user-mapping-scim-filter")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedLocalEntryBaseDN, state.IncludedLocalEntryBaseDN, "included-local-entry-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionCriteria, state.ConnectionCriteria, "connection-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.RequestCriteria, state.RequestCriteria, "request-criteria")
	return ops
}

// Create a ping-one pass-through-authentication-handler
func (r *passThroughAuthenticationHandlerResource) CreatePingOnePassThroughAuthenticationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passThroughAuthenticationHandlerResourceModel) (*passThroughAuthenticationHandlerResourceModel, error) {
	var UserMappingLocalAttributeSlice []string
	plan.UserMappingLocalAttribute.ElementsAs(ctx, &UserMappingLocalAttributeSlice, false)
	var UserMappingRemoteJSONFieldSlice []string
	plan.UserMappingRemoteJSONField.ElementsAs(ctx, &UserMappingRemoteJSONFieldSlice, false)
	addRequest := client.NewAddPingOnePassThroughAuthenticationHandlerRequest(plan.Name.ValueString(),
		[]client.EnumpingOnePassThroughAuthenticationHandlerSchemaUrn{client.ENUMPINGONEPASSTHROUGHAUTHENTICATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASS_THROUGH_AUTHENTICATION_HANDLERPING_ONE},
		plan.ApiURL.ValueString(),
		plan.AuthURL.ValueString(),
		plan.OAuthClientID.ValueString(),
		plan.EnvironmentID.ValueString(),
		UserMappingLocalAttributeSlice,
		UserMappingRemoteJSONFieldSlice)
	err := addOptionalPingOnePassThroughAuthenticationHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Pass Through Authentication Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PassThroughAuthenticationHandlerApi.AddPassThroughAuthenticationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPassThroughAuthenticationHandlerRequest(
		client.AddPingOnePassThroughAuthenticationHandlerRequestAsAddPassThroughAuthenticationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PassThroughAuthenticationHandlerApi.AddPassThroughAuthenticationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Pass Through Authentication Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passThroughAuthenticationHandlerResourceModel
	readPingOnePassThroughAuthenticationHandlerResponse(ctx, addResponse.PingOnePassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a ldap pass-through-authentication-handler
func (r *passThroughAuthenticationHandlerResource) CreateLdapPassThroughAuthenticationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passThroughAuthenticationHandlerResourceModel) (*passThroughAuthenticationHandlerResourceModel, error) {
	var ServerSlice []string
	plan.Server.ElementsAs(ctx, &ServerSlice, false)
	addRequest := client.NewAddLdapPassThroughAuthenticationHandlerRequest(plan.Name.ValueString(),
		[]client.EnumldapPassThroughAuthenticationHandlerSchemaUrn{client.ENUMLDAPPASSTHROUGHAUTHENTICATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASS_THROUGH_AUTHENTICATION_HANDLERLDAP},
		ServerSlice)
	err := addOptionalLdapPassThroughAuthenticationHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Pass Through Authentication Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PassThroughAuthenticationHandlerApi.AddPassThroughAuthenticationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPassThroughAuthenticationHandlerRequest(
		client.AddLdapPassThroughAuthenticationHandlerRequestAsAddPassThroughAuthenticationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PassThroughAuthenticationHandlerApi.AddPassThroughAuthenticationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Pass Through Authentication Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passThroughAuthenticationHandlerResourceModel
	readLdapPassThroughAuthenticationHandlerResponse(ctx, addResponse.LdapPassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a aggregate pass-through-authentication-handler
func (r *passThroughAuthenticationHandlerResource) CreateAggregatePassThroughAuthenticationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passThroughAuthenticationHandlerResourceModel) (*passThroughAuthenticationHandlerResourceModel, error) {
	var SubordinatePassThroughAuthenticationHandlerSlice []string
	plan.SubordinatePassThroughAuthenticationHandler.ElementsAs(ctx, &SubordinatePassThroughAuthenticationHandlerSlice, false)
	addRequest := client.NewAddAggregatePassThroughAuthenticationHandlerRequest(plan.Name.ValueString(),
		[]client.EnumaggregatePassThroughAuthenticationHandlerSchemaUrn{client.ENUMAGGREGATEPASSTHROUGHAUTHENTICATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASS_THROUGH_AUTHENTICATION_HANDLERAGGREGATE},
		SubordinatePassThroughAuthenticationHandlerSlice)
	err := addOptionalAggregatePassThroughAuthenticationHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Pass Through Authentication Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PassThroughAuthenticationHandlerApi.AddPassThroughAuthenticationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPassThroughAuthenticationHandlerRequest(
		client.AddAggregatePassThroughAuthenticationHandlerRequestAsAddPassThroughAuthenticationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PassThroughAuthenticationHandlerApi.AddPassThroughAuthenticationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Pass Through Authentication Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passThroughAuthenticationHandlerResourceModel
	readAggregatePassThroughAuthenticationHandlerResponse(ctx, addResponse.AggregatePassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party pass-through-authentication-handler
func (r *passThroughAuthenticationHandlerResource) CreateThirdPartyPassThroughAuthenticationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passThroughAuthenticationHandlerResourceModel) (*passThroughAuthenticationHandlerResourceModel, error) {
	addRequest := client.NewAddThirdPartyPassThroughAuthenticationHandlerRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyPassThroughAuthenticationHandlerSchemaUrn{client.ENUMTHIRDPARTYPASSTHROUGHAUTHENTICATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASS_THROUGH_AUTHENTICATION_HANDLERTHIRD_PARTY},
		plan.ExtensionClass.ValueString())
	err := addOptionalThirdPartyPassThroughAuthenticationHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Pass Through Authentication Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PassThroughAuthenticationHandlerApi.AddPassThroughAuthenticationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPassThroughAuthenticationHandlerRequest(
		client.AddThirdPartyPassThroughAuthenticationHandlerRequestAsAddPassThroughAuthenticationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PassThroughAuthenticationHandlerApi.AddPassThroughAuthenticationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Pass Through Authentication Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passThroughAuthenticationHandlerResourceModel
	readThirdPartyPassThroughAuthenticationHandlerResponse(ctx, addResponse.ThirdPartyPassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *passThroughAuthenticationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan passThroughAuthenticationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *passThroughAuthenticationHandlerResourceModel
	var err error
	if plan.Type.ValueString() == "ping-one" {
		state, err = r.CreatePingOnePassThroughAuthenticationHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "ldap" {
		state, err = r.CreateLdapPassThroughAuthenticationHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "aggregate" {
		state, err = r.CreateAggregatePassThroughAuthenticationHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyPassThroughAuthenticationHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

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
func (r *defaultPassThroughAuthenticationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan passThroughAuthenticationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PassThroughAuthenticationHandlerApi.GetPassThroughAuthenticationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Pass Through Authentication Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state passThroughAuthenticationHandlerResourceModel
	if readResponse.PingOnePassThroughAuthenticationHandlerResponse != nil {
		readPingOnePassThroughAuthenticationHandlerResponse(ctx, readResponse.PingOnePassThroughAuthenticationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LdapPassThroughAuthenticationHandlerResponse != nil {
		readLdapPassThroughAuthenticationHandlerResponse(ctx, readResponse.LdapPassThroughAuthenticationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AggregatePassThroughAuthenticationHandlerResponse != nil {
		readAggregatePassThroughAuthenticationHandlerResponse(ctx, readResponse.AggregatePassThroughAuthenticationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPassThroughAuthenticationHandlerResponse != nil {
		readThirdPartyPassThroughAuthenticationHandlerResponse(ctx, readResponse.ThirdPartyPassThroughAuthenticationHandlerResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PassThroughAuthenticationHandlerApi.UpdatePassThroughAuthenticationHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createPassThroughAuthenticationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PassThroughAuthenticationHandlerApi.UpdatePassThroughAuthenticationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Pass Through Authentication Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.PingOnePassThroughAuthenticationHandlerResponse != nil {
			readPingOnePassThroughAuthenticationHandlerResponse(ctx, updateResponse.PingOnePassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LdapPassThroughAuthenticationHandlerResponse != nil {
			readLdapPassThroughAuthenticationHandlerResponse(ctx, updateResponse.LdapPassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AggregatePassThroughAuthenticationHandlerResponse != nil {
			readAggregatePassThroughAuthenticationHandlerResponse(ctx, updateResponse.AggregatePassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyPassThroughAuthenticationHandlerResponse != nil {
			readThirdPartyPassThroughAuthenticationHandlerResponse(ctx, updateResponse.ThirdPartyPassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
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
func (r *passThroughAuthenticationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPassThroughAuthenticationHandler(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultPassThroughAuthenticationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPassThroughAuthenticationHandler(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readPassThroughAuthenticationHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state passThroughAuthenticationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.PassThroughAuthenticationHandlerApi.GetPassThroughAuthenticationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Pass Through Authentication Handler", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Pass Through Authentication Handler", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.PingOnePassThroughAuthenticationHandlerResponse != nil {
		readPingOnePassThroughAuthenticationHandlerResponse(ctx, readResponse.PingOnePassThroughAuthenticationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LdapPassThroughAuthenticationHandlerResponse != nil {
		readLdapPassThroughAuthenticationHandlerResponse(ctx, readResponse.LdapPassThroughAuthenticationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AggregatePassThroughAuthenticationHandlerResponse != nil {
		readAggregatePassThroughAuthenticationHandlerResponse(ctx, readResponse.AggregatePassThroughAuthenticationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPassThroughAuthenticationHandlerResponse != nil {
		readThirdPartyPassThroughAuthenticationHandlerResponse(ctx, readResponse.ThirdPartyPassThroughAuthenticationHandlerResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *passThroughAuthenticationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePassThroughAuthenticationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPassThroughAuthenticationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePassThroughAuthenticationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updatePassThroughAuthenticationHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan passThroughAuthenticationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state passThroughAuthenticationHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.PassThroughAuthenticationHandlerApi.UpdatePassThroughAuthenticationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createPassThroughAuthenticationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.PassThroughAuthenticationHandlerApi.UpdatePassThroughAuthenticationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Pass Through Authentication Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.PingOnePassThroughAuthenticationHandlerResponse != nil {
			readPingOnePassThroughAuthenticationHandlerResponse(ctx, updateResponse.PingOnePassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LdapPassThroughAuthenticationHandlerResponse != nil {
			readLdapPassThroughAuthenticationHandlerResponse(ctx, updateResponse.LdapPassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AggregatePassThroughAuthenticationHandlerResponse != nil {
			readAggregatePassThroughAuthenticationHandlerResponse(ctx, updateResponse.AggregatePassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyPassThroughAuthenticationHandlerResponse != nil {
			readThirdPartyPassThroughAuthenticationHandlerResponse(ctx, updateResponse.ThirdPartyPassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
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
func (r *defaultPassThroughAuthenticationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *passThroughAuthenticationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state passThroughAuthenticationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PassThroughAuthenticationHandlerApi.DeletePassThroughAuthenticationHandlerExecute(r.apiClient.PassThroughAuthenticationHandlerApi.DeletePassThroughAuthenticationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Pass Through Authentication Handler", err, httpResp)
		return
	}
}

func (r *passThroughAuthenticationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPassThroughAuthenticationHandler(ctx, req, resp)
}

func (r *defaultPassThroughAuthenticationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPassThroughAuthenticationHandler(ctx, req, resp)
}

func importPassThroughAuthenticationHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
