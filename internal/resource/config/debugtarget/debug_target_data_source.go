// Copyright © 2025 Ping Identity Corporation

package debugtarget

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
	_ datasource.DataSource              = &debugTargetDataSource{}
	_ datasource.DataSourceWithConfigure = &debugTargetDataSource{}
)

// Create a Debug Target data source
func NewDebugTargetDataSource() datasource.DataSource {
	return &debugTargetDataSource{}
}

// debugTargetDataSource is the datasource implementation.
type debugTargetDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *debugTargetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_debug_target"
}

// Configure adds the provider configured client to the data source.
func (r *debugTargetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type debugTargetDataSourceModel struct {
	Id                       types.String `tfsdk:"id"`
	Type                     types.String `tfsdk:"type"`
	LogPublisherName         types.String `tfsdk:"log_publisher_name"`
	DebugScope               types.String `tfsdk:"debug_scope"`
	DebugLevel               types.String `tfsdk:"debug_level"`
	DebugCategory            types.Set    `tfsdk:"debug_category"`
	OmitMethodEntryArguments types.Bool   `tfsdk:"omit_method_entry_arguments"`
	OmitMethodReturnValue    types.Bool   `tfsdk:"omit_method_return_value"`
	IncludeThrowableCause    types.Bool   `tfsdk:"include_throwable_cause"`
	ThrowableStackFrames     types.Int64  `tfsdk:"throwable_stack_frames"`
	Description              types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the datasource.
func (r *debugTargetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Debug Target.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Debug Target resource. Options are ['debug-target']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_publisher_name": schema.StringAttribute{
				Description: "Name of the parent Log Publisher",
				Required:    true,
			},
			"debug_scope": schema.StringAttribute{
				Description: "Specifies the fully-qualified Java package, class, or method affected by the settings in this target definition. Use the number character (#) to separate the class name and the method name (that is, com.unboundid.directory.server.core.DirectoryServer#startUp).",
				Required:    true,
			},
			"debug_level": schema.StringAttribute{
				Description: "Specifies the lowest severity level of debug messages to log.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"debug_category": schema.SetAttribute{
				Description: "Specifies the debug message categories to be logged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"omit_method_entry_arguments": schema.BoolAttribute{
				Description: "Specifies the property to indicate whether to include method arguments in debug messages.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"omit_method_return_value": schema.BoolAttribute{
				Description: "Specifies the property to indicate whether to include the return value in debug messages.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_throwable_cause": schema.BoolAttribute{
				Description: "Specifies the property to indicate whether to include the cause of exceptions in exception thrown and caught messages.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"throwable_stack_frames": schema.Int64Attribute{
				Description: "Specifies the property to indicate the number of stack frames to include in the stack trace for method entry and exception thrown messages.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Debug Target",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a DebugTargetResponse object into the model struct
func readDebugTargetResponseDataSource(ctx context.Context, r *client.DebugTargetResponse, state *debugTargetDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("debug-target")
	state.Id = types.StringValue(r.Id)
	state.DebugScope = types.StringValue(r.DebugScope)
	state.DebugLevel = types.StringValue(r.DebugLevel.String())
	state.DebugCategory = internaltypes.GetStringSet(
		client.StringSliceEnumdebugTargetDebugCategoryProp(r.DebugCategory))
	state.OmitMethodEntryArguments = internaltypes.BoolTypeOrNil(r.OmitMethodEntryArguments)
	state.OmitMethodReturnValue = internaltypes.BoolTypeOrNil(r.OmitMethodReturnValue)
	state.IncludeThrowableCause = internaltypes.BoolTypeOrNil(r.IncludeThrowableCause)
	state.ThrowableStackFrames = internaltypes.Int64TypeOrNil(r.ThrowableStackFrames)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read resource information
func (r *debugTargetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state debugTargetDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DebugTargetAPI.GetDebugTarget(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.DebugScope.ValueString(), state.LogPublisherName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Debug Target", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDebugTargetResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
