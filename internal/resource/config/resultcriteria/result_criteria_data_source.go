package resultcriteria

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10100/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &resultCriteriaDataSource{}
	_ datasource.DataSourceWithConfigure = &resultCriteriaDataSource{}
)

// Create a Result Criteria data source
func NewResultCriteriaDataSource() datasource.DataSource {
	return &resultCriteriaDataSource{}
}

// resultCriteriaDataSource is the datasource implementation.
type resultCriteriaDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *resultCriteriaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_result_criteria"
}

// Configure adds the provider configured client to the data source.
func (r *resultCriteriaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type resultCriteriaDataSourceModel struct {
	Id                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	Type                              types.String `tfsdk:"type"`
	ExtensionClass                    types.String `tfsdk:"extension_class"`
	ExtensionArgument                 types.Set    `tfsdk:"extension_argument"`
	LocalAssuranceLevel               types.Set    `tfsdk:"local_assurance_level"`
	RemoteAssuranceLevel              types.Set    `tfsdk:"remote_assurance_level"`
	AssuranceTimeoutCriteria          types.String `tfsdk:"assurance_timeout_criteria"`
	AssuranceTimeoutValue             types.String `tfsdk:"assurance_timeout_value"`
	ResponseDelayedByAssurance        types.String `tfsdk:"response_delayed_by_assurance"`
	AssuranceBehaviorAlteredByControl types.String `tfsdk:"assurance_behavior_altered_by_control"`
	AssuranceSatisfied                types.String `tfsdk:"assurance_satisfied"`
	AllIncludedResultCriteria         types.Set    `tfsdk:"all_included_result_criteria"`
	AnyIncludedResultCriteria         types.Set    `tfsdk:"any_included_result_criteria"`
	NotAllIncludedResultCriteria      types.Set    `tfsdk:"not_all_included_result_criteria"`
	NoneIncludedResultCriteria        types.Set    `tfsdk:"none_included_result_criteria"`
	RequestCriteria                   types.String `tfsdk:"request_criteria"`
	ResultCodeCriteria                types.String `tfsdk:"result_code_criteria"`
	ResultCodeValue                   types.Set    `tfsdk:"result_code_value"`
	ProcessingTimeCriteria            types.String `tfsdk:"processing_time_criteria"`
	ProcessingTimeValue               types.String `tfsdk:"processing_time_value"`
	QueueTimeCriteria                 types.String `tfsdk:"queue_time_criteria"`
	QueueTimeValue                    types.String `tfsdk:"queue_time_value"`
	ReferralReturned                  types.String `tfsdk:"referral_returned"`
	AllIncludedResponseControl        types.Set    `tfsdk:"all_included_response_control"`
	AnyIncludedResponseControl        types.Set    `tfsdk:"any_included_response_control"`
	NotAllIncludedResponseControl     types.Set    `tfsdk:"not_all_included_response_control"`
	NoneIncludedResponseControl       types.Set    `tfsdk:"none_included_response_control"`
	UsedAlternateAuthzid              types.String `tfsdk:"used_alternate_authzid"`
	UsedAnyPrivilege                  types.String `tfsdk:"used_any_privilege"`
	UsedPrivilege                     types.Set    `tfsdk:"used_privilege"`
	MissingAnyPrivilege               types.String `tfsdk:"missing_any_privilege"`
	MissingPrivilege                  types.Set    `tfsdk:"missing_privilege"`
	RetiredPasswordUsedForBind        types.String `tfsdk:"retired_password_used_for_bind"`
	SearchEntryReturnedCriteria       types.String `tfsdk:"search_entry_returned_criteria"`
	SearchEntryReturnedCount          types.Int64  `tfsdk:"search_entry_returned_count"`
	SearchReferenceReturnedCriteria   types.String `tfsdk:"search_reference_returned_criteria"`
	SearchReferenceReturnedCount      types.Int64  `tfsdk:"search_reference_returned_count"`
	SearchIndexedCriteria             types.String `tfsdk:"search_indexed_criteria"`
	IncludedAuthzUserBaseDN           types.Set    `tfsdk:"included_authz_user_base_dn"`
	ExcludedAuthzUserBaseDN           types.Set    `tfsdk:"excluded_authz_user_base_dn"`
	AllIncludedAuthzUserGroupDN       types.Set    `tfsdk:"all_included_authz_user_group_dn"`
	AnyIncludedAuthzUserGroupDN       types.Set    `tfsdk:"any_included_authz_user_group_dn"`
	NotAllIncludedAuthzUserGroupDN    types.Set    `tfsdk:"not_all_included_authz_user_group_dn"`
	NoneIncludedAuthzUserGroupDN      types.Set    `tfsdk:"none_included_authz_user_group_dn"`
	IncludeAnonymousBinds             types.Bool   `tfsdk:"include_anonymous_binds"`
	IncludedUserBaseDN                types.Set    `tfsdk:"included_user_base_dn"`
	ExcludedUserBaseDN                types.Set    `tfsdk:"excluded_user_base_dn"`
	IncludedUserFilter                types.Set    `tfsdk:"included_user_filter"`
	ExcludedUserFilter                types.Set    `tfsdk:"excluded_user_filter"`
	IncludedUserGroupDN               types.Set    `tfsdk:"included_user_group_dn"`
	ExcludedUserGroupDN               types.Set    `tfsdk:"excluded_user_group_dn"`
	Description                       types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the datasource.
func (r *resultCriteriaDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Result Criteria.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Result Criteria resource. Options are ['successful-bind', 'simple', 'aggregate', 'replication-assurance', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Result Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Result Criteria. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"local_assurance_level": schema.SetAttribute{
				Description: "The local assurance level values that will be allowed to match this Replication Assurance Result Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"remote_assurance_level": schema.SetAttribute{
				Description: "The local assurance level values that will be allowed to match this Replication Assurance Result Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"assurance_timeout_criteria": schema.StringAttribute{
				Description: "The criteria to use when performing matching based on the assurance timeout.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"assurance_timeout_value": schema.StringAttribute{
				Description: "The value to use for performing matching based on the assurance timeout. This will be ignored if the assurance-timeout-criteria is \"any\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"response_delayed_by_assurance": schema.StringAttribute{
				Description: "Indicates whether this Replication Assurance Result Criteria should match operations based on whether the response to the client was delayed by assurance processing.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"assurance_behavior_altered_by_control": schema.StringAttribute{
				Description: "Indicates whether this Replication Assurance Result Criteria should match operations based on whether the assurance requirements were altered by a control included in the request from the client.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"assurance_satisfied": schema.StringAttribute{
				Description: "Indicates whether this Replication Assurance Result Criteria should match operations based on whether the assurance requirements have been satisfied.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"all_included_result_criteria": schema.SetAttribute{
				Description: "Specifies a result criteria object that must match the associated operation result in order to match the aggregate result criteria. If one or more all-included result criteria objects are provided, then an operation result must match all of them in order to match the aggregate result criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_result_criteria": schema.SetAttribute{
				Description: "Specifies a result criteria object that may match the associated operation result in order to match the aggregate result criteria. If one or more any-included result criteria objects are provided, then an operation result must match at least one of them in order to match the aggregate result criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_result_criteria": schema.SetAttribute{
				Description: "Specifies a result criteria object that should not match the associated operation result in order to match the aggregate result criteria. If one or more not-all-included result criteria objects are provided, then an operation result must not match all of them (that is, it may match zero or more of them, but it must not match all of them) in order to match the aggregate result criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_result_criteria": schema.SetAttribute{
				Description: "Specifies a result criteria object that must not match the associated operation result in order to match the aggregate result criteria. If one or more none-included result criteria objects are provided, then an operation result must not match any of them in order to match the aggregate result criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"request_criteria": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `successful-bind`: Specifies a request criteria object that must match the associated request for operations included in this Successful Bind Result Criteria. When the `type` attribute is set to `simple`: Specifies a request criteria object that must match the associated request for operations included in this Simple Result Criteria.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `successful-bind`: Specifies a request criteria object that must match the associated request for operations included in this Successful Bind Result Criteria.\n  - `simple`: Specifies a request criteria object that must match the associated request for operations included in this Simple Result Criteria.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"result_code_criteria": schema.StringAttribute{
				Description: "Specifies which operation result codes are allowed for operations included in this Simple Result Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"result_code_value": schema.SetAttribute{
				Description: "Specifies the operation result code values for results included in this Simple Result Criteria. This will only be taken into account if the \"result-code-criteria\" property has a value of \"selected-result-codes\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"processing_time_criteria": schema.StringAttribute{
				Description: "Indicates whether the time required to process the operation should be taken into consideration when determining whether to include the operation in this Simple Result Criteria. If the processing time should be taken into account, then the \"processing-time-value\" property should contain the boundary value.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"processing_time_value": schema.StringAttribute{
				Description: "Specifies the boundary value to use for the operation processing time when determining whether to include that operation in this Simple Result Criteria. This will be ignored if the \"processing-time-criteria\" property has a value of \"any\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"queue_time_criteria": schema.StringAttribute{
				Description: "Indicates whether the time the operation was required to wait on the work queue should be taken into consideration when determining whether to include the operation in this Simple Result Criteria. If the queue time should be taken into account, then the \"queue-time-value\" property should contain the boundary value. This property should only be given a value other than \"any\" if the work queue has been configured to monitor the time operations have spent on the work queue.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"queue_time_value": schema.StringAttribute{
				Description: "Specifies the boundary value to use for the time an operation spent on the work queue when determining whether to include that operation in this Simple Result Criteria. This will be ignored if the \"queue-time-criteria\" property has a value of \"any\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"referral_returned": schema.StringAttribute{
				Description: "Indicates whether operation results which include one or more referral URLs should be included in this Simple Result Criteria. If no value is provided, then whether an operation includes any referral URLs will not be considered when determining whether it matches this Simple Result Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"all_included_response_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must be present in the response to the client for operations included in this Simple Result Criteria. If any control OIDs are provided, then the response must contain all of those controls.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_response_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that may be present in the response to the client for operations included in this Simple Result Criteria. If any control OIDs are provided, then the response must contain at least one of those controls.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_response_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that should not be present in the response to the client for operations included in this Simple Result Criteria. If any control OIDs are provided, then the response must not contain at least one of those controls (that is, the response may contain zero or more of those controls, but not all of them).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_response_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must not be present in the response to the client for operations included in this Simple Result Criteria. If any control OIDs are provided, then the response must not contain any of those controls.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"used_alternate_authzid": schema.StringAttribute{
				Description: "Indicates whether operation results in which the associated operation used an authorization identity that is different from the authentication identity (e.g., as the result of using a proxied authorization control) should be included in this Simple Result Criteria. If no value is provided, then whether an operation used an alternate authorization identity will not be considered when determining whether it matches this Simple Result Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"used_any_privilege": schema.StringAttribute{
				Description: "Indicates whether operations in which one or more privileges were used should be included in this Simple Result Criteria. If no value is provided, then whether an operation used any privileges will not be considered when determining whether it matches this Simple Result Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"used_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that must have been used during the processing for operations included in this Simple Result Criteria. If any privilege names are provided, then the associated operation must have used at least one of those privileges. If no privilege names were provided, then the set of privileges used will not be considered when determining whether an operation should be included in this Simple Result Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"missing_any_privilege": schema.StringAttribute{
				Description: "Indicates whether operations in which one or more privileges were missing should be included in this Simple Result Criteria. If no value is provided, then whether there were any missing privileges will not be considered when determining whether an operation matches this Simple Result Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"missing_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that must have been missing during the processing for operations included in this Simple Result Criteria. If any privilege names are provided, then the associated operation must have been missing at least one of those privileges. If no privilege names were provided, then the set of privileges missing will not be considered when determining whether an operation should be included in this Simple Result Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"retired_password_used_for_bind": schema.StringAttribute{
				Description: "Indicates whether the use of a retired password for authentication should be considered when determining whether a bind operation should be included in this Simple Result Criteria. This will be ignored for all operations other than bind.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"search_entry_returned_criteria": schema.StringAttribute{
				Description: "Indicates whether the number of entries returned should be considered when determining whether a search operation should be included in this Simple Result Criteria. This will be ignored for all operations other than search.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"search_entry_returned_count": schema.Int64Attribute{
				Description: "Specifies the target number of entries returned for use when determining whether a search operation should be included in this Simple Result Criteria. This will be ignored for all operations other than search, and it will be ignored for search operations if the \"search-entry-criteria\" property has a value of \"any\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"search_reference_returned_criteria": schema.StringAttribute{
				Description: "Indicates whether the number of references returned should be considered when determining whether a search operation should be included in this Simple Result Criteria. This will be ignored for all operations other than search.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"search_reference_returned_count": schema.Int64Attribute{
				Description: "Specifies the target number of references returned for use when determining whether a search operation should be included in this Simple Result Criteria. This will be ignored for all operations other than search, and it will be ignored for search operations if the \"search-reference-criteria\" property has a value of \"any\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"search_indexed_criteria": schema.StringAttribute{
				Description: "Indicates whether a search operation should be matched by this Simple Result Criteria based on whether it is considered indexed by the server. This will be ignored for all operations other than search.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"included_authz_user_base_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which authorization user entries may exist for operations included in this Simple Result Criteria. The authorization user could be the currently authenticated user on the connection (the user that performed the Bind operation), or different if proxied authorization was used to request that the operation be performed under the authorization of another user (as is the case for operations that come through a Directory Proxy Server). This property will be ignored for operations where no authentication or authorization has been performed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_authz_user_base_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which authorization user entries may exist for operations excluded from this Simple Result Criteria. The authorization user could be the currently authenticated user on the connection (the user that performed the Bind operation), or different if proxied authorization was used to request that the operation be performed under the authorization of another user (as is the case for operations that come through a Directory Proxy Server). This property will be ignored for operations where no authentication or authorization has been performed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"all_included_authz_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authorization users must exist for operations included in this Simple Result Criteria. If any group DNs are provided, then the authorization user must be a member of all of those groups. The authorization user could be the currently authenticated user on the connection (the user that performed the Bind operation), or different if proxied authorization was used to request that the operation be performed under the authorization of another user (as is the case for operations that come through a Directory Proxy Server). This property will be ignored for operations where no authentication or authorization has been performed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_authz_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authorization users may exist for operations included in this Simple Result Criteria. If any group DNs are provided, then the authorization user must be a member of at least one of those groups. The authorization user could be the currently authenticated user on the connection (the user that performed the Bind operation), or different if proxied authorization was used to request that the operation be performed under the authorization of another user (as is the case for operations that come through a Directory Proxy Server). This property will be ignored for operations where no authentication or authorization has been performed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_authz_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authorization users should not exist for operations included in this Simple Result Criteria. If any group DNs are provided, then the authorization user must not be a member of at least one of those groups (that is, the user may be a member of zero or more of those groups, but not of all of them). The authorization user could be the currently authenticated user on the connection (the user that performed the Bind operation), or different if proxied authorization was used to request that the operation be performed under the authorization of another user (as is the case for operations that come through a Directory Proxy Server). This property will be ignored for operations where no authentication or authorization has been performed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_authz_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authorization users must not exist for operations included in this Simple Result Criteria. If any group DNs are provided, then the authorization user must not be a member any of those groups. The authorization user could be the currently authenticated user on the connection (the user that performed the Bind operation), or different if proxied authorization was used to request that the operation be performed under the authorization of another user (as is the case for operations that come through a Directory Proxy Server). This property will be ignored for operations where no authentication or authorization has been performed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"include_anonymous_binds": schema.BoolAttribute{
				Description: "Indicates whether this criteria will be permitted to match bind operations that resulted in anonymous authentication.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"included_user_base_dn": schema.SetAttribute{
				Description: "A set of base DNs for authenticated users that will be permitted to match this criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_user_base_dn": schema.SetAttribute{
				Description: "A set of base DNs for authenticated users that will not be permitted to match this criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_user_filter": schema.SetAttribute{
				Description: "A set of filters that may be used to identify entries for authenticated users that will be permitted to match this criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_user_filter": schema.SetAttribute{
				Description: "A set of filters that may be used to identify entries for authenticated users that will not be permitted to match this criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_user_group_dn": schema.SetAttribute{
				Description: "The DNs of the groups whose members will be permitted to match this criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_user_group_dn": schema.SetAttribute{
				Description: "The DNs of the groups whose members will not be permitted to match this criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Result Criteria",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a SuccessfulBindResultCriteriaResponse object into the model struct
func readSuccessfulBindResultCriteriaResponseDataSource(ctx context.Context, r *client.SuccessfulBindResultCriteriaResponse, state *resultCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("successful-bind")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.IncludeAnonymousBinds = internaltypes.BoolTypeOrNil(r.IncludeAnonymousBinds)
	state.IncludedUserBaseDN = internaltypes.GetStringSet(r.IncludedUserBaseDN)
	state.ExcludedUserBaseDN = internaltypes.GetStringSet(r.ExcludedUserBaseDN)
	state.IncludedUserFilter = internaltypes.GetStringSet(r.IncludedUserFilter)
	state.ExcludedUserFilter = internaltypes.GetStringSet(r.ExcludedUserFilter)
	state.IncludedUserGroupDN = internaltypes.GetStringSet(r.IncludedUserGroupDN)
	state.ExcludedUserGroupDN = internaltypes.GetStringSet(r.ExcludedUserGroupDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a SimpleResultCriteriaResponse object into the model struct
func readSimpleResultCriteriaResponseDataSource(ctx context.Context, r *client.SimpleResultCriteriaResponse, state *resultCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("simple")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCodeCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaResultCodeCriteriaProp(r.ResultCodeCriteria), false)
	state.ResultCodeValue = internaltypes.GetStringSet(
		client.StringSliceEnumresultCriteriaResultCodeValueProp(r.ResultCodeValue))
	state.ProcessingTimeCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaProcessingTimeCriteriaProp(r.ProcessingTimeCriteria), false)
	state.ProcessingTimeValue = internaltypes.StringTypeOrNil(r.ProcessingTimeValue, false)
	state.QueueTimeCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaQueueTimeCriteriaProp(r.QueueTimeCriteria), false)
	state.QueueTimeValue = internaltypes.StringTypeOrNil(r.QueueTimeValue, false)
	state.ReferralReturned = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaReferralReturnedProp(r.ReferralReturned), false)
	state.AllIncludedResponseControl = internaltypes.GetStringSet(r.AllIncludedResponseControl)
	state.AnyIncludedResponseControl = internaltypes.GetStringSet(r.AnyIncludedResponseControl)
	state.NotAllIncludedResponseControl = internaltypes.GetStringSet(r.NotAllIncludedResponseControl)
	state.NoneIncludedResponseControl = internaltypes.GetStringSet(r.NoneIncludedResponseControl)
	state.UsedAlternateAuthzid = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaUsedAlternateAuthzidProp(r.UsedAlternateAuthzid), false)
	state.UsedAnyPrivilege = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaUsedAnyPrivilegeProp(r.UsedAnyPrivilege), false)
	state.UsedPrivilege = internaltypes.GetStringSet(
		client.StringSliceEnumresultCriteriaUsedPrivilegeProp(r.UsedPrivilege))
	state.MissingAnyPrivilege = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaMissingAnyPrivilegeProp(r.MissingAnyPrivilege), false)
	state.MissingPrivilege = internaltypes.GetStringSet(
		client.StringSliceEnumresultCriteriaMissingPrivilegeProp(r.MissingPrivilege))
	state.RetiredPasswordUsedForBind = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaRetiredPasswordUsedForBindProp(r.RetiredPasswordUsedForBind), false)
	state.SearchEntryReturnedCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaSearchEntryReturnedCriteriaProp(r.SearchEntryReturnedCriteria), false)
	state.SearchEntryReturnedCount = internaltypes.Int64TypeOrNil(r.SearchEntryReturnedCount)
	state.SearchReferenceReturnedCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaSearchReferenceReturnedCriteriaProp(r.SearchReferenceReturnedCriteria), false)
	state.SearchReferenceReturnedCount = internaltypes.Int64TypeOrNil(r.SearchReferenceReturnedCount)
	state.SearchIndexedCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaSearchIndexedCriteriaProp(r.SearchIndexedCriteria), false)
	state.IncludedAuthzUserBaseDN = internaltypes.GetStringSet(r.IncludedAuthzUserBaseDN)
	state.ExcludedAuthzUserBaseDN = internaltypes.GetStringSet(r.ExcludedAuthzUserBaseDN)
	state.AllIncludedAuthzUserGroupDN = internaltypes.GetStringSet(r.AllIncludedAuthzUserGroupDN)
	state.AnyIncludedAuthzUserGroupDN = internaltypes.GetStringSet(r.AnyIncludedAuthzUserGroupDN)
	state.NotAllIncludedAuthzUserGroupDN = internaltypes.GetStringSet(r.NotAllIncludedAuthzUserGroupDN)
	state.NoneIncludedAuthzUserGroupDN = internaltypes.GetStringSet(r.NoneIncludedAuthzUserGroupDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a AggregateResultCriteriaResponse object into the model struct
func readAggregateResultCriteriaResponseDataSource(ctx context.Context, r *client.AggregateResultCriteriaResponse, state *resultCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("aggregate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllIncludedResultCriteria = internaltypes.GetStringSet(r.AllIncludedResultCriteria)
	state.AnyIncludedResultCriteria = internaltypes.GetStringSet(r.AnyIncludedResultCriteria)
	state.NotAllIncludedResultCriteria = internaltypes.GetStringSet(r.NotAllIncludedResultCriteria)
	state.NoneIncludedResultCriteria = internaltypes.GetStringSet(r.NoneIncludedResultCriteria)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a ReplicationAssuranceResultCriteriaResponse object into the model struct
func readReplicationAssuranceResultCriteriaResponseDataSource(ctx context.Context, r *client.ReplicationAssuranceResultCriteriaResponse, state *resultCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("replication-assurance")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LocalAssuranceLevel = internaltypes.GetStringSet(
		client.StringSliceEnumresultCriteriaLocalAssuranceLevelProp(r.LocalAssuranceLevel))
	state.RemoteAssuranceLevel = internaltypes.GetStringSet(
		client.StringSliceEnumresultCriteriaRemoteAssuranceLevelProp(r.RemoteAssuranceLevel))
	state.AssuranceTimeoutCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaAssuranceTimeoutCriteriaProp(r.AssuranceTimeoutCriteria), false)
	state.AssuranceTimeoutValue = internaltypes.StringTypeOrNil(r.AssuranceTimeoutValue, false)
	state.ResponseDelayedByAssurance = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaResponseDelayedByAssuranceProp(r.ResponseDelayedByAssurance), false)
	state.AssuranceBehaviorAlteredByControl = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaAssuranceBehaviorAlteredByControlProp(r.AssuranceBehaviorAlteredByControl), false)
	state.AssuranceSatisfied = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaAssuranceSatisfiedProp(r.AssuranceSatisfied), false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a ThirdPartyResultCriteriaResponse object into the model struct
func readThirdPartyResultCriteriaResponseDataSource(ctx context.Context, r *client.ThirdPartyResultCriteriaResponse, state *resultCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read resource information
func (r *resultCriteriaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state resultCriteriaDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ResultCriteriaAPI.GetResultCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Result Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.SuccessfulBindResultCriteriaResponse != nil {
		readSuccessfulBindResultCriteriaResponseDataSource(ctx, readResponse.SuccessfulBindResultCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SimpleResultCriteriaResponse != nil {
		readSimpleResultCriteriaResponseDataSource(ctx, readResponse.SimpleResultCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AggregateResultCriteriaResponse != nil {
		readAggregateResultCriteriaResponseDataSource(ctx, readResponse.AggregateResultCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ReplicationAssuranceResultCriteriaResponse != nil {
		readReplicationAssuranceResultCriteriaResponseDataSource(ctx, readResponse.ReplicationAssuranceResultCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyResultCriteriaResponse != nil {
		readThirdPartyResultCriteriaResponseDataSource(ctx, readResponse.ThirdPartyResultCriteriaResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
