package backend

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &localDbBackendResource{}
	_ resource.ResourceWithConfigure   = &localDbBackendResource{}
	_ resource.ResourceWithImportState = &localDbBackendResource{}
	_ resource.Resource                = &defaultLocalDbBackendResource{}
	_ resource.ResourceWithConfigure   = &defaultLocalDbBackendResource{}
	_ resource.ResourceWithImportState = &defaultLocalDbBackendResource{}
)

// Create a Local Db Backend resource
func NewLocalDbBackendResource() resource.Resource {
	return &localDbBackendResource{}
}

func NewDefaultLocalDbBackendResource() resource.Resource {
	return &defaultLocalDbBackendResource{}
}

// localDbBackendResource is the resource implementation.
type localDbBackendResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultLocalDbBackendResource is the resource implementation.
type defaultLocalDbBackendResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *localDbBackendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_db_backend"
}

func (r *defaultLocalDbBackendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_local_db_backend"
}

// Configure adds the provider configured client to the resource.
func (r *localDbBackendResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultLocalDbBackendResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type localDbBackendResourceModel struct {
	Id                                        types.String `tfsdk:"id"`
	LastUpdated                               types.String `tfsdk:"last_updated"`
	Notifications                             types.Set    `tfsdk:"notifications"`
	RequiredActions                           types.Set    `tfsdk:"required_actions"`
	UncachedId2entryCacheMode                 types.String `tfsdk:"uncached_id2entry_cache_mode"`
	UncachedAttributeCriteria                 types.String `tfsdk:"uncached_attribute_criteria"`
	UncachedEntryCriteria                     types.String `tfsdk:"uncached_entry_criteria"`
	WritabilityMode                           types.String `tfsdk:"writability_mode"`
	SetDegradedAlertForUntrustedIndex         types.Bool   `tfsdk:"set_degraded_alert_for_untrusted_index"`
	ReturnUnavailableForUntrustedIndex        types.Bool   `tfsdk:"return_unavailable_for_untrusted_index"`
	ProcessFiltersWithUndefinedAttributeTypes types.Bool   `tfsdk:"process_filters_with_undefined_attribute_types"`
	IsPrivateBackend                          types.Bool   `tfsdk:"is_private_backend"`
	DbDirectory                               types.String `tfsdk:"db_directory"`
	DbDirectoryPermissions                    types.String `tfsdk:"db_directory_permissions"`
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
	DbCachePercent                            types.Int64  `tfsdk:"db_cache_percent"`
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
	BackendID                                 types.String `tfsdk:"backend_id"`
	Description                               types.String `tfsdk:"description"`
	Enabled                                   types.Bool   `tfsdk:"enabled"`
	BaseDN                                    types.Set    `tfsdk:"base_dn"`
	SetDegradedAlertWhenDisabled              types.Bool   `tfsdk:"set_degraded_alert_when_disabled"`
	ReturnUnavailableWhenDisabled             types.Bool   `tfsdk:"return_unavailable_when_disabled"`
	NotificationManager                       types.String `tfsdk:"notification_manager"`
}

// GetSchema defines the schema for the resource.
func (r *localDbBackendResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	localDbBackendSchema(ctx, req, resp, false)
}

func (r *defaultLocalDbBackendResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	localDbBackendSchema(ctx, req, resp, true)
}

func localDbBackendSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Local Db Backend.",
		Attributes: map[string]schema.Attribute{
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
			"writability_mode": schema.StringAttribute{
				Description: "Specifies the behavior that the backend should use when processing write operations.",
				Optional:    true,
				Computed:    true,
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
			"is_private_backend": schema.BoolAttribute{
				Description: "Indicates whether this backend should be considered a private backend in the server. Private backends are meant for storing server-internal information and should not be used for user or application data.",
				Optional:    true,
				Computed:    true,
			},
			"db_directory": schema.StringAttribute{
				Description: "Specifies the path to the filesystem directory that is used to hold the Berkeley DB Java Edition database files containing the data for this backend. The files for this backend are stored in a sub-directory named after the backend-id.",
				Optional:    true,
				Computed:    true,
			},
			"db_directory_permissions": schema.StringAttribute{
				Description: "Specifies the permissions that should be applied to the directory containing the backend database files and to directories and files created during backup or LDIF export of the backend.",
				Optional:    true,
				Computed:    true,
			},
			"compact_common_parent_dn": schema.SetAttribute{
				Description: "Provides a DN of an entry that may be the parent for a large number of entries in the backend. This may be used to help increase the space efficiency when encoding entries for storage.",
				Optional:    true,
				Computed:    true,
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
			},
			"db_background_sync_interval": schema.StringAttribute{
				Description: "Specifies the interval to use when performing background synchronous writes in the database environment in order to smooth overall write performance and increase data durability. A value of \"0 s\" will disable background synchronous writes.",
				Optional:    true,
				Computed:    true,
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
				Description: "Specifies the database and environment properties for the Berkeley DB Java Edition database serving the data for this backend.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"db_cache_percent": schema.Int64Attribute{
				Description: "Specifies the percentage of JVM memory to allocate to the database cache.",
				Optional:    true,
				Computed:    true,
			},
			"default_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used for any database for which the cache mode is not explicitly specified. This includes the id2entry database, which stores encoded entries, and all attribute indexes.",
				Optional:    true,
				Computed:    true,
			},
			"id2entry_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the id2entry database, which provides a mapping between entry IDs and entry contents. Consider configuring uncached entries or uncached attributes in lieu of changing from the \"cache-keys-and-values\" default value.",
				Optional:    true,
				Computed:    true,
			},
			"dn2id_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the dn2id database, which provides a mapping between normalized entry DNs and the corresponding entry IDs.",
				Optional:    true,
				Computed:    true,
			},
			"id2children_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the id2children database, which provides a mapping between the entry ID of a particular entry and the entry IDs of all of its immediate children. This index may be used when performing searches with a single-level scope if the search filter cannot be resolved to a small enough candidate list. The size of this database directly depends on the number of entries that have children.",
				Optional:    true,
				Computed:    true,
			},
			"id2subtree_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the id2subtree database, which provides a mapping between the entry ID of a particular entry and the entry IDs of all of its children to any depth. This index may be used when performing searches with a whole-subtree or subordinate-subtree scope if the search filter cannot be resolved to a small enough candidate list. The size of this database directly depends on the number of entries that have children.",
				Optional:    true,
				Computed:    true,
			},
			"dn2uri_cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the dn2uri database, which provides a mapping between a normalized entry DN and a set of referral URLs contained in the associated smart referral entry.",
				Optional:    true,
				Computed:    true,
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
				ElementType: types.StringType,
			},
			"system_index_to_prime_internal_nodes_only": schema.SetAttribute{
				Description: "Specifies the system index(es) for which internal database nodes only (i.e., the database keys but not values) should be primed when the backend is initialized.",
				Optional:    true,
				Computed:    true,
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
				Computed:    true,
			},
			"id2subtree_index_entry_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entry IDs to maintain for each entry in the id2subtree system index (which keeps track of all descendants below an entry, to assist in otherwise unindexed searches with a whole-subtree or subordinate subtree scope). A value of 0 means there is no limit, however this could have a big impact on database size on disk and on server performance.",
				Optional:    true,
				Computed:    true,
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
			},
			"backend_id": schema.StringAttribute{
				Description: "Specifies a name to identify the associated backend.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
			"base_dn": schema.SetAttribute{
				Description: "Specifies the base DN(s) for the data that the backend handles.",
				Required:    true,
				ElementType: types.StringType,
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
			},
			"notification_manager": schema.StringAttribute{
				Description: "Specifies a notification manager for changes resulting from operations processed through this Backend",
				Optional:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, false)
	if setOptionalToComputed {
		config.SetOptionalAttributesToComputed(&schema)
	}
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalLocalDbBackendFields(ctx context.Context, addRequest *client.AddLocalDbBackendRequest, plan localDbBackendResourceModel) error {
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
		stringVal := plan.UncachedAttributeCriteria.ValueString()
		addRequest.UncachedAttributeCriteria = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UncachedEntryCriteria) {
		stringVal := plan.UncachedEntryCriteria.ValueString()
		addRequest.UncachedEntryCriteria = &stringVal
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
		boolVal := plan.SetDegradedAlertForUntrustedIndex.ValueBool()
		addRequest.SetDegradedAlertForUntrustedIndex = &boolVal
	}
	if internaltypes.IsDefined(plan.ReturnUnavailableForUntrustedIndex) {
		boolVal := plan.ReturnUnavailableForUntrustedIndex.ValueBool()
		addRequest.ReturnUnavailableForUntrustedIndex = &boolVal
	}
	if internaltypes.IsDefined(plan.ProcessFiltersWithUndefinedAttributeTypes) {
		boolVal := plan.ProcessFiltersWithUndefinedAttributeTypes.ValueBool()
		addRequest.ProcessFiltersWithUndefinedAttributeTypes = &boolVal
	}
	if internaltypes.IsDefined(plan.IsPrivateBackend) {
		boolVal := plan.IsPrivateBackend.ValueBool()
		addRequest.IsPrivateBackend = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DbDirectory) {
		stringVal := plan.DbDirectory.ValueString()
		addRequest.DbDirectory = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DbDirectoryPermissions) {
		stringVal := plan.DbDirectoryPermissions.ValueString()
		addRequest.DbDirectoryPermissions = &stringVal
	}
	if internaltypes.IsDefined(plan.CompactCommonParentDN) {
		var slice []string
		plan.CompactCommonParentDN.ElementsAs(ctx, &slice, false)
		addRequest.CompactCommonParentDN = slice
	}
	if internaltypes.IsDefined(plan.CompressEntries) {
		boolVal := plan.CompressEntries.ValueBool()
		addRequest.CompressEntries = &boolVal
	}
	if internaltypes.IsDefined(plan.HashEntries) {
		boolVal := plan.HashEntries.ValueBool()
		addRequest.HashEntries = &boolVal
	}
	if internaltypes.IsDefined(plan.DbNumCleanerThreads) {
		intVal := int32(plan.DbNumCleanerThreads.ValueInt64())
		addRequest.DbNumCleanerThreads = &intVal
	}
	if internaltypes.IsDefined(plan.DbCleanerMinUtilization) {
		intVal := int32(plan.DbCleanerMinUtilization.ValueInt64())
		addRequest.DbCleanerMinUtilization = &intVal
	}
	if internaltypes.IsDefined(plan.DbEvictorCriticalPercentage) {
		intVal := int32(plan.DbEvictorCriticalPercentage.ValueInt64())
		addRequest.DbEvictorCriticalPercentage = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DbCheckpointerWakeupInterval) {
		stringVal := plan.DbCheckpointerWakeupInterval.ValueString()
		addRequest.DbCheckpointerWakeupInterval = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DbBackgroundSyncInterval) {
		stringVal := plan.DbBackgroundSyncInterval.ValueString()
		addRequest.DbBackgroundSyncInterval = &stringVal
	}
	if internaltypes.IsDefined(plan.DbUseThreadLocalHandles) {
		boolVal := plan.DbUseThreadLocalHandles.ValueBool()
		addRequest.DbUseThreadLocalHandles = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DbLogFileMax) {
		stringVal := plan.DbLogFileMax.ValueString()
		addRequest.DbLogFileMax = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DbLoggingLevel) {
		stringVal := plan.DbLoggingLevel.ValueString()
		addRequest.DbLoggingLevel = &stringVal
	}
	if internaltypes.IsDefined(plan.JeProperty) {
		var slice []string
		plan.JeProperty.ElementsAs(ctx, &slice, false)
		addRequest.JeProperty = slice
	}
	if internaltypes.IsDefined(plan.DbCachePercent) {
		intVal := int32(plan.DbCachePercent.ValueInt64())
		addRequest.DbCachePercent = &intVal
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
		intVal := int32(plan.PrimeThreadCount.ValueInt64())
		addRequest.PrimeThreadCount = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PrimeTimeLimit) {
		stringVal := plan.PrimeTimeLimit.ValueString()
		addRequest.PrimeTimeLimit = &stringVal
	}
	if internaltypes.IsDefined(plan.PrimeAllIndexes) {
		boolVal := plan.PrimeAllIndexes.ValueBool()
		addRequest.PrimeAllIndexes = &boolVal
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
		boolVal := plan.BackgroundPrime.ValueBool()
		addRequest.BackgroundPrime = &boolVal
	}
	if internaltypes.IsDefined(plan.IndexEntryLimit) {
		intVal := int32(plan.IndexEntryLimit.ValueInt64())
		addRequest.IndexEntryLimit = &intVal
	}
	if internaltypes.IsDefined(plan.CompositeIndexEntryLimit) {
		intVal := int32(plan.CompositeIndexEntryLimit.ValueInt64())
		addRequest.CompositeIndexEntryLimit = &intVal
	}
	if internaltypes.IsDefined(plan.Id2childrenIndexEntryLimit) {
		intVal := int32(plan.Id2childrenIndexEntryLimit.ValueInt64())
		addRequest.Id2childrenIndexEntryLimit = &intVal
	}
	if internaltypes.IsDefined(plan.Id2subtreeIndexEntryLimit) {
		intVal := int32(plan.Id2subtreeIndexEntryLimit.ValueInt64())
		addRequest.Id2subtreeIndexEntryLimit = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ImportTempDirectory) {
		stringVal := plan.ImportTempDirectory.ValueString()
		addRequest.ImportTempDirectory = &stringVal
	}
	if internaltypes.IsDefined(plan.ImportThreadCount) {
		intVal := int32(plan.ImportThreadCount.ValueInt64())
		addRequest.ImportThreadCount = &intVal
	}
	if internaltypes.IsDefined(plan.ExportThreadCount) {
		intVal := int32(plan.ExportThreadCount.ValueInt64())
		addRequest.ExportThreadCount = &intVal
	}
	if internaltypes.IsDefined(plan.DbImportCachePercent) {
		intVal := int32(plan.DbImportCachePercent.ValueInt64())
		addRequest.DbImportCachePercent = &intVal
	}
	if internaltypes.IsDefined(plan.DbTxnWriteNoSync) {
		boolVal := plan.DbTxnWriteNoSync.ValueBool()
		addRequest.DbTxnWriteNoSync = &boolVal
	}
	if internaltypes.IsDefined(plan.DeadlockRetryLimit) {
		intVal := int32(plan.DeadlockRetryLimit.ValueInt64())
		addRequest.DeadlockRetryLimit = &intVal
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
		intVal := int32(plan.SubtreeDeleteSizeLimit.ValueInt64())
		addRequest.SubtreeDeleteSizeLimit = &intVal
	}
	if internaltypes.IsDefined(plan.NumRecentChanges) {
		intVal := int32(plan.NumRecentChanges.ValueInt64())
		addRequest.NumRecentChanges = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OfflineProcessDatabaseOpenTimeout) {
		stringVal := plan.OfflineProcessDatabaseOpenTimeout.ValueString()
		addRequest.OfflineProcessDatabaseOpenTimeout = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	if internaltypes.IsDefined(plan.SetDegradedAlertWhenDisabled) {
		boolVal := plan.SetDegradedAlertWhenDisabled.ValueBool()
		addRequest.SetDegradedAlertWhenDisabled = &boolVal
	}
	if internaltypes.IsDefined(plan.ReturnUnavailableWhenDisabled) {
		boolVal := plan.ReturnUnavailableWhenDisabled.ValueBool()
		addRequest.ReturnUnavailableWhenDisabled = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.NotificationManager) {
		stringVal := plan.NotificationManager.ValueString()
		addRequest.NotificationManager = &stringVal
	}
	return nil
}

// Read a LocalDbBackendResponse object into the model struct
func readLocalDbBackendResponse(ctx context.Context, r *client.LocalDbBackendResponse, state *localDbBackendResourceModel, expectedValues *localDbBackendResourceModel, diagnostics *diag.Diagnostics) {
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
}

// Create any update operations necessary to make the state match the plan
func createLocalDbBackendOperations(plan localDbBackendResourceModel, state localDbBackendResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.UncachedId2entryCacheMode, state.UncachedId2entryCacheMode, "uncached-id2entry-cache-mode")
	operations.AddStringOperationIfNecessary(&ops, plan.UncachedAttributeCriteria, state.UncachedAttributeCriteria, "uncached-attribute-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.UncachedEntryCriteria, state.UncachedEntryCriteria, "uncached-entry-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.WritabilityMode, state.WritabilityMode, "writability-mode")
	operations.AddBoolOperationIfNecessary(&ops, plan.SetDegradedAlertForUntrustedIndex, state.SetDegradedAlertForUntrustedIndex, "set-degraded-alert-for-untrusted-index")
	operations.AddBoolOperationIfNecessary(&ops, plan.ReturnUnavailableForUntrustedIndex, state.ReturnUnavailableForUntrustedIndex, "return-unavailable-for-untrusted-index")
	operations.AddBoolOperationIfNecessary(&ops, plan.ProcessFiltersWithUndefinedAttributeTypes, state.ProcessFiltersWithUndefinedAttributeTypes, "process-filters-with-undefined-attribute-types")
	operations.AddBoolOperationIfNecessary(&ops, plan.IsPrivateBackend, state.IsPrivateBackend, "is-private-backend")
	operations.AddStringOperationIfNecessary(&ops, plan.DbDirectory, state.DbDirectory, "db-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.DbDirectoryPermissions, state.DbDirectoryPermissions, "db-directory-permissions")
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
	operations.AddInt64OperationIfNecessary(&ops, plan.DbCachePercent, state.DbCachePercent, "db-cache-percent")
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
	operations.AddStringOperationIfNecessary(&ops, plan.BackendID, state.BackendID, "backend-id")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddBoolOperationIfNecessary(&ops, plan.SetDegradedAlertWhenDisabled, state.SetDegradedAlertWhenDisabled, "set-degraded-alert-when-disabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.ReturnUnavailableWhenDisabled, state.ReturnUnavailableWhenDisabled, "return-unavailable-when-disabled")
	operations.AddStringOperationIfNecessary(&ops, plan.NotificationManager, state.NotificationManager, "notification-manager")
	return ops
}

// Create a new resource
func (r *localDbBackendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan localDbBackendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var BaseDNSlice []string
	plan.BaseDN.ElementsAs(ctx, &BaseDNSlice, false)
	addRequest := client.NewAddLocalDbBackendRequest(plan.BackendID.ValueString(),
		[]client.EnumlocalDbBackendSchemaUrn{client.ENUMLOCALDBBACKENDSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0BACKENDLOCAL_DB},
		plan.BackendID.ValueString(),
		plan.Enabled.ValueBool(),
		BaseDNSlice)
	err := addOptionalLocalDbBackendFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Local Db Backend", err.Error())
		return
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Local Db Backend", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state localDbBackendResourceModel
	readLocalDbBackendResponse(ctx, addResponse.LocalDbBackendResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultLocalDbBackendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan localDbBackendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.BackendApi.GetBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.BackendID.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Local Db Backend", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state localDbBackendResourceModel
	readLocalDbBackendResponse(ctx, readResponse.LocalDbBackendResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.BackendApi.UpdateBackend(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.BackendID.ValueString())
	ops := createLocalDbBackendOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.BackendApi.UpdateBackendExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Local Db Backend", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLocalDbBackendResponse(ctx, updateResponse.LocalDbBackendResponse, &state, &plan, &resp.Diagnostics)
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
func (r *localDbBackendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLocalDbBackend(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLocalDbBackendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLocalDbBackend(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readLocalDbBackend(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state localDbBackendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.BackendApi.GetBackend(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.BackendID.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Local Db Backend", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLocalDbBackendResponse(ctx, readResponse.LocalDbBackendResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *localDbBackendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLocalDbBackend(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLocalDbBackendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLocalDbBackend(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateLocalDbBackend(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan localDbBackendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state localDbBackendResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.BackendApi.UpdateBackend(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.BackendID.ValueString())

	// Determine what update operations are necessary
	ops := createLocalDbBackendOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.BackendApi.UpdateBackendExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Local Db Backend", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLocalDbBackendResponse(ctx, updateResponse.LocalDbBackendResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultLocalDbBackendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *localDbBackendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state localDbBackendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.BackendApi.DeleteBackendExecute(r.apiClient.BackendApi.DeleteBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.BackendID.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Local Db Backend", err, httpResp)
		return
	}
}

func (r *localDbBackendResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocalDbBackend(ctx, req, resp)
}

func (r *defaultLocalDbBackendResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocalDbBackend(ctx, req, resp)
}

func importLocalDbBackend(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to backend_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("backend_id"), req, resp)
}
