package backend

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &backendResource{}
	_ resource.ResourceWithConfigure   = &backendResource{}
	_ resource.ResourceWithImportState = &backendResource{}
	_ resource.Resource                = &defaultBackendResource{}
	_ resource.ResourceWithConfigure   = &defaultBackendResource{}
	_ resource.ResourceWithImportState = &defaultBackendResource{}
)

// Create a Backend resource
func NewBackendResource() resource.Resource {
	return &backendResource{}
}

func NewDefaultBackendResource() resource.Resource {
	return &defaultBackendResource{}
}

// backendResource is the resource implementation.
type backendResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultBackendResource is the resource implementation.
type defaultBackendResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *backendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backend"
}

func (r *defaultBackendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_backend"
}

// Configure adds the provider configured client to the resource.
func (r *backendResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultBackendResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type backendResourceModel struct {
	Id                                        types.String `tfsdk:"id"`
	Notifications                             types.Set    `tfsdk:"notifications"`
	RequiredActions                           types.Set    `tfsdk:"required_actions"`
	Type                                      types.String `tfsdk:"type"`
	UncachedId2entryCacheMode                 types.String `tfsdk:"uncached_id2entry_cache_mode"`
	UncachedAttributeCriteria                 types.String `tfsdk:"uncached_attribute_criteria"`
	UncachedEntryCriteria                     types.String `tfsdk:"uncached_entry_criteria"`
	BackendID                                 types.String `tfsdk:"backend_id"`
	SetDegradedAlertForUntrustedIndex         types.Bool   `tfsdk:"set_degraded_alert_for_untrusted_index"`
	ReturnUnavailableForUntrustedIndex        types.Bool   `tfsdk:"return_unavailable_for_untrusted_index"`
	ProcessFiltersWithUndefinedAttributeTypes types.Bool   `tfsdk:"process_filters_with_undefined_attribute_types"`
	DbDirectory                               types.String `tfsdk:"db_directory"`
	DbDirectoryPermissions                    types.String `tfsdk:"db_directory_permissions"`
	DbCachePercent                            types.Int64  `tfsdk:"db_cache_percent"`
	CompactCommonParentDN                     types.Set    `tfsdk:"compact_common_parent_dn"`
	CompressEntries                           types.Bool   `tfsdk:"compress_entries"`
	HashEntries                               types.Bool   `tfsdk:"hash_entries"`
	DbNumCleanerThreads                       types.Int64  `tfsdk:"db_num_cleaner_threads"`
	DbCleanerMinUtilization                   types.Int64  `tfsdk:"db_cleaner_min_utilization"`
	DbEvictorCriticalPercentage               types.Int64  `tfsdk:"db_evictor_critical_percentage"`
	DbCheckpointerWakeupInterval              types.String `tfsdk:"db_checkpointer_wakeup_interval"`
	DbBackgroundSyncInterval                  types.String `tfsdk:"db_background_sync_interval"`
	DbUseThreadLocalHandles                   types.Bool   `tfsdk:"db_use_thread_local_handles"`
	DbLogFileMax                              types.String `tfsdk:"db_log_file_max"`
	DbLoggingLevel                            types.String `tfsdk:"db_logging_level"`
	JeProperty                                types.Set    `tfsdk:"je_property"`
	DefaultCacheMode                          types.String `tfsdk:"default_cache_mode"`
	Id2entryCacheMode                         types.String `tfsdk:"id2entry_cache_mode"`
	Dn2idCacheMode                            types.String `tfsdk:"dn2id_cache_mode"`
	Id2childrenCacheMode                      types.String `tfsdk:"id2children_cache_mode"`
	Id2subtreeCacheMode                       types.String `tfsdk:"id2subtree_cache_mode"`
	Dn2uriCacheMode                           types.String `tfsdk:"dn2uri_cache_mode"`
	PrimeMethod                               types.Set    `tfsdk:"prime_method"`
	PrimeThreadCount                          types.Int64  `tfsdk:"prime_thread_count"`
	PrimeTimeLimit                            types.String `tfsdk:"prime_time_limit"`
	PrimeAllIndexes                           types.Bool   `tfsdk:"prime_all_indexes"`
	SystemIndexToPrime                        types.Set    `tfsdk:"system_index_to_prime"`
	SystemIndexToPrimeInternalNodesOnly       types.Set    `tfsdk:"system_index_to_prime_internal_nodes_only"`
	BackgroundPrime                           types.Bool   `tfsdk:"background_prime"`
	IndexEntryLimit                           types.Int64  `tfsdk:"index_entry_limit"`
	CompositeIndexEntryLimit                  types.Int64  `tfsdk:"composite_index_entry_limit"`
	Id2childrenIndexEntryLimit                types.Int64  `tfsdk:"id2children_index_entry_limit"`
	Id2subtreeIndexEntryLimit                 types.Int64  `tfsdk:"id2subtree_index_entry_limit"`
	ImportTempDirectory                       types.String `tfsdk:"import_temp_directory"`
	ImportThreadCount                         types.Int64  `tfsdk:"import_thread_count"`
	ExportThreadCount                         types.Int64  `tfsdk:"export_thread_count"`
	DbImportCachePercent                      types.Int64  `tfsdk:"db_import_cache_percent"`
	DbTxnWriteNoSync                          types.Bool   `tfsdk:"db_txn_write_no_sync"`
	DeadlockRetryLimit                        types.Int64  `tfsdk:"deadlock_retry_limit"`
	ExternalTxnDefaultBackendLockBehavior     types.String `tfsdk:"external_txn_default_backend_lock_behavior"`
	SingleWriterLockBehavior                  types.String `tfsdk:"single_writer_lock_behavior"`
	SubtreeDeleteSizeLimit                    types.Int64  `tfsdk:"subtree_delete_size_limit"`
	NumRecentChanges                          types.Int64  `tfsdk:"num_recent_changes"`
	OfflineProcessDatabaseOpenTimeout         types.String `tfsdk:"offline_process_database_open_timeout"`
	IsPrivateBackend                          types.Bool   `tfsdk:"is_private_backend"`
	BaseDN                                    types.Set    `tfsdk:"base_dn"`
	WritabilityMode                           types.String `tfsdk:"writability_mode"`
	Description                               types.String `tfsdk:"description"`
	Enabled                                   types.Bool   `tfsdk:"enabled"`
	SetDegradedAlertWhenDisabled              types.Bool   `tfsdk:"set_degraded_alert_when_disabled"`
	ReturnUnavailableWhenDisabled             types.Bool   `tfsdk:"return_unavailable_when_disabled"`
	NotificationManager                       types.String `tfsdk:"notification_manager"`
}

type defaultBackendResourceModel struct {
	Id                                          types.String `tfsdk:"id"`
	Notifications                               types.Set    `tfsdk:"notifications"`
	RequiredActions                             types.Set    `tfsdk:"required_actions"`
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

// GetSchema defines the schema for the resource.
func (r *backendResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	backendSchema(ctx, req, resp, false)
}

func (r *defaultBackendResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	backendSchema(ctx, req, resp, true)
}

func backendSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Backend.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Backend resource. Options are ['schema', 'backup', 'encryption-settings', 'ldif', 'trust-store', 'custom', 'changelog', 'monitor', 'local-db', 'config-file-handler', 'task', 'alert', 'alarm', 'metrics']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"local-db"}...),
				},
			},
			"uncached_id2entry_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the uncached-id2entry database, which provides a way to store complete or partial encoded entries with a different (and presumably less memory-intensive) cache mode than records written to id2entry.",
				Optional:    true,
				Computed:    true,
			},
			"uncached_attribute_criteria": schema.StringAttribute{
				Description: "The criteria that will be used to identify attributes that should be written into the uncached-id2entry database rather than the id2entry database. This will only be used for entries in which the associated uncached-entry-criteria does not indicate that the entire entry should be uncached.",
				Optional:    true,
			},
			"uncached_entry_criteria": schema.StringAttribute{
				Description: "The criteria that will be used to identify entries that should be written into the uncached-id2entry database rather than the id2entry database.",
				Optional:    true,
			},
			"backend_id": schema.StringAttribute{
				Description: "Specifies a name to identify the associated backend.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"set_degraded_alert_for_untrusted_index": schema.BoolAttribute{
				Description: "Determines whether the Directory Server enters a DEGRADED state when this Local DB Backend has an index whose contents cannot be trusted.",
				Optional:    true,
				Computed:    true,
			},
			"return_unavailable_for_untrusted_index": schema.BoolAttribute{
				Description: "Determines whether the Directory Server returns UNAVAILABLE for any LDAP search operation in this Local DB Backend that would use an index whose contents cannot be trusted.",
				Optional:    true,
				Computed:    true,
			},
			"process_filters_with_undefined_attribute_types": schema.BoolAttribute{
				Description: "Determines whether the Directory Server should continue filter processing for LDAP search operations in this Local DB Backend that includes a search filter with an attribute that is not defined in the schema. This will only apply if check-schema is enabled in the global configuration.",
				Optional:    true,
				Computed:    true,
			},
			"db_directory": schema.StringAttribute{
				Description: "Specifies the path to the filesystem directory that is used to hold the Berkeley DB Java Edition database files containing the data for this backend. The files for this backend are stored in a sub-directory named after the backend-id.",
				Optional:    true,
				Computed:    true,
			},
			"db_directory_permissions": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `changelog`: Specifies the permissions that should be applied to the directory containing the backend database files and to directories and files created during backup of the backend. When the `type` attribute is set to `local-db`: Specifies the permissions that should be applied to the directory containing the backend database files and to directories and files created during backup or LDIF export of the backend.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `changelog`: Specifies the permissions that should be applied to the directory containing the backend database files and to directories and files created during backup of the backend.\n  - `local-db`: Specifies the permissions that should be applied to the directory containing the backend database files and to directories and files created during backup or LDIF export of the backend.",
				Optional:            true,
				Computed:            true,
			},
			"db_cache_percent": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to `changelog`: Specifies the percentage of JVM memory to allocate to the changelog database cache. When the `type` attribute is set to `local-db`: Specifies the percentage of JVM memory to allocate to the database cache.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `changelog`: Specifies the percentage of JVM memory to allocate to the changelog database cache.\n  - `local-db`: Specifies the percentage of JVM memory to allocate to the database cache.",
				Optional:            true,
				Computed:            true,
			},
			"compact_common_parent_dn": schema.SetAttribute{
				Description: "Provides a DN of an entry that may be the parent for a large number of entries in the backend. This may be used to help increase the space efficiency when encoding entries for storage.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"compress_entries": schema.BoolAttribute{
				Description: "Indicates whether the backend should attempt to compress entries before storing them in the database.",
				Optional:    true,
				Computed:    true,
			},
			"hash_entries": schema.BoolAttribute{
				Description: "Indicates whether to calculate and store a message digest of the entry contents along with the entry data, in order to provide a means of verifying the integrity of the entry data.",
				Optional:    true,
				Computed:    true,
			},
			"db_num_cleaner_threads": schema.Int64Attribute{
				Description: "Specifies the number of threads that the backend should maintain to keep the database log files at or near the desired utilization. A value of zero indicates that the number of cleaner threads should be automatically configured based on the number of available CPUs.",
				Optional:    true,
				Computed:    true,
			},
			"db_cleaner_min_utilization": schema.Int64Attribute{
				Description: "Specifies the minimum percentage of \"live\" data that the database cleaner attempts to keep in database log files.",
				Optional:    true,
				Computed:    true,
			},
			"db_evictor_critical_percentage": schema.Int64Attribute{
				Description: "Specifies the percentage over the configured maximum that the database cache is allowed to grow. It is recommended to set this value slightly above zero when the database is too large to fully cache in memory. In this case, a dedicated background evictor thread is used to perform evictions once the cache fills up reducing the possibility that server threads are blocked.",
				Optional:    true,
				Computed:    true,
			},
			"db_checkpointer_wakeup_interval": schema.StringAttribute{
				Description: "Specifies the maximum length of time that should pass between checkpoints.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"db_background_sync_interval": schema.StringAttribute{
				Description: "Specifies the interval to use when performing background synchronous writes in the database environment in order to smooth overall write performance and increase data durability. A value of \"0 s\" will disable background synchronous writes.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"db_use_thread_local_handles": schema.BoolAttribute{
				Description: "Indicates whether to use thread-local database handles to reduce contention in the backend.",
				Optional:    true,
				Computed:    true,
			},
			"db_log_file_max": schema.StringAttribute{
				Description: "Specifies the maximum size for a database log file.",
				Optional:    true,
				Computed:    true,
			},
			"db_logging_level": schema.StringAttribute{
				Description: "Specifies the log level that should be used by the database when it is writing information into the je.info file.",
				Optional:    true,
				Computed:    true,
			},
			"je_property": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `changelog`: Specifies the database and environment properties for the Berkeley DB Java Edition database for this changelog backend. When the `type` attribute is set to `local-db`: Specifies the database and environment properties for the Berkeley DB Java Edition database serving the data for this backend.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `changelog`: Specifies the database and environment properties for the Berkeley DB Java Edition database for this changelog backend.\n  - `local-db`: Specifies the database and environment properties for the Berkeley DB Java Edition database serving the data for this backend.",
				Optional:            true,
				Computed:            true,
				Default:             internaltypes.EmptySetDefault(types.StringType),
				ElementType:         types.StringType,
			},
			"default_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used for any database for which the cache mode is not explicitly specified. This includes the id2entry database, which stores encoded entries, and all attribute indexes.",
				Optional:    true,
				Computed:    true,
			},
			"id2entry_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the id2entry database, which provides a mapping between entry IDs and entry contents. Consider configuring uncached entries or uncached attributes in lieu of changing from the \"cache-keys-and-values\" default value.",
				Optional:    true,
			},
			"dn2id_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the dn2id database, which provides a mapping between normalized entry DNs and the corresponding entry IDs.",
				Optional:    true,
			},
			"id2children_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the id2children database, which provides a mapping between the entry ID of a particular entry and the entry IDs of all of its immediate children. This index may be used when performing searches with a single-level scope if the search filter cannot be resolved to a small enough candidate list. The size of this database directly depends on the number of entries that have children.",
				Optional:    true,
			},
			"id2subtree_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the id2subtree database, which provides a mapping between the entry ID of a particular entry and the entry IDs of all of its children to any depth. This index may be used when performing searches with a whole-subtree or subordinate-subtree scope if the search filter cannot be resolved to a small enough candidate list. The size of this database directly depends on the number of entries that have children.",
				Optional:    true,
			},
			"dn2uri_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the dn2uri database, which provides a mapping between a normalized entry DN and a set of referral URLs contained in the associated smart referral entry.",
				Optional:    true,
			},
			"prime_method": schema.SetAttribute{
				Description: "Specifies the method that should be used to prime caches with data for this backend.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"prime_thread_count": schema.Int64Attribute{
				Description: "Specifies the number of threads to use when priming. At present, this applies only to the preload and cursor-across-indexes prime methods.",
				Optional:    true,
				Computed:    true,
			},
			"prime_time_limit": schema.StringAttribute{
				Description: "Specifies the maximum length of time that the backend prime should be allowed to run. A duration of zero seconds indicates that there should not be a time limit.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"prime_all_indexes": schema.BoolAttribute{
				Description: "Indicates whether to prime all indexes associated with this backend, or to only prime the specified set of indexes (as configured with the system-index-to-prime property for the system indexes, and the prime-index property in the attribute index definition for attribute indexes).",
				Optional:    true,
				Computed:    true,
			},
			"system_index_to_prime": schema.SetAttribute{
				Description: "Specifies which system index(es) should be primed when the backend is initialized.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"system_index_to_prime_internal_nodes_only": schema.SetAttribute{
				Description: "Specifies the system index(es) for which internal database nodes only (i.e., the database keys but not values) should be primed when the backend is initialized.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"background_prime": schema.BoolAttribute{
				Description: "Indicates whether to attempt to perform the prime using a background thread if possible. If background priming is enabled, then the Directory Server may be allowed to accept client connections and process requests while the prime is in progress.",
				Optional:    true,
				Computed:    true,
			},
			"index_entry_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that are allowed to match a given index key before that particular index key is no longer maintained.",
				Optional:    true,
				Computed:    true,
			},
			"composite_index_entry_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that are allowed to match a given composite index key before that particular composite index key is no longer maintained.",
				Optional:    true,
				Computed:    true,
			},
			"id2children_index_entry_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entry IDs to maintain for each entry in the id2children system index (which keeps track of the immediate children for an entry, to assist in otherwise unindexed searches with a single-level scope). A value of 0 means there is no limit, however this could have a big impact on database size on disk and on server performance.",
				Optional:    true,
			},
			"id2subtree_index_entry_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entry IDs to maintain for each entry in the id2subtree system index (which keeps track of all descendants below an entry, to assist in otherwise unindexed searches with a whole-subtree or subordinate subtree scope). A value of 0 means there is no limit, however this could have a big impact on database size on disk and on server performance.",
				Optional:    true,
			},
			"import_temp_directory": schema.StringAttribute{
				Description: "Specifies the location of the directory that is used to hold temporary information during the index post-processing phase of an LDIF import.",
				Optional:    true,
				Computed:    true,
			},
			"import_thread_count": schema.Int64Attribute{
				Description: "Specifies the number of threads to use for concurrent processing during an LDIF import.",
				Optional:    true,
				Computed:    true,
			},
			"export_thread_count": schema.Int64Attribute{
				Description: "Specifies the number of threads to use for concurrently retrieving and encoding entries during an LDIF export.",
				Optional:    true,
				Computed:    true,
			},
			"db_import_cache_percent": schema.Int64Attribute{
				Description: "The percentage of JVM memory to allocate to the database cache during import operations.",
				Optional:    true,
				Computed:    true,
			},
			"db_txn_write_no_sync": schema.BoolAttribute{
				Description: "Indicates whether the database should synchronously flush data as it is written to disk.",
				Optional:    true,
				Computed:    true,
			},
			"deadlock_retry_limit": schema.Int64Attribute{
				Description: "Specifies the number of times that the server should retry an attempted operation in the backend if a deadlock results from two concurrent requests that interfere with each other in a conflicting manner.",
				Optional:    true,
				Computed:    true,
			},
			"external_txn_default_backend_lock_behavior": schema.StringAttribute{
				Description: "Specifies the default behavior that should be exhibited by external transactions (e.g., an LDAP transaction or an atomic multi-update operation) with regard to acquiring an exclusive lock in this backend.",
				Optional:    true,
				Computed:    true,
			},
			"single_writer_lock_behavior": schema.StringAttribute{
				Description: "Specifies the condition under which to acquire a single-writer lock to ensure that the associated operation will be the only write in progress at the time the lock is held. The single-writer lock can help avoid problems that result from database lock conflicts that arise between two write operations being processed at the same time in the same backend. This will not have any effect on the read operations processed while the write is in progress.",
				Optional:    true,
				Computed:    true,
			},
			"subtree_delete_size_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that may be deleted from the backend when using the subtree delete control.",
				Optional:    true,
				Computed:    true,
			},
			"num_recent_changes": schema.Int64Attribute{
				Description: "Specifies the number of recent LDAP entry changes per replica for which the backend keeps a record to allow replication to recover in the event that the server is abruptly terminated. Increasing this value can lead to an increased peak server modification rate as well as increased replication throughput.",
				Optional:    true,
				Computed:    true,
			},
			"offline_process_database_open_timeout": schema.StringAttribute{
				Description: "Specifies a timeout duration which will be used for opening the database environment by an offline process, such as export-ldif.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_private_backend": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to `ldif`: Indicates whether the backend should be considered a private backend, which indicates that it is used for storing operational data rather than user-defined information. When the `type` attribute is set to `local-db`: Indicates whether this backend should be considered a private backend in the server. Private backends are meant for storing server-internal information and should not be used for user or application data.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ldif`: Indicates whether the backend should be considered a private backend, which indicates that it is used for storing operational data rather than user-defined information.\n  - `local-db`: Indicates whether this backend should be considered a private backend in the server. Private backends are meant for storing server-internal information and should not be used for user or application data.",
				Optional:            true,
				Computed:            true,
			},
			"base_dn": schema.SetAttribute{
				Description: "Specifies the base DN(s) for the data that the backend handles.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"writability_mode": schema.StringAttribute{
				Description: "Specifies the behavior that the backend should use when processing write operations.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Backend",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the backend is enabled in the server.",
				Required:    true,
			},
			"set_degraded_alert_when_disabled": schema.BoolAttribute{
				Description: "Determines whether the Directory Server enters a DEGRADED state (and sends a corresponding alert) when this Backend is disabled.",
				Optional:    true,
				Computed:    true,
			},
			"return_unavailable_when_disabled": schema.BoolAttribute{
				Description: "Determines whether any LDAP operation that would use this Backend is to return UNAVAILABLE when this Backend is disabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"notification_manager": schema.StringAttribute{
				Description: "Specifies a notification manager for changes resulting from operations processed through this Backend",
				Optional:    true,
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
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"schema", "backup", "encryption-settings", "ldif", "trust-store", "custom", "changelog", "monitor", "local-db", "config-file-handler", "task", "alert", "alarm", "metrics"}...),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		schemaDef.Attributes["storage_dir"] = schema.StringAttribute{
			Description: "Specifies the path to the directory that will be used to store queued samples.",
		}
		schemaDef.Attributes["metrics_dir"] = schema.StringAttribute{
			Description: "Specifies the path to the directory that contains metric definitions.",
		}
		schemaDef.Attributes["sample_flush_interval"] = schema.StringAttribute{
			Description: "Period when samples are flushed to disk.",
		}
		schemaDef.Attributes["retention_policy"] = schema.SetAttribute{
			Description: "The retention policy to use for the Metrics Backend .",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["alarm_retention_time"] = schema.StringAttribute{
			Description: "Specifies the maximum length of time that information about raised alarms should be maintained before they will be purged.",
		}
		schemaDef.Attributes["max_alarms"] = schema.Int64Attribute{
			Description: "Specifies the maximum number of alarms that should be retained. If more alarms than this configured maximum are generated within the alarm retention time, then the oldest alarms will be purged to achieve this maximum. Only alarms at normal severity will be purged.",
		}
		schemaDef.Attributes["alert_retention_time"] = schema.StringAttribute{
			Description: "Specifies the maximum length of time that information about generated alerts should be maintained before they will be purged.",
		}
		schemaDef.Attributes["max_alerts"] = schema.Int64Attribute{
			Description: "Specifies the maximum number of alerts that should be retained. If more alerts than this configured maximum are generated within the alert retention time, then the oldest alerts will be purged to achieve this maximum.",
		}
		schemaDef.Attributes["disabled_alert_type"] = schema.SetAttribute{
			Description: "Specifies the names of the alert types that should not be added to the backend. This can be used to suppress high volume alerts that might trigger hitting the max-alerts limit sooner than desired. Disabled alert types will not be sent out over persistent searches on this backend.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["task_backing_file"] = schema.StringAttribute{
			Description: "Specifies the path to the backing file for storing information about the tasks configured in the server.",
		}
		schemaDef.Attributes["maximum_initial_task_log_messages_to_retain"] = schema.Int64Attribute{
			Description: "The maximum number of log messages to retain in each task entry from the beginning of the processing for that task. If too many messages are logged during task processing, then retaining only a limited number of messages from the beginning and/or end of task processing can reduce the amount of memory that the server consumes by caching information about currently-active and recently-completed tasks.",
		}
		schemaDef.Attributes["maximum_final_task_log_messages_to_retain"] = schema.Int64Attribute{
			Description: "The maximum number of log messages to retain in each task entry from the end of the processing for that task. If too many messages are logged during task processing, then retaining only a limited number of messages from the beginning and/or end of task processing can reduce the amount of memory that the server consumes by caching information about currently-active and recently-completed tasks.",
		}
		schemaDef.Attributes["task_retention_time"] = schema.StringAttribute{
			Description: "Specifies the length of time that task entries should be retained after processing on the associated task has been completed.",
		}
		schemaDef.Attributes["notification_sender_address"] = schema.StringAttribute{
			Description: "Specifies the email address to use as the sender address (that is, the \"From:\" address) for notification mail messages generated when a task completes execution.",
		}
		schemaDef.Attributes["insignificant_config_archive_attribute"] = schema.SetAttribute{
			Description: "The name or OID of an attribute type that is considered insignificant for the purpose of maintaining the configuration archive.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["insignificant_config_archive_base_dn"] = schema.SetAttribute{
			Description: "Supported in PingDirectory product version 9.3.0.0+. The base DN that is considered insignificant for the purpose of maintaining the configuration archive.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["maintain_config_archive"] = schema.BoolAttribute{
			Description: "Supported in PingDirectory product version 9.3.0.0+. Indicates whether the server should maintain the config archive with new changes to the config backend.",
		}
		schemaDef.Attributes["max_config_archive_count"] = schema.Int64Attribute{
			Description: "Supported in PingDirectory product version 9.3.0.0+. Indicates the maximum number of previous config files to keep as part of maintaining the config archive.",
		}
		schemaDef.Attributes["mirrored_subtree_peer_polling_interval"] = schema.StringAttribute{
			Description: "Tells the server component that is responsible for mirroring configuration data across a topology of servers the maximum amount of time to wait before polling the peer servers in the topology to determine if there are any changes in the topology. Mirrored data includes meta-data about the servers in the topology as well as cluster-wide configuration data.",
		}
		schemaDef.Attributes["mirrored_subtree_entry_update_timeout"] = schema.StringAttribute{
			Description: "Tells the server component that is responsible for mirroring configuration data across a topology of servers the maximum amount of time to wait for an update operation (add, delete, modify and modify-dn) on an entry to be applied on all servers in the topology. Mirrored data includes meta-data about the servers in the topology as well as cluster-wide configuration data.",
		}
		schemaDef.Attributes["mirrored_subtree_search_timeout"] = schema.StringAttribute{
			Description: "Tells the server component that is responsible for mirroring configuration data across a topology of servers the maximum amount of time to wait for a search operation to complete. Mirrored data includes meta-data about the servers in the topology as well as cluster-wide configuration data. Search requests that take longer than this timeout will be canceled and considered failures.",
		}
		schemaDef.Attributes["changelog_write_batch_size"] = schema.Int64Attribute{
			Description: "Specifies the number of changelog entries written in a single database transaction.",
		}
		schemaDef.Attributes["changelog_purge_batch_size"] = schema.Int64Attribute{
			Description: "Specifies the number of changelog entries purged in a single database transaction.",
		}
		schemaDef.Attributes["changelog_write_queue_capacity"] = schema.Int64Attribute{
			Description: "Specifies the capacity of the changelog write queue in number of changes.",
		}
		schemaDef.Attributes["index_include_attribute"] = schema.SetAttribute{
			Description: "Specifies which attribute types are to be specifically included in the set of attribute indexes maintained on the changelog. If this property does not have any values then no attribute types are indexed.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["index_exclude_attribute"] = schema.SetAttribute{
			Description: "Specifies which attribute types are to be specifically excluded from the set of attribute indexes maintained on the changelog. This property is useful when the index-include-attribute property contains one of the special values \"*\" and \"+\".",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["changelog_maximum_age"] = schema.StringAttribute{
			Description: "Changes are guaranteed to be maintained in the changelog database for at least this duration. Setting target-database-size can allow additional changes to be maintained up to the configured size on disk.",
		}
		schemaDef.Attributes["target_database_size"] = schema.StringAttribute{
			Description: "The changelog database is allowed to grow up to this size on disk even if changes are older than the configured changelog-maximum-age.",
		}
		schemaDef.Attributes["changelog_entry_include_base_dn"] = schema.SetAttribute{
			Description: "The base DNs for branches in the data for which to record changes in the changelog.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["changelog_entry_exclude_base_dn"] = schema.SetAttribute{
			Description: "The base DNs for branches in the data for which no changelog records should be generated.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["changelog_entry_include_filter"] = schema.SetAttribute{
			Description: "A filter that indicates which changelog entries should actually be stored in the changelog. Note that this filter is evaluated against the changelog entry itself and not against the entry that was the target of the change referenced by the changelog entry. This filter may target any attributes that appear in changelog entries with the exception of the changeNumber and entry-size-bytes attributes, since they will not be known at the time of the filter evaluation.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["changelog_entry_exclude_filter"] = schema.SetAttribute{
			Description: "A filter that indicates which changelog entries should be excluded from the changelog. Note that this filter is evaluated against the changelog entry itself and not against the entry that was the target of the change referenced by the changelog entry. This filter may target any attributes that appear in changelog entries with the exception of the changeNumber and entry-size-bytes attributes, since they will not be known at the time of the filter evaluation.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["changelog_include_attribute"] = schema.SetAttribute{
			Description: "Specifies which attribute types will be included in a changelog entry for ADD and MODIFY operations.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["changelog_exclude_attribute"] = schema.SetAttribute{
			Description: "Specifies a set of attribute types that should be excluded in a changelog entry for ADD and MODIFY operations.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["changelog_deleted_entry_include_attribute"] = schema.SetAttribute{
			Description: "Specifies a set of attribute types that should be included in a changelog entry for DELETE operations.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["changelog_deleted_entry_exclude_attribute"] = schema.SetAttribute{
			Description: "Specifies a set of attribute types that should be excluded from a changelog entry for DELETE operations.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["changelog_include_key_attribute"] = schema.SetAttribute{
			Description: "Specifies which attribute types will be included in a changelog entry on every change.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["changelog_max_before_after_values"] = schema.Int64Attribute{
			Description: "This controls whether all attribute values for a modified attribute (even those values that have not changed) will be included in the changelog entry. If the number of attribute values does not exceed this limit, then all values for the modified attribute will be included in the changelog entry.",
		}
		schemaDef.Attributes["write_lastmod_attributes"] = schema.BoolAttribute{
			Description: "Specifies whether values of creatorsName, createTimestamp, modifiersName and modifyTimestamp attributes will be written to changelog entries.",
		}
		schemaDef.Attributes["use_reversible_form"] = schema.BoolAttribute{
			Description: "Specifies whether the changelog should provide enough information to be able to revert the changes if desired.",
		}
		schemaDef.Attributes["include_virtual_attributes"] = schema.SetAttribute{
			Description: "Specifies the changelog entry elements (if any) in which virtual attributes should be included.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["apply_access_controls_to_changelog_entry_contents"] = schema.BoolAttribute{
			Description: "Indicates whether the contents of changelog entries should be subject to access control and sensitive attribute evaluation such that the contents of attributes like changes, deletedEntryAttrs, ds-changelog-entry-key-attr-values, ds-changelog-before-values, and ds-changelog-after-values may be altered based on attributes the user can see in the target entry.",
		}
		schemaDef.Attributes["report_excluded_changelog_attributes"] = schema.StringAttribute{
			Description: "Indicates whether changelog entries that have been altered by applying access controls should include additional information about any attributes that may have been removed.",
		}
		schemaDef.Attributes["soft_delete_entry_included_operation"] = schema.SetAttribute{
			Description: "Specifies which operations performed on soft-deleted entries will appear in the changelog.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["ldif_file"] = schema.StringAttribute{
			Description:         "When the `type` attribute is set to  one of [`alert`, `alarm`]: Specifies the path to the LDIF file that serves as the backing file for this backend. When the `type` attribute is set to `ldif`: Specifies the path to the LDIF file containing the data for this backend.",
			MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`alert`, `alarm`]: Specifies the path to the LDIF file that serves as the backing file for this backend.\n  - `ldif`: Specifies the path to the LDIF file containing the data for this backend.",
		}
		schemaDef.Attributes["trust_store_file"] = schema.StringAttribute{
			Description: "Specifies the path to the file that stores the trust information.",
		}
		schemaDef.Attributes["trust_store_type"] = schema.StringAttribute{
			Description: "Specifies the format for the data in the key store file.",
		}
		schemaDef.Attributes["trust_store_pin"] = schema.StringAttribute{
			Description: "Specifies the clear-text PIN needed to access the Trust Store Backend.",
			Sensitive:   true,
		}
		schemaDef.Attributes["trust_store_pin_file"] = schema.StringAttribute{
			Description: "Specifies the path to the text file whose only contents should be a single line containing the clear-text PIN needed to access the Trust Store Backend.",
		}
		schemaDef.Attributes["trust_store_pin_passphrase_provider"] = schema.StringAttribute{
			Description: "The passphrase provider to use to obtain the clear-text PIN needed to access the Trust Store Backend.",
		}
		schemaDef.Attributes["backup_directory"] = schema.SetAttribute{
			Description: "Specifies the path to a backup directory containing one or more backups for a particular backend.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["schema_entry_dn"] = schema.SetAttribute{
			Description: "Defines the base DNs of the subtrees in which the schema information is published in addition to the value included in the base-dn property.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["show_all_attributes"] = schema.BoolAttribute{
			Description: "Indicates whether to treat all attributes in the schema entry as if they were user attributes regardless of their configuration.",
		}
		schemaDef.Attributes["read_only_schema_file"] = schema.SetAttribute{
			Description: "Specifies the name of a file (which must exist in the config/schema directory) containing schema elements that should be considered read-only. Any schema definitions contained in read-only files cannot be altered by external clients.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["backup_file_permissions"] = schema.StringAttribute{
			Description: "Specifies the permissions that should be applied to files and directories created by a backup of the backend.",
		}
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type", "backend_id"})
	} else {
		// Add RequiresReplace modifier for read-only attributes
		backendIdAttr := schemaDef.Attributes["backend_id"].(schema.StringAttribute)
		backendIdAttr.PlanModifiers = append(backendIdAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["backend_id"] = backendIdAttr
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *backendResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanBackend(ctx, req, resp, r.apiClient, r.providerConfig)
	var planModel, configModel backendResourceModel
	req.Config.Get(ctx, &configModel)
	req.Plan.Get(ctx, &planModel)
	resourceType := planModel.Type.ValueString()
	anyDefaultsSet := false
	// Set defaults for local-db type
	if resourceType == "local-db" {
		if !internaltypes.IsDefined(configModel.UncachedId2entryCacheMode) {
			defaultVal := types.StringValue("cache-keys-only")
			if !planModel.UncachedId2entryCacheMode.Equal(defaultVal) {
				planModel.UncachedId2entryCacheMode = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.WritabilityMode) {
			defaultVal := types.StringValue("enabled")
			if !planModel.WritabilityMode.Equal(defaultVal) {
				planModel.WritabilityMode = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SetDegradedAlertForUntrustedIndex) {
			defaultVal := types.BoolValue(true)
			if !planModel.SetDegradedAlertForUntrustedIndex.Equal(defaultVal) {
				planModel.SetDegradedAlertForUntrustedIndex = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ReturnUnavailableForUntrustedIndex) {
			defaultVal := types.BoolValue(true)
			if !planModel.ReturnUnavailableForUntrustedIndex.Equal(defaultVal) {
				planModel.ReturnUnavailableForUntrustedIndex = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ProcessFiltersWithUndefinedAttributeTypes) {
			defaultVal := types.BoolValue(false)
			if !planModel.ProcessFiltersWithUndefinedAttributeTypes.Equal(defaultVal) {
				planModel.ProcessFiltersWithUndefinedAttributeTypes = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IsPrivateBackend) {
			defaultVal := types.BoolValue(false)
			if !planModel.IsPrivateBackend.Equal(defaultVal) {
				planModel.IsPrivateBackend = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DbDirectory) {
			defaultVal := types.StringValue("db")
			if !planModel.DbDirectory.Equal(defaultVal) {
				planModel.DbDirectory = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DbDirectoryPermissions) {
			defaultVal := types.StringValue("700")
			if !planModel.DbDirectoryPermissions.Equal(defaultVal) {
				planModel.DbDirectoryPermissions = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CompressEntries) {
			defaultVal := types.BoolValue(false)
			if !planModel.CompressEntries.Equal(defaultVal) {
				planModel.CompressEntries = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.HashEntries) {
			defaultVal := types.BoolValue(false)
			if !planModel.HashEntries.Equal(defaultVal) {
				planModel.HashEntries = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DbNumCleanerThreads) {
			defaultVal := types.Int64Value(0)
			if !planModel.DbNumCleanerThreads.Equal(defaultVal) {
				planModel.DbNumCleanerThreads = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DbCleanerMinUtilization) {
			defaultVal := types.Int64Value(75)
			if !planModel.DbCleanerMinUtilization.Equal(defaultVal) {
				planModel.DbCleanerMinUtilization = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DbEvictorCriticalPercentage) {
			defaultVal := types.Int64Value(0)
			if !planModel.DbEvictorCriticalPercentage.Equal(defaultVal) {
				planModel.DbEvictorCriticalPercentage = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DbUseThreadLocalHandles) {
			defaultVal := types.BoolValue(true)
			if !planModel.DbUseThreadLocalHandles.Equal(defaultVal) {
				planModel.DbUseThreadLocalHandles = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DbLogFileMax) {
			defaultVal := types.StringValue("50 mb")
			if !planModel.DbLogFileMax.Equal(defaultVal) {
				planModel.DbLogFileMax = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DbLoggingLevel) {
			defaultVal := types.StringValue("CONFIG")
			if !planModel.DbLoggingLevel.Equal(defaultVal) {
				planModel.DbLoggingLevel = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DbCachePercent) {
			defaultVal := types.Int64Value(10)
			if !planModel.DbCachePercent.Equal(defaultVal) {
				planModel.DbCachePercent = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefaultCacheMode) {
			defaultVal := types.StringValue("cache-keys-and-values")
			if !planModel.DefaultCacheMode.Equal(defaultVal) {
				planModel.DefaultCacheMode = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.PrimeMethod) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("none")})
			if !planModel.PrimeMethod.Equal(defaultVal) {
				planModel.PrimeMethod = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.PrimeThreadCount) {
			defaultVal := types.Int64Value(2)
			if !planModel.PrimeThreadCount.Equal(defaultVal) {
				planModel.PrimeThreadCount = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.PrimeAllIndexes) {
			defaultVal := types.BoolValue(true)
			if !planModel.PrimeAllIndexes.Equal(defaultVal) {
				planModel.PrimeAllIndexes = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.BackgroundPrime) {
			defaultVal := types.BoolValue(false)
			if !planModel.BackgroundPrime.Equal(defaultVal) {
				planModel.BackgroundPrime = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IndexEntryLimit) {
			defaultVal := types.Int64Value(4000)
			if !planModel.IndexEntryLimit.Equal(defaultVal) {
				planModel.IndexEntryLimit = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CompositeIndexEntryLimit) {
			defaultVal := types.Int64Value(100000)
			if !planModel.CompositeIndexEntryLimit.Equal(defaultVal) {
				planModel.CompositeIndexEntryLimit = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ImportTempDirectory) {
			defaultVal := types.StringValue("import-tmp")
			if !planModel.ImportTempDirectory.Equal(defaultVal) {
				planModel.ImportTempDirectory = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ImportThreadCount) {
			defaultVal := types.Int64Value(16)
			if !planModel.ImportThreadCount.Equal(defaultVal) {
				planModel.ImportThreadCount = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ExportThreadCount) {
			defaultVal := types.Int64Value(0)
			if !planModel.ExportThreadCount.Equal(defaultVal) {
				planModel.ExportThreadCount = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DbImportCachePercent) {
			defaultVal := types.Int64Value(60)
			if !planModel.DbImportCachePercent.Equal(defaultVal) {
				planModel.DbImportCachePercent = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DbTxnWriteNoSync) {
			defaultVal := types.BoolValue(true)
			if !planModel.DbTxnWriteNoSync.Equal(defaultVal) {
				planModel.DbTxnWriteNoSync = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DeadlockRetryLimit) {
			defaultVal := types.Int64Value(3)
			if !planModel.DeadlockRetryLimit.Equal(defaultVal) {
				planModel.DeadlockRetryLimit = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ExternalTxnDefaultBackendLockBehavior) {
			defaultVal := types.StringValue("acquire-after-retries")
			if !planModel.ExternalTxnDefaultBackendLockBehavior.Equal(defaultVal) {
				planModel.ExternalTxnDefaultBackendLockBehavior = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SingleWriterLockBehavior) {
			defaultVal := types.StringValue("acquire-on-retry")
			if !planModel.SingleWriterLockBehavior.Equal(defaultVal) {
				planModel.SingleWriterLockBehavior = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SubtreeDeleteSizeLimit) {
			defaultVal := types.Int64Value(5000)
			if !planModel.SubtreeDeleteSizeLimit.Equal(defaultVal) {
				planModel.SubtreeDeleteSizeLimit = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.NumRecentChanges) {
			defaultVal := types.Int64Value(50000)
			if !planModel.NumRecentChanges.Equal(defaultVal) {
				planModel.NumRecentChanges = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SetDegradedAlertWhenDisabled) {
			defaultVal := types.BoolValue(true)
			if !planModel.SetDegradedAlertWhenDisabled.Equal(defaultVal) {
				planModel.SetDegradedAlertWhenDisabled = defaultVal
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

func (r *defaultBackendResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanBackend(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanBackend(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	compare, err := version.Compare(providerConfig.ProductVersion, version.PingDirectory9300)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model defaultBackendResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsNonEmptySet(model.InsignificantConfigArchiveBaseDN) {
		resp.Diagnostics.AddError("Attribute 'insignificant_config_archive_base_dn' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
	if internaltypes.IsDefined(model.MaintainConfigArchive) {
		resp.Diagnostics.AddError("Attribute 'maintain_config_archive' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
	if internaltypes.IsDefined(model.MaxConfigArchiveCount) {
		resp.Diagnostics.AddError("Attribute 'max_config_archive_count' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsBackend() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("base_dn"),
			path.MatchRoot("type"),
			[]string{"schema", "backup", "encryption-settings", "ldif", "trust-store", "custom", "changelog", "monitor", "local-db", "config-file-handler", "task", "alert", "alarm"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("writability_mode"),
			path.MatchRoot("type"),
			[]string{"schema", "backup", "ldif", "trust-store", "custom", "local-db", "config-file-handler", "task", "alert", "alarm", "metrics"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("set_degraded_alert_when_disabled"),
			path.MatchRoot("type"),
			[]string{"schema", "backup", "encryption-settings", "ldif", "trust-store", "custom", "changelog", "monitor", "local-db", "config-file-handler", "task", "alert", "alarm"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("is_private_backend"),
			path.MatchRoot("type"),
			[]string{"ldif", "local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("db_directory"),
			path.MatchRoot("type"),
			[]string{"changelog", "local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("db_directory_permissions"),
			path.MatchRoot("type"),
			[]string{"changelog", "local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("db_cache_percent"),
			path.MatchRoot("type"),
			[]string{"changelog", "local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("je_property"),
			path.MatchRoot("type"),
			[]string{"changelog", "local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("uncached_id2entry_cache_mode"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("uncached_attribute_criteria"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("uncached_entry_criteria"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("set_degraded_alert_for_untrusted_index"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("return_unavailable_for_untrusted_index"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("process_filters_with_undefined_attribute_types"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("compact_common_parent_dn"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("compress_entries"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("hash_entries"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("db_num_cleaner_threads"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("db_cleaner_min_utilization"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("db_evictor_critical_percentage"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("db_checkpointer_wakeup_interval"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("db_background_sync_interval"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("db_use_thread_local_handles"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("db_log_file_max"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("db_logging_level"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("default_cache_mode"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("id2entry_cache_mode"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("dn2id_cache_mode"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("id2children_cache_mode"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("id2subtree_cache_mode"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("dn2uri_cache_mode"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("prime_method"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("prime_thread_count"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("prime_time_limit"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("prime_all_indexes"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("system_index_to_prime"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("system_index_to_prime_internal_nodes_only"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("background_prime"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("index_entry_limit"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("composite_index_entry_limit"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("id2children_index_entry_limit"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("id2subtree_index_entry_limit"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("import_temp_directory"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("import_thread_count"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("export_thread_count"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("db_import_cache_percent"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("db_txn_write_no_sync"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("deadlock_retry_limit"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("external_txn_default_backend_lock_behavior"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("single_writer_lock_behavior"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("subtree_delete_size_limit"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("num_recent_changes"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("offline_process_database_open_timeout"),
			path.MatchRoot("type"),
			[]string{"local-db"},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"local-db",
			[]path.Expression{path.MatchRoot("backend_id"), path.MatchRoot("base_dn"), path.MatchRoot("enabled")},
		),
	}
}

// Add config validators
func (r backendResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsBackend()
}

// Add config validators
func (r defaultBackendResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	validators := []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("schema_entry_dn"),
			path.MatchRoot("type"),
			[]string{"schema"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("show_all_attributes"),
			path.MatchRoot("type"),
			[]string{"schema"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("read_only_schema_file"),
			path.MatchRoot("type"),
			[]string{"schema"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("backup_file_permissions"),
			path.MatchRoot("type"),
			[]string{"schema", "encryption-settings", "ldif", "trust-store", "custom", "config-file-handler", "task", "alert", "alarm"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("backup_directory"),
			path.MatchRoot("type"),
			[]string{"backup"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("ldif_file"),
			path.MatchRoot("type"),
			[]string{"ldif", "alert", "alarm"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_store_file"),
			path.MatchRoot("type"),
			[]string{"trust-store"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_store_type"),
			path.MatchRoot("type"),
			[]string{"trust-store"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_store_pin"),
			path.MatchRoot("type"),
			[]string{"trust-store"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_store_pin_file"),
			path.MatchRoot("type"),
			[]string{"trust-store"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_store_pin_passphrase_provider"),
			path.MatchRoot("type"),
			[]string{"trust-store"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("changelog_write_batch_size"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("changelog_purge_batch_size"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("changelog_write_queue_capacity"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("index_include_attribute"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("index_exclude_attribute"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("changelog_maximum_age"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("target_database_size"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("changelog_entry_include_base_dn"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("changelog_entry_exclude_base_dn"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("changelog_entry_include_filter"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("changelog_entry_exclude_filter"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("changelog_include_attribute"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("changelog_exclude_attribute"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("changelog_deleted_entry_include_attribute"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("changelog_deleted_entry_exclude_attribute"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("changelog_include_key_attribute"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("changelog_max_before_after_values"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("write_lastmod_attributes"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("use_reversible_form"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_virtual_attributes"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("apply_access_controls_to_changelog_entry_contents"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("report_excluded_changelog_attributes"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("soft_delete_entry_included_operation"),
			path.MatchRoot("type"),
			[]string{"changelog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("insignificant_config_archive_attribute"),
			path.MatchRoot("type"),
			[]string{"config-file-handler"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("insignificant_config_archive_base_dn"),
			path.MatchRoot("type"),
			[]string{"config-file-handler"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("maintain_config_archive"),
			path.MatchRoot("type"),
			[]string{"config-file-handler"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("max_config_archive_count"),
			path.MatchRoot("type"),
			[]string{"config-file-handler"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("mirrored_subtree_peer_polling_interval"),
			path.MatchRoot("type"),
			[]string{"config-file-handler"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("mirrored_subtree_entry_update_timeout"),
			path.MatchRoot("type"),
			[]string{"config-file-handler"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("mirrored_subtree_search_timeout"),
			path.MatchRoot("type"),
			[]string{"config-file-handler"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("task_backing_file"),
			path.MatchRoot("type"),
			[]string{"task"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("maximum_initial_task_log_messages_to_retain"),
			path.MatchRoot("type"),
			[]string{"task"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("maximum_final_task_log_messages_to_retain"),
			path.MatchRoot("type"),
			[]string{"task"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("task_retention_time"),
			path.MatchRoot("type"),
			[]string{"task"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("notification_sender_address"),
			path.MatchRoot("type"),
			[]string{"task"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("alert_retention_time"),
			path.MatchRoot("type"),
			[]string{"alert"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("max_alerts"),
			path.MatchRoot("type"),
			[]string{"alert"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("disabled_alert_type"),
			path.MatchRoot("type"),
			[]string{"alert"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("alarm_retention_time"),
			path.MatchRoot("type"),
			[]string{"alarm"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("max_alarms"),
			path.MatchRoot("type"),
			[]string{"alarm"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("storage_dir"),
			path.MatchRoot("type"),
			[]string{"metrics"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("metrics_dir"),
			path.MatchRoot("type"),
			[]string{"metrics"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("sample_flush_interval"),
			path.MatchRoot("type"),
			[]string{"metrics"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("retention_policy"),
			path.MatchRoot("type"),
			[]string{"metrics"},
		),
	}
	return append(configValidatorsBackend(), validators...)
}

// Add optional fields to create request for local-db backend
func addOptionalLocalDbBackendFields(ctx context.Context, addRequest *client.AddLocalDbBackendRequest, plan backendResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UncachedId2entryCacheMode) {
		uncachedId2entryCacheMode, err := client.NewEnumbackendUncachedId2entryCacheModePropFromValue(plan.UncachedId2entryCacheMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.UncachedId2entryCacheMode = uncachedId2entryCacheMode
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UncachedAttributeCriteria) {
		addRequest.UncachedAttributeCriteria = plan.UncachedAttributeCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UncachedEntryCriteria) {
		addRequest.UncachedEntryCriteria = plan.UncachedEntryCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.WritabilityMode) {
		writabilityMode, err := client.NewEnumbackendWritabilityModePropFromValue(plan.WritabilityMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.WritabilityMode = writabilityMode
	}
	if internaltypes.IsDefined(plan.SetDegradedAlertForUntrustedIndex) {
		addRequest.SetDegradedAlertForUntrustedIndex = plan.SetDegradedAlertForUntrustedIndex.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ReturnUnavailableForUntrustedIndex) {
		addRequest.ReturnUnavailableForUntrustedIndex = plan.ReturnUnavailableForUntrustedIndex.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ProcessFiltersWithUndefinedAttributeTypes) {
		addRequest.ProcessFiltersWithUndefinedAttributeTypes = plan.ProcessFiltersWithUndefinedAttributeTypes.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IsPrivateBackend) {
		addRequest.IsPrivateBackend = plan.IsPrivateBackend.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DbDirectory) {
		addRequest.DbDirectory = plan.DbDirectory.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DbDirectoryPermissions) {
		addRequest.DbDirectoryPermissions = plan.DbDirectoryPermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CompactCommonParentDN) {
		var slice []string
		plan.CompactCommonParentDN.ElementsAs(ctx, &slice, false)
		addRequest.CompactCommonParentDN = slice
	}
	if internaltypes.IsDefined(plan.CompressEntries) {
		addRequest.CompressEntries = plan.CompressEntries.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.HashEntries) {
		addRequest.HashEntries = plan.HashEntries.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.DbNumCleanerThreads) {
		addRequest.DbNumCleanerThreads = plan.DbNumCleanerThreads.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.DbCleanerMinUtilization) {
		addRequest.DbCleanerMinUtilization = plan.DbCleanerMinUtilization.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.DbEvictorCriticalPercentage) {
		addRequest.DbEvictorCriticalPercentage = plan.DbEvictorCriticalPercentage.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DbCheckpointerWakeupInterval) {
		addRequest.DbCheckpointerWakeupInterval = plan.DbCheckpointerWakeupInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DbBackgroundSyncInterval) {
		addRequest.DbBackgroundSyncInterval = plan.DbBackgroundSyncInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.DbUseThreadLocalHandles) {
		addRequest.DbUseThreadLocalHandles = plan.DbUseThreadLocalHandles.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DbLogFileMax) {
		addRequest.DbLogFileMax = plan.DbLogFileMax.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DbLoggingLevel) {
		addRequest.DbLoggingLevel = plan.DbLoggingLevel.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.JeProperty) {
		var slice []string
		plan.JeProperty.ElementsAs(ctx, &slice, false)
		addRequest.JeProperty = slice
	}
	if internaltypes.IsDefined(plan.DbCachePercent) {
		addRequest.DbCachePercent = plan.DbCachePercent.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DefaultCacheMode) {
		defaultCacheMode, err := client.NewEnumbackendDefaultCacheModePropFromValue(plan.DefaultCacheMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.DefaultCacheMode = defaultCacheMode
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Id2entryCacheMode) {
		id2entryCacheMode, err := client.NewEnumbackendId2entryCacheModePropFromValue(plan.Id2entryCacheMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.Id2entryCacheMode = id2entryCacheMode
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Dn2idCacheMode) {
		dn2idCacheMode, err := client.NewEnumbackendDn2idCacheModePropFromValue(plan.Dn2idCacheMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.Dn2idCacheMode = dn2idCacheMode
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Id2childrenCacheMode) {
		id2childrenCacheMode, err := client.NewEnumbackendId2childrenCacheModePropFromValue(plan.Id2childrenCacheMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.Id2childrenCacheMode = id2childrenCacheMode
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Id2subtreeCacheMode) {
		id2subtreeCacheMode, err := client.NewEnumbackendId2subtreeCacheModePropFromValue(plan.Id2subtreeCacheMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.Id2subtreeCacheMode = id2subtreeCacheMode
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Dn2uriCacheMode) {
		dn2uriCacheMode, err := client.NewEnumbackendDn2uriCacheModePropFromValue(plan.Dn2uriCacheMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.Dn2uriCacheMode = dn2uriCacheMode
	}
	if internaltypes.IsDefined(plan.PrimeMethod) {
		var slice []string
		plan.PrimeMethod.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumbackendPrimeMethodProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumbackendPrimeMethodPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.PrimeMethod = enumSlice
	}
	if internaltypes.IsDefined(plan.PrimeThreadCount) {
		addRequest.PrimeThreadCount = plan.PrimeThreadCount.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PrimeTimeLimit) {
		addRequest.PrimeTimeLimit = plan.PrimeTimeLimit.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.PrimeAllIndexes) {
		addRequest.PrimeAllIndexes = plan.PrimeAllIndexes.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SystemIndexToPrime) {
		var slice []string
		plan.SystemIndexToPrime.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumbackendSystemIndexToPrimeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumbackendSystemIndexToPrimePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.SystemIndexToPrime = enumSlice
	}
	if internaltypes.IsDefined(plan.SystemIndexToPrimeInternalNodesOnly) {
		var slice []string
		plan.SystemIndexToPrimeInternalNodesOnly.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumbackendSystemIndexToPrimeInternalNodesOnlyProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumbackendSystemIndexToPrimeInternalNodesOnlyPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.SystemIndexToPrimeInternalNodesOnly = enumSlice
	}
	if internaltypes.IsDefined(plan.BackgroundPrime) {
		addRequest.BackgroundPrime = plan.BackgroundPrime.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IndexEntryLimit) {
		addRequest.IndexEntryLimit = plan.IndexEntryLimit.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.CompositeIndexEntryLimit) {
		addRequest.CompositeIndexEntryLimit = plan.CompositeIndexEntryLimit.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.Id2childrenIndexEntryLimit) {
		addRequest.Id2childrenIndexEntryLimit = plan.Id2childrenIndexEntryLimit.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.Id2subtreeIndexEntryLimit) {
		addRequest.Id2subtreeIndexEntryLimit = plan.Id2subtreeIndexEntryLimit.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ImportTempDirectory) {
		addRequest.ImportTempDirectory = plan.ImportTempDirectory.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.ImportThreadCount) {
		addRequest.ImportThreadCount = plan.ImportThreadCount.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.ExportThreadCount) {
		addRequest.ExportThreadCount = plan.ExportThreadCount.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.DbImportCachePercent) {
		addRequest.DbImportCachePercent = plan.DbImportCachePercent.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.DbTxnWriteNoSync) {
		addRequest.DbTxnWriteNoSync = plan.DbTxnWriteNoSync.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.DeadlockRetryLimit) {
		addRequest.DeadlockRetryLimit = plan.DeadlockRetryLimit.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ExternalTxnDefaultBackendLockBehavior) {
		externalTxnDefaultBackendLockBehavior, err := client.NewEnumbackendExternalTxnDefaultBackendLockBehaviorPropFromValue(plan.ExternalTxnDefaultBackendLockBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.ExternalTxnDefaultBackendLockBehavior = externalTxnDefaultBackendLockBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SingleWriterLockBehavior) {
		singleWriterLockBehavior, err := client.NewEnumbackendSingleWriterLockBehaviorPropFromValue(plan.SingleWriterLockBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.SingleWriterLockBehavior = singleWriterLockBehavior
	}
	if internaltypes.IsDefined(plan.SubtreeDeleteSizeLimit) {
		addRequest.SubtreeDeleteSizeLimit = plan.SubtreeDeleteSizeLimit.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.NumRecentChanges) {
		addRequest.NumRecentChanges = plan.NumRecentChanges.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OfflineProcessDatabaseOpenTimeout) {
		addRequest.OfflineProcessDatabaseOpenTimeout = plan.OfflineProcessDatabaseOpenTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.SetDegradedAlertWhenDisabled) {
		addRequest.SetDegradedAlertWhenDisabled = plan.SetDegradedAlertWhenDisabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ReturnUnavailableWhenDisabled) {
		addRequest.ReturnUnavailableWhenDisabled = plan.ReturnUnavailableWhenDisabled.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.NotificationManager) {
		addRequest.NotificationManager = plan.NotificationManager.ValueStringPointer()
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateBackendUnknownValues(model *backendResourceModel) {
	if model.JeProperty.IsUnknown() || model.JeProperty.IsNull() {
		model.JeProperty, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.BaseDN.IsUnknown() || model.BaseDN.IsNull() {
		model.BaseDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.CompactCommonParentDN.IsUnknown() || model.CompactCommonParentDN.IsNull() {
		model.CompactCommonParentDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.SystemIndexToPrimeInternalNodesOnly.IsUnknown() || model.SystemIndexToPrimeInternalNodesOnly.IsNull() {
		model.SystemIndexToPrimeInternalNodesOnly, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.PrimeMethod.IsUnknown() || model.PrimeMethod.IsNull() {
		model.PrimeMethod, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.SystemIndexToPrime.IsUnknown() || model.SystemIndexToPrime.IsNull() {
		model.SystemIndexToPrime, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateBackendUnknownValuesDefault(model *defaultBackendResourceModel) {
	if model.InsignificantConfigArchiveBaseDN.IsUnknown() || model.InsignificantConfigArchiveBaseDN.IsNull() {
		model.InsignificantConfigArchiveBaseDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ChangelogEntryIncludeFilter.IsUnknown() || model.ChangelogEntryIncludeFilter.IsNull() {
		model.ChangelogEntryIncludeFilter, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludeVirtualAttributes.IsUnknown() || model.IncludeVirtualAttributes.IsNull() {
		model.IncludeVirtualAttributes, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.SystemIndexToPrime.IsUnknown() || model.SystemIndexToPrime.IsNull() {
		model.SystemIndexToPrime, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ReadOnlySchemaFile.IsUnknown() || model.ReadOnlySchemaFile.IsNull() {
		model.ReadOnlySchemaFile, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ChangelogEntryExcludeFilter.IsUnknown() || model.ChangelogEntryExcludeFilter.IsNull() {
		model.ChangelogEntryExcludeFilter, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.BackupDirectory.IsUnknown() || model.BackupDirectory.IsNull() {
		model.BackupDirectory, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.BaseDN.IsUnknown() || model.BaseDN.IsNull() {
		model.BaseDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.SchemaEntryDN.IsUnknown() || model.SchemaEntryDN.IsNull() {
		model.SchemaEntryDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IndexExcludeAttribute.IsUnknown() || model.IndexExcludeAttribute.IsNull() {
		model.IndexExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ChangelogExcludeAttribute.IsUnknown() || model.ChangelogExcludeAttribute.IsNull() {
		model.ChangelogExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ChangelogIncludeAttribute.IsUnknown() || model.ChangelogIncludeAttribute.IsNull() {
		model.ChangelogIncludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.SoftDeleteEntryIncludedOperation.IsUnknown() || model.SoftDeleteEntryIncludedOperation.IsNull() {
		model.SoftDeleteEntryIncludedOperation, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.DisabledAlertType.IsUnknown() || model.DisabledAlertType.IsNull() {
		model.DisabledAlertType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ChangelogEntryIncludeBaseDN.IsUnknown() || model.ChangelogEntryIncludeBaseDN.IsNull() {
		model.ChangelogEntryIncludeBaseDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ChangelogIncludeKeyAttribute.IsUnknown() || model.ChangelogIncludeKeyAttribute.IsNull() {
		model.ChangelogIncludeKeyAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.CompactCommonParentDN.IsUnknown() || model.CompactCommonParentDN.IsNull() {
		model.CompactCommonParentDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IndexIncludeAttribute.IsUnknown() || model.IndexIncludeAttribute.IsNull() {
		model.IndexIncludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.SystemIndexToPrimeInternalNodesOnly.IsUnknown() || model.SystemIndexToPrimeInternalNodesOnly.IsNull() {
		model.SystemIndexToPrimeInternalNodesOnly, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.RetentionPolicy.IsUnknown() || model.RetentionPolicy.IsNull() {
		model.RetentionPolicy, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ChangelogDeletedEntryIncludeAttribute.IsUnknown() || model.ChangelogDeletedEntryIncludeAttribute.IsNull() {
		model.ChangelogDeletedEntryIncludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.JeProperty.IsUnknown() || model.JeProperty.IsNull() {
		model.JeProperty, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ChangelogEntryExcludeBaseDN.IsUnknown() || model.ChangelogEntryExcludeBaseDN.IsNull() {
		model.ChangelogEntryExcludeBaseDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ChangelogDeletedEntryExcludeAttribute.IsUnknown() || model.ChangelogDeletedEntryExcludeAttribute.IsNull() {
		model.ChangelogDeletedEntryExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.InsignificantConfigArchiveAttribute.IsUnknown() || model.InsignificantConfigArchiveAttribute.IsNull() {
		model.InsignificantConfigArchiveAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.PrimeMethod.IsUnknown() || model.PrimeMethod.IsNull() {
		model.PrimeMethod, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Read a SchemaBackendResponse object into the model struct
func readSchemaBackendResponseDefault(ctx context.Context, r *client.SchemaBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("schema")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.SchemaEntryDN = internaltypes.GetStringSet(r.SchemaEntryDN)
	state.ShowAllAttributes = types.BoolValue(r.ShowAllAttributes)
	state.ReadOnlySchemaFile = internaltypes.GetStringSet(r.ReadOnlySchemaFile)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, true)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendUnknownValuesDefault(state)
}

// Read a BackupBackendResponse object into the model struct
func readBackupBackendResponseDefault(ctx context.Context, r *client.BackupBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("backup")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.BackupDirectory = internaltypes.GetStringSet(r.BackupDirectory)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendUnknownValuesDefault(state)
}

// Read a EncryptionSettingsBackendResponse object into the model struct
func readEncryptionSettingsBackendResponseDefault(ctx context.Context, r *client.EncryptionSettingsBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("encryption-settings")
	state.Id = types.StringValue(r.Id)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.BackendID = types.StringValue(r.BackendID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, true)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendUnknownValuesDefault(state)
}

// Read a LdifBackendResponse object into the model struct
func readLdifBackendResponseDefault(ctx context.Context, r *client.LdifBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldif")
	state.Id = types.StringValue(r.Id)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.IsPrivateBackend = internaltypes.BoolTypeOrNil(r.IsPrivateBackend)
	state.LdifFile = types.StringValue(r.LdifFile)
	state.BackendID = types.StringValue(r.BackendID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, true)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendUnknownValuesDefault(state)
}

// Read a TrustStoreBackendResponse object into the model struct
func readTrustStoreBackendResponseDefault(ctx context.Context, r *client.TrustStoreBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("trust-store")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.TrustStoreFile = types.StringValue(r.TrustStoreFile)
	state.TrustStoreType = internaltypes.StringTypeOrNil(r.TrustStoreType, true)
	state.TrustStorePinFile = internaltypes.StringTypeOrNil(r.TrustStorePinFile, true)
	state.TrustStorePinPassphraseProvider = internaltypes.StringTypeOrNil(r.TrustStorePinPassphraseProvider, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, true)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendUnknownValuesDefault(state)
}

// Read a CustomBackendResponse object into the model struct
func readCustomBackendResponseDefault(ctx context.Context, r *client.CustomBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("custom")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, true)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendUnknownValuesDefault(state)
}

// Read a ChangelogBackendResponse object into the model struct
func readChangelogBackendResponseDefault(ctx context.Context, r *client.ChangelogBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("changelog")
	state.Id = types.StringValue(r.Id)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.DbDirectory = internaltypes.StringTypeOrNil(r.DbDirectory, true)
	state.DbDirectoryPermissions = internaltypes.StringTypeOrNil(r.DbDirectoryPermissions, true)
	state.DbCachePercent = internaltypes.Int64TypeOrNil(r.DbCachePercent)
	state.JeProperty = internaltypes.GetStringSet(r.JeProperty)
	state.ChangelogWriteBatchSize = internaltypes.Int64TypeOrNil(r.ChangelogWriteBatchSize)
	state.ChangelogPurgeBatchSize = internaltypes.Int64TypeOrNil(r.ChangelogPurgeBatchSize)
	state.ChangelogWriteQueueCapacity = internaltypes.Int64TypeOrNil(r.ChangelogWriteQueueCapacity)
	state.IndexIncludeAttribute = internaltypes.GetStringSet(r.IndexIncludeAttribute)
	state.IndexExcludeAttribute = internaltypes.GetStringSet(r.IndexExcludeAttribute)
	state.ChangelogMaximumAge = types.StringValue(r.ChangelogMaximumAge)
	config.CheckMismatchedPDFormattedAttributes("changelog_maximum_age",
		expectedValues.ChangelogMaximumAge, state.ChangelogMaximumAge, diagnostics)
	state.TargetDatabaseSize = internaltypes.StringTypeOrNil(r.TargetDatabaseSize, true)
	config.CheckMismatchedPDFormattedAttributes("target_database_size",
		expectedValues.TargetDatabaseSize, state.TargetDatabaseSize, diagnostics)
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
		client.StringPointerEnumbackendReportExcludedChangelogAttributesProp(r.ReportExcludedChangelogAttributes), true)
	state.SoftDeleteEntryIncludedOperation = internaltypes.GetStringSet(
		client.StringSliceEnumbackendSoftDeleteEntryIncludedOperationProp(r.SoftDeleteEntryIncludedOperation))
	state.BackendID = types.StringValue(r.BackendID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendUnknownValuesDefault(state)
}

// Read a MonitorBackendResponse object into the model struct
func readMonitorBackendResponseDefault(ctx context.Context, r *client.MonitorBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("monitor")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendUnknownValuesDefault(state)
}

// Read a LocalDbBackendResponse object into the model struct
func readLocalDbBackendResponse(ctx context.Context, r *client.LocalDbBackendResponse, state *backendResourceModel, expectedValues *backendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("local-db")
	state.Id = types.StringValue(r.Id)
	state.UncachedId2entryCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendUncachedId2entryCacheModeProp(r.UncachedId2entryCacheMode), true)
	state.UncachedAttributeCriteria = internaltypes.StringTypeOrNil(r.UncachedAttributeCriteria, internaltypes.IsEmptyString(expectedValues.UncachedAttributeCriteria))
	state.UncachedEntryCriteria = internaltypes.StringTypeOrNil(r.UncachedEntryCriteria, internaltypes.IsEmptyString(expectedValues.UncachedEntryCriteria))
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.SetDegradedAlertForUntrustedIndex = internaltypes.BoolTypeOrNil(r.SetDegradedAlertForUntrustedIndex)
	state.ReturnUnavailableForUntrustedIndex = internaltypes.BoolTypeOrNil(r.ReturnUnavailableForUntrustedIndex)
	state.ProcessFiltersWithUndefinedAttributeTypes = internaltypes.BoolTypeOrNil(r.ProcessFiltersWithUndefinedAttributeTypes)
	state.IsPrivateBackend = internaltypes.BoolTypeOrNil(r.IsPrivateBackend)
	state.DbDirectory = types.StringValue(r.DbDirectory)
	state.DbDirectoryPermissions = internaltypes.StringTypeOrNil(r.DbDirectoryPermissions, true)
	state.CompactCommonParentDN = internaltypes.GetStringSet(r.CompactCommonParentDN)
	state.CompressEntries = internaltypes.BoolTypeOrNil(r.CompressEntries)
	state.HashEntries = internaltypes.BoolTypeOrNil(r.HashEntries)
	state.DbNumCleanerThreads = internaltypes.Int64TypeOrNil(r.DbNumCleanerThreads)
	state.DbCleanerMinUtilization = internaltypes.Int64TypeOrNil(r.DbCleanerMinUtilization)
	state.DbEvictorCriticalPercentage = internaltypes.Int64TypeOrNil(r.DbEvictorCriticalPercentage)
	state.DbCheckpointerWakeupInterval = internaltypes.StringTypeOrNil(r.DbCheckpointerWakeupInterval, true)
	config.CheckMismatchedPDFormattedAttributes("db_checkpointer_wakeup_interval",
		expectedValues.DbCheckpointerWakeupInterval, state.DbCheckpointerWakeupInterval, diagnostics)
	state.DbBackgroundSyncInterval = internaltypes.StringTypeOrNil(r.DbBackgroundSyncInterval, true)
	config.CheckMismatchedPDFormattedAttributes("db_background_sync_interval",
		expectedValues.DbBackgroundSyncInterval, state.DbBackgroundSyncInterval, diagnostics)
	state.DbUseThreadLocalHandles = internaltypes.BoolTypeOrNil(r.DbUseThreadLocalHandles)
	state.DbLogFileMax = internaltypes.StringTypeOrNil(r.DbLogFileMax, true)
	config.CheckMismatchedPDFormattedAttributes("db_log_file_max",
		expectedValues.DbLogFileMax, state.DbLogFileMax, diagnostics)
	state.DbLoggingLevel = internaltypes.StringTypeOrNil(r.DbLoggingLevel, true)
	state.JeProperty = internaltypes.GetStringSet(r.JeProperty)
	state.DbCachePercent = internaltypes.Int64TypeOrNil(r.DbCachePercent)
	state.DefaultCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendDefaultCacheModeProp(r.DefaultCacheMode), true)
	state.Id2entryCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendId2entryCacheModeProp(r.Id2entryCacheMode), internaltypes.IsEmptyString(expectedValues.Id2entryCacheMode))
	state.Dn2idCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendDn2idCacheModeProp(r.Dn2idCacheMode), internaltypes.IsEmptyString(expectedValues.Dn2idCacheMode))
	state.Id2childrenCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendId2childrenCacheModeProp(r.Id2childrenCacheMode), internaltypes.IsEmptyString(expectedValues.Id2childrenCacheMode))
	state.Id2subtreeCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendId2subtreeCacheModeProp(r.Id2subtreeCacheMode), internaltypes.IsEmptyString(expectedValues.Id2subtreeCacheMode))
	state.Dn2uriCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendDn2uriCacheModeProp(r.Dn2uriCacheMode), internaltypes.IsEmptyString(expectedValues.Dn2uriCacheMode))
	state.PrimeMethod = internaltypes.GetStringSet(
		client.StringSliceEnumbackendPrimeMethodProp(r.PrimeMethod))
	state.PrimeThreadCount = internaltypes.Int64TypeOrNil(r.PrimeThreadCount)
	state.PrimeTimeLimit = internaltypes.StringTypeOrNil(r.PrimeTimeLimit, true)
	config.CheckMismatchedPDFormattedAttributes("prime_time_limit",
		expectedValues.PrimeTimeLimit, state.PrimeTimeLimit, diagnostics)
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
		client.StringPointerEnumbackendExternalTxnDefaultBackendLockBehaviorProp(r.ExternalTxnDefaultBackendLockBehavior), true)
	state.SingleWriterLockBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendSingleWriterLockBehaviorProp(r.SingleWriterLockBehavior), true)
	state.SubtreeDeleteSizeLimit = internaltypes.Int64TypeOrNil(r.SubtreeDeleteSizeLimit)
	state.NumRecentChanges = internaltypes.Int64TypeOrNil(r.NumRecentChanges)
	state.OfflineProcessDatabaseOpenTimeout = internaltypes.StringTypeOrNil(r.OfflineProcessDatabaseOpenTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("offline_process_database_open_timeout",
		expectedValues.OfflineProcessDatabaseOpenTimeout, state.OfflineProcessDatabaseOpenTimeout, diagnostics)
	state.BackendID = types.StringValue(r.BackendID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, internaltypes.IsEmptyString(expectedValues.NotificationManager))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendUnknownValues(state)
}

// Read a LocalDbBackendResponse object into the model struct
func readLocalDbBackendResponseDefault(ctx context.Context, r *client.LocalDbBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("local-db")
	state.Id = types.StringValue(r.Id)
	state.UncachedId2entryCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendUncachedId2entryCacheModeProp(r.UncachedId2entryCacheMode), true)
	state.UncachedAttributeCriteria = internaltypes.StringTypeOrNil(r.UncachedAttributeCriteria, true)
	state.UncachedEntryCriteria = internaltypes.StringTypeOrNil(r.UncachedEntryCriteria, true)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.SetDegradedAlertForUntrustedIndex = internaltypes.BoolTypeOrNil(r.SetDegradedAlertForUntrustedIndex)
	state.ReturnUnavailableForUntrustedIndex = internaltypes.BoolTypeOrNil(r.ReturnUnavailableForUntrustedIndex)
	state.ProcessFiltersWithUndefinedAttributeTypes = internaltypes.BoolTypeOrNil(r.ProcessFiltersWithUndefinedAttributeTypes)
	state.IsPrivateBackend = internaltypes.BoolTypeOrNil(r.IsPrivateBackend)
	state.DbDirectory = types.StringValue(r.DbDirectory)
	state.DbDirectoryPermissions = internaltypes.StringTypeOrNil(r.DbDirectoryPermissions, true)
	state.CompactCommonParentDN = internaltypes.GetStringSet(r.CompactCommonParentDN)
	state.CompressEntries = internaltypes.BoolTypeOrNil(r.CompressEntries)
	state.HashEntries = internaltypes.BoolTypeOrNil(r.HashEntries)
	state.DbNumCleanerThreads = internaltypes.Int64TypeOrNil(r.DbNumCleanerThreads)
	state.DbCleanerMinUtilization = internaltypes.Int64TypeOrNil(r.DbCleanerMinUtilization)
	state.DbEvictorCriticalPercentage = internaltypes.Int64TypeOrNil(r.DbEvictorCriticalPercentage)
	state.DbCheckpointerWakeupInterval = internaltypes.StringTypeOrNil(r.DbCheckpointerWakeupInterval, true)
	config.CheckMismatchedPDFormattedAttributes("db_checkpointer_wakeup_interval",
		expectedValues.DbCheckpointerWakeupInterval, state.DbCheckpointerWakeupInterval, diagnostics)
	state.DbBackgroundSyncInterval = internaltypes.StringTypeOrNil(r.DbBackgroundSyncInterval, true)
	config.CheckMismatchedPDFormattedAttributes("db_background_sync_interval",
		expectedValues.DbBackgroundSyncInterval, state.DbBackgroundSyncInterval, diagnostics)
	state.DbUseThreadLocalHandles = internaltypes.BoolTypeOrNil(r.DbUseThreadLocalHandles)
	state.DbLogFileMax = internaltypes.StringTypeOrNil(r.DbLogFileMax, true)
	config.CheckMismatchedPDFormattedAttributes("db_log_file_max",
		expectedValues.DbLogFileMax, state.DbLogFileMax, diagnostics)
	state.DbLoggingLevel = internaltypes.StringTypeOrNil(r.DbLoggingLevel, true)
	state.JeProperty = internaltypes.GetStringSet(r.JeProperty)
	state.DbCachePercent = internaltypes.Int64TypeOrNil(r.DbCachePercent)
	state.DefaultCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendDefaultCacheModeProp(r.DefaultCacheMode), true)
	state.Id2entryCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendId2entryCacheModeProp(r.Id2entryCacheMode), true)
	state.Dn2idCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendDn2idCacheModeProp(r.Dn2idCacheMode), true)
	state.Id2childrenCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendId2childrenCacheModeProp(r.Id2childrenCacheMode), true)
	state.Id2subtreeCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendId2subtreeCacheModeProp(r.Id2subtreeCacheMode), true)
	state.Dn2uriCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendDn2uriCacheModeProp(r.Dn2uriCacheMode), true)
	state.PrimeMethod = internaltypes.GetStringSet(
		client.StringSliceEnumbackendPrimeMethodProp(r.PrimeMethod))
	state.PrimeThreadCount = internaltypes.Int64TypeOrNil(r.PrimeThreadCount)
	state.PrimeTimeLimit = internaltypes.StringTypeOrNil(r.PrimeTimeLimit, true)
	config.CheckMismatchedPDFormattedAttributes("prime_time_limit",
		expectedValues.PrimeTimeLimit, state.PrimeTimeLimit, diagnostics)
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
		client.StringPointerEnumbackendExternalTxnDefaultBackendLockBehaviorProp(r.ExternalTxnDefaultBackendLockBehavior), true)
	state.SingleWriterLockBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendSingleWriterLockBehaviorProp(r.SingleWriterLockBehavior), true)
	state.SubtreeDeleteSizeLimit = internaltypes.Int64TypeOrNil(r.SubtreeDeleteSizeLimit)
	state.NumRecentChanges = internaltypes.Int64TypeOrNil(r.NumRecentChanges)
	state.OfflineProcessDatabaseOpenTimeout = internaltypes.StringTypeOrNil(r.OfflineProcessDatabaseOpenTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("offline_process_database_open_timeout",
		expectedValues.OfflineProcessDatabaseOpenTimeout, state.OfflineProcessDatabaseOpenTimeout, diagnostics)
	state.BackendID = types.StringValue(r.BackendID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendUnknownValuesDefault(state)
}

// Read a ConfigFileHandlerBackendResponse object into the model struct
func readConfigFileHandlerBackendResponseDefault(ctx context.Context, r *client.ConfigFileHandlerBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("config-file-handler")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.InsignificantConfigArchiveAttribute = internaltypes.GetStringSet(r.InsignificantConfigArchiveAttribute)
	state.InsignificantConfigArchiveBaseDN = internaltypes.GetStringSet(r.InsignificantConfigArchiveBaseDN)
	state.MaintainConfigArchive = internaltypes.BoolTypeOrNil(r.MaintainConfigArchive)
	state.MaxConfigArchiveCount = internaltypes.Int64TypeOrNil(r.MaxConfigArchiveCount)
	state.MirroredSubtreePeerPollingInterval = internaltypes.StringTypeOrNil(r.MirroredSubtreePeerPollingInterval, true)
	config.CheckMismatchedPDFormattedAttributes("mirrored_subtree_peer_polling_interval",
		expectedValues.MirroredSubtreePeerPollingInterval, state.MirroredSubtreePeerPollingInterval, diagnostics)
	state.MirroredSubtreeEntryUpdateTimeout = internaltypes.StringTypeOrNil(r.MirroredSubtreeEntryUpdateTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("mirrored_subtree_entry_update_timeout",
		expectedValues.MirroredSubtreeEntryUpdateTimeout, state.MirroredSubtreeEntryUpdateTimeout, diagnostics)
	state.MirroredSubtreeSearchTimeout = internaltypes.StringTypeOrNil(r.MirroredSubtreeSearchTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("mirrored_subtree_search_timeout",
		expectedValues.MirroredSubtreeSearchTimeout, state.MirroredSubtreeSearchTimeout, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, true)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendUnknownValuesDefault(state)
}

// Read a TaskBackendResponse object into the model struct
func readTaskBackendResponseDefault(ctx context.Context, r *client.TaskBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("task")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.TaskBackingFile = types.StringValue(r.TaskBackingFile)
	state.MaximumInitialTaskLogMessagesToRetain = internaltypes.Int64TypeOrNil(r.MaximumInitialTaskLogMessagesToRetain)
	state.MaximumFinalTaskLogMessagesToRetain = internaltypes.Int64TypeOrNil(r.MaximumFinalTaskLogMessagesToRetain)
	state.TaskRetentionTime = internaltypes.StringTypeOrNil(r.TaskRetentionTime, true)
	config.CheckMismatchedPDFormattedAttributes("task_retention_time",
		expectedValues.TaskRetentionTime, state.TaskRetentionTime, diagnostics)
	state.NotificationSenderAddress = internaltypes.StringTypeOrNil(r.NotificationSenderAddress, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, true)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendUnknownValuesDefault(state)
}

// Read a AlertBackendResponse object into the model struct
func readAlertBackendResponseDefault(ctx context.Context, r *client.AlertBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("alert")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.LdifFile = types.StringValue(r.LdifFile)
	state.AlertRetentionTime = types.StringValue(r.AlertRetentionTime)
	config.CheckMismatchedPDFormattedAttributes("alert_retention_time",
		expectedValues.AlertRetentionTime, state.AlertRetentionTime, diagnostics)
	state.MaxAlerts = internaltypes.Int64TypeOrNil(r.MaxAlerts)
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumbackendDisabledAlertTypeProp(r.DisabledAlertType))
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, true)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendUnknownValuesDefault(state)
}

// Read a AlarmBackendResponse object into the model struct
func readAlarmBackendResponseDefault(ctx context.Context, r *client.AlarmBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("alarm")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.LdifFile = types.StringValue(r.LdifFile)
	state.AlarmRetentionTime = types.StringValue(r.AlarmRetentionTime)
	config.CheckMismatchedPDFormattedAttributes("alarm_retention_time",
		expectedValues.AlarmRetentionTime, state.AlarmRetentionTime, diagnostics)
	state.MaxAlarms = internaltypes.Int64TypeOrNil(r.MaxAlarms)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, true)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendUnknownValuesDefault(state)
}

// Read a MetricsBackendResponse object into the model struct
func readMetricsBackendResponseDefault(ctx context.Context, r *client.MetricsBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("metrics")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.StorageDir = types.StringValue(r.StorageDir)
	state.MetricsDir = types.StringValue(r.MetricsDir)
	state.SampleFlushInterval = internaltypes.StringTypeOrNil(r.SampleFlushInterval, true)
	config.CheckMismatchedPDFormattedAttributes("sample_flush_interval",
		expectedValues.SampleFlushInterval, state.SampleFlushInterval, diagnostics)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendUnknownValuesDefault(state)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *defaultBackendResourceModel) setStateValuesNotReturnedByAPI(expectedValues *defaultBackendResourceModel) {
	if !expectedValues.TrustStorePin.IsUnknown() {
		state.TrustStorePin = expectedValues.TrustStorePin
	}
}

// Create any update operations necessary to make the state match the plan
func createBackendOperations(plan backendResourceModel, state backendResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.UncachedId2entryCacheMode, state.UncachedId2entryCacheMode, "uncached-id2entry-cache-mode")
	operations.AddStringOperationIfNecessary(&ops, plan.UncachedAttributeCriteria, state.UncachedAttributeCriteria, "uncached-attribute-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.UncachedEntryCriteria, state.UncachedEntryCriteria, "uncached-entry-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.BackendID, state.BackendID, "backend-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.SetDegradedAlertForUntrustedIndex, state.SetDegradedAlertForUntrustedIndex, "set-degraded-alert-for-untrusted-index")
	operations.AddBoolOperationIfNecessary(&ops, plan.ReturnUnavailableForUntrustedIndex, state.ReturnUnavailableForUntrustedIndex, "return-unavailable-for-untrusted-index")
	operations.AddBoolOperationIfNecessary(&ops, plan.ProcessFiltersWithUndefinedAttributeTypes, state.ProcessFiltersWithUndefinedAttributeTypes, "process-filters-with-undefined-attribute-types")
	operations.AddStringOperationIfNecessary(&ops, plan.DbDirectory, state.DbDirectory, "db-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.DbDirectoryPermissions, state.DbDirectoryPermissions, "db-directory-permissions")
	operations.AddInt64OperationIfNecessary(&ops, plan.DbCachePercent, state.DbCachePercent, "db-cache-percent")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.CompactCommonParentDN, state.CompactCommonParentDN, "compact-common-parent-dn")
	operations.AddBoolOperationIfNecessary(&ops, plan.CompressEntries, state.CompressEntries, "compress-entries")
	operations.AddBoolOperationIfNecessary(&ops, plan.HashEntries, state.HashEntries, "hash-entries")
	operations.AddInt64OperationIfNecessary(&ops, plan.DbNumCleanerThreads, state.DbNumCleanerThreads, "db-num-cleaner-threads")
	operations.AddInt64OperationIfNecessary(&ops, plan.DbCleanerMinUtilization, state.DbCleanerMinUtilization, "db-cleaner-min-utilization")
	operations.AddInt64OperationIfNecessary(&ops, plan.DbEvictorCriticalPercentage, state.DbEvictorCriticalPercentage, "db-evictor-critical-percentage")
	operations.AddStringOperationIfNecessary(&ops, plan.DbCheckpointerWakeupInterval, state.DbCheckpointerWakeupInterval, "db-checkpointer-wakeup-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.DbBackgroundSyncInterval, state.DbBackgroundSyncInterval, "db-background-sync-interval")
	operations.AddBoolOperationIfNecessary(&ops, plan.DbUseThreadLocalHandles, state.DbUseThreadLocalHandles, "db-use-thread-local-handles")
	operations.AddStringOperationIfNecessary(&ops, plan.DbLogFileMax, state.DbLogFileMax, "db-log-file-max")
	operations.AddStringOperationIfNecessary(&ops, plan.DbLoggingLevel, state.DbLoggingLevel, "db-logging-level")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.JeProperty, state.JeProperty, "je-property")
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultCacheMode, state.DefaultCacheMode, "default-cache-mode")
	operations.AddStringOperationIfNecessary(&ops, plan.Id2entryCacheMode, state.Id2entryCacheMode, "id2entry-cache-mode")
	operations.AddStringOperationIfNecessary(&ops, plan.Dn2idCacheMode, state.Dn2idCacheMode, "dn2id-cache-mode")
	operations.AddStringOperationIfNecessary(&ops, plan.Id2childrenCacheMode, state.Id2childrenCacheMode, "id2children-cache-mode")
	operations.AddStringOperationIfNecessary(&ops, plan.Id2subtreeCacheMode, state.Id2subtreeCacheMode, "id2subtree-cache-mode")
	operations.AddStringOperationIfNecessary(&ops, plan.Dn2uriCacheMode, state.Dn2uriCacheMode, "dn2uri-cache-mode")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.PrimeMethod, state.PrimeMethod, "prime-method")
	operations.AddInt64OperationIfNecessary(&ops, plan.PrimeThreadCount, state.PrimeThreadCount, "prime-thread-count")
	operations.AddStringOperationIfNecessary(&ops, plan.PrimeTimeLimit, state.PrimeTimeLimit, "prime-time-limit")
	operations.AddBoolOperationIfNecessary(&ops, plan.PrimeAllIndexes, state.PrimeAllIndexes, "prime-all-indexes")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SystemIndexToPrime, state.SystemIndexToPrime, "system-index-to-prime")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SystemIndexToPrimeInternalNodesOnly, state.SystemIndexToPrimeInternalNodesOnly, "system-index-to-prime-internal-nodes-only")
	operations.AddBoolOperationIfNecessary(&ops, plan.BackgroundPrime, state.BackgroundPrime, "background-prime")
	operations.AddInt64OperationIfNecessary(&ops, plan.IndexEntryLimit, state.IndexEntryLimit, "index-entry-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.CompositeIndexEntryLimit, state.CompositeIndexEntryLimit, "composite-index-entry-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.Id2childrenIndexEntryLimit, state.Id2childrenIndexEntryLimit, "id2children-index-entry-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.Id2subtreeIndexEntryLimit, state.Id2subtreeIndexEntryLimit, "id2subtree-index-entry-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.ImportTempDirectory, state.ImportTempDirectory, "import-temp-directory")
	operations.AddInt64OperationIfNecessary(&ops, plan.ImportThreadCount, state.ImportThreadCount, "import-thread-count")
	operations.AddInt64OperationIfNecessary(&ops, plan.ExportThreadCount, state.ExportThreadCount, "export-thread-count")
	operations.AddInt64OperationIfNecessary(&ops, plan.DbImportCachePercent, state.DbImportCachePercent, "db-import-cache-percent")
	operations.AddBoolOperationIfNecessary(&ops, plan.DbTxnWriteNoSync, state.DbTxnWriteNoSync, "db-txn-write-no-sync")
	operations.AddInt64OperationIfNecessary(&ops, plan.DeadlockRetryLimit, state.DeadlockRetryLimit, "deadlock-retry-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.ExternalTxnDefaultBackendLockBehavior, state.ExternalTxnDefaultBackendLockBehavior, "external-txn-default-backend-lock-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.SingleWriterLockBehavior, state.SingleWriterLockBehavior, "single-writer-lock-behavior")
	operations.AddInt64OperationIfNecessary(&ops, plan.SubtreeDeleteSizeLimit, state.SubtreeDeleteSizeLimit, "subtree-delete-size-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumRecentChanges, state.NumRecentChanges, "num-recent-changes")
	operations.AddStringOperationIfNecessary(&ops, plan.OfflineProcessDatabaseOpenTimeout, state.OfflineProcessDatabaseOpenTimeout, "offline-process-database-open-timeout")
	operations.AddBoolOperationIfNecessary(&ops, plan.IsPrivateBackend, state.IsPrivateBackend, "is-private-backend")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.WritabilityMode, state.WritabilityMode, "writability-mode")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.SetDegradedAlertWhenDisabled, state.SetDegradedAlertWhenDisabled, "set-degraded-alert-when-disabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.ReturnUnavailableWhenDisabled, state.ReturnUnavailableWhenDisabled, "return-unavailable-when-disabled")
	operations.AddStringOperationIfNecessary(&ops, plan.NotificationManager, state.NotificationManager, "notification-manager")
	return ops
}

// Create any update operations necessary to make the state match the plan
func createBackendOperationsDefault(plan defaultBackendResourceModel, state defaultBackendResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.UncachedId2entryCacheMode, state.UncachedId2entryCacheMode, "uncached-id2entry-cache-mode")
	operations.AddStringOperationIfNecessary(&ops, plan.StorageDir, state.StorageDir, "storage-dir")
	operations.AddStringOperationIfNecessary(&ops, plan.MetricsDir, state.MetricsDir, "metrics-dir")
	operations.AddStringOperationIfNecessary(&ops, plan.SampleFlushInterval, state.SampleFlushInterval, "sample-flush-interval")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RetentionPolicy, state.RetentionPolicy, "retention-policy")
	operations.AddStringOperationIfNecessary(&ops, plan.UncachedAttributeCriteria, state.UncachedAttributeCriteria, "uncached-attribute-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.UncachedEntryCriteria, state.UncachedEntryCriteria, "uncached-entry-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.AlarmRetentionTime, state.AlarmRetentionTime, "alarm-retention-time")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxAlarms, state.MaxAlarms, "max-alarms")
	operations.AddStringOperationIfNecessary(&ops, plan.AlertRetentionTime, state.AlertRetentionTime, "alert-retention-time")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxAlerts, state.MaxAlerts, "max-alerts")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DisabledAlertType, state.DisabledAlertType, "disabled-alert-type")
	operations.AddStringOperationIfNecessary(&ops, plan.TaskBackingFile, state.TaskBackingFile, "task-backing-file")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumInitialTaskLogMessagesToRetain, state.MaximumInitialTaskLogMessagesToRetain, "maximum-initial-task-log-messages-to-retain")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumFinalTaskLogMessagesToRetain, state.MaximumFinalTaskLogMessagesToRetain, "maximum-final-task-log-messages-to-retain")
	operations.AddStringOperationIfNecessary(&ops, plan.TaskRetentionTime, state.TaskRetentionTime, "task-retention-time")
	operations.AddStringOperationIfNecessary(&ops, plan.NotificationSenderAddress, state.NotificationSenderAddress, "notification-sender-address")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.InsignificantConfigArchiveAttribute, state.InsignificantConfigArchiveAttribute, "insignificant-config-archive-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.InsignificantConfigArchiveBaseDN, state.InsignificantConfigArchiveBaseDN, "insignificant-config-archive-base-dn")
	operations.AddBoolOperationIfNecessary(&ops, plan.MaintainConfigArchive, state.MaintainConfigArchive, "maintain-config-archive")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxConfigArchiveCount, state.MaxConfigArchiveCount, "max-config-archive-count")
	operations.AddStringOperationIfNecessary(&ops, plan.MirroredSubtreePeerPollingInterval, state.MirroredSubtreePeerPollingInterval, "mirrored-subtree-peer-polling-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.MirroredSubtreeEntryUpdateTimeout, state.MirroredSubtreeEntryUpdateTimeout, "mirrored-subtree-entry-update-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.MirroredSubtreeSearchTimeout, state.MirroredSubtreeSearchTimeout, "mirrored-subtree-search-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.BackendID, state.BackendID, "backend-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.SetDegradedAlertForUntrustedIndex, state.SetDegradedAlertForUntrustedIndex, "set-degraded-alert-for-untrusted-index")
	operations.AddBoolOperationIfNecessary(&ops, plan.ReturnUnavailableForUntrustedIndex, state.ReturnUnavailableForUntrustedIndex, "return-unavailable-for-untrusted-index")
	operations.AddBoolOperationIfNecessary(&ops, plan.ProcessFiltersWithUndefinedAttributeTypes, state.ProcessFiltersWithUndefinedAttributeTypes, "process-filters-with-undefined-attribute-types")
	operations.AddStringOperationIfNecessary(&ops, plan.DbDirectory, state.DbDirectory, "db-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.DbDirectoryPermissions, state.DbDirectoryPermissions, "db-directory-permissions")
	operations.AddInt64OperationIfNecessary(&ops, plan.DbCachePercent, state.DbCachePercent, "db-cache-percent")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.CompactCommonParentDN, state.CompactCommonParentDN, "compact-common-parent-dn")
	operations.AddBoolOperationIfNecessary(&ops, plan.CompressEntries, state.CompressEntries, "compress-entries")
	operations.AddBoolOperationIfNecessary(&ops, plan.HashEntries, state.HashEntries, "hash-entries")
	operations.AddInt64OperationIfNecessary(&ops, plan.DbNumCleanerThreads, state.DbNumCleanerThreads, "db-num-cleaner-threads")
	operations.AddInt64OperationIfNecessary(&ops, plan.DbCleanerMinUtilization, state.DbCleanerMinUtilization, "db-cleaner-min-utilization")
	operations.AddInt64OperationIfNecessary(&ops, plan.DbEvictorCriticalPercentage, state.DbEvictorCriticalPercentage, "db-evictor-critical-percentage")
	operations.AddStringOperationIfNecessary(&ops, plan.DbCheckpointerWakeupInterval, state.DbCheckpointerWakeupInterval, "db-checkpointer-wakeup-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.DbBackgroundSyncInterval, state.DbBackgroundSyncInterval, "db-background-sync-interval")
	operations.AddBoolOperationIfNecessary(&ops, plan.DbUseThreadLocalHandles, state.DbUseThreadLocalHandles, "db-use-thread-local-handles")
	operations.AddStringOperationIfNecessary(&ops, plan.DbLogFileMax, state.DbLogFileMax, "db-log-file-max")
	operations.AddStringOperationIfNecessary(&ops, plan.DbLoggingLevel, state.DbLoggingLevel, "db-logging-level")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.JeProperty, state.JeProperty, "je-property")
	operations.AddInt64OperationIfNecessary(&ops, plan.ChangelogWriteBatchSize, state.ChangelogWriteBatchSize, "changelog-write-batch-size")
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultCacheMode, state.DefaultCacheMode, "default-cache-mode")
	operations.AddStringOperationIfNecessary(&ops, plan.Id2entryCacheMode, state.Id2entryCacheMode, "id2entry-cache-mode")
	operations.AddStringOperationIfNecessary(&ops, plan.Dn2idCacheMode, state.Dn2idCacheMode, "dn2id-cache-mode")
	operations.AddStringOperationIfNecessary(&ops, plan.Id2childrenCacheMode, state.Id2childrenCacheMode, "id2children-cache-mode")
	operations.AddStringOperationIfNecessary(&ops, plan.Id2subtreeCacheMode, state.Id2subtreeCacheMode, "id2subtree-cache-mode")
	operations.AddStringOperationIfNecessary(&ops, plan.Dn2uriCacheMode, state.Dn2uriCacheMode, "dn2uri-cache-mode")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.PrimeMethod, state.PrimeMethod, "prime-method")
	operations.AddInt64OperationIfNecessary(&ops, plan.PrimeThreadCount, state.PrimeThreadCount, "prime-thread-count")
	operations.AddStringOperationIfNecessary(&ops, plan.PrimeTimeLimit, state.PrimeTimeLimit, "prime-time-limit")
	operations.AddBoolOperationIfNecessary(&ops, plan.PrimeAllIndexes, state.PrimeAllIndexes, "prime-all-indexes")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SystemIndexToPrime, state.SystemIndexToPrime, "system-index-to-prime")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SystemIndexToPrimeInternalNodesOnly, state.SystemIndexToPrimeInternalNodesOnly, "system-index-to-prime-internal-nodes-only")
	operations.AddBoolOperationIfNecessary(&ops, plan.BackgroundPrime, state.BackgroundPrime, "background-prime")
	operations.AddInt64OperationIfNecessary(&ops, plan.IndexEntryLimit, state.IndexEntryLimit, "index-entry-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.CompositeIndexEntryLimit, state.CompositeIndexEntryLimit, "composite-index-entry-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.Id2childrenIndexEntryLimit, state.Id2childrenIndexEntryLimit, "id2children-index-entry-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.Id2subtreeIndexEntryLimit, state.Id2subtreeIndexEntryLimit, "id2subtree-index-entry-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.ImportTempDirectory, state.ImportTempDirectory, "import-temp-directory")
	operations.AddInt64OperationIfNecessary(&ops, plan.ImportThreadCount, state.ImportThreadCount, "import-thread-count")
	operations.AddInt64OperationIfNecessary(&ops, plan.ExportThreadCount, state.ExportThreadCount, "export-thread-count")
	operations.AddInt64OperationIfNecessary(&ops, plan.DbImportCachePercent, state.DbImportCachePercent, "db-import-cache-percent")
	operations.AddBoolOperationIfNecessary(&ops, plan.DbTxnWriteNoSync, state.DbTxnWriteNoSync, "db-txn-write-no-sync")
	operations.AddInt64OperationIfNecessary(&ops, plan.DeadlockRetryLimit, state.DeadlockRetryLimit, "deadlock-retry-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.ExternalTxnDefaultBackendLockBehavior, state.ExternalTxnDefaultBackendLockBehavior, "external-txn-default-backend-lock-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.SingleWriterLockBehavior, state.SingleWriterLockBehavior, "single-writer-lock-behavior")
	operations.AddInt64OperationIfNecessary(&ops, plan.SubtreeDeleteSizeLimit, state.SubtreeDeleteSizeLimit, "subtree-delete-size-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumRecentChanges, state.NumRecentChanges, "num-recent-changes")
	operations.AddStringOperationIfNecessary(&ops, plan.OfflineProcessDatabaseOpenTimeout, state.OfflineProcessDatabaseOpenTimeout, "offline-process-database-open-timeout")
	operations.AddInt64OperationIfNecessary(&ops, plan.ChangelogPurgeBatchSize, state.ChangelogPurgeBatchSize, "changelog-purge-batch-size")
	operations.AddInt64OperationIfNecessary(&ops, plan.ChangelogWriteQueueCapacity, state.ChangelogWriteQueueCapacity, "changelog-write-queue-capacity")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IndexIncludeAttribute, state.IndexIncludeAttribute, "index-include-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IndexExcludeAttribute, state.IndexExcludeAttribute, "index-exclude-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.ChangelogMaximumAge, state.ChangelogMaximumAge, "changelog-maximum-age")
	operations.AddStringOperationIfNecessary(&ops, plan.TargetDatabaseSize, state.TargetDatabaseSize, "target-database-size")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangelogEntryIncludeBaseDN, state.ChangelogEntryIncludeBaseDN, "changelog-entry-include-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangelogEntryExcludeBaseDN, state.ChangelogEntryExcludeBaseDN, "changelog-entry-exclude-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangelogEntryIncludeFilter, state.ChangelogEntryIncludeFilter, "changelog-entry-include-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangelogEntryExcludeFilter, state.ChangelogEntryExcludeFilter, "changelog-entry-exclude-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangelogIncludeAttribute, state.ChangelogIncludeAttribute, "changelog-include-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangelogExcludeAttribute, state.ChangelogExcludeAttribute, "changelog-exclude-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangelogDeletedEntryIncludeAttribute, state.ChangelogDeletedEntryIncludeAttribute, "changelog-deleted-entry-include-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangelogDeletedEntryExcludeAttribute, state.ChangelogDeletedEntryExcludeAttribute, "changelog-deleted-entry-exclude-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangelogIncludeKeyAttribute, state.ChangelogIncludeKeyAttribute, "changelog-include-key-attribute")
	operations.AddInt64OperationIfNecessary(&ops, plan.ChangelogMaxBeforeAfterValues, state.ChangelogMaxBeforeAfterValues, "changelog-max-before-after-values")
	operations.AddBoolOperationIfNecessary(&ops, plan.WriteLastmodAttributes, state.WriteLastmodAttributes, "write-lastmod-attributes")
	operations.AddBoolOperationIfNecessary(&ops, plan.UseReversibleForm, state.UseReversibleForm, "use-reversible-form")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeVirtualAttributes, state.IncludeVirtualAttributes, "include-virtual-attributes")
	operations.AddBoolOperationIfNecessary(&ops, plan.ApplyAccessControlsToChangelogEntryContents, state.ApplyAccessControlsToChangelogEntryContents, "apply-access-controls-to-changelog-entry-contents")
	operations.AddStringOperationIfNecessary(&ops, plan.ReportExcludedChangelogAttributes, state.ReportExcludedChangelogAttributes, "report-excluded-changelog-attributes")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SoftDeleteEntryIncludedOperation, state.SoftDeleteEntryIncludedOperation, "soft-delete-entry-included-operation")
	operations.AddBoolOperationIfNecessary(&ops, plan.IsPrivateBackend, state.IsPrivateBackend, "is-private-backend")
	operations.AddStringOperationIfNecessary(&ops, plan.LdifFile, state.LdifFile, "ldif-file")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStoreFile, state.TrustStoreFile, "trust-store-file")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStoreType, state.TrustStoreType, "trust-store-type")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStorePin, state.TrustStorePin, "trust-store-pin")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStorePinFile, state.TrustStorePinFile, "trust-store-pin-file")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStorePinPassphraseProvider, state.TrustStorePinPassphraseProvider, "trust-store-pin-passphrase-provider")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.WritabilityMode, state.WritabilityMode, "writability-mode")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.BackupDirectory, state.BackupDirectory, "backup-directory")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SchemaEntryDN, state.SchemaEntryDN, "schema-entry-dn")
	operations.AddBoolOperationIfNecessary(&ops, plan.ShowAllAttributes, state.ShowAllAttributes, "show-all-attributes")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ReadOnlySchemaFile, state.ReadOnlySchemaFile, "read-only-schema-file")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.SetDegradedAlertWhenDisabled, state.SetDegradedAlertWhenDisabled, "set-degraded-alert-when-disabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.ReturnUnavailableWhenDisabled, state.ReturnUnavailableWhenDisabled, "return-unavailable-when-disabled")
	operations.AddStringOperationIfNecessary(&ops, plan.BackupFilePermissions, state.BackupFilePermissions, "backup-file-permissions")
	operations.AddStringOperationIfNecessary(&ops, plan.NotificationManager, state.NotificationManager, "notification-manager")
	return ops
}

// Create a local-db backend
func (r *backendResource) CreateLocalDbBackend(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan backendResourceModel) (*backendResourceModel, error) {
	var BaseDNSlice []string
	plan.BaseDN.ElementsAs(ctx, &BaseDNSlice, false)
	addRequest := client.NewAddLocalDbBackendRequest([]client.EnumlocalDbBackendSchemaUrn{client.ENUMLOCALDBBACKENDSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0BACKENDLOCAL_DB},
		plan.BackendID.ValueString(),
		plan.Enabled.ValueBool(),
		BaseDNSlice,
		plan.BackendID.ValueString())
	err := addOptionalLocalDbBackendFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Backend", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.BackendAPI.AddBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLocalDbBackendRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.BackendAPI.AddBackendExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Backend", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state backendResourceModel
	readLocalDbBackendResponse(ctx, addResponse.LocalDbBackendResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *backendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan backendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateLocalDbBackend(ctx, req, resp, plan)
	if err != nil {
		return
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
func (r *defaultBackendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan defaultBackendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.BackendAPI.GetBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.BackendID.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Backend", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state defaultBackendResourceModel
	if readResponse.SchemaBackendResponse != nil {
		readSchemaBackendResponseDefault(ctx, readResponse.SchemaBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.BackupBackendResponse != nil {
		readBackupBackendResponseDefault(ctx, readResponse.BackupBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.EncryptionSettingsBackendResponse != nil {
		readEncryptionSettingsBackendResponseDefault(ctx, readResponse.EncryptionSettingsBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LdifBackendResponse != nil {
		readLdifBackendResponseDefault(ctx, readResponse.LdifBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.TrustStoreBackendResponse != nil {
		readTrustStoreBackendResponseDefault(ctx, readResponse.TrustStoreBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CustomBackendResponse != nil {
		readCustomBackendResponseDefault(ctx, readResponse.CustomBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ChangelogBackendResponse != nil {
		readChangelogBackendResponseDefault(ctx, readResponse.ChangelogBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.MonitorBackendResponse != nil {
		readMonitorBackendResponseDefault(ctx, readResponse.MonitorBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LocalDbBackendResponse != nil {
		readLocalDbBackendResponseDefault(ctx, readResponse.LocalDbBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConfigFileHandlerBackendResponse != nil {
		readConfigFileHandlerBackendResponseDefault(ctx, readResponse.ConfigFileHandlerBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.TaskBackendResponse != nil {
		readTaskBackendResponseDefault(ctx, readResponse.TaskBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AlertBackendResponse != nil {
		readAlertBackendResponseDefault(ctx, readResponse.AlertBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AlarmBackendResponse != nil {
		readAlarmBackendResponseDefault(ctx, readResponse.AlarmBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.MetricsBackendResponse != nil {
		readMetricsBackendResponseDefault(ctx, readResponse.MetricsBackendResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.BackendAPI.UpdateBackend(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.BackendID.ValueString())
	ops := createBackendOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.BackendAPI.UpdateBackendExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Backend", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.SchemaBackendResponse != nil {
			readSchemaBackendResponseDefault(ctx, updateResponse.SchemaBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.BackupBackendResponse != nil {
			readBackupBackendResponseDefault(ctx, updateResponse.BackupBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.EncryptionSettingsBackendResponse != nil {
			readEncryptionSettingsBackendResponseDefault(ctx, updateResponse.EncryptionSettingsBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LdifBackendResponse != nil {
			readLdifBackendResponseDefault(ctx, updateResponse.LdifBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.TrustStoreBackendResponse != nil {
			readTrustStoreBackendResponseDefault(ctx, updateResponse.TrustStoreBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CustomBackendResponse != nil {
			readCustomBackendResponseDefault(ctx, updateResponse.CustomBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ChangelogBackendResponse != nil {
			readChangelogBackendResponseDefault(ctx, updateResponse.ChangelogBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.MonitorBackendResponse != nil {
			readMonitorBackendResponseDefault(ctx, updateResponse.MonitorBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LocalDbBackendResponse != nil {
			readLocalDbBackendResponseDefault(ctx, updateResponse.LocalDbBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConfigFileHandlerBackendResponse != nil {
			readConfigFileHandlerBackendResponseDefault(ctx, updateResponse.ConfigFileHandlerBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.TaskBackendResponse != nil {
			readTaskBackendResponseDefault(ctx, updateResponse.TaskBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AlertBackendResponse != nil {
			readAlertBackendResponseDefault(ctx, updateResponse.AlertBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AlarmBackendResponse != nil {
			readAlarmBackendResponseDefault(ctx, updateResponse.AlarmBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.MetricsBackendResponse != nil {
			readMetricsBackendResponseDefault(ctx, updateResponse.MetricsBackendResponse, &state, &plan, &resp.Diagnostics)
		}
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *backendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state backendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.BackendAPI.GetBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.BackendID.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Backend", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Backend", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.LocalDbBackendResponse != nil {
		readLocalDbBackendResponse(ctx, readResponse.LocalDbBackendResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *defaultBackendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state defaultBackendResourceModel
	diags := req.State.Get(ctx, &state)
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
		readSchemaBackendResponseDefault(ctx, readResponse.SchemaBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.BackupBackendResponse != nil {
		readBackupBackendResponseDefault(ctx, readResponse.BackupBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.EncryptionSettingsBackendResponse != nil {
		readEncryptionSettingsBackendResponseDefault(ctx, readResponse.EncryptionSettingsBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LdifBackendResponse != nil {
		readLdifBackendResponseDefault(ctx, readResponse.LdifBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.TrustStoreBackendResponse != nil {
		readTrustStoreBackendResponseDefault(ctx, readResponse.TrustStoreBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CustomBackendResponse != nil {
		readCustomBackendResponseDefault(ctx, readResponse.CustomBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ChangelogBackendResponse != nil {
		readChangelogBackendResponseDefault(ctx, readResponse.ChangelogBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.MonitorBackendResponse != nil {
		readMonitorBackendResponseDefault(ctx, readResponse.MonitorBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConfigFileHandlerBackendResponse != nil {
		readConfigFileHandlerBackendResponseDefault(ctx, readResponse.ConfigFileHandlerBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.TaskBackendResponse != nil {
		readTaskBackendResponseDefault(ctx, readResponse.TaskBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AlertBackendResponse != nil {
		readAlertBackendResponseDefault(ctx, readResponse.AlertBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AlarmBackendResponse != nil {
		readAlarmBackendResponseDefault(ctx, readResponse.AlarmBackendResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.MetricsBackendResponse != nil {
		readMetricsBackendResponseDefault(ctx, readResponse.MetricsBackendResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *backendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan backendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state backendResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.BackendAPI.UpdateBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.BackendID.ValueString())

	// Determine what update operations are necessary
	ops := createBackendOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.BackendAPI.UpdateBackendExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Backend", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.LocalDbBackendResponse != nil {
			readLocalDbBackendResponse(ctx, updateResponse.LocalDbBackendResponse, &state, &plan, &resp.Diagnostics)
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

func (r *defaultBackendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan defaultBackendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state defaultBackendResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.BackendAPI.UpdateBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.BackendID.ValueString())

	// Determine what update operations are necessary
	ops := createBackendOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.BackendAPI.UpdateBackendExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Backend", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.SchemaBackendResponse != nil {
			readSchemaBackendResponseDefault(ctx, updateResponse.SchemaBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.BackupBackendResponse != nil {
			readBackupBackendResponseDefault(ctx, updateResponse.BackupBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.EncryptionSettingsBackendResponse != nil {
			readEncryptionSettingsBackendResponseDefault(ctx, updateResponse.EncryptionSettingsBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LdifBackendResponse != nil {
			readLdifBackendResponseDefault(ctx, updateResponse.LdifBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.TrustStoreBackendResponse != nil {
			readTrustStoreBackendResponseDefault(ctx, updateResponse.TrustStoreBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CustomBackendResponse != nil {
			readCustomBackendResponseDefault(ctx, updateResponse.CustomBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ChangelogBackendResponse != nil {
			readChangelogBackendResponseDefault(ctx, updateResponse.ChangelogBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.MonitorBackendResponse != nil {
			readMonitorBackendResponseDefault(ctx, updateResponse.MonitorBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LocalDbBackendResponse != nil {
			readLocalDbBackendResponseDefault(ctx, updateResponse.LocalDbBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConfigFileHandlerBackendResponse != nil {
			readConfigFileHandlerBackendResponseDefault(ctx, updateResponse.ConfigFileHandlerBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.TaskBackendResponse != nil {
			readTaskBackendResponseDefault(ctx, updateResponse.TaskBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AlertBackendResponse != nil {
			readAlertBackendResponseDefault(ctx, updateResponse.AlertBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AlarmBackendResponse != nil {
			readAlarmBackendResponseDefault(ctx, updateResponse.AlarmBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.MetricsBackendResponse != nil {
			readMetricsBackendResponseDefault(ctx, updateResponse.MetricsBackendResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultBackendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *backendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state backendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.BackendAPI.DeleteBackendExecute(r.apiClient.BackendAPI.DeleteBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.BackendID.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Backend", err, httpResp)
		return
	}
}

func (r *backendResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importBackend(ctx, req, resp)
}

func (r *defaultBackendResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importBackend(ctx, req, resp)
}

func importBackend(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to backend_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("backend_id"), req, resp)
}
