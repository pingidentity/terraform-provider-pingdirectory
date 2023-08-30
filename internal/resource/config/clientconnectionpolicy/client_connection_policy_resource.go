package clientconnectionpolicy

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &clientConnectionPolicyResource{}
	_ resource.ResourceWithConfigure   = &clientConnectionPolicyResource{}
	_ resource.ResourceWithImportState = &clientConnectionPolicyResource{}
	_ resource.Resource                = &defaultClientConnectionPolicyResource{}
	_ resource.ResourceWithConfigure   = &defaultClientConnectionPolicyResource{}
	_ resource.ResourceWithImportState = &defaultClientConnectionPolicyResource{}
)

// Create a Client Connection Policy resource
func NewClientConnectionPolicyResource() resource.Resource {
	return &clientConnectionPolicyResource{}
}

func NewDefaultClientConnectionPolicyResource() resource.Resource {
	return &defaultClientConnectionPolicyResource{}
}

// clientConnectionPolicyResource is the resource implementation.
type clientConnectionPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultClientConnectionPolicyResource is the resource implementation.
type defaultClientConnectionPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *clientConnectionPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client_connection_policy"
}

func (r *defaultClientConnectionPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_client_connection_policy"
}

// Configure adds the provider configured client to the resource.
func (r *clientConnectionPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultClientConnectionPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type clientConnectionPolicyResourceModel struct {
	Id                                                       types.String `tfsdk:"id"`
	LastUpdated                                              types.String `tfsdk:"last_updated"`
	Notifications                                            types.Set    `tfsdk:"notifications"`
	RequiredActions                                          types.Set    `tfsdk:"required_actions"`
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

// GetSchema defines the schema for the resource.
func (r *clientConnectionPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	clientConnectionPolicySchema(ctx, req, resp, false)
}

func (r *defaultClientConnectionPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	clientConnectionPolicySchema(ctx, req, resp, true)
}

func clientConnectionPolicySchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	allowedOperationDefaults, diags := types.SetValue(types.StringType, []attr.Value{types.StringValue("abandon"), types.StringValue("add"), types.StringValue("bind"), types.StringValue("compare"), types.StringValue("delete"), types.StringValue("extended"), types.StringValue("modify"), types.StringValue("modify-dn"), types.StringValue("search")})
	resp.Diagnostics.Append(diags...)
	allowedAuthTypeDefaults, diags := types.SetValue(types.StringType, []attr.Value{types.StringValue("simple"), types.StringValue("sasl")})
	resp.Diagnostics.Append(diags...)
	allowedFilterTypeDefaults, diags := types.SetValue(types.StringType, []attr.Value{types.StringValue("and"), types.StringValue("or"), types.StringValue("not"), types.StringValue("equality"), types.StringValue("sub-initial"), types.StringValue("sub-any"), types.StringValue("sub-final"), types.StringValue("greater-or-equal"), types.StringValue("less-or-equal"), types.StringValue("present"), types.StringValue("approximate-match"), types.StringValue("extensible-match")})
	resp.Diagnostics.Append(diags...)
	schemaDef := schema.Schema{
		Description: "Manages a Client Connection Policy.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Client Connection Policy resource. Options are ['client-connection-policy']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("client-connection-policy"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"client-connection-policy"}...),
				},
			},
			"policy_id": schema.StringAttribute{
				Description: "Specifies a name which uniquely identifies this Client Connection Policy in the server.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Client Connection Policy",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Client Connection Policy is enabled for use in the server. If a Client Connection Policy is disabled, then no new client connections will be associated with it.",
				Required:    true,
			},
			"evaluation_order_index": schema.Int64Attribute{
				Description: "Specifies the order in which Client Connection Policy definitions will be evaluated. A Client Connection Policy with a lower index will be evaluated before one with a higher index, and the first Client Connection Policy evaluated which may apply to a client connection will be used for that connection. Each Client Connection Policy must be assigned a unique evaluation order index value.",
				Required:    true,
			},
			"connection_criteria": schema.StringAttribute{
				Description: "Specifies a set of connection criteria that must match the associated client connection for it to be associated with this Client Connection Policy.",
				Optional:    true,
			},
			"terminate_connection": schema.BoolAttribute{
				Description: "Indicates whether any client connection for which this Client Connection Policy is selected should be terminated. This makes it possible to define fine-grained criteria for clients that should not be allowed to connect to this Directory Server.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"sensitive_attribute": schema.SetAttribute{
				Description: "Provides the ability to indicate that some attributes should be considered sensitive and additional protection should be in place when interacting with those attributes.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"exclude_global_sensitive_attribute": schema.SetAttribute{
				Description: "Specifies the set of global sensitive attribute definitions that should not apply to this client connection policy.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"result_code_map": schema.StringAttribute{
				Description: "Specifies the result code map that should be used for clients associated with this Client Connection Policy. If a value is defined for this property, then it will override any result code map referenced in the global configuration.",
				Optional:    true,
			},
			"included_backend_base_dn": schema.SetAttribute{
				Description: "Specifies the set of backend base DNs for which subtree views should be included in this Client Connection Policy.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"excluded_backend_base_dn": schema.SetAttribute{
				Description: "Specifies the set of backend base DNs for which subtree views should be excluded from this Client Connection Policy.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"allowed_operation": schema.SetAttribute{
				Description: "Specifies the types of operations that clients associated with this Client Connection Policy will be allowed to request.",
				Optional:    true,
				Computed:    true,
				Default:     setdefault.StaticValue(allowedOperationDefaults),
				ElementType: types.StringType,
			},
			"required_operation_request_criteria": schema.StringAttribute{
				Description: "Specifies a request criteria object that will be required to match all requests submitted by clients associated with this Client Connection Policy. If a client submits a request that does not satisfy this request criteria object, then that request will be rejected.",
				Optional:    true,
			},
			"prohibited_operation_request_criteria": schema.StringAttribute{
				Description: "Specifies a request criteria object that must not match any requests submitted by clients associated with this Client Connection Policy. If a client submits a request that satisfies this request criteria object, then that request will be rejected.",
				Optional:    true,
			},
			"allowed_request_control": schema.SetAttribute{
				Description: "Specifies the OIDs of the controls that clients associated with this Client Connection Policy will be allowed to include in requests.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"denied_request_control": schema.SetAttribute{
				Description: "Specifies the OIDs of the controls that clients associated with this Client Connection Policy will not be allowed to include in requests.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"allowed_extended_operation": schema.SetAttribute{
				Description: "Specifies the OIDs of the extended operations that clients associated with this Client Connection Policy will be allowed to request.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"denied_extended_operation": schema.SetAttribute{
				Description: "Specifies the OIDs of the extended operations that clients associated with this Client Connection Policy will not be allowed to request.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"allowed_auth_type": schema.SetAttribute{
				Description: "Specifies the types of authentication that clients associated with this Client Connection Policy will be allowed to request.",
				Optional:    true,
				Computed:    true,
				Default:     setdefault.StaticValue(allowedAuthTypeDefaults),
				ElementType: types.StringType,
			},
			"allowed_sasl_mechanism": schema.SetAttribute{
				Description: "Specifies the names of the SASL mechanisms that clients associated with this Client Connection Policy will be allowed to request.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"denied_sasl_mechanism": schema.SetAttribute{
				Description: "Specifies the names of the SASL mechanisms that clients associated with this Client Connection Policy will not be allowed to request.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"allowed_filter_type": schema.SetAttribute{
				Description: "Specifies the types of filter components that may be included in search requests from clients associated with this Client Connection Policy which have a non-baseObject scope.",
				Optional:    true,
				Computed:    true,
				Default:     setdefault.StaticValue(allowedFilterTypeDefaults),
				ElementType: types.StringType,
			},
			"allow_unindexed_searches": schema.BoolAttribute{
				Description: "Indicates whether clients will be allowed to request search operations that cannot be efficiently processed using the set of indexes defined in the corresponding backend. Note that even if this is false, some clients may be able to request unindexed searches if the allow-unindexed-searches-with-control property has a value of true and the necessary conditions are satisfied.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"allow_unindexed_searches_with_control": schema.BoolAttribute{
				Description: "Indicates whether clients will be allowed to request search operations that cannot be efficiently processed using the set of indexes defined in the corresponding backend, as long as the search request also includes the permit unindexed search request control and the requester has the unindexed-search-with-control privilege (or that privilege is disabled in the global configuration).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"minimum_substring_length": schema.Int64Attribute{
				Description: "Specifies the minimum number of consecutive bytes that must be present in any subInitial, subAny, or subFinal element of a substring filter component (i.e., the minimum number of consecutive bytes between wildcard characters in a substring filter). Any attempt to use a substring search with an element containing fewer than this number of bytes will be rejected.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
			},
			"maximum_concurrent_connections": schema.Int64Attribute{
				Description: "Specifies the maximum number of client connections which may be associated with this Client Connection Policy at any given time.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"maximum_connection_duration": schema.StringAttribute{
				Description: "Specifies the maximum length of time that a connection associated with this Client Connection Policy may be established. Any connection which is associated with this Client Connection Policy and has been established for longer than this period of time may be terminated.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"maximum_idle_connection_duration": schema.StringAttribute{
				Description: "Specifies the maximum length of time that a connection associated with this Client Connection Policy may remain established after the completion of the last operation processed on that connection. Any new operation requested on the connection will reset this timer. Any connection associated with this Client Connection Policy which has been idle for longer than this length of time may be terminated.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"maximum_operation_count_per_connection": schema.Int64Attribute{
				Description: "Specifies the maximum number of operations that may be requested by any client connection associated with this Client Connection Policy. If an attempt is made to process more than this number of operations on a client connection, then that connection will be terminated.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"maximum_concurrent_operations_per_connection": schema.Int64Attribute{
				Description: "Specifies the maximum number of concurrent operations that can be in progress for any connection. This can help prevent a single client connection from monopolizing server processing resources by sending a large number of concurrent asynchronous requests. A value of zero indicates that no limit will be placed on the number of concurrent requests for a single client.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"maximum_concurrent_operation_wait_time_before_rejecting": schema.StringAttribute{
				Description: "Specifies the maximum length of time that the server should wait for an outstanding operation to complete before rejecting a new request received when the maximum number of outstanding operations are already in progress on that connection. If an existing outstanding operation on the connection completes before this time, then the operation will be processed. Otherwise, the operation will be rejected with a \"busy\" result.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"maximum_concurrent_operations_per_connection_exceeded_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the Directory Server should exhibit if a client attempts to invoke more concurrent operations on a single connection than allowed by the maximum-concurrent-operations-per-connection property.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("reject-busy"),
			},
			"maximum_connection_operation_rate": schema.SetAttribute{
				Description: "Specifies the maximum rate at which a client associated with this Client Connection Policy may issue requests to the Directory Server. If any client attempts to request operations at a rate higher than this limit, then the server will exhibit the behavior described in the connection-operation-rate-exceeded-behavior property.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"connection_operation_rate_exceeded_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the Directory Server should exhibit if a client connection attempts to exceed a rate defined in the maximum-connection-operation-rate property. If the configured behavior is one that will reject requested operations, then that behavior will persist until the end of the corresponding interval. The server will resume allowing that client to perform operations when that interval expires, as long as no other operation rate limits have been exceeded.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("reject-busy"),
			},
			"maximum_policy_operation_rate": schema.SetAttribute{
				Description: "Specifies the maximum rate at which all clients associated with this Client Connection Policy, as a collective set, may issue requests to the Directory Server. If this limit is exceeded, then the server will exhibit the behavior described in the policy-operation-rate-exceeded-behavior property.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"policy_operation_rate_exceeded_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the Directory Server should exhibit if a client connection attempts to exceed a rate defined in the maximum-policy-operation-rate property. If the configured behavior is one that will reject requested operations, then that behavior will persist until the end of the corresponding interval. The server will resume allowing clients associated with this Client Connection Policy to perform operations when that interval expires, as long as no other operation rate limits have been exceeded.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("reject-busy"),
			},
			"maximum_search_size_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that may be returned for a search performed by a client associated with this Client Connection Policy.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"maximum_search_time_limit": schema.StringAttribute{
				Description: "Specifies the maximum length of time that the server should spend processing search operations requested by clients associated with this Client Connection Policy.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"maximum_search_lookthrough_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that may be examined by a backend in the course of processing a search requested by clients associated with this Client Connection Policy.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"maximum_ldap_join_size_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that may be joined with any single search result entry for a search request performed by a client associated with this Client Connection Policy.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"maximum_sort_size_limit_without_vlv_index": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that the server will attempt to sort without the benefit of a VLV index. A value of zero indicates that no limit should be enforced.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
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
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type", "policy_id"})
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Add optional fields to create request for client-connection-policy client-connection-policy
func addOptionalClientConnectionPolicyFields(ctx context.Context, addRequest *client.AddClientConnectionPolicyRequest, plan clientConnectionPolicyResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		addRequest.ConnectionCriteria = plan.ConnectionCriteria.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.TerminateConnection) {
		addRequest.TerminateConnection = plan.TerminateConnection.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SensitiveAttribute) {
		var slice []string
		plan.SensitiveAttribute.ElementsAs(ctx, &slice, false)
		addRequest.SensitiveAttribute = slice
	}
	if internaltypes.IsDefined(plan.ExcludeGlobalSensitiveAttribute) {
		var slice []string
		plan.ExcludeGlobalSensitiveAttribute.ElementsAs(ctx, &slice, false)
		addRequest.ExcludeGlobalSensitiveAttribute = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ResultCodeMap) {
		addRequest.ResultCodeMap = plan.ResultCodeMap.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludedBackendBaseDN) {
		var slice []string
		plan.IncludedBackendBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.IncludedBackendBaseDN = slice
	}
	if internaltypes.IsDefined(plan.ExcludedBackendBaseDN) {
		var slice []string
		plan.ExcludedBackendBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedBackendBaseDN = slice
	}
	if internaltypes.IsDefined(plan.AllowedOperation) {
		var slice []string
		plan.AllowedOperation.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumclientConnectionPolicyAllowedOperationProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumclientConnectionPolicyAllowedOperationPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AllowedOperation = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequiredOperationRequestCriteria) {
		addRequest.RequiredOperationRequestCriteria = plan.RequiredOperationRequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ProhibitedOperationRequestCriteria) {
		addRequest.ProhibitedOperationRequestCriteria = plan.ProhibitedOperationRequestCriteria.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AllowedRequestControl) {
		var slice []string
		plan.AllowedRequestControl.ElementsAs(ctx, &slice, false)
		addRequest.AllowedRequestControl = slice
	}
	if internaltypes.IsDefined(plan.DeniedRequestControl) {
		var slice []string
		plan.DeniedRequestControl.ElementsAs(ctx, &slice, false)
		addRequest.DeniedRequestControl = slice
	}
	if internaltypes.IsDefined(plan.AllowedExtendedOperation) {
		var slice []string
		plan.AllowedExtendedOperation.ElementsAs(ctx, &slice, false)
		addRequest.AllowedExtendedOperation = slice
	}
	if internaltypes.IsDefined(plan.DeniedExtendedOperation) {
		var slice []string
		plan.DeniedExtendedOperation.ElementsAs(ctx, &slice, false)
		addRequest.DeniedExtendedOperation = slice
	}
	if internaltypes.IsDefined(plan.AllowedAuthType) {
		var slice []string
		plan.AllowedAuthType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumclientConnectionPolicyAllowedAuthTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumclientConnectionPolicyAllowedAuthTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AllowedAuthType = enumSlice
	}
	if internaltypes.IsDefined(plan.AllowedSASLMechanism) {
		var slice []string
		plan.AllowedSASLMechanism.ElementsAs(ctx, &slice, false)
		addRequest.AllowedSASLMechanism = slice
	}
	if internaltypes.IsDefined(plan.DeniedSASLMechanism) {
		var slice []string
		plan.DeniedSASLMechanism.ElementsAs(ctx, &slice, false)
		addRequest.DeniedSASLMechanism = slice
	}
	if internaltypes.IsDefined(plan.AllowedFilterType) {
		var slice []string
		plan.AllowedFilterType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumclientConnectionPolicyAllowedFilterTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumclientConnectionPolicyAllowedFilterTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AllowedFilterType = enumSlice
	}
	if internaltypes.IsDefined(plan.AllowUnindexedSearches) {
		addRequest.AllowUnindexedSearches = plan.AllowUnindexedSearches.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AllowUnindexedSearchesWithControl) {
		addRequest.AllowUnindexedSearchesWithControl = plan.AllowUnindexedSearchesWithControl.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MinimumSubstringLength) {
		addRequest.MinimumSubstringLength = plan.MinimumSubstringLength.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaximumConcurrentConnections) {
		addRequest.MaximumConcurrentConnections = plan.MaximumConcurrentConnections.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaximumConnectionDuration) {
		addRequest.MaximumConnectionDuration = plan.MaximumConnectionDuration.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaximumIdleConnectionDuration) {
		addRequest.MaximumIdleConnectionDuration = plan.MaximumIdleConnectionDuration.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.MaximumOperationCountPerConnection) {
		addRequest.MaximumOperationCountPerConnection = plan.MaximumOperationCountPerConnection.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaximumConcurrentOperationsPerConnection) {
		addRequest.MaximumConcurrentOperationsPerConnection = plan.MaximumConcurrentOperationsPerConnection.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaximumConcurrentOperationWaitTimeBeforeRejecting) {
		addRequest.MaximumConcurrentOperationWaitTimeBeforeRejecting = plan.MaximumConcurrentOperationWaitTimeBeforeRejecting.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaximumConcurrentOperationsPerConnectionExceededBehavior) {
		maximumConcurrentOperationsPerConnectionExceededBehavior, err := client.NewEnumclientConnectionPolicyMaximumConcurrentOperationsPerConnectionExceededBehaviorPropFromValue(plan.MaximumConcurrentOperationsPerConnectionExceededBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MaximumConcurrentOperationsPerConnectionExceededBehavior = maximumConcurrentOperationsPerConnectionExceededBehavior
	}
	if internaltypes.IsDefined(plan.MaximumConnectionOperationRate) {
		var slice []string
		plan.MaximumConnectionOperationRate.ElementsAs(ctx, &slice, false)
		addRequest.MaximumConnectionOperationRate = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionOperationRateExceededBehavior) {
		connectionOperationRateExceededBehavior, err := client.NewEnumclientConnectionPolicyConnectionOperationRateExceededBehaviorPropFromValue(plan.ConnectionOperationRateExceededBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConnectionOperationRateExceededBehavior = connectionOperationRateExceededBehavior
	}
	if internaltypes.IsDefined(plan.MaximumPolicyOperationRate) {
		var slice []string
		plan.MaximumPolicyOperationRate.ElementsAs(ctx, &slice, false)
		addRequest.MaximumPolicyOperationRate = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PolicyOperationRateExceededBehavior) {
		policyOperationRateExceededBehavior, err := client.NewEnumclientConnectionPolicyPolicyOperationRateExceededBehaviorPropFromValue(plan.PolicyOperationRateExceededBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.PolicyOperationRateExceededBehavior = policyOperationRateExceededBehavior
	}
	if internaltypes.IsDefined(plan.MaximumSearchSizeLimit) {
		addRequest.MaximumSearchSizeLimit = plan.MaximumSearchSizeLimit.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaximumSearchTimeLimit) {
		addRequest.MaximumSearchTimeLimit = plan.MaximumSearchTimeLimit.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.MaximumSearchLookthroughLimit) {
		addRequest.MaximumSearchLookthroughLimit = plan.MaximumSearchLookthroughLimit.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaximumLDAPJoinSizeLimit) {
		addRequest.MaximumLDAPJoinSizeLimit = plan.MaximumLDAPJoinSizeLimit.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaximumSortSizeLimitWithoutVLVIndex) {
		addRequest.MaximumSortSizeLimitWithoutVLVIndex = plan.MaximumSortSizeLimitWithoutVLVIndex.ValueInt64Pointer()
	}
	return nil
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *clientConnectionPolicyResourceModel) populateAllComputedStringAttributes() {
	if model.MaximumIdleConnectionDuration.IsUnknown() || model.MaximumIdleConnectionDuration.IsNull() {
		model.MaximumIdleConnectionDuration = types.StringValue("")
	}
	if model.ConnectionOperationRateExceededBehavior.IsUnknown() || model.ConnectionOperationRateExceededBehavior.IsNull() {
		model.ConnectionOperationRateExceededBehavior = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.MaximumSearchTimeLimit.IsUnknown() || model.MaximumSearchTimeLimit.IsNull() {
		model.MaximumSearchTimeLimit = types.StringValue("")
	}
	if model.MaximumConcurrentOperationWaitTimeBeforeRejecting.IsUnknown() || model.MaximumConcurrentOperationWaitTimeBeforeRejecting.IsNull() {
		model.MaximumConcurrentOperationWaitTimeBeforeRejecting = types.StringValue("")
	}
	if model.MaximumConnectionDuration.IsUnknown() || model.MaximumConnectionDuration.IsNull() {
		model.MaximumConnectionDuration = types.StringValue("")
	}
	if model.ConnectionCriteria.IsUnknown() || model.ConnectionCriteria.IsNull() {
		model.ConnectionCriteria = types.StringValue("")
	}
	if model.ResultCodeMap.IsUnknown() || model.ResultCodeMap.IsNull() {
		model.ResultCodeMap = types.StringValue("")
	}
	if model.ProhibitedOperationRequestCriteria.IsUnknown() || model.ProhibitedOperationRequestCriteria.IsNull() {
		model.ProhibitedOperationRequestCriteria = types.StringValue("")
	}
	if model.RequiredOperationRequestCriteria.IsUnknown() || model.RequiredOperationRequestCriteria.IsNull() {
		model.RequiredOperationRequestCriteria = types.StringValue("")
	}
	if model.PolicyOperationRateExceededBehavior.IsUnknown() || model.PolicyOperationRateExceededBehavior.IsNull() {
		model.PolicyOperationRateExceededBehavior = types.StringValue("")
	}
	if model.MaximumConcurrentOperationsPerConnectionExceededBehavior.IsUnknown() || model.MaximumConcurrentOperationsPerConnectionExceededBehavior.IsNull() {
		model.MaximumConcurrentOperationsPerConnectionExceededBehavior = types.StringValue("")
	}
	if model.PolicyID.IsUnknown() || model.PolicyID.IsNull() {
		model.PolicyID = types.StringValue("")
	}
}

// Read a ClientConnectionPolicyResponse object into the model struct
func readClientConnectionPolicyResponse(ctx context.Context, r *client.ClientConnectionPolicyResponse, state *clientConnectionPolicyResourceModel, expectedValues *clientConnectionPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("client-connection-policy")
	state.Id = types.StringValue(r.Id)
	state.PolicyID = types.StringValue(r.PolicyID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.EvaluationOrderIndex = types.Int64Value(r.EvaluationOrderIndex)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.TerminateConnection = internaltypes.BoolTypeOrNil(r.TerminateConnection)
	state.SensitiveAttribute = internaltypes.GetStringSet(r.SensitiveAttribute)
	state.ExcludeGlobalSensitiveAttribute = internaltypes.GetStringSet(r.ExcludeGlobalSensitiveAttribute)
	state.ResultCodeMap = internaltypes.StringTypeOrNil(r.ResultCodeMap, internaltypes.IsEmptyString(expectedValues.ResultCodeMap))
	state.IncludedBackendBaseDN = internaltypes.GetStringSet(r.IncludedBackendBaseDN)
	state.ExcludedBackendBaseDN = internaltypes.GetStringSet(r.ExcludedBackendBaseDN)
	state.AllowedOperation = internaltypes.GetStringSet(
		client.StringSliceEnumclientConnectionPolicyAllowedOperationProp(r.AllowedOperation))
	state.RequiredOperationRequestCriteria = internaltypes.StringTypeOrNil(r.RequiredOperationRequestCriteria, internaltypes.IsEmptyString(expectedValues.RequiredOperationRequestCriteria))
	state.ProhibitedOperationRequestCriteria = internaltypes.StringTypeOrNil(r.ProhibitedOperationRequestCriteria, internaltypes.IsEmptyString(expectedValues.ProhibitedOperationRequestCriteria))
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
	state.MaximumConnectionDuration = internaltypes.StringTypeOrNil(r.MaximumConnectionDuration, true)
	config.CheckMismatchedPDFormattedAttributes("maximum_connection_duration",
		expectedValues.MaximumConnectionDuration, state.MaximumConnectionDuration, diagnostics)
	state.MaximumIdleConnectionDuration = internaltypes.StringTypeOrNil(r.MaximumIdleConnectionDuration, true)
	config.CheckMismatchedPDFormattedAttributes("maximum_idle_connection_duration",
		expectedValues.MaximumIdleConnectionDuration, state.MaximumIdleConnectionDuration, diagnostics)
	state.MaximumOperationCountPerConnection = internaltypes.Int64TypeOrNil(r.MaximumOperationCountPerConnection)
	state.MaximumConcurrentOperationsPerConnection = internaltypes.Int64TypeOrNil(r.MaximumConcurrentOperationsPerConnection)
	state.MaximumConcurrentOperationWaitTimeBeforeRejecting = internaltypes.StringTypeOrNil(r.MaximumConcurrentOperationWaitTimeBeforeRejecting, true)
	config.CheckMismatchedPDFormattedAttributes("maximum_concurrent_operation_wait_time_before_rejecting",
		expectedValues.MaximumConcurrentOperationWaitTimeBeforeRejecting, state.MaximumConcurrentOperationWaitTimeBeforeRejecting, diagnostics)
	state.MaximumConcurrentOperationsPerConnectionExceededBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumclientConnectionPolicyMaximumConcurrentOperationsPerConnectionExceededBehaviorProp(r.MaximumConcurrentOperationsPerConnectionExceededBehavior), true)
	state.MaximumConnectionOperationRate = internaltypes.GetStringSet(r.MaximumConnectionOperationRate)
	state.ConnectionOperationRateExceededBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumclientConnectionPolicyConnectionOperationRateExceededBehaviorProp(r.ConnectionOperationRateExceededBehavior), true)
	state.MaximumPolicyOperationRate = internaltypes.GetStringSet(r.MaximumPolicyOperationRate)
	state.PolicyOperationRateExceededBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumclientConnectionPolicyPolicyOperationRateExceededBehaviorProp(r.PolicyOperationRateExceededBehavior), true)
	state.MaximumSearchSizeLimit = internaltypes.Int64TypeOrNil(r.MaximumSearchSizeLimit)
	state.MaximumSearchTimeLimit = internaltypes.StringTypeOrNil(r.MaximumSearchTimeLimit, true)
	config.CheckMismatchedPDFormattedAttributes("maximum_search_time_limit",
		expectedValues.MaximumSearchTimeLimit, state.MaximumSearchTimeLimit, diagnostics)
	state.MaximumSearchLookthroughLimit = internaltypes.Int64TypeOrNil(r.MaximumSearchLookthroughLimit)
	state.MaximumLDAPJoinSizeLimit = internaltypes.Int64TypeOrNil(r.MaximumLDAPJoinSizeLimit)
	state.MaximumSortSizeLimitWithoutVLVIndex = internaltypes.Int64TypeOrNil(r.MaximumSortSizeLimitWithoutVLVIndex)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createClientConnectionPolicyOperations(plan clientConnectionPolicyResourceModel, state clientConnectionPolicyResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.PolicyID, state.PolicyID, "policy-id")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddInt64OperationIfNecessary(&ops, plan.EvaluationOrderIndex, state.EvaluationOrderIndex, "evaluation-order-index")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionCriteria, state.ConnectionCriteria, "connection-criteria")
	operations.AddBoolOperationIfNecessary(&ops, plan.TerminateConnection, state.TerminateConnection, "terminate-connection")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SensitiveAttribute, state.SensitiveAttribute, "sensitive-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeGlobalSensitiveAttribute, state.ExcludeGlobalSensitiveAttribute, "exclude-global-sensitive-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.ResultCodeMap, state.ResultCodeMap, "result-code-map")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedBackendBaseDN, state.IncludedBackendBaseDN, "included-backend-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedBackendBaseDN, state.ExcludedBackendBaseDN, "excluded-backend-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedOperation, state.AllowedOperation, "allowed-operation")
	operations.AddStringOperationIfNecessary(&ops, plan.RequiredOperationRequestCriteria, state.RequiredOperationRequestCriteria, "required-operation-request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.ProhibitedOperationRequestCriteria, state.ProhibitedOperationRequestCriteria, "prohibited-operation-request-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedRequestControl, state.AllowedRequestControl, "allowed-request-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DeniedRequestControl, state.DeniedRequestControl, "denied-request-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedExtendedOperation, state.AllowedExtendedOperation, "allowed-extended-operation")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DeniedExtendedOperation, state.DeniedExtendedOperation, "denied-extended-operation")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedAuthType, state.AllowedAuthType, "allowed-auth-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedSASLMechanism, state.AllowedSASLMechanism, "allowed-sasl-mechanism")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DeniedSASLMechanism, state.DeniedSASLMechanism, "denied-sasl-mechanism")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedFilterType, state.AllowedFilterType, "allowed-filter-type")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowUnindexedSearches, state.AllowUnindexedSearches, "allow-unindexed-searches")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowUnindexedSearchesWithControl, state.AllowUnindexedSearchesWithControl, "allow-unindexed-searches-with-control")
	operations.AddInt64OperationIfNecessary(&ops, plan.MinimumSubstringLength, state.MinimumSubstringLength, "minimum-substring-length")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumConcurrentConnections, state.MaximumConcurrentConnections, "maximum-concurrent-connections")
	operations.AddStringOperationIfNecessary(&ops, plan.MaximumConnectionDuration, state.MaximumConnectionDuration, "maximum-connection-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.MaximumIdleConnectionDuration, state.MaximumIdleConnectionDuration, "maximum-idle-connection-duration")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumOperationCountPerConnection, state.MaximumOperationCountPerConnection, "maximum-operation-count-per-connection")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumConcurrentOperationsPerConnection, state.MaximumConcurrentOperationsPerConnection, "maximum-concurrent-operations-per-connection")
	operations.AddStringOperationIfNecessary(&ops, plan.MaximumConcurrentOperationWaitTimeBeforeRejecting, state.MaximumConcurrentOperationWaitTimeBeforeRejecting, "maximum-concurrent-operation-wait-time-before-rejecting")
	operations.AddStringOperationIfNecessary(&ops, plan.MaximumConcurrentOperationsPerConnectionExceededBehavior, state.MaximumConcurrentOperationsPerConnectionExceededBehavior, "maximum-concurrent-operations-per-connection-exceeded-behavior")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MaximumConnectionOperationRate, state.MaximumConnectionOperationRate, "maximum-connection-operation-rate")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionOperationRateExceededBehavior, state.ConnectionOperationRateExceededBehavior, "connection-operation-rate-exceeded-behavior")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MaximumPolicyOperationRate, state.MaximumPolicyOperationRate, "maximum-policy-operation-rate")
	operations.AddStringOperationIfNecessary(&ops, plan.PolicyOperationRateExceededBehavior, state.PolicyOperationRateExceededBehavior, "policy-operation-rate-exceeded-behavior")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumSearchSizeLimit, state.MaximumSearchSizeLimit, "maximum-search-size-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.MaximumSearchTimeLimit, state.MaximumSearchTimeLimit, "maximum-search-time-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumSearchLookthroughLimit, state.MaximumSearchLookthroughLimit, "maximum-search-lookthrough-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumLDAPJoinSizeLimit, state.MaximumLDAPJoinSizeLimit, "maximum-ldap-join-size-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumSortSizeLimitWithoutVLVIndex, state.MaximumSortSizeLimitWithoutVLVIndex, "maximum-sort-size-limit-without-vlv-index")
	return ops
}

