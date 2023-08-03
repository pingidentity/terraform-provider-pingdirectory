package localdbindex

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
	_ datasource.DataSource              = &localDbIndexDataSource{}
	_ datasource.DataSourceWithConfigure = &localDbIndexDataSource{}
)

// Create a Local Db Index data source
func NewLocalDbIndexDataSource() datasource.DataSource {
	return &localDbIndexDataSource{}
}

// localDbIndexDataSource is the datasource implementation.
type localDbIndexDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *localDbIndexDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_db_index"
}

// Configure adds the provider configured client to the data source.
func (r *localDbIndexDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type localDbIndexDataSourceModel struct {
	Id                                           types.String `tfsdk:"id"`
	BackendName                                  types.String `tfsdk:"backend_name"`
	Attribute                                    types.String `tfsdk:"attribute"`
	IndexEntryLimit                              types.Int64  `tfsdk:"index_entry_limit"`
	SubstringIndexEntryLimit                     types.Int64  `tfsdk:"substring_index_entry_limit"`
	MaintainMatchCountForKeysExceedingEntryLimit types.Bool   `tfsdk:"maintain_match_count_for_keys_exceeding_entry_limit"`
	IndexType                                    types.Set    `tfsdk:"index_type"`
	SubstringLength                              types.Int64  `tfsdk:"substring_length"`
	PrimeIndex                                   types.Bool   `tfsdk:"prime_index"`
	PrimeInternalNodesOnly                       types.Bool   `tfsdk:"prime_internal_nodes_only"`
	EqualityIndexFilter                          types.Set    `tfsdk:"equality_index_filter"`
	MaintainEqualityIndexWithoutFilter           types.Bool   `tfsdk:"maintain_equality_index_without_filter"`
	CacheMode                                    types.String `tfsdk:"cache_mode"`
}

// GetSchema defines the schema for the datasource.
func (r *localDbIndexDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Local Db Index.",
		Attributes: map[string]schema.Attribute{
			"backend_name": schema.StringAttribute{
				Description: "Name of the parent Backend",
				Required:    true,
			},
			"attribute": schema.StringAttribute{
				Description: "Specifies the name of the attribute for which the index is to be maintained.",
				Required:    true,
			},
			"index_entry_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that are allowed to match a given index key before that particular index key is no longer maintained.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"substring_index_entry_limit": schema.Int64Attribute{
				Description: "Specifies, for substring indexes, the maximum number of entries that are allowed to match a given index key before that particular index key is no longer maintained. Setting a large limit can dramatically increase the database size on disk and have a big impact on server performance if the indexed attribute is modified frequently. When a very large limit is required, creating a dedicated composite index with an index-filter-pattern of (attr=*?*) will give the best balance between search and update performance.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maintain_match_count_for_keys_exceeding_entry_limit": schema.BoolAttribute{
				Description: "Indicates whether to continue to maintain a count of the number of matching entries for an index key even after that count exceeds the index entry limit.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"index_type": schema.SetAttribute{
				Description: "Specifies the type(s) of indexing that should be performed for the associated attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"substring_length": schema.Int64Attribute{
				Description: "The length of substrings in a substring index.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"prime_index": schema.BoolAttribute{
				Description: "If this option is enabled and this index's backend is configured to prime indexes, then this index will be loaded at startup.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"prime_internal_nodes_only": schema.BoolAttribute{
				Description: "If this option is enabled and this index's backend is configured to prime indexes using the preload method, then only the internal database nodes (i.e., the database keys but not values) should be primed when the backend is initialized.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"equality_index_filter": schema.SetAttribute{
				Description: "A search filter that may be used in conjunction with an equality component for the associated attribute type. If an equality index filter is defined, then an additional equality index will be maintained for the associated attribute, but only for entries which match the provided filter. Further, the index will be used only for searches containing an equality component with the associated attribute type ANDed with this filter.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"maintain_equality_index_without_filter": schema.BoolAttribute{
				Description: "Indicates whether to maintain a separate equality index for the associated attribute without any filter, in addition to maintaining an index for each equality index filter that is defined. If this is false, then the attribute will not be indexed for equality by itself but only in conjunction with the defined equality index filters.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the database for this index. This controls how much database cache memory can be consumed by this index.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a LocalDbIndexResponse object into the model struct
func readLocalDbIndexResponseDataSource(ctx context.Context, r *client.LocalDbIndexResponse, state *localDbIndexDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Attribute = types.StringValue(r.Attribute)
	state.IndexEntryLimit = internaltypes.Int64TypeOrNil(r.IndexEntryLimit)
	state.SubstringIndexEntryLimit = internaltypes.Int64TypeOrNil(r.SubstringIndexEntryLimit)
	state.MaintainMatchCountForKeysExceedingEntryLimit = internaltypes.BoolTypeOrNil(r.MaintainMatchCountForKeysExceedingEntryLimit)
	state.IndexType = internaltypes.GetStringSet(
		client.StringSliceEnumlocalDbIndexIndexTypeProp(r.IndexType))
	state.SubstringLength = internaltypes.Int64TypeOrNil(r.SubstringLength)
	state.PrimeIndex = internaltypes.BoolTypeOrNil(r.PrimeIndex)
	state.PrimeInternalNodesOnly = internaltypes.BoolTypeOrNil(r.PrimeInternalNodesOnly)
	state.EqualityIndexFilter = internaltypes.GetStringSet(r.EqualityIndexFilter)
	state.MaintainEqualityIndexWithoutFilter = internaltypes.BoolTypeOrNil(r.MaintainEqualityIndexWithoutFilter)
	state.CacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlocalDbIndexCacheModeProp(r.CacheMode), false)
}

// Read resource information
func (r *localDbIndexDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state localDbIndexDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LocalDbIndexApi.GetLocalDbIndex(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Attribute.ValueString(), state.BackendName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Local Db Index", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLocalDbIndexResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
