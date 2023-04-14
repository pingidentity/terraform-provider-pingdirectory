package externalserver

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &smtpExternalServerResource{}
	_ resource.ResourceWithConfigure   = &smtpExternalServerResource{}
	_ resource.ResourceWithImportState = &smtpExternalServerResource{}
	_ resource.Resource                = &defaultSmtpExternalServerResource{}
	_ resource.ResourceWithConfigure   = &defaultSmtpExternalServerResource{}
	_ resource.ResourceWithImportState = &defaultSmtpExternalServerResource{}
)

// Create a Smtp External Server resource
func NewSmtpExternalServerResource() resource.Resource {
	return &smtpExternalServerResource{}
}

func NewDefaultSmtpExternalServerResource() resource.Resource {
	return &defaultSmtpExternalServerResource{}
}

// smtpExternalServerResource is the resource implementation.
type smtpExternalServerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSmtpExternalServerResource is the resource implementation.
type defaultSmtpExternalServerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *smtpExternalServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smtp_external_server"
}

func (r *defaultSmtpExternalServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_smtp_external_server"
}

// Configure adds the provider configured client to the resource.
func (r *smtpExternalServerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultSmtpExternalServerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type smtpExternalServerResourceModel struct {
	Id                       types.String `tfsdk:"id"`
	LastUpdated              types.String `tfsdk:"last_updated"`
	Notifications            types.Set    `tfsdk:"notifications"`
	RequiredActions          types.Set    `tfsdk:"required_actions"`
	ServerHostName           types.String `tfsdk:"server_host_name"`
	ServerPort               types.Int64  `tfsdk:"server_port"`
	SmtpSecurity             types.String `tfsdk:"smtp_security"`
	UserName                 types.String `tfsdk:"user_name"`
	Password                 types.String `tfsdk:"password"`
	PassphraseProvider       types.String `tfsdk:"passphrase_provider"`
	SmtpTimeout              types.String `tfsdk:"smtp_timeout"`
	SmtpConnectionProperties types.Set    `tfsdk:"smtp_connection_properties"`
	Description              types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *smtpExternalServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	smtpExternalServerSchema(ctx, req, resp, false)
}

func (r *defaultSmtpExternalServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	smtpExternalServerSchema(ctx, req, resp, true)
}

func smtpExternalServerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Smtp External Server.",
		Attributes: map[string]schema.Attribute{
			"server_host_name": schema.StringAttribute{
				Description: "The host name of the smtp server.",
				Required:    true,
			},
			"server_port": schema.Int64Attribute{
				Description: "The port number where the smtp server listens for requests.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"smtp_security": schema.StringAttribute{
				Description: "This property specifies type of connection security to use when connecting to the outgoing mail server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_name": schema.StringAttribute{
				Description: "The name of the login account to use when connecting to the smtp server. Both username and password must be supplied if this attribute is set.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "The login password for the specified user name. Both username and password must be supplied if this attribute is set.",
				Optional:    true,
				Sensitive:   true,
			},
			"passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider to use to obtain the login password for the specified user.",
				Optional:    true,
			},
			"smtp_timeout": schema.StringAttribute{
				Description: "Specifies the maximum length of time that a connection or attempted connection to a SMTP server may take.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"smtp_connection_properties": schema.SetAttribute{
				Description: "Specifies the connection properties for the smtp server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this External Server",
				Optional:    true,
			},
		},
	}
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalSmtpExternalServerFields(ctx context.Context, addRequest *client.AddSmtpExternalServerRequest, plan smtpExternalServerResourceModel) error {
	if internaltypes.IsDefined(plan.ServerPort) {
		addRequest.ServerPort = plan.ServerPort.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SmtpSecurity) {
		smtpSecurity, err := client.NewEnumexternalServerSmtpSecurityPropFromValue(plan.SmtpSecurity.ValueString())
		if err != nil {
			return err
		}
		addRequest.SmtpSecurity = smtpSecurity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UserName) {
		addRequest.UserName = plan.UserName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Password) {
		addRequest.Password = plan.Password.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PassphraseProvider) {
		addRequest.PassphraseProvider = plan.PassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SmtpTimeout) {
		addRequest.SmtpTimeout = plan.SmtpTimeout.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.SmtpConnectionProperties) {
		var slice []string
		plan.SmtpConnectionProperties.ElementsAs(ctx, &slice, false)
		addRequest.SmtpConnectionProperties = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Read a SmtpExternalServerResponse object into the model struct
func readSmtpExternalServerResponse(ctx context.Context, r *client.SmtpExternalServerResponse, state *smtpExternalServerResourceModel, expectedValues *smtpExternalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = internaltypes.Int64TypeOrNil(r.ServerPort)
	state.SmtpSecurity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumexternalServerSmtpSecurityProp(r.SmtpSecurity), internaltypes.IsEmptyString(expectedValues.SmtpSecurity))
	state.UserName = internaltypes.StringTypeOrNil(r.UserName, internaltypes.IsEmptyString(expectedValues.UserName))
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.Password = expectedValues.Password
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, internaltypes.IsEmptyString(expectedValues.PassphraseProvider))
	state.SmtpTimeout = internaltypes.StringTypeOrNil(r.SmtpTimeout, internaltypes.IsEmptyString(expectedValues.SmtpTimeout))
	config.CheckMismatchedPDFormattedAttributes("smtp_timeout",
		expectedValues.SmtpTimeout, state.SmtpTimeout, diagnostics)
	state.SmtpConnectionProperties = internaltypes.GetStringSet(r.SmtpConnectionProperties)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSmtpExternalServerOperations(plan smtpExternalServerResourceModel, state smtpExternalServerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ServerHostName, state.ServerHostName, "server-host-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.ServerPort, state.ServerPort, "server-port")
	operations.AddStringOperationIfNecessary(&ops, plan.SmtpSecurity, state.SmtpSecurity, "smtp-security")
	operations.AddStringOperationIfNecessary(&ops, plan.UserName, state.UserName, "user-name")
	operations.AddStringOperationIfNecessary(&ops, plan.Password, state.Password, "password")
	operations.AddStringOperationIfNecessary(&ops, plan.PassphraseProvider, state.PassphraseProvider, "passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.SmtpTimeout, state.SmtpTimeout, "smtp-timeout")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SmtpConnectionProperties, state.SmtpConnectionProperties, "smtp-connection-properties")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *smtpExternalServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan smtpExternalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddSmtpExternalServerRequest(plan.Id.ValueString(),
		[]client.EnumsmtpExternalServerSchemaUrn{client.ENUMSMTPEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERSMTP},
		plan.ServerHostName.ValueString())
	err := addOptionalSmtpExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Smtp External Server", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerApi.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddSmtpExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerApi.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Smtp External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state smtpExternalServerResourceModel
	readSmtpExternalServerResponse(ctx, addResponse.SmtpExternalServerResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultSmtpExternalServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan smtpExternalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ExternalServerApi.GetExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Smtp External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state smtpExternalServerResourceModel
	readSmtpExternalServerResponse(ctx, readResponse.SmtpExternalServerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ExternalServerApi.UpdateExternalServer(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSmtpExternalServerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ExternalServerApi.UpdateExternalServerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Smtp External Server", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSmtpExternalServerResponse(ctx, updateResponse.SmtpExternalServerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *smtpExternalServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSmtpExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSmtpExternalServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSmtpExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSmtpExternalServer(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state smtpExternalServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ExternalServerApi.GetExternalServer(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Smtp External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSmtpExternalServerResponse(ctx, readResponse.SmtpExternalServerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *smtpExternalServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSmtpExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSmtpExternalServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSmtpExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSmtpExternalServer(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan smtpExternalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state smtpExternalServerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ExternalServerApi.UpdateExternalServer(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSmtpExternalServerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ExternalServerApi.UpdateExternalServerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Smtp External Server", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSmtpExternalServerResponse(ctx, updateResponse.SmtpExternalServerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSmtpExternalServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *smtpExternalServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state smtpExternalServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ExternalServerApi.DeleteExternalServerExecute(r.apiClient.ExternalServerApi.DeleteExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Smtp External Server", err, httpResp)
		return
	}
}

func (r *smtpExternalServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSmtpExternalServer(ctx, req, resp)
}

func (r *defaultSmtpExternalServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSmtpExternalServer(ctx, req, resp)
}

func importSmtpExternalServer(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
