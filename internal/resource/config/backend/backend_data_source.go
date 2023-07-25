package backend

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
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
	r.apiClient = providerCfg.ApiClientV9300
}

type backendDataSourceModel struct {
	Id                                        types.String `tfsdk:"id"`
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

// GetSchema defines the schema for the datasource.
func (r *backendDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Describes a Backend.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Name of this object.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
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
				Description: "Specifies the permissions that should be applied to the directory containing the backend database files and to directories and files created during backup of the backend.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"db_cache_percent": schema.Int64Attribute{
				Description: "Specifies the percentage of JVM memory to allocate to the changelog database cache.",
				Required:    false,
				Optional:    false,
				Computed:    true,
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
				Description: "Specifies the database and environment properties for the Berkeley DB Java Edition database for this changelog backend.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
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
			"is_private_backend": schema.BoolAttribute{
				Description: "Indicates whether the backend should be considered a private backend, which indicates that it is used for storing operational data rather than user-defined information.",
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
			"notification_manager": schema.StringAttribute{
				Description: "Specifies a notification manager for changes resulting from operations processed through this Backend",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
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

// Read resource information
func (r *backendDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state backendDataSourceModel
	diags := req.Config.Get(ctx, &state)
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
		readLocalDbBackendResponseDataSource(ctx, readResponse.LocalDbBackendResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
