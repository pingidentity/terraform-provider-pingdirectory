package requestcriteria

import (
	"context"
	"time"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &simpleRequestCriteriaResource{}
	_ resource.ResourceWithConfigure   = &simpleRequestCriteriaResource{}
	_ resource.ResourceWithImportState = &simpleRequestCriteriaResource{}
	_ resource.Resource                = &defaultSimpleRequestCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultSimpleRequestCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultSimpleRequestCriteriaResource{}
)

// Create a Simple Request Criteria resource
func NewSimpleRequestCriteriaResource() resource.Resource {
	return &simpleRequestCriteriaResource{}
}

func NewDefaultSimpleRequestCriteriaResource() resource.Resource {
	return &defaultSimpleRequestCriteriaResource{}
}

// simpleRequestCriteriaResource is the resource implementation.
type simpleRequestCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSimpleRequestCriteriaResource is the resource implementation.
type defaultSimpleRequestCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *simpleRequestCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_simple_request_criteria"
}

func (r *defaultSimpleRequestCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_simple_request_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *simpleRequestCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultSimpleRequestCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type simpleRequestCriteriaResourceModel struct {
	Id                                     types.String `tfsdk:"id"`
	LastUpdated                            types.String `tfsdk:"last_updated"`
	Notifications                          types.Set    `tfsdk:"notifications"`
	RequiredActions                        types.Set    `tfsdk:"required_actions"`
	OperationType                          types.Set    `tfsdk:"operation_type"`
	OperationOrigin                        types.Set    `tfsdk:"operation_origin"`
	ConnectionCriteria                     types.String `tfsdk:"connection_criteria"`
	AllIncludedRequestControl              types.Set    `tfsdk:"all_included_request_control"`
	AnyIncludedRequestControl              types.Set    `tfsdk:"any_included_request_control"`
	NotAllIncludedRequestControl           types.Set    `tfsdk:"not_all_included_request_control"`
	NoneIncludedRequestControl             types.Set    `tfsdk:"none_included_request_control"`
	IncludedTargetEntryDN                  types.Set    `tfsdk:"included_target_entry_dn"`
	ExcludedTargetEntryDN                  types.Set    `tfsdk:"excluded_target_entry_dn"`
	AllIncludedTargetEntryFilter           types.Set    `tfsdk:"all_included_target_entry_filter"`
	AnyIncludedTargetEntryFilter           types.Set    `tfsdk:"any_included_target_entry_filter"`
	NotAllIncludedTargetEntryFilter        types.Set    `tfsdk:"not_all_included_target_entry_filter"`
	NoneIncludedTargetEntryFilter          types.Set    `tfsdk:"none_included_target_entry_filter"`
	AllIncludedTargetEntryGroupDN          types.Set    `tfsdk:"all_included_target_entry_group_dn"`
	AnyIncludedTargetEntryGroupDN          types.Set    `tfsdk:"any_included_target_entry_group_dn"`
	NotAllIncludedTargetEntryGroupDN       types.Set    `tfsdk:"not_all_included_target_entry_group_dn"`
	NoneIncludedTargetEntryGroupDN         types.Set    `tfsdk:"none_included_target_entry_group_dn"`
	TargetBindType                         types.Set    `tfsdk:"target_bind_type"`
	IncludedTargetSASLMechanism            types.Set    `tfsdk:"included_target_sasl_mechanism"`
	ExcludedTargetSASLMechanism            types.Set    `tfsdk:"excluded_target_sasl_mechanism"`
	IncludedTargetAttribute                types.Set    `tfsdk:"included_target_attribute"`
	ExcludedTargetAttribute                types.Set    `tfsdk:"excluded_target_attribute"`
	IncludedExtendedOperationOID           types.Set    `tfsdk:"included_extended_operation_oid"`
	ExcludedExtendedOperationOID           types.Set    `tfsdk:"excluded_extended_operation_oid"`
	IncludedSearchScope                    types.Set    `tfsdk:"included_search_scope"`
	UsingAdministrativeSessionWorkerThread types.String `tfsdk:"using_administrative_session_worker_thread"`
	IncludedApplicationName                types.Set    `tfsdk:"included_application_name"`
	ExcludedApplicationName                types.Set    `tfsdk:"excluded_application_name"`
	Description                            types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *simpleRequestCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	simpleRequestCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultSimpleRequestCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	simpleRequestCriteriaSchema(ctx, req, resp, true)
}

func simpleRequestCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Simple Request Criteria.",
		Attributes: map[string]schema.Attribute{
			"operation_type": schema.SetAttribute{
				Description: "Specifies the operation type(s) for operations that should be included in this Simple Request Criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"operation_origin": schema.SetAttribute{
				Description: "Specifies the origin for operations to be included in this Simple Request Criteria. If no values are provided, then the operation origin will not be taken into consideration when determining whether an operation matches this Simple Request Criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"connection_criteria": schema.StringAttribute{
				Description: "Specifies a connection criteria object that must match the associated client connection for operations included in this Simple Request Criteria.",
				Optional:    true,
			},
			"all_included_request_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must be present in the request from the client for operations included in this Simple Request Criteria. If any control OIDs are provided, then the request must contain all of those controls.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_request_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that may be present in the request from the client for operations included in this Simple Request Criteria. If any control OIDs are provided, then the request must contain at least one of those controls.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_request_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that should not be present in the request from the client for operations included in this Simple Request Criteria. If any control OIDs are provided, then the request must not contain at least one of those controls (that is, the request may contain zero or more of those controls, but not all of them).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_request_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must not be present in the request from the client for operations included in this Simple Request Criteria. If any control OIDs are provided, then the request must not contain any of those controls.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_target_entry_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which targeted entries may exist for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_target_entry_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which targeted entries may not exist for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"all_included_target_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that must match the target entry for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any filters are provided, then the target entry must match all of those filters.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_target_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that may match the target entry for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any filters are provided, then the target entry must match at least one of those filters.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_target_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that should not match the target entry for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any filters are provided, then the target entry must not match at least one of those filters (that is, the request may match zero or more of those filters, but not of all of them).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_target_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that must not match the target entry for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any filters are provided, then the target entry must not match any of those filters.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"all_included_target_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the target entry must be a member for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any group DNs are provided, then the target entry must be a member of all of those groups.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_target_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the target entry may be a member for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any group DNs are provided, then the target entry must be a member of at least one of those groups.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_target_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the target entry should not be a member for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any group DNs are provided, then the target entry must not be a member of at least one of those groups (that is, the target entry may be a member of zero or more of those groups, but not all of them).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_target_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the target entry must not be a member for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any group DNs are provided, then the target entry must not be a member of any of those groups.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"target_bind_type": schema.SetAttribute{
				Description: "Specifies the authentication type for bind requests included in this Simple Request Criteria. This will only be taken into account for bind operations and will be ignored for any other type of operation. If no values are provided, then the authentication type will not be considered when determining whether the request should be included in this Simple Request Criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_target_sasl_mechanism": schema.SetAttribute{
				Description: "Specifies the name of a SASL mechanism for bind requests included in this Simple Request Criteria. This will only be taken into account for SASL bind operations and will be ignored for other types of operations and for bind operations that do not use SASL authentication.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_target_sasl_mechanism": schema.SetAttribute{
				Description: "Specifies the name of a SASL mechanism for bind requests excluded from this Simple Request Criteria. This will only be taken into account for SASL bind operations and will be ignored for other types of operations and for bind operations that do not use SASL authentication.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_target_attribute": schema.SetAttribute{
				Description: "Specifies the name or OID of an attribute type which must be targeted by requests included in this Simple Request Criteria. This will only be taken into account for add, compare, modify, modify DN, and search operations. It will be ignored for abandon, bind, delete, extended, and unbind operations.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_target_attribute": schema.SetAttribute{
				Description: "Specifies the name or OID of an attribute type which must not be targeted by requests included in this Simple Request Criteria. This will only be taken into account for add, compare, modify, modify DN, and search operations. It will be ignored for abandon, bind, delete, extended, and unbind operations.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_extended_operation_oid": schema.SetAttribute{
				Description: "Specifies the request OID for extended requests included in this Simple Request Criteria. This will only be taken into account for extended requests and will be ignored for all other types of requests.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_extended_operation_oid": schema.SetAttribute{
				Description: "Specifies the request OID for extended requests excluded from this Simple Request Criteria. This will only be taken into account for extended requests and will be ignored for all other types of requests.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_search_scope": schema.SetAttribute{
				Description: "Specifies the search scope values included in this Simple Request Criteria. This will only be taken into account for search requests and will be ignored for all other types of requests.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"using_administrative_session_worker_thread": schema.StringAttribute{
				Description: "Indicates whether operations being processed using a dedicated administrative session worker thread should be included in this Simple Request Criteria.",
				Optional:    true,
				Computed:    true,
			},
			"included_application_name": schema.SetAttribute{
				Description: "Specifies an application name for requests included in this Simple Request Criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_application_name": schema.SetAttribute{
				Description: "Specifies an application name for requests excluded from this Simple Request Criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Request Criteria",
				Optional:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalSimpleRequestCriteriaFields(ctx context.Context, addRequest *client.AddSimpleRequestCriteriaRequest, plan simpleRequestCriteriaResourceModel) error {
	if internaltypes.IsDefined(plan.OperationType) {
		var slice []string
		plan.OperationType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumrequestCriteriaSimpleOperationTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumrequestCriteriaSimpleOperationTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.OperationType = enumSlice
	}
	if internaltypes.IsDefined(plan.OperationOrigin) {
		var slice []string
		plan.OperationOrigin.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumrequestCriteriaOperationOriginProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumrequestCriteriaOperationOriginPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.OperationOrigin = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		stringVal := plan.ConnectionCriteria.ValueString()
		addRequest.ConnectionCriteria = &stringVal
	}
	if internaltypes.IsDefined(plan.AllIncludedRequestControl) {
		var slice []string
		plan.AllIncludedRequestControl.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedRequestControl = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedRequestControl) {
		var slice []string
		plan.AnyIncludedRequestControl.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedRequestControl = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedRequestControl) {
		var slice []string
		plan.NotAllIncludedRequestControl.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedRequestControl = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedRequestControl) {
		var slice []string
		plan.NoneIncludedRequestControl.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedRequestControl = slice
	}
	if internaltypes.IsDefined(plan.IncludedTargetEntryDN) {
		var slice []string
		plan.IncludedTargetEntryDN.ElementsAs(ctx, &slice, false)
		addRequest.IncludedTargetEntryDN = slice
	}
	if internaltypes.IsDefined(plan.ExcludedTargetEntryDN) {
		var slice []string
		plan.ExcludedTargetEntryDN.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedTargetEntryDN = slice
	}
	if internaltypes.IsDefined(plan.AllIncludedTargetEntryFilter) {
		var slice []string
		plan.AllIncludedTargetEntryFilter.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedTargetEntryFilter = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedTargetEntryFilter) {
		var slice []string
		plan.AnyIncludedTargetEntryFilter.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedTargetEntryFilter = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedTargetEntryFilter) {
		var slice []string
		plan.NotAllIncludedTargetEntryFilter.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedTargetEntryFilter = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedTargetEntryFilter) {
		var slice []string
		plan.NoneIncludedTargetEntryFilter.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedTargetEntryFilter = slice
	}
	if internaltypes.IsDefined(plan.AllIncludedTargetEntryGroupDN) {
		var slice []string
		plan.AllIncludedTargetEntryGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedTargetEntryGroupDN = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedTargetEntryGroupDN) {
		var slice []string
		plan.AnyIncludedTargetEntryGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedTargetEntryGroupDN = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedTargetEntryGroupDN) {
		var slice []string
		plan.NotAllIncludedTargetEntryGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedTargetEntryGroupDN = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedTargetEntryGroupDN) {
		var slice []string
		plan.NoneIncludedTargetEntryGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedTargetEntryGroupDN = slice
	}
	if internaltypes.IsDefined(plan.TargetBindType) {
		var slice []string
		plan.TargetBindType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumrequestCriteriaTargetBindTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumrequestCriteriaTargetBindTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.TargetBindType = enumSlice
	}
	if internaltypes.IsDefined(plan.IncludedTargetSASLMechanism) {
		var slice []string
		plan.IncludedTargetSASLMechanism.ElementsAs(ctx, &slice, false)
		addRequest.IncludedTargetSASLMechanism = slice
	}
	if internaltypes.IsDefined(plan.ExcludedTargetSASLMechanism) {
		var slice []string
		plan.ExcludedTargetSASLMechanism.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedTargetSASLMechanism = slice
	}
	if internaltypes.IsDefined(plan.IncludedTargetAttribute) {
		var slice []string
		plan.IncludedTargetAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludedTargetAttribute = slice
	}
	if internaltypes.IsDefined(plan.ExcludedTargetAttribute) {
		var slice []string
		plan.ExcludedTargetAttribute.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedTargetAttribute = slice
	}
	if internaltypes.IsDefined(plan.IncludedExtendedOperationOID) {
		var slice []string
		plan.IncludedExtendedOperationOID.ElementsAs(ctx, &slice, false)
		addRequest.IncludedExtendedOperationOID = slice
	}
	if internaltypes.IsDefined(plan.ExcludedExtendedOperationOID) {
		var slice []string
		plan.ExcludedExtendedOperationOID.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedExtendedOperationOID = slice
	}
	if internaltypes.IsDefined(plan.IncludedSearchScope) {
		var slice []string
		plan.IncludedSearchScope.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumrequestCriteriaIncludedSearchScopeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumrequestCriteriaIncludedSearchScopePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.IncludedSearchScope = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UsingAdministrativeSessionWorkerThread) {
		usingAdministrativeSessionWorkerThread, err := client.NewEnumrequestCriteriaUsingAdministrativeSessionWorkerThreadPropFromValue(plan.UsingAdministrativeSessionWorkerThread.ValueString())
		if err != nil {
			return err
		}
		addRequest.UsingAdministrativeSessionWorkerThread = usingAdministrativeSessionWorkerThread
	}
	if internaltypes.IsDefined(plan.IncludedApplicationName) {
		var slice []string
		plan.IncludedApplicationName.ElementsAs(ctx, &slice, false)
		addRequest.IncludedApplicationName = slice
	}
	if internaltypes.IsDefined(plan.ExcludedApplicationName) {
		var slice []string
		plan.ExcludedApplicationName.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedApplicationName = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	return nil
}

// Read a SimpleRequestCriteriaResponse object into the model struct
func readSimpleRequestCriteriaResponse(ctx context.Context, r *client.SimpleRequestCriteriaResponse, state *simpleRequestCriteriaResourceModel, expectedValues *simpleRequestCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.OperationType = internaltypes.GetStringSet(
		client.StringSliceEnumrequestCriteriaSimpleOperationTypeProp(r.OperationType))
	state.OperationOrigin = internaltypes.GetStringSet(
		client.StringSliceEnumrequestCriteriaOperationOriginProp(r.OperationOrigin))
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.AllIncludedRequestControl = internaltypes.GetStringSet(r.AllIncludedRequestControl)
	state.AnyIncludedRequestControl = internaltypes.GetStringSet(r.AnyIncludedRequestControl)
	state.NotAllIncludedRequestControl = internaltypes.GetStringSet(r.NotAllIncludedRequestControl)
	state.NoneIncludedRequestControl = internaltypes.GetStringSet(r.NoneIncludedRequestControl)
	state.IncludedTargetEntryDN = internaltypes.GetStringSet(r.IncludedTargetEntryDN)
	state.ExcludedTargetEntryDN = internaltypes.GetStringSet(r.ExcludedTargetEntryDN)
	state.AllIncludedTargetEntryFilter = internaltypes.GetStringSet(r.AllIncludedTargetEntryFilter)
	state.AnyIncludedTargetEntryFilter = internaltypes.GetStringSet(r.AnyIncludedTargetEntryFilter)
	state.NotAllIncludedTargetEntryFilter = internaltypes.GetStringSet(r.NotAllIncludedTargetEntryFilter)
	state.NoneIncludedTargetEntryFilter = internaltypes.GetStringSet(r.NoneIncludedTargetEntryFilter)
	state.AllIncludedTargetEntryGroupDN = internaltypes.GetStringSet(r.AllIncludedTargetEntryGroupDN)
	state.AnyIncludedTargetEntryGroupDN = internaltypes.GetStringSet(r.AnyIncludedTargetEntryGroupDN)
	state.NotAllIncludedTargetEntryGroupDN = internaltypes.GetStringSet(r.NotAllIncludedTargetEntryGroupDN)
	state.NoneIncludedTargetEntryGroupDN = internaltypes.GetStringSet(r.NoneIncludedTargetEntryGroupDN)
	state.TargetBindType = internaltypes.GetStringSet(
		client.StringSliceEnumrequestCriteriaTargetBindTypeProp(r.TargetBindType))
	state.IncludedTargetSASLMechanism = internaltypes.GetStringSet(r.IncludedTargetSASLMechanism)
	state.ExcludedTargetSASLMechanism = internaltypes.GetStringSet(r.ExcludedTargetSASLMechanism)
	state.IncludedTargetAttribute = internaltypes.GetStringSet(r.IncludedTargetAttribute)
	state.ExcludedTargetAttribute = internaltypes.GetStringSet(r.ExcludedTargetAttribute)
	state.IncludedExtendedOperationOID = internaltypes.GetStringSet(r.IncludedExtendedOperationOID)
	state.ExcludedExtendedOperationOID = internaltypes.GetStringSet(r.ExcludedExtendedOperationOID)
	state.IncludedSearchScope = internaltypes.GetStringSet(
		client.StringSliceEnumrequestCriteriaIncludedSearchScopeProp(r.IncludedSearchScope))
	state.UsingAdministrativeSessionWorkerThread = internaltypes.StringTypeOrNil(
		client.StringPointerEnumrequestCriteriaUsingAdministrativeSessionWorkerThreadProp(r.UsingAdministrativeSessionWorkerThread), internaltypes.IsEmptyString(expectedValues.UsingAdministrativeSessionWorkerThread))
	state.IncludedApplicationName = internaltypes.GetStringSet(r.IncludedApplicationName)
	state.ExcludedApplicationName = internaltypes.GetStringSet(r.ExcludedApplicationName)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSimpleRequestCriteriaOperations(plan simpleRequestCriteriaResourceModel, state simpleRequestCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.OperationType, state.OperationType, "operation-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.OperationOrigin, state.OperationOrigin, "operation-origin")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionCriteria, state.ConnectionCriteria, "connection-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedRequestControl, state.AllIncludedRequestControl, "all-included-request-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedRequestControl, state.AnyIncludedRequestControl, "any-included-request-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedRequestControl, state.NotAllIncludedRequestControl, "not-all-included-request-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedRequestControl, state.NoneIncludedRequestControl, "none-included-request-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedTargetEntryDN, state.IncludedTargetEntryDN, "included-target-entry-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedTargetEntryDN, state.ExcludedTargetEntryDN, "excluded-target-entry-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedTargetEntryFilter, state.AllIncludedTargetEntryFilter, "all-included-target-entry-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedTargetEntryFilter, state.AnyIncludedTargetEntryFilter, "any-included-target-entry-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedTargetEntryFilter, state.NotAllIncludedTargetEntryFilter, "not-all-included-target-entry-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedTargetEntryFilter, state.NoneIncludedTargetEntryFilter, "none-included-target-entry-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedTargetEntryGroupDN, state.AllIncludedTargetEntryGroupDN, "all-included-target-entry-group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedTargetEntryGroupDN, state.AnyIncludedTargetEntryGroupDN, "any-included-target-entry-group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedTargetEntryGroupDN, state.NotAllIncludedTargetEntryGroupDN, "not-all-included-target-entry-group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedTargetEntryGroupDN, state.NoneIncludedTargetEntryGroupDN, "none-included-target-entry-group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.TargetBindType, state.TargetBindType, "target-bind-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedTargetSASLMechanism, state.IncludedTargetSASLMechanism, "included-target-sasl-mechanism")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedTargetSASLMechanism, state.ExcludedTargetSASLMechanism, "excluded-target-sasl-mechanism")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedTargetAttribute, state.IncludedTargetAttribute, "included-target-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedTargetAttribute, state.ExcludedTargetAttribute, "excluded-target-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedExtendedOperationOID, state.IncludedExtendedOperationOID, "included-extended-operation-oid")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedExtendedOperationOID, state.ExcludedExtendedOperationOID, "excluded-extended-operation-oid")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedSearchScope, state.IncludedSearchScope, "included-search-scope")
	operations.AddStringOperationIfNecessary(&ops, plan.UsingAdministrativeSessionWorkerThread, state.UsingAdministrativeSessionWorkerThread, "using-administrative-session-worker-thread")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedApplicationName, state.IncludedApplicationName, "included-application-name")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedApplicationName, state.ExcludedApplicationName, "excluded-application-name")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *simpleRequestCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan simpleRequestCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddSimpleRequestCriteriaRequest(plan.Id.ValueString(),
		[]client.EnumsimpleRequestCriteriaSchemaUrn{client.ENUMSIMPLEREQUESTCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0REQUEST_CRITERIASIMPLE})
	err := addOptionalSimpleRequestCriteriaFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Simple Request Criteria", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RequestCriteriaApi.AddRequestCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRequestCriteriaRequest(
		client.AddSimpleRequestCriteriaRequestAsAddRequestCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RequestCriteriaApi.AddRequestCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Simple Request Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state simpleRequestCriteriaResourceModel
	readSimpleRequestCriteriaResponse(ctx, addResponse.SimpleRequestCriteriaResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultSimpleRequestCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan simpleRequestCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RequestCriteriaApi.GetRequestCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Simple Request Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state simpleRequestCriteriaResourceModel
	readSimpleRequestCriteriaResponse(ctx, readResponse.SimpleRequestCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.RequestCriteriaApi.UpdateRequestCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSimpleRequestCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.RequestCriteriaApi.UpdateRequestCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Simple Request Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSimpleRequestCriteriaResponse(ctx, updateResponse.SimpleRequestCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *simpleRequestCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSimpleRequestCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSimpleRequestCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSimpleRequestCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSimpleRequestCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state simpleRequestCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.RequestCriteriaApi.GetRequestCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Simple Request Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSimpleRequestCriteriaResponse(ctx, readResponse.SimpleRequestCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *simpleRequestCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSimpleRequestCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSimpleRequestCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSimpleRequestCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSimpleRequestCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan simpleRequestCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state simpleRequestCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.RequestCriteriaApi.UpdateRequestCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSimpleRequestCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.RequestCriteriaApi.UpdateRequestCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Simple Request Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSimpleRequestCriteriaResponse(ctx, updateResponse.SimpleRequestCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSimpleRequestCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *simpleRequestCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state simpleRequestCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.RequestCriteriaApi.DeleteRequestCriteriaExecute(r.apiClient.RequestCriteriaApi.DeleteRequestCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Simple Request Criteria", err, httpResp)
		return
	}
}

func (r *simpleRequestCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSimpleRequestCriteria(ctx, req, resp)
}

func (r *defaultSimpleRequestCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSimpleRequestCriteria(ctx, req, resp)
}

func importSimpleRequestCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
