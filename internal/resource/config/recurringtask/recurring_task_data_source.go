// Copyright © 2025 Ping Identity Corporation

package recurringtask

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &recurringTaskDataSource{}
	_ datasource.DataSourceWithConfigure = &recurringTaskDataSource{}
)

// Create a Recurring Task data source
func NewRecurringTaskDataSource() datasource.DataSource {
	return &recurringTaskDataSource{}
}

// recurringTaskDataSource is the datasource implementation.
type recurringTaskDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *recurringTaskDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_recurring_task"
}

// Configure adds the provider configured client to the data source.
func (r *recurringTaskDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type recurringTaskDataSourceModel struct {
	Id                                      types.String `tfsdk:"id"`
	Name                                    types.String `tfsdk:"name"`
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
	PostLDIFExportTaskProcessor             types.Set    `tfsdk:"post_ldif_export_task_processor"`
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

// GetSchema defines the schema for the datasource.
func (r *recurringTaskDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Recurring Task.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Recurring Task resource. Options are ['generate-server-profile', 'leave-lockdown-mode', 'backup', 'delay', 'statically-defined', 'collect-support-data', 'ldif-export', 'enter-lockdown-mode', 'audit-data-security', 'exec', 'file-retention', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Recurring Task.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Recurring Task. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"target_directory": schema.StringAttribute{
				Description: "The path to the directory containing the files to examine. The directory must exist.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"filename_pattern": schema.StringAttribute{
				Description: "A pattern that specifies the names of the files to examine. The pattern may contain zero or more asterisks as wildcards, where each wildcard matches zero or more characters. It may also contain at most one occurrence of the special string \"${timestamp}\", which will match a timestamp with the format specified using the timestamp-format property. All other characters in the pattern will be treated literally.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"timestamp_format": schema.StringAttribute{
				Description: "The format to use for the timestamp represented by the \"${timestamp}\" token in the filename pattern.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"retain_file_count": schema.Int64Attribute{
				Description: "The minimum number of files matching the pattern that will be retained.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"retain_file_age": schema.StringAttribute{
				Description: "The minimum age of files matching the pattern that will be retained.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"retain_aggregate_file_size": schema.StringAttribute{
				Description: "The minimum aggregate size of files that will be retained. The size should be specified as an integer followed by a unit that is one of \"b\" or \"bytes\", \"kb\" or \"kilobytes\", \"mb\" or \"megabytes\", \"gb\" or \"gigabytes\", or \"tb\" or \"terabytes\". For example, a value of \"1 gb\" indicates that at least one gigabyte of files should be retained.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"command_path": schema.StringAttribute{
				Description: "The absolute path to the command to execute. It must be an absolute path, the corresponding file must exist, and it must be listed in the config/exec-command-whitelist.txt file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"command_arguments": schema.StringAttribute{
				Description: "A string containing the arguments to provide to the command. If the command should be run without arguments, this property should be left undefined. If there should be multiple arguments, then they should be separated with spaces.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"command_output_file_base_name": schema.StringAttribute{
				Description: "The path and base name for a file to which the command output (both standard output and standard error) should be written. This may be left undefined if the command output should not be recorded into a file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"retain_previous_output_file_count": schema.Int64Attribute{
				Description: "The minimum number of previous command output files that should be preserved after a new instance of the command is invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"retain_previous_output_file_age": schema.StringAttribute{
				Description: "The minimum age of previous command output files that should be preserved after a new instance of the command is invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_command_output": schema.BoolAttribute{
				Description: "Indicates whether the command's output (both standard output and standard error) should be recorded in the server's error log.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"task_completion_state_for_nonzero_exit_code": schema.StringAttribute{
				Description: "The final task state that a task instance should have if the task executes the specified command and that command completes with a nonzero exit code, which generally means that the command did not complete successfully.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"working_directory": schema.StringAttribute{
				Description: "The absolute path to a working directory where the command should be executed. It must be an absolute path and the corresponding directory must exist.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"base_output_directory": schema.StringAttribute{
				Description: "The base directory below which generated reports will be written. Each invocation of the audit-data-security task will create a new subdirectory below this base directory whose name is a timestamp indicating when the report was generated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"data_security_auditor": schema.SetAttribute{
				Description: "The set of data security auditors that should be invoked. If no auditors are specified, then all auditors defined in the configuration will be used.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"backend": schema.SetAttribute{
				Description: "The set of backends that should be examined. If no backends are specified, then all backends that support this functionality will be included.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"include_filter": schema.SetAttribute{
				Description: "A filter that will be used to identify entries that may be included in the generated report. If multiple filters are specified, then any entry that matches at least one of the filters will be included. If no filters are specified, then all entries will be included.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"retain_previous_report_count": schema.Int64Attribute{
				Description: "The minimum number of previous reports that should be preserved after a new report is generated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"retain_previous_report_age": schema.StringAttribute{
				Description: "The minimum age of previous reports that should be preserved after a new report completes successfully.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ldif_directory": schema.StringAttribute{
				Description: "The directory in which LDIF export files will be placed. The directory must already exist.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"backend_id": schema.SetAttribute{
				Description: "The backend ID for a backend to be exported.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"exclude_backend_id": schema.SetAttribute{
				Description: "The backend ID for a backend to be excluded from the export.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"output_directory": schema.StringAttribute{
				Description: "The directory in which the support data archive files will be placed. The path must be a directory, and that directory must already exist. Relative paths will be interpreted as relative to the server root.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"encryption_passphrase_file": schema.StringAttribute{
				Description: "The path to a file that contains the passphrase to encrypt the contents of the support data archive.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_expensive_data": schema.BoolAttribute{
				Description: "Indicates whether the support data archive should include information that may be expensive to obtain, and that may temporarily affect the server's performance or responsiveness.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_replication_state_dump": schema.BoolAttribute{
				Description: "Indicates whether the support data archive should include a replication state dump, which may be several megabytes in size.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"retain_previous_ldif_export_count": schema.Int64Attribute{
				Description: "The minimum number of previous LDIF exports that should be preserved after a new export completes successfully.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"retain_previous_ldif_export_age": schema.StringAttribute{
				Description: "The minimum age of previous LDIF exports that should be preserved after a new export completes successfully.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_binary_files": schema.BoolAttribute{
				Description: "Indicates whether the support data archive should include binary files that may not have otherwise been included. Note that it may not be possible to obscure or redact sensitive information in binary files.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"post_ldif_export_task_processor": schema.SetAttribute{
				Description: "Supported in PingDirectory product version 10.0.0.0+. An optional set of post-LDIF-export task processors that should be invoked for the resulting LDIF export files.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"include_extension_source": schema.BoolAttribute{
				Description: "Indicates whether the support data archive should include the source code (if available) for any third-party extensions that may be installed in the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"use_sequential_mode": schema.BoolAttribute{
				Description: "Indicates whether to capture support data information sequentially rather than in parallel. Capturing data in sequential mode may reduce the amount of memory that the tool requires to operate, at the cost of taking longer to run.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"security_level": schema.StringAttribute{
				Description: "The security level to use when deciding which information to include in or exclude from the support data archive, and which included data should be obscured or redacted.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"jstack_count": schema.Int64Attribute{
				Description: "The number of times to invoke the jstack utility to obtain a stack trace of all threads running in the JVM. A value of zero indicates that the jstack utility should not be invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"report_count": schema.Int64Attribute{
				Description: "The number of intervals of data to collect from tools that use sample-based reporting, like vmstat, iostat, and mpstat. A value of zero indicates that these kinds of tools should not be used to collect any information.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"report_interval_seconds": schema.Int64Attribute{
				Description: "The duration (in seconds) between each interval of data to collect from tools that use sample-based reporting, like vmstat, iostat, and mpstat.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_duration": schema.StringAttribute{
				Description: "The maximum age (leading up to the time the collect-support-data tool was invoked) for log content to include in the support data archive.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_file_head_collection_size": schema.StringAttribute{
				Description: "The amount of data to collect from the beginning of each log file included in the support data archive.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_file_tail_collection_size": schema.StringAttribute{
				Description: "The amount of data to collect from the end of each log file included in the support data archive.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "An optional comment to include in a README file within the support data archive.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"retain_previous_support_data_archive_count": schema.Int64Attribute{
				Description: "The minimum number of previous support data archives that should be preserved after a new archive is generated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"retain_previous_support_data_archive_age": schema.StringAttribute{
				Description: "The minimum age of previous support data archives that should be preserved after a new archive is generated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"task_java_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class that provides the logic for the task to be invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"task_object_class": schema.SetAttribute{
				Description: "The names or OIDs of the object classes to include in the tasks that are scheduled from this Statically Defined Recurring Task. All object classes must be defined in the server schema, and the combination of object classes must be valid for a task entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"task_attribute_value": schema.SetAttribute{
				Description: "The set of attribute values that should be included in the tasks that are scheduled from this Statically Defined Recurring Task. Each value must be in the form {attribute-type}={value}, where {attribute-type} is the name or OID of an attribute type that is defined in the schema and permitted with the configured set of object classes, and {value} is a value to assign to an attribute with that type. A multivalued attribute can be created by providing multiple name-value pairs with the same name and different values.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"sleep_duration": schema.StringAttribute{
				Description: "The length of time to sleep before the task completes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"duration_to_wait_for_work_queue_idle": schema.StringAttribute{
				Description: "Indicates that task should wait for up to the specified length of time for the work queue to report that all worker threads are idle and there are no pending operations. Note that this primarily monitors operations that use worker threads, which does not include internal operations (for example, those invoked by extensions), and may not include requests from non-LDAP clients (for example, HTTP-based clients).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ldap_url_for_search_expected_to_return_entries": schema.SetAttribute{
				Description: "An LDAP URL that provides the criteria for a search request that is expected to return at least one entry. The search will be performed internally, and only the base DN, scope, and filter from the URL will be used; any host, port, or requested attributes included in the URL will be ignored.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"search_interval": schema.StringAttribute{
				Description: "The length of time the server should sleep between searches performed using the criteria from the ldap-url-for-search-expected-to-return-entries property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"search_time_limit": schema.StringAttribute{
				Description: "The length of time that the server will wait for a response to each internal search performed using the criteria from the ldap-url-for-search-expected-to-return-entries property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"duration_to_wait_for_search_to_return_entries": schema.StringAttribute{
				Description: "The maximum length of time that the server will continue to perform internal searches using the criteria from the ldap-url-for-search-expected-to-return-entries property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"task_return_state_if_timeout_is_encountered": schema.StringAttribute{
				Description: "The return state to use if a timeout is encountered while waiting for the server work queue to become idle (if the duration-to-wait-for-work-queue-idle property has a value), or if the time specified by the duration-to-wait-for-search-to-return-entries elapses without the associated search returning any entries.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"backup_directory": schema.StringAttribute{
				Description: "The directory in which backup files will be placed. When backing up a single backend, the backup files will be placed directly in this directory. When backing up multiple backends, the backup files for each backend will be placed in a subdirectory whose name is the corresponding backend ID.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"included_backend_id": schema.SetAttribute{
				Description: "The backend IDs of any backends that should be included in the backup.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_backend_id": schema.SetAttribute{
				Description: "The backend IDs of any backends that should be excluded from the backup. All backends that support backups and are not listed will be included.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"compress": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to `backup`: Indicates whether to compress the data as it is written into the backup. When the `type` attribute is set to `ldif-export`: Indicates whether to compress the LDIF data as it is exported.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `backup`: Indicates whether to compress the data as it is written into the backup.\n  - `ldif-export`: Indicates whether to compress the LDIF data as it is exported.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"encrypt": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to `backup`: Indicates whether to encrypt the data as it is written into the backup. When the `type` attribute is set to `ldif-export`: Indicates whether to encrypt the LDIF data as it exported.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `backup`: Indicates whether to encrypt the data as it is written into the backup.\n  - `ldif-export`: Indicates whether to encrypt the LDIF data as it exported.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"encryption_settings_definition_id": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `backup`: The ID of an encryption settings definition to use to obtain the backup encryption key. When the `type` attribute is set to `ldif-export`: The ID of an encryption settings definition to use to obtain the LDIF export encryption key.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `backup`: The ID of an encryption settings definition to use to obtain the backup encryption key.\n  - `ldif-export`: The ID of an encryption settings definition to use to obtain the LDIF export encryption key.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"sign": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to `backup`: Indicates whether to cryptographically sign backups, which will make it possible to detect whether the backup has been altered since it was created. When the `type` attribute is set to `ldif-export`: Indicates whether to cryptographically sign the exported data, which will make it possible to detect whether the LDIF data has been altered since it was exported.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `backup`: Indicates whether to cryptographically sign backups, which will make it possible to detect whether the backup has been altered since it was created.\n  - `ldif-export`: Indicates whether to cryptographically sign the exported data, which will make it possible to detect whether the LDIF data has been altered since it was exported.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"retain_previous_full_backup_count": schema.Int64Attribute{
				Description: "The minimum number of previous full backups that should be preserved after a new backup completes successfully.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"retain_previous_full_backup_age": schema.StringAttribute{
				Description: "The minimum age of previous full backups that should be preserved after a new backup completes successfully.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_megabytes_per_second": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to `backup`: The maximum rate, in megabytes per second, at which backups should be written. When the `type` attribute is set to `ldif-export`: The maximum rate, in megabytes per second, at which LDIF exports should be written.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `backup`: The maximum rate, in megabytes per second, at which backups should be written.\n  - `ldif-export`: The maximum rate, in megabytes per second, at which LDIF exports should be written.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"reason": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `leave-lockdown-mode`: The reason that the server is being taken out of in lockdown mode. When the `type` attribute is set to `enter-lockdown-mode`: The reason that the server is being placed in lockdown mode.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `leave-lockdown-mode`: The reason that the server is being taken out of in lockdown mode.\n  - `enter-lockdown-mode`: The reason that the server is being placed in lockdown mode.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"profile_directory": schema.StringAttribute{
				Description: "The directory in which the generated server profiles will be placed. The files will be named with the pattern \"server-profile-{timestamp}.zip\", where \"{timestamp}\" represents the time that the profile was generated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_path": schema.SetAttribute{
				Description: "An optional set of additional paths to files within the instance root that should be included in the generated server profile. All paths must be within the instance root, and relative paths will be relative to the instance root.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"retain_previous_profile_count": schema.Int64Attribute{
				Description: "The minimum number of previous server profile zip files that should be preserved after a new profile is generated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"retain_previous_profile_age": schema.StringAttribute{
				Description: "The minimum age of previous server profile zip files that should be preserved after a new profile is generated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Recurring Task",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"cancel_on_task_dependency_failure": schema.BoolAttribute{
				Description: "Indicates whether an instance of this Recurring Task should be canceled if the task immediately before it in the recurring task chain fails to complete successfully (including if it is canceled by an administrator before it starts or while it is running).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"email_on_start": schema.SetAttribute{
				Description: "The email addresses to which a message should be sent whenever an instance of this Recurring Task starts running. If this option is used, then at least one smtp-server must be configured in the global configuration.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"email_on_success": schema.SetAttribute{
				Description: "The email addresses to which a message should be sent whenever an instance of this Recurring Task completes successfully. If this option is used, then at least one smtp-server must be configured in the global configuration.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"email_on_failure": schema.SetAttribute{
				Description: "The email addresses to which a message should be sent if an instance of this Recurring Task fails to complete successfully. If this option is used, then at least one smtp-server must be configured in the global configuration.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"alert_on_start": schema.BoolAttribute{
				Description: "Indicates whether the server should generate an administrative alert whenever an instance of this Recurring Task starts running.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"alert_on_success": schema.BoolAttribute{
				Description: "Indicates whether the server should generate an administrative alert whenever an instance of this Recurring Task completes successfully.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"alert_on_failure": schema.BoolAttribute{
				Description: "Indicates whether the server should generate an administrative alert whenever an instance of this Recurring Task fails to complete successfully.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a GenerateServerProfileRecurringTaskResponse object into the model struct
func readGenerateServerProfileRecurringTaskResponseDataSource(ctx context.Context, r *client.GenerateServerProfileRecurringTaskResponse, state *recurringTaskDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generate-server-profile")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ProfileDirectory = types.StringValue(r.ProfileDirectory)
	state.IncludePath = internaltypes.GetStringSet(r.IncludePath)
	state.RetainPreviousProfileCount = internaltypes.Int64TypeOrNil(r.RetainPreviousProfileCount)
	state.RetainPreviousProfileAge = internaltypes.StringTypeOrNil(r.RetainPreviousProfileAge, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
}

// Read a LeaveLockdownModeRecurringTaskResponse object into the model struct
func readLeaveLockdownModeRecurringTaskResponseDataSource(ctx context.Context, r *client.LeaveLockdownModeRecurringTaskResponse, state *recurringTaskDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("leave-lockdown-mode")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Reason = internaltypes.StringTypeOrNil(r.Reason, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
}

// Read a BackupRecurringTaskResponse object into the model struct
func readBackupRecurringTaskResponseDataSource(ctx context.Context, r *client.BackupRecurringTaskResponse, state *recurringTaskDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("backup")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BackupDirectory = types.StringValue(r.BackupDirectory)
	state.IncludedBackendID = internaltypes.GetStringSet(r.IncludedBackendID)
	state.ExcludedBackendID = internaltypes.GetStringSet(r.ExcludedBackendID)
	state.Compress = internaltypes.BoolTypeOrNil(r.Compress)
	state.Encrypt = internaltypes.BoolTypeOrNil(r.Encrypt)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Sign = internaltypes.BoolTypeOrNil(r.Sign)
	state.RetainPreviousFullBackupCount = internaltypes.Int64TypeOrNil(r.RetainPreviousFullBackupCount)
	state.RetainPreviousFullBackupAge = internaltypes.StringTypeOrNil(r.RetainPreviousFullBackupAge, false)
	state.MaxMegabytesPerSecond = internaltypes.Int64TypeOrNil(r.MaxMegabytesPerSecond)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
}

// Read a DelayRecurringTaskResponse object into the model struct
func readDelayRecurringTaskResponseDataSource(ctx context.Context, r *client.DelayRecurringTaskResponse, state *recurringTaskDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("delay")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SleepDuration = internaltypes.StringTypeOrNil(r.SleepDuration, false)
	state.DurationToWaitForWorkQueueIdle = internaltypes.StringTypeOrNil(r.DurationToWaitForWorkQueueIdle, false)
	state.LdapURLForSearchExpectedToReturnEntries = internaltypes.GetStringSet(r.LdapURLForSearchExpectedToReturnEntries)
	state.SearchInterval = internaltypes.StringTypeOrNil(r.SearchInterval, false)
	state.SearchTimeLimit = internaltypes.StringTypeOrNil(r.SearchTimeLimit, false)
	state.DurationToWaitForSearchToReturnEntries = internaltypes.StringTypeOrNil(r.DurationToWaitForSearchToReturnEntries, false)
	state.TaskReturnStateIfTimeoutIsEncountered = internaltypes.StringTypeOrNil(
		client.StringPointerEnumrecurringTaskTaskReturnStateIfTimeoutIsEncounteredProp(r.TaskReturnStateIfTimeoutIsEncountered), false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
}

// Read a StaticallyDefinedRecurringTaskResponse object into the model struct
func readStaticallyDefinedRecurringTaskResponseDataSource(ctx context.Context, r *client.StaticallyDefinedRecurringTaskResponse, state *recurringTaskDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("statically-defined")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.TaskJavaClass = types.StringValue(r.TaskJavaClass)
	state.TaskObjectClass = internaltypes.GetStringSet(r.TaskObjectClass)
	state.TaskAttributeValue = internaltypes.GetStringSet(r.TaskAttributeValue)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
}

// Read a CollectSupportDataRecurringTaskResponse object into the model struct
func readCollectSupportDataRecurringTaskResponseDataSource(ctx context.Context, r *client.CollectSupportDataRecurringTaskResponse, state *recurringTaskDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("collect-support-data")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.OutputDirectory = types.StringValue(r.OutputDirectory)
	state.EncryptionPassphraseFile = internaltypes.StringTypeOrNil(r.EncryptionPassphraseFile, false)
	state.IncludeExpensiveData = internaltypes.BoolTypeOrNil(r.IncludeExpensiveData)
	state.IncludeReplicationStateDump = internaltypes.BoolTypeOrNil(r.IncludeReplicationStateDump)
	state.IncludeBinaryFiles = internaltypes.BoolTypeOrNil(r.IncludeBinaryFiles)
	state.IncludeExtensionSource = internaltypes.BoolTypeOrNil(r.IncludeExtensionSource)
	state.UseSequentialMode = internaltypes.BoolTypeOrNil(r.UseSequentialMode)
	state.SecurityLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumrecurringTaskSecurityLevelProp(r.SecurityLevel), false)
	state.JstackCount = internaltypes.Int64TypeOrNil(r.JstackCount)
	state.ReportCount = internaltypes.Int64TypeOrNil(r.ReportCount)
	state.ReportIntervalSeconds = internaltypes.Int64TypeOrNil(r.ReportIntervalSeconds)
	state.LogDuration = internaltypes.StringTypeOrNil(r.LogDuration, false)
	state.LogFileHeadCollectionSize = internaltypes.StringTypeOrNil(r.LogFileHeadCollectionSize, false)
	state.LogFileTailCollectionSize = internaltypes.StringTypeOrNil(r.LogFileTailCollectionSize, false)
	state.Comment = internaltypes.StringTypeOrNil(r.Comment, false)
	state.RetainPreviousSupportDataArchiveCount = internaltypes.Int64TypeOrNil(r.RetainPreviousSupportDataArchiveCount)
	state.RetainPreviousSupportDataArchiveAge = internaltypes.StringTypeOrNil(r.RetainPreviousSupportDataArchiveAge, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
}

// Read a LdifExportRecurringTaskResponse object into the model struct
func readLdifExportRecurringTaskResponseDataSource(ctx context.Context, r *client.LdifExportRecurringTaskResponse, state *recurringTaskDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldif-export")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LdifDirectory = types.StringValue(r.LdifDirectory)
	state.BackendID = internaltypes.GetStringSet(r.BackendID)
	state.ExcludeBackendID = internaltypes.GetStringSet(r.ExcludeBackendID)
	state.Compress = internaltypes.BoolTypeOrNil(r.Compress)
	state.Encrypt = internaltypes.BoolTypeOrNil(r.Encrypt)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Sign = internaltypes.BoolTypeOrNil(r.Sign)
	state.RetainPreviousLDIFExportCount = internaltypes.Int64TypeOrNil(r.RetainPreviousLDIFExportCount)
	state.RetainPreviousLDIFExportAge = internaltypes.StringTypeOrNil(r.RetainPreviousLDIFExportAge, false)
	state.MaxMegabytesPerSecond = internaltypes.Int64TypeOrNil(r.MaxMegabytesPerSecond)
	state.PostLDIFExportTaskProcessor = internaltypes.GetStringSet(r.PostLDIFExportTaskProcessor)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
}

// Read a EnterLockdownModeRecurringTaskResponse object into the model struct
func readEnterLockdownModeRecurringTaskResponseDataSource(ctx context.Context, r *client.EnterLockdownModeRecurringTaskResponse, state *recurringTaskDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("enter-lockdown-mode")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Reason = internaltypes.StringTypeOrNil(r.Reason, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
}

// Read a AuditDataSecurityRecurringTaskResponse object into the model struct
func readAuditDataSecurityRecurringTaskResponseDataSource(ctx context.Context, r *client.AuditDataSecurityRecurringTaskResponse, state *recurringTaskDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("audit-data-security")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BaseOutputDirectory = types.StringValue(r.BaseOutputDirectory)
	state.DataSecurityAuditor = internaltypes.GetStringSet(r.DataSecurityAuditor)
	state.Backend = internaltypes.GetStringSet(r.Backend)
	state.IncludeFilter = internaltypes.GetStringSet(r.IncludeFilter)
	state.RetainPreviousReportCount = internaltypes.Int64TypeOrNil(r.RetainPreviousReportCount)
	state.RetainPreviousReportAge = internaltypes.StringTypeOrNil(r.RetainPreviousReportAge, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
}

// Read a ExecRecurringTaskResponse object into the model struct
func readExecRecurringTaskResponseDataSource(ctx context.Context, r *client.ExecRecurringTaskResponse, state *recurringTaskDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("exec")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.CommandPath = types.StringValue(r.CommandPath)
	state.CommandArguments = internaltypes.StringTypeOrNil(r.CommandArguments, false)
	state.CommandOutputFileBaseName = internaltypes.StringTypeOrNil(r.CommandOutputFileBaseName, false)
	state.RetainPreviousOutputFileCount = internaltypes.Int64TypeOrNil(r.RetainPreviousOutputFileCount)
	state.RetainPreviousOutputFileAge = internaltypes.StringTypeOrNil(r.RetainPreviousOutputFileAge, false)
	state.LogCommandOutput = internaltypes.BoolTypeOrNil(r.LogCommandOutput)
	state.TaskCompletionStateForNonzeroExitCode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumrecurringTaskTaskCompletionStateForNonzeroExitCodeProp(r.TaskCompletionStateForNonzeroExitCode), false)
	state.WorkingDirectory = internaltypes.StringTypeOrNil(r.WorkingDirectory, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
}

// Read a FileRetentionRecurringTaskResponse object into the model struct
func readFileRetentionRecurringTaskResponseDataSource(ctx context.Context, r *client.FileRetentionRecurringTaskResponse, state *recurringTaskDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-retention")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.TargetDirectory = types.StringValue(r.TargetDirectory)
	state.FilenamePattern = types.StringValue(r.FilenamePattern)
	state.TimestampFormat = types.StringValue(r.TimestampFormat.String())
	state.RetainFileCount = internaltypes.Int64TypeOrNil(r.RetainFileCount)
	state.RetainFileAge = internaltypes.StringTypeOrNil(r.RetainFileAge, false)
	state.RetainAggregateFileSize = internaltypes.StringTypeOrNil(r.RetainAggregateFileSize, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
}

// Read a ThirdPartyRecurringTaskResponse object into the model struct
func readThirdPartyRecurringTaskResponseDataSource(ctx context.Context, r *client.ThirdPartyRecurringTaskResponse, state *recurringTaskDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
}

// Read resource information
func (r *recurringTaskDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state recurringTaskDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RecurringTaskAPI.GetRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Recurring Task", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.GenerateServerProfileRecurringTaskResponse != nil {
		readGenerateServerProfileRecurringTaskResponseDataSource(ctx, readResponse.GenerateServerProfileRecurringTaskResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LeaveLockdownModeRecurringTaskResponse != nil {
		readLeaveLockdownModeRecurringTaskResponseDataSource(ctx, readResponse.LeaveLockdownModeRecurringTaskResponse, &state, &resp.Diagnostics)
	}
	if readResponse.BackupRecurringTaskResponse != nil {
		readBackupRecurringTaskResponseDataSource(ctx, readResponse.BackupRecurringTaskResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DelayRecurringTaskResponse != nil {
		readDelayRecurringTaskResponseDataSource(ctx, readResponse.DelayRecurringTaskResponse, &state, &resp.Diagnostics)
	}
	if readResponse.StaticallyDefinedRecurringTaskResponse != nil {
		readStaticallyDefinedRecurringTaskResponseDataSource(ctx, readResponse.StaticallyDefinedRecurringTaskResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CollectSupportDataRecurringTaskResponse != nil {
		readCollectSupportDataRecurringTaskResponseDataSource(ctx, readResponse.CollectSupportDataRecurringTaskResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LdifExportRecurringTaskResponse != nil {
		readLdifExportRecurringTaskResponseDataSource(ctx, readResponse.LdifExportRecurringTaskResponse, &state, &resp.Diagnostics)
	}
	if readResponse.EnterLockdownModeRecurringTaskResponse != nil {
		readEnterLockdownModeRecurringTaskResponseDataSource(ctx, readResponse.EnterLockdownModeRecurringTaskResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AuditDataSecurityRecurringTaskResponse != nil {
		readAuditDataSecurityRecurringTaskResponseDataSource(ctx, readResponse.AuditDataSecurityRecurringTaskResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ExecRecurringTaskResponse != nil {
		readExecRecurringTaskResponseDataSource(ctx, readResponse.ExecRecurringTaskResponse, &state, &resp.Diagnostics)
	}
	if readResponse.FileRetentionRecurringTaskResponse != nil {
		readFileRetentionRecurringTaskResponseDataSource(ctx, readResponse.FileRetentionRecurringTaskResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyRecurringTaskResponse != nil {
		readThirdPartyRecurringTaskResponseDataSource(ctx, readResponse.ThirdPartyRecurringTaskResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
