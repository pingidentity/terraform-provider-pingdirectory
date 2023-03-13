package config

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &topologyAdminUserResource{}
	_ resource.ResourceWithConfigure   = &topologyAdminUserResource{}
	_ resource.ResourceWithImportState = &topologyAdminUserResource{}
	_ resource.Resource                = &defaultTopologyAdminUserResource{}
	_ resource.ResourceWithConfigure   = &defaultTopologyAdminUserResource{}
	_ resource.ResourceWithImportState = &defaultTopologyAdminUserResource{}
)

// Create a Topology Admin User resource
func NewTopologyAdminUserResource() resource.Resource {
	return &topologyAdminUserResource{}
}

func NewDefaultTopologyAdminUserResource() resource.Resource {
	return &defaultTopologyAdminUserResource{}
}

// topologyAdminUserResource is the resource implementation.
type topologyAdminUserResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultTopologyAdminUserResource is the resource implementation.
type defaultTopologyAdminUserResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *topologyAdminUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_topology_admin_user"
}

func (r *defaultTopologyAdminUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_topology_admin_user"
}

// Configure adds the provider configured client to the resource.
func (r *topologyAdminUserResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultTopologyAdminUserResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type topologyAdminUserResourceModel struct {
	Id                             types.String `tfsdk:"id"`
	LastUpdated                    types.String `tfsdk:"last_updated"`
	Notifications                  types.Set    `tfsdk:"notifications"`
	RequiredActions                types.Set    `tfsdk:"required_actions"`
	AlternateBindDN                types.Set    `tfsdk:"alternate_bind_dn"`
	Description                    types.String `tfsdk:"description"`
	Password                       types.String `tfsdk:"password"`
	FirstName                      types.Set    `tfsdk:"first_name"`
	LastName                       types.Set    `tfsdk:"last_name"`
	UserID                         types.String `tfsdk:"user_id"`
	EmailAddress                   types.Set    `tfsdk:"email_address"`
	WorkTelephoneNumber            types.Set    `tfsdk:"work_telephone_number"`
	HomeTelephoneNumber            types.Set    `tfsdk:"home_telephone_number"`
	MobileTelephoneNumber          types.Set    `tfsdk:"mobile_telephone_number"`
	PagerTelephoneNumber           types.Set    `tfsdk:"pager_telephone_number"`
	InheritDefaultRootPrivileges   types.Bool   `tfsdk:"inherit_default_root_privileges"`
	Privilege                      types.Set    `tfsdk:"privilege"`
	SearchResultEntryLimit         types.Int64  `tfsdk:"search_result_entry_limit"`
	TimeLimitSeconds               types.Int64  `tfsdk:"time_limit_seconds"`
	LookThroughEntryLimit          types.Int64  `tfsdk:"look_through_entry_limit"`
	IdleTimeLimitSeconds           types.Int64  `tfsdk:"idle_time_limit_seconds"`
	PasswordPolicy                 types.String `tfsdk:"password_policy"`
	Disabled                       types.Bool   `tfsdk:"disabled"`
	AccountActivationTime          types.String `tfsdk:"account_activation_time"`
	AccountExpirationTime          types.String `tfsdk:"account_expiration_time"`
	RequireSecureAuthentication    types.Bool   `tfsdk:"require_secure_authentication"`
	RequireSecureConnections       types.Bool   `tfsdk:"require_secure_connections"`
	AllowedAuthenticationType      types.Set    `tfsdk:"allowed_authentication_type"`
	AllowedAuthenticationIPAddress types.Set    `tfsdk:"allowed_authentication_ip_address"`
	PreferredOTPDeliveryMechanism  types.Set    `tfsdk:"preferred_otp_delivery_mechanism"`
	IsProxyable                    types.String `tfsdk:"is_proxyable"`
	IsProxyableByDN                types.Set    `tfsdk:"is_proxyable_by_dn"`
	IsProxyableByGroup             types.Set    `tfsdk:"is_proxyable_by_group"`
	IsProxyableByURL               types.Set    `tfsdk:"is_proxyable_by_url"`
	MayProxyAsDN                   types.Set    `tfsdk:"may_proxy_as_dn"`
	MayProxyAsGroup                types.Set    `tfsdk:"may_proxy_as_group"`
	MayProxyAsURL                  types.Set    `tfsdk:"may_proxy_as_url"`
}

// GetSchema defines the schema for the resource.
func (r *topologyAdminUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	topologyAdminUserSchema(ctx, req, resp, false)
}

func (r *defaultTopologyAdminUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	topologyAdminUserSchema(ctx, req, resp, true)
}

func topologyAdminUserSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Topology Admin User.",
		Attributes: map[string]schema.Attribute{
			"alternate_bind_dn": schema.SetAttribute{
				Description: "Specifies one or more alternate DNs that can be used to bind to the server as this User.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this User.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Specifies the user's password. This is stored in the userPassword LDAP attribute. To set a pre-hashed value, the account making the change must have the bypass-pw-policy privilege.",
				Optional:    true,
				Sensitive:   true,
			},
			"first_name": schema.SetAttribute{
				Description: "Specifies the user's first name. This is stored in the givenName LDAP attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"last_name": schema.SetAttribute{
				Description: "Specifies the user's last name. This is stored in the sn LDAP attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"user_id": schema.StringAttribute{
				Description: "Specifies the user's user ID. This is stored in the uid LDAP attribute.",
				Optional:    true,
			},
			"email_address": schema.SetAttribute{
				Description: "Specifies the user's email address. This is stored in the mail LDAP attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"work_telephone_number": schema.SetAttribute{
				Description: "Specifies the user's work telephone number. This is stored in the telephoneNumber LDAP attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"home_telephone_number": schema.SetAttribute{
				Description: "Specifies the user's home telephone number. This is stored in the homePhone LDAP attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"mobile_telephone_number": schema.SetAttribute{
				Description: "Specifies the user's mobile telephone number. This is stored in the mobile LDAP attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"pager_telephone_number": schema.SetAttribute{
				Description: "Specifies the user's pager telephone number. This is stored in the pager LDAP attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"inherit_default_root_privileges": schema.BoolAttribute{
				Description: "Indicates whether this User should be automatically granted the set of privileges defined in the default-root-privilege-name property of the Root DN configuration object.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"privilege": schema.SetAttribute{
				Description: "Privileges that are either explicitly granted or revoked from the root user. Privileges can be revoked by including a minus sign (-) before the privilege name. This is stored in the ds-privilege-name LDAP attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"search_result_entry_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that the server may return to the user in response to any single search request. A value of 0 indicates no limit should be enforced. This is stored in the ds-rlim-size-limit LDAP attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"time_limit_seconds": schema.Int64Attribute{
				Description: "Specifies the maximum length of time (in seconds) that the server may spend processing any single search request. A value of 0 indicates no limit should be enforced. This is stored in the ds-rlim-time-limit LDAP attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"look_through_entry_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of candidate entries that the server may examine in the course of processing any single search request. A value of 0 indicates no limit should be enforced. This is stored in the ds-rlim-lookthrough-limit LDAP attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"idle_time_limit_seconds": schema.Int64Attribute{
				Description: "Specifies the maximum length of time (in seconds) that a connection authenticated as this user may remain established without issuing any requests. A value of 0 indicates no limit should be enforced. This is stored in the ds-rlim-idle-time-limit LDAP attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"password_policy": schema.StringAttribute{
				Description: "Specifies the password policy for the user. This is stored in the ds-pwp-password-policy-dn LDAP attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"disabled": schema.BoolAttribute{
				Description: "Specifies whether the root user account should be disabled. A disabled account is not permitted to authenticate, nor can it be used as an authorization identity. This is stored in the ds-pwp-account-disabled LDAP attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"account_activation_time": schema.StringAttribute{
				Description: "Specifies the time, in generalized time format (e.g., '20160101070000Z'), that the root user account should become active. If an activation time is specified, the user will not be permitted to authenticate, nor can the account be used as an authorization identity, until the activation time has arrived. This is stored in the ds-pwp-account-activation-time LDAP attribute.",
				Optional:    true,
			},
			"account_expiration_time": schema.StringAttribute{
				Description: "Specifies the time, in generalized time format (e.g., '20240101070000Z'), that the root user account should expire. If an expiration time is specified, the user will not be permitted to authenticate, nor can the account be used as an authorization identity, after this time has passed. This is stored in the ds-pwp-account-expiration-time LDAP attribute.",
				Optional:    true,
			},
			"require_secure_authentication": schema.BoolAttribute{
				Description: "Indicates whether this User must authenticate in a secure manner. When set to \"true\", the User will only be allowed to authenticate over a secure connection or using a mechanism that does not expose user credentials (e.g., the CRAM-MD5, DIGEST-MD5, and GSSAPI SASL mechanisms).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"require_secure_connections": schema.BoolAttribute{
				Description: "Indicates whether this User must be required to communicate with the server over a secure connection. When set to \"true\", the User will only be allowed to communicate with the server over a secure connection (i.e., using TLS or the StartTLS extended operation).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_authentication_type": schema.SetAttribute{
				Description: "Indicates that User should only be allowed to authenticate in certain ways. Allowed values include \"simple\" (to indicate that the user should be allowed to bind using simple authentication) or \"sasl {mech}\" (to indicate that the user should be allowed to bind using the specified SASL mechanism, like \"sasl PLAIN\"). The list of available SASL mechanisms can be retrieved by running \"dsconfig --advanced list-sasl-mechanism-handlers\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"allowed_authentication_ip_address": schema.SetAttribute{
				Description: "An IPv4 or IPv6 address mask that controls the set of IP addresses from which this User can authenticate to the server. For instance a value of 127.0.0.1 (or ::1 in IPv6) would restricted access only to localhost connections, whereas 10.6.1.* would restrict access to servers on the 10.6.1.* subnet.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"preferred_otp_delivery_mechanism": schema.SetAttribute{
				Description: "Overrides the default settings for the mechanisms (e.g., email or SMS) that are used to deliver one time passwords to Users.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"is_proxyable": schema.StringAttribute{
				Description: "This can be used to indicate whether the User can be used as an alternate authorization identity (using the proxied authorization v1 or v2 control, the intermediate client control, or a SASL mechanism that allows specifying an alternate authorization identity).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_proxyable_by_dn": schema.SetAttribute{
				Description: "Specifies the DNs of accounts that can proxy as this User using the proxied authorization v1 or v2 control, the intermediate client control, or a SASL mechanism that allows specifying an alternate authorization identity. This property is only applicable if is-proxyable is set to \"allowed\" or \"required\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"is_proxyable_by_group": schema.SetAttribute{
				Description: "Specifies the DNs of groups whose members can proxy as this User using the proxied authorization v1 or v2 control, the intermediate client control, or a SASL mechanism that allows specifying an alternate authorization identity. This property is only applicable if is-proxyable is set to \"allowed\" or \"required\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"is_proxyable_by_url": schema.SetAttribute{
				Description: "Specifies LDAP URLs of accounts that can proxy as this User using the proxied authorization v1 or v2 control, the intermediate client control, or a SASL mechanism that allows specifying an alternate authorization identity. This property is only applicable if is-proxyable is set to \"allowed\" or \"required\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"may_proxy_as_dn": schema.SetAttribute{
				Description: "This restricts the set of accounts that this User can proxy as to entries with the specified DNs.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"may_proxy_as_group": schema.SetAttribute{
				Description: "This restricts the set of accounts that this User can proxy as to entries that are in the group with the specified DN.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"may_proxy_as_url": schema.SetAttribute{
				Description: "This restricts the set of accounts that this User can proxy as to entries that are matched by the specified LDAP URL.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
		},
	}
	if setOptionalToComputed {
		SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalTopologyAdminUserFields(ctx context.Context, addRequest *client.AddTopologyAdminUserRequest, plan topologyAdminUserResourceModel) error {
	if internaltypes.IsDefined(plan.AlternateBindDN) {
		var slice []string
		plan.AlternateBindDN.ElementsAs(ctx, &slice, false)
		addRequest.AlternateBindDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Password) {
		stringVal := plan.Password.ValueString()
		addRequest.Password = &stringVal
	}
	if internaltypes.IsDefined(plan.FirstName) {
		var slice []string
		plan.FirstName.ElementsAs(ctx, &slice, false)
		addRequest.FirstName = slice
	}
	if internaltypes.IsDefined(plan.LastName) {
		var slice []string
		plan.LastName.ElementsAs(ctx, &slice, false)
		addRequest.LastName = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UserID) {
		stringVal := plan.UserID.ValueString()
		addRequest.UserID = &stringVal
	}
	if internaltypes.IsDefined(plan.EmailAddress) {
		var slice []string
		plan.EmailAddress.ElementsAs(ctx, &slice, false)
		addRequest.EmailAddress = slice
	}
	if internaltypes.IsDefined(plan.WorkTelephoneNumber) {
		var slice []string
		plan.WorkTelephoneNumber.ElementsAs(ctx, &slice, false)
		addRequest.WorkTelephoneNumber = slice
	}
	if internaltypes.IsDefined(plan.HomeTelephoneNumber) {
		var slice []string
		plan.HomeTelephoneNumber.ElementsAs(ctx, &slice, false)
		addRequest.HomeTelephoneNumber = slice
	}
	if internaltypes.IsDefined(plan.MobileTelephoneNumber) {
		var slice []string
		plan.MobileTelephoneNumber.ElementsAs(ctx, &slice, false)
		addRequest.MobileTelephoneNumber = slice
	}
	if internaltypes.IsDefined(plan.PagerTelephoneNumber) {
		var slice []string
		plan.PagerTelephoneNumber.ElementsAs(ctx, &slice, false)
		addRequest.PagerTelephoneNumber = slice
	}
	if internaltypes.IsDefined(plan.InheritDefaultRootPrivileges) {
		boolVal := plan.InheritDefaultRootPrivileges.ValueBool()
		addRequest.InheritDefaultRootPrivileges = &boolVal
	}
	if internaltypes.IsDefined(plan.Privilege) {
		var slice []string
		plan.Privilege.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumtopologyAdminUserPrivilegeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumtopologyAdminUserPrivilegePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.Privilege = enumSlice
	}
	if internaltypes.IsDefined(plan.SearchResultEntryLimit) {
		intVal := int32(plan.SearchResultEntryLimit.ValueInt64())
		addRequest.SearchResultEntryLimit = &intVal
	}
	if internaltypes.IsDefined(plan.TimeLimitSeconds) {
		intVal := int32(plan.TimeLimitSeconds.ValueInt64())
		addRequest.TimeLimitSeconds = &intVal
	}
	if internaltypes.IsDefined(plan.LookThroughEntryLimit) {
		intVal := int32(plan.LookThroughEntryLimit.ValueInt64())
		addRequest.LookThroughEntryLimit = &intVal
	}
	if internaltypes.IsDefined(plan.IdleTimeLimitSeconds) {
		intVal := int32(plan.IdleTimeLimitSeconds.ValueInt64())
		addRequest.IdleTimeLimitSeconds = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PasswordPolicy) {
		stringVal := plan.PasswordPolicy.ValueString()
		addRequest.PasswordPolicy = &stringVal
	}
	if internaltypes.IsDefined(plan.Disabled) {
		boolVal := plan.Disabled.ValueBool()
		addRequest.Disabled = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountActivationTime) {
		stringVal := plan.AccountActivationTime.ValueString()
		addRequest.AccountActivationTime = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountExpirationTime) {
		stringVal := plan.AccountExpirationTime.ValueString()
		addRequest.AccountExpirationTime = &stringVal
	}
	if internaltypes.IsDefined(plan.RequireSecureAuthentication) {
		boolVal := plan.RequireSecureAuthentication.ValueBool()
		addRequest.RequireSecureAuthentication = &boolVal
	}
	if internaltypes.IsDefined(plan.RequireSecureConnections) {
		boolVal := plan.RequireSecureConnections.ValueBool()
		addRequest.RequireSecureConnections = &boolVal
	}
	if internaltypes.IsDefined(plan.AllowedAuthenticationType) {
		var slice []string
		plan.AllowedAuthenticationType.ElementsAs(ctx, &slice, false)
		addRequest.AllowedAuthenticationType = slice
	}
	if internaltypes.IsDefined(plan.AllowedAuthenticationIPAddress) {
		var slice []string
		plan.AllowedAuthenticationIPAddress.ElementsAs(ctx, &slice, false)
		addRequest.AllowedAuthenticationIPAddress = slice
	}
	if internaltypes.IsDefined(plan.PreferredOTPDeliveryMechanism) {
		var slice []string
		plan.PreferredOTPDeliveryMechanism.ElementsAs(ctx, &slice, false)
		addRequest.PreferredOTPDeliveryMechanism = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IsProxyable) {
		isProxyable, err := client.NewEnumtopologyAdminUserIsProxyablePropFromValue(plan.IsProxyable.ValueString())
		if err != nil {
			return err
		}
		addRequest.IsProxyable = isProxyable
	}
	if internaltypes.IsDefined(plan.IsProxyableByDN) {
		var slice []string
		plan.IsProxyableByDN.ElementsAs(ctx, &slice, false)
		addRequest.IsProxyableByDN = slice
	}
	if internaltypes.IsDefined(plan.IsProxyableByGroup) {
		var slice []string
		plan.IsProxyableByGroup.ElementsAs(ctx, &slice, false)
		addRequest.IsProxyableByGroup = slice
	}
	if internaltypes.IsDefined(plan.IsProxyableByURL) {
		var slice []string
		plan.IsProxyableByURL.ElementsAs(ctx, &slice, false)
		addRequest.IsProxyableByURL = slice
	}
	if internaltypes.IsDefined(plan.MayProxyAsDN) {
		var slice []string
		plan.MayProxyAsDN.ElementsAs(ctx, &slice, false)
		addRequest.MayProxyAsDN = slice
	}
	if internaltypes.IsDefined(plan.MayProxyAsGroup) {
		var slice []string
		plan.MayProxyAsGroup.ElementsAs(ctx, &slice, false)
		addRequest.MayProxyAsGroup = slice
	}
	if internaltypes.IsDefined(plan.MayProxyAsURL) {
		var slice []string
		plan.MayProxyAsURL.ElementsAs(ctx, &slice, false)
		addRequest.MayProxyAsURL = slice
	}
	return nil
}

// Read a TopologyAdminUserResponse object into the model struct
func readTopologyAdminUserResponse(ctx context.Context, r *client.TopologyAdminUserResponse, state *topologyAdminUserResourceModel, expectedValues *topologyAdminUserResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.AlternateBindDN = internaltypes.GetStringSet(r.AlternateBindDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.Password = expectedValues.Password
	state.FirstName = internaltypes.GetStringSet(r.FirstName)
	state.LastName = internaltypes.GetStringSet(r.LastName)
	state.UserID = internaltypes.StringTypeOrNil(r.UserID, internaltypes.IsEmptyString(expectedValues.UserID))
	state.EmailAddress = internaltypes.GetStringSet(r.EmailAddress)
	state.WorkTelephoneNumber = internaltypes.GetStringSet(r.WorkTelephoneNumber)
	state.HomeTelephoneNumber = internaltypes.GetStringSet(r.HomeTelephoneNumber)
	state.MobileTelephoneNumber = internaltypes.GetStringSet(r.MobileTelephoneNumber)
	state.PagerTelephoneNumber = internaltypes.GetStringSet(r.PagerTelephoneNumber)
	state.InheritDefaultRootPrivileges = types.BoolValue(r.InheritDefaultRootPrivileges)
	state.Privilege = internaltypes.GetStringSet(
		client.StringSliceEnumtopologyAdminUserPrivilegeProp(r.Privilege))
	state.SearchResultEntryLimit = types.Int64Value(int64(r.SearchResultEntryLimit))
	state.TimeLimitSeconds = types.Int64Value(int64(r.TimeLimitSeconds))
	state.LookThroughEntryLimit = types.Int64Value(int64(r.LookThroughEntryLimit))
	state.IdleTimeLimitSeconds = types.Int64Value(int64(r.IdleTimeLimitSeconds))
	state.PasswordPolicy = types.StringValue(r.PasswordPolicy)
	state.Disabled = internaltypes.BoolTypeOrNil(r.Disabled)
	state.AccountActivationTime = internaltypes.StringTypeOrNil(r.AccountActivationTime, internaltypes.IsEmptyString(expectedValues.AccountActivationTime))
	state.AccountExpirationTime = internaltypes.StringTypeOrNil(r.AccountExpirationTime, internaltypes.IsEmptyString(expectedValues.AccountExpirationTime))
	state.RequireSecureAuthentication = types.BoolValue(r.RequireSecureAuthentication)
	state.RequireSecureConnections = types.BoolValue(r.RequireSecureConnections)
	state.AllowedAuthenticationType = internaltypes.GetStringSet(r.AllowedAuthenticationType)
	state.AllowedAuthenticationIPAddress = internaltypes.GetStringSet(r.AllowedAuthenticationIPAddress)
	state.PreferredOTPDeliveryMechanism = internaltypes.GetStringSet(r.PreferredOTPDeliveryMechanism)
	state.IsProxyable = internaltypes.StringTypeOrNil(
		client.StringPointerEnumtopologyAdminUserIsProxyableProp(r.IsProxyable), internaltypes.IsEmptyString(expectedValues.IsProxyable))
	state.IsProxyableByDN = internaltypes.GetStringSet(r.IsProxyableByDN)
	state.IsProxyableByGroup = internaltypes.GetStringSet(r.IsProxyableByGroup)
	state.IsProxyableByURL = internaltypes.GetStringSet(r.IsProxyableByURL)
	state.MayProxyAsDN = internaltypes.GetStringSet(r.MayProxyAsDN)
	state.MayProxyAsGroup = internaltypes.GetStringSet(r.MayProxyAsGroup)
	state.MayProxyAsURL = internaltypes.GetStringSet(r.MayProxyAsURL)
	state.Notifications, state.RequiredActions = ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createTopologyAdminUserOperations(plan topologyAdminUserResourceModel, state topologyAdminUserResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AlternateBindDN, state.AlternateBindDN, "alternate-bind-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.Password, state.Password, "password")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.FirstName, state.FirstName, "first-name")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.LastName, state.LastName, "last-name")
	operations.AddStringOperationIfNecessary(&ops, plan.UserID, state.UserID, "user-id")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EmailAddress, state.EmailAddress, "email-address")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.WorkTelephoneNumber, state.WorkTelephoneNumber, "work-telephone-number")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.HomeTelephoneNumber, state.HomeTelephoneNumber, "home-telephone-number")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MobileTelephoneNumber, state.MobileTelephoneNumber, "mobile-telephone-number")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.PagerTelephoneNumber, state.PagerTelephoneNumber, "pager-telephone-number")
	operations.AddBoolOperationIfNecessary(&ops, plan.InheritDefaultRootPrivileges, state.InheritDefaultRootPrivileges, "inherit-default-root-privileges")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Privilege, state.Privilege, "privilege")
	operations.AddInt64OperationIfNecessary(&ops, plan.SearchResultEntryLimit, state.SearchResultEntryLimit, "search-result-entry-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.TimeLimitSeconds, state.TimeLimitSeconds, "time-limit-seconds")
	operations.AddInt64OperationIfNecessary(&ops, plan.LookThroughEntryLimit, state.LookThroughEntryLimit, "look-through-entry-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.IdleTimeLimitSeconds, state.IdleTimeLimitSeconds, "idle-time-limit-seconds")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordPolicy, state.PasswordPolicy, "password-policy")
	operations.AddBoolOperationIfNecessary(&ops, plan.Disabled, state.Disabled, "disabled")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountActivationTime, state.AccountActivationTime, "account-activation-time")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountExpirationTime, state.AccountExpirationTime, "account-expiration-time")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireSecureAuthentication, state.RequireSecureAuthentication, "require-secure-authentication")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireSecureConnections, state.RequireSecureConnections, "require-secure-connections")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedAuthenticationType, state.AllowedAuthenticationType, "allowed-authentication-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedAuthenticationIPAddress, state.AllowedAuthenticationIPAddress, "allowed-authentication-ip-address")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.PreferredOTPDeliveryMechanism, state.PreferredOTPDeliveryMechanism, "preferred-otp-delivery-mechanism")
	operations.AddStringOperationIfNecessary(&ops, plan.IsProxyable, state.IsProxyable, "is-proxyable")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IsProxyableByDN, state.IsProxyableByDN, "is-proxyable-by-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IsProxyableByGroup, state.IsProxyableByGroup, "is-proxyable-by-group")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IsProxyableByURL, state.IsProxyableByURL, "is-proxyable-by-url")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MayProxyAsDN, state.MayProxyAsDN, "may-proxy-as-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MayProxyAsGroup, state.MayProxyAsGroup, "may-proxy-as-group")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MayProxyAsURL, state.MayProxyAsURL, "may-proxy-as-url")
	return ops
}

