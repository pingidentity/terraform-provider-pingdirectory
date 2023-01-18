package config

import (
	"context"
	"time"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &globalConfigurationResource{}
	_ resource.ResourceWithConfigure   = &globalConfigurationResource{}
	_ resource.ResourceWithImportState = &globalConfigurationResource{}
)

// Create a Global Configuration resource
func NewGlobalConfigurationResource() resource.Resource {
	return &globalConfigurationResource{}
}

// globalConfigurationResource is the resource implementation.
type globalConfigurationResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// globalConfigurationResourceModel maps the resource schema data.
type globalConfigurationResourceModel struct {
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
	TimeLimit                                                      types.String `tfsdk:"time_limit"`
	IdleTimeLimit                                                  types.String `tfsdk:"idle_time_limit"`
	LookthroughLimit                                               types.Int64  `tfsdk:"lookthrough_limit"`
	LdapJoinSizeLimit                                              types.Int64  `tfsdk:"ldap_join_size_limit"`
	MaximumConcurrentConnections                                   types.Int64  `tfsdk:"maximum_concurrent_connections"`
	MaximumConcurrentConnectionsPerIPAddress                       types.Int64  `tfsdk:"maximum_concurrent_connections_per_id_address"`
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
	LastUpdated                                                    types.String `tfsdk:"last_updated"`
	Notifications                                                  types.Set    `tfsdk:"notifications"`
	RequiredActions                                                types.Set    `tfsdk:"required_actions"`
}

// Metadata returns the resource type name.
func (r *globalConfigurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_global_configuration"
}

