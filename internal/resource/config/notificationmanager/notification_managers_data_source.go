// Copyright © 2025 Ping Identity Corporation

package notificationmanager

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &notificationManagersDataSource{}
	_ datasource.DataSourceWithConfigure = &notificationManagersDataSource{}
)

// Create a Notification Managers data source
func NewNotificationManagersDataSource() datasource.DataSource {
	return &notificationManagersDataSource{}
}

// notificationManagersDataSource is the datasource implementation.
type notificationManagersDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *notificationManagersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_managers"
}

// Configure adds the provider configured client to the data source.
func (r *notificationManagersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type notificationManagersDataSourceModel struct {
	Id     types.String `tfsdk:"id"`
	Filter types.String `tfsdk:"filter"`
	Ids    types.Set    `tfsdk:"ids"`
}

// GetSchema defines the schema for the datasource.
func (r *notificationManagersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Notification Manager objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"ids": schema.SetAttribute{
				Description: "Notification Manager IDs found in the configuration",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read resource information
func (r *notificationManagersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state notificationManagersDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.NotificationManagerAPI.ListNotificationManagers(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.NotificationManagerAPI.ListNotificationManagersExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Notification Manager objects", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	ids := []attr.Value{}
	for _, response := range readResponse.Resources {
		ids = append(ids, types.StringValue(response.Id))
	}

	state.Ids, diags = types.SetValue(types.StringType, ids)
	resp.Diagnostics.Append(diags...)
	state.Id = types.StringValue("id")

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
