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
	client "github.com/pingidentity/pingdirectory-go-client/v9100"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &thirdPartyFileBasedAccessLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &thirdPartyFileBasedAccessLogPublisherResource{}
	_ resource.ResourceWithImportState = &thirdPartyFileBasedAccessLogPublisherResource{}
)

// Create a Third Party File Based Access Log Publisher resource
func NewThirdPartyFileBasedAccessLogPublisherResource() resource.Resource {
	return &thirdPartyFileBasedAccessLogPublisherResource{}
}

// thirdPartyFileBasedAccessLogPublisherResource is the resource implementation.
type thirdPartyFileBasedAccessLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *thirdPartyFileBasedAccessLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_third_party_file_based_access_log_publisher"
}

// Configure adds the provider configured client to the resource.
func (r *thirdPartyFileBasedAccessLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type thirdPartyFileBasedAccessLogPublisherResourceModel struct {
	Id                             types.String `tfsdk:"id"`
	LastUpdated                    types.String `tfsdk:"last_updated"`
	Notifications                  types.Set    `tfsdk:"notifications"`
	RequiredActions                types.Set    `tfsdk:"required_actions"`
	LogFile                        types.String `tfsdk:"log_file"`
	LogFilePermissions             types.String `tfsdk:"log_file_permissions"`
	RotationPolicy                 types.Set    `tfsdk:"rotation_policy"`
	RotationListener               types.Set    `tfsdk:"rotation_listener"`
	RetentionPolicy                types.Set    `tfsdk:"retention_policy"`
	CompressionMechanism           types.String `tfsdk:"compression_mechanism"`
	SignLog                        types.Bool   `tfsdk:"sign_log"`
	EncryptLog                     types.Bool   `tfsdk:"encrypt_log"`
	EncryptionSettingsDefinitionID types.String `tfsdk:"encryption_settings_definition_id"`
	Append                         types.Bool   `tfsdk:"append"`
	ExtensionClass                 types.String `tfsdk:"extension_class"`
	ExtensionArgument              types.Set    `tfsdk:"extension_argument"`
	Asynchronous                   types.Bool   `tfsdk:"asynchronous"`
	AutoFlush                      types.Bool   `tfsdk:"auto_flush"`
	BufferSize                     types.String `tfsdk:"buffer_size"`
	QueueSize                      types.Int64  `tfsdk:"queue_size"`
	TimeInterval                   types.String `tfsdk:"time_interval"`
	LogConnects                    types.Bool   `tfsdk:"log_connects"`
	LogDisconnects                 types.Bool   `tfsdk:"log_disconnects"`
	LogSecurityNegotiation         types.Bool   `tfsdk:"log_security_negotiation"`
	LogClientCertificates          types.Bool   `tfsdk:"log_client_certificates"`
	LogRequests                    types.Bool   `tfsdk:"log_requests"`
	LogResults                     types.Bool   `tfsdk:"log_results"`
	LogSearchEntries               types.Bool   `tfsdk:"log_search_entries"`
	LogSearchReferences            types.Bool   `tfsdk:"log_search_references"`
	LogIntermediateResponses       types.Bool   `tfsdk:"log_intermediate_responses"`
	SuppressInternalOperations     types.Bool   `tfsdk:"suppress_internal_operations"`
	SuppressReplicationOperations  types.Bool   `tfsdk:"suppress_replication_operations"`
	CorrelateRequestsAndResults    types.Bool   `tfsdk:"correlate_requests_and_results"`
	ConnectionCriteria             types.String `tfsdk:"connection_criteria"`
	RequestCriteria                types.String `tfsdk:"request_criteria"`
	ResultCriteria                 types.String `tfsdk:"result_criteria"`
	SearchEntryCriteria            types.String `tfsdk:"search_entry_criteria"`
	SearchReferenceCriteria        types.String `tfsdk:"search_reference_criteria"`
	Description                    types.String `tfsdk:"description"`
	Enabled                        types.Bool   `tfsdk:"enabled"`
	LoggingErrorBehavior           types.String `tfsdk:"logging_error_behavior"`
}

// GetSchema defines the schema for the resource.
func (r *thirdPartyFileBasedAccessLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Third Party File Based Access Log Publisher.",
		Attributes: map[string]schema.Attribute{
			"log_file": schema.StringAttribute{
				Description: "The file name to use for the log files generated by the Third Party File Based Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.",
				Required:    true,
			},
			"log_file_permissions": schema.StringAttribute{
				Description: "The UNIX permissions of the log files created by this Third Party File Based Access Log Publisher.",
				Required:    true,
			},
			"rotation_policy": schema.SetAttribute{
				Description: "The rotation policy to use for the Third Party File Based Access Log Publisher .",
				Required:    true,
				ElementType: types.StringType,
			},
			"rotation_listener": schema.SetAttribute{
				Description: "A listener that should be notified whenever a log file is rotated out of service.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"retention_policy": schema.SetAttribute{
				Description: "The retention policy to use for the Third Party File Based Access Log Publisher .",
				Required:    true,
				ElementType: types.StringType,
			},
			"compression_mechanism": schema.StringAttribute{
				Description: "Specifies the type of compression (if any) to use for log files that are written.",
				Optional:    true,
				Computed:    true,
			},
			"sign_log": schema.BoolAttribute{
				Description: "Indicates whether the log should be cryptographically signed so that the log content cannot be altered in an undetectable manner.",
				Optional:    true,
				Computed:    true,
			},
			"encrypt_log": schema.BoolAttribute{
				Description: "Indicates whether log files should be encrypted so that their content is not available to unauthorized users.",
				Optional:    true,
				Computed:    true,
			},
			"encryption_settings_definition_id": schema.StringAttribute{
				Description: "Specifies the ID of the encryption settings definition that should be used to encrypt the data. If this is not provided, the server's preferred encryption settings definition will be used. The \"encryption-settings list\" command can be used to obtain a list of the encryption settings definitions available in the server.",
				Optional:    true,
			},
			"append": schema.BoolAttribute{
				Description: "Specifies whether to append to existing log files.",
				Optional:    true,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party File Based Access Log Publisher.",
				Required:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party File Based Access Log Publisher. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"asynchronous": schema.BoolAttribute{
				Description: "Indicates whether the Third Party File Based Access Log Publisher will publish records asynchronously.",
				Required:    true,
			},
			"auto_flush": schema.BoolAttribute{
				Description: "Specifies whether to flush the writer after every log record.",
				Optional:    true,
				Computed:    true,
			},
			"buffer_size": schema.StringAttribute{
				Description: "Specifies the log file buffer size.",
				Optional:    true,
				Computed:    true,
			},
			"queue_size": schema.Int64Attribute{
				Description: "The maximum number of log records that can be stored in the asynchronous queue.",
				Optional:    true,
				Computed:    true,
			},
			"time_interval": schema.StringAttribute{
				Description: "Specifies the interval at which to check whether the log files need to be rotated.",
				Optional:    true,
				Computed:    true,
			},
			"log_connects": schema.BoolAttribute{
				Description: "Indicates whether to log information about connections established to the server.",
				Optional:    true,
				Computed:    true,
			},
			"log_disconnects": schema.BoolAttribute{
				Description: "Indicates whether to log information about connections that have been closed by the client or terminated by the server.",
				Optional:    true,
				Computed:    true,
			},
			"log_security_negotiation": schema.BoolAttribute{
				Description: "Indicates whether to log information about the result of any security negotiation (e.g., SSL handshake) processing that has been performed.",
				Optional:    true,
				Computed:    true,
			},
			"log_client_certificates": schema.BoolAttribute{
				Description: "Indicates whether to log information about any client certificates presented to the server.",
				Optional:    true,
				Computed:    true,
			},
			"log_requests": schema.BoolAttribute{
				Description: "Indicates whether to log information about requests received from clients.",
				Optional:    true,
				Computed:    true,
			},
			"log_results": schema.BoolAttribute{
				Description: "Indicates whether to log information about the results of client requests.",
				Optional:    true,
				Computed:    true,
			},
			"log_search_entries": schema.BoolAttribute{
				Description: "Indicates whether to log information about search result entries sent to the client.",
				Optional:    true,
				Computed:    true,
			},
			"log_search_references": schema.BoolAttribute{
				Description: "Indicates whether to log information about search result references sent to the client.",
				Optional:    true,
				Computed:    true,
			},
			"log_intermediate_responses": schema.BoolAttribute{
				Description: "Indicates whether to log information about intermediate responses sent to the client.",
				Optional:    true,
				Computed:    true,
			},
			"suppress_internal_operations": schema.BoolAttribute{
				Description: "Indicates whether internal operations (for example, operations that are initiated by plugins) should be logged along with the operations that are requested by users.",
				Optional:    true,
				Computed:    true,
			},
			"suppress_replication_operations": schema.BoolAttribute{
				Description: "Indicates whether access messages that are generated by replication operations should be suppressed.",
				Optional:    true,
				Computed:    true,
			},
			"correlate_requests_and_results": schema.BoolAttribute{
				Description: "Indicates whether to automatically log result messages for any operation in which the corresponding request was logged. In such cases, the result, entry, and reference criteria will be ignored, although the log-responses, log-search-entries, and log-search-references properties will be honored.",
				Optional:    true,
				Computed:    true,
			},
			"connection_criteria": schema.StringAttribute{
				Description: "Specifies a set of connection criteria that must match the associated client connection in order for a connect, disconnect, request, or result message to be logged.",
				Optional:    true,
			},
			"request_criteria": schema.StringAttribute{
				Description: "Specifies a set of request criteria that must match the associated operation request in order for a request or result to be logged by this Access Log Publisher.",
				Optional:    true,
			},
			"result_criteria": schema.StringAttribute{
				Description: "Specifies a set of result criteria that must match the associated operation result in order for that result to be logged by this Access Log Publisher.",
				Optional:    true,
			},
			"search_entry_criteria": schema.StringAttribute{
				Description: "Specifies a set of search entry criteria that must match the associated search result entry in order for that it to be logged by this Access Log Publisher.",
				Optional:    true,
			},
			"search_reference_criteria": schema.StringAttribute{
				Description: "Specifies a set of search reference criteria that must match the associated search result reference in order for that it to be logged by this Access Log Publisher.",
				Optional:    true,
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
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalThirdPartyFileBasedAccessLogPublisherFields(ctx context.Context, addRequest *client.AddThirdPartyFileBasedAccessLogPublisherRequest, plan thirdPartyFileBasedAccessLogPublisherResourceModel) error {
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CompressionMechanism) {
		compressionMechanism, err := client.NewEnumlogPublisherCompressionMechanismPropFromValue(plan.CompressionMechanism.ValueString())
		if err != nil {
			return err
		}
		addRequest.CompressionMechanism = compressionMechanism
	}
	if internaltypes.IsDefined(plan.SignLog) {
		boolVal := plan.SignLog.ValueBool()
		addRequest.SignLog = &boolVal
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		boolVal := plan.EncryptLog.ValueBool()
		addRequest.EncryptLog = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		stringVal := plan.EncryptionSettingsDefinitionID.ValueString()
		addRequest.EncryptionSettingsDefinitionID = &stringVal
	}
	if internaltypes.IsDefined(plan.Append) {
		boolVal := plan.Append.ValueBool()
		addRequest.Append = &boolVal
	}
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		boolVal := plan.AutoFlush.ValueBool()
		addRequest.AutoFlush = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		stringVal := plan.BufferSize.ValueString()
		addRequest.BufferSize = &stringVal
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		intVal := int32(plan.QueueSize.ValueInt64())
		addRequest.QueueSize = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		stringVal := plan.TimeInterval.ValueString()
		addRequest.TimeInterval = &stringVal
	}
	if internaltypes.IsDefined(plan.LogConnects) {
		boolVal := plan.LogConnects.ValueBool()
		addRequest.LogConnects = &boolVal
	}
	if internaltypes.IsDefined(plan.LogDisconnects) {
		boolVal := plan.LogDisconnects.ValueBool()
		addRequest.LogDisconnects = &boolVal
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		boolVal := plan.LogSecurityNegotiation.ValueBool()
		addRequest.LogSecurityNegotiation = &boolVal
	}
	if internaltypes.IsDefined(plan.LogClientCertificates) {
		boolVal := plan.LogClientCertificates.ValueBool()
		addRequest.LogClientCertificates = &boolVal
	}
	if internaltypes.IsDefined(plan.LogRequests) {
		boolVal := plan.LogRequests.ValueBool()
		addRequest.LogRequests = &boolVal
	}
	if internaltypes.IsDefined(plan.LogResults) {
		boolVal := plan.LogResults.ValueBool()
		addRequest.LogResults = &boolVal
	}
	if internaltypes.IsDefined(plan.LogSearchEntries) {
		boolVal := plan.LogSearchEntries.ValueBool()
		addRequest.LogSearchEntries = &boolVal
	}
	if internaltypes.IsDefined(plan.LogSearchReferences) {
		boolVal := plan.LogSearchReferences.ValueBool()
		addRequest.LogSearchReferences = &boolVal
	}
	if internaltypes.IsDefined(plan.LogIntermediateResponses) {
		boolVal := plan.LogIntermediateResponses.ValueBool()
		addRequest.LogIntermediateResponses = &boolVal
	}
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		boolVal := plan.SuppressInternalOperations.ValueBool()
		addRequest.SuppressInternalOperations = &boolVal
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		boolVal := plan.SuppressReplicationOperations.ValueBool()
		addRequest.SuppressReplicationOperations = &boolVal
	}
	if internaltypes.IsDefined(plan.CorrelateRequestsAndResults) {
		boolVal := plan.CorrelateRequestsAndResults.ValueBool()
		addRequest.CorrelateRequestsAndResults = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		stringVal := plan.ConnectionCriteria.ValueString()
		addRequest.ConnectionCriteria = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		stringVal := plan.RequestCriteria.ValueString()
		addRequest.RequestCriteria = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		stringVal := plan.ResultCriteria.ValueString()
		addRequest.ResultCriteria = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchEntryCriteria) {
		stringVal := plan.SearchEntryCriteria.ValueString()
		addRequest.SearchEntryCriteria = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchReferenceCriteria) {
		stringVal := plan.SearchReferenceCriteria.ValueString()
		addRequest.SearchReferenceCriteria = &stringVal
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

// Read a ThirdPartyFileBasedAccessLogPublisherResponse object into the model struct
func readThirdPartyFileBasedAccessLogPublisherResponse(ctx context.Context, r *client.ThirdPartyFileBasedAccessLogPublisherResponse, state *thirdPartyFileBasedAccessLogPublisherResourceModel, expectedValues *thirdPartyFileBasedAccessLogPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), internaltypes.IsEmptyString(expectedValues.CompressionMechanism))
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, internaltypes.IsEmptyString(expectedValues.BufferSize))
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, internaltypes.IsEmptyString(expectedValues.TimeInterval))
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.LogConnects = internaltypes.BoolTypeOrNil(r.LogConnects)
	state.LogDisconnects = internaltypes.BoolTypeOrNil(r.LogDisconnects)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogClientCertificates = internaltypes.BoolTypeOrNil(r.LogClientCertificates)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.LogSearchEntries = internaltypes.BoolTypeOrNil(r.LogSearchEntries)
	state.LogSearchReferences = internaltypes.BoolTypeOrNil(r.LogSearchReferences)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.CorrelateRequestsAndResults = internaltypes.BoolTypeOrNil(r.CorrelateRequestsAndResults)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, internaltypes.IsEmptyString(expectedValues.SearchEntryCriteria))
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, internaltypes.IsEmptyString(expectedValues.SearchReferenceCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createThirdPartyFileBasedAccessLogPublisherOperations(plan thirdPartyFileBasedAccessLogPublisherResourceModel, state thirdPartyFileBasedAccessLogPublisherResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.LogFile, state.LogFile, "log-file")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFilePermissions, state.LogFilePermissions, "log-file-permissions")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RotationPolicy, state.RotationPolicy, "rotation-policy")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RotationListener, state.RotationListener, "rotation-listener")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RetentionPolicy, state.RetentionPolicy, "retention-policy")
	operations.AddStringOperationIfNecessary(&ops, plan.CompressionMechanism, state.CompressionMechanism, "compression-mechanism")
	operations.AddBoolOperationIfNecessary(&ops, plan.SignLog, state.SignLog, "sign-log")
	operations.AddBoolOperationIfNecessary(&ops, plan.EncryptLog, state.EncryptLog, "encrypt-log")
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptionSettingsDefinitionID, state.EncryptionSettingsDefinitionID, "encryption-settings-definition-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.Append, state.Append, "append")
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddBoolOperationIfNecessary(&ops, plan.Asynchronous, state.Asynchronous, "asynchronous")
	operations.AddBoolOperationIfNecessary(&ops, plan.AutoFlush, state.AutoFlush, "auto-flush")
	operations.AddStringOperationIfNecessary(&ops, plan.BufferSize, state.BufferSize, "buffer-size")
	operations.AddInt64OperationIfNecessary(&ops, plan.QueueSize, state.QueueSize, "queue-size")
	operations.AddStringOperationIfNecessary(&ops, plan.TimeInterval, state.TimeInterval, "time-interval")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogConnects, state.LogConnects, "log-connects")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogDisconnects, state.LogDisconnects, "log-disconnects")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogSecurityNegotiation, state.LogSecurityNegotiation, "log-security-negotiation")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogClientCertificates, state.LogClientCertificates, "log-client-certificates")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogRequests, state.LogRequests, "log-requests")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogResults, state.LogResults, "log-results")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogSearchEntries, state.LogSearchEntries, "log-search-entries")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogSearchReferences, state.LogSearchReferences, "log-search-references")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogIntermediateResponses, state.LogIntermediateResponses, "log-intermediate-responses")
	operations.AddBoolOperationIfNecessary(&ops, plan.SuppressInternalOperations, state.SuppressInternalOperations, "suppress-internal-operations")
	operations.AddBoolOperationIfNecessary(&ops, plan.SuppressReplicationOperations, state.SuppressReplicationOperations, "suppress-replication-operations")
	operations.AddBoolOperationIfNecessary(&ops, plan.CorrelateRequestsAndResults, state.CorrelateRequestsAndResults, "correlate-requests-and-results")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionCriteria, state.ConnectionCriteria, "connection-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.RequestCriteria, state.RequestCriteria, "request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.ResultCriteria, state.ResultCriteria, "result-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchEntryCriteria, state.SearchEntryCriteria, "search-entry-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchReferenceCriteria, state.SearchReferenceCriteria, "search-reference-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.LoggingErrorBehavior, state.LoggingErrorBehavior, "logging-error-behavior")
	return ops
}

