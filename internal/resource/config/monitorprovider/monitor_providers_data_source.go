package monitorprovider

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
	_ datasource.DataSource              = &monitorProvidersDataSource{}
	_ datasource.DataSourceWithConfigure = &monitorProvidersDataSource{}
)

// Create a Monitor Providers data source
func NewMonitorProvidersDataSource() datasource.DataSource {
	return &monitorProvidersDataSource{}
}

// monitorProvidersDataSource is the datasource implementation.
type monitorProvidersDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *monitorProvidersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor_providers"
}

// Configure adds the provider configured client to the data source.
func (r *monitorProvidersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type monitorProvidersDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *monitorProvidersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists Monitor Provider objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder name of this object required by Terraform.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Monitor Provider objects found in the configuration",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: internaltypes.ObjectsObjectType(),
			},
		},
	}
}

// Read resource information
func (r *monitorProvidersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state monitorProvidersDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.MonitorProviderApi.ListMonitorProviders(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.MonitorProviderApi.ListMonitorProvidersExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Monitor Provider objects", err, httpResp)
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
		if response.MemoryUsageMonitorProviderResponse != nil {
			attributes["id"] = types.StringValue(response.MemoryUsageMonitorProviderResponse.Id)
			attributes["type"] = types.StringValue("memory-usage")
		}
		if response.StackTraceMonitorProviderResponse != nil {
			attributes["id"] = types.StringValue(response.StackTraceMonitorProviderResponse.Id)
			attributes["type"] = types.StringValue("stack-trace")
		}
		if response.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse != nil {
			attributes["id"] = types.StringValue(response.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse.Id)
			attributes["type"] = types.StringValue("encryption-settings-database-accessibility")
		}
		if response.CustomMonitorProviderResponse != nil {
			attributes["id"] = types.StringValue(response.CustomMonitorProviderResponse.Id)
			attributes["type"] = types.StringValue("custom")
		}
		if response.ActiveOperationsMonitorProviderResponse != nil {
			attributes["id"] = types.StringValue(response.ActiveOperationsMonitorProviderResponse.Id)
			attributes["type"] = types.StringValue("active-operations")
		}
		if response.SslContextMonitorProviderResponse != nil {
			attributes["id"] = types.StringValue(response.SslContextMonitorProviderResponse.Id)
			attributes["type"] = types.StringValue("ssl-context")
		}
		if response.VersionMonitorProviderResponse != nil {
			attributes["id"] = types.StringValue(response.VersionMonitorProviderResponse.Id)
			attributes["type"] = types.StringValue("version")
		}
		if response.HostSystemMonitorProviderResponse != nil {
			attributes["id"] = types.StringValue(response.HostSystemMonitorProviderResponse.Id)
			attributes["type"] = types.StringValue("host-system")
		}
		if response.GeneralMonitorProviderResponse != nil {
			attributes["id"] = types.StringValue(response.GeneralMonitorProviderResponse.Id)
			attributes["type"] = types.StringValue("general")
		}
		if response.DiskSpaceUsageMonitorProviderResponse != nil {
			attributes["id"] = types.StringValue(response.DiskSpaceUsageMonitorProviderResponse.Id)
			attributes["type"] = types.StringValue("disk-space-usage")
		}
		if response.SystemInfoMonitorProviderResponse != nil {
			attributes["id"] = types.StringValue(response.SystemInfoMonitorProviderResponse.Id)
			attributes["type"] = types.StringValue("system-info")
		}
		if response.ClientConnectionMonitorProviderResponse != nil {
			attributes["id"] = types.StringValue(response.ClientConnectionMonitorProviderResponse.Id)
			attributes["type"] = types.StringValue("client-connection")
		}
		if response.ThirdPartyMonitorProviderResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyMonitorProviderResponse.Id)
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
