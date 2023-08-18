package matchingrule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &matchingRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &matchingRuleDataSource{}
)

// Create a Matching Rule data source
func NewMatchingRuleDataSource() datasource.DataSource {
	return &matchingRuleDataSource{}
}

// matchingRuleDataSource is the datasource implementation.
type matchingRuleDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *matchingRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_matching_rule"
}

// Configure adds the provider configured client to the data source.
func (r *matchingRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type matchingRuleDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Type    types.String `tfsdk:"type"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *matchingRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Matching Rule.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Matching Rule resource. Options are ['ordering', 'approximate', 'equality', 'substring', 'generic']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Matching Rule is enabled for use.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a OrderingMatchingRuleResponse object into the model struct
func readOrderingMatchingRuleResponseDataSource(ctx context.Context, r *client.OrderingMatchingRuleResponse, state *matchingRuleDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ordering")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ApproximateMatchingRuleResponse object into the model struct
func readApproximateMatchingRuleResponseDataSource(ctx context.Context, r *client.ApproximateMatchingRuleResponse, state *matchingRuleDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("approximate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a EqualityMatchingRuleResponse object into the model struct
func readEqualityMatchingRuleResponseDataSource(ctx context.Context, r *client.EqualityMatchingRuleResponse, state *matchingRuleDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("equality")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a SubstringMatchingRuleResponse object into the model struct
func readSubstringMatchingRuleResponseDataSource(ctx context.Context, r *client.SubstringMatchingRuleResponse, state *matchingRuleDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("substring")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a GenericMatchingRuleResponse object into the model struct
func readGenericMatchingRuleResponseDataSource(ctx context.Context, r *client.GenericMatchingRuleResponse, state *matchingRuleDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generic")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *matchingRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state matchingRuleDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.MatchingRuleApi.GetMatchingRule(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Matching Rule", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.OrderingMatchingRuleResponse != nil {
		readOrderingMatchingRuleResponseDataSource(ctx, readResponse.OrderingMatchingRuleResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ApproximateMatchingRuleResponse != nil {
		readApproximateMatchingRuleResponseDataSource(ctx, readResponse.ApproximateMatchingRuleResponse, &state, &resp.Diagnostics)
	}
	if readResponse.EqualityMatchingRuleResponse != nil {
		readEqualityMatchingRuleResponseDataSource(ctx, readResponse.EqualityMatchingRuleResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SubstringMatchingRuleResponse != nil {
		readSubstringMatchingRuleResponseDataSource(ctx, readResponse.SubstringMatchingRuleResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GenericMatchingRuleResponse != nil {
		readGenericMatchingRuleResponseDataSource(ctx, readResponse.GenericMatchingRuleResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
