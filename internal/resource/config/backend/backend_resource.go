package backend

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
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
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultBackendResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type backendResourceModel struct {
	Id                                        types.String `tfsdk:"id"`
	LastUpdated                               types.String `tfsdk:"last_updated"`
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
	LastUpdated                                 types.String `tfsdk:"last_updated"`
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"set_degraded_alert_for_untrusted_index": schema.BoolAttribute{
				Description: "Determines whether the Directory Server enters a DEGRADED state when this Local DB Backend has an index whose contents cannot be trusted.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"return_unavailable_for_untrusted_index": schema.BoolAttribute{
				Description: "Determines whether the Directory Server returns UNAVAILABLE for any LDAP search operation in this Local DB Backend that would use an index whose contents cannot be trusted.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"process_filters_with_undefined_attribute_types": schema.BoolAttribute{
				Description: "Determines whether the Directory Server should continue filter processing for LDAP search operations in this Local DB Backend that includes a search filter with an attribute that is not defined in the schema. This will only apply if check-schema is enabled in the global configuration.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"db_directory": schema.StringAttribute{
				Description: "Specifies the path to the filesystem directory that is used to hold the Berkeley DB Java Edition database files containing the data for this backend. The files for this backend are stored in a sub-directory named after the backend-id.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"db_directory_permissions": schema.StringAttribute{
				Description: "Specifies the permissions that should be applied to the directory containing the backend database files and to directories and files created during backup of the backend.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"db_cache_percent": schema.Int64Attribute{
				Description: "Specifies the percentage of JVM memory to allocate to the changelog database cache.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"compact_common_parent_dn": schema.SetAttribute{
				Description: "Provides a DN of an entry that may be the parent for a large number of entries in the backend. This may be used to help increase the space efficiency when encoding entries for storage.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"compress_entries": schema.BoolAttribute{
				Description: "Indicates whether the backend should attempt to compress entries before storing them in the database.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"hash_entries": schema.BoolAttribute{
				Description: "Indicates whether to calculate and store a message digest of the entry contents along with the entry data, in order to provide a means of verifying the integrity of the entry data.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"db_num_cleaner_threads": schema.Int64Attribute{
				Description: "Specifies the number of threads that the backend should maintain to keep the database log files at or near the desired utilization. A value of zero indicates that the number of cleaner threads should be automatically configured based on the number of available CPUs.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"db_cleaner_min_utilization": schema.Int64Attribute{
				Description: "Specifies the minimum percentage of \"live\" data that the database cleaner attempts to keep in database log files.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"db_evictor_critical_percentage": schema.Int64Attribute{
				Description: "Specifies the percentage over the configured maximum that the database cache is allowed to grow. It is recommended to set this value slightly above zero when the database is too large to fully cache in memory. In this case, a dedicated background evictor thread is used to perform evictions once the cache fills up reducing the possibility that server threads are blocked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"db_log_file_max": schema.StringAttribute{
				Description: "Specifies the maximum size for a database log file.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"db_logging_level": schema.StringAttribute{
				Description: "Specifies the log level that should be used by the database when it is writing information into the je.info file.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"je_property": schema.SetAttribute{
				Description: "Specifies the database and environment properties for the Berkeley DB Java Edition database for this changelog backend.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"default_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used for any database for which the cache mode is not explicitly specified. This includes the id2entry database, which stores encoded entries, and all attribute indexes.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id2entry_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the id2entry database, which provides a mapping between entry IDs and entry contents. Consider configuring uncached entries or uncached attributes in lieu of changing from the \"cache-keys-and-values\" default value.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dn2id_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the dn2id database, which provides a mapping between normalized entry DNs and the corresponding entry IDs.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id2children_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the id2children database, which provides a mapping between the entry ID of a particular entry and the entry IDs of all of its immediate children. This index may be used when performing searches with a single-level scope if the search filter cannot be resolved to a small enough candidate list. The size of this database directly depends on the number of entries that have children.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id2subtree_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the id2subtree database, which provides a mapping between the entry ID of a particular entry and the entry IDs of all of its children to any depth. This index may be used when performing searches with a whole-subtree or subordinate-subtree scope if the search filter cannot be resolved to a small enough candidate list. The size of this database directly depends on the number of entries that have children.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dn2uri_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the dn2uri database, which provides a mapping between a normalized entry DN and a set of referral URLs contained in the associated smart referral entry.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"prime_method": schema.SetAttribute{
				Description: "Specifies the method that should be used to prime caches with data for this backend.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"prime_thread_count": schema.Int64Attribute{
				Description: "Specifies the number of threads to use when priming. At present, this applies only to the preload and cursor-across-indexes prime methods.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"system_index_to_prime": schema.SetAttribute{
				Description: "Specifies which system index(es) should be primed when the backend is initialized.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"system_index_to_prime_internal_nodes_only": schema.SetAttribute{
				Description: "Specifies the system index(es) for which internal database nodes only (i.e., the database keys but not values) should be primed when the backend is initialized.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"background_prime": schema.BoolAttribute{
				Description: "Indicates whether to attempt to perform the prime using a background thread if possible. If background priming is enabled, then the Directory Server may be allowed to accept client connections and process requests while the prime is in progress.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"index_entry_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that are allowed to match a given index key before that particular index key is no longer maintained.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"composite_index_entry_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that are allowed to match a given composite index key before that particular composite index key is no longer maintained.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"id2children_index_entry_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entry IDs to maintain for each entry in the id2children system index (which keeps track of the immediate children for an entry, to assist in otherwise unindexed searches with a single-level scope). A value of 0 means there is no limit, however this could have a big impact on database size on disk and on server performance.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"id2subtree_index_entry_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entry IDs to maintain for each entry in the id2subtree system index (which keeps track of all descendants below an entry, to assist in otherwise unindexed searches with a whole-subtree or subordinate subtree scope). A value of 0 means there is no limit, however this could have a big impact on database size on disk and on server performance.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"import_temp_directory": schema.StringAttribute{
				Description: "Specifies the location of the directory that is used to hold temporary information during the index post-processing phase of an LDIF import.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"import_thread_count": schema.Int64Attribute{
				Description: "Specifies the number of threads to use for concurrent processing during an LDIF import.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"export_thread_count": schema.Int64Attribute{
				Description: "Specifies the number of threads to use for concurrently retrieving and encoding entries during an LDIF export.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"db_import_cache_percent": schema.Int64Attribute{
				Description: "The percentage of JVM memory to allocate to the database cache during import operations.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"db_txn_write_no_sync": schema.BoolAttribute{
				Description: "Indicates whether the database should synchronously flush data as it is written to disk.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"deadlock_retry_limit": schema.Int64Attribute{
				Description: "Specifies the number of times that the server should retry an attempted operation in the backend if a deadlock results from two concurrent requests that interfere with each other in a conflicting manner.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"external_txn_default_backend_lock_behavior": schema.StringAttribute{
				Description: "Specifies the default behavior that should be exhibited by external transactions (e.g., an LDAP transaction or an atomic multi-update operation) with regard to acquiring an exclusive lock in this backend.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"single_writer_lock_behavior": schema.StringAttribute{
				Description: "Specifies the condition under which to acquire a single-writer lock to ensure that the associated operation will be the only write in progress at the time the lock is held. The single-writer lock can help avoid problems that result from database lock conflicts that arise between two write operations being processed at the same time in the same backend. This will not have any effect on the read operations processed while the write is in progress.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subtree_delete_size_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that may be deleted from the backend when using the subtree delete control.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"num_recent_changes": schema.Int64Attribute{
				Description: "Specifies the number of recent LDAP entry changes per replica for which the backend keeps a record to allow replication to recover in the event that the server is abruptly terminated. Increasing this value can lead to an increased peak server modification rate as well as increased replication throughput.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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
				Description: "Indicates whether the backend should be considered a private backend, which indicates that it is used for storing operational data rather than user-defined information.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"return_unavailable_when_disabled": schema.BoolAttribute{
				Description: "Determines whether any LDAP operation that would use this Backend is to return UNAVAILABLE when this Backend is disabled.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"notification_manager": schema.StringAttribute{
				Description: "Specifies a notification manager for changes resulting from operations processed through this Backend",
				Optional:    true,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"schema", "backup", "encryption-settings", "ldif", "trust-store", "custom", "changelog", "monitor", "local-db", "config-file-handler", "task", "alert", "alarm", "metrics"}...),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		schemaDef.Attributes["storage_dir"] = schema.StringAttribute{
			Description: "Specifies the path to the directory that will be used to store queued samples.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["metrics_dir"] = schema.StringAttribute{
			Description: "Specifies the path to the directory that contains metric definitions.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["sample_flush_interval"] = schema.StringAttribute{
			Description: "Period when samples are flushed to disk.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["retention_policy"] = schema.SetAttribute{
			Description: "The retention policy to use for the Metrics Backend .",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["alarm_retention_time"] = schema.StringAttribute{
			Description: "Specifies the maximum length of time that information about raised alarms should be maintained before they will be purged.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["max_alarms"] = schema.Int64Attribute{
			Description: "Specifies the maximum number of alarms that should be retained. If more alarms than this configured maximum are generated within the alarm retention time, then the oldest alarms will be purged to achieve this maximum. Only alarms at normal severity will be purged.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["alert_retention_time"] = schema.StringAttribute{
			Description: "Specifies the maximum length of time that information about generated alerts should be maintained before they will be purged.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["max_alerts"] = schema.Int64Attribute{
			Description: "Specifies the maximum number of alerts that should be retained. If more alerts than this configured maximum are generated within the alert retention time, then the oldest alerts will be purged to achieve this maximum.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["disabled_alert_type"] = schema.SetAttribute{
			Description: "Specifies the names of the alert types that should not be added to the backend. This can be used to suppress high volume alerts that might trigger hitting the max-alerts limit sooner than desired. Disabled alert types will not be sent out over persistent searches on this backend.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["task_backing_file"] = schema.StringAttribute{
			Description: "Specifies the path to the backing file for storing information about the tasks configured in the server.",
			Optional:    true,
		}
		schemaDef.Attributes["maximum_initial_task_log_messages_to_retain"] = schema.Int64Attribute{
			Description: "The maximum number of log messages to retain in each task entry from the beginning of the processing for that task. If too many messages are logged during task processing, then retaining only a limited number of messages from the beginning and/or end of task processing can reduce the amount of memory that the server consumes by caching information about currently-active and recently-completed tasks.",
			Optional:    true,
		}
		schemaDef.Attributes["maximum_final_task_log_messages_to_retain"] = schema.Int64Attribute{
			Description: "The maximum number of log messages to retain in each task entry from the end of the processing for that task. If too many messages are logged during task processing, then retaining only a limited number of messages from the beginning and/or end of task processing can reduce the amount of memory that the server consumes by caching information about currently-active and recently-completed tasks.",
			Optional:    true,
		}
		schemaDef.Attributes["task_retention_time"] = schema.StringAttribute{
			Description: "Specifies the length of time that task entries should be retained after processing on the associated task has been completed.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["notification_sender_address"] = schema.StringAttribute{
			Description: "Specifies the email address to use as the sender address (that is, the \"From:\" address) for notification mail messages generated when a task completes execution.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["insignificant_config_archive_attribute"] = schema.SetAttribute{
			Description: "The name or OID of an attribute type that is considered insignificant for the purpose of maintaining the configuration archive.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["mirrored_subtree_peer_polling_interval"] = schema.StringAttribute{
			Description: "Tells the server component that is responsible for mirroring configuration data across a topology of servers the maximum amount of time to wait before polling the peer servers in the topology to determine if there are any changes in the topology. Mirrored data includes meta-data about the servers in the topology as well as cluster-wide configuration data.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["mirrored_subtree_entry_update_timeout"] = schema.StringAttribute{
			Description: "Tells the server component that is responsible for mirroring configuration data across a topology of servers the maximum amount of time to wait for an update operation (add, delete, modify and modify-dn) on an entry to be applied on all servers in the topology. Mirrored data includes meta-data about the servers in the topology as well as cluster-wide configuration data.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["mirrored_subtree_search_timeout"] = schema.StringAttribute{
			Description: "Tells the server component that is responsible for mirroring configuration data across a topology of servers the maximum amount of time to wait for a search operation to complete. Mirrored data includes meta-data about the servers in the topology as well as cluster-wide configuration data. Search requests that take longer than this timeout will be canceled and considered failures.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["changelog_write_batch_size"] = schema.Int64Attribute{
			Description: "Specifies the number of changelog entries written in a single database transaction.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["changelog_purge_batch_size"] = schema.Int64Attribute{
			Description: "Specifies the number of changelog entries purged in a single database transaction.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["changelog_write_queue_capacity"] = schema.Int64Attribute{
			Description: "Specifies the capacity of the changelog write queue in number of changes.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["index_include_attribute"] = schema.SetAttribute{
			Description: "Specifies which attribute types are to be specifically included in the set of attribute indexes maintained on the changelog. If this property does not have any values then no attribute types are indexed.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["index_exclude_attribute"] = schema.SetAttribute{
			Description: "Specifies which attribute types are to be specifically excluded from the set of attribute indexes maintained on the changelog. This property is useful when the index-include-attribute property contains one of the special values \"*\" and \"+\".",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["changelog_maximum_age"] = schema.StringAttribute{
			Description: "Changes are guaranteed to be maintained in the changelog database for at least this duration. Setting target-database-size can allow additional changes to be maintained up to the configured size on disk.",
			Optional:    true,
		}
		schemaDef.Attributes["target_database_size"] = schema.StringAttribute{
			Description: "The changelog database is allowed to grow up to this size on disk even if changes are older than the configured changelog-maximum-age.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["changelog_entry_include_base_dn"] = schema.SetAttribute{
			Description: "The base DNs for branches in the data for which to record changes in the changelog.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["changelog_entry_exclude_base_dn"] = schema.SetAttribute{
			Description: "The base DNs for branches in the data for which no changelog records should be generated.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["changelog_entry_include_filter"] = schema.SetAttribute{
			Description: "A filter that indicates which changelog entries should actually be stored in the changelog. Note that this filter is evaluated against the changelog entry itself and not against the entry that was the target of the change referenced by the changelog entry. This filter may target any attributes that appear in changelog entries with the exception of the changeNumber and entry-size-bytes attributes, since they will not be known at the time of the filter evaluation.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["changelog_entry_exclude_filter"] = schema.SetAttribute{
			Description: "A filter that indicates which changelog entries should be excluded from the changelog. Note that this filter is evaluated against the changelog entry itself and not against the entry that was the target of the change referenced by the changelog entry. This filter may target any attributes that appear in changelog entries with the exception of the changeNumber and entry-size-bytes attributes, since they will not be known at the time of the filter evaluation.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["changelog_include_attribute"] = schema.SetAttribute{
			Description: "Specifies which attribute types will be included in a changelog entry for ADD and MODIFY operations.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["changelog_exclude_attribute"] = schema.SetAttribute{
			Description: "Specifies a set of attribute types that should be excluded in a changelog entry for ADD and MODIFY operations.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["changelog_deleted_entry_include_attribute"] = schema.SetAttribute{
			Description: "Specifies a set of attribute types that should be included in a changelog entry for DELETE operations.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["changelog_deleted_entry_exclude_attribute"] = schema.SetAttribute{
			Description: "Specifies a set of attribute types that should be excluded from a changelog entry for DELETE operations.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["changelog_include_key_attribute"] = schema.SetAttribute{
			Description: "Specifies which attribute types will be included in a changelog entry on every change.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["changelog_max_before_after_values"] = schema.Int64Attribute{
			Description: "This controls whether all attribute values for a modified attribute (even those values that have not changed) will be included in the changelog entry. If the number of attribute values does not exceed this limit, then all values for the modified attribute will be included in the changelog entry.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["write_lastmod_attributes"] = schema.BoolAttribute{
			Description: "Specifies whether values of creatorsName, createTimestamp, modifiersName and modifyTimestamp attributes will be written to changelog entries.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["use_reversible_form"] = schema.BoolAttribute{
			Description: "Specifies whether the changelog should provide enough information to be able to revert the changes if desired.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["include_virtual_attributes"] = schema.SetAttribute{
			Description: "Specifies the changelog entry elements (if any) in which virtual attributes should be included.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["apply_access_controls_to_changelog_entry_contents"] = schema.BoolAttribute{
			Description: "Indicates whether the contents of changelog entries should be subject to access control and sensitive attribute evaluation such that the contents of attributes like changes, deletedEntryAttrs, ds-changelog-entry-key-attr-values, ds-changelog-before-values, and ds-changelog-after-values may be altered based on attributes the user can see in the target entry.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["report_excluded_changelog_attributes"] = schema.StringAttribute{
			Description: "Indicates whether changelog entries that have been altered by applying access controls should include additional information about any attributes that may have been removed.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["soft_delete_entry_included_operation"] = schema.SetAttribute{
			Description: "Specifies which operations performed on soft-deleted entries will appear in the changelog.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["ldif_file"] = schema.StringAttribute{
			Description: "Specifies the path to the LDIF file containing the data for this backend.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["trust_store_file"] = schema.StringAttribute{
			Description: "Specifies the path to the file that stores the trust information.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["trust_store_type"] = schema.StringAttribute{
			Description: "Specifies the format for the data in the key store file.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["trust_store_pin"] = schema.StringAttribute{
			Description: "Specifies the clear-text PIN needed to access the Trust Store Backend.",
			Optional:    true,
			Sensitive:   true,
		}
		schemaDef.Attributes["trust_store_pin_file"] = schema.StringAttribute{
			Description: "Specifies the path to the text file whose only contents should be a single line containing the clear-text PIN needed to access the Trust Store Backend.",
			Optional:    true,
		}
		schemaDef.Attributes["trust_store_pin_passphrase_provider"] = schema.StringAttribute{
			Description: "The passphrase provider to use to obtain the clear-text PIN needed to access the Trust Store Backend.",
			Optional:    true,
		}
		schemaDef.Attributes["backup_directory"] = schema.SetAttribute{
			Description: "Specifies the path to a backup directory containing one or more backups for a particular backend.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["schema_entry_dn"] = schema.SetAttribute{
			Description: "Defines the base DNs of the subtrees in which the schema information is published in addition to the value included in the base-dn property.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["show_all_attributes"] = schema.BoolAttribute{
			Description: "Indicates whether to treat all attributes in the schema entry as if they were user attributes regardless of their configuration.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["read_only_schema_file"] = schema.SetAttribute{
			Description: "Specifies the name of a file (which must exist in the config/schema directory) containing schema elements that should be considered read-only. Any schema definitions contained in read-only files cannot be altered by external clients.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["backup_file_permissions"] = schema.StringAttribute{
			Description: "Specifies the permissions that should be applied to files and directories created by a backup of the backend.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		config.SetAllAttributesToOptionalAndComputed(&schemaDef, []string{"backend_id"})
	}
	config.AddCommonSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *backendResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanBackend(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultBackendResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanBackend(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanBackend(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var model defaultBackendResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.DbDirectoryPermissions) && model.Type.ValueString() != "changelog" && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'db_directory_permissions' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'db_directory_permissions', the 'type' attribute must be one of ['changelog', 'local-db']")
	}
	if internaltypes.IsDefined(model.Id2childrenIndexEntryLimit) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'id2children_index_entry_limit' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'id2children_index_entry_limit', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.CompositeIndexEntryLimit) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'composite_index_entry_limit' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'composite_index_entry_limit', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.ChangelogMaximumAge) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'changelog_maximum_age' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'changelog_maximum_age', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.ChangelogEntryExcludeBaseDN) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'changelog_entry_exclude_base_dn' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'changelog_entry_exclude_base_dn', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.ProcessFiltersWithUndefinedAttributeTypes) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'process_filters_with_undefined_attribute_types' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'process_filters_with_undefined_attribute_types', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.PrimeMethod) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'prime_method' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'prime_method', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.PrimeAllIndexes) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'prime_all_indexes' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'prime_all_indexes', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.MirroredSubtreeSearchTimeout) && model.Type.ValueString() != "config-file-handler" {
		resp.Diagnostics.AddError("Attribute 'mirrored_subtree_search_timeout' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'mirrored_subtree_search_timeout', the 'type' attribute must be one of ['config-file-handler']")
	}
	if internaltypes.IsDefined(model.SingleWriterLockBehavior) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'single_writer_lock_behavior' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'single_writer_lock_behavior', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.PrimeThreadCount) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'prime_thread_count' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'prime_thread_count', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.MaxAlarms) && model.Type.ValueString() != "alarm" {
		resp.Diagnostics.AddError("Attribute 'max_alarms' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'max_alarms', the 'type' attribute must be one of ['alarm']")
	}
	if internaltypes.IsDefined(model.CompactCommonParentDN) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'compact_common_parent_dn' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'compact_common_parent_dn', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.Id2entryCacheMode) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'id2entry_cache_mode' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'id2entry_cache_mode', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.IndexEntryLimit) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'index_entry_limit' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'index_entry_limit', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.WritabilityMode) && model.Type.ValueString() != "schema" && model.Type.ValueString() != "config-file-handler" && model.Type.ValueString() != "backup" && model.Type.ValueString() != "task" && model.Type.ValueString() != "alert" && model.Type.ValueString() != "ldif" && model.Type.ValueString() != "trust-store" && model.Type.ValueString() != "custom" && model.Type.ValueString() != "alarm" && model.Type.ValueString() != "local-db" && model.Type.ValueString() != "metrics" {
		resp.Diagnostics.AddError("Attribute 'writability_mode' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'writability_mode', the 'type' attribute must be one of ['schema', 'config-file-handler', 'backup', 'task', 'alert', 'ldif', 'trust-store', 'custom', 'alarm', 'local-db', 'metrics']")
	}
	if internaltypes.IsDefined(model.IncludeVirtualAttributes) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'include_virtual_attributes' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_virtual_attributes', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.ReturnUnavailableForUntrustedIndex) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'return_unavailable_for_untrusted_index' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'return_unavailable_for_untrusted_index', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.BaseDN) && model.Type.ValueString() != "schema" && model.Type.ValueString() != "backup" && model.Type.ValueString() != "encryption-settings" && model.Type.ValueString() != "ldif" && model.Type.ValueString() != "trust-store" && model.Type.ValueString() != "custom" && model.Type.ValueString() != "changelog" && model.Type.ValueString() != "monitor" && model.Type.ValueString() != "local-db" && model.Type.ValueString() != "config-file-handler" && model.Type.ValueString() != "task" && model.Type.ValueString() != "alert" && model.Type.ValueString() != "alarm" {
		resp.Diagnostics.AddError("Attribute 'base_dn' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'base_dn', the 'type' attribute must be one of ['schema', 'backup', 'encryption-settings', 'ldif', 'trust-store', 'custom', 'changelog', 'monitor', 'local-db', 'config-file-handler', 'task', 'alert', 'alarm']")
	}
	if internaltypes.IsDefined(model.DbCleanerMinUtilization) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'db_cleaner_min_utilization' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'db_cleaner_min_utilization', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.TrustStorePinPassphraseProvider) && model.Type.ValueString() != "trust-store" {
		resp.Diagnostics.AddError("Attribute 'trust_store_pin_passphrase_provider' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'trust_store_pin_passphrase_provider', the 'type' attribute must be one of ['trust-store']")
	}
	if internaltypes.IsDefined(model.ImportThreadCount) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'import_thread_count' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'import_thread_count', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.ExportThreadCount) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'export_thread_count' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'export_thread_count', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.MaximumInitialTaskLogMessagesToRetain) && model.Type.ValueString() != "task" {
		resp.Diagnostics.AddError("Attribute 'maximum_initial_task_log_messages_to_retain' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'maximum_initial_task_log_messages_to_retain', the 'type' attribute must be one of ['task']")
	}
	if internaltypes.IsDefined(model.CompressEntries) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'compress_entries' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'compress_entries', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.AlarmRetentionTime) && model.Type.ValueString() != "alarm" {
		resp.Diagnostics.AddError("Attribute 'alarm_retention_time' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'alarm_retention_time', the 'type' attribute must be one of ['alarm']")
	}
	if internaltypes.IsDefined(model.OfflineProcessDatabaseOpenTimeout) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'offline_process_database_open_timeout' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'offline_process_database_open_timeout', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.TrustStorePinFile) && model.Type.ValueString() != "trust-store" {
		resp.Diagnostics.AddError("Attribute 'trust_store_pin_file' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'trust_store_pin_file', the 'type' attribute must be one of ['trust-store']")
	}
	if internaltypes.IsDefined(model.ReadOnlySchemaFile) && model.Type.ValueString() != "schema" {
		resp.Diagnostics.AddError("Attribute 'read_only_schema_file' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'read_only_schema_file', the 'type' attribute must be one of ['schema']")
	}
	if internaltypes.IsDefined(model.DeadlockRetryLimit) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'deadlock_retry_limit' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'deadlock_retry_limit', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.TrustStorePin) && model.Type.ValueString() != "trust-store" {
		resp.Diagnostics.AddError("Attribute 'trust_store_pin' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'trust_store_pin', the 'type' attribute must be one of ['trust-store']")
	}
	if internaltypes.IsDefined(model.SoftDeleteEntryIncludedOperation) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'soft_delete_entry_included_operation' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'soft_delete_entry_included_operation', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.PrimeTimeLimit) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'prime_time_limit' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'prime_time_limit', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.NotificationSenderAddress) && model.Type.ValueString() != "task" {
		resp.Diagnostics.AddError("Attribute 'notification_sender_address' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'notification_sender_address', the 'type' attribute must be one of ['task']")
	}
	if internaltypes.IsDefined(model.ChangelogExcludeAttribute) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'changelog_exclude_attribute' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'changelog_exclude_attribute', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.ExternalTxnDefaultBackendLockBehavior) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'external_txn_default_backend_lock_behavior' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'external_txn_default_backend_lock_behavior', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.SystemIndexToPrime) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'system_index_to_prime' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'system_index_to_prime', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.LdifFile) && model.Type.ValueString() != "alert" && model.Type.ValueString() != "ldif" && model.Type.ValueString() != "alarm" {
		resp.Diagnostics.AddError("Attribute 'ldif_file' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'ldif_file', the 'type' attribute must be one of ['alert', 'ldif', 'alarm']")
	}
	if internaltypes.IsDefined(model.MaximumFinalTaskLogMessagesToRetain) && model.Type.ValueString() != "task" {
		resp.Diagnostics.AddError("Attribute 'maximum_final_task_log_messages_to_retain' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'maximum_final_task_log_messages_to_retain', the 'type' attribute must be one of ['task']")
	}
	if internaltypes.IsDefined(model.Dn2uriCacheMode) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'dn2uri_cache_mode' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'dn2uri_cache_mode', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.DbImportCachePercent) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'db_import_cache_percent' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'db_import_cache_percent', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.IsPrivateBackend) && model.Type.ValueString() != "ldif" && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'is_private_backend' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'is_private_backend', the 'type' attribute must be one of ['ldif', 'local-db']")
	}
	if internaltypes.IsDefined(model.TaskBackingFile) && model.Type.ValueString() != "task" {
		resp.Diagnostics.AddError("Attribute 'task_backing_file' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'task_backing_file', the 'type' attribute must be one of ['task']")
	}
	if internaltypes.IsDefined(model.WriteLastmodAttributes) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'write_lastmod_attributes' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'write_lastmod_attributes', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.SchemaEntryDN) && model.Type.ValueString() != "schema" {
		resp.Diagnostics.AddError("Attribute 'schema_entry_dn' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'schema_entry_dn', the 'type' attribute must be one of ['schema']")
	}
	if internaltypes.IsDefined(model.BackupDirectory) && model.Type.ValueString() != "backup" {
		resp.Diagnostics.AddError("Attribute 'backup_directory' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'backup_directory', the 'type' attribute must be one of ['backup']")
	}
	if internaltypes.IsDefined(model.UseReversibleForm) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'use_reversible_form' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'use_reversible_form', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.DbNumCleanerThreads) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'db_num_cleaner_threads' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'db_num_cleaner_threads', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.SampleFlushInterval) && model.Type.ValueString() != "metrics" {
		resp.Diagnostics.AddError("Attribute 'sample_flush_interval' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'sample_flush_interval', the 'type' attribute must be one of ['metrics']")
	}
	if internaltypes.IsDefined(model.BackupFilePermissions) && model.Type.ValueString() != "schema" && model.Type.ValueString() != "config-file-handler" && model.Type.ValueString() != "task" && model.Type.ValueString() != "encryption-settings" && model.Type.ValueString() != "alert" && model.Type.ValueString() != "ldif" && model.Type.ValueString() != "trust-store" && model.Type.ValueString() != "custom" && model.Type.ValueString() != "alarm" {
		resp.Diagnostics.AddError("Attribute 'backup_file_permissions' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'backup_file_permissions', the 'type' attribute must be one of ['schema', 'config-file-handler', 'task', 'encryption-settings', 'alert', 'ldif', 'trust-store', 'custom', 'alarm']")
	}
	if internaltypes.IsDefined(model.ChangelogIncludeKeyAttribute) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'changelog_include_key_attribute' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'changelog_include_key_attribute', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.Id2subtreeCacheMode) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'id2subtree_cache_mode' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'id2subtree_cache_mode', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.InsignificantConfigArchiveAttribute) && model.Type.ValueString() != "config-file-handler" {
		resp.Diagnostics.AddError("Attribute 'insignificant_config_archive_attribute' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'insignificant_config_archive_attribute', the 'type' attribute must be one of ['config-file-handler']")
	}
	if internaltypes.IsDefined(model.ChangelogDeletedEntryIncludeAttribute) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'changelog_deleted_entry_include_attribute' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'changelog_deleted_entry_include_attribute', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.DisabledAlertType) && model.Type.ValueString() != "alert" {
		resp.Diagnostics.AddError("Attribute 'disabled_alert_type' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'disabled_alert_type', the 'type' attribute must be one of ['alert']")
	}
	if internaltypes.IsDefined(model.DbLogFileMax) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'db_log_file_max' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'db_log_file_max', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.MirroredSubtreeEntryUpdateTimeout) && model.Type.ValueString() != "config-file-handler" {
		resp.Diagnostics.AddError("Attribute 'mirrored_subtree_entry_update_timeout' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'mirrored_subtree_entry_update_timeout', the 'type' attribute must be one of ['config-file-handler']")
	}
	if internaltypes.IsDefined(model.HashEntries) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'hash_entries' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'hash_entries', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.DefaultCacheMode) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'default_cache_mode' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'default_cache_mode', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.ChangelogIncludeAttribute) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'changelog_include_attribute' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'changelog_include_attribute', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.UncachedAttributeCriteria) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'uncached_attribute_criteria' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'uncached_attribute_criteria', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.DbCheckpointerWakeupInterval) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'db_checkpointer_wakeup_interval' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'db_checkpointer_wakeup_interval', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.DbUseThreadLocalHandles) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'db_use_thread_local_handles' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'db_use_thread_local_handles', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.ChangelogEntryIncludeFilter) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'changelog_entry_include_filter' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'changelog_entry_include_filter', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.TrustStoreFile) && model.Type.ValueString() != "trust-store" {
		resp.Diagnostics.AddError("Attribute 'trust_store_file' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'trust_store_file', the 'type' attribute must be one of ['trust-store']")
	}
	if internaltypes.IsDefined(model.IndexIncludeAttribute) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'index_include_attribute' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'index_include_attribute', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.ReportExcludedChangelogAttributes) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'report_excluded_changelog_attributes' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'report_excluded_changelog_attributes', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.ChangelogPurgeBatchSize) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'changelog_purge_batch_size' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'changelog_purge_batch_size', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.StorageDir) && model.Type.ValueString() != "metrics" {
		resp.Diagnostics.AddError("Attribute 'storage_dir' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'storage_dir', the 'type' attribute must be one of ['metrics']")
	}
	if internaltypes.IsDefined(model.ShowAllAttributes) && model.Type.ValueString() != "schema" {
		resp.Diagnostics.AddError("Attribute 'show_all_attributes' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'show_all_attributes', the 'type' attribute must be one of ['schema']")
	}
	if internaltypes.IsDefined(model.DbDirectory) && model.Type.ValueString() != "changelog" && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'db_directory' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'db_directory', the 'type' attribute must be one of ['changelog', 'local-db']")
	}
	if internaltypes.IsDefined(model.JeProperty) && model.Type.ValueString() != "changelog" && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'je_property' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'je_property', the 'type' attribute must be one of ['changelog', 'local-db']")
	}
	if internaltypes.IsDefined(model.UncachedId2entryCacheMode) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'uncached_id2entry_cache_mode' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'uncached_id2entry_cache_mode', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.ChangelogWriteBatchSize) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'changelog_write_batch_size' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'changelog_write_batch_size', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.ChangelogEntryExcludeFilter) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'changelog_entry_exclude_filter' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'changelog_entry_exclude_filter', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.DbLoggingLevel) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'db_logging_level' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'db_logging_level', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.SubtreeDeleteSizeLimit) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'subtree_delete_size_limit' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'subtree_delete_size_limit', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.MirroredSubtreePeerPollingInterval) && model.Type.ValueString() != "config-file-handler" {
		resp.Diagnostics.AddError("Attribute 'mirrored_subtree_peer_polling_interval' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'mirrored_subtree_peer_polling_interval', the 'type' attribute must be one of ['config-file-handler']")
	}
	if internaltypes.IsDefined(model.ChangelogDeletedEntryExcludeAttribute) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'changelog_deleted_entry_exclude_attribute' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'changelog_deleted_entry_exclude_attribute', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.SetDegradedAlertWhenDisabled) && model.Type.ValueString() != "schema" && model.Type.ValueString() != "backup" && model.Type.ValueString() != "encryption-settings" && model.Type.ValueString() != "ldif" && model.Type.ValueString() != "trust-store" && model.Type.ValueString() != "custom" && model.Type.ValueString() != "changelog" && model.Type.ValueString() != "monitor" && model.Type.ValueString() != "local-db" && model.Type.ValueString() != "config-file-handler" && model.Type.ValueString() != "task" && model.Type.ValueString() != "alert" && model.Type.ValueString() != "alarm" {
		resp.Diagnostics.AddError("Attribute 'set_degraded_alert_when_disabled' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'set_degraded_alert_when_disabled', the 'type' attribute must be one of ['schema', 'backup', 'encryption-settings', 'ldif', 'trust-store', 'custom', 'changelog', 'monitor', 'local-db', 'config-file-handler', 'task', 'alert', 'alarm']")
	}
	if internaltypes.IsDefined(model.DbCachePercent) && model.Type.ValueString() != "changelog" && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'db_cache_percent' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'db_cache_percent', the 'type' attribute must be one of ['changelog', 'local-db']")
	}
	if internaltypes.IsDefined(model.MaxAlerts) && model.Type.ValueString() != "alert" {
		resp.Diagnostics.AddError("Attribute 'max_alerts' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'max_alerts', the 'type' attribute must be one of ['alert']")
	}
	if internaltypes.IsDefined(model.DbEvictorCriticalPercentage) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'db_evictor_critical_percentage' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'db_evictor_critical_percentage', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.TargetDatabaseSize) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'target_database_size' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'target_database_size', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.ChangelogWriteQueueCapacity) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'changelog_write_queue_capacity' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'changelog_write_queue_capacity', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.TaskRetentionTime) && model.Type.ValueString() != "task" {
		resp.Diagnostics.AddError("Attribute 'task_retention_time' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'task_retention_time', the 'type' attribute must be one of ['task']")
	}
	if internaltypes.IsDefined(model.IndexExcludeAttribute) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'index_exclude_attribute' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'index_exclude_attribute', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.AlertRetentionTime) && model.Type.ValueString() != "alert" {
		resp.Diagnostics.AddError("Attribute 'alert_retention_time' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'alert_retention_time', the 'type' attribute must be one of ['alert']")
	}
	if internaltypes.IsDefined(model.SystemIndexToPrimeInternalNodesOnly) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'system_index_to_prime_internal_nodes_only' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'system_index_to_prime_internal_nodes_only', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.TrustStoreType) && model.Type.ValueString() != "trust-store" {
		resp.Diagnostics.AddError("Attribute 'trust_store_type' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'trust_store_type', the 'type' attribute must be one of ['trust-store']")
	}
	if internaltypes.IsDefined(model.RetentionPolicy) && model.Type.ValueString() != "metrics" {
		resp.Diagnostics.AddError("Attribute 'retention_policy' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'retention_policy', the 'type' attribute must be one of ['metrics']")
	}
	if internaltypes.IsDefined(model.DbTxnWriteNoSync) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'db_txn_write_no_sync' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'db_txn_write_no_sync', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.UncachedEntryCriteria) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'uncached_entry_criteria' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'uncached_entry_criteria', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.MetricsDir) && model.Type.ValueString() != "metrics" {
		resp.Diagnostics.AddError("Attribute 'metrics_dir' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'metrics_dir', the 'type' attribute must be one of ['metrics']")
	}
	if internaltypes.IsDefined(model.NumRecentChanges) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'num_recent_changes' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'num_recent_changes', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.BackgroundPrime) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'background_prime' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'background_prime', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.DbBackgroundSyncInterval) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'db_background_sync_interval' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'db_background_sync_interval', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.Id2subtreeIndexEntryLimit) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'id2subtree_index_entry_limit' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'id2subtree_index_entry_limit', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.ChangelogEntryIncludeBaseDN) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'changelog_entry_include_base_dn' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'changelog_entry_include_base_dn', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.ApplyAccessControlsToChangelogEntryContents) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'apply_access_controls_to_changelog_entry_contents' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'apply_access_controls_to_changelog_entry_contents', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.ChangelogMaxBeforeAfterValues) && model.Type.ValueString() != "changelog" {
		resp.Diagnostics.AddError("Attribute 'changelog_max_before_after_values' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'changelog_max_before_after_values', the 'type' attribute must be one of ['changelog']")
	}
	if internaltypes.IsDefined(model.ImportTempDirectory) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'import_temp_directory' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'import_temp_directory', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.Dn2idCacheMode) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'dn2id_cache_mode' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'dn2id_cache_mode', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.Id2childrenCacheMode) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'id2children_cache_mode' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'id2children_cache_mode', the 'type' attribute must be one of ['local-db']")
	}
	if internaltypes.IsDefined(model.SetDegradedAlertForUntrustedIndex) && model.Type.ValueString() != "local-db" {
		resp.Diagnostics.AddError("Attribute 'set_degraded_alert_for_untrusted_index' not supported by pingdirectory_backend resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'set_degraded_alert_for_untrusted_index', the 'type' attribute must be one of ['local-db']")
	}
}

