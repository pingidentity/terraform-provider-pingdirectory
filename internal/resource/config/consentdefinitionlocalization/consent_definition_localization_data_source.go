package consentdefinitionlocalization

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
	_ datasource.DataSource              = &consentDefinitionLocalizationDataSource{}
	_ datasource.DataSourceWithConfigure = &consentDefinitionLocalizationDataSource{}
)

// Create a Consent Definition Localization data source
func NewConsentDefinitionLocalizationDataSource() datasource.DataSource {
	return &consentDefinitionLocalizationDataSource{}
}

// consentDefinitionLocalizationDataSource is the datasource implementation.
type consentDefinitionLocalizationDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *consentDefinitionLocalizationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_consent_definition_localization"
}

// Configure adds the provider configured client to the data source.
func (r *consentDefinitionLocalizationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type consentDefinitionLocalizationDataSourceModel struct {
	Id                    types.String `tfsdk:"id"`
	ConsentDefinitionName types.String `tfsdk:"consent_definition_name"`
	Locale                types.String `tfsdk:"locale"`
	Version               types.String `tfsdk:"version"`
	TitleText             types.String `tfsdk:"title_text"`
	DataText              types.String `tfsdk:"data_text"`
	PurposeText           types.String `tfsdk:"purpose_text"`
}

// GetSchema defines the schema for the datasource.
func (r *consentDefinitionLocalizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Describes a Consent Definition Localization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Name of this object.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"consent_definition_name": schema.StringAttribute{
				Description: "Name of the parent Consent Definition",
				Required:    true,
			},
			"locale": schema.StringAttribute{
				Description: "The locale of this Consent Definition Localization.",
				Required:    true,
			},
			"version": schema.StringAttribute{
				Description: "The version of this Consent Definition Localization, using the format MAJOR.MINOR.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"title_text": schema.StringAttribute{
				Description: "Localized text that may be used to provide a title or summary for a consent request or a granted consent.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"data_text": schema.StringAttribute{
				Description: "Localized text describing the data to be shared.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"purpose_text": schema.StringAttribute{
				Description: "Localized text describing how the data is to be used.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
}

// Read a ConsentDefinitionLocalizationResponse object into the model struct
func readConsentDefinitionLocalizationResponseDataSource(ctx context.Context, r *client.ConsentDefinitionLocalizationResponse, state *consentDefinitionLocalizationDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Locale = types.StringValue(r.Locale)
	state.Version = types.StringValue(r.Version)
	state.TitleText = internaltypes.StringTypeOrNil(r.TitleText, false)
	state.DataText = types.StringValue(r.DataText)
	state.PurposeText = types.StringValue(r.PurposeText)
}

// Read resource information
func (r *consentDefinitionLocalizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state consentDefinitionLocalizationDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConsentDefinitionLocalizationApi.GetConsentDefinitionLocalization(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Locale.ValueString(), state.ConsentDefinitionName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Consent Definition Localization", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readConsentDefinitionLocalizationResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
