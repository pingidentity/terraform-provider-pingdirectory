// Copyright © 2025 Ping Identity Corporation

package entrycache

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &entryCacheDataSource{}
	_ datasource.DataSourceWithConfigure = &entryCacheDataSource{}
)

// Create a Entry Cache data source
func NewEntryCacheDataSource() datasource.DataSource {
	return &entryCacheDataSource{}
}

// entryCacheDataSource is the datasource implementation.
type entryCacheDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *entryCacheDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_cache"
}

// Configure adds the provider configured client to the data source.
func (r *entryCacheDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type entryCacheDataSourceModel struct {
	Id                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	Type                        types.String `tfsdk:"type"`
	MaxMemoryPercent            types.Int64  `tfsdk:"max_memory_percent"`
	MaxEntries                  types.Int64  `tfsdk:"max_entries"`
	OnlyCacheFrequentlyAccessed types.Bool   `tfsdk:"only_cache_frequently_accessed"`
	IncludeFilter               types.Set    `tfsdk:"include_filter"`
	ExcludeFilter               types.Set    `tfsdk:"exclude_filter"`
	MinCacheEntryValueCount     types.Int64  `tfsdk:"min_cache_entry_value_count"`
	MinCacheEntryAttribute      types.Set    `tfsdk:"min_cache_entry_attribute"`
	Description                 types.String `tfsdk:"description"`
	Enabled                     types.Bool   `tfsdk:"enabled"`
	CacheLevel                  types.Int64  `tfsdk:"cache_level"`
	CacheUnindexedSearchResults types.Bool   `tfsdk:"cache_unindexed_search_results"`
}

// GetSchema defines the schema for the datasource.
func (r *entryCacheDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Entry Cache.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Entry Cache resource. Options are ['fifo']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_memory_percent": schema.Int64Attribute{
				Description: "Specifies the maximum amount of memory, as a percentage of the total maximum JVM heap size, that this cache should occupy when full. If the amount of memory the cache is using is greater than this amount, then an attempt to put a new entry in the cache will be ignored and will cause the oldest entry to be purged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_entries": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that will be allowed in the cache. Once the cache reaches this size, then adding new entries will cause existing entries to be purged, starting with the oldest.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"only_cache_frequently_accessed": schema.BoolAttribute{
				Description: "Specifies that the cache should only store entries which are accessed much more frequently than the average entry. The cache will observe attempts to place entries in the cache and compare an entry's accesses to the average entry's.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_filter": schema.SetAttribute{
				Description: "The set of filters that define the entries that should be included in the cache.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"exclude_filter": schema.SetAttribute{
				Description: "The set of filters that define the entries that should be excluded from the cache.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"min_cache_entry_value_count": schema.Int64Attribute{
				Description: "Specifies the minimum number of attribute values (optionally across a specified subset of attributes as defined in the min-cache-entry-attributes property) for entries that should be held in the cache. Entries with fewer than this number of attribute values will be excluded from the cache.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"min_cache_entry_attribute": schema.SetAttribute{
				Description: "Specifies the names of the attribute types for which the min-cache-entry-value-count property should apply. If no attribute types are specified, then all user attributes will be examined.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Entry Cache",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Entry Cache is enabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"cache_level": schema.Int64Attribute{
				Description: "Specifies the cache level in the cache order if more than one instance of the cache is configured.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"cache_unindexed_search_results": schema.BoolAttribute{
				Description: "Indicates whether the entry cache should be updated with entries that have been returned to the client during the course of processing an unindexed search.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a FifoEntryCacheResponse object into the model struct
func readFifoEntryCacheResponseDataSource(ctx context.Context, r *client.FifoEntryCacheResponse, state *entryCacheDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("fifo")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MaxMemoryPercent = internaltypes.Int64TypeOrNil(r.MaxMemoryPercent)
	state.MaxEntries = internaltypes.Int64TypeOrNil(r.MaxEntries)
	state.OnlyCacheFrequentlyAccessed = internaltypes.BoolTypeOrNil(r.OnlyCacheFrequentlyAccessed)
	state.IncludeFilter = internaltypes.GetStringSet(r.IncludeFilter)
	state.ExcludeFilter = internaltypes.GetStringSet(r.ExcludeFilter)
	state.MinCacheEntryValueCount = internaltypes.Int64TypeOrNil(r.MinCacheEntryValueCount)
	state.MinCacheEntryAttribute = internaltypes.GetStringSet(r.MinCacheEntryAttribute)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.CacheLevel = types.Int64Value(r.CacheLevel)
	state.CacheUnindexedSearchResults = internaltypes.BoolTypeOrNil(r.CacheUnindexedSearchResults)
}

// Read resource information
func (r *entryCacheDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state entryCacheDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.EntryCacheAPI.GetEntryCache(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Entry Cache", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readFifoEntryCacheResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
