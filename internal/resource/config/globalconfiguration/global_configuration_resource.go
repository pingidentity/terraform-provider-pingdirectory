package globalconfiguration

import (
	"context"

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
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/planmodifiers"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
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

// Metadata returns the resource type name.
func (r *globalConfigurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_global_configuration"
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

type globalConfigurationResourceModel struct {
	Id                                                             types.String `tfsdk:"id"`
	Notifications                                                  types.Set    `tfsdk:"notifications"`
	RequiredActions                                                types.Set    `tfsdk:"required_actions"`
	Type                                                           types.String `tfsdk:"type"`
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
	UseSharedDatabaseCacheAcrossAllLocalDBBackends                 types.Bool   `tfsdk:"use_shared_database_cache_across_all_local_db_backends"`
	SharedLocalDBBackendDatabaseCachePercent                       types.Int64  `tfsdk:"shared_local_db_backend_database_cache_percent"`
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

// GetSchema defines the schema for the resource.
func (r *globalConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a Global Configuration.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Global Configuration resource. Options are ['global-configuration']",
				Optional:    false,
				Required:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"global-configuration"}...),
				},
			},
			"instance_name": schema.StringAttribute{
				Description: "Specifies a name that may be used to uniquely identify this Directory Server instance among other instances in the environment.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"location": schema.StringAttribute{
				Description: "Specifies the location for this Directory Server. Operations performed which involve communication with other servers may prefer servers in the same location to help ensure low-latency responses.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"configuration_server_group": schema.StringAttribute{
				Description: "When this property is set, changes made to this server using the console or dsconfig can be automatically applied to all servers in the specified server group.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"force_as_master_for_mirrored_data": schema.BoolAttribute{
				Description: "Indicates whether this server should be forced to assume the master role if no other suitable server is found to act as master or if multiple masters are detected. A master is only needed when changes are made to mirrored data, i.e. data specific to the topology itself and cluster-wide configuration data.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"encrypt_data": schema.BoolAttribute{
				Description: "Indicates whether the Directory Server should encrypt the data that it stores in all components that support it. This may include certain types of backends (including local DB and large attribute backends), the LDAP changelog, and the replication server database.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"encryption_settings_cipher_stream_provider": schema.StringAttribute{
				Description: "Specifies the cipher stream provider that should be used to protect the contents of the encryption settings database.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"encrypt_backups_by_default": schema.BoolAttribute{
				Description: "Indicates whether the server should encrypt backups by default.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"backup_encryption_settings_definition_id": schema.StringAttribute{
				Description: "The unique identifier for the encryption settings definition to use to generate the encryption key for encrypted backups by default.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"encrypt_ldif_exports_by_default": schema.BoolAttribute{
				Description: "Indicates whether the server should encrypt LDIF exports by default.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ldif_export_encryption_settings_definition_id": schema.StringAttribute{
				Description: "The unique identifier for the encryption settings definition to use to generate the encryption key for encrypted LDIF exports by default.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"automatically_compress_encrypted_ldif_exports": schema.BoolAttribute{
				Description: "Indicates whether to automatically compress LDIF exports that are also encrypted.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"redact_sensitive_values_in_config_logs": schema.BoolAttribute{
				Description: "Indicates whether the values of sensitive configuration properties should be redacted when logging configuration changes, including in the configuration audit log, the error log, and the server.out log file.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sensitive_attribute": schema.SetAttribute{
				Description: "Provides the ability to indicate that some attributes should be considered sensitive and additional protection should be in place when interacting with those attributes.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"reject_insecure_requests": schema.BoolAttribute{
				Description: "Indicates whether the Directory Server should reject any LDAP request (other than StartTLS) received from a client that is not using an encrypted connection.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_insecure_request_criteria": schema.StringAttribute{
				Description: "A set of criteria that may be used to match LDAP requests that may be permitted over an insecure connection even if reject-insecure-requests is true. Note that some types of requests will always be permitted, including StartTLS and start administrative session requests.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"reject_unauthenticated_requests": schema.BoolAttribute{
				Description: "Indicates whether the Directory Server should reject any LDAP request (other than bind or StartTLS requests) received from a client that has not yet been authenticated, whose last authentication attempt was unsuccessful, or whose last authentication attempt used anonymous authentication.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_unauthenticated_request_criteria": schema.StringAttribute{
				Description: "A set of criteria that may be used to match LDAP requests that may be permitted over an unauthenticated connection even if reject-unauthenticated-requests is true. Note that some types of requests will always be permitted, including bind, StartTLS, and start administrative session requests.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bind_with_dn_requires_password": schema.BoolAttribute{
				Description: "Indicates whether the Directory Server should reject any simple bind request that contains a DN but no password.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"disabled_privilege": schema.SetAttribute{
				Description: "Specifies the name of a privilege that should not be evaluated by the server.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"default_password_policy": schema.StringAttribute{
				Description: "Specifies the name of the password policy that is in effect for users whose entries do not specify an alternate password policy (either via a real or virtual attribute).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"maximum_user_data_password_policies_to_cache": schema.Int64Attribute{
				Description: "Specifies the maximum number of password policies that are defined in the user data (that is, outside of the configuration) that the server should cache in memory for faster access. A value of zero indicates that the server should not cache any user data password policies.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"proxied_authorization_identity_mapper": schema.StringAttribute{
				Description: "Specifies the name of the identity mapper to map authorization ID values (using the \"u:\" form) provided in the proxied authorization control to the corresponding user entry.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"verify_entry_digests": schema.BoolAttribute{
				Description: "Indicates whether the digest should always be verified whenever an entry containing a digest is decoded. If this is \"true\", then if a digest exists, it will always be verified. Otherwise, the digest will be written when encoding entries but ignored when decoding entries but may still be available for other verification processing.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_insecure_tls_protocol": schema.SetAttribute{
				Description: "Specifies a set of TLS protocols that will be permitted for use in the server even though there may be known vulnerabilities that could cause their use to be unsafe in some conditions. Enabling support for insecure TLS protocols is discouraged, and is generally recommended only as a short-term measure to permit legacy clients to interact with the server until they can be updated to support more secure communication protocols.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_insecure_local_jmx_connections": schema.BoolAttribute{
				Description: "Indicates that processes attaching to this server's local JVM are allowed to access internal data through JMX without the authentication requirements that remote JMX connections are subject to. Please review and understand the data that this option will expose (such as cn=monitor) to client applications to ensure there are no security concerns.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"default_internal_operation_client_connection_policy": schema.StringAttribute{
				Description: "Specifies the client connection policy that will be used by default for internal operations.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"size_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that the Directory Server should return to clients by default when processing a search operation.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"unauthenticated_size_limit": schema.Int64Attribute{
				Description: "Supported in PingDirectory product version 9.2.0.0+. The size limit value that will apply for connections from unauthenticated clients. If this is not specified, then the value of the size-limit property will be applied for both authenticated and unauthenticated connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"time_limit": schema.StringAttribute{
				Description: "Specifies the maximum length of time that the Directory Server should be allowed to spend processing a search operation.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"unauthenticated_time_limit": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 9.2.0.0+. The time limit value that will apply for connections from unauthenticated clients. If this is not specified, then the value of the time-limit property will be applied for both authenticated and unauthenticated connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"idle_time_limit": schema.StringAttribute{
				Description: "Specifies the maximum length of time that a client connection may remain established since its last completed operation.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"unauthenticated_idle_time_limit": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 9.2.0.0+. The idle-time-limit limit value that will apply for connections from unauthenticated clients. If this is not specified, then the value of the idle-time-limit property will be applied for both authenticated and unauthenticated connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"lookthrough_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that the Directory Server should \"look through\" in the course of processing a search request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"unauthenticated_lookthrough_limit": schema.Int64Attribute{
				Description: "Supported in PingDirectory product version 9.2.0.0+. The lookthrough limit value that will apply for connections from unauthenticated clients. If this is not specified, then the value of the lookthrough-limit property will be applied for both authenticated and unauthenticated connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"ldap_join_size_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that may be directly joined with any individual search result entry.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"maximum_concurrent_connections": schema.Int64Attribute{
				Description: "Specifies the maximum number of LDAP client connections which may be established to this Directory Server at the same time.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"maximum_concurrent_connections_per_ip_address": schema.Int64Attribute{
				Description: "Specifies the maximum number of LDAP client connections originating from the same IP address which may be established to this Directory Server at the same time.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"maximum_concurrent_connections_per_bind_dn": schema.Int64Attribute{
				Description: "Specifies the maximum number of LDAP client connections which may be established to this Directory Server at the same time and authenticated as the same user.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"maximum_concurrent_unindexed_searches": schema.Int64Attribute{
				Description: "Specifies the maximum number of unindexed searches that may be in progress in this backend at any given time. Any unindexed searches requested while the maximum number of unindexed searches are already being processed will be rejected. A value of zero indicates that no limit will be enforced.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"maximum_attributes_per_add_request": schema.Int64Attribute{
				Description: "Specifies the maximum number of attributes that may be included in an add request. This property does not impose any limit on the number of values that an attribute may have.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"maximum_modifications_per_modify_request": schema.Int64Attribute{
				Description: "Specifies the maximum number of modifications that may be included in a modify request. This property does not impose any limit on the number of attribute values that a modification may have.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"background_thread_for_each_persistent_search": schema.BoolAttribute{
				Description: "Indicates whether the server should use a separate background thread for each persistent search.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_attribute_name_exceptions": schema.BoolAttribute{
				Description: "Indicates whether the Directory Server should allow underscores in attribute names and allow attribute names to begin with numeric digits (both of which are violations of the LDAP standards).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"invalid_attribute_syntax_behavior": schema.StringAttribute{
				Description: "Specifies how the Directory Server should handle operations whenever an attribute value violates the associated attribute syntax.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"accept", "reject", "warn"}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					planmodifiers.ToLowercasePlanModifier(),
				},
			},
			"permit_syntax_violations_for_attribute": schema.SetAttribute{
				Description: "Specifies a set of attribute types for which the server will permit values that do not conform to the associated attribute syntax.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"single_structural_objectclass_behavior": schema.StringAttribute{
				Description: "Specifies how the Directory Server should handle operations for an entry does not contain a structural object class, or for an entry that contains multiple structural classes.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"accept", "reject", "warn"}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					planmodifiers.ToLowercasePlanModifier(),
				},
			},
			"attributes_modifiable_with_ignore_no_user_modification_request_control": schema.SetAttribute{
				Description: "Specifies the operational attribute types that are defined in the schema with the NO-USER-MODIFICATION constraint that the server will allow to be altered if the associated request contains the ignore NO-USER-MODIFICATION request control.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"maximum_server_out_log_file_size": schema.StringAttribute{
				Description: "The maximum allowed size that the server.out log file will be allowed to have. If a write would cause the file to exceed this size, then the current file will be rotated out of place and a new empty file will be created and the message written to it.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"maximum_server_out_log_file_count": schema.Int64Attribute{
				Description: "The maximum number of server.out log files (including the current active log file) that should be retained. When rotating the log file, if the total number of files exceeds this count, then the oldest file(s) will be removed so that the total number of log files is within this limit.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"startup_error_logger_output_location": schema.StringAttribute{
				Description: "Specifies how the server should handle error log messages (which may include errors, warnings, and notices) generated during startup. All of these messages will be written to all configured error loggers, but they may also be written to other locations (like standard output, standard error, or the server.out log file) so that they are displayed on the console when the server is starting.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"standard_output", "standard_error", "server_out_file", "standard_output_and_server_out_file", "standard_error_and_server_out_file", "disabled"}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					planmodifiers.ToLowercasePlanModifier(),
				},
			},
			"exit_on_jvm_error": schema.BoolAttribute{
				Description: "Indicates whether the Directory Server should be shut down if a severe error is raised (e.g., an out of memory error) which may prevent the JVM from continuing to run properly.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"server_error_result_code": schema.Int64Attribute{
				Description: "Specifies the numeric value of the result code when request processing fails due to an internal server error.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"result_code_map": schema.StringAttribute{
				Description: "Specifies a result code map that should be used for clients that do not have a map associated with their client connection policy. If the associated client connection policy has a result code map, then that map will be used instead. If no map is associated either with the client connection policy or the global configuration, then an internal default will be used.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"return_bind_error_messages": schema.BoolAttribute{
				Description: "Indicates whether responses for failed bind operations should include a message string providing the reason for the authentication failure.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"notify_abandoned_operations": schema.BoolAttribute{
				Description: "Indicates whether the Directory Server should send a response to any operation that is interrupted via an abandon request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"duplicate_error_log_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of duplicate error log messages that should be logged in the time window specified by the duplicate-error-log-time-limit property.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"duplicate_error_log_time_limit": schema.StringAttribute{
				Description: "Specifies the length of time that must expire before duplicate log messages above the duplicate-error-log-limit threshold are logged again to the error log.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"duplicate_alert_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of duplicate alert messages that should be sent via the administrative alert framework in the time window specified by the duplicate-alert-time-limit property.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"duplicate_alert_time_limit": schema.StringAttribute{
				Description: "Specifies the length of time that must expire before duplicate messages are sent via the administrative alert framework.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"writability_mode": schema.StringAttribute{
				Description: "Specifies the kinds of write operations the Directory Server can process.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"enabled", "disabled", "internal_only"}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					planmodifiers.ToLowercasePlanModifier(),
				},
			},
			"use_shared_database_cache_across_all_local_db_backends": schema.BoolAttribute{
				Description: "Supported in PingDirectory product version 10.0.0.0+. Indicates whether the server should use a common database cache that is shared across all local DB backends instead of maintaining a separate cache for each backend.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"shared_local_db_backend_database_cache_percent": schema.Int64Attribute{
				Description: "Supported in PingDirectory product version 10.0.0.0+. Specifies the percentage of the JVM memory to allocate to the database cache that is shared across all local DB backends.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"unrecoverable_database_error_mode": schema.StringAttribute{
				Description: "Specifies the action which should be taken for any database that experiences an unrecoverable error. Action applies to local database backends and the replication recent changes database.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"enter_lockdown_mode", "raise_unavailable_alarm", "initiate_server_shutdown"}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					planmodifiers.ToLowercasePlanModifier(),
				},
			},
			"database_on_virtualized_or_network_storage": schema.BoolAttribute{
				Description: "This setting provides data integrity options when the Directory Server is installed with a database on a network storage device. A storage device may be accessed directly by a physical server, or indirectly through a virtual machine running on a hypervisor. Enabling this setting will apply changes to all Local DB Backends, the LDAP Changelog Backend, and the replication changelog database.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"auto_name_with_entry_uuid_connection_criteria": schema.StringAttribute{
				Description: "Connection criteria that may be used to identify clients whose add requests should use entryUUID as the naming attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"auto_name_with_entry_uuid_request_criteria": schema.StringAttribute{
				Description: "Request criteria that may be used to identify add requests that should use entryUUID as the naming attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"soft_delete_policy": schema.StringAttribute{
				Description: "Specifies the soft delete policy that will be used by default for delete operations. Soft delete operations introduce the ability to control the server behavior of the delete operation. Instead of performing a permanent delete of an entry, deleted entries can be retained as soft deleted entries by their entryUUID values and are available for undelete at a later time. In addition to a soft delete policy enabling soft deletes, delete operations sent to the server must have the soft delete request control present with sufficient access privileges to access the soft delete request control.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subtree_accessibility_alert_time_limit": schema.StringAttribute{
				Description: "Specifies the length of time that a subtree may remain hidden or read-only before an administrative alert is sent.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"warn_for_backends_with_multiple_base_dns": schema.BoolAttribute{
				Description: "Indicates whether the server should issue a warning when enabling a backend that contains multiple base DNs.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"forced_gc_prime_duration": schema.StringAttribute{
				Description: "Specifies the minimum length of time required for backend or request processor initialization that will trigger the server to force an explicit garbage collection. A value of \"0 seconds\" indicates that the server should never invoke an explicit garbage collection regardless of the length of time required to initialize the server backends.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"replication_set_name": schema.StringAttribute{
				Description: "The name of the replication set assigned to this Directory Server. Restricted domains are only replicated within instances using the same replication set name.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"startup_min_replication_backlog_count": schema.Int64Attribute{
				Description: "The number of outstanding changes any replica can have before the Directory Server will start accepting connections. The Directory Server may never accept connections if this setting is too low. If you are unsure which value to use, you can use the number of expected updates within a five second interval.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"replication_backlog_count_alert_threshold": schema.Int64Attribute{
				Description: "An alert is sent when the number of outstanding replication changes for the Directory Server has exceeded this threshold for longer than the replication backlog duration alert threshold.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"replication_backlog_duration_alert_threshold": schema.StringAttribute{
				Description: "An alert is sent when the number of outstanding replication changes for the Directory Server has exceeded the replication backlog count alert threshold for longer than this duration.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"replication_assurance_source_timeout_suspend_duration": schema.StringAttribute{
				Description: "The amount of time a replication assurance source (i.e. a peer Directory Server) will be suspended from assurance requirements on this Directory Server if it experiences an assurance timeout.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"replication_assurance_source_backlog_fast_start_threshold": schema.Int64Attribute{
				Description: "The maximum number of replication backlog updates a replication assurance source (i.e. a peer Directory Server) can have and be immediately recognized as an available assurance source by this Directory Server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"replication_history_limit": schema.Int64Attribute{
				Description: "Specifies the size limit for historical information.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"allow_inherited_replication_of_subordinate_backends": schema.BoolAttribute{
				Description: "Allow replication to be inherited by subordinate/child backends.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"replication_purge_obsolete_replicas": schema.BoolAttribute{
				Description: "Indicates whether state about obsolete replicas is automatically purged.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"smtp_server": schema.SetAttribute{
				Description: "Specifies the set of servers that will be used to send email messages. The order in which the servers are listed indicates the order in which the Directory Server will attempt to use them in the course of sending a message. The first attempt will always go to the server at the top of the list, and servers further down the list will only be used if none of the servers listed above it were able to successfully send the message.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"max_smtp_connection_count": schema.Int64Attribute{
				Description: "The maximum number of SMTP connections that will be maintained for delivering email messages.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"max_smtp_connection_age": schema.StringAttribute{
				Description: "The maximum length of time that a connection to an SMTP server should be considered valid.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"smtp_connection_health_check_interval": schema.StringAttribute{
				Description: "The length of time between checks to ensure that available SMTP connections are still valid.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_task": schema.SetAttribute{
				Description: "Specifies the fully-qualified name of a Java class that may be invoked in the server.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"enable_sub_operation_timer": schema.BoolAttribute{
				Description: "Indicates whether the Directory Server should attempt to record information about the length of time required to process various phases of an operation. Enabling this feature may impact performance, but could make it easier to identify potential bottlenecks in operation processing.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"maximum_shutdown_time": schema.StringAttribute{
				Description: "Specifies the maximum amount of time the shutdown of Directory Server may take.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network_address_cache_ttl": schema.StringAttribute{
				Description: "Specifies the length of time that the Directory Server should cache the IP addresses associated with the names of systems with which it interacts.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network_address_outage_cache_enabled": schema.BoolAttribute{
				Description: "Specifies whether the Directory Server should cache the last valid IP addresses associated with the names of systems with which it interacts with when the domain name service returns an unknown host exception. Java may return an unknown host exception when there is unexpected interruption in domain name service so this setting protects the Directory Server from temporary DNS server outages if previous results have been cached.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"tracked_application": schema.SetAttribute{
				Description: "Specifies criteria for identifying specific applications that access the server to enable tracking throughput and latency of LDAP operations issued by an application.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"jmx_value_behavior": schema.StringAttribute{
				Description: "Specifies how a Java type is chosen for monitor attributes exposed as JMX attribute values.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"inferred", "string"}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					planmodifiers.ToLowercasePlanModifier(),
				},
			},
			"jmx_use_legacy_mbean_names": schema.BoolAttribute{
				Description: "When set to true, the server will use its original, non-standard JMX MBean names for the monitoring MBeans. These include RDN keys of \"Rdn1\" and \"Rdn2\" instead of the recommended \"type\" and \"name\" keys. This should option should only be enabled for installations that have monitoring infrastructure that depends on the old keys.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *globalConfigurationResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingDirectory10000)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model globalConfigurationResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.UseSharedDatabaseCacheAcrossAllLocalDBBackends) {
		resp.Diagnostics.AddError("Attribute 'use_shared_database_cache_across_all_local_db_backends' not supported by PingDirectory version "+r.providerConfig.ProductVersion, "")
	}
	if internaltypes.IsDefined(model.SharedLocalDBBackendDatabaseCachePercent) {
		resp.Diagnostics.AddError("Attribute 'shared_local_db_backend_database_cache_percent' not supported by PingDirectory version "+r.providerConfig.ProductVersion, "")
	}
	compare, err = version.Compare(r.providerConfig.ProductVersion, version.PingDirectory9200)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	if internaltypes.IsDefined(model.UnauthenticatedSizeLimit) {
		resp.Diagnostics.AddError("Attribute 'unauthenticated_size_limit' not supported by PingDirectory version "+r.providerConfig.ProductVersion, "")
	}
	if internaltypes.IsNonEmptyString(model.UnauthenticatedTimeLimit) {
		resp.Diagnostics.AddError("Attribute 'unauthenticated_time_limit' not supported by PingDirectory version "+r.providerConfig.ProductVersion, "")
	}
	if internaltypes.IsNonEmptyString(model.UnauthenticatedIdleTimeLimit) {
		resp.Diagnostics.AddError("Attribute 'unauthenticated_idle_time_limit' not supported by PingDirectory version "+r.providerConfig.ProductVersion, "")
	}
	if internaltypes.IsDefined(model.UnauthenticatedLookthroughLimit) {
		resp.Diagnostics.AddError("Attribute 'unauthenticated_lookthrough_limit' not supported by PingDirectory version "+r.providerConfig.ProductVersion, "")
	}
}

// Read a GlobalConfigurationResponse object into the model struct
func readGlobalConfigurationResponse(ctx context.Context, r *client.GlobalConfigurationResponse, state *globalConfigurationResourceModel, expectedValues *globalConfigurationResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("global-configuration")
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
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
	state.DisabledPrivilege = internaltypes.GetStringSet(
		client.StringSliceEnumglobalConfigurationDisabledPrivilegeProp(r.DisabledPrivilege))
	state.DefaultPasswordPolicy = types.StringValue(r.DefaultPasswordPolicy)
	state.MaximumUserDataPasswordPoliciesToCache = internaltypes.Int64TypeOrNil(r.MaximumUserDataPasswordPoliciesToCache)
	state.ProxiedAuthorizationIdentityMapper = types.StringValue(r.ProxiedAuthorizationIdentityMapper)
	state.VerifyEntryDigests = internaltypes.BoolTypeOrNil(r.VerifyEntryDigests)
	state.AllowedInsecureTLSProtocol = internaltypes.GetStringSet(
		client.StringSliceEnumglobalConfigurationAllowedInsecureTLSProtocolProp(r.AllowedInsecureTLSProtocol))
	state.AllowInsecureLocalJMXConnections = internaltypes.BoolTypeOrNil(r.AllowInsecureLocalJMXConnections)
	state.DefaultInternalOperationClientConnectionPolicy = internaltypes.StringTypeOrNil(r.DefaultInternalOperationClientConnectionPolicy, true)
	state.SizeLimit = internaltypes.Int64TypeOrNil(r.SizeLimit)
	state.UnauthenticatedSizeLimit = internaltypes.Int64TypeOrNil(r.UnauthenticatedSizeLimit)
	state.TimeLimit = internaltypes.StringTypeOrNil(r.TimeLimit, true)
	config.CheckMismatchedPDFormattedAttributes("time_limit",
		expectedValues.TimeLimit, state.TimeLimit, diagnostics)
	state.UnauthenticatedTimeLimit = internaltypes.StringTypeOrNil(r.UnauthenticatedTimeLimit, true)
	config.CheckMismatchedPDFormattedAttributes("unauthenticated_time_limit",
		expectedValues.UnauthenticatedTimeLimit, state.UnauthenticatedTimeLimit, diagnostics)
	state.IdleTimeLimit = internaltypes.StringTypeOrNil(r.IdleTimeLimit, true)
	config.CheckMismatchedPDFormattedAttributes("idle_time_limit",
		expectedValues.IdleTimeLimit, state.IdleTimeLimit, diagnostics)
	state.UnauthenticatedIdleTimeLimit = internaltypes.StringTypeOrNil(r.UnauthenticatedIdleTimeLimit, true)
	config.CheckMismatchedPDFormattedAttributes("unauthenticated_idle_time_limit",
		expectedValues.UnauthenticatedIdleTimeLimit, state.UnauthenticatedIdleTimeLimit, diagnostics)
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
		client.StringPointerEnumglobalConfigurationInvalidAttributeSyntaxBehaviorProp(r.InvalidAttributeSyntaxBehavior), true)
	state.PermitSyntaxViolationsForAttribute = internaltypes.GetStringSet(r.PermitSyntaxViolationsForAttribute)
	state.SingleStructuralObjectclassBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumglobalConfigurationSingleStructuralObjectclassBehaviorProp(r.SingleStructuralObjectclassBehavior), true)
	state.AttributesModifiableWithIgnoreNoUserModificationRequestControl = internaltypes.GetStringSet(
		client.StringSliceEnumglobalConfigurationAttributesModifiableWithIgnoreNoUserModificationRequestControlProp(r.AttributesModifiableWithIgnoreNoUserModificationRequestControl))
	state.MaximumServerOutLogFileSize = internaltypes.StringTypeOrNil(r.MaximumServerOutLogFileSize, true)
	config.CheckMismatchedPDFormattedAttributes("maximum_server_out_log_file_size",
		expectedValues.MaximumServerOutLogFileSize, state.MaximumServerOutLogFileSize, diagnostics)
	state.MaximumServerOutLogFileCount = internaltypes.Int64TypeOrNil(r.MaximumServerOutLogFileCount)
	state.StartupErrorLoggerOutputLocation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumglobalConfigurationStartupErrorLoggerOutputLocationProp(r.StartupErrorLoggerOutputLocation), true)
	state.ExitOnJVMError = internaltypes.BoolTypeOrNil(r.ExitOnJVMError)
	state.ServerErrorResultCode = internaltypes.Int64TypeOrNil(r.ServerErrorResultCode)
	state.ResultCodeMap = internaltypes.StringTypeOrNil(r.ResultCodeMap, true)
	state.ReturnBindErrorMessages = internaltypes.BoolTypeOrNil(r.ReturnBindErrorMessages)
	state.NotifyAbandonedOperations = internaltypes.BoolTypeOrNil(r.NotifyAbandonedOperations)
	state.DuplicateErrorLogLimit = types.Int64Value(r.DuplicateErrorLogLimit)
	state.DuplicateErrorLogTimeLimit = types.StringValue(r.DuplicateErrorLogTimeLimit)
	config.CheckMismatchedPDFormattedAttributes("duplicate_error_log_time_limit",
		expectedValues.DuplicateErrorLogTimeLimit, state.DuplicateErrorLogTimeLimit, diagnostics)
	state.DuplicateAlertLimit = types.Int64Value(r.DuplicateAlertLimit)
	state.DuplicateAlertTimeLimit = types.StringValue(r.DuplicateAlertTimeLimit)
	config.CheckMismatchedPDFormattedAttributes("duplicate_alert_time_limit",
		expectedValues.DuplicateAlertTimeLimit, state.DuplicateAlertTimeLimit, diagnostics)
	state.WritabilityMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumglobalConfigurationWritabilityModeProp(r.WritabilityMode), true)
	state.UseSharedDatabaseCacheAcrossAllLocalDBBackends = internaltypes.BoolTypeOrNil(r.UseSharedDatabaseCacheAcrossAllLocalDBBackends)
	state.SharedLocalDBBackendDatabaseCachePercent = internaltypes.Int64TypeOrNil(r.SharedLocalDBBackendDatabaseCachePercent)
	state.UnrecoverableDatabaseErrorMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumglobalConfigurationUnrecoverableDatabaseErrorModeProp(r.UnrecoverableDatabaseErrorMode), true)
	state.DatabaseOnVirtualizedOrNetworkStorage = internaltypes.BoolTypeOrNil(r.DatabaseOnVirtualizedOrNetworkStorage)
	state.AutoNameWithEntryUUIDConnectionCriteria = internaltypes.StringTypeOrNil(r.AutoNameWithEntryUUIDConnectionCriteria, true)
	state.AutoNameWithEntryUUIDRequestCriteria = internaltypes.StringTypeOrNil(r.AutoNameWithEntryUUIDRequestCriteria, true)
	state.SoftDeletePolicy = internaltypes.StringTypeOrNil(r.SoftDeletePolicy, true)
	state.SubtreeAccessibilityAlertTimeLimit = internaltypes.StringTypeOrNil(r.SubtreeAccessibilityAlertTimeLimit, true)
	config.CheckMismatchedPDFormattedAttributes("subtree_accessibility_alert_time_limit",
		expectedValues.SubtreeAccessibilityAlertTimeLimit, state.SubtreeAccessibilityAlertTimeLimit, diagnostics)
	state.WarnForBackendsWithMultipleBaseDns = internaltypes.BoolTypeOrNil(r.WarnForBackendsWithMultipleBaseDns)
	state.ForcedGCPrimeDuration = internaltypes.StringTypeOrNil(r.ForcedGCPrimeDuration, true)
	config.CheckMismatchedPDFormattedAttributes("forced_gc_prime_duration",
		expectedValues.ForcedGCPrimeDuration, state.ForcedGCPrimeDuration, diagnostics)
	state.ReplicationSetName = internaltypes.StringTypeOrNil(r.ReplicationSetName, true)
	state.StartupMinReplicationBacklogCount = types.Int64Value(r.StartupMinReplicationBacklogCount)
	state.ReplicationBacklogCountAlertThreshold = types.Int64Value(r.ReplicationBacklogCountAlertThreshold)
	state.ReplicationBacklogDurationAlertThreshold = types.StringValue(r.ReplicationBacklogDurationAlertThreshold)
	config.CheckMismatchedPDFormattedAttributes("replication_backlog_duration_alert_threshold",
		expectedValues.ReplicationBacklogDurationAlertThreshold, state.ReplicationBacklogDurationAlertThreshold, diagnostics)
	state.ReplicationAssuranceSourceTimeoutSuspendDuration = types.StringValue(r.ReplicationAssuranceSourceTimeoutSuspendDuration)
	config.CheckMismatchedPDFormattedAttributes("replication_assurance_source_timeout_suspend_duration",
		expectedValues.ReplicationAssuranceSourceTimeoutSuspendDuration, state.ReplicationAssuranceSourceTimeoutSuspendDuration, diagnostics)
	state.ReplicationAssuranceSourceBacklogFastStartThreshold = types.Int64Value(r.ReplicationAssuranceSourceBacklogFastStartThreshold)
	state.ReplicationHistoryLimit = internaltypes.Int64TypeOrNil(r.ReplicationHistoryLimit)
	state.AllowInheritedReplicationOfSubordinateBackends = types.BoolValue(r.AllowInheritedReplicationOfSubordinateBackends)
	state.ReplicationPurgeObsoleteReplicas = internaltypes.BoolTypeOrNil(r.ReplicationPurgeObsoleteReplicas)
	state.SmtpServer = internaltypes.GetStringSet(r.SmtpServer)
	state.MaxSMTPConnectionCount = internaltypes.Int64TypeOrNil(r.MaxSMTPConnectionCount)
	state.MaxSMTPConnectionAge = internaltypes.StringTypeOrNil(r.MaxSMTPConnectionAge, true)
	config.CheckMismatchedPDFormattedAttributes("max_smtp_connection_age",
		expectedValues.MaxSMTPConnectionAge, state.MaxSMTPConnectionAge, diagnostics)
	state.SmtpConnectionHealthCheckInterval = internaltypes.StringTypeOrNil(r.SmtpConnectionHealthCheckInterval, true)
	config.CheckMismatchedPDFormattedAttributes("smtp_connection_health_check_interval",
		expectedValues.SmtpConnectionHealthCheckInterval, state.SmtpConnectionHealthCheckInterval, diagnostics)
	state.AllowedTask = internaltypes.GetStringSet(r.AllowedTask)
	state.EnableSubOperationTimer = internaltypes.BoolTypeOrNil(r.EnableSubOperationTimer)
	state.MaximumShutdownTime = internaltypes.StringTypeOrNil(r.MaximumShutdownTime, true)
	config.CheckMismatchedPDFormattedAttributes("maximum_shutdown_time",
		expectedValues.MaximumShutdownTime, state.MaximumShutdownTime, diagnostics)
	state.NetworkAddressCacheTTL = internaltypes.StringTypeOrNil(r.NetworkAddressCacheTTL, true)
	config.CheckMismatchedPDFormattedAttributes("network_address_cache_ttl",
		expectedValues.NetworkAddressCacheTTL, state.NetworkAddressCacheTTL, diagnostics)
	state.NetworkAddressOutageCacheEnabled = internaltypes.BoolTypeOrNil(r.NetworkAddressOutageCacheEnabled)
	state.TrackedApplication = internaltypes.GetStringSet(r.TrackedApplication)
	state.JmxValueBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumglobalConfigurationJmxValueBehaviorProp(r.JmxValueBehavior), true)
	state.JmxUseLegacyMbeanNames = internaltypes.BoolTypeOrNil(r.JmxUseLegacyMbeanNames)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createGlobalConfigurationOperations(plan globalConfigurationResourceModel, state globalConfigurationResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.InstanceName, state.InstanceName, "instance-name")
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
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SensitiveAttribute, state.SensitiveAttribute, "sensitive-attribute")
	operations.AddBoolOperationIfNecessary(&ops, plan.RejectInsecureRequests, state.RejectInsecureRequests, "reject-insecure-requests")
	operations.AddStringOperationIfNecessary(&ops, plan.AllowedInsecureRequestCriteria, state.AllowedInsecureRequestCriteria, "allowed-insecure-request-criteria")
	operations.AddBoolOperationIfNecessary(&ops, plan.RejectUnauthenticatedRequests, state.RejectUnauthenticatedRequests, "reject-unauthenticated-requests")
	operations.AddStringOperationIfNecessary(&ops, plan.AllowedUnauthenticatedRequestCriteria, state.AllowedUnauthenticatedRequestCriteria, "allowed-unauthenticated-request-criteria")
	operations.AddBoolOperationIfNecessary(&ops, plan.BindWithDNRequiresPassword, state.BindWithDNRequiresPassword, "bind-with-dn-requires-password")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DisabledPrivilege, state.DisabledPrivilege, "disabled-privilege")
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultPasswordPolicy, state.DefaultPasswordPolicy, "default-password-policy")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumUserDataPasswordPoliciesToCache, state.MaximumUserDataPasswordPoliciesToCache, "maximum-user-data-password-policies-to-cache")
	operations.AddStringOperationIfNecessary(&ops, plan.ProxiedAuthorizationIdentityMapper, state.ProxiedAuthorizationIdentityMapper, "proxied-authorization-identity-mapper")
	operations.AddBoolOperationIfNecessary(&ops, plan.VerifyEntryDigests, state.VerifyEntryDigests, "verify-entry-digests")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedInsecureTLSProtocol, state.AllowedInsecureTLSProtocol, "allowed-insecure-tls-protocol")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowInsecureLocalJMXConnections, state.AllowInsecureLocalJMXConnections, "allow-insecure-local-jmx-connections")
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultInternalOperationClientConnectionPolicy, state.DefaultInternalOperationClientConnectionPolicy, "default-internal-operation-client-connection-policy")
	operations.AddInt64OperationIfNecessary(&ops, plan.SizeLimit, state.SizeLimit, "size-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.UnauthenticatedSizeLimit, state.UnauthenticatedSizeLimit, "unauthenticated-size-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.TimeLimit, state.TimeLimit, "time-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.UnauthenticatedTimeLimit, state.UnauthenticatedTimeLimit, "unauthenticated-time-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.IdleTimeLimit, state.IdleTimeLimit, "idle-time-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.UnauthenticatedIdleTimeLimit, state.UnauthenticatedIdleTimeLimit, "unauthenticated-idle-time-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.LookthroughLimit, state.LookthroughLimit, "lookthrough-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.UnauthenticatedLookthroughLimit, state.UnauthenticatedLookthroughLimit, "unauthenticated-lookthrough-limit")
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
	operations.AddStringSetOperationsIfNecessary(&ops, plan.PermitSyntaxViolationsForAttribute, state.PermitSyntaxViolationsForAttribute, "permit-syntax-violations-for-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.SingleStructuralObjectclassBehavior, state.SingleStructuralObjectclassBehavior, "single-structural-objectclass-behavior")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AttributesModifiableWithIgnoreNoUserModificationRequestControl, state.AttributesModifiableWithIgnoreNoUserModificationRequestControl, "attributes-modifiable-with-ignore-no-user-modification-request-control")
	operations.AddStringOperationIfNecessary(&ops, plan.MaximumServerOutLogFileSize, state.MaximumServerOutLogFileSize, "maximum-server-out-log-file-size")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumServerOutLogFileCount, state.MaximumServerOutLogFileCount, "maximum-server-out-log-file-count")
	operations.AddStringOperationIfNecessary(&ops, plan.StartupErrorLoggerOutputLocation, state.StartupErrorLoggerOutputLocation, "startup-error-logger-output-location")
	operations.AddBoolOperationIfNecessary(&ops, plan.ExitOnJVMError, state.ExitOnJVMError, "exit-on-jvm-error")
	operations.AddInt64OperationIfNecessary(&ops, plan.ServerErrorResultCode, state.ServerErrorResultCode, "server-error-result-code")
	operations.AddStringOperationIfNecessary(&ops, plan.ResultCodeMap, state.ResultCodeMap, "result-code-map")
	operations.AddBoolOperationIfNecessary(&ops, plan.ReturnBindErrorMessages, state.ReturnBindErrorMessages, "return-bind-error-messages")
	operations.AddBoolOperationIfNecessary(&ops, plan.NotifyAbandonedOperations, state.NotifyAbandonedOperations, "notify-abandoned-operations")
	operations.AddInt64OperationIfNecessary(&ops, plan.DuplicateErrorLogLimit, state.DuplicateErrorLogLimit, "duplicate-error-log-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.DuplicateErrorLogTimeLimit, state.DuplicateErrorLogTimeLimit, "duplicate-error-log-time-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.DuplicateAlertLimit, state.DuplicateAlertLimit, "duplicate-alert-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.DuplicateAlertTimeLimit, state.DuplicateAlertTimeLimit, "duplicate-alert-time-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.WritabilityMode, state.WritabilityMode, "writability-mode")
	operations.AddBoolOperationIfNecessary(&ops, plan.UseSharedDatabaseCacheAcrossAllLocalDBBackends, state.UseSharedDatabaseCacheAcrossAllLocalDBBackends, "use-shared-database-cache-across-all-local-db-backends")
	operations.AddInt64OperationIfNecessary(&ops, plan.SharedLocalDBBackendDatabaseCachePercent, state.SharedLocalDBBackendDatabaseCachePercent, "shared-local-db-backend-database-cache-percent")
	operations.AddStringOperationIfNecessary(&ops, plan.UnrecoverableDatabaseErrorMode, state.UnrecoverableDatabaseErrorMode, "unrecoverable-database-error-mode")
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
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SmtpServer, state.SmtpServer, "smtp-server")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxSMTPConnectionCount, state.MaxSMTPConnectionCount, "max-smtp-connection-count")
	operations.AddStringOperationIfNecessary(&ops, plan.MaxSMTPConnectionAge, state.MaxSMTPConnectionAge, "max-smtp-connection-age")
	operations.AddStringOperationIfNecessary(&ops, plan.SmtpConnectionHealthCheckInterval, state.SmtpConnectionHealthCheckInterval, "smtp-connection-health-check-interval")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedTask, state.AllowedTask, "allowed-task")
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableSubOperationTimer, state.EnableSubOperationTimer, "enable-sub-operation-timer")
	operations.AddStringOperationIfNecessary(&ops, plan.MaximumShutdownTime, state.MaximumShutdownTime, "maximum-shutdown-time")
	operations.AddStringOperationIfNecessary(&ops, plan.NetworkAddressCacheTTL, state.NetworkAddressCacheTTL, "network-address-cache-ttl")
	operations.AddBoolOperationIfNecessary(&ops, plan.NetworkAddressOutageCacheEnabled, state.NetworkAddressOutageCacheEnabled, "network-address-outage-cache-enabled")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.TrackedApplication, state.TrackedApplication, "tracked-application")
	operations.AddStringOperationIfNecessary(&ops, plan.JmxValueBehavior, state.JmxValueBehavior, "jmx-value-behavior")
	operations.AddBoolOperationIfNecessary(&ops, plan.JmxUseLegacyMbeanNames, state.JmxUseLegacyMbeanNames, "jmx-use-legacy-mbean-names")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *globalConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan globalConfigurationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.GlobalConfigurationAPI.GetGlobalConfiguration(
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

	// Read the existing configuration
	var state globalConfigurationResourceModel
	readGlobalConfigurationResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.GlobalConfigurationAPI.UpdateGlobalConfiguration(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	ops := createGlobalConfigurationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.GlobalConfigurationAPI.UpdateGlobalConfigurationExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Global Configuration", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGlobalConfigurationResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
	}

	diags = resp.State.Set(ctx, state)
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

	readResponse, httpResp, err := r.apiClient.GlobalConfigurationAPI.GetGlobalConfiguration(
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
	readGlobalConfigurationResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *globalConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan globalConfigurationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state globalConfigurationResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.GlobalConfigurationAPI.UpdateGlobalConfiguration(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))

	// Determine what update operations are necessary
	ops := createGlobalConfigurationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.GlobalConfigurationAPI.UpdateGlobalConfigurationExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Global Configuration", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGlobalConfigurationResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *globalConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *globalConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Set a placeholder id value to appease terraform.
	// The real attributes will be imported when terraform performs a read after the import.
	// If no value is set here, Terraform will error out when importing.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), "id")...)
}
