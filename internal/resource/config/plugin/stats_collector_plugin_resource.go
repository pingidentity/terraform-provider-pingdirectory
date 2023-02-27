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
	_ resource.Resource                = &statsCollectorPluginResource{}
	_ resource.ResourceWithConfigure   = &statsCollectorPluginResource{}
	_ resource.ResourceWithImportState = &statsCollectorPluginResource{}
)

// Create a Stats Collector Plugin resource
func NewStatsCollectorPluginResource() resource.Resource {
	return &statsCollectorPluginResource{}
}

// statsCollectorPluginResource is the resource implementation.
type statsCollectorPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *statsCollectorPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stats_collector_plugin"
}

// Configure adds the provider configured client to the resource.
func (r *statsCollectorPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type statsCollectorPluginResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	LastUpdated             types.String `tfsdk:"last_updated"`
	Notifications           types.Set    `tfsdk:"notifications"`
	RequiredActions         types.Set    `tfsdk:"required_actions"`
	SampleInterval          types.String `tfsdk:"sample_interval"`
	CollectionInterval      types.String `tfsdk:"collection_interval"`
	LdapInfo                types.String `tfsdk:"ldap_info"`
	ServerInfo              types.String `tfsdk:"server_info"`
	PerApplicationLDAPStats types.String `tfsdk:"per_application_ldap_stats"`
	LdapChangelogInfo       types.String `tfsdk:"ldap_changelog_info"`
	StatusSummaryInfo       types.String `tfsdk:"status_summary_info"`
	GenerateCollectorFiles  types.Bool   `tfsdk:"generate_collector_files"`
	LocalDBBackendInfo      types.String `tfsdk:"local_db_backend_info"`
	ReplicationInfo         types.String `tfsdk:"replication_info"`
	EntryCacheInfo          types.String `tfsdk:"entry_cache_info"`
	HostInfo                types.Set    `tfsdk:"host_info"`
	IncludedLDAPApplication types.Set    `tfsdk:"included_ldap_application"`
	Description             types.String `tfsdk:"description"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *statsCollectorPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Stats Collector Plugin.",
		Attributes: map[string]schema.Attribute{
			"sample_interval": schema.StringAttribute{
				Description: "The duration between statistics collections. Setting this value too small can have an impact on performance. This value should be a multiple of collection-interval.",
				Optional:    true,
				Computed:    true,
			},
			"collection_interval": schema.StringAttribute{
				Description: "Some of the calculated statistics, such as the average and maximum queue sizes, can use multiple samples within a log interval. This value controls how often samples are gathered, and setting this value too small can have an adverse impact on performance.",
				Optional:    true,
				Computed:    true,
			},
			"ldap_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include about the LDAP connection handlers.",
				Optional:    true,
				Computed:    true,
			},
			"server_info": schema.StringAttribute{
				Description: "Specifies whether statistics related to resource utilization such as JVM memory and CPU/Network/Disk utilization.",
				Optional:    true,
				Computed:    true,
			},
			"per_application_ldap_stats": schema.StringAttribute{
				Description: "Controls whether per application LDAP statistics are included in the output for selected LDAP operation statistics.",
				Optional:    true,
				Computed:    true,
			},
			"ldap_changelog_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include for the LDAP changelog.",
				Optional:    true,
				Computed:    true,
			},
			"status_summary_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include about the status summary monitor entry.",
				Optional:    true,
				Computed:    true,
			},
			"generate_collector_files": schema.BoolAttribute{
				Description: "Indicates whether this plugin should store metric samples on disk for use by the Data Metrics Server. If the Stats Collector Plugin is only being used to collect metrics for one or more StatsD Monitoring Endpoints, then this can be set to false to prevent unnecessary I/O.",
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
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the plug-in is enabled for use.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Read a StatsCollectorPluginResponse object into the model struct
func readStatsCollectorPluginResponse(ctx context.Context, r *client.StatsCollectorPluginResponse, state *statsCollectorPluginResourceModel, expectedValues *statsCollectorPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.SampleInterval = types.StringValue(r.SampleInterval)
	config.CheckMismatchedPDFormattedAttributes("sample_interval",
		expectedValues.SampleInterval, state.SampleInterval, diagnostics)
	state.CollectionInterval = types.StringValue(r.CollectionInterval)
	config.CheckMismatchedPDFormattedAttributes("collection_interval",
		expectedValues.CollectionInterval, state.CollectionInterval, diagnostics)
	state.LdapInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLdapInfoProp(r.LdapInfo), true)
	state.ServerInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginServerInfoProp(r.ServerInfo), true)
	state.PerApplicationLDAPStats = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginStatsCollectorPerApplicationLDAPStatsProp(r.PerApplicationLDAPStats), true)
	state.LdapChangelogInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLdapChangelogInfoProp(r.LdapChangelogInfo), true)
	state.StatusSummaryInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginStatusSummaryInfoProp(r.StatusSummaryInfo), true)
	state.GenerateCollectorFiles = internaltypes.BoolTypeOrNil(r.GenerateCollectorFiles)
	state.LocalDBBackendInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLocalDBBackendInfoProp(r.LocalDBBackendInfo), true)
	state.ReplicationInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginReplicationInfoProp(r.ReplicationInfo), true)
	state.EntryCacheInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginEntryCacheInfoProp(r.EntryCacheInfo), true)
	state.HostInfo = internaltypes.GetStringSet(
		client.StringSliceEnumpluginHostInfoProp(r.HostInfo))
	state.IncludedLDAPApplication = internaltypes.GetStringSet(r.IncludedLDAPApplication)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createStatsCollectorPluginOperations(plan statsCollectorPluginResourceModel, state statsCollectorPluginResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.SampleInterval, state.SampleInterval, "sample-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.CollectionInterval, state.CollectionInterval, "collection-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.LdapInfo, state.LdapInfo, "ldap-info")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerInfo, state.ServerInfo, "server-info")
	operations.AddStringOperationIfNecessary(&ops, plan.PerApplicationLDAPStats, state.PerApplicationLDAPStats, "per-application-ldap-stats")
	operations.AddStringOperationIfNecessary(&ops, plan.LdapChangelogInfo, state.LdapChangelogInfo, "ldap-changelog-info")
	operations.AddStringOperationIfNecessary(&ops, plan.StatusSummaryInfo, state.StatusSummaryInfo, "status-summary-info")
	operations.AddBoolOperationIfNecessary(&ops, plan.GenerateCollectorFiles, state.GenerateCollectorFiles, "generate-collector-files")
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
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *statsCollectorPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan statsCollectorPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Stats Collector Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state statsCollectorPluginResourceModel
	readStatsCollectorPluginResponse(ctx, readResponse.StatsCollectorPluginResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PluginApi.UpdatePlugin(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createStatsCollectorPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Stats Collector Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readStatsCollectorPluginResponse(ctx, updateResponse.StatsCollectorPluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *statsCollectorPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state statsCollectorPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Stats Collector Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readStatsCollectorPluginResponse(ctx, readResponse.StatsCollectorPluginResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *statsCollectorPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan statsCollectorPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state statsCollectorPluginResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.PluginApi.UpdatePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createStatsCollectorPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Stats Collector Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readStatsCollectorPluginResponse(ctx, updateResponse.StatsCollectorPluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *statsCollectorPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *statsCollectorPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
