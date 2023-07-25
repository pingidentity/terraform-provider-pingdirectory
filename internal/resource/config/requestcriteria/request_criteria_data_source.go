package requestcriteria

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
	_ datasource.DataSource              = &requestCriteriaDataSource{}
	_ datasource.DataSourceWithConfigure = &requestCriteriaDataSource{}
)

// Create a Request Criteria data source
func NewRequestCriteriaDataSource() datasource.DataSource {
	return &requestCriteriaDataSource{}
}

// requestCriteriaDataSource is the datasource implementation.
type requestCriteriaDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *requestCriteriaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_criteria"
}

// Configure adds the provider configured client to the data source.
func (r *requestCriteriaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type requestCriteriaDataSourceModel struct {
	Id                                     types.String `tfsdk:"id"`
	Type                                   types.String `tfsdk:"type"`
	ExtensionClass                         types.String `tfsdk:"extension_class"`
	ExtensionArgument                      types.Set    `tfsdk:"extension_argument"`
	AllIncludedRequestCriteria             types.Set    `tfsdk:"all_included_request_criteria"`
	AnyIncludedRequestCriteria             types.Set    `tfsdk:"any_included_request_criteria"`
	NotAllIncludedRequestCriteria          types.Set    `tfsdk:"not_all_included_request_criteria"`
	NoneIncludedRequestCriteria            types.Set    `tfsdk:"none_included_request_criteria"`
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

// GetSchema defines the schema for the datasource.
func (r *requestCriteriaDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Describes a Request Criteria.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Name of this object.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of Request Criteria resource. Options are ['root-dse', 'simple', 'aggregate', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Request Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Request Criteria. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"all_included_request_criteria": schema.SetAttribute{
				Description: "Specifies a request criteria object that must match the associated operation request in order to match the aggregate request criteria. If one or more all-included request criteria objects are provided, then an operation request must match all of them in order to match the aggregate request criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_request_criteria": schema.SetAttribute{
				Description: "Specifies a request criteria object that may match the associated operation request in order to the this aggregate request criteria. If one or more any-included request criteria objects are provided, then an operation request must match at least one of them in order to match the aggregate request criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_request_criteria": schema.SetAttribute{
				Description: "Specifies a request criteria object that should not match the associated operation request in order to match the aggregate request criteria. If one or more not-all-included request criteria objects are provided, then an operation request must not match all of them (that is, it may match zero or more of them, but it must not match all of them) in order to match the aggregate request criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_request_criteria": schema.SetAttribute{
				Description: "Specifies a request criteria object that must not match the associated operation request in order to match the aggregate request criteria. If one or more none-included request criteria objects are provided, then an operation request must not match any of them in order to match the aggregate request criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"operation_type": schema.SetAttribute{
				Description: "The types of operations that may be matched by this Root DSE Request Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"operation_origin": schema.SetAttribute{
				Description: "Specifies the origin for operations to be included in this Simple Request Criteria. If no values are provided, then the operation origin will not be taken into consideration when determining whether an operation matches this Simple Request Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"connection_criteria": schema.StringAttribute{
				Description: "Specifies a connection criteria object that must match the associated client connection for operations included in this Simple Request Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"all_included_request_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must be present in the request from the client for operations included in this Simple Request Criteria. If any control OIDs are provided, then the request must contain all of those controls.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_request_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that may be present in the request from the client for operations included in this Simple Request Criteria. If any control OIDs are provided, then the request must contain at least one of those controls.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_request_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that should not be present in the request from the client for operations included in this Simple Request Criteria. If any control OIDs are provided, then the request must not contain at least one of those controls (that is, the request may contain zero or more of those controls, but not all of them).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_request_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must not be present in the request from the client for operations included in this Simple Request Criteria. If any control OIDs are provided, then the request must not contain any of those controls.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_target_entry_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which targeted entries may exist for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_target_entry_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which targeted entries may not exist for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"all_included_target_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that must match the target entry for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any filters are provided, then the target entry must match all of those filters.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_target_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that may match the target entry for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any filters are provided, then the target entry must match at least one of those filters.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_target_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that should not match the target entry for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any filters are provided, then the target entry must not match at least one of those filters (that is, the request may match zero or more of those filters, but not of all of them).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_target_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that must not match the target entry for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any filters are provided, then the target entry must not match any of those filters.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"all_included_target_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the target entry must be a member for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any group DNs are provided, then the target entry must be a member of all of those groups.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_target_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the target entry may be a member for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any group DNs are provided, then the target entry must be a member of at least one of those groups.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_target_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the target entry should not be a member for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any group DNs are provided, then the target entry must not be a member of at least one of those groups (that is, the target entry may be a member of zero or more of those groups, but not all of them).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_target_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the target entry must not be a member for requests included in this Simple Request Criteria. This will only be taken into account for add, simple bind, compare, delete, modify, modify DN, and search operations. It will be ignored for abandon, SASL bind, extended, and unbind operations. If any group DNs are provided, then the target entry must not be a member of any of those groups.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"target_bind_type": schema.SetAttribute{
				Description: "Specifies the authentication type for bind requests included in this Simple Request Criteria. This will only be taken into account for bind operations and will be ignored for any other type of operation. If no values are provided, then the authentication type will not be considered when determining whether the request should be included in this Simple Request Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_target_sasl_mechanism": schema.SetAttribute{
				Description: "Specifies the name of a SASL mechanism for bind requests included in this Simple Request Criteria. This will only be taken into account for SASL bind operations and will be ignored for other types of operations and for bind operations that do not use SASL authentication.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_target_sasl_mechanism": schema.SetAttribute{
				Description: "Specifies the name of a SASL mechanism for bind requests excluded from this Simple Request Criteria. This will only be taken into account for SASL bind operations and will be ignored for other types of operations and for bind operations that do not use SASL authentication.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_target_attribute": schema.SetAttribute{
				Description: "Specifies the name or OID of an attribute type which must be targeted by requests included in this Simple Request Criteria. This will only be taken into account for add, compare, modify, modify DN, and search operations. It will be ignored for abandon, bind, delete, extended, and unbind operations.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_target_attribute": schema.SetAttribute{
				Description: "Specifies the name or OID of an attribute type which must not be targeted by requests included in this Simple Request Criteria. This will only be taken into account for add, compare, modify, modify DN, and search operations. It will be ignored for abandon, bind, delete, extended, and unbind operations.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_extended_operation_oid": schema.SetAttribute{
				Description: "Specifies the request OID for extended requests included in this Simple Request Criteria. This will only be taken into account for extended requests and will be ignored for all other types of requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_extended_operation_oid": schema.SetAttribute{
				Description: "Specifies the request OID for extended requests excluded from this Simple Request Criteria. This will only be taken into account for extended requests and will be ignored for all other types of requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_search_scope": schema.SetAttribute{
				Description: "Specifies the search scope values included in this Simple Request Criteria. This will only be taken into account for search requests and will be ignored for all other types of requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"using_administrative_session_worker_thread": schema.StringAttribute{
				Description: "Indicates whether operations being processed using a dedicated administrative session worker thread should be included in this Simple Request Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"included_application_name": schema.SetAttribute{
				Description: "Specifies an application name for requests included in this Simple Request Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_application_name": schema.SetAttribute{
				Description: "Specifies an application name for requests excluded from this Simple Request Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Request Criteria",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
}

// Read a RootDseRequestCriteriaResponse object into the model struct
func readRootDseRequestCriteriaResponseDataSource(ctx context.Context, r *client.RootDseRequestCriteriaResponse, state *requestCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("root-dse")
	state.Id = types.StringValue(r.Id)
	state.OperationType = internaltypes.GetStringSet(
		client.StringSliceEnumrequestCriteriaRootDseOperationTypeProp(r.OperationType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a SimpleRequestCriteriaResponse object into the model struct
func readSimpleRequestCriteriaResponseDataSource(ctx context.Context, r *client.SimpleRequestCriteriaResponse, state *requestCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("simple")
	state.Id = types.StringValue(r.Id)
	state.OperationType = internaltypes.GetStringSet(
		client.StringSliceEnumrequestCriteriaSimpleOperationTypeProp(r.OperationType))
	state.OperationOrigin = internaltypes.GetStringSet(
		client.StringSliceEnumrequestCriteriaOperationOriginProp(r.OperationOrigin))
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
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
		client.StringPointerEnumrequestCriteriaUsingAdministrativeSessionWorkerThreadProp(r.UsingAdministrativeSessionWorkerThread), false)
	state.IncludedApplicationName = internaltypes.GetStringSet(r.IncludedApplicationName)
	state.ExcludedApplicationName = internaltypes.GetStringSet(r.ExcludedApplicationName)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a AggregateRequestCriteriaResponse object into the model struct
func readAggregateRequestCriteriaResponseDataSource(ctx context.Context, r *client.AggregateRequestCriteriaResponse, state *requestCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("aggregate")
	state.Id = types.StringValue(r.Id)
	state.AllIncludedRequestCriteria = internaltypes.GetStringSet(r.AllIncludedRequestCriteria)
	state.AnyIncludedRequestCriteria = internaltypes.GetStringSet(r.AnyIncludedRequestCriteria)
	state.NotAllIncludedRequestCriteria = internaltypes.GetStringSet(r.NotAllIncludedRequestCriteria)
	state.NoneIncludedRequestCriteria = internaltypes.GetStringSet(r.NoneIncludedRequestCriteria)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a ThirdPartyRequestCriteriaResponse object into the model struct
func readThirdPartyRequestCriteriaResponseDataSource(ctx context.Context, r *client.ThirdPartyRequestCriteriaResponse, state *requestCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read resource information
func (r *requestCriteriaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state requestCriteriaDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RequestCriteriaApi.GetRequestCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Request Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.RootDseRequestCriteriaResponse != nil {
		readRootDseRequestCriteriaResponseDataSource(ctx, readResponse.RootDseRequestCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SimpleRequestCriteriaResponse != nil {
		readSimpleRequestCriteriaResponseDataSource(ctx, readResponse.SimpleRequestCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AggregateRequestCriteriaResponse != nil {
		readAggregateRequestCriteriaResponseDataSource(ctx, readResponse.AggregateRequestCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyRequestCriteriaResponse != nil {
		readThirdPartyRequestCriteriaResponseDataSource(ctx, readResponse.ThirdPartyRequestCriteriaResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
