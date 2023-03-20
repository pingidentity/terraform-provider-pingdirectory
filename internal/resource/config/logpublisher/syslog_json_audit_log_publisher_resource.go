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
	_ resource.Resource                = &syslogJsonAuditLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &syslogJsonAuditLogPublisherResource{}
	_ resource.ResourceWithImportState = &syslogJsonAuditLogPublisherResource{}
	_ resource.Resource                = &defaultSyslogJsonAuditLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &defaultSyslogJsonAuditLogPublisherResource{}
	_ resource.ResourceWithImportState = &defaultSyslogJsonAuditLogPublisherResource{}
)

// Create a Syslog Json Audit Log Publisher resource
func NewSyslogJsonAuditLogPublisherResource() resource.Resource {
	return &syslogJsonAuditLogPublisherResource{}
}

func NewDefaultSyslogJsonAuditLogPublisherResource() resource.Resource {
	return &defaultSyslogJsonAuditLogPublisherResource{}
}

// syslogJsonAuditLogPublisherResource is the resource implementation.
type syslogJsonAuditLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSyslogJsonAuditLogPublisherResource is the resource implementation.
type defaultSyslogJsonAuditLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *syslogJsonAuditLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_syslog_json_audit_log_publisher"
}

func (r *defaultSyslogJsonAuditLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_syslog_json_audit_log_publisher"
}

