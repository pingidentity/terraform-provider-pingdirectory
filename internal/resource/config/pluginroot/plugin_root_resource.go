package pluginroot

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &pluginRootResource{}
	_ resource.ResourceWithConfigure   = &pluginRootResource{}
	_ resource.ResourceWithImportState = &pluginRootResource{}
)

// Create a Plugin Root resource
func NewPluginRootResource() resource.Resource {
	return &pluginRootResource{}
}

// pluginRootResource is the resource implementation.
type pluginRootResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *pluginRootResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_plugin_root"
}

// Configure adds the provider configured client to the resource.
func (r *pluginRootResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type pluginRootResourceModel struct {
	Id                                     types.String `tfsdk:"id"`
	LastUpdated                            types.String `tfsdk:"last_updated"`
	Notifications                          types.Set    `tfsdk:"notifications"`
	RequiredActions                        types.Set    `tfsdk:"required_actions"`
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

// GetSchema defines the schema for the resource.
func (r *pluginRootResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a Plugin Root.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Plugin Root resource. Options are ['plugin-root']",
				Optional:    false,
				Required:    false,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"plugin-root"}...),
				},
			},
			"plugin_order_startup": schema.StringAttribute{
				Description: "Specifies the order in which startup plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_shutdown": schema.StringAttribute{
				Description: "Specifies the order in which shutdown plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_connect": schema.StringAttribute{
				Description: "Specifies the order in which post-connect plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_disconnect": schema.StringAttribute{
				Description: "Specifies the order in which post-disconnect plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_ldif_import": schema.StringAttribute{
				Description: "Specifies the order in which LDIF import plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_ldif_export": schema.StringAttribute{
				Description: "Specifies the order in which LDIF export plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_pre_parse_abandon": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse abandon plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_pre_parse_add": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse add plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_pre_parse_bind": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse bind plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_pre_parse_compare": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse compare plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_pre_parse_delete": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse delete plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_pre_parse_extended": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse extended operation plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_pre_parse_modify": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse modify plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_pre_parse_modify_dn": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse modify DN plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_pre_parse_search": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse search plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_pre_parse_unbind": schema.StringAttribute{
				Description: "Specifies the order in which pre-parse unbind plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_pre_operation_add": schema.StringAttribute{
				Description: "Specifies the order in which pre-operation add plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_pre_operation_bind": schema.StringAttribute{
				Description: "Specifies the order in which pre-operation bind plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_pre_operation_compare": schema.StringAttribute{
				Description: "Specifies the order in which pre-operation compare plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_pre_operation_delete": schema.StringAttribute{
				Description: "Specifies the order in which pre-operation delete plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_pre_operation_extended": schema.StringAttribute{
				Description: "Specifies the order in which pre-operation extended operation plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_pre_operation_modify": schema.StringAttribute{
				Description: "Specifies the order in which pre-operation modify plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_pre_operation_modify_dn": schema.StringAttribute{
				Description: "Specifies the order in which pre-operation modify DN plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_pre_operation_search": schema.StringAttribute{
				Description: "Specifies the order in which pre-operation search plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_operation_abandon": schema.StringAttribute{
				Description: "Specifies the order in which post-operation abandon plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_operation_add": schema.StringAttribute{
				Description: "Specifies the order in which post-operation add plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_operation_bind": schema.StringAttribute{
				Description: "Specifies the order in which post-operation bind plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_operation_compare": schema.StringAttribute{
				Description: "Specifies the order in which post-operation compare plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_operation_delete": schema.StringAttribute{
				Description: "Specifies the order in which post-operation delete plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_operation_extended": schema.StringAttribute{
				Description: "Specifies the order in which post-operation extended operation plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_operation_modify": schema.StringAttribute{
				Description: "Specifies the order in which post-operation modify plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_operation_modify_dn": schema.StringAttribute{
				Description: "Specifies the order in which post-operation modify DN plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_operation_search": schema.StringAttribute{
				Description: "Specifies the order in which post-operation search plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_operation_unbind": schema.StringAttribute{
				Description: "Specifies the order in which post-operation unbind plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_response_add": schema.StringAttribute{
				Description: "Specifies the order in which post-response add plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_response_bind": schema.StringAttribute{
				Description: "Specifies the order in which post-response bind plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_response_compare": schema.StringAttribute{
				Description: "Specifies the order in which post-response compare plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_response_delete": schema.StringAttribute{
				Description: "Specifies the order in which post-response delete plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_response_extended": schema.StringAttribute{
				Description: "Specifies the order in which post-response extended operation plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_response_modify": schema.StringAttribute{
				Description: "Specifies the order in which post-response modify plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_response_modify_dn": schema.StringAttribute{
				Description: "Specifies the order in which post-response modify DN plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_synchronization_add": schema.StringAttribute{
				Description: "Specifies the order in which post-synchronization add plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_synchronization_delete": schema.StringAttribute{
				Description: "Specifies the order in which post-synchronization delete plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_synchronization_modify": schema.StringAttribute{
				Description: "Specifies the order in which post-synchronization modify plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_synchronization_modify_dn": schema.StringAttribute{
				Description: "Specifies the order in which post-synchronization modify DN plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_post_response_search": schema.StringAttribute{
				Description: "Specifies the order in which post-response search plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_search_result_entry": schema.StringAttribute{
				Description: "Specifies the order in which search result entry plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_search_result_reference": schema.StringAttribute{
				Description: "Specifies the order in which search result reference plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_subordinate_modify_dn": schema.StringAttribute{
				Description: "Specifies the order in which subordinate modify DN plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_order_intermediate_response": schema.StringAttribute{
				Description: "Specifies the order in which intermediate response plug-ins are to be loaded and invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a PluginRootResponse object into the model struct
func readPluginRootResponse(ctx context.Context, r *client.PluginRootResponse, state *pluginRootResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("plugin-root")
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
	state.PluginOrderStartup = internaltypes.StringTypeOrNil(r.PluginOrderStartup, true)
	state.PluginOrderShutdown = internaltypes.StringTypeOrNil(r.PluginOrderShutdown, true)
	state.PluginOrderPostConnect = internaltypes.StringTypeOrNil(r.PluginOrderPostConnect, true)
	state.PluginOrderPostDisconnect = internaltypes.StringTypeOrNil(r.PluginOrderPostDisconnect, true)
	state.PluginOrderLDIFImport = internaltypes.StringTypeOrNil(r.PluginOrderLDIFImport, true)
	state.PluginOrderLDIFExport = internaltypes.StringTypeOrNil(r.PluginOrderLDIFExport, true)
	state.PluginOrderPreParseAbandon = internaltypes.StringTypeOrNil(r.PluginOrderPreParseAbandon, true)
	state.PluginOrderPreParseAdd = internaltypes.StringTypeOrNil(r.PluginOrderPreParseAdd, true)
	state.PluginOrderPreParseBind = internaltypes.StringTypeOrNil(r.PluginOrderPreParseBind, true)
	state.PluginOrderPreParseCompare = internaltypes.StringTypeOrNil(r.PluginOrderPreParseCompare, true)
	state.PluginOrderPreParseDelete = internaltypes.StringTypeOrNil(r.PluginOrderPreParseDelete, true)
	state.PluginOrderPreParseExtended = internaltypes.StringTypeOrNil(r.PluginOrderPreParseExtended, true)
	state.PluginOrderPreParseModify = internaltypes.StringTypeOrNil(r.PluginOrderPreParseModify, true)
	state.PluginOrderPreParseModifyDN = internaltypes.StringTypeOrNil(r.PluginOrderPreParseModifyDN, true)
	state.PluginOrderPreParseSearch = internaltypes.StringTypeOrNil(r.PluginOrderPreParseSearch, true)
	state.PluginOrderPreParseUnbind = internaltypes.StringTypeOrNil(r.PluginOrderPreParseUnbind, true)
	state.PluginOrderPreOperationAdd = internaltypes.StringTypeOrNil(r.PluginOrderPreOperationAdd, true)
	state.PluginOrderPreOperationBind = internaltypes.StringTypeOrNil(r.PluginOrderPreOperationBind, true)
	state.PluginOrderPreOperationCompare = internaltypes.StringTypeOrNil(r.PluginOrderPreOperationCompare, true)
	state.PluginOrderPreOperationDelete = internaltypes.StringTypeOrNil(r.PluginOrderPreOperationDelete, true)
	state.PluginOrderPreOperationExtended = internaltypes.StringTypeOrNil(r.PluginOrderPreOperationExtended, true)
	state.PluginOrderPreOperationModify = internaltypes.StringTypeOrNil(r.PluginOrderPreOperationModify, true)
	state.PluginOrderPreOperationModifyDN = internaltypes.StringTypeOrNil(r.PluginOrderPreOperationModifyDN, true)
	state.PluginOrderPreOperationSearch = internaltypes.StringTypeOrNil(r.PluginOrderPreOperationSearch, true)
	state.PluginOrderPostOperationAbandon = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationAbandon, true)
	state.PluginOrderPostOperationAdd = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationAdd, true)
	state.PluginOrderPostOperationBind = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationBind, true)
	state.PluginOrderPostOperationCompare = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationCompare, true)
	state.PluginOrderPostOperationDelete = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationDelete, true)
	state.PluginOrderPostOperationExtended = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationExtended, true)
	state.PluginOrderPostOperationModify = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationModify, true)
	state.PluginOrderPostOperationModifyDN = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationModifyDN, true)
	state.PluginOrderPostOperationSearch = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationSearch, true)
	state.PluginOrderPostOperationUnbind = internaltypes.StringTypeOrNil(r.PluginOrderPostOperationUnbind, true)
	state.PluginOrderPostResponseAdd = internaltypes.StringTypeOrNil(r.PluginOrderPostResponseAdd, true)
	state.PluginOrderPostResponseBind = internaltypes.StringTypeOrNil(r.PluginOrderPostResponseBind, true)
	state.PluginOrderPostResponseCompare = internaltypes.StringTypeOrNil(r.PluginOrderPostResponseCompare, true)
	state.PluginOrderPostResponseDelete = internaltypes.StringTypeOrNil(r.PluginOrderPostResponseDelete, true)
	state.PluginOrderPostResponseExtended = internaltypes.StringTypeOrNil(r.PluginOrderPostResponseExtended, true)
	state.PluginOrderPostResponseModify = internaltypes.StringTypeOrNil(r.PluginOrderPostResponseModify, true)
	state.PluginOrderPostResponseModifyDN = internaltypes.StringTypeOrNil(r.PluginOrderPostResponseModifyDN, true)
	state.PluginOrderPostSynchronizationAdd = internaltypes.StringTypeOrNil(r.PluginOrderPostSynchronizationAdd, true)
	state.PluginOrderPostSynchronizationDelete = internaltypes.StringTypeOrNil(r.PluginOrderPostSynchronizationDelete, true)
	state.PluginOrderPostSynchronizationModify = internaltypes.StringTypeOrNil(r.PluginOrderPostSynchronizationModify, true)
	state.PluginOrderPostSynchronizationModifyDN = internaltypes.StringTypeOrNil(r.PluginOrderPostSynchronizationModifyDN, true)
	state.PluginOrderPostResponseSearch = internaltypes.StringTypeOrNil(r.PluginOrderPostResponseSearch, true)
	state.PluginOrderSearchResultEntry = internaltypes.StringTypeOrNil(r.PluginOrderSearchResultEntry, true)
	state.PluginOrderSearchResultReference = internaltypes.StringTypeOrNil(r.PluginOrderSearchResultReference, true)
	state.PluginOrderSubordinateModifyDN = internaltypes.StringTypeOrNil(r.PluginOrderSubordinateModifyDN, true)
	state.PluginOrderIntermediateResponse = internaltypes.StringTypeOrNil(r.PluginOrderIntermediateResponse, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createPluginRootOperations(plan pluginRootResourceModel, state pluginRootResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderStartup, state.PluginOrderStartup, "plugin-order-startup")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderShutdown, state.PluginOrderShutdown, "plugin-order-shutdown")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostConnect, state.PluginOrderPostConnect, "plugin-order-post-connect")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostDisconnect, state.PluginOrderPostDisconnect, "plugin-order-post-disconnect")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderLDIFImport, state.PluginOrderLDIFImport, "plugin-order-ldif-import")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderLDIFExport, state.PluginOrderLDIFExport, "plugin-order-ldif-export")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPreParseAbandon, state.PluginOrderPreParseAbandon, "plugin-order-pre-parse-abandon")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPreParseAdd, state.PluginOrderPreParseAdd, "plugin-order-pre-parse-add")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPreParseBind, state.PluginOrderPreParseBind, "plugin-order-pre-parse-bind")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPreParseCompare, state.PluginOrderPreParseCompare, "plugin-order-pre-parse-compare")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPreParseDelete, state.PluginOrderPreParseDelete, "plugin-order-pre-parse-delete")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPreParseExtended, state.PluginOrderPreParseExtended, "plugin-order-pre-parse-extended")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPreParseModify, state.PluginOrderPreParseModify, "plugin-order-pre-parse-modify")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPreParseModifyDN, state.PluginOrderPreParseModifyDN, "plugin-order-pre-parse-modify-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPreParseSearch, state.PluginOrderPreParseSearch, "plugin-order-pre-parse-search")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPreParseUnbind, state.PluginOrderPreParseUnbind, "plugin-order-pre-parse-unbind")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPreOperationAdd, state.PluginOrderPreOperationAdd, "plugin-order-pre-operation-add")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPreOperationBind, state.PluginOrderPreOperationBind, "plugin-order-pre-operation-bind")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPreOperationCompare, state.PluginOrderPreOperationCompare, "plugin-order-pre-operation-compare")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPreOperationDelete, state.PluginOrderPreOperationDelete, "plugin-order-pre-operation-delete")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPreOperationExtended, state.PluginOrderPreOperationExtended, "plugin-order-pre-operation-extended")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPreOperationModify, state.PluginOrderPreOperationModify, "plugin-order-pre-operation-modify")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPreOperationModifyDN, state.PluginOrderPreOperationModifyDN, "plugin-order-pre-operation-modify-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPreOperationSearch, state.PluginOrderPreOperationSearch, "plugin-order-pre-operation-search")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostOperationAbandon, state.PluginOrderPostOperationAbandon, "plugin-order-post-operation-abandon")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostOperationAdd, state.PluginOrderPostOperationAdd, "plugin-order-post-operation-add")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostOperationBind, state.PluginOrderPostOperationBind, "plugin-order-post-operation-bind")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostOperationCompare, state.PluginOrderPostOperationCompare, "plugin-order-post-operation-compare")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostOperationDelete, state.PluginOrderPostOperationDelete, "plugin-order-post-operation-delete")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostOperationExtended, state.PluginOrderPostOperationExtended, "plugin-order-post-operation-extended")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostOperationModify, state.PluginOrderPostOperationModify, "plugin-order-post-operation-modify")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostOperationModifyDN, state.PluginOrderPostOperationModifyDN, "plugin-order-post-operation-modify-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostOperationSearch, state.PluginOrderPostOperationSearch, "plugin-order-post-operation-search")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostOperationUnbind, state.PluginOrderPostOperationUnbind, "plugin-order-post-operation-unbind")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostResponseAdd, state.PluginOrderPostResponseAdd, "plugin-order-post-response-add")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostResponseBind, state.PluginOrderPostResponseBind, "plugin-order-post-response-bind")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostResponseCompare, state.PluginOrderPostResponseCompare, "plugin-order-post-response-compare")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostResponseDelete, state.PluginOrderPostResponseDelete, "plugin-order-post-response-delete")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostResponseExtended, state.PluginOrderPostResponseExtended, "plugin-order-post-response-extended")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostResponseModify, state.PluginOrderPostResponseModify, "plugin-order-post-response-modify")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostResponseModifyDN, state.PluginOrderPostResponseModifyDN, "plugin-order-post-response-modify-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostSynchronizationAdd, state.PluginOrderPostSynchronizationAdd, "plugin-order-post-synchronization-add")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostSynchronizationDelete, state.PluginOrderPostSynchronizationDelete, "plugin-order-post-synchronization-delete")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostSynchronizationModify, state.PluginOrderPostSynchronizationModify, "plugin-order-post-synchronization-modify")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostSynchronizationModifyDN, state.PluginOrderPostSynchronizationModifyDN, "plugin-order-post-synchronization-modify-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderPostResponseSearch, state.PluginOrderPostResponseSearch, "plugin-order-post-response-search")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderSearchResultEntry, state.PluginOrderSearchResultEntry, "plugin-order-search-result-entry")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderSearchResultReference, state.PluginOrderSearchResultReference, "plugin-order-search-result-reference")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderSubordinateModifyDN, state.PluginOrderSubordinateModifyDN, "plugin-order-subordinate-modify-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.PluginOrderIntermediateResponse, state.PluginOrderIntermediateResponse, "plugin-order-intermediate-response")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *pluginRootResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan pluginRootResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginRootApi.GetPluginRoot(
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

	// Read the existing configuration
	var state pluginRootResourceModel
	readPluginRootResponse(ctx, readResponse, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PluginRootApi.UpdatePluginRoot(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	ops := createPluginRootOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginRootApi.UpdatePluginRootExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Plugin Root", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPluginRootResponse(ctx, updateResponse, &state, &resp.Diagnostics)
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
func (r *pluginRootResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state pluginRootResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginRootApi.GetPluginRoot(
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
	readPluginRootResponse(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *pluginRootResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan pluginRootResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state pluginRootResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.PluginRootApi.UpdatePluginRoot(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))

	// Determine what update operations are necessary
	ops := createPluginRootOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginRootApi.UpdatePluginRootExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Plugin Root", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPluginRootResponse(ctx, updateResponse, &state, &resp.Diagnostics)
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
func (r *pluginRootResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *pluginRootResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Set a placeholder id value to appease terraform.
	// The real attributes will be imported when terraform performs a read after the import.
	// If no value is set here, Terraform will error out when importing.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), "id")...)
}
