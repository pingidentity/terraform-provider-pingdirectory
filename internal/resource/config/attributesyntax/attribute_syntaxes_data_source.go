package attributesyntax

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &attributeSyntaxesDataSource{}
	_ datasource.DataSourceWithConfigure = &attributeSyntaxesDataSource{}
)

// Create a Attribute Syntaxes data source
func NewAttributeSyntaxesDataSource() datasource.DataSource {
	return &attributeSyntaxesDataSource{}
}

// attributeSyntaxesDataSource is the datasource implementation.
type attributeSyntaxesDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *attributeSyntaxesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_attribute_syntaxes"
}

// Configure adds the provider configured client to the data source.
func (r *attributeSyntaxesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type attributeSyntaxesDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *attributeSyntaxesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Attribute Syntax objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Attribute Syntax objects found in the configuration",
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
func (r *attributeSyntaxesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state attributeSyntaxesDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.AttributeSyntaxAPI.ListAttributeSyntaxes(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.AttributeSyntaxAPI.ListAttributeSyntaxesExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Attribute Syntax objects", err, httpResp)
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
		if response.AttributeTypeDescriptionAttributeSyntaxResponse != nil {
			attributes["id"] = types.StringValue(response.AttributeTypeDescriptionAttributeSyntaxResponse.Id)
			attributes["type"] = types.StringValue("attribute-type-description")
		}
		if response.DirectoryStringAttributeSyntaxResponse != nil {
			attributes["id"] = types.StringValue(response.DirectoryStringAttributeSyntaxResponse.Id)
			attributes["type"] = types.StringValue("directory-string")
		}
		if response.TelephoneNumberAttributeSyntaxResponse != nil {
			attributes["id"] = types.StringValue(response.TelephoneNumberAttributeSyntaxResponse.Id)
			attributes["type"] = types.StringValue("telephone-number")
		}
		if response.DistinguishedNameAttributeSyntaxResponse != nil {
			attributes["id"] = types.StringValue(response.DistinguishedNameAttributeSyntaxResponse.Id)
			attributes["type"] = types.StringValue("distinguished-name")
		}
		if response.GeneralizedTimeAttributeSyntaxResponse != nil {
			attributes["id"] = types.StringValue(response.GeneralizedTimeAttributeSyntaxResponse.Id)
			attributes["type"] = types.StringValue("generalized-time")
		}
		if response.IntegerAttributeSyntaxResponse != nil {
			attributes["id"] = types.StringValue(response.IntegerAttributeSyntaxResponse.Id)
			attributes["type"] = types.StringValue("integer")
		}
		if response.UuidAttributeSyntaxResponse != nil {
			attributes["id"] = types.StringValue(response.UuidAttributeSyntaxResponse.Id)
			attributes["type"] = types.StringValue("uuid")
		}
		if response.GenericAttributeSyntaxResponse != nil {
			attributes["id"] = types.StringValue(response.GenericAttributeSyntaxResponse.Id)
			attributes["type"] = types.StringValue("generic")
		}
		if response.JsonObjectAttributeSyntaxResponse != nil {
			attributes["id"] = types.StringValue(response.JsonObjectAttributeSyntaxResponse.Id)
			attributes["type"] = types.StringValue("json-object")
		}
		if response.UserPasswordAttributeSyntaxResponse != nil {
			attributes["id"] = types.StringValue(response.UserPasswordAttributeSyntaxResponse.Id)
			attributes["type"] = types.StringValue("user-password")
		}
		if response.BooleanAttributeSyntaxResponse != nil {
			attributes["id"] = types.StringValue(response.BooleanAttributeSyntaxResponse.Id)
			attributes["type"] = types.StringValue("boolean")
		}
		if response.HexStringAttributeSyntaxResponse != nil {
			attributes["id"] = types.StringValue(response.HexStringAttributeSyntaxResponse.Id)
			attributes["type"] = types.StringValue("hex-string")
		}
		if response.BitStringAttributeSyntaxResponse != nil {
			attributes["id"] = types.StringValue(response.BitStringAttributeSyntaxResponse.Id)
			attributes["type"] = types.StringValue("bit-string")
		}
		if response.LdapUrlAttributeSyntaxResponse != nil {
			attributes["id"] = types.StringValue(response.LdapUrlAttributeSyntaxResponse.Id)
			attributes["type"] = types.StringValue("ldap-url")
		}
		if response.NameAndOptionalUidAttributeSyntaxResponse != nil {
			attributes["id"] = types.StringValue(response.NameAndOptionalUidAttributeSyntaxResponse.Id)
			attributes["type"] = types.StringValue("name-and-optional-uid")
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
