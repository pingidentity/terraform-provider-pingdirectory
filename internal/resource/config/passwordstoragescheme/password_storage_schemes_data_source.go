// Copyright © 2025 Ping Identity Corporation

package passwordstoragescheme

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
	_ datasource.DataSource              = &passwordStorageSchemesDataSource{}
	_ datasource.DataSourceWithConfigure = &passwordStorageSchemesDataSource{}
)

// Create a Password Storage Schemes data source
func NewPasswordStorageSchemesDataSource() datasource.DataSource {
	return &passwordStorageSchemesDataSource{}
}

// passwordStorageSchemesDataSource is the datasource implementation.
type passwordStorageSchemesDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *passwordStorageSchemesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_storage_schemes"
}

// Configure adds the provider configured client to the data source.
func (r *passwordStorageSchemesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type passwordStorageSchemesDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *passwordStorageSchemesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Password Storage Scheme objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Password Storage Scheme objects found in the configuration",
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
func (r *passwordStorageSchemesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state passwordStorageSchemesDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.PasswordStorageSchemeAPI.ListPasswordStorageSchemes(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.PasswordStorageSchemeAPI.ListPasswordStorageSchemesExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Password Storage Scheme objects", err, httpResp)
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
		if response.SaltedSha256PasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.SaltedSha256PasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("salted-sha256")
		}
		if response.Argon2dPasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.Argon2dPasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("argon2d")
		}
		if response.CryptPasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.CryptPasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("crypt")
		}
		if response.Argon2iPasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.Argon2iPasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("argon2i")
		}
		if response.Base64PasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.Base64PasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("base64")
		}
		if response.SaltedMd5PasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.SaltedMd5PasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("salted-md5")
		}
		if response.AesPasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.AesPasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("aes")
		}
		if response.Argon2idPasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.Argon2idPasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("argon2id")
		}
		if response.VaultPasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.VaultPasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("vault")
		}
		if response.ThirdPartyPasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyPasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("third-party")
		}
		if response.Argon2PasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.Argon2PasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("argon2")
		}
		if response.ThirdPartyEnhancedPasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyEnhancedPasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("third-party-enhanced")
		}
		if response.Pbkdf2PasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.Pbkdf2PasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("pbkdf2")
		}
		if response.Rc4PasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.Rc4PasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("rc4")
		}
		if response.SaltedSha384PasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.SaltedSha384PasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("salted-sha384")
		}
		if response.TripleDesPasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.TripleDesPasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("triple-des")
		}
		if response.ClearPasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.ClearPasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("clear")
		}
		if response.Aes256PasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.Aes256PasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("aes-256")
		}
		if response.BcryptPasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.BcryptPasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("bcrypt")
		}
		if response.BlowfishPasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.BlowfishPasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("blowfish")
		}
		if response.Sha1PasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.Sha1PasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("sha1")
		}
		if response.AmazonSecretsManagerPasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.AmazonSecretsManagerPasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("amazon-secrets-manager")
		}
		if response.AzureKeyVaultPasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.AzureKeyVaultPasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("azure-key-vault")
		}
		if response.ConjurPasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.ConjurPasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("conjur")
		}
		if response.SaltedSha1PasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.SaltedSha1PasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("salted-sha1")
		}
		if response.SaltedSha512PasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.SaltedSha512PasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("salted-sha512")
		}
		if response.ScryptPasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.ScryptPasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("scrypt")
		}
		if response.Md5PasswordStorageSchemeResponse != nil {
			attributes["id"] = types.StringValue(response.Md5PasswordStorageSchemeResponse.Id)
			attributes["type"] = types.StringValue("md5")
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
