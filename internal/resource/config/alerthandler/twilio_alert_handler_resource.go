package alerthandler

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
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
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &twilioAlertHandlerResource{}
	_ resource.ResourceWithConfigure   = &twilioAlertHandlerResource{}
	_ resource.ResourceWithImportState = &twilioAlertHandlerResource{}
	_ resource.Resource                = &defaultTwilioAlertHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultTwilioAlertHandlerResource{}
	_ resource.ResourceWithImportState = &defaultTwilioAlertHandlerResource{}
)

// Create a Twilio Alert Handler resource
func NewTwilioAlertHandlerResource() resource.Resource {
	return &twilioAlertHandlerResource{}
}

func NewDefaultTwilioAlertHandlerResource() resource.Resource {
	return &defaultTwilioAlertHandlerResource{}
}

// twilioAlertHandlerResource is the resource implementation.
type twilioAlertHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultTwilioAlertHandlerResource is the resource implementation.
type defaultTwilioAlertHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *twilioAlertHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_twilio_alert_handler"
}

func (r *defaultTwilioAlertHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_twilio_alert_handler"
}

// Configure adds the provider configured client to the resource.
func (r *twilioAlertHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultTwilioAlertHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type twilioAlertHandlerResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	LastUpdated                       types.String `tfsdk:"last_updated"`
	Notifications                     types.Set    `tfsdk:"notifications"`
	RequiredActions                   types.Set    `tfsdk:"required_actions"`
	Asynchronous                      types.Bool   `tfsdk:"asynchronous"`
	HttpProxyExternalServer           types.String `tfsdk:"http_proxy_external_server"`
	TwilioAccountSID                  types.String `tfsdk:"twilio_account_sid"`
	TwilioAuthToken                   types.String `tfsdk:"twilio_auth_token"`
	TwilioAuthTokenPassphraseProvider types.String `tfsdk:"twilio_auth_token_passphrase_provider"`
	SenderPhoneNumber                 types.Set    `tfsdk:"sender_phone_number"`
	RecipientPhoneNumber              types.Set    `tfsdk:"recipient_phone_number"`
	LongMessageBehavior               types.String `tfsdk:"long_message_behavior"`
	Description                       types.String `tfsdk:"description"`
	Enabled                           types.Bool   `tfsdk:"enabled"`
	EnabledAlertSeverity              types.Set    `tfsdk:"enabled_alert_severity"`
	EnabledAlertType                  types.Set    `tfsdk:"enabled_alert_type"`
	DisabledAlertType                 types.Set    `tfsdk:"disabled_alert_type"`
}

// GetSchema defines the schema for the resource.
func (r *twilioAlertHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	twilioAlertHandlerSchema(ctx, req, resp, false)
}

func (r *defaultTwilioAlertHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	twilioAlertHandlerSchema(ctx, req, resp, true)
}

func twilioAlertHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Twilio Alert Handler.",
		Attributes: map[string]schema.Attribute{
			"asynchronous": schema.BoolAttribute{
				Description: "Indicates whether the server should attempt to invoke this Twilio Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"http_proxy_external_server": schema.StringAttribute{
				Description: "A reference to an HTTP proxy server that should be used for requests sent to the Twilio service. Supported in PingDirectory product version 9.2.0.0+.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"twilio_account_sid": schema.StringAttribute{
				Description: "The unique identifier assigned to the Twilio account that will be used.",
				Required:    true,
			},
			"twilio_auth_token": schema.StringAttribute{
				Description: "The auth token for the Twilio account that will be used.",
				Optional:    true,
				Sensitive:   true,
			},
			"twilio_auth_token_passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider that may be used to obtain the auth token for the Twilio account that will be used.",
				Optional:    true,
			},
			"sender_phone_number": schema.SetAttribute{
				Description: "The outgoing phone number to use for the messages. Values must be phone numbers you have obtained for use with your Twilio account.",
				Required:    true,
				ElementType: types.StringType,
			},
			"recipient_phone_number": schema.SetAttribute{
				Description: "The phone number to which alert notifications should be delivered.",
				Required:    true,
				ElementType: types.StringType,
			},
			"long_message_behavior": schema.StringAttribute{
				Description: "The behavior to use for alert messages that are longer than the 160-character size limit for SMS messages.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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

// Validate that any version restrictions are met in the plan
func (r *twilioAlertHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanTwilioAlertHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultTwilioAlertHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanTwilioAlertHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanTwilioAlertHandler(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	compare, err := version.Compare(providerConfig.ProductVersion, version.PingDirectory9200)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model twilioAlertHandlerResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsNonEmptyString(model.HttpProxyExternalServer) {
		resp.Diagnostics.AddError("Attribute 'http_proxy_external_server' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
}

// Add config validators
func (r twilioAlertHandlerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("twilio_auth_token_passphrase_provider"),
			path.MatchRoot("twilio_auth_token"),
		),
	}
}

// Add optional fields to create request
func addOptionalTwilioAlertHandlerFields(ctx context.Context, addRequest *client.AddTwilioAlertHandlerRequest, plan twilioAlertHandlerResourceModel) error {
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HttpProxyExternalServer) {
		addRequest.HttpProxyExternalServer = plan.HttpProxyExternalServer.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TwilioAuthToken) {
		addRequest.TwilioAuthToken = plan.TwilioAuthToken.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TwilioAuthTokenPassphraseProvider) {
		addRequest.TwilioAuthTokenPassphraseProvider = plan.TwilioAuthTokenPassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LongMessageBehavior) {
		longMessageBehavior, err := client.NewEnumalertHandlerLongMessageBehaviorPropFromValue(plan.LongMessageBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LongMessageBehavior = longMessageBehavior
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

// Read a TwilioAlertHandlerResponse object into the model struct
func readTwilioAlertHandlerResponse(ctx context.Context, r *client.TwilioAlertHandlerResponse, state *twilioAlertHandlerResourceModel, expectedValues *twilioAlertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, internaltypes.IsEmptyString(expectedValues.HttpProxyExternalServer))
	state.TwilioAccountSID = types.StringValue(r.TwilioAccountSID)
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.TwilioAuthToken = expectedValues.TwilioAuthToken
	state.TwilioAuthTokenPassphraseProvider = internaltypes.StringTypeOrNil(r.TwilioAuthTokenPassphraseProvider, internaltypes.IsEmptyString(expectedValues.TwilioAuthTokenPassphraseProvider))
	state.SenderPhoneNumber = internaltypes.GetStringSet(r.SenderPhoneNumber)
	state.RecipientPhoneNumber = internaltypes.GetStringSet(r.RecipientPhoneNumber)
	state.LongMessageBehavior = types.StringValue(r.LongMessageBehavior.String())
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
func createTwilioAlertHandlerOperations(plan twilioAlertHandlerResourceModel, state twilioAlertHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.Asynchronous, state.Asynchronous, "asynchronous")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpProxyExternalServer, state.HttpProxyExternalServer, "http-proxy-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.TwilioAccountSID, state.TwilioAccountSID, "twilio-account-sid")
	operations.AddStringOperationIfNecessary(&ops, plan.TwilioAuthToken, state.TwilioAuthToken, "twilio-auth-token")
	operations.AddStringOperationIfNecessary(&ops, plan.TwilioAuthTokenPassphraseProvider, state.TwilioAuthTokenPassphraseProvider, "twilio-auth-token-passphrase-provider")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SenderPhoneNumber, state.SenderPhoneNumber, "sender-phone-number")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RecipientPhoneNumber, state.RecipientPhoneNumber, "recipient-phone-number")
	operations.AddStringOperationIfNecessary(&ops, plan.LongMessageBehavior, state.LongMessageBehavior, "long-message-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EnabledAlertSeverity, state.EnabledAlertSeverity, "enabled-alert-severity")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EnabledAlertType, state.EnabledAlertType, "enabled-alert-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DisabledAlertType, state.DisabledAlertType, "disabled-alert-type")
	return ops
}

// Create a new resource
func (r *twilioAlertHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan twilioAlertHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var SenderPhoneNumberSlice []string
	plan.SenderPhoneNumber.ElementsAs(ctx, &SenderPhoneNumberSlice, false)
	var RecipientPhoneNumberSlice []string
	plan.RecipientPhoneNumber.ElementsAs(ctx, &RecipientPhoneNumberSlice, false)
	addRequest := client.NewAddTwilioAlertHandlerRequest(plan.Id.ValueString(),
		[]client.EnumtwilioAlertHandlerSchemaUrn{client.ENUMTWILIOALERTHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ALERT_HANDLERTWILIO},
		plan.TwilioAccountSID.ValueString(),
		SenderPhoneNumberSlice,
		RecipientPhoneNumberSlice,
		plan.Enabled.ValueBool())
	err := addOptionalTwilioAlertHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Twilio Alert Handler", err.Error())
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
		client.AddTwilioAlertHandlerRequestAsAddAlertHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AlertHandlerApi.AddAlertHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Twilio Alert Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state twilioAlertHandlerResourceModel
	readTwilioAlertHandlerResponse(ctx, addResponse.TwilioAlertHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultTwilioAlertHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan twilioAlertHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AlertHandlerApi.GetAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Twilio Alert Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state twilioAlertHandlerResourceModel
	readTwilioAlertHandlerResponse(ctx, readResponse.TwilioAlertHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AlertHandlerApi.UpdateAlertHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createTwilioAlertHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AlertHandlerApi.UpdateAlertHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Twilio Alert Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readTwilioAlertHandlerResponse(ctx, updateResponse.TwilioAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *twilioAlertHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readTwilioAlertHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultTwilioAlertHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readTwilioAlertHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readTwilioAlertHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state twilioAlertHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.AlertHandlerApi.GetAlertHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Twilio Alert Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readTwilioAlertHandlerResponse(ctx, readResponse.TwilioAlertHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *twilioAlertHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateTwilioAlertHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultTwilioAlertHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateTwilioAlertHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateTwilioAlertHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan twilioAlertHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state twilioAlertHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.AlertHandlerApi.UpdateAlertHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createTwilioAlertHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.AlertHandlerApi.UpdateAlertHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Twilio Alert Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readTwilioAlertHandlerResponse(ctx, updateResponse.TwilioAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultTwilioAlertHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *twilioAlertHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state twilioAlertHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.AlertHandlerApi.DeleteAlertHandlerExecute(r.apiClient.AlertHandlerApi.DeleteAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Twilio Alert Handler", err, httpResp)
		return
	}
}

func (r *twilioAlertHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importTwilioAlertHandler(ctx, req, resp)
}

func (r *defaultTwilioAlertHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importTwilioAlertHandler(ctx, req, resp)
}

func importTwilioAlertHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
