package connectioncriteria

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &connectionCriteriaDataSource{}
	_ datasource.DataSourceWithConfigure = &connectionCriteriaDataSource{}
)

// Create a Connection Criteria data source
func NewConnectionCriteriaDataSource() datasource.DataSource {
	return &connectionCriteriaDataSource{}
}

// connectionCriteriaDataSource is the datasource implementation.
type connectionCriteriaDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *connectionCriteriaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connection_criteria"
}

// Configure adds the provider configured client to the data source.
func (r *connectionCriteriaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type connectionCriteriaDataSourceModel struct {
	Id                               types.String `tfsdk:"id"`
	Name                             types.String `tfsdk:"name"`
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

// GetSchema defines the schema for the datasource.
func (r *connectionCriteriaDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Connection Criteria.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Connection Criteria resource. Options are ['simple', 'aggregate', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Connection Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Connection Criteria. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"all_included_connection_criteria": schema.SetAttribute{
				Description: "Specifies a connection criteria object that must match the associated client connection in order to match the aggregate connection criteria. If one or more all-included connection criteria objects are provided, then a client connection must match all of them in order to match the aggregate connection criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_connection_criteria": schema.SetAttribute{
				Description: "Specifies a connection criteria object that may match the associated client connection in order to match the aggregate connection criteria. If one or more any-included connection criteria objects are provided, then a client connection must match at least one of them in order to match the aggregate connection criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_connection_criteria": schema.SetAttribute{
				Description: "Specifies a connection criteria object that should not match the associated client connection in order to match the aggregate connection criteria. If one or more not-all-included connection criteria objects are provided, then a client connection must not match all of them (that is, it may match zero or more of them, but it must not match all of them) in order to match the aggregate connection criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_connection_criteria": schema.SetAttribute{
				Description: "Specifies a connection criteria object that must not match the associated client connection in order to match the aggregate connection criteria. If one or more none-included connection criteria objects are provided, then a client connection must not match any of them in order to match the aggregate connection criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_client_address": schema.SetAttribute{
				Description: "Specifies an address mask that may be used to specify a set of clients that should be included in this Simple Connection Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_client_address": schema.SetAttribute{
				Description: "Specifies an address mask that may be used to specify a set of clients that should be excluded from this Simple Connection Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_connection_handler": schema.SetAttribute{
				Description: "Specifies a connection handler for clients that should be included in this Simple Connection Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_connection_handler": schema.SetAttribute{
				Description: "Specifies a connection handler for clients that should be excluded from this Simple Connection Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_protocol": schema.SetAttribute{
				Description: "Specifies the name of a communication protocol that should be used by clients included in this Simple Connection Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_protocol": schema.SetAttribute{
				Description: "Specifies the name of a communication protocol that should be used by clients excluded from this Simple Connection Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"communication_security_level": schema.StringAttribute{
				Description: "Indicates whether this Simple Connection Criteria should require or allow clients using a secure communication channel.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"user_auth_type": schema.SetAttribute{
				Description: "Specifies the authentication types for client connections that may be included in this Simple Connection Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"authentication_security_level": schema.StringAttribute{
				Description: "Indicates whether this Simple Connection Criteria should require or allow clients that authenticated using a secure manner. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"included_user_sasl_mechanism": schema.SetAttribute{
				Description: "Specifies the name of a SASL mechanism that should be used by clients included in this Simple Connection Criteria. This will only be taken into account for client connections that have authenticated to the server using a SASL mechanism and will be ignored for unauthenticated client connections and for client connections that authenticated using some other method (e.g., those performing simple or internal authentication).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_user_sasl_mechanism": schema.SetAttribute{
				Description: "Specifies the name of a SASL mechanism that should be used by clients excluded from this Simple Connection Criteria. This will only be taken into account for client connections that have authenticated to the server using a SASL mechanism and will be ignored for unauthenticated client connections and for client connections that authenticated using some other method (e.g., those performing simple or internal authentication).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_user_base_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which authenticated user entries may exist for clients included in this Simple Connection Criteria. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections. Refer to the authz version of this property in Simple Result Criteria if operations are being proxied (performed using proxied authorization), and you need to match the originating user of the operation rather than the proxy user (the user the proxy authenticated as).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_user_base_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which authenticated user entries may exist for clients excluded from this Simple Connection Criteria. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections. Refer to the authz version of this property in Simple Result Criteria if operations are being proxied (performed using proxied authorization), and you need to match the originating user of the operation rather than the proxy user (the user the proxy authenticated as).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"all_included_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authenticated users must exist for clients included in this Simple Connection Criteria. If any group DNs are provided, then the authenticated user must be a member of all of those groups. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections. Refer to the authz version of this property in Simple Result Criteria if operations are being proxied (performed using proxied authorization), and you need to match the originating user of the operation rather than the proxy user (the user the proxy authenticated as).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authenticated users may exist for clients included in this Simple Connection Criteria. If any group DNs are provided, then the authenticated user must be a member of at least one of those groups. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections. Refer to the authz version of this property in Simple Result Criteria if operations are being proxied (performed using proxied authorization), and you need to match the originating user of the operation rather than the proxy user (the user the proxy authenticated as).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authenticated users should not exist for clients included in this Simple Connection Criteria. If any group DNs are provided, then the authenticated user must not be a member of at least one of those groups (that is, the user may be a member of zero or more of those groups, but not of all of them). This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections. Refer to the authz version of this property in Simple Result Criteria if operations are being proxied (performed using proxied authorization), and you need to match the originating user of the operation rather than the proxy user (the user the proxy authenticated as).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authenticated users must not exist for clients included in this Simple Connection Criteria. If any group DNs are provided, then the authenticated user must not be a member any of those groups. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections. Refer to the authz version of this property in Simple Result Criteria if operations are being proxied (performed using proxied authorization), and you need to match the originating user of the operation rather than the proxy user (the user the proxy authenticated as).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"all_included_user_filter": schema.SetAttribute{
				Description: "Specifies a search filter that must match the entry of the authenticated user for clients included in this Simple Connection Criteria. If any filters are provided, then all of those filters must match the authenticated user entry. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_user_filter": schema.SetAttribute{
				Description: "Specifies a search filter that may match the entry of the authenticated user for clients included in this Simple Connection Criteria. If any filters are provided, then at least one of those filters must match the authenticated user entry. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_user_filter": schema.SetAttribute{
				Description: "Specifies a search filter that should not match the entry of the authenticated user for clients included in this Simple Connection Criteria. If any filters are provided, then at least one of those filters must not match the authenticated user entry (that is, the user entry may match zero or more of those filters, but not all of them). This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_user_filter": schema.SetAttribute{
				Description: "Specifies a search filter that must not match the entry of the authenticated user for clients included in this Simple Connection Criteria. If any filters are provided, then none of those filters may match the authenticated user entry. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"all_included_user_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that must be held by the authenticated user for clients included in this Simple Connection Criteria. If any privilege names are provided, then the authenticated user must have all of those privileges. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_user_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that may be held by the authenticated user for clients included in this Simple Connection Criteria. If any privilege names are provided, then the authenticated user must have at least one of those privileges. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_user_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that should not be held by the authenticated user for clients included in this Simple Connection Criteria. If any privilege names are provided, then the authenticated user must not have at least one of those privileges (that is, the user may hold zero or more of those privileges, but not all of them). This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_user_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that must not be held by the authenticated user for clients included in this Simple Connection Criteria. If any privilege names are provided, then the authenticated user must not have any of those privileges. This will only be taken into account for client connections that have authenticated to the server and will be ignored for unauthenticated client connections.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Connection Criteria",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a SimpleConnectionCriteriaResponse object into the model struct
func readSimpleConnectionCriteriaResponseDataSource(ctx context.Context, r *client.SimpleConnectionCriteriaResponse, state *connectionCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("simple")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IncludedClientAddress = internaltypes.GetStringSet(r.IncludedClientAddress)
	state.ExcludedClientAddress = internaltypes.GetStringSet(r.ExcludedClientAddress)
	state.IncludedConnectionHandler = internaltypes.GetStringSet(r.IncludedConnectionHandler)
	state.ExcludedConnectionHandler = internaltypes.GetStringSet(r.ExcludedConnectionHandler)
	state.IncludedProtocol = internaltypes.GetStringSet(r.IncludedProtocol)
	state.ExcludedProtocol = internaltypes.GetStringSet(r.ExcludedProtocol)
	state.CommunicationSecurityLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumconnectionCriteriaCommunicationSecurityLevelProp(r.CommunicationSecurityLevel), false)
	state.UserAuthType = internaltypes.GetStringSet(
		client.StringSliceEnumconnectionCriteriaUserAuthTypeProp(r.UserAuthType))
	state.AuthenticationSecurityLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumconnectionCriteriaAuthenticationSecurityLevelProp(r.AuthenticationSecurityLevel), false)
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
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a AggregateConnectionCriteriaResponse object into the model struct
func readAggregateConnectionCriteriaResponseDataSource(ctx context.Context, r *client.AggregateConnectionCriteriaResponse, state *connectionCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("aggregate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllIncludedConnectionCriteria = internaltypes.GetStringSet(r.AllIncludedConnectionCriteria)
	state.AnyIncludedConnectionCriteria = internaltypes.GetStringSet(r.AnyIncludedConnectionCriteria)
	state.NotAllIncludedConnectionCriteria = internaltypes.GetStringSet(r.NotAllIncludedConnectionCriteria)
	state.NoneIncludedConnectionCriteria = internaltypes.GetStringSet(r.NoneIncludedConnectionCriteria)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a ThirdPartyConnectionCriteriaResponse object into the model struct
func readThirdPartyConnectionCriteriaResponseDataSource(ctx context.Context, r *client.ThirdPartyConnectionCriteriaResponse, state *connectionCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read resource information
func (r *connectionCriteriaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state connectionCriteriaDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConnectionCriteriaAPI.GetConnectionCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
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
		readSimpleConnectionCriteriaResponseDataSource(ctx, readResponse.SimpleConnectionCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AggregateConnectionCriteriaResponse != nil {
		readAggregateConnectionCriteriaResponseDataSource(ctx, readResponse.AggregateConnectionCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyConnectionCriteriaResponse != nil {
		readThirdPartyConnectionCriteriaResponseDataSource(ctx, readResponse.ThirdPartyConnectionCriteriaResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
