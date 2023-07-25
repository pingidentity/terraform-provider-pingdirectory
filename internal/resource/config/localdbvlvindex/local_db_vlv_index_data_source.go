package localdbvlvindex

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
	_ datasource.DataSource              = &localDbVlvIndexDataSource{}
	_ datasource.DataSourceWithConfigure = &localDbVlvIndexDataSource{}
)

// Create a Local Db Vlv Index data source
func NewLocalDbVlvIndexDataSource() datasource.DataSource {
	return &localDbVlvIndexDataSource{}
}

// localDbVlvIndexDataSource is the datasource implementation.
type localDbVlvIndexDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *localDbVlvIndexDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_db_vlv_index"
}

// Configure adds the provider configured client to the data source.
func (r *localDbVlvIndexDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type localDbVlvIndexDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	BackendName  types.String `tfsdk:"backend_name"`
	BaseDN       types.String `tfsdk:"base_dn"`
	Scope        types.String `tfsdk:"scope"`
	Filter       types.String `tfsdk:"filter"`
	SortOrder    types.String `tfsdk:"sort_order"`
	Name         types.String `tfsdk:"name"`
	MaxBlockSize types.Int64  `tfsdk:"max_block_size"`
	CacheMode    types.String `tfsdk:"cache_mode"`
}

// GetSchema defines the schema for the datasource.
func (r *localDbVlvIndexDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Describes a Local Db Vlv Index.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Name of this object.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"backend_name": schema.StringAttribute{
				Description: "Name of the parent Backend",
				Required:    true,
			},
			"base_dn": schema.StringAttribute{
				Description: "Specifies the base DN used in the search query that is being indexed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"scope": schema.StringAttribute{
				Description: "Specifies the LDAP scope of the query that is being indexed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"filter": schema.StringAttribute{
				Description: "Specifies the LDAP filter used in the query that is being indexed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"sort_order": schema.StringAttribute{
				Description: "Specifies the names of the attributes that are used to sort the entries for the query being indexed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Specifies a unique name for this VLV index.",
				Required:    true,
			},
			"max_block_size": schema.Int64Attribute{
				Description: "Specifies the number of entry IDs to store in a single sorted set before it must be split.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the database for this index.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
}

// Read a LocalDbVlvIndexResponse object into the model struct
func readLocalDbVlvIndexResponseDataSource(ctx context.Context, r *client.LocalDbVlvIndexResponse, state *localDbVlvIndexDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.BaseDN = types.StringValue(r.BaseDN)
	state.Scope = types.StringValue(r.Scope.String())
	state.Filter = types.StringValue(r.Filter)
	state.SortOrder = types.StringValue(r.SortOrder)
	state.Name = types.StringValue(r.Name)
	state.MaxBlockSize = internaltypes.Int64TypeOrNil(r.MaxBlockSize)
	state.CacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlocalDbVlvIndexCacheModeProp(r.CacheMode), false)
}

// Read resource information
func (r *localDbVlvIndexDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state localDbVlvIndexDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LocalDbVlvIndexApi.GetLocalDbVlvIndex(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.BackendName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Local Db Vlv Index", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLocalDbVlvIndexResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