// Create a new resource
func (r *thirdPartyFileBasedAccessLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan thirdPartyFileBasedAccessLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var RotationPolicySlice []string
	plan.RotationPolicy.ElementsAs(ctx, &RotationPolicySlice, false)
	var RetentionPolicySlice []string
	plan.RetentionPolicy.ElementsAs(ctx, &RetentionPolicySlice, false)
	addRequest := client.NewAddThirdPartyFileBasedAccessLogPublisherRequest(plan.Id.ValueString(),
		[]client.EnumthirdPartyFileBasedAccessLogPublisherSchemaUrn{client.ENUMTHIRDPARTYFILEBASEDACCESSLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERTHIRD_PARTY_FILE_BASED_ACCESS},
		plan.LogFile.ValueString(),
		plan.LogFilePermissions.ValueString(),
		RotationPolicySlice,
		RetentionPolicySlice,
		plan.ExtensionClass.ValueString(),
		plan.Asynchronous.ValueBool(),
		plan.Enabled.ValueBool())
	err := addOptionalThirdPartyFileBasedAccessLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Third Party File Based Access Log Publisher", err.Error())
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
		client.AddThirdPartyFileBasedAccessLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Third Party File Based Access Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state thirdPartyFileBasedAccessLogPublisherResourceModel
	readThirdPartyFileBasedAccessLogPublisherResponse(ctx, addResponse.ThirdPartyFileBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)

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
func (r *thirdPartyFileBasedAccessLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state thirdPartyFileBasedAccessLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Third Party File Based Access Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readThirdPartyFileBasedAccessLogPublisherResponse(ctx, readResponse.ThirdPartyFileBasedAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *thirdPartyFileBasedAccessLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan thirdPartyFileBasedAccessLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state thirdPartyFileBasedAccessLogPublisherResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.LogPublisherApi.UpdateLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createThirdPartyFileBasedAccessLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Third Party File Based Access Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readThirdPartyFileBasedAccessLogPublisherResponse(ctx, updateResponse.ThirdPartyFileBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *thirdPartyFileBasedAccessLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state thirdPartyFileBasedAccessLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogPublisherApi.DeleteLogPublisherExecute(r.apiClient.LogPublisherApi.DeleteLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Third Party File Based Access Log Publisher", err, httpResp)
		return
	}
}

func (r *thirdPartyFileBasedAccessLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
