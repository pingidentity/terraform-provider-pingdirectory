package attributesyntax

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
	_ datasource.DataSource              = &attributeSyntaxDataSource{}
	_ datasource.DataSourceWithConfigure = &attributeSyntaxDataSource{}
)

// Create a Attribute Syntax data source
func NewAttributeSyntaxDataSource() datasource.DataSource {
	return &attributeSyntaxDataSource{}
}

// attributeSyntaxDataSource is the datasource implementation.
type attributeSyntaxDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *attributeSyntaxDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_attribute_syntax"
}

// Configure adds the provider configured client to the data source.
func (r *attributeSyntaxDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type attributeSyntaxDataSourceModel struct {
	Id                             types.String `tfsdk:"id"`
	Name                           types.String `tfsdk:"name"`
	Type                           types.String `tfsdk:"type"`
	EnableCompaction               types.Bool   `tfsdk:"enable_compaction"`
	IncludeAttributeInCompaction   types.Set    `tfsdk:"include_attribute_in_compaction"`
	ExcludeAttributeFromCompaction types.Set    `tfsdk:"exclude_attribute_from_compaction"`
	StrictFormat                   types.Bool   `tfsdk:"strict_format"`
	AllowZeroLengthValues          types.Bool   `tfsdk:"allow_zero_length_values"`
	StripSyntaxMinUpperBound       types.Bool   `tfsdk:"strip_syntax_min_upper_bound"`
	Enabled                        types.Bool   `tfsdk:"enabled"`
	RequireBinaryTransfer          types.Bool   `tfsdk:"require_binary_transfer"`
}

// GetSchema defines the schema for the datasource.
func (r *attributeSyntaxDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Attribute Syntax.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Attribute Syntax resource. Options are ['attribute-type-description', 'directory-string', 'telephone-number', 'distinguished-name', 'generalized-time', 'integer', 'uuid', 'generic', 'json-object', 'user-password', 'boolean', 'hex-string', 'bit-string', 'ldap-url', 'name-and-optional-uid']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enable_compaction": schema.BoolAttribute{
				Description: "Indicates whether values of attributes with this syntax should be compacted when stored in a local DB database.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_attribute_in_compaction": schema.SetAttribute{
				Description: "Specifies the specific attributes (which should be associated with this syntax) whose values should be compacted. If one or more include attributes are specified, then only those attributes will have their values compacted. If not set then all attributes will have their values compacted. The exclude-attribute-from-compaction property takes precedence over this property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"exclude_attribute_from_compaction": schema.SetAttribute{
				Description: "Specifies the specific attributes (which should be associated with this syntax) whose values should not be compacted. If one or more exclude attributes are specified, then values of those attributes will not have their values compacted. This property takes precedence over the include-attribute-in-compaction property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"strict_format": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to `telephone-number`: Indicates whether to require telephone number values to strictly comply with the standard definition for this syntax. When the `type` attribute is set to `ldap-url`: Indicates whether values for attributes with this syntax will be required to be in the valid LDAP URL format. If this is set to false, then arbitrary strings will be allowed.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `telephone-number`: Indicates whether to require telephone number values to strictly comply with the standard definition for this syntax.\n  - `ldap-url`: Indicates whether values for attributes with this syntax will be required to be in the valid LDAP URL format. If this is set to false, then arbitrary strings will be allowed.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"allow_zero_length_values": schema.BoolAttribute{
				Description: "Indicates whether zero-length (that is, an empty string) values are allowed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"strip_syntax_min_upper_bound": schema.BoolAttribute{
				Description: "Indicates whether the suggested minimum upper bound appended to an attribute's syntax OID in its schema definition Attribute Type Description should be stripped.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Attribute Syntax is enabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"require_binary_transfer": schema.BoolAttribute{
				Description: "Indicates whether values of this attribute are required to have a \"binary\" transfer option as described in RFC 4522. Attributes with this syntax will generally be referenced with names including \";binary\" (e.g., \"userCertificate;binary\").",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a AttributeTypeDescriptionAttributeSyntaxResponse object into the model struct
func readAttributeTypeDescriptionAttributeSyntaxResponseDataSource(ctx context.Context, r *client.AttributeTypeDescriptionAttributeSyntaxResponse, state *attributeSyntaxDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("attribute-type-description")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.StripSyntaxMinUpperBound = internaltypes.BoolTypeOrNil(r.StripSyntaxMinUpperBound)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
}

// Read a DirectoryStringAttributeSyntaxResponse object into the model struct
func readDirectoryStringAttributeSyntaxResponseDataSource(ctx context.Context, r *client.DirectoryStringAttributeSyntaxResponse, state *attributeSyntaxDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("directory-string")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllowZeroLengthValues = internaltypes.BoolTypeOrNil(r.AllowZeroLengthValues)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
}

// Read a TelephoneNumberAttributeSyntaxResponse object into the model struct
func readTelephoneNumberAttributeSyntaxResponseDataSource(ctx context.Context, r *client.TelephoneNumberAttributeSyntaxResponse, state *attributeSyntaxDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("telephone-number")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.StrictFormat = internaltypes.BoolTypeOrNil(r.StrictFormat)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
}

// Read a DistinguishedNameAttributeSyntaxResponse object into the model struct
func readDistinguishedNameAttributeSyntaxResponseDataSource(ctx context.Context, r *client.DistinguishedNameAttributeSyntaxResponse, state *attributeSyntaxDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("distinguished-name")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
}

// Read a GeneralizedTimeAttributeSyntaxResponse object into the model struct
func readGeneralizedTimeAttributeSyntaxResponseDataSource(ctx context.Context, r *client.GeneralizedTimeAttributeSyntaxResponse, state *attributeSyntaxDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generalized-time")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
}

// Read a IntegerAttributeSyntaxResponse object into the model struct
func readIntegerAttributeSyntaxResponseDataSource(ctx context.Context, r *client.IntegerAttributeSyntaxResponse, state *attributeSyntaxDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("integer")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
}

// Read a UuidAttributeSyntaxResponse object into the model struct
func readUuidAttributeSyntaxResponseDataSource(ctx context.Context, r *client.UuidAttributeSyntaxResponse, state *attributeSyntaxDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("uuid")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
}

// Read a GenericAttributeSyntaxResponse object into the model struct
func readGenericAttributeSyntaxResponseDataSource(ctx context.Context, r *client.GenericAttributeSyntaxResponse, state *attributeSyntaxDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generic")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
}

// Read a JsonObjectAttributeSyntaxResponse object into the model struct
func readJsonObjectAttributeSyntaxResponseDataSource(ctx context.Context, r *client.JsonObjectAttributeSyntaxResponse, state *attributeSyntaxDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("json-object")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
}

// Read a UserPasswordAttributeSyntaxResponse object into the model struct
func readUserPasswordAttributeSyntaxResponseDataSource(ctx context.Context, r *client.UserPasswordAttributeSyntaxResponse, state *attributeSyntaxDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("user-password")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
}

// Read a BooleanAttributeSyntaxResponse object into the model struct
func readBooleanAttributeSyntaxResponseDataSource(ctx context.Context, r *client.BooleanAttributeSyntaxResponse, state *attributeSyntaxDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("boolean")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
}

// Read a HexStringAttributeSyntaxResponse object into the model struct
func readHexStringAttributeSyntaxResponseDataSource(ctx context.Context, r *client.HexStringAttributeSyntaxResponse, state *attributeSyntaxDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("hex-string")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
}

// Read a BitStringAttributeSyntaxResponse object into the model struct
func readBitStringAttributeSyntaxResponseDataSource(ctx context.Context, r *client.BitStringAttributeSyntaxResponse, state *attributeSyntaxDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("bit-string")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
}

// Read a LdapUrlAttributeSyntaxResponse object into the model struct
func readLdapUrlAttributeSyntaxResponseDataSource(ctx context.Context, r *client.LdapUrlAttributeSyntaxResponse, state *attributeSyntaxDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap-url")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.StrictFormat = internaltypes.BoolTypeOrNil(r.StrictFormat)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
}

// Read a NameAndOptionalUidAttributeSyntaxResponse object into the model struct
func readNameAndOptionalUidAttributeSyntaxResponseDataSource(ctx context.Context, r *client.NameAndOptionalUidAttributeSyntaxResponse, state *attributeSyntaxDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("name-and-optional-uid")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
}

// Read resource information
func (r *attributeSyntaxDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state attributeSyntaxDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AttributeSyntaxApi.GetAttributeSyntax(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Attribute Syntax", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.AttributeTypeDescriptionAttributeSyntaxResponse != nil {
		readAttributeTypeDescriptionAttributeSyntaxResponseDataSource(ctx, readResponse.AttributeTypeDescriptionAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DirectoryStringAttributeSyntaxResponse != nil {
		readDirectoryStringAttributeSyntaxResponseDataSource(ctx, readResponse.DirectoryStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.TelephoneNumberAttributeSyntaxResponse != nil {
		readTelephoneNumberAttributeSyntaxResponseDataSource(ctx, readResponse.TelephoneNumberAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DistinguishedNameAttributeSyntaxResponse != nil {
		readDistinguishedNameAttributeSyntaxResponseDataSource(ctx, readResponse.DistinguishedNameAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GeneralizedTimeAttributeSyntaxResponse != nil {
		readGeneralizedTimeAttributeSyntaxResponseDataSource(ctx, readResponse.GeneralizedTimeAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.IntegerAttributeSyntaxResponse != nil {
		readIntegerAttributeSyntaxResponseDataSource(ctx, readResponse.IntegerAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.UuidAttributeSyntaxResponse != nil {
		readUuidAttributeSyntaxResponseDataSource(ctx, readResponse.UuidAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GenericAttributeSyntaxResponse != nil {
		readGenericAttributeSyntaxResponseDataSource(ctx, readResponse.GenericAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.JsonObjectAttributeSyntaxResponse != nil {
		readJsonObjectAttributeSyntaxResponseDataSource(ctx, readResponse.JsonObjectAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.UserPasswordAttributeSyntaxResponse != nil {
		readUserPasswordAttributeSyntaxResponseDataSource(ctx, readResponse.UserPasswordAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.BooleanAttributeSyntaxResponse != nil {
		readBooleanAttributeSyntaxResponseDataSource(ctx, readResponse.BooleanAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.HexStringAttributeSyntaxResponse != nil {
		readHexStringAttributeSyntaxResponseDataSource(ctx, readResponse.HexStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.BitStringAttributeSyntaxResponse != nil {
		readBitStringAttributeSyntaxResponseDataSource(ctx, readResponse.BitStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LdapUrlAttributeSyntaxResponse != nil {
		readLdapUrlAttributeSyntaxResponseDataSource(ctx, readResponse.LdapUrlAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.NameAndOptionalUidAttributeSyntaxResponse != nil {
		readNameAndOptionalUidAttributeSyntaxResponseDataSource(ctx, readResponse.NameAndOptionalUidAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
