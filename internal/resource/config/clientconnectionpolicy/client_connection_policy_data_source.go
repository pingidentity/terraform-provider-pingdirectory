package clientconnectionpolicy

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
	_ datasource.DataSource              = &clientConnectionPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &clientConnectionPolicyDataSource{}
)

// Create a Client Connection Policy data source
func NewClientConnectionPolicyDataSource() datasource.DataSource {
	return &clientConnectionPolicyDataSource{}
}

// clientConnectionPolicyDataSource is the datasource implementation.
type clientConnectionPolicyDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *clientConnectionPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client_connection_policy"
}

// Configure adds the provider configured client to the data source.
func (r *clientConnectionPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type clientConnectionPolicyDataSourceModel struct {
	Id                                                       types.String `tfsdk:"id"`
	Type                                                     types.String `tfsdk:"type"`
	PolicyID                                                 types.String `tfsdk:"policy_id"`
	Description                                              types.String `tfsdk:"description"`
	Enabled                                                  types.Bool   `tfsdk:"enabled"`
	EvaluationOrderIndex                                     types.Int64  `tfsdk:"evaluation_order_index"`
	ConnectionCriteria                                       types.String `tfsdk:"connection_criteria"`
	TerminateConnection                                      types.Bool   `tfsdk:"terminate_connection"`
	SensitiveAttribute                                       types.Set    `tfsdk:"sensitive_attribute"`
	ExcludeGlobalSensitiveAttribute                          types.Set    `tfsdk:"exclude_global_sensitive_attribute"`
	ResultCodeMap                                            types.String `tfsdk:"result_code_map"`
	IncludedBackendBaseDN                                    types.Set    `tfsdk:"included_backend_base_dn"`
	ExcludedBackendBaseDN                                    types.Set    `tfsdk:"excluded_backend_base_dn"`
	AllowedOperation                                         types.Set    `tfsdk:"allowed_operation"`
	RequiredOperationRequestCriteria                         types.String `tfsdk:"required_operation_request_criteria"`
	ProhibitedOperationRequestCriteria                       types.String `tfsdk:"prohibited_operation_request_criteria"`
	AllowedRequestControl                                    types.Set    `tfsdk:"allowed_request_control"`
	DeniedRequestControl                                     types.Set    `tfsdk:"denied_request_control"`
	AllowedExtendedOperation                                 types.Set    `tfsdk:"allowed_extended_operation"`
	DeniedExtendedOperation                                  types.Set    `tfsdk:"denied_extended_operation"`
	AllowedAuthType                                          types.Set    `tfsdk:"allowed_auth_type"`
	AllowedSASLMechanism                                     types.Set    `tfsdk:"allowed_sasl_mechanism"`
	DeniedSASLMechanism                                      types.Set    `tfsdk:"denied_sasl_mechanism"`
	AllowedFilterType                                        types.Set    `tfsdk:"allowed_filter_type"`
	AllowUnindexedSearches                                   types.Bool   `tfsdk:"allow_unindexed_searches"`
	AllowUnindexedSearchesWithControl                        types.Bool   `tfsdk:"allow_unindexed_searches_with_control"`
	MinimumSubstringLength                                   types.Int64  `tfsdk:"minimum_substring_length"`
	MaximumConcurrentConnections                             types.Int64  `tfsdk:"maximum_concurrent_connections"`
	MaximumConnectionDuration                                types.String `tfsdk:"maximum_connection_duration"`
	MaximumIdleConnectionDuration                            types.String `tfsdk:"maximum_idle_connection_duration"`
	MaximumOperationCountPerConnection                       types.Int64  `tfsdk:"maximum_operation_count_per_connection"`
	MaximumConcurrentOperationsPerConnection                 types.Int64  `tfsdk:"maximum_concurrent_operations_per_connection"`
	MaximumConcurrentOperationWaitTimeBeforeRejecting        types.String `tfsdk:"maximum_concurrent_operation_wait_time_before_rejecting"`
	MaximumConcurrentOperationsPerConnectionExceededBehavior types.String `tfsdk:"maximum_concurrent_operations_per_connection_exceeded_behavior"`
	MaximumConnectionOperationRate                           types.Set    `tfsdk:"maximum_connection_operation_rate"`
	ConnectionOperationRateExceededBehavior                  types.String `tfsdk:"connection_operation_rate_exceeded_behavior"`
	MaximumPolicyOperationRate                               types.Set    `tfsdk:"maximum_policy_operation_rate"`
	PolicyOperationRateExceededBehavior                      types.String `tfsdk:"policy_operation_rate_exceeded_behavior"`
	MaximumSearchSizeLimit                                   types.Int64  `tfsdk:"maximum_search_size_limit"`
	MaximumSearchTimeLimit                                   types.String `tfsdk:"maximum_search_time_limit"`
	MaximumSearchLookthroughLimit                            types.Int64  `tfsdk:"maximum_search_lookthrough_limit"`
	MaximumLDAPJoinSizeLimit                                 types.Int64  `tfsdk:"maximum_ldap_join_size_limit"`
	MaximumSortSizeLimitWithoutVLVIndex                      types.Int64  `tfsdk:"maximum_sort_size_limit_without_vlv_index"`
}

// GetSchema defines the schema for the datasource.
func (r *clientConnectionPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Client Connection Policy.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Client Connection Policy resource. Options are ['client-connection-policy']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"policy_id": schema.StringAttribute{
				Description: "Specifies a name which uniquely identifies this Client Connection Policy in the server.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Client Connection Policy",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Client Connection Policy is enabled for use in the server. If a Client Connection Policy is disabled, then no new client connections will be associated with it.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"evaluation_order_index": schema.Int64Attribute{
				Description: "Specifies the order in which Client Connection Policy definitions will be evaluated. A Client Connection Policy with a lower index will be evaluated before one with a higher index, and the first Client Connection Policy evaluated which may apply to a client connection will be used for that connection. Each Client Connection Policy must be assigned a unique evaluation order index value.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"connection_criteria": schema.StringAttribute{
				Description: "Specifies a set of connection criteria that must match the associated client connection for it to be associated with this Client Connection Policy.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"terminate_connection": schema.BoolAttribute{
				Description: "Indicates whether any client connection for which this Client Connection Policy is selected should be terminated. This makes it possible to define fine-grained criteria for clients that should not be allowed to connect to this Directory Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"sensitive_attribute": schema.SetAttribute{
				Description: "Provides the ability to indicate that some attributes should be considered sensitive and additional protection should be in place when interacting with those attributes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"exclude_global_sensitive_attribute": schema.SetAttribute{
				Description: "Specifies the set of global sensitive attribute definitions that should not apply to this client connection policy.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"result_code_map": schema.StringAttribute{
				Description: "Specifies the result code map that should be used for clients associated with this Client Connection Policy. If a value is defined for this property, then it will override any result code map referenced in the global configuration.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"included_backend_base_dn": schema.SetAttribute{
				Description: "Specifies the set of backend base DNs for which subtree views should be included in this Client Connection Policy.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_backend_base_dn": schema.SetAttribute{
				Description: "Specifies the set of backend base DNs for which subtree views should be excluded from this Client Connection Policy.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"allowed_operation": schema.SetAttribute{
				Description: "Specifies the types of operations that clients associated with this Client Connection Policy will be allowed to request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"required_operation_request_criteria": schema.StringAttribute{
				Description: "Specifies a request criteria object that will be required to match all requests submitted by clients associated with this Client Connection Policy. If a client submits a request that does not satisfy this request criteria object, then that request will be rejected.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"prohibited_operation_request_criteria": schema.StringAttribute{
				Description: "Specifies a request criteria object that must not match any requests submitted by clients associated with this Client Connection Policy. If a client submits a request that satisfies this request criteria object, then that request will be rejected.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allowed_request_control": schema.SetAttribute{
				Description: "Specifies the OIDs of the controls that clients associated with this Client Connection Policy will be allowed to include in requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"denied_request_control": schema.SetAttribute{
				Description: "Specifies the OIDs of the controls that clients associated with this Client Connection Policy will not be allowed to include in requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"allowed_extended_operation": schema.SetAttribute{
				Description: "Specifies the OIDs of the extended operations that clients associated with this Client Connection Policy will be allowed to request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"denied_extended_operation": schema.SetAttribute{
				Description: "Specifies the OIDs of the extended operations that clients associated with this Client Connection Policy will not be allowed to request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"allowed_auth_type": schema.SetAttribute{
				Description: "Specifies the types of authentication that clients associated with this Client Connection Policy will be allowed to request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"allowed_sasl_mechanism": schema.SetAttribute{
				Description: "Specifies the names of the SASL mechanisms that clients associated with this Client Connection Policy will be allowed to request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"denied_sasl_mechanism": schema.SetAttribute{
				Description: "Specifies the names of the SASL mechanisms that clients associated with this Client Connection Policy will not be allowed to request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"allowed_filter_type": schema.SetAttribute{
				Description: "Specifies the types of filter components that may be included in search requests from clients associated with this Client Connection Policy which have a non-baseObject scope.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"allow_unindexed_searches": schema.BoolAttribute{
				Description: "Indicates whether clients will be allowed to request search operations that cannot be efficiently processed using the set of indexes defined in the corresponding backend. Note that even if this is false, some clients may be able to request unindexed searches if the allow-unindexed-searches-with-control property has a value of true and the necessary conditions are satisfied.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_unindexed_searches_with_control": schema.BoolAttribute{
				Description: "Indicates whether clients will be allowed to request search operations that cannot be efficiently processed using the set of indexes defined in the corresponding backend, as long as the search request also includes the permit unindexed search request control and the requester has the unindexed-search-with-control privilege (or that privilege is disabled in the global configuration).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"minimum_substring_length": schema.Int64Attribute{
				Description: "Specifies the minimum number of consecutive bytes that must be present in any subInitial, subAny, or subFinal element of a substring filter component (i.e., the minimum number of consecutive bytes between wildcard characters in a substring filter). Any attempt to use a substring search with an element containing fewer than this number of bytes will be rejected.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_concurrent_connections": schema.Int64Attribute{
				Description: "Specifies the maximum number of client connections which may be associated with this Client Connection Policy at any given time.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_connection_duration": schema.StringAttribute{
				Description: "Specifies the maximum length of time that a connection associated with this Client Connection Policy may be established. Any connection which is associated with this Client Connection Policy and has been established for longer than this period of time may be terminated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_idle_connection_duration": schema.StringAttribute{
				Description: "Specifies the maximum length of time that a connection associated with this Client Connection Policy may remain established after the completion of the last operation processed on that connection. Any new operation requested on the connection will reset this timer. Any connection associated with this Client Connection Policy which has been idle for longer than this length of time may be terminated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_operation_count_per_connection": schema.Int64Attribute{
				Description: "Specifies the maximum number of operations that may be requested by any client connection associated with this Client Connection Policy. If an attempt is made to process more than this number of operations on a client connection, then that connection will be terminated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_concurrent_operations_per_connection": schema.Int64Attribute{
				Description: "Specifies the maximum number of concurrent operations that can be in progress for any connection. This can help prevent a single client connection from monopolizing server processing resources by sending a large number of concurrent asynchronous requests. A value of zero indicates that no limit will be placed on the number of concurrent requests for a single client.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_concurrent_operation_wait_time_before_rejecting": schema.StringAttribute{
				Description: "Specifies the maximum length of time that the server should wait for an outstanding operation to complete before rejecting a new request received when the maximum number of outstanding operations are already in progress on that connection. If an existing outstanding operation on the connection completes before this time, then the operation will be processed. Otherwise, the operation will be rejected with a \"busy\" result.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_concurrent_operations_per_connection_exceeded_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the Directory Server should exhibit if a client attempts to invoke more concurrent operations on a single connection than allowed by the maximum-concurrent-operations-per-connection property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_connection_operation_rate": schema.SetAttribute{
				Description: "Specifies the maximum rate at which a client associated with this Client Connection Policy may issue requests to the Directory Server. If any client attempts to request operations at a rate higher than this limit, then the server will exhibit the behavior described in the connection-operation-rate-exceeded-behavior property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"connection_operation_rate_exceeded_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the Directory Server should exhibit if a client connection attempts to exceed a rate defined in the maximum-connection-operation-rate property. If the configured behavior is one that will reject requested operations, then that behavior will persist until the end of the corresponding interval. The server will resume allowing that client to perform operations when that interval expires, as long as no other operation rate limits have been exceeded.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_policy_operation_rate": schema.SetAttribute{
				Description: "Specifies the maximum rate at which all clients associated with this Client Connection Policy, as a collective set, may issue requests to the Directory Server. If this limit is exceeded, then the server will exhibit the behavior described in the policy-operation-rate-exceeded-behavior property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"policy_operation_rate_exceeded_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the Directory Server should exhibit if a client connection attempts to exceed a rate defined in the maximum-policy-operation-rate property. If the configured behavior is one that will reject requested operations, then that behavior will persist until the end of the corresponding interval. The server will resume allowing clients associated with this Client Connection Policy to perform operations when that interval expires, as long as no other operation rate limits have been exceeded.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_search_size_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that may be returned for a search performed by a client associated with this Client Connection Policy.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_search_time_limit": schema.StringAttribute{
				Description: "Specifies the maximum length of time that the server should spend processing search operations requested by clients associated with this Client Connection Policy.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_search_lookthrough_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that may be examined by a backend in the course of processing a search requested by clients associated with this Client Connection Policy.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_ldap_join_size_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that may be joined with any single search result entry for a search request performed by a client associated with this Client Connection Policy.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_sort_size_limit_without_vlv_index": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that the server will attempt to sort without the benefit of a VLV index. A value of zero indicates that no limit should be enforced.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a ClientConnectionPolicyResponse object into the model struct
func readClientConnectionPolicyResponseDataSource(ctx context.Context, r *client.ClientConnectionPolicyResponse, state *clientConnectionPolicyDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("client-connection-policy")
	state.Id = types.StringValue(r.Id)
	state.PolicyID = types.StringValue(r.PolicyID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.EvaluationOrderIndex = types.Int64Value(r.EvaluationOrderIndex)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.TerminateConnection = internaltypes.BoolTypeOrNil(r.TerminateConnection)
	state.SensitiveAttribute = internaltypes.GetStringSet(r.SensitiveAttribute)
	state.ExcludeGlobalSensitiveAttribute = internaltypes.GetStringSet(r.ExcludeGlobalSensitiveAttribute)
	state.ResultCodeMap = internaltypes.StringTypeOrNil(r.ResultCodeMap, false)
	state.IncludedBackendBaseDN = internaltypes.GetStringSet(r.IncludedBackendBaseDN)
	state.ExcludedBackendBaseDN = internaltypes.GetStringSet(r.ExcludedBackendBaseDN)
	state.AllowedOperation = internaltypes.GetStringSet(
		client.StringSliceEnumclientConnectionPolicyAllowedOperationProp(r.AllowedOperation))
	state.RequiredOperationRequestCriteria = internaltypes.StringTypeOrNil(r.RequiredOperationRequestCriteria, false)
	state.ProhibitedOperationRequestCriteria = internaltypes.StringTypeOrNil(r.ProhibitedOperationRequestCriteria, false)
	state.AllowedRequestControl = internaltypes.GetStringSet(r.AllowedRequestControl)
	state.DeniedRequestControl = internaltypes.GetStringSet(r.DeniedRequestControl)
	state.AllowedExtendedOperation = internaltypes.GetStringSet(r.AllowedExtendedOperation)
	state.DeniedExtendedOperation = internaltypes.GetStringSet(r.DeniedExtendedOperation)
	state.AllowedAuthType = internaltypes.GetStringSet(
		client.StringSliceEnumclientConnectionPolicyAllowedAuthTypeProp(r.AllowedAuthType))
	state.AllowedSASLMechanism = internaltypes.GetStringSet(r.AllowedSASLMechanism)
	state.DeniedSASLMechanism = internaltypes.GetStringSet(r.DeniedSASLMechanism)
	state.AllowedFilterType = internaltypes.GetStringSet(
		client.StringSliceEnumclientConnectionPolicyAllowedFilterTypeProp(r.AllowedFilterType))
	state.AllowUnindexedSearches = internaltypes.BoolTypeOrNil(r.AllowUnindexedSearches)
	state.AllowUnindexedSearchesWithControl = internaltypes.BoolTypeOrNil(r.AllowUnindexedSearchesWithControl)
	state.MinimumSubstringLength = internaltypes.Int64TypeOrNil(r.MinimumSubstringLength)
	state.MaximumConcurrentConnections = internaltypes.Int64TypeOrNil(r.MaximumConcurrentConnections)
	state.MaximumConnectionDuration = internaltypes.StringTypeOrNil(r.MaximumConnectionDuration, false)
	state.MaximumIdleConnectionDuration = internaltypes.StringTypeOrNil(r.MaximumIdleConnectionDuration, false)
	state.MaximumOperationCountPerConnection = internaltypes.Int64TypeOrNil(r.MaximumOperationCountPerConnection)
	state.MaximumConcurrentOperationsPerConnection = internaltypes.Int64TypeOrNil(r.MaximumConcurrentOperationsPerConnection)
	state.MaximumConcurrentOperationWaitTimeBeforeRejecting = internaltypes.StringTypeOrNil(r.MaximumConcurrentOperationWaitTimeBeforeRejecting, false)
	state.MaximumConcurrentOperationsPerConnectionExceededBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumclientConnectionPolicyMaximumConcurrentOperationsPerConnectionExceededBehaviorProp(r.MaximumConcurrentOperationsPerConnectionExceededBehavior), false)
	state.MaximumConnectionOperationRate = internaltypes.GetStringSet(r.MaximumConnectionOperationRate)
	state.ConnectionOperationRateExceededBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumclientConnectionPolicyConnectionOperationRateExceededBehaviorProp(r.ConnectionOperationRateExceededBehavior), false)
	state.MaximumPolicyOperationRate = internaltypes.GetStringSet(r.MaximumPolicyOperationRate)
	state.PolicyOperationRateExceededBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumclientConnectionPolicyPolicyOperationRateExceededBehaviorProp(r.PolicyOperationRateExceededBehavior), false)
	state.MaximumSearchSizeLimit = internaltypes.Int64TypeOrNil(r.MaximumSearchSizeLimit)
	state.MaximumSearchTimeLimit = internaltypes.StringTypeOrNil(r.MaximumSearchTimeLimit, false)
	state.MaximumSearchLookthroughLimit = internaltypes.Int64TypeOrNil(r.MaximumSearchLookthroughLimit)
	state.MaximumLDAPJoinSizeLimit = internaltypes.Int64TypeOrNil(r.MaximumLDAPJoinSizeLimit)
	state.MaximumSortSizeLimitWithoutVLVIndex = internaltypes.Int64TypeOrNil(r.MaximumSortSizeLimitWithoutVLVIndex)
}

// Read resource information
func (r *clientConnectionPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state clientConnectionPolicyDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ClientConnectionPolicyApi.GetClientConnectionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.PolicyID.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Client Connection Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readClientConnectionPolicyResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
