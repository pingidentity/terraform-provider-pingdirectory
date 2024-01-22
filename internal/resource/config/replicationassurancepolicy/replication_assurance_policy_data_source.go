package replicationassurancepolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &replicationAssurancePolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &replicationAssurancePolicyDataSource{}
)

// Create a Replication Assurance Policy data source
func NewReplicationAssurancePolicyDataSource() datasource.DataSource {
	return &replicationAssurancePolicyDataSource{}
}

// replicationAssurancePolicyDataSource is the datasource implementation.
type replicationAssurancePolicyDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *replicationAssurancePolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication_assurance_policy"
}

// Configure adds the provider configured client to the data source.
func (r *replicationAssurancePolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type replicationAssurancePolicyDataSourceModel struct {
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Type                 types.String `tfsdk:"type"`
	Description          types.String `tfsdk:"description"`
	Enabled              types.Bool   `tfsdk:"enabled"`
	EvaluationOrderIndex types.Int64  `tfsdk:"evaluation_order_index"`
	LocalLevel           types.String `tfsdk:"local_level"`
	RemoteLevel          types.String `tfsdk:"remote_level"`
	Timeout              types.String `tfsdk:"timeout"`
	ConnectionCriteria   types.String `tfsdk:"connection_criteria"`
	RequestCriteria      types.String `tfsdk:"request_criteria"`
}

// GetSchema defines the schema for the datasource.
func (r *replicationAssurancePolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Replication Assurance Policy.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Replication Assurance Policy resource. Options are ['replication-assurance-policy']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the Replication Assurance Policy.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Replication Assurance Policy is enabled for use in the server. If a Replication Assurance Policy is disabled, then no new operations will be associated with it.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"evaluation_order_index": schema.Int64Attribute{
				Description: "When multiple Replication Assurance Policies are defined, this property determines the evaluation order for finding a Replication Assurance Policy match against an operation. Policies are evaluated based on this index from least to greatest. Values of this property must be unique but not necessarily contiguous.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"local_level": schema.StringAttribute{
				Description: "Specifies the assurance level used to replicate to local servers. A local server is defined as one with the same value for the location setting in the global configuration.  The local-level must be set to an assurance level at least as strict as the remote-level. In other words, if remote-level is set to \"received-any-remote-location\" or \"received-all-remote-locations\", then local-level must be either \"received-any-server\" or \"processed-all-servers\". If remote-level is \"processed-all-remote-servers\", then local-level must be \"processed-all-servers\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"remote_level": schema.StringAttribute{
				Description: "Specifies the assurance level used to replicate to remote servers. A remote server is defined as one with a different value for the location setting in the global configuration.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"timeout": schema.StringAttribute{
				Description: "Specifies the maximum length of time to wait for the replication assurance requirements to be met before timing out and replying to the client.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"connection_criteria": schema.StringAttribute{
				Description: "Specifies a connection criteria used to indicate which operations from clients matching this criteria use this policy. If both a connection criteria and a request criteria are specified for a policy, then both must match an operation for the policy to be assigned.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"request_criteria": schema.StringAttribute{
				Description: "Specifies a request criteria used to indicate which operations from clients matching this criteria use this policy. If both a connection criteria and a request criteria are specified for a policy, then both must match an operation for the policy to be assigned.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a ReplicationAssurancePolicyResponse object into the model struct
func readReplicationAssurancePolicyResponseDataSource(ctx context.Context, r *client.ReplicationAssurancePolicyResponse, state *replicationAssurancePolicyDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("replication-assurance-policy")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.EvaluationOrderIndex = types.Int64Value(r.EvaluationOrderIndex)
	state.LocalLevel = types.StringValue(r.LocalLevel.String())
	state.RemoteLevel = types.StringValue(r.RemoteLevel.String())
	state.Timeout = types.StringValue(r.Timeout)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
}

// Read resource information
func (r *replicationAssurancePolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state replicationAssurancePolicyDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ReplicationAssurancePolicyAPI.GetReplicationAssurancePolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Replication Assurance Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readReplicationAssurancePolicyResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
