package logpublisher

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
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &fileBasedAuditLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &fileBasedAuditLogPublisherResource{}
	_ resource.ResourceWithImportState = &fileBasedAuditLogPublisherResource{}
	_ resource.Resource                = &defaultFileBasedAuditLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &defaultFileBasedAuditLogPublisherResource{}
	_ resource.ResourceWithImportState = &defaultFileBasedAuditLogPublisherResource{}
)

// Create a File Based Audit Log Publisher resource
func NewFileBasedAuditLogPublisherResource() resource.Resource {
	return &fileBasedAuditLogPublisherResource{}
}

func NewDefaultFileBasedAuditLogPublisherResource() resource.Resource {
	return &defaultFileBasedAuditLogPublisherResource{}
}

// fileBasedAuditLogPublisherResource is the resource implementation.
type fileBasedAuditLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultFileBasedAuditLogPublisherResource is the resource implementation.
type defaultFileBasedAuditLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *fileBasedAuditLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_based_audit_log_publisher"
}

func (r *defaultFileBasedAuditLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_file_based_audit_log_publisher"
}

// Configure adds the provider configured client to the resource.
func (r *fileBasedAuditLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultFileBasedAuditLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type fileBasedAuditLogPublisherResourceModel struct {
	Id                                      types.String `tfsdk:"id"`
	LastUpdated                             types.String `tfsdk:"last_updated"`
	Notifications                           types.Set    `tfsdk:"notifications"`
	RequiredActions                         types.Set    `tfsdk:"required_actions"`
	SuppressInternalOperations              types.Bool   `tfsdk:"suppress_internal_operations"`
	LogFile                                 types.String `tfsdk:"log_file"`
	LogFilePermissions                      types.String `tfsdk:"log_file_permissions"`
	RotationPolicy                          types.Set    `tfsdk:"rotation_policy"`
	RotationListener                        types.Set    `tfsdk:"rotation_listener"`
	RetentionPolicy                         types.Set    `tfsdk:"retention_policy"`
	CompressionMechanism                    types.String `tfsdk:"compression_mechanism"`
	SignLog                                 types.Bool   `tfsdk:"sign_log"`
	EncryptLog                              types.Bool   `tfsdk:"encrypt_log"`
	EncryptionSettingsDefinitionID          types.String `tfsdk:"encryption_settings_definition_id"`
	Append                                  types.Bool   `tfsdk:"append"`
	IncludeProductName                      types.Bool   `tfsdk:"include_product_name"`
	IncludeInstanceName                     types.Bool   `tfsdk:"include_instance_name"`
	IncludeStartupID                        types.Bool   `tfsdk:"include_startup_id"`
	IncludeThreadID                         types.Bool   `tfsdk:"include_thread_id"`
	IncludeRequesterIPAddress               types.Bool   `tfsdk:"include_requester_ip_address"`
	IncludeRequesterDN                      types.Bool   `tfsdk:"include_requester_dn"`
	IncludeReplicationChangeID              types.Bool   `tfsdk:"include_replication_change_id"`
	UseReversibleForm                       types.Bool   `tfsdk:"use_reversible_form"`
	SoftDeleteEntryAuditBehavior            types.String `tfsdk:"soft_delete_entry_audit_behavior"`
	IncludeRequestControls                  types.Bool   `tfsdk:"include_request_controls"`
	IncludeOperationPurposeRequestControl   types.Bool   `tfsdk:"include_operation_purpose_request_control"`
	IncludeIntermediateClientRequestControl types.Bool   `tfsdk:"include_intermediate_client_request_control"`
	ObscureAttribute                        types.Set    `tfsdk:"obscure_attribute"`
	ExcludeAttribute                        types.Set    `tfsdk:"exclude_attribute"`
	Asynchronous                            types.Bool   `tfsdk:"asynchronous"`
	AutoFlush                               types.Bool   `tfsdk:"auto_flush"`
	BufferSize                              types.String `tfsdk:"buffer_size"`
	QueueSize                               types.Int64  `tfsdk:"queue_size"`
	TimeInterval                            types.String `tfsdk:"time_interval"`
	TimestampPrecision                      types.String `tfsdk:"timestamp_precision"`
	LogSecurityNegotiation                  types.Bool   `tfsdk:"log_security_negotiation"`
	LogIntermediateResponses                types.Bool   `tfsdk:"log_intermediate_responses"`
	SuppressReplicationOperations           types.Bool   `tfsdk:"suppress_replication_operations"`
	ConnectionCriteria                      types.String `tfsdk:"connection_criteria"`
	RequestCriteria                         types.String `tfsdk:"request_criteria"`
	ResultCriteria                          types.String `tfsdk:"result_criteria"`
	Description                             types.String `tfsdk:"description"`
	Enabled                                 types.Bool   `tfsdk:"enabled"`
	LoggingErrorBehavior                    types.String `tfsdk:"logging_error_behavior"`
}

// GetSchema defines the schema for the resource.
func (r *fileBasedAuditLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	fileBasedAuditLogPublisherSchema(ctx, req, resp, false)
}

func (r *defaultFileBasedAuditLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	fileBasedAuditLogPublisherSchema(ctx, req, resp, true)
}

func fileBasedAuditLogPublisherSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a File Based Audit Log Publisher.",
		Attributes: map[string]schema.Attribute{
			"suppress_internal_operations": schema.BoolAttribute{
				Description: "Indicates whether internal operations (for example, operations that are initiated by plugins) should be logged along with the operations that are requested by users.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"log_file": schema.StringAttribute{
				Description: "The file name to use for the log files generated by the File Based Audit Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.",
				Required:    true,
			},
			"log_file_permissions": schema.StringAttribute{
				Description: "The UNIX permissions of the log files created by this File Based Audit Log Publisher.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rotation_policy": schema.SetAttribute{
				Description: "The rotation policy to use for the File Based Audit Log Publisher .",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"rotation_listener": schema.SetAttribute{
				Description: "A listener that should be notified whenever a log file is rotated out of service.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"retention_policy": schema.SetAttribute{
				Description: "The retention policy to use for the File Based Audit Log Publisher .",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"compression_mechanism": schema.StringAttribute{
				Description: "Specifies the type of compression (if any) to use for log files that are written.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sign_log": schema.BoolAttribute{
				Description: "Indicates whether the log should be cryptographically signed so that the log content cannot be altered in an undetectable manner.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"encrypt_log": schema.BoolAttribute{
				Description: "Indicates whether log files should be encrypted so that their content is not available to unauthorized users.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"encryption_settings_definition_id": schema.StringAttribute{
				Description: "Specifies the ID of the encryption settings definition that should be used to encrypt the data. If this is not provided, the server's preferred encryption settings definition will be used. The \"encryption-settings list\" command can be used to obtain a list of the encryption settings definitions available in the server.",
				Optional:    true,
			},
			"append": schema.BoolAttribute{
				Description: "Specifies whether to append to existing log files.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_product_name": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the product name for the Directory Server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_instance_name": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the instance name for the Directory Server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_startup_id": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the startup ID for the Directory Server, which is a value assigned to the server instance at startup and may be used to identify when the server has been restarted.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_thread_id": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the thread ID for the Directory Server in each log message. This ID can be used to correlate log messages from the same thread within a single log as well as generated by the same thread across different types of log files. More information about the thread with a specific ID can be obtained using the cn=JVM Stack Trace,cn=monitor entry.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_requester_ip_address": schema.BoolAttribute{
				Description: "Indicates whether log messages for operation requests should include the IP address of the client that requested the operation.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_requester_dn": schema.BoolAttribute{
				Description: "Indicates whether log messages for operation requests should include the DN of the authenticated user for the client connection on which the operation was requested.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_replication_change_id": schema.BoolAttribute{
				Description: "Indicates whether to log information about the replication change ID.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"use_reversible_form": schema.BoolAttribute{
				Description: "Indicates whether the audit log should be written in reversible form so that it is possible to revert the changes if desired.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"soft_delete_entry_audit_behavior": schema.StringAttribute{
				Description: "Specifies the audit behavior for delete and modify operations on soft-deleted entries.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"include_request_controls": schema.BoolAttribute{
				Description: "Indicates whether log messages for operation requests should include a list of the OIDs of any controls included in the request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_operation_purpose_request_control": schema.BoolAttribute{
				Description: "Indicates whether to include information about any operation purpose request control that may have been included in the request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_intermediate_client_request_control": schema.BoolAttribute{
				Description: "Indicates whether to include information about any intermediate client request control that may have been included in the request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"obscure_attribute": schema.SetAttribute{
				Description: "Specifies the names of any attribute types that should have their values obscured in the audit log because they may be considered sensitive.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"exclude_attribute": schema.SetAttribute{
				Description: "Specifies the names of any attribute types that should be excluded from the audit log.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"asynchronous": schema.BoolAttribute{
				Description: "Indicates whether the File Based Audit Log Publisher will publish records asynchronously.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"auto_flush": schema.BoolAttribute{
				Description: "Specifies whether to flush the writer after every log record.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"buffer_size": schema.StringAttribute{
				Description: "Specifies the log file buffer size.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"queue_size": schema.Int64Attribute{
				Description: "The maximum number of log records that can be stored in the asynchronous queue.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"time_interval": schema.StringAttribute{
				Description: "Specifies the interval at which to check whether the log files need to be rotated.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"timestamp_precision": schema.StringAttribute{
				Description: "Specifies the smallest time unit to be included in timestamps.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"log_security_negotiation": schema.BoolAttribute{
				Description: "Indicates whether to log information about the result of any security negotiation (e.g., SSL handshake) processing that has been performed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"log_intermediate_responses": schema.BoolAttribute{
				Description: "Indicates whether to log information about intermediate responses sent to the client.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"suppress_replication_operations": schema.BoolAttribute{
				Description: "Indicates whether access messages that are generated by replication operations should be suppressed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
func addOptionalFileBasedAuditLogPublisherFields(ctx context.Context, addRequest *client.AddFileBasedAuditLogPublisherRequest, plan fileBasedAuditLogPublisherResourceModel) error {
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		addRequest.SuppressInternalOperations = plan.SuppressInternalOperations.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		addRequest.LogFilePermissions = plan.LogFilePermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
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
		addRequest.SignLog = plan.SignLog.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		addRequest.EncryptLog = plan.EncryptLog.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Append) {
		addRequest.Append = plan.Append.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterIPAddress) {
		addRequest.IncludeRequesterIPAddress = plan.IncludeRequesterIPAddress.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterDN) {
		addRequest.IncludeRequesterDN = plan.IncludeRequesterDN.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeReplicationChangeID) {
		addRequest.IncludeReplicationChangeID = plan.IncludeReplicationChangeID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.UseReversibleForm) {
		addRequest.UseReversibleForm = plan.UseReversibleForm.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SoftDeleteEntryAuditBehavior) {
		softDeleteEntryAuditBehavior, err := client.NewEnumlogPublisherFileBasedAuditSoftDeleteEntryAuditBehaviorPropFromValue(plan.SoftDeleteEntryAuditBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.SoftDeleteEntryAuditBehavior = softDeleteEntryAuditBehavior
	}
	if internaltypes.IsDefined(plan.IncludeRequestControls) {
		addRequest.IncludeRequestControls = plan.IncludeRequestControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeOperationPurposeRequestControl) {
		addRequest.IncludeOperationPurposeRequestControl = plan.IncludeOperationPurposeRequestControl.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeIntermediateClientRequestControl) {
		addRequest.IncludeIntermediateClientRequestControl = plan.IncludeIntermediateClientRequestControl.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ObscureAttribute) {
		var slice []string
		plan.ObscureAttribute.ElementsAs(ctx, &slice, false)
		addRequest.ObscureAttribute = slice
	}
	if internaltypes.IsDefined(plan.ExcludeAttribute) {
		var slice []string
		plan.ExcludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.ExcludeAttribute = slice
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		addRequest.BufferSize = plan.BufferSize.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		addRequest.TimeInterval = plan.TimeInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimestampPrecision) {
		timestampPrecision, err := client.NewEnumlogPublisherTimestampPrecisionPropFromValue(plan.TimestampPrecision.ValueString())
		if err != nil {
			return err
		}
		addRequest.TimestampPrecision = timestampPrecision
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		addRequest.LogSecurityNegotiation = plan.LogSecurityNegotiation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogIntermediateResponses) {
		addRequest.LogIntermediateResponses = plan.LogIntermediateResponses.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		addRequest.SuppressReplicationOperations = plan.SuppressReplicationOperations.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		addRequest.ConnectionCriteria = plan.ConnectionCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		addRequest.ResultCriteria = plan.ResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
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

// Read a FileBasedAuditLogPublisherResponse object into the model struct
func readFileBasedAuditLogPublisherResponse(ctx context.Context, r *client.FileBasedAuditLogPublisherResponse, state *fileBasedAuditLogPublisherResourceModel, expectedValues *fileBasedAuditLogPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
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
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequesterIPAddress = internaltypes.BoolTypeOrNil(r.IncludeRequesterIPAddress)
	state.IncludeRequesterDN = internaltypes.BoolTypeOrNil(r.IncludeRequesterDN)
	state.IncludeReplicationChangeID = internaltypes.BoolTypeOrNil(r.IncludeReplicationChangeID)
	state.UseReversibleForm = internaltypes.BoolTypeOrNil(r.UseReversibleForm)
	state.SoftDeleteEntryAuditBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherFileBasedAuditSoftDeleteEntryAuditBehaviorProp(r.SoftDeleteEntryAuditBehavior), internaltypes.IsEmptyString(expectedValues.SoftDeleteEntryAuditBehavior))
	state.IncludeRequestControls = internaltypes.BoolTypeOrNil(r.IncludeRequestControls)
	state.IncludeOperationPurposeRequestControl = internaltypes.BoolTypeOrNil(r.IncludeOperationPurposeRequestControl)
	state.IncludeIntermediateClientRequestControl = internaltypes.BoolTypeOrNil(r.IncludeIntermediateClientRequestControl)
	state.ObscureAttribute = internaltypes.GetStringSet(r.ObscureAttribute)
	state.ExcludeAttribute = internaltypes.GetStringSet(r.ExcludeAttribute)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, internaltypes.IsEmptyString(expectedValues.BufferSize))
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, internaltypes.IsEmptyString(expectedValues.TimeInterval))
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.TimestampPrecision = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherTimestampPrecisionProp(r.TimestampPrecision), internaltypes.IsEmptyString(expectedValues.TimestampPrecision))
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createFileBasedAuditLogPublisherOperations(plan fileBasedAuditLogPublisherResourceModel, state fileBasedAuditLogPublisherResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.SuppressInternalOperations, state.SuppressInternalOperations, "suppress-internal-operations")
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
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeProductName, state.IncludeProductName, "include-product-name")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeInstanceName, state.IncludeInstanceName, "include-instance-name")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeStartupID, state.IncludeStartupID, "include-startup-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeThreadID, state.IncludeThreadID, "include-thread-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequesterIPAddress, state.IncludeRequesterIPAddress, "include-requester-ip-address")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequesterDN, state.IncludeRequesterDN, "include-requester-dn")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeReplicationChangeID, state.IncludeReplicationChangeID, "include-replication-change-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.UseReversibleForm, state.UseReversibleForm, "use-reversible-form")
	operations.AddStringOperationIfNecessary(&ops, plan.SoftDeleteEntryAuditBehavior, state.SoftDeleteEntryAuditBehavior, "soft-delete-entry-audit-behavior")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequestControls, state.IncludeRequestControls, "include-request-controls")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeOperationPurposeRequestControl, state.IncludeOperationPurposeRequestControl, "include-operation-purpose-request-control")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeIntermediateClientRequestControl, state.IncludeIntermediateClientRequestControl, "include-intermediate-client-request-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ObscureAttribute, state.ObscureAttribute, "obscure-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeAttribute, state.ExcludeAttribute, "exclude-attribute")
	operations.AddBoolOperationIfNecessary(&ops, plan.Asynchronous, state.Asynchronous, "asynchronous")
	operations.AddBoolOperationIfNecessary(&ops, plan.AutoFlush, state.AutoFlush, "auto-flush")
	operations.AddStringOperationIfNecessary(&ops, plan.BufferSize, state.BufferSize, "buffer-size")
	operations.AddInt64OperationIfNecessary(&ops, plan.QueueSize, state.QueueSize, "queue-size")
	operations.AddStringOperationIfNecessary(&ops, plan.TimeInterval, state.TimeInterval, "time-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.TimestampPrecision, state.TimestampPrecision, "timestamp-precision")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogSecurityNegotiation, state.LogSecurityNegotiation, "log-security-negotiation")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogIntermediateResponses, state.LogIntermediateResponses, "log-intermediate-responses")
	operations.AddBoolOperationIfNecessary(&ops, plan.SuppressReplicationOperations, state.SuppressReplicationOperations, "suppress-replication-operations")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionCriteria, state.ConnectionCriteria, "connection-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.RequestCriteria, state.RequestCriteria, "request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.ResultCriteria, state.ResultCriteria, "result-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.LoggingErrorBehavior, state.LoggingErrorBehavior, "logging-error-behavior")
	return ops
}

// Create a new resource
func (r *fileBasedAuditLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan fileBasedAuditLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddFileBasedAuditLogPublisherRequest(plan.Id.ValueString(),
		[]client.EnumfileBasedAuditLogPublisherSchemaUrn{client.ENUMFILEBASEDAUDITLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERFILE_BASED_AUDIT},
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalFileBasedAuditLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for File Based Audit Log Publisher", err.Error())
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
		client.AddFileBasedAuditLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the File Based Audit Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state fileBasedAuditLogPublisherResourceModel
	readFileBasedAuditLogPublisherResponse(ctx, addResponse.FileBasedAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultFileBasedAuditLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan fileBasedAuditLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the File Based Audit Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state fileBasedAuditLogPublisherResourceModel
	readFileBasedAuditLogPublisherResponse(ctx, readResponse.FileBasedAuditLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogPublisherApi.UpdateLogPublisher(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createFileBasedAuditLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the File Based Audit Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFileBasedAuditLogPublisherResponse(ctx, updateResponse.FileBasedAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *fileBasedAuditLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFileBasedAuditLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFileBasedAuditLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFileBasedAuditLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readFileBasedAuditLogPublisher(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state fileBasedAuditLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the File Based Audit Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readFileBasedAuditLogPublisherResponse(ctx, readResponse.FileBasedAuditLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *fileBasedAuditLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFileBasedAuditLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFileBasedAuditLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFileBasedAuditLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateFileBasedAuditLogPublisher(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan fileBasedAuditLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state fileBasedAuditLogPublisherResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogPublisherApi.UpdateLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createFileBasedAuditLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the File Based Audit Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFileBasedAuditLogPublisherResponse(ctx, updateResponse.FileBasedAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultFileBasedAuditLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *fileBasedAuditLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state fileBasedAuditLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogPublisherApi.DeleteLogPublisherExecute(r.apiClient.LogPublisherApi.DeleteLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the File Based Audit Log Publisher", err, httpResp)
		return
	}
}

func (r *fileBasedAuditLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFileBasedAuditLogPublisher(ctx, req, resp)
}

func (r *defaultFileBasedAuditLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFileBasedAuditLogPublisher(ctx, req, resp)
}

func importFileBasedAuditLogPublisher(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
