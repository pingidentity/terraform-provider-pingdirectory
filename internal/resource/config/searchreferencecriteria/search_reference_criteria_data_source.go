package searchreferencecriteria

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
	_ datasource.DataSource              = &searchReferenceCriteriaDataSource{}
	_ datasource.DataSourceWithConfigure = &searchReferenceCriteriaDataSource{}
)

// Create a Search Reference Criteria data source
func NewSearchReferenceCriteriaDataSource() datasource.DataSource {
	return &searchReferenceCriteriaDataSource{}
}

// searchReferenceCriteriaDataSource is the datasource implementation.
type searchReferenceCriteriaDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *searchReferenceCriteriaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_search_reference_criteria"
}

// Configure adds the provider configured client to the data source.
func (r *searchReferenceCriteriaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type searchReferenceCriteriaDataSourceModel struct {
	Id                                    types.String `tfsdk:"id"`
	Type                                  types.String `tfsdk:"type"`
	ExtensionClass                        types.String `tfsdk:"extension_class"`
	ExtensionArgument                     types.Set    `tfsdk:"extension_argument"`
	AllIncludedSearchReferenceCriteria    types.Set    `tfsdk:"all_included_search_reference_criteria"`
	AnyIncludedSearchReferenceCriteria    types.Set    `tfsdk:"any_included_search_reference_criteria"`
	NotAllIncludedSearchReferenceCriteria types.Set    `tfsdk:"not_all_included_search_reference_criteria"`
	NoneIncludedSearchReferenceCriteria   types.Set    `tfsdk:"none_included_search_reference_criteria"`
	RequestCriteria                       types.String `tfsdk:"request_criteria"`
	AllIncludedReferenceControl           types.Set    `tfsdk:"all_included_reference_control"`
	AnyIncludedReferenceControl           types.Set    `tfsdk:"any_included_reference_control"`
	NotAllIncludedReferenceControl        types.Set    `tfsdk:"not_all_included_reference_control"`
	NoneIncludedReferenceControl          types.Set    `tfsdk:"none_included_reference_control"`
	Description                           types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the datasource.
func (r *searchReferenceCriteriaDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Describes a Search Reference Criteria.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Name of this object.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of Search Reference Criteria resource. Options are ['simple', 'aggregate', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Search Reference Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Search Reference Criteria. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"all_included_search_reference_criteria": schema.SetAttribute{
				Description: "Specifies a search reference criteria object that must match the associated search result reference in order to match the aggregate search reference criteria. If one or more all-included search reference criteria objects are provided, then a search result reference must match all of them in order to match the aggregate search reference criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_search_reference_criteria": schema.SetAttribute{
				Description: "Specifies a search reference criteria object that may match the associated search result reference in order to match the aggregate search reference criteria. If one or more any-included search reference criteria objects are provided, then a search result reference must match at least one of them in order to match the aggregate search reference criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_search_reference_criteria": schema.SetAttribute{
				Description: "Specifies a search reference criteria object that should not match the associated search result reference in order to match the aggregate search reference criteria. If one or more not-all-included search reference criteria objects are provided, then a search result reference must not match all of them (that is, it may match zero or more of them, but it must not match all of them) in order to match the aggregate search reference criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_search_reference_criteria": schema.SetAttribute{
				Description: "Specifies a search reference criteria object that must not match the associated search result reference in order to match the aggregate search reference criteria. If one or more none-included search reference criteria objects are provided, then a search result reference must not match any of them in order to match the aggregate search reference criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"request_criteria": schema.StringAttribute{
				Description: "Specifies a request criteria object that must match the associated request for references included in this Simple Search Reference Criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"all_included_reference_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must be present in search result references included in this Simple Search Reference Criteria. If any control OIDs are provided, then the reference must contain all of those controls.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_reference_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that may be present in search result references included in this Simple Search Reference Criteria. If any control OIDs are provided, then the reference must contain at least one of those controls.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"not_all_included_reference_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that should not be present in search result references included in this Simple Search Reference Criteria. If any control OIDs are provided, then the reference must not contain at least one of those controls (that is, it may contain zero or more of those controls, but not all of them).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"none_included_reference_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must not be present in search result references included in this Simple Search Reference Criteria. If any control OIDs are provided, then the reference must not contain any of those controls.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Search Reference Criteria",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
}

// Read a SimpleSearchReferenceCriteriaResponse object into the model struct
func readSimpleSearchReferenceCriteriaResponseDataSource(ctx context.Context, r *client.SimpleSearchReferenceCriteriaResponse, state *searchReferenceCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("simple")
	state.Id = types.StringValue(r.Id)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.AllIncludedReferenceControl = internaltypes.GetStringSet(r.AllIncludedReferenceControl)
	state.AnyIncludedReferenceControl = internaltypes.GetStringSet(r.AnyIncludedReferenceControl)
	state.NotAllIncludedReferenceControl = internaltypes.GetStringSet(r.NotAllIncludedReferenceControl)
	state.NoneIncludedReferenceControl = internaltypes.GetStringSet(r.NoneIncludedReferenceControl)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a AggregateSearchReferenceCriteriaResponse object into the model struct
func readAggregateSearchReferenceCriteriaResponseDataSource(ctx context.Context, r *client.AggregateSearchReferenceCriteriaResponse, state *searchReferenceCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("aggregate")
	state.Id = types.StringValue(r.Id)
	state.AllIncludedSearchReferenceCriteria = internaltypes.GetStringSet(r.AllIncludedSearchReferenceCriteria)
	state.AnyIncludedSearchReferenceCriteria = internaltypes.GetStringSet(r.AnyIncludedSearchReferenceCriteria)
	state.NotAllIncludedSearchReferenceCriteria = internaltypes.GetStringSet(r.NotAllIncludedSearchReferenceCriteria)
	state.NoneIncludedSearchReferenceCriteria = internaltypes.GetStringSet(r.NoneIncludedSearchReferenceCriteria)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a ThirdPartySearchReferenceCriteriaResponse object into the model struct
func readThirdPartySearchReferenceCriteriaResponseDataSource(ctx context.Context, r *client.ThirdPartySearchReferenceCriteriaResponse, state *searchReferenceCriteriaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read resource information
func (r *searchReferenceCriteriaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state searchReferenceCriteriaDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SearchReferenceCriteriaApi.GetSearchReferenceCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Search Reference Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.SimpleSearchReferenceCriteriaResponse != nil {
		readSimpleSearchReferenceCriteriaResponseDataSource(ctx, readResponse.SimpleSearchReferenceCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AggregateSearchReferenceCriteriaResponse != nil {
		readAggregateSearchReferenceCriteriaResponseDataSource(ctx, readResponse.AggregateSearchReferenceCriteriaResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartySearchReferenceCriteriaResponse != nil {
		readThirdPartySearchReferenceCriteriaResponseDataSource(ctx, readResponse.ThirdPartySearchReferenceCriteriaResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