// Configure adds the provider configured client to the resource.
func (r *syslogJsonAuditLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultSyslogJsonAuditLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type syslogJsonAuditLogPublisherResourceModel struct {
	Id                                      types.String `tfsdk:"id"`
	LastUpdated                             types.String `tfsdk:"last_updated"`
	Notifications                           types.Set    `tfsdk:"notifications"`
	RequiredActions                         types.Set    `tfsdk:"required_actions"`
	SyslogExternalServer                    types.Set    `tfsdk:"syslog_external_server"`
	SyslogFacility                          types.String `tfsdk:"syslog_facility"`
	SyslogSeverity                          types.String `tfsdk:"syslog_severity"`
	SyslogMessageHostName                   types.String `tfsdk:"syslog_message_host_name"`
	SyslogMessageApplicationName            types.String `tfsdk:"syslog_message_application_name"`
	QueueSize                               types.Int64  `tfsdk:"queue_size"`
	WriteMultiLineMessages                  types.Bool   `tfsdk:"write_multi_line_messages"`
	UseReversibleForm                       types.Bool   `tfsdk:"use_reversible_form"`
	SoftDeleteEntryAuditBehavior            types.String `tfsdk:"soft_delete_entry_audit_behavior"`
	IncludeOperationPurposeRequestControl   types.Bool   `tfsdk:"include_operation_purpose_request_control"`
	IncludeIntermediateClientRequestControl types.Bool   `tfsdk:"include_intermediate_client_request_control"`
	ObscureAttribute                        types.Set    `tfsdk:"obscure_attribute"`
	ExcludeAttribute                        types.Set    `tfsdk:"exclude_attribute"`
	SuppressInternalOperations              types.Bool   `tfsdk:"suppress_internal_operations"`
	IncludeProductName                      types.Bool   `tfsdk:"include_product_name"`
	IncludeInstanceName                     types.Bool   `tfsdk:"include_instance_name"`
	IncludeStartupID                        types.Bool   `tfsdk:"include_startup_id"`
	IncludeThreadID                         types.Bool   `tfsdk:"include_thread_id"`
	IncludeRequesterDN                      types.Bool   `tfsdk:"include_requester_dn"`
	IncludeRequesterIPAddress               types.Bool   `tfsdk:"include_requester_ip_address"`
	IncludeRequestControls                  types.Bool   `tfsdk:"include_request_controls"`
	IncludeResponseControls                 types.Bool   `tfsdk:"include_response_controls"`
	IncludeReplicationChangeID              types.Bool   `tfsdk:"include_replication_change_id"`
	LogSecurityNegotiation                  types.Bool   `tfsdk:"log_security_negotiation"`
	SuppressReplicationOperations           types.Bool   `tfsdk:"suppress_replication_operations"`
	ConnectionCriteria                      types.String `tfsdk:"connection_criteria"`
	RequestCriteria                         types.String `tfsdk:"request_criteria"`
	ResultCriteria                          types.String `tfsdk:"result_criteria"`
	Description                             types.String `tfsdk:"description"`
	Enabled                                 types.Bool   `tfsdk:"enabled"`
	LoggingErrorBehavior                    types.String `tfsdk:"logging_error_behavior"`
}

// GetSchema defines the schema for the resource.
func (r *syslogJsonAuditLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	syslogJsonAuditLogPublisherSchema(ctx, req, resp, false)
}

func (r *defaultSyslogJsonAuditLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	syslogJsonAuditLogPublisherSchema(ctx, req, resp, true)
}

func syslogJsonAuditLogPublisherSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Syslog Json Audit Log Publisher.",
		Attributes: map[string]schema.Attribute{
			"syslog_external_server": schema.SetAttribute{
				Description: "The syslog server to which messages should be sent.",
				Required:    true,
				ElementType: types.StringType,
			},
			"syslog_facility": schema.StringAttribute{
				Description: "The syslog facility to use for the messages that are logged by this Syslog JSON Audit Log Publisher.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"syslog_severity": schema.StringAttribute{
				Description: "The syslog severity to use for the messages that are logged by this Syslog JSON Audit Log Publisher.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"syslog_message_host_name": schema.StringAttribute{
				Description: "The local host name that will be included in syslog messages that are logged by this Syslog JSON Audit Log Publisher.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"syslog_message_application_name": schema.StringAttribute{
				Description: "The application name that will be included in syslog messages that are logged by this Syslog JSON Audit Log Publisher.",
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
			"write_multi_line_messages": schema.BoolAttribute{
				Description: "Indicates whether the JSON objects should use a multi-line representation (with each object field and array value on its own line) that may be easier for administrators to read, but each message will be larger (because of additional spaces and end-of-line markers), and it may be more difficult to consume and parse through some text-oriented tools.",
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
			"suppress_internal_operations": schema.BoolAttribute{
				Description: "Indicates whether internal operations (for example, operations that are initiated by plugins) should be logged along with the operations that are requested by users.",
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
			"include_requester_dn": schema.BoolAttribute{
				Description: "Indicates whether log messages for operation requests should include the DN of the authenticated user for the client connection on which the operation was requested.",
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
			"include_request_controls": schema.BoolAttribute{
				Description: "Indicates whether log messages for operation requests should include a list of the OIDs of any controls included in the request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_response_controls": schema.BoolAttribute{
				Description: "Indicates whether log messages for operation results should include a list of the OIDs of any controls included in the result.",
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
			"log_security_negotiation": schema.BoolAttribute{
				Description: "Indicates whether to log information about the result of any security negotiation (e.g., SSL handshake) processing that has been performed.",
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
func addOptionalSyslogJsonAuditLogPublisherFields(ctx context.Context, addRequest *client.AddSyslogJsonAuditLogPublisherRequest, plan syslogJsonAuditLogPublisherResourceModel) error {
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
	if internaltypes.IsDefined(plan.WriteMultiLineMessages) {
		boolVal := plan.WriteMultiLineMessages.ValueBool()
		addRequest.WriteMultiLineMessages = &boolVal
	}
	if internaltypes.IsDefined(plan.UseReversibleForm) {
		boolVal := plan.UseReversibleForm.ValueBool()
		addRequest.UseReversibleForm = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SoftDeleteEntryAuditBehavior) {
		softDeleteEntryAuditBehavior, err := client.NewEnumlogPublisherSyslogJsonAuditSoftDeleteEntryAuditBehaviorPropFromValue(plan.SoftDeleteEntryAuditBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.SoftDeleteEntryAuditBehavior = softDeleteEntryAuditBehavior
	}
	if internaltypes.IsDefined(plan.IncludeOperationPurposeRequestControl) {
		boolVal := plan.IncludeOperationPurposeRequestControl.ValueBool()
		addRequest.IncludeOperationPurposeRequestControl = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeIntermediateClientRequestControl) {
		boolVal := plan.IncludeIntermediateClientRequestControl.ValueBool()
		addRequest.IncludeIntermediateClientRequestControl = &boolVal
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
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		boolVal := plan.SuppressInternalOperations.ValueBool()
		addRequest.SuppressInternalOperations = &boolVal
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
	if internaltypes.IsDefined(plan.IncludeRequesterDN) {
		boolVal := plan.IncludeRequesterDN.ValueBool()
		addRequest.IncludeRequesterDN = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeRequesterIPAddress) {
		boolVal := plan.IncludeRequesterIPAddress.ValueBool()
		addRequest.IncludeRequesterIPAddress = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeRequestControls) {
		boolVal := plan.IncludeRequestControls.ValueBool()
		addRequest.IncludeRequestControls = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeResponseControls) {
		boolVal := plan.IncludeResponseControls.ValueBool()
		addRequest.IncludeResponseControls = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeReplicationChangeID) {
		boolVal := plan.IncludeReplicationChangeID.ValueBool()
		addRequest.IncludeReplicationChangeID = &boolVal
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		boolVal := plan.LogSecurityNegotiation.ValueBool()
		addRequest.LogSecurityNegotiation = &boolVal
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		boolVal := plan.SuppressReplicationOperations.ValueBool()
		addRequest.SuppressReplicationOperations = &boolVal
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

// Read a SyslogJsonAuditLogPublisherResponse object into the model struct
func readSyslogJsonAuditLogPublisherResponse(ctx context.Context, r *client.SyslogJsonAuditLogPublisherResponse, state *syslogJsonAuditLogPublisherResourceModel, expectedValues *syslogJsonAuditLogPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.SyslogExternalServer = internaltypes.GetStringSet(r.SyslogExternalServer)
	state.SyslogFacility = types.StringValue(r.SyslogFacility.String())
	state.SyslogSeverity = types.StringValue(r.SyslogSeverity.String())
	state.SyslogMessageHostName = internaltypes.StringTypeOrNil(r.SyslogMessageHostName, internaltypes.IsEmptyString(expectedValues.SyslogMessageHostName))
	state.SyslogMessageApplicationName = internaltypes.StringTypeOrNil(r.SyslogMessageApplicationName, internaltypes.IsEmptyString(expectedValues.SyslogMessageApplicationName))
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.UseReversibleForm = internaltypes.BoolTypeOrNil(r.UseReversibleForm)
	state.SoftDeleteEntryAuditBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherSyslogJsonAuditSoftDeleteEntryAuditBehaviorProp(r.SoftDeleteEntryAuditBehavior), internaltypes.IsEmptyString(expectedValues.SoftDeleteEntryAuditBehavior))
	state.IncludeOperationPurposeRequestControl = internaltypes.BoolTypeOrNil(r.IncludeOperationPurposeRequestControl)
	state.IncludeIntermediateClientRequestControl = internaltypes.BoolTypeOrNil(r.IncludeIntermediateClientRequestControl)
	state.ObscureAttribute = internaltypes.GetStringSet(r.ObscureAttribute)
	state.ExcludeAttribute = internaltypes.GetStringSet(r.ExcludeAttribute)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequesterDN = internaltypes.BoolTypeOrNil(r.IncludeRequesterDN)
	state.IncludeRequesterIPAddress = internaltypes.BoolTypeOrNil(r.IncludeRequesterIPAddress)
	state.IncludeRequestControls = internaltypes.BoolTypeOrNil(r.IncludeRequestControls)
	state.IncludeResponseControls = internaltypes.BoolTypeOrNil(r.IncludeResponseControls)
	state.IncludeReplicationChangeID = internaltypes.BoolTypeOrNil(r.IncludeReplicationChangeID)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
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
func createSyslogJsonAuditLogPublisherOperations(plan syslogJsonAuditLogPublisherResourceModel, state syslogJsonAuditLogPublisherResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SyslogExternalServer, state.SyslogExternalServer, "syslog-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.SyslogFacility, state.SyslogFacility, "syslog-facility")
	operations.AddStringOperationIfNecessary(&ops, plan.SyslogSeverity, state.SyslogSeverity, "syslog-severity")
	operations.AddStringOperationIfNecessary(&ops, plan.SyslogMessageHostName, state.SyslogMessageHostName, "syslog-message-host-name")
	operations.AddStringOperationIfNecessary(&ops, plan.SyslogMessageApplicationName, state.SyslogMessageApplicationName, "syslog-message-application-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.QueueSize, state.QueueSize, "queue-size")
	operations.AddBoolOperationIfNecessary(&ops, plan.WriteMultiLineMessages, state.WriteMultiLineMessages, "write-multi-line-messages")
	operations.AddBoolOperationIfNecessary(&ops, plan.UseReversibleForm, state.UseReversibleForm, "use-reversible-form")
	operations.AddStringOperationIfNecessary(&ops, plan.SoftDeleteEntryAuditBehavior, state.SoftDeleteEntryAuditBehavior, "soft-delete-entry-audit-behavior")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeOperationPurposeRequestControl, state.IncludeOperationPurposeRequestControl, "include-operation-purpose-request-control")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeIntermediateClientRequestControl, state.IncludeIntermediateClientRequestControl, "include-intermediate-client-request-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ObscureAttribute, state.ObscureAttribute, "obscure-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeAttribute, state.ExcludeAttribute, "exclude-attribute")
	operations.AddBoolOperationIfNecessary(&ops, plan.SuppressInternalOperations, state.SuppressInternalOperations, "suppress-internal-operations")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeProductName, state.IncludeProductName, "include-product-name")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeInstanceName, state.IncludeInstanceName, "include-instance-name")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeStartupID, state.IncludeStartupID, "include-startup-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeThreadID, state.IncludeThreadID, "include-thread-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequesterDN, state.IncludeRequesterDN, "include-requester-dn")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequesterIPAddress, state.IncludeRequesterIPAddress, "include-requester-ip-address")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequestControls, state.IncludeRequestControls, "include-request-controls")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeResponseControls, state.IncludeResponseControls, "include-response-controls")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeReplicationChangeID, state.IncludeReplicationChangeID, "include-replication-change-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogSecurityNegotiation, state.LogSecurityNegotiation, "log-security-negotiation")
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
func (r *syslogJsonAuditLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan syslogJsonAuditLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var SyslogExternalServerSlice []string
	plan.SyslogExternalServer.ElementsAs(ctx, &SyslogExternalServerSlice, false)
	addRequest := client.NewAddSyslogJsonAuditLogPublisherRequest(plan.Id.ValueString(),
		[]client.EnumsyslogJsonAuditLogPublisherSchemaUrn{client.ENUMSYSLOGJSONAUDITLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERSYSLOG_JSON_AUDIT},
		SyslogExternalServerSlice,
		plan.Enabled.ValueBool())
	err := addOptionalSyslogJsonAuditLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Syslog Json Audit Log Publisher", err.Error())
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
		client.AddSyslogJsonAuditLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Syslog Json Audit Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state syslogJsonAuditLogPublisherResourceModel
	readSyslogJsonAuditLogPublisherResponse(ctx, addResponse.SyslogJsonAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultSyslogJsonAuditLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan syslogJsonAuditLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Syslog Json Audit Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state syslogJsonAuditLogPublisherResourceModel
	readSyslogJsonAuditLogPublisherResponse(ctx, readResponse.SyslogJsonAuditLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogPublisherApi.UpdateLogPublisher(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSyslogJsonAuditLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Syslog Json Audit Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSyslogJsonAuditLogPublisherResponse(ctx, updateResponse.SyslogJsonAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *syslogJsonAuditLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSyslogJsonAuditLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSyslogJsonAuditLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSyslogJsonAuditLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSyslogJsonAuditLogPublisher(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state syslogJsonAuditLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Syslog Json Audit Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSyslogJsonAuditLogPublisherResponse(ctx, readResponse.SyslogJsonAuditLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *syslogJsonAuditLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSyslogJsonAuditLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSyslogJsonAuditLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSyslogJsonAuditLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSyslogJsonAuditLogPublisher(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan syslogJsonAuditLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state syslogJsonAuditLogPublisherResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogPublisherApi.UpdateLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSyslogJsonAuditLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Syslog Json Audit Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSyslogJsonAuditLogPublisherResponse(ctx, updateResponse.SyslogJsonAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSyslogJsonAuditLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *syslogJsonAuditLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state syslogJsonAuditLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogPublisherApi.DeleteLogPublisherExecute(r.apiClient.LogPublisherApi.DeleteLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Syslog Json Audit Log Publisher", err, httpResp)
		return
	}
}

func (r *syslogJsonAuditLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSyslogJsonAuditLogPublisher(ctx, req, resp)
}

func (r *defaultSyslogJsonAuditLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSyslogJsonAuditLogPublisher(ctx, req, resp)
}

func importSyslogJsonAuditLogPublisher(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
