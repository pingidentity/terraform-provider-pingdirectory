package pluginroot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10100/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &pluginRootDataSource{}
	_ datasource.DataSourceWithConfigure = &pluginRootDataSource{}
)

// Create a Plugin Root data source
func NewPluginRootDataSource() datasource.DataSource {
	return &pluginRootDataSource{}
}

// pluginRootDataSource is the datasource implementation.
type pluginRootDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *pluginRootDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_plugin_root"
}

// Configure adds the provider configured client to the data source.
func (r *pluginRootDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type pluginRootDataSourceModel struct {
	Id                                     types.String `tfsdk:"id"`
	Type                                   types.String `tfsdk:"type"`
	PluginOrderStartup                     types.String `tfsdk:"plugin_order_startup"`
	PluginOrderShutdown                    types.String `tfsdk:"plugin_order_shutdown"`
	PluginOrderPostConnect                 types.String `tfsdk:"plugin_order_post_connect"`
	PluginOrderPostDisconnect              types.String `tfsdk:"plugin_order_post_disconnect"`
	PluginOrderLDIFImport                  types.String `tfsdk:"plugin_order_ldif_import"`
	PluginOrderLDIFExport                  types.String `tfsdk:"plugin_order_ldif_export"`
	PluginOrderPreParseAbandon             types.String `tfsdk:"plugin_order_pre_parse_abandon"`
	PluginOrderPreParseAdd                 types.String `tfsdk:"plugin_order_pre_parse_add"`
	PluginOrderPreParseBind                types.String `tfsdk:"plugin_order_pre_parse_bind"`
	PluginOrderPreParseCompare             types.String `tfsdk:"plugin_order_pre_parse_compare"`
	PluginOrderPreParseDelete              types.String `tfsdk:"plugin_order_pre_parse_delete"`
	PluginOrderPreParseExtended            types.String `tfsdk:"plugin_order_pre_parse_extended"`
	PluginOrderPreParseModify              types.String `tfsdk:"plugin_order_pre_parse_modify"`
	PluginOrderPreParseModifyDN            types.String `tfsdk:"plugin_order_pre_parse_modify_dn"`
	PluginOrderPreParseSearch              types.String `tfsdk:"plugin_order_pre_parse_search"`
	PluginOrderPreParseUnbind              types.String `tfsdk:"plugin_order_pre_parse_unbind"`
	PluginOrderPreOperationAdd             types.String `tfsdk:"plugin_order_pre_operation_add"`
	PluginOrderPreOperationBind            types.String `tfsdk:"plugin_order_pre_operation_bind"`
	PluginOrderPreOperationCompare         types.String `tfsdk:"plugin_order_pre_operation_compare"`
	PluginOrderPreOperationDelete          types.String `tfsdk:"plugin_order_pre_operation_delete"`
	PluginOrderPreOperationExtended        types.String `tfsdk:"plugin_order_pre_operation_extended"`
	PluginOrderPreOperationModify          types.String `tfsdk:"plugin_order_pre_operation_modify"`
	PluginOrderPreOperationModifyDN        types.String `tfsdk:"plugin_order_pre_operation_modify_dn"`
	PluginOrderPreOperationSearch          types.String `tfsdk:"plugin_order_pre_operation_search"`
	PluginOrderPostOperationAbandon        types.String `tfsdk:"plugin_order_post_operation_abandon"`
	PluginOrderPostOperationAdd            types.String `tfsdk:"plugin_order_post_operation_add"`
	PluginOrderPostOperationBind           types.String `tfsdk:"plugin_order_post_operation_bind"`
	PluginOrderPostOperationCompare        types.String `tfsdk:"plugin_order_post_operation_compare"`
	PluginOrderPostOperationDelete         types.String `tfsdk:"plugin_order_post_operation_delete"`
	PluginOrderPostOperationExtended       types.String `tfsdk:"plugin_order_post_operation_extended"`
	PluginOrderPostOperationModify         types.String `tfsdk:"plugin_order_post_operation_modify"`
	PluginOrderPostOperationModifyDN       types.String `tfsdk:"plugin_order_post_operation_modify_dn"`
	PluginOrderPostOperationSearch         types.String `tfsdk:"plugin_order_post_operation_search"`
	PluginOrderPostOperationUnbind         types.String `tfsdk:"plugin_order_post_operation_unbind"`
	PluginOrderPostResponseAdd             types.String `tfsdk:"plugin_order_post_response_add"`
	PluginOrderPostResponseBind            types.String `tfsdk:"plugin_order_post_response_bind"`
	PluginOrderPostResponseCompare         types.String `tfsdk:"plugin_order_post_response_compare"`
	PluginOrderPostResponseDelete          types.String `tfsdk:"plugin_order_post_response_delete"`
	PluginOrderPostResponseExtended        types.String `tfsdk:"plugin_order_post_response_extended"`
	PluginOrderPostResponseModify          types.String `tfsdk:"plugin_order_post_response_modify"`
	PluginOrderPostResponseModifyDN        types.String `tfsdk:"plugin_order_post_response_modify_dn"`
	PluginOrderPostSynchronizationAdd      types.String `tfsdk:"plugin_order_post_synchronization_add"`
	PluginOrderPostSynchronizationDelete   types.String `tfsdk:"plugin_order_post_synchronization_delete"`
	PluginOrderPostSynchronizationModify   types.String `tfsdk:"plugin_order_post_synchronization_modify"`
	PluginOrderPostSynchronizationModifyDN types.String `tfsdk:"plugin_order_post_synchronization_modify_dn"`
	PluginOrderPostResponseSearch          types.String `tfsdk:"plugin_order_post_response_search"`
	PluginOrderSearchResultEntry           types.String `tfsdk:"plugin_order_search_result_entry"`
	PluginOrderSearchResultReference       types.String `tfsdk:"plugin_order_search_result_reference"`
	PluginOrderSubordinateModifyDN         types.String `tfsdk:"plugin_order_subordinate_modify_dn"`
	PluginOrderIntermediateResponse        types.String `tfsdk:"plugin_order_intermediate_response"`
}

// GetSchema defines the schema for the datasource.
func (r *pluginRootDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Plugin Root.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Plugin Root resource. Options are ['plugin-root']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_startup": schema.StringAttribute{
				Description: "Specifies the order in which startup plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_shutdown": schema.StringAttribute{
				Description: "Specifies the order in which shutdown plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_connect": schema.StringAttribute{
				Description: "Specifies the order in which post-connect plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_disconnect": schema.StringAttribute{
				Description: "Specifies the order in which post-disconnect plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_ldif_import": schema.StringAttribute{
				Description: "Specifies the order in which LDIF import plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_ldif_export": schema.StringAttribute{
				Description: "Specifies the order in which LDIF export plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_pre_parse_abandon": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse abandon plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_pre_parse_add": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse add plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_pre_parse_bind": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse bind plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_pre_parse_compare": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse compare plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_pre_parse_delete": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse delete plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_pre_parse_extended": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse extended operation plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_pre_parse_modify": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse modify plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_pre_parse_modify_dn": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse modify DN plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_pre_parse_search": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse search plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_pre_parse_unbind": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse unbind plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_pre_operation_add": schema.StringAttribute{
				Description: "Specifies the order in which pre-operation add plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_pre_operation_bind": schema.StringAttribute{
				Description: "Specifies the order in which pre-operation bind plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_pre_operation_compare": schema.StringAttribute{
				Description: "Specifies the order in which pre-operation compare plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_pre_operation_delete": schema.StringAttribute{
				Description: "Specifies the order in which pre-operation delete plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_pre_operation_extended": schema.StringAttribute{
				Description: "Specifies the order in which pre-operation extended operation plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_pre_operation_modify": schema.StringAttribute{
				Description: "Specifies the order in which pre-operation modify plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_pre_operation_modify_dn": schema.StringAttribute{
				Description: "Specifies the order in which pre-operation modify DN plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_pre_operation_search": schema.StringAttribute{
				Description: "Specifies the order in which pre-operation search plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_operation_abandon": schema.StringAttribute{
				Description: "Specifies the order in which post-operation abandon plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_operation_add": schema.StringAttribute{
				Description: "Specifies the order in which post-operation add plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_operation_bind": schema.StringAttribute{
				Description: "Specifies the order in which post-operation bind plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_operation_compare": schema.StringAttribute{
				Description: "Specifies the order in which post-operation compare plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_operation_delete": schema.StringAttribute{
				Description: "Specifies the order in which post-operation delete plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_operation_extended": schema.StringAttribute{
				Description: "Specifies the order in which post-operation extended operation plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_operation_modify": schema.StringAttribute{
				Description: "Specifies the order in which post-operation modify plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_operation_modify_dn": schema.StringAttribute{
				Description: "Specifies the order in which post-operation modify DN plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_operation_search": schema.StringAttribute{
				Description: "Specifies the order in which post-operation search plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_operation_unbind": schema.StringAttribute{
				Description: "Specifies the order in which post-operation unbind plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_response_add": schema.StringAttribute{
				Description: "Specifies the order in which post-response add plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_response_bind": schema.StringAttribute{
				Description: "Specifies the order in which post-response bind plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_response_compare": schema.StringAttribute{
				Description: "Specifies the order in which post-response compare plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_response_delete": schema.StringAttribute{
				Description: "Specifies the order in which post-response delete plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_response_extended": schema.StringAttribute{
				Description: "Specifies the order in which post-response extended operation plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_response_modify": schema.StringAttribute{
				Description: "Specifies the order in which post-response modify plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_response_modify_dn": schema.StringAttribute{
				Description: "Specifies the order in which post-response modify DN plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_synchronization_add": schema.StringAttribute{
				Description: "Specifies the order in which post-synchronization add plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_synchronization_delete": schema.StringAttribute{
				Description: "Specifies the order in which post-synchronization delete plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_synchronization_modify": schema.StringAttribute{
				Description: "Specifies the order in which post-synchronization modify plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_synchronization_modify_dn": schema.StringAttribute{
				Description: "Specifies the order in which post-synchronization modify DN plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_post_response_search": schema.StringAttribute{
				Description: "Specifies the order in which post-response search plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_search_result_entry": schema.StringAttribute{
				Description: "Specifies the order in which search result entry plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_search_result_reference": schema.StringAttribute{
				Description: "Specifies the order in which search result reference plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_subordinate_modify_dn": schema.StringAttribute{
				Description: "Specifies the order in which subordinate modify DN plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_order_intermediate_response": schema.StringAttribute{
				Description: "Specifies the order in which intermediate response plug-ins are to be loaded and invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a PluginRootResponse object into the model struct
func readPluginRootResponseDataSource(ctx context.Context, r *client.PluginRootResponse, state *pluginRootDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("plugin-root")
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
	state.PluginOrderStartup = internaltypes.StringTypeOrNil(r.PluginOrderStartup, false)
	state.PluginOrderShutdown = internaltypes.StringTypeOrNil(r.PluginOrderShutdown, false)
	state.PluginOrderPostConnect = internaltypes.StringTypeOrNil(r.PluginOrderPostConnect, false)
	state.PluginOrderPostDisconnect = internaltypes.StringTypeOrNil(r.PluginOrderPostDisconnect, false)
	state.PluginOrderLDIFImport = internaltypes.StringTypeOrNil(r.PluginOrderLDIFImport, false)
	state.PluginOrderLDIFExport = internaltypes.StringTypeOrNil(r.PluginOrderLDIFExport, false)
	state.PluginOrderPreParseAbandon = internaltypes.StringTypeOrNil(r.PluginOrderPreParseAbandon, false)
	state.PluginOrderPreParseAdd = internaltypes.StringTypeOrNil(r.PluginOrderPreParseAdd, false)
	state.PluginOrderPreParseBind = internaltypes.StringTypeOrNil(r.PluginOrderPreParseBind, false)
	state.PluginOrderPreParseCompare = internaltypes.StringTypeOrNil(r.PluginOrderPreParseCompare, false)
	state.PluginOrderPreParseDelete = internaltypes.StringTypeOrNil(r.PluginOrderPreParseDelete, false)
	state.PluginOrderPreParseExtended = internaltypes.StringTypeOrNil(r.PluginOrderPreParseExtended, false)
	state.PluginOrderPreParseModify = internaltypes.StringTypeOrNil(r.PluginOrderPreParseModify, false)
	state.PluginOrderPreParseModifyDN = internaltypes.StringTypeOrNil(r.PluginOrderPreParseModifyDN, false)
	state.PluginOrderPreParseSearch = internaltypes.StringTypeOrNil(r.PluginOrderPreParseSearch, false)
	state.PluginOrderPreParseUnbind = internaltypes.StringTypeOrNil(r.PluginOrderPreParseUnbind, false)
	state.PluginOrderPreOperationAdd = internaltypes.StringTypeOrNil(r.PluginOrderPreOperationAdd, false)
	state.PluginOrderPreOperationBind = internaltypes.StringTypeOrNil(r.PluginOrderPreOperationBind, false)
	state.PluginOrderPreOperationCompare = internaltypes.StringTypeOrNil(r.PluginOrderPreOperationCompare, false)
	state.PluginOrderPreOperationDelete = internaltypes.StringTypeOrNil(r.PluginOrderPreOperationDelete, false)
	state.PluginOrderPreOperationExtended = internaltypes.StringTypeOrNil(r.PluginOrderPreOperationExtended, false)
	state.PluginOrderPreOperationModify = internaltypes.StringTypeOrNil(r.PluginOrderPreOperationModify, false)
	state.PluginOrderPreOperationModifyDN = internaltypes.StringTypeOrNil(r.PluginOrderPreOperationModifyDN, false)
	state.PluginOrderPreOperationSearch = internaltypes.StringTypeOrNil(r.PluginOrderPreOperationSearch, false)
	state.PluginOrderPostOperationAbandon = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationAbandon, false)
	state.PluginOrderPostOperationAdd = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationAdd, false)
	state.PluginOrderPostOperationBind = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationBind, false)
	state.PluginOrderPostOperationCompare = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationCompare, false)
	state.PluginOrderPostOperationDelete = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationDelete, false)
	state.PluginOrderPostOperationExtended = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationExtended, false)
	state.PluginOrderPostOperationModify = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationModify, false)
	state.PluginOrderPostOperationModifyDN = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationModifyDN, false)
	state.PluginOrderPostOperationSearch = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationSearch, false)
	state.PluginOrderPostOperationUnbind = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationUnbind, false)
	state.PluginOrderPostResponseAdd = internaltypes.StringTypeOrNil(r.PluginOrderPostResponseAdd, false)
	state.PluginOrderPostResponseBind = internaltypes.StringTypeOrNil(r.PluginOrderPostResponseBind, false)
	state.PluginOrderPostResponseCompare = internaltypes.StringTypeOrNil(r.PluginOrderPostResponseCompare, false)
	state.PluginOrderPostResponseDelete = internaltypes.StringTypeOrNil(r.PluginOrderPostResponseDelete, false)
	state.PluginOrderPostResponseExtended = internaltypes.StringTypeOrNil(r.PluginOrderPostResponseExtended, false)
	state.PluginOrderPostResponseModify = internaltypes.StringTypeOrNil(r.PluginOrderPostResponseModify, false)
	state.PluginOrderPostResponseModifyDN = internaltypes.StringTypeOrNil(r.PluginOrderPostResponseModifyDN, false)
	state.PluginOrderPostSynchronizationAdd = internaltypes.StringTypeOrNil(r.PluginOrderPostSynchronizationAdd, false)
	state.PluginOrderPostSynchronizationDelete = internaltypes.StringTypeOrNil(r.PluginOrderPostSynchronizationDelete, false)
	state.PluginOrderPostSynchronizationModify = internaltypes.StringTypeOrNil(r.PluginOrderPostSynchronizationModify, false)
	state.PluginOrderPostSynchronizationModifyDN = internaltypes.StringTypeOrNil(r.PluginOrderPostSynchronizationModifyDN, false)
	state.PluginOrderPostResponseSearch = internaltypes.StringTypeOrNil(r.PluginOrderPostResponseSearch, false)
	state.PluginOrderSearchResultEntry = internaltypes.StringTypeOrNil(r.PluginOrderSearchResultEntry, false)
	state.PluginOrderSearchResultReference = internaltypes.StringTypeOrNil(r.PluginOrderSearchResultReference, false)
	state.PluginOrderSubordinateModifyDN = internaltypes.StringTypeOrNil(r.PluginOrderSubordinateModifyDN, false)
	state.PluginOrderIntermediateResponse = internaltypes.StringTypeOrNil(r.PluginOrderIntermediateResponse, false)
}

// Read resource information
func (r *pluginRootDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state pluginRootDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginRootAPI.GetPluginRoot(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Plugin Root", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readPluginRootResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
