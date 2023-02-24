package externalserver

import (
	"context"
	"time"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &syslogExternalServerResource{}
	_ resource.ResourceWithConfigure   = &syslogExternalServerResource{}
	_ resource.ResourceWithImportState = &syslogExternalServerResource{}
)

// Create a Syslog External Server resource
func NewSyslogExternalServerResource() resource.Resource {
	return &syslogExternalServerResource{}
}

// syslogExternalServerResource is the resource implementation.
type syslogExternalServerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *syslogExternalServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_syslog_external_server"
}

// Configure adds the provider configured client to the resource.
func (r *syslogExternalServerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type syslogExternalServerResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	LastUpdated          types.String `tfsdk:"last_updated"`
	Notifications        types.Set    `tfsdk:"notifications"`
	RequiredActions      types.Set    `tfsdk:"required_actions"`
	ServerHostName       types.String `tfsdk:"server_host_name"`
	ServerPort           types.Int64  `tfsdk:"server_port"`
	TransportMechanism   types.String `tfsdk:"transport_mechanism"`
	ConnectTimeout       types.String `tfsdk:"connect_timeout"`
	MaxConnectionAge     types.String `tfsdk:"max_connection_age"`
	TrustManagerProvider types.String `tfsdk:"trust_manager_provider"`
	Description          types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *syslogExternalServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Syslog External Server.",
		Attributes: map[string]schema.Attribute{
			"server_host_name": schema.StringAttribute{
				Description: "The address of the syslog server.",
				Required:    true,
			},
			"server_port": schema.Int64Attribute{
				Description: "The port on which the syslog server accepts connections.",
				Optional:    true,
				Computed:    true,
			},
			"transport_mechanism": schema.StringAttribute{
				Description: "The transport mechanism that should be used when communicating with the syslog server.",
				Required:    true,
			},
			"connect_timeout": schema.StringAttribute{
				Description: "Specifies the maximum length of time to wait for a connection to be established before giving up and considering the server unavailable. This will only be used when communicating with the syslog server over TCP (with or without TLS encryption).",
				Optional:    true,
				Computed:    true,
			},
			"max_connection_age": schema.StringAttribute{
				Description: "The maximum length of time that TCP connections should remain established. This will be ignored for UDP-based connections. A zero duration indicates that no maximum age will be imposed.",
				Optional:    true,
				Computed:    true,
			},
			"trust_manager_provider": schema.StringAttribute{
				Description: "A trust manager provider that will be used to determine whether to trust the certificate chain presented by the syslog server when communication is encrypted with TLS. This property will be ignored when not using TLS encryption.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this External Server",
				Optional:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalSyslogExternalServerFields(ctx context.Context, addRequest *client.AddSyslogExternalServerRequest, plan syslogExternalServerResourceModel) {
	if internaltypes.IsDefined(plan.ServerPort) {
		intVal := int32(plan.ServerPort.ValueInt64())
		addRequest.ServerPort = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectTimeout) {
		stringVal := plan.ConnectTimeout.ValueString()
		addRequest.ConnectTimeout = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxConnectionAge) {
		stringVal := plan.MaxConnectionAge.ValueString()
		addRequest.MaxConnectionAge = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustManagerProvider) {
		stringVal := plan.TrustManagerProvider.ValueString()
		addRequest.TrustManagerProvider = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
}

// Read a SyslogExternalServerResponse object into the model struct
func readSyslogExternalServerResponse(ctx context.Context, r *client.SyslogExternalServerResponse, state *syslogExternalServerResourceModel, expectedValues *syslogExternalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = internaltypes.Int64TypeOrNil(r.ServerPort)
	state.TransportMechanism = types.StringValue(r.TransportMechanism.String())
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	config.CheckMismatchedPDFormattedAttributes("connect_timeout",
		expectedValues.ConnectTimeout, state.ConnectTimeout, diagnostics)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	config.CheckMismatchedPDFormattedAttributes("max_connection_age",
		expectedValues.MaxConnectionAge, state.MaxConnectionAge, diagnostics)
	state.TrustManagerProvider = types.StringValue(r.TrustManagerProvider)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSyslogExternalServerOperations(plan syslogExternalServerResourceModel, state syslogExternalServerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ServerHostName, state.ServerHostName, "server-host-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.ServerPort, state.ServerPort, "server-port")
	operations.AddStringOperationIfNecessary(&ops, plan.TransportMechanism, state.TransportMechanism, "transport-mechanism")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectTimeout, state.ConnectTimeout, "connect-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.MaxConnectionAge, state.MaxConnectionAge, "max-connection-age")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustManagerProvider, state.TrustManagerProvider, "trust-manager-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *syslogExternalServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan syslogExternalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	transportMechanism, err := client.NewEnumexternalServerTransportMechanismPropFromValue(plan.TransportMechanism.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for TransportMechanism", err.Error())
		return
	}
	addRequest := client.NewAddSyslogExternalServerRequest(plan.Id.ValueString(),
		[]client.EnumsyslogExternalServerSchemaUrn{client.ENUMSYSLOGEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERSYSLOG},
		plan.ServerHostName.ValueString(),
		*transportMechanism)
	addOptionalSyslogExternalServerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerApi.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddSyslogExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerApi.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Syslog External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state syslogExternalServerResourceModel
	readSyslogExternalServerResponse(ctx, addResponse.SyslogExternalServerResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *syslogExternalServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state syslogExternalServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ExternalServerApi.GetExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Syslog External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSyslogExternalServerResponse(ctx, readResponse.SyslogExternalServerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *syslogExternalServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan syslogExternalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state syslogExternalServerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.ExternalServerApi.UpdateExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSyslogExternalServerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ExternalServerApi.UpdateExternalServerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Syslog External Server", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSyslogExternalServerResponse(ctx, updateResponse.SyslogExternalServerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *syslogExternalServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state syslogExternalServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ExternalServerApi.DeleteExternalServerExecute(r.apiClient.ExternalServerApi.DeleteExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Syslog External Server", err, httpResp)
		return
	}
}

func (r *syslogExternalServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
