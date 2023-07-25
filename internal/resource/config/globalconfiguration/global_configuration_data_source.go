package globalconfiguration

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
	_ datasource.DataSource              = &globalConfigurationDataSource{}
	_ datasource.DataSourceWithConfigure = &globalConfigurationDataSource{}
)

// Create a Global Configuration data source
func NewGlobalConfigurationDataSource() datasource.DataSource {
	return &globalConfigurationDataSource{}
}

// globalConfigurationDataSource is the datasource implementation.
type globalConfigurationDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *globalConfigurationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_global_configuration"
}

// Configure adds the provider configured client to the data source.
func (r *globalConfigurationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type globalConfigurationDataSourceModel struct {
	// Id field required for acceptance testing framework
	Id                                                             types.String `tfsdk:"id"`
	InstanceName                                                   types.String `tfsdk:"instance_name"`
	Location                                                       types.String `tfsdk:"location"`
	ConfigurationServerGroup                                       types.String `tfsdk:"configuration_server_group"`
	ForceAsMasterForMirroredData                                   types.Bool   `tfsdk:"force_as_master_for_mirrored_data"`
	EncryptData                                                    types.Bool   `tfsdk:"encrypt_data"`
	EncryptionSettingsCipherStreamProvider                         types.String `tfsdk:"encryption_settings_cipher_stream_provider"`
	EncryptBackupsByDefault                                        types.Bool   `tfsdk:"encrypt_backups_by_default"`
	BackupEncryptionSettingsDefinitionID                           types.String `tfsdk:"backup_encryption_settings_definition_id"`
	EncryptLDIFExportsByDefault                                    types.Bool   `tfsdk:"encrypt_ldif_exports_by_default"`
	LdifExportEncryptionSettingsDefinitionID                       types.String `tfsdk:"ldif_export_encryption_settings_definition_id"`
	AutomaticallyCompressEncryptedLDIFExports                      types.Bool   `tfsdk:"automatically_compress_encrypted_ldif_exports"`
	RedactSensitiveValuesInConfigLogs                              types.Bool   `tfsdk:"redact_sensitive_values_in_config_logs"`
	SensitiveAttribute                                             types.Set    `tfsdk:"sensitive_attribute"`
	RejectInsecureRequests                                         types.Bool   `tfsdk:"reject_insecure_requests"`
	AllowedInsecureRequestCriteria                                 types.String `tfsdk:"allowed_insecure_request_criteria"`
	RejectUnauthenticatedRequests                                  types.Bool   `tfsdk:"reject_unauthenticated_requests"`
	AllowedUnauthenticatedRequestCriteria                          types.String `tfsdk:"allowed_unauthenticated_request_criteria"`
	BindWithDNRequiresPassword                                     types.Bool   `tfsdk:"bind_with_dn_requires_password"`
	DisabledPrivilege                                              types.Set    `tfsdk:"disabled_privilege"`
	DefaultPasswordPolicy                                          types.String `tfsdk:"default_password_policy"`
	MaximumUserDataPasswordPoliciesToCache                         types.Int64  `tfsdk:"maximum_user_data_password_policies_to_cache"`
	ProxiedAuthorizationIdentityMapper                             types.String `tfsdk:"proxied_authorization_identity_mapper"`
	VerifyEntryDigests                                             types.Bool   `tfsdk:"verify_entry_digests"`
	AllowedInsecureTLSProtocol                                     types.Set    `tfsdk:"allowed_insecure_tls_protocol"`
	AllowInsecureLocalJMXConnections                               types.Bool   `tfsdk:"allow_insecure_local_jmx_connections"`
	DefaultInternalOperationClientConnectionPolicy                 types.String `tfsdk:"default_internal_operation_client_connection_policy"`
	SizeLimit                                                      types.Int64  `tfsdk:"size_limit"`
	UnauthenticatedSizeLimit                                       types.Int64  `tfsdk:"unauthenticated_size_limit"`
	TimeLimit                                                      types.String `tfsdk:"time_limit"`
	UnauthenticatedTimeLimit                                       types.String `tfsdk:"unauthenticated_time_limit"`
	IdleTimeLimit                                                  types.String `tfsdk:"idle_time_limit"`
	UnauthenticatedIdleTimeLimit                                   types.String `tfsdk:"unauthenticated_idle_time_limit"`
	LookthroughLimit                                               types.Int64  `tfsdk:"lookthrough_limit"`
	UnauthenticatedLookthroughLimit                                types.Int64  `tfsdk:"unauthenticated_lookthrough_limit"`
	LdapJoinSizeLimit                                              types.Int64  `tfsdk:"ldap_join_size_limit"`
	MaximumConcurrentConnections                                   types.Int64  `tfsdk:"maximum_concurrent_connections"`
	MaximumConcurrentConnectionsPerIPAddress                       types.Int64  `tfsdk:"maximum_concurrent_connections_per_ip_address"`
	MaximumConcurrentConnectionsPerBindDN                          types.Int64  `tfsdk:"maximum_concurrent_connections_per_bind_dn"`
	MaximumConcurrentUnindexedSearches                             types.Int64  `tfsdk:"maximum_concurrent_unindexed_searches"`
	MaximumAttributesPerAddRequest                                 types.Int64  `tfsdk:"maximum_attributes_per_add_request"`
	MaximumModificationsPerModifyRequest                           types.Int64  `tfsdk:"maximum_modifications_per_modify_request"`
	BackgroundThreadForEachPersistentSearch                        types.Bool   `tfsdk:"background_thread_for_each_persistent_search"`
	AllowAttributeNameExceptions                                   types.Bool   `tfsdk:"allow_attribute_name_exceptions"`
	InvalidAttributeSyntaxBehavior                                 types.String `tfsdk:"invalid_attribute_syntax_behavior"`
	PermitSyntaxViolationsForAttribute                             types.Set    `tfsdk:"permit_syntax_violations_for_attribute"`
	SingleStructuralObjectclassBehavior                            types.String `tfsdk:"single_structural_objectclass_behavior"`
	AttributesModifiableWithIgnoreNoUserModificationRequestControl types.Set    `tfsdk:"attributes_modifiable_with_ignore_no_user_modification_request_control"`
	MaximumServerOutLogFileSize                                    types.String `tfsdk:"maximum_server_out_log_file_size"`
	MaximumServerOutLogFileCount                                   types.Int64  `tfsdk:"maximum_server_out_log_file_count"`
	StartupErrorLoggerOutputLocation                               types.String `tfsdk:"startup_error_logger_output_location"`
	ExitOnJVMError                                                 types.Bool   `tfsdk:"exit_on_jvm_error"`
	ServerErrorResultCode                                          types.Int64  `tfsdk:"server_error_result_code"`
	ResultCodeMap                                                  types.String `tfsdk:"result_code_map"`
	ReturnBindErrorMessages                                        types.Bool   `tfsdk:"return_bind_error_messages"`
	NotifyAbandonedOperations                                      types.Bool   `tfsdk:"notify_abandoned_operations"`
	DuplicateErrorLogLimit                                         types.Int64  `tfsdk:"duplicate_error_log_limit"`
	DuplicateErrorLogTimeLimit                                     types.String `tfsdk:"duplicate_error_log_time_limit"`
	DuplicateAlertLimit                                            types.Int64  `tfsdk:"duplicate_alert_limit"`
	DuplicateAlertTimeLimit                                        types.String `tfsdk:"duplicate_alert_time_limit"`
	WritabilityMode                                                types.String `tfsdk:"writability_mode"`
	UnrecoverableDatabaseErrorMode                                 types.String `tfsdk:"unrecoverable_database_error_mode"`
	DatabaseOnVirtualizedOrNetworkStorage                          types.Bool   `tfsdk:"database_on_virtualized_or_network_storage"`
	AutoNameWithEntryUUIDConnectionCriteria                        types.String `tfsdk:"auto_name_with_entry_uuid_connection_criteria"`
	AutoNameWithEntryUUIDRequestCriteria                           types.String `tfsdk:"auto_name_with_entry_uuid_request_criteria"`
	SoftDeletePolicy                                               types.String `tfsdk:"soft_delete_policy"`
	SubtreeAccessibilityAlertTimeLimit                             types.String `tfsdk:"subtree_accessibility_alert_time_limit"`
	WarnForBackendsWithMultipleBaseDns                             types.Bool   `tfsdk:"warn_for_backends_with_multiple_base_dns"`
	ForcedGCPrimeDuration                                          types.String `tfsdk:"forced_gc_prime_duration"`
	ReplicationSetName                                             types.String `tfsdk:"replication_set_name"`
	StartupMinReplicationBacklogCount                              types.Int64  `tfsdk:"startup_min_replication_backlog_count"`
	ReplicationBacklogCountAlertThreshold                          types.Int64  `tfsdk:"replication_backlog_count_alert_threshold"`
	ReplicationBacklogDurationAlertThreshold                       types.String `tfsdk:"replication_backlog_duration_alert_threshold"`
	ReplicationAssuranceSourceTimeoutSuspendDuration               types.String `tfsdk:"replication_assurance_source_timeout_suspend_duration"`
	ReplicationAssuranceSourceBacklogFastStartThreshold            types.Int64  `tfsdk:"replication_assurance_source_backlog_fast_start_threshold"`
	ReplicationHistoryLimit                                        types.Int64  `tfsdk:"replication_history_limit"`
	AllowInheritedReplicationOfSubordinateBackends                 types.Bool   `tfsdk:"allow_inherited_replication_of_subordinate_backends"`
	ReplicationPurgeObsoleteReplicas                               types.Bool   `tfsdk:"replication_purge_obsolete_replicas"`
	SmtpServer                                                     types.Set    `tfsdk:"smtp_server"`
	MaxSMTPConnectionCount                                         types.Int64  `tfsdk:"max_smtp_connection_count"`
	MaxSMTPConnectionAge                                           types.String `tfsdk:"max_smtp_connection_age"`
	SmtpConnectionHealthCheckInterval                              types.String `tfsdk:"smtp_connection_health_check_interval"`
	AllowedTask                                                    types.Set    `tfsdk:"allowed_task"`
	EnableSubOperationTimer                                        types.Bool   `tfsdk:"enable_sub_operation_timer"`
	MaximumShutdownTime                                            types.String `tfsdk:"maximum_shutdown_time"`
	NetworkAddressCacheTTL                                         types.String `tfsdk:"network_address_cache_ttl"`
	NetworkAddressOutageCacheEnabled                               types.Bool   `tfsdk:"network_address_outage_cache_enabled"`
	TrackedApplication                                             types.Set    `tfsdk:"tracked_application"`
	JmxValueBehavior                                               types.String `tfsdk:"jmx_value_behavior"`
	JmxUseLegacyMbeanNames                                         types.Bool   `tfsdk:"jmx_use_legacy_mbean_names"`
}

// GetSchema defines the schema for the datasource.
func (r *globalConfigurationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Describes a Global Configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Name of this object.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"instance_name": schema.StringAttribute{
				Description: "Specifies a name that may be used to uniquely identify this Directory Server instance among other instances in the environment.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"location": schema.StringAttribute{
				Description: "Specifies the location for this Directory Server. Operations performed which involve communication with other servers may prefer servers in the same location to help ensure low-latency responses.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"configuration_server_group": schema.StringAttribute{
				Description: "When this property is set, changes made to this server using the console or dsconfig can be automatically applied to all servers in the specified server group.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"force_as_master_for_mirrored_data": schema.BoolAttribute{
				Description: "Indicates whether this server should be forced to assume the master role if no other suitable server is found to act as master or if multiple masters are detected. A master is only needed when changes are made to mirrored data, i.e. data specific to the topology itself and cluster-wide configuration data.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"encrypt_data": schema.BoolAttribute{
				Description: "Indicates whether the Directory Server should encrypt the data that it stores in all components that support it. This may include certain types of backends (including local DB and large attribute backends), the LDAP changelog, and the replication server database.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"encryption_settings_cipher_stream_provider": schema.StringAttribute{
				Description: "Specifies the cipher stream provider that should be used to protect the contents of the encryption settings database.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"encrypt_backups_by_default": schema.BoolAttribute{
				Description: "Indicates whether the server should encrypt backups by default.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"backup_encryption_settings_definition_id": schema.StringAttribute{
				Description: "The unique identifier for the encryption settings definition to use to generate the encryption key for encrypted backups by default.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"encrypt_ldif_exports_by_default": schema.BoolAttribute{
				Description: "Indicates whether the server should encrypt LDIF exports by default.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ldif_export_encryption_settings_definition_id": schema.StringAttribute{
				Description: "The unique identifier for the encryption settings definition to use to generate the encryption key for encrypted LDIF exports by default.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"automatically_compress_encrypted_ldif_exports": schema.BoolAttribute{
				Description: "Indicates whether to automatically compress LDIF exports that are also encrypted.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"redact_sensitive_values_in_config_logs": schema.BoolAttribute{
				Description: "Indicates whether the values of sensitive configuration properties should be redacted when logging configuration changes, including in the configuration audit log, the error log, and the server.out log file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"sensitive_attribute": schema.SetAttribute{
				Description: "Provides the ability to indicate that some attributes should be considered sensitive and additional protection should be in place when interacting with those attributes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"reject_insecure_requests": schema.BoolAttribute{
				Description: "Indicates whether the Directory Server should reject any LDAP request (other than StartTLS) received from a client that is not using an encrypted connection.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allowed_insecure_request_criteria": schema.StringAttribute{
				Description: "A set of criteria that may be used to match LDAP requests that may be permitted over an insecure connection even if reject-insecure-requests is true. Note that some types of requests will always be permitted, including StartTLS and start administrative session requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"reject_unauthenticated_requests": schema.BoolAttribute{
				Description: "Indicates whether the Directory Server should reject any LDAP request (other than bind or StartTLS requests) received from a client that has not yet been authenticated, whose last authentication attempt was unsuccessful, or whose last authentication attempt used anonymous authentication.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allowed_unauthenticated_request_criteria": schema.StringAttribute{
				Description: "A set of criteria that may be used to match LDAP requests that may be permitted over an unauthenticated connection even if reject-unauthenticated-requests is true. Note that some types of requests will always be permitted, including bind, StartTLS, and start administrative session requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"bind_with_dn_requires_password": schema.BoolAttribute{
				Description: "Indicates whether the Directory Server should reject any simple bind request that contains a DN but no password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"disabled_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that should not be evaluated by the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"default_password_policy": schema.StringAttribute{
				Description: "Specifies the name of the password policy that is in effect for users whose entries do not specify an alternate password policy (either via a real or virtual attribute).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_user_data_password_policies_to_cache": schema.Int64Attribute{
				Description: "Specifies the maximum number of password policies that are defined in the user data (that is, outside of the configuration) that the server should cache in memory for faster access. A value of zero indicates that the server should not cache any user data password policies.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"proxied_authorization_identity_mapper": schema.StringAttribute{
				Description: "Specifies the name of the identity mapper to map authorization ID values (using the \"u:\" form) provided in the proxied authorization control to the corresponding user entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"verify_entry_digests": schema.BoolAttribute{
				Description: "Indicates whether the digest should always be verified whenever an entry containing a digest is decoded. If this is \"true\", then if a digest exists, it will always be verified. Otherwise, the digest will be written when encoding entries but ignored when decoding entries but may still be available for other verification processing.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allowed_insecure_tls_protocol": schema.SetAttribute{
				Description: "Specifies a set of TLS protocols that will be permitted for use in the server even though there may be known vulnerabilities that could cause their use to be unsafe in some conditions. Enabling support for insecure TLS protocols is discouraged, and is generally recommended only as a short-term measure to permit legacy clients to interact with the server until they can be updated to support more secure communication protocols.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"allow_insecure_local_jmx_connections": schema.BoolAttribute{
				Description: "Indicates that processes attaching to this server's local JVM are allowed to access internal data through JMX without the authentication requirements that remote JMX connections are subject to. Please review and understand the data that this option will expose (such as cn=monitor) to client applications to ensure there are no security concerns.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"default_internal_operation_client_connection_policy": schema.StringAttribute{
				Description: "Specifies the client connection policy that will be used by default for internal operations.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"size_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that the Directory Server should return to clients by default when processing a search operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"unauthenticated_size_limit": schema.Int64Attribute{
				Description: "The size limit value that will apply for connections from unauthenticated clients. If this is not specified, then the value of the size-limit property will be applied for both authenticated and unauthenticated connections. Supported in PingDirectory product version 9.2.0.0+.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"time_limit": schema.StringAttribute{
				Description: "Specifies the maximum length of time that the Directory Server should be allowed to spend processing a search operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"unauthenticated_time_limit": schema.StringAttribute{
				Description: "The time limit value that will apply for connections from unauthenticated clients. If this is not specified, then the value of the time-limit property will be applied for both authenticated and unauthenticated connections. Supported in PingDirectory product version 9.2.0.0+.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"idle_time_limit": schema.StringAttribute{
				Description: "Specifies the maximum length of time that a client connection may remain established since its last completed operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"unauthenticated_idle_time_limit": schema.StringAttribute{
				Description: "The idle-time-limit limit value that will apply for connections from unauthenticated clients. If this is not specified, then the value of the idle-time-limit property will be applied for both authenticated and unauthenticated connections. Supported in PingDirectory product version 9.2.0.0+.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"lookthrough_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that the Directory Server should \"look through\" in the course of processing a search request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"unauthenticated_lookthrough_limit": schema.Int64Attribute{
				Description: "The lookthrough limit value that will apply for connections from unauthenticated clients. If this is not specified, then the value of the lookthrough-limit property will be applied for both authenticated and unauthenticated connections. Supported in PingDirectory product version 9.2.0.0+.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ldap_join_size_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that may be directly joined with any individual search result entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_concurrent_connections": schema.Int64Attribute{
				Description: "Specifies the maximum number of LDAP client connections which may be established to this Directory Server at the same time.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_concurrent_connections_per_ip_address": schema.Int64Attribute{
				Description: "Specifies the maximum number of LDAP client connections originating from the same IP address which may be established to this Directory Server at the same time.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_concurrent_connections_per_bind_dn": schema.Int64Attribute{
				Description: "Specifies the maximum number of LDAP client connections which may be established to this Directory Server at the same time and authenticated as the same user.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_concurrent_unindexed_searches": schema.Int64Attribute{
				Description: "Specifies the maximum number of unindexed searches that may be in progress in this backend at any given time. Any unindexed searches requested while the maximum number of unindexed searches are already being processed will be rejected. A value of zero indicates that no limit will be enforced.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_attributes_per_add_request": schema.Int64Attribute{
				Description: "Specifies the maximum number of attributes that may be included in an add request. This property does not impose any limit on the number of values that an attribute may have.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_modifications_per_modify_request": schema.Int64Attribute{
				Description: "Specifies the maximum number of modifications that may be included in a modify request. This property does not impose any limit on the number of attribute values that a modification may have.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"background_thread_for_each_persistent_search": schema.BoolAttribute{
				Description: "Indicates whether the server should use a separate background thread for each persistent search.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_attribute_name_exceptions": schema.BoolAttribute{
				Description: "Indicates whether the Directory Server should allow underscores in attribute names and allow attribute names to begin with numeric digits (both of which are violations of the LDAP standards).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"invalid_attribute_syntax_behavior": schema.StringAttribute{
				Description: "Specifies how the Directory Server should handle operations whenever an attribute value violates the associated attribute syntax.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"permit_syntax_violations_for_attribute": schema.SetAttribute{
				Description: "Specifies a set of attribute types for which the server will permit values that do not conform to the associated attribute syntax.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"single_structural_objectclass_behavior": schema.StringAttribute{
				Description: "Specifies how the Directory Server should handle operations for an entry does not contain a structural object class, or for an entry that contains multiple structural classes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"attributes_modifiable_with_ignore_no_user_modification_request_control": schema.SetAttribute{
				Description: "Specifies the operational attribute types that are defined in the schema with the NO-USER-MODIFICATION constraint that the server will allow to be altered if the associated request contains the ignore NO-USER-MODIFICATION request control.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"maximum_server_out_log_file_size": schema.StringAttribute{
				Description: "The maximum allowed size that the server.out log file will be allowed to have. If a write would cause the file to exceed this size, then the current file will be rotated out of place and a new empty file will be created and the message written to it.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_server_out_log_file_count": schema.Int64Attribute{
				Description: "The maximum number of server.out log files (including the current active log file) that should be retained. When rotating the log file, if the total number of files exceeds this count, then the oldest file(s) will be removed so that the total number of log files is within this limit.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"startup_error_logger_output_location": schema.StringAttribute{
				Description: "Specifies how the server should handle error log messages (which may include errors, warnings, and notices) generated during startup. All of these messages will be written to all configured error loggers, but they may also be written to other locations (like standard output, standard error, or the server.out log file) so that they are displayed on the console when the server is starting.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"exit_on_jvm_error": schema.BoolAttribute{
				Description: "Indicates whether the Directory Server should be shut down if a severe error is raised (e.g., an out of memory error) which may prevent the JVM from continuing to run properly.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_error_result_code": schema.Int64Attribute{
				Description: "Specifies the numeric value of the result code when request processing fails due to an internal server error.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"result_code_map": schema.StringAttribute{
				Description: "Specifies a result code map that should be used for clients that do not have a map associated with their client connection policy. If the associated client connection policy has a result code map, then that map will be used instead. If no map is associated either with the client connection policy or the global configuration, then an internal default will be used.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"return_bind_error_messages": schema.BoolAttribute{
				Description: "Indicates whether responses for failed bind operations should include a message string providing the reason for the authentication failure.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"notify_abandoned_operations": schema.BoolAttribute{
				Description: "Indicates whether the Directory Server should send a response to any operation that is interrupted via an abandon request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"duplicate_error_log_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of duplicate error log messages that should be logged in the time window specified by the duplicate-error-log-time-limit property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"duplicate_error_log_time_limit": schema.StringAttribute{
				Description: "Specifies the length of time that must expire before duplicate log messages above the duplicate-error-log-limit threshold are logged again to the error log.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"duplicate_alert_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of duplicate alert messages that should be sent via the administrative alert framework in the time window specified by the duplicate-alert-time-limit property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"duplicate_alert_time_limit": schema.StringAttribute{
				Description: "Specifies the length of time that must expire before duplicate messages are sent via the administrative alert framework.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"writability_mode": schema.StringAttribute{
				Description: "Specifies the kinds of write operations the Directory Server can process.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"unrecoverable_database_error_mode": schema.StringAttribute{
				Description: "Specifies the action which should be taken for any database that experiences an unrecoverable error. Action applies to local database backends and the replication recent changes database.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"database_on_virtualized_or_network_storage": schema.BoolAttribute{
				Description: "This setting provides data integrity options when the Directory Server is installed with a database on a network storage device. A storage device may be accessed directly by a physical server, or indirectly through a virtual machine running on a hypervisor. Enabling this setting will apply changes to all Local DB Backends, the LDAP Changelog Backend, and the replication changelog database.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"auto_name_with_entry_uuid_connection_criteria": schema.StringAttribute{
				Description: "Connection criteria that may be used to identify clients whose add requests should use entryUUID as the naming attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"auto_name_with_entry_uuid_request_criteria": schema.StringAttribute{
				Description: "Request criteria that may be used to identify add requests that should use entryUUID as the naming attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"soft_delete_policy": schema.StringAttribute{
				Description: "Specifies the soft delete policy that will be used by default for delete operations. Soft delete operations introduce the ability to control the server behavior of the delete operation. Instead of performing a permanent delete of an entry, deleted entries can be retained as soft deleted entries by their entryUUID values and are available for undelete at a later time. In addition to a soft delete policy enabling soft deletes, delete operations sent to the server must have the soft delete request control present with sufficient access privileges to access the soft delete request control.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"subtree_accessibility_alert_time_limit": schema.StringAttribute{
				Description: "Specifies the length of time that a subtree may remain hidden or read-only before an administrative alert is sent.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"warn_for_backends_with_multiple_base_dns": schema.BoolAttribute{
				Description: "Indicates whether the server should issue a warning when enabling a backend that contains multiple base DNs.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"forced_gc_prime_duration": schema.StringAttribute{
				Description: "Specifies the minimum length of time required for backend or request processor initialization that will trigger the server to force an explicit garbage collection. A value of \"0 seconds\" indicates that the server should never invoke an explicit garbage collection regardless of the length of time required to initialize the server backends.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"replication_set_name": schema.StringAttribute{
				Description: "The name of the replication set assigned to this Directory Server. Restricted domains are only replicated within instances using the same replication set name.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"startup_min_replication_backlog_count": schema.Int64Attribute{
				Description: "The number of outstanding changes any replica can have before the Directory Server will start accepting connections. The Directory Server may never accept connections if this setting is too low. If you are unsure which value to use, you can use the number of expected updates within a five second interval.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"replication_backlog_count_alert_threshold": schema.Int64Attribute{
				Description: "An alert is sent when the number of outstanding replication changes for the Directory Server has exceeded this threshold for longer than the replication backlog duration alert threshold.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"replication_backlog_duration_alert_threshold": schema.StringAttribute{
				Description: "An alert is sent when the number of outstanding replication changes for the Directory Server has exceeded the replication backlog count alert threshold for longer than this duration.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"replication_assurance_source_timeout_suspend_duration": schema.StringAttribute{
				Description: "The amount of time a replication assurance source (i.e. a peer Directory Server) will be suspended from assurance requirements on this Directory Server if it experiences an assurance timeout.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"replication_assurance_source_backlog_fast_start_threshold": schema.Int64Attribute{
				Description: "The maximum number of replication backlog updates a replication assurance source (i.e. a peer Directory Server) can have and be immediately recognized as an available assurance source by this Directory Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"replication_history_limit": schema.Int64Attribute{
				Description: "Specifies the size limit for historical information.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_inherited_replication_of_subordinate_backends": schema.BoolAttribute{
				Description: "Allow replication to be inherited by subordinate/child backends.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"replication_purge_obsolete_replicas": schema.BoolAttribute{
				Description: "Indicates whether state about obsolete replicas is automatically purged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"smtp_server": schema.SetAttribute{
				Description: "Specifies the set of servers that will be used to send email messages. The order in which the servers are listed indicates the order in which the Directory Server will attempt to use them in the course of sending a message. The first attempt will always go to the server at the top of the list, and servers further down the list will only be used if none of the servers listed above it were able to successfully send the message.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"max_smtp_connection_count": schema.Int64Attribute{
				Description: "The maximum number of SMTP connections that will be maintained for delivering email messages.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_smtp_connection_age": schema.StringAttribute{
				Description: "The maximum length of time that a connection to an SMTP server should be considered valid.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"smtp_connection_health_check_interval": schema.StringAttribute{
				Description: "The length of time between checks to ensure that available SMTP connections are still valid.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allowed_task": schema.SetAttribute{
				Description: "Specifies the fully-qualified name of a Java class that may be invoked in the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"enable_sub_operation_timer": schema.BoolAttribute{
				Description: "Indicates whether the Directory Server should attempt to record information about the length of time required to process various phases of an operation. Enabling this feature may impact performance, but could make it easier to identify potential bottlenecks in operation processing.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_shutdown_time": schema.StringAttribute{
				Description: "Specifies the maximum amount of time the shutdown of Directory Server may take.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"network_address_cache_ttl": schema.StringAttribute{
				Description: "Specifies the length of time that the Directory Server should cache the IP addresses associated with the names of systems with which it interacts.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"network_address_outage_cache_enabled": schema.BoolAttribute{
				Description: "Specifies whether the Directory Server should cache the last valid IP addresses associated with the names of systems with which it interacts with when the domain name service returns an unknown host exception. Java may return an unknown host exception when there is unexpected interruption in domain name service so this setting protects the Directory Server from temporary DNS server outages if previous results have been cached.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"tracked_application": schema.SetAttribute{
				Description: "Specifies criteria for identifying specific applications that access the server to enable tracking throughput and latency of LDAP operations issued by an application.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"jmx_value_behavior": schema.StringAttribute{
				Description: "Specifies how a Java type is chosen for monitor attributes exposed as JMX attribute values.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"jmx_use_legacy_mbean_names": schema.BoolAttribute{
				Description: "When set to true, the server will use its original, non-standard JMX MBean names for the monitoring MBeans. These include RDN keys of \"Rdn1\" and \"Rdn2\" instead of the recommended \"type\" and \"name\" keys. This should option should only be enabled for installations that have monitoring infrastructure that depends on the old keys.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
}

// Read a GlobalConfigurationResponse object into the model struct
func readGlobalConfigurationResponseDataSource(ctx context.Context, r *client.GlobalConfigurationResponse, state *globalConfigurationDataSourceModel, diagnostics *diag.Diagnostics) {
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
	state.InstanceName = types.StringValue(r.InstanceName)
	state.Location = internaltypes.StringTypeOrNil(r.Location, false)
	state.ConfigurationServerGroup = internaltypes.StringTypeOrNil(r.ConfigurationServerGroup, false)
	state.ForceAsMasterForMirroredData = internaltypes.BoolTypeOrNil(r.ForceAsMasterForMirroredData)
	state.EncryptData = internaltypes.BoolTypeOrNil(r.EncryptData)
	state.EncryptionSettingsCipherStreamProvider = internaltypes.StringTypeOrNil(r.EncryptionSettingsCipherStreamProvider, false)
	state.EncryptBackupsByDefault = internaltypes.BoolTypeOrNil(r.EncryptBackupsByDefault)
	state.BackupEncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.BackupEncryptionSettingsDefinitionID, false)
	state.EncryptLDIFExportsByDefault = internaltypes.BoolTypeOrNil(r.EncryptLDIFExportsByDefault)
	state.LdifExportEncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.LdifExportEncryptionSettingsDefinitionID, false)
	state.AutomaticallyCompressEncryptedLDIFExports = internaltypes.BoolTypeOrNil(r.AutomaticallyCompressEncryptedLDIFExports)
	state.RedactSensitiveValuesInConfigLogs = internaltypes.BoolTypeOrNil(r.RedactSensitiveValuesInConfigLogs)
	state.SensitiveAttribute = internaltypes.GetStringSet(r.SensitiveAttribute)
	state.RejectInsecureRequests = internaltypes.BoolTypeOrNil(r.RejectInsecureRequests)
	state.AllowedInsecureRequestCriteria = internaltypes.StringTypeOrNil(r.AllowedInsecureRequestCriteria, false)
	state.RejectUnauthenticatedRequests = internaltypes.BoolTypeOrNil(r.RejectUnauthenticatedRequests)
	state.AllowedUnauthenticatedRequestCriteria = internaltypes.StringTypeOrNil(r.AllowedUnauthenticatedRequestCriteria, false)
	state.BindWithDNRequiresPassword = internaltypes.BoolTypeOrNil(r.BindWithDNRequiresPassword)
	state.DisabledPrivilege = internaltypes.GetStringSet(
		client.StringSliceEnumglobalConfigurationDisabledPrivilegeProp(r.DisabledPrivilege))
	state.DefaultPasswordPolicy = types.StringValue(r.DefaultPasswordPolicy)
	state.MaximumUserDataPasswordPoliciesToCache = internaltypes.Int64TypeOrNil(r.MaximumUserDataPasswordPoliciesToCache)
	state.ProxiedAuthorizationIdentityMapper = types.StringValue(r.ProxiedAuthorizationIdentityMapper)
	state.VerifyEntryDigests = internaltypes.BoolTypeOrNil(r.VerifyEntryDigests)
	state.AllowedInsecureTLSProtocol = internaltypes.GetStringSet(
		client.StringSliceEnumglobalConfigurationAllowedInsecureTLSProtocolProp(r.AllowedInsecureTLSProtocol))
	state.AllowInsecureLocalJMXConnections = internaltypes.BoolTypeOrNil(r.AllowInsecureLocalJMXConnections)
	state.DefaultInternalOperationClientConnectionPolicy = internaltypes.StringTypeOrNil(r.DefaultInternalOperationClientConnectionPolicy, false)
	state.SizeLimit = internaltypes.Int64TypeOrNil(r.SizeLimit)
	state.UnauthenticatedSizeLimit = internaltypes.Int64TypeOrNil(r.UnauthenticatedSizeLimit)
	state.TimeLimit = internaltypes.StringTypeOrNil(r.TimeLimit, false)
	state.UnauthenticatedTimeLimit = internaltypes.StringTypeOrNil(r.UnauthenticatedTimeLimit, false)
	state.IdleTimeLimit = internaltypes.StringTypeOrNil(r.IdleTimeLimit, false)
	state.UnauthenticatedIdleTimeLimit = internaltypes.StringTypeOrNil(r.UnauthenticatedIdleTimeLimit, false)
	state.LookthroughLimit = internaltypes.Int64TypeOrNil(r.LookthroughLimit)
	state.UnauthenticatedLookthroughLimit = internaltypes.Int64TypeOrNil(r.UnauthenticatedLookthroughLimit)
	state.LdapJoinSizeLimit = internaltypes.Int64TypeOrNil(r.LdapJoinSizeLimit)
	state.MaximumConcurrentConnections = internaltypes.Int64TypeOrNil(r.MaximumConcurrentConnections)
	state.MaximumConcurrentConnectionsPerIPAddress = internaltypes.Int64TypeOrNil(r.MaximumConcurrentConnectionsPerIPAddress)
	state.MaximumConcurrentConnectionsPerBindDN = internaltypes.Int64TypeOrNil(r.MaximumConcurrentConnectionsPerBindDN)
	state.MaximumConcurrentUnindexedSearches = internaltypes.Int64TypeOrNil(r.MaximumConcurrentUnindexedSearches)
	state.MaximumAttributesPerAddRequest = internaltypes.Int64TypeOrNil(r.MaximumAttributesPerAddRequest)
	state.MaximumModificationsPerModifyRequest = internaltypes.Int64TypeOrNil(r.MaximumModificationsPerModifyRequest)
	state.BackgroundThreadForEachPersistentSearch = internaltypes.BoolTypeOrNil(r.BackgroundThreadForEachPersistentSearch)
	state.AllowAttributeNameExceptions = internaltypes.BoolTypeOrNil(r.AllowAttributeNameExceptions)
	state.InvalidAttributeSyntaxBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumglobalConfigurationInvalidAttributeSyntaxBehaviorProp(r.InvalidAttributeSyntaxBehavior), false)
	state.PermitSyntaxViolationsForAttribute = internaltypes.GetStringSet(r.PermitSyntaxViolationsForAttribute)
	state.SingleStructuralObjectclassBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumglobalConfigurationSingleStructuralObjectclassBehaviorProp(r.SingleStructuralObjectclassBehavior), false)
	state.AttributesModifiableWithIgnoreNoUserModificationRequestControl = internaltypes.GetStringSet(
		client.StringSliceEnumglobalConfigurationAttributesModifiableWithIgnoreNoUserModificationRequestControlProp(r.AttributesModifiableWithIgnoreNoUserModificationRequestControl))
	state.MaximumServerOutLogFileSize = internaltypes.StringTypeOrNil(r.MaximumServerOutLogFileSize, false)
	state.MaximumServerOutLogFileCount = internaltypes.Int64TypeOrNil(r.MaximumServerOutLogFileCount)
	state.StartupErrorLoggerOutputLocation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumglobalConfigurationStartupErrorLoggerOutputLocationProp(r.StartupErrorLoggerOutputLocation), false)
	state.ExitOnJVMError = internaltypes.BoolTypeOrNil(r.ExitOnJVMError)
	state.ServerErrorResultCode = internaltypes.Int64TypeOrNil(r.ServerErrorResultCode)
	state.ResultCodeMap = internaltypes.StringTypeOrNil(r.ResultCodeMap, false)
	state.ReturnBindErrorMessages = internaltypes.BoolTypeOrNil(r.ReturnBindErrorMessages)
	state.NotifyAbandonedOperations = internaltypes.BoolTypeOrNil(r.NotifyAbandonedOperations)
	state.DuplicateErrorLogLimit = types.Int64Value(r.DuplicateErrorLogLimit)
	state.DuplicateErrorLogTimeLimit = types.StringValue(r.DuplicateErrorLogTimeLimit)
	state.DuplicateAlertLimit = types.Int64Value(r.DuplicateAlertLimit)
	state.DuplicateAlertTimeLimit = types.StringValue(r.DuplicateAlertTimeLimit)
	state.WritabilityMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumglobalConfigurationWritabilityModeProp(r.WritabilityMode), false)
	state.UnrecoverableDatabaseErrorMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumglobalConfigurationUnrecoverableDatabaseErrorModeProp(r.UnrecoverableDatabaseErrorMode), false)
	state.DatabaseOnVirtualizedOrNetworkStorage = internaltypes.BoolTypeOrNil(r.DatabaseOnVirtualizedOrNetworkStorage)
	state.AutoNameWithEntryUUIDConnectionCriteria = internaltypes.StringTypeOrNil(r.AutoNameWithEntryUUIDConnectionCriteria, false)
	state.AutoNameWithEntryUUIDRequestCriteria = internaltypes.StringTypeOrNil(r.AutoNameWithEntryUUIDRequestCriteria, false)
	state.SoftDeletePolicy = internaltypes.StringTypeOrNil(r.SoftDeletePolicy, false)
	state.SubtreeAccessibilityAlertTimeLimit = internaltypes.StringTypeOrNil(r.SubtreeAccessibilityAlertTimeLimit, false)
	state.WarnForBackendsWithMultipleBaseDns = internaltypes.BoolTypeOrNil(r.WarnForBackendsWithMultipleBaseDns)
	state.ForcedGCPrimeDuration = internaltypes.StringTypeOrNil(r.ForcedGCPrimeDuration, false)
	state.ReplicationSetName = internaltypes.StringTypeOrNil(r.ReplicationSetName, false)
	state.StartupMinReplicationBacklogCount = types.Int64Value(r.StartupMinReplicationBacklogCount)
	state.ReplicationBacklogCountAlertThreshold = types.Int64Value(r.ReplicationBacklogCountAlertThreshold)
	state.ReplicationBacklogDurationAlertThreshold = types.StringValue(r.ReplicationBacklogDurationAlertThreshold)
	state.ReplicationAssuranceSourceTimeoutSuspendDuration = types.StringValue(r.ReplicationAssuranceSourceTimeoutSuspendDuration)
	state.ReplicationAssuranceSourceBacklogFastStartThreshold = types.Int64Value(r.ReplicationAssuranceSourceBacklogFastStartThreshold)
	state.ReplicationHistoryLimit = internaltypes.Int64TypeOrNil(r.ReplicationHistoryLimit)
	state.AllowInheritedReplicationOfSubordinateBackends = types.BoolValue(r.AllowInheritedReplicationOfSubordinateBackends)
	state.ReplicationPurgeObsoleteReplicas = internaltypes.BoolTypeOrNil(r.ReplicationPurgeObsoleteReplicas)
	state.SmtpServer = internaltypes.GetStringSet(r.SmtpServer)
	state.MaxSMTPConnectionCount = internaltypes.Int64TypeOrNil(r.MaxSMTPConnectionCount)
	state.MaxSMTPConnectionAge = internaltypes.StringTypeOrNil(r.MaxSMTPConnectionAge, false)
	state.SmtpConnectionHealthCheckInterval = internaltypes.StringTypeOrNil(r.SmtpConnectionHealthCheckInterval, false)
	state.AllowedTask = internaltypes.GetStringSet(r.AllowedTask)
	state.EnableSubOperationTimer = internaltypes.BoolTypeOrNil(r.EnableSubOperationTimer)
	state.MaximumShutdownTime = internaltypes.StringTypeOrNil(r.MaximumShutdownTime, false)
	state.NetworkAddressCacheTTL = internaltypes.StringTypeOrNil(r.NetworkAddressCacheTTL, false)
	state.NetworkAddressOutageCacheEnabled = internaltypes.BoolTypeOrNil(r.NetworkAddressOutageCacheEnabled)
	state.TrackedApplication = internaltypes.GetStringSet(r.TrackedApplication)
	state.JmxValueBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumglobalConfigurationJmxValueBehaviorProp(r.JmxValueBehavior), false)
	state.JmxUseLegacyMbeanNames = internaltypes.BoolTypeOrNil(r.JmxUseLegacyMbeanNames)
}

// Read resource information
func (r *globalConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state globalConfigurationDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.GlobalConfigurationApi.GetGlobalConfiguration(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Global Configuration", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readGlobalConfigurationResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
