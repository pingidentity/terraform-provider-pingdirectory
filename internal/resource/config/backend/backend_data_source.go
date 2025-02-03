// Copyright © 2025 Ping Identity Corporation

package backend

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
	_ datasource.DataSource              = &backendDataSource{}
	_ datasource.DataSourceWithConfigure = &backendDataSource{}
)

// Create a Backend data source
func NewBackendDataSource() datasource.DataSource {
	return &backendDataSource{}
}

// backendDataSource is the datasource implementation.
type backendDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *backendDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backend"
}

// Configure adds the provider configured client to the data source.
func (r *backendDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type backendDataSourceModel struct {
	Id                                          types.String `tfsdk:"id"`
	Type                                        types.String `tfsdk:"type"`
	UncachedId2entryCacheMode                   types.String `tfsdk:"uncached_id2entry_cache_mode"`
	StorageDir                                  types.String `tfsdk:"storage_dir"`
	MetricsDir                                  types.String `tfsdk:"metrics_dir"`
	SampleFlushInterval                         types.String `tfsdk:"sample_flush_interval"`
	RetentionPolicy                             types.Set    `tfsdk:"retention_policy"`
	UncachedAttributeCriteria                   types.String `tfsdk:"uncached_attribute_criteria"`
	UncachedEntryCriteria                       types.String `tfsdk:"uncached_entry_criteria"`
	AlarmRetentionTime                          types.String `tfsdk:"alarm_retention_time"`
	MaxAlarms                                   types.Int64  `tfsdk:"max_alarms"`
	AlertRetentionTime                          types.String `tfsdk:"alert_retention_time"`
	MaxAlerts                                   types.Int64  `tfsdk:"max_alerts"`
	DisabledAlertType                           types.Set    `tfsdk:"disabled_alert_type"`
	TaskBackingFile                             types.String `tfsdk:"task_backing_file"`
	MaximumInitialTaskLogMessagesToRetain       types.Int64  `tfsdk:"maximum_initial_task_log_messages_to_retain"`
	MaximumFinalTaskLogMessagesToRetain         types.Int64  `tfsdk:"maximum_final_task_log_messages_to_retain"`
	TaskRetentionTime                           types.String `tfsdk:"task_retention_time"`
	NotificationSenderAddress                   types.String `tfsdk:"notification_sender_address"`
	InsignificantConfigArchiveAttribute         types.Set    `tfsdk:"insignificant_config_archive_attribute"`
	InsignificantConfigArchiveBaseDN            types.Set    `tfsdk:"insignificant_config_archive_base_dn"`
	MaintainConfigArchive                       types.Bool   `tfsdk:"maintain_config_archive"`
	MaxConfigArchiveCount                       types.Int64  `tfsdk:"max_config_archive_count"`
	MirroredSubtreePeerPollingInterval          types.String `tfsdk:"mirrored_subtree_peer_polling_interval"`
	MirroredSubtreeEntryUpdateTimeout           types.String `tfsdk:"mirrored_subtree_entry_update_timeout"`
	MirroredSubtreeSearchTimeout                types.String `tfsdk:"mirrored_subtree_search_timeout"`
	BackendID                                   types.String `tfsdk:"backend_id"`
	SetDegradedAlertForUntrustedIndex           types.Bool   `tfsdk:"set_degraded_alert_for_untrusted_index"`
	ReturnUnavailableForUntrustedIndex          types.Bool   `tfsdk:"return_unavailable_for_untrusted_index"`
	ProcessFiltersWithUndefinedAttributeTypes   types.Bool   `tfsdk:"process_filters_with_undefined_attribute_types"`
	DbDirectory                                 types.String `tfsdk:"db_directory"`
	DbDirectoryPermissions                      types.String `tfsdk:"db_directory_permissions"`
	DbCachePercent                              types.Int64  `tfsdk:"db_cache_percent"`
	CompactCommonParentDN                       types.Set    `tfsdk:"compact_common_parent_dn"`
	CompressEntries                             types.Bool   `tfsdk:"compress_entries"`
	HashEntries                                 types.Bool   `tfsdk:"hash_entries"`
	DbNumCleanerThreads                         types.Int64  `tfsdk:"db_num_cleaner_threads"`
	DbCleanerMinUtilization                     types.Int64  `tfsdk:"db_cleaner_min_utilization"`
	DbEvictorCriticalPercentage                 types.Int64  `tfsdk:"db_evictor_critical_percentage"`
	DbCheckpointerWakeupInterval                types.String `tfsdk:"db_checkpointer_wakeup_interval"`
	DbBackgroundSyncInterval                    types.String `tfsdk:"db_background_sync_interval"`
	DbUseThreadLocalHandles                     types.Bool   `tfsdk:"db_use_thread_local_handles"`
	DbLogFileMax                                types.String `tfsdk:"db_log_file_max"`
	DbLoggingLevel                              types.String `tfsdk:"db_logging_level"`
	JeProperty                                  types.Set    `tfsdk:"je_property"`
	ChangelogWriteBatchSize                     types.Int64  `tfsdk:"changelog_write_batch_size"`
	DefaultCacheMode                            types.String `tfsdk:"default_cache_mode"`
	Id2entryCacheMode                           types.String `tfsdk:"id2entry_cache_mode"`
	Dn2idCacheMode                              types.String `tfsdk:"dn2id_cache_mode"`
	Id2childrenCacheMode                        types.String `tfsdk:"id2children_cache_mode"`
	Id2subtreeCacheMode                         types.String `tfsdk:"id2subtree_cache_mode"`
	Dn2uriCacheMode                             types.String `tfsdk:"dn2uri_cache_mode"`
	SimplePagedResultsIDSetCacheDuration        types.String `tfsdk:"simple_paged_results_id_set_cache_duration"`
	PrimeMethod                                 types.Set    `tfsdk:"prime_method"`
	PrimeThreadCount                            types.Int64  `tfsdk:"prime_thread_count"`
	PrimeTimeLimit                              types.String `tfsdk:"prime_time_limit"`
	PrimeAllIndexes                             types.Bool   `tfsdk:"prime_all_indexes"`
	SystemIndexToPrime                          types.Set    `tfsdk:"system_index_to_prime"`
	SystemIndexToPrimeInternalNodesOnly         types.Set    `tfsdk:"system_index_to_prime_internal_nodes_only"`
	BackgroundPrime                             types.Bool   `tfsdk:"background_prime"`
	IndexEntryLimit                             types.Int64  `tfsdk:"index_entry_limit"`
	CompositeIndexEntryLimit                    types.Int64  `tfsdk:"composite_index_entry_limit"`
	Id2childrenIndexEntryLimit                  types.Int64  `tfsdk:"id2children_index_entry_limit"`
	Id2subtreeIndexEntryLimit                   types.Int64  `tfsdk:"id2subtree_index_entry_limit"`
	ImportTempDirectory                         types.String `tfsdk:"import_temp_directory"`
	ImportThreadCount                           types.Int64  `tfsdk:"import_thread_count"`
	ExportThreadCount                           types.Int64  `tfsdk:"export_thread_count"`
	DbImportCachePercent                        types.Int64  `tfsdk:"db_import_cache_percent"`
	DbTxnWriteNoSync                            types.Bool   `tfsdk:"db_txn_write_no_sync"`
	DeadlockRetryLimit                          types.Int64  `tfsdk:"deadlock_retry_limit"`
	ExternalTxnDefaultBackendLockBehavior       types.String `tfsdk:"external_txn_default_backend_lock_behavior"`
	SingleWriterLockBehavior                    types.String `tfsdk:"single_writer_lock_behavior"`
	SubtreeModifyDNSizeLimit                    types.Int64  `tfsdk:"subtree_modify_dn_size_limit"`
	SubtreeDeleteSizeLimit                      types.Int64  `tfsdk:"subtree_delete_size_limit"`
	NumRecentChanges                            types.Int64  `tfsdk:"num_recent_changes"`
	OfflineProcessDatabaseOpenTimeout           types.String `tfsdk:"offline_process_database_open_timeout"`
	ChangelogPurgeBatchSize                     types.Int64  `tfsdk:"changelog_purge_batch_size"`
	ChangelogWriteQueueCapacity                 types.Int64  `tfsdk:"changelog_write_queue_capacity"`
	IndexIncludeAttribute                       types.Set    `tfsdk:"index_include_attribute"`
	IndexExcludeAttribute                       types.Set    `tfsdk:"index_exclude_attribute"`
	ChangelogMaximumAge                         types.String `tfsdk:"changelog_maximum_age"`
	TargetDatabaseSize                          types.String `tfsdk:"target_database_size"`
	ChangelogEntryIncludeBaseDN                 types.Set    `tfsdk:"changelog_entry_include_base_dn"`
	ChangelogEntryExcludeBaseDN                 types.Set    `tfsdk:"changelog_entry_exclude_base_dn"`
	ChangelogEntryIncludeFilter                 types.Set    `tfsdk:"changelog_entry_include_filter"`
	ChangelogEntryExcludeFilter                 types.Set    `tfsdk:"changelog_entry_exclude_filter"`
	ChangelogIncludeAttribute                   types.Set    `tfsdk:"changelog_include_attribute"`
	ChangelogExcludeAttribute                   types.Set    `tfsdk:"changelog_exclude_attribute"`
	ChangelogDeletedEntryIncludeAttribute       types.Set    `tfsdk:"changelog_deleted_entry_include_attribute"`
	ChangelogDeletedEntryExcludeAttribute       types.Set    `tfsdk:"changelog_deleted_entry_exclude_attribute"`
	ChangelogIncludeKeyAttribute                types.Set    `tfsdk:"changelog_include_key_attribute"`
	ChangelogMaxBeforeAfterValues               types.Int64  `tfsdk:"changelog_max_before_after_values"`
	WriteLastmodAttributes                      types.Bool   `tfsdk:"write_lastmod_attributes"`
	UseReversibleForm                           types.Bool   `tfsdk:"use_reversible_form"`
	IncludeVirtualAttributes                    types.Set    `tfsdk:"include_virtual_attributes"`
	ApplyAccessControlsToChangelogEntryContents types.Bool   `tfsdk:"apply_access_controls_to_changelog_entry_contents"`
	ReportExcludedChangelogAttributes           types.String `tfsdk:"report_excluded_changelog_attributes"`
	SoftDeleteEntryIncludedOperation            types.Set    `tfsdk:"soft_delete_entry_included_operation"`
	IsPrivateBackend                            types.Bool   `tfsdk:"is_private_backend"`
	LdifFile                                    types.String `tfsdk:"ldif_file"`
	TrustStoreFile                              types.String `tfsdk:"trust_store_file"`
	TrustStoreType                              types.String `tfsdk:"trust_store_type"`
	TrustStorePin                               types.String `tfsdk:"trust_store_pin"`
	TrustStorePinFile                           types.String `tfsdk:"trust_store_pin_file"`
	TrustStorePinPassphraseProvider             types.String `tfsdk:"trust_store_pin_passphrase_provider"`
	BaseDN                                      types.Set    `tfsdk:"base_dn"`
	WritabilityMode                             types.String `tfsdk:"writability_mode"`
	BackupDirectory                             types.Set    `tfsdk:"backup_directory"`
	SchemaEntryDN                               types.Set    `tfsdk:"schema_entry_dn"`
	ShowAllAttributes                           types.Bool   `tfsdk:"show_all_attributes"`
	ReadOnlySchemaFile                          types.Set    `tfsdk:"read_only_schema_file"`
	Description                                 types.String `tfsdk:"description"`
	Enabled                                     types.Bool   `tfsdk:"enabled"`
	SetDegradedAlertWhenDisabled                types.Bool   `tfsdk:"set_degraded_alert_when_disabled"`
	ReturnUnavailableWhenDisabled               types.Bool   `tfsdk:"return_unavailable_when_disabled"`
	BackupFilePermissions                       types.String `tfsdk:"backup_file_permissions"`
	NotificationManager                         types.String `tfsdk:"notification_manager"`
}

// GetSchema defines the schema for the datasource.
func (r *backendDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Backend.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Backend resource. Options are ['schema', 'backup', 'encryption-settings', 'ldif', 'trust-store', 'custom', 'changelog', 'monitor', 'local-db', 'config-file-handler', 'task', 'alert', 'alarm', 'metrics']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"uncached_id2entry_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the uncached-id2entry database, which provides a way to store complete or partial encoded entries with a different (and presumably less memory-intensive) cache mode than records written to id2entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"storage_dir": schema.StringAttribute{
				Description: "Specifies the path to the directory that will be used to store queued samples.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"metrics_dir": schema.StringAttribute{
				Description: "Specifies the path to the directory that contains metric definitions.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"sample_flush_interval": schema.StringAttribute{
				Description: "Period when samples are flushed to disk.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"retention_policy": schema.SetAttribute{
				Description: "The retention policy to use for the Metrics Backend .",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"uncached_attribute_criteria": schema.StringAttribute{
				Description: "The criteria that will be used to identify attributes that should be written into the uncached-id2entry database rather than the id2entry database. This will only be used for entries in which the associated uncached-entry-criteria does not indicate that the entire entry should be uncached.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"uncached_entry_criteria": schema.StringAttribute{
				Description: "The criteria that will be used to identify entries that should be written into the uncached-id2entry database rather than the id2entry database.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"alarm_retention_time": schema.StringAttribute{
				Description: "Specifies the maximum length of time that information about raised alarms should be maintained before they will be purged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_alarms": schema.Int64Attribute{
				Description: "Specifies the maximum number of alarms that should be retained. If more alarms than this configured maximum are generated within the alarm retention time, then the oldest alarms will be purged to achieve this maximum. Only alarms at normal severity will be purged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"alert_retention_time": schema.StringAttribute{
				Description: "Specifies the maximum length of time that information about generated alerts should be maintained before they will be purged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_alerts": schema.Int64Attribute{
				Description: "Specifies the maximum number of alerts that should be retained. If more alerts than this configured maximum are generated within the alert retention time, then the oldest alerts will be purged to achieve this maximum.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"disabled_alert_type": schema.SetAttribute{
				Description: "Specifies the names of the alert types that should not be added to the backend. This can be used to suppress high volume alerts that might trigger hitting the max-alerts limit sooner than desired. Disabled alert types will not be sent out over persistent searches on this backend.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"task_backing_file": schema.StringAttribute{
				Description: "Specifies the path to the backing file for storing information about the tasks configured in the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_initial_task_log_messages_to_retain": schema.Int64Attribute{
				Description: "The maximum number of log messages to retain in each task entry from the beginning of the processing for that task. If too many messages are logged during task processing, then retaining only a limited number of messages from the beginning and/or end of task processing can reduce the amount of memory that the server consumes by caching information about currently-active and recently-completed tasks.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_final_task_log_messages_to_retain": schema.Int64Attribute{
				Description: "The maximum number of log messages to retain in each task entry from the end of the processing for that task. If too many messages are logged during task processing, then retaining only a limited number of messages from the beginning and/or end of task processing can reduce the amount of memory that the server consumes by caching information about currently-active and recently-completed tasks.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"task_retention_time": schema.StringAttribute{
				Description: "Specifies the length of time that task entries should be retained after processing on the associated task has been completed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"notification_sender_address": schema.StringAttribute{
				Description: "Specifies the email address to use as the sender address (that is, the \"From:\" address) for notification mail messages generated when a task completes execution.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"insignificant_config_archive_attribute": schema.SetAttribute{
				Description: "The name or OID of an attribute type that is considered insignificant for the purpose of maintaining the configuration archive.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"insignificant_config_archive_base_dn": schema.SetAttribute{
				Description: "Supported in PingDirectory product version 9.2.0.3+. The base DN that is considered insignificant for the purpose of maintaining the configuration archive.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"maintain_config_archive": schema.BoolAttribute{
				Description: "Supported in PingDirectory product version 9.3.0.0+. Indicates whether the server should maintain the config archive with new changes to the config backend.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_config_archive_count": schema.Int64Attribute{
				Description: "Supported in PingDirectory product version 9.3.0.0+. Indicates the maximum number of previous config files to keep as part of maintaining the config archive.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"mirrored_subtree_peer_polling_interval": schema.StringAttribute{
				Description: "Tells the server component that is responsible for mirroring configuration data across a topology of servers the maximum amount of time to wait before polling the peer servers in the topology to determine if there are any changes in the topology. Mirrored data includes meta-data about the servers in the topology as well as cluster-wide configuration data.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"mirrored_subtree_entry_update_timeout": schema.StringAttribute{
				Description: "Tells the server component that is responsible for mirroring configuration data across a topology of servers the maximum amount of time to wait for an update operation (add, delete, modify and modify-dn) on an entry to be applied on all servers in the topology. Mirrored data includes meta-data about the servers in the topology as well as cluster-wide configuration data.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"mirrored_subtree_search_timeout": schema.StringAttribute{
				Description: "Tells the server component that is responsible for mirroring configuration data across a topology of servers the maximum amount of time to wait for a search operation to complete. Mirrored data includes meta-data about the servers in the topology as well as cluster-wide configuration data. Search requests that take longer than this timeout will be canceled and considered failures.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"backend_id": schema.StringAttribute{
				Description: "Specifies a name to identify the associated backend.",
				Required:    true,
			},
			"set_degraded_alert_for_untrusted_index": schema.BoolAttribute{
				Description: "Determines whether the Directory Server enters a DEGRADED state when this Local DB Backend has an index whose contents cannot be trusted.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"return_unavailable_for_untrusted_index": schema.BoolAttribute{
				Description: "Determines whether the Directory Server returns UNAVAILABLE for any LDAP search operation in this Local DB Backend that would use an index whose contents cannot be trusted.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"process_filters_with_undefined_attribute_types": schema.BoolAttribute{
				Description: "Determines whether the Directory Server should continue filter processing for LDAP search operations in this Local DB Backend that includes a search filter with an attribute that is not defined in the schema. This will only apply if check-schema is enabled in the global configuration.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"db_directory": schema.StringAttribute{
				Description: "Specifies the path to the filesystem directory that is used to hold the Berkeley DB Java Edition database files containing the data for this backend. The files for this backend are stored in a sub-directory named after the backend-id.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"db_directory_permissions": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `changelog`: Specifies the permissions that should be applied to the directory containing the backend database files and to directories and files created during backup of the backend. When the `type` attribute is set to `local-db`: Specifies the permissions that should be applied to the directory containing the backend database files and to directories and files created during backup or LDIF export of the backend.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `changelog`: Specifies the permissions that should be applied to the directory containing the backend database files and to directories and files created during backup of the backend.\n  - `local-db`: Specifies the permissions that should be applied to the directory containing the backend database files and to directories and files created during backup or LDIF export of the backend.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"db_cache_percent": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to `changelog`: Specifies the percentage of JVM memory to allocate to the changelog database cache. When the `type` attribute is set to `local-db`: Specifies the percentage of JVM memory to allocate to the database cache.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `changelog`: Specifies the percentage of JVM memory to allocate to the changelog database cache.\n  - `local-db`: Specifies the percentage of JVM memory to allocate to the database cache.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"compact_common_parent_dn": schema.SetAttribute{
				Description: "Provides a DN of an entry that may be the parent for a large number of entries in the backend. This may be used to help increase the space efficiency when encoding entries for storage.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"compress_entries": schema.BoolAttribute{
				Description: "Indicates whether the backend should attempt to compress entries before storing them in the database.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"hash_entries": schema.BoolAttribute{
				Description: "Indicates whether to calculate and store a message digest of the entry contents along with the entry data, in order to provide a means of verifying the integrity of the entry data.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"db_num_cleaner_threads": schema.Int64Attribute{
				Description: "Specifies the number of threads that the backend should maintain to keep the database log files at or near the desired utilization. A value of zero indicates that the number of cleaner threads should be automatically configured based on the number of available CPUs.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"db_cleaner_min_utilization": schema.Int64Attribute{
				Description: "Specifies the minimum percentage of \"live\" data that the database cleaner attempts to keep in database log files.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"db_evictor_critical_percentage": schema.Int64Attribute{
				Description: "Specifies the percentage over the configured maximum that the database cache is allowed to grow. It is recommended to set this value slightly above zero when the database is too large to fully cache in memory. In this case, a dedicated background evictor thread is used to perform evictions once the cache fills up reducing the possibility that server threads are blocked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"db_checkpointer_wakeup_interval": schema.StringAttribute{
				Description: "Specifies the maximum length of time that should pass between checkpoints.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"db_background_sync_interval": schema.StringAttribute{
				Description: "Specifies the interval to use when performing background synchronous writes in the database environment in order to smooth overall write performance and increase data durability. A value of \"0 s\" will disable background synchronous writes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"db_use_thread_local_handles": schema.BoolAttribute{
				Description: "Indicates whether to use thread-local database handles to reduce contention in the backend.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"db_log_file_max": schema.StringAttribute{
				Description: "Specifies the maximum size for a database log file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"db_logging_level": schema.StringAttribute{
				Description: "Specifies the log level that should be used by the database when it is writing information into the je.info file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"je_property": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `changelog`: Specifies the database and environment properties for the Berkeley DB Java Edition database for this changelog backend. When the `type` attribute is set to `local-db`: Specifies the database and environment properties for the Berkeley DB Java Edition database serving the data for this backend.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `changelog`: Specifies the database and environment properties for the Berkeley DB Java Edition database for this changelog backend.\n  - `local-db`: Specifies the database and environment properties for the Berkeley DB Java Edition database serving the data for this backend.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"changelog_write_batch_size": schema.Int64Attribute{
				Description: "Specifies the number of changelog entries written in a single database transaction.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"default_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used for any database for which the cache mode is not explicitly specified. This includes the id2entry database, which stores encoded entries, and all attribute indexes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"id2entry_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the id2entry database, which provides a mapping between entry IDs and entry contents. Consider configuring uncached entries or uncached attributes in lieu of changing from the \"cache-keys-and-values\" default value.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"dn2id_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the dn2id database, which provides a mapping between normalized entry DNs and the corresponding entry IDs.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"id2children_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the id2children database, which provides a mapping between the entry ID of a particular entry and the entry IDs of all of its immediate children. This index may be used when performing searches with a single-level scope if the search filter cannot be resolved to a small enough candidate list. The size of this database directly depends on the number of entries that have children.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"id2subtree_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the id2subtree database, which provides a mapping between the entry ID of a particular entry and the entry IDs of all of its children to any depth. This index may be used when performing searches with a whole-subtree or subordinate-subtree scope if the search filter cannot be resolved to a small enough candidate list. The size of this database directly depends on the number of entries that have children.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"dn2uri_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the dn2uri database, which provides a mapping between a normalized entry DN and a set of referral URLs contained in the associated smart referral entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"simple_paged_results_id_set_cache_duration": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 10.1.0.0+. Specifies the length of time to cache the candidate ID set used for indexed search operations including the simple paged results control.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"prime_method": schema.SetAttribute{
				Description: "Specifies the method that should be used to prime caches with data for this backend.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"prime_thread_count": schema.Int64Attribute{
				Description: "Specifies the number of threads to use when priming. At present, this applies only to the preload and cursor-across-indexes prime methods.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"prime_time_limit": schema.StringAttribute{
				Description: "Specifies the maximum length of time that the backend prime should be allowed to run. A duration of zero seconds indicates that there should not be a time limit.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"prime_all_indexes": schema.BoolAttribute{
				Description: "Indicates whether to prime all indexes associated with this backend, or to only prime the specified set of indexes (as configured with the system-index-to-prime property for the system indexes, and the prime-index property in the attribute index definition for attribute indexes).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"system_index_to_prime": schema.SetAttribute{
				Description: "Specifies which system index(es) should be primed when the backend is initialized.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"system_index_to_prime_internal_nodes_only": schema.SetAttribute{
				Description: "Specifies the system index(es) for which internal database nodes only (i.e., the database keys but not values) should be primed when the backend is initialized.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"background_prime": schema.BoolAttribute{
				Description: "Indicates whether to attempt to perform the prime using a background thread if possible. If background priming is enabled, then the Directory Server may be allowed to accept client connections and process requests while the prime is in progress.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"index_entry_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that are allowed to match a given index key before that particular index key is no longer maintained.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"composite_index_entry_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that are allowed to match a given composite index key before that particular composite index key is no longer maintained.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"id2children_index_entry_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entry IDs to maintain for each entry in the id2children system index (which keeps track of the immediate children for an entry, to assist in otherwise unindexed searches with a single-level scope). A value of 0 means there is no limit, however this could have a big impact on database size on disk and on server performance.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"id2subtree_index_entry_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entry IDs to maintain for each entry in the id2subtree system index (which keeps track of all descendants below an entry, to assist in otherwise unindexed searches with a whole-subtree or subordinate subtree scope). A value of 0 means there is no limit, however this could have a big impact on database size on disk and on server performance.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"import_temp_directory": schema.StringAttribute{
				Description: "Specifies the location of the directory that is used to hold temporary information during the index post-processing phase of an LDIF import.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"import_thread_count": schema.Int64Attribute{
				Description: "Specifies the number of threads to use for concurrent processing during an LDIF import.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"export_thread_count": schema.Int64Attribute{
				Description: "Specifies the number of threads to use for concurrently retrieving and encoding entries during an LDIF export.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"db_import_cache_percent": schema.Int64Attribute{
				Description: "The percentage of JVM memory to allocate to the database cache during import operations.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"db_txn_write_no_sync": schema.BoolAttribute{
				Description: "Indicates whether the database should synchronously flush data as it is written to disk.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"deadlock_retry_limit": schema.Int64Attribute{
				Description: "Specifies the number of times that the server should retry an attempted operation in the backend if a deadlock results from two concurrent requests that interfere with each other in a conflicting manner.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"external_txn_default_backend_lock_behavior": schema.StringAttribute{
				Description: "Specifies the default behavior that should be exhibited by external transactions (e.g., an LDAP transaction or an atomic multi-update operation) with regard to acquiring an exclusive lock in this backend.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"single_writer_lock_behavior": schema.StringAttribute{
				Description: "Specifies the condition under which to acquire a single-writer lock to ensure that the associated operation will be the only write in progress at the time the lock is held. The single-writer lock can help avoid problems that result from database lock conflicts that arise between two write operations being processed at the same time in the same backend. This will not have any effect on the read operations processed while the write is in progress.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"subtree_modify_dn_size_limit": schema.Int64Attribute{
				Description: "Supported in PingDirectory product version 10.1.0.0+. Specifies the maximum number of entries that may exist below an entry targeted by a modify DN operation. This includes both direct and indirect subordinates (to any depth), although the entry at the top of the subtree (the one directly targeted by the modify DN operation) is not included in this count.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"subtree_delete_size_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that may be deleted from the backend when using the subtree delete control.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"num_recent_changes": schema.Int64Attribute{
				Description: "Specifies the number of recent LDAP entry changes per replica for which the backend keeps a record to allow replication to recover in the event that the server is abruptly terminated. Increasing this value can lead to an increased peak server modification rate as well as increased replication throughput.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"offline_process_database_open_timeout": schema.StringAttribute{
				Description: "Specifies a timeout duration which will be used for opening the database environment by an offline process, such as export-ldif.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"changelog_purge_batch_size": schema.Int64Attribute{
				Description: "Specifies the number of changelog entries purged in a single database transaction.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"changelog_write_queue_capacity": schema.Int64Attribute{
				Description: "Specifies the capacity of the changelog write queue in number of changes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"index_include_attribute": schema.SetAttribute{
				Description: "Specifies which attribute types are to be specifically included in the set of attribute indexes maintained on the changelog. If this property does not have any values then no attribute types are indexed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"index_exclude_attribute": schema.SetAttribute{
				Description: "Specifies which attribute types are to be specifically excluded from the set of attribute indexes maintained on the changelog. This property is useful when the index-include-attribute property contains one of the special values \"*\" and \"+\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"changelog_maximum_age": schema.StringAttribute{
				Description: "Changes are guaranteed to be maintained in the changelog database for at least this duration. Setting target-database-size can allow additional changes to be maintained up to the configured size on disk.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"target_database_size": schema.StringAttribute{
				Description: "The changelog database is allowed to grow up to this size on disk even if changes are older than the configured changelog-maximum-age.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"changelog_entry_include_base_dn": schema.SetAttribute{
				Description: "The base DNs for branches in the data for which to record changes in the changelog.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"changelog_entry_exclude_base_dn": schema.SetAttribute{
				Description: "The base DNs for branches in the data for which no changelog records should be generated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"changelog_entry_include_filter": schema.SetAttribute{
				Description: "A filter that indicates which changelog entries should actually be stored in the changelog. Note that this filter is evaluated against the changelog entry itself and not against the entry that was the target of the change referenced by the changelog entry. This filter may target any attributes that appear in changelog entries with the exception of the changeNumber and entry-size-bytes attributes, since they will not be known at the time of the filter evaluation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"changelog_entry_exclude_filter": schema.SetAttribute{
				Description: "A filter that indicates which changelog entries should be excluded from the changelog. Note that this filter is evaluated against the changelog entry itself and not against the entry that was the target of the change referenced by the changelog entry. This filter may target any attributes that appear in changelog entries with the exception of the changeNumber and entry-size-bytes attributes, since they will not be known at the time of the filter evaluation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"changelog_include_attribute": schema.SetAttribute{
				Description: "Specifies which attribute types will be included in a changelog entry for ADD and MODIFY operations.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"changelog_exclude_attribute": schema.SetAttribute{
				Description: "Specifies a set of attribute types that should be excluded in a changelog entry for ADD and MODIFY operations.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"changelog_deleted_entry_include_attribute": schema.SetAttribute{
				Description: "Specifies a set of attribute types that should be included in a changelog entry for DELETE operations.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"changelog_deleted_entry_exclude_attribute": schema.SetAttribute{
				Description: "Specifies a set of attribute types that should be excluded from a changelog entry for DELETE operations.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"changelog_include_key_attribute": schema.SetAttribute{
				Description: "Specifies which attribute types will be included in a changelog entry on every change.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"changelog_max_before_after_values": schema.Int64Attribute{
				Description: "This controls whether all attribute values for a modified attribute (even those values that have not changed) will be included in the changelog entry. If the number of attribute values does not exceed this limit, then all values for the modified attribute will be included in the changelog entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"write_lastmod_attributes": schema.BoolAttribute{
				Description: "Specifies whether values of creatorsName, createTimestamp, modifiersName and modifyTimestamp attributes will be written to changelog entries.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"use_reversible_form": schema.BoolAttribute{
				Description: "Specifies whether the changelog should provide enough information to be able to revert the changes if desired.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_virtual_attributes": schema.SetAttribute{
				Description: "Specifies the changelog entry elements (if any) in which virtual attributes should be included.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"apply_access_controls_to_changelog_entry_contents": schema.BoolAttribute{
				Description: "Indicates whether the contents of changelog entries should be subject to access control and sensitive attribute evaluation such that the contents of attributes like changes, deletedEntryAttrs, ds-changelog-entry-key-attr-values, ds-changelog-before-values, and ds-changelog-after-values may be altered based on attributes the user can see in the target entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"report_excluded_changelog_attributes": schema.StringAttribute{
				Description: "Indicates whether changelog entries that have been altered by applying access controls should include additional information about any attributes that may have been removed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"soft_delete_entry_included_operation": schema.SetAttribute{
				Description: "Specifies which operations performed on soft-deleted entries will appear in the changelog.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"is_private_backend": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to `ldif`: Indicates whether the backend should be considered a private backend, which indicates that it is used for storing operational data rather than user-defined information. When the `type` attribute is set to `local-db`: Indicates whether this backend should be considered a private backend in the server. Private backends are meant for storing server-internal information and should not be used for user or application data.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ldif`: Indicates whether the backend should be considered a private backend, which indicates that it is used for storing operational data rather than user-defined information.\n  - `local-db`: Indicates whether this backend should be considered a private backend in the server. Private backends are meant for storing server-internal information and should not be used for user or application data.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"ldif_file": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`alert`, `alarm`]: Specifies the path to the LDIF file that serves as the backing file for this backend. When the `type` attribute is set to `ldif`: Specifies the path to the LDIF file containing the data for this backend.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`alert`, `alarm`]: Specifies the path to the LDIF file that serves as the backing file for this backend.\n  - `ldif`: Specifies the path to the LDIF file containing the data for this backend.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"trust_store_file": schema.StringAttribute{
				Description: "Specifies the path to the file that stores the trust information.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"trust_store_type": schema.StringAttribute{
				Description: "Specifies the format for the data in the key store file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"trust_store_pin": schema.StringAttribute{
				Description: "Specifies the clear-text PIN needed to access the Trust Store Backend.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"trust_store_pin_file": schema.StringAttribute{
				Description: "Specifies the path to the text file whose only contents should be a single line containing the clear-text PIN needed to access the Trust Store Backend.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"trust_store_pin_passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider to use to obtain the clear-text PIN needed to access the Trust Store Backend.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"base_dn": schema.SetAttribute{
				Description: "Specifies the base DN(s) for the data that the backend handles.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"writability_mode": schema.StringAttribute{
				Description: "Specifies the behavior that the backend should use when processing write operations.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"backup_directory": schema.SetAttribute{
				Description: "Specifies the path to a backup directory containing one or more backups for a particular backend.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"schema_entry_dn": schema.SetAttribute{
				Description: "Defines the base DNs of the subtrees in which the schema information is published in addition to the value included in the base-dn property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"show_all_attributes": schema.BoolAttribute{
				Description: "Indicates whether to treat all attributes in the schema entry as if they were user attributes regardless of their configuration.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"read_only_schema_file": schema.SetAttribute{
				Description: "Specifies the name of a file (which must exist in the config/schema directory) containing schema elements that should be considered read-only. Any schema definitions contained in read-only files cannot be altered by external clients.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Backend",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the backend is enabled in the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"set_degraded_alert_when_disabled": schema.BoolAttribute{
				Description: "Determines whether the Directory Server enters a DEGRADED state (and sends a corresponding alert) when this Backend is disabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"return_unavailable_when_disabled": schema.BoolAttribute{
				Description: "Determines whether any LDAP operation that would use this Backend is to return UNAVAILABLE when this Backend is disabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"backup_file_permissions": schema.StringAttribute{
				Description: "Specifies the permissions that should be applied to files and directories created by a backup of the backend.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"notification_manager": schema.StringAttribute{
				Description: "Specifies a notification manager for changes resulting from operations processed through this Backend",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a SchemaBackendResponse object into the model struct
func readSchemaBackendResponseDataSource(ctx context.Context, r *client.SchemaBackendResponse, state *backendDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("schema")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.SchemaEntryDN = internaltypes.GetStringSet(r.SchemaEntryDN)
	state.ShowAllAttributes = types.BoolValue(r.ShowAllAttributes)
	state.ReadOnlySchemaFile = internaltypes.GetStringSet(r.ReadOnlySchemaFile)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, false)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, false)
}

// Read a BackupBackendResponse object into the model struct
func readBackupBackendResponseDataSource(ctx context.Context, r *client.BackupBackendResponse, state *backendDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("backup")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.BackupDirectory = internaltypes.GetStringSet(r.BackupDirectory)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, false)
}

// Read a EncryptionSettingsBackendResponse object into the model struct
func readEncryptionSettingsBackendResponseDataSource(ctx context.Context, r *client.EncryptionSettingsBackendResponse, state *backendDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("encryption-settings")
	state.Id = types.StringValue(r.Id)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.BackendID = types.StringValue(r.BackendID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, false)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, false)
}

// Read a LdifBackendResponse object into the model struct
func readLdifBackendResponseDataSource(ctx context.Context, r *client.LdifBackendResponse, state *backendDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldif")
	state.Id = types.StringValue(r.Id)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.IsPrivateBackend = internaltypes.BoolTypeOrNil(r.IsPrivateBackend)
	state.LdifFile = types.StringValue(r.LdifFile)
	state.BackendID = types.StringValue(r.BackendID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, false)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, false)
}

// Read a TrustStoreBackendResponse object into the model struct
func readTrustStoreBackendResponseDataSource(ctx context.Context, r *client.TrustStoreBackendResponse, state *backendDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("trust-store")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.TrustStoreFile = types.StringValue(r.TrustStoreFile)
	state.TrustStoreType = internaltypes.StringTypeOrNil(r.TrustStoreType, false)
	state.TrustStorePinFile = internaltypes.StringTypeOrNil(r.TrustStorePinFile, false)
	state.TrustStorePinPassphraseProvider = internaltypes.StringTypeOrNil(r.TrustStorePinPassphraseProvider, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, false)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, false)
}

// Read a CustomBackendResponse object into the model struct
func readCustomBackendResponseDataSource(ctx context.Context, r *client.CustomBackendResponse, state *backendDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("custom")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, false)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, false)
}

// Read a ChangelogBackendResponse object into the model struct
func readChangelogBackendResponseDataSource(ctx context.Context, r *client.ChangelogBackendResponse, state *backendDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("changelog")
	state.Id = types.StringValue(r.Id)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.DbDirectory = internaltypes.StringTypeOrNil(r.DbDirectory, false)
	state.DbDirectoryPermissions = internaltypes.StringTypeOrNil(r.DbDirectoryPermissions, false)
	state.DbCachePercent = internaltypes.Int64TypeOrNil(r.DbCachePercent)
	state.JeProperty = internaltypes.GetStringSet(r.JeProperty)
	state.ChangelogWriteBatchSize = internaltypes.Int64TypeOrNil(r.ChangelogWriteBatchSize)
	state.ChangelogPurgeBatchSize = internaltypes.Int64TypeOrNil(r.ChangelogPurgeBatchSize)
	state.ChangelogWriteQueueCapacity = internaltypes.Int64TypeOrNil(r.ChangelogWriteQueueCapacity)
	state.IndexIncludeAttribute = internaltypes.GetStringSet(r.IndexIncludeAttribute)
	state.IndexExcludeAttribute = internaltypes.GetStringSet(r.IndexExcludeAttribute)
	state.ChangelogMaximumAge = types.StringValue(r.ChangelogMaximumAge)
	state.TargetDatabaseSize = internaltypes.StringTypeOrNil(r.TargetDatabaseSize, false)
	state.ChangelogEntryIncludeBaseDN = internaltypes.GetStringSet(r.ChangelogEntryIncludeBaseDN)
	state.ChangelogEntryExcludeBaseDN = internaltypes.GetStringSet(r.ChangelogEntryExcludeBaseDN)
	state.ChangelogEntryIncludeFilter = internaltypes.GetStringSet(r.ChangelogEntryIncludeFilter)
	state.ChangelogEntryExcludeFilter = internaltypes.GetStringSet(r.ChangelogEntryExcludeFilter)
	state.ChangelogIncludeAttribute = internaltypes.GetStringSet(r.ChangelogIncludeAttribute)
	state.ChangelogExcludeAttribute = internaltypes.GetStringSet(r.ChangelogExcludeAttribute)
	state.ChangelogDeletedEntryIncludeAttribute = internaltypes.GetStringSet(r.ChangelogDeletedEntryIncludeAttribute)
	state.ChangelogDeletedEntryExcludeAttribute = internaltypes.GetStringSet(r.ChangelogDeletedEntryExcludeAttribute)
	state.ChangelogIncludeKeyAttribute = internaltypes.GetStringSet(r.ChangelogIncludeKeyAttribute)
	state.ChangelogMaxBeforeAfterValues = internaltypes.Int64TypeOrNil(r.ChangelogMaxBeforeAfterValues)
	state.WriteLastmodAttributes = internaltypes.BoolTypeOrNil(r.WriteLastmodAttributes)
	state.UseReversibleForm = internaltypes.BoolTypeOrNil(r.UseReversibleForm)
	state.IncludeVirtualAttributes = internaltypes.GetStringSet(
		client.StringSliceEnumbackendIncludeVirtualAttributesProp(r.IncludeVirtualAttributes))
	state.ApplyAccessControlsToChangelogEntryContents = internaltypes.BoolTypeOrNil(r.ApplyAccessControlsToChangelogEntryContents)
	state.ReportExcludedChangelogAttributes = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendReportExcludedChangelogAttributesProp(r.ReportExcludedChangelogAttributes), false)
	state.SoftDeleteEntryIncludedOperation = internaltypes.GetStringSet(
		client.StringSliceEnumbackendSoftDeleteEntryIncludedOperationProp(r.SoftDeleteEntryIncludedOperation))
	state.BackendID = types.StringValue(r.BackendID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, false)
}

// Read a MonitorBackendResponse object into the model struct
func readMonitorBackendResponseDataSource(ctx context.Context, r *client.MonitorBackendResponse, state *backendDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("monitor")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, false)
}

// Read a LocalDbBackendResponse object into the model struct
func readLocalDbBackendResponseDataSource(ctx context.Context, r *client.LocalDbBackendResponse, state *backendDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("local-db")
	state.Id = types.StringValue(r.Id)
	state.UncachedId2entryCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendUncachedId2entryCacheModeProp(r.UncachedId2entryCacheMode), false)
	state.UncachedAttributeCriteria = internaltypes.StringTypeOrNil(r.UncachedAttributeCriteria, false)
	state.UncachedEntryCriteria = internaltypes.StringTypeOrNil(r.UncachedEntryCriteria, false)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.SetDegradedAlertForUntrustedIndex = internaltypes.BoolTypeOrNil(r.SetDegradedAlertForUntrustedIndex)
	state.ReturnUnavailableForUntrustedIndex = internaltypes.BoolTypeOrNil(r.ReturnUnavailableForUntrustedIndex)
	state.ProcessFiltersWithUndefinedAttributeTypes = internaltypes.BoolTypeOrNil(r.ProcessFiltersWithUndefinedAttributeTypes)
	state.IsPrivateBackend = internaltypes.BoolTypeOrNil(r.IsPrivateBackend)
	state.DbDirectory = types.StringValue(r.DbDirectory)
	state.DbDirectoryPermissions = internaltypes.StringTypeOrNil(r.DbDirectoryPermissions, false)
	state.CompactCommonParentDN = internaltypes.GetStringSet(r.CompactCommonParentDN)
	state.CompressEntries = internaltypes.BoolTypeOrNil(r.CompressEntries)
	state.HashEntries = internaltypes.BoolTypeOrNil(r.HashEntries)
	state.DbNumCleanerThreads = internaltypes.Int64TypeOrNil(r.DbNumCleanerThreads)
	state.DbCleanerMinUtilization = internaltypes.Int64TypeOrNil(r.DbCleanerMinUtilization)
	state.DbEvictorCriticalPercentage = internaltypes.Int64TypeOrNil(r.DbEvictorCriticalPercentage)
	state.DbCheckpointerWakeupInterval = internaltypes.StringTypeOrNil(r.DbCheckpointerWakeupInterval, false)
	state.DbBackgroundSyncInterval = internaltypes.StringTypeOrNil(r.DbBackgroundSyncInterval, false)
	state.DbUseThreadLocalHandles = internaltypes.BoolTypeOrNil(r.DbUseThreadLocalHandles)
	state.DbLogFileMax = internaltypes.StringTypeOrNil(r.DbLogFileMax, false)
	state.DbLoggingLevel = internaltypes.StringTypeOrNil(r.DbLoggingLevel, false)
	state.JeProperty = internaltypes.GetStringSet(r.JeProperty)
	state.DbCachePercent = internaltypes.Int64TypeOrNil(r.DbCachePercent)
	state.DefaultCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendDefaultCacheModeProp(r.DefaultCacheMode), false)
	state.Id2entryCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendId2entryCacheModeProp(r.Id2entryCacheMode), false)
	state.Dn2idCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendDn2idCacheModeProp(r.Dn2idCacheMode), false)
	state.Id2childrenCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendId2childrenCacheModeProp(r.Id2childrenCacheMode), false)
	state.Id2subtreeCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendId2subtreeCacheModeProp(r.Id2subtreeCacheMode), false)
	state.Dn2uriCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendDn2uriCacheModeProp(r.Dn2uriCacheMode), false)
	state.SimplePagedResultsIDSetCacheDuration = internaltypes.StringTypeOrNil(r.SimplePagedResultsIDSetCacheDuration, false)
	state.PrimeMethod = internaltypes.GetStringSet(
		client.StringSliceEnumbackendPrimeMethodProp(r.PrimeMethod))
	state.PrimeThreadCount = internaltypes.Int64TypeOrNil(r.PrimeThreadCount)
	state.PrimeTimeLimit = internaltypes.StringTypeOrNil(r.PrimeTimeLimit, false)
	state.PrimeAllIndexes = internaltypes.BoolTypeOrNil(r.PrimeAllIndexes)
	state.SystemIndexToPrime = internaltypes.GetStringSet(
		client.StringSliceEnumbackendSystemIndexToPrimeProp(r.SystemIndexToPrime))
	state.SystemIndexToPrimeInternalNodesOnly = internaltypes.GetStringSet(
		client.StringSliceEnumbackendSystemIndexToPrimeInternalNodesOnlyProp(r.SystemIndexToPrimeInternalNodesOnly))
	state.BackgroundPrime = internaltypes.BoolTypeOrNil(r.BackgroundPrime)
	state.IndexEntryLimit = internaltypes.Int64TypeOrNil(r.IndexEntryLimit)
	state.CompositeIndexEntryLimit = internaltypes.Int64TypeOrNil(r.CompositeIndexEntryLimit)
	state.Id2childrenIndexEntryLimit = internaltypes.Int64TypeOrNil(r.Id2childrenIndexEntryLimit)
	state.Id2subtreeIndexEntryLimit = internaltypes.Int64TypeOrNil(r.Id2subtreeIndexEntryLimit)
	state.ImportTempDirectory = types.StringValue(r.ImportTempDirectory)
	state.ImportThreadCount = internaltypes.Int64TypeOrNil(r.ImportThreadCount)
	state.ExportThreadCount = internaltypes.Int64TypeOrNil(r.ExportThreadCount)
	state.DbImportCachePercent = internaltypes.Int64TypeOrNil(r.DbImportCachePercent)
	state.DbTxnWriteNoSync = internaltypes.BoolTypeOrNil(r.DbTxnWriteNoSync)
	state.DeadlockRetryLimit = internaltypes.Int64TypeOrNil(r.DeadlockRetryLimit)
	state.ExternalTxnDefaultBackendLockBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendExternalTxnDefaultBackendLockBehaviorProp(r.ExternalTxnDefaultBackendLockBehavior), false)
	state.SingleWriterLockBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendSingleWriterLockBehaviorProp(r.SingleWriterLockBehavior), false)
	state.SubtreeModifyDNSizeLimit = internaltypes.Int64TypeOrNil(r.SubtreeModifyDNSizeLimit)
	state.SubtreeDeleteSizeLimit = internaltypes.Int64TypeOrNil(r.SubtreeDeleteSizeLimit)
	state.NumRecentChanges = internaltypes.Int64TypeOrNil(r.NumRecentChanges)
	state.OfflineProcessDatabaseOpenTimeout = internaltypes.StringTypeOrNil(r.OfflineProcessDatabaseOpenTimeout, false)
	state.BackendID = types.StringValue(r.BackendID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, false)
}

