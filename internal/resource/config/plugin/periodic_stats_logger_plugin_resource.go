package plugin

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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &periodicStatsLoggerPluginResource{}
	_ resource.ResourceWithConfigure   = &periodicStatsLoggerPluginResource{}
	_ resource.ResourceWithImportState = &periodicStatsLoggerPluginResource{}
)

// Create a Periodic Stats Logger Plugin resource
func NewPeriodicStatsLoggerPluginResource() resource.Resource {
	return &periodicStatsLoggerPluginResource{}
}

// periodicStatsLoggerPluginResource is the resource implementation.
type periodicStatsLoggerPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *periodicStatsLoggerPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_periodic_stats_logger_plugin"
}

// Configure adds the provider configured client to the resource.
func (r *periodicStatsLoggerPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type periodicStatsLoggerPluginResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	LastUpdated             types.String `tfsdk:"last_updated"`
	Notifications           types.Set    `tfsdk:"notifications"`
	RequiredActions         types.Set    `tfsdk:"required_actions"`
	LogInterval             types.String `tfsdk:"log_interval"`
	CollectionInterval      types.String `tfsdk:"collection_interval"`
	SuppressIfIdle          types.Bool   `tfsdk:"suppress_if_idle"`
	HeaderPrefixPerColumn   types.Bool   `tfsdk:"header_prefix_per_column"`
	EmptyInsteadOfZero      types.Bool   `tfsdk:"empty_instead_of_zero"`
	LinesBetweenHeader      types.Int64  `tfsdk:"lines_between_header"`
	IncludedLDAPStat        types.Set    `tfsdk:"included_ldap_stat"`
	IncludedResourceStat    types.Set    `tfsdk:"included_resource_stat"`
	HistogramFormat         types.String `tfsdk:"histogram_format"`
	HistogramOpType         types.Set    `tfsdk:"histogram_op_type"`
	PerApplicationLDAPStats types.String `tfsdk:"per_application_ldap_stats"`
	StatusSummaryInfo       types.String `tfsdk:"status_summary_info"`
	LdapChangelogInfo       types.String `tfsdk:"ldap_changelog_info"`
	GaugeInfo               types.String `tfsdk:"gauge_info"`
	LogFileFormat           types.String `tfsdk:"log_file_format"`
	LogFile                 types.String `tfsdk:"log_file"`
	LogFilePermissions      types.String `tfsdk:"log_file_permissions"`
	Append                  types.Bool   `tfsdk:"append"`
	RotationPolicy          types.Set    `tfsdk:"rotation_policy"`
	RotationListener        types.Set    `tfsdk:"rotation_listener"`
	RetentionPolicy         types.Set    `tfsdk:"retention_policy"`
	LoggingErrorBehavior    types.String `tfsdk:"logging_error_behavior"`
	LocalDBBackendInfo      types.String `tfsdk:"local_db_backend_info"`
	ReplicationInfo         types.String `tfsdk:"replication_info"`
	EntryCacheInfo          types.String `tfsdk:"entry_cache_info"`
	HostInfo                types.Set    `tfsdk:"host_info"`
	IncludedLDAPApplication types.Set    `tfsdk:"included_ldap_application"`
	Description             types.String `tfsdk:"description"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *periodicStatsLoggerPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Periodic Stats Logger Plugin.",
		Attributes: map[string]schema.Attribute{
			"log_interval": schema.StringAttribute{
				Description: "The duration between statistics collection and logging. A new line is logged to the output for each interval. Setting this value too small can have an impact on performance.",
				Optional:    true,
				Computed:    true,
			},
			"collection_interval": schema.StringAttribute{
				Description: "Some of the calculated statistics, such as the average and maximum queue sizes, can use multiple samples within a log interval. This value controls how often samples are gathered. It should be a multiple of the log-interval.",
				Optional:    true,
				Computed:    true,
			},
			"suppress_if_idle": schema.BoolAttribute{
				Description: "If the server is idle during the specified interval, then do not log any output if this property is set to true. The server is idle if during the interval, no new connections were established, no operations were processed, and no operations are pending.",
				Optional:    true,
				Computed:    true,
			},
			"header_prefix_per_column": schema.BoolAttribute{
				Description: "This property controls whether the header prefix, which applies to a group of columns, appears at the start of each column header or only the first column in a group.",
				Optional:    true,
				Computed:    true,
			},
			"empty_instead_of_zero": schema.BoolAttribute{
				Description: "This property controls whether a value in the output is shown as empty if the value is zero.",
				Optional:    true,
				Computed:    true,
			},
			"lines_between_header": schema.Int64Attribute{
				Description: "The number of lines to log between logging the header line that summarizes the columns in the table.",
				Optional:    true,
				Computed:    true,
			},
			"included_ldap_stat": schema.SetAttribute{
				Description: "Specifies the types of statistics related to LDAP connections and operation processing that should be included in the output.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_resource_stat": schema.SetAttribute{
				Description: "Specifies whether statistics related to resource utilization such as JVM memory.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"histogram_format": schema.StringAttribute{
				Description: "The format of the data in the processing time histogram.",
				Optional:    true,
				Computed:    true,
			},
			"histogram_op_type": schema.SetAttribute{
				Description: "Specifies the operation type(s) to use when outputting the response time histogram data. The order of the operations here determines the order of the columns in the output. Use the per-application-ldap-stats setting to further control this.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"per_application_ldap_stats": schema.StringAttribute{
				Description: "Controls whether per application LDAP statistics are included in the output for selected LDAP operation statistics.",
				Optional:    true,
				Computed:    true,
			},
			"status_summary_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include about the status summary monitor entry.",
				Optional:    true,
				Computed:    true,
			},
			"ldap_changelog_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include for the LDAP changelog.",
				Optional:    true,
				Computed:    true,
			},
			"gauge_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include for Gauges.",
				Optional:    true,
				Computed:    true,
			},
			"log_file_format": schema.StringAttribute{
				Description: "Specifies the format to use when logging server statistics.",
				Optional:    true,
				Computed:    true,
			},
			"log_file": schema.StringAttribute{
				Description: "The file name to use for the log files generated by the Periodic Stats Logger Plugin. The path to the file can be specified either as relative to the server root or as an absolute path.",
				Required:    true,
			},
			"log_file_permissions": schema.StringAttribute{
				Description: "The UNIX permissions of the log files created by this Periodic Stats Logger Plugin.",
				Optional:    true,
				Computed:    true,
			},
			"append": schema.BoolAttribute{
				Description: "Specifies whether to append to existing log files.",
				Optional:    true,
				Computed:    true,
			},
			"rotation_policy": schema.SetAttribute{
				Description: "The rotation policy to use for the Periodic Stats Logger Plugin .",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"rotation_listener": schema.SetAttribute{
				Description: "A listener that should be notified whenever a log file is rotated out of service.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"retention_policy": schema.SetAttribute{
				Description: "The retention policy to use for the Periodic Stats Logger Plugin .",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"logging_error_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the server should exhibit if an error occurs during logging processing.",
				Optional:    true,
				Computed:    true,
			},
			"local_db_backend_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include about the Local DB Backends.",
				Optional:    true,
				Computed:    true,
			},
			"replication_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include about replication.",
				Optional:    true,
				Computed:    true,
			},
			"entry_cache_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include for each entry cache.",
				Optional:    true,
				Computed:    true,
			},
			"host_info": schema.SetAttribute{
				Description: "Specifies the level of detail to include about the host system resource utilization including CPU, memory, disk and network activity.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_ldap_application": schema.SetAttribute{
				Description: "If statistics should not be included for all applications, this property names the subset of applications that should be included.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Plugin",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the plug-in is enabled for use.",
				Required:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalPeriodicStatsLoggerPluginFields(ctx context.Context, addRequest *client.AddPeriodicStatsLoggerPluginRequest, plan periodicStatsLoggerPluginResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogInterval) {
		stringVal := plan.LogInterval.ValueString()
		addRequest.LogInterval = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CollectionInterval) {
		stringVal := plan.CollectionInterval.ValueString()
		addRequest.CollectionInterval = &stringVal
	}
	if internaltypes.IsDefined(plan.SuppressIfIdle) {
		boolVal := plan.SuppressIfIdle.ValueBool()
		addRequest.SuppressIfIdle = &boolVal
	}
	if internaltypes.IsDefined(plan.HeaderPrefixPerColumn) {
		boolVal := plan.HeaderPrefixPerColumn.ValueBool()
		addRequest.HeaderPrefixPerColumn = &boolVal
	}
	if internaltypes.IsDefined(plan.EmptyInsteadOfZero) {
		boolVal := plan.EmptyInsteadOfZero.ValueBool()
		addRequest.EmptyInsteadOfZero = &boolVal
	}
	if internaltypes.IsDefined(plan.LinesBetweenHeader) {
		intVal := int32(plan.LinesBetweenHeader.ValueInt64())
		addRequest.LinesBetweenHeader = &intVal
	}
	if internaltypes.IsDefined(plan.IncludedLDAPStat) {
		var slice []string
		plan.IncludedLDAPStat.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginIncludedLDAPStatProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginIncludedLDAPStatPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.IncludedLDAPStat = enumSlice
	}
	if internaltypes.IsDefined(plan.IncludedResourceStat) {
		var slice []string
		plan.IncludedResourceStat.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginIncludedResourceStatProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginIncludedResourceStatPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.IncludedResourceStat = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HistogramFormat) {
		histogramFormat, err := client.NewEnumpluginHistogramFormatPropFromValue(plan.HistogramFormat.ValueString())
		if err != nil {
			return err
		}
		addRequest.HistogramFormat = histogramFormat
	}
	if internaltypes.IsDefined(plan.HistogramOpType) {
		var slice []string
		plan.HistogramOpType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginHistogramOpTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginHistogramOpTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.HistogramOpType = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PerApplicationLDAPStats) {
		perApplicationLDAPStats, err := client.NewEnumpluginPeriodicStatsLoggerPerApplicationLDAPStatsPropFromValue(plan.PerApplicationLDAPStats.ValueString())
		if err != nil {
			return err
		}
		addRequest.PerApplicationLDAPStats = perApplicationLDAPStats
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.StatusSummaryInfo) {
		statusSummaryInfo, err := client.NewEnumpluginStatusSummaryInfoPropFromValue(plan.StatusSummaryInfo.ValueString())
		if err != nil {
			return err
		}
		addRequest.StatusSummaryInfo = statusSummaryInfo
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LdapChangelogInfo) {
		ldapChangelogInfo, err := client.NewEnumpluginLdapChangelogInfoPropFromValue(plan.LdapChangelogInfo.ValueString())
		if err != nil {
			return err
		}
		addRequest.LdapChangelogInfo = ldapChangelogInfo
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.GaugeInfo) {
		gaugeInfo, err := client.NewEnumpluginGaugeInfoPropFromValue(plan.GaugeInfo.ValueString())
		if err != nil {
			return err
		}
		addRequest.GaugeInfo = gaugeInfo
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFileFormat) {
		logFileFormat, err := client.NewEnumpluginLogFileFormatPropFromValue(plan.LogFileFormat.ValueString())
		if err != nil {
			return err
		}
		addRequest.LogFileFormat = logFileFormat
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		stringVal := plan.LogFilePermissions.ValueString()
		addRequest.LogFilePermissions = &stringVal
	}
	if internaltypes.IsDefined(plan.Append) {
		boolVal := plan.Append.ValueBool()
		addRequest.Append = &boolVal
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumpluginLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LocalDBBackendInfo) {
		localDBBackendInfo, err := client.NewEnumpluginLocalDBBackendInfoPropFromValue(plan.LocalDBBackendInfo.ValueString())
		if err != nil {
			return err
		}
		addRequest.LocalDBBackendInfo = localDBBackendInfo
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReplicationInfo) {
		replicationInfo, err := client.NewEnumpluginReplicationInfoPropFromValue(plan.ReplicationInfo.ValueString())
		if err != nil {
			return err
		}
		addRequest.ReplicationInfo = replicationInfo
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EntryCacheInfo) {
		entryCacheInfo, err := client.NewEnumpluginEntryCacheInfoPropFromValue(plan.EntryCacheInfo.ValueString())
		if err != nil {
			return err
		}
		addRequest.EntryCacheInfo = entryCacheInfo
	}
	if internaltypes.IsDefined(plan.HostInfo) {
		var slice []string
		plan.HostInfo.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginHostInfoProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginHostInfoPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.HostInfo = enumSlice
	}
	if internaltypes.IsDefined(plan.IncludedLDAPApplication) {
		var slice []string
		plan.IncludedLDAPApplication.ElementsAs(ctx, &slice, false)
		addRequest.IncludedLDAPApplication = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	return nil
}

// Read a PeriodicStatsLoggerPluginResponse object into the model struct
func readPeriodicStatsLoggerPluginResponse(ctx context.Context, r *client.PeriodicStatsLoggerPluginResponse, state *periodicStatsLoggerPluginResourceModel, expectedValues *periodicStatsLoggerPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.LogInterval = types.StringValue(r.LogInterval)
	config.CheckMismatchedPDFormattedAttributes("log_interval",
		expectedValues.LogInterval, state.LogInterval, diagnostics)
	state.CollectionInterval = types.StringValue(r.CollectionInterval)
	config.CheckMismatchedPDFormattedAttributes("collection_interval",
		expectedValues.CollectionInterval, state.CollectionInterval, diagnostics)
	state.SuppressIfIdle = types.BoolValue(r.SuppressIfIdle)
	state.HeaderPrefixPerColumn = internaltypes.BoolTypeOrNil(r.HeaderPrefixPerColumn)
	state.EmptyInsteadOfZero = internaltypes.BoolTypeOrNil(r.EmptyInsteadOfZero)
	state.LinesBetweenHeader = types.Int64Value(int64(r.LinesBetweenHeader))
	state.IncludedLDAPStat = internaltypes.GetStringSet(
		client.StringSliceEnumpluginIncludedLDAPStatProp(r.IncludedLDAPStat))
	state.IncludedResourceStat = internaltypes.GetStringSet(
		client.StringSliceEnumpluginIncludedResourceStatProp(r.IncludedResourceStat))
	state.HistogramFormat = types.StringValue(r.HistogramFormat.String())
	state.HistogramOpType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginHistogramOpTypeProp(r.HistogramOpType))
	state.PerApplicationLDAPStats = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginPeriodicStatsLoggerPerApplicationLDAPStatsProp(r.PerApplicationLDAPStats), internaltypes.IsEmptyString(expectedValues.PerApplicationLDAPStats))
	state.StatusSummaryInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginStatusSummaryInfoProp(r.StatusSummaryInfo), internaltypes.IsEmptyString(expectedValues.StatusSummaryInfo))
	state.LdapChangelogInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLdapChangelogInfoProp(r.LdapChangelogInfo), internaltypes.IsEmptyString(expectedValues.LdapChangelogInfo))
	state.GaugeInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginGaugeInfoProp(r.GaugeInfo), internaltypes.IsEmptyString(expectedValues.GaugeInfo))
	state.LogFileFormat = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLogFileFormatProp(r.LogFileFormat), internaltypes.IsEmptyString(expectedValues.LogFileFormat))
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.LocalDBBackendInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLocalDBBackendInfoProp(r.LocalDBBackendInfo), internaltypes.IsEmptyString(expectedValues.LocalDBBackendInfo))
	state.ReplicationInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginReplicationInfoProp(r.ReplicationInfo), internaltypes.IsEmptyString(expectedValues.ReplicationInfo))
	state.EntryCacheInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginEntryCacheInfoProp(r.EntryCacheInfo), internaltypes.IsEmptyString(expectedValues.EntryCacheInfo))
	state.HostInfo = internaltypes.GetStringSet(
		client.StringSliceEnumpluginHostInfoProp(r.HostInfo))
	state.IncludedLDAPApplication = internaltypes.GetStringSet(r.IncludedLDAPApplication)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createPeriodicStatsLoggerPluginOperations(plan periodicStatsLoggerPluginResourceModel, state periodicStatsLoggerPluginResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.LogInterval, state.LogInterval, "log-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.CollectionInterval, state.CollectionInterval, "collection-interval")
	operations.AddBoolOperationIfNecessary(&ops, plan.SuppressIfIdle, state.SuppressIfIdle, "suppress-if-idle")
	operations.AddBoolOperationIfNecessary(&ops, plan.HeaderPrefixPerColumn, state.HeaderPrefixPerColumn, "header-prefix-per-column")
	operations.AddBoolOperationIfNecessary(&ops, plan.EmptyInsteadOfZero, state.EmptyInsteadOfZero, "empty-instead-of-zero")
	operations.AddInt64OperationIfNecessary(&ops, plan.LinesBetweenHeader, state.LinesBetweenHeader, "lines-between-header")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedLDAPStat, state.IncludedLDAPStat, "included-ldap-stat")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedResourceStat, state.IncludedResourceStat, "included-resource-stat")
	operations.AddStringOperationIfNecessary(&ops, plan.HistogramFormat, state.HistogramFormat, "histogram-format")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.HistogramOpType, state.HistogramOpType, "histogram-op-type")
	operations.AddStringOperationIfNecessary(&ops, plan.PerApplicationLDAPStats, state.PerApplicationLDAPStats, "per-application-ldap-stats")
	operations.AddStringOperationIfNecessary(&ops, plan.StatusSummaryInfo, state.StatusSummaryInfo, "status-summary-info")
	operations.AddStringOperationIfNecessary(&ops, plan.LdapChangelogInfo, state.LdapChangelogInfo, "ldap-changelog-info")
	operations.AddStringOperationIfNecessary(&ops, plan.GaugeInfo, state.GaugeInfo, "gauge-info")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFileFormat, state.LogFileFormat, "log-file-format")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFile, state.LogFile, "log-file")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFilePermissions, state.LogFilePermissions, "log-file-permissions")
	operations.AddBoolOperationIfNecessary(&ops, plan.Append, state.Append, "append")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RotationPolicy, state.RotationPolicy, "rotation-policy")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RotationListener, state.RotationListener, "rotation-listener")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RetentionPolicy, state.RetentionPolicy, "retention-policy")
	operations.AddStringOperationIfNecessary(&ops, plan.LoggingErrorBehavior, state.LoggingErrorBehavior, "logging-error-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.LocalDBBackendInfo, state.LocalDBBackendInfo, "local-db-backend-info")
	operations.AddStringOperationIfNecessary(&ops, plan.ReplicationInfo, state.ReplicationInfo, "replication-info")
	operations.AddStringOperationIfNecessary(&ops, plan.EntryCacheInfo, state.EntryCacheInfo, "entry-cache-info")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.HostInfo, state.HostInfo, "host-info")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedLDAPApplication, state.IncludedLDAPApplication, "included-ldap-application")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *periodicStatsLoggerPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan periodicStatsLoggerPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddPeriodicStatsLoggerPluginRequest(plan.Id.ValueString(),
		[]client.EnumperiodicStatsLoggerPluginSchemaUrn{client.ENUMPERIODICSTATSLOGGERPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINPERIODIC_STATS_LOGGER},
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalPeriodicStatsLoggerPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Periodic Stats Logger Plugin", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddPeriodicStatsLoggerPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Periodic Stats Logger Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state periodicStatsLoggerPluginResourceModel
	readPeriodicStatsLoggerPluginResponse(ctx, addResponse.PeriodicStatsLoggerPluginResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *periodicStatsLoggerPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state periodicStatsLoggerPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Periodic Stats Logger Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readPeriodicStatsLoggerPluginResponse(ctx, readResponse.PeriodicStatsLoggerPluginResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *periodicStatsLoggerPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan periodicStatsLoggerPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state periodicStatsLoggerPluginResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.PluginApi.UpdatePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createPeriodicStatsLoggerPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Periodic Stats Logger Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPeriodicStatsLoggerPluginResponse(ctx, updateResponse.PeriodicStatsLoggerPluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *periodicStatsLoggerPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state periodicStatsLoggerPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PluginApi.DeletePluginExecute(r.apiClient.PluginApi.DeletePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Periodic Stats Logger Plugin", err, httpResp)
		return
	}
}

func (r *periodicStatsLoggerPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
