package tokenclaimvalidation

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10100/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &tokenClaimValidationDataSource{}
	_ datasource.DataSourceWithConfigure = &tokenClaimValidationDataSource{}
)

// Create a Token Claim Validation data source
func NewTokenClaimValidationDataSource() datasource.DataSource {
	return &tokenClaimValidationDataSource{}
}

// tokenClaimValidationDataSource is the datasource implementation.
type tokenClaimValidationDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *tokenClaimValidationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_token_claim_validation"
}

// Configure adds the provider configured client to the data source.
func (r *tokenClaimValidationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type tokenClaimValidationDataSourceModel struct {
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Type                 types.String `tfsdk:"type"`
	IdTokenValidatorName types.String `tfsdk:"id_token_validator_name"`
	RequiredValue        types.String `tfsdk:"required_value"`
	AllRequiredValue     types.Set    `tfsdk:"all_required_value"`
	AnyRequiredValue     types.Set    `tfsdk:"any_required_value"`
	Description          types.String `tfsdk:"description"`
	ClaimName            types.String `tfsdk:"claim_name"`
}

// GetSchema defines the schema for the datasource.
func (r *tokenClaimValidationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Token Claim Validation.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Token Claim Validation resource. Options are ['string-array', 'boolean', 'string']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"id_token_validator_name": schema.StringAttribute{
				Description: "Name of the parent ID Token Validator",
				Required:    true,
			},
			"required_value": schema.StringAttribute{
				Description: "Specifies the boolean claim's required value.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"all_required_value": schema.SetAttribute{
				Description: "The set of all values that the claim must have to be considered valid.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_required_value": schema.SetAttribute{
				Description: "The set of values that the claim may have to be considered valid.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Token Claim Validation",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"claim_name": schema.StringAttribute{
				Description: "The name of the claim to be validated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a StringArrayTokenClaimValidationResponse object into the model struct
func readStringArrayTokenClaimValidationResponseDataSource(ctx context.Context, r *client.StringArrayTokenClaimValidationResponse, state *tokenClaimValidationDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("string-array")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllRequiredValue = internaltypes.GetStringSet(r.AllRequiredValue)
	state.AnyRequiredValue = internaltypes.GetStringSet(r.AnyRequiredValue)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.ClaimName = types.StringValue(r.ClaimName)
}

// Read a BooleanTokenClaimValidationResponse object into the model struct
func readBooleanTokenClaimValidationResponseDataSource(ctx context.Context, r *client.BooleanTokenClaimValidationResponse, state *tokenClaimValidationDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("boolean")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.RequiredValue = types.StringValue(r.RequiredValue.String())
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.ClaimName = types.StringValue(r.ClaimName)
}

// Read a StringTokenClaimValidationResponse object into the model struct
func readStringTokenClaimValidationResponseDataSource(ctx context.Context, r *client.StringTokenClaimValidationResponse, state *tokenClaimValidationDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("string")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AnyRequiredValue = internaltypes.GetStringSet(r.AnyRequiredValue)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.ClaimName = types.StringValue(r.ClaimName)
}

// Read resource information
func (r *tokenClaimValidationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state tokenClaimValidationDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.TokenClaimValidationAPI.GetTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.IdTokenValidatorName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Token Claim Validation", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.StringArrayTokenClaimValidationResponse != nil {
		readStringArrayTokenClaimValidationResponseDataSource(ctx, readResponse.StringArrayTokenClaimValidationResponse, &state, &resp.Diagnostics)
	}
	if readResponse.BooleanTokenClaimValidationResponse != nil {
		readBooleanTokenClaimValidationResponseDataSource(ctx, readResponse.BooleanTokenClaimValidationResponse, &state, &resp.Diagnostics)
	}
	if readResponse.StringTokenClaimValidationResponse != nil {
		readStringTokenClaimValidationResponseDataSource(ctx, readResponse.StringTokenClaimValidationResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