// Create a new resource
func (r *topologyAdminUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan topologyAdminUserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddTopologyAdminUserRequest(plan.Id.ValueString())
	err := addOptionalTopologyAdminUserFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Topology Admin User", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.TopologyAdminUserApi.AddTopologyAdminUser(
		ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddTopologyAdminUserRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.TopologyAdminUserApi.AddTopologyAdminUserExecute(apiAddRequest)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Topology Admin User", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state topologyAdminUserResourceModel
	readTopologyAdminUserResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultTopologyAdminUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan topologyAdminUserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.TopologyAdminUserApi.GetTopologyAdminUser(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Topology Admin User", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state topologyAdminUserResourceModel
	readTopologyAdminUserResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.TopologyAdminUserApi.UpdateTopologyAdminUser(ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createTopologyAdminUserOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.TopologyAdminUserApi.UpdateTopologyAdminUserExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Topology Admin User", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readTopologyAdminUserResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *topologyAdminUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readTopologyAdminUser(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultTopologyAdminUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readTopologyAdminUser(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readTopologyAdminUser(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state topologyAdminUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.TopologyAdminUserApi.GetTopologyAdminUser(
		ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Topology Admin User", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readTopologyAdminUserResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *topologyAdminUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateTopologyAdminUser(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultTopologyAdminUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateTopologyAdminUser(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateTopologyAdminUser(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan topologyAdminUserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state topologyAdminUserResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.TopologyAdminUserApi.UpdateTopologyAdminUser(
		ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createTopologyAdminUserOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.TopologyAdminUserApi.UpdateTopologyAdminUserExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Topology Admin User", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readTopologyAdminUserResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultTopologyAdminUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *topologyAdminUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state topologyAdminUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.TopologyAdminUserApi.DeleteTopologyAdminUserExecute(r.apiClient.TopologyAdminUserApi.DeleteTopologyAdminUser(
		ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Topology Admin User", err, httpResp)
		return
	}
}

func (r *topologyAdminUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importTopologyAdminUser(ctx, req, resp)
}

func (r *defaultTopologyAdminUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importTopologyAdminUser(ctx, req, resp)
}

func importTopologyAdminUser(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
