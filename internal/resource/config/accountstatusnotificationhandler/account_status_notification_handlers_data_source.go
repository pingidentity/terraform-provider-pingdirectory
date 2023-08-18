package accountstatusnotificationhandler

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
	_ datasource.DataSource              = &accountStatusNotificationHandlersDataSource{}
	_ datasource.DataSourceWithConfigure = &accountStatusNotificationHandlersDataSource{}
)

// Create a Account Status Notification Handlers data source
func NewAccountStatusNotificationHandlersDataSource() datasource.DataSource {
	return &accountStatusNotificationHandlersDataSource{}
}

// accountStatusNotificationHandlersDataSource is the datasource implementation.
type accountStatusNotificationHandlersDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *accountStatusNotificationHandlersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_status_notification_handlers"
}

// Configure adds the provider configured client to the data source.
func (r *accountStatusNotificationHandlersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type accountStatusNotificationHandlersDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *accountStatusNotificationHandlersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Account Status Notification Handler objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Account Status Notification Handler objects found in the configuration",
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
func (r *accountStatusNotificationHandlersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state accountStatusNotificationHandlersDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.AccountStatusNotificationHandlerApi.ListAccountStatusNotificationHandlers(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.ListAccountStatusNotificationHandlersExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Account Status Notification Handler objects", err, httpResp)
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
		if response.SmtpAccountStatusNotificationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.SmtpAccountStatusNotificationHandlerResponse.Id)
			attributes["type"] = types.StringValue("smtp")
		}
		if response.GroovyScriptedAccountStatusNotificationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.GroovyScriptedAccountStatusNotificationHandlerResponse.Id)
			attributes["type"] = types.StringValue("groovy-scripted")
		}
		if response.AdminAlertAccountStatusNotificationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.AdminAlertAccountStatusNotificationHandlerResponse.Id)
			attributes["type"] = types.StringValue("admin-alert")
		}
		if response.ErrorLogAccountStatusNotificationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.ErrorLogAccountStatusNotificationHandlerResponse.Id)
			attributes["type"] = types.StringValue("error-log")
		}
		if response.MultiPartEmailAccountStatusNotificationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.MultiPartEmailAccountStatusNotificationHandlerResponse.Id)
			attributes["type"] = types.StringValue("multi-part-email")
		}
		if response.ThirdPartyAccountStatusNotificationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyAccountStatusNotificationHandlerResponse.Id)
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