// Add optional fields to create request for local-db backend
func addOptionalLocalDbBackendFields(ctx context.Context, addRequest *client.AddLocalDbBackendRequest, plan backendResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UncachedId2entryCacheMode) {
		uncachedId2EntryCacheMode, err := client.NewEnumbackendUncachedId2entryCacheModePropFromValue(plan.UncachedId2entryCacheMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.UncachedId2entryCacheMode = uncachedId2EntryCacheMode
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
		id2EntryCacheMode, err := client.NewEnumbackendId2entryCacheModePropFromValue(plan.Id2entryCacheMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.Id2entryCacheMode = id2EntryCacheMode
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Dn2idCacheMode) {
		dn2IdCacheMode, err := client.NewEnumbackendDn2idCacheModePropFromValue(plan.Dn2idCacheMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.Dn2idCacheMode = dn2IdCacheMode
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Id2childrenCacheMode) {
		id2ChildrenCacheMode, err := client.NewEnumbackendId2childrenCacheModePropFromValue(plan.Id2childrenCacheMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.Id2childrenCacheMode = id2ChildrenCacheMode
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Id2subtreeCacheMode) {
		id2SubtreeCacheMode, err := client.NewEnumbackendId2subtreeCacheModePropFromValue(plan.Id2subtreeCacheMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.Id2subtreeCacheMode = id2SubtreeCacheMode
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Dn2uriCacheMode) {
		dn2UriCacheMode, err := client.NewEnumbackendDn2uriCacheModePropFromValue(plan.Dn2uriCacheMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.Dn2uriCacheMode = dn2UriCacheMode
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

// Populate any sets that have a nil ElementType, to avoid a nil pointer when setting the state
func populateBackendNilSets(ctx context.Context, model *backendResourceModel) {
	if model.JeProperty.ElementType(ctx) == nil {
		model.JeProperty = types.SetNull(types.StringType)
	}
	if model.BaseDN.ElementType(ctx) == nil {
		model.BaseDN = types.SetNull(types.StringType)
	}
	if model.CompactCommonParentDN.ElementType(ctx) == nil {
		model.CompactCommonParentDN = types.SetNull(types.StringType)
	}
	if model.SystemIndexToPrimeInternalNodesOnly.ElementType(ctx) == nil {
		model.SystemIndexToPrimeInternalNodesOnly = types.SetNull(types.StringType)
	}
	if model.PrimeMethod.ElementType(ctx) == nil {
		model.PrimeMethod = types.SetNull(types.StringType)
	}
	if model.SystemIndexToPrime.ElementType(ctx) == nil {
		model.SystemIndexToPrime = types.SetNull(types.StringType)
	}
}

// Populate any sets that have a nil ElementType, to avoid a nil pointer when setting the state
func populateBackendNilSetsDefault(ctx context.Context, model *defaultBackendResourceModel) {
	if model.ChangelogEntryIncludeFilter.ElementType(ctx) == nil {
		model.ChangelogEntryIncludeFilter = types.SetNull(types.StringType)
	}
	if model.IncludeVirtualAttributes.ElementType(ctx) == nil {
		model.IncludeVirtualAttributes = types.SetNull(types.StringType)
	}
	if model.SystemIndexToPrime.ElementType(ctx) == nil {
		model.SystemIndexToPrime = types.SetNull(types.StringType)
	}
	if model.ReadOnlySchemaFile.ElementType(ctx) == nil {
		model.ReadOnlySchemaFile = types.SetNull(types.StringType)
	}
	if model.ChangelogEntryExcludeFilter.ElementType(ctx) == nil {
		model.ChangelogEntryExcludeFilter = types.SetNull(types.StringType)
	}
	if model.BackupDirectory.ElementType(ctx) == nil {
		model.BackupDirectory = types.SetNull(types.StringType)
	}
	if model.BaseDN.ElementType(ctx) == nil {
		model.BaseDN = types.SetNull(types.StringType)
	}
	if model.SchemaEntryDN.ElementType(ctx) == nil {
		model.SchemaEntryDN = types.SetNull(types.StringType)
	}
	if model.IndexExcludeAttribute.ElementType(ctx) == nil {
		model.IndexExcludeAttribute = types.SetNull(types.StringType)
	}
	if model.ChangelogExcludeAttribute.ElementType(ctx) == nil {
		model.ChangelogExcludeAttribute = types.SetNull(types.StringType)
	}
	if model.ChangelogIncludeAttribute.ElementType(ctx) == nil {
		model.ChangelogIncludeAttribute = types.SetNull(types.StringType)
	}
	if model.SoftDeleteEntryIncludedOperation.ElementType(ctx) == nil {
		model.SoftDeleteEntryIncludedOperation = types.SetNull(types.StringType)
	}
	if model.DisabledAlertType.ElementType(ctx) == nil {
		model.DisabledAlertType = types.SetNull(types.StringType)
	}
	if model.ChangelogEntryIncludeBaseDN.ElementType(ctx) == nil {
		model.ChangelogEntryIncludeBaseDN = types.SetNull(types.StringType)
	}
	if model.ChangelogIncludeKeyAttribute.ElementType(ctx) == nil {
		model.ChangelogIncludeKeyAttribute = types.SetNull(types.StringType)
	}
	if model.CompactCommonParentDN.ElementType(ctx) == nil {
		model.CompactCommonParentDN = types.SetNull(types.StringType)
	}
	if model.IndexIncludeAttribute.ElementType(ctx) == nil {
		model.IndexIncludeAttribute = types.SetNull(types.StringType)
	}
	if model.SystemIndexToPrimeInternalNodesOnly.ElementType(ctx) == nil {
		model.SystemIndexToPrimeInternalNodesOnly = types.SetNull(types.StringType)
	}
	if model.RetentionPolicy.ElementType(ctx) == nil {
		model.RetentionPolicy = types.SetNull(types.StringType)
	}
	if model.ChangelogDeletedEntryIncludeAttribute.ElementType(ctx) == nil {
		model.ChangelogDeletedEntryIncludeAttribute = types.SetNull(types.StringType)
	}
	if model.JeProperty.ElementType(ctx) == nil {
		model.JeProperty = types.SetNull(types.StringType)
	}
	if model.ChangelogEntryExcludeBaseDN.ElementType(ctx) == nil {
		model.ChangelogEntryExcludeBaseDN = types.SetNull(types.StringType)
	}
	if model.ChangelogDeletedEntryExcludeAttribute.ElementType(ctx) == nil {
		model.ChangelogDeletedEntryExcludeAttribute = types.SetNull(types.StringType)
	}
	if model.InsignificantConfigArchiveAttribute.ElementType(ctx) == nil {
		model.InsignificantConfigArchiveAttribute = types.SetNull(types.StringType)
	}
	if model.PrimeMethod.ElementType(ctx) == nil {
		model.PrimeMethod = types.SetNull(types.StringType)
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
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, internaltypes.IsEmptyString(expectedValues.BackupFilePermissions))
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, internaltypes.IsEmptyString(expectedValues.NotificationManager))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendNilSetsDefault(ctx, state)
}

// Read a BackupBackendResponse object into the model struct
func readBackupBackendResponseDefault(ctx context.Context, r *client.BackupBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("backup")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.BackupDirectory = internaltypes.GetStringSet(r.BackupDirectory)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, internaltypes.IsEmptyString(expectedValues.NotificationManager))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendNilSetsDefault(ctx, state)
}

// Read a EncryptionSettingsBackendResponse object into the model struct
func readEncryptionSettingsBackendResponseDefault(ctx context.Context, r *client.EncryptionSettingsBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("encryption-settings")
	state.Id = types.StringValue(r.Id)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.BackendID = types.StringValue(r.BackendID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, internaltypes.IsEmptyString(expectedValues.BackupFilePermissions))
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, internaltypes.IsEmptyString(expectedValues.NotificationManager))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendNilSetsDefault(ctx, state)
}

// Read a LdifBackendResponse object into the model struct
func readLdifBackendResponseDefault(ctx context.Context, r *client.LdifBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldif")
	state.Id = types.StringValue(r.Id)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.IsPrivateBackend = internaltypes.BoolTypeOrNil(r.IsPrivateBackend)
	state.LdifFile = types.StringValue(r.LdifFile)
	state.BackendID = types.StringValue(r.BackendID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, internaltypes.IsEmptyString(expectedValues.BackupFilePermissions))
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, internaltypes.IsEmptyString(expectedValues.NotificationManager))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendNilSetsDefault(ctx, state)
}

// Read a TrustStoreBackendResponse object into the model struct
func readTrustStoreBackendResponseDefault(ctx context.Context, r *client.TrustStoreBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("trust-store")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.TrustStoreFile = types.StringValue(r.TrustStoreFile)
	state.TrustStoreType = internaltypes.StringTypeOrNil(r.TrustStoreType, internaltypes.IsEmptyString(expectedValues.TrustStoreType))
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.TrustStorePin = expectedValues.TrustStorePin
	state.TrustStorePinFile = internaltypes.StringTypeOrNil(r.TrustStorePinFile, internaltypes.IsEmptyString(expectedValues.TrustStorePinFile))
	state.TrustStorePinPassphraseProvider = internaltypes.StringTypeOrNil(r.TrustStorePinPassphraseProvider, internaltypes.IsEmptyString(expectedValues.TrustStorePinPassphraseProvider))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, internaltypes.IsEmptyString(expectedValues.BackupFilePermissions))
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, internaltypes.IsEmptyString(expectedValues.NotificationManager))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendNilSetsDefault(ctx, state)
}

// Read a CustomBackendResponse object into the model struct
func readCustomBackendResponseDefault(ctx context.Context, r *client.CustomBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("custom")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, internaltypes.IsEmptyString(expectedValues.BackupFilePermissions))
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, internaltypes.IsEmptyString(expectedValues.NotificationManager))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendNilSetsDefault(ctx, state)
}

