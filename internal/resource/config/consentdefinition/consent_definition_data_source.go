package consentdefinition

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
	_ datasource.DataSource              = &consentDefinitionDataSource{}
	_ datasource.DataSourceWithConfigure = &consentDefinitionDataSource{}
)

// Create a Consent Definition data source
func NewConsentDefinitionDataSource() datasource.DataSource {
	return &consentDefinitionDataSource{}
}

// consentDefinitionDataSource is the datasource implementation.
type consentDefinitionDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *consentDefinitionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_consent_definition"
}

// Configure adds the provider configured client to the data source.
func (r *consentDefinitionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type consentDefinitionDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Type        types.String `tfsdk:"type"`
	UniqueID    types.String `tfsdk:"unique_id"`
	DisplayName types.String `tfsdk:"display_name"`
	Parameter   types.Set    `tfsdk:"parameter"`
	Description types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the datasource.
func (r *consentDefinitionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Consent Definition.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Consent Definition resource. Options are ['consent-definition']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"unique_id": schema.StringAttribute{
				Description: "A version-independent unique identifier for this Consent Definition.",
				Required:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "A human-readable display name for this Consent Definition.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"parameter": schema.SetAttribute{
				Description: "Optional parameters for this Consent Definition.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Consent Definition",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a ConsentDefinitionResponse object into the model struct
func readConsentDefinitionResponseDataSource(ctx context.Context, r *client.ConsentDefinitionResponse, state *consentDefinitionDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("consent-definition")
	state.Id = types.StringValue(r.Id)
	state.UniqueID = types.StringValue(r.UniqueID)
	state.DisplayName = internaltypes.StringTypeOrNil(r.DisplayName, false)
	state.Parameter = internaltypes.GetStringSet(r.Parameter)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read resource information
func (r *consentDefinitionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state consentDefinitionDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConsentDefinitionAPI.GetConsentDefinition(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.UniqueID.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Consent Definition", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readConsentDefinitionResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