// Create a client-connection-policy client-connection-policy
func (r *clientConnectionPolicyResource) CreateClientConnectionPolicy(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan clientConnectionPolicyResourceModel) (*clientConnectionPolicyResourceModel, error) {
	addRequest := client.NewAddClientConnectionPolicyRequest(plan.PolicyID.ValueString(),
		plan.PolicyID.ValueString(),
		plan.Enabled.ValueBool(),
		plan.EvaluationOrderIndex.ValueInt64())
	err := addOptionalClientConnectionPolicyFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Client Connection Policy", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ClientConnectionPolicyApi.AddClientConnectionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddClientConnectionPolicyRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.ClientConnectionPolicyApi.AddClientConnectionPolicyExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Client Connection Policy", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state clientConnectionPolicyResourceModel
	readClientConnectionPolicyResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *clientConnectionPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan clientConnectionPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateClientConnectionPolicy(ctx, req, resp, plan)
	if err != nil {
		return
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
func (r *defaultClientConnectionPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan clientConnectionPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ClientConnectionPolicyApi.GetClientConnectionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.PolicyID.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Client Connection Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state clientConnectionPolicyResourceModel
	readClientConnectionPolicyResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ClientConnectionPolicyApi.UpdateClientConnectionPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.PolicyID.ValueString())
	ops := createClientConnectionPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ClientConnectionPolicyApi.UpdateClientConnectionPolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Client Connection Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readClientConnectionPolicyResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	}

	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *clientConnectionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readClientConnectionPolicy(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultClientConnectionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readClientConnectionPolicy(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readClientConnectionPolicy(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state clientConnectionPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ClientConnectionPolicyApi.GetClientConnectionPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.PolicyID.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Client Connection Policy", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Client Connection Policy", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readClientConnectionPolicyResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *clientConnectionPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateClientConnectionPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultClientConnectionPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateClientConnectionPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateClientConnectionPolicy(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan clientConnectionPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state clientConnectionPolicyResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ClientConnectionPolicyApi.UpdateClientConnectionPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.PolicyID.ValueString())

	// Determine what update operations are necessary
	ops := createClientConnectionPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ClientConnectionPolicyApi.UpdateClientConnectionPolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Client Connection Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readClientConnectionPolicyResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultClientConnectionPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *clientConnectionPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state clientConnectionPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ClientConnectionPolicyApi.DeleteClientConnectionPolicyExecute(r.apiClient.ClientConnectionPolicyApi.DeleteClientConnectionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.PolicyID.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Client Connection Policy", err, httpResp)
		return
	}
}

func (r *clientConnectionPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importClientConnectionPolicy(ctx, req, resp)
}

func (r *defaultClientConnectionPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importClientConnectionPolicy(ctx, req, resp)
}

func importClientConnectionPolicy(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to policy_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("policy_id"), req, resp)
}
