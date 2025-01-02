package recurringtask

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
	_ datasource.DataSource              = &recurringTasksDataSource{}
	_ datasource.DataSourceWithConfigure = &recurringTasksDataSource{}
)

// Create a Recurring Tasks data source
func NewRecurringTasksDataSource() datasource.DataSource {
	return &recurringTasksDataSource{}
}

// recurringTasksDataSource is the datasource implementation.
type recurringTasksDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *recurringTasksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_recurring_tasks"
}

// Configure adds the provider configured client to the data source.
func (r *recurringTasksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type recurringTasksDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *recurringTasksDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Recurring Task objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Recurring Task objects found in the configuration",
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
func (r *recurringTasksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state recurringTasksDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.RecurringTaskAPI.ListRecurringTasks(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.RecurringTaskAPI.ListRecurringTasksExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Recurring Task objects", err, httpResp)
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
		if response.GenerateServerProfileRecurringTaskResponse != nil {
			attributes["id"] = types.StringValue(response.GenerateServerProfileRecurringTaskResponse.Id)
			attributes["type"] = types.StringValue("generate-server-profile")
		}
		if response.LeaveLockdownModeRecurringTaskResponse != nil {
			attributes["id"] = types.StringValue(response.LeaveLockdownModeRecurringTaskResponse.Id)
			attributes["type"] = types.StringValue("leave-lockdown-mode")
		}
		if response.BackupRecurringTaskResponse != nil {
			attributes["id"] = types.StringValue(response.BackupRecurringTaskResponse.Id)
			attributes["type"] = types.StringValue("backup")
		}
		if response.StaticallyDefinedRecurringTaskResponse != nil {
			attributes["id"] = types.StringValue(response.StaticallyDefinedRecurringTaskResponse.Id)
			attributes["type"] = types.StringValue("statically-defined")
		}
		if response.CollectSupportDataRecurringTaskResponse != nil {
			attributes["id"] = types.StringValue(response.CollectSupportDataRecurringTaskResponse.Id)
			attributes["type"] = types.StringValue("collect-support-data")
		}
		if response.AuditDataSecurityRecurringTaskResponse != nil {
			attributes["id"] = types.StringValue(response.AuditDataSecurityRecurringTaskResponse.Id)
			attributes["type"] = types.StringValue("audit-data-security")
		}
		if response.DelayRecurringTaskResponse != nil {
			attributes["id"] = types.StringValue(response.DelayRecurringTaskResponse.Id)
			attributes["type"] = types.StringValue("delay")
		}
		if response.LdifExportRecurringTaskResponse != nil {
			attributes["id"] = types.StringValue(response.LdifExportRecurringTaskResponse.Id)
			attributes["type"] = types.StringValue("ldif-export")
		}
		if response.EnterLockdownModeRecurringTaskResponse != nil {
			attributes["id"] = types.StringValue(response.EnterLockdownModeRecurringTaskResponse.Id)
			attributes["type"] = types.StringValue("enter-lockdown-mode")
		}
		if response.ExecRecurringTaskResponse != nil {
			attributes["id"] = types.StringValue(response.ExecRecurringTaskResponse.Id)
			attributes["type"] = types.StringValue("exec")
		}
		if response.FileRetentionRecurringTaskResponse != nil {
			attributes["id"] = types.StringValue(response.FileRetentionRecurringTaskResponse.Id)
			attributes["type"] = types.StringValue("file-retention")
		}
		if response.ThirdPartyRecurringTaskResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyRecurringTaskResponse.Id)
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
