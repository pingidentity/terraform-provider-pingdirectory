package plugin

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &pluginsDataSource{}
	_ datasource.DataSourceWithConfigure = &pluginsDataSource{}
)

// Create a Plugins data source
func NewPluginsDataSource() datasource.DataSource {
	return &pluginsDataSource{}
}

// pluginsDataSource is the datasource implementation.
type pluginsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *pluginsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_plugins"
}

// Configure adds the provider configured client to the data source.
func (r *pluginsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type pluginsDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *pluginsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Plugin objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Plugin objects found in the configuration",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: internaltypes.ObjectsObjectType(),
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read resource information
func (r *pluginsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state pluginsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.PluginAPI.ListPlugins(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.PluginAPI.ListPluginsExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Plugin objects", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	objects := []attr.Value{}
	for _, response := range readResponse.Resources {
		attributes := map[string]attr.Value{}
		if response.LastAccessTimePluginResponse != nil {
			attributes["id"] = types.StringValue(response.LastAccessTimePluginResponse.Id)
			attributes["type"] = types.StringValue("last-access-time")
		}
		if response.StatsCollectorPluginResponse != nil {
			attributes["id"] = types.StringValue(response.StatsCollectorPluginResponse.Id)
			attributes["type"] = types.StringValue("stats-collector")
		}
		if response.TraditionalStaticGroupSupportForInvertedStaticGroupsPluginResponse != nil {
			attributes["id"] = types.StringValue(response.TraditionalStaticGroupSupportForInvertedStaticGroupsPluginResponse.Id)
			attributes["type"] = types.StringValue("traditional-static-group-support-for-inverted-static-groups")
		}
		if response.InternalSearchRatePluginResponse != nil {
			attributes["id"] = types.StringValue(response.InternalSearchRatePluginResponse.Id)
			attributes["type"] = types.StringValue("internal-search-rate")
		}
		if response.ModifiablePasswordPolicyStatePluginResponse != nil {
			attributes["id"] = types.StringValue(response.ModifiablePasswordPolicyStatePluginResponse.Id)
			attributes["type"] = types.StringValue("modifiable-password-policy-state")
		}
		if response.SevenBitCleanPluginResponse != nil {
			attributes["id"] = types.StringValue(response.SevenBitCleanPluginResponse.Id)
			attributes["type"] = types.StringValue("seven-bit-clean")
		}
		if response.CleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse != nil {
			attributes["id"] = types.StringValue(response.CleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse.Id)
			attributes["type"] = types.StringValue("clean-up-expired-pingfederate-persistent-access-grants")
		}
		if response.PeriodicGcPluginResponse != nil {
			attributes["id"] = types.StringValue(response.PeriodicGcPluginResponse.Id)
			attributes["type"] = types.StringValue("periodic-gc")
		}
		if response.PingOnePassThroughAuthenticationPluginResponse != nil {
			attributes["id"] = types.StringValue(response.PingOnePassThroughAuthenticationPluginResponse.Id)
			attributes["type"] = types.StringValue("ping-one-pass-through-authentication")
		}
		if response.SecretKeyDeleteAlertPluginResponse != nil {
			attributes["id"] = types.StringValue(response.SecretKeyDeleteAlertPluginResponse.Id)
			attributes["type"] = types.StringValue("secret-key-delete-alert")
		}
		if response.ChangelogPasswordEncryptionPluginResponse != nil {
			attributes["id"] = types.StringValue(response.ChangelogPasswordEncryptionPluginResponse.Id)
			attributes["type"] = types.StringValue("changelog-password-encryption")
		}
		if response.GloballyUniqueAttributePluginResponse != nil {
			attributes["id"] = types.StringValue(response.GloballyUniqueAttributePluginResponse.Id)
			attributes["type"] = types.StringValue("globally-unique-attribute")
		}
		if response.ProcessingTimeHistogramPluginResponse != nil {
			attributes["id"] = types.StringValue(response.ProcessingTimeHistogramPluginResponse.Id)
			attributes["type"] = types.StringValue("processing-time-histogram")
		}
		if response.SearchShutdownPluginResponse != nil {
			attributes["id"] = types.StringValue(response.SearchShutdownPluginResponse.Id)
			attributes["type"] = types.StringValue("search-shutdown")
		}
		if response.PeriodicStatsLoggerPluginResponse != nil {
			attributes["id"] = types.StringValue(response.PeriodicStatsLoggerPluginResponse.Id)
			attributes["type"] = types.StringValue("periodic-stats-logger")
		}
		if response.PurgeExpiredDataPluginResponse != nil {
			attributes["id"] = types.StringValue(response.PurgeExpiredDataPluginResponse.Id)
			attributes["type"] = types.StringValue("purge-expired-data")
		}
		if response.ChangeSubscriptionNotificationPluginResponse != nil {
			attributes["id"] = types.StringValue(response.ChangeSubscriptionNotificationPluginResponse.Id)
			attributes["type"] = types.StringValue("change-subscription-notification")
		}
		if response.LdapAttributeDescriptionListPluginResponse != nil {
			attributes["id"] = types.StringValue(response.LdapAttributeDescriptionListPluginResponse.Id)
			attributes["type"] = types.StringValue("ldap-attribute-description-list")
		}
		if response.SubOperationTimingPluginResponse != nil {
			attributes["id"] = types.StringValue(response.SubOperationTimingPluginResponse.Id)
			attributes["type"] = types.StringValue("sub-operation-timing")
		}
		if response.ThirdPartyPluginResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyPluginResponse.Id)
			attributes["type"] = types.StringValue("third-party")
		}
		if response.EncryptAttributeValuesPluginResponse != nil {
			attributes["id"] = types.StringValue(response.EncryptAttributeValuesPluginResponse.Id)
			attributes["type"] = types.StringValue("encrypt-attribute-values")
		}
		if response.PassThroughAuthenticationPluginResponse != nil {
			attributes["id"] = types.StringValue(response.PassThroughAuthenticationPluginResponse.Id)
			attributes["type"] = types.StringValue("pass-through-authentication")
		}
		if response.GlobalReferentialIntegrityPluginResponse != nil {
			attributes["id"] = types.StringValue(response.GlobalReferentialIntegrityPluginResponse.Id)
			attributes["type"] = types.StringValue("global-referential-integrity")
		}
		if response.DnMapperPluginResponse != nil {
			attributes["id"] = types.StringValue(response.DnMapperPluginResponse.Id)
			attributes["type"] = types.StringValue("dn-mapper")
		}
		if response.MonitorHistoryPluginResponse != nil {
			attributes["id"] = types.StringValue(response.MonitorHistoryPluginResponse.Id)
			attributes["type"] = types.StringValue("monitor-history")
		}
		if response.ReferralOnUpdatePluginResponse != nil {
			attributes["id"] = types.StringValue(response.ReferralOnUpdatePluginResponse.Id)
			attributes["type"] = types.StringValue("referral-on-update")
		}
		if response.SimpleToExternalBindPluginResponse != nil {
			attributes["id"] = types.StringValue(response.SimpleToExternalBindPluginResponse.Id)
			attributes["type"] = types.StringValue("simple-to-external-bind")
		}
		if response.CustomPluginResponse != nil {
			attributes["id"] = types.StringValue(response.CustomPluginResponse.Id)
			attributes["type"] = types.StringValue("custom")
		}
		if response.SnmpSubagentPluginResponse != nil {
			attributes["id"] = types.StringValue(response.SnmpSubagentPluginResponse.Id)
			attributes["type"] = types.StringValue("snmp-subagent")
		}
		if response.CoalesceModificationsPluginResponse != nil {
			attributes["id"] = types.StringValue(response.CoalesceModificationsPluginResponse.Id)
			attributes["type"] = types.StringValue("coalesce-modifications")
		}
		if response.PasswordPolicyImportPluginResponse != nil {
			attributes["id"] = types.StringValue(response.PasswordPolicyImportPluginResponse.Id)
			attributes["type"] = types.StringValue("password-policy-import")
		}
		if response.ProfilerPluginResponse != nil {
			attributes["id"] = types.StringValue(response.ProfilerPluginResponse.Id)
			attributes["type"] = types.StringValue("profiler")
		}
		if response.EntryUuidPluginResponse != nil {
			attributes["id"] = types.StringValue(response.EntryUuidPluginResponse.Id)
			attributes["type"] = types.StringValue("entry-uuid")
		}
		if response.CleanUpInactivePingfederatePersistentSessionsPluginResponse != nil {
			attributes["id"] = types.StringValue(response.CleanUpInactivePingfederatePersistentSessionsPluginResponse.Id)
			attributes["type"] = types.StringValue("clean-up-inactive-pingfederate-persistent-sessions")
		}
		if response.ComposedAttributePluginResponse != nil {
			attributes["id"] = types.StringValue(response.ComposedAttributePluginResponse.Id)
			attributes["type"] = types.StringValue("composed-attribute")
		}
		if response.LdapResultCodeTrackerPluginResponse != nil {
			attributes["id"] = types.StringValue(response.LdapResultCodeTrackerPluginResponse.Id)
			attributes["type"] = types.StringValue("ldap-result-code-tracker")
		}
		if response.AttributeMapperPluginResponse != nil {
			attributes["id"] = types.StringValue(response.AttributeMapperPluginResponse.Id)
			attributes["type"] = types.StringValue("attribute-mapper")
		}
		if response.DelayPluginResponse != nil {
			attributes["id"] = types.StringValue(response.DelayPluginResponse.Id)
			attributes["type"] = types.StringValue("delay")
		}
		if response.PreUpdateConfigPluginResponse != nil {
			attributes["id"] = types.StringValue(response.PreUpdateConfigPluginResponse.Id)
			attributes["type"] = types.StringValue("pre-update-config")
		}
		if response.CleanUpExpiredPingfederatePersistentSessionsPluginResponse != nil {
			attributes["id"] = types.StringValue(response.CleanUpExpiredPingfederatePersistentSessionsPluginResponse.Id)
			attributes["type"] = types.StringValue("clean-up-expired-pingfederate-persistent-sessions")
		}
		if response.GroovyScriptedPluginResponse != nil {
			attributes["id"] = types.StringValue(response.GroovyScriptedPluginResponse.Id)
			attributes["type"] = types.StringValue("groovy-scripted")
		}
		if response.LastModPluginResponse != nil {
			attributes["id"] = types.StringValue(response.LastModPluginResponse.Id)
			attributes["type"] = types.StringValue("last-mod")
		}
		if response.PluggablePassThroughAuthenticationPluginResponse != nil {
			attributes["id"] = types.StringValue(response.PluggablePassThroughAuthenticationPluginResponse.Id)
			attributes["type"] = types.StringValue("pluggable-pass-through-authentication")
		}
		if response.ReferentialIntegrityPluginResponse != nil {
			attributes["id"] = types.StringValue(response.ReferentialIntegrityPluginResponse.Id)
			attributes["type"] = types.StringValue("referential-integrity")
		}
		if response.UniqueAttributePluginResponse != nil {
			attributes["id"] = types.StringValue(response.UniqueAttributePluginResponse.Id)
			attributes["type"] = types.StringValue("unique-attribute")
		}
		if response.SnmpMasterAgentPluginResponse != nil {
			attributes["id"] = types.StringValue(response.SnmpMasterAgentPluginResponse.Id)
			attributes["type"] = types.StringValue("snmp-master-agent")
		}
		if response.InvertedStaticGroupReferentialIntegrityPluginResponse != nil {
			attributes["id"] = types.StringValue(response.InvertedStaticGroupReferentialIntegrityPluginResponse.Id)
			attributes["type"] = types.StringValue("inverted-static-group-referential-integrity")
		}
		obj, diags := types.ObjectValue(internaltypes.ObjectsAttrTypes(), attributes)
		resp.Diagnostics.Append(diags...)
		objects = append(objects, obj)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	state.Objects, diags = types.SetValue(internaltypes.ObjectsObjectType(), objects)
	resp.Diagnostics.Append(diags...)
	state.Id = types.StringValue("id")

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
