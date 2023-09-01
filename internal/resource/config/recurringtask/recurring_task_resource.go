package recurringtask

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &recurringTaskResource{}
	_ resource.ResourceWithConfigure   = &recurringTaskResource{}
	_ resource.ResourceWithImportState = &recurringTaskResource{}
	_ resource.Resource                = &defaultRecurringTaskResource{}
	_ resource.ResourceWithConfigure   = &defaultRecurringTaskResource{}
	_ resource.ResourceWithImportState = &defaultRecurringTaskResource{}
)

// Create a Recurring Task resource
func NewRecurringTaskResource() resource.Resource {
	return &recurringTaskResource{}
}

func NewDefaultRecurringTaskResource() resource.Resource {
	return &defaultRecurringTaskResource{}
}

// recurringTaskResource is the resource implementation.
type recurringTaskResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultRecurringTaskResource is the resource implementation.
type defaultRecurringTaskResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *recurringTaskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_recurring_task"
}

func (r *defaultRecurringTaskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_recurring_task"
}

// Configure adds the provider configured client to the resource.
func (r *recurringTaskResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultRecurringTaskResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type recurringTaskResourceModel struct {
	Id                                      types.String `tfsdk:"id"`
	Name                                    types.String `tfsdk:"name"`
	Notifications                           types.Set    `tfsdk:"notifications"`
	RequiredActions                         types.Set    `tfsdk:"required_actions"`
	Type                                    types.String `tfsdk:"type"`
	ExtensionClass                          types.String `tfsdk:"extension_class"`
	ExtensionArgument                       types.Set    `tfsdk:"extension_argument"`
	TargetDirectory                         types.String `tfsdk:"target_directory"`
	FilenamePattern                         types.String `tfsdk:"filename_pattern"`
	TimestampFormat                         types.String `tfsdk:"timestamp_format"`
	RetainFileCount                         types.Int64  `tfsdk:"retain_file_count"`
	RetainFileAge                           types.String `tfsdk:"retain_file_age"`
	RetainAggregateFileSize                 types.String `tfsdk:"retain_aggregate_file_size"`
	CommandPath                             types.String `tfsdk:"command_path"`
	CommandArguments                        types.String `tfsdk:"command_arguments"`
	CommandOutputFileBaseName               types.String `tfsdk:"command_output_file_base_name"`
	RetainPreviousOutputFileCount           types.Int64  `tfsdk:"retain_previous_output_file_count"`
	RetainPreviousOutputFileAge             types.String `tfsdk:"retain_previous_output_file_age"`
	LogCommandOutput                        types.Bool   `tfsdk:"log_command_output"`
	TaskCompletionStateForNonzeroExitCode   types.String `tfsdk:"task_completion_state_for_nonzero_exit_code"`
	WorkingDirectory                        types.String `tfsdk:"working_directory"`
	BaseOutputDirectory                     types.String `tfsdk:"base_output_directory"`
	DataSecurityAuditor                     types.Set    `tfsdk:"data_security_auditor"`
	Backend                                 types.Set    `tfsdk:"backend"`
	IncludeFilter                           types.Set    `tfsdk:"include_filter"`
	RetainPreviousReportCount               types.Int64  `tfsdk:"retain_previous_report_count"`
	RetainPreviousReportAge                 types.String `tfsdk:"retain_previous_report_age"`
	LdifDirectory                           types.String `tfsdk:"ldif_directory"`
	BackendID                               types.Set    `tfsdk:"backend_id"`
	ExcludeBackendID                        types.Set    `tfsdk:"exclude_backend_id"`
	OutputDirectory                         types.String `tfsdk:"output_directory"`
	EncryptionPassphraseFile                types.String `tfsdk:"encryption_passphrase_file"`
	IncludeExpensiveData                    types.Bool   `tfsdk:"include_expensive_data"`
	IncludeReplicationStateDump             types.Bool   `tfsdk:"include_replication_state_dump"`
	RetainPreviousLDIFExportCount           types.Int64  `tfsdk:"retain_previous_ldif_export_count"`
	RetainPreviousLDIFExportAge             types.String `tfsdk:"retain_previous_ldif_export_age"`
	IncludeBinaryFiles                      types.Bool   `tfsdk:"include_binary_files"`
	IncludeExtensionSource                  types.Bool   `tfsdk:"include_extension_source"`
	UseSequentialMode                       types.Bool   `tfsdk:"use_sequential_mode"`
	SecurityLevel                           types.String `tfsdk:"security_level"`
	JstackCount                             types.Int64  `tfsdk:"jstack_count"`
	ReportCount                             types.Int64  `tfsdk:"report_count"`
	ReportIntervalSeconds                   types.Int64  `tfsdk:"report_interval_seconds"`
	LogDuration                             types.String `tfsdk:"log_duration"`
	LogFileHeadCollectionSize               types.String `tfsdk:"log_file_head_collection_size"`
	LogFileTailCollectionSize               types.String `tfsdk:"log_file_tail_collection_size"`
	Comment                                 types.String `tfsdk:"comment"`
	RetainPreviousSupportDataArchiveCount   types.Int64  `tfsdk:"retain_previous_support_data_archive_count"`
	RetainPreviousSupportDataArchiveAge     types.String `tfsdk:"retain_previous_support_data_archive_age"`
	TaskJavaClass                           types.String `tfsdk:"task_java_class"`
	TaskObjectClass                         types.Set    `tfsdk:"task_object_class"`
	TaskAttributeValue                      types.Set    `tfsdk:"task_attribute_value"`
	SleepDuration                           types.String `tfsdk:"sleep_duration"`
	DurationToWaitForWorkQueueIdle          types.String `tfsdk:"duration_to_wait_for_work_queue_idle"`
	LdapURLForSearchExpectedToReturnEntries types.Set    `tfsdk:"ldap_url_for_search_expected_to_return_entries"`
	SearchInterval                          types.String `tfsdk:"search_interval"`
	SearchTimeLimit                         types.String `tfsdk:"search_time_limit"`
	DurationToWaitForSearchToReturnEntries  types.String `tfsdk:"duration_to_wait_for_search_to_return_entries"`
	TaskReturnStateIfTimeoutIsEncountered   types.String `tfsdk:"task_return_state_if_timeout_is_encountered"`
	BackupDirectory                         types.String `tfsdk:"backup_directory"`
	IncludedBackendID                       types.Set    `tfsdk:"included_backend_id"`
	ExcludedBackendID                       types.Set    `tfsdk:"excluded_backend_id"`
	Compress                                types.Bool   `tfsdk:"compress"`
	Encrypt                                 types.Bool   `tfsdk:"encrypt"`
	EncryptionSettingsDefinitionID          types.String `tfsdk:"encryption_settings_definition_id"`
	Sign                                    types.Bool   `tfsdk:"sign"`
	RetainPreviousFullBackupCount           types.Int64  `tfsdk:"retain_previous_full_backup_count"`
	RetainPreviousFullBackupAge             types.String `tfsdk:"retain_previous_full_backup_age"`
	MaxMegabytesPerSecond                   types.Int64  `tfsdk:"max_megabytes_per_second"`
	Reason                                  types.String `tfsdk:"reason"`
	ProfileDirectory                        types.String `tfsdk:"profile_directory"`
	IncludePath                             types.Set    `tfsdk:"include_path"`
	RetainPreviousProfileCount              types.Int64  `tfsdk:"retain_previous_profile_count"`
	RetainPreviousProfileAge                types.String `tfsdk:"retain_previous_profile_age"`
	Description                             types.String `tfsdk:"description"`
	CancelOnTaskDependencyFailure           types.Bool   `tfsdk:"cancel_on_task_dependency_failure"`
	EmailOnStart                            types.Set    `tfsdk:"email_on_start"`
	EmailOnSuccess                          types.Set    `tfsdk:"email_on_success"`
	EmailOnFailure                          types.Set    `tfsdk:"email_on_failure"`
	AlertOnStart                            types.Bool   `tfsdk:"alert_on_start"`
	AlertOnSuccess                          types.Bool   `tfsdk:"alert_on_success"`
	AlertOnFailure                          types.Bool   `tfsdk:"alert_on_failure"`
}

// GetSchema defines the schema for the resource.
func (r *recurringTaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	recurringTaskSchema(ctx, req, resp, false)
}

func (r *defaultRecurringTaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	recurringTaskSchema(ctx, req, resp, true)
}

func recurringTaskSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Recurring Task.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Recurring Task resource. Options are ['generate-server-profile', 'leave-lockdown-mode', 'backup', 'delay', 'statically-defined', 'collect-support-data', 'ldif-export', 'enter-lockdown-mode', 'audit-data-security', 'exec', 'file-retention', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"generate-server-profile", "leave-lockdown-mode", "backup", "delay", "statically-defined", "collect-support-data", "ldif-export", "enter-lockdown-mode", "audit-data-security", "exec", "file-retention", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Recurring Task.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Recurring Task. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"target_directory": schema.StringAttribute{
				Description: "The path to the directory containing the files to examine. The directory must exist.",
				Optional:    true,
			},
			"filename_pattern": schema.StringAttribute{
				Description: "A pattern that specifies the names of the files to examine. The pattern may contain zero or more asterisks as wildcards, where each wildcard matches zero or more characters. It may also contain at most one occurrence of the special string \"${timestamp}\", which will match a timestamp with the format specified using the timestamp-format property. All other characters in the pattern will be treated literally.",
				Optional:    true,
			},
			"timestamp_format": schema.StringAttribute{
				Description: "The format to use for the timestamp represented by the \"${timestamp}\" token in the filename pattern.",
				Optional:    true,
			},
			"retain_file_count": schema.Int64Attribute{
				Description: "The minimum number of files matching the pattern that will be retained.",
				Optional:    true,
			},
			"retain_file_age": schema.StringAttribute{
				Description: "The minimum age of files matching the pattern that will be retained.",
				Optional:    true,
			},
			"retain_aggregate_file_size": schema.StringAttribute{
				Description: "The minimum aggregate size of files that will be retained. The size should be specified as an integer followed by a unit that is one of \"b\" or \"bytes\", \"kb\" or \"kilobytes\", \"mb\" or \"megabytes\", \"gb\" or \"gigabytes\", or \"tb\" or \"terabytes\". For example, a value of \"1 gb\" indicates that at least one gigabyte of files should be retained.",
				Optional:    true,
			},
			"command_path": schema.StringAttribute{
				Description: "The absolute path to the command to execute. It must be an absolute path, the corresponding file must exist, and it must be listed in the config/exec-command-whitelist.txt file.",
				Optional:    true,
			},
			"command_arguments": schema.StringAttribute{
				Description: "A string containing the arguments to provide to the command. If the command should be run without arguments, this property should be left undefined. If there should be multiple arguments, then they should be separated with spaces.",
				Optional:    true,
			},
			"command_output_file_base_name": schema.StringAttribute{
				Description: "The path and base name for a file to which the command output (both standard output and standard error) should be written. This may be left undefined if the command output should not be recorded into a file.",
				Optional:    true,
			},
			"retain_previous_output_file_count": schema.Int64Attribute{
				Description: "The minimum number of previous command output files that should be preserved after a new instance of the command is invoked.",
				Optional:    true,
			},
			"retain_previous_output_file_age": schema.StringAttribute{
				Description: "The minimum age of previous command output files that should be preserved after a new instance of the command is invoked.",
				Optional:    true,
			},
			"log_command_output": schema.BoolAttribute{
				Description: "Indicates whether the command's output (both standard output and standard error) should be recorded in the server's error log.",
				Optional:    true,
				Computed:    true,
			},
			"task_completion_state_for_nonzero_exit_code": schema.StringAttribute{
				Description: "The final task state that a task instance should have if the task executes the specified command and that command completes with a nonzero exit code, which generally means that the command did not complete successfully.",
				Optional:    true,
				Computed:    true,
			},
			"working_directory": schema.StringAttribute{
				Description: "The absolute path to a working directory where the command should be executed. It must be an absolute path and the corresponding directory must exist.",
				Optional:    true,
			},
			"base_output_directory": schema.StringAttribute{
				Description: "The base directory below which generated reports will be written. Each invocation of the audit-data-security task will create a new subdirectory below this base directory whose name is a timestamp indicating when the report was generated.",
				Optional:    true,
				Computed:    true,
			},
			"data_security_auditor": schema.SetAttribute{
				Description: "The set of data security auditors that should be invoked. If no auditors are specified, then all auditors defined in the configuration will be used.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"backend": schema.SetAttribute{
				Description: "The set of backends that should be examined. If no backends are specified, then all backends that support this functionality will be included.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"include_filter": schema.SetAttribute{
				Description: "A filter that will be used to identify entries that may be included in the generated report. If multiple filters are specified, then any entry that matches at least one of the filters will be included. If no filters are specified, then all entries will be included.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"retain_previous_report_count": schema.Int64Attribute{
				Description: "The minimum number of previous reports that should be preserved after a new report is generated.",
				Optional:    true,
			},
			"retain_previous_report_age": schema.StringAttribute{
				Description: "The minimum age of previous reports that should be preserved after a new report completes successfully.",
				Optional:    true,
			},
			"ldif_directory": schema.StringAttribute{
				Description: "The directory in which LDIF export files will be placed. The directory must already exist.",
				Optional:    true,
				Computed:    true,
			},
			"backend_id": schema.SetAttribute{
				Description: "The backend ID for a backend to be exported.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"exclude_backend_id": schema.SetAttribute{
				Description: "The backend ID for a backend to be excluded from the export.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"output_directory": schema.StringAttribute{
				Description: "The directory in which the support data archive files will be placed. The path must be a directory, and that directory must already exist. Relative paths will be interpreted as relative to the server root.",
				Optional:    true,
			},
			"encryption_passphrase_file": schema.StringAttribute{
				Description: "The path to a file that contains the passphrase to encrypt the contents of the support data archive.",
				Optional:    true,
			},
			"include_expensive_data": schema.BoolAttribute{
				Description: "Indicates whether the support data archive should include information that may be expensive to obtain, and that may temporarily affect the server's performance or responsiveness.",
				Optional:    true,
				Computed:    true,
			},
			"include_replication_state_dump": schema.BoolAttribute{
				Description: "Indicates whether the support data archive should include a replication state dump, which may be several megabytes in size.",
				Optional:    true,
				Computed:    true,
			},
			"retain_previous_ldif_export_count": schema.Int64Attribute{
				Description: "The minimum number of previous LDIF exports that should be preserved after a new export completes successfully.",
				Optional:    true,
			},
			"retain_previous_ldif_export_age": schema.StringAttribute{
				Description: "The minimum age of previous LDIF exports that should be preserved after a new export completes successfully.",
				Optional:    true,
			},
			"include_binary_files": schema.BoolAttribute{
				Description: "Indicates whether the support data archive should include binary files that may not have otherwise been included. Note that it may not be possible to obscure or redact sensitive information in binary files.",
				Optional:    true,
				Computed:    true,
			},
			"include_extension_source": schema.BoolAttribute{
				Description: "Indicates whether the support data archive should include the source code (if available) for any third-party extensions that may be installed in the server.",
				Optional:    true,
				Computed:    true,
			},
			"use_sequential_mode": schema.BoolAttribute{
				Description: "Indicates whether to capture support data information sequentially rather than in parallel. Capturing data in sequential mode may reduce the amount of memory that the tool requires to operate, at the cost of taking longer to run.",
				Optional:    true,
				Computed:    true,
			},
			"security_level": schema.StringAttribute{
				Description: "The security level to use when deciding which information to include in or exclude from the support data archive, and which included data should be obscured or redacted.",
				Optional:    true,
				Computed:    true,
			},
			"jstack_count": schema.Int64Attribute{
				Description: "The number of times to invoke the jstack utility to obtain a stack trace of all threads running in the JVM. A value of zero indicates that the jstack utility should not be invoked.",
				Optional:    true,
				Computed:    true,
			},
			"report_count": schema.Int64Attribute{
				Description: "The number of intervals of data to collect from tools that use sample-based reporting, like vmstat, iostat, and mpstat. A value of zero indicates that these kinds of tools should not be used to collect any information.",
				Optional:    true,
				Computed:    true,
			},
			"report_interval_seconds": schema.Int64Attribute{
				Description: "The duration (in seconds) between each interval of data to collect from tools that use sample-based reporting, like vmstat, iostat, and mpstat.",
				Optional:    true,
				Computed:    true,
			},
			"log_duration": schema.StringAttribute{
				Description: "The maximum age (leading up to the time the collect-support-data tool was invoked) for log content to include in the support data archive.",
				Optional:    true,
			},
			"log_file_head_collection_size": schema.StringAttribute{
				Description: "The amount of data to collect from the beginning of each log file included in the support data archive.",
				Optional:    true,
			},
			"log_file_tail_collection_size": schema.StringAttribute{
				Description: "The amount of data to collect from the end of each log file included in the support data archive.",
				Optional:    true,
			},
			"comment": schema.StringAttribute{
				Description: "An optional comment to include in a README file within the support data archive.",
				Optional:    true,
			},
			"retain_previous_support_data_archive_count": schema.Int64Attribute{
				Description: "The minimum number of previous support data archives that should be preserved after a new archive is generated.",
				Optional:    true,
			},
			"retain_previous_support_data_archive_age": schema.StringAttribute{
				Description: "The minimum age of previous support data archives that should be preserved after a new archive is generated.",
				Optional:    true,
			},
			"task_java_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class that provides the logic for the task to be invoked.",
				Optional:    true,
			},
			"task_object_class": schema.SetAttribute{
				Description: "The names or OIDs of the object classes to include in the tasks that are scheduled from this Statically Defined Recurring Task. All object classes must be defined in the server schema, and the combination of object classes must be valid for a task entry.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"task_attribute_value": schema.SetAttribute{
				Description: "The set of attribute values that should be included in the tasks that are scheduled from this Statically Defined Recurring Task. Each value must be in the form {attribute-type}={value}, where {attribute-type} is the name or OID of an attribute type that is defined in the schema and permitted with the configured set of object classes, and {value} is a value to assign to an attribute with that type. A multivalued attribute can be created by providing multiple name-value pairs with the same name and different values.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"sleep_duration": schema.StringAttribute{
				Description: "The length of time to sleep before the task completes.",
				Optional:    true,
			},
			"duration_to_wait_for_work_queue_idle": schema.StringAttribute{
				Description: "Indicates that task should wait for up to the specified length of time for the work queue to report that all worker threads are idle and there are no pending operations. Note that this primarily monitors operations that use worker threads, which does not include internal operations (for example, those invoked by extensions), and may not include requests from non-LDAP clients (for example, HTTP-based clients).",
				Optional:    true,
			},
			"ldap_url_for_search_expected_to_return_entries": schema.SetAttribute{
				Description: "An LDAP URL that provides the criteria for a search request that is expected to return at least one entry. The search will be performed internally, and only the base DN, scope, and filter from the URL will be used; any host, port, or requested attributes included in the URL will be ignored.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"search_interval": schema.StringAttribute{
				Description: "The length of time the server should sleep between searches performed using the criteria from the ldap-url-for-search-expected-to-return-entries property.",
				Optional:    true,
			},
			"search_time_limit": schema.StringAttribute{
				Description: "The length of time that the server will wait for a response to each internal search performed using the criteria from the ldap-url-for-search-expected-to-return-entries property.",
				Optional:    true,
			},
			"duration_to_wait_for_search_to_return_entries": schema.StringAttribute{
				Description: "The maximum length of time that the server will continue to perform internal searches using the criteria from the ldap-url-for-search-expected-to-return-entries property.",
				Optional:    true,
			},
			"task_return_state_if_timeout_is_encountered": schema.StringAttribute{
				Description: "The return state to use if a timeout is encountered while waiting for the server work queue to become idle (if the duration-to-wait-for-work-queue-idle property has a value), or if the time specified by the duration-to-wait-for-search-to-return-entries elapses without the associated search returning any entries.",
				Optional:    true,
				Computed:    true,
			},
			"backup_directory": schema.StringAttribute{
				Description: "The directory in which backup files will be placed. When backing up a single backend, the backup files will be placed directly in this directory. When backing up multiple backends, the backup files for each backend will be placed in a subdirectory whose name is the corresponding backend ID.",
				Optional:    true,
				Computed:    true,
			},
			"included_backend_id": schema.SetAttribute{
				Description: "The backend IDs of any backends that should be included in the backup.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"excluded_backend_id": schema.SetAttribute{
				Description: "The backend IDs of any backends that should be excluded from the backup. All backends that support backups and are not listed will be included.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"compress": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to `backup`: Indicates whether to compress the data as it is written into the backup. When the `type` attribute is set to `ldif-export`: Indicates whether to compress the LDIF data as it is exported.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `backup`: Indicates whether to compress the data as it is written into the backup.\n  - `ldif-export`: Indicates whether to compress the LDIF data as it is exported.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"encrypt": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to `backup`: Indicates whether to encrypt the data as it is written into the backup. When the `type` attribute is set to `ldif-export`: Indicates whether to encrypt the LDIF data as it exported.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `backup`: Indicates whether to encrypt the data as it is written into the backup.\n  - `ldif-export`: Indicates whether to encrypt the LDIF data as it exported.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"encryption_settings_definition_id": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `backup`: The ID of an encryption settings definition to use to obtain the backup encryption key. When the `type` attribute is set to `ldif-export`: The ID of an encryption settings definition to use to obtain the LDIF export encryption key.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `backup`: The ID of an encryption settings definition to use to obtain the backup encryption key.\n  - `ldif-export`: The ID of an encryption settings definition to use to obtain the LDIF export encryption key.",
				Optional:            true,
			},
			"sign": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to `backup`: Indicates whether to cryptographically sign backups, which will make it possible to detect whether the backup has been altered since it was created. When the `type` attribute is set to `ldif-export`: Indicates whether to cryptographically sign the exported data, which will make it possible to detect whether the LDIF data has been altered since it was exported.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `backup`: Indicates whether to cryptographically sign backups, which will make it possible to detect whether the backup has been altered since it was created.\n  - `ldif-export`: Indicates whether to cryptographically sign the exported data, which will make it possible to detect whether the LDIF data has been altered since it was exported.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"retain_previous_full_backup_count": schema.Int64Attribute{
				Description: "The minimum number of previous full backups that should be preserved after a new backup completes successfully.",
				Optional:    true,
			},
			"retain_previous_full_backup_age": schema.StringAttribute{
				Description: "The minimum age of previous full backups that should be preserved after a new backup completes successfully.",
				Optional:    true,
			},
			"max_megabytes_per_second": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to `backup`: The maximum rate, in megabytes per second, at which backups should be written. When the `type` attribute is set to `ldif-export`: The maximum rate, in megabytes per second, at which LDIF exports should be written.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `backup`: The maximum rate, in megabytes per second, at which backups should be written.\n  - `ldif-export`: The maximum rate, in megabytes per second, at which LDIF exports should be written.",
				Optional:            true,
			},
			"reason": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `leave-lockdown-mode`: The reason that the server is being taken out of in lockdown mode. When the `type` attribute is set to `enter-lockdown-mode`: The reason that the server is being placed in lockdown mode.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `leave-lockdown-mode`: The reason that the server is being taken out of in lockdown mode.\n  - `enter-lockdown-mode`: The reason that the server is being placed in lockdown mode.",
				Optional:            true,
			},
			"profile_directory": schema.StringAttribute{
				Description: "The directory in which the generated server profiles will be placed. The files will be named with the pattern \"server-profile-{timestamp}.zip\", where \"{timestamp}\" represents the time that the profile was generated.",
				Optional:    true,
			},
			"include_path": schema.SetAttribute{
				Description: "An optional set of additional paths to files within the instance root that should be included in the generated server profile. All paths must be within the instance root, and relative paths will be relative to the instance root.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"retain_previous_profile_count": schema.Int64Attribute{
				Description: "The minimum number of previous server profile zip files that should be preserved after a new profile is generated.",
				Optional:    true,
			},
			"retain_previous_profile_age": schema.StringAttribute{
				Description: "The minimum age of previous server profile zip files that should be preserved after a new profile is generated.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Recurring Task",
				Optional:    true,
			},
			"cancel_on_task_dependency_failure": schema.BoolAttribute{
				Description: "Indicates whether an instance of this Recurring Task should be canceled if the task immediately before it in the recurring task chain fails to complete successfully (including if it is canceled by an administrator before it starts or while it is running).",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"email_on_start": schema.SetAttribute{
				Description: "The email addresses to which a message should be sent whenever an instance of this Recurring Task starts running. If this option is used, then at least one smtp-server must be configured in the global configuration.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"email_on_success": schema.SetAttribute{
				Description: "The email addresses to which a message should be sent whenever an instance of this Recurring Task completes successfully. If this option is used, then at least one smtp-server must be configured in the global configuration.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"email_on_failure": schema.SetAttribute{
				Description: "The email addresses to which a message should be sent if an instance of this Recurring Task fails to complete successfully. If this option is used, then at least one smtp-server must be configured in the global configuration.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"alert_on_start": schema.BoolAttribute{
				Description: "Indicates whether the server should generate an administrative alert whenever an instance of this Recurring Task starts running.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"alert_on_success": schema.BoolAttribute{
				Description: "Indicates whether the server should generate an administrative alert whenever an instance of this Recurring Task completes successfully.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"alert_on_failure": schema.BoolAttribute{
				Description: "Indicates whether the server should generate an administrative alert whenever an instance of this Recurring Task fails to complete successfully.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *recurringTaskResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanRecurringTask(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_recurring_task")
	var planModel, configModel recurringTaskResourceModel
	req.Config.Get(ctx, &configModel)
	req.Plan.Get(ctx, &planModel)
	resourceType := planModel.Type.ValueString()
	anyDefaultsSet := false
	// Set defaults for backup type
	if resourceType == "backup" {
		if !internaltypes.IsDefined(configModel.BackupDirectory) {
			defaultVal := types.StringValue("bak")
			if !planModel.BackupDirectory.Equal(defaultVal) {
				planModel.BackupDirectory = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for delay type
	if resourceType == "delay" {
		if !internaltypes.IsDefined(configModel.TaskReturnStateIfTimeoutIsEncountered) {
			defaultVal := types.StringValue("stopped-by-error")
			if !planModel.TaskReturnStateIfTimeoutIsEncountered.Equal(defaultVal) {
				planModel.TaskReturnStateIfTimeoutIsEncountered = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for collect-support-data type
	if resourceType == "collect-support-data" {
		if !internaltypes.IsDefined(configModel.IncludeExpensiveData) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeExpensiveData.Equal(defaultVal) {
				planModel.IncludeExpensiveData = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeReplicationStateDump) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeReplicationStateDump.Equal(defaultVal) {
				planModel.IncludeReplicationStateDump = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeBinaryFiles) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeBinaryFiles.Equal(defaultVal) {
				planModel.IncludeBinaryFiles = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeExtensionSource) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeExtensionSource.Equal(defaultVal) {
				planModel.IncludeExtensionSource = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.UseSequentialMode) {
			defaultVal := types.BoolValue(true)
			if !planModel.UseSequentialMode.Equal(defaultVal) {
				planModel.UseSequentialMode = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SecurityLevel) {
			defaultVal := types.StringValue("obscure-secrets")
			if !planModel.SecurityLevel.Equal(defaultVal) {
				planModel.SecurityLevel = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.JstackCount) {
			defaultVal := types.Int64Value(10)
			if !planModel.JstackCount.Equal(defaultVal) {
				planModel.JstackCount = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ReportCount) {
			defaultVal := types.Int64Value(10)
			if !planModel.ReportCount.Equal(defaultVal) {
				planModel.ReportCount = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ReportIntervalSeconds) {
			defaultVal := types.Int64Value(1)
			if !planModel.ReportIntervalSeconds.Equal(defaultVal) {
				planModel.ReportIntervalSeconds = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for ldif-export type
	if resourceType == "ldif-export" {
		if !internaltypes.IsDefined(configModel.LdifDirectory) {
			defaultVal := types.StringValue("ldif")
			if !planModel.LdifDirectory.Equal(defaultVal) {
				planModel.LdifDirectory = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for audit-data-security type
	if resourceType == "audit-data-security" {
		if !internaltypes.IsDefined(configModel.BaseOutputDirectory) {
			defaultVal := types.StringValue("reports/audit-data-security")
			if !planModel.BaseOutputDirectory.Equal(defaultVal) {
				planModel.BaseOutputDirectory = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for exec type
	if resourceType == "exec" {
		if !internaltypes.IsDefined(configModel.LogCommandOutput) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogCommandOutput.Equal(defaultVal) {
				planModel.LogCommandOutput = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.TaskCompletionStateForNonzeroExitCode) {
			defaultVal := types.StringValue("stopped-by-error")
			if !planModel.TaskCompletionStateForNonzeroExitCode.Equal(defaultVal) {
				planModel.TaskCompletionStateForNonzeroExitCode = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	if anyDefaultsSet {
		planModel.Notifications = types.SetUnknown(types.StringType)
		planModel.RequiredActions = types.SetUnknown(config.GetRequiredActionsObjectType())
	}
	resp.Plan.Set(ctx, &planModel)
}

func (r *defaultRecurringTaskResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanRecurringTask(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_default_recurring_task")
}

func modifyPlanRecurringTask(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, resourceName string) {
	compare, err := version.Compare(providerConfig.ProductVersion, version.PingDirectory9200)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model recurringTaskResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.Type) && model.Type.ValueString() == "audit-data-security" {
		version.CheckResourceSupported(&resp.Diagnostics, version.PingDirectory9200,
			providerConfig.ProductVersion, resourceName+" with type \"audit_data_security\"")
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsRecurringTask() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"backup"},
			resourcevalidator.Conflicting(
				path.MatchRoot("included_backend_id"),
				path.MatchRoot("excluded_backend_id"),
			),
		),
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"delay"},
			resourcevalidator.AtLeastOneOf(
				path.MatchRoot("sleep_duration"),
				path.MatchRoot("duration_to_wait_for_work_queue_idle"),
				path.MatchRoot("ldap_url_for_search_expected_to_return_entries"),
			),
		),
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"ldif-export"},
			resourcevalidator.Conflicting(
				path.MatchRoot("backend_id"),
				path.MatchRoot("exclude_backend_id"),
			),
		),
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"backup", "ldif-export"},
			resourcevalidator.Conflicting(
				path.MatchRoot("encryption_passphrase_file"),
				path.MatchRoot("encryption_settings_definition_id"),
			),
		),
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"file-retention"},
			resourcevalidator.AtLeastOneOf(
				path.MatchRoot("retain_file_count"),
				path.MatchRoot("retain_file_age"),
				path.MatchRoot("retain_aggregate_file_size"),
			),
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("profile_directory"),
			path.MatchRoot("type"),
			[]string{"generate-server-profile"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_path"),
			path.MatchRoot("type"),
			[]string{"generate-server-profile"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("retain_previous_profile_count"),
			path.MatchRoot("type"),
			[]string{"generate-server-profile"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("retain_previous_profile_age"),
			path.MatchRoot("type"),
			[]string{"generate-server-profile"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("reason"),
			path.MatchRoot("type"),
			[]string{"leave-lockdown-mode", "enter-lockdown-mode"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("backup_directory"),
			path.MatchRoot("type"),
			[]string{"backup"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("included_backend_id"),
			path.MatchRoot("type"),
			[]string{"backup"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("excluded_backend_id"),
			path.MatchRoot("type"),
			[]string{"backup"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("compress"),
			path.MatchRoot("type"),
			[]string{"backup", "ldif-export"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("encrypt"),
			path.MatchRoot("type"),
			[]string{"backup", "ldif-export"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("encryption_settings_definition_id"),
			path.MatchRoot("type"),
			[]string{"backup", "ldif-export"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("sign"),
			path.MatchRoot("type"),
			[]string{"backup", "ldif-export"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("retain_previous_full_backup_count"),
			path.MatchRoot("type"),
			[]string{"backup"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("retain_previous_full_backup_age"),
			path.MatchRoot("type"),
			[]string{"backup"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("max_megabytes_per_second"),
			path.MatchRoot("type"),
			[]string{"backup", "ldif-export"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("sleep_duration"),
			path.MatchRoot("type"),
			[]string{"delay"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("duration_to_wait_for_work_queue_idle"),
			path.MatchRoot("type"),
			[]string{"delay"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("ldap_url_for_search_expected_to_return_entries"),
			path.MatchRoot("type"),
			[]string{"delay"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("search_interval"),
			path.MatchRoot("type"),
			[]string{"delay"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("search_time_limit"),
			path.MatchRoot("type"),
			[]string{"delay"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("duration_to_wait_for_search_to_return_entries"),
			path.MatchRoot("type"),
			[]string{"delay"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("task_return_state_if_timeout_is_encountered"),
			path.MatchRoot("type"),
			[]string{"delay"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("task_java_class"),
			path.MatchRoot("type"),
			[]string{"statically-defined"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("task_object_class"),
			path.MatchRoot("type"),
			[]string{"statically-defined"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("task_attribute_value"),
			path.MatchRoot("type"),
			[]string{"statically-defined"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("output_directory"),
			path.MatchRoot("type"),
			[]string{"collect-support-data"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("encryption_passphrase_file"),
			path.MatchRoot("type"),
			[]string{"collect-support-data"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_expensive_data"),
			path.MatchRoot("type"),
			[]string{"collect-support-data"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_replication_state_dump"),
			path.MatchRoot("type"),
			[]string{"collect-support-data"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_binary_files"),
			path.MatchRoot("type"),
			[]string{"collect-support-data"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_extension_source"),
			path.MatchRoot("type"),
			[]string{"collect-support-data"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("use_sequential_mode"),
			path.MatchRoot("type"),
			[]string{"collect-support-data"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("security_level"),
			path.MatchRoot("type"),
			[]string{"collect-support-data"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("jstack_count"),
			path.MatchRoot("type"),
			[]string{"collect-support-data"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("report_count"),
			path.MatchRoot("type"),
			[]string{"collect-support-data"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("report_interval_seconds"),
			path.MatchRoot("type"),
			[]string{"collect-support-data"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_duration"),
			path.MatchRoot("type"),
			[]string{"collect-support-data"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_file_head_collection_size"),
			path.MatchRoot("type"),
			[]string{"collect-support-data"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_file_tail_collection_size"),
			path.MatchRoot("type"),
			[]string{"collect-support-data"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("comment"),
			path.MatchRoot("type"),
			[]string{"collect-support-data"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("retain_previous_support_data_archive_count"),
			path.MatchRoot("type"),
			[]string{"collect-support-data"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("retain_previous_support_data_archive_age"),
			path.MatchRoot("type"),
			[]string{"collect-support-data"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("ldif_directory"),
			path.MatchRoot("type"),
			[]string{"ldif-export"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("backend_id"),
			path.MatchRoot("type"),
			[]string{"ldif-export"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("exclude_backend_id"),
			path.MatchRoot("type"),
			[]string{"ldif-export"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("retain_previous_ldif_export_count"),
			path.MatchRoot("type"),
			[]string{"ldif-export"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("retain_previous_ldif_export_age"),
			path.MatchRoot("type"),
			[]string{"ldif-export"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("base_output_directory"),
			path.MatchRoot("type"),
			[]string{"audit-data-security"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("data_security_auditor"),
			path.MatchRoot("type"),
			[]string{"audit-data-security"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("backend"),
			path.MatchRoot("type"),
			[]string{"audit-data-security"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_filter"),
			path.MatchRoot("type"),
			[]string{"audit-data-security"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("retain_previous_report_count"),
			path.MatchRoot("type"),
			[]string{"audit-data-security"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("retain_previous_report_age"),
			path.MatchRoot("type"),
			[]string{"audit-data-security"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("command_path"),
			path.MatchRoot("type"),
			[]string{"exec"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("command_arguments"),
			path.MatchRoot("type"),
			[]string{"exec"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("command_output_file_base_name"),
			path.MatchRoot("type"),
			[]string{"exec"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("retain_previous_output_file_count"),
			path.MatchRoot("type"),
			[]string{"exec"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("retain_previous_output_file_age"),
			path.MatchRoot("type"),
			[]string{"exec"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_command_output"),
			path.MatchRoot("type"),
			[]string{"exec"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("task_completion_state_for_nonzero_exit_code"),
			path.MatchRoot("type"),
			[]string{"exec"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("working_directory"),
			path.MatchRoot("type"),
			[]string{"exec"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("target_directory"),
			path.MatchRoot("type"),
			[]string{"file-retention"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("filename_pattern"),
			path.MatchRoot("type"),
			[]string{"file-retention"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("timestamp_format"),
			path.MatchRoot("type"),
			[]string{"file-retention"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("retain_file_count"),
			path.MatchRoot("type"),
			[]string{"file-retention"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("retain_file_age"),
			path.MatchRoot("type"),
			[]string{"file-retention"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("retain_aggregate_file_size"),
			path.MatchRoot("type"),
			[]string{"file-retention"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_class"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_argument"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"generate-server-profile",
			[]path.Expression{path.MatchRoot("profile_directory")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"statically-defined",
			[]path.Expression{path.MatchRoot("task_java_class"), path.MatchRoot("task_object_class")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"collect-support-data",
			[]path.Expression{path.MatchRoot("output_directory")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"exec",
			[]path.Expression{path.MatchRoot("command_path")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"file-retention",
			[]path.Expression{path.MatchRoot("target_directory"), path.MatchRoot("filename_pattern"), path.MatchRoot("timestamp_format")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party",
			[]path.Expression{path.MatchRoot("extension_class")},
		),
	}
}

// Add config validators
func (r recurringTaskResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsRecurringTask()
}

// Add config validators
func (r defaultRecurringTaskResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsRecurringTask()
}

// Add optional fields to create request for generate-server-profile recurring-task
func addOptionalGenerateServerProfileRecurringTaskFields(ctx context.Context, addRequest *client.AddGenerateServerProfileRecurringTaskRequest, plan recurringTaskResourceModel) error {
	if internaltypes.IsDefined(plan.IncludePath) {
		var slice []string
		plan.IncludePath.ElementsAs(ctx, &slice, false)
		addRequest.IncludePath = slice
	}
	if internaltypes.IsDefined(plan.RetainPreviousProfileCount) {
		addRequest.RetainPreviousProfileCount = plan.RetainPreviousProfileCount.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RetainPreviousProfileAge) {
		addRequest.RetainPreviousProfileAge = plan.RetainPreviousProfileAge.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CancelOnTaskDependencyFailure) {
		addRequest.CancelOnTaskDependencyFailure = plan.CancelOnTaskDependencyFailure.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EmailOnStart) {
		var slice []string
		plan.EmailOnStart.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnStart = slice
	}
	if internaltypes.IsDefined(plan.EmailOnSuccess) {
		var slice []string
		plan.EmailOnSuccess.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnSuccess = slice
	}
	if internaltypes.IsDefined(plan.EmailOnFailure) {
		var slice []string
		plan.EmailOnFailure.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnFailure = slice
	}
	if internaltypes.IsDefined(plan.AlertOnStart) {
		addRequest.AlertOnStart = plan.AlertOnStart.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnSuccess) {
		addRequest.AlertOnSuccess = plan.AlertOnSuccess.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnFailure) {
		addRequest.AlertOnFailure = plan.AlertOnFailure.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for leave-lockdown-mode recurring-task
func addOptionalLeaveLockdownModeRecurringTaskFields(ctx context.Context, addRequest *client.AddLeaveLockdownModeRecurringTaskRequest, plan recurringTaskResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Reason) {
		addRequest.Reason = plan.Reason.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CancelOnTaskDependencyFailure) {
		addRequest.CancelOnTaskDependencyFailure = plan.CancelOnTaskDependencyFailure.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EmailOnStart) {
		var slice []string
		plan.EmailOnStart.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnStart = slice
	}
	if internaltypes.IsDefined(plan.EmailOnSuccess) {
		var slice []string
		plan.EmailOnSuccess.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnSuccess = slice
	}
	if internaltypes.IsDefined(plan.EmailOnFailure) {
		var slice []string
		plan.EmailOnFailure.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnFailure = slice
	}
	if internaltypes.IsDefined(plan.AlertOnStart) {
		addRequest.AlertOnStart = plan.AlertOnStart.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnSuccess) {
		addRequest.AlertOnSuccess = plan.AlertOnSuccess.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnFailure) {
		addRequest.AlertOnFailure = plan.AlertOnFailure.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for backup recurring-task
func addOptionalBackupRecurringTaskFields(ctx context.Context, addRequest *client.AddBackupRecurringTaskRequest, plan recurringTaskResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BackupDirectory) {
		addRequest.BackupDirectory = plan.BackupDirectory.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludedBackendID) {
		var slice []string
		plan.IncludedBackendID.ElementsAs(ctx, &slice, false)
		addRequest.IncludedBackendID = slice
	}
	if internaltypes.IsDefined(plan.ExcludedBackendID) {
		var slice []string
		plan.ExcludedBackendID.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedBackendID = slice
	}
	if internaltypes.IsDefined(plan.Compress) {
		addRequest.Compress = plan.Compress.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Encrypt) {
		addRequest.Encrypt = plan.Encrypt.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Sign) {
		addRequest.Sign = plan.Sign.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.RetainPreviousFullBackupCount) {
		addRequest.RetainPreviousFullBackupCount = plan.RetainPreviousFullBackupCount.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RetainPreviousFullBackupAge) {
		addRequest.RetainPreviousFullBackupAge = plan.RetainPreviousFullBackupAge.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.MaxMegabytesPerSecond) {
		addRequest.MaxMegabytesPerSecond = plan.MaxMegabytesPerSecond.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CancelOnTaskDependencyFailure) {
		addRequest.CancelOnTaskDependencyFailure = plan.CancelOnTaskDependencyFailure.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EmailOnStart) {
		var slice []string
		plan.EmailOnStart.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnStart = slice
	}
	if internaltypes.IsDefined(plan.EmailOnSuccess) {
		var slice []string
		plan.EmailOnSuccess.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnSuccess = slice
	}
	if internaltypes.IsDefined(plan.EmailOnFailure) {
		var slice []string
		plan.EmailOnFailure.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnFailure = slice
	}
	if internaltypes.IsDefined(plan.AlertOnStart) {
		addRequest.AlertOnStart = plan.AlertOnStart.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnSuccess) {
		addRequest.AlertOnSuccess = plan.AlertOnSuccess.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnFailure) {
		addRequest.AlertOnFailure = plan.AlertOnFailure.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for delay recurring-task
func addOptionalDelayRecurringTaskFields(ctx context.Context, addRequest *client.AddDelayRecurringTaskRequest, plan recurringTaskResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SleepDuration) {
		addRequest.SleepDuration = plan.SleepDuration.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DurationToWaitForWorkQueueIdle) {
		addRequest.DurationToWaitForWorkQueueIdle = plan.DurationToWaitForWorkQueueIdle.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.LdapURLForSearchExpectedToReturnEntries) {
		var slice []string
		plan.LdapURLForSearchExpectedToReturnEntries.ElementsAs(ctx, &slice, false)
		addRequest.LdapURLForSearchExpectedToReturnEntries = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchInterval) {
		addRequest.SearchInterval = plan.SearchInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchTimeLimit) {
		addRequest.SearchTimeLimit = plan.SearchTimeLimit.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DurationToWaitForSearchToReturnEntries) {
		addRequest.DurationToWaitForSearchToReturnEntries = plan.DurationToWaitForSearchToReturnEntries.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TaskReturnStateIfTimeoutIsEncountered) {
		taskReturnStateIfTimeoutIsEncountered, err := client.NewEnumrecurringTaskTaskReturnStateIfTimeoutIsEncounteredPropFromValue(plan.TaskReturnStateIfTimeoutIsEncountered.ValueString())
		if err != nil {
			return err
		}
		addRequest.TaskReturnStateIfTimeoutIsEncountered = taskReturnStateIfTimeoutIsEncountered
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CancelOnTaskDependencyFailure) {
		addRequest.CancelOnTaskDependencyFailure = plan.CancelOnTaskDependencyFailure.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EmailOnStart) {
		var slice []string
		plan.EmailOnStart.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnStart = slice
	}
	if internaltypes.IsDefined(plan.EmailOnSuccess) {
		var slice []string
		plan.EmailOnSuccess.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnSuccess = slice
	}
	if internaltypes.IsDefined(plan.EmailOnFailure) {
		var slice []string
		plan.EmailOnFailure.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnFailure = slice
	}
	if internaltypes.IsDefined(plan.AlertOnStart) {
		addRequest.AlertOnStart = plan.AlertOnStart.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnSuccess) {
		addRequest.AlertOnSuccess = plan.AlertOnSuccess.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnFailure) {
		addRequest.AlertOnFailure = plan.AlertOnFailure.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for statically-defined recurring-task
func addOptionalStaticallyDefinedRecurringTaskFields(ctx context.Context, addRequest *client.AddStaticallyDefinedRecurringTaskRequest, plan recurringTaskResourceModel) error {
	if internaltypes.IsDefined(plan.TaskAttributeValue) {
		var slice []string
		plan.TaskAttributeValue.ElementsAs(ctx, &slice, false)
		addRequest.TaskAttributeValue = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CancelOnTaskDependencyFailure) {
		addRequest.CancelOnTaskDependencyFailure = plan.CancelOnTaskDependencyFailure.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EmailOnStart) {
		var slice []string
		plan.EmailOnStart.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnStart = slice
	}
	if internaltypes.IsDefined(plan.EmailOnSuccess) {
		var slice []string
		plan.EmailOnSuccess.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnSuccess = slice
	}
	if internaltypes.IsDefined(plan.EmailOnFailure) {
		var slice []string
		plan.EmailOnFailure.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnFailure = slice
	}
	if internaltypes.IsDefined(plan.AlertOnStart) {
		addRequest.AlertOnStart = plan.AlertOnStart.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnSuccess) {
		addRequest.AlertOnSuccess = plan.AlertOnSuccess.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnFailure) {
		addRequest.AlertOnFailure = plan.AlertOnFailure.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for collect-support-data recurring-task
func addOptionalCollectSupportDataRecurringTaskFields(ctx context.Context, addRequest *client.AddCollectSupportDataRecurringTaskRequest, plan recurringTaskResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionPassphraseFile) {
		addRequest.EncryptionPassphraseFile = plan.EncryptionPassphraseFile.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludeExpensiveData) {
		addRequest.IncludeExpensiveData = plan.IncludeExpensiveData.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeReplicationStateDump) {
		addRequest.IncludeReplicationStateDump = plan.IncludeReplicationStateDump.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeBinaryFiles) {
		addRequest.IncludeBinaryFiles = plan.IncludeBinaryFiles.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeExtensionSource) {
		addRequest.IncludeExtensionSource = plan.IncludeExtensionSource.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.UseSequentialMode) {
		addRequest.UseSequentialMode = plan.UseSequentialMode.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SecurityLevel) {
		securityLevel, err := client.NewEnumrecurringTaskSecurityLevelPropFromValue(plan.SecurityLevel.ValueString())
		if err != nil {
			return err
		}
		addRequest.SecurityLevel = securityLevel
	}
	if internaltypes.IsDefined(plan.JstackCount) {
		addRequest.JstackCount = plan.JstackCount.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.ReportCount) {
		addRequest.ReportCount = plan.ReportCount.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.ReportIntervalSeconds) {
		addRequest.ReportIntervalSeconds = plan.ReportIntervalSeconds.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogDuration) {
		addRequest.LogDuration = plan.LogDuration.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFileHeadCollectionSize) {
		addRequest.LogFileHeadCollectionSize = plan.LogFileHeadCollectionSize.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFileTailCollectionSize) {
		addRequest.LogFileTailCollectionSize = plan.LogFileTailCollectionSize.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Comment) {
		addRequest.Comment = plan.Comment.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RetainPreviousSupportDataArchiveCount) {
		addRequest.RetainPreviousSupportDataArchiveCount = plan.RetainPreviousSupportDataArchiveCount.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RetainPreviousSupportDataArchiveAge) {
		addRequest.RetainPreviousSupportDataArchiveAge = plan.RetainPreviousSupportDataArchiveAge.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CancelOnTaskDependencyFailure) {
		addRequest.CancelOnTaskDependencyFailure = plan.CancelOnTaskDependencyFailure.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EmailOnStart) {
		var slice []string
		plan.EmailOnStart.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnStart = slice
	}
	if internaltypes.IsDefined(plan.EmailOnSuccess) {
		var slice []string
		plan.EmailOnSuccess.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnSuccess = slice
	}
	if internaltypes.IsDefined(plan.EmailOnFailure) {
		var slice []string
		plan.EmailOnFailure.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnFailure = slice
	}
	if internaltypes.IsDefined(plan.AlertOnStart) {
		addRequest.AlertOnStart = plan.AlertOnStart.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnSuccess) {
		addRequest.AlertOnSuccess = plan.AlertOnSuccess.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnFailure) {
		addRequest.AlertOnFailure = plan.AlertOnFailure.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for ldif-export recurring-task
func addOptionalLdifExportRecurringTaskFields(ctx context.Context, addRequest *client.AddLdifExportRecurringTaskRequest, plan recurringTaskResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LdifDirectory) {
		addRequest.LdifDirectory = plan.LdifDirectory.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.BackendID) {
		var slice []string
		plan.BackendID.ElementsAs(ctx, &slice, false)
		addRequest.BackendID = slice
	}
	if internaltypes.IsDefined(plan.ExcludeBackendID) {
		var slice []string
		plan.ExcludeBackendID.ElementsAs(ctx, &slice, false)
		addRequest.ExcludeBackendID = slice
	}
	if internaltypes.IsDefined(plan.Compress) {
		addRequest.Compress = plan.Compress.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Encrypt) {
		addRequest.Encrypt = plan.Encrypt.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Sign) {
		addRequest.Sign = plan.Sign.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.RetainPreviousLDIFExportCount) {
		addRequest.RetainPreviousLDIFExportCount = plan.RetainPreviousLDIFExportCount.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RetainPreviousLDIFExportAge) {
		addRequest.RetainPreviousLDIFExportAge = plan.RetainPreviousLDIFExportAge.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.MaxMegabytesPerSecond) {
		addRequest.MaxMegabytesPerSecond = plan.MaxMegabytesPerSecond.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CancelOnTaskDependencyFailure) {
		addRequest.CancelOnTaskDependencyFailure = plan.CancelOnTaskDependencyFailure.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EmailOnStart) {
		var slice []string
		plan.EmailOnStart.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnStart = slice
	}
	if internaltypes.IsDefined(plan.EmailOnSuccess) {
		var slice []string
		plan.EmailOnSuccess.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnSuccess = slice
	}
	if internaltypes.IsDefined(plan.EmailOnFailure) {
		var slice []string
		plan.EmailOnFailure.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnFailure = slice
	}
	if internaltypes.IsDefined(plan.AlertOnStart) {
		addRequest.AlertOnStart = plan.AlertOnStart.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnSuccess) {
		addRequest.AlertOnSuccess = plan.AlertOnSuccess.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnFailure) {
		addRequest.AlertOnFailure = plan.AlertOnFailure.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for enter-lockdown-mode recurring-task
func addOptionalEnterLockdownModeRecurringTaskFields(ctx context.Context, addRequest *client.AddEnterLockdownModeRecurringTaskRequest, plan recurringTaskResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Reason) {
		addRequest.Reason = plan.Reason.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CancelOnTaskDependencyFailure) {
		addRequest.CancelOnTaskDependencyFailure = plan.CancelOnTaskDependencyFailure.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EmailOnStart) {
		var slice []string
		plan.EmailOnStart.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnStart = slice
	}
	if internaltypes.IsDefined(plan.EmailOnSuccess) {
		var slice []string
		plan.EmailOnSuccess.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnSuccess = slice
	}
	if internaltypes.IsDefined(plan.EmailOnFailure) {
		var slice []string
		plan.EmailOnFailure.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnFailure = slice
	}
	if internaltypes.IsDefined(plan.AlertOnStart) {
		addRequest.AlertOnStart = plan.AlertOnStart.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnSuccess) {
		addRequest.AlertOnSuccess = plan.AlertOnSuccess.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnFailure) {
		addRequest.AlertOnFailure = plan.AlertOnFailure.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for audit-data-security recurring-task
func addOptionalAuditDataSecurityRecurringTaskFields(ctx context.Context, addRequest *client.AddAuditDataSecurityRecurringTaskRequest, plan recurringTaskResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BaseOutputDirectory) {
		addRequest.BaseOutputDirectory = plan.BaseOutputDirectory.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.DataSecurityAuditor) {
		var slice []string
		plan.DataSecurityAuditor.ElementsAs(ctx, &slice, false)
		addRequest.DataSecurityAuditor = slice
	}
	if internaltypes.IsDefined(plan.Backend) {
		var slice []string
		plan.Backend.ElementsAs(ctx, &slice, false)
		addRequest.Backend = slice
	}
	if internaltypes.IsDefined(plan.IncludeFilter) {
		var slice []string
		plan.IncludeFilter.ElementsAs(ctx, &slice, false)
		addRequest.IncludeFilter = slice
	}
	if internaltypes.IsDefined(plan.RetainPreviousReportCount) {
		addRequest.RetainPreviousReportCount = plan.RetainPreviousReportCount.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RetainPreviousReportAge) {
		addRequest.RetainPreviousReportAge = plan.RetainPreviousReportAge.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CancelOnTaskDependencyFailure) {
		addRequest.CancelOnTaskDependencyFailure = plan.CancelOnTaskDependencyFailure.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EmailOnStart) {
		var slice []string
		plan.EmailOnStart.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnStart = slice
	}
	if internaltypes.IsDefined(plan.EmailOnSuccess) {
		var slice []string
		plan.EmailOnSuccess.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnSuccess = slice
	}
	if internaltypes.IsDefined(plan.EmailOnFailure) {
		var slice []string
		plan.EmailOnFailure.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnFailure = slice
	}
	if internaltypes.IsDefined(plan.AlertOnStart) {
		addRequest.AlertOnStart = plan.AlertOnStart.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnSuccess) {
		addRequest.AlertOnSuccess = plan.AlertOnSuccess.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnFailure) {
		addRequest.AlertOnFailure = plan.AlertOnFailure.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for exec recurring-task
func addOptionalExecRecurringTaskFields(ctx context.Context, addRequest *client.AddExecRecurringTaskRequest, plan recurringTaskResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CommandArguments) {
		addRequest.CommandArguments = plan.CommandArguments.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CommandOutputFileBaseName) {
		addRequest.CommandOutputFileBaseName = plan.CommandOutputFileBaseName.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RetainPreviousOutputFileCount) {
		addRequest.RetainPreviousOutputFileCount = plan.RetainPreviousOutputFileCount.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RetainPreviousOutputFileAge) {
		addRequest.RetainPreviousOutputFileAge = plan.RetainPreviousOutputFileAge.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.LogCommandOutput) {
		addRequest.LogCommandOutput = plan.LogCommandOutput.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TaskCompletionStateForNonzeroExitCode) {
		taskCompletionStateForNonzeroExitCode, err := client.NewEnumrecurringTaskTaskCompletionStateForNonzeroExitCodePropFromValue(plan.TaskCompletionStateForNonzeroExitCode.ValueString())
		if err != nil {
			return err
		}
		addRequest.TaskCompletionStateForNonzeroExitCode = taskCompletionStateForNonzeroExitCode
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.WorkingDirectory) {
		addRequest.WorkingDirectory = plan.WorkingDirectory.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CancelOnTaskDependencyFailure) {
		addRequest.CancelOnTaskDependencyFailure = plan.CancelOnTaskDependencyFailure.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EmailOnStart) {
		var slice []string
		plan.EmailOnStart.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnStart = slice
	}
	if internaltypes.IsDefined(plan.EmailOnSuccess) {
		var slice []string
		plan.EmailOnSuccess.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnSuccess = slice
	}
	if internaltypes.IsDefined(plan.EmailOnFailure) {
		var slice []string
		plan.EmailOnFailure.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnFailure = slice
	}
	if internaltypes.IsDefined(plan.AlertOnStart) {
		addRequest.AlertOnStart = plan.AlertOnStart.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnSuccess) {
		addRequest.AlertOnSuccess = plan.AlertOnSuccess.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnFailure) {
		addRequest.AlertOnFailure = plan.AlertOnFailure.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for file-retention recurring-task
func addOptionalFileRetentionRecurringTaskFields(ctx context.Context, addRequest *client.AddFileRetentionRecurringTaskRequest, plan recurringTaskResourceModel) error {
	if internaltypes.IsDefined(plan.RetainFileCount) {
		addRequest.RetainFileCount = plan.RetainFileCount.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RetainFileAge) {
		addRequest.RetainFileAge = plan.RetainFileAge.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RetainAggregateFileSize) {
		addRequest.RetainAggregateFileSize = plan.RetainAggregateFileSize.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CancelOnTaskDependencyFailure) {
		addRequest.CancelOnTaskDependencyFailure = plan.CancelOnTaskDependencyFailure.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EmailOnStart) {
		var slice []string
		plan.EmailOnStart.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnStart = slice
	}
	if internaltypes.IsDefined(plan.EmailOnSuccess) {
		var slice []string
		plan.EmailOnSuccess.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnSuccess = slice
	}
	if internaltypes.IsDefined(plan.EmailOnFailure) {
		var slice []string
		plan.EmailOnFailure.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnFailure = slice
	}
	if internaltypes.IsDefined(plan.AlertOnStart) {
		addRequest.AlertOnStart = plan.AlertOnStart.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnSuccess) {
		addRequest.AlertOnSuccess = plan.AlertOnSuccess.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnFailure) {
		addRequest.AlertOnFailure = plan.AlertOnFailure.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for third-party recurring-task
func addOptionalThirdPartyRecurringTaskFields(ctx context.Context, addRequest *client.AddThirdPartyRecurringTaskRequest, plan recurringTaskResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CancelOnTaskDependencyFailure) {
		addRequest.CancelOnTaskDependencyFailure = plan.CancelOnTaskDependencyFailure.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EmailOnStart) {
		var slice []string
		plan.EmailOnStart.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnStart = slice
	}
	if internaltypes.IsDefined(plan.EmailOnSuccess) {
		var slice []string
		plan.EmailOnSuccess.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnSuccess = slice
	}
	if internaltypes.IsDefined(plan.EmailOnFailure) {
		var slice []string
		plan.EmailOnFailure.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnFailure = slice
	}
	if internaltypes.IsDefined(plan.AlertOnStart) {
		addRequest.AlertOnStart = plan.AlertOnStart.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnSuccess) {
		addRequest.AlertOnSuccess = plan.AlertOnSuccess.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnFailure) {
		addRequest.AlertOnFailure = plan.AlertOnFailure.ValueBoolPointer()
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateRecurringTaskUnknownValues(model *recurringTaskResourceModel) {
	if model.IncludedBackendID.IsUnknown() || model.IncludedBackendID.IsNull() {
		model.IncludedBackendID, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.Backend.IsUnknown() || model.Backend.IsNull() {
		model.Backend, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludePath.IsUnknown() || model.IncludePath.IsNull() {
		model.IncludePath, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.LdapURLForSearchExpectedToReturnEntries.IsUnknown() || model.LdapURLForSearchExpectedToReturnEntries.IsNull() {
		model.LdapURLForSearchExpectedToReturnEntries, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.TaskAttributeValue.IsUnknown() || model.TaskAttributeValue.IsNull() {
		model.TaskAttributeValue, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.BackendID.IsUnknown() || model.BackendID.IsNull() {
		model.BackendID, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.DataSecurityAuditor.IsUnknown() || model.DataSecurityAuditor.IsNull() {
		model.DataSecurityAuditor, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExcludeBackendID.IsUnknown() || model.ExcludeBackendID.IsNull() {
		model.ExcludeBackendID, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExcludedBackendID.IsUnknown() || model.ExcludedBackendID.IsNull() {
		model.ExcludedBackendID, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.TaskObjectClass.IsUnknown() || model.TaskObjectClass.IsNull() {
		model.TaskObjectClass, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludeFilter.IsUnknown() || model.IncludeFilter.IsNull() {
		model.IncludeFilter, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.BaseOutputDirectory.IsUnknown() || model.BaseOutputDirectory.IsNull() {
		model.BaseOutputDirectory = types.StringValue("")
	}
	if model.BackupDirectory.IsUnknown() || model.BackupDirectory.IsNull() {
		model.BackupDirectory = types.StringValue("")
	}
	if model.LdifDirectory.IsUnknown() || model.LdifDirectory.IsNull() {
		model.LdifDirectory = types.StringValue("")
	}
	if model.SecurityLevel.IsUnknown() || model.SecurityLevel.IsNull() {
		model.SecurityLevel = types.StringValue("")
	}
	if model.TaskCompletionStateForNonzeroExitCode.IsUnknown() || model.TaskCompletionStateForNonzeroExitCode.IsNull() {
		model.TaskCompletionStateForNonzeroExitCode = types.StringValue("")
	}
	if model.TaskReturnStateIfTimeoutIsEncountered.IsUnknown() || model.TaskReturnStateIfTimeoutIsEncountered.IsNull() {
		model.TaskReturnStateIfTimeoutIsEncountered = types.StringValue("")
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *recurringTaskResourceModel) populateAllComputedStringAttributes() {
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.LogFileHeadCollectionSize.IsUnknown() || model.LogFileHeadCollectionSize.IsNull() {
		model.LogFileHeadCollectionSize = types.StringValue("")
	}
	if model.OutputDirectory.IsUnknown() || model.OutputDirectory.IsNull() {
		model.OutputDirectory = types.StringValue("")
	}
	if model.RetainAggregateFileSize.IsUnknown() || model.RetainAggregateFileSize.IsNull() {
		model.RetainAggregateFileSize = types.StringValue("")
	}
	if model.CommandArguments.IsUnknown() || model.CommandArguments.IsNull() {
		model.CommandArguments = types.StringValue("")
	}
	if model.EncryptionSettingsDefinitionID.IsUnknown() || model.EncryptionSettingsDefinitionID.IsNull() {
		model.EncryptionSettingsDefinitionID = types.StringValue("")
	}
	if model.DurationToWaitForWorkQueueIdle.IsUnknown() || model.DurationToWaitForWorkQueueIdle.IsNull() {
		model.DurationToWaitForWorkQueueIdle = types.StringValue("")
	}
	if model.Reason.IsUnknown() || model.Reason.IsNull() {
		model.Reason = types.StringValue("")
	}
	if model.DurationToWaitForSearchToReturnEntries.IsUnknown() || model.DurationToWaitForSearchToReturnEntries.IsNull() {
		model.DurationToWaitForSearchToReturnEntries = types.StringValue("")
	}
	if model.TaskJavaClass.IsUnknown() || model.TaskJavaClass.IsNull() {
		model.TaskJavaClass = types.StringValue("")
	}
	if model.TargetDirectory.IsUnknown() || model.TargetDirectory.IsNull() {
		model.TargetDirectory = types.StringValue("")
	}
	if model.RetainPreviousReportAge.IsUnknown() || model.RetainPreviousReportAge.IsNull() {
		model.RetainPreviousReportAge = types.StringValue("")
	}
	if model.SleepDuration.IsUnknown() || model.SleepDuration.IsNull() {
		model.SleepDuration = types.StringValue("")
	}
	if model.FilenamePattern.IsUnknown() || model.FilenamePattern.IsNull() {
		model.FilenamePattern = types.StringValue("")
	}
	if model.RetainPreviousProfileAge.IsUnknown() || model.RetainPreviousProfileAge.IsNull() {
		model.RetainPreviousProfileAge = types.StringValue("")
	}
	if model.RetainPreviousOutputFileAge.IsUnknown() || model.RetainPreviousOutputFileAge.IsNull() {
		model.RetainPreviousOutputFileAge = types.StringValue("")
	}
	if model.RetainPreviousSupportDataArchiveAge.IsUnknown() || model.RetainPreviousSupportDataArchiveAge.IsNull() {
		model.RetainPreviousSupportDataArchiveAge = types.StringValue("")
	}
	if model.Comment.IsUnknown() || model.Comment.IsNull() {
		model.Comment = types.StringValue("")
	}
	if model.WorkingDirectory.IsUnknown() || model.WorkingDirectory.IsNull() {
		model.WorkingDirectory = types.StringValue("")
	}
	if model.SearchInterval.IsUnknown() || model.SearchInterval.IsNull() {
		model.SearchInterval = types.StringValue("")
	}
	if model.ExtensionClass.IsUnknown() || model.ExtensionClass.IsNull() {
		model.ExtensionClass = types.StringValue("")
	}
	if model.ProfileDirectory.IsUnknown() || model.ProfileDirectory.IsNull() {
		model.ProfileDirectory = types.StringValue("")
	}
	if model.RetainFileAge.IsUnknown() || model.RetainFileAge.IsNull() {
		model.RetainFileAge = types.StringValue("")
	}
	if model.EncryptionPassphraseFile.IsUnknown() || model.EncryptionPassphraseFile.IsNull() {
		model.EncryptionPassphraseFile = types.StringValue("")
	}
	if model.TimestampFormat.IsUnknown() || model.TimestampFormat.IsNull() {
		model.TimestampFormat = types.StringValue("")
	}
	if model.RetainPreviousFullBackupAge.IsUnknown() || model.RetainPreviousFullBackupAge.IsNull() {
		model.RetainPreviousFullBackupAge = types.StringValue("")
	}
	if model.SearchTimeLimit.IsUnknown() || model.SearchTimeLimit.IsNull() {
		model.SearchTimeLimit = types.StringValue("")
	}
	if model.LogFileTailCollectionSize.IsUnknown() || model.LogFileTailCollectionSize.IsNull() {
		model.LogFileTailCollectionSize = types.StringValue("")
	}
	if model.CommandOutputFileBaseName.IsUnknown() || model.CommandOutputFileBaseName.IsNull() {
		model.CommandOutputFileBaseName = types.StringValue("")
	}
	if model.LogDuration.IsUnknown() || model.LogDuration.IsNull() {
		model.LogDuration = types.StringValue("")
	}
	if model.RetainPreviousLDIFExportAge.IsUnknown() || model.RetainPreviousLDIFExportAge.IsNull() {
		model.RetainPreviousLDIFExportAge = types.StringValue("")
	}
	if model.CommandPath.IsUnknown() || model.CommandPath.IsNull() {
		model.CommandPath = types.StringValue("")
	}
}

// Read a GenerateServerProfileRecurringTaskResponse object into the model struct
func readGenerateServerProfileRecurringTaskResponse(ctx context.Context, r *client.GenerateServerProfileRecurringTaskResponse, state *recurringTaskResourceModel, expectedValues *recurringTaskResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generate-server-profile")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ProfileDirectory = types.StringValue(r.ProfileDirectory)
	state.IncludePath = internaltypes.GetStringSet(r.IncludePath)
	state.RetainPreviousProfileCount = internaltypes.Int64TypeOrNil(r.RetainPreviousProfileCount)
	state.RetainPreviousProfileAge = internaltypes.StringTypeOrNil(r.RetainPreviousProfileAge, internaltypes.IsEmptyString(expectedValues.RetainPreviousProfileAge))
	config.CheckMismatchedPDFormattedAttributes("retain_previous_profile_age",
		expectedValues.RetainPreviousProfileAge, state.RetainPreviousProfileAge, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateRecurringTaskUnknownValues(state)
}

// Read a LeaveLockdownModeRecurringTaskResponse object into the model struct
func readLeaveLockdownModeRecurringTaskResponse(ctx context.Context, r *client.LeaveLockdownModeRecurringTaskResponse, state *recurringTaskResourceModel, expectedValues *recurringTaskResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("leave-lockdown-mode")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Reason = internaltypes.StringTypeOrNil(r.Reason, internaltypes.IsEmptyString(expectedValues.Reason))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateRecurringTaskUnknownValues(state)
}

// Read a BackupRecurringTaskResponse object into the model struct
func readBackupRecurringTaskResponse(ctx context.Context, r *client.BackupRecurringTaskResponse, state *recurringTaskResourceModel, expectedValues *recurringTaskResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("backup")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BackupDirectory = types.StringValue(r.BackupDirectory)
	state.IncludedBackendID = internaltypes.GetStringSet(r.IncludedBackendID)
	state.ExcludedBackendID = internaltypes.GetStringSet(r.ExcludedBackendID)
	state.Compress = internaltypes.BoolTypeOrNil(r.Compress)
	state.Encrypt = internaltypes.BoolTypeOrNil(r.Encrypt)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Sign = internaltypes.BoolTypeOrNil(r.Sign)
	state.RetainPreviousFullBackupCount = internaltypes.Int64TypeOrNil(r.RetainPreviousFullBackupCount)
	state.RetainPreviousFullBackupAge = internaltypes.StringTypeOrNil(r.RetainPreviousFullBackupAge, internaltypes.IsEmptyString(expectedValues.RetainPreviousFullBackupAge))
	config.CheckMismatchedPDFormattedAttributes("retain_previous_full_backup_age",
		expectedValues.RetainPreviousFullBackupAge, state.RetainPreviousFullBackupAge, diagnostics)
	state.MaxMegabytesPerSecond = internaltypes.Int64TypeOrNil(r.MaxMegabytesPerSecond)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateRecurringTaskUnknownValues(state)
}

// Read a DelayRecurringTaskResponse object into the model struct
func readDelayRecurringTaskResponse(ctx context.Context, r *client.DelayRecurringTaskResponse, state *recurringTaskResourceModel, expectedValues *recurringTaskResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("delay")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SleepDuration = internaltypes.StringTypeOrNil(r.SleepDuration, internaltypes.IsEmptyString(expectedValues.SleepDuration))
	config.CheckMismatchedPDFormattedAttributes("sleep_duration",
		expectedValues.SleepDuration, state.SleepDuration, diagnostics)
	state.DurationToWaitForWorkQueueIdle = internaltypes.StringTypeOrNil(r.DurationToWaitForWorkQueueIdle, internaltypes.IsEmptyString(expectedValues.DurationToWaitForWorkQueueIdle))
	config.CheckMismatchedPDFormattedAttributes("duration_to_wait_for_work_queue_idle",
		expectedValues.DurationToWaitForWorkQueueIdle, state.DurationToWaitForWorkQueueIdle, diagnostics)
	state.LdapURLForSearchExpectedToReturnEntries = internaltypes.GetStringSet(r.LdapURLForSearchExpectedToReturnEntries)
	state.SearchInterval = internaltypes.StringTypeOrNil(r.SearchInterval, internaltypes.IsEmptyString(expectedValues.SearchInterval))
	config.CheckMismatchedPDFormattedAttributes("search_interval",
		expectedValues.SearchInterval, state.SearchInterval, diagnostics)
	state.SearchTimeLimit = internaltypes.StringTypeOrNil(r.SearchTimeLimit, internaltypes.IsEmptyString(expectedValues.SearchTimeLimit))
	config.CheckMismatchedPDFormattedAttributes("search_time_limit",
		expectedValues.SearchTimeLimit, state.SearchTimeLimit, diagnostics)
	state.DurationToWaitForSearchToReturnEntries = internaltypes.StringTypeOrNil(r.DurationToWaitForSearchToReturnEntries, internaltypes.IsEmptyString(expectedValues.DurationToWaitForSearchToReturnEntries))
	config.CheckMismatchedPDFormattedAttributes("duration_to_wait_for_search_to_return_entries",
		expectedValues.DurationToWaitForSearchToReturnEntries, state.DurationToWaitForSearchToReturnEntries, diagnostics)
	state.TaskReturnStateIfTimeoutIsEncountered = internaltypes.StringTypeOrNil(
		client.StringPointerEnumrecurringTaskTaskReturnStateIfTimeoutIsEncounteredProp(r.TaskReturnStateIfTimeoutIsEncountered), true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateRecurringTaskUnknownValues(state)
}

// Read a StaticallyDefinedRecurringTaskResponse object into the model struct
func readStaticallyDefinedRecurringTaskResponse(ctx context.Context, r *client.StaticallyDefinedRecurringTaskResponse, state *recurringTaskResourceModel, expectedValues *recurringTaskResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("statically-defined")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.TaskJavaClass = types.StringValue(r.TaskJavaClass)
	state.TaskObjectClass = internaltypes.GetStringSet(r.TaskObjectClass)
	state.TaskAttributeValue = internaltypes.GetStringSet(r.TaskAttributeValue)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateRecurringTaskUnknownValues(state)
}

// Read a CollectSupportDataRecurringTaskResponse object into the model struct
func readCollectSupportDataRecurringTaskResponse(ctx context.Context, r *client.CollectSupportDataRecurringTaskResponse, state *recurringTaskResourceModel, expectedValues *recurringTaskResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("collect-support-data")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.OutputDirectory = types.StringValue(r.OutputDirectory)
	state.EncryptionPassphraseFile = internaltypes.StringTypeOrNil(r.EncryptionPassphraseFile, internaltypes.IsEmptyString(expectedValues.EncryptionPassphraseFile))
	state.IncludeExpensiveData = internaltypes.BoolTypeOrNil(r.IncludeExpensiveData)
	state.IncludeReplicationStateDump = internaltypes.BoolTypeOrNil(r.IncludeReplicationStateDump)
	state.IncludeBinaryFiles = internaltypes.BoolTypeOrNil(r.IncludeBinaryFiles)
	state.IncludeExtensionSource = internaltypes.BoolTypeOrNil(r.IncludeExtensionSource)
	state.UseSequentialMode = internaltypes.BoolTypeOrNil(r.UseSequentialMode)
	state.SecurityLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumrecurringTaskSecurityLevelProp(r.SecurityLevel), true)
	state.JstackCount = internaltypes.Int64TypeOrNil(r.JstackCount)
	state.ReportCount = internaltypes.Int64TypeOrNil(r.ReportCount)
	state.ReportIntervalSeconds = internaltypes.Int64TypeOrNil(r.ReportIntervalSeconds)
	state.LogDuration = internaltypes.StringTypeOrNil(r.LogDuration, internaltypes.IsEmptyString(expectedValues.LogDuration))
	config.CheckMismatchedPDFormattedAttributes("log_duration",
		expectedValues.LogDuration, state.LogDuration, diagnostics)
	state.LogFileHeadCollectionSize = internaltypes.StringTypeOrNil(r.LogFileHeadCollectionSize, internaltypes.IsEmptyString(expectedValues.LogFileHeadCollectionSize))
	config.CheckMismatchedPDFormattedAttributes("log_file_head_collection_size",
		expectedValues.LogFileHeadCollectionSize, state.LogFileHeadCollectionSize, diagnostics)
	state.LogFileTailCollectionSize = internaltypes.StringTypeOrNil(r.LogFileTailCollectionSize, internaltypes.IsEmptyString(expectedValues.LogFileTailCollectionSize))
	config.CheckMismatchedPDFormattedAttributes("log_file_tail_collection_size",
		expectedValues.LogFileTailCollectionSize, state.LogFileTailCollectionSize, diagnostics)
	state.Comment = internaltypes.StringTypeOrNil(r.Comment, internaltypes.IsEmptyString(expectedValues.Comment))
	state.RetainPreviousSupportDataArchiveCount = internaltypes.Int64TypeOrNil(r.RetainPreviousSupportDataArchiveCount)
	state.RetainPreviousSupportDataArchiveAge = internaltypes.StringTypeOrNil(r.RetainPreviousSupportDataArchiveAge, internaltypes.IsEmptyString(expectedValues.RetainPreviousSupportDataArchiveAge))
	config.CheckMismatchedPDFormattedAttributes("retain_previous_support_data_archive_age",
		expectedValues.RetainPreviousSupportDataArchiveAge, state.RetainPreviousSupportDataArchiveAge, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateRecurringTaskUnknownValues(state)
}

// Read a LdifExportRecurringTaskResponse object into the model struct
func readLdifExportRecurringTaskResponse(ctx context.Context, r *client.LdifExportRecurringTaskResponse, state *recurringTaskResourceModel, expectedValues *recurringTaskResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldif-export")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LdifDirectory = types.StringValue(r.LdifDirectory)
	state.BackendID = internaltypes.GetStringSet(r.BackendID)
	state.ExcludeBackendID = internaltypes.GetStringSet(r.ExcludeBackendID)
	state.Compress = internaltypes.BoolTypeOrNil(r.Compress)
	state.Encrypt = internaltypes.BoolTypeOrNil(r.Encrypt)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Sign = internaltypes.BoolTypeOrNil(r.Sign)
	state.RetainPreviousLDIFExportCount = internaltypes.Int64TypeOrNil(r.RetainPreviousLDIFExportCount)
	state.RetainPreviousLDIFExportAge = internaltypes.StringTypeOrNil(r.RetainPreviousLDIFExportAge, internaltypes.IsEmptyString(expectedValues.RetainPreviousLDIFExportAge))
	config.CheckMismatchedPDFormattedAttributes("retain_previous_ldif_export_age",
		expectedValues.RetainPreviousLDIFExportAge, state.RetainPreviousLDIFExportAge, diagnostics)
	state.MaxMegabytesPerSecond = internaltypes.Int64TypeOrNil(r.MaxMegabytesPerSecond)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateRecurringTaskUnknownValues(state)
}

// Read a EnterLockdownModeRecurringTaskResponse object into the model struct
func readEnterLockdownModeRecurringTaskResponse(ctx context.Context, r *client.EnterLockdownModeRecurringTaskResponse, state *recurringTaskResourceModel, expectedValues *recurringTaskResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("enter-lockdown-mode")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Reason = internaltypes.StringTypeOrNil(r.Reason, internaltypes.IsEmptyString(expectedValues.Reason))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateRecurringTaskUnknownValues(state)
}

// Read a AuditDataSecurityRecurringTaskResponse object into the model struct
func readAuditDataSecurityRecurringTaskResponse(ctx context.Context, r *client.AuditDataSecurityRecurringTaskResponse, state *recurringTaskResourceModel, expectedValues *recurringTaskResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("audit-data-security")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BaseOutputDirectory = types.StringValue(r.BaseOutputDirectory)
	state.DataSecurityAuditor = internaltypes.GetStringSet(r.DataSecurityAuditor)
	state.Backend = internaltypes.GetStringSet(r.Backend)
	state.IncludeFilter = internaltypes.GetStringSet(r.IncludeFilter)
	state.RetainPreviousReportCount = internaltypes.Int64TypeOrNil(r.RetainPreviousReportCount)
	state.RetainPreviousReportAge = internaltypes.StringTypeOrNil(r.RetainPreviousReportAge, internaltypes.IsEmptyString(expectedValues.RetainPreviousReportAge))
	config.CheckMismatchedPDFormattedAttributes("retain_previous_report_age",
		expectedValues.RetainPreviousReportAge, state.RetainPreviousReportAge, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateRecurringTaskUnknownValues(state)
}

// Read a ExecRecurringTaskResponse object into the model struct
func readExecRecurringTaskResponse(ctx context.Context, r *client.ExecRecurringTaskResponse, state *recurringTaskResourceModel, expectedValues *recurringTaskResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("exec")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.CommandPath = types.StringValue(r.CommandPath)
	state.CommandArguments = internaltypes.StringTypeOrNil(r.CommandArguments, internaltypes.IsEmptyString(expectedValues.CommandArguments))
	state.CommandOutputFileBaseName = internaltypes.StringTypeOrNil(r.CommandOutputFileBaseName, internaltypes.IsEmptyString(expectedValues.CommandOutputFileBaseName))
	state.RetainPreviousOutputFileCount = internaltypes.Int64TypeOrNil(r.RetainPreviousOutputFileCount)
	state.RetainPreviousOutputFileAge = internaltypes.StringTypeOrNil(r.RetainPreviousOutputFileAge, internaltypes.IsEmptyString(expectedValues.RetainPreviousOutputFileAge))
	config.CheckMismatchedPDFormattedAttributes("retain_previous_output_file_age",
		expectedValues.RetainPreviousOutputFileAge, state.RetainPreviousOutputFileAge, diagnostics)
	state.LogCommandOutput = internaltypes.BoolTypeOrNil(r.LogCommandOutput)
	state.TaskCompletionStateForNonzeroExitCode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumrecurringTaskTaskCompletionStateForNonzeroExitCodeProp(r.TaskCompletionStateForNonzeroExitCode), true)
	state.WorkingDirectory = internaltypes.StringTypeOrNil(r.WorkingDirectory, internaltypes.IsEmptyString(expectedValues.WorkingDirectory))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateRecurringTaskUnknownValues(state)
}

// Read a FileRetentionRecurringTaskResponse object into the model struct
func readFileRetentionRecurringTaskResponse(ctx context.Context, r *client.FileRetentionRecurringTaskResponse, state *recurringTaskResourceModel, expectedValues *recurringTaskResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-retention")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.TargetDirectory = types.StringValue(r.TargetDirectory)
	state.FilenamePattern = types.StringValue(r.FilenamePattern)
	state.TimestampFormat = types.StringValue(r.TimestampFormat.String())
	state.RetainFileCount = internaltypes.Int64TypeOrNil(r.RetainFileCount)
	state.RetainFileAge = internaltypes.StringTypeOrNil(r.RetainFileAge, internaltypes.IsEmptyString(expectedValues.RetainFileAge))
	config.CheckMismatchedPDFormattedAttributes("retain_file_age",
		expectedValues.RetainFileAge, state.RetainFileAge, diagnostics)
	state.RetainAggregateFileSize = internaltypes.StringTypeOrNil(r.RetainAggregateFileSize, internaltypes.IsEmptyString(expectedValues.RetainAggregateFileSize))
	config.CheckMismatchedPDFormattedAttributes("retain_aggregate_file_size",
		expectedValues.RetainAggregateFileSize, state.RetainAggregateFileSize, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateRecurringTaskUnknownValues(state)
}

// Read a ThirdPartyRecurringTaskResponse object into the model struct
func readThirdPartyRecurringTaskResponse(ctx context.Context, r *client.ThirdPartyRecurringTaskResponse, state *recurringTaskResourceModel, expectedValues *recurringTaskResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateRecurringTaskUnknownValues(state)
}

// Create any update operations necessary to make the state match the plan
func createRecurringTaskOperations(plan recurringTaskResourceModel, state recurringTaskResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.TargetDirectory, state.TargetDirectory, "target-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.FilenamePattern, state.FilenamePattern, "filename-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.TimestampFormat, state.TimestampFormat, "timestamp-format")
	operations.AddInt64OperationIfNecessary(&ops, plan.RetainFileCount, state.RetainFileCount, "retain-file-count")
	operations.AddStringOperationIfNecessary(&ops, plan.RetainFileAge, state.RetainFileAge, "retain-file-age")
	operations.AddStringOperationIfNecessary(&ops, plan.RetainAggregateFileSize, state.RetainAggregateFileSize, "retain-aggregate-file-size")
	operations.AddStringOperationIfNecessary(&ops, plan.CommandPath, state.CommandPath, "command-path")
	operations.AddStringOperationIfNecessary(&ops, plan.CommandArguments, state.CommandArguments, "command-arguments")
	operations.AddStringOperationIfNecessary(&ops, plan.CommandOutputFileBaseName, state.CommandOutputFileBaseName, "command-output-file-base-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.RetainPreviousOutputFileCount, state.RetainPreviousOutputFileCount, "retain-previous-output-file-count")
	operations.AddStringOperationIfNecessary(&ops, plan.RetainPreviousOutputFileAge, state.RetainPreviousOutputFileAge, "retain-previous-output-file-age")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogCommandOutput, state.LogCommandOutput, "log-command-output")
	operations.AddStringOperationIfNecessary(&ops, plan.TaskCompletionStateForNonzeroExitCode, state.TaskCompletionStateForNonzeroExitCode, "task-completion-state-for-nonzero-exit-code")
	operations.AddStringOperationIfNecessary(&ops, plan.WorkingDirectory, state.WorkingDirectory, "working-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.BaseOutputDirectory, state.BaseOutputDirectory, "base-output-directory")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DataSecurityAuditor, state.DataSecurityAuditor, "data-security-auditor")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Backend, state.Backend, "backend")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeFilter, state.IncludeFilter, "include-filter")
	operations.AddInt64OperationIfNecessary(&ops, plan.RetainPreviousReportCount, state.RetainPreviousReportCount, "retain-previous-report-count")
	operations.AddStringOperationIfNecessary(&ops, plan.RetainPreviousReportAge, state.RetainPreviousReportAge, "retain-previous-report-age")
	operations.AddStringOperationIfNecessary(&ops, plan.LdifDirectory, state.LdifDirectory, "ldif-directory")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.BackendID, state.BackendID, "backend-id")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeBackendID, state.ExcludeBackendID, "exclude-backend-id")
	operations.AddStringOperationIfNecessary(&ops, plan.OutputDirectory, state.OutputDirectory, "output-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptionPassphraseFile, state.EncryptionPassphraseFile, "encryption-passphrase-file")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeExpensiveData, state.IncludeExpensiveData, "include-expensive-data")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeReplicationStateDump, state.IncludeReplicationStateDump, "include-replication-state-dump")
	operations.AddInt64OperationIfNecessary(&ops, plan.RetainPreviousLDIFExportCount, state.RetainPreviousLDIFExportCount, "retain-previous-ldif-export-count")
	operations.AddStringOperationIfNecessary(&ops, plan.RetainPreviousLDIFExportAge, state.RetainPreviousLDIFExportAge, "retain-previous-ldif-export-age")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeBinaryFiles, state.IncludeBinaryFiles, "include-binary-files")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeExtensionSource, state.IncludeExtensionSource, "include-extension-source")
	operations.AddBoolOperationIfNecessary(&ops, plan.UseSequentialMode, state.UseSequentialMode, "use-sequential-mode")
	operations.AddStringOperationIfNecessary(&ops, plan.SecurityLevel, state.SecurityLevel, "security-level")
	operations.AddInt64OperationIfNecessary(&ops, plan.JstackCount, state.JstackCount, "jstack-count")
	operations.AddInt64OperationIfNecessary(&ops, plan.ReportCount, state.ReportCount, "report-count")
	operations.AddInt64OperationIfNecessary(&ops, plan.ReportIntervalSeconds, state.ReportIntervalSeconds, "report-interval-seconds")
	operations.AddStringOperationIfNecessary(&ops, plan.LogDuration, state.LogDuration, "log-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFileHeadCollectionSize, state.LogFileHeadCollectionSize, "log-file-head-collection-size")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFileTailCollectionSize, state.LogFileTailCollectionSize, "log-file-tail-collection-size")
	operations.AddStringOperationIfNecessary(&ops, plan.Comment, state.Comment, "comment")
	operations.AddInt64OperationIfNecessary(&ops, plan.RetainPreviousSupportDataArchiveCount, state.RetainPreviousSupportDataArchiveCount, "retain-previous-support-data-archive-count")
	operations.AddStringOperationIfNecessary(&ops, plan.RetainPreviousSupportDataArchiveAge, state.RetainPreviousSupportDataArchiveAge, "retain-previous-support-data-archive-age")
	operations.AddStringOperationIfNecessary(&ops, plan.TaskJavaClass, state.TaskJavaClass, "task-java-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.TaskObjectClass, state.TaskObjectClass, "task-object-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.TaskAttributeValue, state.TaskAttributeValue, "task-attribute-value")
	operations.AddStringOperationIfNecessary(&ops, plan.SleepDuration, state.SleepDuration, "sleep-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.DurationToWaitForWorkQueueIdle, state.DurationToWaitForWorkQueueIdle, "duration-to-wait-for-work-queue-idle")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.LdapURLForSearchExpectedToReturnEntries, state.LdapURLForSearchExpectedToReturnEntries, "ldap-url-for-search-expected-to-return-entries")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchInterval, state.SearchInterval, "search-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchTimeLimit, state.SearchTimeLimit, "search-time-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.DurationToWaitForSearchToReturnEntries, state.DurationToWaitForSearchToReturnEntries, "duration-to-wait-for-search-to-return-entries")
	operations.AddStringOperationIfNecessary(&ops, plan.TaskReturnStateIfTimeoutIsEncountered, state.TaskReturnStateIfTimeoutIsEncountered, "task-return-state-if-timeout-is-encountered")
	operations.AddStringOperationIfNecessary(&ops, plan.BackupDirectory, state.BackupDirectory, "backup-directory")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedBackendID, state.IncludedBackendID, "included-backend-id")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedBackendID, state.ExcludedBackendID, "excluded-backend-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.Compress, state.Compress, "compress")
	operations.AddBoolOperationIfNecessary(&ops, plan.Encrypt, state.Encrypt, "encrypt")
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptionSettingsDefinitionID, state.EncryptionSettingsDefinitionID, "encryption-settings-definition-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.Sign, state.Sign, "sign")
	operations.AddInt64OperationIfNecessary(&ops, plan.RetainPreviousFullBackupCount, state.RetainPreviousFullBackupCount, "retain-previous-full-backup-count")
	operations.AddStringOperationIfNecessary(&ops, plan.RetainPreviousFullBackupAge, state.RetainPreviousFullBackupAge, "retain-previous-full-backup-age")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxMegabytesPerSecond, state.MaxMegabytesPerSecond, "max-megabytes-per-second")
	operations.AddStringOperationIfNecessary(&ops, plan.Reason, state.Reason, "reason")
	operations.AddStringOperationIfNecessary(&ops, plan.ProfileDirectory, state.ProfileDirectory, "profile-directory")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludePath, state.IncludePath, "include-path")
	operations.AddInt64OperationIfNecessary(&ops, plan.RetainPreviousProfileCount, state.RetainPreviousProfileCount, "retain-previous-profile-count")
	operations.AddStringOperationIfNecessary(&ops, plan.RetainPreviousProfileAge, state.RetainPreviousProfileAge, "retain-previous-profile-age")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.CancelOnTaskDependencyFailure, state.CancelOnTaskDependencyFailure, "cancel-on-task-dependency-failure")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EmailOnStart, state.EmailOnStart, "email-on-start")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EmailOnSuccess, state.EmailOnSuccess, "email-on-success")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EmailOnFailure, state.EmailOnFailure, "email-on-failure")
	operations.AddBoolOperationIfNecessary(&ops, plan.AlertOnStart, state.AlertOnStart, "alert-on-start")
	operations.AddBoolOperationIfNecessary(&ops, plan.AlertOnSuccess, state.AlertOnSuccess, "alert-on-success")
	operations.AddBoolOperationIfNecessary(&ops, plan.AlertOnFailure, state.AlertOnFailure, "alert-on-failure")
	return ops
}

// Create a generate-server-profile recurring-task
func (r *recurringTaskResource) CreateGenerateServerProfileRecurringTask(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan recurringTaskResourceModel) (*recurringTaskResourceModel, error) {
	addRequest := client.NewAddGenerateServerProfileRecurringTaskRequest(plan.Name.ValueString(),
		[]client.EnumgenerateServerProfileRecurringTaskSchemaUrn{client.ENUMGENERATESERVERPROFILERECURRINGTASKSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RECURRING_TASKGENERATE_SERVER_PROFILE},
		plan.ProfileDirectory.ValueString())
	err := addOptionalGenerateServerProfileRecurringTaskFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Recurring Task", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RecurringTaskApi.AddRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRecurringTaskRequest(
		client.AddGenerateServerProfileRecurringTaskRequestAsAddRecurringTaskRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RecurringTaskApi.AddRecurringTaskExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Recurring Task", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state recurringTaskResourceModel
	readGenerateServerProfileRecurringTaskResponse(ctx, addResponse.GenerateServerProfileRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a leave-lockdown-mode recurring-task
func (r *recurringTaskResource) CreateLeaveLockdownModeRecurringTask(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan recurringTaskResourceModel) (*recurringTaskResourceModel, error) {
	addRequest := client.NewAddLeaveLockdownModeRecurringTaskRequest(plan.Name.ValueString(),
		[]client.EnumleaveLockdownModeRecurringTaskSchemaUrn{client.ENUMLEAVELOCKDOWNMODERECURRINGTASKSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RECURRING_TASKLEAVE_LOCKDOWN_MODE})
	err := addOptionalLeaveLockdownModeRecurringTaskFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Recurring Task", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RecurringTaskApi.AddRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRecurringTaskRequest(
		client.AddLeaveLockdownModeRecurringTaskRequestAsAddRecurringTaskRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RecurringTaskApi.AddRecurringTaskExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Recurring Task", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state recurringTaskResourceModel
	readLeaveLockdownModeRecurringTaskResponse(ctx, addResponse.LeaveLockdownModeRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a backup recurring-task
func (r *recurringTaskResource) CreateBackupRecurringTask(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan recurringTaskResourceModel) (*recurringTaskResourceModel, error) {
	addRequest := client.NewAddBackupRecurringTaskRequest(plan.Name.ValueString(),
		[]client.EnumbackupRecurringTaskSchemaUrn{client.ENUMBACKUPRECURRINGTASKSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RECURRING_TASKBACKUP})
	err := addOptionalBackupRecurringTaskFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Recurring Task", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RecurringTaskApi.AddRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRecurringTaskRequest(
		client.AddBackupRecurringTaskRequestAsAddRecurringTaskRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RecurringTaskApi.AddRecurringTaskExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Recurring Task", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state recurringTaskResourceModel
	readBackupRecurringTaskResponse(ctx, addResponse.BackupRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a delay recurring-task
func (r *recurringTaskResource) CreateDelayRecurringTask(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan recurringTaskResourceModel) (*recurringTaskResourceModel, error) {
	addRequest := client.NewAddDelayRecurringTaskRequest(plan.Name.ValueString(),
		[]client.EnumdelayRecurringTaskSchemaUrn{client.ENUMDELAYRECURRINGTASKSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RECURRING_TASKDELAY})
	err := addOptionalDelayRecurringTaskFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Recurring Task", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RecurringTaskApi.AddRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRecurringTaskRequest(
		client.AddDelayRecurringTaskRequestAsAddRecurringTaskRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RecurringTaskApi.AddRecurringTaskExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Recurring Task", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state recurringTaskResourceModel
	readDelayRecurringTaskResponse(ctx, addResponse.DelayRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a statically-defined recurring-task
func (r *recurringTaskResource) CreateStaticallyDefinedRecurringTask(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan recurringTaskResourceModel) (*recurringTaskResourceModel, error) {
	var TaskObjectClassSlice []string
	plan.TaskObjectClass.ElementsAs(ctx, &TaskObjectClassSlice, false)
	addRequest := client.NewAddStaticallyDefinedRecurringTaskRequest(plan.Name.ValueString(),
		[]client.EnumstaticallyDefinedRecurringTaskSchemaUrn{client.ENUMSTATICALLYDEFINEDRECURRINGTASKSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RECURRING_TASKSTATICALLY_DEFINED},
		plan.TaskJavaClass.ValueString(),
		TaskObjectClassSlice)
	err := addOptionalStaticallyDefinedRecurringTaskFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Recurring Task", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RecurringTaskApi.AddRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRecurringTaskRequest(
		client.AddStaticallyDefinedRecurringTaskRequestAsAddRecurringTaskRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RecurringTaskApi.AddRecurringTaskExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Recurring Task", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state recurringTaskResourceModel
	readStaticallyDefinedRecurringTaskResponse(ctx, addResponse.StaticallyDefinedRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a collect-support-data recurring-task
func (r *recurringTaskResource) CreateCollectSupportDataRecurringTask(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan recurringTaskResourceModel) (*recurringTaskResourceModel, error) {
	addRequest := client.NewAddCollectSupportDataRecurringTaskRequest(plan.Name.ValueString(),
		[]client.EnumcollectSupportDataRecurringTaskSchemaUrn{client.ENUMCOLLECTSUPPORTDATARECURRINGTASKSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RECURRING_TASKCOLLECT_SUPPORT_DATA},
		plan.OutputDirectory.ValueString())
	err := addOptionalCollectSupportDataRecurringTaskFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Recurring Task", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RecurringTaskApi.AddRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRecurringTaskRequest(
		client.AddCollectSupportDataRecurringTaskRequestAsAddRecurringTaskRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RecurringTaskApi.AddRecurringTaskExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Recurring Task", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state recurringTaskResourceModel
	readCollectSupportDataRecurringTaskResponse(ctx, addResponse.CollectSupportDataRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a ldif-export recurring-task
func (r *recurringTaskResource) CreateLdifExportRecurringTask(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan recurringTaskResourceModel) (*recurringTaskResourceModel, error) {
	addRequest := client.NewAddLdifExportRecurringTaskRequest(plan.Name.ValueString(),
		[]client.EnumldifExportRecurringTaskSchemaUrn{client.ENUMLDIFEXPORTRECURRINGTASKSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RECURRING_TASKLDIF_EXPORT})
	err := addOptionalLdifExportRecurringTaskFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Recurring Task", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RecurringTaskApi.AddRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRecurringTaskRequest(
		client.AddLdifExportRecurringTaskRequestAsAddRecurringTaskRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RecurringTaskApi.AddRecurringTaskExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Recurring Task", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state recurringTaskResourceModel
	readLdifExportRecurringTaskResponse(ctx, addResponse.LdifExportRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a enter-lockdown-mode recurring-task
func (r *recurringTaskResource) CreateEnterLockdownModeRecurringTask(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan recurringTaskResourceModel) (*recurringTaskResourceModel, error) {
	addRequest := client.NewAddEnterLockdownModeRecurringTaskRequest(plan.Name.ValueString(),
		[]client.EnumenterLockdownModeRecurringTaskSchemaUrn{client.ENUMENTERLOCKDOWNMODERECURRINGTASKSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RECURRING_TASKENTER_LOCKDOWN_MODE})
	err := addOptionalEnterLockdownModeRecurringTaskFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Recurring Task", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RecurringTaskApi.AddRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRecurringTaskRequest(
		client.AddEnterLockdownModeRecurringTaskRequestAsAddRecurringTaskRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RecurringTaskApi.AddRecurringTaskExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Recurring Task", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state recurringTaskResourceModel
	readEnterLockdownModeRecurringTaskResponse(ctx, addResponse.EnterLockdownModeRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a audit-data-security recurring-task
func (r *recurringTaskResource) CreateAuditDataSecurityRecurringTask(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan recurringTaskResourceModel) (*recurringTaskResourceModel, error) {
	addRequest := client.NewAddAuditDataSecurityRecurringTaskRequest(plan.Name.ValueString(),
		[]client.EnumauditDataSecurityRecurringTaskSchemaUrn{client.ENUMAUDITDATASECURITYRECURRINGTASKSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RECURRING_TASKAUDIT_DATA_SECURITY})
	err := addOptionalAuditDataSecurityRecurringTaskFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Recurring Task", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RecurringTaskApi.AddRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRecurringTaskRequest(
		client.AddAuditDataSecurityRecurringTaskRequestAsAddRecurringTaskRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RecurringTaskApi.AddRecurringTaskExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Recurring Task", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state recurringTaskResourceModel
	readAuditDataSecurityRecurringTaskResponse(ctx, addResponse.AuditDataSecurityRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a exec recurring-task
func (r *recurringTaskResource) CreateExecRecurringTask(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan recurringTaskResourceModel) (*recurringTaskResourceModel, error) {
	addRequest := client.NewAddExecRecurringTaskRequest(plan.Name.ValueString(),
		[]client.EnumexecRecurringTaskSchemaUrn{client.ENUMEXECRECURRINGTASKSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RECURRING_TASKEXEC},
		plan.CommandPath.ValueString())
	err := addOptionalExecRecurringTaskFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Recurring Task", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RecurringTaskApi.AddRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRecurringTaskRequest(
		client.AddExecRecurringTaskRequestAsAddRecurringTaskRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RecurringTaskApi.AddRecurringTaskExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Recurring Task", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state recurringTaskResourceModel
	readExecRecurringTaskResponse(ctx, addResponse.ExecRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a file-retention recurring-task
func (r *recurringTaskResource) CreateFileRetentionRecurringTask(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan recurringTaskResourceModel) (*recurringTaskResourceModel, error) {
	timestampFormat, err := client.NewEnumrecurringTaskTimestampFormatPropFromValue(plan.TimestampFormat.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for TimestampFormat", err.Error())
		return nil, err
	}
	addRequest := client.NewAddFileRetentionRecurringTaskRequest(plan.Name.ValueString(),
		[]client.EnumfileRetentionRecurringTaskSchemaUrn{client.ENUMFILERETENTIONRECURRINGTASKSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RECURRING_TASKFILE_RETENTION},
		plan.TargetDirectory.ValueString(),
		plan.FilenamePattern.ValueString(),
		*timestampFormat)
	err = addOptionalFileRetentionRecurringTaskFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Recurring Task", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RecurringTaskApi.AddRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRecurringTaskRequest(
		client.AddFileRetentionRecurringTaskRequestAsAddRecurringTaskRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RecurringTaskApi.AddRecurringTaskExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Recurring Task", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state recurringTaskResourceModel
	readFileRetentionRecurringTaskResponse(ctx, addResponse.FileRetentionRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party recurring-task
func (r *recurringTaskResource) CreateThirdPartyRecurringTask(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan recurringTaskResourceModel) (*recurringTaskResourceModel, error) {
	addRequest := client.NewAddThirdPartyRecurringTaskRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyRecurringTaskSchemaUrn{client.ENUMTHIRDPARTYRECURRINGTASKSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RECURRING_TASKTHIRD_PARTY},
		plan.ExtensionClass.ValueString())
	err := addOptionalThirdPartyRecurringTaskFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Recurring Task", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RecurringTaskApi.AddRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRecurringTaskRequest(
		client.AddThirdPartyRecurringTaskRequestAsAddRecurringTaskRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RecurringTaskApi.AddRecurringTaskExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Recurring Task", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state recurringTaskResourceModel
	readThirdPartyRecurringTaskResponse(ctx, addResponse.ThirdPartyRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *recurringTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan recurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *recurringTaskResourceModel
	var err error
	if plan.Type.ValueString() == "generate-server-profile" {
		state, err = r.CreateGenerateServerProfileRecurringTask(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "leave-lockdown-mode" {
		state, err = r.CreateLeaveLockdownModeRecurringTask(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "backup" {
		state, err = r.CreateBackupRecurringTask(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "delay" {
		state, err = r.CreateDelayRecurringTask(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "statically-defined" {
		state, err = r.CreateStaticallyDefinedRecurringTask(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "collect-support-data" {
		state, err = r.CreateCollectSupportDataRecurringTask(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "ldif-export" {
		state, err = r.CreateLdifExportRecurringTask(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "enter-lockdown-mode" {
		state, err = r.CreateEnterLockdownModeRecurringTask(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "audit-data-security" {
		state, err = r.CreateAuditDataSecurityRecurringTask(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "exec" {
		state, err = r.CreateExecRecurringTask(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "file-retention" {
		state, err = r.CreateFileRetentionRecurringTask(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyRecurringTask(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, *state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultRecurringTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan recurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RecurringTaskApi.GetRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Recurring Task", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state recurringTaskResourceModel
	if readResponse.GenerateServerProfileRecurringTaskResponse != nil {
		readGenerateServerProfileRecurringTaskResponse(ctx, readResponse.GenerateServerProfileRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LeaveLockdownModeRecurringTaskResponse != nil {
		readLeaveLockdownModeRecurringTaskResponse(ctx, readResponse.LeaveLockdownModeRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.BackupRecurringTaskResponse != nil {
		readBackupRecurringTaskResponse(ctx, readResponse.BackupRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DelayRecurringTaskResponse != nil {
		readDelayRecurringTaskResponse(ctx, readResponse.DelayRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.StaticallyDefinedRecurringTaskResponse != nil {
		readStaticallyDefinedRecurringTaskResponse(ctx, readResponse.StaticallyDefinedRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CollectSupportDataRecurringTaskResponse != nil {
		readCollectSupportDataRecurringTaskResponse(ctx, readResponse.CollectSupportDataRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LdifExportRecurringTaskResponse != nil {
		readLdifExportRecurringTaskResponse(ctx, readResponse.LdifExportRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.EnterLockdownModeRecurringTaskResponse != nil {
		readEnterLockdownModeRecurringTaskResponse(ctx, readResponse.EnterLockdownModeRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AuditDataSecurityRecurringTaskResponse != nil {
		readAuditDataSecurityRecurringTaskResponse(ctx, readResponse.AuditDataSecurityRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ExecRecurringTaskResponse != nil {
		readExecRecurringTaskResponse(ctx, readResponse.ExecRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileRetentionRecurringTaskResponse != nil {
		readFileRetentionRecurringTaskResponse(ctx, readResponse.FileRetentionRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyRecurringTaskResponse != nil {
		readThirdPartyRecurringTaskResponse(ctx, readResponse.ThirdPartyRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.RecurringTaskApi.UpdateRecurringTask(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createRecurringTaskOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.RecurringTaskApi.UpdateRecurringTaskExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Recurring Task", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.GenerateServerProfileRecurringTaskResponse != nil {
			readGenerateServerProfileRecurringTaskResponse(ctx, updateResponse.GenerateServerProfileRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LeaveLockdownModeRecurringTaskResponse != nil {
			readLeaveLockdownModeRecurringTaskResponse(ctx, updateResponse.LeaveLockdownModeRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.BackupRecurringTaskResponse != nil {
			readBackupRecurringTaskResponse(ctx, updateResponse.BackupRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DelayRecurringTaskResponse != nil {
			readDelayRecurringTaskResponse(ctx, updateResponse.DelayRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.StaticallyDefinedRecurringTaskResponse != nil {
			readStaticallyDefinedRecurringTaskResponse(ctx, updateResponse.StaticallyDefinedRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CollectSupportDataRecurringTaskResponse != nil {
			readCollectSupportDataRecurringTaskResponse(ctx, updateResponse.CollectSupportDataRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LdifExportRecurringTaskResponse != nil {
			readLdifExportRecurringTaskResponse(ctx, updateResponse.LdifExportRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.EnterLockdownModeRecurringTaskResponse != nil {
			readEnterLockdownModeRecurringTaskResponse(ctx, updateResponse.EnterLockdownModeRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AuditDataSecurityRecurringTaskResponse != nil {
			readAuditDataSecurityRecurringTaskResponse(ctx, updateResponse.AuditDataSecurityRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ExecRecurringTaskResponse != nil {
			readExecRecurringTaskResponse(ctx, updateResponse.ExecRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileRetentionRecurringTaskResponse != nil {
			readFileRetentionRecurringTaskResponse(ctx, updateResponse.FileRetentionRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyRecurringTaskResponse != nil {
			readThirdPartyRecurringTaskResponse(ctx, updateResponse.ThirdPartyRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
	}

	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *recurringTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readRecurringTask(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultRecurringTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readRecurringTask(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readRecurringTask(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state recurringTaskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.RecurringTaskApi.GetRecurringTask(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Recurring Task", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Recurring Task", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.GenerateServerProfileRecurringTaskResponse != nil {
		readGenerateServerProfileRecurringTaskResponse(ctx, readResponse.GenerateServerProfileRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LeaveLockdownModeRecurringTaskResponse != nil {
		readLeaveLockdownModeRecurringTaskResponse(ctx, readResponse.LeaveLockdownModeRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.BackupRecurringTaskResponse != nil {
		readBackupRecurringTaskResponse(ctx, readResponse.BackupRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DelayRecurringTaskResponse != nil {
		readDelayRecurringTaskResponse(ctx, readResponse.DelayRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.StaticallyDefinedRecurringTaskResponse != nil {
		readStaticallyDefinedRecurringTaskResponse(ctx, readResponse.StaticallyDefinedRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CollectSupportDataRecurringTaskResponse != nil {
		readCollectSupportDataRecurringTaskResponse(ctx, readResponse.CollectSupportDataRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LdifExportRecurringTaskResponse != nil {
		readLdifExportRecurringTaskResponse(ctx, readResponse.LdifExportRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.EnterLockdownModeRecurringTaskResponse != nil {
		readEnterLockdownModeRecurringTaskResponse(ctx, readResponse.EnterLockdownModeRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AuditDataSecurityRecurringTaskResponse != nil {
		readAuditDataSecurityRecurringTaskResponse(ctx, readResponse.AuditDataSecurityRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ExecRecurringTaskResponse != nil {
		readExecRecurringTaskResponse(ctx, readResponse.ExecRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileRetentionRecurringTaskResponse != nil {
		readFileRetentionRecurringTaskResponse(ctx, readResponse.FileRetentionRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyRecurringTaskResponse != nil {
		readThirdPartyRecurringTaskResponse(ctx, readResponse.ThirdPartyRecurringTaskResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *recurringTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateRecurringTask(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultRecurringTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateRecurringTask(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateRecurringTask(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan recurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state recurringTaskResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.RecurringTaskApi.UpdateRecurringTask(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createRecurringTaskOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.RecurringTaskApi.UpdateRecurringTaskExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Recurring Task", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.GenerateServerProfileRecurringTaskResponse != nil {
			readGenerateServerProfileRecurringTaskResponse(ctx, updateResponse.GenerateServerProfileRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LeaveLockdownModeRecurringTaskResponse != nil {
			readLeaveLockdownModeRecurringTaskResponse(ctx, updateResponse.LeaveLockdownModeRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.BackupRecurringTaskResponse != nil {
			readBackupRecurringTaskResponse(ctx, updateResponse.BackupRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DelayRecurringTaskResponse != nil {
			readDelayRecurringTaskResponse(ctx, updateResponse.DelayRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.StaticallyDefinedRecurringTaskResponse != nil {
			readStaticallyDefinedRecurringTaskResponse(ctx, updateResponse.StaticallyDefinedRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CollectSupportDataRecurringTaskResponse != nil {
			readCollectSupportDataRecurringTaskResponse(ctx, updateResponse.CollectSupportDataRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LdifExportRecurringTaskResponse != nil {
			readLdifExportRecurringTaskResponse(ctx, updateResponse.LdifExportRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.EnterLockdownModeRecurringTaskResponse != nil {
			readEnterLockdownModeRecurringTaskResponse(ctx, updateResponse.EnterLockdownModeRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AuditDataSecurityRecurringTaskResponse != nil {
			readAuditDataSecurityRecurringTaskResponse(ctx, updateResponse.AuditDataSecurityRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ExecRecurringTaskResponse != nil {
			readExecRecurringTaskResponse(ctx, updateResponse.ExecRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileRetentionRecurringTaskResponse != nil {
			readFileRetentionRecurringTaskResponse(ctx, updateResponse.FileRetentionRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyRecurringTaskResponse != nil {
			readThirdPartyRecurringTaskResponse(ctx, updateResponse.ThirdPartyRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *defaultRecurringTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *recurringTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state recurringTaskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.RecurringTaskApi.DeleteRecurringTaskExecute(r.apiClient.RecurringTaskApi.DeleteRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Recurring Task", err, httpResp)
		return
	}
}

func (r *recurringTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importRecurringTask(ctx, req, resp)
}

func (r *defaultRecurringTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importRecurringTask(ctx, req, resp)
}

func importRecurringTask(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
