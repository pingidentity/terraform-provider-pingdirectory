package logretentionpolicy

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
	_ datasource.DataSource              = &logRetentionPoliciesDataSource{}
	_ datasource.DataSourceWithConfigure = &logRetentionPoliciesDataSource{}
)

// Create a Log Retention Policies data source
func NewLogRetentionPoliciesDataSource() datasource.DataSource {
	return &logRetentionPoliciesDataSource{}
}

// logRetentionPoliciesDataSource is the datasource implementation.
type logRetentionPoliciesDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *logRetentionPoliciesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_retention_policies"
}

// Configure adds the provider configured client to the data source.
func (r *logRetentionPoliciesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type logRetentionPoliciesDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *logRetentionPoliciesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Log Retention Policy objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Log Retention Policy objects found in the configuration",
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
func (r *logRetentionPoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state logRetentionPoliciesDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.LogRetentionPolicyAPI.ListLogRetentionPolicies(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.LogRetentionPolicyAPI.ListLogRetentionPoliciesExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Log Retention Policy objects", err, httpResp)
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
		if response.TimeLimitLogRetentionPolicyResponse != nil {
			attributes["id"] = types.StringValue(response.TimeLimitLogRetentionPolicyResponse.Id)
			attributes["type"] = types.StringValue("time-limit")
		}
		if response.NeverDeleteLogRetentionPolicyResponse != nil {
			attributes["id"] = types.StringValue(response.NeverDeleteLogRetentionPolicyResponse.Id)
			attributes["type"] = types.StringValue("never-delete")
		}
		if response.FileCountLogRetentionPolicyResponse != nil {
			attributes["id"] = types.StringValue(response.FileCountLogRetentionPolicyResponse.Id)
			attributes["type"] = types.StringValue("file-count")
		}
		if response.FreeDiskSpaceLogRetentionPolicyResponse != nil {
			attributes["id"] = types.StringValue(response.FreeDiskSpaceLogRetentionPolicyResponse.Id)
			attributes["type"] = types.StringValue("free-disk-space")
		}
		if response.SizeLimitLogRetentionPolicyResponse != nil {
			attributes["id"] = types.StringValue(response.SizeLimitLogRetentionPolicyResponse.Id)
			attributes["type"] = types.StringValue("size-limit")
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
