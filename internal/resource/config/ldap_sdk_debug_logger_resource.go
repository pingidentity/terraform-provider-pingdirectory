package config

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ldapSdkDebugLoggerResource{}
	_ resource.ResourceWithConfigure   = &ldapSdkDebugLoggerResource{}
	_ resource.ResourceWithImportState = &ldapSdkDebugLoggerResource{}
)

// Create a Ldap Sdk Debug Logger resource
func NewLdapSdkDebugLoggerResource() resource.Resource {
	return &ldapSdkDebugLoggerResource{}
}

// ldapSdkDebugLoggerResource is the resource implementation.
type ldapSdkDebugLoggerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *ldapSdkDebugLoggerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_ldap_sdk_debug_logger"
}

// Configure adds the provider configured client to the resource.
func (r *ldapSdkDebugLoggerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type ldapSdkDebugLoggerResourceModel struct {
	// Id field required for acceptance testing framework
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
}

type defaultLdapSdkDebugLoggerResourceModel struct {
	// Id field required for acceptance testing framework
	Id                             types.String `tfsdk:"id"`
	LastUpdated                    types.String `tfsdk:"last_updated"`
	Notifications                  types.Set    `tfsdk:"notifications"`
	RequiredActions                types.Set    `tfsdk:"required_actions"`
	Description                    types.String `tfsdk:"description"`
	Enabled                        types.Bool   `tfsdk:"enabled"`
	LogFile                        types.String `tfsdk:"log_file"`
	DebugLevel                     types.String `tfsdk:"debug_level"`
	DebugType                      types.Set    `tfsdk:"debug_type"`
	IncludeStackTrace              types.Bool   `tfsdk:"include_stack_trace"`
	LogFilePermissions             types.String `tfsdk:"log_file_permissions"`
	TimeInterval                   types.String `tfsdk:"time_interval"`
	AutoFlush                      types.Bool   `tfsdk:"auto_flush"`
	Asynchronous                   types.Bool   `tfsdk:"asynchronous"`
	QueueSize                      types.Int64  `tfsdk:"queue_size"`
	BufferSize                     types.String `tfsdk:"buffer_size"`
	Append                         types.Bool   `tfsdk:"append"`
	RotationPolicy                 types.Set    `tfsdk:"rotation_policy"`
	RotationListener               types.Set    `tfsdk:"rotation_listener"`
	RetentionPolicy                types.Set    `tfsdk:"retention_policy"`
	CompressionMechanism           types.String `tfsdk:"compression_mechanism"`
	SignLog                        types.Bool   `tfsdk:"sign_log"`
	EncryptLog                     types.Bool   `tfsdk:"encrypt_log"`
	EncryptionSettingsDefinitionID types.String `tfsdk:"encryption_settings_definition_id"`
	TimestampPrecision             types.String `tfsdk:"timestamp_precision"`
	LoggingErrorBehavior           types.String `tfsdk:"logging_error_behavior"`
}

// GetSchema defines the schema for the resource.
func (r *ldapSdkDebugLoggerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a Ldap Sdk Debug Logger.",
		Attributes:  map[string]schema.Attribute{},
	}
	AddCommonSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a LdapSdkDebugLoggerResponse object into the model struct
func readLdapSdkDebugLoggerResponseDefault(ctx context.Context, r *client.LdapSdkDebugLoggerResponse, state *defaultLdapSdkDebugLoggerResourceModel, expectedValues *defaultLdapSdkDebugLoggerResourceModel, diagnostics *diag.Diagnostics) {
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LogFile = types.StringValue(r.LogFile)
	state.DebugLevel = types.StringValue(r.DebugLevel.String())
	state.DebugType = internaltypes.GetStringSet(
		client.StringSliceEnumldapSdkDebugLoggerDebugTypeProp(r.DebugType))
	state.IncludeStackTrace = types.BoolValue(r.IncludeStackTrace)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, true)
	CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, true)
	CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumldapSdkDebugLoggerCompressionMechanismProp(r.CompressionMechanism), true)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, true)
	state.TimestampPrecision = internaltypes.StringTypeOrNil(
		client.StringPointerEnumldapSdkDebugLoggerTimestampPrecisionProp(r.TimestampPrecision), true)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumldapSdkDebugLoggerLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createLdapSdkDebugLoggerOperations(plan ldapSdkDebugLoggerResourceModel, state ldapSdkDebugLoggerResourceModel) []client.Operation {
	var ops []client.Operation
	return ops
}

