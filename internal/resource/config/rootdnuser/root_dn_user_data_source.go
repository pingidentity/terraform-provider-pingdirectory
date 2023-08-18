package rootdnuser

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &rootDnUserDataSource{}
	_ datasource.DataSourceWithConfigure = &rootDnUserDataSource{}
)

// Create a Root Dn User data source
func NewRootDnUserDataSource() datasource.DataSource {
	return &rootDnUserDataSource{}
}

// rootDnUserDataSource is the datasource implementation.
type rootDnUserDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *rootDnUserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_root_dn_user"
}

// Configure adds the provider configured client to the data source.
func (r *rootDnUserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type rootDnUserDataSourceModel struct {
	Id                             types.String `tfsdk:"id"`
	Name                           types.String `tfsdk:"name"`
	Type                           types.String `tfsdk:"type"`
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

// GetSchema defines the schema for the datasource.
func (r *rootDnUserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Root Dn User.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Root DN User resource. Options are ['root-dn-user']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"alternate_bind_dn": schema.SetAttribute{
				Description: "Specifies one or more alternate DNs that can be used to bind to the server as this User.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this User.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "Specifies the user's password. This is stored in the userPassword LDAP attribute. To set a pre-hashed value, the account making the change must have the bypass-pw-policy privilege.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"first_name": schema.SetAttribute{
				Description: "Specifies the user's first name. This is stored in the givenName LDAP attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"last_name": schema.SetAttribute{
				Description: "Specifies the user's last name. This is stored in the sn LDAP attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"user_id": schema.StringAttribute{
				Description: "Specifies the user's user ID. This is stored in the uid LDAP attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"email_address": schema.SetAttribute{
				Description: "Specifies the user's email address. This is stored in the mail LDAP attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"work_telephone_number": schema.SetAttribute{
				Description: "Specifies the user's work telephone number. This is stored in the telephoneNumber LDAP attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"home_telephone_number": schema.SetAttribute{
				Description: "Specifies the user's home telephone number. This is stored in the homePhone LDAP attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"mobile_telephone_number": schema.SetAttribute{
				Description: "Specifies the user's mobile telephone number. This is stored in the mobile LDAP attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"pager_telephone_number": schema.SetAttribute{
				Description: "Specifies the user's pager telephone number. This is stored in the pager LDAP attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"inherit_default_root_privileges": schema.BoolAttribute{
				Description: "Indicates whether this User should be automatically granted the set of privileges defined in the default-root-privilege-name property of the Root DN configuration object.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"privilege": schema.SetAttribute{
				Description: "Privileges that are either explicitly granted or revoked from the root user. Privileges can be revoked by including a minus sign (-) before the privilege name. This is stored in the ds-privilege-name LDAP attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"search_result_entry_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that the server may return to the user in response to any single search request. A value of 0 indicates no limit should be enforced. This is stored in the ds-rlim-size-limit LDAP attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"time_limit_seconds": schema.Int64Attribute{
				Description: "Specifies the maximum length of time (in seconds) that the server may spend processing any single search request. A value of 0 indicates no limit should be enforced. This is stored in the ds-rlim-time-limit LDAP attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"look_through_entry_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of candidate entries that the server may examine in the course of processing any single search request. A value of 0 indicates no limit should be enforced. This is stored in the ds-rlim-lookthrough-limit LDAP attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"idle_time_limit_seconds": schema.Int64Attribute{
				Description: "Specifies the maximum length of time (in seconds) that a connection authenticated as this user may remain established without issuing any requests. A value of 0 indicates no limit should be enforced. This is stored in the ds-rlim-idle-time-limit LDAP attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password_policy": schema.StringAttribute{
				Description: "Specifies the password policy for the user. This is stored in the ds-pwp-password-policy-dn LDAP attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: "Specifies whether the root user account should be disabled. A disabled account is not permitted to authenticate, nor can it be used as an authorization identity. This is stored in the ds-pwp-account-disabled LDAP attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_activation_time": schema.StringAttribute{
				Description: "Specifies the time, in generalized time format (e.g., '20160101070000Z'), that the root user account should become active. If an activation time is specified, the user will not be permitted to authenticate, nor can the account be used as an authorization identity, until the activation time has arrived. This is stored in the ds-pwp-account-activation-time LDAP attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_expiration_time": schema.StringAttribute{
				Description: "Specifies the time, in generalized time format (e.g., '20240101070000Z'), that the root user account should expire. If an expiration time is specified, the user will not be permitted to authenticate, nor can the account be used as an authorization identity, after this time has passed. This is stored in the ds-pwp-account-expiration-time LDAP attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"require_secure_authentication": schema.BoolAttribute{
				Description: "Indicates whether this User must authenticate in a secure manner. When set to \"true\", the User will only be allowed to authenticate over a secure connection or using a mechanism that does not expose user credentials (e.g., the CRAM-MD5, DIGEST-MD5, and GSSAPI SASL mechanisms).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"require_secure_connections": schema.BoolAttribute{
				Description: "Indicates whether this User must be required to communicate with the server over a secure connection. When set to \"true\", the User will only be allowed to communicate with the server over a secure connection (i.e., using TLS or the StartTLS extended operation).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allowed_authentication_type": schema.SetAttribute{
				Description: "Indicates that User should only be allowed to authenticate in certain ways. Allowed values include \"simple\" (to indicate that the user should be allowed to bind using simple authentication) or \"sasl {mech}\" (to indicate that the user should be allowed to bind using the specified SASL mechanism, like \"sasl PLAIN\"). The list of available SASL mechanisms can be retrieved by running \"dsconfig --advanced list-sasl-mechanism-handlers\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"allowed_authentication_ip_address": schema.SetAttribute{
				Description: "An IPv4 or IPv6 address mask that controls the set of IP addresses from which this User can authenticate to the server. For instance a value of 127.0.0.1 (or ::1 in IPv6) would restricted access only to localhost connections, whereas 10.6.1.* would restrict access to servers on the 10.6.1.* subnet.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"preferred_otp_delivery_mechanism": schema.SetAttribute{
				Description: "Overrides the default settings for the mechanisms (e.g., email or SMS) that are used to deliver one time passwords to Users.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"is_proxyable": schema.StringAttribute{
				Description: "This can be used to indicate whether the User can be used as an alternate authorization identity (using the proxied authorization v1 or v2 control, the intermediate client control, or a SASL mechanism that allows specifying an alternate authorization identity).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"is_proxyable_by_dn": schema.SetAttribute{
				Description: "Specifies the DNs of accounts that can proxy as this User using the proxied authorization v1 or v2 control, the intermediate client control, or a SASL mechanism that allows specifying an alternate authorization identity. This property is only applicable if is-proxyable is set to \"allowed\" or \"required\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"is_proxyable_by_group": schema.SetAttribute{
				Description: "Specifies the DNs of groups whose members can proxy as this User using the proxied authorization v1 or v2 control, the intermediate client control, or a SASL mechanism that allows specifying an alternate authorization identity. This property is only applicable if is-proxyable is set to \"allowed\" or \"required\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"is_proxyable_by_url": schema.SetAttribute{
				Description: "Specifies LDAP URLs of accounts that can proxy as this User using the proxied authorization v1 or v2 control, the intermediate client control, or a SASL mechanism that allows specifying an alternate authorization identity. This property is only applicable if is-proxyable is set to \"allowed\" or \"required\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"may_proxy_as_dn": schema.SetAttribute{
				Description: "This restricts the set of accounts that this User can proxy as to entries with the specified DNs.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"may_proxy_as_group": schema.SetAttribute{
				Description: "This restricts the set of accounts that this User can proxy as to entries that are in the group with the specified DN.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"may_proxy_as_url": schema.SetAttribute{
				Description: "This restricts the set of accounts that this User can proxy as to entries that are matched by the specified LDAP URL.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a RootDnUserResponse object into the model struct
func readRootDnUserResponseDataSource(ctx context.Context, r *client.RootDnUserResponse, state *rootDnUserDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("root-dn-user")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AlternateBindDN = internaltypes.GetStringSet(r.AlternateBindDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.FirstName = internaltypes.GetStringSet(r.FirstName)
	state.LastName = internaltypes.GetStringSet(r.LastName)
	state.UserID = internaltypes.StringTypeOrNil(r.UserID, false)
	state.EmailAddress = internaltypes.GetStringSet(r.EmailAddress)
	state.WorkTelephoneNumber = internaltypes.GetStringSet(r.WorkTelephoneNumber)
	state.HomeTelephoneNumber = internaltypes.GetStringSet(r.HomeTelephoneNumber)
	state.MobileTelephoneNumber = internaltypes.GetStringSet(r.MobileTelephoneNumber)
	state.PagerTelephoneNumber = internaltypes.GetStringSet(r.PagerTelephoneNumber)
	state.InheritDefaultRootPrivileges = types.BoolValue(r.InheritDefaultRootPrivileges)
	state.Privilege = internaltypes.GetStringSet(
		client.StringSliceEnumrootDnUserPrivilegeProp(r.Privilege))
	state.SearchResultEntryLimit = types.Int64Value(r.SearchResultEntryLimit)
	state.TimeLimitSeconds = types.Int64Value(r.TimeLimitSeconds)
	state.LookThroughEntryLimit = types.Int64Value(r.LookThroughEntryLimit)
	state.IdleTimeLimitSeconds = types.Int64Value(r.IdleTimeLimitSeconds)
	state.PasswordPolicy = types.StringValue(r.PasswordPolicy)
	state.Disabled = internaltypes.BoolTypeOrNil(r.Disabled)
	state.AccountActivationTime = internaltypes.StringTypeOrNil(r.AccountActivationTime, false)
	state.AccountExpirationTime = internaltypes.StringTypeOrNil(r.AccountExpirationTime, false)
	state.RequireSecureAuthentication = types.BoolValue(r.RequireSecureAuthentication)
	state.RequireSecureConnections = types.BoolValue(r.RequireSecureConnections)
	state.AllowedAuthenticationType = internaltypes.GetStringSet(r.AllowedAuthenticationType)
	state.AllowedAuthenticationIPAddress = internaltypes.GetStringSet(r.AllowedAuthenticationIPAddress)
	state.PreferredOTPDeliveryMechanism = internaltypes.GetStringSet(r.PreferredOTPDeliveryMechanism)
	state.IsProxyable = internaltypes.StringTypeOrNil(
		client.StringPointerEnumrootDnUserIsProxyableProp(r.IsProxyable), false)
	state.IsProxyableByDN = internaltypes.GetStringSet(r.IsProxyableByDN)
	state.IsProxyableByGroup = internaltypes.GetStringSet(r.IsProxyableByGroup)
	state.IsProxyableByURL = internaltypes.GetStringSet(r.IsProxyableByURL)
	state.MayProxyAsDN = internaltypes.GetStringSet(r.MayProxyAsDN)
	state.MayProxyAsGroup = internaltypes.GetStringSet(r.MayProxyAsGroup)
	state.MayProxyAsURL = internaltypes.GetStringSet(r.MayProxyAsURL)
}

// Read resource information
func (r *rootDnUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state rootDnUserDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RootDnUserApi.GetRootDnUser(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Root Dn User", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readRootDnUserResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