// Read a ChangelogBackendResponse object into the model struct
func readChangelogBackendResponseDefault(ctx context.Context, r *client.ChangelogBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("changelog")
	state.Id = types.StringValue(r.Id)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.DbDirectory = internaltypes.StringTypeOrNil(r.DbDirectory, internaltypes.IsEmptyString(expectedValues.DbDirectory))
	state.DbDirectoryPermissions = internaltypes.StringTypeOrNil(r.DbDirectoryPermissions, internaltypes.IsEmptyString(expectedValues.DbDirectoryPermissions))
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
	state.TargetDatabaseSize = internaltypes.StringTypeOrNil(r.TargetDatabaseSize, internaltypes.IsEmptyString(expectedValues.TargetDatabaseSize))
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
		client.StringPointerEnumbackendReportExcludedChangelogAttributesProp(r.ReportExcludedChangelogAttributes), internaltypes.IsEmptyString(expectedValues.ReportExcludedChangelogAttributes))
	state.SoftDeleteEntryIncludedOperation = internaltypes.GetStringSet(
		client.StringSliceEnumbackendSoftDeleteEntryIncludedOperationProp(r.SoftDeleteEntryIncludedOperation))
	state.BackendID = types.StringValue(r.BackendID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, internaltypes.IsEmptyString(expectedValues.NotificationManager))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendNilSetsDefault(ctx, state)
}

