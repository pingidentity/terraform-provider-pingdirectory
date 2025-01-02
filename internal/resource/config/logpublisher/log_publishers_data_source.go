package logpublisher

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &logPublishersDataSource{}
	_ datasource.DataSourceWithConfigure = &logPublishersDataSource{}
)

// Create a Log Publishers data source
func NewLogPublishersDataSource() datasource.DataSource {
	return &logPublishersDataSource{}
}

// logPublishersDataSource is the datasource implementation.
type logPublishersDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *logPublishersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_publishers"
}

// Configure adds the provider configured client to the data source.
func (r *logPublishersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type logPublishersDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *logPublishersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Log Publisher objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Log Publisher objects found in the configuration",
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
func (r *logPublishersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state logPublishersDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.LogPublisherAPI.ListLogPublishers(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.LogPublisherAPI.ListLogPublishersExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Log Publisher objects", err, httpResp)
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
		if response.SyslogJsonAuditLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.SyslogJsonAuditLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("syslog-json-audit")
		}
		if response.SyslogBasedErrorLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.SyslogBasedErrorLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("syslog-based-error")
		}
		if response.ThirdPartyHttpOperationLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyHttpOperationLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("third-party-http-operation")
		}
		if response.FileBasedTraceLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.FileBasedTraceLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("file-based-trace")
		}
		if response.JdbcBasedAccessLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.JdbcBasedAccessLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("jdbc-based-access")
		}
		if response.SyslogTextErrorLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.SyslogTextErrorLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("syslog-text-error")
		}
		if response.ThirdPartyPolicyDecisionLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyPolicyDecisionLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("third-party-policy-decision")
		}
		if response.SyslogBasedAccessLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.SyslogBasedAccessLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("syslog-based-access")
		}
		if response.FileBasedDebugLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.FileBasedDebugLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("file-based-debug")
		}
		if response.ConsoleJsonSyncFailedOpsLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.ConsoleJsonSyncFailedOpsLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("console-json-sync-failed-ops")
		}
		if response.FileBasedErrorLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.FileBasedErrorLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("file-based-error")
		}
		if response.ConsoleJsonSyncLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.ConsoleJsonSyncLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("console-json-sync")
		}
		if response.FileBasedPolicyDecisionLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.FileBasedPolicyDecisionLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("file-based-policy-decision")
		}
		if response.DebugAccessLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.DebugAccessLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("debug-access")
		}
		if response.SyncFailedOpsLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.SyncFailedOpsLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("sync-failed-ops")
		}
		if response.ThirdPartyAccessLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyAccessLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("third-party-access")
		}
		if response.FileBasedAuditLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.FileBasedAuditLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("file-based-audit")
		}
		if response.JsonErrorLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.JsonErrorLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("json-error")
		}
		if response.SyslogJsonSyncLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.SyslogJsonSyncLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("syslog-json-sync")
		}
		if response.GroovyScriptedFileBasedAccessLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.GroovyScriptedFileBasedAccessLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("groovy-scripted-file-based-access")
		}
		if response.SyslogJsonAccessLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.SyslogJsonAccessLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("syslog-json-access")
		}
		if response.SyslogJsonSyncFailedOpsLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.SyslogJsonSyncFailedOpsLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("syslog-json-sync-failed-ops")
		}
		if response.GroovyScriptedAccessLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.GroovyScriptedAccessLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("groovy-scripted-access")
		}
		if response.ConsoleJsonHttpOperationLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.ConsoleJsonHttpOperationLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("console-json-http-operation")
		}
		if response.FileBasedAccessLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.FileBasedAccessLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("file-based-access")
		}
		if response.FileBasedSyncLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.FileBasedSyncLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("file-based-sync")
		}
		if response.SyslogJsonErrorLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.SyslogJsonErrorLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("syslog-json-error")
		}
		if response.ThirdPartyFileBasedAccessLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyFileBasedAccessLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("third-party-file-based-access")
		}
		if response.OperationTimingAccessLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.OperationTimingAccessLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("operation-timing-access")
		}
		if response.AdminAlertAccessLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.AdminAlertAccessLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("admin-alert-access")
		}
		if response.JdbcBasedErrorLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.JdbcBasedErrorLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("jdbc-based-error")
		}
		if response.CommonLogFileHttpOperationLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.CommonLogFileHttpOperationLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("common-log-file-http-operation")
		}
		if response.ConsoleJsonErrorLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.ConsoleJsonErrorLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("console-json-error")
		}
		if response.FileBasedPolicyQueryLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.FileBasedPolicyQueryLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("file-based-policy-query")
		}
		if response.FileBasedJsonAuditLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.FileBasedJsonAuditLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("file-based-json-audit")
		}
		if response.ThirdPartyErrorLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyErrorLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("third-party-error")
		}
		if response.SyslogTextAccessLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.SyslogTextAccessLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("syslog-text-access")
		}
		if response.DetailedHttpOperationLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.DetailedHttpOperationLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("detailed-http-operation")
		}
		if response.JsonAccessLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.JsonAccessLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("json-access")
		}
		if response.SyslogJsonHttpOperationLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.SyslogJsonHttpOperationLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("syslog-json-http-operation")
		}
		if response.GroovyScriptedFileBasedErrorLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.GroovyScriptedFileBasedErrorLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("groovy-scripted-file-based-error")
		}
		if response.ThirdPartyFileBasedErrorLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyFileBasedErrorLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("third-party-file-based-error")
		}
		if response.ConsoleJsonAuditLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.ConsoleJsonAuditLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("console-json-audit")
		}
		if response.ConsoleJsonAccessLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.ConsoleJsonAccessLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("console-json-access")
		}
		if response.FileBasedJsonSyncLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.FileBasedJsonSyncLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("file-based-json-sync")
		}
		if response.GroovyScriptedErrorLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.GroovyScriptedErrorLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("groovy-scripted-error")
		}
		if response.FileBasedJsonHttpOperationLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.FileBasedJsonHttpOperationLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("file-based-json-http-operation")
		}
		if response.FileBasedJsonSyncFailedOpsLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.FileBasedJsonSyncFailedOpsLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("file-based-json-sync-failed-ops")
		}
		if response.GroovyScriptedHttpOperationLogPublisherResponse != nil {
			attributes["id"] = types.StringValue(response.GroovyScriptedHttpOperationLogPublisherResponse.Id)
			attributes["type"] = types.StringValue("groovy-scripted-http-operation")
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
