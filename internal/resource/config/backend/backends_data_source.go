// Copyright © 2025 Ping Identity Corporation

package backend

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
	_ datasource.DataSource              = &backendsDataSource{}
	_ datasource.DataSourceWithConfigure = &backendsDataSource{}
)

// Create a Backends data source
func NewBackendsDataSource() datasource.DataSource {
	return &backendsDataSource{}
}

// backendsDataSource is the datasource implementation.
type backendsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *backendsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backends"
}

// Configure adds the provider configured client to the data source.
func (r *backendsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type backendsDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *backendsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Backend objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Backend objects found in the configuration",
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
func (r *backendsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state backendsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.BackendAPI.ListBackends(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.BackendAPI.ListBackendsExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Backend objects", err, httpResp)
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
		if response.SchemaBackendResponse != nil {
			attributes["id"] = types.StringValue(response.SchemaBackendResponse.Id)
			attributes["type"] = types.StringValue("schema")
		}
		if response.BackupBackendResponse != nil {
			attributes["id"] = types.StringValue(response.BackupBackendResponse.Id)
			attributes["type"] = types.StringValue("backup")
		}
		if response.MemoryBackendResponse != nil {
			attributes["id"] = types.StringValue(response.MemoryBackendResponse.Id)
			attributes["type"] = types.StringValue("memory")
		}
		if response.EncryptionSettingsBackendResponse != nil {
			attributes["id"] = types.StringValue(response.EncryptionSettingsBackendResponse.Id)
			attributes["type"] = types.StringValue("encryption-settings")
		}
		if response.LdifBackendResponse != nil {
			attributes["id"] = types.StringValue(response.LdifBackendResponse.Id)
			attributes["type"] = types.StringValue("ldif")
		}
		if response.TrustStoreBackendResponse != nil {
			attributes["id"] = types.StringValue(response.TrustStoreBackendResponse.Id)
			attributes["type"] = types.StringValue("trust-store")
		}
		if response.CustomBackendResponse != nil {
			attributes["id"] = types.StringValue(response.CustomBackendResponse.Id)
			attributes["type"] = types.StringValue("custom")
		}
		if response.ChangelogBackendResponse != nil {
			attributes["id"] = types.StringValue(response.ChangelogBackendResponse.Id)
			attributes["type"] = types.StringValue("changelog")
		}
		if response.MonitorBackendResponse != nil {
			attributes["id"] = types.StringValue(response.MonitorBackendResponse.Id)
			attributes["type"] = types.StringValue("monitor")
		}
		if response.LocalDbBackendResponse != nil {
			attributes["id"] = types.StringValue(response.LocalDbBackendResponse.Id)
			attributes["type"] = types.StringValue("local-db")
		}
		if response.MirroredLdifBackendResponse != nil {
			attributes["id"] = types.StringValue(response.MirroredLdifBackendResponse.Id)
			attributes["type"] = types.StringValue("mirrored-ldif")
		}
		if response.ConfigFileHandlerBackendResponse != nil {
			attributes["id"] = types.StringValue(response.ConfigFileHandlerBackendResponse.Id)
			attributes["type"] = types.StringValue("config-file-handler")
		}
		if response.TaskBackendResponse != nil {
			attributes["id"] = types.StringValue(response.TaskBackendResponse.Id)
			attributes["type"] = types.StringValue("task")
		}
		if response.AlertBackendResponse != nil {
			attributes["id"] = types.StringValue(response.AlertBackendResponse.Id)
			attributes["type"] = types.StringValue("alert")
		}
		if response.AlarmBackendResponse != nil {
			attributes["id"] = types.StringValue(response.AlarmBackendResponse.Id)
			attributes["type"] = types.StringValue("alarm")
		}
		if response.MetricsBackendResponse != nil {
			attributes["id"] = types.StringValue(response.MetricsBackendResponse.Id)
			attributes["type"] = types.StringValue("metrics")
		}
		if response.LargeAttributeBackendResponse != nil {
			attributes["id"] = types.StringValue(response.LargeAttributeBackendResponse.Id)
			attributes["type"] = types.StringValue("large-attribute")
		}
		if response.CannedResponseBackendResponse != nil {
			attributes["id"] = types.StringValue(response.CannedResponseBackendResponse.Id)
			attributes["type"] = types.StringValue("canned-response")
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