// Read a ConfigFileHandlerBackendResponse object into the model struct
func readConfigFileHandlerBackendResponseDataSource(ctx context.Context, r *client.ConfigFileHandlerBackendResponse, state *backendDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("config-file-handler")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.InsignificantConfigArchiveAttribute = internaltypes.GetStringSet(r.InsignificantConfigArchiveAttribute)
	state.InsignificantConfigArchiveBaseDN = internaltypes.GetStringSet(r.InsignificantConfigArchiveBaseDN)
	state.MaintainConfigArchive = internaltypes.BoolTypeOrNil(r.MaintainConfigArchive)
	state.MaxConfigArchiveCount = internaltypes.Int64TypeOrNil(r.MaxConfigArchiveCount)
	state.MirroredSubtreePeerPollingInterval = internaltypes.StringTypeOrNil(r.MirroredSubtreePeerPollingInterval, false)
	state.MirroredSubtreeEntryUpdateTimeout = internaltypes.StringTypeOrNil(r.MirroredSubtreeEntryUpdateTimeout, false)
	state.MirroredSubtreeSearchTimeout = internaltypes.StringTypeOrNil(r.MirroredSubtreeSearchTimeout, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, false)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, false)
}

// Read a TaskBackendResponse object into the model struct
func readTaskBackendResponseDataSource(ctx context.Context, r *client.TaskBackendResponse, state *backendDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("task")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.TaskBackingFile = types.StringValue(r.TaskBackingFile)
	state.MaximumInitialTaskLogMessagesToRetain = internaltypes.Int64TypeOrNil(r.MaximumInitialTaskLogMessagesToRetain)
	state.MaximumFinalTaskLogMessagesToRetain = internaltypes.Int64TypeOrNil(r.MaximumFinalTaskLogMessagesToRetain)
	state.TaskRetentionTime = internaltypes.StringTypeOrNil(r.TaskRetentionTime, false)
	state.NotificationSenderAddress = internaltypes.StringTypeOrNil(r.NotificationSenderAddress, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, false)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, false)
}