// Read a MonitorBackendResponse object into the model struct
func readMonitorBackendResponseDefault(ctx context.Context, r *client.MonitorBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("monitor")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, internaltypes.IsEmptyString(expectedValues.NotificationManager))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendNilSetsDefault(ctx, state)
}

// Read a LocalDbBackendResponse object into the model struct
func readLocalDbBackendResponse(ctx context.Context, r *client.LocalDbBackendResponse, state *backendResourceModel, expectedValues *backendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("local-db")
	state.Id = types.StringValue(r.Id)
	state.UncachedId2entryCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendUncachedId2entryCacheModeProp(r.UncachedId2entryCacheMode), internaltypes.IsEmptyString(expectedValues.UncachedId2entryCacheMode))
	state.UncachedAttributeCriteria = internaltypes.StringTypeOrNil(r.UncachedAttributeCriteria, internaltypes.IsEmptyString(expectedValues.UncachedAttributeCriteria))
	state.UncachedEntryCriteria = internaltypes.StringTypeOrNil(r.UncachedEntryCriteria, internaltypes.IsEmptyString(expectedValues.UncachedEntryCriteria))
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.SetDegradedAlertForUntrustedIndex = internaltypes.BoolTypeOrNil(r.SetDegradedAlertForUntrustedIndex)
	state.ReturnUnavailableForUntrustedIndex = internaltypes.BoolTypeOrNil(r.ReturnUnavailableForUntrustedIndex)
	state.ProcessFiltersWithUndefinedAttributeTypes = internaltypes.BoolTypeOrNil(r.ProcessFiltersWithUndefinedAttributeTypes)
	state.IsPrivateBackend = internaltypes.BoolTypeOrNil(r.IsPrivateBackend)
	state.DbDirectory = types.StringValue(r.DbDirectory)
	state.DbDirectoryPermissions = internaltypes.StringTypeOrNil(r.DbDirectoryPermissions, internaltypes.IsEmptyString(expectedValues.DbDirectoryPermissions))
	state.CompactCommonParentDN = internaltypes.GetStringSet(r.CompactCommonParentDN)
	state.CompressEntries = internaltypes.BoolTypeOrNil(r.CompressEntries)
	state.HashEntries = internaltypes.BoolTypeOrNil(r.HashEntries)
	state.DbNumCleanerThreads = internaltypes.Int64TypeOrNil(r.DbNumCleanerThreads)
	state.DbCleanerMinUtilization = internaltypes.Int64TypeOrNil(r.DbCleanerMinUtilization)
	state.DbEvictorCriticalPercentage = internaltypes.Int64TypeOrNil(r.DbEvictorCriticalPercentage)
	state.DbCheckpointerWakeupInterval = internaltypes.StringTypeOrNil(r.DbCheckpointerWakeupInterval, internaltypes.IsEmptyString(expectedValues.DbCheckpointerWakeupInterval))
	config.CheckMismatchedPDFormattedAttributes("db_checkpointer_wakeup_interval",
		expectedValues.DbCheckpointerWakeupInterval, state.DbCheckpointerWakeupInterval, diagnostics)
	state.DbBackgroundSyncInterval = internaltypes.StringTypeOrNil(r.DbBackgroundSyncInterval, internaltypes.IsEmptyString(expectedValues.DbBackgroundSyncInterval))
	config.CheckMismatchedPDFormattedAttributes("db_background_sync_interval",
		expectedValues.DbBackgroundSyncInterval, state.DbBackgroundSyncInterval, diagnostics)
	state.DbUseThreadLocalHandles = internaltypes.BoolTypeOrNil(r.DbUseThreadLocalHandles)
	state.DbLogFileMax = internaltypes.StringTypeOrNil(r.DbLogFileMax, internaltypes.IsEmptyString(expectedValues.DbLogFileMax))
	config.CheckMismatchedPDFormattedAttributes("db_log_file_max",
		expectedValues.DbLogFileMax, state.DbLogFileMax, diagnostics)
	state.DbLoggingLevel = internaltypes.StringTypeOrNil(r.DbLoggingLevel, internaltypes.IsEmptyString(expectedValues.DbLoggingLevel))
	state.JeProperty = internaltypes.GetStringSet(r.JeProperty)
	state.DbCachePercent = internaltypes.Int64TypeOrNil(r.DbCachePercent)
	state.DefaultCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendDefaultCacheModeProp(r.DefaultCacheMode), internaltypes.IsEmptyString(expectedValues.DefaultCacheMode))
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
	state.PrimeTimeLimit = internaltypes.StringTypeOrNil(r.PrimeTimeLimit, internaltypes.IsEmptyString(expectedValues.PrimeTimeLimit))
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
		client.StringPointerEnumbackendExternalTxnDefaultBackendLockBehaviorProp(r.ExternalTxnDefaultBackendLockBehavior), internaltypes.IsEmptyString(expectedValues.ExternalTxnDefaultBackendLockBehavior))
	state.SingleWriterLockBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendSingleWriterLockBehaviorProp(r.SingleWriterLockBehavior), internaltypes.IsEmptyString(expectedValues.SingleWriterLockBehavior))
	state.SubtreeDeleteSizeLimit = internaltypes.Int64TypeOrNil(r.SubtreeDeleteSizeLimit)
	state.NumRecentChanges = internaltypes.Int64TypeOrNil(r.NumRecentChanges)
	state.OfflineProcessDatabaseOpenTimeout = internaltypes.StringTypeOrNil(r.OfflineProcessDatabaseOpenTimeout, internaltypes.IsEmptyString(expectedValues.OfflineProcessDatabaseOpenTimeout))
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
	populateBackendNilSets(ctx, state)
}

