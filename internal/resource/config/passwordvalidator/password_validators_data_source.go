// Copyright © 2025 Ping Identity Corporation

package passwordvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &passwordValidatorsDataSource{}
	_ datasource.DataSourceWithConfigure = &passwordValidatorsDataSource{}
)

// Create a Password Validators data source
func NewPasswordValidatorsDataSource() datasource.DataSource {
	return &passwordValidatorsDataSource{}
}

// passwordValidatorsDataSource is the datasource implementation.
type passwordValidatorsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *passwordValidatorsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_validators"
}

// Configure adds the provider configured client to the data source.
func (r *passwordValidatorsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type passwordValidatorsDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *passwordValidatorsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Password Validator objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Password Validator objects found in the configuration",
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
func (r *passwordValidatorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state passwordValidatorsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.PasswordValidatorAPI.ListPasswordValidators(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.PasswordValidatorAPI.ListPasswordValidatorsExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Password Validator objects", err, httpResp)
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
		if response.CharacterSetPasswordValidatorResponse != nil {
			attributes["id"] = types.StringValue(response.CharacterSetPasswordValidatorResponse.Id)
			attributes["type"] = types.StringValue("character-set")
		}
		if response.SimilarityBasedPasswordValidatorResponse != nil {
			attributes["id"] = types.StringValue(response.SimilarityBasedPasswordValidatorResponse.Id)
			attributes["type"] = types.StringValue("similarity-based")
		}
		if response.AttributeValuePasswordValidatorResponse != nil {
			attributes["id"] = types.StringValue(response.AttributeValuePasswordValidatorResponse.Id)
			attributes["type"] = types.StringValue("attribute-value")
		}
		if response.CustomPasswordValidatorResponse != nil {
			attributes["id"] = types.StringValue(response.CustomPasswordValidatorResponse.Id)
			attributes["type"] = types.StringValue("custom")
		}
		if response.RepeatedCharactersPasswordValidatorResponse != nil {
			attributes["id"] = types.StringValue(response.RepeatedCharactersPasswordValidatorResponse.Id)
			attributes["type"] = types.StringValue("repeated-characters")
		}
		if response.DictionaryPasswordValidatorResponse != nil {
			attributes["id"] = types.StringValue(response.DictionaryPasswordValidatorResponse.Id)
			attributes["type"] = types.StringValue("dictionary")
		}
		if response.HaystackPasswordValidatorResponse != nil {
			attributes["id"] = types.StringValue(response.HaystackPasswordValidatorResponse.Id)
			attributes["type"] = types.StringValue("haystack")
		}
		if response.Utf8PasswordValidatorResponse != nil {
			attributes["id"] = types.StringValue(response.Utf8PasswordValidatorResponse.Id)
			attributes["type"] = types.StringValue("utf-8")
		}
		if response.GroovyScriptedPasswordValidatorResponse != nil {
			attributes["id"] = types.StringValue(response.GroovyScriptedPasswordValidatorResponse.Id)
			attributes["type"] = types.StringValue("groovy-scripted")
		}
		if response.PwnedPasswordsPasswordValidatorResponse != nil {
			attributes["id"] = types.StringValue(response.PwnedPasswordsPasswordValidatorResponse.Id)
			attributes["type"] = types.StringValue("pwned-passwords")
		}
		if response.DisallowedCharactersPasswordValidatorResponse != nil {
			attributes["id"] = types.StringValue(response.DisallowedCharactersPasswordValidatorResponse.Id)
			attributes["type"] = types.StringValue("disallowed-characters")
		}
		if response.LengthBasedPasswordValidatorResponse != nil {
			attributes["id"] = types.StringValue(response.LengthBasedPasswordValidatorResponse.Id)
			attributes["type"] = types.StringValue("length-based")
		}
		if response.RegularExpressionPasswordValidatorResponse != nil {
			attributes["id"] = types.StringValue(response.RegularExpressionPasswordValidatorResponse.Id)
			attributes["type"] = types.StringValue("regular-expression")
		}
		if response.UniqueCharactersPasswordValidatorResponse != nil {
			attributes["id"] = types.StringValue(response.UniqueCharactersPasswordValidatorResponse.Id)
			attributes["type"] = types.StringValue("unique-characters")
		}
		if response.ThirdPartyPasswordValidatorResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyPasswordValidatorResponse.Id)
			attributes["type"] = types.StringValue("third-party")
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
