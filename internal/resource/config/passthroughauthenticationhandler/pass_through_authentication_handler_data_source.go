// Copyright © 2025 Ping Identity Corporation

package passthroughauthenticationhandler

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &passThroughAuthenticationHandlerDataSource{}
	_ datasource.DataSourceWithConfigure = &passThroughAuthenticationHandlerDataSource{}
)

// Create a Pass Through Authentication Handler data source
func NewPassThroughAuthenticationHandlerDataSource() datasource.DataSource {
	return &passThroughAuthenticationHandlerDataSource{}
}

// passThroughAuthenticationHandlerDataSource is the datasource implementation.
type passThroughAuthenticationHandlerDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *passThroughAuthenticationHandlerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pass_through_authentication_handler"
}

// Configure adds the provider configured client to the data source.
func (r *passThroughAuthenticationHandlerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type passThroughAuthenticationHandlerDataSourceModel struct {
	Id                                          types.String `tfsdk:"id"`
	Name                                        types.String `tfsdk:"name"`
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

// GetSchema defines the schema for the datasource.
func (r *passThroughAuthenticationHandlerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Pass Through Authentication Handler.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Pass Through Authentication Handler resource. Options are ['ping-one', 'ldap', 'aggregate', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Pass Through Authentication Handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Pass Through Authentication Handler. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"subordinate_pass_through_authentication_handler": schema.SetAttribute{
				Description: "The set of subordinate pass-through authentication handlers that may be used to perform the authentication processing. Handlers will be invoked in order until one is found for which the bind operation matches the associated criteria and either succeeds or fails in a manner that should not be ignored.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"continue_on_failure_type": schema.SetAttribute{
				Description: "The set of pass-through authentication failure types that should not result in an immediate failure, but should instead allow the aggregate handler to proceed with the next configured subordinate handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"server": schema.SetAttribute{
				Description: "Specifies the LDAP external server(s) to which authentication attempts should be forwarded.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"server_access_mode": schema.StringAttribute{
				Description: "Specifies the manner in which external servers should be used for pass-through authentication attempts if multiple servers are defined.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"dn_map": schema.SetAttribute{
				Description: "Specifies one or more DN mappings that may be used to transform bind DNs before attempting to bind to the external servers.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"bind_dn_pattern": schema.StringAttribute{
				Description: "A pattern to use to construct the bind DN for the simple bind request to send to the remote server. This may consist of a combination of static text and attribute values and other directives enclosed in curly braces.  For example, the value \"cn={cn},ou=People,dc=example,dc=com\" indicates that the remote bind DN should be constructed from the text \"cn=\" followed by the value of the local entry's cn attribute followed by the text \"ou=People,dc=example,dc=com\". If an attribute contains the value to use as the bind DN for pass-through authentication, then the pattern may simply be the name of that attribute in curly braces (e.g., if the seeAlso attribute contains the bind DN for the target user, then a bind DN pattern of \"{seeAlso}\" would be appropriate).  Note that a bind DN pattern can be used to construct a bind DN that is not actually a valid LDAP distinguished name. For example, if authentication is being passed through to a Microsoft Active Directory server, then a bind DN pattern could be used to construct a user principal name (UPN) as an alternative to a distinguished name.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"search_base_dn": schema.StringAttribute{
				Description: "The base DN to use when searching for the user entry using a filter constructed from the pattern defined in the search-filter-pattern property. If no base DN is specified, the null DN will be used as the search base DN.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"search_filter_pattern": schema.StringAttribute{
				Description: "A pattern to use to construct a filter to use when searching an external server for the entry of the user as whom to bind. For example, \"(mail={uid:ldapFilterEscape}@example.com)\" would construct a search filter to search for a user whose entry in the local server contains a uid attribute whose value appears before \"@example.com\" in the mail attribute in the external server. Note that the \"ldapFilterEscape\" modifier should almost always be used with attributes specified in the pattern.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"initial_connections": schema.Int64Attribute{
				Description: "Specifies the initial number of connections to establish to each external server against which authentication may be attempted.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_connections": schema.Int64Attribute{
				Description: "Specifies the maximum number of connections to maintain to each external server against which authentication may be attempted. This value must be greater than or equal to the value for the initial-connections property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"use_location": schema.BoolAttribute{
				Description: "Indicates whether to take server locations into account when prioritizing the servers to use for pass-through authentication attempts.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_allowed_local_response_time": schema.StringAttribute{
				Description: "The maximum length of time to wait for a response from an external server in the same location as this Directory Server before considering it unavailable.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_allowed_nonlocal_response_time": schema.StringAttribute{
				Description: "The maximum length of time to wait for a response from an external server in a different location from this Directory Server before considering it unavailable.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"use_password_policy_control": schema.BoolAttribute{
				Description: "Indicates whether to include the password policy request control (as defined in draft-behera-ldap-password-policy-10) in bind requests sent to the external server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"api_url": schema.StringAttribute{
				Description: "Specifies the API endpoint for the PingOne web service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"auth_url": schema.StringAttribute{
				Description: "Specifies the API endpoint for the PingOne authentication service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"oauth_client_id": schema.StringAttribute{
				Description: "Specifies the OAuth Client ID used to authenticate connections to the PingOne API.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"oauth_client_secret": schema.StringAttribute{
				Description: "Specifies the OAuth Client Secret used to authenticate connections to the PingOne API.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"oauth_client_secret_passphrase_provider": schema.StringAttribute{
				Description: "Specifies a passphrase provider that can be used to obtain the OAuth Client Secret used to authenticate connections to the PingOne API.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"environment_id": schema.StringAttribute{
				Description: "Specifies the PingOne Environment that will be associated with this PingOne Pass Through Authentication Handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"http_proxy_external_server": schema.StringAttribute{
				Description: "A reference to an HTTP proxy server that should be used for requests sent to the PingOne service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"user_mapping_local_attribute": schema.SetAttribute{
				Description: "The names of the attributes in the local user entry whose values must match the values of the corresponding fields in the PingOne service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"user_mapping_remote_json_field": schema.SetAttribute{
				Description: "The names of the fields in the PingOne service whose values must match the values of the corresponding attributes in the local user entry, as specified in the user-mapping-local-attribute property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"additional_user_mapping_scim_filter": schema.StringAttribute{
				Description: "An optional SCIM filter that will be ANDed with the filter created to identify the account in the PingOne service that corresponds to the local entry. Only the \"eq\", \"sw\", \"and\", and \"or\" filter types may be used.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Pass Through Authentication Handler",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"included_local_entry_base_dn": schema.SetAttribute{
				Description: "The base DNs for the local users whose authentication attempts may be passed through to the external authentication service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"connection_criteria": schema.StringAttribute{
				Description: "A reference to connection criteria that will be used to indicate which bind requests should be passed through to the external authentication service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"request_criteria": schema.StringAttribute{
				Description: "A reference to request criteria that will be used to indicate which bind requests should be passed through to the external authentication service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a PingOnePassThroughAuthenticationHandlerResponse object into the model struct
func readPingOnePassThroughAuthenticationHandlerResponseDataSource(ctx context.Context, r *client.PingOnePassThroughAuthenticationHandlerResponse, state *passThroughAuthenticationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ping-one")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ApiURL = types.StringValue(r.ApiURL)
	state.AuthURL = types.StringValue(r.AuthURL)
	state.OAuthClientID = types.StringValue(r.OAuthClientID)
	state.OAuthClientSecretPassphraseProvider = internaltypes.StringTypeOrNil(r.OAuthClientSecretPassphraseProvider, false)
	state.EnvironmentID = types.StringValue(r.EnvironmentID)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, false)
	state.UserMappingLocalAttribute = internaltypes.GetStringSet(r.UserMappingLocalAttribute)
	state.UserMappingRemoteJSONField = internaltypes.GetStringSet(r.UserMappingRemoteJSONField)
	state.AdditionalUserMappingSCIMFilter = internaltypes.StringTypeOrNil(r.AdditionalUserMappingSCIMFilter, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.IncludedLocalEntryBaseDN = internaltypes.GetStringSet(r.IncludedLocalEntryBaseDN)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
}

// Read a LdapPassThroughAuthenticationHandlerResponse object into the model struct
func readLdapPassThroughAuthenticationHandlerResponseDataSource(ctx context.Context, r *client.LdapPassThroughAuthenticationHandlerResponse, state *passThroughAuthenticationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Server = internaltypes.GetStringSet(r.Server)
	state.ServerAccessMode = types.StringValue(r.ServerAccessMode.String())
	state.DnMap = internaltypes.GetStringSet(r.DnMap)
	state.BindDNPattern = internaltypes.StringTypeOrNil(r.BindDNPattern, false)
	state.SearchBaseDN = internaltypes.StringTypeOrNil(r.SearchBaseDN, false)
	state.SearchFilterPattern = internaltypes.StringTypeOrNil(r.SearchFilterPattern, false)
	state.InitialConnections = types.Int64Value(r.InitialConnections)
	state.MaxConnections = types.Int64Value(r.MaxConnections)
	state.UseLocation = internaltypes.BoolTypeOrNil(r.UseLocation)
	state.MaximumAllowedLocalResponseTime = internaltypes.StringTypeOrNil(r.MaximumAllowedLocalResponseTime, false)
	state.MaximumAllowedNonlocalResponseTime = internaltypes.StringTypeOrNil(r.MaximumAllowedNonlocalResponseTime, false)
	state.UsePasswordPolicyControl = internaltypes.BoolTypeOrNil(r.UsePasswordPolicyControl)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.IncludedLocalEntryBaseDN = internaltypes.GetStringSet(r.IncludedLocalEntryBaseDN)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
}

// Read a AggregatePassThroughAuthenticationHandlerResponse object into the model struct
func readAggregatePassThroughAuthenticationHandlerResponseDataSource(ctx context.Context, r *client.AggregatePassThroughAuthenticationHandlerResponse, state *passThroughAuthenticationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("aggregate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SubordinatePassThroughAuthenticationHandler = internaltypes.GetStringSet(r.SubordinatePassThroughAuthenticationHandler)
	state.ContinueOnFailureType = internaltypes.GetStringSet(
		client.StringSliceEnumpassThroughAuthenticationHandlerContinueOnFailureTypeProp(r.ContinueOnFailureType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.IncludedLocalEntryBaseDN = internaltypes.GetStringSet(r.IncludedLocalEntryBaseDN)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
}

// Read a ThirdPartyPassThroughAuthenticationHandlerResponse object into the model struct
func readThirdPartyPassThroughAuthenticationHandlerResponseDataSource(ctx context.Context, r *client.ThirdPartyPassThroughAuthenticationHandlerResponse, state *passThroughAuthenticationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.IncludedLocalEntryBaseDN = internaltypes.GetStringSet(r.IncludedLocalEntryBaseDN)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
}

// Read resource information
func (r *passThroughAuthenticationHandlerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state passThroughAuthenticationHandlerDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PassThroughAuthenticationHandlerAPI.GetPassThroughAuthenticationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Pass Through Authentication Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.PingOnePassThroughAuthenticationHandlerResponse != nil {
		readPingOnePassThroughAuthenticationHandlerResponseDataSource(ctx, readResponse.PingOnePassThroughAuthenticationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LdapPassThroughAuthenticationHandlerResponse != nil {
		readLdapPassThroughAuthenticationHandlerResponseDataSource(ctx, readResponse.LdapPassThroughAuthenticationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AggregatePassThroughAuthenticationHandlerResponse != nil {
		readAggregatePassThroughAuthenticationHandlerResponseDataSource(ctx, readResponse.AggregatePassThroughAuthenticationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPassThroughAuthenticationHandlerResponse != nil {
		readThirdPartyPassThroughAuthenticationHandlerResponseDataSource(ctx, readResponse.ThirdPartyPassThroughAuthenticationHandlerResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