// Read a LocalDbBackendResponse object into the model struct
func readLocalDbBackendResponseDefault(ctx context.Context, r *client.LocalDbBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("local-db")
	state.Id = types.StringValue(r.Id)
	state.UncachedId2entryCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendUncachedId2entryCacheModeProp(r.UncachedId2entryCacheMode), internaltypes.IsEmptyString(expectedValues.UncachedId2entryCacheMode))
	state.UncachedAttributeCriteria = internaltypes.StringTypeOrNil(r.UncachedAttributeCriteria, internaltypes.IsEmptyString(expectedValues.UncachedAttributeCriteria))
	state.UncachedEntryCriteria = internaltypes.StringTypeOrNil(r.UncachedEntryCriteria, internaltypes.IsEmptyString(expectedValues.UncachedEntryCriteria))
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.SetDegradedAlertForUntrustedIndex = internaltypes.BoolTypeOrNil(r.SetDegradedAlertForUntrustedIndex)
	state.ReturnUnavailableForUntrustedIndex = internaltypes.BoolTypeOrNil(r.ReturnUnavailableForUntrustedIndex)
	state.ProcessFiltersWithUndefinedAttributeTypes = internaltypes.BoolTypeOrNil(r.ProcessFiltersWithUndefinedAttributeTypes)
	state.IsPrivateBackend = internaltypes.BoolTypeOrNil(r.IsPrivateBackend)
	state.DbDirectory = types.StringValue(r.DbDirectory)
	state.DbDirectoryPermissions = internaltypes.StringTypeOrNil(r.DbDirectoryPermissions, internaltypes.IsEmptyString(expectedValues.DbDirectoryPermissions))
	state.CompactCommonParentDN = internaltypes.GetStringSet(r.CompactCommonParentDN)
	state.CompressEntries = internaltypes.BoolTypeOrNil(r.CompressEntries)
	state.HashEntries = internaltypes.BoolTypeOrNil(r.HashEntries)
	state.DbNumCleanerThreads = internaltypes.Int64TypeOrNil(r.DbNumCleanerThreads)
	state.DbCleanerMinUtilization = internaltypes.Int64TypeOrNil(r.DbCleanerMinUtilization)
	state.DbEvictorCriticalPercentage = internaltypes.Int64TypeOrNil(r.DbEvictorCriticalPercentage)
	state.DbCheckpointerWakeupInterval = internaltypes.StringTypeOrNil(r.DbCheckpointerWakeupInterval, internaltypes.IsEmptyString(expectedValues.DbCheckpointerWakeupInterval))
	config.CheckMismatchedPDFormattedAttributes("db_checkpointer_wakeup_interval",
		expectedValues.DbCheckpointerWakeupInterval, state.DbCheckpointerWakeupInterval, diagnostics)
	state.DbBackgroundSyncInterval = internaltypes.StringTypeOrNil(r.DbBackgroundSyncInterval, internaltypes.IsEmptyString(expectedValues.DbBackgroundSyncInterval))
	config.CheckMismatchedPDFormattedAttributes("db_background_sync_interval",
		expectedValues.DbBackgroundSyncInterval, state.DbBackgroundSyncInterval, diagnostics)
	state.DbUseThreadLocalHandles = internaltypes.BoolTypeOrNil(r.DbUseThreadLocalHandles)
	state.DbLogFileMax = internaltypes.StringTypeOrNil(r.DbLogFileMax, internaltypes.IsEmptyString(expectedValues.DbLogFileMax))
	config.CheckMismatchedPDFormattedAttributes("db_log_file_max",
		expectedValues.DbLogFileMax, state.DbLogFileMax, diagnostics)
	state.DbLoggingLevel = internaltypes.StringTypeOrNil(r.DbLoggingLevel, internaltypes.IsEmptyString(expectedValues.DbLoggingLevel))
	state.JeProperty = internaltypes.GetStringSet(r.JeProperty)
	state.DbCachePercent = internaltypes.Int64TypeOrNil(r.DbCachePercent)
	state.DefaultCacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendDefaultCacheModeProp(r.DefaultCacheMode), internaltypes.IsEmptyString(expectedValues.DefaultCacheMode))
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
	state.PrimeTimeLimit = internaltypes.StringTypeOrNil(r.PrimeTimeLimit, internaltypes.IsEmptyString(expectedValues.PrimeTimeLimit))
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
		client.StringPointerEnumbackendExternalTxnDefaultBackendLockBehaviorProp(r.ExternalTxnDefaultBackendLockBehavior), internaltypes.IsEmptyString(expectedValues.ExternalTxnDefaultBackendLockBehavior))
	state.SingleWriterLockBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendSingleWriterLockBehaviorProp(r.SingleWriterLockBehavior), internaltypes.IsEmptyString(expectedValues.SingleWriterLockBehavior))
	state.SubtreeDeleteSizeLimit = internaltypes.Int64TypeOrNil(r.SubtreeDeleteSizeLimit)
	state.NumRecentChanges = internaltypes.Int64TypeOrNil(r.NumRecentChanges)
	state.OfflineProcessDatabaseOpenTimeout = internaltypes.StringTypeOrNil(r.OfflineProcessDatabaseOpenTimeout, internaltypes.IsEmptyString(expectedValues.OfflineProcessDatabaseOpenTimeout))
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
	populateBackendNilSetsDefault(ctx, state)
}

