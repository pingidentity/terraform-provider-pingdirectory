package changesubscriptionhandler

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
	_ datasource.DataSource              = &changeSubscriptionHandlerDataSource{}
	_ datasource.DataSourceWithConfigure = &changeSubscriptionHandlerDataSource{}
)

// Create a Change Subscription Handler data source
func NewChangeSubscriptionHandlerDataSource() datasource.DataSource {
	return &changeSubscriptionHandlerDataSource{}
}

// changeSubscriptionHandlerDataSource is the datasource implementation.
type changeSubscriptionHandlerDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *changeSubscriptionHandlerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_change_subscription_handler"
}

// Configure adds the provider configured client to the data source.
func (r *changeSubscriptionHandlerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type changeSubscriptionHandlerDataSourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Type               types.String `tfsdk:"type"`
	ExtensionClass     types.String `tfsdk:"extension_class"`
	ExtensionArgument  types.Set    `tfsdk:"extension_argument"`
	LogFile            types.String `tfsdk:"log_file"`
	ScriptClass        types.String `tfsdk:"script_class"`
	ScriptArgument     types.Set    `tfsdk:"script_argument"`
	Description        types.String `tfsdk:"description"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	ChangeSubscription types.Set    `tfsdk:"change_subscription"`
}

// GetSchema defines the schema for the datasource.
func (r *changeSubscriptionHandlerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Change Subscription Handler.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Change Subscription Handler resource. Options are ['groovy-scripted', 'logging', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Change Subscription Handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Change Subscription Handler. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"log_file": schema.StringAttribute{
				Description: "Specifies the log file in which the change notification messages will be written.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Change Subscription Handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Change Subscription Handler. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Change Subscription Handler",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this change subscription handler is enabled within the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"change_subscription": schema.SetAttribute{
				Description: "The set of change subscriptions for which this change subscription handler should be notified. If no values are provided then it will be notified for all change subscriptions defined in the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a GroovyScriptedChangeSubscriptionHandlerResponse object into the model struct
func readGroovyScriptedChangeSubscriptionHandlerResponseDataSource(ctx context.Context, r *client.GroovyScriptedChangeSubscriptionHandlerResponse, state *changeSubscriptionHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ChangeSubscription = internaltypes.GetStringSet(r.ChangeSubscription)
}

// Read a LoggingChangeSubscriptionHandlerResponse object into the model struct
func readLoggingChangeSubscriptionHandlerResponseDataSource(ctx context.Context, r *client.LoggingChangeSubscriptionHandlerResponse, state *changeSubscriptionHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("logging")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ChangeSubscription = internaltypes.GetStringSet(r.ChangeSubscription)
}

// Read a ThirdPartyChangeSubscriptionHandlerResponse object into the model struct
func readThirdPartyChangeSubscriptionHandlerResponseDataSource(ctx context.Context, r *client.ThirdPartyChangeSubscriptionHandlerResponse, state *changeSubscriptionHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ChangeSubscription = internaltypes.GetStringSet(r.ChangeSubscription)
}

// Read resource information
func (r *changeSubscriptionHandlerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state changeSubscriptionHandlerDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ChangeSubscriptionHandlerAPI.GetChangeSubscriptionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Change Subscription Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.GroovyScriptedChangeSubscriptionHandlerResponse != nil {
		readGroovyScriptedChangeSubscriptionHandlerResponseDataSource(ctx, readResponse.GroovyScriptedChangeSubscriptionHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LoggingChangeSubscriptionHandlerResponse != nil {
		readLoggingChangeSubscriptionHandlerResponseDataSource(ctx, readResponse.LoggingChangeSubscriptionHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyChangeSubscriptionHandlerResponse != nil {
		readThirdPartyChangeSubscriptionHandlerResponseDataSource(ctx, readResponse.ThirdPartyChangeSubscriptionHandlerResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