// GetSchema defines the schema for the resource.
func (r *globalConfigurationResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	schema := tfsdk.Schema{
		Description: "Manages the global configuration.",
		// All are considered computed, since we are importing the existing global
		// configuration from a server, rather than "creating" the global configuration
		// like a typical Terraform resource.
		Attributes: map[string]tfsdk.Attribute{
			"instance_name": {
				Description: "A name that may be used to uniquely identify this Directory Server instance among other instances in the environment.",
				Type:        types.StringType,
				// instance name is read-only after setup, so Terraform can't change it
				Required: false,
				Optional: false,
				Computed: true,
			},
			"location": {
				Description: "Specifies the location for this Directory Server. Operations performed which involve communication with other servers may prefer servers in the same location to help ensure low-latency responses.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"configuration_server_group": {
				Description: "When this property is set, changes made to this server using the console or dsconfig can be automatically applied to all servers in the specified server group.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"force_as_master_for_mirrored_data": {
				Description: "Indicates whether this server should be forced to assume the master role if no other suitable server is found to act as master or if multiple masters are detected. A master is only needed when changes are made to mirrored data, i.e. data specific to the topology itself and cluster-wide configuration data.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"encrypt_data": {
				Description: "Indicates whether the Directory Server should encrypt the data that it stores in all components that support it. This may include certain types of backends (including local DB and large attribute backends), the LDAP changelog, and the replication server database.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"encryption_settings_cipher_stream_provider": {
				Description: "Specifies the cipher stream provider that should be used to protect the contents of the encryption settings database.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"encrypt_backups_by_default": {
				Description: "Indicates whether the server should encrypt backups by default.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"backup_encryption_settings_definition_id": {
				Description: "The unique identifier for the encryption settings definition to use to generate the encryption key for encrypted backups by default.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"encrypt_ldif_exports_by_default": {
				Description: "Indicates whether the server should encrypt LDIF exports by default.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"ldif_export_encryption_settings_definition_id": {
				Description: "The unique identifier for the encryption settings definition to use to generate the encryption key for encrypted LDIF exports by default.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"automatically_compress_encrypted_ldif_exports": {
				Description: "Indicates whether to automatically compress LDIF exports that are also encrypted.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"redact_sensitive_values_in_config_logs": {
				Description: "Indicates whether the values of sensitive configuration properties should be redacted when logging configuration changes, including in the configuration audit log, the error log, and the server.out log file.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"sensitive_attribute": {
				Description: "Provides the ability to indicate that some attributes should be considered sensitive and additional protection should be in place when interacting with those attributes.",
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Optional: true,
				Computed: true,
			},
			"reject_insecure_requests": {
				Description: "Indicates whether the Directory Server should reject any LDAP request (other than StartTLS) received from a client that is not using an encrypted connection.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"allowed_insecure_request_criteria": {
				Description: "A set of criteria that may be used to match LDAP requests that may be permitted over an insecure connection even if reject-insecure-requests is true. Note that some types of requests will always be permitted, including StartTLS and start administrative session requests.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"reject_unauthenticated_requests": {
				Description: "Indicates whether the Directory Server should reject any LDAP request (other than bind or StartTLS requests) received from a client that has not yet been authenticated, whose last authentication attempt was unsuccessful, or whose last authentication attempt used anonymous authentication.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"allowed_unauthenticated_request_criteria": {
				Description: "A set of criteria that may be used to match LDAP requests that may be permitted over an unauthenticated connection even if reject-unauthenticated-requests is true. Note that some types of requests will always be permitted, including bind, StartTLS, and start administrative session requests.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"bind_with_dn_requires_password": {
				Description: "Indicates whether the Directory Server should reject any simple bind request that contains a DN but no password.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"disabled_privilege": {
				Description: "Specifies the name of a privilege that should not be evaluated by the server.",
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Optional: true,
				Computed: true,
			},
			"default_password_policy": {
				Description: "Specifies the name of the password policy that is in effect for users whose entries do not specify an alternate password policy (either via a real or virtual attribute).",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"maximum_user_data_password_policies_to_cache": {
				Description: "Specifies the maximum number of password policies that are defined in the user data (that is, outside of the configuration) that the server should cache in memory for faster access. A value of zero indicates that the server should not cache any user data password policies.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"proxied_authorization_identity_mapper": {
				Description: "Specifies the name of the identity mapper to map authorization ID values (using the \"u:\" form) provided in the proxied authorization control to the corresponding user entry.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"verify_entry_digests": {
				Description: "Indicates whether the digest should always be verified whenever an entry containing a digest is decoded. If this is \"true\", then if a digest exists, it will always be verified. Otherwise, the digest will be written when encoding entries but ignored when decoding entries but may still be available for other verification processing.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"allowed_insecure_tls_protocol": {
				Description: "Specifies a set of TLS protocols that will be permitted for use in the server even though there may be known vulnerabilities that could cause their use to be unsafe in some conditions. Enabling support for insecure TLS protocols is discouraged, and is generally recommended only as a short-term measure to permit legacy clients to interact with the server until they can be updated to support more secure communication protocols.",
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Optional: true,
				Computed: true,
			},
			"allow_insecure_local_jmx_connections": {
				Description: "Indicates that processes attaching to this server's local JVM are allowed to access internal data through JMX without the authentication requirements that remote JMX connections are subject to. Please review and understand the data that this option will expose (such as cn=monitor) to client applications to ensure there are no security concerns.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"default_internal_operation_client_connection_policy": {
				Description: "Specifies the client connection policy that will be used by default for internal operations.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"size_limit": {
				Description: "Specifies the maximum number of entries that the Directory Server should return to the client during a search operation.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"time_limit": {
				Description: "Specifies the maximum length of time that the Directory Server should be allowed to spend processing a search operation.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"idle_time_limit": {
				Description: "Specifies the maximum length of time that a client connection may remain established since its last completed operation.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"lookthrough_limit": {
				Description: "Specifies the maximum number of entries that the Directory Server should \"look through\" in the course of processing a search request.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"ldap_join_size_limit": {
				Description: "Specifies the maximum number of entries that may be directly joined with any individual search result entry.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"maximum_concurrent_connections": {
				Description: "Specifies the maximum number of LDAP client connections which may be established to this Directory Server at the same time.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"maximum_concurrent_connections_per_id_address": {
				Description: "Specifies the maximum number of LDAP client connections originating from the same IP address which may be established to this Directory Server at the same time.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"maximum_concurrent_connections_per_bind_dn": {
				Description: "Specifies the maximum number of LDAP client connections which may be established to this Directory Server at the same time and authenticated as the same user.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"maximum_concurrent_unindexed_searches": {
				Description: "Specifies the maximum number of unindexed searches that may be in progress in this backend at any given time. Any unindexed searches requested while the maximum number of unindexed searches are already being processed will be rejected. A value of zero indicates that no limit will be enforced.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"maximum_attributes_per_add_request": {
				Description: "Specifies the maximum number of attributes that may be included in an add request. This property does not impose any limit on the number of values that an attribute may have.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"maximum_modifications_per_modify_request": {
				Description: "Specifies the maximum number of modifications that may be included in a modify request. This property does not impose any limit on the number of attribute values that a modification may have.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"background_thread_for_each_persistent_search": {
				Description: "Indicates whether the server should use a separate background thread for each persistent search.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"allow_attribute_name_exceptions": {
				Description: "Indicates whether the Directory Server should allow underscores in attribute names and allow attribute names to begin with numeric digits (both of which are violations of the LDAP standards).",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"invalid_attribute_syntax_behavior": {
				Description: "Specifies how the Directory Server should handle operations whenever an attribute value violates the associated attribute syntax.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"permit_syntax_violations_for_attribute": {
				Description: "Specifies a set of attribute types for which the server will permit values that do not conform to the associated attribute syntax.",
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Optional: true,
				Computed: true,
			},
			"single_structural_objectclass_behavior": {
				Description: "Specifies how the Directory Server should handle operations for an entry does not contain a structural object class, or for an entry that contains multiple structural classes.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"attributes_modifiable_with_ignore_no_user_modification_request_control": {
				Description: "Specifies the operational attribute types that are defined in the schema with the NO-USER-MODIFICATION constraint that the server will allow to be altered if the associated request contains the ignore NO-USER-MODIFICATION request control.",
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Optional: true,
				Computed: true,
			},
			"maximum_server_out_log_file_size": {
				Description: "The maximum allowed size that the server.out log file will be allowed to have. If a write would cause the file to exceed this size, then the current file will be rotated out of place and a new empty file will be created and the message written to it.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"maximum_server_out_log_file_count": {
				Description: "The maximum number of server.out log files (including the current active log file) that should be retained. When rotating the log file, if the total number of files exceeds this count, then the oldest file(s) will be removed so that the total number of log files is within this limit.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"startup_error_logger_output_location": {
				Description: "Specifies how the server should handle error log messages (which may include errors, warnings, and notices) generated during startup. All of these messages will be written to all configured error loggers, but they may also be written to other locations (like standard output, standard error, or the server.out log file) so that they are displayed on the console when the server is starting.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"exit_on_jvm_error": {
				Description: "Indicates whether the Directory Server should be shut down if a severe error is raised (e.g., an out of memory error) which may prevent the JVM from continuing to run properly.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"server_error_result_code": {
				Description: "Specifies the numeric value of the result code when request processing fails due to an internal server error.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"result_code_map": {
				Description: "Specifies a result code map that should be used for clients that do not have a map associated with their client connection policy. If the associated client connection policy has a result code map, then that map will be used instead. If no map is associated either with the client connection policy or the global configuration, then an internal default will be used.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"return_bind_error_messages": {
				Description: "Indicates whether responses for failed bind operations should include a message string providing the reason for the authentication failure.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"notify_abandoned_operations": {
				Description: "Indicates whether the Directory Server should send a response to any operation that is interrupted via an abandon request.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"duplicate_error_log_limit": {
				Description: "Specifies the maximum number of duplicate error log messages that should be logged in the time window specified by the duplicate-error-log-time-limit property.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"duplicate_error_log_time_limit": {
				Description: "Specifies the length of time that must expire before duplicate log messages above the duplicate-error-log-limit threshold are logged again to the error log.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"duplicate_alert_limit": {
				Description: "Specifies the maximum number of duplicate alert messages that should be sent via the administrative alert framework in the time window specified by the duplicate-alert-time-limit property.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"duplicate_alert_time_limit": {
				Description: "Specifies the length of time that must expire before duplicate messages are sent via the administrative alert framework.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"writability_mode": {
				Description: "Specifies the kinds of write operations the Directory Server can process.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"unrecoverable_database_error_mode": {
				Description: "Specifies the action which should be taken for any database that experiences an unrecoverable error. Action applies to local database backends and the replication recent changes database.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"database_on_virtualized_or_network_storage": {
				Description: "This setting provides data integrity options when the Directory Server is installed with a database on a network storage device. A storage device may be accessed directly by a physical server, or indirectly through a virtual machine running on a hypervisor. Enabling this setting will apply changes to all Local DB Backends, the LDAP Changelog Backend, and the replication changelog database.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"auto_name_with_entry_uuid_connection_criteria": {
				Description: "Connection criteria that may be used to identify clients whose add requests should use entryUUID as the naming attribute.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"auto_name_with_entry_uuid_request_criteria": {
				Description: "Request criteria that may be used to identify add requests that should use entryUUID as the naming attribute.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"soft_delete_policy": {
				Description: "Specifies the soft delete policy that will be used by default for delete operations. Soft delete operations introduce the ability to control the server behavior of the delete operation. Instead of performing a permanent delete of an entry, deleted entries can be retained as soft deleted entries by their entryUUID values and are available for undelete at a later time. In addition to a soft delete policy enabling soft deletes, delete operations sent to the server must have the soft delete request control present with sufficient access privileges to access the soft delete request control.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"subtree_accessibility_alert_time_limit": {
				Description: "Specifies the length of time that a subtree may remain hidden or read-only before an administrative alert is sent.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"warn_for_backends_with_multiple_base_dns": {
				Description: "Indicates whether the server should issue a warning when enabling a backend that contains multiple base DNs.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"forced_gc_prime_duration": {
				Description: "Specifies the minimum length of time required for backend or request processor initialization that will trigger the server to force an explicit garbage collection. A value of \"0 seconds\" indicates that the server should never invoke an explicit garbage collection regardless of the length of time required to initialize the server backends.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"replication_set_name": {
				Description: "The name of the replication set assigned to this Directory Server. Restricted domains are only replicated within instances using the same replication set name.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"startup_min_replication_backlog_count": {
				Description: "The number of outstanding changes any replica can have before the Directory Server will start accepting connections. The Directory Server may never accept connections if this setting is too low. If you are unsure which value to use, you can use the number of expected updates within a five second interval.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"replication_backlog_count_alert_threshold": {
				Description: "An alert is sent when the number of outstanding replication changes for the Directory Server has exceeded this threshold for longer than the replication backlog duration alert threshold.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"replication_backlog_duration_alert_threshold": {
				Description: "An alert is sent when the number of outstanding replication changes for the Directory Server has exceeded the replication backlog count alert threshold for longer than this duration.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"replication_assurance_source_timeout_suspend_duration": {
				Description: "The amount of time a replication assurance source (i.e. a peer Directory Server) will be suspended from assurance requirements on this Directory Server if it experiences an assurance timeout.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"replication_assurance_source_backlog_fast_start_threshold": {
				Description: "The maximum number of replication backlog updates a replication assurance source (i.e. a peer Directory Server) can have and be immediately recognized as an available assurance source by this Directory Server.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"replication_history_limit": {
				Description: "Specifies the size limit for historical information.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"allow_inherited_replication_of_subordinate_backends": {
				Description: "Allow replication to be inherited by subordinate/child backends.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"replication_purge_obsolete_replicas": {
				Description: "Indicates whether state about obsolete replicas is automatically purged.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"smtp_server": {
				Description: "Specifies the set of servers that will be used to send email messages. The order in which the servers are listed indicates the order in which the Directory Server will attempt to use them in the course of sending a message. The first attempt will always go to the server at the top of the list, and servers further down the list will only be used if none of the servers listed above it were able to successfully send the message.",
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Optional: true,
				Computed: true,
			},
			"max_smtp_connection_count": {
				Description: "The maximum number of SMTP connections that will be maintained for delivering email messages.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"max_smtp_connection_age": {
				Description: "The maximum length of time that a connection to an SMTP server should be considered valid.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"smtp_connection_health_check_interval": {
				Description: "The length of time between checks to ensure that available SMTP connections are still valid.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"allowed_task": {
				Description: "Specifies the fully-qualified name of a Java class that may be invoked in the server.",
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Optional: true,
				Computed: true,
			},
			"enable_sub_operation_timer": {
				Description: "Indicates whether the Directory Server should attempt to record information about the length of time required to process various phases of an operation. Enabling this feature may impact performance, but could make it easier to identify potential bottlenecks in operation processing.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"maximum_shutdown_time": {
				Description: "Specifies the maximum amount of time the shutdown of Directory Server may take.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"network_address_cache_ttl": {
				Description: "Specifies the length of time that the Directory Server should cache the IP addresses associated with the names of systems with which it interacts.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"network_address_outage_cache_enabled": {
				Description: "Specifies whether the Directory Server should cache the last valid IP addresses associated with the names of systems with which it interacts with when the domain name service returns an unknown host exception. Java may return an unknown host exception when there is unexpected interruption in domain name service so this setting protects the Directory Server from temporary DNS server outages if previous results have been cached.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"tracked_application": {
				Description: "Specifies criteria for identifying specific applications that access the server to enable tracking throughput and latency of LDAP operations issued by an application.",
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Optional: true,
				Computed: true,
			},
			"jmx_value_behavior": {
				Description: "Specifies how a Java type is chosen for monitor attributes exposed as JMX attribute values.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"jmx_use_legacy_mbean_names": {
				Description: "When set to true, the server will use its original, non-standard JMX MBean names for the monitoring MBeans. These include RDN keys of \"Rdn1\" and \"Rdn2\" instead of the recommended \"type\" and \"name\" keys. This should option should only be enabled for installations that have monitoring infrastructure that depends on the old keys.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
		},
	}
	AddCommonSchema(&schema)
	return schema, nil
}

// Configure adds the provider configured client to the resource.
func (r *globalConfigurationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Create a new resource
// For global config, create doesn't actually "create" anything - it "adopts" the servers existing
// global configuration into management by terraform. This method reads the existing global config
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *globalConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan globalConfigurationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	getResp, httpResp, err := r.apiClient.GlobalConfigurationApi.GetGlobalConfiguration(ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the global configuration", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := getResp.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read existing global config
	var state globalConfigurationResourceModel
	readGlobalConfigurationResponse(ctx, getResp, &state)

	// Determine what changes need to be made to match the plan
	updateGCRequest := r.apiClient.GlobalConfigurationApi.UpdateGlobalConfiguration(ProviderBasicAuthContext(ctx, r.providerConfig))
	ops := createGlobalConfigurationOperations(plan, state)

	if len(ops) > 0 {
		updateGCRequest = updateGCRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)
		globalResp, httpResp, err := r.apiClient.GlobalConfigurationApi.UpdateGlobalConfigurationExecute(updateGCRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the global configuration", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := globalResp.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGlobalConfigurationResponse(ctx, globalResp, &plan)
		// Populate Computed attribute values
		plan.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		// Just put the initial read into the plan
		plan = state
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *globalConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state globalConfigurationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	getResp, httpResp, err := r.apiClient.GlobalConfigurationApi.GetGlobalConfiguration(ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the global configuration", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := getResp.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readGlobalConfigurationResponse(ctx, getResp, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read a GlobalConfigurationRespnse object into the model struct
func readGlobalConfigurationResponse(ctx context.Context, r *client.GlobalConfigurationResponse, state *globalConfigurationResourceModel) {
	// Placeholder Id value for acceptance test framework
	state.Id = types.StringValue(r.InstanceName)
	state.InstanceName = types.StringValue(r.InstanceName)
	state.Location = internaltypes.StringTypeOrNil(r.Location, true)
	state.ConfigurationServerGroup = internaltypes.StringTypeOrNil(r.ConfigurationServerGroup, true)
	state.ForceAsMasterForMirroredData = internaltypes.BoolTypeOrNil(r.ForceAsMasterForMirroredData)
	state.EncryptData = internaltypes.BoolTypeOrNil(r.EncryptData)
	state.EncryptionSettingsCipherStreamProvider = internaltypes.StringTypeOrNil(r.EncryptionSettingsCipherStreamProvider, true)
	state.EncryptBackupsByDefault = internaltypes.BoolTypeOrNil(r.EncryptBackupsByDefault)
	state.BackupEncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.BackupEncryptionSettingsDefinitionID, true)
	state.EncryptLDIFExportsByDefault = internaltypes.BoolTypeOrNil(r.EncryptLDIFExportsByDefault)
	state.LdifExportEncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.LdifExportEncryptionSettingsDefinitionID, true)
	state.AutomaticallyCompressEncryptedLDIFExports = internaltypes.BoolTypeOrNil(r.AutomaticallyCompressEncryptedLDIFExports)
	state.RedactSensitiveValuesInConfigLogs = internaltypes.BoolTypeOrNil(r.RedactSensitiveValuesInConfigLogs)
	state.SensitiveAttribute = internaltypes.GetStringSet(r.SensitiveAttribute)
	state.RejectInsecureRequests = internaltypes.BoolTypeOrNil(r.RejectInsecureRequests)
	state.AllowedInsecureRequestCriteria = internaltypes.StringTypeOrNil(r.AllowedInsecureRequestCriteria, true)
	state.RejectUnauthenticatedRequests = internaltypes.BoolTypeOrNil(r.RejectUnauthenticatedRequests)
	state.AllowedUnauthenticatedRequestCriteria = internaltypes.StringTypeOrNil(r.AllowedUnauthenticatedRequestCriteria, true)
	state.BindWithDNRequiresPassword = internaltypes.BoolTypeOrNil(r.BindWithDNRequiresPassword)
	state.DisabledPrivilege = internaltypes.GetEnumSet(r.DisabledPrivilege)
	state.DefaultPasswordPolicy = types.StringValue(r.DefaultPasswordPolicy)
	state.MaximumUserDataPasswordPoliciesToCache = internaltypes.Int64TypeOrNil(r.MaximumUserDataPasswordPoliciesToCache)
	state.ProxiedAuthorizationIdentityMapper = types.StringValue(r.ProxiedAuthorizationIdentityMapper)
	state.VerifyEntryDigests = internaltypes.BoolTypeOrNil(r.VerifyEntryDigests)
	state.AllowedInsecureTLSProtocol = internaltypes.GetEnumSet(r.AllowedInsecureTLSProtocol)
	state.AllowInsecureLocalJMXConnections = internaltypes.BoolTypeOrNil(r.AllowInsecureLocalJMXConnections)
	state.DefaultInternalOperationClientConnectionPolicy = internaltypes.StringTypeOrNil(r.DefaultInternalOperationClientConnectionPolicy, true)
	state.SizeLimit = internaltypes.Int64TypeOrNil(r.SizeLimit)
	state.TimeLimit = internaltypes.StringTypeOrNil(r.TimeLimit, true)
	state.IdleTimeLimit = internaltypes.StringTypeOrNil(r.IdleTimeLimit, true)
	state.LookthroughLimit = internaltypes.Int64TypeOrNil(r.LookthroughLimit)
	state.LdapJoinSizeLimit = internaltypes.Int64TypeOrNil(r.LdapJoinSizeLimit)
	state.MaximumConcurrentConnections = internaltypes.Int64TypeOrNil(r.MaximumConcurrentConnections)
	state.MaximumConcurrentConnectionsPerIPAddress = internaltypes.Int64TypeOrNil(r.MaximumConcurrentConnectionsPerIPAddress)
	state.MaximumConcurrentConnectionsPerBindDN = internaltypes.Int64TypeOrNil(r.MaximumConcurrentConnectionsPerBindDN)
	state.MaximumConcurrentUnindexedSearches = internaltypes.Int64TypeOrNil(r.MaximumConcurrentUnindexedSearches)
	state.MaximumAttributesPerAddRequest = internaltypes.Int64TypeOrNil(r.MaximumAttributesPerAddRequest)
	state.MaximumModificationsPerModifyRequest = internaltypes.Int64TypeOrNil(r.MaximumModificationsPerModifyRequest)
	state.BackgroundThreadForEachPersistentSearch = internaltypes.BoolTypeOrNil(r.BackgroundThreadForEachPersistentSearch)
	state.AllowAttributeNameExceptions = internaltypes.BoolTypeOrNil(r.AllowAttributeNameExceptions)
	state.InvalidAttributeSyntaxBehavior = internaltypes.StringerStringTypeOrNil(r.InvalidAttributeSyntaxBehavior)
	state.PermitSyntaxViolationsForAttribute = internaltypes.GetStringSet(r.PermitSyntaxViolationsForAttribute)
	state.SingleStructuralObjectclassBehavior = internaltypes.StringerStringTypeOrNil(r.SingleStructuralObjectclassBehavior)
	state.AttributesModifiableWithIgnoreNoUserModificationRequestControl = internaltypes.GetEnumSet(r.AttributesModifiableWithIgnoreNoUserModificationRequestControl)
	state.MaximumServerOutLogFileSize = internaltypes.StringTypeOrNil(r.MaximumServerOutLogFileSize, true)
	state.MaximumServerOutLogFileCount = internaltypes.Int64TypeOrNil(r.MaximumServerOutLogFileCount)
	state.StartupErrorLoggerOutputLocation = internaltypes.StringerStringTypeOrNil(r.StartupErrorLoggerOutputLocation)
	state.ExitOnJVMError = internaltypes.BoolTypeOrNil(r.ExitOnJVMError)
	state.ServerErrorResultCode = internaltypes.Int64TypeOrNil(r.ServerErrorResultCode)
	state.ResultCodeMap = internaltypes.StringTypeOrNil(r.ResultCodeMap, true)
	state.ReturnBindErrorMessages = internaltypes.BoolTypeOrNil(r.ReturnBindErrorMessages)
	state.NotifyAbandonedOperations = internaltypes.BoolTypeOrNil(r.NotifyAbandonedOperations)
	state.DuplicateErrorLogLimit = types.Int64Value(int64(r.DuplicateErrorLogLimit))
	state.DuplicateErrorLogTimeLimit = types.StringValue(r.DuplicateErrorLogTimeLimit)
	state.DuplicateAlertLimit = types.Int64Value(int64(r.DuplicateAlertLimit))
	state.DuplicateAlertTimeLimit = types.StringValue(r.DuplicateAlertTimeLimit)
	state.WritabilityMode = internaltypes.StringerStringTypeOrNil(r.WritabilityMode)
	state.UnrecoverableDatabaseErrorMode = internaltypes.StringerStringTypeOrNil(r.UnrecoverableDatabaseErrorMode)
	state.DatabaseOnVirtualizedOrNetworkStorage = internaltypes.BoolTypeOrNil(r.DatabaseOnVirtualizedOrNetworkStorage)
	state.AutoNameWithEntryUUIDConnectionCriteria = internaltypes.StringTypeOrNil(r.AutoNameWithEntryUUIDConnectionCriteria, true)
	state.AutoNameWithEntryUUIDRequestCriteria = internaltypes.StringTypeOrNil(r.AutoNameWithEntryUUIDRequestCriteria, true)
	state.SoftDeletePolicy = internaltypes.StringTypeOrNil(r.SoftDeletePolicy, true)
	state.SubtreeAccessibilityAlertTimeLimit = internaltypes.StringTypeOrNil(r.SubtreeAccessibilityAlertTimeLimit, true)
	state.WarnForBackendsWithMultipleBaseDns = internaltypes.BoolTypeOrNil(r.WarnForBackendsWithMultipleBaseDns)
	state.ForcedGCPrimeDuration = internaltypes.StringTypeOrNil(r.ForcedGCPrimeDuration, true)
	state.ReplicationSetName = internaltypes.StringTypeOrNil(r.ReplicationSetName, true)
	state.StartupMinReplicationBacklogCount = types.Int64Value(int64(r.StartupMinReplicationBacklogCount))
	state.ReplicationBacklogCountAlertThreshold = types.Int64Value(int64(r.ReplicationBacklogCountAlertThreshold))
	state.ReplicationBacklogDurationAlertThreshold = types.StringValue(r.ReplicationBacklogDurationAlertThreshold)
	state.ReplicationAssuranceSourceTimeoutSuspendDuration = types.StringValue(r.ReplicationAssuranceSourceTimeoutSuspendDuration)
	state.ReplicationAssuranceSourceBacklogFastStartThreshold = types.Int64Value(int64(r.ReplicationAssuranceSourceBacklogFastStartThreshold))
	state.ReplicationHistoryLimit = internaltypes.Int64TypeOrNil(r.ReplicationHistoryLimit)
	state.AllowInheritedReplicationOfSubordinateBackends = types.BoolValue(r.AllowInheritedReplicationOfSubordinateBackends)
	state.ReplicationPurgeObsoleteReplicas = internaltypes.BoolTypeOrNil(r.ReplicationPurgeObsoleteReplicas)
	state.SmtpServer = internaltypes.GetStringSet(r.SmtpServer)
	state.MaxSMTPConnectionCount = internaltypes.Int64TypeOrNil(r.MaxSMTPConnectionCount)
	state.MaxSMTPConnectionAge = internaltypes.StringTypeOrNil(r.MaxSMTPConnectionAge, true)
	state.SmtpConnectionHealthCheckInterval = internaltypes.StringTypeOrNil(r.SmtpConnectionHealthCheckInterval, true)
	state.AllowedTask = internaltypes.GetStringSet(r.AllowedTask)
	state.EnableSubOperationTimer = internaltypes.BoolTypeOrNil(r.EnableSubOperationTimer)
	state.MaximumShutdownTime = internaltypes.StringTypeOrNil(r.MaximumShutdownTime, true)
	state.NetworkAddressCacheTTL = internaltypes.StringTypeOrNil(r.NetworkAddressCacheTTL, true)
	state.NetworkAddressOutageCacheEnabled = internaltypes.BoolTypeOrNil(r.NetworkAddressOutageCacheEnabled)
	state.TrackedApplication = internaltypes.GetStringSet(r.TrackedApplication)
	state.JmxValueBehavior = internaltypes.StringerStringTypeOrNil(r.JmxValueBehavior)
	state.JmxUseLegacyMbeanNames = internaltypes.BoolTypeOrNil(r.JmxUseLegacyMbeanNames)
	state.Notifications, state.RequiredActions = ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20)
}

// Create any update operations necessary to make the state match the plan
func createGlobalConfigurationOperations(plan globalConfigurationResourceModel, state globalConfigurationResourceModel) []client.Operation {
	var ops []client.Operation

	operations.AddStringOperationIfNecessary(&ops, plan.Location, state.Location, "location")
	operations.AddStringOperationIfNecessary(&ops, plan.ConfigurationServerGroup, state.ConfigurationServerGroup, "configuration-server-group")
	operations.AddBoolOperationIfNecessary(&ops, plan.ForceAsMasterForMirroredData, state.ForceAsMasterForMirroredData, "force-as-master-for-mirrored-data")
	operations.AddBoolOperationIfNecessary(&ops, plan.EncryptData, state.EncryptData, "encrypt-data")
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptionSettingsCipherStreamProvider, state.EncryptionSettingsCipherStreamProvider, "encryption-settings-cipher-stream-provider")
	operations.AddBoolOperationIfNecessary(&ops, plan.EncryptBackupsByDefault, state.EncryptBackupsByDefault, "encrypt-backups-by-default")
	operations.AddStringOperationIfNecessary(&ops, plan.BackupEncryptionSettingsDefinitionID, state.BackupEncryptionSettingsDefinitionID, "backup-encryption-settings-definition-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.EncryptLDIFExportsByDefault, state.EncryptLDIFExportsByDefault, "encrypt-ldif-exports-by-default")
	operations.AddStringOperationIfNecessary(&ops, plan.LdifExportEncryptionSettingsDefinitionID, state.LdifExportEncryptionSettingsDefinitionID, "ldif-export-encryption-settings-definition-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.AutomaticallyCompressEncryptedLDIFExports, state.AutomaticallyCompressEncryptedLDIFExports, "automatically-compress-encrypted-ldif-exports")
	operations.AddBoolOperationIfNecessary(&ops, plan.RedactSensitiveValuesInConfigLogs, state.RedactSensitiveValuesInConfigLogs, "redact-sensitive-values-in-config-logs")
	operations.AddBoolOperationIfNecessary(&ops, plan.RejectInsecureRequests, state.RejectInsecureRequests, "reject-insecure-requests")
	operations.AddStringOperationIfNecessary(&ops, plan.AllowedInsecureRequestCriteria, state.AllowedInsecureRequestCriteria, "allowed-insecure-request-criteria")
	operations.AddBoolOperationIfNecessary(&ops, plan.RejectUnauthenticatedRequests, state.RejectUnauthenticatedRequests, "reject-unauthenticated-requests")
	operations.AddStringOperationIfNecessary(&ops, plan.AllowedUnauthenticatedRequestCriteria, state.AllowedUnauthenticatedRequestCriteria, "allowed-unauthenticated-request-criteria")
	operations.AddBoolOperationIfNecessary(&ops, plan.BindWithDNRequiresPassword, state.BindWithDNRequiresPassword, "bind-with-dn-requires-password")
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultPasswordPolicy, state.DefaultPasswordPolicy, "default-password-policy")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumUserDataPasswordPoliciesToCache, state.MaximumUserDataPasswordPoliciesToCache, "maximum-user-data-password-policies-cache")
	operations.AddStringOperationIfNecessary(&ops, plan.ProxiedAuthorizationIdentityMapper, state.ProxiedAuthorizationIdentityMapper, "proxied-authorization-identity-mapper")
	operations.AddBoolOperationIfNecessary(&ops, plan.VerifyEntryDigests, state.VerifyEntryDigests, "verify-entry-digests")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowInsecureLocalJMXConnections, state.AllowInsecureLocalJMXConnections, "allow-insecure-local-jmx-connections")
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultInternalOperationClientConnectionPolicy, state.DefaultInternalOperationClientConnectionPolicy, "default-internal-operation-client-connection-policy")
	operations.AddInt64OperationIfNecessary(&ops, plan.SizeLimit, state.SizeLimit, "size-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.TimeLimit, state.TimeLimit, "time-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.IdleTimeLimit, state.IdleTimeLimit, "idle-time-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.LookthroughLimit, state.LookthroughLimit, "lookthrough-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.LdapJoinSizeLimit, state.LdapJoinSizeLimit, "ldap-join-size-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumConcurrentConnections, state.MaximumConcurrentConnections, "maximum-concurrent-connections")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumConcurrentConnectionsPerIPAddress, state.MaximumConcurrentConnectionsPerIPAddress, "maximum-concurrent-connections-per-ip-address")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumConcurrentConnectionsPerBindDN, state.MaximumConcurrentConnectionsPerBindDN, "maximum-concurrent-connections-per-bind-dn")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumConcurrentUnindexedSearches, state.MaximumConcurrentUnindexedSearches, "maximum-concurrent-unindexed-searches")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumAttributesPerAddRequest, state.MaximumAttributesPerAddRequest, "maximum-attributes-per-add-request")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumModificationsPerModifyRequest, state.MaximumModificationsPerModifyRequest, "maximum-modifications-per-modify-request")
	operations.AddBoolOperationIfNecessary(&ops, plan.BackgroundThreadForEachPersistentSearch, state.BackgroundThreadForEachPersistentSearch, "background-thread-for-each-persistent-search")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowAttributeNameExceptions, state.AllowAttributeNameExceptions, "allow-attribute-name-exceptions")
	operations.AddStringOperationIfNecessary(&ops, plan.InvalidAttributeSyntaxBehavior, state.InvalidAttributeSyntaxBehavior, "invalid-attribute-syntax-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.SingleStructuralObjectclassBehavior, state.SingleStructuralObjectclassBehavior, "single-structural-objectclass-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.MaximumServerOutLogFileSize, state.MaximumServerOutLogFileSize, "maximum-server-out-log-file-size")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumServerOutLogFileCount, state.MaximumServerOutLogFileCount, "maximum-server-out-log-file-count")
	operations.AddStringOperationIfNecessary(&ops, plan.StartupErrorLoggerOutputLocation, state.StartupErrorLoggerOutputLocation, "startup-error-logger-output-location")
	operations.AddBoolOperationIfNecessary(&ops, plan.ExitOnJVMError, state.ExitOnJVMError, "exit-on-jvm-error")
	operations.AddStringOperationIfNecessary(&ops, plan.ResultCodeMap, state.ResultCodeMap, "result-code-map")
	operations.AddBoolOperationIfNecessary(&ops, plan.ReturnBindErrorMessages, state.ReturnBindErrorMessages, "return-bind-error-messages")
	operations.AddBoolOperationIfNecessary(&ops, plan.NotifyAbandonedOperations, state.NotifyAbandonedOperations, "notify-abandoned-operations")
	operations.AddInt64OperationIfNecessary(&ops, plan.DuplicateErrorLogLimit, state.DuplicateErrorLogLimit, "duplicate-error-log-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.DuplicateErrorLogTimeLimit, state.DuplicateErrorLogTimeLimit, "duplicate-error-log-time-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.DuplicateAlertLimit, state.DuplicateAlertLimit, "duplicate-alert-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.DuplicateAlertTimeLimit, state.DuplicateAlertTimeLimit, "duplicate-alert-time-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.WritabilityMode, state.WritabilityMode, "writability-mode")
	operations.AddStringOperationIfNecessary(&ops, plan.UnrecoverableDatabaseErrorMode, state.UnrecoverableDatabaseErrorMode, "unrecoverable-database-error")
	operations.AddBoolOperationIfNecessary(&ops, plan.DatabaseOnVirtualizedOrNetworkStorage, state.DatabaseOnVirtualizedOrNetworkStorage, "database-on-virtualized-or-network-storage")
	operations.AddStringOperationIfNecessary(&ops, plan.AutoNameWithEntryUUIDConnectionCriteria, state.AutoNameWithEntryUUIDConnectionCriteria, "auto-name-with-entry-uuid-connection-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.AutoNameWithEntryUUIDRequestCriteria, state.AutoNameWithEntryUUIDRequestCriteria, "auto-name-with-entry-uuid-request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.SoftDeletePolicy, state.SoftDeletePolicy, "soft-delete-policy")
	operations.AddStringOperationIfNecessary(&ops, plan.SubtreeAccessibilityAlertTimeLimit, state.SubtreeAccessibilityAlertTimeLimit, "subtree-accessibility-alert-time-limit")
	operations.AddBoolOperationIfNecessary(&ops, plan.WarnForBackendsWithMultipleBaseDns, state.WarnForBackendsWithMultipleBaseDns, "warn-for-backends-with-multiple-base-dns")
	operations.AddStringOperationIfNecessary(&ops, plan.ForcedGCPrimeDuration, state.ForcedGCPrimeDuration, "forced-gc-prime-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.ReplicationSetName, state.ReplicationSetName, "replication-set-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.StartupMinReplicationBacklogCount, state.StartupMinReplicationBacklogCount, "startup-min-replication-backlog-count")
	operations.AddInt64OperationIfNecessary(&ops, plan.ReplicationBacklogCountAlertThreshold, state.ReplicationBacklogCountAlertThreshold, "replication-backlog-count-alert-threshold")
	operations.AddStringOperationIfNecessary(&ops, plan.ReplicationBacklogDurationAlertThreshold, state.ReplicationBacklogDurationAlertThreshold, "replication-backlog-duration-alert-threshold")
	operations.AddStringOperationIfNecessary(&ops, plan.ReplicationAssuranceSourceTimeoutSuspendDuration, state.ReplicationAssuranceSourceTimeoutSuspendDuration, "replication-assurance-source-timeout-suspend-duration")
	operations.AddInt64OperationIfNecessary(&ops, plan.ReplicationAssuranceSourceBacklogFastStartThreshold, state.ReplicationAssuranceSourceBacklogFastStartThreshold, "replication-assurance-source-backlog-fast-start-threshold")
	operations.AddInt64OperationIfNecessary(&ops, plan.ReplicationHistoryLimit, state.ReplicationHistoryLimit, "replication-history-limit")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowInheritedReplicationOfSubordinateBackends, state.AllowInheritedReplicationOfSubordinateBackends, "allow-inherited-replication-of-subordinate-backends")
	operations.AddBoolOperationIfNecessary(&ops, plan.ReplicationPurgeObsoleteReplicas, state.ReplicationPurgeObsoleteReplicas, "replication-purge-obsolete-replicas")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxSMTPConnectionCount, state.MaxSMTPConnectionCount, "max-smtp-connection-count")
	operations.AddStringOperationIfNecessary(&ops, plan.MaxSMTPConnectionAge, state.MaxSMTPConnectionAge, "max-smtp-connection-age")
	operations.AddStringOperationIfNecessary(&ops, plan.SmtpConnectionHealthCheckInterval, state.SmtpConnectionHealthCheckInterval, "smtp-connection-health-check-interval")
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableSubOperationTimer, state.EnableSubOperationTimer, "enable-sub-operation-timer")
	operations.AddStringOperationIfNecessary(&ops, plan.MaximumShutdownTime, state.MaximumShutdownTime, "maximum-shutdown-time")
	operations.AddStringOperationIfNecessary(&ops, plan.NetworkAddressCacheTTL, state.NetworkAddressCacheTTL, "network-address-cache-ttl")
	operations.AddBoolOperationIfNecessary(&ops, plan.NetworkAddressOutageCacheEnabled, state.NetworkAddressOutageCacheEnabled, "network-address-outage-cache-enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.JmxValueBehavior, state.JmxValueBehavior, "jmx-value-behavior")
	operations.AddBoolOperationIfNecessary(&ops, plan.JmxUseLegacyMbeanNames, state.JmxUseLegacyMbeanNames, "jmx-use-legacy-mbean-names")

	// Multi-valued attributes
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SensitiveAttribute, state.SensitiveAttribute, "sensitive-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DisabledPrivilege, state.DisabledPrivilege, "disabled-privilege")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedInsecureTLSProtocol, state.AllowedInsecureTLSProtocol, "allowed-insecure-tls-protocol")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.PermitSyntaxViolationsForAttribute, state.PermitSyntaxViolationsForAttribute, "permit-syntax-violations-for-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AttributesModifiableWithIgnoreNoUserModificationRequestControl, state.AttributesModifiableWithIgnoreNoUserModificationRequestControl, "attributes-modifiable-with-ignore-no-user-modification-request-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SmtpServer, state.SmtpServer, "smtp-server")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedTask, state.AllowedTask, "allowed-task")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.TrackedApplication, state.TrackedApplication, "tracked-application")

	return ops
}

// Update the global configuration - similar to the Create method since the config is just adopted
func (r *globalConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan globalConfigurationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state
	var state globalConfigurationResourceModel
	req.State.Get(ctx, &state)
	updateGCRequest := r.apiClient.GlobalConfigurationApi.UpdateGlobalConfiguration(ProviderBasicAuthContext(ctx, r.providerConfig))

	// Determine what update operations are necessary
	ops := createGlobalConfigurationOperations(plan, state)
	if len(ops) > 0 {
		updateGCRequest = updateGCRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		globalResp, httpResp, err := r.apiClient.GlobalConfigurationApi.UpdateGlobalConfigurationExecute(updateGCRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the global configuration", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := globalResp.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGlobalConfigurationResponse(ctx, globalResp, &plan)
		// Populate Computed attribute values
		plan.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// Terraform can't actually delete the global configuration, so this method does nothing.
// Terraform will just "forget" about the global config, and it can be managed elsewhere.
func (r *globalConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *globalConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Set an arbitrary state value to appease terraform - the placeholder will immediately be
	// replaced with the actual instance name when terraform performs a read after the import.
	// If no value is set here, Terraform will error out when importing.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_name"), "placeholder")...)
}
