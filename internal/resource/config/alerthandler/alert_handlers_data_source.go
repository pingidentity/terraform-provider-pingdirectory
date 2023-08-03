package alerthandler

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &alertHandlersDataSource{}
	_ datasource.DataSourceWithConfigure = &alertHandlersDataSource{}
)

// Create a Alert Handlers data source
func NewAlertHandlersDataSource() datasource.DataSource {
	return &alertHandlersDataSource{}
}

// alertHandlersDataSource is the datasource implementation.
type alertHandlersDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *alertHandlersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_handlers"
}

// Configure adds the provider configured client to the data source.
func (r *alertHandlersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type alertHandlersDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *alertHandlersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Alert Handler objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Alert Handler objects found in the configuration",
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
func (r *alertHandlersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state alertHandlersDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.AlertHandlerApi.ListAlertHandlers(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.AlertHandlerApi.ListAlertHandlersExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Alert Handler objects", err, httpResp)
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
		if response.OutputAlertHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.OutputAlertHandlerResponse.Id)
			attributes["type"] = types.StringValue("output")
		}
		if response.SmtpAlertHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.SmtpAlertHandlerResponse.Id)
			attributes["type"] = types.StringValue("smtp")
		}
		if response.JmxAlertHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.JmxAlertHandlerResponse.Id)
			attributes["type"] = types.StringValue("jmx")
		}
		if response.GroovyScriptedAlertHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.GroovyScriptedAlertHandlerResponse.Id)
			attributes["type"] = types.StringValue("groovy-scripted")
		}
		if response.CustomAlertHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.CustomAlertHandlerResponse.Id)
			attributes["type"] = types.StringValue("custom")
		}
		if response.SnmpAlertHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.SnmpAlertHandlerResponse.Id)
			attributes["type"] = types.StringValue("snmp")
		}
		if response.TwilioAlertHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.TwilioAlertHandlerResponse.Id)
			attributes["type"] = types.StringValue("twilio")
		}
		if response.ErrorLogAlertHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.ErrorLogAlertHandlerResponse.Id)
			attributes["type"] = types.StringValue("error-log")
		}
		if response.SnmpSubAgentAlertHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.SnmpSubAgentAlertHandlerResponse.Id)
			attributes["type"] = types.StringValue("snmp-sub-agent")
		}
		if response.ExecAlertHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.ExecAlertHandlerResponse.Id)
			attributes["type"] = types.StringValue("exec")
		}
		if response.ThirdPartyAlertHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyAlertHandlerResponse.Id)
			attributes["type"] = types.StringValue("third-party")
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