// Create any update operations necessary to make the state match the plan
func createLdapSdkDebugLoggerOperationsDefault(plan defaultLdapSdkDebugLoggerResourceModel, state defaultLdapSdkDebugLoggerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFile, state.LogFile, "log-file")
	operations.AddStringOperationIfNecessary(&ops, plan.DebugLevel, state.DebugLevel, "debug-level")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DebugType, state.DebugType, "debug-type")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeStackTrace, state.IncludeStackTrace, "include-stack-trace")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFilePermissions, state.LogFilePermissions, "log-file-permissions")
	operations.AddStringOperationIfNecessary(&ops, plan.TimeInterval, state.TimeInterval, "time-interval")
	operations.AddBoolOperationIfNecessary(&ops, plan.AutoFlush, state.AutoFlush, "auto-flush")
	operations.AddBoolOperationIfNecessary(&ops, plan.Asynchronous, state.Asynchronous, "asynchronous")
	operations.AddInt64OperationIfNecessary(&ops, plan.QueueSize, state.QueueSize, "queue-size")
	operations.AddStringOperationIfNecessary(&ops, plan.BufferSize, state.BufferSize, "buffer-size")
	operations.AddBoolOperationIfNecessary(&ops, plan.Append, state.Append, "append")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RotationPolicy, state.RotationPolicy, "rotation-policy")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RotationListener, state.RotationListener, "rotation-listener")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RetentionPolicy, state.RetentionPolicy, "retention-policy")
	operations.AddStringOperationIfNecessary(&ops, plan.CompressionMechanism, state.CompressionMechanism, "compression-mechanism")
	operations.AddBoolOperationIfNecessary(&ops, plan.SignLog, state.SignLog, "sign-log")
	operations.AddBoolOperationIfNecessary(&ops, plan.EncryptLog, state.EncryptLog, "encrypt-log")
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptionSettingsDefinitionID, state.EncryptionSettingsDefinitionID, "encryption-settings-definition-id")
	operations.AddStringOperationIfNecessary(&ops, plan.TimestampPrecision, state.TimestampPrecision, "timestamp-precision")
	operations.AddStringOperationIfNecessary(&ops, plan.LoggingErrorBehavior, state.LoggingErrorBehavior, "logging-error-behavior")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *ldapSdkDebugLoggerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan defaultLdapSdkDebugLoggerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LdapSdkDebugLoggerApi.GetLdapSdkDebugLogger(
		ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ldap Sdk Debug Logger", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state defaultLdapSdkDebugLoggerResourceModel
	readLdapSdkDebugLoggerResponseDefault(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LdapSdkDebugLoggerApi.UpdateLdapSdkDebugLogger(ProviderBasicAuthContext(ctx, r.providerConfig))
	ops := createLdapSdkDebugLoggerOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LdapSdkDebugLoggerApi.UpdateLdapSdkDebugLoggerExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ldap Sdk Debug Logger", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLdapSdkDebugLoggerResponseDefault(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *ldapSdkDebugLoggerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ldapSdkDebugLoggerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LdapSdkDebugLoggerApi.GetLdapSdkDebugLogger(
		ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ldap Sdk Debug Logger", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLdapSdkDebugLoggerResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *ldapSdkDebugLoggerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan ldapSdkDebugLoggerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state ldapSdkDebugLoggerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.LdapSdkDebugLoggerApi.UpdateLdapSdkDebugLogger(
		ProviderBasicAuthContext(ctx, r.providerConfig))

	// Determine what update operations are necessary
	ops := createLdapSdkDebugLoggerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LdapSdkDebugLoggerApi.UpdateLdapSdkDebugLoggerExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ldap Sdk Debug Logger", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLdapSdkDebugLoggerResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *ldapSdkDebugLoggerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *ldapSdkDebugLoggerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Set a placeholder id value to appease terraform.
	// The real attributes will be imported when terraform performs a read after the import.
	// If no value is set here, Terraform will error out when importing.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), "id")...)
}
