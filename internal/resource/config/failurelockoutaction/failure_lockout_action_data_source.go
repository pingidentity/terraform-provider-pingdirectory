package failurelockoutaction

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
	_ datasource.DataSource              = &failureLockoutActionDataSource{}
	_ datasource.DataSourceWithConfigure = &failureLockoutActionDataSource{}
)

// Create a Failure Lockout Action data source
func NewFailureLockoutActionDataSource() datasource.DataSource {
	return &failureLockoutActionDataSource{}
}

// failureLockoutActionDataSource is the datasource implementation.
type failureLockoutActionDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *failureLockoutActionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_failure_lockout_action"
}

// Configure adds the provider configured client to the data source.
func (r *failureLockoutActionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type failureLockoutActionDataSourceModel struct {
	Id                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	Type                              types.String `tfsdk:"type"`
	Delay                             types.String `tfsdk:"delay"`
	AllowBlockingDelay                types.Bool   `tfsdk:"allow_blocking_delay"`
	GenerateAccountStatusNotification types.Bool   `tfsdk:"generate_account_status_notification"`
	Description                       types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the datasource.
func (r *failureLockoutActionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Failure Lockout Action.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Failure Lockout Action resource. Options are ['delay-bind-response', 'no-operation', 'lock-account']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"delay": schema.StringAttribute{
				Description: "The length of time to delay the bind response for accounts with too many failed authentication attempts.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_blocking_delay": schema.BoolAttribute{
				Description: "Indicates whether to delay the response for authentication attempts even if that delay may block the thread being used to process the attempt.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"generate_account_status_notification": schema.BoolAttribute{
				Description: " When the `type` value is one of [`delay-bind-response`]: Indicates whether to generate an account status notification for cases in which a bind response is delayed because of failure lockout. When the `type` value is one of [`no-operation`]: Indicates whether to generate an account status notification for cases in which this failure lockout action is invoked for a bind attempt with too many outstanding authentication failures.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Failure Lockout Action",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a DelayBindResponseFailureLockoutActionResponse object into the model struct
func readDelayBindResponseFailureLockoutActionResponseDataSource(ctx context.Context, r *client.DelayBindResponseFailureLockoutActionResponse, state *failureLockoutActionDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("delay-bind-response")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Delay = types.StringValue(r.Delay)
	state.AllowBlockingDelay = internaltypes.BoolTypeOrNil(r.AllowBlockingDelay)
	state.GenerateAccountStatusNotification = internaltypes.BoolTypeOrNil(r.GenerateAccountStatusNotification)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a NoOperationFailureLockoutActionResponse object into the model struct
func readNoOperationFailureLockoutActionResponseDataSource(ctx context.Context, r *client.NoOperationFailureLockoutActionResponse, state *failureLockoutActionDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("no-operation")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.GenerateAccountStatusNotification = internaltypes.BoolTypeOrNil(r.GenerateAccountStatusNotification)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a LockAccountFailureLockoutActionResponse object into the model struct
func readLockAccountFailureLockoutActionResponseDataSource(ctx context.Context, r *client.LockAccountFailureLockoutActionResponse, state *failureLockoutActionDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("lock-account")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read resource information
func (r *failureLockoutActionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state failureLockoutActionDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.FailureLockoutActionApi.GetFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Failure Lockout Action", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.DelayBindResponseFailureLockoutActionResponse != nil {
		readDelayBindResponseFailureLockoutActionResponseDataSource(ctx, readResponse.DelayBindResponseFailureLockoutActionResponse, &state, &resp.Diagnostics)
	}
	if readResponse.NoOperationFailureLockoutActionResponse != nil {
		readNoOperationFailureLockoutActionResponseDataSource(ctx, readResponse.NoOperationFailureLockoutActionResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LockAccountFailureLockoutActionResponse != nil {
		readLockAccountFailureLockoutActionResponseDataSource(ctx, readResponse.LockAccountFailureLockoutActionResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
