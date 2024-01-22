package delegatedadminattributecategory

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
	_ datasource.DataSource              = &delegatedAdminAttributeCategoryDataSource{}
	_ datasource.DataSourceWithConfigure = &delegatedAdminAttributeCategoryDataSource{}
)

// Create a Delegated Admin Attribute Category data source
func NewDelegatedAdminAttributeCategoryDataSource() datasource.DataSource {
	return &delegatedAdminAttributeCategoryDataSource{}
}

// delegatedAdminAttributeCategoryDataSource is the datasource implementation.
type delegatedAdminAttributeCategoryDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *delegatedAdminAttributeCategoryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delegated_admin_attribute_category"
}

// Configure adds the provider configured client to the data source.
func (r *delegatedAdminAttributeCategoryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type delegatedAdminAttributeCategoryDataSourceModel struct {
	Id                types.String `tfsdk:"id"`
	Type              types.String `tfsdk:"type"`
	Description       types.String `tfsdk:"description"`
	DisplayName       types.String `tfsdk:"display_name"`
	DisplayOrderIndex types.Int64  `tfsdk:"display_order_index"`
}

// GetSchema defines the schema for the datasource.
func (r *delegatedAdminAttributeCategoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Delegated Admin Attribute Category.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Delegated Admin Attribute Category resource. Options are ['delegated-admin-attribute-category']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Delegated Admin Attribute Category",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "A human readable display name for this Delegated Admin Attribute Category.",
				Required:    true,
			},
			"display_order_index": schema.Int64Attribute{
				Description: "Delegated Admin Attribute Categories are ordered for display based on this index from least to greatest.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a DelegatedAdminAttributeCategoryResponse object into the model struct
func readDelegatedAdminAttributeCategoryResponseDataSource(ctx context.Context, r *client.DelegatedAdminAttributeCategoryResponse, state *delegatedAdminAttributeCategoryDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("delegated-admin-attribute-category")
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.DisplayName = types.StringValue(r.DisplayName)
	state.DisplayOrderIndex = types.Int64Value(r.DisplayOrderIndex)
}

// Read resource information
func (r *delegatedAdminAttributeCategoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state delegatedAdminAttributeCategoryDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DelegatedAdminAttributeCategoryAPI.GetDelegatedAdminAttributeCategory(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.DisplayName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Attribute Category", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDelegatedAdminAttributeCategoryResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