// Read a AlertBackendResponse object into the model struct
func readAlertBackendResponseDataSource(ctx context.Context, r *client.AlertBackendResponse, state *backendDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("alert")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.LdifFile = types.StringValue(r.LdifFile)
	state.AlertRetentionTime = types.StringValue(r.AlertRetentionTime)
	state.MaxAlerts = internaltypes.Int64TypeOrNil(r.MaxAlerts)
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumbackendDisabledAlertTypeProp(r.DisabledAlertType))
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, false)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, false)
}

// Read a AlarmBackendResponse object into the model struct
func readAlarmBackendResponseDataSource(ctx context.Context, r *client.AlarmBackendResponse, state *backendDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("alarm")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.LdifFile = types.StringValue(r.LdifFile)
	state.AlarmRetentionTime = types.StringValue(r.AlarmRetentionTime)
	state.MaxAlarms = internaltypes.Int64TypeOrNil(r.MaxAlarms)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, false)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, false)
}

// Read a MetricsBackendResponse object into the model struct
func readMetricsBackendResponseDataSource(ctx context.Context, r *client.MetricsBackendResponse, state *backendDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("metrics")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.StorageDir = types.StringValue(r.StorageDir)
	state.MetricsDir = types.StringValue(r.MetricsDir)
	state.SampleFlushInterval = internaltypes.StringTypeOrNil(r.SampleFlushInterval, false)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, false)
}

