package alerthandler

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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
	_ resource.Resource                = &smtpAlertHandlerResource{}
	_ resource.ResourceWithConfigure   = &smtpAlertHandlerResource{}
	_ resource.ResourceWithImportState = &smtpAlertHandlerResource{}
	_ resource.Resource                = &defaultSmtpAlertHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultSmtpAlertHandlerResource{}
	_ resource.ResourceWithImportState = &defaultSmtpAlertHandlerResource{}
)

// Create a Smtp Alert Handler resource
func NewSmtpAlertHandlerResource() resource.Resource {
	return &smtpAlertHandlerResource{}
}

func NewDefaultSmtpAlertHandlerResource() resource.Resource {
	return &defaultSmtpAlertHandlerResource{}
}

// smtpAlertHandlerResource is the resource implementation.
type smtpAlertHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSmtpAlertHandlerResource is the resource implementation.
type defaultSmtpAlertHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *smtpAlertHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smtp_alert_handler"
}

func (r *defaultSmtpAlertHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_smtp_alert_handler"
}

// Configure adds the provider configured client to the resource.
func (r *smtpAlertHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultSmtpAlertHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type smtpAlertHandlerResourceModel struct {
	Id                       types.String `tfsdk:"id"`
	LastUpdated              types.String `tfsdk:"last_updated"`
	Notifications            types.Set    `tfsdk:"notifications"`
	RequiredActions          types.Set    `tfsdk:"required_actions"`
	Asynchronous             types.Bool   `tfsdk:"asynchronous"`
	SenderAddress            types.String `tfsdk:"sender_address"`
	RecipientAddress         types.Set    `tfsdk:"recipient_address"`
	MessageSubject           types.String `tfsdk:"message_subject"`
	MessageBody              types.String `tfsdk:"message_body"`
	IncludeMonitorDataFilter types.String `tfsdk:"include_monitor_data_filter"`
	Description              types.String `tfsdk:"description"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	EnabledAlertSeverity     types.Set    `tfsdk:"enabled_alert_severity"`
	EnabledAlertType         types.Set    `tfsdk:"enabled_alert_type"`
	DisabledAlertType        types.Set    `tfsdk:"disabled_alert_type"`
}

// GetSchema defines the schema for the resource.
func (r *smtpAlertHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	smtpAlertHandlerSchema(ctx, req, resp, false)
}

func (r *defaultSmtpAlertHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	smtpAlertHandlerSchema(ctx, req, resp, true)
}

func smtpAlertHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Smtp Alert Handler.",
		Attributes: map[string]schema.Attribute{
			"asynchronous": schema.BoolAttribute{
				Description: "Indicates whether the server should attempt to invoke this SMTP Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sender_address": schema.StringAttribute{
				Description: "Specifies the email address to use as the sender for messages generated by this alert handler.",
				Required:    true,
			},
			"recipient_address": schema.SetAttribute{
				Description: "Specifies an email address to which the messages should be sent.",
				Required:    true,
				ElementType: types.StringType,
			},
			"message_subject": schema.StringAttribute{
				Description: "Specifies the subject that should be used for email messages generated by this alert handler.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"message_body": schema.StringAttribute{
				Description: "Specifies the body that should be used for email messages generated by this alert handler.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"include_monitor_data_filter": schema.StringAttribute{
				Description: "Include monitor entries that match this filter.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Alert Handler",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Alert Handler is enabled.",
				Required:    true,
			},
			"enabled_alert_severity": schema.SetAttribute{
				Description: "Specifies the alert severities for which this alert handler should be used. If no values are provided, then this alert handler will be enabled for alerts with any severity.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"enabled_alert_type": schema.SetAttribute{
				Description: "Specifies the names of the alert types that are enabled for this alert handler.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"disabled_alert_type": schema.SetAttribute{
				Description: "Specifies the names of the alert types that are disabled for this alert handler.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
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
func addOptionalSmtpAlertHandlerFields(ctx context.Context, addRequest *client.AddSmtpAlertHandlerRequest, plan smtpAlertHandlerResourceModel) error {
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MessageSubject) {
		addRequest.MessageSubject = plan.MessageSubject.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MessageBody) {
		addRequest.MessageBody = plan.MessageBody.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IncludeMonitorDataFilter) {
		addRequest.IncludeMonitorDataFilter = plan.IncludeMonitorDataFilter.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.EnabledAlertSeverity) {
		var slice []string
		plan.EnabledAlertSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.EnabledAlertType) {
		var slice []string
		plan.EnabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertType = enumSlice
	}
	if internaltypes.IsDefined(plan.DisabledAlertType) {
		var slice []string
		plan.DisabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerDisabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerDisabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DisabledAlertType = enumSlice
	}
	return nil
}

// Read a SmtpAlertHandlerResponse object into the model struct
func readSmtpAlertHandlerResponse(ctx context.Context, r *client.SmtpAlertHandlerResponse, state *smtpAlertHandlerResourceModel, expectedValues *smtpAlertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.SenderAddress = types.StringValue(r.SenderAddress)
	state.RecipientAddress = internaltypes.GetStringSet(r.RecipientAddress)
	state.MessageSubject = types.StringValue(r.MessageSubject)
	state.MessageBody = types.StringValue(r.MessageBody)
	state.IncludeMonitorDataFilter = internaltypes.StringTypeOrNil(r.IncludeMonitorDataFilter, internaltypes.IsEmptyString(expectedValues.IncludeMonitorDataFilter))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSmtpAlertHandlerOperations(plan smtpAlertHandlerResourceModel, state smtpAlertHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.Asynchronous, state.Asynchronous, "asynchronous")
	operations.AddStringOperationIfNecessary(&ops, plan.SenderAddress, state.SenderAddress, "sender-address")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RecipientAddress, state.RecipientAddress, "recipient-address")
	operations.AddStringOperationIfNecessary(&ops, plan.MessageSubject, state.MessageSubject, "message-subject")
	operations.AddStringOperationIfNecessary(&ops, plan.MessageBody, state.MessageBody, "message-body")
	operations.AddStringOperationIfNecessary(&ops, plan.IncludeMonitorDataFilter, state.IncludeMonitorDataFilter, "include-monitor-data-filter")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EnabledAlertSeverity, state.EnabledAlertSeverity, "enabled-alert-severity")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EnabledAlertType, state.EnabledAlertType, "enabled-alert-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DisabledAlertType, state.DisabledAlertType, "disabled-alert-type")
	return ops
}

// Create a new resource
func (r *smtpAlertHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan smtpAlertHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var RecipientAddressSlice []string
	plan.RecipientAddress.ElementsAs(ctx, &RecipientAddressSlice, false)
	addRequest := client.NewAddSmtpAlertHandlerRequest(plan.Id.ValueString(),
		[]client.EnumsmtpAlertHandlerSchemaUrn{client.ENUMSMTPALERTHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ALERT_HANDLERSMTP},
		plan.SenderAddress.ValueString(),
		RecipientAddressSlice,
		plan.Enabled.ValueBool())
	err := addOptionalSmtpAlertHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Smtp Alert Handler", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AlertHandlerApi.AddAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAlertHandlerRequest(
		client.AddSmtpAlertHandlerRequestAsAddAlertHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AlertHandlerApi.AddAlertHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Smtp Alert Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state smtpAlertHandlerResourceModel
	readSmtpAlertHandlerResponse(ctx, addResponse.SmtpAlertHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultSmtpAlertHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan smtpAlertHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AlertHandlerApi.GetAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Smtp Alert Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state smtpAlertHandlerResourceModel
	readSmtpAlertHandlerResponse(ctx, readResponse.SmtpAlertHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AlertHandlerApi.UpdateAlertHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSmtpAlertHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AlertHandlerApi.UpdateAlertHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Smtp Alert Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSmtpAlertHandlerResponse(ctx, updateResponse.SmtpAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *smtpAlertHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSmtpAlertHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSmtpAlertHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSmtpAlertHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSmtpAlertHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state smtpAlertHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.AlertHandlerApi.GetAlertHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Smtp Alert Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSmtpAlertHandlerResponse(ctx, readResponse.SmtpAlertHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *smtpAlertHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSmtpAlertHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSmtpAlertHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSmtpAlertHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSmtpAlertHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan smtpAlertHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state smtpAlertHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.AlertHandlerApi.UpdateAlertHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSmtpAlertHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.AlertHandlerApi.UpdateAlertHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Smtp Alert Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSmtpAlertHandlerResponse(ctx, updateResponse.SmtpAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSmtpAlertHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *smtpAlertHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state smtpAlertHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.AlertHandlerApi.DeleteAlertHandlerExecute(r.apiClient.AlertHandlerApi.DeleteAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Smtp Alert Handler", err, httpResp)
		return
	}
}

func (r *smtpAlertHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSmtpAlertHandler(ctx, req, resp)
}

func (r *defaultSmtpAlertHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSmtpAlertHandler(ctx, req, resp)
}

func importSmtpAlertHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
