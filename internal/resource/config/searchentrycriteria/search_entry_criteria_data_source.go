package searchentrycriteria

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
	_ datasource.DataSource              = &searchEntryCriteriaDataSource{}
	_ datasource.DataSourceWithConfigure = &searchEntryCriteriaDataSource{}
)

// Create a Search Entry Criteria data source
func NewSearchEntryCriteriaDataSource() datasource.DataSource {
	return &searchEntryCriteriaDataSource{}
}

// searchEntryCriteriaDataSource is the datasource implementation.
type searchEntryCriteriaDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *searchEntryCriteriaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_search_entry_criteria"
}

// Configure adds the provider configured client to the data source.
func (r *searchEntryCriteriaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type searchEntryCriteriaDataSourceModel struct {
	Id                                types.String `tfsdk:"id"`
	Type                              types.String `tfsdk:"type"`
	ExtensionClass                    types.String `tfsdk:"extension_class"`
	ExtensionArgument                 types.Set    `tfsdk:"extension_argument"`
	AllIncludedSearchEntryCriteria    types.Set    `tfsdk:"all_included_search_entry_criteria"`
	AnyIncludedSearchEntryCriteria    types.Set    `tfsdk:"any_included_search_entry_criteria"`
	NotAllIncludedSearchEntryCriteria types.Set    `tfsdk:"not_all_included_search_entry_criteria"`
	NoneIncludedSearchEntryCriteria   types.Set    `tfsdk:"none_included_search_entry_criteria"`
	RequestCriteria                   types.String `tfsdk:"request_criteria"`
	AllIncludedEntryControl           types.Set    `tfsdk:"all_included_entry_control"`
	AnyIncludedEntryControl           types.Set    `tfsdk:"any_included_entry_control"`
	NotAllIncludedEntryControl        types.Set    `tfsdk:"not_all_included_entry_control"`
	NoneIncludedEntryControl          types.Set    `tfsdk:"none_included_entry_control"`
	IncludedEntryBaseDN               types.Set    `tfsdk:"included_entry_base_dn"`
	ExcludedEntryBaseDN               types.Set    `tfsdk:"excluded_entry_base_dn"`
	AllIncludedEntryFilter            types.Set    `tfsdk:"all_included_entry_filter"`
	AnyIncludedEntryFilter            types.Set    `tfsdk:"any_included_entry_filter"`
	NotAllIncludedEntryFilter         types.Set    `tfsdk:"not_all_included_entry_filter"`
	NoneIncludedEntryFilter           types.Set    `tfsdk:"none_included_entry_filter"`
	AllIncludedEntryGroupDN           types.Set    `tfsdk:"all_included_entry_group_dn"`
	AnyIncludedEntryGroupDN           types.Set    `tfsdk:"any_included_entry_group_dn"`
	NotAllIncludedEntryGroupDN        types.Set    `tfsdk:"not_all_included_entry_group_dn"`
	NoneIncludedEntryGroupDN          types.Set    `tfsdk:"none_included_entry_group_dn"`
	Description                       types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the datasource.
func (r *searchEntryCriteriaDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Describes a Search Entry Criteria.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Name of this object.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of Search Entry Criteria resource. Options are ['simple', 'aggregate', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Search Entry Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Search Entry Criteria. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"all_included_search_entry_criteria": schema.SetAttribute{
				Description: "Specifies a search entry criteria object that must match the associated search result entry in order to match the aggregate search entry criteria. If one or more all-included search entry criteria objects are provided, then a search result entry must match all of them in order to match the aggregate search entry criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_search_entry_criteria": schema.SetAttribute{
				Description: "Specifies a search entry criteria object that may match the associated search result entry in order to match the aggregate search entry criteria. If one or more any-included search entry criteria objects are provided, then a search result entry must match at least one of them in order to match the aggregate search entry criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_search_entry_criteria": schema.SetAttribute{
				Description: "Specifies a search entry criteria object that should not match the associated search result entry in order to match the aggregate search entry criteria. If one or more not-all-included search entry criteria objects are provided, then a search result entry must not match all of them (that is, it may match zero or more of them, but it must not match all of them) in order to match the aggregate search entry criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_search_entry_criteria": schema.SetAttribute{
				Description: "Specifies a search entry criteria object that must not match the associated search result entry in order to match the aggregate search entry criteria. If one or more none-included search entry criteria objects are provided, then a search result entry must not match any of them in order to match the aggregate search entry criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"request_criteria": schema.StringAttribute{
				Description: "Specifies a request criteria object that must match the associated request for entries included in this Simple Search Entry Criteria. of them.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"all_included_entry_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must be present in search result entries included in this Simple Search Entry Criteria. If any control OIDs are provided, then the entry must contain all of those controls.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_entry_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that may be present in search result entries included in this Simple Search Entry Criteria. If any control OIDs are provided, then the entry must contain at least one of those controls.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_entry_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that should not be present in search result entries included in this Simple Search Entry Criteria. If any control OIDs are provided, then the entry must not contain at least one of those controls (that is, it may contain zero or more of those controls, but not all of them).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_entry_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must not be present in search result entries included in this Simple Search Entry Criteria. If any control OIDs are provided, then the entry must not contain any of those controls.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_entry_base_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which entries included in this Simple Search Entry Criteria may exist.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_entry_base_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which entries included in this Simple Search Entry Criteria may not exist.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"all_included_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that must match search result entries included in this Simple Search Entry Criteria. Note that this matching will be performed against the entry that is actually returned to the client and may not reflect the complete entry stored in the server. If any filters are provided, then the returned entry must match all of those filters.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that may match search result entries included in this Simple Search Entry Criteria. Note that this matching will be performed against the entry that is actually returned to the client and may not reflect the complete entry stored in the server. If any filters are provided, then the entry must match at least one of those filters.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that should not match search result entries included in this Simple Search Entry Criteria. Note that this matching will be performed against the entry that is actually returned to the client and may not reflect the complete entry stored in the server. If any filters are provided, then the entry must not match at least one of those filters (that is, the entry may match zero or more of those filters, but not of all of them).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that must not match search result entries included in this Simple Search Entry Criteria. Note that this matching will be performed against the entry that is actually returned to the client and may not reflect the complete entry stored in the server. If any filters are provided, then the entry must not match any of those filters.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"all_included_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the entry must be a member to be included in this Simple Search Entry Criteria. If any group DNs are provided, then the entry must be a member of all of them.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the entry may be a member to be included in this Simple Search Entry Criteria. If any group DNs are provided, then the entry must be a member of at least one of them.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the entry should not be a member to be included in this Simple Search Entry Criteria. If any group DNs are provided, then the entry must not be a member of at least one of them (that is, the entry may be a member of zero or more of the specified groups, but not of all of them).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the entry must not be a member to be included in this Simple Search Entry Criteria. If any group DNs are provided, then the entry must not be a member of any of them.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Search Entry Criteria",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
}

// Read a SimpleSearchEntryCriteriaResponse object into the model struct
func readSimpleSearchEntryCriteriaResponseDataSource(ctx context.Context, r *client.SimpleSearchEntryCriteriaResponse, state *searchEntryCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("simple")
	state.Id = types.StringValue(r.Id)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.AllIncludedEntryControl = internaltypes.GetStringSet(r.AllIncludedEntryControl)
	state.AnyIncludedEntryControl = internaltypes.GetStringSet(r.AnyIncludedEntryControl)
	state.NotAllIncludedEntryControl = internaltypes.GetStringSet(r.NotAllIncludedEntryControl)
	state.NoneIncludedEntryControl = internaltypes.GetStringSet(r.NoneIncludedEntryControl)
	state.IncludedEntryBaseDN = internaltypes.GetStringSet(r.IncludedEntryBaseDN)
	state.ExcludedEntryBaseDN = internaltypes.GetStringSet(r.ExcludedEntryBaseDN)
	state.AllIncludedEntryFilter = internaltypes.GetStringSet(r.AllIncludedEntryFilter)
	state.AnyIncludedEntryFilter = internaltypes.GetStringSet(r.AnyIncludedEntryFilter)
	state.NotAllIncludedEntryFilter = internaltypes.GetStringSet(r.NotAllIncludedEntryFilter)
	state.NoneIncludedEntryFilter = internaltypes.GetStringSet(r.NoneIncludedEntryFilter)
	state.AllIncludedEntryGroupDN = internaltypes.GetStringSet(r.AllIncludedEntryGroupDN)
	state.AnyIncludedEntryGroupDN = internaltypes.GetStringSet(r.AnyIncludedEntryGroupDN)
	state.NotAllIncludedEntryGroupDN = internaltypes.GetStringSet(r.NotAllIncludedEntryGroupDN)
	state.NoneIncludedEntryGroupDN = internaltypes.GetStringSet(r.NoneIncludedEntryGroupDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a AggregateSearchEntryCriteriaResponse object into the model struct
func readAggregateSearchEntryCriteriaResponseDataSource(ctx context.Context, r *client.AggregateSearchEntryCriteriaResponse, state *searchEntryCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("aggregate")
	state.Id = types.StringValue(r.Id)
	state.AllIncludedSearchEntryCriteria = internaltypes.GetStringSet(r.AllIncludedSearchEntryCriteria)
	state.AnyIncludedSearchEntryCriteria = internaltypes.GetStringSet(r.AnyIncludedSearchEntryCriteria)
	state.NotAllIncludedSearchEntryCriteria = internaltypes.GetStringSet(r.NotAllIncludedSearchEntryCriteria)
	state.NoneIncludedSearchEntryCriteria = internaltypes.GetStringSet(r.NoneIncludedSearchEntryCriteria)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a ThirdPartySearchEntryCriteriaResponse object into the model struct
func readThirdPartySearchEntryCriteriaResponseDataSource(ctx context.Context, r *client.ThirdPartySearchEntryCriteriaResponse, state *searchEntryCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read resource information
func (r *searchEntryCriteriaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state searchEntryCriteriaDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SearchEntryCriteriaApi.GetSearchEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Search Entry Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.SimpleSearchEntryCriteriaResponse != nil {
		readSimpleSearchEntryCriteriaResponseDataSource(ctx, readResponse.SimpleSearchEntryCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AggregateSearchEntryCriteriaResponse != nil {
		readAggregateSearchEntryCriteriaResponseDataSource(ctx, readResponse.AggregateSearchEntryCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartySearchEntryCriteriaResponse != nil {
		readThirdPartySearchEntryCriteriaResponseDataSource(ctx, readResponse.ThirdPartySearchEntryCriteriaResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