// Read resource information
func (r *backendDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state backendDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.BackendAPI.GetBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.BackendID.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Backend", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.SchemaBackendResponse != nil {
		readSchemaBackendResponseDataSource(ctx, readResponse.SchemaBackendResponse, &state, &resp.Diagnostics)
	}
	if readResponse.BackupBackendResponse != nil {
		readBackupBackendResponseDataSource(ctx, readResponse.BackupBackendResponse, &state, &resp.Diagnostics)
	}
	if readResponse.EncryptionSettingsBackendResponse != nil {
		readEncryptionSettingsBackendResponseDataSource(ctx, readResponse.EncryptionSettingsBackendResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LdifBackendResponse != nil {
		readLdifBackendResponseDataSource(ctx, readResponse.LdifBackendResponse, &state, &resp.Diagnostics)
	}
	if readResponse.TrustStoreBackendResponse != nil {
		readTrustStoreBackendResponseDataSource(ctx, readResponse.TrustStoreBackendResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CustomBackendResponse != nil {
		readCustomBackendResponseDataSource(ctx, readResponse.CustomBackendResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ChangelogBackendResponse != nil {
		readChangelogBackendResponseDataSource(ctx, readResponse.ChangelogBackendResponse, &state, &resp.Diagnostics)
	}
	if readResponse.MonitorBackendResponse != nil {
		readMonitorBackendResponseDataSource(ctx, readResponse.MonitorBackendResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LocalDbBackendResponse != nil {
		readLocalDbBackendResponseDataSource(ctx, readResponse.LocalDbBackendResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ConfigFileHandlerBackendResponse != nil {
		readConfigFileHandlerBackendResponseDataSource(ctx, readResponse.ConfigFileHandlerBackendResponse, &state, &resp.Diagnostics)
	}
	if readResponse.TaskBackendResponse != nil {
		readTaskBackendResponseDataSource(ctx, readResponse.TaskBackendResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AlertBackendResponse != nil {
		readAlertBackendResponseDataSource(ctx, readResponse.AlertBackendResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AlarmBackendResponse != nil {
		readAlarmBackendResponseDataSource(ctx, readResponse.AlarmBackendResponse, &state, &resp.Diagnostics)
	}
	if readResponse.MetricsBackendResponse != nil {
		readMetricsBackendResponseDataSource(ctx, readResponse.MetricsBackendResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
