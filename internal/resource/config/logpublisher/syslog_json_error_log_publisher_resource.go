package logpublisher

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
	_ resource.Resource                = &syslogJsonErrorLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &syslogJsonErrorLogPublisherResource{}
	_ resource.ResourceWithImportState = &syslogJsonErrorLogPublisherResource{}
	_ resource.Resource                = &defaultSyslogJsonErrorLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &defaultSyslogJsonErrorLogPublisherResource{}
	_ resource.ResourceWithImportState = &defaultSyslogJsonErrorLogPublisherResource{}
)

// Create a Syslog Json Error Log Publisher resource
func NewSyslogJsonErrorLogPublisherResource() resource.Resource {
	return &syslogJsonErrorLogPublisherResource{}
}

func NewDefaultSyslogJsonErrorLogPublisherResource() resource.Resource {
	return &defaultSyslogJsonErrorLogPublisherResource{}
}

// syslogJsonErrorLogPublisherResource is the resource implementation.
type syslogJsonErrorLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSyslogJsonErrorLogPublisherResource is the resource implementation.
type defaultSyslogJsonErrorLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *syslogJsonErrorLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_syslog_json_error_log_publisher"
}

func (r *defaultSyslogJsonErrorLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_syslog_json_error_log_publisher"
}

// Configure adds the provider configured client to the resource.
func (r *syslogJsonErrorLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultSyslogJsonErrorLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type syslogJsonErrorLogPublisherResourceModel struct {
	Id                                 types.String `tfsdk:"id"`
	LastUpdated                        types.String `tfsdk:"last_updated"`
	Notifications                      types.Set    `tfsdk:"notifications"`
	RequiredActions                    types.Set    `tfsdk:"required_actions"`
	DefaultSeverity                    types.Set    `tfsdk:"default_severity"`
	SyslogExternalServer               types.Set    `tfsdk:"syslog_external_server"`
	SyslogFacility                     types.String `tfsdk:"syslog_facility"`
	SyslogSeverity                     types.String `tfsdk:"syslog_severity"`
	SyslogMessageHostName              types.String `tfsdk:"syslog_message_host_name"`
	SyslogMessageApplicationName       types.String `tfsdk:"syslog_message_application_name"`
	QueueSize                          types.Int64  `tfsdk:"queue_size"`
	IncludeProductName                 types.Bool   `tfsdk:"include_product_name"`
	IncludeInstanceName                types.Bool   `tfsdk:"include_instance_name"`
	IncludeStartupID                   types.Bool   `tfsdk:"include_startup_id"`
	IncludeThreadID                    types.Bool   `tfsdk:"include_thread_id"`
	GenerifyMessageStringsWhenPossible types.Bool   `tfsdk:"generify_message_strings_when_possible"`
	OverrideSeverity                   types.Set    `tfsdk:"override_severity"`
	Description                        types.String `tfsdk:"description"`
	Enabled                            types.Bool   `tfsdk:"enabled"`
	LoggingErrorBehavior               types.String `tfsdk:"logging_error_behavior"`
}

// GetSchema defines the schema for the resource.
func (r *syslogJsonErrorLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	syslogJsonErrorLogPublisherSchema(ctx, req, resp, false)
}

func (r *defaultSyslogJsonErrorLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	syslogJsonErrorLogPublisherSchema(ctx, req, resp, true)
}

func syslogJsonErrorLogPublisherSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Syslog Json Error Log Publisher.",
		Attributes: map[string]schema.Attribute{
			"default_severity": schema.SetAttribute{
				Description: "Specifies the default severity levels for the logger.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"syslog_external_server": schema.SetAttribute{
				Description: "The syslog server to which messages should be sent.",
				Required:    true,
				ElementType: types.StringType,
			},
			"syslog_facility": schema.StringAttribute{
				Description: "The syslog facility to use for the messages that are logged by this Syslog JSON Error Log Publisher.",
				Optional:    true,
				Computed:    true,
			},
			"syslog_severity": schema.StringAttribute{
				Description: "The syslog severity to use for the messages that are logged by this Syslog JSON Error Log Publisher. If this is not specified, then the severity for each syslog message will be automatically based on the severity for the associated log message.",
				Optional:    true,
				Computed:    true,
			},
			"syslog_message_host_name": schema.StringAttribute{
				Description: "The local host name that will be included in syslog messages that are logged by this Syslog JSON Error Log Publisher.",
				Optional:    true,
				Computed:    true,
			},
			"syslog_message_application_name": schema.StringAttribute{
				Description: "The application name that will be included in syslog messages that are logged by this Syslog JSON Error Log Publisher.",
				Optional:    true,
				Computed:    true,
			},
			"queue_size": schema.Int64Attribute{
				Description: "The maximum number of log records that can be stored in the asynchronous queue.",
				Optional:    true,
				Computed:    true,
			},
			"include_product_name": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the product name for the Directory Server.",
				Optional:    true,
				Computed:    true,
			},
			"include_instance_name": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the instance name for the Directory Server.",
				Optional:    true,
				Computed:    true,
			},
			"include_startup_id": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the startup ID for the Directory Server, which is a value assigned to the server instance at startup and may be used to identify when the server has been restarted.",
				Optional:    true,
				Computed:    true,
			},
			"include_thread_id": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the thread ID for the Directory Server in each log message. This ID can be used to correlate log messages from the same thread within a single log as well as generated by the same thread across different types of log files. More information about the thread with a specific ID can be obtained using the cn=JVM Stack Trace,cn=monitor entry.",
				Optional:    true,
				Computed:    true,
			},
			"generify_message_strings_when_possible": schema.BoolAttribute{
				Description: "Indicates whether to use the generified version of the log message string (which may use placeholders like %s for a string or %d for an integer), rather than the version of the message with those placeholders replaced with specific values that would normally be written to the log.",
				Optional:    true,
				Computed:    true,
			},
			"override_severity": schema.SetAttribute{
				Description: "Specifies the override severity levels for the logger based on the category of the messages.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Publisher",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Log Publisher is enabled for use.",
				Required:    true,
			},
			"logging_error_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the server should exhibit if an error occurs during logging processing.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	if setOptionalToComputed {
		config.SetOptionalAttributesToComputed(&schema)
	}
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalSyslogJsonErrorLogPublisherFields(ctx context.Context, addRequest *client.AddSyslogJsonErrorLogPublisherRequest, plan syslogJsonErrorLogPublisherResourceModel) error {
	if internaltypes.IsDefined(plan.DefaultSeverity) {
		var slice []string
		plan.DefaultSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherDefaultSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherDefaultSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefaultSeverity = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogFacility) {
		syslogFacility, err := client.NewEnumlogPublisherSyslogFacilityPropFromValue(plan.SyslogFacility.ValueString())
		if err != nil {
			return err
		}
		addRequest.SyslogFacility = syslogFacility
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogSeverity) {
		syslogSeverity, err := client.NewEnumlogPublisherSyslogSeverityPropFromValue(plan.SyslogSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.SyslogSeverity = syslogSeverity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogMessageHostName) {
		stringVal := plan.SyslogMessageHostName.ValueString()
		addRequest.SyslogMessageHostName = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogMessageApplicationName) {
		stringVal := plan.SyslogMessageApplicationName.ValueString()
		addRequest.SyslogMessageApplicationName = &stringVal
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		intVal := int32(plan.QueueSize.ValueInt64())
		addRequest.QueueSize = &intVal
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		boolVal := plan.IncludeProductName.ValueBool()
		addRequest.IncludeProductName = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		boolVal := plan.IncludeInstanceName.ValueBool()
		addRequest.IncludeInstanceName = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		boolVal := plan.IncludeStartupID.ValueBool()
		addRequest.IncludeStartupID = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		boolVal := plan.IncludeThreadID.ValueBool()
		addRequest.IncludeThreadID = &boolVal
	}
	if internaltypes.IsDefined(plan.GenerifyMessageStringsWhenPossible) {
		boolVal := plan.GenerifyMessageStringsWhenPossible.ValueBool()
		addRequest.GenerifyMessageStringsWhenPossible = &boolVal
	}
	if internaltypes.IsDefined(plan.OverrideSeverity) {
		var slice []string
		plan.OverrideSeverity.ElementsAs(ctx, &slice, false)
		addRequest.OverrideSeverity = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Read a SyslogJsonErrorLogPublisherResponse object into the model struct
func readSyslogJsonErrorLogPublisherResponse(ctx context.Context, r *client.SyslogJsonErrorLogPublisherResponse, state *syslogJsonErrorLogPublisherResourceModel, expectedValues *syslogJsonErrorLogPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.SyslogExternalServer = internaltypes.GetStringSet(r.SyslogExternalServer)
	state.SyslogFacility = types.StringValue(r.SyslogFacility.String())
	state.SyslogSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherSyslogSeverityProp(r.SyslogSeverity), internaltypes.IsEmptyString(expectedValues.SyslogSeverity))
	state.SyslogMessageHostName = internaltypes.StringTypeOrNil(r.SyslogMessageHostName, internaltypes.IsEmptyString(expectedValues.SyslogMessageHostName))
	state.SyslogMessageApplicationName = internaltypes.StringTypeOrNil(r.SyslogMessageApplicationName, internaltypes.IsEmptyString(expectedValues.SyslogMessageApplicationName))
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSyslogJsonErrorLogPublisherOperations(plan syslogJsonErrorLogPublisherResourceModel, state syslogJsonErrorLogPublisherResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefaultSeverity, state.DefaultSeverity, "default-severity")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SyslogExternalServer, state.SyslogExternalServer, "syslog-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.SyslogFacility, state.SyslogFacility, "syslog-facility")
	operations.AddStringOperationIfNecessary(&ops, plan.SyslogSeverity, state.SyslogSeverity, "syslog-severity")
	operations.AddStringOperationIfNecessary(&ops, plan.SyslogMessageHostName, state.SyslogMessageHostName, "syslog-message-host-name")
	operations.AddStringOperationIfNecessary(&ops, plan.SyslogMessageApplicationName, state.SyslogMessageApplicationName, "syslog-message-application-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.QueueSize, state.QueueSize, "queue-size")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeProductName, state.IncludeProductName, "include-product-name")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeInstanceName, state.IncludeInstanceName, "include-instance-name")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeStartupID, state.IncludeStartupID, "include-startup-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeThreadID, state.IncludeThreadID, "include-thread-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.GenerifyMessageStringsWhenPossible, state.GenerifyMessageStringsWhenPossible, "generify-message-strings-when-possible")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.OverrideSeverity, state.OverrideSeverity, "override-severity")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.LoggingErrorBehavior, state.LoggingErrorBehavior, "logging-error-behavior")
	return ops
}

// Create a new resource
func (r *syslogJsonErrorLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan syslogJsonErrorLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var SyslogExternalServerSlice []string
	plan.SyslogExternalServer.ElementsAs(ctx, &SyslogExternalServerSlice, false)
	addRequest := client.NewAddSyslogJsonErrorLogPublisherRequest(plan.Id.ValueString(),
		[]client.EnumsyslogJsonErrorLogPublisherSchemaUrn{client.ENUMSYSLOGJSONERRORLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERSYSLOG_JSON_ERROR},
		SyslogExternalServerSlice,
		plan.Enabled.ValueBool())
	err := addOptionalSyslogJsonErrorLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Syslog Json Error Log Publisher", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddSyslogJsonErrorLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Syslog Json Error Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state syslogJsonErrorLogPublisherResourceModel
	readSyslogJsonErrorLogPublisherResponse(ctx, addResponse.SyslogJsonErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultSyslogJsonErrorLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan syslogJsonErrorLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Syslog Json Error Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state syslogJsonErrorLogPublisherResourceModel
	readSyslogJsonErrorLogPublisherResponse(ctx, readResponse.SyslogJsonErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogPublisherApi.UpdateLogPublisher(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSyslogJsonErrorLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Syslog Json Error Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSyslogJsonErrorLogPublisherResponse(ctx, updateResponse.SyslogJsonErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *syslogJsonErrorLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSyslogJsonErrorLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSyslogJsonErrorLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSyslogJsonErrorLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSyslogJsonErrorLogPublisher(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state syslogJsonErrorLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Syslog Json Error Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSyslogJsonErrorLogPublisherResponse(ctx, readResponse.SyslogJsonErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *syslogJsonErrorLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSyslogJsonErrorLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSyslogJsonErrorLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSyslogJsonErrorLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSyslogJsonErrorLogPublisher(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan syslogJsonErrorLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state syslogJsonErrorLogPublisherResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogPublisherApi.UpdateLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSyslogJsonErrorLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Syslog Json Error Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSyslogJsonErrorLogPublisherResponse(ctx, updateResponse.SyslogJsonErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSyslogJsonErrorLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *syslogJsonErrorLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state syslogJsonErrorLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogPublisherApi.DeleteLogPublisherExecute(r.apiClient.LogPublisherApi.DeleteLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Syslog Json Error Log Publisher", err, httpResp)
		return
	}
}

func (r *syslogJsonErrorLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSyslogJsonErrorLogPublisher(ctx, req, resp)
}

func (r *defaultSyslogJsonErrorLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSyslogJsonErrorLogPublisher(ctx, req, resp)
}

func importSyslogJsonErrorLogPublisher(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
