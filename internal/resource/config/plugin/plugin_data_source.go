package plugin

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
	_ datasource.DataSource              = &pluginDataSource{}
	_ datasource.DataSourceWithConfigure = &pluginDataSource{}
)

// Create a Plugin data source
func NewPluginDataSource() datasource.DataSource {
	return &pluginDataSource{}
}

// pluginDataSource is the datasource implementation.
type pluginDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *pluginDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_plugin"
}

// Configure adds the provider configured client to the data source.
func (r *pluginDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type pluginDataSourceModel struct {
	Id                                                   types.String `tfsdk:"id"`
	Name                                                 types.String `tfsdk:"name"`
	ResourceType                                         types.String `tfsdk:"resource_type"`
	PassThroughAuthenticationHandler                     types.String `tfsdk:"pass_through_authentication_handler"`
	Type                                                 types.Set    `tfsdk:"type"`
	MultipleAttributeBehavior                            types.String `tfsdk:"multiple_attribute_behavior"`
	ScriptClass                                          types.String `tfsdk:"script_class"`
	PreventConflictsWithSoftDeletedEntries               types.Bool   `tfsdk:"prevent_conflicts_with_soft_deleted_entries"`
	ProfileSampleInterval                                types.String `tfsdk:"profile_sample_interval"`
	ExcludeAttribute                                     types.Set    `tfsdk:"exclude_attribute"`
	UpdateInterval                                       types.String `tfsdk:"update_interval"`
	ScriptArgument                                       types.Set    `tfsdk:"script_argument"`
	Delay                                                types.String `tfsdk:"delay"`
	SourceAttribute                                      types.String `tfsdk:"source_attribute"`
	TargetAttribute                                      types.String `tfsdk:"target_attribute"`
	ProfileDirectory                                     types.String `tfsdk:"profile_directory"`
	ValuePattern                                         types.Set    `tfsdk:"value_pattern"`
	MultipleValuePatternBehavior                         types.String `tfsdk:"multiple_value_pattern_behavior"`
	MultiValuedAttributeBehavior                         types.String `tfsdk:"multi_valued_attribute_behavior"`
	TargetAttributeExistsDuringInitialPopulationBehavior types.String `tfsdk:"target_attribute_exists_during_initial_population_behavior"`
	UpdateSourceAttributeBehavior                        types.String `tfsdk:"update_source_attribute_behavior"`
	SourceAttributeRemovalBehavior                       types.String `tfsdk:"source_attribute_removal_behavior"`
	UpdateTargetAttributeBehavior                        types.String `tfsdk:"update_target_attribute_behavior"`
	IncludeBaseDN                                        types.Set    `tfsdk:"include_base_dn"`
	ExcludeBaseDN                                        types.Set    `tfsdk:"exclude_base_dn"`
	IncludeFilter                                        types.Set    `tfsdk:"include_filter"`
	ExcludeFilter                                        types.Set    `tfsdk:"exclude_filter"`
	UpdatedEntryNewlyMatchesCriteriaBehavior             types.String `tfsdk:"updated_entry_newly_matches_criteria_behavior"`
	UpdatedEntryNoLongerMatchesCriteriaBehavior          types.String `tfsdk:"updated_entry_no_longer_matches_criteria_behavior"`
	EnableProfilingOnStartup                             types.Bool   `tfsdk:"enable_profiling_on_startup"`
	ProfileAction                                        types.String `tfsdk:"profile_action"`
	ContextName                                          types.String `tfsdk:"context_name"`
	DefaultUserPasswordStorageScheme                     types.Set    `tfsdk:"default_user_password_storage_scheme"`
	DefaultAuthPasswordStorageScheme                     types.Set    `tfsdk:"default_auth_password_storage_scheme"`
	AllowedRequestControl                                types.Set    `tfsdk:"allowed_request_control"`
	AgentxAddress                                        types.String `tfsdk:"agentx_address"`
	AgentxPort                                           types.Int64  `tfsdk:"agentx_port"`
	NumWorkerThreads                                     types.Int64  `tfsdk:"num_worker_threads"`
	SessionTimeout                                       types.String `tfsdk:"session_timeout"`
	ConnectRetryMaxWait                                  types.String `tfsdk:"connect_retry_max_wait"`
	PingInterval                                         types.String `tfsdk:"ping_interval"`
	ExtensionClass                                       types.String `tfsdk:"extension_class"`
	ReferralBaseURL                                      types.Set    `tfsdk:"referral_base_url"`
	SourceDN                                             types.String `tfsdk:"source_dn"`
	TargetDN                                             types.String `tfsdk:"target_dn"`
	EnableAttributeMapping                               types.Bool   `tfsdk:"enable_attribute_mapping"`
	MapAttribute                                         types.Set    `tfsdk:"map_attribute"`
	RetainFilesSparselyByAge                             types.Bool   `tfsdk:"retain_files_sparsely_by_age"`
	Sanitize                                             types.Bool   `tfsdk:"sanitize"`
	EnableControlMapping                                 types.Bool   `tfsdk:"enable_control_mapping"`
	AlwaysMapResponses                                   types.Bool   `tfsdk:"always_map_responses"`
	Server                                               types.Set    `tfsdk:"server"`
	ExtensionArgument                                    types.Set    `tfsdk:"extension_argument"`
	EncryptionSettingsDefinitionID                       types.String `tfsdk:"encryption_settings_definition_id"`
	DatetimeAttribute                                    types.String `tfsdk:"datetime_attribute"`
	DatetimeJSONField                                    types.String `tfsdk:"datetime_json_field"`
	ServerAccessMode                                     types.String `tfsdk:"server_access_mode"`
	NumMostExpensivePhasesShown                          types.Int64  `tfsdk:"num_most_expensive_phases_shown"`
	DatetimeFormat                                       types.String `tfsdk:"datetime_format"`
	CustomDatetimeFormat                                 types.String `tfsdk:"custom_datetime_format"`
	DnMap                                                types.Set    `tfsdk:"dn_map"`
	BindDNPattern                                        types.String `tfsdk:"bind_dn_pattern"`
	SearchBaseDN                                         types.String `tfsdk:"search_base_dn"`
	SearchFilterPattern                                  types.String `tfsdk:"search_filter_pattern"`
	InitialConnections                                   types.Int64  `tfsdk:"initial_connections"`
	MaxConnections                                       types.Int64  `tfsdk:"max_connections"`
	CustomTimezone                                       types.String `tfsdk:"custom_timezone"`
	ExpirationOffset                                     types.String `tfsdk:"expiration_offset"`
	PurgeBehavior                                        types.String `tfsdk:"purge_behavior"`
	LogInterval                                          types.String `tfsdk:"log_interval"`
	ChangelogPasswordEncryptionKey                       types.String `tfsdk:"changelog_password_encryption_key"`
	SuppressIfIdle                                       types.Bool   `tfsdk:"suppress_if_idle"`
	HeaderPrefixPerColumn                                types.Bool   `tfsdk:"header_prefix_per_column"`
	EmptyInsteadOfZero                                   types.Bool   `tfsdk:"empty_instead_of_zero"`
	LinesBetweenHeader                                   types.Int64  `tfsdk:"lines_between_header"`
	IncludedLDAPStat                                     types.Set    `tfsdk:"included_ldap_stat"`
	IncludedResourceStat                                 types.Set    `tfsdk:"included_resource_stat"`
	HistogramFormat                                      types.String `tfsdk:"histogram_format"`
	HistogramOpType                                      types.Set    `tfsdk:"histogram_op_type"`
	Scope                                                types.String `tfsdk:"scope"`
	HistogramCategoryBoundary                            types.Set    `tfsdk:"histogram_category_boundary"`
	IncludeAttribute                                     types.Set    `tfsdk:"include_attribute"`
	GaugeInfo                                            types.String `tfsdk:"gauge_info"`
	LogFileFormat                                        types.String `tfsdk:"log_file_format"`
	LogFile                                              types.String `tfsdk:"log_file"`
	LogFilePermissions                                   types.String `tfsdk:"log_file_permissions"`
	Append                                               types.Bool   `tfsdk:"append"`
	RotationPolicy                                       types.Set    `tfsdk:"rotation_policy"`
	RotationListener                                     types.Set    `tfsdk:"rotation_listener"`
	RetentionPolicy                                      types.Set    `tfsdk:"retention_policy"`
	LoggingErrorBehavior                                 types.String `tfsdk:"logging_error_behavior"`
	OutputFile                                           types.String `tfsdk:"output_file"`
	PreviousFileExtension                                types.String `tfsdk:"previous_file_extension"`
	IncludeQueueTime                                     types.Bool   `tfsdk:"include_queue_time"`
	SeparateMonitorEntryPerTrackedApplication            types.Bool   `tfsdk:"separate_monitor_entry_per_tracked_application"`
	ChangelogPasswordEncryptionKeyPassphraseProvider     types.String `tfsdk:"changelog_password_encryption_key_passphrase_provider"`
	ApiURL                                               types.String `tfsdk:"api_url"`
	AuthURL                                              types.String `tfsdk:"auth_url"`
	OAuthClientID                                        types.String `tfsdk:"oauth_client_id"`
	OAuthClientSecret                                    types.String `tfsdk:"oauth_client_secret"`
	OAuthClientSecretPassphraseProvider                  types.String `tfsdk:"oauth_client_secret_passphrase_provider"`
	EnvironmentID                                        types.String `tfsdk:"environment_id"`
	HttpProxyExternalServer                              types.String `tfsdk:"http_proxy_external_server"`
	IncludedLocalEntryBaseDN                             types.Set    `tfsdk:"included_local_entry_base_dn"`
	ConnectionCriteria                                   types.String `tfsdk:"connection_criteria"`
	PollingInterval                                      types.String `tfsdk:"polling_interval"`
	TryLocalBind                                         types.Bool   `tfsdk:"try_local_bind"`
	OverrideLocalPassword                                types.Bool   `tfsdk:"override_local_password"`
	UpdateLocalPassword                                  types.Bool   `tfsdk:"update_local_password"`
	UpdateLocalPasswordDN                                types.String `tfsdk:"update_local_password_dn"`
	AllowLaxPassThroughAuthenticationPasswords           types.Bool   `tfsdk:"allow_lax_pass_through_authentication_passwords"`
	IgnoredPasswordPolicyStateErrorCondition             types.Set    `tfsdk:"ignored_password_policy_state_error_condition"`
	UserMappingLocalAttribute                            types.Set    `tfsdk:"user_mapping_local_attribute"`
	UserMappingRemoteJSONField                           types.Set    `tfsdk:"user_mapping_remote_json_field"`
	AdditionalUserMappingSCIMFilter                      types.String `tfsdk:"additional_user_mapping_scim_filter"`
	InvokeGCDayOfWeek                                    types.Set    `tfsdk:"invoke_gc_day_of_week"`
	InvokeGCTimeUtc                                      types.Set    `tfsdk:"invoke_gc_time_utc"`
	DelayAfterAlert                                      types.String `tfsdk:"delay_after_alert"`
	DelayPostGC                                          types.String `tfsdk:"delay_post_gc"`
	PeerServerPriorityIndex                              types.Int64  `tfsdk:"peer_server_priority_index"`
	PluginType                                           types.Set    `tfsdk:"plugin_type"`
	MaxUpdatesPerSecond                                  types.Int64  `tfsdk:"max_updates_per_second"`
	NumDeleteThreads                                     types.Int64  `tfsdk:"num_delete_threads"`
	AttributeType                                        types.Set    `tfsdk:"attribute_type"`
	Filter                                               types.Set    `tfsdk:"filter"`
	NumThreads                                           types.Int64  `tfsdk:"num_threads"`
	BaseDN                                               types.Set    `tfsdk:"base_dn"`
	LowerBound                                           types.Int64  `tfsdk:"lower_bound"`
	UpperBound                                           types.Int64  `tfsdk:"upper_bound"`
	FilterPrefix                                         types.String `tfsdk:"filter_prefix"`
	FilterSuffix                                         types.String `tfsdk:"filter_suffix"`
	SampleInterval                                       types.String `tfsdk:"sample_interval"`
	CollectionInterval                                   types.String `tfsdk:"collection_interval"`
	LdapInfo                                             types.String `tfsdk:"ldap_info"`
	ServerInfo                                           types.String `tfsdk:"server_info"`
	PerApplicationLDAPStats                              types.String `tfsdk:"per_application_ldap_stats"`
	LdapChangelogInfo                                    types.String `tfsdk:"ldap_changelog_info"`
	StatusSummaryInfo                                    types.String `tfsdk:"status_summary_info"`
	GenerateCollectorFiles                               types.Bool   `tfsdk:"generate_collector_files"`
	LocalDBBackendInfo                                   types.String `tfsdk:"local_db_backend_info"`
	ReplicationInfo                                      types.String `tfsdk:"replication_info"`
	EntryCacheInfo                                       types.String `tfsdk:"entry_cache_info"`
	HostInfo                                             types.Set    `tfsdk:"host_info"`
	IncludedLDAPApplication                              types.Set    `tfsdk:"included_ldap_application"`
	MaxUpdateFrequency                                   types.String `tfsdk:"max_update_frequency"`
	OperationType                                        types.Set    `tfsdk:"operation_type"`
	InvokeForFailedBinds                                 types.Bool   `tfsdk:"invoke_for_failed_binds"`
	MaxSearchResultEntriesToUpdate                       types.Int64  `tfsdk:"max_search_result_entries_to_update"`
	RequestCriteria                                      types.String `tfsdk:"request_criteria"`
	InvokeForInternalOperations                          types.Bool   `tfsdk:"invoke_for_internal_operations"`
	Description                                          types.String `tfsdk:"description"`
	Enabled                                              types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *pluginDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Plugin.",
		Attributes: map[string]schema.Attribute{
			"resource_type": schema.StringAttribute{
				Description: "The type of Plugin resource. Options are ['last-access-time', 'stats-collector', 'internal-search-rate', 'modifiable-password-policy-state', 'seven-bit-clean', 'clean-up-expired-pingfederate-persistent-access-grants', 'periodic-gc', 'ping-one-pass-through-authentication', 'changelog-password-encryption', 'processing-time-histogram', 'search-shutdown', 'periodic-stats-logger', 'purge-expired-data', 'change-subscription-notification', 'sub-operation-timing', 'third-party', 'encrypt-attribute-values', 'pass-through-authentication', 'dn-mapper', 'monitor-history', 'referral-on-update', 'simple-to-external-bind', 'custom', 'snmp-subagent', 'coalesce-modifications', 'password-policy-import', 'profiler', 'clean-up-inactive-pingfederate-persistent-sessions', 'composed-attribute', 'ldap-result-code-tracker', 'attribute-mapper', 'delay', 'clean-up-expired-pingfederate-persistent-sessions', 'groovy-scripted', 'last-mod', 'pluggable-pass-through-authentication', 'referential-integrity', 'unique-attribute']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"pass_through_authentication_handler": schema.StringAttribute{
				Description: "The component used to manage authentication with the external authentication service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"type": schema.SetAttribute{
				Description: "Specifies the type of attributes to check for value uniqueness.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"multiple_attribute_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit if multiple attribute types are specified.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Plugin.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"prevent_conflicts_with_soft_deleted_entries": schema.BoolAttribute{
				Description: "Indicates whether this Unique Attribute Plugin should reject a change that would result in one or more conflicts, even if those conflicts only exist in soft-deleted entries.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"profile_sample_interval": schema.StringAttribute{
				Description: "Specifies the sample interval in milliseconds to be used when capturing profiling information in the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"exclude_attribute": schema.SetAttribute{
				Description: "Specifies the name or OID of an attribute type which may be updated in a modify or modify DN operation without causing the modifiersName and modifyTimestamp values to be updated for that entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"update_interval": schema.StringAttribute{
				Description: "Specifies the interval in seconds when referential integrity updates are made.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Plugin. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"delay": schema.StringAttribute{
				Description: "The delay to inject for operations matching the associated criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"source_attribute": schema.StringAttribute{
				Description: "Specifies the source attribute type that may appear in client requests which should be remapped to the target attribute. Note that the source attribute type must be defined in the server schema and must not be equal to the target attribute type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"target_attribute": schema.StringAttribute{
				Description: "Specifies the target attribute type to which the source attribute type should be mapped. Note that the target attribute type must be defined in the server schema and must not be equal to the source attribute type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"profile_directory": schema.StringAttribute{
				Description: "Specifies the path to the directory where profile information is to be written. This path may be either an absolute path or a path that is relative to the root of the Directory Server instance.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"value_pattern": schema.SetAttribute{
				Description: "Specifies a pattern for constructing the values to use for the target attribute type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"multiple_value_pattern_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit if the plugin is configured with multiple value patterns.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"multi_valued_attribute_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit for source attributes that have multiple values.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"target_attribute_exists_during_initial_population_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit if the target attribute exists when initially populating the entry with composed values (whether during an LDIF import, an add operation, or an invocation of the populate composed attribute values task).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"update_source_attribute_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit for modify and modify DN operations that update one or more of the source attributes used in any of the value patterns.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"source_attribute_removal_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit for modify and modify DN operations that update an entry to remove source attributes in such a way that this plugin would no longer generate any composed values for that entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"update_target_attribute_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit for modify and modify DN operations that attempt to update the set of values for the target attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_base_dn": schema.SetAttribute{
				Description: "The set of base DNs below which composed values may be generated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"exclude_base_dn": schema.SetAttribute{
				Description: "The set of base DNs below which composed values will not be generated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"include_filter": schema.SetAttribute{
				Description: "The set of search filters that identify entries for which composed values may be generated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"exclude_filter": schema.SetAttribute{
				Description: "The set of search filters that identify entries for which composed values will not be generated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"updated_entry_newly_matches_criteria_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit for modify or modify DN operations that update an entry that previously did not satisfy either the base DN or filter criteria, but now do satisfy that criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"updated_entry_no_longer_matches_criteria_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit for modify or modify DN operations that update an entry that previously satisfied the base DN and filter criteria, but now no longer satisfies that criteria.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enable_profiling_on_startup": schema.BoolAttribute{
				Description: "Indicates whether the profiler plug-in is to start collecting data automatically when the Directory Server is started.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"profile_action": schema.StringAttribute{
				Description: "Specifies the action that should be taken by the profiler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"context_name": schema.StringAttribute{
				Description: "The SNMP context name for this sub-agent. The context name must not be longer than 30 ASCII characters. Each server in a topology must have a unique SNMP context name.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"default_user_password_storage_scheme": schema.SetAttribute{
				Description: "Specifies the names of the password storage schemes to be used for encoding passwords contained in attributes with the user password syntax for entries that do not include the ds-pwp-password-policy-dn attribute specifying which password policy is to be used to govern them.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"default_auth_password_storage_scheme": schema.SetAttribute{
				Description: "Specifies the names of password storage schemes that to be used for encoding passwords contained in attributes with the auth password syntax for entries that do not include the ds-pwp-password-policy-dn attribute specifying which password policy should be used to govern them.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"allowed_request_control": schema.SetAttribute{
				Description: "Specifies the OIDs of the controls that are allowed to be present in operations to coalesce. These controls are passed through when the request is validated, but they will not be included when the background thread applies the coalesced modify requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"agentx_address": schema.StringAttribute{
				Description: "The hostname or IP address of the SNMP master agent.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"agentx_port": schema.Int64Attribute{
				Description: "The port number on which the SNMP master agent will be contacted.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"num_worker_threads": schema.Int64Attribute{
				Description: "The number of worker threads to use to handle SNMP requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"session_timeout": schema.StringAttribute{
				Description: "Specifies the maximum amount of time to wait for a session to the master agent to be established.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"connect_retry_max_wait": schema.StringAttribute{
				Description: "The maximum amount of time to wait between attempts to establish a connection to the master agent.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ping_interval": schema.StringAttribute{
				Description: "The amount of time between consecutive pings sent by the sub-agent on its connection to the master agent. A value of zero disables the sending of pings by the sub-agent.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Plugin.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"referral_base_url": schema.SetAttribute{
				Description: "Specifies the base URL to use for the referrals generated by this plugin. It should include only the scheme, address, and port to use to communicate with the target server (e.g., \"ldap://server.example.com:389/\").",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"source_dn": schema.StringAttribute{
				Description: "Specifies the source DN that may appear in client requests which should be remapped to the target DN. Note that the source DN must not be equal to the target DN.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"target_dn": schema.StringAttribute{
				Description: "Specifies the DN to which the source DN should be mapped. Note that the target DN must not be equal to the source DN.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enable_attribute_mapping": schema.BoolAttribute{
				Description: "Indicates whether DN mapping should be applied to the values of attributes with appropriate syntaxes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"map_attribute": schema.SetAttribute{
				Description: "Specifies a set of specific attributes for which DN mapping should be applied. This will only be applicable if the enable-attribute-mapping property has a value of \"true\". Any attributes listed must be defined in the server schema with either the distinguished name syntax or the name and optional UID syntax.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"retain_files_sparsely_by_age": schema.BoolAttribute{
				Description: "Retain some older files to give greater perspective on how monitoring information has changed over time.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"sanitize": schema.BoolAttribute{
				Description: "Server monitoring data can include a small amount of personally identifiable information in the form of LDAP DNs and search filters. Setting this property to true will redact this information from the monitor files. This should only be used when necessary, as it reduces the information available in the archive and can increase the time to find the source of support issues.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enable_control_mapping": schema.BoolAttribute{
				Description: "Indicates whether DN mapping should be applied to DNs that may be present in specific controls. DN mapping will only be applied for control types which are specifically supported by the DN mapper plugin.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"always_map_responses": schema.BoolAttribute{
				Description: "Indicates whether DNs in response messages containing the target DN should always be remapped back to the source DN. If this is \"false\", then mapping will be performed for a response message only if one or more elements of the associated request are mapped. Otherwise, the mapping will be performed for all responses regardless of whether the mapping was applied to the request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server": schema.SetAttribute{
				Description: "Specifies the LDAP external server(s) to which authentication attempts should be forwarded.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Plugin. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"encryption_settings_definition_id": schema.StringAttribute{
				Description: "Specifies the ID of the encryption settings definition that should be used to encrypt the data. If this is not provided, the server's preferred encryption settings definition will be used. The \"encryption-settings list\" command can be used to obtain a list of the encryption settings definitions available in the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"datetime_attribute": schema.StringAttribute{
				Description: "The LDAP attribute that determines when data should be deleted. This could store the expiration time, or it could store the creation time and the expiration-offset property specifies the duration before data is deleted.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"datetime_json_field": schema.StringAttribute{
				Description: "The top-level JSON field within the configured datetime-attribute that determines when data should be deleted. This could store the expiration time, or it could store the creation time and the expiration-offset property specifies the duration before data is deleted.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_access_mode": schema.StringAttribute{
				Description: "Specifies the manner in which external servers should be used for pass-through authentication attempts if multiple servers are defined.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"num_most_expensive_phases_shown": schema.Int64Attribute{
				Description: "This controls how many of the most expensive phases are included per operation type in the monitor entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"datetime_format": schema.StringAttribute{
				Description: "Specifies the format of the datetime stored within the entry that determines when data should be purged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"custom_datetime_format": schema.StringAttribute{
				Description: "When the datetime-format property is configured with a value of \"custom\", this specifies the format (using a string compatible with the java.text.SimpleDateFormat class) that will be used to search for expired data.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"dn_map": schema.SetAttribute{
				Description: "Specifies one or more DN mappings that may be used to transform bind DNs before attempting to bind to the external servers.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"bind_dn_pattern": schema.StringAttribute{
				Description: "A pattern to use to construct the bind DN for the simple bind request to send to the remote server. This may consist of a combination of static text and attribute values and other directives enclosed in curly braces.  For example, the value \"cn={cn},ou=People,dc=example,dc=com\" indicates that the remote bind DN should be constructed from the text \"cn=\" followed by the value of the local entry's cn attribute followed by the text \"ou=People,dc=example,dc=com\". If an attribute contains the value to use as the bind DN for pass-through authentication, then the pattern may simply be the name of that attribute in curly braces (e.g., if the seeAlso attribute contains the bind DN for the target user, then a bind DN pattern of \"{seeAlso}\" would be appropriate).  Note that a bind DN pattern can be used to construct a bind DN that is not actually a valid LDAP distinguished name. For example, if authentication is being passed through to a Microsoft Active Directory server, then a bind DN pattern could be used to construct a user principal name (UPN) as an alternative to a distinguished name.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"search_base_dn": schema.StringAttribute{
				Description: "The base DN to use when searching for the user entry using a filter constructed from the pattern defined in the search-filter-pattern property. If no base DN is specified, the null DN will be used as the search base DN.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"search_filter_pattern": schema.StringAttribute{
				Description: "A pattern to use to construct a filter to use when searching an external server for the entry of the user as whom to bind. For example, \"(mail={uid:ldapFilterEscape}@example.com)\" would construct a search filter to search for a user whose entry in the local server contains a uid attribute whose value appears before \"@example.com\" in the mail attribute in the external server. Note that the \"ldapFilterEscape\" modifier should almost always be used with attributes specified in the pattern.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"initial_connections": schema.Int64Attribute{
				Description: "Specifies the initial number of connections to establish to each external server against which authentication may be attempted.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_connections": schema.Int64Attribute{
				Description: "Specifies the maximum number of connections to maintain to each external server against which authentication may be attempted. This value must be greater than or equal to the value for the initial-connections property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"custom_timezone": schema.StringAttribute{
				Description: "Specifies the time zone to use when generating a date string using the configured custom-datetime-format value. The provided value must be accepted by java.util.TimeZone.getTimeZone.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"expiration_offset": schema.StringAttribute{
				Description: "The duration to wait after the value specified in datetime-attribute (and optionally datetime-json-field) before purging the data.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"purge_behavior": schema.StringAttribute{
				Description: "Specifies whether to delete expired entries or attribute values. By default entries are deleted.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_interval": schema.StringAttribute{
				Description: "The duration between statistics collection and logging. A new line is logged to the output for each interval. Setting this value too small can have an impact on performance.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"changelog_password_encryption_key": schema.StringAttribute{
				Description: "A passphrase that may be used to generate the key for encrypting passwords stored in the changelog. The same passphrase also needs to be set (either through the \"changelog-password-decryption-key\" property or the \"changelog-password-decryption-key-passphrase-provider\" property) in the Global Sync Configuration in the Data Sync Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"suppress_if_idle": schema.BoolAttribute{
				Description: "If the server is idle during the specified interval, then do not log any output if this property is set to true. The server is idle if during the interval, no new connections were established, no operations were processed, and no operations are pending.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"header_prefix_per_column": schema.BoolAttribute{
				Description: "This property controls whether the header prefix, which applies to a group of columns, appears at the start of each column header or only the first column in a group.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"empty_instead_of_zero": schema.BoolAttribute{
				Description: "This property controls whether a value in the output is shown as empty if the value is zero.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"lines_between_header": schema.Int64Attribute{
				Description: "The number of lines to log between logging the header line that summarizes the columns in the table.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"included_ldap_stat": schema.SetAttribute{
				Description: "Specifies the types of statistics related to LDAP connections and operation processing that should be included in the output.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_resource_stat": schema.SetAttribute{
				Description: "Specifies whether statistics related to resource utilization such as JVM memory.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"histogram_format": schema.StringAttribute{
				Description: "The format of the data in the processing time histogram.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"histogram_op_type": schema.SetAttribute{
				Description: "Specifies the operation type(s) to use when outputting the response time histogram data. The order of the operations here determines the order of the columns in the output. Use the per-application-ldap-stats setting to further control this.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"scope": schema.StringAttribute{
				Description: "The scope to use for the search.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"histogram_category_boundary": schema.SetAttribute{
				Description: "Specifies the boundary values that will be used to separate the processing times into categories. Values should be specified as durations, and all values must be greater than zero.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"include_attribute": schema.SetAttribute{
				Description: "The name of an attribute that should be included in the results. This may include any token which is allowed as a requested attribute in search requests, including the name of an attribute, an asterisk (to indicate all user attributes), a plus sign (to indicate all operational attributes), an object class name preceded with an at symbol (to indicate all attributes associated with that object class), an attribute name preceded by a caret (to indicate that attribute should be excluded), or an object class name preceded by a caret and an at symbol (to indicate that all attributes associated with that object class should be excluded).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"gauge_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include for Gauges.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_file_format": schema.StringAttribute{
				Description: "Specifies the format to use when logging server statistics.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_file": schema.StringAttribute{
				Description: "The file name to use for the log files generated by the Periodic Stats Logger Plugin. The path to the file can be specified either as relative to the server root or as an absolute path.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_file_permissions": schema.StringAttribute{
				Description: "The UNIX permissions of the log files created by this Periodic Stats Logger Plugin.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"append": schema.BoolAttribute{
				Description: "Specifies whether to append to existing log files.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"rotation_policy": schema.SetAttribute{
				Description: "The rotation policy to use for the Periodic Stats Logger Plugin .",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"rotation_listener": schema.SetAttribute{
				Description: "A listener that should be notified whenever a log file is rotated out of service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"retention_policy": schema.SetAttribute{
				Description: "The retention policy to use for the Periodic Stats Logger Plugin .",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"logging_error_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the server should exhibit if an error occurs during logging processing.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"output_file": schema.StringAttribute{
				Description: "The path of an LDIF file that should be created with the results of the search.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"previous_file_extension": schema.StringAttribute{
				Description: "An extension that should be appended to the name of an existing output file rather than deleting it. If a file already exists with the full previous file name, then it will be deleted before the current file is renamed to become the previous file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_queue_time": schema.BoolAttribute{
				Description: "Indicates whether operation processing times should include the time spent waiting on the work queue. This will only be available if the work queue is configured to monitor the queue time.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"separate_monitor_entry_per_tracked_application": schema.BoolAttribute{
				Description: "When enabled, separate monitor entries will be included for each application defined in the Global Configuration's tracked-application property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"changelog_password_encryption_key_passphrase_provider": schema.StringAttribute{
				Description: "A passphrase provider that may be used to obtain the passphrase that will be used to generate the key for encrypting passwords stored in the changelog. The same passphrase also needs to be set (either through the \"changelog-password-decryption-key\" property or the \"changelog-password-decryption-key-passphrase-provider\" property) in the Global Sync Configuration in the Data Sync Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"api_url": schema.StringAttribute{
				Description: "Specifies the API endpoint for the PingOne web service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"auth_url": schema.StringAttribute{
				Description: "Specifies the API endpoint for the PingOne authentication service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"oauth_client_id": schema.StringAttribute{
				Description: "Specifies the OAuth Client ID used to authenticate connections to the PingOne API.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"oauth_client_secret": schema.StringAttribute{
				Description: "Specifies the OAuth Client Secret used to authenticate connections to the PingOne API.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"oauth_client_secret_passphrase_provider": schema.StringAttribute{
				Description: "Specifies a passphrase provider that can be used to obtain the OAuth Client Secret used to authenticate connections to the PingOne API.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"environment_id": schema.StringAttribute{
				Description: "Specifies the PingOne Environment that will be associated with this PingOne Pass Through Authentication Plugin.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"http_proxy_external_server": schema.StringAttribute{
				Description: "A reference to an HTTP proxy server that should be used for requests sent to the PingOne service. Supported in PingDirectory product version 9.2.0.0+.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"included_local_entry_base_dn": schema.SetAttribute{
				Description: "The base DNs for the local users whose authentication attempts may be passed through to the PingOne service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"connection_criteria": schema.StringAttribute{
				Description: "A reference to connection criteria that will be used to indicate which bind requests should be passed through to the PingOne service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"polling_interval": schema.StringAttribute{
				Description: "This specifies how often the plugin should check for expired data. It also controls the offset of peer servers (see the peer-server-priority-index for more information).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"try_local_bind": schema.BoolAttribute{
				Description: "Indicates whether to attempt the bind in the local server first, or to only send it to the PingOne service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"override_local_password": schema.BoolAttribute{
				Description: "Indicates whether to attempt the authentication in the PingOne service if the local user entry includes a password. This property will only be used if try-local-bind is true.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"update_local_password": schema.BoolAttribute{
				Description: "Indicates whether to overwrite the user's local password if the local bind fails but the authentication attempt succeeds when attempted in the PingOne service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"update_local_password_dn": schema.StringAttribute{
				Description: "This is the DN of the user that will be used to overwrite the user's local password if update-local-password is set. The DN put here should be added to 'ignore-changes-by-dn' in the appropriate Sync Source.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_lax_pass_through_authentication_passwords": schema.BoolAttribute{
				Description: "Indicates whether to overwrite the user's local password even if the password used to authenticate to the PingOne service would have failed validation if the user attempted to set it directly.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ignored_password_policy_state_error_condition": schema.SetAttribute{
				Description: "A set of password policy state error conditions that should not be enforced when authentication succeeds when attempted in the PingOne service. This option can only be used if try-local-bind is true.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"user_mapping_local_attribute": schema.SetAttribute{
				Description: "The names of the attributes in the local user entry whose values must match the values of the corresponding fields in the PingOne service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"user_mapping_remote_json_field": schema.SetAttribute{
				Description: "The names of the fields in the PingOne service whose values must match the values of the corresponding attributes in the local user entry, as specified in the user-mapping-local-attribute property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"additional_user_mapping_scim_filter": schema.StringAttribute{
				Description: "An optional SCIM filter that will be ANDed with the filter created to identify the account in the PingOne service that corresponds to the local entry. Only the \"eq\", \"sw\", \"and\", and \"or\" filter types may be used.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"invoke_gc_day_of_week": schema.SetAttribute{
				Description: "Specifies the days of the week which the Periodic GC Plugin should run. If no values are provided, then the plugin will run every day at the specified time.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"invoke_gc_time_utc": schema.SetAttribute{
				Description: "Specifies the times of the day at which garbage collection may be explicitly invoked. The times should be specified in \"HH:MM\" format, with \"HH\" as a two-digit numeric value between 00 and 23 representing the hour of the day, and MM as a two-digit numeric value between 00 and 59 representing the minute of the hour. All times will be interpreted in the UTC time zone.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"delay_after_alert": schema.StringAttribute{
				Description: "Specifies the length of time that the Directory Server should wait after sending the \"force-gc-starting\" administrative alert before actually invoking the garbage collection processing.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"delay_post_gc": schema.StringAttribute{
				Description: "Specifies the length of time that the Directory Server should wait after successfully completing the garbage collection processing, before removing the \"force-gc-starting\" administrative alert, which marks the server as unavailable.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"peer_server_priority_index": schema.Int64Attribute{
				Description: "In a replicated environment, this determines the order in which peer servers should attempt to purge data.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_type": schema.SetAttribute{
				Description: "Specifies the set of plug-in types for the plug-in, which specifies the times at which the plug-in is invoked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"max_updates_per_second": schema.Int64Attribute{
				Description: "This setting smooths out the performance impact on the server by throttling the purging to the specified maximum number of updates per second. To avoid a large backlog, this value should be set comfortably above the average rate that expired data is generated. When purge-behavior is set to subtree-delete-entries, then deletion of the entire subtree is considered a single update for the purposes of throttling.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"num_delete_threads": schema.Int64Attribute{
				Description: "The number of threads used to delete expired entries.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"attribute_type": schema.SetAttribute{
				Description: "Specifies the name or OID of an attribute type for which values should be checked to ensure that they are 7-bit clean.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"filter": schema.SetAttribute{
				Description: "A filter that may be used to identify entries that should support the ds-pwp-modifiable-state-json operational attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"num_threads": schema.Int64Attribute{
				Description: "Specifies the number of concurrent threads that should be used to process the search operations.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"base_dn": schema.SetAttribute{
				Description: "A base DN that may be used to identify entries that should support the ds-pwp-modifiable-state-json operational attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"lower_bound": schema.Int64Attribute{
				Description: "Specifies the lower bound for the numeric value which will be inserted into the search filter.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"upper_bound": schema.Int64Attribute{
				Description: "Specifies the upper bound for the numeric value which will be inserted into the search filter.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"filter_prefix": schema.StringAttribute{
				Description: "Specifies a prefix which will be used in front of the randomly-selected numeric value in all search filters used. If no upper bound is defined, then this should contain the entire filter string.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"filter_suffix": schema.StringAttribute{
				Description: "Specifies a suffix which will be used after of the randomly-selected numeric value in all search filters used. If no upper bound is defined, then this should be omitted.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"sample_interval": schema.StringAttribute{
				Description: "The duration between statistics collections. Setting this value too small can have an impact on performance. This value should be a multiple of collection-interval.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"collection_interval": schema.StringAttribute{
				Description: "Some of the calculated statistics, such as the average and maximum queue sizes, can use multiple samples within a log interval. This value controls how often samples are gathered, and setting this value too small can have an adverse impact on performance.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ldap_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include about the LDAP connection handlers.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_info": schema.StringAttribute{
				Description: "Specifies whether statistics related to resource utilization such as JVM memory and CPU/Network/Disk utilization.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"per_application_ldap_stats": schema.StringAttribute{
				Description: "Controls whether per application LDAP statistics are included in the output for selected LDAP operation statistics.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ldap_changelog_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include for the LDAP changelog.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"status_summary_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include about the status summary monitor entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"generate_collector_files": schema.BoolAttribute{
				Description: "Indicates whether this plugin should store metric samples on disk for use by the Data Metrics Server. If the Stats Collector Plugin is only being used to collect metrics for one or more StatsD Monitoring Endpoints, then this can be set to false to prevent unnecessary I/O.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"local_db_backend_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include about the Local DB Backends.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"replication_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include about replication.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"entry_cache_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include for each entry cache.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"host_info": schema.SetAttribute{
				Description: "Specifies the level of detail to include about the host system resource utilization including CPU, memory, disk and network activity.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_ldap_application": schema.SetAttribute{
				Description: "If statistics should not be included for all applications, this property names the subset of applications that should be included.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"max_update_frequency": schema.StringAttribute{
				Description: "Specifies the maximum frequency with which last access time values should be written for an entry. This may help limit the rate of internal write operations processed in the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"operation_type": schema.SetAttribute{
				Description: "Specifies the types of operations that should result in access time updates.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"invoke_for_failed_binds": schema.BoolAttribute{
				Description: "Indicates whether to update the last access time for an entry targeted by a bind operation if the bind is unsuccessful.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_search_result_entries_to_update": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that should be updated in a search operation. Only search result entries actually returned to the client may have their last access time updated, but because a single search operation may return a very large number of entries, the plugin will only update entries if no more than a specified number of entries are updated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"request_criteria": schema.StringAttribute{
				Description: "Specifies a set of request criteria that may be used to indicate whether to apply access time updates for the associated operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"invoke_for_internal_operations": schema.BoolAttribute{
				Description: "Indicates whether the plug-in should be invoked for internal operations.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Plugin",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the plug-in is enabled for use.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a LastAccessTimePluginResponse object into the model struct
func readLastAccessTimePluginResponseDataSource(ctx context.Context, r *client.LastAccessTimePluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("last-access-time")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MaxUpdateFrequency = internaltypes.StringTypeOrNil(r.MaxUpdateFrequency, false)
	state.OperationType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginOperationTypeProp(r.OperationType))
	state.InvokeForFailedBinds = internaltypes.BoolTypeOrNil(r.InvokeForFailedBinds)
	state.MaxSearchResultEntriesToUpdate = internaltypes.Int64TypeOrNil(r.MaxSearchResultEntriesToUpdate)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a StatsCollectorPluginResponse object into the model struct
func readStatsCollectorPluginResponseDataSource(ctx context.Context, r *client.StatsCollectorPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("stats-collector")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SampleInterval = types.StringValue(r.SampleInterval)
	state.CollectionInterval = types.StringValue(r.CollectionInterval)
	state.LdapInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLdapInfoProp(r.LdapInfo), false)
	state.ServerInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginServerInfoProp(r.ServerInfo), false)
	state.PerApplicationLDAPStats = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginStatsCollectorPerApplicationLDAPStatsProp(r.PerApplicationLDAPStats), false)
	state.LdapChangelogInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLdapChangelogInfoProp(r.LdapChangelogInfo), false)
	state.StatusSummaryInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginStatusSummaryInfoProp(r.StatusSummaryInfo), false)
	state.GenerateCollectorFiles = internaltypes.BoolTypeOrNil(r.GenerateCollectorFiles)
	state.LocalDBBackendInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLocalDBBackendInfoProp(r.LocalDBBackendInfo), false)
	state.ReplicationInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginReplicationInfoProp(r.ReplicationInfo), false)
	state.EntryCacheInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginEntryCacheInfoProp(r.EntryCacheInfo), false)
	state.HostInfo = internaltypes.GetStringSet(
		client.StringSliceEnumpluginHostInfoProp(r.HostInfo))
	state.IncludedLDAPApplication = internaltypes.GetStringSet(r.IncludedLDAPApplication)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a InternalSearchRatePluginResponse object into the model struct
func readInternalSearchRatePluginResponseDataSource(ctx context.Context, r *client.InternalSearchRatePluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("internal-search-rate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.NumThreads = types.Int64Value(r.NumThreads)
	baseDNValues := []string{r.BaseDN}
	state.BaseDN = internaltypes.GetStringSet(baseDNValues)
	state.LowerBound = internaltypes.Int64TypeOrNil(r.LowerBound)
	state.UpperBound = internaltypes.Int64TypeOrNil(r.UpperBound)
	state.FilterPrefix = types.StringValue(r.FilterPrefix)
	state.FilterSuffix = internaltypes.StringTypeOrNil(r.FilterSuffix, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a ModifiablePasswordPolicyStatePluginResponse object into the model struct
func readModifiablePasswordPolicyStatePluginResponseDataSource(ctx context.Context, r *client.ModifiablePasswordPolicyStatePluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("modifiable-password-policy-state")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a SevenBitCleanPluginResponse object into the model struct
func readSevenBitCleanPluginResponseDataSource(ctx context.Context, r *client.SevenBitCleanPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("seven-bit-clean")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.AttributeType = internaltypes.GetStringSet(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a CleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse object into the model struct
func readCleanUpExpiredPingfederatePersistentAccessGrantsPluginResponseDataSource(ctx context.Context, r *client.CleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("clean-up-expired-pingfederate-persistent-access-grants")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PollingInterval = types.StringValue(r.PollingInterval)
	state.PeerServerPriorityIndex = internaltypes.Int64TypeOrNil(r.PeerServerPriorityIndex)
	baseDNValues := []string{}
	baseDNType := internaltypes.StringTypeOrNil(r.BaseDN, false)
	if !baseDNType.IsNull() {
		baseDNValues = append(baseDNValues, baseDNType.ValueString())
	}
	state.BaseDN = internaltypes.GetStringSet(baseDNValues)
	state.MaxUpdatesPerSecond = types.Int64Value(r.MaxUpdatesPerSecond)
	state.NumDeleteThreads = types.Int64Value(r.NumDeleteThreads)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a PeriodicGcPluginResponse object into the model struct
func readPeriodicGcPluginResponseDataSource(ctx context.Context, r *client.PeriodicGcPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("periodic-gc")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.InvokeGCDayOfWeek = internaltypes.GetStringSet(
		client.StringSliceEnumpluginInvokeGCDayOfWeekProp(r.InvokeGCDayOfWeek))
	state.InvokeGCTimeUtc = internaltypes.GetStringSet(r.InvokeGCTimeUtc)
	state.DelayAfterAlert = internaltypes.StringTypeOrNil(r.DelayAfterAlert, false)
	state.DelayPostGC = internaltypes.StringTypeOrNil(r.DelayPostGC, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a PingOnePassThroughAuthenticationPluginResponse object into the model struct
func readPingOnePassThroughAuthenticationPluginResponseDataSource(ctx context.Context, r *client.PingOnePassThroughAuthenticationPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("ping-one-pass-through-authentication")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ApiURL = types.StringValue(r.ApiURL)
	state.AuthURL = types.StringValue(r.AuthURL)
	state.OAuthClientID = types.StringValue(r.OAuthClientID)
	state.OAuthClientSecretPassphraseProvider = internaltypes.StringTypeOrNil(r.OAuthClientSecretPassphraseProvider, false)
	state.EnvironmentID = types.StringValue(r.EnvironmentID)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, false)
	state.IncludedLocalEntryBaseDN = internaltypes.GetStringSet(r.IncludedLocalEntryBaseDN)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.TryLocalBind = internaltypes.BoolTypeOrNil(r.TryLocalBind)
	state.OverrideLocalPassword = internaltypes.BoolTypeOrNil(r.OverrideLocalPassword)
	state.UpdateLocalPassword = internaltypes.BoolTypeOrNil(r.UpdateLocalPassword)
	state.UpdateLocalPasswordDN = internaltypes.StringTypeOrNil(r.UpdateLocalPasswordDN, false)
	state.AllowLaxPassThroughAuthenticationPasswords = internaltypes.BoolTypeOrNil(r.AllowLaxPassThroughAuthenticationPasswords)
	state.IgnoredPasswordPolicyStateErrorCondition = internaltypes.GetStringSet(
		client.StringSliceEnumpluginIgnoredPasswordPolicyStateErrorConditionProp(r.IgnoredPasswordPolicyStateErrorCondition))
	state.UserMappingLocalAttribute = internaltypes.GetStringSet(r.UserMappingLocalAttribute)
	state.UserMappingRemoteJSONField = internaltypes.GetStringSet(r.UserMappingRemoteJSONField)
	state.AdditionalUserMappingSCIMFilter = internaltypes.StringTypeOrNil(r.AdditionalUserMappingSCIMFilter, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a ChangelogPasswordEncryptionPluginResponse object into the model struct
func readChangelogPasswordEncryptionPluginResponseDataSource(ctx context.Context, r *client.ChangelogPasswordEncryptionPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("changelog-password-encryption")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ChangelogPasswordEncryptionKeyPassphraseProvider = internaltypes.StringTypeOrNil(r.ChangelogPasswordEncryptionKeyPassphraseProvider, false)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a ProcessingTimeHistogramPluginResponse object into the model struct
func readProcessingTimeHistogramPluginResponseDataSource(ctx context.Context, r *client.ProcessingTimeHistogramPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("processing-time-histogram")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.HistogramCategoryBoundary = internaltypes.GetStringSet(r.HistogramCategoryBoundary)
	state.IncludeQueueTime = internaltypes.BoolTypeOrNil(r.IncludeQueueTime)
	state.SeparateMonitorEntryPerTrackedApplication = internaltypes.BoolTypeOrNil(r.SeparateMonitorEntryPerTrackedApplication)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a SearchShutdownPluginResponse object into the model struct
func readSearchShutdownPluginResponseDataSource(ctx context.Context, r *client.SearchShutdownPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("search-shutdown")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	baseDNValues := []string{}
	baseDNType := internaltypes.StringTypeOrNil(r.BaseDN, false)
	if !baseDNType.IsNull() {
		baseDNValues = append(baseDNValues, baseDNType.ValueString())
	}
	state.BaseDN = internaltypes.GetStringSet(baseDNValues)
	state.Scope = types.StringValue(r.Scope.String())
	filterValues := []string{r.Filter}
	state.Filter = internaltypes.GetStringSet(filterValues)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.OutputFile = types.StringValue(r.OutputFile)
	state.PreviousFileExtension = internaltypes.StringTypeOrNil(r.PreviousFileExtension, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a PeriodicStatsLoggerPluginResponse object into the model struct
func readPeriodicStatsLoggerPluginResponseDataSource(ctx context.Context, r *client.PeriodicStatsLoggerPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("periodic-stats-logger")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogInterval = types.StringValue(r.LogInterval)
	state.CollectionInterval = types.StringValue(r.CollectionInterval)
	state.SuppressIfIdle = types.BoolValue(r.SuppressIfIdle)
	state.HeaderPrefixPerColumn = internaltypes.BoolTypeOrNil(r.HeaderPrefixPerColumn)
	state.EmptyInsteadOfZero = internaltypes.BoolTypeOrNil(r.EmptyInsteadOfZero)
	state.LinesBetweenHeader = types.Int64Value(r.LinesBetweenHeader)
	state.IncludedLDAPStat = internaltypes.GetStringSet(
		client.StringSliceEnumpluginIncludedLDAPStatProp(r.IncludedLDAPStat))
	state.IncludedResourceStat = internaltypes.GetStringSet(
		client.StringSliceEnumpluginIncludedResourceStatProp(r.IncludedResourceStat))
	state.HistogramFormat = types.StringValue(r.HistogramFormat.String())
	state.HistogramOpType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginHistogramOpTypeProp(r.HistogramOpType))
	state.PerApplicationLDAPStats = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginPeriodicStatsLoggerPerApplicationLDAPStatsProp(r.PerApplicationLDAPStats), false)
	state.StatusSummaryInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginStatusSummaryInfoProp(r.StatusSummaryInfo), false)
	state.LdapChangelogInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLdapChangelogInfoProp(r.LdapChangelogInfo), false)
	state.GaugeInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginGaugeInfoProp(r.GaugeInfo), false)
	state.LogFileFormat = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLogFileFormatProp(r.LogFileFormat), false)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
	state.LocalDBBackendInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLocalDBBackendInfoProp(r.LocalDBBackendInfo), false)
	state.ReplicationInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginReplicationInfoProp(r.ReplicationInfo), false)
	state.EntryCacheInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginEntryCacheInfoProp(r.EntryCacheInfo), false)
	state.HostInfo = internaltypes.GetStringSet(
		client.StringSliceEnumpluginHostInfoProp(r.HostInfo))
	state.IncludedLDAPApplication = internaltypes.GetStringSet(r.IncludedLDAPApplication)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a PurgeExpiredDataPluginResponse object into the model struct
func readPurgeExpiredDataPluginResponseDataSource(ctx context.Context, r *client.PurgeExpiredDataPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("purge-expired-data")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.DatetimeAttribute = types.StringValue(r.DatetimeAttribute)
	state.DatetimeJSONField = internaltypes.StringTypeOrNil(r.DatetimeJSONField, false)
	state.DatetimeFormat = types.StringValue(r.DatetimeFormat.String())
	state.CustomDatetimeFormat = internaltypes.StringTypeOrNil(r.CustomDatetimeFormat, false)
	state.CustomTimezone = internaltypes.StringTypeOrNil(r.CustomTimezone, false)
	state.ExpirationOffset = types.StringValue(r.ExpirationOffset)
	state.PurgeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginPurgeBehaviorProp(r.PurgeBehavior), false)
	baseDNValues := []string{}
	baseDNType := internaltypes.StringTypeOrNil(r.BaseDN, false)
	if !baseDNType.IsNull() {
		baseDNValues = append(baseDNValues, baseDNType.ValueString())
	}
	state.BaseDN = internaltypes.GetStringSet(baseDNValues)
	filterValues := []string{}
	filterType := internaltypes.StringTypeOrNil(r.Filter, false)
	if !filterType.IsNull() {
		filterValues = append(filterValues, filterType.ValueString())
	}
	state.Filter = internaltypes.GetStringSet(filterValues)
	state.PollingInterval = types.StringValue(r.PollingInterval)
	state.MaxUpdatesPerSecond = types.Int64Value(r.MaxUpdatesPerSecond)
	state.PeerServerPriorityIndex = internaltypes.Int64TypeOrNil(r.PeerServerPriorityIndex)
	state.NumDeleteThreads = types.Int64Value(r.NumDeleteThreads)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ChangeSubscriptionNotificationPluginResponse object into the model struct
func readChangeSubscriptionNotificationPluginResponseDataSource(ctx context.Context, r *client.ChangeSubscriptionNotificationPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("change-subscription-notification")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a SubOperationTimingPluginResponse object into the model struct
func readSubOperationTimingPluginResponseDataSource(ctx context.Context, r *client.SubOperationTimingPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("sub-operation-timing")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.NumMostExpensivePhasesShown = types.Int64Value(r.NumMostExpensivePhasesShown)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartyPluginResponse object into the model struct
func readThirdPartyPluginResponseDataSource(ctx context.Context, r *client.ThirdPartyPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a EncryptAttributeValuesPluginResponse object into the model struct
func readEncryptAttributeValuesPluginResponseDataSource(ctx context.Context, r *client.EncryptAttributeValuesPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("encrypt-attribute-values")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.AttributeType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginAttributeTypeProp(r.AttributeType))
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a PassThroughAuthenticationPluginResponse object into the model struct
func readPassThroughAuthenticationPluginResponseDataSource(ctx context.Context, r *client.PassThroughAuthenticationPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("pass-through-authentication")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.Server = internaltypes.GetStringSet(r.Server)
	state.TryLocalBind = types.BoolValue(r.TryLocalBind)
	state.OverrideLocalPassword = types.BoolValue(r.OverrideLocalPassword)
	state.UpdateLocalPassword = types.BoolValue(r.UpdateLocalPassword)
	state.AllowLaxPassThroughAuthenticationPasswords = internaltypes.BoolTypeOrNil(r.AllowLaxPassThroughAuthenticationPasswords)
	state.ServerAccessMode = types.StringValue(r.ServerAccessMode.String())
	state.IncludedLocalEntryBaseDN = internaltypes.GetStringSet(r.IncludedLocalEntryBaseDN)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.DnMap = internaltypes.GetStringSet(r.DnMap)
	state.BindDNPattern = internaltypes.StringTypeOrNil(r.BindDNPattern, false)
	state.SearchBaseDN = internaltypes.StringTypeOrNil(r.SearchBaseDN, false)
	state.SearchFilterPattern = internaltypes.StringTypeOrNil(r.SearchFilterPattern, false)
	state.InitialConnections = types.Int64Value(r.InitialConnections)
	state.MaxConnections = types.Int64Value(r.MaxConnections)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a DnMapperPluginResponse object into the model struct
func readDnMapperPluginResponseDataSource(ctx context.Context, r *client.DnMapperPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("dn-mapper")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.SourceDN = types.StringValue(r.SourceDN)
	state.TargetDN = types.StringValue(r.TargetDN)
	state.EnableAttributeMapping = types.BoolValue(r.EnableAttributeMapping)
	state.MapAttribute = internaltypes.GetStringSet(r.MapAttribute)
	state.EnableControlMapping = types.BoolValue(r.EnableControlMapping)
	state.AlwaysMapResponses = types.BoolValue(r.AlwaysMapResponses)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a MonitorHistoryPluginResponse object into the model struct
func readMonitorHistoryPluginResponseDataSource(ctx context.Context, r *client.MonitorHistoryPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("monitor-history")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogInterval = types.StringValue(r.LogInterval)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.RetainFilesSparselyByAge = internaltypes.BoolTypeOrNil(r.RetainFilesSparselyByAge)
	state.Sanitize = internaltypes.BoolTypeOrNil(r.Sanitize)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ReferralOnUpdatePluginResponse object into the model struct
func readReferralOnUpdatePluginResponseDataSource(ctx context.Context, r *client.ReferralOnUpdatePluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("referral-on-update")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.ReferralBaseURL = internaltypes.GetStringSet(r.ReferralBaseURL)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a SimpleToExternalBindPluginResponse object into the model struct
func readSimpleToExternalBindPluginResponseDataSource(ctx context.Context, r *client.SimpleToExternalBindPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("simple-to-external-bind")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a CustomPluginResponse object into the model struct
func readCustomPluginResponseDataSource(ctx context.Context, r *client.CustomPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("custom")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a SnmpSubagentPluginResponse object into the model struct
func readSnmpSubagentPluginResponseDataSource(ctx context.Context, r *client.SnmpSubagentPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("snmp-subagent")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ContextName = internaltypes.StringTypeOrNil(r.ContextName, false)
	state.AgentxAddress = types.StringValue(r.AgentxAddress)
	state.AgentxPort = types.Int64Value(r.AgentxPort)
	state.NumWorkerThreads = internaltypes.Int64TypeOrNil(r.NumWorkerThreads)
	state.SessionTimeout = internaltypes.StringTypeOrNil(r.SessionTimeout, false)
	state.ConnectRetryMaxWait = internaltypes.StringTypeOrNil(r.ConnectRetryMaxWait, false)
	state.PingInterval = internaltypes.StringTypeOrNil(r.PingInterval, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a CoalesceModificationsPluginResponse object into the model struct
func readCoalesceModificationsPluginResponseDataSource(ctx context.Context, r *client.CoalesceModificationsPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("coalesce-modifications")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.RequestCriteria = types.StringValue(r.RequestCriteria)
	state.AllowedRequestControl = internaltypes.GetStringSet(r.AllowedRequestControl)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a PasswordPolicyImportPluginResponse object into the model struct
func readPasswordPolicyImportPluginResponseDataSource(ctx context.Context, r *client.PasswordPolicyImportPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("password-policy-import")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.DefaultUserPasswordStorageScheme = internaltypes.GetStringSet(r.DefaultUserPasswordStorageScheme)
	state.DefaultAuthPasswordStorageScheme = internaltypes.GetStringSet(r.DefaultAuthPasswordStorageScheme)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ProfilerPluginResponse object into the model struct
func readProfilerPluginResponseDataSource(ctx context.Context, r *client.ProfilerPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("profiler")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ProfileSampleInterval = types.StringValue(r.ProfileSampleInterval)
	state.ProfileDirectory = types.StringValue(r.ProfileDirectory)
	state.EnableProfilingOnStartup = types.BoolValue(r.EnableProfilingOnStartup)
	state.ProfileAction = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginProfileActionProp(r.ProfileAction), false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a CleanUpInactivePingfederatePersistentSessionsPluginResponse object into the model struct
func readCleanUpInactivePingfederatePersistentSessionsPluginResponseDataSource(ctx context.Context, r *client.CleanUpInactivePingfederatePersistentSessionsPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("clean-up-inactive-pingfederate-persistent-sessions")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExpirationOffset = types.StringValue(r.ExpirationOffset)
	state.PollingInterval = types.StringValue(r.PollingInterval)
	state.PeerServerPriorityIndex = internaltypes.Int64TypeOrNil(r.PeerServerPriorityIndex)
	baseDNValues := []string{}
	baseDNType := internaltypes.StringTypeOrNil(r.BaseDN, false)
	if !baseDNType.IsNull() {
		baseDNValues = append(baseDNValues, baseDNType.ValueString())
	}
	state.BaseDN = internaltypes.GetStringSet(baseDNValues)
	state.MaxUpdatesPerSecond = types.Int64Value(r.MaxUpdatesPerSecond)
	state.NumDeleteThreads = types.Int64Value(r.NumDeleteThreads)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ComposedAttributePluginResponse object into the model struct
func readComposedAttributePluginResponseDataSource(ctx context.Context, r *client.ComposedAttributePluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("composed-attribute")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	attributeTypeValues := []string{r.AttributeType}
	state.AttributeType = internaltypes.GetStringSet(attributeTypeValues)
	state.ValuePattern = internaltypes.GetStringSet(r.ValuePattern)
	state.MultipleValuePatternBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginMultipleValuePatternBehaviorProp(r.MultipleValuePatternBehavior), false)
	state.MultiValuedAttributeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginMultiValuedAttributeBehaviorProp(r.MultiValuedAttributeBehavior), false)
	state.TargetAttributeExistsDuringInitialPopulationBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginTargetAttributeExistsDuringInitialPopulationBehaviorProp(r.TargetAttributeExistsDuringInitialPopulationBehavior), false)
	state.UpdateSourceAttributeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginUpdateSourceAttributeBehaviorProp(r.UpdateSourceAttributeBehavior), false)
	state.SourceAttributeRemovalBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginSourceAttributeRemovalBehaviorProp(r.SourceAttributeRemovalBehavior), false)
	state.UpdateTargetAttributeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginUpdateTargetAttributeBehaviorProp(r.UpdateTargetAttributeBehavior), false)
	state.IncludeBaseDN = internaltypes.GetStringSet(r.IncludeBaseDN)
	state.ExcludeBaseDN = internaltypes.GetStringSet(r.ExcludeBaseDN)
	state.IncludeFilter = internaltypes.GetStringSet(r.IncludeFilter)
	state.ExcludeFilter = internaltypes.GetStringSet(r.ExcludeFilter)
	state.UpdatedEntryNewlyMatchesCriteriaBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginUpdatedEntryNewlyMatchesCriteriaBehaviorProp(r.UpdatedEntryNewlyMatchesCriteriaBehavior), false)
	state.UpdatedEntryNoLongerMatchesCriteriaBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginUpdatedEntryNoLongerMatchesCriteriaBehaviorProp(r.UpdatedEntryNoLongerMatchesCriteriaBehavior), false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a LdapResultCodeTrackerPluginResponse object into the model struct
func readLdapResultCodeTrackerPluginResponseDataSource(ctx context.Context, r *client.LdapResultCodeTrackerPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("ldap-result-code-tracker")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a AttributeMapperPluginResponse object into the model struct
func readAttributeMapperPluginResponseDataSource(ctx context.Context, r *client.AttributeMapperPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("attribute-mapper")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.SourceAttribute = types.StringValue(r.SourceAttribute)
	state.TargetAttribute = types.StringValue(r.TargetAttribute)
	state.EnableControlMapping = types.BoolValue(r.EnableControlMapping)
	state.AlwaysMapResponses = types.BoolValue(r.AlwaysMapResponses)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a DelayPluginResponse object into the model struct
func readDelayPluginResponseDataSource(ctx context.Context, r *client.DelayPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("delay")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.Delay = types.StringValue(r.Delay)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a CleanUpExpiredPingfederatePersistentSessionsPluginResponse object into the model struct
func readCleanUpExpiredPingfederatePersistentSessionsPluginResponseDataSource(ctx context.Context, r *client.CleanUpExpiredPingfederatePersistentSessionsPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("clean-up-expired-pingfederate-persistent-sessions")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PollingInterval = types.StringValue(r.PollingInterval)
	state.PeerServerPriorityIndex = internaltypes.Int64TypeOrNil(r.PeerServerPriorityIndex)
	baseDNValues := []string{}
	baseDNType := internaltypes.StringTypeOrNil(r.BaseDN, false)
	if !baseDNType.IsNull() {
		baseDNValues = append(baseDNValues, baseDNType.ValueString())
	}
	state.BaseDN = internaltypes.GetStringSet(baseDNValues)
	state.MaxUpdatesPerSecond = types.Int64Value(r.MaxUpdatesPerSecond)
	state.NumDeleteThreads = types.Int64Value(r.NumDeleteThreads)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a GroovyScriptedPluginResponse object into the model struct
func readGroovyScriptedPluginResponseDataSource(ctx context.Context, r *client.GroovyScriptedPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a LastModPluginResponse object into the model struct
func readLastModPluginResponseDataSource(ctx context.Context, r *client.LastModPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("last-mod")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.ExcludeAttribute = internaltypes.GetStringSet(r.ExcludeAttribute)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a PluggablePassThroughAuthenticationPluginResponse object into the model struct
func readPluggablePassThroughAuthenticationPluginResponseDataSource(ctx context.Context, r *client.PluggablePassThroughAuthenticationPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("pluggable-pass-through-authentication")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PassThroughAuthenticationHandler = types.StringValue(r.PassThroughAuthenticationHandler)
	state.IncludedLocalEntryBaseDN = internaltypes.GetStringSet(r.IncludedLocalEntryBaseDN)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.TryLocalBind = internaltypes.BoolTypeOrNil(r.TryLocalBind)
	state.OverrideLocalPassword = internaltypes.BoolTypeOrNil(r.OverrideLocalPassword)
	state.UpdateLocalPassword = internaltypes.BoolTypeOrNil(r.UpdateLocalPassword)
	state.UpdateLocalPasswordDN = internaltypes.StringTypeOrNil(r.UpdateLocalPasswordDN, false)
	state.AllowLaxPassThroughAuthenticationPasswords = internaltypes.BoolTypeOrNil(r.AllowLaxPassThroughAuthenticationPasswords)
	state.IgnoredPasswordPolicyStateErrorCondition = internaltypes.GetStringSet(
		client.StringSliceEnumpluginIgnoredPasswordPolicyStateErrorConditionProp(r.IgnoredPasswordPolicyStateErrorCondition))
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a ReferentialIntegrityPluginResponse object into the model struct
func readReferentialIntegrityPluginResponseDataSource(ctx context.Context, r *client.ReferentialIntegrityPluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("referential-integrity")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.AttributeType = internaltypes.GetStringSet(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.LogFile = internaltypes.StringTypeOrNil(r.LogFile, false)
	state.UpdateInterval = internaltypes.StringTypeOrNil(r.UpdateInterval, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read a UniqueAttributePluginResponse object into the model struct
func readUniqueAttributePluginResponseDataSource(ctx context.Context, r *client.UniqueAttributePluginResponse, state *pluginDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("unique-attribute")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.Type = internaltypes.GetStringSet(r.Type)
	state.MultipleAttributeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginUniqueAttributeMultipleAttributeBehaviorProp(r.MultipleAttributeBehavior), false)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.PreventConflictsWithSoftDeletedEntries = internaltypes.BoolTypeOrNil(r.PreventConflictsWithSoftDeletedEntries)
	filterValues := []string{}
	filterType := internaltypes.StringTypeOrNil(r.Filter, false)
	if !filterType.IsNull() {
		filterValues = append(filterValues, filterType.ValueString())
	}
	state.Filter = internaltypes.GetStringSet(filterValues)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
}

// Read resource information
func (r *pluginDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state pluginDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.LastAccessTimePluginResponse != nil {
		readLastAccessTimePluginResponseDataSource(ctx, readResponse.LastAccessTimePluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.StatsCollectorPluginResponse != nil {
		readStatsCollectorPluginResponseDataSource(ctx, readResponse.StatsCollectorPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.InternalSearchRatePluginResponse != nil {
		readInternalSearchRatePluginResponseDataSource(ctx, readResponse.InternalSearchRatePluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ModifiablePasswordPolicyStatePluginResponse != nil {
		readModifiablePasswordPolicyStatePluginResponseDataSource(ctx, readResponse.ModifiablePasswordPolicyStatePluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SevenBitCleanPluginResponse != nil {
		readSevenBitCleanPluginResponseDataSource(ctx, readResponse.SevenBitCleanPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse != nil {
		readCleanUpExpiredPingfederatePersistentAccessGrantsPluginResponseDataSource(ctx, readResponse.CleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PeriodicGcPluginResponse != nil {
		readPeriodicGcPluginResponseDataSource(ctx, readResponse.PeriodicGcPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PingOnePassThroughAuthenticationPluginResponse != nil {
		readPingOnePassThroughAuthenticationPluginResponseDataSource(ctx, readResponse.PingOnePassThroughAuthenticationPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ChangelogPasswordEncryptionPluginResponse != nil {
		readChangelogPasswordEncryptionPluginResponseDataSource(ctx, readResponse.ChangelogPasswordEncryptionPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ProcessingTimeHistogramPluginResponse != nil {
		readProcessingTimeHistogramPluginResponseDataSource(ctx, readResponse.ProcessingTimeHistogramPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SearchShutdownPluginResponse != nil {
		readSearchShutdownPluginResponseDataSource(ctx, readResponse.SearchShutdownPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PeriodicStatsLoggerPluginResponse != nil {
		readPeriodicStatsLoggerPluginResponseDataSource(ctx, readResponse.PeriodicStatsLoggerPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PurgeExpiredDataPluginResponse != nil {
		readPurgeExpiredDataPluginResponseDataSource(ctx, readResponse.PurgeExpiredDataPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ChangeSubscriptionNotificationPluginResponse != nil {
		readChangeSubscriptionNotificationPluginResponseDataSource(ctx, readResponse.ChangeSubscriptionNotificationPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SubOperationTimingPluginResponse != nil {
		readSubOperationTimingPluginResponseDataSource(ctx, readResponse.SubOperationTimingPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPluginResponse != nil {
		readThirdPartyPluginResponseDataSource(ctx, readResponse.ThirdPartyPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.EncryptAttributeValuesPluginResponse != nil {
		readEncryptAttributeValuesPluginResponseDataSource(ctx, readResponse.EncryptAttributeValuesPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PassThroughAuthenticationPluginResponse != nil {
		readPassThroughAuthenticationPluginResponseDataSource(ctx, readResponse.PassThroughAuthenticationPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DnMapperPluginResponse != nil {
		readDnMapperPluginResponseDataSource(ctx, readResponse.DnMapperPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.MonitorHistoryPluginResponse != nil {
		readMonitorHistoryPluginResponseDataSource(ctx, readResponse.MonitorHistoryPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ReferralOnUpdatePluginResponse != nil {
		readReferralOnUpdatePluginResponseDataSource(ctx, readResponse.ReferralOnUpdatePluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SimpleToExternalBindPluginResponse != nil {
		readSimpleToExternalBindPluginResponseDataSource(ctx, readResponse.SimpleToExternalBindPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CustomPluginResponse != nil {
		readCustomPluginResponseDataSource(ctx, readResponse.CustomPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SnmpSubagentPluginResponse != nil {
		readSnmpSubagentPluginResponseDataSource(ctx, readResponse.SnmpSubagentPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CoalesceModificationsPluginResponse != nil {
		readCoalesceModificationsPluginResponseDataSource(ctx, readResponse.CoalesceModificationsPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PasswordPolicyImportPluginResponse != nil {
		readPasswordPolicyImportPluginResponseDataSource(ctx, readResponse.PasswordPolicyImportPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ProfilerPluginResponse != nil {
		readProfilerPluginResponseDataSource(ctx, readResponse.ProfilerPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CleanUpInactivePingfederatePersistentSessionsPluginResponse != nil {
		readCleanUpInactivePingfederatePersistentSessionsPluginResponseDataSource(ctx, readResponse.CleanUpInactivePingfederatePersistentSessionsPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ComposedAttributePluginResponse != nil {
		readComposedAttributePluginResponseDataSource(ctx, readResponse.ComposedAttributePluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LdapResultCodeTrackerPluginResponse != nil {
		readLdapResultCodeTrackerPluginResponseDataSource(ctx, readResponse.LdapResultCodeTrackerPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AttributeMapperPluginResponse != nil {
		readAttributeMapperPluginResponseDataSource(ctx, readResponse.AttributeMapperPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DelayPluginResponse != nil {
		readDelayPluginResponseDataSource(ctx, readResponse.DelayPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CleanUpExpiredPingfederatePersistentSessionsPluginResponse != nil {
		readCleanUpExpiredPingfederatePersistentSessionsPluginResponseDataSource(ctx, readResponse.CleanUpExpiredPingfederatePersistentSessionsPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedPluginResponse != nil {
		readGroovyScriptedPluginResponseDataSource(ctx, readResponse.GroovyScriptedPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LastModPluginResponse != nil {
		readLastModPluginResponseDataSource(ctx, readResponse.LastModPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PluggablePassThroughAuthenticationPluginResponse != nil {
		readPluggablePassThroughAuthenticationPluginResponseDataSource(ctx, readResponse.PluggablePassThroughAuthenticationPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ReferentialIntegrityPluginResponse != nil {
		readReferentialIntegrityPluginResponseDataSource(ctx, readResponse.ReferentialIntegrityPluginResponse, &state, &resp.Diagnostics)
	}
	if readResponse.UniqueAttributePluginResponse != nil {
		readUniqueAttributePluginResponseDataSource(ctx, readResponse.UniqueAttributePluginResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
