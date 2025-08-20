// Copyright © 2025 Ping Identity Corporation

package tokenclaimvalidation

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &tokenClaimValidationsDataSource{}
	_ datasource.DataSourceWithConfigure = &tokenClaimValidationsDataSource{}
)

// Create a Token Claim Validations data source
func NewTokenClaimValidationsDataSource() datasource.DataSource {
	return &tokenClaimValidationsDataSource{}
}

// tokenClaimValidationsDataSource is the datasource implementation.
type tokenClaimValidationsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *tokenClaimValidationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_token_claim_validations"
}

// Configure adds the provider configured client to the data source.
func (r *tokenClaimValidationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type tokenClaimValidationsDataSourceModel struct {
	Id                   types.String `tfsdk:"id"`
	Filter               types.String `tfsdk:"filter"`
	Objects              types.Set    `tfsdk:"objects"`
	IdTokenValidatorName types.String `tfsdk:"id_token_validator_name"`
}

// GetSchema defines the schema for the datasource.
func (r *tokenClaimValidationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Token Claim Validation objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"id_token_validator_name": schema.StringAttribute{
				Description: "Name of the parent ID Token Validator",
				Required:    true,
			},
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Token Claim Validation objects found in the configuration",
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
func (r *tokenClaimValidationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state tokenClaimValidationsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.TokenClaimValidationAPI.ListTokenClaimValidations(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.IdTokenValidatorName.ValueString())
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.TokenClaimValidationAPI.ListTokenClaimValidationsExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Token Claim Validation objects", err, httpResp)
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
		if response.StringArrayTokenClaimValidationResponse != nil {
			attributes["id"] = types.StringValue(response.StringArrayTokenClaimValidationResponse.Id)
			attributes["type"] = types.StringValue("string-array")
		}
		if response.BooleanTokenClaimValidationResponse != nil {
			attributes["id"] = types.StringValue(response.BooleanTokenClaimValidationResponse.Id)
			attributes["type"] = types.StringValue("boolean")
		}
		if response.StringTokenClaimValidationResponse != nil {
			attributes["id"] = types.StringValue(response.StringTokenClaimValidationResponse.Id)
			attributes["type"] = types.StringValue("string")
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
