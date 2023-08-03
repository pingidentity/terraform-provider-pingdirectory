package scimschema

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
	_ datasource.DataSource              = &scimSchemaDataSource{}
	_ datasource.DataSourceWithConfigure = &scimSchemaDataSource{}
)

// Create a Scim Schema data source
func NewScimSchemaDataSource() datasource.DataSource {
	return &scimSchemaDataSource{}
}

// scimSchemaDataSource is the datasource implementation.
type scimSchemaDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *scimSchemaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scim_schema"
}

// Configure adds the provider configured client to the data source.
func (r *scimSchemaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type scimSchemaDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	SchemaURN   types.String `tfsdk:"schema_urn"`
	DisplayName types.String `tfsdk:"display_name"`
}

// GetSchema defines the schema for the datasource.
func (r *scimSchemaDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Scim Schema.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description: "A description for this SCIM Schema",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"schema_urn": schema.StringAttribute{
				Description: "The URN which identifies this SCIM Schema.",
				Required:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "The human readable name for this SCIM Schema.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a ScimSchemaResponse object into the model struct
func readScimSchemaResponseDataSource(ctx context.Context, r *client.ScimSchemaResponse, state *scimSchemaDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.SchemaURN = types.StringValue(r.SchemaURN)
	state.DisplayName = internaltypes.StringTypeOrNil(r.DisplayName, false)
}

// Read resource information
func (r *scimSchemaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state scimSchemaDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ScimSchemaApi.GetScimSchema(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.SchemaURN.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Scim Schema", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readScimSchemaResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
