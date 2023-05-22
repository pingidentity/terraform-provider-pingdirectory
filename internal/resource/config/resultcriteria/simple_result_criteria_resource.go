package resultcriteria

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &simpleResultCriteriaResource{}
	_ resource.ResourceWithConfigure   = &simpleResultCriteriaResource{}
	_ resource.ResourceWithImportState = &simpleResultCriteriaResource{}
	_ resource.Resource                = &defaultSimpleResultCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultSimpleResultCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultSimpleResultCriteriaResource{}
)

// Create a Simple Result Criteria resource
func NewSimpleResultCriteriaResource() resource.Resource {
	return &simpleResultCriteriaResource{}
}

func NewDefaultSimpleResultCriteriaResource() resource.Resource {
	return &defaultSimpleResultCriteriaResource{}
}

// simpleResultCriteriaResource is the resource implementation.
type simpleResultCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSimpleResultCriteriaResource is the resource implementation.
type defaultSimpleResultCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *simpleResultCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_simple_result_criteria"
}

func (r *defaultSimpleResultCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_simple_result_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *simpleResultCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultSimpleResultCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type simpleResultCriteriaResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	LastUpdated                     types.String `tfsdk:"last_updated"`
	Notifications                   types.Set    `tfsdk:"notifications"`
	RequiredActions                 types.Set    `tfsdk:"required_actions"`
	RequestCriteria                 types.String `tfsdk:"request_criteria"`
	ResultCodeCriteria              types.String `tfsdk:"result_code_criteria"`
	ResultCodeValue                 types.Set    `tfsdk:"result_code_value"`
	ProcessingTimeCriteria          types.String `tfsdk:"processing_time_criteria"`
	ProcessingTimeValue             types.String `tfsdk:"processing_time_value"`
	QueueTimeCriteria               types.String `tfsdk:"queue_time_criteria"`
	QueueTimeValue                  types.String `tfsdk:"queue_time_value"`
	ReferralReturned                types.String `tfsdk:"referral_returned"`
	AllIncludedResponseControl      types.Set    `tfsdk:"all_included_response_control"`
	AnyIncludedResponseControl      types.Set    `tfsdk:"any_included_response_control"`
	NotAllIncludedResponseControl   types.Set    `tfsdk:"not_all_included_response_control"`
	NoneIncludedResponseControl     types.Set    `tfsdk:"none_included_response_control"`
	UsedAlternateAuthzid            types.String `tfsdk:"used_alternate_authzid"`
	UsedAnyPrivilege                types.String `tfsdk:"used_any_privilege"`
	UsedPrivilege                   types.Set    `tfsdk:"used_privilege"`
	MissingAnyPrivilege             types.String `tfsdk:"missing_any_privilege"`
	MissingPrivilege                types.Set    `tfsdk:"missing_privilege"`
	RetiredPasswordUsedForBind      types.String `tfsdk:"retired_password_used_for_bind"`
	SearchEntryReturnedCriteria     types.String `tfsdk:"search_entry_returned_criteria"`
	SearchEntryReturnedCount        types.Int64  `tfsdk:"search_entry_returned_count"`
	SearchReferenceReturnedCriteria types.String `tfsdk:"search_reference_returned_criteria"`
	SearchReferenceReturnedCount    types.Int64  `tfsdk:"search_reference_returned_count"`
	SearchIndexedCriteria           types.String `tfsdk:"search_indexed_criteria"`
	IncludedAuthzUserBaseDN         types.Set    `tfsdk:"included_authz_user_base_dn"`
	ExcludedAuthzUserBaseDN         types.Set    `tfsdk:"excluded_authz_user_base_dn"`
	AllIncludedAuthzUserGroupDN     types.Set    `tfsdk:"all_included_authz_user_group_dn"`
	AnyIncludedAuthzUserGroupDN     types.Set    `tfsdk:"any_included_authz_user_group_dn"`
	NotAllIncludedAuthzUserGroupDN  types.Set    `tfsdk:"not_all_included_authz_user_group_dn"`
	NoneIncludedAuthzUserGroupDN    types.Set    `tfsdk:"none_included_authz_user_group_dn"`
	Description                     types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *simpleResultCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	simpleResultCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultSimpleResultCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	simpleResultCriteriaSchema(ctx, req, resp, true)
}

func simpleResultCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Simple Result Criteria.",
		Attributes: map[string]schema.Attribute{
			"request_criteria": schema.StringAttribute{
				Description: "Specifies a request criteria object that must match the associated request for operations included in this Simple Result Criteria.",
				Optional:    true,
			},
			"result_code_criteria": schema.StringAttribute{
				Description: "Specifies which operation result codes are allowed for operations included in this Simple Result Criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"result_code_value": schema.SetAttribute{
				Description: "Specifies the operation result code values for results included in this Simple Result Criteria. This will only be taken into account if the \"result-code-criteria\" property has a value of \"selected-result-codes\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"processing_time_criteria": schema.StringAttribute{
				Description: "Indicates whether the time required to process the operation should be taken into consideration when determining whether to include the operation in this Simple Result Criteria. If the processing time should be taken into account, then the \"processing-time-value\" property should contain the boundary value.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"all_included_response_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must be present in the response to the client for operations included in this Simple Result Criteria. If any control OIDs are provided, then the response must contain all of those controls.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"any_included_response_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that may be present in the response to the client for operations included in this Simple Result Criteria. If any control OIDs are provided, then the response must contain at least one of those controls.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"not_all_included_response_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that should not be present in the response to the client for operations included in this Simple Result Criteria. If any control OIDs are provided, then the response must not contain at least one of those controls (that is, the response may contain zero or more of those controls, but not all of them).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"none_included_response_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must not be present in the response to the client for operations included in this Simple Result Criteria. If any control OIDs are provided, then the response must not contain any of those controls.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"used_alternate_authzid": schema.StringAttribute{
				Description: "Indicates whether operation results in which the associated operation used an authorization identity that is different from the authentication identity (e.g., as the result of using a proxied authorization control) should be included in this Simple Result Criteria. If no value is provided, then whether an operation used an alternate authorization identity will not be considered when determining whether it matches this Simple Result Criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"used_any_privilege": schema.StringAttribute{
				Description: "Indicates whether operations in which one or more privileges were used should be included in this Simple Result Criteria. If no value is provided, then whether an operation used any privileges will not be considered when determining whether it matches this Simple Result Criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"used_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that must have been used during the processing for operations included in this Simple Result Criteria. If any privilege names are provided, then the associated operation must have used at least one of those privileges. If no privilege names were provided, then the set of privileges used will not be considered when determining whether an operation should be included in this Simple Result Criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"missing_any_privilege": schema.StringAttribute{
				Description: "Indicates whether operations in which one or more privileges were missing should be included in this Simple Result Criteria. If no value is provided, then whether there were any missing privileges will not be considered when determining whether an operation matches this Simple Result Criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"missing_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that must have been missing during the processing for operations included in this Simple Result Criteria. If any privilege names are provided, then the associated operation must have been missing at least one of those privileges. If no privilege names were provided, then the set of privileges missing will not be considered when determining whether an operation should be included in this Simple Result Criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"retired_password_used_for_bind": schema.StringAttribute{
				Description: "Indicates whether the use of a retired password for authentication should be considered when determining whether a bind operation should be included in this Simple Result Criteria. This will be ignored for all operations other than bind.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"search_entry_returned_criteria": schema.StringAttribute{
				Description: "Indicates whether the number of entries returned should be considered when determining whether a search operation should be included in this Simple Result Criteria. This will be ignored for all operations other than search.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"search_entry_returned_count": schema.Int64Attribute{
				Description: "Specifies the target number of entries returned for use when determining whether a search operation should be included in this Simple Result Criteria. This will be ignored for all operations other than search, and it will be ignored for search operations if the \"search-entry-criteria\" property has a value of \"any\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"search_reference_returned_criteria": schema.StringAttribute{
				Description: "Indicates whether the number of references returned should be considered when determining whether a search operation should be included in this Simple Result Criteria. This will be ignored for all operations other than search.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"search_reference_returned_count": schema.Int64Attribute{
				Description: "Specifies the target number of references returned for use when determining whether a search operation should be included in this Simple Result Criteria. This will be ignored for all operations other than search, and it will be ignored for search operations if the \"search-reference-criteria\" property has a value of \"any\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"search_indexed_criteria": schema.StringAttribute{
				Description: "Indicates whether a search operation should be matched by this Simple Result Criteria based on whether it is considered indexed by the server. This will be ignored for all operations other than search.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"included_authz_user_base_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which authorization user entries may exist for operations included in this Simple Result Criteria. The authorization user could be the currently authenticated user on the connection (the user that performed the Bind operation), or different if proxied authorization was used to request that the operation be performed under the authorization of another user (as is the case for operations that come through a Directory Proxy Server). This property will be ignored for operations where no authentication or authorization has been performed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"excluded_authz_user_base_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which authorization user entries may exist for operations excluded from this Simple Result Criteria. The authorization user could be the currently authenticated user on the connection (the user that performed the Bind operation), or different if proxied authorization was used to request that the operation be performed under the authorization of another user (as is the case for operations that come through a Directory Proxy Server). This property will be ignored for operations where no authentication or authorization has been performed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"all_included_authz_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authorization users must exist for operations included in this Simple Result Criteria. If any group DNs are provided, then the authorization user must be a member of all of those groups. The authorization user could be the currently authenticated user on the connection (the user that performed the Bind operation), or different if proxied authorization was used to request that the operation be performed under the authorization of another user (as is the case for operations that come through a Directory Proxy Server). This property will be ignored for operations where no authentication or authorization has been performed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"any_included_authz_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authorization users may exist for operations included in this Simple Result Criteria. If any group DNs are provided, then the authorization user must be a member of at least one of those groups. The authorization user could be the currently authenticated user on the connection (the user that performed the Bind operation), or different if proxied authorization was used to request that the operation be performed under the authorization of another user (as is the case for operations that come through a Directory Proxy Server). This property will be ignored for operations where no authentication or authorization has been performed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"not_all_included_authz_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authorization users should not exist for operations included in this Simple Result Criteria. If any group DNs are provided, then the authorization user must not be a member of at least one of those groups (that is, the user may be a member of zero or more of those groups, but not of all of them). The authorization user could be the currently authenticated user on the connection (the user that performed the Bind operation), or different if proxied authorization was used to request that the operation be performed under the authorization of another user (as is the case for operations that come through a Directory Proxy Server). This property will be ignored for operations where no authentication or authorization has been performed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"none_included_authz_user_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which authorization users must not exist for operations included in this Simple Result Criteria. If any group DNs are provided, then the authorization user must not be a member any of those groups. The authorization user could be the currently authenticated user on the connection (the user that performed the Bind operation), or different if proxied authorization was used to request that the operation be performed under the authorization of another user (as is the case for operations that come through a Directory Proxy Server). This property will be ignored for operations where no authentication or authorization has been performed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Result Criteria",
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
func addOptionalSimpleResultCriteriaFields(ctx context.Context, addRequest *client.AddSimpleResultCriteriaRequest, plan simpleResultCriteriaResourceModel) error {
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

// Read a SimpleResultCriteriaResponse object into the model struct
func readSimpleResultCriteriaResponse(ctx context.Context, r *client.SimpleResultCriteriaResponse, state *simpleResultCriteriaResourceModel, expectedValues *simpleResultCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCodeCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaResultCodeCriteriaProp(r.ResultCodeCriteria), internaltypes.IsEmptyString(expectedValues.ResultCodeCriteria))
	state.ResultCodeValue = internaltypes.GetStringSet(
		client.StringSliceEnumresultCriteriaResultCodeValueProp(r.ResultCodeValue))
	state.ProcessingTimeCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaProcessingTimeCriteriaProp(r.ProcessingTimeCriteria), internaltypes.IsEmptyString(expectedValues.ProcessingTimeCriteria))
	state.ProcessingTimeValue = internaltypes.StringTypeOrNil(r.ProcessingTimeValue, internaltypes.IsEmptyString(expectedValues.ProcessingTimeValue))
	config.CheckMismatchedPDFormattedAttributes("processing_time_value",
		expectedValues.ProcessingTimeValue, state.ProcessingTimeValue, diagnostics)
	state.QueueTimeCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaQueueTimeCriteriaProp(r.QueueTimeCriteria), internaltypes.IsEmptyString(expectedValues.QueueTimeCriteria))
	state.QueueTimeValue = internaltypes.StringTypeOrNil(r.QueueTimeValue, internaltypes.IsEmptyString(expectedValues.QueueTimeValue))
	config.CheckMismatchedPDFormattedAttributes("queue_time_value",
		expectedValues.QueueTimeValue, state.QueueTimeValue, diagnostics)
	state.ReferralReturned = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaReferralReturnedProp(r.ReferralReturned), internaltypes.IsEmptyString(expectedValues.ReferralReturned))
	state.AllIncludedResponseControl = internaltypes.GetStringSet(r.AllIncludedResponseControl)
	state.AnyIncludedResponseControl = internaltypes.GetStringSet(r.AnyIncludedResponseControl)
	state.NotAllIncludedResponseControl = internaltypes.GetStringSet(r.NotAllIncludedResponseControl)
	state.NoneIncludedResponseControl = internaltypes.GetStringSet(r.NoneIncludedResponseControl)
	state.UsedAlternateAuthzid = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaUsedAlternateAuthzidProp(r.UsedAlternateAuthzid), internaltypes.IsEmptyString(expectedValues.UsedAlternateAuthzid))
	state.UsedAnyPrivilege = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaUsedAnyPrivilegeProp(r.UsedAnyPrivilege), internaltypes.IsEmptyString(expectedValues.UsedAnyPrivilege))
	state.UsedPrivilege = internaltypes.GetStringSet(
		client.StringSliceEnumresultCriteriaUsedPrivilegeProp(r.UsedPrivilege))
	state.MissingAnyPrivilege = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaMissingAnyPrivilegeProp(r.MissingAnyPrivilege), internaltypes.IsEmptyString(expectedValues.MissingAnyPrivilege))
	state.MissingPrivilege = internaltypes.GetStringSet(
		client.StringSliceEnumresultCriteriaMissingPrivilegeProp(r.MissingPrivilege))
	state.RetiredPasswordUsedForBind = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaRetiredPasswordUsedForBindProp(r.RetiredPasswordUsedForBind), internaltypes.IsEmptyString(expectedValues.RetiredPasswordUsedForBind))
	state.SearchEntryReturnedCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaSearchEntryReturnedCriteriaProp(r.SearchEntryReturnedCriteria), internaltypes.IsEmptyString(expectedValues.SearchEntryReturnedCriteria))
	state.SearchEntryReturnedCount = internaltypes.Int64TypeOrNil(r.SearchEntryReturnedCount)
	state.SearchReferenceReturnedCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaSearchReferenceReturnedCriteriaProp(r.SearchReferenceReturnedCriteria), internaltypes.IsEmptyString(expectedValues.SearchReferenceReturnedCriteria))
	state.SearchReferenceReturnedCount = internaltypes.Int64TypeOrNil(r.SearchReferenceReturnedCount)
	state.SearchIndexedCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumresultCriteriaSearchIndexedCriteriaProp(r.SearchIndexedCriteria), internaltypes.IsEmptyString(expectedValues.SearchIndexedCriteria))
	state.IncludedAuthzUserBaseDN = internaltypes.GetStringSet(r.IncludedAuthzUserBaseDN)
	state.ExcludedAuthzUserBaseDN = internaltypes.GetStringSet(r.ExcludedAuthzUserBaseDN)
	state.AllIncludedAuthzUserGroupDN = internaltypes.GetStringSet(r.AllIncludedAuthzUserGroupDN)
	state.AnyIncludedAuthzUserGroupDN = internaltypes.GetStringSet(r.AnyIncludedAuthzUserGroupDN)
	state.NotAllIncludedAuthzUserGroupDN = internaltypes.GetStringSet(r.NotAllIncludedAuthzUserGroupDN)
	state.NoneIncludedAuthzUserGroupDN = internaltypes.GetStringSet(r.NoneIncludedAuthzUserGroupDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSimpleResultCriteriaOperations(plan simpleResultCriteriaResourceModel, state simpleResultCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
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
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *simpleResultCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan simpleResultCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddSimpleResultCriteriaRequest(plan.Id.ValueString(),
		[]client.EnumsimpleResultCriteriaSchemaUrn{client.ENUMSIMPLERESULTCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RESULT_CRITERIASIMPLE})
	err := addOptionalSimpleResultCriteriaFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Simple Result Criteria", err.Error())
		return
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Simple Result Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state simpleResultCriteriaResourceModel
	readSimpleResultCriteriaResponse(ctx, addResponse.SimpleResultCriteriaResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultSimpleResultCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan simpleResultCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ResultCriteriaApi.GetResultCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Simple Result Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state simpleResultCriteriaResourceModel
	readSimpleResultCriteriaResponse(ctx, readResponse.SimpleResultCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ResultCriteriaApi.UpdateResultCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSimpleResultCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ResultCriteriaApi.UpdateResultCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Simple Result Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSimpleResultCriteriaResponse(ctx, updateResponse.SimpleResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *simpleResultCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSimpleResultCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSimpleResultCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSimpleResultCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSimpleResultCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state simpleResultCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ResultCriteriaApi.GetResultCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Simple Result Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSimpleResultCriteriaResponse(ctx, readResponse.SimpleResultCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *simpleResultCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSimpleResultCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSimpleResultCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSimpleResultCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSimpleResultCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan simpleResultCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state simpleResultCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ResultCriteriaApi.UpdateResultCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSimpleResultCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ResultCriteriaApi.UpdateResultCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Simple Result Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSimpleResultCriteriaResponse(ctx, updateResponse.SimpleResultCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSimpleResultCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *simpleResultCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state simpleResultCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ResultCriteriaApi.DeleteResultCriteriaExecute(r.apiClient.ResultCriteriaApi.DeleteResultCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Simple Result Criteria", err, httpResp)
		return
	}
}

func (r *simpleResultCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSimpleResultCriteria(ctx, req, resp)
}

func (r *defaultSimpleResultCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSimpleResultCriteria(ctx, req, resp)
}

func importSimpleResultCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
