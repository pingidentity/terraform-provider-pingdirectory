package resultcriteria

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	_ resource.Resource                = &resultCriteriaResource{}
	_ resource.ResourceWithConfigure   = &resultCriteriaResource{}
	_ resource.ResourceWithImportState = &resultCriteriaResource{}
	_ resource.Resource                = &defaultResultCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultResultCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultResultCriteriaResource{}
)

// Create a Result Criteria resource
func NewResultCriteriaResource() resource.Resource {
	return &resultCriteriaResource{}
}

func NewDefaultResultCriteriaResource() resource.Resource {
	return &defaultResultCriteriaResource{}
}

// resultCriteriaResource is the resource implementation.
type resultCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultResultCriteriaResource is the resource implementation.
type defaultResultCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *resultCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_result_criteria"
}

func (r *defaultResultCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_result_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *resultCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultResultCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type resultCriteriaResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	Notifications                     types.Set    `tfsdk:"notifications"`
	RequiredActions                   types.Set    `tfsdk:"required_actions"`
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

// GetSchema defines the schema for the resource.
func (r *resultCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resultCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultResultCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resultCriteriaSchema(ctx, req, resp, true)
}

func resultCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Result Criteria.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Result Criteria resource. Options are ['successful-bind', 'simple', 'aggregate', 'replication-assurance', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"successful-bind", "simple", "aggregate", "replication-assurance", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Result Criteria.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Result Criteria. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"local_assurance_level": schema.SetAttribute{
				Description: "The local assurance level values that will be allowed to match this Replication Assurance Result Criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"remote_assurance_level": schema.SetAttribute{
				Description: "The local assurance level values that will be allowed to match this Replication Assurance Result Criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"assurance_timeout_criteria": schema.StringAttribute{
				Description: "The criteria to use when performing matching based on the assurance timeout.",
				Optional:    true,
				Computed:    true,
			},
			"assurance_timeout_value": schema.StringAttribute{
				Description: "The value to use for performing matching based on the assurance timeout. This will be ignored if the assurance-timeout-criteria is \"any\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"response_delayed_by_assurance": schema.StringAttribute{
				Description: "Indicates whether this Replication Assurance Result Criteria should match operations based on whether the response to the client was delayed by assurance processing.",
				Optional:    true,
				Computed:    true,
			},
			"assurance_behavior_altered_by_control": schema.StringAttribute{
				Description: "Indicates whether this Replication Assurance Result Criteria should match operations based on whether the assurance requirements were altered by a control included in the request from the client.",
				Optional:    true,
				Computed:    true,
			},
			"assurance_satisfied": schema.StringAttribute{
				Description: "Indicates whether this Replication Assurance Result Criteria should match operations based on whether the assurance requirements have been satisfied.",
				Optional:    true,
				Computed:    true,
			},
			"all_included_result_criteria": schema.SetAttribute{
				Description: "Specifies a result criteria object that must match the associated operation result in order to match the aggregate result criteria. If one or more all-included result criteria objects are provided, then an operation result must match all of them in order to match the aggregate result criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"any_included_result_criteria": schema.SetAttribute{
				Description: "Specifies a result criteria object that may match the associated operation result in order to match the aggregate result criteria. If one or more any-included result criteria objects are provided, then an operation result must match at least one of them in order to match the aggregate result criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"not_all_included_result_criteria": schema.SetAttribute{
				Description: "Specifies a result criteria object that should not match the associated operation result in order to match the aggregate result criteria. If one or more not-all-included result criteria objects are provided, then an operation result must not match all of them (that is, it may match zero or more of them, but it must not match all of them) in order to match the aggregate result criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"none_included_result_criteria": schema.SetAttribute{
				Description: "Specifies a result criteria object that must not match the associated operation result in order to match the aggregate result criteria. If one or more none-included result criteria objects are provided, then an operation result must not match any of them in order to match the aggregate result criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"request_criteria": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `successful-bind`: Specifies a request criteria object that must match the associated request for operations included in this Successful Bind Result Criteria. When the `type` attribute is set to `simple`: Specifies a request criteria object that must match the associated request for operations included in this Simple Result Criteria.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `successful-bind`: Specifies a request criteria object that must match the associated request for operations included in this Successful Bind Result Criteria.\n  - `simple`: Specifies a request criteria object that must match the associated request for operations included in this Simple Result Criteria.",
				Optional:            true,
			},
			"result_code_criteria": schema.StringAttribute{
				Description: "Specifies which operation result codes are allowed for operations included in this Simple Result Criteria.",
				Optional:    true,
				Computed:    true,
			},
			"result_code_value": schema.SetAttribute{
				Description: "Specifies the operation result code values for results included in this Simple Result Criteria. This will only be taken into account if the \"result-code-criteria\" property has a value of \"selected-result-codes\".",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"processing_time_criteria": schema.StringAttribute{
				Description: "Indicates whether the time required to process the operation should be taken into consideration when determining whether to include the operation in this Simple Result Criteria. If the processing time should be taken into account, then the \"processing-time-value\" property should contain the boundary value.",
				Optional:    true,
				Computed:    true,
			},
			"processing_time_value": schema.StringAttribute{
				Description: "Specifies the boundary value to use for the operation processing time when determining whether to include that operation in this Simple Result Criteria. This will be ignored if the \"processing-time-criteria\" property has a value of \"any\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"queue_time_criteria": schema.StringAttribute{
				Description: "Indicates whether the time the operation was required to wait on the work queue should be taken into consideration when determining whether to include the operation in this Simple Result Criteria. If the queue time should be taken into account, then the \"queue-time-value\" property should contain the boundary value. This property should only be given a value other than \"any\" if the work queue has been configured to monitor the time operations have spent on the work queue.",
				Optional:    true,
				Computed:    true,
			},
			"queue_time_value": schema.StringAttribute{
				Description: "Specifies the boundary value to use for the time an operation spent on the work queue when determining whether to include that operation in this Simple Result Criteria. This will be ignored if the \"queue-time-criteria\" property has a value of \"any\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"referral_returned": schema.StringAttribute{
				Description: "Indicates whether operation results which include one or more referral URLs should be included in this Simple Result Criteria. If no value is provided, then whether an operation includes any referral URLs will not be considered when determining whether it matches this Simple Result Criteria.",
				Optional:    true,
				Computed:    true,
			},
			"all_included_response_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must be present in the response to the client for operations included in this Simple Result Criteria. If any control OIDs are provided, then the response must contain all of those controls.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"any_included_response_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that may be present in the response to the client for operations included in this Simple Result Criteria. If any control OIDs are provided, then the response must contain at least one of those controls.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"not_all_included_response_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that should not be present in the response to the client for operations included in this Simple Result Criteria. If any control OIDs are provided, then the response must not contain at least one of those controls (that is, the response may contain zero or more of those controls, but not all of them).",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"none_included_response_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must not be present in the response to the client for operations included in this Simple Result Criteria. If any control OIDs are provided, then the response must not contain any of those controls.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"used_alternate_authzid": schema.StringAttribute{
				Description: "Indicates whether operation results in which the associated operation used an authorization identity that is different from the authentication identity (e.g., as the result of using a proxied authorization control) should be included in this Simple Result Criteria. If no value is provided, then whether an operation used an alternate authorization identity will not be considered when determining whether it matches this Simple Result Criteria.",
				Optional:    true,
				Computed:    true,
			},
			"used_any_privilege": schema.StringAttribute{
				Description: "Indicates whether operations in which one or more privileges were used should be included in this Simple Result Criteria. If no value is provided, then whether an operation used any privileges will not be considered when determining whether it matches this Simple Result Criteria.",
				Optional:    true,
				Computed:    true,
			},
			"used_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that must have been used during the processing for operations included in this Simple Result Criteria. If any privilege names are provided, then the associated operation must have used at least one of those privileges. If no privilege names were provided, then the set of privileges used will not be considered when determining whether an operation should be included in this Simple Result Criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"missing_any_privilege": schema.StringAttribute{
				Description: "Indicates whether operations in which one or more privileges were missing should be included in this Simple Result Criteria. If no value is provided, then whether there were any missing privileges will not be considered when determining whether an operation matches this Simple Result Criteria.",
				Optional:    true,
				Computed:    true,
			},
			"missing_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that must have been missing during the processing for operations included in this Simple Result Criteria. If any privilege names are provided, then the associated operation must have been missing at least one of those privileges. If no privilege names were provided, then the set of privileges missing will not be considered when determining whether an operation should be included in this Simple Result Criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"retired_password_used_for_bind": schema.StringAttribute{
				Description: "Indicates whether the use of a retired password for authentication should be considered when determining whether a bind operation should be included in this Simple Result Criteria. This will be ignored for all operations other than bind.",
				Optional:    true,
				Computed:    true,
			},
			"search_entry_returned_criteria": schema.StringAttribute{
				Description: "Indicates whether the number of entries returned should be considered when determining whether a search operation should be included in this Simple Result Criteria. This will be ignored for all operations other than search.",
				Optional:    true,
				Computed:    true,
			},
			"search_entry_returned_count": schema.Int64Attribute{
				Description: "Specifies the target number of entries returned for use when determining whether a search operation should be included in this Simple Result Criteria. This will be ignored for all operations other than search, and it will be ignored for search operations if the \"search-entry-criteria\" property has a value of \"any\".",
				Optional:    true,
				Computed:    true,
			},
			"search_reference_returned_criteria": schema.StringAttribute{
				Description: "Indicates whether the number of references returned should be considered when determining whether a search operation should be included in this Simple Result Criteria. This will be ignored for all operations other than search.",
				Optional:    true,
				Computed:    true,
			},
			"search_reference_returned_count": schema.Int64Attribute{
				Description: "Specifies the target number of references returned for use when determining whether a search operation should be included in this Simple Result Criteria. This will be ignored for all operations other than search, and it will be ignored for search operations if the \"search-reference-criteria\" property has a value of \"any\".",
				Optional:    true,
				Computed:    true,
			},
			"search_indexed_criteria": schema.StringAttribute{
				Description: "Indicates whether a search operation should be matched by this Simple Result Criteria based on whether it is considered indexed by the server. This will be ignored for all operations other than search.",
				Optional:    true,
				Computed:    true,
			},
			"included_authz_user_base_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which authorization user entries may exist for operations included in this Simple Result Criteria. The authorization user could be the currently authenticated user on the connection (the user that performed the Bind operation), or different if proxied authorization was used to request that the operation be performed under the authorization of another user (as is the case for operations that come through a Directory Proxy Server). This property will be ignored for operations where no authentication or authorization has been performed.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"excluded_authz_user_base_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which authorization user entries may exist for operations excluded from this Simple Result Criteria. The authorization user could be the currently authenticated user on the connection (the user that performed the Bind operation), or different if proxied authorization was used to request that the operation be performed under the authorization of another user (as is the case for operations that come through a Directory Proxy Server). This property will be ignored for operations where no authentication or authorization has been performed.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"all_included_authz_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authorization users must exist for operations included in this Simple Result Criteria. If any group DNs are provided, then the authorization user must be a member of all of those groups. The authorization user could be the currently authenticated user on the connection (the user that performed the Bind operation), or different if proxied authorization was used to request that the operation be performed under the authorization of another user (as is the case for operations that come through a Directory Proxy Server). This property will be ignored for operations where no authentication or authorization has been performed.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"any_included_authz_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authorization users may exist for operations included in this Simple Result Criteria. If any group DNs are provided, then the authorization user must be a member of at least one of those groups. The authorization user could be the currently authenticated user on the connection (the user that performed the Bind operation), or different if proxied authorization was used to request that the operation be performed under the authorization of another user (as is the case for operations that come through a Directory Proxy Server). This property will be ignored for operations where no authentication or authorization has been performed.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"not_all_included_authz_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authorization users should not exist for operations included in this Simple Result Criteria. If any group DNs are provided, then the authorization user must not be a member of at least one of those groups (that is, the user may be a member of zero or more of those groups, but not of all of them). The authorization user could be the currently authenticated user on the connection (the user that performed the Bind operation), or different if proxied authorization was used to request that the operation be performed under the authorization of another user (as is the case for operations that come through a Directory Proxy Server). This property will be ignored for operations where no authentication or authorization has been performed.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"none_included_authz_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authorization users must not exist for operations included in this Simple Result Criteria. If any group DNs are provided, then the authorization user must not be a member any of those groups. The authorization user could be the currently authenticated user on the connection (the user that performed the Bind operation), or different if proxied authorization was used to request that the operation be performed under the authorization of another user (as is the case for operations that come through a Directory Proxy Server). This property will be ignored for operations where no authentication or authorization has been performed.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"include_anonymous_binds": schema.BoolAttribute{
				Description: "Indicates whether this criteria will be permitted to match bind operations that resulted in anonymous authentication.",
				Optional:    true,
				Computed:    true,
			},
			"included_user_base_dn": schema.SetAttribute{
				Description: "A set of base DNs for authenticated users that will be permitted to match this criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"excluded_user_base_dn": schema.SetAttribute{
				Description: "A set of base DNs for authenticated users that will not be permitted to match this criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"included_user_filter": schema.SetAttribute{
				Description: "A set of filters that may be used to identify entries for authenticated users that will be permitted to match this criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"excluded_user_filter": schema.SetAttribute{
				Description: "A set of filters that may be used to identify entries for authenticated users that will not be permitted to match this criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"included_user_group_dn": schema.SetAttribute{
				Description: "The DNs of the groups whose members will be permitted to match this criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"excluded_user_group_dn": schema.SetAttribute{
				Description: "The DNs of the groups whose members will not be permitted to match this criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Result Criteria",
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
	} else {
		// Add RequiresReplace modifier for read-only attributes
		extensionClassAttr := schemaDef.Attributes["extension_class"].(schema.StringAttribute)
		extensionClassAttr.PlanModifiers = append(extensionClassAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["extension_class"] = extensionClassAttr
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *resultCriteriaResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanResultCriteria(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_result_criteria")
	var planModel, configModel resultCriteriaResourceModel
	req.Config.Get(ctx, &configModel)
	req.Plan.Get(ctx, &planModel)
	resourceType := planModel.Type.ValueString()
	anyDefaultsSet := false
	// Set defaults for successful-bind type
	if resourceType == "successful-bind" {
		if !internaltypes.IsDefined(configModel.IncludeAnonymousBinds) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeAnonymousBinds.Equal(defaultVal) {
				planModel.IncludeAnonymousBinds = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for simple type
	if resourceType == "simple" {
		if !internaltypes.IsDefined(configModel.ResultCodeCriteria) {
			defaultVal := types.StringValue("all-result-codes")
			if !planModel.ResultCodeCriteria.Equal(defaultVal) {
				planModel.ResultCodeCriteria = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ProcessingTimeCriteria) {
			defaultVal := types.StringValue("any")
			if !planModel.ProcessingTimeCriteria.Equal(defaultVal) {
				planModel.ProcessingTimeCriteria = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueTimeCriteria) {
			defaultVal := types.StringValue("any")
			if !planModel.QueueTimeCriteria.Equal(defaultVal) {
				planModel.QueueTimeCriteria = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ReferralReturned) {
			defaultVal := types.StringValue("optional")
			if !planModel.ReferralReturned.Equal(defaultVal) {
				planModel.ReferralReturned = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.UsedAlternateAuthzid) {
			defaultVal := types.StringValue("optional")
			if !planModel.UsedAlternateAuthzid.Equal(defaultVal) {
				planModel.UsedAlternateAuthzid = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.UsedAnyPrivilege) {
			defaultVal := types.StringValue("optional")
			if !planModel.UsedAnyPrivilege.Equal(defaultVal) {
				planModel.UsedAnyPrivilege = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MissingAnyPrivilege) {
			defaultVal := types.StringValue("optional")
			if !planModel.MissingAnyPrivilege.Equal(defaultVal) {
				planModel.MissingAnyPrivilege = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.RetiredPasswordUsedForBind) {
			defaultVal := types.StringValue("any")
			if !planModel.RetiredPasswordUsedForBind.Equal(defaultVal) {
				planModel.RetiredPasswordUsedForBind = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SearchEntryReturnedCriteria) {
			defaultVal := types.StringValue("any")
			if !planModel.SearchEntryReturnedCriteria.Equal(defaultVal) {
				planModel.SearchEntryReturnedCriteria = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SearchEntryReturnedCount) {
			defaultVal := types.Int64Value(0)
			if !planModel.SearchEntryReturnedCount.Equal(defaultVal) {
				planModel.SearchEntryReturnedCount = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SearchReferenceReturnedCriteria) {
			defaultVal := types.StringValue("any")
			if !planModel.SearchReferenceReturnedCriteria.Equal(defaultVal) {
				planModel.SearchReferenceReturnedCriteria = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SearchReferenceReturnedCount) {
			defaultVal := types.Int64Value(0)
			if !planModel.SearchReferenceReturnedCount.Equal(defaultVal) {
				planModel.SearchReferenceReturnedCount = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SearchIndexedCriteria) {
			defaultVal := types.StringValue("any")
			if !planModel.SearchIndexedCriteria.Equal(defaultVal) {
				planModel.SearchIndexedCriteria = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for replication-assurance type
	if resourceType == "replication-assurance" {
		if !internaltypes.IsDefined(configModel.LocalAssuranceLevel) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("none"), types.StringValue("received-any-server"), types.StringValue("processed-all-servers")})
			if !planModel.LocalAssuranceLevel.Equal(defaultVal) {
				planModel.LocalAssuranceLevel = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.RemoteAssuranceLevel) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("none"), types.StringValue("received-any-remote-location"), types.StringValue("received-all-remote-locations"), types.StringValue("processed-all-remote-servers")})
			if !planModel.RemoteAssuranceLevel.Equal(defaultVal) {
				planModel.RemoteAssuranceLevel = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AssuranceTimeoutCriteria) {
			defaultVal := types.StringValue("any")
			if !planModel.AssuranceTimeoutCriteria.Equal(defaultVal) {
				planModel.AssuranceTimeoutCriteria = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ResponseDelayedByAssurance) {
			defaultVal := types.StringValue("any")
			if !planModel.ResponseDelayedByAssurance.Equal(defaultVal) {
				planModel.ResponseDelayedByAssurance = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AssuranceBehaviorAlteredByControl) {
			defaultVal := types.StringValue("any")
			if !planModel.AssuranceBehaviorAlteredByControl.Equal(defaultVal) {
				planModel.AssuranceBehaviorAlteredByControl = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AssuranceSatisfied) {
			defaultVal := types.StringValue("any")
			if !planModel.AssuranceSatisfied.Equal(defaultVal) {
				planModel.AssuranceSatisfied = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	if anyDefaultsSet {
		planModel.Notifications = types.SetUnknown(types.StringType)
		planModel.RequiredActions = types.SetUnknown(config.GetRequiredActionsObjectType())
	}
	planModel.setNotApplicableAttrsNull()
	resp.Plan.Set(ctx, &planModel)
}

func (r *defaultResultCriteriaResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanResultCriteria(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_default_result_criteria")
}

func modifyPlanResultCriteria(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, resourceName string) {
	compare, err := version.Compare(providerConfig.ProductVersion, version.PingDirectory9300)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model resultCriteriaResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.Type) && model.Type.ValueString() == "successful-bind" {
		version.CheckResourceSupported(&resp.Diagnostics, version.PingDirectory9300,
			providerConfig.ProductVersion, resourceName+" with type \"successful_bind\"")
	}
}

func (model *resultCriteriaResourceModel) setNotApplicableAttrsNull() {
	resourceType := model.Type.ValueString()
	// Set any not applicable computed attributes to null for each type
	if resourceType == "successful-bind" {
		model.SearchEntryReturnedCriteria = types.StringNull()
		model.AssuranceBehaviorAlteredByControl = types.StringNull()
		model.SearchReferenceReturnedCriteria = types.StringNull()
		model.ReferralReturned = types.StringNull()
		model.AssuranceSatisfied = types.StringNull()
		model.UsedAnyPrivilege = types.StringNull()
		model.ProcessingTimeValue = types.StringNull()
		model.QueueTimeCriteria = types.StringNull()
		model.MissingAnyPrivilege = types.StringNull()
		model.ResultCodeCriteria = types.StringNull()
		model.ProcessingTimeCriteria = types.StringNull()
		model.SearchIndexedCriteria = types.StringNull()
		model.UsedAlternateAuthzid = types.StringNull()
		model.AssuranceTimeoutCriteria = types.StringNull()
		model.LocalAssuranceLevel, _ = types.SetValue(types.StringType, []attr.Value{})
		model.ResponseDelayedByAssurance = types.StringNull()
		model.QueueTimeValue = types.StringNull()
		model.RemoteAssuranceLevel, _ = types.SetValue(types.StringType, []attr.Value{})
		model.RetiredPasswordUsedForBind = types.StringNull()
		model.AssuranceTimeoutValue = types.StringNull()
		model.SearchEntryReturnedCount = types.Int64Null()
		model.SearchReferenceReturnedCount = types.Int64Null()
	}
	if resourceType == "simple" {
		model.AssuranceBehaviorAlteredByControl = types.StringNull()
		model.AssuranceSatisfied = types.StringNull()
		model.IncludeAnonymousBinds = types.BoolNull()
		model.AssuranceTimeoutCriteria = types.StringNull()
		model.LocalAssuranceLevel, _ = types.SetValue(types.StringType, []attr.Value{})
		model.ResponseDelayedByAssurance = types.StringNull()
		model.RemoteAssuranceLevel, _ = types.SetValue(types.StringType, []attr.Value{})
		model.AssuranceTimeoutValue = types.StringNull()
	}
	if resourceType == "aggregate" {
		model.SearchEntryReturnedCriteria = types.StringNull()
		model.AssuranceBehaviorAlteredByControl = types.StringNull()
		model.SearchReferenceReturnedCriteria = types.StringNull()
		model.ReferralReturned = types.StringNull()
		model.AssuranceSatisfied = types.StringNull()
		model.UsedAnyPrivilege = types.StringNull()
		model.ProcessingTimeValue = types.StringNull()
		model.QueueTimeCriteria = types.StringNull()
		model.MissingAnyPrivilege = types.StringNull()
		model.IncludeAnonymousBinds = types.BoolNull()
		model.ResultCodeCriteria = types.StringNull()
		model.ProcessingTimeCriteria = types.StringNull()
		model.SearchIndexedCriteria = types.StringNull()
		model.UsedAlternateAuthzid = types.StringNull()
		model.AssuranceTimeoutCriteria = types.StringNull()
		model.LocalAssuranceLevel, _ = types.SetValue(types.StringType, []attr.Value{})
		model.ResponseDelayedByAssurance = types.StringNull()
		model.QueueTimeValue = types.StringNull()
		model.RemoteAssuranceLevel, _ = types.SetValue(types.StringType, []attr.Value{})
		model.RetiredPasswordUsedForBind = types.StringNull()
		model.AssuranceTimeoutValue = types.StringNull()
		model.SearchEntryReturnedCount = types.Int64Null()
		model.SearchReferenceReturnedCount = types.Int64Null()
	}
	if resourceType == "replication-assurance" {
		model.SearchEntryReturnedCriteria = types.StringNull()
		model.SearchReferenceReturnedCriteria = types.StringNull()
		model.ReferralReturned = types.StringNull()
		model.UsedAnyPrivilege = types.StringNull()
		model.ProcessingTimeValue = types.StringNull()
		model.QueueTimeCriteria = types.StringNull()
		model.MissingAnyPrivilege = types.StringNull()
		model.IncludeAnonymousBinds = types.BoolNull()
		model.ResultCodeCriteria = types.StringNull()
		model.ProcessingTimeCriteria = types.StringNull()
		model.SearchIndexedCriteria = types.StringNull()
		model.UsedAlternateAuthzid = types.StringNull()
		model.QueueTimeValue = types.StringNull()
		model.RetiredPasswordUsedForBind = types.StringNull()
		model.SearchEntryReturnedCount = types.Int64Null()
		model.SearchReferenceReturnedCount = types.Int64Null()
	}
	if resourceType == "third-party" {
		model.SearchEntryReturnedCriteria = types.StringNull()
		model.AssuranceBehaviorAlteredByControl = types.StringNull()
		model.SearchReferenceReturnedCriteria = types.StringNull()
		model.ReferralReturned = types.StringNull()
		model.AssuranceSatisfied = types.StringNull()
		model.UsedAnyPrivilege = types.StringNull()
		model.ProcessingTimeValue = types.StringNull()
		model.QueueTimeCriteria = types.StringNull()
		model.MissingAnyPrivilege = types.StringNull()
		model.IncludeAnonymousBinds = types.BoolNull()
		model.ResultCodeCriteria = types.StringNull()
		model.ProcessingTimeCriteria = types.StringNull()
		model.SearchIndexedCriteria = types.StringNull()
		model.UsedAlternateAuthzid = types.StringNull()
		model.AssuranceTimeoutCriteria = types.StringNull()
		model.LocalAssuranceLevel, _ = types.SetValue(types.StringType, []attr.Value{})
		model.ResponseDelayedByAssurance = types.StringNull()
		model.QueueTimeValue = types.StringNull()
		model.RemoteAssuranceLevel, _ = types.SetValue(types.StringType, []attr.Value{})
		model.RetiredPasswordUsedForBind = types.StringNull()
		model.AssuranceTimeoutValue = types.StringNull()
		model.SearchEntryReturnedCount = types.Int64Null()
		model.SearchReferenceReturnedCount = types.Int64Null()
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsResultCriteria() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("request_criteria"),
			path.MatchRoot("type"),
			[]string{"successful-bind", "simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_anonymous_binds"),
			path.MatchRoot("type"),
			[]string{"successful-bind"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("included_user_base_dn"),
			path.MatchRoot("type"),
			[]string{"successful-bind"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("excluded_user_base_dn"),
			path.MatchRoot("type"),
			[]string{"successful-bind"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("included_user_filter"),
			path.MatchRoot("type"),
			[]string{"successful-bind"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("excluded_user_filter"),
			path.MatchRoot("type"),
			[]string{"successful-bind"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("included_user_group_dn"),
			path.MatchRoot("type"),
			[]string{"successful-bind"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("excluded_user_group_dn"),
			path.MatchRoot("type"),
			[]string{"successful-bind"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("result_code_criteria"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("result_code_value"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("processing_time_criteria"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("processing_time_value"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("queue_time_criteria"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("queue_time_value"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("referral_returned"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("all_included_response_control"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("any_included_response_control"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("not_all_included_response_control"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("none_included_response_control"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("used_alternate_authzid"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("used_any_privilege"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("used_privilege"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("missing_any_privilege"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("missing_privilege"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("retired_password_used_for_bind"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("search_entry_returned_criteria"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("search_entry_returned_count"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("search_reference_returned_criteria"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("search_reference_returned_count"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("search_indexed_criteria"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("included_authz_user_base_dn"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("excluded_authz_user_base_dn"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("all_included_authz_user_group_dn"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("any_included_authz_user_group_dn"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("not_all_included_authz_user_group_dn"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("none_included_authz_user_group_dn"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("all_included_result_criteria"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("any_included_result_criteria"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("not_all_included_result_criteria"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("none_included_result_criteria"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("local_assurance_level"),
			path.MatchRoot("type"),
			[]string{"replication-assurance"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("remote_assurance_level"),
			path.MatchRoot("type"),
			[]string{"replication-assurance"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("assurance_timeout_criteria"),
			path.MatchRoot("type"),
			[]string{"replication-assurance"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("assurance_timeout_value"),
			path.MatchRoot("type"),
			[]string{"replication-assurance"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("response_delayed_by_assurance"),
			path.MatchRoot("type"),
			[]string{"replication-assurance"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("assurance_behavior_altered_by_control"),
			path.MatchRoot("type"),
			[]string{"replication-assurance"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("assurance_satisfied"),
			path.MatchRoot("type"),
			[]string{"replication-assurance"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_class"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_argument"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party",
			[]path.Expression{path.MatchRoot("extension_class")},
		),
	}
}

// Add config validators
func (r resultCriteriaResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsResultCriteria()
}

// Add config validators
func (r defaultResultCriteriaResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsResultCriteria()
}

// Add optional fields to create request for successful-bind result-criteria
func addOptionalSuccessfulBindResultCriteriaFields(ctx context.Context, addRequest *client.AddSuccessfulBindResultCriteriaRequest, plan resultCriteriaResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAnonymousBinds) {
		addRequest.IncludeAnonymousBinds = plan.IncludeAnonymousBinds.ValueBoolPointer()
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
	if internaltypes.IsDefined(plan.IncludedUserFilter) {
		var slice []string
		plan.IncludedUserFilter.ElementsAs(ctx, &slice, false)
		addRequest.IncludedUserFilter = slice
	}
	if internaltypes.IsDefined(plan.ExcludedUserFilter) {
		var slice []string
		plan.ExcludedUserFilter.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedUserFilter = slice
	}
	if internaltypes.IsDefined(plan.IncludedUserGroupDN) {
		var slice []string
		plan.IncludedUserGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.IncludedUserGroupDN = slice
	}
	if internaltypes.IsDefined(plan.ExcludedUserGroupDN) {
		var slice []string
		plan.ExcludedUserGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedUserGroupDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for simple result-criteria
func addOptionalSimpleResultCriteriaFields(ctx context.Context, addRequest *client.AddSimpleResultCriteriaRequest, plan resultCriteriaResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ResultCodeCriteria) {
		resultCodeCriteria, err := client.NewEnumresultCriteriaResultCodeCriteriaPropFromValue(plan.ResultCodeCriteria.ValueString())
		if err != nil {
			return err
		}
		addRequest.ResultCodeCriteria = resultCodeCriteria
	}
	if internaltypes.IsDefined(plan.ResultCodeValue) {
		var slice []string
		plan.ResultCodeValue.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumresultCriteriaResultCodeValueProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumresultCriteriaResultCodeValuePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.ResultCodeValue = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ProcessingTimeCriteria) {
		processingTimeCriteria, err := client.NewEnumresultCriteriaProcessingTimeCriteriaPropFromValue(plan.ProcessingTimeCriteria.ValueString())
		if err != nil {
			return err
		}
		addRequest.ProcessingTimeCriteria = processingTimeCriteria
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ProcessingTimeValue) {
		addRequest.ProcessingTimeValue = plan.ProcessingTimeValue.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.QueueTimeCriteria) {
		queueTimeCriteria, err := client.NewEnumresultCriteriaQueueTimeCriteriaPropFromValue(plan.QueueTimeCriteria.ValueString())
		if err != nil {
			return err
		}
		addRequest.QueueTimeCriteria = queueTimeCriteria
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.QueueTimeValue) {
		addRequest.QueueTimeValue = plan.QueueTimeValue.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReferralReturned) {
		referralReturned, err := client.NewEnumresultCriteriaReferralReturnedPropFromValue(plan.ReferralReturned.ValueString())
		if err != nil {
			return err
		}
		addRequest.ReferralReturned = referralReturned
	}
	if internaltypes.IsDefined(plan.AllIncludedResponseControl) {
		var slice []string
		plan.AllIncludedResponseControl.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedResponseControl = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedResponseControl) {
		var slice []string
		plan.AnyIncludedResponseControl.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedResponseControl = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedResponseControl) {
		var slice []string
		plan.NotAllIncludedResponseControl.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedResponseControl = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedResponseControl) {
		var slice []string
		plan.NoneIncludedResponseControl.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedResponseControl = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UsedAlternateAuthzid) {
		usedAlternateAuthzid, err := client.NewEnumresultCriteriaUsedAlternateAuthzidPropFromValue(plan.UsedAlternateAuthzid.ValueString())
		if err != nil {
			return err
		}
		addRequest.UsedAlternateAuthzid = usedAlternateAuthzid
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UsedAnyPrivilege) {
		usedAnyPrivilege, err := client.NewEnumresultCriteriaUsedAnyPrivilegePropFromValue(plan.UsedAnyPrivilege.ValueString())
		if err != nil {
			return err
		}
		addRequest.UsedAnyPrivilege = usedAnyPrivilege
	}
	if internaltypes.IsDefined(plan.UsedPrivilege) {
		var slice []string
		plan.UsedPrivilege.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumresultCriteriaUsedPrivilegeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumresultCriteriaUsedPrivilegePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.UsedPrivilege = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MissingAnyPrivilege) {
		missingAnyPrivilege, err := client.NewEnumresultCriteriaMissingAnyPrivilegePropFromValue(plan.MissingAnyPrivilege.ValueString())
		if err != nil {
			return err
		}
		addRequest.MissingAnyPrivilege = missingAnyPrivilege
	}
	if internaltypes.IsDefined(plan.MissingPrivilege) {
		var slice []string
		plan.MissingPrivilege.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumresultCriteriaMissingPrivilegeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumresultCriteriaMissingPrivilegePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.MissingPrivilege = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RetiredPasswordUsedForBind) {
		retiredPasswordUsedForBind, err := client.NewEnumresultCriteriaRetiredPasswordUsedForBindPropFromValue(plan.RetiredPasswordUsedForBind.ValueString())
		if err != nil {
			return err
		}
		addRequest.RetiredPasswordUsedForBind = retiredPasswordUsedForBind
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchEntryReturnedCriteria) {
		searchEntryReturnedCriteria, err := client.NewEnumresultCriteriaSearchEntryReturnedCriteriaPropFromValue(plan.SearchEntryReturnedCriteria.ValueString())
		if err != nil {
			return err
		}
		addRequest.SearchEntryReturnedCriteria = searchEntryReturnedCriteria
	}
	if internaltypes.IsDefined(plan.SearchEntryReturnedCount) {
		addRequest.SearchEntryReturnedCount = plan.SearchEntryReturnedCount.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchReferenceReturnedCriteria) {
		searchReferenceReturnedCriteria, err := client.NewEnumresultCriteriaSearchReferenceReturnedCriteriaPropFromValue(plan.SearchReferenceReturnedCriteria.ValueString())
		if err != nil {
			return err
		}
		addRequest.SearchReferenceReturnedCriteria = searchReferenceReturnedCriteria
	}
	if internaltypes.IsDefined(plan.SearchReferenceReturnedCount) {
		addRequest.SearchReferenceReturnedCount = plan.SearchReferenceReturnedCount.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchIndexedCriteria) {
		searchIndexedCriteria, err := client.NewEnumresultCriteriaSearchIndexedCriteriaPropFromValue(plan.SearchIndexedCriteria.ValueString())
		if err != nil {
			return err
		}
		addRequest.SearchIndexedCriteria = searchIndexedCriteria
	}
	if internaltypes.IsDefined(plan.IncludedAuthzUserBaseDN) {
		var slice []string
		plan.IncludedAuthzUserBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.IncludedAuthzUserBaseDN = slice
	}
	if internaltypes.IsDefined(plan.ExcludedAuthzUserBaseDN) {
		var slice []string
		plan.ExcludedAuthzUserBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedAuthzUserBaseDN = slice
	}
	if internaltypes.IsDefined(plan.AllIncludedAuthzUserGroupDN) {
		var slice []string
		plan.AllIncludedAuthzUserGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedAuthzUserGroupDN = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedAuthzUserGroupDN) {
		var slice []string
		plan.AnyIncludedAuthzUserGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedAuthzUserGroupDN = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedAuthzUserGroupDN) {
		var slice []string
		plan.NotAllIncludedAuthzUserGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedAuthzUserGroupDN = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedAuthzUserGroupDN) {
		var slice []string
		plan.NoneIncludedAuthzUserGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedAuthzUserGroupDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for aggregate result-criteria
func addOptionalAggregateResultCriteriaFields(ctx context.Context, addRequest *client.AddAggregateResultCriteriaRequest, plan resultCriteriaResourceModel) error {
	if internaltypes.IsDefined(plan.AllIncludedResultCriteria) {
		var slice []string
		plan.AllIncludedResultCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedResultCriteria = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedResultCriteria) {
		var slice []string
		plan.AnyIncludedResultCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedResultCriteria = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedResultCriteria) {
		var slice []string
		plan.NotAllIncludedResultCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedResultCriteria = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedResultCriteria) {
		var slice []string
		plan.NoneIncludedResultCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedResultCriteria = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for replication-assurance result-criteria
func addOptionalReplicationAssuranceResultCriteriaFields(ctx context.Context, addRequest *client.AddReplicationAssuranceResultCriteriaRequest, plan resultCriteriaResourceModel) error {
	if internaltypes.IsDefined(plan.LocalAssuranceLevel) {
		var slice []string
		plan.LocalAssuranceLevel.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumresultCriteriaLocalAssuranceLevelProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumresultCriteriaLocalAssuranceLevelPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.LocalAssuranceLevel = enumSlice
	}
	if internaltypes.IsDefined(plan.RemoteAssuranceLevel) {
		var slice []string
		plan.RemoteAssuranceLevel.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumresultCriteriaRemoteAssuranceLevelProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumresultCriteriaRemoteAssuranceLevelPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.RemoteAssuranceLevel = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AssuranceTimeoutCriteria) {
		assuranceTimeoutCriteria, err := client.NewEnumresultCriteriaAssuranceTimeoutCriteriaPropFromValue(plan.AssuranceTimeoutCriteria.ValueString())
		if err != nil {
			return err
		}
		addRequest.AssuranceTimeoutCriteria = assuranceTimeoutCriteria
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AssuranceTimeoutValue) {
		addRequest.AssuranceTimeoutValue = plan.AssuranceTimeoutValue.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ResponseDelayedByAssurance) {
		responseDelayedByAssurance, err := client.NewEnumresultCriteriaResponseDelayedByAssurancePropFromValue(plan.ResponseDelayedByAssurance.ValueString())
		if err != nil {
			return err
		}
		addRequest.ResponseDelayedByAssurance = responseDelayedByAssurance
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AssuranceBehaviorAlteredByControl) {
		assuranceBehaviorAlteredByControl, err := client.NewEnumresultCriteriaAssuranceBehaviorAlteredByControlPropFromValue(plan.AssuranceBehaviorAlteredByControl.ValueString())
		if err != nil {
			return err
		}
		addRequest.AssuranceBehaviorAlteredByControl = assuranceBehaviorAlteredByControl
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AssuranceSatisfied) {
		assuranceSatisfied, err := client.NewEnumresultCriteriaAssuranceSatisfiedPropFromValue(plan.AssuranceSatisfied.ValueString())
		if err != nil {
			return err
		}
		addRequest.AssuranceSatisfied = assuranceSatisfied
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for third-party result-criteria
func addOptionalThirdPartyResultCriteriaFields(ctx context.Context, addRequest *client.AddThirdPartyResultCriteriaRequest, plan resultCriteriaResourceModel) error {
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
func populateResultCriteriaUnknownValues(model *resultCriteriaResourceModel) {
	if model.NoneIncludedResponseControl.IsUnknown() || model.NoneIncludedResponseControl.IsNull() {
		model.NoneIncludedResponseControl, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.NotAllIncludedResultCriteria.IsUnknown() || model.NotAllIncludedResultCriteria.IsNull() {
		model.NotAllIncludedResultCriteria, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AnyIncludedAuthzUserGroupDN.IsUnknown() || model.AnyIncludedAuthzUserGroupDN.IsNull() {
		model.AnyIncludedAuthzUserGroupDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.LocalAssuranceLevel.IsUnknown() || model.LocalAssuranceLevel.IsNull() {
		model.LocalAssuranceLevel, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExcludedUserGroupDN.IsUnknown() || model.ExcludedUserGroupDN.IsNull() {
		model.ExcludedUserGroupDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AllIncludedResponseControl.IsUnknown() || model.AllIncludedResponseControl.IsNull() {
		model.AllIncludedResponseControl, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.NotAllIncludedAuthzUserGroupDN.IsUnknown() || model.NotAllIncludedAuthzUserGroupDN.IsNull() {
		model.NotAllIncludedAuthzUserGroupDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ResultCodeValue.IsUnknown() || model.ResultCodeValue.IsNull() {
		model.ResultCodeValue, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AnyIncludedResultCriteria.IsUnknown() || model.AnyIncludedResultCriteria.IsNull() {
		model.AnyIncludedResultCriteria, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.NotAllIncludedResponseControl.IsUnknown() || model.NotAllIncludedResponseControl.IsNull() {
		model.NotAllIncludedResponseControl, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludedAuthzUserBaseDN.IsUnknown() || model.IncludedAuthzUserBaseDN.IsNull() {
		model.IncludedAuthzUserBaseDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExcludedUserBaseDN.IsUnknown() || model.ExcludedUserBaseDN.IsNull() {
		model.ExcludedUserBaseDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExcludedUserFilter.IsUnknown() || model.ExcludedUserFilter.IsNull() {
		model.ExcludedUserFilter, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AllIncludedAuthzUserGroupDN.IsUnknown() || model.AllIncludedAuthzUserGroupDN.IsNull() {
		model.AllIncludedAuthzUserGroupDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludedUserGroupDN.IsUnknown() || model.IncludedUserGroupDN.IsNull() {
		model.IncludedUserGroupDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludedUserFilter.IsUnknown() || model.IncludedUserFilter.IsNull() {
		model.IncludedUserFilter, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AllIncludedResultCriteria.IsUnknown() || model.AllIncludedResultCriteria.IsNull() {
		model.AllIncludedResultCriteria, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludedUserBaseDN.IsUnknown() || model.IncludedUserBaseDN.IsNull() {
		model.IncludedUserBaseDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.NoneIncludedAuthzUserGroupDN.IsUnknown() || model.NoneIncludedAuthzUserGroupDN.IsNull() {
		model.NoneIncludedAuthzUserGroupDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.UsedPrivilege.IsUnknown() || model.UsedPrivilege.IsNull() {
		model.UsedPrivilege, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.NoneIncludedResultCriteria.IsUnknown() || model.NoneIncludedResultCriteria.IsNull() {
		model.NoneIncludedResultCriteria, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AnyIncludedResponseControl.IsUnknown() || model.AnyIncludedResponseControl.IsNull() {
		model.AnyIncludedResponseControl, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExcludedAuthzUserBaseDN.IsUnknown() || model.ExcludedAuthzUserBaseDN.IsNull() {
		model.ExcludedAuthzUserBaseDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.MissingPrivilege.IsUnknown() || model.MissingPrivilege.IsNull() {
		model.MissingPrivilege, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.RemoteAssuranceLevel.IsUnknown() || model.RemoteAssuranceLevel.IsNull() {
		model.RemoteAssuranceLevel, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *resultCriteriaResourceModel) populateAllComputedStringAttributes() {
	if model.ResultCodeCriteria.IsUnknown() || model.ResultCodeCriteria.IsNull() {
		model.ResultCodeCriteria = types.StringValue("")
	}
	if model.ReferralReturned.IsUnknown() || model.ReferralReturned.IsNull() {
		model.ReferralReturned = types.StringValue("")
	}
	if model.SearchIndexedCriteria.IsUnknown() || model.SearchIndexedCriteria.IsNull() {
		model.SearchIndexedCriteria = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.AssuranceTimeoutValue.IsUnknown() || model.AssuranceTimeoutValue.IsNull() {
		model.AssuranceTimeoutValue = types.StringValue("")
	}
	if model.ExtensionClass.IsUnknown() || model.ExtensionClass.IsNull() {
		model.ExtensionClass = types.StringValue("")
	}
	if model.QueueTimeCriteria.IsUnknown() || model.QueueTimeCriteria.IsNull() {
		model.QueueTimeCriteria = types.StringValue("")
	}
	if model.UsedAlternateAuthzid.IsUnknown() || model.UsedAlternateAuthzid.IsNull() {
		model.UsedAlternateAuthzid = types.StringValue("")
	}
	if model.AssuranceBehaviorAlteredByControl.IsUnknown() || model.AssuranceBehaviorAlteredByControl.IsNull() {
		model.AssuranceBehaviorAlteredByControl = types.StringValue("")
	}
	if model.ProcessingTimeValue.IsUnknown() || model.ProcessingTimeValue.IsNull() {
		model.ProcessingTimeValue = types.StringValue("")
	}
	if model.QueueTimeValue.IsUnknown() || model.QueueTimeValue.IsNull() {
		model.QueueTimeValue = types.StringValue("")
	}
	if model.SearchEntryReturnedCriteria.IsUnknown() || model.SearchEntryReturnedCriteria.IsNull() {
		model.SearchEntryReturnedCriteria = types.StringValue("")
	}
	if model.AssuranceTimeoutCriteria.IsUnknown() || model.AssuranceTimeoutCriteria.IsNull() {
		model.AssuranceTimeoutCriteria = types.StringValue("")
	}
	if model.ProcessingTimeCriteria.IsUnknown() || model.ProcessingTimeCriteria.IsNull() {
		model.ProcessingTimeCriteria = types.StringValue("")
	}
	if model.SearchReferenceReturnedCriteria.IsUnknown() || model.SearchReferenceReturnedCriteria.IsNull() {
		model.SearchReferenceReturnedCriteria = types.StringValue("")
	}
	if model.RequestCriteria.IsUnknown() || model.RequestCriteria.IsNull() {
		model.RequestCriteria = types.StringValue("")
	}
	if model.UsedAnyPrivilege.IsUnknown() || model.UsedAnyPrivilege.IsNull() {
		model.UsedAnyPrivilege = types.StringValue("")
	}
	if model.AssuranceSatisfied.IsUnknown() || model.AssuranceSatisfied.IsNull() {
		model.AssuranceSatisfied = types.StringValue("")
	}
	if model.RetiredPasswordUsedForBind.IsUnknown() || model.RetiredPasswordUsedForBind.IsNull() {
		model.RetiredPasswordUsedForBind = types.StringValue("")
	}
	if model.MissingAnyPrivilege.IsUnknown() || model.MissingAnyPrivilege.IsNull() {
		model.MissingAnyPrivilege = types.StringValue("")
	}
	if model.ResponseDelayedByAssurance.IsUnknown() || model.ResponseDelayedByAssurance.IsNull() {
		model.ResponseDelayedByAssurance = types.StringValue("")
	}
}

// Read a SuccessfulBindResultCriteriaResponse object into the model struct
func readSuccessfulBindResultCriteriaResponse(ctx context.Context, r *client.SuccessfulBindResultCriteriaResponse, state *resultCriteriaResourceModel, expectedValues *resultCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("successful-bind")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.IncludeAnonymousBinds = internaltypes.BoolTypeOrNil(r.IncludeAnonymousBinds)
	state.IncludedUserBaseDN = internaltypes.GetStringSet(r.IncludedUserBaseDN)
	state.ExcludedUserBaseDN = internaltypes.GetStringSet(r.ExcludedUserBaseDN)
	state.IncludedUserFilter = internaltypes.GetStringSet(r.IncludedUserFilter)
	state.ExcludedUserFilter = internaltypes.GetStringSet(r.ExcludedUserFilter)
	state.IncludedUserGroupDN = internaltypes.GetStringSet(r.IncludedUserGroupDN)
	state.ExcludedUserGroupDN = internaltypes.GetStringSet(r.ExcludedUserGroupDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateResultCriteriaUnknownValues(state)
}

// Read a SimpleResultCriteriaResponse object into the model struct
func readSimpleResultCriteriaResponse(ctx context.Context, r *client.SimpleResultCriteriaResponse, state *resultCriteriaResourceModel, expectedValues *resultCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("simple")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCodeCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaResultCodeCriteriaProp(r.ResultCodeCriteria), true)
	state.ResultCodeValue = internaltypes.GetStringSet(
		client.StringSliceEnumresultCriteriaResultCodeValueProp(r.ResultCodeValue))
	state.ProcessingTimeCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaProcessingTimeCriteriaProp(r.ProcessingTimeCriteria), true)
	state.ProcessingTimeValue = internaltypes.StringTypeOrNil(r.ProcessingTimeValue, true)
	config.CheckMismatchedPDFormattedAttributes("processing_time_value",
		expectedValues.ProcessingTimeValue, state.ProcessingTimeValue, diagnostics)
	state.QueueTimeCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaQueueTimeCriteriaProp(r.QueueTimeCriteria), true)
	state.QueueTimeValue = internaltypes.StringTypeOrNil(r.QueueTimeValue, true)
	config.CheckMismatchedPDFormattedAttributes("queue_time_value",
		expectedValues.QueueTimeValue, state.QueueTimeValue, diagnostics)
	state.ReferralReturned = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaReferralReturnedProp(r.ReferralReturned), true)
	state.AllIncludedResponseControl = internaltypes.GetStringSet(r.AllIncludedResponseControl)
	state.AnyIncludedResponseControl = internaltypes.GetStringSet(r.AnyIncludedResponseControl)
	state.NotAllIncludedResponseControl = internaltypes.GetStringSet(r.NotAllIncludedResponseControl)
	state.NoneIncludedResponseControl = internaltypes.GetStringSet(r.NoneIncludedResponseControl)
	state.UsedAlternateAuthzid = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaUsedAlternateAuthzidProp(r.UsedAlternateAuthzid), true)
	state.UsedAnyPrivilege = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaUsedAnyPrivilegeProp(r.UsedAnyPrivilege), true)
	state.UsedPrivilege = internaltypes.GetStringSet(
		client.StringSliceEnumresultCriteriaUsedPrivilegeProp(r.UsedPrivilege))
	state.MissingAnyPrivilege = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaMissingAnyPrivilegeProp(r.MissingAnyPrivilege), true)
	state.MissingPrivilege = internaltypes.GetStringSet(
		client.StringSliceEnumresultCriteriaMissingPrivilegeProp(r.MissingPrivilege))
	state.RetiredPasswordUsedForBind = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaRetiredPasswordUsedForBindProp(r.RetiredPasswordUsedForBind), true)
	state.SearchEntryReturnedCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaSearchEntryReturnedCriteriaProp(r.SearchEntryReturnedCriteria), true)
	state.SearchEntryReturnedCount = internaltypes.Int64TypeOrNil(r.SearchEntryReturnedCount)
	state.SearchReferenceReturnedCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaSearchReferenceReturnedCriteriaProp(r.SearchReferenceReturnedCriteria), true)
	state.SearchReferenceReturnedCount = internaltypes.Int64TypeOrNil(r.SearchReferenceReturnedCount)
	state.SearchIndexedCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaSearchIndexedCriteriaProp(r.SearchIndexedCriteria), true)
	state.IncludedAuthzUserBaseDN = internaltypes.GetStringSet(r.IncludedAuthzUserBaseDN)
	state.ExcludedAuthzUserBaseDN = internaltypes.GetStringSet(r.ExcludedAuthzUserBaseDN)
	state.AllIncludedAuthzUserGroupDN = internaltypes.GetStringSet(r.AllIncludedAuthzUserGroupDN)
	state.AnyIncludedAuthzUserGroupDN = internaltypes.GetStringSet(r.AnyIncludedAuthzUserGroupDN)
	state.NotAllIncludedAuthzUserGroupDN = internaltypes.GetStringSet(r.NotAllIncludedAuthzUserGroupDN)
	state.NoneIncludedAuthzUserGroupDN = internaltypes.GetStringSet(r.NoneIncludedAuthzUserGroupDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateResultCriteriaUnknownValues(state)
}

// Read a AggregateResultCriteriaResponse object into the model struct
func readAggregateResultCriteriaResponse(ctx context.Context, r *client.AggregateResultCriteriaResponse, state *resultCriteriaResourceModel, expectedValues *resultCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("aggregate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllIncludedResultCriteria = internaltypes.GetStringSet(r.AllIncludedResultCriteria)
	state.AnyIncludedResultCriteria = internaltypes.GetStringSet(r.AnyIncludedResultCriteria)
	state.NotAllIncludedResultCriteria = internaltypes.GetStringSet(r.NotAllIncludedResultCriteria)
	state.NoneIncludedResultCriteria = internaltypes.GetStringSet(r.NoneIncludedResultCriteria)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateResultCriteriaUnknownValues(state)
}

// Read a ReplicationAssuranceResultCriteriaResponse object into the model struct
func readReplicationAssuranceResultCriteriaResponse(ctx context.Context, r *client.ReplicationAssuranceResultCriteriaResponse, state *resultCriteriaResourceModel, expectedValues *resultCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("replication-assurance")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LocalAssuranceLevel = internaltypes.GetStringSet(
		client.StringSliceEnumresultCriteriaLocalAssuranceLevelProp(r.LocalAssuranceLevel))
	state.RemoteAssuranceLevel = internaltypes.GetStringSet(
		client.StringSliceEnumresultCriteriaRemoteAssuranceLevelProp(r.RemoteAssuranceLevel))
	state.AssuranceTimeoutCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaAssuranceTimeoutCriteriaProp(r.AssuranceTimeoutCriteria), true)
	state.AssuranceTimeoutValue = internaltypes.StringTypeOrNil(r.AssuranceTimeoutValue, true)
	config.CheckMismatchedPDFormattedAttributes("assurance_timeout_value",
		expectedValues.AssuranceTimeoutValue, state.AssuranceTimeoutValue, diagnostics)
	state.ResponseDelayedByAssurance = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaResponseDelayedByAssuranceProp(r.ResponseDelayedByAssurance), true)
	state.AssuranceBehaviorAlteredByControl = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaAssuranceBehaviorAlteredByControlProp(r.AssuranceBehaviorAlteredByControl), true)
	state.AssuranceSatisfied = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaAssuranceSatisfiedProp(r.AssuranceSatisfied), true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateResultCriteriaUnknownValues(state)
}

// Read a ThirdPartyResultCriteriaResponse object into the model struct
func readThirdPartyResultCriteriaResponse(ctx context.Context, r *client.ThirdPartyResultCriteriaResponse, state *resultCriteriaResourceModel, expectedValues *resultCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateResultCriteriaUnknownValues(state)
}

// Create any update operations necessary to make the state match the plan
func createResultCriteriaOperations(plan resultCriteriaResourceModel, state resultCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.LocalAssuranceLevel, state.LocalAssuranceLevel, "local-assurance-level")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RemoteAssuranceLevel, state.RemoteAssuranceLevel, "remote-assurance-level")
	operations.AddStringOperationIfNecessary(&ops, plan.AssuranceTimeoutCriteria, state.AssuranceTimeoutCriteria, "assurance-timeout-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.AssuranceTimeoutValue, state.AssuranceTimeoutValue, "assurance-timeout-value")
	operations.AddStringOperationIfNecessary(&ops, plan.ResponseDelayedByAssurance, state.ResponseDelayedByAssurance, "response-delayed-by-assurance")
	operations.AddStringOperationIfNecessary(&ops, plan.AssuranceBehaviorAlteredByControl, state.AssuranceBehaviorAlteredByControl, "assurance-behavior-altered-by-control")
	operations.AddStringOperationIfNecessary(&ops, plan.AssuranceSatisfied, state.AssuranceSatisfied, "assurance-satisfied")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedResultCriteria, state.AllIncludedResultCriteria, "all-included-result-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedResultCriteria, state.AnyIncludedResultCriteria, "any-included-result-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedResultCriteria, state.NotAllIncludedResultCriteria, "not-all-included-result-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedResultCriteria, state.NoneIncludedResultCriteria, "none-included-result-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.RequestCriteria, state.RequestCriteria, "request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.ResultCodeCriteria, state.ResultCodeCriteria, "result-code-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ResultCodeValue, state.ResultCodeValue, "result-code-value")
	operations.AddStringOperationIfNecessary(&ops, plan.ProcessingTimeCriteria, state.ProcessingTimeCriteria, "processing-time-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.ProcessingTimeValue, state.ProcessingTimeValue, "processing-time-value")
	operations.AddStringOperationIfNecessary(&ops, plan.QueueTimeCriteria, state.QueueTimeCriteria, "queue-time-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.QueueTimeValue, state.QueueTimeValue, "queue-time-value")
	operations.AddStringOperationIfNecessary(&ops, plan.ReferralReturned, state.ReferralReturned, "referral-returned")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedResponseControl, state.AllIncludedResponseControl, "all-included-response-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedResponseControl, state.AnyIncludedResponseControl, "any-included-response-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedResponseControl, state.NotAllIncludedResponseControl, "not-all-included-response-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedResponseControl, state.NoneIncludedResponseControl, "none-included-response-control")
	operations.AddStringOperationIfNecessary(&ops, plan.UsedAlternateAuthzid, state.UsedAlternateAuthzid, "used-alternate-authzid")
	operations.AddStringOperationIfNecessary(&ops, plan.UsedAnyPrivilege, state.UsedAnyPrivilege, "used-any-privilege")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.UsedPrivilege, state.UsedPrivilege, "used-privilege")
	operations.AddStringOperationIfNecessary(&ops, plan.MissingAnyPrivilege, state.MissingAnyPrivilege, "missing-any-privilege")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MissingPrivilege, state.MissingPrivilege, "missing-privilege")
	operations.AddStringOperationIfNecessary(&ops, plan.RetiredPasswordUsedForBind, state.RetiredPasswordUsedForBind, "retired-password-used-for-bind")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchEntryReturnedCriteria, state.SearchEntryReturnedCriteria, "search-entry-returned-criteria")
	operations.AddInt64OperationIfNecessary(&ops, plan.SearchEntryReturnedCount, state.SearchEntryReturnedCount, "search-entry-returned-count")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchReferenceReturnedCriteria, state.SearchReferenceReturnedCriteria, "search-reference-returned-criteria")
	operations.AddInt64OperationIfNecessary(&ops, plan.SearchReferenceReturnedCount, state.SearchReferenceReturnedCount, "search-reference-returned-count")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchIndexedCriteria, state.SearchIndexedCriteria, "search-indexed-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedAuthzUserBaseDN, state.IncludedAuthzUserBaseDN, "included-authz-user-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedAuthzUserBaseDN, state.ExcludedAuthzUserBaseDN, "excluded-authz-user-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedAuthzUserGroupDN, state.AllIncludedAuthzUserGroupDN, "all-included-authz-user-group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedAuthzUserGroupDN, state.AnyIncludedAuthzUserGroupDN, "any-included-authz-user-group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedAuthzUserGroupDN, state.NotAllIncludedAuthzUserGroupDN, "not-all-included-authz-user-group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedAuthzUserGroupDN, state.NoneIncludedAuthzUserGroupDN, "none-included-authz-user-group-dn")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeAnonymousBinds, state.IncludeAnonymousBinds, "include-anonymous-binds")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedUserBaseDN, state.IncludedUserBaseDN, "included-user-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedUserBaseDN, state.ExcludedUserBaseDN, "excluded-user-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedUserFilter, state.IncludedUserFilter, "included-user-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedUserFilter, state.ExcludedUserFilter, "excluded-user-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedUserGroupDN, state.IncludedUserGroupDN, "included-user-group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedUserGroupDN, state.ExcludedUserGroupDN, "excluded-user-group-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a successful-bind result-criteria
func (r *resultCriteriaResource) CreateSuccessfulBindResultCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan resultCriteriaResourceModel) (*resultCriteriaResourceModel, error) {
	addRequest := client.NewAddSuccessfulBindResultCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumsuccessfulBindResultCriteriaSchemaUrn{client.ENUMSUCCESSFULBINDRESULTCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RESULT_CRITERIASUCCESSFUL_BIND})
	err := addOptionalSuccessfulBindResultCriteriaFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Result Criteria", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ResultCriteriaApi.AddResultCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddResultCriteriaRequest(
		client.AddSuccessfulBindResultCriteriaRequestAsAddResultCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ResultCriteriaApi.AddResultCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Result Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state resultCriteriaResourceModel
	readSuccessfulBindResultCriteriaResponse(ctx, addResponse.SuccessfulBindResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a simple result-criteria
func (r *resultCriteriaResource) CreateSimpleResultCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan resultCriteriaResourceModel) (*resultCriteriaResourceModel, error) {
	addRequest := client.NewAddSimpleResultCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumsimpleResultCriteriaSchemaUrn{client.ENUMSIMPLERESULTCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RESULT_CRITERIASIMPLE})
	err := addOptionalSimpleResultCriteriaFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Result Criteria", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ResultCriteriaApi.AddResultCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddResultCriteriaRequest(
		client.AddSimpleResultCriteriaRequestAsAddResultCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ResultCriteriaApi.AddResultCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Result Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state resultCriteriaResourceModel
	readSimpleResultCriteriaResponse(ctx, addResponse.SimpleResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a aggregate result-criteria
func (r *resultCriteriaResource) CreateAggregateResultCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan resultCriteriaResourceModel) (*resultCriteriaResourceModel, error) {
	addRequest := client.NewAddAggregateResultCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumaggregateResultCriteriaSchemaUrn{client.ENUMAGGREGATERESULTCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RESULT_CRITERIAAGGREGATE})
	err := addOptionalAggregateResultCriteriaFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Result Criteria", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ResultCriteriaApi.AddResultCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddResultCriteriaRequest(
		client.AddAggregateResultCriteriaRequestAsAddResultCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ResultCriteriaApi.AddResultCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Result Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state resultCriteriaResourceModel
	readAggregateResultCriteriaResponse(ctx, addResponse.AggregateResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a replication-assurance result-criteria
func (r *resultCriteriaResource) CreateReplicationAssuranceResultCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan resultCriteriaResourceModel) (*resultCriteriaResourceModel, error) {
	addRequest := client.NewAddReplicationAssuranceResultCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumreplicationAssuranceResultCriteriaSchemaUrn{client.ENUMREPLICATIONASSURANCERESULTCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RESULT_CRITERIAREPLICATION_ASSURANCE})
	err := addOptionalReplicationAssuranceResultCriteriaFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Result Criteria", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ResultCriteriaApi.AddResultCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddResultCriteriaRequest(
		client.AddReplicationAssuranceResultCriteriaRequestAsAddResultCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ResultCriteriaApi.AddResultCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Result Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state resultCriteriaResourceModel
	readReplicationAssuranceResultCriteriaResponse(ctx, addResponse.ReplicationAssuranceResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party result-criteria
func (r *resultCriteriaResource) CreateThirdPartyResultCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan resultCriteriaResourceModel) (*resultCriteriaResourceModel, error) {
	addRequest := client.NewAddThirdPartyResultCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyResultCriteriaSchemaUrn{client.ENUMTHIRDPARTYRESULTCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RESULT_CRITERIATHIRD_PARTY},
		plan.ExtensionClass.ValueString())
	err := addOptionalThirdPartyResultCriteriaFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Result Criteria", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ResultCriteriaApi.AddResultCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddResultCriteriaRequest(
		client.AddThirdPartyResultCriteriaRequestAsAddResultCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ResultCriteriaApi.AddResultCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Result Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state resultCriteriaResourceModel
	readThirdPartyResultCriteriaResponse(ctx, addResponse.ThirdPartyResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *resultCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan resultCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *resultCriteriaResourceModel
	var err error
	if plan.Type.ValueString() == "successful-bind" {
		state, err = r.CreateSuccessfulBindResultCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "simple" {
		state, err = r.CreateSimpleResultCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "aggregate" {
		state, err = r.CreateAggregateResultCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "replication-assurance" {
		state, err = r.CreateReplicationAssuranceResultCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyResultCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

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
func (r *defaultResultCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan resultCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ResultCriteriaApi.GetResultCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Result Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state resultCriteriaResourceModel
	if readResponse.SuccessfulBindResultCriteriaResponse != nil {
		readSuccessfulBindResultCriteriaResponse(ctx, readResponse.SuccessfulBindResultCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SimpleResultCriteriaResponse != nil {
		readSimpleResultCriteriaResponse(ctx, readResponse.SimpleResultCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AggregateResultCriteriaResponse != nil {
		readAggregateResultCriteriaResponse(ctx, readResponse.AggregateResultCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ReplicationAssuranceResultCriteriaResponse != nil {
		readReplicationAssuranceResultCriteriaResponse(ctx, readResponse.ReplicationAssuranceResultCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyResultCriteriaResponse != nil {
		readThirdPartyResultCriteriaResponse(ctx, readResponse.ThirdPartyResultCriteriaResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ResultCriteriaApi.UpdateResultCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createResultCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ResultCriteriaApi.UpdateResultCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Result Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.SuccessfulBindResultCriteriaResponse != nil {
			readSuccessfulBindResultCriteriaResponse(ctx, updateResponse.SuccessfulBindResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SimpleResultCriteriaResponse != nil {
			readSimpleResultCriteriaResponse(ctx, updateResponse.SimpleResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AggregateResultCriteriaResponse != nil {
			readAggregateResultCriteriaResponse(ctx, updateResponse.AggregateResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ReplicationAssuranceResultCriteriaResponse != nil {
			readReplicationAssuranceResultCriteriaResponse(ctx, updateResponse.ReplicationAssuranceResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyResultCriteriaResponse != nil {
			readThirdPartyResultCriteriaResponse(ctx, updateResponse.ThirdPartyResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
	}

	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *resultCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readResultCriteria(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultResultCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readResultCriteria(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readResultCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state resultCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ResultCriteriaApi.GetResultCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Result Criteria", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Result Criteria", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.SuccessfulBindResultCriteriaResponse != nil {
		readSuccessfulBindResultCriteriaResponse(ctx, readResponse.SuccessfulBindResultCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SimpleResultCriteriaResponse != nil {
		readSimpleResultCriteriaResponse(ctx, readResponse.SimpleResultCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AggregateResultCriteriaResponse != nil {
		readAggregateResultCriteriaResponse(ctx, readResponse.AggregateResultCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ReplicationAssuranceResultCriteriaResponse != nil {
		readReplicationAssuranceResultCriteriaResponse(ctx, readResponse.ReplicationAssuranceResultCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyResultCriteriaResponse != nil {
		readThirdPartyResultCriteriaResponse(ctx, readResponse.ThirdPartyResultCriteriaResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *resultCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateResultCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultResultCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateResultCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateResultCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan resultCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state resultCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ResultCriteriaApi.UpdateResultCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createResultCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ResultCriteriaApi.UpdateResultCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Result Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.SuccessfulBindResultCriteriaResponse != nil {
			readSuccessfulBindResultCriteriaResponse(ctx, updateResponse.SuccessfulBindResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SimpleResultCriteriaResponse != nil {
			readSimpleResultCriteriaResponse(ctx, updateResponse.SimpleResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AggregateResultCriteriaResponse != nil {
			readAggregateResultCriteriaResponse(ctx, updateResponse.AggregateResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ReplicationAssuranceResultCriteriaResponse != nil {
			readReplicationAssuranceResultCriteriaResponse(ctx, updateResponse.ReplicationAssuranceResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyResultCriteriaResponse != nil {
			readThirdPartyResultCriteriaResponse(ctx, updateResponse.ThirdPartyResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *defaultResultCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *resultCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state resultCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ResultCriteriaApi.DeleteResultCriteriaExecute(r.apiClient.ResultCriteriaApi.DeleteResultCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Result Criteria", err, httpResp)
		return
	}
}

func (r *resultCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importResultCriteria(ctx, req, resp)
}

func (r *defaultResultCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importResultCriteria(ctx, req, resp)
}

func importResultCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
