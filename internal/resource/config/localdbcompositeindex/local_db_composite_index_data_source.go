// Copyright © 2025 Ping Identity Corporation

package localdbcompositeindex

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &localDbCompositeIndexDataSource{}
	_ datasource.DataSourceWithConfigure = &localDbCompositeIndexDataSource{}
)

// Create a Local Db Composite Index data source
func NewLocalDbCompositeIndexDataSource() datasource.DataSource {
	return &localDbCompositeIndexDataSource{}
}

// localDbCompositeIndexDataSource is the datasource implementation.
type localDbCompositeIndexDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *localDbCompositeIndexDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_db_composite_index"
}

// Configure adds the provider configured client to the data source.
func (r *localDbCompositeIndexDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type localDbCompositeIndexDataSourceModel struct {
	Id                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Type                   types.String `tfsdk:"type"`
	BackendName            types.String `tfsdk:"backend_name"`
	Description            types.String `tfsdk:"description"`
	IndexFilterPattern     types.String `tfsdk:"index_filter_pattern"`
	IndexBaseDNPattern     types.String `tfsdk:"index_base_dn_pattern"`
	IndexEntryLimit        types.Int64  `tfsdk:"index_entry_limit"`
	PrimeIndex             types.Bool   `tfsdk:"prime_index"`
	PrimeInternalNodesOnly types.Bool   `tfsdk:"prime_internal_nodes_only"`
	CacheMode              types.String `tfsdk:"cache_mode"`
}

// GetSchema defines the schema for the datasource.
func (r *localDbCompositeIndexDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Local Db Composite Index.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Local DB Composite Index resource. Options are ['local-db-composite-index']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"backend_name": schema.StringAttribute{
				Description: "Name of the parent Backend",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Local DB Composite Index",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"index_filter_pattern": schema.StringAttribute{
				Description: "A filter pattern that identifies which entries to include in the index.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"index_base_dn_pattern": schema.StringAttribute{
				Description: "An optional base DN pattern that identifies portions of the DIT in which entries to index may exist.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"index_entry_limit": schema.Int64Attribute{
				Description: "The maximum number of entries that any single index key will be allowed to match before the server stops maintaining the ID set for that index key.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"prime_index": schema.BoolAttribute{
				Description: "Indicates whether the server should load the contents of this index into memory when the backend is being opened.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"prime_internal_nodes_only": schema.BoolAttribute{
				Description: "Indicates whether to only prime the internal nodes of the index database, rather than priming both internal and leaf nodes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"cache_mode": schema.StringAttribute{
				Description: "The behavior that the server should exhibit when storing information from this index in the database cache.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a LocalDbCompositeIndexResponse object into the model struct
func readLocalDbCompositeIndexResponseDataSource(ctx context.Context, r *client.LocalDbCompositeIndexResponse, state *localDbCompositeIndexDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("local-db-composite-index")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.IndexFilterPattern = types.StringValue(r.IndexFilterPattern)
	state.IndexBaseDNPattern = internaltypes.StringTypeOrNil(r.IndexBaseDNPattern, false)
	state.IndexEntryLimit = internaltypes.Int64TypeOrNil(r.IndexEntryLimit)
	state.PrimeIndex = internaltypes.BoolTypeOrNil(r.PrimeIndex)
	state.PrimeInternalNodesOnly = internaltypes.BoolTypeOrNil(r.PrimeInternalNodesOnly)
	state.CacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlocalDbCompositeIndexCacheModeProp(r.CacheMode), false)
}

// Read resource information
func (r *localDbCompositeIndexDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state localDbCompositeIndexDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LocalDbCompositeIndexAPI.GetLocalDbCompositeIndex(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.BackendName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Local Db Composite Index", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLocalDbCompositeIndexResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