// Read a ConfigFileHandlerBackendResponse object into the model struct
func readConfigFileHandlerBackendResponseDefault(ctx context.Context, r *client.ConfigFileHandlerBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("config-file-handler")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.InsignificantConfigArchiveAttribute = internaltypes.GetStringSet(r.InsignificantConfigArchiveAttribute)
	state.MirroredSubtreePeerPollingInterval = internaltypes.StringTypeOrNil(r.MirroredSubtreePeerPollingInterval, internaltypes.IsEmptyString(expectedValues.MirroredSubtreePeerPollingInterval))
	config.CheckMismatchedPDFormattedAttributes("mirrored_subtree_peer_polling_interval",
		expectedValues.MirroredSubtreePeerPollingInterval, state.MirroredSubtreePeerPollingInterval, diagnostics)
	state.MirroredSubtreeEntryUpdateTimeout = internaltypes.StringTypeOrNil(r.MirroredSubtreeEntryUpdateTimeout, internaltypes.IsEmptyString(expectedValues.MirroredSubtreeEntryUpdateTimeout))
	config.CheckMismatchedPDFormattedAttributes("mirrored_subtree_entry_update_timeout",
		expectedValues.MirroredSubtreeEntryUpdateTimeout, state.MirroredSubtreeEntryUpdateTimeout, diagnostics)
	state.MirroredSubtreeSearchTimeout = internaltypes.StringTypeOrNil(r.MirroredSubtreeSearchTimeout, internaltypes.IsEmptyString(expectedValues.MirroredSubtreeSearchTimeout))
	config.CheckMismatchedPDFormattedAttributes("mirrored_subtree_search_timeout",
		expectedValues.MirroredSubtreeSearchTimeout, state.MirroredSubtreeSearchTimeout, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, internaltypes.IsEmptyString(expectedValues.BackupFilePermissions))
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, internaltypes.IsEmptyString(expectedValues.NotificationManager))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendNilSetsDefault(ctx, state)
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
	state.TaskRetentionTime = internaltypes.StringTypeOrNil(r.TaskRetentionTime, internaltypes.IsEmptyString(expectedValues.TaskRetentionTime))
	config.CheckMismatchedPDFormattedAttributes("task_retention_time",
		expectedValues.TaskRetentionTime, state.TaskRetentionTime, diagnostics)
	state.NotificationSenderAddress = internaltypes.StringTypeOrNil(r.NotificationSenderAddress, internaltypes.IsEmptyString(expectedValues.NotificationSenderAddress))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, internaltypes.IsEmptyString(expectedValues.BackupFilePermissions))
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, internaltypes.IsEmptyString(expectedValues.NotificationManager))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendNilSetsDefault(ctx, state)
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
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, internaltypes.IsEmptyString(expectedValues.BackupFilePermissions))
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, internaltypes.IsEmptyString(expectedValues.NotificationManager))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendNilSetsDefault(ctx, state)
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
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, internaltypes.IsEmptyString(expectedValues.BackupFilePermissions))
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, internaltypes.IsEmptyString(expectedValues.NotificationManager))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendNilSetsDefault(ctx, state)
}

