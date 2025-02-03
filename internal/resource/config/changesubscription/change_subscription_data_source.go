// Copyright © 2025 Ping Identity Corporation

package changesubscription

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &changeSubscriptionDataSource{}
	_ datasource.DataSourceWithConfigure = &changeSubscriptionDataSource{}
)

// Create a Change Subscription data source
func NewChangeSubscriptionDataSource() datasource.DataSource {
	return &changeSubscriptionDataSource{}
}

// changeSubscriptionDataSource is the datasource implementation.
type changeSubscriptionDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *changeSubscriptionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_change_subscription"
}

// Configure adds the provider configured client to the data source.
func (r *changeSubscriptionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type changeSubscriptionDataSourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Type               types.String `tfsdk:"type"`
	Description        types.String `tfsdk:"description"`
	ConnectionCriteria types.String `tfsdk:"connection_criteria"`
	RequestCriteria    types.String `tfsdk:"request_criteria"`
	ResultCriteria     types.String `tfsdk:"result_criteria"`
	ExpirationTime     types.String `tfsdk:"expiration_time"`
}

// GetSchema defines the schema for the datasource.
func (r *changeSubscriptionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Change Subscription.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Change Subscription resource. Options are ['change-subscription']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Change Subscription",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"connection_criteria": schema.StringAttribute{
				Description: "Specifies a set of connection criteria that must match the client connection associated with an operation in order for that operation to be processed by a change subscription handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"request_criteria": schema.StringAttribute{
				Description: "Specifies a set of request criteria that must match the request associated with an operation in order for that operation to be processed by a change subscription handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"result_criteria": schema.StringAttribute{
				Description: "Specifies a set of result criteria that must match the result associated with an operation in order for that operation to be processed by a change subscription handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"expiration_time": schema.StringAttribute{
				Description: "Specifies a timestamp that provides an expiration time for this change subscription. If an expiration time is provided, then the change subscription will not be active after that time has passed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a ChangeSubscriptionResponse object into the model struct
func readChangeSubscriptionResponseDataSource(ctx context.Context, r *client.ChangeSubscriptionResponse, state *changeSubscriptionDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("change-subscription")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, false)
	state.ExpirationTime = internaltypes.StringTypeOrNil(r.ExpirationTime, false)
}

// Read resource information
func (r *changeSubscriptionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state changeSubscriptionDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ChangeSubscriptionAPI.GetChangeSubscription(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Change Subscription", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readChangeSubscriptionResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