// Read a MetricsBackendResponse object into the model struct
func readMetricsBackendResponseDefault(ctx context.Context, r *client.MetricsBackendResponse, state *defaultBackendResourceModel, expectedValues *defaultBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("metrics")
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.StorageDir = types.StringValue(r.StorageDir)
	state.MetricsDir = types.StringValue(r.MetricsDir)
	state.SampleFlushInterval = internaltypes.StringTypeOrNil(r.SampleFlushInterval, internaltypes.IsEmptyString(expectedValues.SampleFlushInterval))
	config.CheckMismatchedPDFormattedAttributes("sample_flush_interval",
		expectedValues.SampleFlushInterval, state.SampleFlushInterval, diagnostics)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, internaltypes.IsEmptyString(expectedValues.NotificationManager))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateBackendNilSetsDefault(ctx, state)
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
	addRequest := client.NewAddLocalDbBackendRequest(plan.BackendID.ValueString(),
		[]client.EnumlocalDbBackendSchemaUrn{client.ENUMLOCALDBBACKENDSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0BACKENDLOCAL_DB},
		plan.BackendID.ValueString(),
		plan.Enabled.ValueBool(),
		BaseDNSlice)
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
	apiAddRequest := r.apiClient.BackendApi.AddBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLocalDbBackendRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.BackendApi.AddBackendExecute(apiAddRequest)
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

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

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

	readResponse, httpResp, err := r.apiClient.BackendApi.GetBackend(
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
	if plan.Type.ValueString() == "schema" {
		readSchemaBackendResponseDefault(ctx, readResponse.SchemaBackendResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "backup" {
		readBackupBackendResponseDefault(ctx, readResponse.BackupBackendResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "encryption-settings" {
		readEncryptionSettingsBackendResponseDefault(ctx, readResponse.EncryptionSettingsBackendResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "ldif" {
		readLdifBackendResponseDefault(ctx, readResponse.LdifBackendResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "trust-store" {
		readTrustStoreBackendResponseDefault(ctx, readResponse.TrustStoreBackendResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "custom" {
		readCustomBackendResponseDefault(ctx, readResponse.CustomBackendResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "changelog" {
		readChangelogBackendResponseDefault(ctx, readResponse.ChangelogBackendResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "monitor" {
		readMonitorBackendResponseDefault(ctx, readResponse.MonitorBackendResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "local-db" {
		readLocalDbBackendResponseDefault(ctx, readResponse.LocalDbBackendResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "config-file-handler" {
		readConfigFileHandlerBackendResponseDefault(ctx, readResponse.ConfigFileHandlerBackendResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "task" {
		readTaskBackendResponseDefault(ctx, readResponse.TaskBackendResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "alert" {
		readAlertBackendResponseDefault(ctx, readResponse.AlertBackendResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "alarm" {
		readAlarmBackendResponseDefault(ctx, readResponse.AlarmBackendResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "metrics" {
		readMetricsBackendResponseDefault(ctx, readResponse.MetricsBackendResponse, &state, &plan, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.BackendApi.UpdateBackend(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.BackendID.ValueString())
	ops := createBackendOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.BackendApi.UpdateBackendExecute(updateRequest)
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
		if plan.Type.ValueString() == "schema" {
			readSchemaBackendResponseDefault(ctx, updateResponse.SchemaBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "backup" {
			readBackupBackendResponseDefault(ctx, updateResponse.BackupBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "encryption-settings" {
			readEncryptionSettingsBackendResponseDefault(ctx, updateResponse.EncryptionSettingsBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "ldif" {
			readLdifBackendResponseDefault(ctx, updateResponse.LdifBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "trust-store" {
			readTrustStoreBackendResponseDefault(ctx, updateResponse.TrustStoreBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "custom" {
			readCustomBackendResponseDefault(ctx, updateResponse.CustomBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "changelog" {
			readChangelogBackendResponseDefault(ctx, updateResponse.ChangelogBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "monitor" {
			readMonitorBackendResponseDefault(ctx, updateResponse.MonitorBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "local-db" {
			readLocalDbBackendResponseDefault(ctx, updateResponse.LocalDbBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "config-file-handler" {
			readConfigFileHandlerBackendResponseDefault(ctx, updateResponse.ConfigFileHandlerBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "task" {
			readTaskBackendResponseDefault(ctx, updateResponse.TaskBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "alert" {
			readAlertBackendResponseDefault(ctx, updateResponse.AlertBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "alarm" {
			readAlarmBackendResponseDefault(ctx, updateResponse.AlarmBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "metrics" {
			readMetricsBackendResponseDefault(ctx, updateResponse.MetricsBackendResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *backendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state backendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.BackendApi.GetBackend(
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
	if readResponse.LocalDbBackendResponse != nil {
		readLocalDbBackendResponse(ctx, readResponse.LocalDbBackendResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *defaultBackendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state defaultBackendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.BackendApi.GetBackend(
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
	if resp.Diagnostics.HasError() {
		return
	}
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
	updateRequest := r.apiClient.BackendApi.UpdateBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.BackendID.ValueString())

	// Determine what update operations are necessary
	ops := createBackendOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.BackendApi.UpdateBackendExecute(updateRequest)
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
		if plan.Type.ValueString() == "local-db" {
			readLocalDbBackendResponse(ctx, updateResponse.LocalDbBackendResponse, &state, &plan, &resp.Diagnostics)
		}
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
	updateRequest := r.apiClient.BackendApi.UpdateBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.BackendID.ValueString())

	// Determine what update operations are necessary
	ops := createBackendOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.BackendApi.UpdateBackendExecute(updateRequest)
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
		if plan.Type.ValueString() == "schema" {
			readSchemaBackendResponseDefault(ctx, updateResponse.SchemaBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "backup" {
			readBackupBackendResponseDefault(ctx, updateResponse.BackupBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "encryption-settings" {
			readEncryptionSettingsBackendResponseDefault(ctx, updateResponse.EncryptionSettingsBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "ldif" {
			readLdifBackendResponseDefault(ctx, updateResponse.LdifBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "trust-store" {
			readTrustStoreBackendResponseDefault(ctx, updateResponse.TrustStoreBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "custom" {
			readCustomBackendResponseDefault(ctx, updateResponse.CustomBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "changelog" {
			readChangelogBackendResponseDefault(ctx, updateResponse.ChangelogBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "monitor" {
			readMonitorBackendResponseDefault(ctx, updateResponse.MonitorBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "local-db" {
			readLocalDbBackendResponseDefault(ctx, updateResponse.LocalDbBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "config-file-handler" {
			readConfigFileHandlerBackendResponseDefault(ctx, updateResponse.ConfigFileHandlerBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "task" {
			readTaskBackendResponseDefault(ctx, updateResponse.TaskBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "alert" {
			readAlertBackendResponseDefault(ctx, updateResponse.AlertBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "alarm" {
			readAlarmBackendResponseDefault(ctx, updateResponse.AlarmBackendResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "metrics" {
			readMetricsBackendResponseDefault(ctx, updateResponse.MetricsBackendResponse, &state, &plan, &resp.Diagnostics)
		}
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

	httpResp, err := r.apiClient.BackendApi.DeleteBackendExecute(r.apiClient.BackendApi.DeleteBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.BackendID.ValueString()))
	if err != nil {
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
